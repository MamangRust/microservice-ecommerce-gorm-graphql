package handler

import (
	"encoding/json"

	"github.com/MamangRust/microservice-ecommerce-grpc-category/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)


func (h *Handler) mapToCategoryResponseFromModel(v *models.Category) *pb.CategoryResponse {
	if v == nil { return nil }
	return &pb.CategoryResponse{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatDatePtr(v.UpdatedAt),
	}
}

func (h *Handler) mapToCategoryResponseFromResult(v *repository.CategoryResult) *pb.CategoryResponse {
	if v == nil { return nil }
	return &pb.CategoryResponse{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatDatePtr(v.UpdatedAt),
	}
}

func (h *Handler) mapToDeleteAtResponseFromModel(v *models.Category) *pb.CategoryResponseDeleteAt {
	if v == nil { return nil }
	res := &pb.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatDatePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func (h *Handler) mapToDeleteAtResponseFromResult(v *repository.CategoryResult) *pb.CategoryResponseDeleteAt {
	if v == nil { return nil }
	res := &pb.CategoryResponseDeleteAt{
		Id:            v.CategoryID,
		Name:          v.Name,
		Description:   ptrStr(v.Description),
		SlugCategory:  ptrStr(v.SlugCategory),
		ImageCategory: ptrStr(v.ImageCategory),
		CreatedAt:     convert.FormatDatePtr(v.CreatedAt),
		UpdatedAt:     convert.FormatDatePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil { return "" }
	return string(jsonData)
}

func ptrStr(s *string) string {
	if s != nil { return *s }
	return ""
}
