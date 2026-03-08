package mqtt

import (
	"log"

	"github.com/savvyinsight/agrisenseiot/internal/mqtt/handlers"
	"github.com/savvyinsight/agrisenseiot/internal/service/data"
)

type Service struct {
	client      *Client
	dataService *data.Service
}

func NewService(cfg Config, dataService *data.Service) (*Service, error) {
	// Initialize handlers with the data service
	handlers.Init(dataService)
	handlers := &Handlers{
		TelemetryHandler: handlers.HandleTelemetry,
		HeartbeatHandler: handlers.HandleHeartbeat,
		ResponseHandler:  handlers.HandleResponse,
	}

	client, err := NewClient(cfg, handlers)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

func (s *Service) Start() error {
	log.Println("Starting MQTT service...")
	return s.client.Subscribe()
}

func (s *Service) Stop() {
	if s.client != nil {
		s.client.Disconnect()
	}
}

func (s *Service) SendCommand(deviceID string, command []byte) error {
	return s.client.PublishCommand(deviceID, command)
}

func (s *Service) SendConfig(deviceID string, config []byte) error {
	return s.client.PublishConfig(deviceID, config)
}
