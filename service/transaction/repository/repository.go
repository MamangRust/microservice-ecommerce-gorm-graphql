package repository

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"gorm.io/gorm"
)

type Repositories struct {
	TransactionCommand TransactionCommandRepository
	TransactionQuery   TransactionQueryRepository
	OrderItem          OrderItemRepository
	OrderQuery         OrderQueryRepository
	MerchantQuery      MerchantQueryRepository
	ShippingAddress    ShippingAddressQueryRepository
	UserQuery          UserQueryRepository
	Outbox             OutboxRepository
}

type Deps struct {
	DB             *gorm.DB
	UserQuery      pb.UserQueryServiceClient
	MerchantQuery  pb.MerchantQueryServiceClient
	OrderQuery     pb.OrderQueryServiceClient
	OrderItemQuery pb.OrderItemQueryServiceClient
	ShippingQuery  pb.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.DB),
		TransactionQuery:   NewTransactionQueryRepository(deps.DB),
		OrderItem:          NewOrderItemRepository(deps.OrderItemQuery),
		OrderQuery:         NewOrderQueryRepository(deps.OrderQuery),
		MerchantQuery:      NewMerchantQueryRepository(deps.MerchantQuery),
		ShippingAddress:    NewShippingAddressQueryRepository(deps.ShippingQuery),
		UserQuery:          NewUserQueryRepository(deps.UserQuery),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
