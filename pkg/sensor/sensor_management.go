package sensor

import (
    "errors"
    "fmt"
    "time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	
)

// ConnectSensor connects a sensor to the network.
func ConnectSensor(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.IsConnected = true
    sensor.LastConnected = time.Now()
    return ledger.UpdateSensor(sensor)
}

// DisconnectSensor disconnects a sensor from the network.
func DisconnectSensor(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.IsConnected = false
    sensor.LastDisconnected = time.Now()
    return ledger.UpdateSensor(sensor)
}

// ReadSensorData retrieves the latest data from a sensor.
func ReadSensorData(sensorID string) (common.SensorData, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return common.SensorData{}, fmt.Errorf("sensor not found: %v", err)
    }

    if len(sensor.Data) == 0 {
        return common.SensorData{}, errors.New("no data available for this sensor")
    }

    latestData := sensor.Data[len(sensor.Data)-1]
    decryptedData, err := encryption.DecryptData(latestData.Content)
    if err != nil {
        return common.SensorData{}, fmt.Errorf("data decryption failed: %v", err)
    }

    latestData.Content = decryptedData
    return latestData, nil
}

// WriteSensorConfiguration updates the configuration for a sensor.
func WriteSensorConfiguration(sensorID string, config common.SensorConfiguration) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Configuration = config
    return ledger.UpdateSensor(sensor)
}

// CalibrateSensor calibrates a sensor to ensure accurate readings.
func CalibrateSensor(sensorID string, calibrationParams common.CalibrationParams) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Calibration = calibrationParams
    sensor.LastCalibration = time.Now()
    return ledger.UpdateSensor(sensor)
}

// CheckSensorStatus verifies if a sensor is operational.
func CheckSensorStatus(sensorID string) (string, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return "", fmt.Errorf("sensor not found: %v", err)
    }

    if !sensor.IsConnected {
        return "Disconnected", nil
    }

    if sensor.IsFaulty {
        return "Faulty", nil
    }

    return "Operational", nil
}

// ListConnectedSensors lists all sensors currently connected to the network.
func ListConnectedSensors() ([]string, error) {
    sensors, err := ledger.GetAllSensors()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve sensors: %v", err)
    }

    var connectedSensors []string
    for _, sensor := range sensors {
        if sensor.IsConnected {
            connectedSensors = append(connectedSensors, sensor.ID)
        }
    }
    return connectedSensors, nil
}

// RebootSensor restarts a sensor to resolve potential issues.
func RebootSensor(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    if !sensor.IsConnected {
        return errors.New("sensor must be connected to reboot")
    }

    sensor.LastReboot = time.Now()
    return ledger.UpdateSensor(sensor)
}

// SetSensorFrequency configures the frequency at which a sensor collects data.
func SetSensorFrequency(sensorID string, frequency int) error {
    if frequency <= 0 {
        return errors.New("frequency must be positive")
    }

    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Frequency = frequency
    return ledger.UpdateSensor(sensor)
}

// GetSensorFrequency retrieves the current data collection frequency of a sensor.
func GetSensorFrequency(sensorID string) (int, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return 0, fmt.Errorf("sensor not found: %v", err)
    }

    return sensor.Frequency, nil
}

// InitiateSensorDiagnostics runs diagnostics on a sensor to check for issues.
func InitiateSensorDiagnostics(sensorID string) (common.DiagnosticReport, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return common.DiagnosticReport{}, fmt.Errorf("sensor not found: %v", err)
    }

    report := common.DiagnosticReport{
        SensorID:   sensorID,
        Status:     "OK",
        CheckedAt:  time.Now(),
        Details:    "All parameters within normal range.",
    }

    if sensor.IsFaulty {
        report.Status = "Faulty"
        report.Details = "Fault detected in sensor. Requires maintenance."
    }

    sensor.LastDiagnostic = report
    return report, ledger.UpdateSensor(sensor)
}

// LogSensorEvent records an event or anomaly detected in the sensor.
func LogSensorEvent(sensorID string, event common.EventLog) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    event.Timestamp = time.Now()
    sensor.EventLogs = append(sensor.EventLogs, event)
    return ledger.UpdateSensor(sensor)
}

// SetSensorDataThreshold defines the threshold values for sensor data parameters.
func SetSensorDataThreshold(sensorID string, threshold common.DataThreshold) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Threshold = threshold
    return ledger.UpdateSensor(sensor)
}

// ResetSensorDataThreshold resets a sensor’s data threshold to default values.
func ResetSensorDataThreshold(sensorID string) error {
    defaultThreshold := common.DataThreshold{}  // Assume default is zero or unset values
    return SetSensorDataThreshold(sensorID, defaultThreshold)
}

// EnableSensorNotifications enables notifications for sensor alerts.
func EnableSensorNotifications(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.NotificationsEnabled = true
    return ledger.UpdateSensor(sensor)
}

// DisableSensorNotifications disables notifications for sensor alerts.
func DisableSensorNotifications(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.NotificationsEnabled = false
    return ledger.UpdateSensor(sensor)
}
