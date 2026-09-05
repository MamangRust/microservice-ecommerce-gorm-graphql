package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_detail/repository"
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
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}


func mapToProtoMerchantDetailResponseFromResult(v *repository.MerchantDetailResult) *pb.MerchantDetailResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantDetailResponse{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      getString(v.DisplayName),
		CoverImageUrl:    getString(v.CoverImageUrl),
		LogoUrl:          getString(v.LogoUrl),
		ShortDescription: getString(v.ShortDescription),
		WebsiteUrl:       getString(v.WebsiteUrl),
		CreatedAt:        getString(v.CreatedAt),
		UpdatedAt:        getString(v.UpdatedAt),
	}
}

func mapToProtoMerchantDetailResponseFromModel(v *models.MerchantDetail) *pb.MerchantDetailResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantDetailResponse{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      getString(v.DisplayName),
		CoverImageUrl:    getString(v.CoverImageUrl),
		LogoUrl:          getString(v.LogoUrl),
		ShortDescription: getString(v.ShortDescription),
		WebsiteUrl:       getString(v.WebsiteUrl),
		CreatedAt:        convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt:        convert.FormatTimePtrMS(v.UpdatedAt),
	}
}

func mapToProtoMerchantDetailResponseDeleteAtFromResult(v *repository.MerchantDetailResult) *pb.MerchantDetailResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantDetailResponseDeleteAt{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      getString(v.DisplayName),
		CoverImageUrl:    getString(v.CoverImageUrl),
		LogoUrl:          getString(v.LogoUrl),
		ShortDescription: getString(v.ShortDescription),
		WebsiteUrl:       getString(v.WebsiteUrl),
		CreatedAt:        getString(v.CreatedAt),
		UpdatedAt:        getString(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantDetailResponseDeleteAtFromModel(v *models.MerchantDetail) *pb.MerchantDetailResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantDetailResponseDeleteAt{
		Id:               v.MerchantDetailID,
		MerchantId:       v.MerchantID,
		DisplayName:      getString(v.DisplayName),
		CoverImageUrl:    getString(v.CoverImageUrl),
		LogoUrl:          getString(v.LogoUrl),
		ShortDescription: getString(v.ShortDescription),
		WebsiteUrl:       getString(v.WebsiteUrl),
		CreatedAt:        convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt:        convert.FormatTimePtrMS(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantDetailResponse(m interface{}) *pb.MerchantDetailResponse {
	switch v := m.(type) {
	case *repository.MerchantDetailResult:
		return mapToProtoMerchantDetailResponseFromResult(v)
	case *models.MerchantDetail:
		return mapToProtoMerchantDetailResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantDetailResponseDeleteAt(m interface{}) *pb.MerchantDetailResponseDeleteAt {
	switch v := m.(type) {
	case *repository.MerchantDetailResult:
		return mapToProtoMerchantDetailResponseDeleteAtFromResult(v)
	case *models.MerchantDetail:
		return mapToProtoMerchantDetailResponseDeleteAtFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantSocialLinkResponse(m interface{}) *pb.MerchantSocialMediaLinkResponse {
	switch v := m.(type) {
	case *models.MerchantSocialMediaLink:
		return &pb.MerchantSocialMediaLinkResponse{
			Id:               v.MerchantSocialID,
			MerchantDetailId: v.MerchantDetailID,
			Platform:         v.Platform,
			Url:              v.Url,
			CreatedAt:        convert.FormatTimePtrMS(v.CreatedAt),
			UpdatedAt:        convert.FormatTimePtrMS(v.UpdatedAt),
		}
	default:
		return nil
	}
}
