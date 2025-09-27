package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/billing/domain"
	"github.com/yosuarichel/billing-engine/biz/billing/iface"
	loanDomain "github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
	loanService "github.com/yosuarichel/billing-engine/biz/loan/service"
	loanScheduleDomain "github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
	loanScheduleService "github.com/yosuarichel/billing-engine/biz/loan_schedule/service"
	paymentDomain "github.com/yosuarichel/billing-engine/biz/payment/domain"
	paymentRepo "github.com/yosuarichel/billing-engine/biz/payment/repo"
	paymentService "github.com/yosuarichel/billing-engine/biz/payment/service"
	"gorm.io/gorm"
)

type PaymentService struct {
	db                  *gorm.DB
	LoanService         iface.ILoanService
	LoanScheduleService iface.ILoanScheduleService
	PaymentService      iface.IPaymentService

	// repo
	PaymentRepository      iface.IPaymentRepository
	LoanRepository         iface.ILoanRepository
	LoanScheduleRepository iface.ILoanScheduleRepository
}

func NewPaymentService(
	db *gorm.DB,
	loanService *loanService.LoanService,
	paymentService *paymentService.PaymentService,
	loanScheduleService *loanScheduleService.LoanScheduleService,
	paymentRepo *paymentRepo.PaymentRepository,
	loanRepo *loanRepo.LoanRepository,
	loanScheduleRepo *loanScheduleRepo.LoanScheduleRepository,
) *PaymentService {
	return &PaymentService{
		db:                     db,
		LoanService:            loanService,
		PaymentService:         paymentService,
		LoanScheduleService:    loanScheduleService,
		PaymentRepository:      paymentRepo,
		LoanRepository:         loanRepo,
		LoanScheduleRepository: loanScheduleRepo,
	}
}

func (s *PaymentService) MakePayment(ctx context.Context, param *domain.MakePaymentParam) (data *domain.PaymentData, err error) {
	if param == nil {
		return nil, errors.New("invalid parameter")
	}
	if param.Amount <= 0 || param.LoanID <= 0 {
		return nil, errors.New("missing parameter")
	}

	data = &domain.PaymentData{}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Get loan detail
		loanData, err2 := s.LoanService.GetLoanDetail(ctx, &loanDomain.FindOneLoanParam{
			LoanID: param.LoanID,
			Status: loanDomain.LoanStatus_Ongoing,
		})
		if err2 != nil {
			klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call GetLoanDetail", map[string]interface{}{
				"error": err2,
				"param": param,
			})
			return err2
		}
		if loanData == nil {
			return errors.New("loan data not found")
		}

		// set to data response
		data.Amount = loanData.TotalAmount
		data.PayAmount = param.Amount
		data.TermWeeks = loanData.TermWeeks

		days := int(time.Since(loanData.StartDate).Hours() / 24)
		currentWeek := (days / 7) + 1

		// Get loan schedules
		loanScheduleData, err2 := s.LoanScheduleService.GetLoanScheduleList(ctx, &loanScheduleDomain.FindAllScheduleParam{
			LoanID:     param.LoanID,
			WeekNumber: currentWeek,
			IsPaid:     gptr.Of(false),
		})
		if err2 != nil {
			klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call GetLoanScheduleList", map[string]interface{}{
				"error": err2,
				"param": param,
			})
			return err2
		}
		if len(loanScheduleData) == 0 {
			return fmt.Errorf("no dues pending for loan %d up to week %d", param.LoanID, currentWeek)
		}

		// Sum total due amount
		var dueTotal int64
		for _, schedule := range loanScheduleData {
			dueTotal += schedule.Amount
		}

		data.Outstanding = data.Amount - dueTotal

		// Check dueTotal with pay amount
		if dueTotal != param.Amount {
			return fmt.Errorf("invalid payment: expected %d, got %d", dueTotal, param.Amount)
		}

		// change to tx
		paymentRepoTx := s.PaymentRepository.WithTx(tx)
		loanRepoTx := s.LoanRepository.WithTx(tx)
		loanScheduleRepoTx := s.LoanScheduleRepository.WithTx(tx)

		now := time.Now()

		// Create payment
		paymentParam := gslice.Map(loanScheduleData, func(schedule *loanScheduleDomain.LoanSchedule) *paymentDomain.Payment {
			return &paymentDomain.Payment{
				LoanID:      param.LoanID,
				ScheduleID:  schedule.ID,
				PaymentDate: now,
				Amount:      schedule.Amount,
			}
		})
		_, err2 = paymentRepoTx.BulkSavePayment(ctx, paymentParam)
		if err2 != nil {
			klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call BulkSavePayment", map[string]interface{}{
				"error": err2,
				"param": param,
			})
			return err2
		}

		// Update loan schedules data to PAID
		isUpdated, err2 := loanScheduleRepoTx.UpdateLoanSchedulesToPaid(ctx, &loanScheduleDomain.UpdateLoanSchedulesToPaidParam{
			LoanID:     param.LoanID,
			WeekNumber: currentWeek,
			IsPaid:     gptr.Of(false),
		})
		if err2 != nil {
			klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call UpdateLoanSchedulesToPaid", map[string]interface{}{
				"error": err2,
				"param": param,
			})
			return err2
		}
		if !isUpdated {
			return errors.New("failed to updated loan schedules")
		}

		remainingWeeks, err2 := loanScheduleRepoTx.CountRemainingWeeks(ctx, gptr.Of(param.LoanID))
		if err2 != nil {
			klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call CountRemainingWeeks", map[string]interface{}{
				"error": err2,
				"param": param,
			})
			return err2
		}

		if remainingWeeks == 0 {
			err2 := loanRepoTx.UpdateLoanToPaid(ctx, gptr.Of(param.LoanID))
			if err2 != nil {
				klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error call UpdateLoanToPaid", map[string]interface{}{
					"error": err2,
					"param": param,
				})
				return err2
			}
		}

		data.WeeksRemain = int(remainingWeeks)

		// Commit
		return nil
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][MakePayment] Error save processing make payment", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}

	return data, nil
}
