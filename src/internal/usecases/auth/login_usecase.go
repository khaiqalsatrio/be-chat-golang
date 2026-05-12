package auth

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"chat-golang/src/internal/infrastructure/services"
	"chat-golang/src/pkg/utils"
	"errors"
)

type LoginUsecase struct {
	userRepo   repositories.UserRepository
	jwtService *services.JWTService
}

func NewLoginUsecase(userRepo repositories.UserRepository, jwtService *services.JWTService) *LoginUsecase {
	return &LoginUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

type LoginResult struct {
	User  *entities.User `json:"user"`
	Token string         `json:"token"`
}

func (u *LoginUsecase) Execute(email, password string) (*LoginResult, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}
