package main

import (
	"net/http"
)

func handlerError(writer http.ResponseWriter, req *http.Request) {
	respondWithError(writer, http.StatusBadRequest, "explodey time!")
}
