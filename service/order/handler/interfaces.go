package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderQueryHandler interface {
	pb.OrderQueryServiceServer
}

type OrderCommandHandler interface {
	pb.OrderCommandServiceServer
}

type OrderHandleGrpc interface {
	FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error)
	FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error)

	FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error)

	Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error)
	Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error)
	TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error)
	RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error)
	DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error)
	RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error)
	DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error)
}
