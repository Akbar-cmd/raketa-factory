package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/converter"
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filters := converter.PartsFilterToModel(req.GetFilter())

	parts, err := a.service.ListParts(ctx, filters)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return nil, status.Error(codes.NotFound, "parts not found")
		}
		if errors.Is(err, model.ErrPartsInternalError) {
			return nil, status.Error(codes.Internal, "Internal inventory service error")
		}
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
