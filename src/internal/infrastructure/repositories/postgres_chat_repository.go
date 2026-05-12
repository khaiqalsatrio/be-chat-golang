package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresChatRepository struct {
	db *gorm.DB
}

func NewPostgresChatRepository(db *gorm.DB) repositories.ChatRepository {
	return &postgresChatRepository{db: db}
}

func (r *postgresChatRepository) SaveMessage(message *entities.Message) error {
	return r.db.Create(message).Error
}

func (r *postgresChatRepository) GetMessagesByRoomID(roomID uuid.UUID, limit int, offset int) ([]entities.Message, error) {
	var messages []entities.Message
	err := r.db.Where("room_id = ?", roomID).
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Preload("Sender").
		Find(&messages).Error
	return messages, err
}

func (r *postgresChatRepository) GetMessageByID(id uuid.UUID) (*entities.Message, error) {
	var message entities.Message
	err := r.db.First(&message, "id = ?", id).Error
	return &message, err
}

func (r *postgresChatRepository) DeleteMessage(id uuid.UUID) error {
	return r.db.Delete(&entities.Message{}, "id = ?", id).Error
}
