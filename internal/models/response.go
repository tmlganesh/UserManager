package models

// ErrorResponse is the standard error payload returned by all endpoints.
type ErrorResponse struct {
	Error string `json:"error"`
}

// PaginatedResponse wraps a list of items with pagination metadata.
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
}
