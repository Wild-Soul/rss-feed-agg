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

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	log.Println("[Starting server] Loaded environment variables")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
		os.Exit(1)
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Invalid value for db connection")
		os.Exit(1)
	}

	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to DB", err)
		os.Exit(1)
	}

	apiCfg := apiConfig{
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
	v1Router.Post("/users", apiCfg.handleCreateUser)

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
