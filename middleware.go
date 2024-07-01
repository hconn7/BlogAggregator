package main

import (
	"github.com/hconn7/BlogAggregator/helpers"
	"github.com/hconn7/BlogAggregator/internal/database"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Missing authorization")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "ApiKey ")
		if tokenStr == authHeader {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Authorization header missing key")
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), tokenStr)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Error fetching user data from DB")
			return
		}

		next(w, r, user)
	}
}
