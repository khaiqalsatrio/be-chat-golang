package status

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
)

type GetStatusesUsecase struct {
	statusRepo repositories.StatusRepository
}

func NewGetStatusesUsecase(statusRepo repositories.StatusRepository) *GetStatusesUsecase {
	return &GetStatusesUsecase{
		statusRepo: statusRepo,
	}
}

// Execute mengambil status aktif dari pengguna yang terhubung, dikelompokkan per user
func (u *GetStatusesUsecase) Execute(userID uuid.UUID) (map[uuid.UUID][]entities.Status, error) {
	statuses, err := u.statusRepo.GetActiveStatusesByConcern(userID)
	if err != nil {
		return nil, err
	}

	// Kelompokkan status berdasarkan user_id
	statusMap := make(map[uuid.UUID][]entities.Status)
	for _, s := range statuses {
		statusMap[s.UserID] = append(statusMap[s.UserID], s)
	}

	return statusMap, nil
}
