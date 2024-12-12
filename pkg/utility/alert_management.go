package utility

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// SendEmailAlert sends an email alert to the specified recipients with encrypted content.
func SendEmailAlert(recipients []string, subject, message string) error {
	encryptedMessage, err := common.EncryptMessage(message)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}
	// Placeholder for sending email logic; replace with actual email sending service integration
	fmt.Printf("Sending email to %v with subject: %s\n", recipients, subject)
	ledger.RecordAlert("EmailAlert", recipients, subject)
	return nil
}

// SendSMSAlert sends an SMS alert to the specified recipients with encrypted content.
func SendSMSAlert(recipients []string, message string) error {
	encryptedMessage, err := common.EncryptMessage(message)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}
	// Placeholder for SMS sending logic; replace with actual SMS sending service integration
	fmt.Printf("Sending SMS to %v\n", recipients)
	ledger.RecordAlert("SMSAlert", recipients, "")
	return nil
}

// SetNotificationPreferences sets the alert preferences for a user, allowing them to specify delivery channels.
func SetNotificationPreferences(userID string, preferences map[string]bool) error {
	if err := ledger.StoreNotificationPreferences(userID, preferences); err != nil {
		return fmt.Errorf("failed to set preferences: %v", err)
	}
	return nil
}

// GetNotificationPreferences retrieves the alert preferences for a specified user.
func GetNotificationPreferences(userID string) (map[string]bool, error) {
	preferences, err := ledger.FetchNotificationPreferences(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %v", err)
	}
	return preferences, nil
}

// ScheduleNotification schedules a notification for a future time, storing it in the ledger.
func ScheduleNotification(userID, alertType, message string, sendTime time.Time) error {
	if time.Now().After(sendTime) {
		return errors.New("send time must be in the future")
	}
	ledger.RecordScheduledNotification(userID, alertType, message, sendTime)
	return nil
}

// CancelScheduledNotification cancels a scheduled notification.
func CancelScheduledNotification(userID, alertType string) error {
	if err := ledger.CancelNotification(userID, alertType); err != nil {
		return fmt.Errorf("failed to cancel notification: %v", err)
	}
	return nil
}

// SetUrgentAlertFlag marks a notification as urgent, which can affect delivery speed and priority.
func SetUrgentAlertFlag(notificationID string) error {
	if err := ledger.UpdateAlertFlag(notificationID, "urgent", true); err != nil {
		return fmt.Errorf("failed to set urgent flag: %v", err)
	}
	return nil
}

// ClearUrgentAlertFlag clears the urgent flag on a notification.
func ClearUrgentAlertFlag(notificationID string) error {
	if err := ledger.UpdateAlertFlag(notificationID, "urgent", false); err != nil {
		return fmt.Errorf("failed to clear urgent flag: %v", err)
	}
	return nil
}

// LogAlertAcknowledgement records when an alert has been acknowledged by the recipient.
func LogAlertAcknowledgement(notificationID, userID string) error {
	if err := ledger.LogAcknowledgement(notificationID, userID); err != nil {
		return fmt.Errorf("failed to log acknowledgement: %v", err)
	}
	return nil
}

// SendSystemWarning sends a system warning alert with encrypted content.
func SendSystemWarning(recipients []string, message string) error {
	encryptedMessage, err := common.EncryptMessage(message)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}
	// Placeholder for system warning logic; replace with actual alerting system
	fmt.Printf("Sending system warning to %v\n", recipients)
	ledger.RecordAlert("SystemWarning", recipients, "")
	return nil
}

// SendSystemInfo sends a system info alert with encrypted content.
func SendSystemInfo(recipients []string, message string) error {
	encryptedMessage, err := common.EncryptMessage(message)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}
	// Placeholder for system info logic; replace with actual alerting system
	fmt.Printf("Sending system info to %v\n", recipients)
	ledger.RecordAlert("SystemInfo", recipients, "")
	return nil
}

// QueryNotificationLog retrieves the log of notifications sent to a user.
func QueryNotificationLog(userID string) ([]ledger.NotificationLog, error) {
	logs, err := ledger.FetchNotificationLog(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notification log: %v", err)
	}
	return logs, nil
}

// TrackAlertResponseTime measures the time taken for an alert to be acknowledged by the recipient.
func TrackAlertResponseTime(notificationID string, acknowledgmentTime time.Time) error {
	sentTime, err := ledger.GetNotificationSentTime(notificationID)
	if err != nil {
		return fmt.Errorf("failed to get sent time: %v", err)
	}
	responseTime := acknowledgmentTime.Sub(sentTime)
	ledger.RecordResponseTime(notificationID, responseTime)
	return nil
}

// SetAlertPriority sets the priority level of a notification to adjust delivery speed and user attention.
func SetAlertPriority(notificationID string, priorityLevel int) error {
	if priorityLevel < 0 || priorityLevel > 5 {
		return errors.New("priority level must be between 0 and 5")
	}
	if err := ledger.UpdateAlertPriority(notificationID, priorityLevel); err != nil {
		return fmt.Errorf("failed to set alert priority: %v", err)
	}
	return nil
}
