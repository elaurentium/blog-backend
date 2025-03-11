package entities

import (
	"time"

	"github.com/google/uuid"
)

type Sub struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rules       []string   `json:"rules"`
	CreatorID   uuid.UUID  `json:"creator_id"`
	IsPrivate   bool       `json:"is_private"`
	BannerURL   string     `json:"banner_url"`
	IconURL     string     `json:"icon_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}