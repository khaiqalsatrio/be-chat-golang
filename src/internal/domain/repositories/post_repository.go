package repositories

import (
	"chat-golang/src/internal/domain/entities"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(post *entities.Post) error
	FindAll(limit, offset int, currentUserID uuid.UUID) ([]*entities.Post, error)
	FindByID(id uuid.UUID) (*entities.Post, error)
	Delete(id uuid.UUID) error
	ToggleLike(postID, userID uuid.UUID) (bool, error)
	AddComment(comment *entities.Comment) error
	GetCommentsByPostID(postID uuid.UUID) ([]*entities.Comment, error)
}
