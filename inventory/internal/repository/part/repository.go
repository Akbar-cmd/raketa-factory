package part

import (
	"log"
	"sync"

	def "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository"
	repoModel "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	repo := &repository{
		data: make(map[string]repoModel.Part),
	}
	repo.initParts()
	return repo
}

// initParts загружает сгенерированные данные в память репозитория.
func (r *repository) initParts() {
	initial := generateParts()

	for _, part := range initial {
		r.data[part.Uuid] = part
	}

	log.Printf("✅ Инициализировано %d запчастей в inventory", len(initial))
}
