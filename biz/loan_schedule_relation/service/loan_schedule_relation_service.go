package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	loanDomain "github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
	loanScheduleDomain "github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
	"github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/domain"
	"github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/iface"
	"gorm.io/gorm"
)

type LoanScheduleRelationService struct {
	db                     *gorm.DB
	LoanRepository         iface.ILoanRepository
	LoanScheduleRepository iface.ILoanScheduleRepository
}

func NewLoanScheduleRelationService(
	db *gorm.DB,
	loanRepo *loanRepo.LoanRepository,
	loanScheduleRepo *loanScheduleRepo.LoanScheduleRepository,
) *LoanScheduleRelationService {
	return &LoanScheduleRelationService{
		db:                     db,
		LoanRepository:         loanRepo,
		LoanScheduleRepository: loanScheduleRepo,
	}
}

func (s *LoanScheduleRelationService) CreateLoanWithSchedules(ctx context.Context, loan *loanDomain.Loan) (res *domain.LoanWithSchedulesData, err error) {
	// Do a transaction
	schedules := []*loanScheduleDomain.LoanSchedule{}
	err = s.db.Transaction(func(tx *gorm.DB) error {
		loanRepoTx := s.LoanRepository.WithTx(tx)
		scheduleRepoTx := s.LoanScheduleRepository.WithTx(tx)

		// Create loan first
		loanID, err2 := loanRepoTx.SaveLoan(ctx, loan)
		if err2 != nil {
			return err2
		}

		// Create loan schedules
		weeklyAmount := loan.TotalAmount / int64(loan.TermWeeks)
		remainder := loan.TotalAmount % int64(loan.TermWeeks)

		for i := 1; i <= loan.TermWeeks; i++ {
			schedules = append(schedules, &loanScheduleDomain.LoanSchedule{
				LoanID:     loanID,
				WeekNumber: i,
				DueDate:    loan.StartDate.AddDate(0, 0, 7*i),
				Amount:     weeklyAmount,
			})
		}

		if remainder > 0 && len(schedules) == loan.TermWeeks {
			schedules[loan.TermWeeks-1].Amount += remainder
		}

		_, err2 = scheduleRepoTx.SaveLoanSchedules(ctx, schedules)
		if err2 != nil {
			return err2
		}

		// Commit
		return nil
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanScheduleRelation][Service][CreateLoanWithSchedules] Error save loan with schedules", map[string]interface{}{
			"error": err,
			"loan":  loan,
		})
		return
	}

	res = &domain.LoanWithSchedulesData{
		Loan:      loan,
		Schedules: schedules,
	}

	return res, nil

}
