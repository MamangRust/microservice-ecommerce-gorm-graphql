package repository

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"gorm.io/gorm"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

type categoryQueryRepository struct {
	client pb.CategoryQueryServiceClient
}

func NewCategoryQueryRepository(client pb.CategoryQueryServiceClient) *categoryQueryRepository {
	return &categoryQueryRepository{client: client}
}

type merchantQueryRepository struct {
	client pb.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pb.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{client: client}
}

func NewRepositories(db *gorm.DB, categoryClient pb.CategoryQueryServiceClient, merchantClient pb.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db, categoryClient),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  NewCategoryQueryRepository(categoryClient),
		MerchantQuery:  NewMerchantQueryRepository(merchantClient),
	}
}
