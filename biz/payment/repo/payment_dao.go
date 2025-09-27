package infra

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/payment/domain"
	"github.com/yosuarichel/billing-engine/pkg/utils"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) WithTx(tx *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: tx}
}

func (r *PaymentRepository) GetSumAmount(ctx context.Context, loanID *int64) (sumAmount *int64, err error) {
	if loanID == nil {
		return
	}

	db := r.db.WithContext(ctx).Table((&domain.Payment{}).TableName())
	db.Select("COALESCE(SUM(amount), 0)")
	db.Where("deleted_at IS NULL")
	db.Where("loan_id = ?", loanID)

	err = db.Scan(&sumAmount).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Payment][Repo][GetSumAmount] Error get sum ammount of payment from db", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}
	return sumAmount, nil
}

func (r *PaymentRepository) FindAll(ctx context.Context, param *domain.FindAllParam) (paymentData []*domain.Payment, err error) {
	if param == nil {
		return
	}

	db := r.db.WithContext(ctx).Table((&domain.Payment{}).TableName())
	db.Where("deleted_at IS NULL")

	if param.PaymentID > 0 {
		db.Where("id = ?", param.PaymentID)
	}
	if param.LoanID > 0 {
		db.Where("loan_id = ?", param.LoanID)
	}
	if param.ScheduleID > 0 {
		db.Where("schedule_id = ?", param.ScheduleID)
	}

	date, _ := time.Parse(time.DateOnly, param.PaymentDate)
	if !date.IsZero() {
		db.Where("payment_date = ?", date)
	}

	err = db.Find(&paymentData).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Payment][Repo][FindAll] Error payments data from db", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}
	return paymentData, nil
}

func (r *PaymentRepository) BulkSavePayment(ctx context.Context, payments []*domain.Payment) (ids []int64, err error) {
	now := time.Now()
	ids = make([]int64, len(payments))

	for i, payment := range payments {
		if payment.ID == 0 {
			payment.ID = utils.GenerateSonyflakeID()
		}
		payment.CreatedAt = now

		if payment.CreatedBy == "" {
			payment.CreatedBy = "1122334455"
		}

		ids[i] = payment.ID
	}

	db := r.db.WithContext(ctx).Table((&domain.Payment{}).TableName())
	err = db.Create(&payments).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Customer][Repo][BulkSavePayment] Error saving payments to db", map[string]interface{}{
			"error":    err,
			"payments": payments,
		})
		return nil, err
	}
	return ids, nil
}
