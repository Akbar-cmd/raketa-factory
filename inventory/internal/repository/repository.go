package repository

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
