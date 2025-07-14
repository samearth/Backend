package handlers

import (
	"encoding/json"
	"github.com/MentorsPath/Backend/internal/api/profile"
	auth "github.com/MentorsPath/Backend/internal/user"
	"github.com/MentorsPath/Backend/pkg/utils"
	"net/http"

	"github.com/MentorsPath/Backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProfileHandler struct {
	service     *profile.Service
	userService *auth.UserService
}

func NewProfileHandler(service *profile.Service, userService *auth.UserService) *ProfileHandler {
	return &ProfileHandler{
		service:     service,
		userService: userService,
	}
}

func getUserIDFromPath(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	return uuid.Parse(userIDStr)
}

// ---------- Profile ----------

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(utils.UserIDKey)
	userIDStr, ok := raw.(string)
	if !ok {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", userIDStr)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID format", userIDStr)
		return
	}

	// ✅ Step 1: Get the user
	user, err := h.userService.GetByID(userID)
	if err != nil || user.ProfileID == nil {
		models.JSON(w, http.StatusNotFound, "user or profile not found", nil)
		return
	}

	// ✅ Step 2: Use ProfileID to fetch profile
	profile, err := h.service.GetProfile(*user.ProfileID)
	if err != nil {
		models.JSON(w, http.StatusNotFound, "profile not found", nil)
		return
	}

	models.JSON(w, http.StatusOK, "success", profile)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(utils.UserIDKey)
	userIDStr, ok := raw.(string)
	if !ok {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID format", nil)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil || user.ProfileID == nil {
		models.JSON(w, http.StatusNotFound, "user or profile not found", nil)
		return
	}

	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid body", nil)
		return
	}

	profile.ID = *user.ProfileID

	if err := h.service.UpdateProfile(&profile); err != nil {
		models.JSON(w, http.StatusInternalServerError, "failed to update", nil)
		return
	}
	models.JSON(w, http.StatusOK, "updated successfully", nil)
}

func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(utils.UserIDKey)
	userIDStr, ok := raw.(string)
	if !ok {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID format", nil)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil || user.ProfileID == nil {
		models.JSON(w, http.StatusNotFound, "user or profile not found", nil)
		return
	}

	if err := h.service.DeleteProfile(*user.ProfileID); err != nil {
		models.JSON(w, http.StatusInternalServerError, "failed to delete", nil)
		return
	}
	models.JSON(w, http.StatusOK, "deleted successfully", nil)
}

// ---------- Mentor Profile ----------

func (h *ProfileHandler) GetMentorProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	profile, err := h.service.GetMentorProfile(userID)
	if err != nil {
		models.JSON(w, http.StatusNotFound, "mentor profile not found", nil)
		return
	}
	models.JSON(w, http.StatusOK, "success", profile)
}

func (h *ProfileHandler) UpdateMentorProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	var profile models.MentorProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid body", nil)
		return
	}

	profile.UserID = userID

	if err := h.service.UpdateMentorProfile(&profile); err != nil {
		models.JSON(w, http.StatusInternalServerError, "failed to update", nil)
		return
	}
	models.JSON(w, http.StatusOK, "updated successfully", nil)
}

// ---------- Mentee Profile ----------

func (h *ProfileHandler) GetMenteeProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	profile, err := h.service.GetMenteeProfile(userID)
	if err != nil {
		models.JSON(w, http.StatusNotFound, "mentee profile not found", nil)
		return
	}
	models.JSON(w, http.StatusOK, "success", profile)
}

func (h *ProfileHandler) UpdateMenteeProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromPath(r)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	var profile models.MenteeProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid body", nil)
		return
	}

	profile.UserID = userID

	if err := h.service.UpdateMenteeProfile(&profile); err != nil {
		models.JSON(w, http.StatusInternalServerError, "failed to update", nil)
		return
	}
	models.JSON(w, http.StatusOK, "updated successfully", nil)
}
