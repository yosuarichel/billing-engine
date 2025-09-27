package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/kitex_gen/base"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/data/billing_data"
)

func (a *BillingApp) GetOutstanding(ctx context.Context, req *billing_engine.GetOutstandingRequest) (res *billing_engine.GetOutstandingResponse) {
	if req == nil {
		return &billing_engine.GetOutstandingResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: errors.New("invalid request").Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	outstandingData, err := a.BillingService.GetOutstanding(ctx, gptr.Of(gconv.To[int64](req.GetLoanId())))
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][App][GetOutstanding] Error call GetOutstanding", map[string]interface{}{
			"error":  err.Error(),
			"params": req,
		})
		return &billing_engine.GetOutstandingResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	res = &billing_engine.GetOutstandingResponse{
		BaseResp: &base.BaseResp{
			StatusMessage: "success",
			StatusCode:    http.StatusOK,
		},
	}

	if outstandingData != nil {
		res.Data = &billing_data.OutstandingData{
			LoanId:      gptr.Of(gconv.To[string](outstandingData.LoanID)),
			CustomerId:  gptr.Of(gconv.To[string](outstandingData.CustomerID)),
			Outstanding: gptr.Of(outstandingData.Outstanding),
		}
	}

	return
}
