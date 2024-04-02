package contracts

import "time"

// UserResponse represents the response body for user operations
type UserResponse struct {
	AuthToken string    `json:"auth_token"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommonResponse represents the response body for string responses
type CommonResponse struct {
	Message string `json:"message"`
}
