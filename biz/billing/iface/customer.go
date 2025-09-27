package iface

import (
	"context"

	"github.com/yosuarichel/billing-engine/biz/customer/domain"
)

type ICustomerService interface {
	GetCustomerDetail(ctx context.Context, customerID *int64) (customer *domain.Customer, err error)
}
