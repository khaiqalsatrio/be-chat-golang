package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"room_id"`
	SenderID  uuid.UUID      `gorm:"type:uuid;not null" json:"sender_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      string         `gorm:"default:'TEXT'" json:"type"` // TEXT, IMAGE, FILE
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Sender User `gorm:"foreignKey:SenderID" json:"sender"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
