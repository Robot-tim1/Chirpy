package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/Robot-tim1/Chirpy/internal/auth"
	"github.com/Robot-tim1/Chirpy/internal/database"
	"github.com/google/uuid"
)

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
