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

func (apiCfg *ApiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := parameters{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		fmt.Printf("Error while parsing feed: %v\n", err)
		respondWithError(w, 400, "Unable to parse request body.")
		return
	}

	feed, err := apiCfg.DB.CreatFeed(r.Context(), database.CreatFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
		Url:       reqBody.Url,
		UserID:    user.ID,
	})
	if err != nil {
		fmt.Printf("Error while creating feed: %v\n", err)
		respondWithError(w, 500, "Something went wrong, please try again.")
		return
	}

	response := dto.Feed{}
	response.FromDbFeed(feed)
	respondWithJSON(w, 201, response)
}

// Gets all feeds from database.
func (apiCfg *ApiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a get user details request")
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		fmt.Printf("Error while fetching feed details: %v\n", err)
		return
	}

	feedsDto := make([]dto.Feed, len(feeds))
	for i, feed := range feeds {
		feedDto := &dto.Feed{}
		feedDto.FromDbFeed(feed)
		feedsDto[i] = *feedDto
	}

	respondWithJSON(w, 200, feedsDto)
}
