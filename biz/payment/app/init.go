package app

import (
	paymentService "github.com/yosuarichel/billing-engine/biz/payment/service"
)

type PaymentApp struct {
	PaymentService *paymentService.PaymentService
}

func NewPaymentApp(
	paymentService *paymentService.PaymentService,
) *PaymentApp {
	return &PaymentApp{
		PaymentService: paymentService,
	}
}
