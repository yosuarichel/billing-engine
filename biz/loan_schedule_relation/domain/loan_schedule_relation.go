package domain

import (
	loanDomain "github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanScheduleDomain "github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
)

type LoanWithSchedulesData struct {
	*loanDomain.Loan
	Schedules []*loanScheduleDomain.LoanSchedule `json:"schedules"`
}
