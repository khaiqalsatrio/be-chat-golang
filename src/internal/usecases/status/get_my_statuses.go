package status

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetMyStatusesUsecase struct {
	statusRepo repositories.StatusRepository
}

func NewGetMyStatusesUsecase(statusRepo repositories.StatusRepository) *GetMyStatusesUsecase {
	return &GetMyStatusesUsecase{
		statusRepo: statusRepo,
	}
}

// Execute mengambil riwayat status milik user sendiri
func (u *GetMyStatusesUsecase) Execute(userID uuid.UUID) ([]entities.Status, error) {
	return u.statusRepo.GetActiveStatusesByUserID(userID)
}
