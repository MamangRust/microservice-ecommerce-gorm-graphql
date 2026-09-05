package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-user/repository"
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


// mapToProtoUserResponse maps a *models.User to UserResponse
func mapToProtoUserResponse(m *models.User) *pb.UserResponse {
	if m == nil {
		return nil
	}
	return &pb.UserResponse{
		Id:        m.UserID,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Email:     m.Email,
		CreatedAt: convert.FormatTimeRFC3339(m.CreatedAt),
		UpdatedAt: convert.FormatTimeRFC3339(m.UpdatedAt),
	}
}

// mapToProtoUserResult maps a *repository.UserResult to UserResponse
func mapToProtoUserResult(r *repository.UserResult) *pb.UserResponse {
	if r == nil {
		return nil
	}
	return &pb.UserResponse{
		Id:        r.UserID,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Email:     r.Email,
		CreatedAt: convert.FormatTimeRFC3339(r.CreatedAt),
		UpdatedAt: convert.FormatTimeRFC3339(r.UpdatedAt),
	}
}

// mapToProtoUserResponseDeleteAt maps a *models.User to UserResponseDeleteAt
func mapToProtoUserResponseDeleteAt(m *models.User) *pb.UserResponseDeleteAt {
	if m == nil {
		return nil
	}
	res := &pb.UserResponseDeleteAt{
		Id:        m.UserID,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Email:     m.Email,
		CreatedAt: convert.FormatTimeRFC3339(m.CreatedAt),
		UpdatedAt: convert.FormatTimeRFC3339(m.UpdatedAt),
	}
	res.DeletedAt = convert.TimeToWrappers(m.DeletedAt)
	return res
}

// mapToProtoUserResultDeleteAt maps a *repository.UserResult to UserResponseDeleteAt
func mapToProtoUserResultDeleteAt(r *repository.UserResult) *pb.UserResponseDeleteAt {
	if r == nil {
		return nil
	}
	res := &pb.UserResponseDeleteAt{
		Id:        r.UserID,
		Firstname: r.Firstname,
		Lastname:  r.Lastname,
		Email:     r.Email,
		CreatedAt: convert.FormatTimeRFC3339(r.CreatedAt),
		UpdatedAt: convert.FormatTimeRFC3339(r.UpdatedAt),
	}
	res.DeletedAt = convert.TimeToWrappers(r.DeletedAt)
	return res
}

func mapToProtoUserResponseWithPassword(m *models.User) *pb.UserResponseWithPassword {
	if m == nil {
		return nil
	}
	return &pb.UserResponseWithPassword{
		Id:        m.UserID,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: convert.FormatTimeRFC3339(m.CreatedAt),
		UpdatedAt: convert.FormatTimeRFC3339(m.UpdatedAt),
	}
}
