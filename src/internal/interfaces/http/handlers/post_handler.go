package handlers

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/internal/usecases/post"
	"chat-golang/src/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postUsecase       post.PostUsecase
	fileUploadService *services.FileUploadService
}

func NewPostHandler(postUsecase post.PostUsecase, fileUploadService *services.FileUploadService) *PostHandler {
	return &PostHandler{
		postUsecase:       postUsecase,
		fileUploadService: fileUploadService,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with optional image upload
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param image formData file false "Post image"
// @Param caption formData string false "Post caption"
// @Security BearerAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDRaw.(string))

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	caption := c.PostForm("caption")
	var imageURL string

	header, err := c.FormFile("image")
	if err == nil {
		file, _ := header.Open()
		defer file.Close()
		uploadResult, uploadErr := h.fileUploadService.UploadFile(file, header.Filename)
		if uploadErr != nil {
			response.Error(c, http.StatusBadRequest, uploadErr.Error())
			return
		}
		imageURL = uploadResult.FileURL
	}

	newPost := &entities.Post{
		UserID:   userID,
		Caption:  caption,
		ImageURL: imageURL,
	}

	if err := h.postUsecase.CreatePost(newPost); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Post created successfully", newPost)
}

// GetGlobalFeed godoc
// @Summary Get global feed
// @Description Retrieve latest posts with pagination
// @Tags posts
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts [get]
func (h *PostHandler) GetGlobalFeed(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	var currentUserID uuid.UUID
	if exists {
		currentUserID, _ = uuid.Parse(userIDRaw.(string))
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	posts, err := h.postUsecase.GetGlobalFeed(limit, offset, currentUserID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Feed retrieved successfully", posts)
}

// ToggleLike godoc
// @Summary Toggle like on a post
// @Description Like or unlike a post
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts/{id}/like [post]
func (h *PostHandler) ToggleLike(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDRaw.(string))

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	isLiked, err := h.postUsecase.ToggleLike(postID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	msg := "Post liked"
	if !isLiked {
		msg = "Post unliked"
	}

	response.JSON(c, http.StatusOK, msg, gin.H{"is_liked": isLiked})
}

// AddComment godoc
// @Summary Add a comment to a post
// @Description Add a new comment to the specified post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param input body object true "Comment content"
// @Security BearerAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts/{id}/comments [post]
func (h *PostHandler) AddComment(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDRaw.(string))

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	comment := &entities.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: input.Content,
	}

	if err := h.postUsecase.AddComment(comment); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, "Comment added successfully", comment)
}

// GetComments godoc
// @Summary Get comments for a post
// @Description Retrieve all comments for the specified post
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts/{id}/comments [get]
func (h *PostHandler) GetComments(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	comments, err := h.postUsecase.GetComments(postID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(c, http.StatusOK, "Comments retrieved successfully", comments)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post owned by the user
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDRaw.(string))

	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	if err := h.postUsecase.DeletePost(postID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete post or unauthorized")
		return
	}

	response.JSON(c, http.StatusOK, "Post deleted successfully", nil)
}
