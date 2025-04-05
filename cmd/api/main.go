package main

import (
	"log"
	"time"
)

func main() {
	cfg := config{
		BASE_PATH:    "/api/v1",
		ADDRESS:      ":8080",
		writeTimeout: time.Second * 10,
		readTimeout:  time.Second * 5,
	}

	app := &application{
		config: cfg,
	}

	log.Fatal(app.run())
}
