package services

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type GoogleAuthService struct {
	clientID string
}

func NewGoogleAuthService(clientID string) *GoogleAuthService {
	return &GoogleAuthService{clientID: clientID}
}

type GoogleUser struct {
	Email   string
	Name    string
	Picture string
	Sub     string
}

func (s *GoogleAuthService) VerifyToken(idToken string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, s.clientID)
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in token")
	}

	name := payload.Claims["name"].(string)
	picture := payload.Claims["picture"].(string)
	sub := payload.Claims["sub"].(string)

	return &GoogleUser{
		Email:   email,
		Name:    name,
		Picture: picture,
		Sub:     sub,
	}, nil
}
