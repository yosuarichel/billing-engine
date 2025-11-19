package external

import (
	"context"
	"fmt"
	"sync"

	"github.com/bytedance/gg/gconv"
	"github.com/bytedance/gg/gslice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	dns "github.com/kitex-contrib/resolver-dns"
	"github.com/yosuarichel/billing-engine/pkg/config"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/billing_customer_service"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/billing_customer_service/billingcustomerservice"
	"github.com/yosuarichel/idl_gen_billing_customer_service/kitex_gen/billing/billing_customer/data/customer_data"
)

var (
	customerClient billingcustomerservice.Client
	once           sync.Once
	cfg            = config.GetAppCfg()
)

func InitBillingCustomerClient() billingcustomerservice.Client {
	once.Do(func() {
		host := gslice.Find(cfg.Upstreams, func(u config.UpstreamConfig) bool {
			return u.Name == "billing-customer-rpc"
		})
		if !host.IsOK() {
			klog.Fatal("billing-customer-rpc upstream not found")
		}
		c := billingcustomerservice.MustNewClient(
			"billing-customer-rpc",
			// client.WithHostPorts(fmt.Sprintf("%s:%d", host.Value().Host, host.Value().Port)),
			// client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
			// client.WithRetryPolicy(retry.NewFailurePolicy()),
			// client.WithConnectTimeout(3*time.Second),
			client.WithResolver(dns.NewDNSResolver()),
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

func GetCustomers(ctx context.Context, req *billing_customer_service.GetCustomerListRequest) (custumers []*customer_data.CustomerData, err error) {
	resp, err := customerClient.GetCustomerList(ctx, req)
	if err != nil || resp.GetBaseResp().GetStatusCode() != 200 {
		klog.CtxErrorf(ctx, "[Customer][Repo][SaveCustomer] Error saving customer to db", map[string]interface{}{
			"error": err,
			"req":   fmt.Sprintf("%+v", req),
		})
		return nil, err
	}
	if resp.GetData() == nil {
		return nil, nil
	}

	custumers = resp.GetData()
	return
}
