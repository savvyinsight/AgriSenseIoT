package main

import (
	"log"

	"github.com/savvyinsight/agrisenseiot/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("AgriSenseIoT server starting on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Env)

	// TODO: Add server logic here

	select {} // Block forever
}
