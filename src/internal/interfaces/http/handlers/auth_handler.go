package handlers

import (
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/usecases/auth"
	"chat-golang/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registerUsecase    *auth.RegisterUsecase
	loginUsecase       *auth.LoginUsecase
	getMeUsecase       *auth.GetMeUsecase
	googleLoginUsecase *auth.GoogleLoginUsecase
	jwtService         *services.JWTService
}

func NewAuthHandler(
	registerUsecase *auth.RegisterUsecase,
	loginUsecase *auth.LoginUsecase,
	getMeUsecase *auth.GetMeUsecase,
	googleLoginUsecase *auth.GoogleLoginUsecase,
	jwtService *services.JWTService,
) *AuthHandler {
	return &AuthHandler{
		registerUsecase:    registerUsecase,
		loginUsecase:       loginUsecase,
		getMeUsecase:       getMeUsecase,
		googleLoginUsecase: googleLoginUsecase,
		jwtService:         jwtService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} response.Response{data=auth.LoginResult}
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.registerUsecase.Execute(req.Username, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate token for auto-login
	token, err := h.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	result := auth.LoginResult{
		Token: token,
		User:  user,
	}

	response.JSON(c, http.StatusCreated, "User registered successfully", result)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} response.Response{data=auth.LoginResult}
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.loginUsecase.Execute(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Login successful", result)
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get current user profile using JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entities.User}
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.getMeUsecase.Execute(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "User profile fetched", user)
}

type GoogleLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

// GoogleLogin godoc
// @Summary Google Login
// @Description Authenticate user using Google ID Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body GoogleLoginRequest true "Google Login Request"
// @Success 200 {object} response.Response{data=auth.LoginResult}
// @Failure 401 {object} response.Response
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.googleLoginUsecase.Execute(req.IDToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Google login successful", result)
}
