package order

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (info model.OrderCreationInfo, err error) {
	// Создаем таймаут на обращение
	clientCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// получаем список запчастей по uuid
	filter := model.PartsFilter{
		Uuids: partsUUIDs,
	}
	partsList, err := s.inventoryClient.ListParts(clientCtx, filter)
	if err != nil {
		return model.OrderCreationInfo{}, err
	}
	if len(partsList) != len(partsUUIDs) {
		return model.OrderCreationInfo{}, model.ErrOrderConflict
	}

	// Создаем базовую информацию о заказе
	order := model.OrderData{
		OrderUUID: uuid.NewString(),
		UserUUID:  userUUID,
		Status:    model.OrderStatusPendingPayment,
	}

	// Проверяем на наличие всех необходимых запчастей, при нахождении добавляем в заказ и плюсуем цену,
	// при не находе падаем в ошибку
	for _, partUuid := range partsUUIDs {
		var part *model.Part
		for _, p := range partsList {
			if p.Uuid == partUuid && p.StockQuantity > 0 {
				part = &p
			}
		}

		if part == nil {
			return model.OrderCreationInfo{
				OrderUUID:  "",
				TotalPrice: 0,
			}, fmt.Errorf("%s - %w", partUuid, model.ErrPartNotFound)
		}

		order.PartUuids = append(order.PartUuids, partUuid)
		order.TotalPrice += part.Price
	}

	// Сохраняем заказ
	orderInfo, createrOrderErr := s.orderRepository.CreateOrder(ctx, order)
	if createrOrderErr != nil {
		return model.OrderCreationInfo{}, createrOrderErr
	}

	log.Printf(`
💳 [Order Created]
• 🆔 Order UUID: %s
• 👤 User UUID: %s
• 💰 Part UUIDs: %v
• 💰 Total Price: %f
• 💰 Status: %s
`, order.OrderUUID, order.UserUUID, order.PartUuids, order.TotalPrice, order.Status,
	)

	return model.OrderCreationInfo{
		OrderUUID:  orderInfo.OrderUUID,
		TotalPrice: orderInfo.TotalPrice,
	}, nil
}
