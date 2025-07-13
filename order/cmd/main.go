package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort         = "8080"
	inventoryAddress = "localhost:50051"
	paymentAddress   = "localhost:50052"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// OrderStorage представляет хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderStorage создает новое хранилище данных о заказах
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// UpdateOrder сохраняет заказ в хранилище
func (s *OrderStorage) UpdateOrder(uuid string, order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = order
}

// GetOrder возвращает информацию о заказе по его uuid
func (s *OrderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}
	return order
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказов
type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API заказов
func NewOrderHandler(storage *OrderStorage, inventoryClient inventoryV1.InventoryServiceClient, paymentClient paymentV1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// CreateOrder обрабатывает создание нового заказа
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	// Валидируем входные данные

	// Проверка UserUUID
	if req.UserUUID == uuid.Nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "user_uuid must be a valid non-empty UUID",
		}, nil
	}

	// Проверка наличия хотя бы одного элемента в PartUuids
	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "part_uuids must contain at least one UUID",
		}, nil
	}

	// Проверка каждого uuid в списке PartUuids
	partUuidsStr := make([]string, len(req.PartUuids))
	for i, pu := range req.PartUuids {
		if pu == uuid.Nil {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: fmt.Sprintf("part_uuids[%d] is not a valid UUID", i),
			}, nil
		}
		partUuidsStr[i] = pu.String()
	}

	// Запрос деталей из InventoryService
	listPartsResp, err := h.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: partUuidsStr,
		},
	})
	if err != nil {
		return &orderV1.BadGatewayError{
			Code:    502,
			Message: "Failed to fetch parts from inventory service",
		}, nil
	}

	// Проверяем, что все запрошенные детали найдены
	foundParts := make(map[string]bool)
	for _, part := range listPartsResp.Parts {
		foundParts[part.Uuid] = true
	}

	for _, partUuidStr := range partUuidsStr {
		if !foundParts[partUuidStr] {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: fmt.Sprintf("Part withUUID %s not found", partUuidStr),
			}, nil
		}
	}

	// Подсчет общей стоимссти
	var totalPrice float64
	for _, part := range listPartsResp.Parts {
		totalPrice += part.Price
	}

	// Генерация UUID для заказа
	orderUUID := uuid.New()

	// Создание заказа
	order := &orderV1.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: orderV1.NilUUID{},
		PaymentMethod:   orderV1.PaymentMethod(""),
		OrderStatus:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	// Сохранение заказа в хранилище
	h.storage.UpdateOrder(orderUUID.String(), order)

	log.Printf("✅ Создан заказ %s для пользователя %s на сумму %.2f",
		orderUUID.String(), req.UserUUID.String(), totalPrice)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder обрабатывает оплату заказа
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID: `" + params.OrderUUID.String() + "` not found",
		}, nil
	}

	// Проверяем статус заказа
	if order.OrderStatus != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("Order with UUID %s is already processed (status: %s)", params.OrderUUID.String(), order.OrderStatus),
		}, nil
	}

	// Конвертация PaymentMethod для gRPC-запроса
	var grpcPaymentMethod paymentV1.PaymentMethod
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_CARD
	case orderV1.PayOrderRequestPaymentMethodSBP:
		grpcPaymentMethod = paymentV1.PaymentMethod_SBP
	case orderV1.PayOrderRequestPaymentMethodCREDITCARD:
		grpcPaymentMethod = paymentV1.PaymentMethod_CREDIT_CARD
	case orderV1.PayOrderRequestPaymentMethodINVESTORMONEY:
		grpcPaymentMethod = paymentV1.PaymentMethod_INVESTOR_MONEY
	default:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid payment method",
		}, nil
	}

	// Вызов PaymentService
	paymentResp, err := h.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: grpcPaymentMethod,
	})
	if err != nil {
		return &orderV1.BadGatewayError{
			Code:    502,
			Message: "Failed to process payment",
		}, nil
	}

	// Парсинг transaction_uuid
	transactionUUID, err := uuid.Parse(paymentResp.TransactionUuid)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Invalid transaction UUID from payment service",
		}, nil
	}

	// Обновление заказа
	order.OrderStatus = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewNilUUID(transactionUUID)
	order.PaymentMethod = orderV1.PaymentMethod(req.PaymentMethod)

	h.storage.UpdateOrder(params.OrderUUID.String(), order)

	log.Printf("✅ Заказ %s успешно оплачен, transaction_uuid: %s",
		params.OrderUUID.String(), transactionUUID.String())

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// GetOrderByUUID возвращает информацию о заказе по его UUID
func (h *OrderHandler) GetOrderByUuid(_ context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("Order with UUID %s not found", params.OrderUUID.String()),
		}, nil
	}

	return order, nil
}

// CancelOrder обрабатывает отмену заказа
func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order with UUID:`" + params.OrderUUID.String() + "` not found",
		}, nil
	}

	// Проверяем статус заказа
	switch order.OrderStatus {
	case orderV1.OrderStatusPAID:
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order is already paid and cannot be cancelled",
		}, nil
	case orderV1.OrderStatusCANCELLED:
		return &orderV1.ConflictError{
			Code:    409,
			Message: "Order is already cancelled",
		}, nil
	case orderV1.OrderStatusPENDINGPAYMENT:
		// Разрешаем отмену
	default:
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Unknown order status",
		}, nil
	}

	// Обновляем статус заказа
	order.OrderStatus = orderV1.OrderStatusCANCELLED
	h.storage.UpdateOrder(params.OrderUUID.String(), order)

	log.Printf("✅ Заказ %s успешно отменён", params.OrderUUID.String())

	return &orderV1.CancelOrderNoContent{}, nil
}

// NewError создает новую ошибку в формате GenericError
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func connectToInventoryService() (inventoryV1.InventoryServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	client := inventoryV1.NewInventoryServiceClient(conn)
	return client, conn, nil
}

func connectToPaymentService() (paymentV1.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	client := paymentV1.NewPaymentServiceClient(conn)
	return client, conn, nil
}

func main() {
	// Подключение к gRPC-сервисам
	inventoryClient, inventoryConn, err := connectToInventoryService()
	if err != nil {
		log.Printf("❌ Ошибка подключения к InventoryService: %v", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с InventoryService: %v", err)
		}
	}()

	paymentClient, paymentConn, err := connectToPaymentService()
	if err != nil {
		log.Printf("❌ Ошибка подключения к PaymentService: %v", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с PaymentService: %v", err)
		}
	}()

	// Создаем хранилище данных о заказах
	storage := NewOrderStorage()

	// Создаем обработчик API заказов
	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	// Созадем OpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавялем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
