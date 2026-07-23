package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"unicode/utf8"

	"github.com/Robot-tim1/Chirpy/internal/auth"
	"github.com/Robot-tim1/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerChirpPost(w http.ResponseWriter, r *http.Request) {
	var req chirpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", nil)
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
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
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
	var dbChirps []database.Chirp
	var err error
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, ParseErr := uuid.Parse(authorIDString)
		if ParseErr != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing id into uuid", nil)
			return
		}
		dbChirps, err = c.db.GetChirpsAuthorID(r.Context(), authorID)
	} else {
		dbChirps, err = c.db.GetChirps(r.Context())
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "author has no posts", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	sortQuery := r.URL.Query().Get("sort")
	if sortQuery != "" {
		if sortQuery != "asc" && sortQuery != "desc" {
			respondWithError(w, http.StatusBadRequest, "no sort query of that kind", nil)
			return
		}
		if sortQuery == "desc" {
			sort.Slice(dbChirps, func(i, j int) bool { return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt) })
		}
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
		respondWithError(w, http.StatusBadRequest, "error parsing id into uuid", nil)
		return
	}

	dbChirp, err := c.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp could not be found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	resp := chirpResp(dbChirp)
	respondWithJSON(w, http.StatusOK, resp)
}

func (c *apiConfig) handlerChirpDeleteID(w http.ResponseWriter, r *http.Request) {
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

	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing id into uuid", nil)
		return
	}

	dbChirp, err := c.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp could not be found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can not delete chirp post from other user", nil)
		return
	}

	err = c.db.DeleteChirp(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
