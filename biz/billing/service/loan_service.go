package service

import (
	"context"
	"errors"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/billing/iface"
	customerService "github.com/yosuarichel/billing-engine/biz/customer/service"
	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanService "github.com/yosuarichel/billing-engine/biz/loan/service"
	loanScheduleService "github.com/yosuarichel/billing-engine/biz/loan_schedule/service"
	loanScheduleRelationDomain "github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/domain"
	loanScheduleRelationService "github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/service"
	sharedDomain "github.com/yosuarichel/billing-engine/pkg/shared/domain"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

type LoanService struct {
	CustomerService             iface.ICustomerService
	LoanScheduleRelationService iface.ILoanScheduleRelationService
	LoanService                 iface.ILoanService
	LoanScheduleService         iface.ILoanScheduleService
}

func NewLoanService(
	customerService *customerService.CustomerService,
	loanScheduleRelationService *loanScheduleRelationService.LoanScheduleRelationService,
	loanService *loanService.LoanService,
	loanScheduleService *loanScheduleService.LoanScheduleService,
) *LoanService {
	return &LoanService{
		CustomerService:             customerService,
		LoanScheduleRelationService: loanScheduleRelationService,
		LoanService:                 loanService,
		LoanScheduleService:         loanScheduleService,
	}
}

func (s *LoanService) CreateNewLoan(ctx context.Context, input *domain.Loan) (data *loanScheduleRelationDomain.LoanWithSchedulesData, err error) {
	if input.CustomerID <= 0 || input.Principal <= 0 || input.StartDate.IsZero() || input.TermWeeks <= 0 {
		return nil, errors.New("missing parameters")
	}

	customerData, err := s.CustomerService.GetCustomerDetail(ctx, gptr.Of(input.CustomerID))
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][CreateNewLoan] Error call GetCustomerDetail", map[string]interface{}{
			"error": err,
			"input": input,
		})
		return
	}
	if customerData == nil {
		return nil, errors.New("customer not found")
	}

	loanData, err := s.LoanService.GetLoanDetail(ctx, &domain.FindOneLoanParam{
		CustomerID: customerData.ID,
		Status:     domain.LoanStatus_Ongoing,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][CreateNewLoan] Error call GetLoanDetail", map[string]interface{}{
			"error": err,
			"input": input,
		})
		return
	}
	if loanData != nil {
		return nil, errors.New("customer still have ongoing loan")
	}

	input.InterestRate = sharedDomain.FlatAnnualInterestRate
	input.TotalAmount = utils.CalculateTotalAmount(input.Principal, input.InterestRate, input.TermWeeks)

	createdLoan, err := s.LoanScheduleRelationService.CreateLoanWithSchedules(ctx, input)
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][Service][CreateNewLoan] Error call CreateLoanWithSchedules", map[string]interface{}{
			"error": err,
			"input": input,
		})
		return
	}
	return createdLoan, nil
}

func (s *LoanService) IsDelinquent(ctx context.Context, loanID *int64) (isDelinquent bool, err error) {
	return s.LoanScheduleService.IsDelinquentLoan(ctx, loanID)
}
