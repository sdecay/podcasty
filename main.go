package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

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
		MaxAge:           300,
	}))
}

func main() {
	err := loadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	router := chi.NewRouter()

	setupCors(router)

	v1Router := chi.NewRouter()

	v1Router.Get("/health", handlerReady)
	v1Router.Get("/error", handlerError)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	// add IP addr from utils.go
	log.Printf("Listening on port %s\n", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
