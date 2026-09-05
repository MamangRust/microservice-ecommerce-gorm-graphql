package repository

import (
	"context"

	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
)

type merchantQueryRepository struct {
	client pb.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pb.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{client: client}
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, merchantID int) (*MerchantByIDResult, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdMerchantRequest{Id: int32(merchantID)})
	if err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return &MerchantByIDResult{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Address:      &res.Data.Address,
		ContactEmail: &res.Data.ContactEmail,
		ContactPhone: &res.Data.ContactPhone,
		Status:       res.Data.Status,
	}, nil
}
