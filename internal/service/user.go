package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/ganesh/ainyx/db/sqlc"
	"github.com/ganesh/ainyx/internal/models"
	"github.com/ganesh/ainyx/internal/repository"
	"github.com/ganesh/ainyx/internal/utils"
)

// dobLayout is the canonical date format used across the application.
const dobLayout = "2006-01-02"

// ErrUserNotFound signals that the requested user does not exist.
var ErrUserNotFound = fmt.Errorf("user not found")

// UserService encapsulates all business logic for user operations.
type UserService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

// NewUserService constructs a UserService with its dependencies injected.
func NewUserService(repo repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// Create persists a new user and returns the enriched response with calculated age.
func (s *UserService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	user, err := s.repo.Create(ctx, sqlc.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user")
	}

	s.logger.Info("user created", zap.Int32("user_id", user.ID), zap.String("name", user.Name))
	return toUserResponse(user), nil
}

// GetByID retrieves a single user by primary key.
func (s *UserService) GetByID(ctx context.Context, id int32) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		s.logger.Error("failed to get user", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get user")
	}

	return toUserResponse(user), nil
}

// Update modifies an existing user's name and dob.
func (s *UserService) Update(ctx context.Context, id int32, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	user, err := s.repo.Update(ctx, sqlc.UpdateUserParams{
		Name: req.Name,
		Dob:  dob,
		ID:   id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		s.logger.Error("failed to update user", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to update user")
	}

	s.logger.Info("user updated", zap.Int32("user_id", user.ID), zap.String("name", user.Name))
	return toUserResponse(user), nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id int32) error {
	// Verify the user exists before attempting deletion.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		s.logger.Error("failed to check user before delete", zap.Int32("id", id), zap.Error(err))
		return fmt.Errorf("failed to delete user")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Int32("id", id), zap.Error(err))
		return fmt.Errorf("failed to delete user")
	}

	s.logger.Info("user deleted", zap.Int32("user_id", id))
	return nil
}

// List returns a paginated slice of users with total count.
func (s *UserService) List(ctx context.Context, page, limit int) ([]models.UserResponse, int64, error) {
	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list users")
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error("failed to count users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count users")
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toUserResponse(u)
	}

	return responses, total, nil
}

// toUserResponse converts a database user row to an API response,
// dynamically calculating age from dob.
func toUserResponse(user sqlc.User) *models.UserResponse {
	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format(dobLayout),
		Age:  utils.CalculateAge(user.Dob),
	}
}
