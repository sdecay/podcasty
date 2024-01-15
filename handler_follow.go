package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sdecay/podcasty/internal/database"
)

func (config *apiConfig) handlerFollowFeed(writer http.ResponseWriter, req *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(req.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	follow, err := config.DB.CreateFollow(req.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not follow feed %s", err))
		return
	}

	respondWithJson(writer, http.StatusCreated, dbFollowToFollow(follow))
}
