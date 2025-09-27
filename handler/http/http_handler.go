package http_handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/klog"

	billingApp "github.com/yosuarichel/billing-engine/biz/billing/app"
	customerApp "github.com/yosuarichel/billing-engine/biz/customer/app"
	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
)

type HttpHandler struct {
	CustomerApp *customerApp.CustomerApp
	BillingApp  *billingApp.BillingApp
}

func NewHttpHandler(
	customerApp *customerApp.CustomerApp,
	billingApp *billingApp.BillingApp,
) *HttpHandler {
	return &HttpHandler{
		CustomerApp: customerApp,
		BillingApp:  billingApp,
	}
}

func (h *HttpHandler) CreateCustomer(c context.Context, ctx *app.RequestContext) {
	var req billing_engine.CreateCustomerRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	klog.CtxInfof(c, "[HTTP Handler GetTestData]")

	res := h.CustomerApp.CreateCustomer(c, &req)

	ctx.JSON(http.StatusOK, &res)
}

func (h *HttpHandler) IsDelinquent(c context.Context, ctx *app.RequestContext) {
	var req billing_engine.IsDelinquentRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	klog.CtxInfof(c, "[HTTP Handler GetTestData]")

	res := h.BillingApp.IsDelinquent(c, &req)

	ctx.JSON(http.StatusOK, &res)
}

func (h *HttpHandler) CreateLoan(c context.Context, ctx *app.RequestContext) {
	var req billing_engine.CreateLoanRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	klog.CtxInfof(c, "[HTTP Handler GetTestData]")

	res := h.BillingApp.CreateLoan(c, &req)

	ctx.JSON(http.StatusOK, &res)
}

func (h *HttpHandler) GetOutstanding(c context.Context, ctx *app.RequestContext) {
	var req billing_engine.GetOutstandingRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	klog.CtxInfof(c, "[HTTP Handler GetTestData]")

	res := h.BillingApp.GetOutstanding(c, &req)

	ctx.JSON(http.StatusOK, &res)
}

func (h *HttpHandler) MakePayment(c context.Context, ctx *app.RequestContext) {
	var req billing_engine.MakePaymentRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	klog.CtxInfof(c, "[HTTP Handler GetTestData]")

	res := h.BillingApp.MakePayment(c, &req)

	ctx.JSON(http.StatusOK, &res)
}
