package infra

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	"github.com/yosuarichel/billing-engine/pkg/utils"
	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{
		db: db,
	}
}

func (r *LoanRepository) WithTx(tx *gorm.DB) *LoanRepository {
	return &LoanRepository{db: tx}
}

func (r *LoanRepository) SaveLoan(ctx context.Context, loan *domain.Loan) (id int64, err error) {
	now := time.Now()
	if loan.ID == 0 {
		loan.ID = utils.GenerateSonyflakeID()
	}
	loan.CreatedAt = now

	// Set default
	if loan.CreatedBy == "" {
		loan.CreatedBy = "1122334455"
	}

	db := r.db.WithContext(ctx).Table(loan.TableName())
	err = db.Create(&loan).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Loan][Repo][SaveLoan] Error saving loan to db", map[string]interface{}{
			"error": err,
			"loan":  loan,
		})
		return 0, err
	}
	return loan.ID, nil
}

func (r *LoanRepository) FindOne(ctx context.Context, param *domain.FindOneLoanParam) (loan *domain.Loan, err error) {
	if param == nil {
		return
	}

	klog.Info("MASUUUUKKKK SINIIII star get loan data")

	db := r.db.WithContext(ctx).Table(loan.TableName())
	db.Where("deleted_at IS NULL")

	if param.CustomerID > 0 {
		klog.Info("MASUUUUKKKK SINIIII customer")
		db.Where("customer_id = ?", param.CustomerID)
	}
	if param.LoanID > 0 {
		klog.Info("MASUUUUKKKK SINIIII loan")
		db.Where("id = ?", param.LoanID)
	}
	if param.Status != "" {
		klog.Info("MASUUUUKKKK SINIIII status")
		db.Where("status = ?", param.Status)
	}

	klog.Info("MASUUUUKKKK SINIIII get loan data")

	err = db.First(&loan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		klog.CtxErrorf(ctx, "[Loan][Repo][FindOne] Error finding loan from db", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}

	klog.Info("MASUUUUKKKK SINIIII berhasil loan data")
	return loan, nil
}

func (r *LoanRepository) UpdateLoanToPaid(ctx context.Context, loanID *int64) (err error) {
	if loanID == nil {
		return
	}

	db := r.db.WithContext(ctx).Table((&domain.Loan{}).TableName())
	db.Where("deleted_at IS NULL")
	db.Where("id = ?", loanID)

	err = db.Update("status", domain.LoanStatus_Paid).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Loan][Repo][UpdateLoanToPaid] Error update loan data to db", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}

	return nil
}
