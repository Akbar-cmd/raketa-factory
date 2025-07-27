package payment

import "github.com/brianvoe/gofakeit/v7"

func (s *ServiceSuite) TestPayOrder() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = gofakeit.RandomString([]string{"CARD", "SBP", "CREDIT_CARD"})
	)
	trxnUUID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, paymentMethod)
	s.NoError(err)
	s.NotEmpty(trxnUUID)
}
