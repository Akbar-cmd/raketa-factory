package service

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.OrderData) (info model.OrderCreationInfo, err error)
	GetOrderByUuid(ctx context.Context, uuid string) (model.OrderData, error)
	UpdateOrder(ctx context.Context, uuid string, updateOrder model.OrderUpdateInfo) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (info model.OrderCreationInfo, err error)
	GetOrderByUuid(ctx context.Context, uuid string) (model.OrderData, error)
	CancelOrder(ctx context.Context, uuid string) error
	PayOrder(ctx context.Context, uuid, paymentMethod string) (transactionUUID string, err error)
}
