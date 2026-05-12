package room

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type CreateRoomUsecase struct {
	roomRepo repositories.RoomRepository
}

func NewCreateRoomUsecase(roomRepo repositories.RoomRepository) *CreateRoomUsecase {
	return &CreateRoomUsecase{roomRepo: roomRepo}
}

type CreateRoomRequest struct {
	Name         string   `json:"name"`
	Type         string   `json:"type" binding:"required"` // "private" or "group"
	Participants []string `json:"participants" binding:"required"`
}

func (u *CreateRoomUsecase) Execute(req CreateRoomRequest, creatorIDStr string) (*entities.Room, error) {
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return nil, err
	}

	room := &entities.Room{
		Name: req.Name,
		Type: req.Type,
	}

	err = u.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	// Add creator as participant
	_ = u.roomRepo.AddParticipant(room.ID, creatorID)

	// Add other participants
	for _, pIDStr := range req.Participants {
		pID, err := uuid.Parse(pIDStr)
		if err == nil && pID != creatorID {
			_ = u.roomRepo.AddParticipant(room.ID, pID)
		}
	}

	return u.roomRepo.GetByID(room.ID)
}
