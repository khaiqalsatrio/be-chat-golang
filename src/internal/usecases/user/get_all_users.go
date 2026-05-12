package user

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
)

type GetAllUsersUsecase struct {
	userRepo repositories.UserRepository
}

func NewGetAllUsersUsecase(userRepo repositories.UserRepository) *GetAllUsersUsecase {
	return &GetAllUsersUsecase{userRepo: userRepo}
}

func (u *GetAllUsersUsecase) Execute() ([]entities.User, error) {
	return u.userRepo.GetAll()
}
