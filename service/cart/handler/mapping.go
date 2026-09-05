package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoCartResponseFromResult(v *repository.CartResult) *pb.CartResponse {
	return &pb.CartResponse{
		Id:        v.CartID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Price:     v.Price,
		Image:     v.Image,
		Quantity:  v.Quantity,
		Weight:    v.Weight,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func mapToProtoCartResponseFromCreate(v *repository.CartCreateResult) *pb.CartResponse {
	return &pb.CartResponse{
		Id:        v.CartID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Price:     v.Price,
		Image:     v.Image,
		Quantity:  v.Quantity,
		Weight:    v.Weight,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
