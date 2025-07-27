package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	"github.com/Akbar-cmd/raketa-factory/order/internal/repository/converter"
)

func (s *RepositorySuite) TestCreateOrder() {
	type args struct {
		part model.OrderData
	}

	part := model.OrderData{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: gofakeit.Float64(),
		Status:     model.OrderStatusPendingPayment,
	}

	partWithoutID := part
	partWithoutID.OrderUUID = ""
	partWithoutID.TotalPrice = gofakeit.Float64Range(200, 2_000)

	tests := []struct {
		name  string
		args  args
		err   error
		newID bool
	}{
		{
			name: "Case with UUID",
			args: args{
				part: part,
			},
			err:   nil,
			newID: false,
		},
		{
			name: "Case without UUID",
			args: args{
				part: partWithoutID,
			},
			err:   nil,
			newID: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.repository.CreateOrder(s.ctx, tt.args.part)
			s.Require().Equal(tt.err, err)
			s.Require().NotEmpty(res.OrderUUID, "возвращённый UUID не должен быть пустым")
			s.Equal(tt.args.part.TotalPrice, res.TotalPrice)

			stored, ok := s.repository.data[res.OrderUUID]
			s.Require().True(ok, "заказ должен быть сохранён в репозитории")
			s.Equal(converter.OrderDataToRepoModel(tt.args.part).TotalPrice,
				stored.TotalPrice)
		})
	}
}
