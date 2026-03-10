package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/savvyinsight/agrisenseiot/internal/service/data"
)

type DataHandler struct {
	dataService *data.Service
}

func NewDataHandler(dataService *data.Service) *DataHandler {
	return &DataHandler{
		dataService: dataService,
	}
}

func (h *DataHandler) GetLatest(c *gin.Context) {
	deviceID := c.Param("deviceId")
	sensorType := c.Query("sensor_type")
	if sensorType == "" {
		sensorType = "temperature" // Default
	}

	reading, err := h.dataService.GetLatestReading(deviceID, sensorType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reading)
}

func (h *DataHandler) GetHistorical(c *gin.Context) {
	deviceID := c.Param("deviceId")
	sensorType := c.Query("sensor_type")
	if sensorType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensor_type is required"})
		return
	}

	// Parse time range
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time format"})
			return
		}
	} else {
		start = time.Now().Add(-24 * time.Hour) // Default: last 24h
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time format"})
			return
		}
	} else {
		end = time.Now()
	}

	data, err := h.dataService.GetHistoricalData(deviceID, sensorType, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *DataHandler) GetAggregated(c *gin.Context) {
	deviceID := c.Param("deviceId")
	sensorType := c.Query("sensor_type")
	interval := c.DefaultQuery("interval", "1h")

	if sensorType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensor_type is required"})
		return
	}

	// Parse time range (default: last 7 days)
	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time format"})
			return
		}
	} else {
		start = time.Now().Add(-7 * 24 * time.Hour)
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time format"})
			return
		}
	} else {
		end = time.Now()
	}

	data, err := h.dataService.GetAggregatedData(deviceID, sensorType, start, end, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
