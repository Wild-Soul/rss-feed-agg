package main

import (
	"fmt"
	"net/http"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/auth"
	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ExtractApiKey(r.Header)

		if err != nil {
			fmt.Printf("[Error]:[getUserHandler]: %v\n", err)
			respondWithError(w, 400, fmt.Sprintf("Invalid api key: %v\n", err.Error()))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			fmt.Printf("[Error]:[getUserHandler]: %v\n", err)
			// TODO:: need to handle different users.
			respondWithError(w, 404, fmt.Sprintf("User not found: %v\n", err.Error()))
			return
		}
		handler(w, r, user)
	}
}
