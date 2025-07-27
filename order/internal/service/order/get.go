package order

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *service) GetOrderByUuid(ctx context.Context, uuid string) (model.OrderData, error) {
	order, err := s.orderRepository.GetOrderByUuid(ctx, uuid)
	if err != nil {
		return model.OrderData{}, err
	}
	return order, nil
}
