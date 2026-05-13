package repositories

import (
	"chat-golang/src/internal/domain/entities"

	"github.com/google/uuid"
)

type StatusRepository interface {
	CreateStatus(status *entities.Status) error
	GetStatusByID(id uuid.UUID) (*entities.Status, error)
	GetActiveStatusesByUserID(userID uuid.UUID) ([]entities.Status, error)
	GetActiveStatusesByConcern(userID uuid.UUID) ([]entities.Status, error) // Status dari teman/kontak
	DeleteExpiredStatuses() error
	DeleteStatusByID(id uuid.UUID) error
	GetAllExpiredStatuses() ([]entities.Status, error) // Untuk cleanup
}
