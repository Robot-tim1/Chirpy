package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Robot-tim1/Chirpy/internal/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	apiConfigParams := api.ApiConfigParams{
		Db:       db,
		Platform: os.Getenv("PLATFORM"),
		Secret:   os.Getenv("SECRET"),
	}
	server := api.NewServer(apiConfigParams)

	httpServer := &http.Server{
		Handler: server.Handler(),
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()
}
