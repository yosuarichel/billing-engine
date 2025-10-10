package service

import (
	"context"
	"errors"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/kerrors"
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
		return 0, kerrors.ErrPayloadValidation.WithCause(errors.New("input name is required"))
	}

	return s.CustomerRepository.SaveCustomer(ctx, input)
}

func (s *CustomerService) GetCustomerDetail(ctx context.Context, customerID *int64) (customer *domain.Customer, err error) {
	if customerID == nil && gptr.Indirect(customerID) <= 0 {
		return nil, errors.New("customer id is required")
	}

	return s.CustomerRepository.FindOne(ctx, customerID)
}
