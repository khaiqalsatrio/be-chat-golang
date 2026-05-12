package repositories

import (
	"chat-golang/src/internal/domain/entities"

	"github.com/google/uuid"
)

type ChatRepository interface {
	SaveMessage(message *entities.Message) error
	GetMessagesByRoomID(roomID uuid.UUID, limit int, offset int) ([]entities.Message, error)
	GetMessageByID(id uuid.UUID) (*entities.Message, error)
	DeleteMessage(id uuid.UUID) error
}
