package service

import (
	"Goal-Storage/utils"
	"errors"
	"net/http"
)

// GetUserIDFromRequest is an HTTP handler that returns the user ID from the Authentication API
func GetUserIDFromRequest(r *http.Request) (int, error) {
	// Extract the token from the Authorization header
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, errors.New("authorization token missing")
	}
	// Validate token format (should start with "Bearer")
	if len(token) < 7 || token[:7] != "Bearer " {
		return 0, errors.New("invalid authorization token format")
	}
	// Call the FetchUserIdFromAuthAPI to get the user ID
	userID, err := utils.FetchUserIdFromAuthAPI(token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
