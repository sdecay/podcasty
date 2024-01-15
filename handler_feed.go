package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sdecay/podcasty/internal/database"
)

func (config *apiConfig) handlerCreateFeed(writer http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	feed, err := config.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not create feed %s", err))
		return
	}

	respondWithJson(writer, http.StatusCreated, dbFeedtoFeed(feed))
}
