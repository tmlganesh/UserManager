package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/ganesh/ainyx/internal/models"
	"github.com/ganesh/ainyx/internal/service"
	"github.com/ganesh/ainyx/internal/validator"
)

// UserHandler handles HTTP requests for user operations.
// It delegates all business logic to the service layer.
type UserHandler struct {
	service   *service.UserService
	validator *validator.Validator
	logger    *zap.Logger
}

// NewUserHandler constructs a handler with its dependencies injected.
func NewUserHandler(svc *service.UserService, v *validator.Validator, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service:   svc,
		validator: v,
		logger:    logger,
	}
}

// RegisterRoutes sets up all user-related routes on the Fiber app.
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	users := app.Group("/users")
	users.Post("/", h.Create)
	users.Get("/", h.List)
	users.Get("/:id", h.Get)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}

// Create handles POST /users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if _, err := h.validator.ValidateCreateUser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	resp, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Get handles GET /users/:id
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	resp, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Update handles PUT /users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid request body",
		})
	}

	if _, err := h.validator.ValidateUpdateUser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	resp, err := h.service.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Delete handles DELETE /users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "invalid user id",
		})
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List handles GET /users with pagination support.
func (h *UserHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Clamp pagination values to sensible bounds.
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, total, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.PaginatedResponse{
		Data:  users,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}

// parseID extracts and validates the :id path parameter.
func parseID(c *fiber.Ctx) (int32, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
