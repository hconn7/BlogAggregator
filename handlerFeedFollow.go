package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hconn7/BlogAggregator/helpers"
	"github.com/hconn7/BlogAggregator/internal/database"
)

func (cfg *apiConfig) handlderCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	type respone struct {
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Couldn't decode parametrs")
		return
	}

	feedID, err := uuid.Parse(params.FeedID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid FeedID format")
		return
	}
	err = cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{

		FeedID: feedID,
		UserID: user.ID,
	})
	feedParams, err := cfg.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, feedParams)
}

func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't get feed follows from database")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, feedFollows)
}
func (cfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := r.PathValue("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete feed follow from database")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, "Feed follow deleted successfully")
}
