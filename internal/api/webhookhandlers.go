package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Robot-tim1/Chirpy/internal/auth"
)

func (c *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no Authorization header found", nil)
		return
	}
	if apiKey != c.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "wrong apikey", nil)
		return
	}

	var req polkaReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding request body", err)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = c.db.UpgradeUserRed(r.Context(), req.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
