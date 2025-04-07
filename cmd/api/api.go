package main

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/khorzhenwin/go-chujang/docs" // <-- this is required for Swagger to embed docs
	dbconfig "github.com/khorzhenwin/go-chujang/internal/config"
	"github.com/khorzhenwin/go-chujang/internal/db"
	"github.com/khorzhenwin/go-chujang/internal/health"
	"github.com/khorzhenwin/go-chujang/internal/watchlist"
	migration "github.com/khorzhenwin/go-chujang/scripts"
	_ "github.com/swaggo/files"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func (app *application) run() error {
	// 1. Load DB Config
	cfg, err := dbconfig.LoadDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize DB
	conn, err := db.New(cfg.GetFormattedDSN())
	if err != nil {
		log.Fatal(err)
	}
	repo := watchlist.NewRepository(conn)

	// 3. Run Migrations
	println(cfg.DSN)
	migration.RunMigration(cfg.DSN)

	// 4. Setup Router config
	r := chi.NewRouter()
	server := &http.Server{
		Addr:         app.config.ADDRESS,
		Handler:      r,
		WriteTimeout: app.config.writeTimeout,
		ReadTimeout:  app.config.readTimeout,
	}

	// 5. Register all API routes
	r.Route(app.config.BASE_PATH, func(r chi.Router) {
		health.RegisterRoutes(r)
		watchlist.RegisterRoutes(r, repo)
	})

	// 6. Serve Swagger (if generated)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Starting server on", app.config.ADDRESS)
	return server.ListenAndServe()
}
