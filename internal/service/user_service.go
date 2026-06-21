package service

import (
	"context"

	"github.com/0xSangeet/user-api/internal/domain"
	"github.com/google/uuid"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) Create(ctx context.Context, name, email string) (*domain.User, error) {
	if name == "" || email == "" {
		return nil, domain.ErrInvalidInput
	}

	user := domain.User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
	}

	err := u.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.repo.GetByID(ctx, id)
}

func (u *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *UserService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return u.repo.Delete(ctx, id)
}
