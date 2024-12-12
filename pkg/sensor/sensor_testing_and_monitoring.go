package sensor

import (
    "errors"
    "fmt"
    "time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	
)

// PerformSensorSelfTest runs a diagnostic self-test on the sensor to check for functionality.
func PerformSensorSelfTest(sensorID string) (common.SelfTestReport, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return common.SelfTestReport{}, fmt.Errorf("sensor not found: %v", err)
    }

    report := common.SelfTestReport{
        SensorID:  sensorID,
        Status:    "OK",
        Timestamp: time.Now(),
    }

    if sensor.IsFaulty {
        report.Status = "Faulty"
        report.Details = "Sensor failed self-test and requires maintenance."
    } else {
        report.Details = "All systems operational."
    }

    sensor.LastSelfTestReport = report
    return report, ledger.UpdateSensor(sensor)
}

// ScheduleSensorDataCollection sets a schedule for automatic data collection.
func ScheduleSensorDataCollection(sensorID string, interval time.Duration) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataCollectionInterval = interval
    return ledger.UpdateSensor(sensor)
}

// ConfigureSensorDataFormat sets the data format for a sensor.
func ConfigureSensorDataFormat(sensorID string, format common.DataFormat) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataFormat = format
    return ledger.UpdateSensor(sensor)
}

// EnableSensorDataCompression activates data compression for storage efficiency.
func EnableSensorDataCompression(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataCompressionEnabled = true
    return ledger.UpdateSensor(sensor)
}

// DisableSensorDataCompression deactivates data compression.
func DisableSensorDataCompression(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.DataCompressionEnabled = false
    return ledger.UpdateSensor(sensor)
}

// TrackSensorResponseTime monitors the response time of a sensor.
func TrackSensorResponseTime(sensorID string) (time.Duration, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return 0, fmt.Errorf("sensor not found: %v", err)
    }

    responseTime := time.Since(sensor.LastResponseCheck)
    sensor.LastResponseTime = responseTime
    sensor.LastResponseCheck = time.Now()
    return responseTime, ledger.UpdateSensor(sensor)
}

// LogSensorCalibration records calibration events for the sensor.
func LogSensorCalibration(sensorID string, calibration common.CalibrationLog) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    calibration.Timestamp = time.Now()
    sensor.CalibrationLogs = append(sensor.CalibrationLogs, calibration)
    return ledger.UpdateSensor(sensor)
}

// MonitorSensorPowerConsumption tracks the power consumption of the sensor.
func MonitorSensorPowerConsumption(sensorID string) (float64, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return 0, fmt.Errorf("sensor not found: %v", err)
    }

    powerConsumption := sensor.CalculateCurrentPowerConsumption()
    sensor.LastPowerConsumption = powerConsumption
    return powerConsumption, ledger.UpdateSensor(sensor)
}

// DefineSensorAlerts sets up alert criteria for the sensor.
func DefineSensorAlerts(sensorID string, alerts []common.Alert) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Alerts = alerts
    return ledger.UpdateSensor(sensor)
}

// ClearSensorAlerts clears all alert settings on the sensor.
func ClearSensorAlerts(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.Alerts = nil
    return ledger.UpdateSensor(sensor)
}

// EnableSensorAutoShutdown activates auto-shutdown for low battery or idle states.
func EnableSensorAutoShutdown(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.AutoShutdownEnabled = true
    return ledger.UpdateSensor(sensor)
}

// DisableSensorAutoShutdown deactivates auto-shutdown mode.
func DisableSensorAutoShutdown(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.AutoShutdownEnabled = false
    return ledger.UpdateSensor(sensor)
}

// SetSensorFailsafeMode enables a failsafe mode for critical situations.
func SetSensorFailsafeMode(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.FailsafeModeEnabled = true
    return ledger.UpdateSensor(sensor)
}

// ClearSensorFailsafeMode disables the failsafe mode.
func ClearSensorFailsafeMode(sensorID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    sensor.FailsafeModeEnabled = false
    return ledger.UpdateSensor(sensor)
}

// QuerySensorUptime calculates the uptime of a sensor since the last reboot.
func QuerySensorUptime(sensorID string) (time.Duration, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return 0, fmt.Errorf("sensor not found: %v", err)
    }

    uptime := time.Since(sensor.LastReboot)
    return uptime, nil
}

// VerifySensorDataIntegrity checks the integrity of the sensor’s stored data.
func VerifySensorDataIntegrity(sensorID string) (bool, error) {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return false, fmt.Errorf("sensor not found: %v", err)
    }

    for _, data := range sensor.Data {
        decryptedContent, err := encryption.DecryptData(data.Content)
        if err != nil {
            return false, fmt.Errorf("data integrity check failed for sensor %s: %v", sensorID, err)
        }
        if !encryption.VerifyChecksum(decryptedContent, data.Checksum) {
            return false, errors.New("data checksum mismatch detected")
        }
    }

    return true, nil
}
