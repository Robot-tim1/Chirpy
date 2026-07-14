package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

func handlerHealthzEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}

func handlerValidateChirpEnd(w http.ResponseWriter, r *http.Request) {
	var chirpPost chirpPost
	if err := json.NewDecoder(r.Body).Decode(&chirpPost); err != nil {
		log.Printf("Error decoding chirpPost: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding request body")
		return
	}

	if utf8.RuneCountInString(chirpPost.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpPost.Body = cleanProfane(chirpPost.Body)
	respondWithJSON(w, http.StatusOK, cleanBodyResp{CleanedBody: chirpPost.Body})
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

func (c *apiConfig) handlerResetNum(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
}
