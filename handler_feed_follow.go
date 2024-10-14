package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/dto"
	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) handleCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := parameters{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		fmt.Printf("Error while parsing feed: %v\n", err)
		respondWithError(w, 400, "Unable to parse request body.")
		return
	}

	feedFollow, err := apiCfg.DB.CreatFeedFollow(r.Context(), database.CreatFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    reqBody.FeedId,
	})
	if err != nil {
		fmt.Printf("Error while creating feed: %v\n", err)
		respondWithError(w, 500, "Something went wrong, please try again.")
		return
	}

	response := dto.FeedFollows{}
	response.FromDbFeed(feedFollow)
	respondWithJSON(w, 201, response)
}
