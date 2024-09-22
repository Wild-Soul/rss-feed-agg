package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/Wild-Soul/go-rss-feed-agg/dto"
	"github.com/Wild-Soul/go-rss-feed-agg/internal/auth"
	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
)

// Adds a new user to database.
// TODO:: Custom json parse so as to enforce required field (Name) constraint.
func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request to add user")
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Printf("[Error]:[createUserHandler]: %v\n", err)
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		fmt.Printf("[Error]:[createUserHandler]: %v\n", err)
		respondWithError(w, 500, fmt.Sprintf("Failed to create user %v", err))
		return
	}
	userdto := &dto.UserDTO{}
	userdto.FromDbUser(user)

	fmt.Printf("User added: %v\n", userdto.Name)
	respondWithJSON(w, 201, userdto)
}

// Gets a user from database based on ApiKey provided in auth header.
func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a get user details request")
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

	userdto := &dto.UserDTO{}
	userdto.FromDbUser(user)
	respondWithJSON(w, 200, userdto)
}
