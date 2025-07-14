package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MentorsPath/Backend/internal/auth"
	"github.com/MentorsPath/Backend/models"
	_ "github.com/google/uuid"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

type SignupRequest struct {
	Email      string                `json:"email"`
	Password   string                `json:"password"`
	Role       string                `json:"role"` // "mentor" or "mentee"
	Profile    models.Profile        `json:"profile"`
	MentorData *models.MentorProfile `json:"mentor_profile,omitempty"`
	MenteeData *models.MenteeProfile `json:"mentee_profile,omitempty"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         models.UserRole(req.Role),
	}

	var roleData interface{}
	if req.Role == "mentor" {
		roleData = req.MentorData
	} else if req.Role == "mentee" {
		roleData = req.MenteeData
	} else {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	createdUser, accessToken, refreshToken, err := h.authService.Register(user, &req.Profile, roleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message":       "signup successful",
		"user_id":       createdUser.ID.String(),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	models.JSON(w, http.StatusOK, "success", resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       user.ID.String(),
		"role":          user.Role,
	}
	models.JSON(w, http.StatusOK, "success", resp)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		models.JSON(w, http.StatusBadRequest, "invalid refresh token request", nil)
		return
	}

	user, newAccessToken, newRefreshToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		models.JSON(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	models.JSON(w, http.StatusOK, "success", map[string]interface{}{
		"user_id":       user.ID.String(),
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
