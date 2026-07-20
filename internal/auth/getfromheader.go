package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", errors.New("error no Authorization header found")
	}
	authValue = strings.TrimPrefix(authValue, "Bearer")
	authValue = strings.TrimSpace(authValue)
	return authValue, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", errors.New("error no Authorization header found")
	}
	authValue = strings.TrimPrefix(authValue, "ApiKey")
	authValue = strings.TrimSpace(authValue)
	return authValue, nil
}
