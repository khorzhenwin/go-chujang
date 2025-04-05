package main

import (
	"log"
	"net/http"
	"time"

	health_api "github.com/khorzhenwin/go-chujang/cmd/api/health-api"
)

type application struct {
	config config
}

type config struct {
	BASE_PATH    string
	ADDRESS      string
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func (app *application) run() error {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         app.config.ADDRESS,
		Handler:      mux,
		WriteTimeout: app.config.writeTimeout,
		ReadTimeout:  app.config.readTimeout,
	}

	mux.HandleFunc("GET "+app.config.BASE_PATH+"/health", health_api.HealthHandler)

	log.Println("Starting server on", app.config.ADDRESS)
	return server.ListenAndServe()
}
