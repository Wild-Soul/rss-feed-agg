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
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			fmt.Printf("[Error]:[getUserHandler]: %v\n", err)
			// TODO:: need to handle different users.
			respondWithError(w, 400, fmt.Sprintf("Failed to get user: %v", err.Error()))
			return
		}
		handler(w, r, user)
	}
}
