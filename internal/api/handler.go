package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/Robot-tim1/Chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerHealthzEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}

func (c *apiConfig) handlerChirpPost(w http.ResponseWriter, r *http.Request) {
	var chirpPost chirpPost
	if err := json.NewDecoder(r.Body).Decode(&chirpPost); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	if utf8.RuneCountInString(chirpPost.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirpPost.Body = cleanProfane(chirpPost.Body)
	params := database.CreateChirpParams{
		Body:   chirpPost.Body,
		UserID: chirpPost.UserID,
	}

	dbChirp, err := c.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp in database", err)
		return
	}

	respChirp := chirpResp(dbChirp)
	respondWithJSON(w, http.StatusCreated, respChirp)
}

func (c *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := c.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirps from database", err)
		return
	}

	respChirps := make([]chirpResp, 0, len(dbChirps))
	for _, c := range dbChirps {
		respChirps = append(respChirps, chirpResp(c))
	}

	respondWithJSON(w, http.StatusOK, respChirps)
}

func (c *apiConfig) handlerChirpGetID(w http.ResponseWriter, r *http.Request) {
	chirpString := r.PathValue("chirpID")
	chirpId, _ := uuid.Parse(chirpString)

	dbChirp, err := c.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("error chirp could not be found at id: %s", chirpString), err)
		return
	}

	respChirp := chirpResp(dbChirp)
	respondWithJSON(w, http.StatusOK, respChirp)
}

func (c *apiConfig) handlerUserEnd(w http.ResponseWriter, r *http.Request) {
	var userEmail createUser
	if err := json.NewDecoder(r.Body).Decode(&userEmail); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	dbuser, err := c.db.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user in database", err)
		return
	}

	user := userResp(dbuser)

	respondWithJSON(w, http.StatusCreated, user)
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
