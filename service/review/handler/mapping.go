package handler

import (
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-review/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)


func parseOptionalTime(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func (h *reviewHandleGrpc) mapResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *models.Review:
		return h.mapReview(v)
	case *repository.ReviewResult:
		return h.mapReviewResult(v)
	case []*repository.ReviewResult:
		if len(v) > 0 && v[0].ReviewDetails != nil {
			res := make([]*pb.ReviewsDetailResponse, len(v))
			for i, r := range v {
				res[i] = h.mapReviewResultDetail(r)
			}
			return res
		}
		// Check if we should return ReviewResponse or ReviewResponseDeleteAt
		// by examining whether the items have deleted_at
		if len(v) > 0 && v[0].DeletedAt != nil && *v[0].DeletedAt != "" {
			res := make([]*pb.ReviewResponseDeleteAt, len(v))
			for i, r := range v {
				res[i] = h.mapReviewResultDeleteAt(r)
			}
			return res
		}
		res := make([]*pb.ReviewResponse, len(v))
		for i, r := range v {
			res[i] = h.mapReviewResultResponse(r)
		}
		return res
	default:
		_ = fmt.Sprintf("%T", data)
		return nil
	}
}

func (h *reviewHandleGrpc) mapReview(v *models.Review) *pb.ReviewResponseDeleteAt {
	res := &pb.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt: convert.FormatDatePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func (h *reviewHandleGrpc) mapReviewResult(v *repository.ReviewResult) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: parseOptionalTime(v.CreatedAt),
		UpdatedAt: parseOptionalTime(v.UpdatedAt),
	}
}

func (h *reviewHandleGrpc) mapReviewResultResponse(v *repository.ReviewResult) *pb.ReviewResponse {
	return &pb.ReviewResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: parseOptionalTime(v.CreatedAt),
		UpdatedAt: parseOptionalTime(v.UpdatedAt),
	}
}

func (h *reviewHandleGrpc) mapReviewResultDeleteAt(v *repository.ReviewResult) *pb.ReviewResponseDeleteAt {
	res := &pb.ReviewResponseDeleteAt{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: parseOptionalTime(v.CreatedAt),
		UpdatedAt: parseOptionalTime(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}

func (h *reviewHandleGrpc) mapReviewResultDetail(v *repository.ReviewResult) *pb.ReviewsDetailResponse {
	res := &pb.ReviewsDetailResponse{
		Id:        v.ReviewID,
		UserId:    v.UserID,
		ProductId: v.ProductID,
		Name:      v.Name,
		Comment:   v.Comment,
		Rating:    v.Rating,
		CreatedAt: parseOptionalTime(v.CreatedAt),
		UpdatedAt: parseOptionalTime(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = *v.DeletedAt
	}
	return res
}
