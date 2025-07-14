package auth

import (
	"github.com/MentorsPath/Backend/database/repositories"
	"github.com/MentorsPath/Backend/models"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
