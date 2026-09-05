package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
)

type reviewDetailQueryGrpc struct {
	pb.UnimplementedReviewDetailQueryServiceServer
	reviewDetailQuery service.ReviewDetailQueryService
	logger            logger.LoggerInterface
}

func NewReviewDetailQueryHandler(svc service.ReviewDetailQueryService, log logger.LoggerInterface) pb.ReviewDetailQueryServiceServer {
	return &reviewDetailQueryGrpc{
		reviewDetailQuery: svc,
		logger:            log,
	}
}

func (h *reviewDetailQueryGrpc) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetails, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := h.reviewDetailQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponse, len(reviewDetails))
	for i, detail := range reviewDetails {
		protoReviewDetails[i] = mapToReviewDetailResponseFromResult(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationReviewDetails{
		Status:     "success",
		Message:    "Successfully fetched review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (h *reviewDetailQueryGrpc) FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, review_detail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := h.reviewDetailQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully fetched review detail",
		Data:    mapToReviewDetailResponseFromResult(reviewDetail),
	}, nil
}

func (h *reviewDetailQueryGrpc) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := h.reviewDetailQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, detail := range reviewDetails {
		protoReviewDetails[i] = mapToReviewDetailResponseDeleteAtFromResult(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (h *reviewDetailQueryGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := h.reviewDetailQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, detail := range reviewDetails {
		protoReviewDetails[i] = mapToReviewDetailResponseDeleteAtFromResult(detail)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}
