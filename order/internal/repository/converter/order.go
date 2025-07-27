package converter

import (
	"github.com/samber/lo"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	repoModel "github.com/Akbar-cmd/raketa-factory/order/internal/repository/model"
)

func OrderDataToModel(order repoModel.OrderData) model.OrderData {
	paymentMethod := lo.ToPtr(model.PaymentMethod(lo.FromPtr(order.PaymentMethod)))

	return model.OrderData{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
	}
}

func OrderUpdateInfoToRepoModel(order model.OrderUpdateInfo) repoModel.OrderUpdateInfo {
	paymentMethod := lo.ToPtr(repoModel.PaymentMethod(lo.FromPtr(order.PaymentMethod)))
	status := lo.ToPtr(repoModel.OrderStatus(lo.FromPtr(order.Status)))

	return repoModel.OrderUpdateInfo{
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}
}

func OrderDataToRepoModel(order model.OrderData) repoModel.OrderData {
	paymentMethod := lo.ToPtr(repoModel.PaymentMethod(lo.FromPtr(order.PaymentMethod)))

	return repoModel.OrderData{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          repoModel.OrderStatus(order.Status),
	}
}
