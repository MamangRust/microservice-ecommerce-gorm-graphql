package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)


func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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

func mapToReviewDetailResponseFromResult(v *repository.ReviewDetailResult) *pb.ReviewDetailsResponse {
	if v == nil {
		return nil
	}
	return &pb.ReviewDetailsResponse{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   getString(v.Caption),
		CreatedAt: getString(v.CreatedAt),
		UpdatedAt: getString(v.UpdatedAt),
	}
}

func mapToReviewDetailResponseFromModel(v *models.ReviewDetail) *pb.ReviewDetailsResponse {
	if v == nil {
		return nil
	}
	return &pb.ReviewDetailsResponse{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   getString(v.Caption),
		CreatedAt: convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt: convert.FormatDatePtr(v.UpdatedAt),
	}
}

func mapToReviewDetailResponseDeleteAtFromResult(v *repository.ReviewDetailResult) *pb.ReviewDetailsResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.ReviewDetailsResponseDeleteAt{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   getString(v.Caption),
		CreatedAt: getString(v.CreatedAt),
		UpdatedAt: getString(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}

func mapToReviewDetailResponseDeleteAtFromModel(v *models.ReviewDetail) *pb.ReviewDetailsResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.ReviewDetailsResponseDeleteAt{
		Id:        v.ReviewDetailID,
		ReviewId:  v.ReviewID,
		Type:      v.Type,
		Url:       v.Url,
		Caption:   getString(v.Caption),
		CreatedAt: convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt: convert.FormatDatePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}
