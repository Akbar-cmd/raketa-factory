package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Akbar-cmd/raketa-factory/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, orderUUID, userUUID, paymentMethod string) (transactionUUID string, err error) {
	if orderUUID == "" || userUUID == "" {
		return "", model.ErrInvalidArgument
	}

	log.Printf(`
💳 [Order Paid]
• 🆔 Order UUID: %s
• 👤 User UUID: %s
• 💰 Payment Method: %s
`, orderUUID, userUUID, paymentMethod,
	)

	trxnUUID := uuid.NewString()

	log.Printf("✅ Оплата прошла успешно, transaction_uuid: %v\n", trxnUUID)

	return trxnUUID, nil
}
