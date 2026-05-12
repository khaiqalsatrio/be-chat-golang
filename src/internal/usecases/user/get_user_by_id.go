package user

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetUserByIDUsecase struct {
	userRepo repositories.UserRepository
}

func NewGetUserByIDUsecase(userRepo repositories.UserRepository) *GetUserByIDUsecase {
	return &GetUserByIDUsecase{userRepo: userRepo}
}

func (u *GetUserByIDUsecase) Execute(idStr string) (*entities.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return u.userRepo.GetByID(id)
}
