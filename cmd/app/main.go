package main

import (
	"log"

	"github.com/ElOtro/auction-go/config"
	"github.com/ElOtro/auction-go/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
