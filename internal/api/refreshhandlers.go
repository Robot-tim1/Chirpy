package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Robot-tim1/Chirpy/internal/auth"
)

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
