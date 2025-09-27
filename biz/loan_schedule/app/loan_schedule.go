package app

// import (
// 	"context"
// 	"net/http"

// 	"github.com/bytedance/gg/gconv"
// 	"github.com/bytedance/gg/gptr"
// 	"github.com/cloudwego/kitex/pkg/klog"
// 	"github.com/yosuarichel/billing-engine/biz/customer/domain"
// 	"github.com/yosuarichel/billing-engine/kitex_gen/base"
// 	"github.com/yosuarichel/billing-engine/kitex_gen/billing_engine"
// )

// func (a *CustomerApp) CreateCustomer(ctx context.Context, req *billing_engine.CreateCustomerRequest) (res *billing_engine.CreateCustomerResponse) {
// 	customerID, err := a.CustomerService.CreateCustomer(ctx, &domain.Customer{
// 		Name:        req.GetName(),
// 		PhoneNumber: gptr.Of(req.GetPhoneNumber()),
// 	})
// 	if err != nil {
// 		klog.CtxErrorf(ctx, "[Customer][App][CreateCustomer] Error call CreateCustomer service", map[string]interface{}{
// 			"error":  err.Error(),
// 			"params": req,
// 		})
// 		return &billing_engine.CreateCustomerResponse{
// 			BaseResp: &base.BaseResp{
// 				StatusMessage: err.Error(),
// 				StatusCode:    http.StatusInternalServerError,
// 			},
// 		}
// 	}

// 	res = &billing_engine.CreateCustomerResponse{
// 		CustomerId: gconv.To[string](customerID),
// 		BaseResp: &base.BaseResp{
// 			StatusMessage: "success",
// 			StatusCode:    http.StatusOK,
// 		},
// 	}
// 	return
// }

// // func (a *ProductApp) GetProductDetail(ctx context.Context, req *billing_engine.GetProductDetailRequest) (res *billing_engine.GetProductDetailResponse) {
// // 	customerID := gconv.To[int64](req.GetProductId())
// // 	productData, err := a.CustomerService.GetProductDetail(ctx, gptr.Of(customerID))
// // 	if err != nil {
// // 		klog.CtxErrorf(ctx, "[Product][App][GetProductDetail] Error call GetProductDetail service", map[string]interface{}{
// // 			"error":  err.Error(),
// // 			"params": req,
// // 		})
// // 		return &billing_engine.GetProductDetailResponse{
// // 			BaseResp: &base.BaseResp{
// // 				StatusMessage: err.Error(),
// // 				StatusCode:    http.StatusInternalServerError,
// // 			},
// // 		}
// // 	}

// // 	res = &billing_engine.GetProductDetailResponse{
// // 		BaseResp: &base.BaseResp{
// // 			StatusMessage: "success",
// // 			StatusCode:    http.StatusOK,
// // 		},
// // 	}

// // 	if productData != nil {
// // 		res.Data = &product_data.ProductData{
// // 			ProductId:   gconv.To[string](productData.ID),
// // 			Name:        productData.Name,
// // 			Description: productData.Description,
// // 		}
// // 	}
// // 	return
// // }

// // func (a *ProductApp) GetProductList(ctx context.Context, req *billing_engine.GetProductListRequest) (res *billing_engine.GetProductListResponse) {
// // 	klog.CtxInfof(ctx, "[App GetProductList]")

// // 	var customerIDs []int64

// // 	if req.GetProductIds() != nil {
// // 		customerIDs = gslice.Map(req.GetProductIds(), func(id string) int64 {
// // 			return gconv.To[int64](id)
// // 		})
// // 	}

// // 	filter := &domain.GetProductListParam{}
// // 	if len(customerIDs) > 0 {
// // 		filter.ProductIDs = customerIDs
// // 	}
// // 	if req.GetProductName() != "" {
// // 		filter.Name = gptr.Of(req.GetProductName())
// // 	}

// // 	productData, err := a.CustomerService.GetProductList(ctx, filter)
// // 	if err != nil {
// // 		klog.CtxErrorf(ctx, "[Product][App][GetProductList] Error call GetProductList service", map[string]interface{}{
// // 			"error":  err.Error(),
// // 			"params": req,
// // 		})
// // 		return &billing_engine.GetProductListResponse{
// // 			BaseResp: &base.BaseResp{
// // 				StatusMessage: err.Error(),
// // 				StatusCode:    http.StatusInternalServerError,
// // 			},
// // 		}
// // 	}

// // 	res = &billing_engine.GetProductListResponse{
// // 		BaseResp: &base.BaseResp{
// // 			StatusMessage: "success",
// // 			StatusCode:    http.StatusOK,
// // 		},
// // 	}

// // 	var pageNum, pageSize int64
// // 	if req.GetPagination() != nil {
// // 		pageNum = int64(req.GetPagination().GetPageNum())
// // 		pageSize = int64(req.GetPagination().GetPageSize())
// // 	}

// // 	// Paginate the product list
// // 	paginatedProducts, total := utils.Paginate(productData, pageNum, pageSize)

// // 	res.Pagination = &common.PageNumberPagination{
// // 		PageNum:  int32(pageNum),
// // 		PageSize: int32(pageSize),
// // 		Total:    int32(total),
// // 	}

// // 	// Map paginated products to ProductData
// // 	if len(paginatedProducts) > 0 {
// // 		res.Data = gslice.Map(paginatedProducts, func(product *domain.Product) *product_data.ProductData {
// // 			return &product_data.ProductData{
// // 				ProductId:   gconv.To[string](product.ID),
// // 				Name:        product.Name,
// // 				Description: product.Description,
// // 			}
// // 		})
// // 	}

// // 	return res
// // }
