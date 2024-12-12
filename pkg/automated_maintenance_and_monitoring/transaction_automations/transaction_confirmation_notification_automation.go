package automations

import (
	"fmt"
	"log"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/common"
	"synnergy_network_demo/notifications"
)

// TransactionConfirmationNotificationAutomation automates transaction confirmation notifications.
type TransactionConfirmationNotificationAutomation struct {
	ledgerInstance          *ledger.Ledger
	mutex                   sync.Mutex
	stopChan                chan bool
	confirmationCheckInterval time.Duration
}

// NewTransactionConfirmationNotificationAutomation initializes a new TransactionConfirmationNotificationAutomation.
func NewTransactionConfirmationNotificationAutomation(ledgerInstance *ledger.Ledger) *TransactionConfirmationNotificationAutomation {
	return &TransactionConfirmationNotificationAutomation{
		ledgerInstance:          ledgerInstance,
		confirmationCheckInterval: 500 * time.Millisecond, // Check confirmation status every 0.5 seconds
		stopChan:                make(chan bool),
	}
}

// Start begins the continuous monitoring of transaction confirmations.
func (t *TransactionConfirmationNotificationAutomation) Start() {
	go t.runConfirmationLoop()
	log.Println("Transaction Confirmation Notification Automation started.")
}

// Stop stops the continuous confirmation monitoring process.
func (t *TransactionConfirmationNotificationAutomation) Stop() {
	t.stopChan <- true
	log.Println("Transaction Confirmation Notification Automation stopped.")
}

// runConfirmationLoop continuously monitors and sends notifications when transactions reach specific confirmation milestones.
func (t *TransactionConfirmationNotificationAutomation) runConfirmationLoop() {
	ticker := time.NewTicker(t.confirmationCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.checkTransactionConfirmations()
		case <-t.stopChan:
			return
		}
	}
}

// checkTransactionConfirmations checks the confirmation status of all pending transactions.
func (t *TransactionConfirmationNotificationAutomation) checkTransactionConfirmations() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Fetch all pending transactions from the ledger
	pendingTransactions, err := t.ledgerInstance.GetPendingTransactions()
	if err != nil {
		log.Printf("Failed to fetch pending transactions: %v", err)
		return
	}

	// Loop through all pending transactions and check their confirmation status
	for _, tx := range pendingTransactions {
		confirmations, err := t.ledgerInstance.GetTransactionConfirmations(tx.ID)
		if err != nil {
			log.Printf("Failed to get confirmations for transaction %s: %v", tx.ID, err)
			continue
		}

		// Notify based on confirmation milestones
		if confirmations >= 6 && confirmations < 85 {
			err := t.sendNotification(tx.ID, "Your transaction has reached 6 confirmations.")
			if err != nil {
				log.Printf("Failed to send 6-confirmation notification: %v", err)
			}
		} else if confirmations >= 85 {
			err := t.sendNotification(tx.ID, "Your transaction has been fully confirmed (85 confirmations).")
			if err != nil {
				log.Printf("Failed to send 85-confirmation notification: %v", err)
			}
		}
	}
}

// sendNotification sends a notification when a transaction reaches a confirmation milestone.
func (t *TransactionConfirmationNotificationAutomation) sendNotification(transactionID, message string) error {
	// Retrieve the recipient's contact details from the ledger
	recipient, err := t.ledgerInstance.GetTransactionRecipient(transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve recipient for transaction %s: %v", transactionID, err)
	}

	// Construct the notification
	notification := notifications.Notification{
		Recipient: recipient.Email, // Assuming the recipient has an email address
		Message:   message,
		Type:      notifications.EmailNotification, // Assuming email notifications for now
	}

	// Send the notification using the notifications package
	err = notifications.SendNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	log.Printf("Notification sent to %s for transaction %s", recipient.Email, transactionID)
	return nil
}
