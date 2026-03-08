package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/savvyinsight/agrisenseiot/internal/config"
	"github.com/savvyinsight/agrisenseiot/internal/mqtt"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create MQTT service
	mqttCfg := mqtt.Config{
		Broker:   cfg.MQTTBroker,
		ClientID: "agrisense-mqtt-handler",
		Username: cfg.MQTTUsername,
		Password: cfg.MQTTPassword,
	}

	service, err := mqtt.NewService(mqttCfg)
	if err != nil {
		log.Fatalf("Failed to create MQTT service: %v", err)
	}

	// Start MQTT service
	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start MQTT service: %v", err)
	}

	log.Println("MQTT handler started. Waiting for messages...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down MQTT handler...")
	service.Stop()
}
