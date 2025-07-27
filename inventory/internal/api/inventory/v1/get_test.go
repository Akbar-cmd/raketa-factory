package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/converter"
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPart() {
	type args struct {
		req *inventoryV1.GetPartRequest
	}

	var (
		partUuid = uuid.NewString()

		validReq = &inventoryV1.GetPartRequest{
			Uuid: partUuid,
		}
		emptyUuidReq = &inventoryV1.GetPartRequest{
			Uuid: "",
		}

		part = model.Part{
			Uuid:          partUuid,
			Name:          gofakeit.Name(),
			Description:   "Primary propulsion unit",
			Price:         gofakeit.Float64Range(100, 10_000),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      "ENGINE",
			Dimensions: model.Dimensions{
				Width:  gofakeit.Float64Range(0.1, 10.0),
				Height: gofakeit.Float64Range(0.1, 10.0),
				Length: gofakeit.Float64Range(0.1, 10.0),
				Weight: gofakeit.Float64Range(0.1, 10.0),
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Name(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags: []string{gofakeit.EmojiTag(), gofakeit.EmojiTag()},
			Metadata: model.Metadata{
				StringValue: lo.ToPtr(gofakeit.Word()),
				Int64Value:  lo.ToPtr(gofakeit.Int64()),
				DoubleValue: lo.ToPtr(gofakeit.Float64()),
				BoolValue:   lo.ToPtr(gofakeit.Bool()),
			},
			CreatedAt: timestamppb.Now().AsTime(),
			UpdatedAt: lo.ToPtr(time.Now()),
		}

		ErrPartNotFound  = status.Error(codes.NotFound, "part UUID not found")
		ErrPartsInternal = status.Error(codes.Internal, "internal error while getting part UUID")

		res = &inventoryV1.GetPartResponse{
			Part: converter.PartToProto(part),
		}
	)

	tests := []struct {
		name                          string
		args                          args
		want                          *inventoryV1.GetPartResponse
		err                           error
		inventoryServiceMockConfigure func()
	}{
		{
			name: "Success Case",
			args: args{req: validReq},
			want: res,
			err:  nil,
			inventoryServiceMockConfigure: func() {
				s.inventoryService.On(
					"GetPart",
					s.ctx,
					partUuid,
				).Return(part, nil).Once()
			},
		},
		{
			name: "Part Not Found",
			args: args{req: validReq},
			want: nil,
			err:  ErrPartNotFound,
			inventoryServiceMockConfigure: func() {
				s.inventoryService.On(
					"GetPart",
					s.ctx,
					partUuid,
				).Return(part, ErrPartNotFound).Once()
			},
		},
		{
			name: "Empty UUID - Part Not Found",
			args: args{emptyUuidReq},
			want: nil,
			err:  ErrPartNotFound,
			inventoryServiceMockConfigure: func() {
				s.inventoryService.On(
					"GetPart",
					s.ctx,
					"",
				).Return(part, ErrPartNotFound).Once()
			},
		},
		{
			name: "Internal Error",
			args: args{validReq},
			want: nil,
			err:  ErrPartsInternal,
			inventoryServiceMockConfigure: func() {
				s.inventoryService.On(
					"GetPart",
					s.ctx,
					partUuid,
				).Return(part, ErrPartsInternal).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.inventoryServiceMockConfigure()
			res, err := s.api.GetPart(s.ctx, tt.args.req)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
