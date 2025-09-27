package service

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
)

type LoanScheduleService struct {
	LoanScheduleRepository *loanScheduleRepo.LoanScheduleRepository
}

func NewLoanScheduleService(loanScheduleRepo *loanScheduleRepo.LoanScheduleRepository) *LoanScheduleService {
	return &LoanScheduleService{
		LoanScheduleRepository: loanScheduleRepo,
	}
}

func (s *LoanScheduleService) CreateLoanSchedule(ctx context.Context, input *domain.LoanSchedule) (ids []int64, err error) {
	if input.LoanID <= 0 || input.WeekNumber <= 0 || input.DueDate.IsZero() || input.Amount <= 0 {
		return nil, errors.New("missing parameters")
	}

	return s.LoanScheduleRepository.SaveLoanSchedules(ctx, []*domain.LoanSchedule{input})
}

func (s *LoanScheduleService) IsDelinquentLoan(ctx context.Context, loanID *int64) (isDelinquent bool, err error) {
	if loanID == nil {
		return false, errors.New("missing parameters")
	}

	loanSchedules, err := s.LoanScheduleRepository.GetDelinquentLoanSchedules(ctx, loanID)
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Service][IsDelinquentLoan] Error call GetLoanSchedules", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}
	if len(loanSchedules) == 0 {
		return false, errors.New("loan schedule not found")
	}

	now := time.Now()
	consecutiveMiss := 0

	for _, schedule := range loanSchedules {
		if !schedule.IsPaid && schedule.DueDate.Before(now) {
			consecutiveMiss++
			if consecutiveMiss >= 2 {
				return true, nil
			}
		} else {
			consecutiveMiss = 0
		}
	}

	return false, nil
}

func (s *LoanScheduleService) GetLoanScheduleList(ctx context.Context, param *domain.FindAllScheduleParam) (loanSchedules []*domain.LoanSchedule, err error) {
	if param == nil {
		return nil, errors.New("invalid parameters")
	}

	return s.LoanScheduleRepository.FindAll(ctx, param)
}
