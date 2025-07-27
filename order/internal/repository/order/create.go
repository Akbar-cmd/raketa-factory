package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	"github.com/Akbar-cmd/raketa-factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(_ context.Context, order model.OrderData) (model.OrderCreationInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.OrderUUID == "" {
		order.OrderUUID = uuid.NewString()
	}

	r.data[order.OrderUUID] = converter.OrderDataToRepoModel(order)

	return model.OrderCreationInfo{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
