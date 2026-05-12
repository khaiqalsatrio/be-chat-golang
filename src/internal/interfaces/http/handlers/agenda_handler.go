package handlers

import (
	"chat-golang/src/internal/usecases/agenda"
	"chat-golang/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgendaHandler struct {
	createAgendaUsecase *agenda.CreateAgendaUsecase
	getAgendasUsecase   *agenda.GetAgendasUsecase
}

func NewAgendaHandler(create *agenda.CreateAgendaUsecase, get *agenda.GetAgendasUsecase) *AgendaHandler {
	return &AgendaHandler{
		createAgendaUsecase: create,
		getAgendasUsecase:   get,
	}
}

// Create godoc
// @Summary Create a new agenda
// @Description Create a new agenda for a room
// @Tags agenda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param request body agenda.CreateAgendaRequest true "Create Agenda Request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /rooms/{roomId}/agendas [post]
func (h *AgendaHandler) Create(c *gin.Context) {
	var req agenda.CreateAgendaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set creator ID from context (from middleware)
	req.CreatorID = c.GetString("user_id")

	err := h.createAgendaUsecase.Execute(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Agenda created successfully", nil)
}

// GetByRoom godoc
// @Summary Get agendas for a room
// @Description Get all agendas for a specific room
// @Tags agenda
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Success 200 {object} response.Response{data=[]entities.Agenda}
// @Failure 500 {object} response.Response
// @Router /rooms/{roomId}/agendas [get]
func (h *AgendaHandler) GetByRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	agendas, err := h.getAgendasUsecase.Execute(roomID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Agendas fetched", agendas)
}
