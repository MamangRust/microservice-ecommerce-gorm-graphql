package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-banner/repository"
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


func fmtTimeOnlyStr(t *string) string {
	if t == nil {
		return ""
	}
	return *t
}

func boolVal(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func mapToProtoBannerResponse(m interface{}) *pb.BannerResponse {
	switch v := m.(type) {
	case *models.Banner:
		return &pb.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: convert.FormatDatePtr(v.StartDate),
			EndDate:   convert.FormatDatePtr(v.EndDate),
			StartTime: fmtTimeOnlyStr(v.StartTime),
			EndTime:   fmtTimeOnlyStr(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
			UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
		}
	case *repository.BannerResult:
		return &pb.BannerResponse{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: strVal(v.StartDate),
			EndDate:   strVal(v.EndDate),
			StartTime: strVal(v.StartTime),
			EndTime:   strVal(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: strVal(v.CreatedAt),
			UpdatedAt: strVal(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapToProtoBannerResponseDeleteAt(m interface{}) *pb.BannerResponseDeleteAt {
	switch v := m.(type) {
	case *models.Banner:
		res := &pb.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: convert.FormatDatePtr(v.StartDate),
			EndDate:   convert.FormatDatePtr(v.EndDate),
			StartTime: fmtTimeOnlyStr(v.StartTime),
			EndTime:   fmtTimeOnlyStr(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: convert.FormatTimePtrMS(v.CreatedAt),
			UpdatedAt: convert.FormatTimePtrMS(v.UpdatedAt),
		}
		if v.DeletedAt != nil {
			res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
		}
		return res
	case *repository.BannerResult:
		res := &pb.BannerResponseDeleteAt{
			BannerId:  v.BannerID,
			Name:      v.Name,
			StartDate: strVal(v.StartDate),
			EndDate:   strVal(v.EndDate),
			StartTime: strVal(v.StartTime),
			EndTime:   strVal(v.EndTime),
			IsActive:  boolVal(v.IsActive),
			CreatedAt: strVal(v.CreatedAt),
			UpdatedAt: strVal(v.UpdatedAt),
		}
		if v.DeletedAt != nil && *v.DeletedAt != "" {
			res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
		}
		return res
	default:
		return nil
	}
}
