package api

import (
	"time"

	"github.com/google/uuid"
)

type apiError struct {
	Error string `json:"error"`
}

type chirpReq struct {
	Body string `json:"body"`
}

type polkaReq struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

type userReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResp struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type authResp struct {
	userResp
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type tokenResp struct {
	Token string `json:"token"`
}

type chirpResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
