package chat

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetMessagesUsecase struct {
	chatRepo repositories.ChatRepository
}

func NewGetMessagesUsecase(chatRepo repositories.ChatRepository) *GetMessagesUsecase {
	return &GetMessagesUsecase{chatRepo: chatRepo}
}

func (u *GetMessagesUsecase) Execute(roomID string, limit, offset int) ([]entities.Message, error) {
	id, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	return u.chatRepo.GetMessagesByRoomID(id, limit, offset)
}
