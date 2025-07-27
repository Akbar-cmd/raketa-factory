package order

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.GetOrderByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	// Проверяем статус заказа
	switch order.Status {
	case model.OrderStatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return model.ErrOrderAlreadyCancelled
	}

	// Обновляем статус заказа
	status := model.OrderStatusCancelled
	err = s.orderRepository.UpdateOrder(ctx, order.OrderUUID, model.OrderUpdateInfo{
		Status: &status,
	})
	if err != nil {
		return err
	}

	return nil
}
