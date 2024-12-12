package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/transactions"
)

// ReversalRequest represents a request to reverse a transaction.
type ReversalRequest struct {
    TransactionID string    // Unique transaction ID
    WalletID      string    // Wallet ID requesting the reversal
    Reason        string    // Reason for reversal
    Status        string    // Current status of the request (e.g., pending, approved, rejected)
    Timestamp     time.Time // Time when the request was created
}


// WalletTransactionReversalRequest handles the wallet side of requesting a transaction reversal.
type WalletTransactionReversalRequest struct {
    mutex            sync.Mutex
    WalletID         string                     // Unique wallet ID
    OwnerID          string                     // Owner's ID associated with the wallet
    Ledger           *ledger.Ledger             // Reference to the ledger for transaction lookups
    TransactionMgr   *transactions.TransactionReversalManager // Manager to handle transaction reversals
    Encryption       *common.Encryption         // Encryption service for secure data transmission
    Notification     WalletNotificationService  // Service for sending notifications
    ReversalRequests map[string]ReversalRequest // Track reversal requests by transaction ID
}

// NewWalletTransactionReversalRequest initializes a new reversal request manager for a wallet.
func NewWalletTransactionReversalRequest(walletID string, ownerID string, ledger *ledger.Ledger, transactionMgr *transactions.TransactionReversalManager, encryptionService *common.Encryption, notificationService WalletNotificationService) *WalletTransactionReversalRequest {
    return &WalletTransactionReversalRequest{
        WalletID:       walletID,
        OwnerID:        ownerID,
        Ledger:         ledger,
        TransactionMgr: transactionMgr,
        Encryption:     encryptionService,
        Notification:   notificationService,
        ReversalRequests: make(map[string]ReversalRequest), // Initialize map to track requests
    }
}


// RequestReversal allows a wallet owner to request the reversal of a transaction.
func (wrr *WalletTransactionReversalRequest) RequestReversal(txnID, reason, contactDetails string, evidence []byte) (*ReversalRequest, error) {
    wrr.mutex.Lock()
    defer wrr.mutex.Unlock()

    // Validate that the transaction exists.
    txn, err := wrr.Ledger.GetTransactionByID(txnID) // Only pass txnID
    if err != nil || txn.From != wrr.OwnerID { // Replace OwnerID with the correct field (e.g., WalletID)
        return nil, errors.New("transaction not found or unauthorized request")
    }

    // Validate that the transaction is within the reversal window (between 24 hours and 28 days after the transaction).
    if time.Since(txn.Timestamp) < 24*time.Hour || time.Since(txn.Timestamp) > 28*24*time.Hour {
        return nil, errors.New("reversal request must be made after 24 hours and within 28 days of the transaction")
    }

    // Freeze the transaction amount to prevent further changes while the request is processed.
    err = wrr.Ledger.FreezeTransactionAmount(txn.ID, txn.Amount) // Use correct field names
    if err != nil {
        return nil, fmt.Errorf("failed to freeze transaction: %v", err)
    }

    // Encrypt the reversal evidence (such as document proof).
    encryptedEvidence, err := wrr.Encryption.EncryptData("AES", evidence, common.EncryptionKey) // Added encryption method name
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt evidence: %v", err)
    }

    // Submit the reversal request to the transaction reversal manager.
    reversalRequest, err := wrr.TransactionMgr.RequestTransactionReversal(txnID, reason, contactDetails, encryptedEvidence)
    if err != nil {
        return nil, fmt.Errorf("failed to submit reversal request: %v", err)
    }

    // Notify the user of the successful submission.
    err = wrr.Notification.SendTransactionAlert(txnID, wrr.OwnerID) // Fixed to use SendTransactionAlert method
    if err != nil {
        return nil, fmt.Errorf("failed to notify user: %v", err)
    }

    // Assuming transactions.ReversalRequest is the correct type, you need to convert it.
    localReversalRequest := &ReversalRequest{
        TransactionID: reversalRequest.TransactionID, // Only map fields that exist in your struct
        WalletID:      reversalRequest.WalletID,      // Adjust to WalletID instead of UserID
        Reason:        reversalRequest.Reason,
        Status:        reversalRequest.Status,
        Timestamp:     reversalRequest.Timestamp,
    }

    return localReversalRequest, nil
}






// GetReversalStatus allows the user to check the status of their transaction reversal request.
func (wrr *WalletTransactionReversalRequest) GetReversalStatus(txnID string) (string, error) {
	wrr.mutex.Lock()
	defer wrr.mutex.Unlock()

	// Retrieve the reversal request from the ledger.
	request, err := wrr.Ledger.GetReversalRequestByTransactionID(txnID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve reversal request: %v", err)
	}

	// Return the current status of the reversal request.
	return request.Status, nil
}


// ReversalRequestRejected notifies the wallet owner that their reversal request was rejected.
func (wrr *WalletTransactionReversalRequest) ReversalRequestRejected(txnID string) error {
	// Notify the user that the reversal request was rejected.
	err := wrr.Notification.SendTransactionAlert(txnID, wrr.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to notify user of rejection: %v", err)
	}

	return nil
}

// ReversalRequestApproved notifies the wallet owner that their reversal request was approved.
func (wrr *WalletTransactionReversalRequest) ReversalRequestApproved(txnID string) error {
	// Notify the user that the reversal request was approved.
	err := wrr.Notification.SendTransactionAlert(txnID, wrr.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to notify user of approval: %v", err)
	}

	return nil
}

