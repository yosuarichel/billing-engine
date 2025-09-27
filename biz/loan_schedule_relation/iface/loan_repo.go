package iface

import (
	"context"

	"gorm.io/gorm"

	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
)

type ILoanRepository interface {
	WithTx(tx *gorm.DB) *loanRepo.LoanRepository
	SaveLoan(ctx context.Context, loan *domain.Loan) (id int64, err error)
}
