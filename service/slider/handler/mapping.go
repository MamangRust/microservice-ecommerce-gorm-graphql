package handler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)


func strVal(s *string) string {
	if s == nil { return "" }
	return *s
}

func mapToSliderResponse(slider interface{}) *pb.SliderResponse {
	switch v := slider.(type) {
	case *models.Slider:
		return &pb.SliderResponse{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.FormatDatePtr(v.CreatedAt), UpdatedAt: convert.FormatDatePtr(v.UpdatedAt)}
	case *repository.SliderResult:
		return &pb.SliderResponse{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: strVal(v.CreatedAt), UpdatedAt: strVal(v.UpdatedAt)}
	default:
		return nil
	}
}

func mapToSliderResponseDeleteAt(slider interface{}) *pb.SliderResponseDeleteAt {
	switch v := slider.(type) {
	case *models.Slider:
		res := &pb.SliderResponseDeleteAt{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: convert.FormatDatePtr(v.CreatedAt), UpdatedAt: convert.FormatDatePtr(v.UpdatedAt)}
		if v.DeletedAt != nil {
			res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
		}
		return res
	case *repository.SliderResult:
		res := &pb.SliderResponseDeleteAt{Id: v.SliderID, Name: v.Name, Image: v.Image, CreatedAt: strVal(v.CreatedAt), UpdatedAt: strVal(v.UpdatedAt)}
		if v.DeletedAt != nil && *v.DeletedAt != "" {
			res.DeletedAt = convert.StrValToWrappers(v.DeletedAt)
		}
		return res
	default:
		return nil
	}
}
