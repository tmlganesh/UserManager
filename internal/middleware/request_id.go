package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	// HeaderRequestID is the HTTP header used for request tracing.
	HeaderRequestID = "X-Request-ID"

	// LocalRequestID is the Fiber locals key to access the request ID downstream.
	LocalRequestID = "requestID"
)

// RequestID middleware ensures every request has a unique trace identifier.
// If the client provides X-Request-ID, it is reused; otherwise a UUID is generated.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get(HeaderRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Store in Fiber locals for downstream access (logging, error context).
		c.Locals(LocalRequestID, reqID)

		// Echo the request ID back in the response header.
		c.Set(HeaderRequestID, reqID)

		return c.Next()
	}
}
