package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/hconn7/BlogAggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	const filePathRoot = "."
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	dbURL := os.Getenv("SQLSTRING")
	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/healthz", ReadinessCheck)
	mux.HandleFunc("GET /v1/err", ErrorCheck)
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Printf("Serving running on %s from path %s\n", port, filePathRoot)
	log.Fatal(srv.ListenAndServe())
}
