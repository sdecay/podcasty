package main

import (
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	if code > http.StatusUnavailableForLegalReasons {
		log.Println("500+ error:", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(writer, code, errorResponse{
		Error: msg,
	})
}

func handlerReady(writer http.ResponseWriter, req *http.Request) {
	respondWithJson(writer, http.StatusAccepted, struct{}{})
}
