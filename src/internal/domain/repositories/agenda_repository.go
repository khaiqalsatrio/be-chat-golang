package repositories

import (
	"chat-golang/src/internal/domain/entities"

	"github.com/google/uuid"
)

type AgendaRepository interface {
	Create(agenda *entities.Agenda) error
	GetByID(id uuid.UUID) (*entities.Agenda, error)
	GetByRoomID(roomID uuid.UUID) ([]entities.Agenda, error)
	Update(agenda *entities.Agenda) error
	Delete(id uuid.UUID) error
}
