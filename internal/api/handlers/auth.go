package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
	Email    string               `json:"mailer"`
	Password string               `json:"password"`
	Role     string               `json:"role"` // "mentor" or "mentee"
	Profile  *models.ProfileInput `json:"profile"`
}

func generateAvatarURL(name string) string {
	escapedName := url.QueryEscape(name)
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random&rounded=true", escapedName)
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

	profile := &models.Profile{
		FirstName:   req.Profile.FirstName,
		LastName:    req.Profile.LastName,
		Bio:         req.Profile.Bio,
		AvatarURL:   generateAvatarURL(fmt.Sprintf("%s %s", req.Profile.FirstName, req.Profile.LastName)),
		ImageURL:    req.Profile.ImageURL,
		Headline:    req.Profile.Headline,
		WebsiteURL:  req.Profile.WebsiteURL,
		LinkedInURL: req.Profile.LinkedInURL,
		Twitter:     req.Profile.Twitter,
		Timezone:    req.Profile.Timezone,
	}

	createdUser, accessToken, refreshToken, err := h.authService.Register(user, profile, roleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message":       "signup successful",
		"user_id":       createdUser,
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
		"user":          user,
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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	models.JSON(w, http.StatusOK, "token generated", "")
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	models.JSON(w, http.StatusOK, "password reset successful", nil)
}
