package auth

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type UploadProfilePhotoUsecase struct {
	userRepo repositories.UserRepository
}

func NewUploadProfilePhotoUsecase(userRepo repositories.UserRepository) *UploadProfilePhotoUsecase {
	return &UploadProfilePhotoUsecase{userRepo: userRepo}
}

func (u *UploadProfilePhotoUsecase) Execute(userID string, avatarURL string) (*entities.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.AvatarURL = avatarURL
	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
