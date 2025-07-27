package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/Akbar-cmd/raketa-factory/order/internal/model"
	repoModel "github.com/Akbar-cmd/raketa-factory/order/internal/repository/model"
)

type RepositorySuite struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	repository *repository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()

	s.repository = NewRepository()

	for i := 0; i < 2; i++ {
		uuid := gofakeit.UUID()
		s.repository.data[uuid] = repoModel.OrderData{
			OrderUUID:  uuid,
			UserUUID:   gofakeit.UUID(),
			PartUuids:  []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice: gofakeit.Float64(),
			Status:     repoModel.OrderStatus(model.OrderStatusPendingPayment),
		}
	}
}

func (s *RepositorySuite) TearDownTest() {}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
