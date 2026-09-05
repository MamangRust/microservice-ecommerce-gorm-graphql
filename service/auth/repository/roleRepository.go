package repository

import (
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-shared/pb"
)

type roleRepository struct {
	client pb.RoleQueryServiceClient
}

func NewRoleRepository(client pb.RoleQueryServiceClient) RoleRepository {
	return &roleRepository{
		client: client,
	}
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*AuthRole, error) {
	res, err := r.client.FindByIdRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}

	return &AuthRole{
		RoleID:   res.Data.Id,
		RoleName: res.Data.Name,
	}, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*AuthRole, error) {
	res, err := r.client.FindByNameRole(ctx, &pb.FindByNameRoleRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by name %s: %w", name, err)
	}

	return &AuthRole{
		RoleID:   res.Data.Id,
		RoleName: res.Data.Name,
	}, nil
}
