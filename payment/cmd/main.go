package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

// paymentService реализует контракт PaymentService
type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// NewPaymentService конструктор сервиса оплаты
func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// PayOrder обрабатывает запрос оплаты заказа, генерирует transaction_uuid и логирует результат
func (s *PaymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Валидация полей запроса
	if req.GetOrderUuid() == "" || req.GetUserUuid() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "order_uuid and user_uuid must be set")
	}

	// Генерация UUID транзакции
	txnUUID := uuid.NewString()

	// Логирование
	log.Printf("Оплата прошла успешно, \n transaction_uuid: %s\n order_uuid: %s\n user_uuid: %s\n method: %s\n", txnUUID, req.GetOrderUuid(), req.GetUserUuid(), req.GetPaymentMethod().String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: txnUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создание и регистрация gRPC-сервера
	s := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(s, NewPaymentService())
	reflection.Register(s)

	// Запуск сервера в горутине
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
