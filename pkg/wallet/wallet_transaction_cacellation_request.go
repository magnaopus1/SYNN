package wallet

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/transactions"
)

// CancellationRequest represents a request to cancel a transaction.
type CancellationRequest struct {
	ID                string            // Unique request ID
	TransactionID     string            // ID of the transaction to cancel
	UserID            string            // ID of the user requesting cancellation
	Timestamp         time.Time         // Timestamp of the cancellation request
	Status            string            // Current status (e.g., "Pending", "Approved", "Rejected")
	Reason            string            // Reason for the cancellation request
	DocumentEvidence  string            // Any evidence provided for the cancellation
	ContactDetails    string            // Contact details of the user requesting cancellation
	RequiredApprovals int               // Number of required approvals for cancellation
	ApprovalNodes     map[string]bool   // Nodes that have approved
	RejectionNodes    map[string]bool   // Nodes that have rejected
}



// WalletTransactionCancellationRequest handles the wallet side of requesting a transaction cancellation.
type WalletTransactionCancellationRequest struct {
	mutex            sync.Mutex
	WalletID         string                     // Unique wallet ID
	OwnerID          string                     // Owner's ID associated with the wallet
	Ledger           *ledger.Ledger             // Reference to the ledger for transaction lookups
	TransactionMgr   *transactions.TransactionCancellationManager // Transaction manager for handling cancellations
	Encryption       *common.Encryption         // Encryption service for secure data transmission
	Notification     WalletNotificationService  // Notification service for sending alerts
    CancellationReqs map[string]*transactions.CancellationRequest
}

// RequestCancellation allows a wallet owner to request the cancellation of a transaction.
func (wtcr *WalletTransactionCancellationRequest) RequestCancellation(transactionID, reason, documentEvidence, contactDetails string) (*transactions.CancellationRequest, error) {
    wtcr.mutex.Lock()
    defer wtcr.mutex.Unlock()

    // Check if a cancellation request for this transaction already exists
    if _, exists := wtcr.CancellationReqs[transactionID]; exists {
        return nil, fmt.Errorf("cancellation request already exists for transaction %s", transactionID)
    }

    // Request transaction cancellation using the TransactionCancellationManager
    cancellationRequest, err := wtcr.TransactionMgr.RequestTransactionCancellation(transactionID, reason, documentEvidence, contactDetails)
    if err != nil {
        return nil, fmt.Errorf("failed to process cancellation request: %v", err)
    }

    // Store the cancellation request in the map, using the correct type from the transactions package
    wtcr.CancellationReqs[transactionID] = cancellationRequest

    // Send a notification about the cancellation request
    err = wtcr.Notification.SendTransactionAlert(transactionID, wtcr.WalletID)
    if err != nil {
        return nil, fmt.Errorf("failed to send transaction cancellation notification: %v", err)
    }

    // Return the cancellation request, using the correct type
    return cancellationRequest, nil
}




// GetCancellationStatus allows the user to check the status of their transaction cancellation request.
func (wcr *WalletTransactionCancellationRequest) GetCancellationStatus(txnID string) (string, error) {
	wcr.mutex.Lock()
	defer wcr.mutex.Unlock()

	// Check if the cancellation request exists locally
	request, exists := wcr.CancellationReqs[txnID]
	if !exists {
		return "", fmt.Errorf("no cancellation request found for transaction %s", txnID)
	}

	// Return the current status of the cancellation request.
	return request.Status, nil
}


// CancelRequestRejected notifies the wallet owner that their cancellation request was rejected.
func (wcr *WalletTransactionCancellationRequest) CancelRequestRejected(txnID string) error {
	// Notify the user that the cancellation request was rejected using SendTransactionAlert.
	err := wcr.Notification.SendTransactionAlert(txnID, wcr.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to notify user of rejection: %v", err)
	}

	return nil
}

// CancelRequestApproved notifies the wallet owner that their cancellation request was approved.
func (wcr *WalletTransactionCancellationRequest) CancelRequestApproved(txnID string) error {
	// Notify the user that the cancellation request was approved using SendTransactionAlert.
	err := wcr.Notification.SendTransactionAlert(txnID, wcr.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to notify user of approval: %v", err)
	}

	return nil
}


