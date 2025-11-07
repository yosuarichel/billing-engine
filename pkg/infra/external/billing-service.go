package external

import (
	"context"
	"fmt"
	"sync"

	"github.com/bytedance/gg/gconv"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/billing_customer_service"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/billing_customer_service/billingcustomerservice"
)

var (
	customerClient billingcustomerservice.Client
	once           sync.Once
)

func InitBillingCustomerClient() billingcustomerservice.Client {
	once.Do(func() {
		r, err := consul.NewConsulResolver("consul:8500")
		if err != nil {
			klog.Fatal(err)
		}

		c := billingcustomerservice.MustNewClient(
			"billing-customer-rpc",
			client.WithResolver(r),
		)
		customerClient = c
	})

	return customerClient
}

func CreateCustomer(ctx context.Context, req *billing_customer_service.CreateCustomerRequest) (customerID int64, err error) {
	resp, err := customerClient.CreateCustomer(ctx, req)
	if err != nil || resp.GetBaseResp().GetStatusCode() != 200 {
		klog.CtxErrorf(ctx, "[Customer][Repo][SaveCustomer] Error saving customer to db", map[string]interface{}{
			"error": err,
			"req":   fmt.Sprintf("%+v", req),
		})
		return 0, err
	}

	customerID = gconv.To[int64](resp.GetCustomerId())

	return
}
