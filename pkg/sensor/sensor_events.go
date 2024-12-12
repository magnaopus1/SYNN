package sensor

import (
    "errors"
    "fmt"
    "time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	
)

// TriggerSensorEvent triggers an event on a sensor, recording it to the ledger.
func TriggerSensorEvent(sensorID string, event common.SensorEvent) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    event.Timestamp = time.Now()
    sensor.Events = append(sensor.Events, event)
    return ledger.UpdateSensor(sensor)
}

// ClearSensorData clears all non-essential data stored in a sensor.
func ClearSensorData(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Data = nil
    return ledger.UpdateSensor(sensor)
}

// StoreSensorData securely stores data from a sensor.
func StoreSensorData(sensorID string, data common.SensorData) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    encryptedData, err := encryption.EncryptData(data.Content)
    if err != nil {
        return fmt.Errorf("data encryption failed: %v", err)
    }

    data.Content = encryptedData
    sensor.Data = append(sensor.Data, data)
    return ledger.UpdateSensor(sensor)
}

// LoadSensorPresets loads predefined settings into a sensor.
func LoadSensorPresets(sensorID string, presets common.SensorPresets) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Presets = presets
    return ledger.UpdateSensor(sensor)
}

// DefineSensorDataRetention sets a data retention policy for a sensor.
func DefineSensorDataRetention(sensorID string, retentionRule common.DataRetentionRule) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataRetentionRule = retentionRule
    return ledger.UpdateSensor(sensor)
}

// MonitorSensorHealth continuously monitors the health of a sensor.
func MonitorSensorHealth(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    if sensor.IsFaulty {
        return fmt.Errorf("sensor %s requires maintenance", sensorID)
    }

    sensor.LastHealthCheck = time.Now()
    return ledger.UpdateSensor(sensor)
}

// ActivateSensorMode activates a specific operational mode on a sensor.
func ActivateSensorMode(sensorID string, mode common.SensorMode) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.ActiveMode = mode
    return ledger.UpdateSensor(sensor)
}

// DeactivateSensorMode deactivates the currently active operational mode on a sensor.
func DeactivateSensorMode(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.ActiveMode = common.SensorMode{}
    return ledger.UpdateSensor(sensor)
}

// SetSensorOperationTimeout sets an operational timeout for a sensor.
func SetSensorOperationTimeout(sensorID string, timeout time.Duration) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.OperationTimeout = time.Now().Add(timeout)
    return ledger.UpdateSensor(sensor)
}

// ExtendSensorOperationTimeout extends the current operational timeout for a sensor.
func ExtendSensorOperationTimeout(sensorID string, additionalTime time.Duration) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.OperationTimeout = sensor.OperationTimeout.Add(additionalTime)
    return ledger.UpdateSensor(sensor)
}

// RestrictSensorAccess temporarily restricts access to a sensor.
func RestrictSensorAccess(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.AccessRestricted = true
    return ledger.UpdateSensor(sensor)
}

// LiftSensorAccessRestriction removes the access restriction from a sensor.
func LiftSensorAccessRestriction(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.AccessRestricted = false
    return ledger.UpdateSensor(sensor)
}

// CheckSensorBattery checks the current battery level of a sensor.
func CheckSensorBattery(sensorID string) (float64, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return 0, fmt.Errorf("sensor not found: %v", err)
    }

    return sensor.BatteryLevel, nil
}

// ScheduleSensorMaintenance schedules maintenance for a sensor.
func ScheduleSensorMaintenance(sensorID string, maintenanceTime time.Time) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.NextMaintenance = maintenanceTime
    return ledger.UpdateSensor(sensor)
}

// RetrieveSensorLogs retrieves operational logs for a sensor.
func RetrieveSensorLogs(sensorID string, startTime, endTime time.Time) ([]common.SensorLog, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return nil, fmt.Errorf("sensor not found: %v", err)
    }

    return ledger.FetchSensorLogs(sensorID, startTime, endTime)
}

// SetSensorSecurityLevel sets the security level for a sensor.
func SetSensorSecurityLevel(sensorID string, securityLevel string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.SecurityLevel = securityLevel
    return ledger.UpdateSensor(sensor)
}
