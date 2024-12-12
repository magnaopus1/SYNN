package dao

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewDAOFundVault initializes a new DAO Fund Vault.
func NewDAOFundVault(daoID string, initialBalance float64, ledgerInstance *ledger.Ledger, encryptionService *common.Encryption, syn900Verifier *common.Syn900Verifier) *DAOFundVault {
	return &DAOFundVault{
		DAOID:            daoID,
		Balance:          initialBalance,
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
		Syn900Verifier:   syn900Verifier,
		TransactionLimit: 10000, // Example transaction limit
		Admins:           make(map[string]bool),
	}
}

// AddAdmin adds an admin address to the DAO vault, enabling them to approve transactions.
func (vault *DAOFundVault) AddAdmin(adminAddress string) {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	vault.Admins[adminAddress] = true
	fmt.Printf("Admin %s added to DAO %s fund vault.\n", adminAddress, vault.DAOID)
}

// RemoveAdmin removes an admin address from the DAO vault.
func (vault *DAOFundVault) RemoveAdmin(adminAddress string) {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	delete(vault.Admins, adminAddress)
	fmt.Printf("Admin %s removed from DAO %s fund vault.\n", adminAddress, vault.DAOID)
}

// SubmitTransaction submits a transaction for approval by DAO admins.
func (vault *DAOFundVault) SubmitTransaction(amount float64, recipient string, submittedBy string) (*common.VaultTransaction, error) {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	// Check if the submitter is an admin
	if !vault.Admins[submittedBy] {
		return nil, errors.New("only admins can submit transactions")
	}

	// Ensure transaction amount does not exceed the current balance
	if amount > vault.Balance {
		return nil, errors.New("insufficient funds")
	}

	// Ensure the transaction does not exceed the daily transaction limit
	if time.Since(vault.LastTransactionAt) < 24*time.Hour && amount > vault.TransactionLimit {
		return nil, errors.New("transaction amount exceeds daily limit")
	}

	// Create the transaction
	transaction := &VaultTransaction{
		TransactionID:  GenerateUniqueID(),
		Amount:         amount,
		Recipient:      recipient,
		Timestamp:      time.Now(),
		Status:         "Pending",
	}

	// Add the transaction to the queue
	vault.TransactionQueue = append(vault.TransactionQueue, *transaction)
	vault.LastTransactionAt = time.Now()

	// Log transaction submission in the ledger
	err := vault.Ledger.DAOLedger.RecordVaultTransaction(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to record transaction in ledger: %v", err)
	}

	fmt.Printf("Transaction submitted: %s, Amount: %.2f, Recipient: %s\n", transaction.TransactionID, amount, recipient)
	return transaction, nil
}

// ApproveTransaction approves a pending transaction.
func (vault *DAOFundVault) ApproveTransaction(transactionID, adminAddress string) error {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	// Check if the admin exists
	if !vault.Admins[adminAddress] {
		return errors.New("only admins can approve transactions")
	}

	// Find the transaction
	for i, transaction := range vault.TransactionQueue {
		if transaction.TransactionID == transactionID {
			// Add admin approval
			transaction.ApprovedBy = append(transaction.ApprovedBy, adminAddress)

			// Check if a majority approval is reached (e.g., 2 approvals)
			if len(transaction.ApprovedBy) >= 2 {
				// Update transaction status
				transaction.Status = "Approved"
				vault.TransactionQueue[i] = transaction

				// Execute the transaction
				vault.Balance -= transaction.Amount
				err := vault.Ledger.BlockchainConsensusCoinLedger.RecordTransactionExecution(transaction)
				if err != nil {
					return fmt.Errorf("failed to execute transaction: %v", err)
				}

				fmt.Printf("Transaction %s approved and executed. Amount: %.2f\n", transactionID, transaction.Amount)
				return nil
			}

			// Log the approval
			err := vault.Ledger.DAOLedger.RecordTransactionApproval(transactionID, adminAddress)
			if err != nil {
				return fmt.Errorf("failed to record approval in ledger: %v", err)
			}

			fmt.Printf("Transaction %s approved by %s. Waiting for additional approvals...\n", transactionID, adminAddress)
			return nil
		}
	}

	return errors.New("transaction not found")
}

// EmergencyAccess triggers an emergency fund access request through Syn900.
func (vault *DAOFundVault) EmergencyAccess(requestedBy, reason string) (*EmergencyAccessRequest, error) {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	// Create an emergency access request
	request := &EmergencyAccessRequest{
		RequestID:   GenerateUniqueID(),
		RequestedBy: requestedBy,
		Reason:      reason,
		Timestamp:   time.Now(),
		Status:      "Pending",
	}

	// Log the request in the ledger
	err := vault.Ledger.DAOLedger.RecordEmergencyAccessRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to record emergency access request: %v", err)
	}

	// Verify through Syn900 if emergency access can be granted
	verified, err := vault.Syn900Verifier.VerifyEmergencyAccess(requestedBy, reason)
	if err != nil || !verified {
		return nil, errors.New("emergency access denied by Syn900")
	}

	// If verified, the request will be approved and the funds released
	request.Status = "Approved"
	request.ApprovalConfirm = append(request.ApprovalConfirm, requestedBy)

	// Log the approval in the ledger
	err = vault.Ledger.DAOLedger.RecordEmergencyAccessApproval(request)
	if err != nil {
		return nil, fmt.Errorf("failed to record emergency access approval in ledger: %v", err)
	}

	fmt.Printf("Emergency access granted for DAO %s by %s for reason: %s\n", vault.DAOID, requestedBy, reason)
	return request, nil
}

// RejectTransaction rejects a pending transaction.
func (vault *DAOFundVault) RejectTransaction(transactionID, adminAddress string) error {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	// Check if the admin exists
	if !vault.Admins[adminAddress] {
		return errors.New("only admins can reject transactions")
	}

	// Find and reject the transaction
	for i, transaction := range vault.TransactionQueue {
		if transaction.TransactionID == transactionID {
			// Update transaction status
			transaction.Status = "Rejected"
			vault.TransactionQueue[i] = transaction

			// Log the rejection in the ledger
			err := vault.Ledger.DAOLedger.RecordTransactionRejection(transactionID, adminAddress)
			if err != nil {
				return fmt.Errorf("failed to record rejection in ledger: %v", err)
			}

			fmt.Printf("Transaction %s rejected by %s\n", transactionID, adminAddress)
			return nil
		}
	}

	return errors.New("transaction not found")
}

// ViewBalance returns the current balance of the DAO vault.
func (vault *DAOFundVault) ViewBalance() float64 {
	vault.mutex.Lock()
	defer vault.mutex.Unlock()

	return vault.Balance
}
