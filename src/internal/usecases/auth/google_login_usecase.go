package auth

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"chat-golang/src/internal/infrastructure/services"
	"gorm.io/gorm"
)

type GoogleLoginUsecase struct {
	userRepo          repositories.UserRepository
	googleAuthService *services.GoogleAuthService
	jwtService        *services.JWTService
}

func NewGoogleLoginUsecase(
	userRepo repositories.UserRepository,
	googleAuthService *services.GoogleAuthService,
	jwtService *services.JWTService,
) *GoogleLoginUsecase {
	return &GoogleLoginUsecase{
		userRepo:          userRepo,
		googleAuthService: googleAuthService,
		jwtService:        jwtService,
	}
}

func (u *GoogleLoginUsecase) Execute(idToken string) (*LoginResult, error) {
	googleUser, err := u.googleAuthService.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByEmail(googleUser.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user
			user = &entities.User{
				Username:  googleUser.Email, // Default username
				Email:     googleUser.Email,
				GoogleID:  googleUser.Sub,
				AvatarURL: googleUser.Picture,
				Status:    "ONLINE",
			}
			err = u.userRepo.Create(user)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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
