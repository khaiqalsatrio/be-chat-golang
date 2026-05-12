package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresAgendaRepository struct {
	db *gorm.DB
}

func NewPostgresAgendaRepository(db *gorm.DB) repositories.AgendaRepository {
	return &postgresAgendaRepository{db: db}
}

func (r *postgresAgendaRepository) Create(agenda *entities.Agenda) error {
	return r.db.Create(agenda).Error
}

func (r *postgresAgendaRepository) GetByID(id uuid.UUID) (*entities.Agenda, error) {
	var agenda entities.Agenda
	err := r.db.Preload("Creator").First(&agenda, "id = ?", id).Error
	return &agenda, err
}

func (r *postgresAgendaRepository) GetByRoomID(roomID uuid.UUID) ([]entities.Agenda, error) {
	var agendas []entities.Agenda
	err := r.db.Where("room_id = ?", roomID).
		Order("scheduled_at asc").
		Preload("Creator").
		Find(&agendas).Error
	return agendas, err
}

func (r *postgresAgendaRepository) Update(agenda *entities.Agenda) error {
	return r.db.Save(agenda).Error
}

func (r *postgresAgendaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Agenda{}, "id = ?", id).Error
}
