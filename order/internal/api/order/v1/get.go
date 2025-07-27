package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Akbar-cmd/raketa-factory/order/internal/converter"
	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	orderV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, err := a.service.GetOrderByUuid(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "order by this UUID %s not found", params.OrderUUID.String())
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, status.Errorf(codes.Unavailable, "Order service timeout")
		}
		if errors.Is(err, model.ErrOrderInternalError) {
			return nil, status.Errorf(codes.Internal, "Order service internal error")
		}
		return nil, err
	}

	return converter.OrderDataToAPI(order), nil
}
