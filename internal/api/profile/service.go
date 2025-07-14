package profile

import (
	_ "errors"

	"github.com/MentorsPath/Backend/database/repositories"
	"github.com/MentorsPath/Backend/models"
	"github.com/google/uuid"
)

type Service struct {
	repo repositories.ProfileRepository
}

func NewProfileService(repo repositories.ProfileRepository) *Service {
	return &Service{repo: repo}
}

// --- Profile Methods ---

func (s *Service) CreateProfile(profile *models.Profile) error {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	return s.repo.CreateProfile(profile)
}

func (s *Service) GetProfile(userID uuid.UUID) (*models.Profile, error) {
	return s.repo.GetProfile(userID)
}

func (s *Service) UpdateProfile(profile *models.Profile) error {
	return s.repo.UpdateProfile(profile)
}

func (s *Service) DeleteProfile(userID uuid.UUID) error {
	return s.repo.DeleteProfile(userID)
}

// --- Mentor Profile Methods ---

func (s *Service) CreateMentorProfile(profile *models.MentorProfile) error {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	return s.repo.CreateMentorProfile(profile)
}

func (s *Service) GetMentorProfile(userID uuid.UUID) (*models.MentorProfile, error) {
	return s.repo.GetMentorProfile(userID)
}

func (s *Service) UpdateMentorProfile(profile *models.MentorProfile) error {
	return s.repo.UpdateMentorProfile(profile)
}

func (s *Service) DeleteMentorProfile(userID uuid.UUID) error {
	return s.repo.DeleteMentorProfile(userID)
}

// --- Mentee Profile Methods ---

func (s *Service) CreateMenteeProfile(profile *models.MenteeProfile) error {
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	return s.repo.CreateMenteeProfile(profile)
}

func (s *Service) GetMenteeProfile(userID uuid.UUID) (*models.MenteeProfile, error) {
	return s.repo.GetMenteeProfile(userID)
}

func (s *Service) UpdateMenteeProfile(profile *models.MenteeProfile) error {
	return s.repo.UpdateMenteeProfile(profile)
}

func (s *Service) DeleteMenteeProfile(userID uuid.UUID) error {
	return s.repo.DeleteMenteeProfile(userID)
}
