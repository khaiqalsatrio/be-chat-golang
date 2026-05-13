package handlers

import (
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/usecases/auth"
	"chat-golang/src/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	registerUsecase         *auth.RegisterUsecase
	loginUsecase            *auth.LoginUsecase
	getMeUsecase            *auth.GetMeUsecase
	googleLoginUsecase      *auth.GoogleLoginUsecase
	uploadProfilePhotoUsecase *auth.UploadProfilePhotoUsecase
	deleteProfilePhotoUsecase *auth.DeleteProfilePhotoUsecase
	jwtService              *services.JWTService
	fileUploadService       *services.FileUploadService
}

func NewAuthHandler(
	registerUsecase *auth.RegisterUsecase,
	loginUsecase *auth.LoginUsecase,
	getMeUsecase *auth.GetMeUsecase,
	googleLoginUsecase *auth.GoogleLoginUsecase,
	uploadProfilePhotoUsecase *auth.UploadProfilePhotoUsecase,
	deleteProfilePhotoUsecase *auth.DeleteProfilePhotoUsecase,
	jwtService *services.JWTService,
	fileUploadService *services.FileUploadService,
) *AuthHandler {
	return &AuthHandler{
		registerUsecase:         registerUsecase,
		loginUsecase:            loginUsecase,
		getMeUsecase:            getMeUsecase,
		googleLoginUsecase:      googleLoginUsecase,
		uploadProfilePhotoUsecase: uploadProfilePhotoUsecase,
		deleteProfilePhotoUsecase: deleteProfilePhotoUsecase,
		jwtService:              jwtService,
		fileUploadService:       fileUploadService,
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

// Logout godoc
// @Summary Logout current user
// @Description Blacklist current JWT token and logout the current user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	token, err := h.jwtService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid token expiration")
		return
	}

	expiresAt := time.Unix(int64(exp), 0)
	h.jwtService.BlacklistToken(tokenString, expiresAt)

	response.JSON(c, http.StatusOK, "Logged out successfully", nil)
}

// UploadProfilePhoto godoc
// @Summary Upload or update profile photo
// @Description Upload a new profile photo and update user avatar URL
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Avatar image file"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entities.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/profile/photo [post]
func (h *AuthHandler) UploadProfilePhoto(c *gin.Context) {
	userID := c.GetString("user_id")

	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file field is required")
		return
	}

	file, err := header.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to open file")
		return
	}
	defer file.Close()

	uploadResult, err := h.fileUploadService.UploadFile(file, header.Filename)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	avatarURL := uploadResult.FileURL
	user, err := h.uploadProfilePhotoUsecase.Execute(userID, avatarURL)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Profile photo updated", user)
}

// DeleteProfilePhoto godoc
// @Summary Delete profile photo
// @Description Remove the current user profile photo
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/profile/photo [delete]
func (h *AuthHandler) DeleteProfilePhoto(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.deleteProfilePhotoUsecase.Execute(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Profile photo deleted", nil)
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
