package service

import (
	"context"
	"errors"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/billing/domain"
	"github.com/yosuarichel/billing-engine/biz/billing/iface"
	loanDomain "github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanService "github.com/yosuarichel/billing-engine/biz/loan/service"
	paymentService "github.com/yosuarichel/billing-engine/biz/payment/service"
)

type BillingService struct {
	LoanService    iface.ILoanService
	PaymentService iface.IPaymentService
}

func NewBillingService(
	loanService *loanService.LoanService,
	paymentService *paymentService.PaymentService,
) *BillingService {
	return &BillingService{
		LoanService:    loanService,
		PaymentService: paymentService,
	}
}

func (s *BillingService) GetOutstanding(ctx context.Context, loanID *int64) (data *domain.OutstandingData, err error) {
	if gptr.Indirect(loanID) <= 0 {
		return nil, errors.New("missing parameter loanID")
	}

	loanData, err := s.LoanService.GetLoanDetail(ctx, &loanDomain.FindOneLoanParam{
		LoanID: gptr.Indirect(loanID),
		Status: loanDomain.LoanStatus_Ongoing,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][GetOutstanding] Error call GetLoanDetail", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}
	if loanData == nil {
		return nil, errors.New("loan data not found")
	}

	totalLoan := loanData.TotalAmount

	sumAmountPayment, err := s.PaymentService.GetSumAmount(ctx, loanID)
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][GetOutstanding] Error call GetSumAmount", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}

	outstanding := totalLoan - gptr.Indirect(sumAmountPayment)
	if outstanding < 0 {
		outstanding = 0
	}

	data = &domain.OutstandingData{
		LoanID:      loanData.ID,
		CustomerID:  loanData.CustomerID,
		Outstanding: outstanding,
	}

	return data, nil
}
