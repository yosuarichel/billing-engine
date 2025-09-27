package iface

import (
	"context"

	"github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
	"gorm.io/gorm"
)

type ILoanScheduleService interface {
	IsDelinquentLoan(ctx context.Context, loanID *int64) (isDelinquent bool, err error)
	GetLoanScheduleList(ctx context.Context, param *domain.FindAllScheduleParam) (loanSchedules []*domain.LoanSchedule, err error)
}

type ILoanScheduleRepository interface {
	WithTx(tx *gorm.DB) *loanScheduleRepo.LoanScheduleRepository
	UpdateLoanSchedulesToPaid(ctx context.Context, param *domain.UpdateLoanSchedulesToPaidParam) (isUpdated bool, err error)
	CountRemainingWeeks(ctx context.Context, loanID *int64) (remainingWeeks int64, err error)
}
