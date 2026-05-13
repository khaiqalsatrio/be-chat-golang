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
	lastMessageSubquery := r.db.Select("content").Model(&entities.Message{}).
		Where("messages.room_id = rooms.id").
		Order("created_at desc").
		Limit(1)

	err := r.db.Model(&entities.Room{}).
		Select("rooms.*, (?) as last_message", lastMessageSubquery).
		Joins("JOIN room_participants ON room_participants.room_id = rooms.id").
		Where("room_participants.user_id = ?", userID).
		Preload("Participants").
		Order("rooms.updated_at desc").
		Find(&rooms).Error
	return rooms, err
}

func (r *postgresRoomRepository) UpdateUpdatedAt(roomID uuid.UUID) error {
	return r.db.Model(&entities.Room{}).
		Where("id = ?", roomID).
		UpdateColumn("updated_at", gorm.Expr("NOW()")).Error
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
