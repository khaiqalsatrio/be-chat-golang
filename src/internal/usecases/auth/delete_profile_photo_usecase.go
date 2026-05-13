package auth

import (
	"chat-golang/src/internal/domain/repositories"
	"github.com/google/uuid"
)

type DeleteProfilePhotoUsecase struct {
	userRepo repositories.UserRepository
}

func NewDeleteProfilePhotoUsecase(userRepo repositories.UserRepository) *DeleteProfilePhotoUsecase {
	return &DeleteProfilePhotoUsecase{userRepo: userRepo}
}

func (u *DeleteProfilePhotoUsecase) Execute(userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	user.AvatarURL = ""
	return u.userRepo.Update(user)
}
