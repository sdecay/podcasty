package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sdecay/podcasty/internal/database"
)

func (config *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(req.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	user, err := config.DB.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not create user %s", err))
		return
	}

	respondWithJson(writer, http.StatusCreated, dbUserToUser(user))
}

func (config *apiConfig) handlerGetUser(writer http.ResponseWriter, req *http.Request, user database.User) {
	respondWithJson(writer, http.StatusOK, dbUserToUser(user))
}

func (config *apiConfig) handlerGetUsersPosts(writer http.ResponseWriter, req *http.Request, user database.User) {
	posts, err := config.DB.GetUsersPosts(req.Context(), database.GetUsersPostsParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Sprintf("could not get your latest posts: %s", err))
		return
	}

	respondWithJson(writer, http.StatusOK, dbPostsToPosts(posts))
}
