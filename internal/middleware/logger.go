package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger middleware measures request duration and logs structured request metadata.
// Fields logged: request_id, method, path, status_code, duration.
func Logger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process the request.
		err := c.Next()

		duration := time.Since(start)
		reqID, _ := c.Locals(LocalRequestID).(string)

		logger.Info("request completed",
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)

		return err
	}
}
