package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	"github.com/Akbar-cmd/raketa-factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(_ context.Context, uuid string, updateOrder model.OrderUpdateInfo) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// находим текущий заказ в памяти
	order, ok := r.data[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	// Конвертируем model.OrderUpdateInfo в repo
	repoUpd := converter.OrderUpdateInfoToRepoModel(updateOrder)

	// Обновляем поля, только если они были установлены в запросе
	if repoUpd.PartUuids != nil {
		order.PartUuids = *repoUpd.PartUuids
	}
	if repoUpd.TotalPrice != nil {
		order.TotalPrice = *repoUpd.TotalPrice
	}
	if repoUpd.TransactionUUID != nil {
		order.TransactionUUID = repoUpd.TransactionUUID
	}
	if repoUpd.PaymentMethod != nil {
		order.PaymentMethod = lo.ToPtr(lo.FromPtr(repoUpd.PaymentMethod))
	}
	if repoUpd.Status != nil {
		order.Status = *repoUpd.Status
	}

	// Сохраняем
	r.data[uuid] = order

	return nil
}
