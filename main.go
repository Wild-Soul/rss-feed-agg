package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Wild-Soul/go-rss-feed-agg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	log.Println("[Starting server] Loaded environment variables")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalln("Invalid value for db connection")
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln("Unable to connect to DB", err)
	}

	apiCfg := ApiConfig{
		DB: database.New(dbConn),
	}

	router := chi.NewRouter()

	// CORS Setup
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handleReadiness)
	v1Router.Post("/users", apiCfg.createUserHandler)
	v1Router.Get("/users", apiCfg.authMiddleware(apiCfg.getUserHandler))
	v1Router.Post("/feeds", apiCfg.authMiddleware(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.getFeedsHandler)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("HTTP server starting on port: %v\n", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("Unable to server http requests:", err)
	}
}
