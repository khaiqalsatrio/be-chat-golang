package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRoomRepository struct {
	db *gorm.DB
}

func NewPostgresRoomRepository(db *gorm.DB) repositories.RoomRepository {
	return &postgresRoomRepository{db: db}
}

func (r *postgresRoomRepository) Create(room *entities.Room) error {
	return r.db.Create(room).Error
}

func (r *postgresRoomRepository) GetByID(id uuid.UUID) (*entities.Room, error) {
	var room entities.Room
	err := r.db.Preload("Participants").First(&room, "id = ?", id).Error
	return &room, err
}

func (r *postgresRoomRepository) GetByUserID(userID uuid.UUID) ([]entities.Room, error) {
	var rooms []entities.Room
	// Simple join to get rooms where user is participant
	err := r.db.Joins("JOIN room_participants ON room_participants.room_id = rooms.id").
		Where("room_participants.user_id = ?", userID).
		Preload("Participants").
		Find(&rooms).Error
	return rooms, err
}

func (r *postgresRoomRepository) AddParticipant(roomID uuid.UUID, userID uuid.UUID) error {
	return r.db.Table("room_participants").Create(map[string]interface{}{
		"room_id": roomID,
		"user_id": userID,
	}).Error
}

func (r *postgresRoomRepository) RemoveParticipant(roomID uuid.UUID, userID uuid.UUID) error {
	return r.db.Table("room_participants").
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(nil).Error
}
