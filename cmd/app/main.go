package main

import (
	"log"

	"github.com/lovelydaemon/url-shortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
