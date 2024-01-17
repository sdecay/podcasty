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
	"github.com/sdecay/podcasty/internal/database"
)

type apiConfig struct {
	DB         *database.Queries
	serverPort string
	dbUrl      string
}

func loadEnvironment() error {
	err := godotenv.Load(".env")
	return err
}

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
	// testes testes 1, 2... 3?
	// feed, err := UrlToRssFeed("https://www.sevatimassage.com/blog?format=rss")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

	config := apiConfig{}

	err := loadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	config.serverPort = os.Getenv("PORT")
	config.dbUrl = os.Getenv("DB_URL")

	conn, err := sql.Open("postgres", config.dbUrl)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}

	config.DB = database.New(conn)

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
