package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	UserID      uuid.UUID  `json:"user_id"`
	SubID 		uuid.UUID  `json:"sub_id"`
	Upvotes     int        `json:"upvotes"`
	Downvotes   int        `json:"downvotes"`
	IsLocked    bool       `json:"is_locked"`
	IsPinned    bool       `json:"is_pinned"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}