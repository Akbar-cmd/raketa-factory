package converter

import (
	"log"

	"github.com/google/uuid"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	orderV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/openapi/order/v1"
)

func OrderDataToAPI(order model.OrderData) *orderV1.OrderDto {
	var paymentMethod orderV1.PaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.PaymentMethod(*order.PaymentMethod)
	} else {
		paymentMethod = orderV1.PaymentMethodUNKNOWN
	}

	return &orderV1.OrderDto{
		OrderUUID:       StringToUUID(order.OrderUUID),
		UserUUID:        StringToUUID(order.UserUUID),
		PartUuids:       StringToUUIDs(order.PartUuids),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: ToNilUUID(order.TransactionUUID),
		PaymentMethod:   paymentMethod,
		OrderStatus:     orderV1.OrderStatus(order.Status),
	}
}

func StringToUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
	}

	return u
}

func StringToUUIDs(arr []string) []uuid.UUID {
	uuids := make([]uuid.UUID, len(arr))
	for i, s := range arr {
		uuids[i] = StringToUUID(s)
	}
	return uuids
}

func ToNilUUID(s *string) orderV1.NilUUID {
	if s == nil {
		return orderV1.NilUUID{}
	}
	return orderV1.NilUUID{Value: StringToUUID(*s)}
}

func UUIDsToStrings(arr []uuid.UUID) []string {
	uuids := make([]string, len(arr))
	for i, s := range arr {
		uuids[i] = s.String()
	}
	return uuids
}
