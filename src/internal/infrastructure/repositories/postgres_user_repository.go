package repositories

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repositories.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *postgresUserRepository) GetByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *postgresUserRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *postgresUserRepository) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *postgresUserRepository) GetAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *postgresUserRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *postgresUserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.User{}, "id = ?", id).Error
}
