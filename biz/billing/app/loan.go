package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/loan/domain"
	loanScheduleDomain "github.com/yosuarichel/billing-engine/biz/loan_schedule/domain"
	"github.com/yosuarichel/billing-engine/kitex_gen/base"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/data/loan_schedule_data"
)

func (a *BillingApp) CreateLoan(ctx context.Context, req *billing_engine.CreateLoanRequest) (res *billing_engine.CreateLoanResponse) {
	req, err := a.validateCreateLoan(ctx, req)
	if err != nil {
		return &billing_engine.CreateLoanResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
			Schedules: nil,
		}
	}

	now := time.Now()
	loanWithSchedulesData, err := a.BillingLoanService.CreateNewLoan(ctx, &domain.Loan{
		CustomerID: gconv.To[int64](req.GetCustomerId()),
		Principal:  req.GetPrincipal(),
		StartDate:  now,
		TermWeeks:  int(req.GetTermWeeks()),
		Status:     domain.LoanStatus_Ongoing,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][App][CreateLoan] Error call CreateNewLoan", map[string]interface{}{
			"error":  err.Error(),
			"params": req,
		})
		return &billing_engine.CreateLoanResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
			Schedules: nil,
		}
	}

	res = &billing_engine.CreateLoanResponse{
		LoanId:      gptr.Of(gconv.To[string](loanWithSchedulesData.ID)),
		CustomerId:  gptr.Of(gconv.To[string](loanWithSchedulesData.CustomerID)),
		Principal:   gptr.Of(loanWithSchedulesData.Principal),
		TotalAmount: gptr.Of(loanWithSchedulesData.TotalAmount),
		TermWeeks:   gptr.Of(int32(loanWithSchedulesData.TermWeeks)),
		StartDate:   gptr.Of(loanWithSchedulesData.StartDate.Unix()),
		Status:      gptr.Of(loanWithSchedulesData.Status),
		BaseResp: &base.BaseResp{
			StatusMessage: "success",
			StatusCode:    http.StatusOK,
		},
	}

	schedules := gslice.Map(loanWithSchedulesData.Schedules, func(schedule *loanScheduleDomain.LoanSchedule) *loan_schedule_data.LoanScheduleSummaryData {
		return &loan_schedule_data.LoanScheduleSummaryData{
			ScheduleId: gptr.Of(schedule.ID),
			WeekNumber: gptr.Of(int32(schedule.WeekNumber)),
			DueDate:    gptr.Of(schedule.DueDate.Unix()),
			Amount:     gptr.Of(schedule.Amount),
			IsPaid:     gptr.Of(schedule.IsPaid),
		}
	})

	if len(schedules) > 0 {
		res.Schedules = schedules
	}
	return
}

func (a *BillingApp) validateCreateLoan(ctx context.Context, req *billing_engine.CreateLoanRequest) (newReq *billing_engine.CreateLoanRequest, err error) {
	if req == nil {
		return nil, errors.New("invalid request")
	}

	customerID := gconv.To[int64](req.GetCustomerId())
	if customerID <= 0 {
		return nil, errors.New("missing customer id")
	}
	if req.GetPrincipal() <= 0 {
		return nil, errors.New("missing principal")
	}
	if req.GetTermWeeks() <= 0 {
		return nil, errors.New("missing term weeks")
	}

	return req, nil
}

func (a *BillingApp) IsDelinquent(ctx context.Context, req *billing_engine.IsDelinquentRequest) (res *billing_engine.IsDelinquentResponse) {
	if req == nil {
		return &billing_engine.IsDelinquentResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: errors.New("invalid request").Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	isDelinquent, err := a.BillingLoanService.IsDelinquent(ctx, gptr.Of(gconv.To[int64](req.GetLoanId())))
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][App][IsDelinquent] Error call IsDelinquent", map[string]interface{}{
			"error":  err.Error(),
			"params": req,
		})
		return &billing_engine.IsDelinquentResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	return &billing_engine.IsDelinquentResponse{
		BaseResp: &base.BaseResp{
			StatusMessage: "success",
			StatusCode:    http.StatusOK,
		},
		IsDelinquent: gptr.Of(isDelinquent),
	}
}
