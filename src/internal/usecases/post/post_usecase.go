package post

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type PostUsecase interface {
	CreatePost(post *entities.Post) error
	GetGlobalFeed(limit, offset int, currentUserID uuid.UUID) ([]*entities.Post, error)
	ToggleLike(postID, userID uuid.UUID) (bool, error)
	AddComment(comment *entities.Comment) error
	GetComments(postID uuid.UUID) ([]*entities.Comment, error)
	DeletePost(postID, userID uuid.UUID) error
}

type postUsecase struct {
	postRepo repositories.PostRepository
}

func NewPostUsecase(postRepo repositories.PostRepository) PostUsecase {
	return &postUsecase{postRepo: postRepo}
}

func (u *postUsecase) CreatePost(post *entities.Post) error {
	return u.postRepo.Create(post)
}

func (u *postUsecase) GetGlobalFeed(limit, offset int, currentUserID uuid.UUID) ([]*entities.Post, error) {
	return u.postRepo.FindAll(limit, offset, currentUserID)
}

func (u *postUsecase) ToggleLike(postID, userID uuid.UUID) (bool, error) {
	return u.postRepo.ToggleLike(postID, userID)
}

func (u *postUsecase) AddComment(comment *entities.Comment) error {
	return u.postRepo.AddComment(comment)
}

func (u *postUsecase) GetComments(postID uuid.UUID) ([]*entities.Comment, error) {
	return u.postRepo.GetCommentsByPostID(postID)
}

func (u *postUsecase) DeletePost(postID, userID uuid.UUID) error {
	post, err := u.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post.UserID != userID {
		return interface{}(nil).(error) // Should return an unauthorized error, but for now simple check
	}
	return u.postRepo.Delete(postID)
}
