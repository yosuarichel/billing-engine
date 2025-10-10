package service

import (
	"context"
	"errors"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/biz/payment/domain"
	paymentRepo "github.com/yosuarichel/billing-engine/biz/payment/repo"
)

type PaymentService struct {
	PaymentRepository *paymentRepo.PaymentRepository
}

func NewPaymentService(paymentRepo *paymentRepo.PaymentRepository) *PaymentService {
	return &PaymentService{
		PaymentRepository: paymentRepo,
	}
}

func (s *PaymentService) GetSumAmount(ctx context.Context, loanID *int64) (sumAmount *int64, err error) {
	if gptr.Indirect(loanID) <= 0 {
		return nil, errors.New("missing loanID")
	}

	paymentData, err := s.PaymentRepository.FindAll(ctx, &domain.FindAllParam{
		LoanID: gptr.Indirect(loanID),
	})
	if err != nil {
		klog.CtxErrorf(ctx, "[Payment][Service][GetSumAmount] Error call FindAll", map[string]interface{}{
			"error":  err,
			"loanID": loanID,
		})
		return
	}

	sumAmout := int64(0)
	for _, payment := range paymentData {
		sumAmout += payment.Amount
	}

	return gptr.Of(sumAmout), nil
}

func (s *PaymentService) GetPaymentList(ctx context.Context, param *domain.FindAllParam) (payments []*domain.Payment, err error) {
	if param == nil {
		return nil, errors.New("invalid param")
	}

	paymentData, err := s.PaymentRepository.FindAll(ctx, param)
	if err != nil {
		klog.CtxErrorf(ctx, "[Payment][Service][GetPaymentList] Error call FindAll", map[string]interface{}{
			"error": err,
			"param": param,
		})
		return
	}

	return paymentData, nil

}
