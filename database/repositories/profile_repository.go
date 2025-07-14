package repositories

import (
	"github.com/MentorsPath/Backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfile(profile *models.Profile) error
	GetProfile(profileID uuid.UUID) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
	DeleteProfile(profileID uuid.UUID) error

	CreateMentorProfile(profile *models.MentorProfile) error
	GetMentorProfile(userID uuid.UUID) (*models.MentorProfile, error)
	UpdateMentorProfile(profile *models.MentorProfile) error
	DeleteMentorProfile(userID uuid.UUID) error

	CreateMenteeProfile(profile *models.MenteeProfile) error
	GetMenteeProfile(userID uuid.UUID) (*models.MenteeProfile, error)
	UpdateMenteeProfile(profile *models.MenteeProfile) error
	DeleteMenteeProfile(userID uuid.UUID) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

// ---------- Generic Profile ----------

func (r *profileRepository) CreateProfile(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) GetProfile(profileID uuid.UUID) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.First(&profile, "id = ?", profileID).Error
	return &profile, err
}

func (r *profileRepository) UpdateProfile(profile *models.Profile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) DeleteProfile(profileID uuid.UUID) error {
	return r.db.Delete(&models.Profile{}, "id = ?", profileID).Error
}

// ---------- Mentor Profile ----------

func (r *profileRepository) CreateMentorProfile(profile *models.MentorProfile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) GetMentorProfile(userID uuid.UUID) (*models.MentorProfile, error) {
	var profile models.MentorProfile
	err := r.db.First(&profile, "user_id = ?", userID).Error
	return &profile, err
}

func (r *profileRepository) UpdateMentorProfile(profile *models.MentorProfile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) DeleteMentorProfile(userID uuid.UUID) error {
	return r.db.Delete(&models.MentorProfile{}, "user_id = ?", userID).Error
}

// ---------- Mentee Profile ----------

func (r *profileRepository) CreateMenteeProfile(profile *models.MenteeProfile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) GetMenteeProfile(userID uuid.UUID) (*models.MenteeProfile, error) {
	var profile models.MenteeProfile
	err := r.db.First(&profile, "user_id = ?", userID).Error
	return &profile, err
}

func (r *profileRepository) UpdateMenteeProfile(profile *models.MenteeProfile) error {
	return r.db.Save(profile).Error
}

func (r *profileRepository) DeleteMenteeProfile(userID uuid.UUID) error {
	return r.db.Delete(&models.MenteeProfile{}, "user_id = ?", userID).Error
}
