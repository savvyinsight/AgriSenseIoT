package handlers

import (
    "encoding/json"
    "log"
    "time"
)

type CommandResponse struct {
    CommandID int       `json:"command_id"`
    Status    string    `json:"status"` // executed, failed
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message,omitempty"`
}

func HandleResponse(deviceID string, payload []byte) {
    log.Printf("Received command response from device %s", deviceID)
    
    var resp CommandResponse
    if err := json.Unmarshal(payload, &resp); err != nil {
        log.Printf("Failed to parse response from device %s: %v", deviceID, err)
        return
    }
    
    // TODO: Update command status in database
    log.Printf("Device %s command %d: %s", deviceID, resp.CommandID, resp.Status)
    if resp.Message != "" {
        log.Printf("  Message: %s", resp.Message)
    }
}
