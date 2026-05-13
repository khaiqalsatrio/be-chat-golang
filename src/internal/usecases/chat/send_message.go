package chat

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"chat-golang/src/internal/interfaces/websocket"

	"github.com/google/uuid"
)

type SendMessageUsecase struct {
	chatRepo repositories.ChatRepository
	roomRepo repositories.RoomRepository
	hub      *websocket.Hub
}

func NewSendMessageUsecase(chatRepo repositories.ChatRepository, roomRepo repositories.RoomRepository, hub *websocket.Hub) *SendMessageUsecase {
	return &SendMessageUsecase{
		chatRepo: chatRepo,
		roomRepo: roomRepo,
		hub:      hub,
	}
}

type SendMessageRequest struct {
	RoomID   string `json:"room_id" form:"room_id"`
	ConversationID string `json:"conversation_id" form:"conversation_id"`
	SenderID string `json:"sender_id"`
	Content  string `json:"content" binding:"required"`
	Type     string `json:"type"`
}

func (u *SendMessageUsecase) Execute(req SendMessageRequest) (*entities.Message, error) {
	effectiveRoomID := req.RoomID
	if effectiveRoomID == "" {
		effectiveRoomID = req.ConversationID
	}

	roomID, err := uuid.Parse(effectiveRoomID)
	if err != nil {
		return nil, err
	}

	senderID, err := uuid.Parse(req.SenderID)
	if err != nil {
		return nil, err
	}

	message := &entities.Message{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  req.Content,
		Type:     req.Type,
	}

	err = u.chatRepo.SaveMessage(message)
	if err != nil {
		return nil, err
	}

	// Update room timestamp after successful message save
	if err = u.roomRepo.UpdateUpdatedAt(roomID); err != nil {
		return nil, err
	}

	// Re-fetch to get Sender info
	fullMessage, err := u.chatRepo.GetMessageByID(message.ID)
	if err == nil && u.hub != nil {
		// Broadcast message real-time
		u.hub.Broadcast <- fullMessage
	}

	return fullMessage, err
}
