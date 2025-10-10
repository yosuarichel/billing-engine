package infra

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	"github.com/yosuarichel/billing-engine/pkg/utils"
	"gorm.io/gorm"
)

type LoanScheduleRepository struct {
	db *gorm.DB
}

func NewLoanScheduleRepository(db *gorm.DB) *LoanScheduleRepository {
	return &LoanScheduleRepository{
		db: db,
	}
}

func (r *LoanScheduleRepository) WithTx(tx *gorm.DB) *LoanScheduleRepository {
	return &LoanScheduleRepository{db: tx}
}

func (r *LoanScheduleRepository) SaveLoanSchedules(ctx context.Context, schedules []*domain.LoanSchedule) (ids []int64, err error) {
	now := time.Now()
	ids = make([]int64, len(schedules))

	for i, sched := range schedules {
		if sched.ID == 0 {
			sched.ID = utils.GenerateSonyflakeID()
		}
		sched.CreatedAt = now

		if sched.CreatedBy == "" {
			sched.CreatedBy = "1122334455"
		}

		ids[i] = sched.ID
	}

	db := r.db.WithContext(ctx).Table((&domain.LoanSchedule{}).TableName())
	err = db.Create(&schedules).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Repo][SaveLoanSchedules] Error saving loan schedules to db", map[string]interface{}{
			"error":     err,
			"schedules": schedules,
		})
		return nil, err
	}
	return ids, nil
}

func (r *LoanScheduleRepository) GetDelinquentLoanSchedules(ctx context.Context, loanID *int64) (loanSchedules []*domain.LoanSchedule, err error) {
	if loanID == nil {
		return
	}
	db := r.db.WithContext(ctx).Table((&domain.LoanSchedule{}).TableName())
	db.Where("deleted_at IS NULL")
	db.Where("loan_id = ?", loanID)
	db.Order("week_number ASC")

	err = db.Find(&loanSchedules).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Repo][GetLoanSchedules] Error get loan schedules from db", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}

	return loanSchedules, nil
}

func (r *LoanScheduleRepository) FindAll(ctx context.Context, param *domain.FindAllScheduleParam) (loanSchedules []*domain.LoanSchedule, err error) {
	if param == nil {
		return
	}

	db := r.db.WithContext(ctx).Table((&domain.LoanSchedule{}).TableName())
	db.Where("deleted_at IS NULL")

	if param.LoanScheduleID > 0 {
		db.Where("id = ?", param.LoanScheduleID)
	}
	if param.LoanID > 0 {
		db.Where("loan_id = ?", param.LoanID)
	}
	if param.WeekNumber > 0 {
		db.Where("week_number = ?", param.WeekNumber)
	}
	if param.IsPaid != nil {
		db.Where("is_paid = ?", param.IsPaid)
	}

	err = db.Find(&loanSchedules).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Repo][FindAll] Error finding loan schedules from db", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}
	return loanSchedules, nil
}

func (r *LoanScheduleRepository) UpdateLoanSchedulesToPaid(ctx context.Context, param *domain.UpdateLoanSchedulesToPaidParam) (isUpdated bool, err error) {
	if param == nil {
		return
	}

	db := r.db.WithContext(ctx).Table((&domain.LoanSchedule{}).TableName())
	db.Where("deleted_at IS NULL")

	if param.LoanID > 0 {
		db.Where("loan_id = ?", param.LoanID)
	}
	if param.WeekNumber > 0 {
		db.Where("week_number = ?", param.WeekNumber)
	}
	if param.IsPaid != nil {
		db.Where("is_paid = ?", param.IsPaid)
	}
	err = db.Update("is_paid", true).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Repo][UpdateLoanSchedulesToPaid] Error update loan schedules data to db", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}

	return true, nil
}

func (r *LoanScheduleRepository) CountRemainingWeeks(ctx context.Context, loanID *int64) (remainingWeeks int64, err error) {
	if loanID == nil {
		return 0, errors.New("missing param loanID")
	}

	db := r.db.WithContext(ctx).Table((&domain.LoanSchedule{}).TableName())
	db.Where("deleted_at IS NULL")
	db.Where("loan_id = ?", loanID)
	db.Where("is_paid = ?", false)

	err = db.Count(&remainingWeeks).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[LoanSchedule][Repo][GetLoanSchedules] Error get loan schedules from db", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}

	return remainingWeeks, nil
}
