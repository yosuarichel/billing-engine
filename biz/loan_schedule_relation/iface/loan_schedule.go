package iface

import (
	"context"

	"gorm.io/gorm"

	"github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
)

type ILoanScheduleRepository interface {
	WithTx(tx *gorm.DB) *loanScheduleRepo.LoanScheduleRepository
	SaveLoanSchedules(ctx context.Context, schedules []*domain.LoanSchedule) (ids []int64, err error)
}
