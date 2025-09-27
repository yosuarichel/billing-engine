package iface

import (
	"context"

	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
	"gorm.io/gorm"
)

type ILoanService interface {
	GetLoanDetail(ctx context.Context, param *domain.FindOneLoanParam) (loan *domain.Loan, err error)
}

type ILoanRepository interface {
	WithTx(tx *gorm.DB) *loanRepo.LoanRepository
	UpdateLoanToPaid(ctx context.Context, loanID *int64) (err error)
}
