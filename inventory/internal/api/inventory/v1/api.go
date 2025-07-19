package v1

import (
	"github.com/Akbar-cmd/raketa-factory/inventory/internal/service"
	inventoryV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	service service.InventoryService
}

func NewAPI(service service.InventoryService) *api {
	return &api{
		service: service,
	}
}
