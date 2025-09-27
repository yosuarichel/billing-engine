package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/billing/domain"
	"github.com/yosuarichel/billing-engine/kitex_gen/base"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine/data/payment_data"
)

func (a *BillingApp) MakePayment(ctx context.Context, req *billing_engine.MakePaymentRequest) (res *billing_engine.MakePaymentResponse) {
	if req == nil {
		return &billing_engine.MakePaymentResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: errors.New("invalid request").Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	paymentData, err := a.BillingPaymentService.MakePayment(ctx, &domain.MakePaymentParam{
		LoanID: gconv.To[int64](req.GetLoanId()),
		Amount: req.GetAmount(),
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Billing][App][MakePayment] Error call MakePayment", map[string]interface{}{
			"error":  err.Error(),
			"params": req,
		})
		return &billing_engine.MakePaymentResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	res = &billing_engine.MakePaymentResponse{
		BaseResp: &base.BaseResp{
			StatusMessage: "success",
			StatusCode:    http.StatusOK,
		},
	}

	if paymentData != nil {
		res.Data = &payment_data.PaymentData{
			Amount:      gptr.Of(paymentData.Amount),
			PayAmount:   gptr.Of(paymentData.PayAmount),
			Outstanding: gptr.Of(paymentData.Outstanding),
			TermWeeks:   gptr.Of(int32(paymentData.TermWeeks)),
			WeeksRemain: gptr.Of(int32(paymentData.WeeksRemain)),
		}
	}

	return
}
