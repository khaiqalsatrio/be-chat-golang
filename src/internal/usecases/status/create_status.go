package status

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type CreateStatusUsecase struct {
	statusRepo repositories.StatusRepository
}

func NewCreateStatusUsecase(statusRepo repositories.StatusRepository) *CreateStatusUsecase {
	return &CreateStatusUsecase{
		statusRepo: statusRepo,
	}
}

type CreateStatusRequest struct {
	MediaURL  string `json:"media_url" binding:"required"`
	MediaType string `json:"media_type" binding:"required,oneof=image video"`
	Caption   string `json:"caption"`
}

func (u *CreateStatusUsecase) Execute(userID uuid.UUID, req CreateStatusRequest) (*entities.Status, error) {
	status := &entities.Status{
		ID:        uuid.New(),
		UserID:    userID,
		MediaURL:  req.MediaURL,
		MediaType: req.MediaType,
		Caption:   req.Caption,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := u.statusRepo.CreateStatus(status)
	if err != nil {
		return nil, err
	}

	return status, nil
}
