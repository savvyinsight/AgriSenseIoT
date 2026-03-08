package handlers

import (
    "encoding/json"
    "log"
    "time"
)

type TelemetryData struct {
    Timestamp time.Time                `json:"timestamp"`
    Readings  []SensorReading          `json:"readings"`
    Metadata  map[string]interface{}   `json:"metadata,omitempty"`
}

type SensorReading struct {
    Sensor string  `json:"sensor"` // temperature, humidity, etc.
    Value  float64 `json:"value"`
}

func HandleTelemetry(deviceID string, payload []byte) {
    log.Printf("Received telemetry from device %s", deviceID)
    
    var data TelemetryData
    if err := json.Unmarshal(payload, &data); err != nil {
        log.Printf("Failed to parse telemetry from device %s: %v", deviceID, err)
        return
    }
    
    // TODO: Pass to data service for processing
    log.Printf("Device %s sent %d readings", deviceID, len(data.Readings))
    
    for _, reading := range data.Readings {
        log.Printf("  - %s: %.2f", reading.Sensor, reading.Value)
    }
}
