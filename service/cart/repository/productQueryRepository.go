package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
)

type productQueryRepository struct {
	client pb.ProductQueryServiceClient
}

func NewProductQueryRepository(client pb.ProductQueryServiceClient) ProductQueryRepository {
	return &productQueryRepository{client: client}
}

func (r *productQueryRepository) FindById(ctx context.Context, id int) (*ProductResult, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}

	return &ProductResult{
		ProductID:    res.Data.Id,
		Name:         res.Data.Name,
		Price:        res.Data.Price,
		CountInStock: res.Data.CountInStock,
		ImageProduct: &res.Data.ImageProduct,
		Weight:       &res.Data.Weight,
	}, nil
}
