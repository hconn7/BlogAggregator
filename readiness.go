package main

import (
	"net/http"
)

func ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}

	RespondWithJson(w, http.StatusOK, response{
		Status: "ok",
	})
}
func ErrorCheck(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Status string `json:"error"`
	}
	type response struct {
		Error string `json:"error"`
	}
	RespondWithError(w, http.StatusInternalServerError, "Error working")
	return
}
