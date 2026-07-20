package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Robot-tim1/Chirpy/internal/auth"
	"github.com/Robot-tim1/Chirpy/internal/database"
)

func (c *apiConfig) handlerUserPost(w http.ResponseWriter, r *http.Request) {
	var req userReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "please enter a password", nil)
		return
	}

	hashedpassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	params := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedpassword,
	}

	dbUser, err := c.db.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user in database", err)
		return
	}

	resp := userResp{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (c *apiConfig) handlerUserPut(w http.ResponseWriter, r *http.Request) {
	var req userReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Authorization header found", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT token is invalid", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}

	params := database.UpdateUserEmailPasswordParams{
		Email:          req.Email,
		HashedPassword: newHashedPassword,
		ID:             userID,
	}

	dbUser, err := c.db.UpdateUserEmailPassword(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user record", err)
		return
	}

	resp := userResp{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerLoginEnd(w http.ResponseWriter, r *http.Request) {
	var req userReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	dbUser, err := c.db.GetUserFromEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, c.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making JWT token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	}

	_, err = c.db.CreateRefreshToken(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token record", err)
		return
	}

	resp := authResp{
		userResp: userResp{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        tokenString,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
