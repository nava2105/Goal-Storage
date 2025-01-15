package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

// Fetch userId from the Authentication API using the given token
func FetchUserIdFromAuthAPI(token string) (int, error) {
	authUrl := os.Getenv("AUTH_URL") // Authentication URL from .env

	client := &http.Client{}

	req, err := http.NewRequest("GET", strings.Trim(authUrl, "'"), bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token) // Pass the token as received

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("authentication API returned error")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	userId, ok := result["userId"].(float64) // JSON numbers are float64
	if !ok {
		return 0, errors.New("invalid response from authentication API")
	}

	return int(userId), nil
}