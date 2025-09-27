package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytedance/gg/gptr"
	"github.com/yosuarichel/billing-engine/biz/customer/domain"
	customerRepo "github.com/yosuarichel/billing-engine/biz/customer/repo"
)

type CustomerService struct {
	CustomerRepository *customerRepo.CustomerRepository
}

func NewCustomerService(customerRepo *customerRepo.CustomerRepository) *CustomerService {
	return &CustomerService{
		CustomerRepository: customerRepo,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, input *domain.Customer) (id int64, err error) {
	if input.Name == "" {
		return 0, fmt.Errorf("input name is required")
	}

	return s.CustomerRepository.SaveCustomer(ctx, input)
}

func (s *CustomerService) GetCustomerDetail(ctx context.Context, customerID *int64) (customer *domain.Customer, err error) {
	if customerID == nil && gptr.Indirect(customerID) <= 0 {
		return nil, errors.New("customer id is required")
	}

	return s.CustomerRepository.FindOne(ctx, customerID)
}

// func (s *CustomerService) GetCustomerList(ctx context.Context, param *domain.GetCustomerListParam) (customers []*domain.Customer, err error) {
// 	if param == nil {
// 		return
// 	}

// 	customerList, err := s.CustomerRepository.FindAllCustomers(ctx, &domain.FindAllCustomerFilterParam{
// 		CustomerIDs: param.CustomerIDs,
// 		Name:       param.Name,
// 	})
// 	if err != nil {
// 		klog.CtxErrorf(ctx, "[Customer][Service][GetCustomerList] Error call FindAllCustomers repo", map[string]interface{}{
// 			"error": err,
// 			"param": param,
// 		})
// 		return
// 	}

// 	return customerList, nil
// }
