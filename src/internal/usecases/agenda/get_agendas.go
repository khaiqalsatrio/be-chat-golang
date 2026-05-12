package agenda

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetAgendasUsecase struct {
	agendaRepo repositories.AgendaRepository
}

func NewGetAgendasUsecase(agendaRepo repositories.AgendaRepository) *GetAgendasUsecase {
	return &GetAgendasUsecase{agendaRepo: agendaRepo}
}

func (u *GetAgendasUsecase) Execute(roomID string) ([]entities.Agenda, error) {
	id, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	return u.agendaRepo.GetByRoomID(id)
}
