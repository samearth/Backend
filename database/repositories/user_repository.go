package repositories

import (
	"errors"
	"time"

	"github.com/MentorsPath/Backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	EmailExists(email string) (bool, error)
	UpdateLoginTime(id uuid.UUID, loginTime time.Time) error
	UpdatePassword(userID uuid.UUID, hashedPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateLoginTime(userID uuid.UUID, loginTime time.Time) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", loginTime).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Profile").
		Preload("MentorProfile").
		Preload("MenteeProfile").
		Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Profile").
		Preload("MentorProfile").
		Preload("MenteeProfile").
		Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("mailer = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword).Error
}
