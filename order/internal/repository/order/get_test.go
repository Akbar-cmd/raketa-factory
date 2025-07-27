package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	"github.com/Akbar-cmd/raketa-factory/order/internal/repository/converter"
)

func (s *RepositorySuite) TestGetOrderByUuid() {
	type args struct {
		uuid string
	}

	var (
		uuid         string
		expectedPart model.OrderData
	)

	k := gofakeit.Number(0, len(s.repository.data)-1)
	i := 0
	for _, v := range s.repository.data {
		if i == k {
			uuid = v.OrderUUID
			expectedPart = converter.OrderDataToModel(v)
			break
		}
		i++
	}

	tests := []struct {
		name string
		args args
		want model.OrderData
		err  error
	}{
		{
			name: "success case",
			args: args{
				uuid: uuid,
			},
			want: expectedPart,
			err:  nil,
		},
		{
			name: "error case",
			args: args{
				uuid: gofakeit.UUID(),
			},
			want: model.OrderData{},
			err:  model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.repository.GetOrderByUuid(s.ctx, tt.args.uuid)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
