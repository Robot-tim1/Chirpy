package api

import (
	"time"

	"github.com/google/uuid"
)

type apiError struct {
	Error string `json:"error"`
}

type chirpReq struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type userReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type chirpResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
