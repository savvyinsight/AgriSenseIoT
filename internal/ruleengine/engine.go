package ruleengine

import (
	"log"
	"sync"
	"time"

	"github.com/savvyinsight/agrisenseiot/internal/domain"
)

type Engine struct {
	rules      map[int]*domain.AlertRule
	rulesMutex sync.RWMutex
	evaluator  *Evaluator
	alertSvc   domain.AlertRepository
	ruleRepo   domain.AlertRuleRepository
	deviceRepo domain.DeviceRepository
	stopChan   chan struct{}
}

func NewEngine(ruleRepo domain.AlertRuleRepository,
	alertSvc domain.AlertRepository,
	deviceRepo domain.DeviceRepository) *Engine {
	return &Engine{
		rules:      make(map[int]*domain.AlertRule),
		evaluator:  NewEvaluator(),
		ruleRepo:   ruleRepo,
		alertSvc:   alertSvc,
		deviceRepo: deviceRepo,
		stopChan:   make(chan struct{}),
	}
}

func (e *Engine) Start() error {
	log.Println("Starting rule engine...")

	// Load rules initially
	if err := e.loadRules(); err != nil {
		return err
	}

	// Start background rule refresh (every 5 minutes)
	go e.refreshRulesPeriodically()

	return nil
}

func (e *Engine) Stop() {
	close(e.stopChan)
}

func (e *Engine) loadRules() error {
	rules, err := e.ruleRepo.GetEnabledRules()
	if err != nil {
		return err
	}

	e.rulesMutex.Lock()
	defer e.rulesMutex.Unlock()

	e.rules = make(map[int]*domain.AlertRule)
	for _, rule := range rules {
		e.rules[rule.ID] = &rule
	}

	log.Printf("Loaded %d enabled rules", len(e.rules))
	return nil
}

func (e *Engine) refreshRulesPeriodically() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := e.loadRules(); err != nil {
				log.Printf("Failed to refresh rules: %v", err)
			}
		case <-e.stopChan:
			return
		}
	}
}

func (e *Engine) Evaluate(data *domain.SensorData) {
	e.rulesMutex.RLock()
	defer e.rulesMutex.RUnlock()

	for _, rule := range e.rules {
		// Skip if rule is for a different device (unless rule applies to all devices)
		if rule.DeviceID != nil && *rule.DeviceID != 0 {
			// Need to map device_id string to device ID - this requires deviceRepo
			// For now, we'll handle this in the evaluator with more context
			continue
		}

		// Skip if rule is for different sensor type
		if rule.SensorTypeID != e.getSensorTypeID(data.SensorType) {
			continue
		}

		// Evaluate the rule
		if e.evaluator.Evaluate(rule, data) {
			e.triggerAlert(rule, data)
		}
	}
}

func (e *Engine) getSensorTypeID(sensorType string) int {
	// This should be cached/mapped from sensor_type_repo
	// For now, hardcode common types
	switch sensorType {
	case "temperature":
		return 1
	case "humidity":
		return 2
	case "soil_moisture":
		return 3
	case "light_intensity":
		return 4
	default:
		return 0
	}
}

func (e *Engine) triggerAlert(rule *domain.AlertRule, data *domain.SensorData) {
	// Get device ID from database using the string device_id
	device, err := e.deviceRepo.GetByDeviceID(data.DeviceID)
	if err != nil {
		log.Printf("Failed to find device %s: %v", data.DeviceID, err)
		return
	}

	alert := &domain.Alert{
		RuleID:      rule.ID,
		DeviceID:    device.ID, // Use the retrieved device ID
		SensorValue: data.Value,
		Message:     e.evaluator.FormatMessage(rule, data),
		Severity:    rule.Severity,
		Status:      domain.AlertStatusTriggered,
		TriggeredAt: time.Now(),
		Metadata: map[string]interface{}{
			"sensor_type": data.SensorType,
			"value":       data.Value,
			"rule_name":   rule.Name,
		},
	}

	// Save alert
	if err := e.alertSvc.Create(alert); err != nil {
		log.Printf("Failed to save alert: %v", err)
		return
	}

	log.Printf("🚨 Alert triggered: %s (rule: %s)", alert.Message, rule.Name)

	// TODO: Send WebSocket notification
	// TODO: Send email notification
}
