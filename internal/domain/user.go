package domain

import (
	"context"
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type User struct {
	ID    string `json:"id" db:"id" bson:"_id"`
	Name  string `json:"name" db:"name" bson:"name"`
	Email string `json:"email" db:"email" bson:"email"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
}
