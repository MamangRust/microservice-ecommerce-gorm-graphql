package handler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pb.PaginationMeta {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}


func mapToProtoOrderItemResponseFromModel(v *models.OrderItem) *pb.OrderItemResponse {
	if v == nil {
		return nil
	}
	return &pb.OrderItemResponse{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
	}
}

func mapToProtoOrderItemResponseFromResult(v *repository.OrderItemResult) *pb.OrderItemResponse {
	if v == nil {
		return nil
	}
	return &pb.OrderItemResponse{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
	}
}

func mapToProtoOrderItemResponseDeleteAtFromModel(v *models.OrderItem) *pb.OrderItemResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.OrderItemResponseDeleteAt{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoOrderItemResponseDeleteAtFromResult(v *repository.OrderItemResult) *pb.OrderItemResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.OrderItemResponseDeleteAt{
		Id:        v.OrderItemID,
		OrderId:   v.OrderID,
		ProductId: v.ProductID,
		Quantity:  v.Quantity,
		Price:     v.Price,
		CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

