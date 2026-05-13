package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresStatusRepository struct {
	db *gorm.DB
}

func NewPostgresStatusRepository(db *gorm.DB) repositories.StatusRepository {
	return &postgresStatusRepository{db: db}
}

func (r *postgresStatusRepository) CreateStatus(status *entities.Status) error {
	return r.db.Create(status).Error
}

func (r *postgresStatusRepository) GetStatusByID(id uuid.UUID) (*entities.Status, error) {
	var status entities.Status
	err := r.db.Preload("User").First(&status, "id = ?", id).Error
	return &status, err
}

func (r *postgresStatusRepository) GetActiveStatusesByUserID(userID uuid.UUID) ([]entities.Status, error) {
	var statuses []entities.Status
	err := r.db.
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Preload("User").
		Order("created_at DESC").
		Find(&statuses).Error
	return statuses, err
}

// GetActiveStatusesByConcern mengambil status dari pengguna yang terhubung (teman/kontak)
// Ini adalah placeholder - implementasi sebenarnya tergantung pada struktur koneksi di database
func (r *postgresStatusRepository) GetActiveStatusesByConcern(userID uuid.UUID) ([]entities.Status, error) {
	var statuses []entities.Status
	
	// Ambil status dari semua user (dalam implementasi real, gunakan JOIN dengan tabel koneksi/friendships)
	// Untuk sekarang, kita ambil status aktif dari semua user kecuali user sendiri
	err := r.db.
		Where("user_id != ? AND expires_at > ?", userID, time.Now()).
		Preload("User").
		Order("created_at DESC").
		Find(&statuses).Error
	
	return statuses, err
}

func (r *postgresStatusRepository) DeleteStatusByID(id uuid.UUID) error {
	return r.db.Delete(&entities.Status{}, "id = ?", id).Error
}

func (r *postgresStatusRepository) DeleteExpiredStatuses() error {
	return r.db.Delete(&entities.Status{}, "expires_at <= ?", time.Now()).Error
}

func (r *postgresStatusRepository) GetAllExpiredStatuses() ([]entities.Status, error) {
	var statuses []entities.Status
	err := r.db.
		Where("expires_at <= ?", time.Now()).
		Find(&statuses).Error
	return statuses, err
}
