package sensor

import (
    "errors"
    "fmt"
    "time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	
)

// InitiateSensorNetworkScan performs a network scan to identify active and available sensors.
func InitiateSensorNetworkScan() ([]common.NetworkScanResult, error) {
    // Perform network scan, simulate retrieval of active sensors
    scanResults := []common.NetworkScanResult{}

    sensors, err := ledger.GetAllSensors()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve sensors: %v", err)
    }

    for _, sensor := range sensors {
        // Check connectivity and activity status
        isConnected := sensor.IsConnected
        healthStatus := "Operational"
        if sensor.IsFaulty {
            healthStatus = "Faulty"
        }

        scanResults = append(scanResults, common.NetworkScanResult{
            SensorID:     sensor.ID,
            IsConnected:  isConnected,
            HealthStatus: healthStatus,
            LastChecked:  time.Now(),
        })
    }

    return scanResults, nil
}

// RegisterSensorEventHandler registers an event handler to respond to specific sensor events.
func RegisterSensorEventHandler(sensorID string, handler common.EventHandler) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    // Encrypt the handler details before saving
    handler.EncryptedDetails, err = encryption.EncryptData(handler.Details)
    if err != nil {
        return fmt.Errorf("failed to encrypt event handler details: %v", err)
    }
    handler.RegisteredAt = time.Now()

    sensor.EventHandlers = append(sensor.EventHandlers, handler)
    return ledger.UpdateSensor(sensor)
}

// RemoveSensorEventHandler removes a previously registered event handler from a sensor.
func RemoveSensorEventHandler(sensorID string, handlerID string) error {
    sensor, err := ledger.GetSensor(sensorID)
    if err != nil {
        return fmt.Errorf("sensor not found: %v", err)
    }

    found := false
    for i, handler := range sensor.EventHandlers {
        if handler.ID == handlerID {
            // Remove handler and update ledger
            sensor.EventHandlers = append(sensor.EventHandlers[:i], sensor.EventHandlers[i+1:]...)
            found = true
            break
        }
    }

    if !found {
        return errors.New("event handler not found")
    }

    return ledger.UpdateSensor(sensor)
}
