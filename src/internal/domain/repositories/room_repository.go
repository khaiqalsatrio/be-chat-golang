package repositories

import (
	"chat-golang/src/internal/domain/entities"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(room *entities.Room) error
	GetByID(id uuid.UUID) (*entities.Room, error)
	GetByUserID(userID uuid.UUID) ([]entities.Room, error)
	UpdateUpdatedAt(roomID uuid.UUID) error
	AddParticipant(roomID uuid.UUID, userID uuid.UUID) error
	RemoveParticipant(roomID uuid.UUID, userID uuid.UUID) error
}
