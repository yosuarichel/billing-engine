package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/kitex_gen/base"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
	"github.com/yosuarichel/billing-engine/pkg/infra/external"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/billing_customer_service"
)

func (a *CustomerApp) CreateCustomer(ctx context.Context, req *billing_engine.CreateCustomerRequest) (res *billing_engine.CreateCustomerResponse) {
	if req == nil {
		return &billing_engine.CreateCustomerResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: errors.New("invalid request").Error(),
				StatusCode:    http.StatusBadRequest,
			},
		}
	}

	customerID, err := external.CreateCustomer(ctx, &billing_customer_service.CreateCustomerRequest{
		Name:        req.GetName(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Customer][App][CreateCustomer] Error call CreateCustomer service", map[string]interface{}{
			"error":  err.Error(),
			"params": fmt.Sprintf("%+v", req),
		})

		return &billing_engine.CreateCustomerResponse{
			BaseResp: &base.BaseResp{
				StatusMessage: err.Error(),
				StatusCode:    http.StatusBadRequest,
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
