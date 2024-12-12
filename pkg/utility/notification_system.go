package utility

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var alertThreshold int
var alertHistory []string
var alertHandlers sync.Map
var pushNotificationsEnabled bool
var customAlerts sync.Map

// SendNotification: Sends an encrypted notification message to a recipient
func SendNotification(recipientID, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    err = common.ledger.SendNotification(recipientID, encryptedMessage)
    if err != nil {
        return err
    }
    LogNotificationEvent("SendNotification", "Notification sent to " + recipientID)
    return nil
}

// SetAlertThreshold: Sets a threshold level for triggering system alerts
func SetAlertThreshold(threshold int) {
    alertThreshold = threshold
    LogNotificationEvent("AlertThreshold", "Alert threshold set to " + string(threshold))
}

// CancelAlert: Cancels a specific alert if it is active
func CancelAlert(alertID string) error {
    if _, exists := customAlerts.Load(alertID); !exists {
        return errors.New("alert not found")
    }
    customAlerts.Delete(alertID)
    LogNotificationEvent("CancelAlert", "Alert " + alertID + " cancelled")
    return nil
}

// TriggerSystemAlert: Triggers a system alert if conditions meet the threshold
func TriggerSystemAlert(alertMessage string) {
    if len(alertHistory) >= alertThreshold {
        SendBroadcastMessage(alertMessage)
        alertHistory = append(alertHistory, alertMessage)
        LogNotificationEvent("SystemAlert", "System alert triggered: " + alertMessage)
    }
}

// LogNotificationEvent: Logs notification events in the ledger
func LogNotificationEvent(context, message string) error {
    encryptedMessage, err := encryption.Encrypt([]byte(message))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent(context, encryptedMessage)
}

// RegisterAlertHandler: Registers a handler for specific alert types
func RegisterAlertHandler(alertType string, handler func(string)) {
    alertHandlers.Store(alertType, handler)
    LogNotificationEvent("RegisterHandler", "Handler registered for alert type " + alertType)
}

// UnregisterAlertHandler: Unregisters a previously registered alert handler
func UnregisterAlertHandler(alertType string) {
    alertHandlers.Delete(alertType)
    LogNotificationEvent("UnregisterHandler", "Handler unregistered for alert type " + alertType)
}

// CheckAlertStatus: Checks if a specific alert is active or has been triggered
func CheckAlertStatus(alertID string) bool {
    _, exists := customAlerts.Load(alertID)
    LogNotificationEvent("CheckAlertStatus", "Checked status for alert " + alertID)
    return exists
}

// GetAlertHistory: Retrieves the history of triggered alerts
func GetAlertHistory(limit int) []string {
    if limit > len(alertHistory) {
        limit = len(alertHistory)
    }
    return alertHistory[len(alertHistory)-limit:]
}

// SendBroadcastMessage: Sends a broadcast message to all registered handlers
func SendBroadcastMessage(message string) {
    alertHandlers.Range(func(key, value interface{}) bool {
        handler := value.(func(string))
        handler(message)
        return true
    })
    LogNotificationEvent("BroadcastMessage", "Broadcast message sent: " + message)
}

// CreateCustomAlert: Creates a custom alert with specific conditions
func CreateCustomAlert(alertID, alertMessage string) {
    customAlerts.Store(alertID, alertMessage)
    LogNotificationEvent("CreateCustomAlert", "Custom alert " + alertID + " created")
}

// DeleteCustomAlert: Deletes a custom alert
func DeleteCustomAlert(alertID string) error {
    if _, exists := customAlerts.Load(alertID); !exists {
        return errors.New("custom alert not found")
    }
    customAlerts.Delete(alertID)
    LogNotificationEvent("DeleteCustomAlert", "Custom alert " + alertID + " deleted")
    return nil
}

// EnablePushNotification: Enables push notifications for the system
func EnablePushNotification() {
    pushNotificationsEnabled = true
    LogNotificationEvent("PushNotification", "Push notifications enabled")
}

// DisablePushNotification: Disables push notifications for the system
func DisablePushNotification() {
    pushNotificationsEnabled = false
    LogNotificationEvent("PushNotification", "Push notifications disabled")
}
