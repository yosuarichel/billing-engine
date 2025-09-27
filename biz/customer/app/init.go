package app

import (
	customerService "github.com/yosuarichel/billing-engine/biz/customer/service"
)

type CustomerApp struct {
	CustomerService *customerService.CustomerService
}

func NewCustomerApp(
	customerService *customerService.CustomerService,
) *CustomerApp {
	return &CustomerApp{
		CustomerService: customerService,
	}
}
