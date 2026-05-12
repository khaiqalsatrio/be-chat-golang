package auth

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"chat-golang/src/pkg/utils"
	"errors"
)

type RegisterUsecase struct {
	userRepo repositories.UserRepository
}

func NewRegisterUsecase(userRepo repositories.UserRepository) *RegisterUsecase {
	return &RegisterUsecase{userRepo: userRepo}
}

func (u *RegisterUsecase) Execute(username, email, password string) (*entities.User, error) {
	// Check if user already exists
	_, err := u.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	_, err = u.userRepo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	err = u.userRepo.Create(user)
	return user, err
}
