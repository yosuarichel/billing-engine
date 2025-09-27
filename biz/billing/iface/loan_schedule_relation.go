package iface

import (
	"context"

	loanDomain "github.com/yosuarichel/billing-engine/biz/loan/domain"
	"github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/domain"
)

type ILoanScheduleRelationService interface {
	CreateLoanWithSchedules(ctx context.Context, loan *loanDomain.Loan) (res *domain.LoanWithSchedulesData, err error)
}
