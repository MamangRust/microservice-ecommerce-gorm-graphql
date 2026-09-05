package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-grpc-order/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pb.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pb.PaginationMeta{CurrentPage: int32(page), PageSize: int32(pageSize), TotalPages: int32(totalPages), TotalRecords: int32(totalRecords)}
}


func mapToProtoOrderResponseFromModel(v *models.Order) *pb.OrderResponse {
	if v == nil { return nil }
	return &pb.OrderResponse{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: convert.FormatTimePtrMS(v.CreatedAt), UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt)}
}

func mapToProtoOrderResponseFromResult(v *repository.OrderResult) *pb.OrderResponse {
	if v == nil { return nil }
	return &pb.OrderResponse{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: convert.FormatTimePtrMS(v.CreatedAt), UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt)}
}

func mapToProtoOrderResponseDeleteAtFromModel(v *models.Order) *pb.OrderResponseDeleteAt {
	if v == nil { return nil }
	res := &pb.OrderResponseDeleteAt{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: convert.FormatTimePtrMS(v.CreatedAt), UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt)}
	res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	return res
}

func mapToProtoOrderResponseDeleteAtFromResult(v *repository.OrderResult) *pb.OrderResponseDeleteAt {
	if v == nil { return nil }
	res := &pb.OrderResponseDeleteAt{Id: v.OrderID, MerchantId: v.MerchantID, UserId: v.UserID, TotalPrice: int32(v.TotalPrice), CreatedAt: convert.FormatTimePtrMS(v.CreatedAt), UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt)}
	res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	return res
}
