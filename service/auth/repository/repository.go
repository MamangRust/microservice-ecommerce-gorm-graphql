package repository

import (
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"gorm.io/gorm"
)

func NewRepositories(DB *gorm.DB,
	userQuery pb.UserQueryServiceClient,
	userCommand pb.UserCommandServiceClient,
	roleQuery pb.RoleQueryServiceClient,
	roleCommand pb.RoleCommandServiceClient,
) *Repositories {
	return &Repositories{
		User:         NewUserRepository(userQuery, userCommand),
		RefreshToken: NewRefreshTokenRepository(DB),
		UserRole:     NewUserRoleRepository(roleCommand),
		Role:         NewRoleRepository(roleQuery),
		ResetToken:   NewResetTokenRepository(DB),
	}
}
