package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"

	customerRepo "github.com/yosuarichel/billing-engine/biz/customer/repo"
	loanRepo "github.com/yosuarichel/billing-engine/biz/loan/repo"
	loanScheduleRepo "github.com/yosuarichel/billing-engine/biz/loan_schedule/repo"
	paymentRepo "github.com/yosuarichel/billing-engine/biz/payment/repo"

	billingSvc "github.com/yosuarichel/billing-engine/biz/billing/service"
	customerSvc "github.com/yosuarichel/billing-engine/biz/customer/service"
	loanSvc "github.com/yosuarichel/billing-engine/biz/loan/service"
	loanScheduleSvc "github.com/yosuarichel/billing-engine/biz/loan_schedule/service"
	loanScheduleRelationSvc "github.com/yosuarichel/billing-engine/biz/loan_schedule_relation/service"
	paymentSvc "github.com/yosuarichel/billing-engine/biz/payment/service"

	billingApp "github.com/yosuarichel/billing-engine/biz/billing/app"
	customerApp "github.com/yosuarichel/billing-engine/biz/customer/app"

	httpHandler "github.com/yosuarichel/billing-engine/handler/http"
	rpcHandler "github.com/yosuarichel/billing-engine/handler/rpc"
	"github.com/yosuarichel/billing-engine/internal/bootstrap"
	"github.com/yosuarichel/billing-engine/internal/server"
	"github.com/yosuarichel/billing-engine/pkg/infra/db"
	"github.com/yosuarichel/billing-engine/pkg/utils"
)

func main() {
	klog.SetLogger(kitexlogrus.NewLogger())
	klog.SetLevel(klog.LevelDebug)

	deps := bootstrap.Init()

	// Repo
	customerRepo := customerRepo.NewCustomerRepository(db.GetDB())
	loanRepo := loanRepo.NewLoanRepository(db.GetDB())
	loanScheduleRepo := loanScheduleRepo.NewLoanScheduleRepository(db.GetDB())
	paymentRepo := paymentRepo.NewPaymentRepository(db.GetDB())

	// Service
	customerService := customerSvc.NewCustomerService(customerRepo)
	loanScheduleRelationService := loanScheduleRelationSvc.NewLoanScheduleRelationService(db.GetDB(), loanRepo, loanScheduleRepo)
	loanService := loanSvc.NewLoanService(loanRepo)
	loanScheduleService := loanScheduleSvc.NewLoanScheduleService(loanScheduleRepo)
	billingLoanService := billingSvc.NewLoanService(customerService, loanScheduleRelationService, loanService, loanScheduleService)
	paymentService := paymentSvc.NewPaymentService(paymentRepo)
	billingService := billingSvc.NewBillingService(loanService, paymentService)
	billingPaymentService := billingSvc.NewPaymentService(db.GetDB(), loanService, paymentService, loanScheduleService, paymentRepo, loanRepo, loanScheduleRepo)

	// App
	customerApp := customerApp.NewCustomerApp(customerService)
	billingApp := billingApp.NewBillingApp(billingLoanService, billingService, billingPaymentService)

	switch utils.GetAppType() {
	case utils.RPC_APP_TYPE:
		handler := rpcHandler.NewRpcHandler(customerApp, billingApp)
		server.StartRPC(deps.Ctx, deps.Cfg, handler)
	case utils.HTTP_APP_TYPE:
		handler := httpHandler.NewHttpHandler(customerApp, billingApp)
		server.StartHTTP(deps.Ctx, deps.Cfg, handler)
	case utils.CRON_APP_TYPE:
		klog.CtxInfof(deps.Ctx, "CRON APP TYPE")
	case utils.CONSUMER_APP_TYPE:
		klog.CtxInfof(deps.Ctx, "CONSUMER APP TYPE")
	default:
		klog.CtxInfof(deps.Ctx, "no matching app type, service not started")
	}
}
