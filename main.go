package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlderCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerGetFeeds)
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/healthz", ReadinessCheck)
	mux.HandleFunc("GET /v1/err", ErrorCheck)
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	log.Printf("Serving running on %s from path %s\n", port, filePathRoot)
	log.Fatal(srv.ListenAndServe())
}
