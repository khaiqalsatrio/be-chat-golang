package status

import (
	"chat-golang/src/internal/domain/repositories"
	"errors"

	"github.com/google/uuid"
)

type DeleteStatusUsecase struct {
	statusRepo repositories.StatusRepository
}

func NewDeleteStatusUsecase(statusRepo repositories.StatusRepository) *DeleteStatusUsecase {
	return &DeleteStatusUsecase{
		statusRepo: statusRepo,
	}
}

// Execute menghapus status, hanya user pemilik yang bisa menghapus
func (u *DeleteStatusUsecase) Execute(statusID uuid.UUID, userID uuid.UUID) error {
	status, err := u.statusRepo.GetStatusByID(statusID)
	if err != nil {
		return errors.New("status not found")
	}

	// Verify owner
	if status.UserID != userID {
		return errors.New("unauthorized: only owner can delete status")
	}

	return u.statusRepo.DeleteStatusByID(statusID)
}
