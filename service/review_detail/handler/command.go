package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-review-detail/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailCommandGrpc struct {
	pb.UnimplementedReviewDetailCommandServiceServer
	reviewDetailCommand service.ReviewDetailCommandService
	logger              logger.LoggerInterface
}

func NewReviewDetailCommandHandler(svc service.ReviewDetailCommandService, log logger.LoggerInterface) pb.ReviewDetailCommandServiceServer {
	return &reviewDetailCommandGrpc{
		reviewDetailCommand: svc,
		logger:              log,
	}
}

func (s *reviewDetailCommandGrpc) Create(ctx context.Context, request *pb.CreateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	req := &requests.CreateReviewDetailRequest{
		ReviewID: int(request.GetReviewId()),
		Type:     request.GetType(),
		Url:      request.GetUrl(),
		Caption:  request.GetCaption(),
	}
	if err := req.Validate(); err != nil {
		return nil, review_detail_errors.ErrGrpcValidateCreateReviewDetail
	}
	reviewDetail, err := s.reviewDetailCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully created review detail",
		Data:    mapToReviewDetailResponseFromModel(reviewDetail),
	}, nil
}

func (s *reviewDetailCommandGrpc) Update(ctx context.Context, request *pb.UpdateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetReviewDetailId())
	req := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &id,
		Type:           request.GetType(),
		Url:            request.GetUrl(),
		Caption:        request.GetCaption(),
	}
	if err := req.Validate(); err != nil {
		return nil, review_detail_errors.ErrGrpcValidateUpdateReviewDetail
	}
	reviewDetail, err := s.reviewDetailCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully updated review detail",
		Data:    mapToReviewDetailResponseFromModel(reviewDetail),
	}, nil
}

func (s *reviewDetailCommandGrpc) TrashedReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, review_detail_errors.ErrGrpcInvalidID
	}
	reviewDetail, err := s.reviewDetailCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully trashed review detail",
		Data:    mapToReviewDetailResponseDeleteAtFromModel(reviewDetail),
	}, nil
}

func (s *reviewDetailCommandGrpc) RestoreReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, review_detail_errors.ErrGrpcInvalidID
	}
	reviewDetail, err := s.reviewDetailCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully restored review detail",
		Data:    mapToReviewDetailResponseDeleteAtFromModel(reviewDetail),
	}, nil
}

func (s *reviewDetailCommandGrpc) DeleteReviewDetailPermanent(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, review_detail_errors.ErrGrpcInvalidID
	}
	_, err := s.reviewDetailCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewDelete{
		Status:  "success",
		Message: "Successfully deleted review detail permanently",
	}, nil
}

func (s *reviewDetailCommandGrpc) RestoreAllReviewDetail(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully restored all review details",
	}, nil
}

func (s *reviewDetailCommandGrpc) DeleteAllReviewDetailPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.reviewDetailCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}
	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully deleted all review details permanently",
	}, nil
}
