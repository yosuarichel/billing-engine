package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/customer/domain"
	"github.com/yosuarichel/billing-engine/kitex_gen/base"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
)

func (a *CustomerApp) CreateCustomer(ctx context.Context, req *billing_engine.CreateCustomerRequest) (res *billing_engine.CreateCustomerResponse) {
	if req == nil {
		return &billing_engine.CreateCustomerResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: errors.New("invalid request").Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	customerID, err := a.CustomerService.CreateCustomer(ctx, &domain.Customer{
		Name:        req.GetName(),
		PhoneNumber: gptr.Of(req.GetPhoneNumber()),
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Customer][App][CreateCustomer] Error call CreateCustomer service", map[string]interface{}{
			"error":  err.Error(),
			"params": req,
		})
		return &billing_engine.CreateCustomerResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	}

	res = &billing_engine.CreateCustomerResponse{
		CustomerId: gptr.Of(gconv.To[string](customerID)),
		BaseResp: &base.BaseResp{
			StatusMessage: "success",
			StatusCode:    http.StatusOK,
		},
	}
	return
}
