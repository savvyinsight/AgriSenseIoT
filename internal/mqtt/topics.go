package mqtt

import "fmt"

// Topic patterns
const (
    // Device → Platform
    TelemetryTopic = "device/+/telemetry"  // device/{id}/telemetry
    HeartbeatTopic = "device/+/heartbeat"  // device/{id}/heartbeat
    ResponseTopic  = "device/+/response"   // device/{id}/response
    
    // Platform → Device
    CommandTopic = "device/%s/commands"    // device/{id}/commands
    ConfigTopic  = "device/%s/config"      // device/{id}/config
)

// GetCommandTopic returns the command topic for a specific device
func GetCommandTopic(deviceID string) string {
    return fmt.Sprintf("device/%s/commands", deviceID)
}

// GetConfigTopic returns the config topic for a specific device
func GetConfigTopic(deviceID string) string {
    return fmt.Sprintf("device/%s/config", deviceID)
}

// ExtractDeviceIDFromTopic extracts device ID from subscription topics
func ExtractDeviceIDFromTopic(topic string, prefix string) string {
    // Example: device/esp32_001/telemetry -> esp32_001
    return topic[len(prefix)+1 : len(topic)-len("/telemetry")-1]
}
