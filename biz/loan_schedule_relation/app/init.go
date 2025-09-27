package app

import (
	"github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/service"
)

type LoanScheduleRelationApp struct {
	LoanScheduleRelationService *service.LoanScheduleRelationService
}

func NewLoanScheduleRelationApp(
	loanScheduleRelationService *service.LoanScheduleRelationService,
) *LoanScheduleRelationApp {
	return &LoanScheduleRelationApp{
		LoanScheduleRelationService: loanScheduleRelationService,
	}
}
