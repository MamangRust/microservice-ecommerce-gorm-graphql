package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_award/repository"
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

func mapToProtoMerchantAwardResponse(v *repository.MerchantCertResult) *pb.MerchantAwardResponse {
	return mapToProtoMerchantAwardResponseFromResult(v)
}

func mapToProtoMerchantAwardResponseDeleteAt(v *repository.MerchantCertResult) *pb.MerchantAwardResponseDeleteAt {
	return mapToProtoMerchantAwardResponseDeleteAtFromResult(v)
}


func getStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapToProtoMerchantAwardResponseFromResult(v *repository.MerchantCertResult) *pb.MerchantAwardResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantAwardResponse{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    getStringPtr(v.Description),
		IssuedBy:       getStringPtr(v.IssuedBy),
		CertificateUrl: getStringPtr(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantAwardResponseFromModel(v *models.MerchantCertificationsAndAward) *pb.MerchantAwardResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantAwardResponse{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    getStringPtr(v.Description),
		IssuedBy:       getStringPtr(v.IssuedBy),
		CertificateUrl: getStringPtr(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantAwardResponseDeleteAtFromResult(v *repository.MerchantCertResult) *pb.MerchantAwardResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantAwardResponseDeleteAt{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    getStringPtr(v.Description),
		IssuedBy:       getStringPtr(v.IssuedBy),
		CertificateUrl: getStringPtr(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantAwardResponseDeleteAtFromModel(v *models.MerchantCertificationsAndAward) *pb.MerchantAwardResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantAwardResponseDeleteAt{
		Id:             v.MerchantCertificationID,
		MerchantId:     v.MerchantID,
		Title:          v.Title,
		Description:    getStringPtr(v.Description),
		IssuedBy:       getStringPtr(v.IssuedBy),
		CertificateUrl: getStringPtr(v.CertificateUrl),
		IssueDate:      convert.FormatTimePtr(v.IssueDate),
		ExpiryDate:     convert.FormatTimePtr(v.ExpiryDate),
		CreatedAt:      convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:      convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}
