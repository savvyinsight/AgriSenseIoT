package postgres

import (
	"database/sql"
	"time"

	"github.com/savvyinsight/agrisenseiot/internal/domain"
)

type AlertRepository struct {
	DB *sql.DB
}

func (r *AlertRepository) Create(alert *domain.Alert) error {
	query := `
        INSERT INTO alerts (
            rule_id, device_id, sensor_value, message, severity, 
            status, triggered_at, metadata
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `

	err := r.DB.QueryRow(
		query,
		alert.RuleID,
		alert.DeviceID,
		alert.SensorValue,
		alert.Message,
		alert.Severity,
		alert.Status,
		alert.TriggeredAt,
		alert.Metadata,
	).Scan(&alert.ID)

	return err
}

func (r *AlertRepository) GetActive() ([]domain.Alert, error) {
	query := `
        SELECT id, rule_id, device_id, sensor_value, message, severity, 
               status, triggered_at, acknowledged_at, resolved_at, metadata
        FROM alerts 
        WHERE status IN ('triggered', 'acknowledged')
        ORDER BY triggered_at DESC
    `

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.RuleID,
			&alert.DeviceID,
			&alert.SensorValue,
			&alert.Message,
			&alert.Severity,
			&alert.Status,
			&alert.TriggeredAt,
			&alert.AcknowledgedAt,
			&alert.ResolvedAt,
			&alert.Metadata,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) GetByDeviceID(deviceID int) ([]domain.Alert, error) {
	query := `
        SELECT id, rule_id, device_id, sensor_value, message, severity, 
               status, triggered_at, acknowledged_at, resolved_at, metadata
        FROM alerts 
        WHERE device_id = $1
        ORDER BY triggered_at DESC
    `

	rows, err := r.DB.Query(query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.RuleID,
			&alert.DeviceID,
			&alert.SensorValue,
			&alert.Message,
			&alert.Severity,
			&alert.Status,
			&alert.TriggeredAt,
			&alert.AcknowledgedAt,
			&alert.ResolvedAt,
			&alert.Metadata,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) GetByRuleID(ruleID int) ([]domain.Alert, error) {
	query := `
        SELECT id, rule_id, device_id, sensor_value, message, severity, 
               status, triggered_at, acknowledged_at, resolved_at, metadata
        FROM alerts 
        WHERE rule_id = $1
        ORDER BY triggered_at DESC
    `

	rows, err := r.DB.Query(query, ruleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.RuleID,
			&alert.DeviceID,
			&alert.SensorValue,
			&alert.Message,
			&alert.Severity,
			&alert.Status,
			&alert.TriggeredAt,
			&alert.AcknowledgedAt,
			&alert.ResolvedAt,
			&alert.Metadata,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *AlertRepository) Acknowledge(id int) error {
	query := `UPDATE alerts SET status = $1, acknowledged_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, domain.AlertStatusAcknowledged, time.Now(), id)
	return err
}

func (r *AlertRepository) Resolve(id int) error {
	query := `UPDATE alerts SET status = $1, resolved_at = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, domain.AlertStatusResolved, time.Now(), id)
	return err
}

func (r *AlertRepository) List(limit, offset int) ([]domain.Alert, int64, error) {
	query := `
        SELECT id, rule_id, device_id, sensor_value, message, severity, 
               status, triggered_at, acknowledged_at, resolved_at, metadata
        FROM alerts 
        ORDER BY triggered_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.RuleID,
			&alert.DeviceID,
			&alert.SensorValue,
			&alert.Message,
			&alert.Severity,
			&alert.Status,
			&alert.TriggeredAt,
			&alert.AcknowledgedAt,
			&alert.ResolvedAt,
			&alert.Metadata,
		)
		if err != nil {
			return nil, 0, err
		}
		alerts = append(alerts, alert)
	}

	var total int64
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM alerts`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}
