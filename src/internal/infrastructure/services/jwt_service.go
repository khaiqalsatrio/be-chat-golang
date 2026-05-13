package services

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret    string
	blacklist map[string]time.Time
	mu        sync.RWMutex
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret:    secret,
		blacklist: make(map[string]time.Time),
	}
}

func (s *JWTService) GenerateToken(userID uuid.UUID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID.String(),
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if s.isBlacklisted(tokenString) {
		return nil, errors.New("token is blacklisted")
	}

	return token, nil
}

func (s *JWTService) BlacklistToken(tokenString string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blacklist[tokenString] = expiresAt
}

func (s *JWTService) isBlacklisted(tokenString string) bool {
	s.mu.RLock()
	expiresAt, ok := s.blacklist[tokenString]
	s.mu.RUnlock()

	if !ok {
		return false
	}

	if time.Now().After(expiresAt) {
		s.mu.Lock()
		delete(s.blacklist, tokenString)
		s.mu.Unlock()
		return false
	}

	return true
}
