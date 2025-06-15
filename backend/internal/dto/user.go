package dto

import "github.com/google/uuid"

// UserCreateRequest for creating a new user (includes password).
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"` // Add min length for password
}

// UserUpdateRequest for updating user profile (excluding password).
type UserUpdateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	// Password updates would typically go through a separate endpoint like /users/{id}/password
	// If you include it here, ensure proper validation and hashing in the service.
}

// UserLoginRequest for user authentication.
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// FollowUserRequest for following another user.
type FollowUserRequest struct {
	FollowedUserID uuid.UUID `json:"followedUserId" validate:"required"`
}
