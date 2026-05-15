package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index:idx_post_user,unique" json:"post_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_post_user,unique" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (l *Like) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return
}
