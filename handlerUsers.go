package main

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hconn7/BlogAggregator/helpers"

	"net/http"
	"time"

	"github.com/hconn7/BlogAggregator/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't update DB")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	helpers.RespondWithJson(w, http.StatusOK, databaseUserToUser(user))
}
