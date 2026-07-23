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
		respondWithError(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	if req.Password == "" || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "please fill out both the email and password fields", nil)
		return
	}

	hashedpassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	params := database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedpassword,
	}

	dbUser, err := c.db.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
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
		respondWithError(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	if req.Password == "" && req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "at least the email or password fields must be filled out", nil)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Authorization header found", nil)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, c.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT token is invalid", nil)
		return
	}

	newHashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	params := database.UpdateUserEmailPasswordParams{
		Email:          req.Email,
		HashedPassword: newHashedPassword,
		ID:             userID,
	}

	dbUser, err := c.db.UpdateUserEmailPassword(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
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
		respondWithError(w, http.StatusBadRequest, "error decoding request body", nil)
		return
	}

	dbUser, err := c.db.GetUserFromEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	tokenString, err := auth.MakeJWT(dbUser.ID, c.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
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
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
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
