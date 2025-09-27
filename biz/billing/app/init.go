package app

import (
	"github.com/yosuarichel/billing-engine/biz/billing/service"
)

type BillingApp struct {
	BillingLoanService    *service.LoanService
	BillingService        *service.BillingService
	BillingPaymentService *service.PaymentService
}

func NewBillingApp(
	billingLoanService *service.LoanService,
	billingService *service.BillingService,
	billingPaymentService *service.PaymentService,
) *BillingApp {
	return &BillingApp{
		BillingLoanService:    billingLoanService,
		BillingService:        billingService,
		BillingPaymentService: billingPaymentService,
	}
}
