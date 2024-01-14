package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/lius-new/rssagg/internal/database"
)

var (
	portString string
	apiCfg     apiConfig
)

type apiConfig struct {
	DB *database.Queries
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if portString = os.Getenv("PORT"); portString == "" {
		log.Fatal("Port is not found in the environment")
	}
	var (
		dbURLTemp string
		conn      *sql.DB
		err       error
	)
	if dbURLTemp = os.Getenv("DB_URL"); portString == "" {
		log.Fatal("Port is not found in the environment")
	}

	if conn, err = sql.Open("postgres", dbURLTemp); err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	queries := database.New(conn)

	apiCfg = apiConfig{
		DB: queries,
	}
}

func main() {
	start()
}

func start() {
	router := chi.NewRouter()

	registerCors(router)

	v1Rotuer := chi.NewRouter()
	v1Rotuer.Get("/healthz", handlerReadiness)
	v1Rotuer.Get("/err", handlerErr)
	v1Rotuer.Post("/users", apiCfg.handlerCreateUser)
	router.Mount("/v1", v1Rotuer)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// registerCors: register basic cors
func registerCors(router *chi.Mux) {
	log.Println("Application Register middleware for cors ")

	// basic cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"Get", "Post", "Put", "Delete", "Options"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}
