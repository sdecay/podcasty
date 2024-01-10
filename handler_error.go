package main

import (
	"log"
	"net/http"
)

func handlerError(writer http.ResponseWriter, req *http.Request) {
	log.Print("got an error")
	respondWithError(writer, http.StatusBadRequest, "Explodey time!")
}
