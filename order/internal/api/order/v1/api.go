package v1

import "github.com/Akbar-cmd/raketa-factory/order/internal/service"

type api struct {
	service service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
