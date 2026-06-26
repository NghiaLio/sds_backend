package model

import (
	"errors"
	"strings"
)

// RegisterRequest holds registration details.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate checks registration inputs.
func (r *RegisterRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	if len(r.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

// LoginRequest holds credentials for logging in.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse returns the token to the mobile client.
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Username    string `json:"username"`
}
