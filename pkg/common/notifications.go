package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// Notification represents a structure for sending notifications
type Notification struct {
	Recipient  string           // Email address or webhook URL
	Message    string           // Message to be sent
	Type       NotificationType // Type of notification (Email, Webhook, etc.)
	SentAt     time.Time        // Timestamp of when the notification was sent
	IsDelivered bool            // Status whether the notification was successfully delivered
}

// NotificationType represents different notification types
type NotificationType string

const (
	EmailNotification   NotificationType = "email"
	WebhookNotification NotificationType = "webhook"
	SMSNotification     NotificationType = "sms"
	PushNotification    NotificationType = "push"
)

// NotificationLog records the history of sent notifications
type NotificationLog struct {
	NotificationID string    // Unique identifier for the notification
	Recipient      string    // Recipient of the notification
	Type           NotificationType // Type of notification sent
	Message        string    // Content of the notification
	SentAt         time.Time // Time when the notification was sent
	Status         string    // Status of the notification (Sent, Failed, Delivered)
}

// NotificationManager manages the sending and tracking of notifications
type NotificationManager struct {
	Notifications     map[string]*Notification  // Active notifications by their unique ID
	NotificationLog   map[string]*NotificationLog // Log of all sent notifications
	mu                sync.Mutex               // Mutex for thread-safe access
	Ledger            *ledger.Ledger           // Ledger for tracking notification actions and status
	EncryptionService *Encryption              // Encryption service to secure notification data
}


// Logger represents a simple logging mechanism for the system
type Logger struct {
    Logs      []string    // List of log messages
    LogLevel  string      // Current logging level (e.g., "INFO", "ERROR", "DEBUG")
    LogTime   time.Time   // Timestamp for the log
    mutex     sync.Mutex  // Mutex for thread-safe log access
}

// Log adds a new log entry to the Logger
func (l *Logger) Log(level string, message string) {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    logEntry := time.Now().Format(time.RFC3339) + " [" + level + "] " + message
    l.Logs = append(l.Logs, logEntry)
    l.LogLevel = level
    l.LogTime = time.Now()
}

// GetLogs returns all log entries
func (l *Logger) GetLogs() []string {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    return l.Logs
}


// SendNotification sends a notification based on the type of notification (Email, Webhook, etc.)
func SendNotification(notification Notification) error {
	switch notification.Type {
	case EmailNotification:
		// Simulate sending an email notification
		return sendEmailNotification(notification.Recipient, notification.Message)
	case WebhookNotification:
		// Simulate sending a webhook notification
		return sendWebhookNotification(notification.Recipient, notification.Message)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

// sendEmailNotification simulates sending an email
func sendEmailNotification(recipient, message string) error {
	// In a real-world scenario, integrate with an email API such as SendGrid, SES, or any SMTP server
	log.Printf("Email sent to %s: %s", recipient, message)
	return nil
}

// sendWebhookNotification simulates sending a webhook
func sendWebhookNotification(webhookURL, message string) error {
	// Simulate sending a webhook POST request with the message payload
	payload := map[string]string{"message": message}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status: %s", resp.Status)
	}

	log.Printf("Webhook notification sent to %s: %s", webhookURL, message)
	return nil
}


