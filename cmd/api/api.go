package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/khorzhenwin/go-chujang/docs" // <-- this is required for Swagger to embed docs
	applicationConfig "github.com/khorzhenwin/go-chujang/internal/config"
	"github.com/khorzhenwin/go-chujang/internal/db"
	"github.com/khorzhenwin/go-chujang/internal/health"
	"github.com/khorzhenwin/go-chujang/internal/ticker-price"
	"github.com/khorzhenwin/go-chujang/internal/watchlist"
	_ "github.com/swaggo/files"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func (app *application) run() error {
	// 1. Load Configs
	_ = godotenv.Load() // Loads from .env file

	dbCfg, dbErr := applicationConfig.LoadDBConfig()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	vantageCfg, vErr := applicationConfig.LoadVantageConfig()
	if vErr != nil {
		log.Fatal(vErr)
	}

	// 2. Initialize DB
	conn, err := db.New(dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Run Migrations & initialize Repository
	if err := conn.AutoMigrate(&watchlist.Ticker{}); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	watchlistRepo := watchlist.NewRepository(conn)
	watchlistService := watchlist.NewService(watchlistRepo)
	tickerPriceService := ticker_price.NewService(watchlistService, vantageCfg)

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
		watchlist.RegisterRoutes(r, watchlistService)
		ticker_price.RegisterRoutes(r, tickerPriceService)
	})

	// 6. Serve Swagger (if generated)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Println("Starting server on", app.config.ADDRESS)
	return server.ListenAndServe()
}
