package handler

import (
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 { page = 1 }
	if pageSize <= 0 { pageSize = 10 }
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

func fmtTime(t *time.Time) string {
	if t == nil { return "" }
	return t.Format("2006-01-02")
}

func strVal(s *string) string {
	if s == nil { return "" }
	return *s
}

func mapToProtoShippingResponse(shipping interface{}) *pb.ShippingResponse {
	switch s := shipping.(type) {
	case *models.ShippingAddress:
		return &pb.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      fmtTime(s.CreatedAt),
			UpdatedAt:      fmtTime(s.UpdatedAt),
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		return &pb.ShippingResponse{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      strVal(s.CreatedAt),
			UpdatedAt:      strVal(s.UpdatedAt),
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}

func mapToProtoShippingResponseDeleteAt(shipping interface{}) *pb.ShippingResponseDeleteAt {
	switch s := shipping.(type) {
	case *models.ShippingAddress:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt != nil {
			deletedAt = wrapperspb.String(fmtTime(s.DeletedAt))
		}
		return &pb.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      fmtTime(s.CreatedAt),
			UpdatedAt:      fmtTime(s.UpdatedAt),
			DeletedAt:      deletedAt,
			Courier:        s.Courier,
		}
	case *repository.ShippingAddressResult:
		var deletedAt *wrapperspb.StringValue
		if s.DeletedAt != nil && *s.DeletedAt != "" {
			deletedAt = wrapperspb.String(*s.DeletedAt)
		}
		return &pb.ShippingResponseDeleteAt{
			Id:             s.ShippingAddressID,
			OrderId:        s.OrderID,
			Alamat:         s.Alamat,
			Provinsi:       s.Provinsi,
			Negara:         s.Negara,
			Kota:           s.Kota,
			ShippingMethod: s.ShippingMethod,
			ShippingCost:   int32(s.ShippingCost),
			CreatedAt:      strVal(s.CreatedAt),
			UpdatedAt:      strVal(s.UpdatedAt),
			DeletedAt:      deletedAt,
			Courier:        s.Courier,
		}
	default:
		return nil
	}
}
