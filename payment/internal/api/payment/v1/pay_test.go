package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrder() {
	type args struct {
		req *paymentV1.PayOrderRequest
	}

	var (
		trxnUUID      = uuid.NewString()
		orderUuid     = uuid.NewString()
		userUuid      = uuid.NewString()
		paymentMethod = gofakeit.IntRange(0, 4)

		internalServerErr  = status.Error(codes.Internal, "internal error")
		invalidArgumentErr = status.Error(codes.InvalidArgument, "order_uuid and user_uuid must be set")

		validReq = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentV1.PaymentMethod(paymentMethod),
		}

		invalidOrderReq = &paymentV1.PayOrderRequest{
			OrderUuid:     "",
			UserUuid:      userUuid,
			PaymentMethod: paymentV1.PaymentMethod(paymentMethod),
		}

		invalidUserReq = &paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      "",
			PaymentMethod: paymentV1.PaymentMethod(paymentMethod),
		}

		res = &paymentV1.PayOrderResponse{
			TransactionUuid: trxnUUID,
		}
	)

	tests := []struct {
		name                        string
		args                        args
		want                        *paymentV1.PayOrderResponse
		err                         error
		paymentServiceMockConfigure func()
	}{
		{
			name: "Success Case",
			args: args{req: validReq},
			want: res,
			err:  nil,
			paymentServiceMockConfigure: func() {
				s.paymentService.On(
					"PayOrder",
					s.ctx,
					orderUuid,
					userUuid,
					paymentV1.PaymentMethod(paymentMethod).String(),
				).Return(trxnUUID, nil).Once()
			},
		},
		{
			name: "Empty order_uuid",
			args: args{req: invalidOrderReq},
			want: nil,
			err:  invalidArgumentErr,
			paymentServiceMockConfigure: func() {
				s.paymentService.On(
					"PayOrder",
					s.ctx,
					"",
					userUuid,
					paymentV1.PaymentMethod(paymentMethod).String(),
				).Return("", invalidArgumentErr).Once()
			},
		},
		{
			name: "Empty user_uuid",
			args: args{req: invalidUserReq},
			want: nil,
			err:  invalidArgumentErr,
			paymentServiceMockConfigure: func() {
				s.paymentService.On(
					"PayOrder",
					s.ctx,
					orderUuid,
					"",
					paymentV1.PaymentMethod(paymentMethod).String(),
				).Return("", invalidArgumentErr).Once()
			},
		},
		{
			name: "Internal Error",
			args: args{req: validReq},
			want: nil,
			err:  internalServerErr,
			paymentServiceMockConfigure: func() {
				s.paymentService.On(
					"PayOrder",
					s.ctx,
					orderUuid,
					userUuid,
					paymentV1.PaymentMethod(paymentMethod).String(),
				).Return("", internalServerErr).Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.paymentServiceMockConfigure()
			res, err := s.api.PayOrder(s.ctx, tt.args.req)
			s.Require().Equal(tt.want, res)
			s.Require().Equal(tt.err, err)
		})
	}
}
