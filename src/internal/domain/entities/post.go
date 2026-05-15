package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user"`
	Caption       string         `gorm:"type:text" json:"caption"`
	ImageURL      string         `gorm:"type:varchar(255)" json:"image_url"`
	LikesCount    int            `gorm:"-" json:"likes_count"`
	CommentsCount int            `gorm:"-" json:"comments_count"`
	IsLiked       bool           `gorm:"-" json:"is_liked"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
