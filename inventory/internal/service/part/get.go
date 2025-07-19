package part

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
)

// GetPart получает деталь из хранилища по uuid и возвращает ее
func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.inventoryRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
