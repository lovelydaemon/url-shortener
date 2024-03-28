package main

import (
	"log"

	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/app"
)

func main() {
	cfg := config.NewConfig()

	log.Fatal(app.Run(cfg))
}
