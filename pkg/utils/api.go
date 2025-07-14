package utils

import (
	"context"
	"encoding/json"
	"github.com/MentorsPath/Backend/models"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{"error": message})
}

func ParseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

type contextKey string

const userClaimsKey contextKey = "userClaims"

// SetUserClaims stores user claims in the context.
func SetUserClaims(ctx context.Context, claims *models.AccessTokenClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserClaims retrieves user claims from the context.
func GetUserClaims(ctx context.Context) *models.AccessTokenClaims {
	if claims, ok := ctx.Value(userClaimsKey).(*models.AccessTokenClaims); ok {
		return claims
	}
	return nil
}
