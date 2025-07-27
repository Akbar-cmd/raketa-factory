package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByUuid() {
	type args struct {
		uuid string
	}

	var (
		uuid = gofakeit.UUID()

		order = model.OrderData{
			OrderUUID:  uuid,
			UserUUID:   gofakeit.UUID(),
			PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice: gofakeit.Float64(),
			Status:     model.OrderStatusCancelled,
		}
	)

	tests := []struct {
		name                         string
		args                         args
		want                         model.OrderData
		err                          error
		orderRepositoryMockConfigure func()
	}{
		{
			name: "Success Case",
			args: args{
				uuid: uuid,
			},
			want: order,
			err:  nil,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("GetOrderByUuid", s.ctx, uuid).Return(order, nil).Once()
			},
		},
		{
			name: "Not Found Error",
			args: args{
				uuid: uuid,
			},
			want: model.OrderData{},
			err:  model.ErrOrderNotFound,
			orderRepositoryMockConfigure: func() {
				s.orderRepository.On("GetOrderByUuid", s.ctx, uuid).Return(model.OrderData{}, model.ErrOrderNotFound).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.orderRepositoryMockConfigure()
			res, err := s.service.GetOrderByUuid(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
