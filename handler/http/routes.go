package http_handler

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, handler *HttpHandler) {
	h.POST("/api/v1/customer/create", handler.CreateCustomer)
	h.POST("/api/v1/billing/is_delinquent", handler.IsDelinquent)
	h.POST("/api/v1/billing/loan/create", handler.CreateLoan)
	h.POST("/api/v1/billing/get_outstanding", handler.GetOutstanding)
	h.POST("api/v1/billing/payment/create", handler.MakePayment)
}
