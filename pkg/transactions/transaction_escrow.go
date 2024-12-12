package transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewEscrowManager initializes a new EscrowManager.
func NewEscrowManager(ledgerInstance *ledger.Ledger) *EscrowManager {
	return &EscrowManager{
		ledgerInstance: ledgerInstance,
		escrows:        make(map[string]*EscrowTransaction),
	}
}

// CreateEscrow initializes a new escrow transaction.
func (em *EscrowManager) CreateEscrow(senderID, receiverID string, amount float64, condition string) (*EscrowTransaction, error) {
    em.mutex.Lock()
    defer em.mutex.Unlock()

    escrowID := generateEscrowID(senderID, receiverID)
    escrow := &EscrowTransaction{
        EscrowID:     escrowID,
        SenderID:     senderID,
        ReceiverID:   receiverID,
        Amount:       amount,
        Status:       EscrowStatusPending,  // Use EscrowStatusPending directly (no need to convert to string)
        CreationTime: time.Now(),
        Condition:    condition,
    }

    // Create an encryption instance within the function
    encryptionInstance := &common.Encryption{}

    // Encrypt escrow details using the encryption instance and AES algorithm
    encryptedEscrow, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%+v", escrow)), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt escrow details: %v", err)
    }

    // Hash the encrypted escrow details to create a unique ID (since the TransactionRecord requires a Hash field)
    hash := sha256.Sum256(encryptedEscrow)
    transactionHash := fmt.Sprintf("%x", hash)

    // Create a transaction record for the escrow
    record := ledger.TransactionRecord{
        From:        senderID,
        To:          receiverID,
        Amount:      amount,
        Hash:        transactionHash,       
        Status:      string(EscrowStatusPending),    // Use EscrowStatusPending directly
        Timestamp:   time.Now(),
        Action:      "EscrowCreated",       
        Details:     fmt.Sprintf("Encrypted escrow details: %x", encryptedEscrow),
    }

    // Log the escrow in the ledger, passing both the escrowID and the transaction record
    err = em.ledgerInstance.RecordEscrowTransaction(escrowID, record) // Pass escrowID and record
    if err != nil {
        return nil, fmt.Errorf("failed to record escrow transaction: %v", err)
    }

    // Store the escrow transaction
    em.escrows[escrow.EscrowID] = escrow
    return escrow, nil
}




// ReleaseEscrow releases the funds held in escrow to the receiver.
func (em *EscrowManager) ReleaseEscrow(escrowID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	escrow, exists := em.escrows[escrowID]
	if !exists {
		return fmt.Errorf("escrow transaction not found")
	}

	if escrow.Status != EscrowStatusPending {
		return fmt.Errorf("escrow transaction is not in a pending state")
	}

	// Convert escrow amount from float64 to uint64 for fund transfer
	amountToTransfer := uint64(escrow.Amount)

	// Convert uint64 to float64
	amountToTransferFloat := float64(amountToTransfer)

	// Release funds to the receiver
	err := em.ledgerInstance.TransferFunds(escrow.SenderID, escrow.ReceiverID, amountToTransferFloat)
	if err != nil {
		return fmt.Errorf("failed to release escrow funds: %v", err)
	}

	escrow.Status = EscrowStatusReleased
	escrow.ReleaseTime = time.Now()

	// Log the release in the ledger
	err = em.logEscrowAction(escrowID, "Funds released")
	if err != nil {
		return fmt.Errorf("failed to log escrow release: %v", err)
	}

	return nil
}

// CancelEscrow cancels the escrow transaction and returns the funds to the sender.
func (em *EscrowManager) CancelEscrow(escrowID string) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	escrow, exists := em.escrows[escrowID]
	if !exists {
		return fmt.Errorf("escrow transaction not found")
	}

	if escrow.Status != EscrowStatusPending {
		return fmt.Errorf("escrow transaction is not in a pending state")
	}

	// Convert escrow amount from float64 to uint64
	amountToTransfer := uint64(escrow.Amount)

	// Convert uint64 to float64
	amountToTransferFloat := float64(amountToTransfer)

	// Return funds to the sender
	err := em.ledgerInstance.TransferFunds(escrow.ReceiverID, escrow.SenderID, amountToTransferFloat)
	if err != nil {
		return fmt.Errorf("failed to cancel escrow and return funds: %v", err)
	}

	escrow.Status = EscrowStatusCancelled

	// Log the cancellation in the ledger
	err = em.logEscrowAction(escrowID, "Escrow cancelled")
	if err != nil {
		return fmt.Errorf("failed to log escrow cancellation: %v", err)
	}

	return nil
}


// GetEscrow retrieves the details of an escrow transaction.
func (em *EscrowManager) GetEscrow(escrowID string) (*EscrowTransaction, error) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	escrow, exists := em.escrows[escrowID]
	if !exists {
		return nil, fmt.Errorf("escrow transaction not found")
	}

	return escrow, nil
}


// logEscrowAction logs escrow actions such as release or cancellation into the ledger.
func (em *EscrowManager) logEscrowAction(escrowID, action string) error {
	// Create an encryption instance within the function
	encryptionInstance := &common.Encryption{}

	logMessage := fmt.Sprintf("Escrow %s: %s", escrowID, action)
	
	// Encrypt the log message using the encryption instance
	encryptedLogMessage, err := encryptionInstance.EncryptData("AES", []byte(logMessage), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt escrow log: %v", err)
	}

	// Convert encryptedLogMessage (from []byte to string) using hex encoding
	encryptedLogMessageStr := hex.EncodeToString(encryptedLogMessage)

	// Log the encrypted message into the ledger
	return em.ledgerInstance.RecordEscrowLog(escrowID, encryptedLogMessageStr)
}



// generateEscrowID generates a unique escrow ID based on sender and receiver details.
func generateEscrowID(senderID, receiverID string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%s-%d", senderID, receiverID, timestamp)
}
