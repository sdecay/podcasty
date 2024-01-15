package main

import (
	"fmt"
	"net/http"

	auth "github.com/sdecay/podcasty/internal"
	"github.com/sdecay/podcasty/internal/database"
)

type authorizedHandler func(http.ResponseWriter, *http.Request, database.User)

// TODO: pass in config instead of method it so i can break handlers out
// of main
func (config *apiConfig) middlewareAuth(handler authorizedHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		apiKey, err := auth.GetAPIKey(req.Header)
		if err != nil {
			respondWithError(writer, http.StatusUnauthorized, fmt.Sprintf("could not get API key: %s", err))
			return
		}

		user, err := config.DB.GetUserByAPIKey(req.Context(), apiKey)
		if err != nil {
			respondWithError(writer, http.StatusUnauthorized, fmt.Sprintf("could not get user: %s", err))
			return
		}

		handler(writer, req, user)
	}
}
