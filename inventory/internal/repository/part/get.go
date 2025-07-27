package part

import (
	"context"

	"github.com/Akbar-cmd/raketa-factory/inventory/internal/model"
	repoConverter "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/converter"
)

// GetPart возвращает информацию по детали по ее UUID
func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}
	return repoConverter.PartToModel(repoPart), nil
}
