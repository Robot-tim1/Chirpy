package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/Robot-tim1/Chirpy/internal/auth"
	"github.com/Robot-tim1/Chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerHealthzEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}

func (c *apiConfig) handlerChirpPost(w http.ResponseWriter, r *http.Request) {
	var req chirpReq
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

	if utf8.RuneCountInString(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	req.Body = cleanProfane(req.Body)
	params := database.CreateChirpParams{
		Body:   req.Body,
		UserID: userID,
	}

	dbChirp, err := c.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp in database", err)
		return
	}

	resp := chirpResp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (c *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := c.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps from database", err)
		return
	}

	resp := make([]chirpResp, 0, len(dbChirps))
	for _, c := range dbChirps {
		resp = append(resp, chirpResp(c))
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerChirpGetID(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing id into uuid", err)
		return
	}

	dbChirp, err := c.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp could not be found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	resp := chirpResp(dbChirp)
	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerChirpDeleteID(w http.ResponseWriter, r *http.Request) {
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

	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing id into uuid", err)
		return
	}

	dbChirp, err := c.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp could not be found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can not delete chirp post form other user", nil)
		return
	}

	err = c.db.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp record", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

	dbuser, err := c.db.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user in database", err)
		return
	}

	resp := userResp{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email:     dbuser.Email,
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
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerLoginEnd(w http.ResponseWriter, r *http.Request) {
	var req userReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	dbuser, err := c.db.GetUserFromEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, dbuser.HashedPassword)
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenString, err := auth.MakeJWT(dbuser.ID, c.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making JWT token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbuser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	}

	_, err = c.db.CreateRefreshToken(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token record", err)
		return
	}

	resp := authResp{
		userResp: userResp{
			ID:        dbuser.ID,
			CreatedAt: dbuser.CreatedAt,
			UpdatedAt: dbuser.UpdatedAt,
			Email:     dbuser.Email,
		},
		Token:        tokenString,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerRefreshEnd(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token could not be found in header", err)
		return
	}

	dbRefreshToken, err := c.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "refresh token could not be found in database", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	if dbRefreshToken.RevokedAt.Valid || time.Now().UTC().After(dbRefreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token is invalid", nil)
		return
	}

	newTokenString, err := auth.MakeJWT(dbRefreshToken.UserID, c.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making JWT token", err)
		return
	}

	resp := tokenResp{Token: newTokenString}

	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerRevokeEnd(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error refresh token could not be found", err)
		return
	}

	err = c.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *apiConfig) handlerRequestNum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	text := fmt.Sprintf(`<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>`, c.fileserverHits.Load())
	w.Write([]byte(text))
}

func (c *apiConfig) handlerResetEnd(w http.ResponseWriter, r *http.Request) {
	if c.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "no access to endpoint", nil)
	}
	c.fileserverHits.Store(0)
	c.db.DeleteUsers(r.Context())
}
