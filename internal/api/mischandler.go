package api

import (
	"fmt"
	"net/http"
)

func handlerHealthzEnd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
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
