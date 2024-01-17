package main

import (
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	if code >= http.StatusInternalServerError {
		log.Println("bad Thing error", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(writer, code, errorResponse{
		Error: msg,
	})
}

func handlerReady(writer http.ResponseWriter, req *http.Request) {
	respondWithJson(writer, http.StatusOK, struct{}{})
}
