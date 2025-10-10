package service

import (
	"context"
	"errors"

	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

type LoanService struct {
	LoanRepository *loanRepo.LoanRepository
}

func NewLoanService(loanRepo *loanRepo.LoanRepository) *LoanService {
	return &LoanService{
		LoanRepository: loanRepo,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, input *domain.Loan) (id int64, err error) {
	if input.CustomerID <= 0 || input.Principal <= 0 || input.StartDate.IsZero() || input.TermWeeks <= 0 {
		return 0, errors.New("missing parameters")
	}

	input.TotalAmount = utils.CalculateTotalAmount(input.Principal, input.InterestRate, input.TermWeeks)

	return s.LoanRepository.SaveLoan(ctx, input)
}

func (s *LoanService) GetLoanDetail(ctx context.Context, param *domain.FindOneLoanParam) (loan *domain.Loan, err error) {
	if param == nil {
		return
	}

	return s.LoanRepository.FindOne(ctx, param)
}
