package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Birthday       string     `json:"birthday"`
	HashedPassword string     `json:"-"`
	Salt           string     `json:"-"`
	FullName       string     `json:"full_name"`
	Bio            string     `json:"bio"`
	AvatarURL      string     `json:"avatar_url"`
	Role           string     `json:"role"`
	IsActive       bool       `json:"is_active"`
	EmailVerified  bool       `json:"email_verified"`
	LastLogin      *time.Time `json:"last_login"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}