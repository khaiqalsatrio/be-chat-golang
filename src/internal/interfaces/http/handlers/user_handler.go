package handlers

import (
	"chat-golang/src/internal/usecases/user"
	"chat-golang/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	getAllUsersUsecase *user.GetAllUsersUsecase
	getUserByIDUsecase *user.GetUserByIDUsecase
}

func NewUserHandler(getAll *user.GetAllUsersUsecase, getByID *user.GetUserByIDUsecase) *UserHandler {
	return &UserHandler{
		getAllUsersUsecase: getAll,
		getUserByIDUsecase: getByID,
	}
}

// GetAll godoc
// @Summary Get all users
// @Description Get a list of all registered users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]entities.User}
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.getAllUsersUsecase.Execute()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Users fetched", users)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get details of a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=entities.User}
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.getUserByIDUsecase.Execute(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.JSON(c, http.StatusOK, "User fetched", user)
}
