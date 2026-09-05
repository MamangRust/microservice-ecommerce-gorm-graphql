package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/microservice-ecommerce-shared/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepository struct {
	queryClient   pb.UserQueryServiceClient
	commandClient pb.UserCommandServiceClient
}

func NewUserRepository(queryClient pb.UserQueryServiceClient, commandClient pb.UserCommandServiceClient) *userRepository {
	return &userRepository{
		queryClient:   queryClient,
		commandClient: commandClient,
	}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*AuthUser, error) {
	res, err := r.queryClient.FindById(ctx, &pb.FindByIdUserRequest{Id: int32(user_id)})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*AuthUser, error) {
	res, err := r.queryClient.FindByEmail(ctx, &pb.FindByEmailRequest{Email: email})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
	}, nil
}

func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*AuthUser, error) {
	res, err := r.queryClient.FindByEmail(ctx, &pb.FindByEmailRequest{Email: email})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
	}, nil
}

func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*AuthUser, error) {
	res, err := r.queryClient.FindByVerificationCode(ctx, &pb.FindByVerificationCodeRequest{VerificationCode: verification_code})
	if err != nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*AuthUser, error) {
	res, err := r.commandClient.Create(ctx, &pb.CreateUserRequest{
		Firstname:       request.FirstName,
		Lastname:        request.LastName,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
	})

	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*AuthUser, error) {
	res, err := r.commandClient.UpdateIsVerified(ctx, &pb.UpdateUserIsVerifiedRequest{
		Id:         int32(user_id),
		IsVerified: is_verified,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*AuthUser, error) {
	res, err := r.commandClient.UpdatePassword(ctx, &pb.UpdateUserPasswordRequest{
		Id:       int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	return &AuthUser{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
	}, nil
}
