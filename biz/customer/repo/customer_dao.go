package infra

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/customer/domain"
	"github.com/yosuarichel/billing-engine/pkg/utils"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) SaveCustomer(ctx context.Context, customer *domain.Customer) (id int64, err error) {
	now := time.Now()
	if customer.ID == 0 {
		customer.ID = utils.GenerateSonyflakeID()
	}
	customer.CreatedAt = now
	customer.CreatedBy = "1122334455"
	db := r.db.WithContext(ctx).Table(customer.TableName())
	err = db.Create(&customer).Error
	if err != nil {
		klog.CtxErrorf(ctx, "[Customer][Repo][SaveCustomer] Error saving customer to db", map[string]interface{}{
			"error":    err,
			"customer": customer,
		})
		return 0, err
	}
	return customer.ID, nil
}

func (r *CustomerRepository) FindOne(ctx context.Context, customerID *int64) (customer *domain.Customer, err error) {
	if customerID == nil {
		return
	}

	db := r.db.WithContext(ctx).Table(customer.TableName())
	db.Where("deleted_at IS NULL")
	db.Where("id = ?", customerID)

	err = db.First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		klog.CtxErrorf(ctx, "[Customer][Repo][FindOne] Error finding customer from db", map[string]interface{}{
			"error":      err,
			"customerID": customerID,
		})
		return
	}
	return customer, nil
}
