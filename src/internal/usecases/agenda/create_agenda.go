package agenda

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type CreateAgendaUsecase struct {
	agendaRepo repositories.AgendaRepository
}

func NewCreateAgendaUsecase(agendaRepo repositories.AgendaRepository) *CreateAgendaUsecase {
	return &CreateAgendaUsecase{agendaRepo: agendaRepo}
}

type CreateAgendaRequest struct {
	RoomID      string    `json:"room_id" binding:"required"`
	CreatorID   string    `json:"creator_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

func (u *CreateAgendaUsecase) Execute(req CreateAgendaRequest) error {
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		return err
	}

	creatorID, err := uuid.Parse(req.CreatorID)
	if err != nil {
		return err
	}

	agenda := &entities.Agenda{
		RoomID:      roomID,
		CreatorID:   creatorID,
		Title:       req.Title,
		Description: req.Description,
		ScheduledAt: req.ScheduledAt,
	}

	return u.agendaRepo.Create(agenda)
}
