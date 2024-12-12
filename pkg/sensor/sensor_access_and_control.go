package sensor

import (
    "errors"
    "fmt"
    "time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// RemoveSensorSecurityLevel removes a specific security level from a sensor.
func RemoveSensorSecurityLevel(sensorID string, securityLevel string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    if _, exists := sensor.SecurityLevels[securityLevel]; exists {
        delete(sensor.SecurityLevels, securityLevel)
    } else {
        return errors.New("specified security level not found")
    }
    
    return ledger.UpdateSensor(sensor)
}

// DefineSensorDataProcessingRule adds a data processing rule for a sensor.
func DefineSensorDataProcessingRule(sensorID string, rule common.DataProcessingRule) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataProcessingRules = append(sensor.DataProcessingRules, rule)
    return ledger.UpdateSensor(sensor)
}

// RemoveSensorDataProcessingRule removes a data processing rule from a sensor.
func RemoveSensorDataProcessingRule(sensorID string, ruleID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    for i, rule := range sensor.DataProcessingRules {
        if rule.ID == ruleID {
            sensor.DataProcessingRules = append(sensor.DataProcessingRules[:i], sensor.DataProcessingRules[i+1:]...)
            return ledger.UpdateSensor(sensor)
        }
    }
    return errors.New("data processing rule not found")
}

// EnableSensorRealTimeMonitoring activates real-time monitoring on a sensor.
func EnableSensorRealTimeMonitoring(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.RealTimeMonitoringEnabled = true
    return ledger.UpdateSensor(sensor)
}

// DisableSensorRealTimeMonitoring deactivates real-time monitoring on a sensor.
func DisableSensorRealTimeMonitoring(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.RealTimeMonitoringEnabled = false
    return ledger.UpdateSensor(sensor)
}

// QuerySensorDataHistory retrieves the historical data for a sensor.
func QuerySensorDataHistory(sensorID string, startTime, endTime time.Time) ([]common.SensorData, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return nil, fmt.Errorf("sensor not found: %v", err)
    }

    return ledger.FetchSensorDataHistory(sensorID, startTime, endTime)
}

// PerformSensorSoftwareUpdate applies a software update to a sensor.
func PerformSensorSoftwareUpdate(sensorID string, updateDetails common.SoftwareUpdate) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.LastUpdate = updateDetails
    sensor.UpdateTimestamp = time.Now()
    return ledger.UpdateSensor(sensor)
}

// SetSensorLocation sets the physical location of a sensor.
func SetSensorLocation(sensorID string, location common.Location) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Location = location
    return ledger.UpdateSensor(sensor)
}

// GetSensorLocation retrieves the location details of a sensor.
func GetSensorLocation(sensorID string) (common.Location, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return common.Location{}, fmt.Errorf("sensor not found: %v", err)
    }

    return sensor.Location, nil
}

// DefineSensorPriority sets the operational priority level for a sensor.
func DefineSensorPriority(sensorID string, priorityLevel int) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.PriorityLevel = priorityLevel
    return ledger.UpdateSensor(sensor)
}

// ResetSensorPriority resets the priority level of a sensor to its default.
func ResetSensorPriority(sensorID string) error {
    return DefineSensorPriority(sensorID, common.DefaultPriority)
}

// CheckSensorConnectivity verifies the connectivity status of a sensor.
func CheckSensorConnectivity(sensorID string) (bool, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return false, fmt.Errorf("sensor not found: %v", err)
    }

    return sensor.IsConnected, nil
}

// LogSensorDisconnection logs the disconnection of a sensor.
func LogSensorDisconnection(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.IsConnected = false
    sensor.DisconnectionTimestamp = time.Now()
    return ledger.UpdateSensor(sensor)
}

// SetSensorAccessPermissions grants access permissions to a user for a sensor.
func SetSensorAccessPermissions(sensorID string, userID string, permissions []string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    access := common.AccessPermission{
        UserID:      userID,
        Permissions: permissions,
        GrantedAt:   time.Now(),
    }

    sensor.AccessPermissions = append(sensor.AccessPermissions, access)
    return ledger.UpdateSensor(sensor)
}

// RevokeSensorAccessPermissions revokes access permissions for a user from a sensor.
func RevokeSensorAccessPermissions(sensorID string, userID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    for i, access := range sensor.AccessPermissions {
        if access.UserID == userID {
            sensor.AccessPermissions = append(sensor.AccessPermissions[:i], sensor.AccessPermissions[i+1:]...)
            return ledger.UpdateSensor(sensor)
        }
    }
    return errors.New("user access permission not found")
}
