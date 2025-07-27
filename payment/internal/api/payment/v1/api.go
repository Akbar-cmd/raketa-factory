package v1

import (
	"github.com/Akbar-cmd/raketa-factory/payment/internal/service"
	paymentV1 "github.com/Akbar-cmd/raketa-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	service service.PaymentService
}

func NewAPI(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
