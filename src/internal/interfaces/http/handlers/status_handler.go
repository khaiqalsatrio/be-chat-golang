package handlers

import (
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/usecases/status"
	"chat-golang/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StatusHandler struct {
	createStatusUsecase   *status.CreateStatusUsecase
	getStatusesUsecase    *status.GetStatusesUsecase
	getMyStatusesUsecase  *status.GetMyStatusesUsecase
	deleteStatusUsecase   *status.DeleteStatusUsecase
	fileUploadService     *services.FileUploadService
}

func NewStatusHandler(
	createStatusUsecase *status.CreateStatusUsecase,
	getStatusesUsecase *status.GetStatusesUsecase,
	getMyStatusesUsecase *status.GetMyStatusesUsecase,
	deleteStatusUsecase *status.DeleteStatusUsecase,
	fileUploadService *services.FileUploadService,
) *StatusHandler {
	return &StatusHandler{
		createStatusUsecase:   createStatusUsecase,
		getStatusesUsecase:    getStatusesUsecase,
		getMyStatusesUsecase:  getMyStatusesUsecase,
		deleteStatusUsecase:   deleteStatusUsecase,
		fileUploadService:     fileUploadService,
	}
}

// CreateStatus godoc
// @Summary Create a new status (image/video)
// @Description Upload a new status with image or video media file
// @Tags status
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Media file (image or video)"
// @Param caption formData string false "Status caption"
// @Security BearerAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /status [post]
func (h *StatusHandler) CreateStatus(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		response.Error(c, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	// Get file from form
	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file field is required")
		return
	}

	// Open file reader
	file, err := header.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to read file")
		return
	}
	defer file.Close()

	// Upload file
	uploadResult, err := h.fileUploadService.UploadFile(file, header.Filename)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Determine media type from file extension
	mediaType := "image"
	filename := header.Filename
	if len(filename) > 0 {
		ext := filename[len(filename)-3:]
		if ext == "mp4" || ext == "avi" || ext == "mov" || ext == "mkv" || ext == "ebm" {
			mediaType = "video"
		}
	}

	// Get caption from form
	caption := c.PostForm("caption")

	// Create request
	req := status.CreateStatusRequest{
		MediaURL:  uploadResult.FileURL,
		MediaType: mediaType,
		Caption:   caption,
	}

	// Call usecase
	result, err := h.createStatusUsecase.Execute(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Status created successfully", result)
}

// GetStatuses godoc
// @Summary Get active statuses from contacts
// @Description Retrieve active statuses from connected users, grouped by user
// @Tags status
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /status [get]
func (h *StatusHandler) GetStatuses(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.getStatusesUsecase.Execute(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Statuses retrieved successfully", result)
}

// GetMyStatuses godoc
// @Summary Get my status history
// @Description Retrieve all active statuses of the current user
// @Tags status
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /status/me [get]
func (h *StatusHandler) GetMyStatuses(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.getMyStatusesUsecase.Execute(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "My statuses retrieved successfully", result)
}

// DeleteStatus godoc
// @Summary Delete a status
// @Description Delete a status before it expires (only owner can delete)
// @Tags status
// @Produce json
// @Param id path string true "Status ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /status/{id} [delete]
func (h *StatusHandler) DeleteStatus(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	statusID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid status id")
		return
	}

	err = h.deleteStatusUsecase.Execute(statusID, userID)
	if err != nil {
		if err.Error() == "status not found" {
			response.Error(c, http.StatusNotFound, err.Error())
		} else {
			response.Error(c, http.StatusForbidden, err.Error())
		}
		return
	}

	response.JSON(c, http.StatusOK, "Status deleted successfully", nil)
}
