package db

import (
    "context"
    "github.com/boldyrevn/mod-example/internal/model"
)

//go:generate mockery --name DB
type DB interface {
    GetUser(ctx context.Context, id string) (model.User, error)
    UpdateUser(ctx context.Context, user model.User) error
    DeleteUser(ctx context.Context, id string) error
    CreateUser(ctx context.Context, user model.User) error
}
