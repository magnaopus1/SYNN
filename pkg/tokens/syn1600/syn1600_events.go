package syn1600

import (
	"errors"
	"time"
)

// EventManager handles event logging and management for SYN1600 tokens.
type EventManager struct {
	Ledger ledger.Ledger
}

// LogEvent logs an event related to a SYN1600 token.
func (em *EventManager) LogEvent(tokenID string, eventType string, description string, performedBy string) error {
	// Retrieve the token from the ledger
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Create an event log
	newEvent := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   eventType,
		Description: description,
		EventDate:   time.Now(),
		PerformedBy: performedBy,
	}

	// Append the new event to the token's event log
	token.(*common.SYN1600Token).EventLogs = append(token.(*common.SYN1600Token).EventLogs, newEvent)

	// Update the ledger with the new event log
	return em.Ledger.UpdateToken(tokenID, token)
}

// NotifyRealTimeEvent sends real-time notifications when important events occur for a SYN1600 token.
func (em *EventManager) NotifyRealTimeEvent(tokenID string, eventType string, notificationMessage string) error {
	// Log the event before sending a notification
	err := em.LogEvent(tokenID, eventType, notificationMessage, "System")
	if err != nil {
		return err
	}

	// Example: Simulate sending a real-time notification (you can integrate with a messaging or notification service)
	// Notify the owner of the token
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	ownerID := token.(*common.SYN1600Token).Owner
	// Here, you would send the notification to the owner's system (email, app notification, etc.)
	// For this example, we will simulate a simple message
	sendNotification(ownerID, notificationMessage)

	return nil
}

// HandleOwnershipChange logs and notifies when ownership of a SYN1600 token is transferred.
func (em *EventManager) HandleOwnershipChange(tokenID string, newOwnerID string, performedBy string) error {
	// Retrieve the token from the ledger
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Log the ownership change
	err = em.LogEvent(tokenID, "OwnershipChange", "Ownership changed to "+newOwnerID, performedBy)
	if err != nil {
		return err
	}

	// Update the owner of the token
	token.(*common.SYN1600Token).Owner = newOwnerID

	// Update the ledger with the new owner information
	err = em.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return err
	}

	// Send real-time notification about the ownership change
	notificationMessage := "The ownership of token " + tokenID + " has been transferred to " + newOwnerID
	return em.NotifyRealTimeEvent(tokenID, "OwnershipChange", notificationMessage)
}

// HandleRevenueDistributionEvent logs and notifies when royalties are distributed for a SYN1600 token.
func (em *EventManager) HandleRevenueDistributionEvent(tokenID string, distributionType string, amount float64, recipientID string) error {
	// Log the revenue distribution event
	description := "Distributed " + distributionType + " royalties of " + formatAmount(amount) + " to recipient " + recipientID
	err := em.LogEvent(tokenID, "RevenueDistribution", description, "RoyaltyDistributionModule")
	if err != nil {
		return err
	}

	// Notify the recipient of the royalties
	notificationMessage := "You have received " + formatAmount(amount) + " in " + distributionType + " royalties."
	return em.NotifyRealTimeEvent(tokenID, "RevenueDistribution", notificationMessage)
}

// HandleComplianceAuditEvent logs and notifies when a compliance audit occurs for a SYN1600 token.
func (em *EventManager) HandleComplianceAuditEvent(tokenID string, auditDescription string, performedBy string) error {
	// Log the compliance audit
	err := em.LogEvent(tokenID, "ComplianceAudit", auditDescription, performedBy)
	if err != nil {
		return err
	}

	// Notify the token owner of the audit
	notificationMessage := "A compliance audit has been performed for token " + tokenID + ". " + auditDescription
	return em.NotifyRealTimeEvent(tokenID, "ComplianceAudit", notificationMessage)
}

// EncryptEventData encrypts sensitive event data for a SYN1600 token event log.
func (em *EventManager) EncryptEventData(tokenID string, eventID string, key []byte) error {
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return err
	}

	// Find the event by eventID
	var eventLog *common.EventLog
	for i, event := range token.(*common.SYN1600Token).EventLogs {
		if event.EventID == eventID {
			eventLog = &token.(*common.SYN1600Token).EventLogs[i]
			break
		}
	}
	if eventLog == nil {
		return errors.New("event not found")
	}

	// Encrypt the event description
	encryptedDescription, err := encryption.Encrypt([]byte(eventLog.Description), key)
	if err != nil {
		return err
	}
	eventLog.Description = string(encryptedDescription)

	// Update the ledger with the encrypted event log
	return em.Ledger.UpdateToken(tokenID, token)
}

// DecryptEventData decrypts sensitive event data for a SYN1600 token event log.
func (em *EventManager) DecryptEventData(tokenID string, eventID string, key []byte) (string, error) {
	token, err := em.Ledger.GetToken(tokenID)
	if err != nil {
		return "", err
	}

	// Find the event by eventID
	var eventLog *common.EventLog
	for i, event := range token.(*common.SYN1600Token).EventLogs {
		if event.EventID == eventID {
			eventLog = &token.(*common.SYN1600Token).EventLogs[i]
			break
		}
	}
	if eventLog == nil {
		return "", errors.New("event not found")
	}

	// Decrypt the event description
	decryptedDescription, err := encryption.Decrypt([]byte(eventLog.Description), key)
	if err != nil {
		return "", err
	}

	return string(decryptedDescription), nil
}

// Helper function to simulate sending notifications.
func sendNotification(userID string, message string) {
	// Simulate sending a notification
	println("Notification sent to user:", userID, "Message:", message)
}

// Helper function to generate a unique ID for events.
func generateUniqueID() string {
	return "EVENT_" + time.Now().Format("20060102150405")
}

// Helper function to format amounts into a readable format.
func formatAmount(amount float64) string {
	return "$" + fmt.Sprintf("%.2f", amount)
}
