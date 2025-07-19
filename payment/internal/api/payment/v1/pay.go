package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Akbar-cmd/raketa-factory/payment/internal/model"
	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	trxnUUID, err := a.service.PayOrder(ctx, req.GetOrderUuid(), req.GetUserUuid(), req.GetPaymentMethod().String())
	if err != nil {
		if errors.Is(err, model.ErrInvalidArgument) {
			return nil, status.Error(codes.InvalidArgument, "order_uuid and user_uuid must be set")
		}
		if errors.Is(err, model.ErrPaymentInternalServer) {
			return nil, status.Error(codes.Internal, "Internal payment service error")
		}
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: trxnUUID,
	}, nil
}
