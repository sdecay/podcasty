package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

func (config *apiConfig) handlerGetFollowed(writer http.ResponseWriter, req *http.Request, user database.User) {
	followed, err := config.DB.GetFollowed(req.Context(), user.ID)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not get who you follow: %s", err))
		return
	}

	respondWithJson(writer, http.StatusOK, dbFollowedToFollowed(followed))
}

func (config *apiConfig) handlerDeleteFollow(writer http.ResponseWriter, req *http.Request, user database.User) {
	feedID := chi.URLParam(req, "followID")

	feed, err := uuid.Parse(feedID)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("invalid uuid format: %s", err))
		return
	}

	err = config.DB.DeleteFollowed(req.Context(), database.DeleteFollowedParams{
		ID:     feed,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not delete feed %s", err))
		return
	}

	respondWithJson(writer, http.StatusTeapot, struct{}{})
}
