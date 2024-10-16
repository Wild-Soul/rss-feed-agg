package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wild-Soul/go-rss-feed-agg/dto"
	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/go-chi/chi/v5"
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
		respondWithError(w, 500, err.Error())
		return
	}

	response := dto.FeedFollows{}
	response.FromDbFeed(feedFollow)
	respondWithJSON(w, 201, response)
}

func (apiCfg *ApiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		fmt.Printf("Couldn't get feed follows: %v\n", err)
		respondWithError(w, 500, fmt.Sprintf("Couldn't get feed follows: %v", err))
		return
	}

	feedFollowsDto := make([]dto.FeedFollows, len(feedFollows))
	for i, feedFollow := range feedFollows {
		response := dto.FeedFollows{}
		response.FromDbFeed(feedFollow)
		feedFollowsDto[i] = response
	}
	fmt.Printf("Retrieved feed_follows for user: %v", user.ID)
	respondWithJSON(w, 200, feedFollowsDto)
}

func (apiCfg *ApiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		fmt.Printf("Couldn't parse feed follow id: %v\n", err)
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		fmt.Printf("Couldn't delete feed follow, error: %v", err)
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}
	fmt.Printf("Deleted feed_follow with id: %v for user: %v", feedFollowId, user.ID)
	respondWithJSON(w, 200, struct{}{})
}
