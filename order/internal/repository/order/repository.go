package order

import (
	"sync"

	def "github.com/Akbar-cmd/raketa-factory/order/internal/repository"
	repoModel "github.com/Akbar-cmd/raketa-factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.OrderData
}

func NewRepository() *repository {
	repo := &repository{
		data: make(map[string]repoModel.OrderData),
	}
	return repo
}
