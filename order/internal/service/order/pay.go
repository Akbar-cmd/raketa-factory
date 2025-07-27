package order

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, paymentMethod string) (transactionUUID string, err error) {
	order, err := s.orderRepository.GetOrderByUuid(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	// проверяем статус заказа
	if resp, ok := canPayOrder(order); ok {
		return "", resp
	}

	// Создаем таймаут на обращение
	clientCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// оплачиваем посредством PaymenTservice
	trxnUUID, err := s.paymentClient.PayOrder(clientCtx, order.UserUUID, orderUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	// обновление заказа
	status := model.OrderStatusPaid
	updateErr := s.orderRepository.UpdateOrder(ctx, order.OrderUUID, model.OrderUpdateInfo{
		Status:          &status,
		PaymentMethod:   lo.ToPtr(model.PaymentMethod(paymentMethod)),
		TransactionUUID: &trxnUUID,
	})

	if updateErr != nil {
		return "", err
	}

	return trxnUUID, nil
}

func canPayOrder(order model.OrderData) (error, bool) {
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrPaymentConflict, true
	case model.OrderStatusCancelled:
		return model.ErrPaymentConflict, true
	case model.OrderStatusPendingPayment:
		return nil, false
	default:
		return model.ErrPaymentInternalError, true
	}
}
