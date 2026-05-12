package auth

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetMeUsecase struct {
	userRepo repositories.UserRepository
}

func NewGetMeUsecase(userRepo repositories.UserRepository) *GetMeUsecase {
	return &GetMeUsecase{userRepo: userRepo}
}

func (u *GetMeUsecase) Execute(userID string) (*entities.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return u.userRepo.GetByID(id)
}
