package handler

import (
	"math"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func mapToProtoMerchantResponseFromResult(v *repository.MerchantResult) *pb.MerchantResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantResponse{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  getString(v.Description),
		Address:      getString(v.Address),
		ContactEmail: getString(v.ContactEmail),
		ContactPhone: getString(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    getString(v.CreatedAt),
		UpdatedAt:    getString(v.UpdatedAt),
	}
}

func mapToProtoMerchantResponseFromModel(v *models.Merchant) *pb.MerchantResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantResponse{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  getString(v.Description),
		Address:      getString(v.Address),
		ContactEmail: getString(v.ContactEmail),
		ContactPhone: getString(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantResponse(m interface{}) *pb.MerchantResponse {
	switch v := m.(type) {
	case *repository.MerchantResult:
		return mapToProtoMerchantResponseFromResult(v)
	case *models.Merchant:
		return mapToProtoMerchantResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantResponseDeleteAtFromResult(v *repository.MerchantResult) *pb.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  getString(v.Description),
		Address:      getString(v.Address),
		ContactEmail: getString(v.ContactEmail),
		ContactPhone: getString(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    getString(v.CreatedAt),
		UpdatedAt:    getString(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
	}
	return res
}

func mapToProtoMerchantResponseDeleteAtFromModel(v *models.Merchant) *pb.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  getString(v.Description),
		Address:      getString(v.Address),
		ContactEmail: getString(v.ContactEmail),
		ContactPhone: getString(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToProtoMerchantResponseDeleteAt(m interface{}) *pb.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *repository.MerchantResult:
		return mapToProtoMerchantResponseDeleteAtFromResult(v)
	case *models.Merchant:
		return mapToProtoMerchantResponseDeleteAtFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantResponseTrashed(v *repository.MerchantResult) *pb.MerchantResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantResponseDeleteAt{
		Id:           v.MerchantID,
		UserId:       v.UserID,
		Name:         v.Name,
		Description:  getString(v.Description),
		Address:      getString(v.Address),
		ContactEmail: getString(v.ContactEmail),
		ContactPhone: getString(v.ContactPhone),
		Status:       v.Status,
		CreatedAt:    getString(v.CreatedAt),
		UpdatedAt:    getString(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
	}
	return res
}

func mapToProtoMerchantDocumentResponseFromResult(v *repository.MerchantDocumentResult) *pb.MerchantDocument {
	if v == nil {
		return nil
	}
	return &pb.MerchantDocument{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         getString(v.Note),
		UploadedAt:   getString(v.UploadedAt),
		UpdatedAt:    getString(v.UpdatedAt),
	}
}

func mapToProtoMerchantDocumentResponseFromModel(v *models.MerchantDocument) *pb.MerchantDocument {
	if v == nil {
		return nil
	}
	return &pb.MerchantDocument{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         getString(v.Note),
		UploadedAt:   convert.FormatTimePtr(v.UploadedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToProtoMerchantDocumentResponse(m interface{}) *pb.MerchantDocument {
	switch v := m.(type) {
	case *repository.MerchantDocumentResult:
		return mapToProtoMerchantDocumentResponseFromResult(v)
	case *models.MerchantDocument:
		return mapToProtoMerchantDocumentResponseFromModel(v)
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponseAt(v *repository.MerchantDocumentResult) *pb.MerchantDocumentDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantDocumentDeleteAt{
		DocumentId:   v.DocumentID,
		MerchantId:   v.MerchantID,
		DocumentType: v.DocumentType,
		DocumentUrl:  v.DocumentUrl,
		Status:       v.Status,
		Note:         getString(v.Note),
		UploadedAt:   getString(v.UploadedAt),
		UpdatedAt:    getString(v.UpdatedAt),
	}
	if v.DeletedAt != nil && *v.DeletedAt != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: *v.DeletedAt}
	}
	return res
}
