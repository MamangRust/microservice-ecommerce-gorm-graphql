package handler

import (
	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-shared/convert"
)


func mapToSingleResponseFromModel(data *models.MerchantPolicy) *pb.ApiResponseMerchantPolicies {
	if data == nil {
		return nil
	}
	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data: &pb.MerchantPoliciesResponse{
			Id:         data.MerchantPolicyID,
			MerchantId: data.MerchantID,
			PolicyType: data.PolicyType,
			Title:      data.Title,
			Description: data.Description,
			CreatedAt:  convert.FormatTimePtr(data.CreatedAt),
			UpdatedAt:  convert.FormatTimePtr(data.UpdatedAt),
		},
	}
}

func mapToSingleResponse(data *repository.MerchantPolicyResult) *pb.ApiResponseMerchantPolicies {
	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data:    mapToMerchantPolicyResponseFromResult(data),
	}
}

func mapToPaginationResponse(data []*repository.MerchantPolicyResult, total *int) *pb.ApiResponsePaginationMerchantPolicies {
	var policies []*pb.MerchantPoliciesResponse
	for _, v := range data {
		policies = append(policies, mapToMerchantPolicyResponseFromResult(v))
	}
	return &pb.ApiResponsePaginationMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pb.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToPaginationDeleteAtResponse(data []*repository.MerchantPolicyResult, total *int) *pb.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var policies []*pb.MerchantPoliciesResponseDeleteAt
	for _, item := range data {
		policies = append(policies, mapToMerchantPolicyResponseDeleteAtFromResult(item))
	}
	return &pb.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pb.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToSingleDeleteAtResponse(data *models.MerchantPolicy) *pb.ApiResponseMerchantPoliciesDeleteAt {
	return &pb.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully processed merchant policy",
		Data:    mapToMerchantPolicyResponseDeleteAtFromModel(data),
	}
}

func mapToMerchantPolicyResponseFromResult(v *repository.MerchantPolicyResult) *pb.MerchantPoliciesResponse {
	if v == nil {
		return nil
	}
	return &pb.MerchantPoliciesResponse{
		Id:           v.MerchantPolicyID,
		MerchantId:   v.MerchantID,
		PolicyType:   v.PolicyType,
		Title:        v.Title,
		Description:  v.Description,
		MerchantName: v.MerchantName,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
}

func mapToMerchantPolicyResponseDeleteAtFromResult(v *repository.MerchantPolicyResult) *pb.MerchantPoliciesResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantPoliciesResponseDeleteAt{
		Id:           v.MerchantPolicyID,
		MerchantId:   v.MerchantID,
		PolicyType:   v.PolicyType,
		Title:        v.Title,
		Description:  v.Description,
		MerchantName: v.MerchantName,
		CreatedAt:    convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:    convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}

func mapToMerchantPolicyResponseDeleteAtFromModel(v *models.MerchantPolicy) *pb.MerchantPoliciesResponseDeleteAt {
	if v == nil {
		return nil
	}
	res := &pb.MerchantPoliciesResponseDeleteAt{
		Id:         v.MerchantPolicyID,
		MerchantId: v.MerchantID,
		PolicyType: v.PolicyType,
		Title:      v.Title,
		Description: v.Description,
		CreatedAt:  convert.FormatTimePtr(v.CreatedAt),
		UpdatedAt:  convert.FormatTimePtr(v.UpdatedAt),
	}
	if v.DeletedAt != nil {
		res.DeletedAt = convert.TimeToWrappers(v.DeletedAt)
	}
	return res
}
