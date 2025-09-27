package app

import (
	"github.com/yosuarichel/billing-engine/biz/loan/service"
)

type LoanApp struct {
	LoanService *service.LoanService
}

func NewLoanApp(
	loanService *service.LoanService,
) *LoanApp {
	return &LoanApp{
		LoanService: loanService,
	}
}
