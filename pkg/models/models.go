package models

import (
	"time"

	"github.com/google/uuid"
)

// APIKey struct representing the API key
type APIKey struct {
	ID        uuid.UUID  `json:"api_key_id"`
	Key       string     `json:"api_key"`
	TenantID  uuid.UUID  `json:"tenant_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Tenant represents the structure of a tenant in the system
type Tenant struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"Name"`
}

// CreateTenantRequest represents the request body for creating a new tenant
type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

// AuthResponse represents the response format for authentication status
type AuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}
