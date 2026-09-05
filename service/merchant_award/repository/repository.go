package repository

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        MerchantQueryRepository
}

func NewRepositories(DB *gorm.DB, merchantQuery pb.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(DB),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(DB),
		MerchantQuery:        NewMerchantQueryRepository(merchantQuery),
	}
}
