package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Agenda struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID      uuid.UUID      `gorm:"type:uuid;index;not null" json:"room_id"`
	CreatorID   uuid.UUID      `gorm:"type:uuid;not null" json:"creator_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ScheduledAt time.Time      `gorm:"not null" json:"scheduled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Room    Room `gorm:"foreignKey:RoomID" json:"-"`
	Creator User `gorm:"foreignKey:CreatorID" json:"creator"`
}

func (a *Agenda) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
