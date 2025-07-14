package handlers

import (
	auth "github.com/MentorsPath/Backend/internal/user"
	"net/http"

	"github.com/MentorsPath/Backend/models"
	"github.com/MentorsPath/Backend/pkg/utils"
	"github.com/google/uuid"
)

type SessionHandler struct {
	userService *auth.UserService
}

func NewSessionHandler(userService *auth.UserService) *SessionHandler {
	return &SessionHandler{userService: userService}
}

func (h *SessionHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	raw := r.Context().Value(utils.UserIDKey)
	userIDStr, ok := raw.(string)
	if !ok {
		models.JSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		models.JSON(w, http.StatusBadRequest, "invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		models.JSON(w, http.StatusInternalServerError, "failed to fetch user", nil)
		return
	}

	user.PasswordHash = "" // hide password
	models.JSON(w, http.StatusOK, "success", user)
}
