package auth

import (
	"errors"
	"time"

	"github.com/MentorsPath/Backend/database/repositories"
	"github.com/MentorsPath/Backend/models"
	"github.com/MentorsPath/Backend/pkg/jwt"
	"github.com/MentorsPath/Backend/pkg/password"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    repositories.UserRepository
	profileRepo repositories.ProfileRepository
	jwtGen      *jwt.Generator
	refreshGen  *jwt.Generator
}

func NewAuthService(
	userRepo repositories.UserRepository,
	profileRepo repositories.ProfileRepository,
	jwtGen, refreshGen *jwt.Generator,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		jwtGen:      jwtGen,
		refreshGen:  refreshGen,
	}
}

// Register creates a user with associated profile and optional mentor/mentee data
func (s *AuthService) Register(user *models.User, profile *models.Profile, roleData interface{}) (*models.User, string, string, error) {
	exists, err := s.userRepo.EmailExists(user.Email)
	if err != nil {
		return nil, "", "", err
	}
	if exists {
		return nil, "", "", errors.New("email already registered")
	}

	// Hash password
	hashed, err := password.Hash(user.PasswordHash)
	if err != nil {
		return nil, "", "", err
	}
	user.PasswordHash = hashed

	// Save profile first
	profile.ID = uuid.New()
	if err := s.profileRepo.CreateProfile(profile); err != nil {
		return nil, "", "", err
	}
	user.ID = uuid.New()
	user.ProfileID = &profile.ID

	// Create user
	createdUser, err := s.userRepo.Create(user)
	createdUser.Profile = profile
	if err != nil {
		return nil, "", "", err
	}

	switch user.Role {
	case models.RoleMentor:
		mp, ok := roleData.(*models.MentorProfile)
		if !ok {
			return nil, "", "", errors.New("invalid mentor profile")
		}
		mp.ID = uuid.New()
		mp.UserID = createdUser.ID
		mp.ProfileID = profile.ID
		if err := s.profileRepo.CreateMentorProfile(mp); err != nil {
			return nil, "", "", err
		}
		createdUser.MentorProfile = mp

	case models.RoleMentee:
		mp, ok := roleData.(*models.MenteeProfile)
		if !ok {
			return nil, "", "", errors.New("invalid mentee profile")
		}
		mp.ID = uuid.New()
		mp.UserID = createdUser.ID
		mp.ProfileID = profile.ID
		if err := s.profileRepo.CreateMenteeProfile(mp); err != nil {
			return nil, "", "", err
		}
		createdUser.MenteeProfile = mp

	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return createdUser, accessToken, refreshToken, nil
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(email, plainPassword string) (*models.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", "", errors.New("account inactive")
	}

	if err := password.Verify(user.PasswordHash, plainPassword); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.UpdateLoginTime(user.ID, now); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) generateTokens(user *models.User) (string, string, error) {
	accessClaims := jwtv5.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshClaims := jwtv5.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	at, err := s.jwtGen.GenerateToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	rt, err := s.refreshGen.GenerateToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (s *AuthService) Refresh(refreshToken string) (*models.User, string, string, error) {
	claims, err := s.refreshGen.ValidateToken(refreshToken)
	if err != nil {
		return nil, "", "", errors.New("invalid refresh token")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, "", "", errors.New("invalid token payload")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, "", "", err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, "", "", errors.New("user not found")
	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
