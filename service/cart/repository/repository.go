package repository

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"gorm.io/gorm"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *gorm.DB,
	userQuery pb.UserQueryServiceClient,
	productQuery pb.ProductQueryServiceClient,
) *Repositories {
	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    NewUserQueryRepository(userQuery),
		ProductQuery: NewProductQueryRepository(productQuery),
	}
}
