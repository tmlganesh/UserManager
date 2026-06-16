package repository

import (
	"context"

	"github.com/ganesh/ainyx/db/sqlc"
)

// UserRepository defines the contract for user persistence operations.
// Using an interface enables dependency injection and makes testing straightforward.
type UserRepository interface {
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetByID(ctx context.Context, id int32) (sqlc.User, error)
	List(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error)
	Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, id int32) error
	Count(ctx context.Context) (int64, error)
}

// userRepository is the concrete implementation backed by SQLC-generated queries.
type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a UserRepository backed by the given SQLC Queries.
func NewUserRepository(q *sqlc.Queries) UserRepository {
	return &userRepository{queries: q}
}

func (r *userRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, arg sqlc.ListUsersParams) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, arg)
}

func (r *userRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
