package repository

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.OrderData) (info model.OrderCreationInfo, err error)
	GetOrderByUuid(ctx context.Context, uuid string) (model.OrderData, error)
	UpdateOrder(ctx context.Context, uuid string, updateOrder model.OrderUpdateInfo) error
}
