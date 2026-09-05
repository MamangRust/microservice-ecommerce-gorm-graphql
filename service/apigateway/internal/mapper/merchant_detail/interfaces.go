package merchant_detailgraphqlmapper

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"github.com/MamangRust/microservice-ecommerce-grpc/service/apigateway/internal/model"
)

type MerchantDetailGraphqlMapper interface {
	ToGraphqlResponseMerchantDetailRelation(res *pb.ApiResponseMerchantDetail) *model.APIResponseMerchantDetailRelation
	ToGraphqlResponseMerchantDetailDelete(res *pb.ApiResponseMerchantDelete) *model.APIResponseMerchantDetailDelete
	ToGraphqlResponseMerchantDetailAll(res *pb.ApiResponseMerchantAll) *model.APIResponseMerchantDetailAll
	ToGraphqlResponseMerchantDetail(res *pb.ApiResponseMerchantDetail) *model.APIResponseMerchantDetail
	ToGraphqlResponseMerchantDetailDeleteAt(res *pb.ApiResponseMerchantDetailDeleteAt) *model.APIResponseMerchantDetailDeleteAt
	ToGraphqlResponsePaginationMerchantDetail(res *pb.ApiResponsePaginationMerchantDetail) *model.APIResponsePaginationMerchantDetail
	ToGraphqlResponsePaginationMerchantDetailDeleteAt(res *pb.ApiResponsePaginationMerchantDetailDeleteAt) *model.APIResponsePaginationMerchantDetailDeleteAt
}
