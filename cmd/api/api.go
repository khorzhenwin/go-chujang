package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/khorzhenwin/go-chujang/docs" // <-- this is required for Swagger to embed docs
	healthapi "github.com/khorzhenwin/go-chujang/internal/health"
	watchlistapi "github.com/khorzhenwin/go-chujang/internal/watchlist"
	_ "github.com/swaggo/files"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func (app *application) run() error {
	r := chi.NewRouter()

	server := &http.Server{
		Addr:         app.config.ADDRESS,
		Handler:      r,
		WriteTimeout: app.config.writeTimeout,
		ReadTimeout:  app.config.readTimeout,
	}

	// Register all API routes
	r.Route(app.config.BASE_PATH, func(r chi.Router) {
		healthapi.RegisterRoutes(r)
		watchlistapi.RegisterRoutes(r)
	})

	// Serve Swagger (if generated)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Starting server on", app.config.ADDRESS)
	return server.ListenAndServe()
}
