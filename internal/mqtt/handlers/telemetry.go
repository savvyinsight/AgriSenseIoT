package handlers

import (
	"log"

	"github.com/savvyinsight/agrisenseiot/internal/service/data"
)

var dataService *data.Service

// Init sets the data service for handlers
func Init(ds *data.Service) {
	dataService = ds
}

func HandleTelemetry(deviceID string, payload []byte) {
	if dataService == nil {
		log.Println("ERROR: Data service not initialized")
		return
	}

	if err := dataService.ProcessTelemetry(deviceID, payload); err != nil {
		log.Printf("Failed to process telemetry from device %s: %v", deviceID, err)
	}
}

// type TelemetryData struct {
//     Timestamp time.Time                `json:"timestamp"`
//     Readings  []SensorReading          `json:"readings"`
//     Metadata  map[string]interface{}   `json:"metadata,omitempty"`
// }

// type SensorReading struct {
//     Sensor string  `json:"sensor"` // temperature, humidity, etc.
//     Value  float64 `json:"value"`
// }

// func HandleTelemetry(deviceID string, payload []byte) {
//     log.Printf("Received telemetry from device %s", deviceID)

//     var data TelemetryData
//     if err := json.Unmarshal(payload, &data); err != nil {
//         log.Printf("Failed to parse telemetry from device %s: %v", deviceID, err)
//         return
//     }

//     // TODO: Pass to data service for processing
//     log.Printf("Device %s sent %d readings", deviceID, len(data.Readings))

//     for _, reading := range data.Readings {
//         log.Printf("  - %s: %.2f", reading.Sensor, reading.Value)
//     }
// }
