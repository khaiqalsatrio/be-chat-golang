package room

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetRoomsUsecase struct {
	roomRepo repositories.RoomRepository
}

func NewGetRoomsUsecase(roomRepo repositories.RoomRepository) *GetRoomsUsecase {
	return &GetRoomsUsecase{roomRepo: roomRepo}
}

func (u *GetRoomsUsecase) Execute(userIDStr string) ([]entities.Room, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return u.roomRepo.GetByUserID(userID)
}
