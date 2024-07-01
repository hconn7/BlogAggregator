package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hconn7/BlogAggregator/helpers"
	"github.com/hconn7/BlogAggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      params.Name,
		Url:       params.Url,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}
	cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: feed.ID,
		UserID: user.ID,
	})
	helpers.RespondWithJson(w, http.StatusOK, databaseFeedToFeed(feed))
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't get feeds from database")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, feeds)
}
