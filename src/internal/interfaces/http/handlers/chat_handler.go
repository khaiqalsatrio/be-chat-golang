package handlers

import (
	"chat-golang/src/internal/usecases/chat"
	"chat-golang/src/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	sendMessageUsecase *chat.SendMessageUsecase
	getMessagesUsecase *chat.GetMessagesUsecase
}

func NewChatHandler(send *chat.SendMessageUsecase, get *chat.GetMessagesUsecase) *ChatHandler {
	return &ChatHandler{
		sendMessageUsecase: send,
		getMessagesUsecase: get,
	}
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a message to a room
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param request body chat.SendMessageRequest true "Send Message Request"
// @Success 201 {object} response.Response{data=entities.Message}
// @Failure 400 {object} response.Response
// @Router /rooms/{roomId}/messages [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req chat.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.RoomID == "" && req.ConversationID == "" {
		req.RoomID = c.Param("roomId")
	}
	req.SenderID = c.GetString("user_id")

	msg, err := h.sendMessageUsecase.Execute(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Message sent", msg)
}

// GetMessages godoc
// @Summary Get messages for a room
// @Description Get message history for a specific room
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.Response{data=[]entities.Message}
// @Failure 500 {object} response.Response
// @Router /rooms/{roomId}/messages [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Param("roomId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.getMessagesUsecase.Execute(roomID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Messages fetched", messages)
}
