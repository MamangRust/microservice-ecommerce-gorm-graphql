package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/repository"
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


// mapToProtoRoleResult maps a RoleResult to RoleResponse
func mapToProtoRoleResult(r *repository.RoleResult) *pb.RoleResponse {
	if r == nil {
		return nil
	}
	return &pb.RoleResponse{
		Id:        r.RoleID,
		Name:      r.RoleName,
		CreatedAt: convert.FormatTimePtrMS(r.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(r.UpdatedAt),
	}
}

// mapToProtoRoleResultDeleteAt maps a RoleResult to RoleResponseDeleteAt
func mapToProtoRoleResultDeleteAt(r *repository.RoleResult) *pb.RoleResponseDeleteAt {
	if r == nil {
		return nil
	}
	res := &pb.RoleResponseDeleteAt{
		Id:        r.RoleID,
		Name:      r.RoleName,
		CreatedAt: convert.FormatTimePtrMS(r.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(r.UpdatedAt),
	}
	if r.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(r.DeletedAt)
	}
	return res
}

// mapToProtoRoleResponse maps a *models.Role to RoleResponse
func mapToProtoRoleResponse(m *models.Role) *pb.RoleResponse {
	if m == nil {
		return nil
	}
	return &pb.RoleResponse{
		Id:        m.RoleID,
		Name:      m.RoleName,
		CreatedAt: convert.FormatTimePtrMS(m.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(m.UpdatedAt),
	}
}

// mapToProtoRoleResponseDeleteAt maps a *models.Role to RoleResponseDeleteAt
func mapToProtoRoleResponseDeleteAt(m *models.Role) *pb.RoleResponseDeleteAt {
	if m == nil {
		return nil
	}
	res := &pb.RoleResponseDeleteAt{
		Id:        m.RoleID,
		Name:      m.RoleName,
		CreatedAt: convert.FormatTimePtrMS(m.CreatedAt),
		UpdatedAt: convert.FormatTimePtrMS(m.UpdatedAt),
	}
	if m.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(m.DeletedAt)
	}
	return res
}

func mapToProtoUserRoleResponse(v *models.UserRole) *pb.UserRoleResponse {
	if v == nil {
		return nil
	}
	return &pb.UserRoleResponse{
		UserRoleId: v.UserRoleID,
		UserId:     v.UserID,
		RoleId:     v.RoleID,
		CreatedAt:  convert.FormatTimePtrMS(v.CreatedAt),
		UpdatedAt:  convert.FormatTimePtrMS(v.UpdatedAt),
	}
}
