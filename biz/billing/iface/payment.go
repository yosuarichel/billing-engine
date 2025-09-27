package iface

import (
	"context"

	"github.com/yosuarichel/billing-engine/biz/payment/domain"
	paymentRepo "github.com/yosuarichel/billing-engine/biz/payment/repo"
	"gorm.io/gorm"
)

type IPaymentService interface {
	GetSumAmount(ctx context.Context, loanID *int64) (sumAmount *int64, err error)
	GetPaymentList(ctx context.Context, param *domain.FindAllParam) (payments []*domain.Payment, err error)
}

type IPaymentRepository interface {
	WithTx(tx *gorm.DB) *paymentRepo.PaymentRepository
	BulkSavePayment(ctx context.Context, payments []*domain.Payment) (ids []int64, err error)
}
