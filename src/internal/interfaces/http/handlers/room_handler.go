package handlers

import (
	"chat-golang/src/internal/usecases/room"
	"chat-golang/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	createRoomUsecase *room.CreateRoomUsecase
	getRoomsUsecase   *room.GetRoomsUsecase
}

func NewRoomHandler(create *room.CreateRoomUsecase, get *room.GetRoomsUsecase) *RoomHandler {
	return &RoomHandler{
		createRoomUsecase: create,
		getRoomsUsecase:   get,
	}
}

// Create godoc
// @Summary Create a new room
// @Description Create a private or group chat room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body room.CreateRoomRequest true "Create Room Request"
// @Success 201 {object} response.Response{data=entities.Room}
// @Failure 400 {object} response.Response
// @Router /rooms [post]
func (h *RoomHandler) Create(c *gin.Context) {
	var req room.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	creatorID := c.GetString("user_id")
	res, err := h.createRoomUsecase.Execute(req, creatorID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Room created successfully", res)
}

// GetAll godoc
// @Summary Get my rooms
// @Description Get all chat rooms that the current user belongs to
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entities.Room}
// @Router /rooms [get]
func (h *RoomHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")
	res, err := h.getRoomsUsecase.Execute(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Rooms fetched", res)
}
