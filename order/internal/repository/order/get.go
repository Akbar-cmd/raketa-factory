package order

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	"github.com/Akbar-cmd/raketa-factory/order/internal/repository/converter"
)

func (r *repository) GetOrderByUuid(_ context.Context, uuid string) (model.OrderData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.OrderData{}, model.ErrOrderNotFound
	}

	return converter.OrderDataToModel(order), nil
}
