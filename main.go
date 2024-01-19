package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sdecay/podcasty/internal/database"
)

// TODO: cap
type apiConfig struct {
	DB               *database.Queries
	serverPort       string
	scrapeMaxThreads int
	dbUrl            string
	scrapeDelay      time.Duration
}

func loadEnvironment() error {
	return godotenv.Load(".env")
}

// cors is something i'm not looking forward to learning about.
// seems dreadful.
func setupCors(router *chi.Mux) {
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           30,
	}))
}

// TODO: break it up at some point.  too much shit in here
func main() {
	config := apiConfig{}

	err := loadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	config.serverPort = os.Getenv("PORT")
	config.dbUrl = os.Getenv("DB_URL")
	config.scrapeMaxThreads, _ = strconv.Atoi(os.Getenv("SCRAPE_MAX_THREADS"))
	scrapeDelay, _ := strconv.Atoi(os.Getenv("SCRAPE_DELAY")) // ugh
	config.scrapeDelay = time.Second * time.Duration(scrapeDelay)

	conn, err := sql.Open("postgres", config.dbUrl)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}

	config.DB = database.New(conn)

	go scrape(config.DB, config.scrapeMaxThreads, config.scrapeDelay)

	router := chi.NewRouter()
	setupCors(router)

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handlerReady)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", config.handlerCreateUser)
	v1Router.Get("/users", config.middlewareAuth(config.handlerGetUser))

	v1Router.Get("/feeds", config.handlerGetFeeds)
	v1Router.Post("/feeds", config.middlewareAuth(config.handlerCreateFeed))

	v1Router.Post("/follow", config.middlewareAuth(config.handlerFollowFeed))
	v1Router.Get("/follow", config.middlewareAuth(config.handlerGetFollowed))
	v1Router.Delete("/follow/{followID}", config.middlewareAuth(config.handlerDeleteFollow))

	v1Router.Get("/latest", config.middlewareAuth(config.handlerGetUsersPosts))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + config.serverPort,
	}

	log.Printf("Listening at %s:%s\n", GetLocalIP(), config.serverPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
