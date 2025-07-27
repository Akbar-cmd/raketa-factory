package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *RepositorySuite) TestUpdateOrder() {
	type args struct {
		uuid        string
		updateOrder model.OrderUpdateInfo
	}

	var (
		uuid            string
		partUuids       = []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}
		totalPrice      = gofakeit.Float64()
		transactionUuid = gofakeit.UUID()
		paymentMethod   = model.PaymentMethodCard
		status          = model.OrderStatusPaid

		orderUpdate = model.OrderUpdateInfo{
			PartUuids:       &partUuids,
			TotalPrice:      &totalPrice,
			TransactionUUID: &transactionUuid,
			PaymentMethod:   &paymentMethod,
			Status:          &status,
		}
	)

	k := gofakeit.Number(0, len(s.repository.data)-1)
	i := 0
	for _, v := range s.repository.data {
		if i == k {
			uuid = v.OrderUUID
			break
		}
		i++
	}

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Success Case",
			args: args{
				uuid:        uuid,
				updateOrder: orderUpdate,
			},
			err: nil,
		},
		{
			name: "Not Found Error",
			args: args{
				uuid:        gofakeit.UUID(),
				updateOrder: orderUpdate,
			},
			err: model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repository.UpdateOrder(s.ctx, tt.args.uuid, tt.args.updateOrder)
			s.Require().Equal(tt.err, err)
		})
	}
}
