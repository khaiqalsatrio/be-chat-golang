package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresPostRepository struct {
	db *gorm.DB
}

func NewPostgresPostRepository(db *gorm.DB) repositories.PostRepository {
	return &postgresPostRepository{db: db}
}

func (r *postgresPostRepository) Create(post *entities.Post) error {
	return r.db.Create(post).Error
}

func (r *postgresPostRepository) FindAll(limit, offset int, currentUserID uuid.UUID) ([]*entities.Post, error) {
	var posts []*entities.Post
	err := r.db.Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Preload("User").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		// Count likes
		var likesCount int64
		r.db.Model(&entities.Like{}).Where("post_id = ?", post.ID).Count(&likesCount)
		post.LikesCount = int(likesCount)

		// Count comments
		var commentsCount int64
		r.db.Model(&entities.Comment{}).Where("post_id = ?", post.ID).Count(&commentsCount)
		post.CommentsCount = int(commentsCount)

		// Check if liked by current user
		if currentUserID != uuid.Nil {
			var count int64
			r.db.Model(&entities.Like{}).Where("post_id = ? AND user_id = ?", post.ID, currentUserID).Count(&count)
			post.IsLiked = count > 0
		}
	}

	return posts, nil
}

func (r *postgresPostRepository) FindByID(id uuid.UUID) (*entities.Post, error) {
	var post entities.Post
	err := r.db.Preload("User").First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	var likesCount int64
	r.db.Model(&entities.Like{}).Where("post_id = ?", post.ID).Count(&likesCount)
	post.LikesCount = int(likesCount)

	var commentsCount int64
	r.db.Model(&entities.Comment{}).Where("post_id = ?", post.ID).Count(&commentsCount)
	post.CommentsCount = int(commentsCount)

	return &post, nil
}

func (r *postgresPostRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Post{}, "id = ?", id).Error
}

func (r *postgresPostRepository) ToggleLike(postID, userID uuid.UUID) (bool, error) {
	var like entities.Like
	result := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like)

	if result.Error == gorm.ErrRecordNotFound {
		// Like it
		newLike := entities.Like{
			PostID: postID,
			UserID: userID,
		}
		if err := r.db.Create(&newLike).Error; err != nil {
			return false, err
		}
		return true, nil
	} else if result.Error != nil {
		return false, result.Error
	}

	// Unlike it
	if err := r.db.Delete(&like).Error; err != nil {
		return false, err
	}
	return false, nil
}

func (r *postgresPostRepository) AddComment(comment *entities.Comment) error {
	return r.db.Create(comment).Error
}

func (r *postgresPostRepository) GetCommentsByPostID(postID uuid.UUID) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	err := r.db.Where("post_id = ?", postID).
		Order("created_at asc").
		Preload("User").
		Find(&comments).Error
	return comments, err
}
