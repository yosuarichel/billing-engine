package rpc_handler

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	billingApp "github.com/yosuarichel/billing-engine/biz/billing/app"
	customerApp "github.com/yosuarichel/billing-engine/biz/customer/app"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
)

type RpcHandler struct {
	CustomerApp *customerApp.CustomerApp
	BillingApp  *billingApp.BillingApp
}

func NewRpcHandler(
	customerApp *customerApp.CustomerApp,
	billingApp *billingApp.BillingApp,
) *RpcHandler {
	return &RpcHandler{
		CustomerApp: customerApp,
		BillingApp:  billingApp,
	}
}

func (h *RpcHandler) CreateCustomer(ctx context.Context, req *billing_engine.CreateCustomerRequest) (res *billing_engine.CreateCustomerResponse, err error) {
	klog.CtxInfof(ctx, "[RPC Handler CreateCustomer]")
	return h.CustomerApp.CreateCustomer(ctx, req), nil
}

func (h *RpcHandler) IsDelinquent(ctx context.Context, req *billing_engine.IsDelinquentRequest) (res *billing_engine.IsDelinquentResponse, err error) {
	klog.CtxInfof(ctx, "[RPC Handler IsDelinquent]")
	return h.BillingApp.IsDelinquent(ctx, req), nil
}

func (h *RpcHandler) CreateLoan(ctx context.Context, req *billing_engine.CreateLoanRequest) (res *billing_engine.CreateLoanResponse, err error) {
	klog.CtxInfof(ctx, "[RPC Handler CreateLoan]")
	return h.BillingApp.CreateLoan(ctx, req), nil
}

func (h *RpcHandler) GetOutstanding(ctx context.Context, req *billing_engine.GetOutstandingRequest) (res *billing_engine.GetOutstandingResponse, err error) {
	klog.CtxInfof(ctx, "[RPC Handler GetOutstanding]")
	return h.BillingApp.GetOutstanding(ctx, req), nil
}

func (h *RpcHandler) MakePayment(ctx context.Context, req *billing_engine.MakePaymentRequest) (res *billing_engine.MakePaymentResponse, err error) {
	klog.CtxInfof(ctx, "[RPC Handler MakePayment]")
	return h.BillingApp.MakePayment(ctx, req), nil
}

// func (h *RpcHandler) GetProductList(ctx context.Context, req *billing_engine.GetProductListRequest) (res *billing_engine.GetProductListResponse, err error) {
// 	klog.CtxInfof(ctx, "[RPC Handler GetProductList]")
// 	return h.CustomerApp.GetProductList(ctx, req), nil
// }
