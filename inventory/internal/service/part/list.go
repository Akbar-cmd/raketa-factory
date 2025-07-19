package part

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
)

// ListParts получает отфильтрованные данные из хранилища
func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.inventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
