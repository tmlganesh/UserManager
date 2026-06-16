package models

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required"`
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required"`
}

// UserResponse is the API response for a single user.
// Age is always calculated dynamically from Dob — never stored.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}
