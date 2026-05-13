package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"-"`
	MediaURL  string         `gorm:"not null" json:"media_url"`
	MediaType string         `gorm:"default:'image'" json:"media_type"` // 'image' atau 'video'
	Caption   string         `json:"caption"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `gorm:"index" json:"expires_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Status) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	// Set expires_at to 24 hours from now if not already set
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = time.Now().Add(24 * time.Hour)
	}
	return
}

func (s *Status) IsActive() bool {
	return time.Now().Before(s.ExpiresAt)
}
