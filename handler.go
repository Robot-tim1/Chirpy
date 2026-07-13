package main

import (
	"fmt"
	"net/http"
)

func handlerEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
}

func (c *apiConfig) handlerRequestNum(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
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
