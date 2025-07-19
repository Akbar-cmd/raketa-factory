package part

import (
	repo "github.com/Akbar-cmd/raketa-factory/inventory/internal/repository"
	def "github.com/Akbar-cmd/raketa-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repo.InventoryRepository
}

func NewService(inventoryRepository repo.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
