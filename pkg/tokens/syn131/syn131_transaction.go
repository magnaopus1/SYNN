package syn131

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewTransactionManager creates a new instance of TransactionManager.
func NewTransactionManager(ledger *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *TransactionManager {
	return &TransactionManager{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}

// createOwnershipTransaction creates a new ownership transfer transaction.
// createOwnershipTransaction creates a new ownership transfer transaction.
func (tm *TransactionManager) createOwnershipTransaction(assetID, fromOwner, toOwner string, amount float64) (*OwnershipTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Input validation
	if assetID == "" || fromOwner == "" || toOwner == "" {
		return nil, errors.New("assetID, fromOwner, and toOwner must be provided")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Encrypt transaction data
	encryptedData, err := tm.EncryptionService.EncryptData([]byte(assetID + fromOwner + toOwner + fmt.Sprintf("%f", amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt transaction data")
	}

	// Create the transaction
	transaction := &OwnershipTransaction{
		TransactionID: generateTransactionID(assetID),
		AssetID:       assetID,
		FromOwner:     fromOwner,
		ToOwner:       toOwner,
		Amount:        amount,
		Timestamp:     time.Now(),
		EncryptedData: string(encryptedData), // Convert encrypted data to string
		Status:        "pending",
		Fee:           0.01, // Example fee logic
	}

	// Validate the transaction with Synnergy Consensus
	err = tm.ConsensusEngine.ProcessTransactions([]common.Transaction{
		{
			TransactionID:   transaction.TransactionID,
			TransactionType: "Ownership",
			Amount:          transaction.Amount,
			Timestamp:       transaction.Timestamp,
			Details:         "Ownership Transfer",
		},
	}, nil) // No cross-chain transactions in this case
	if err != nil {
		return nil, fmt.Errorf("transaction validation failed via Synnergy Consensus: %v", err)
	}

	// Record the transaction in the ledger
	ledgerTransaction := &ledger.OwnershipTransaction{
		TransactionID: transaction.TransactionID,
		AssetID:       transaction.AssetID,
		FromOwner:     transaction.FromOwner,
		ToOwner:       transaction.ToOwner,
		Amount:        transaction.Amount,
		Timestamp:     transaction.Timestamp,
		EncryptedData: transaction.EncryptedData,
		Status:        transaction.Status,
		Fee:           transaction.Fee,
	}

	err = tm.Ledger.RecordOwnershipTransaction(transaction.TransactionID, ledgerTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to record ownership transaction in the ledger: %v", err)
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}



// createRentalTransaction handles rental payment transactions.
func (tm *TransactionManager) createRentalTransaction(rentalAgreementID, renter, lessor string, amount float64) (*RentalTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Input validation
	if rentalAgreementID == "" || renter == "" || lessor == "" {
		return nil, errors.New("rentalAgreementID, renter, and lessor must be provided")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(rentalAgreementID + renter + lessor + fmt.Sprintf("%f", amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt rental transaction data")
	}

	transaction := &RentalTransaction{
		TransactionID:    generateTransactionID(rentalAgreementID),
		RentalAgreementID: rentalAgreementID,
		Renter:           renter,
		Lessor:           lessor,
		Amount:           amount,
		PaymentDate:      time.Now(),
		NextPaymentDue:   time.Now().AddDate(0, 1, 0), // Assume monthly payments
		EncryptedData:    encryptedData,
		Status:           "pending",
		Fee:              0.01, // Example fee logic
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ProcessTransactions([]common.Transaction{transaction}); err != nil {
		return nil, errors.New("rental transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.recordRentalTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record rental transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// createLeaseTransaction handles lease payment transactions.
func (tm *TransactionManager) createLeaseTransaction(leaseAgreementID, lessee, lessor string, amount float64) (*LeaseTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Input validation
	if leaseAgreementID == "" || lessee == "" || lessor == "" {
		return nil, errors.New("leaseAgreementID, lessee, and lessor must be provided")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(leaseAgreementID + lessee + lessor + fmt.Sprintf("%f", amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt lease transaction data")
	}

	transaction := &LeaseTransaction{
		TransactionID:    generateTransactionID(leaseAgreementID),
		LeaseAgreementID: leaseAgreementID,
		Lessee:           lessee,
		Lessor:           lessor,
		Amount:           amount,
		PaymentDate:      time.Now(),
		NextPaymentDue:   time.Now().AddDate(0, 1, 0), // Assume monthly payments
		EncryptedData:    encryptedData,
		Status:           "pending",
		Fee:              0.01, // Example fee logic
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ProcessTransactions([]common.Transaction{transaction}); err != nil {
		return nil, errors.New("lease transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.recordLeaseTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record lease transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// createPurchaseTransaction handles purchase transactions for assets.
func (tm *TransactionManager) createPurchaseTransaction(assetID, buyer, seller string, amount float64) (*PurchaseTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Input validation
	if assetID == "" || buyer == "" || seller == "" {
		return nil, errors.New("assetID, buyer, and seller must be provided")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(assetID + buyer + seller + fmt.Sprintf("%f", amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt purchase transaction data")
	}

	transaction := &PurchaseTransaction{
		TransactionID: generateTransactionID(assetID),
		AssetID:       assetID,
		Buyer:         buyer,
		Seller:        seller,
		Amount:        amount,
		Timestamp:     time.Now(),
		EncryptedData: encryptedData,
		Status:        "pending",
		Fee:           0.01, // Example fee logic
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ProcessTransactions([]common.Transaction{transaction}); err != nil {
		return nil, errors.New("purchase transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordPurchaseTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record purchase transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}


// CreateOwnershipTransaction creates a new ownership transfer transaction.
func (tm *TransactionManager) CreateOwnershipTransaction(assetID, fromOwner, toOwner string, amount float64) (*OwnershipTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt transaction data
	encryptedData, encryptionKey, err := tm.EncryptionService.EncryptData([]byte(assetID + fromOwner + toOwner + string(amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt transaction data")
	}

	transaction := &OwnershipTransaction{
		TransactionID: generateTransactionID(assetID),
		AssetID:       assetID,
		FromOwner:     fromOwner,
		ToOwner:       toOwner,
		Amount:        amount,
		Timestamp:     time.Now(),
		EncryptedData: encryptedData,
		Status:        "pending",
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateOwnershipTransaction(transaction); err != nil {
		return nil, errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordOwnershipTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record ownership transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// CreateShardedOwnershipTransaction creates a transaction for fractional or sharded ownership transfer.
func (tm *TransactionManager) CreateShardedOwnershipTransaction(assetID string, fromOwners, toOwners map[string]float64, totalAmount float64) (*ShardedOwnershipTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(assetID + string(totalAmount)))
	if err != nil {
		return nil, errors.New("failed to encrypt transaction data")
	}

	transaction := &ShardedOwnershipTransaction{
		TransactionID: generateTransactionID(assetID),
		AssetID:       assetID,
		FromOwners:    fromOwners,
		ToOwners:      toOwners,
		TotalAmount:   totalAmount,
		Timestamp:     time.Now(),
		EncryptedData: encryptedData,
		Status:        "pending",
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateShardedOwnershipTransaction(transaction); err != nil {
		return nil, errors.New("transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordShardedOwnershipTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record sharded ownership transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// CreateRentalTransaction handles rental payment transactions.
func (tm *TransactionManager) CreateRentalTransaction(rentalAgreementID, renter, lessor string, amount float64) (*RentalTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(rentalAgreementID + renter + lessor + string(amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt rental transaction data")
	}

	transaction := &RentalTransaction{
		TransactionID:   generateTransactionID(rentalAgreementID),
		RentalAgreementID: rentalAgreementID,
		Renter:          renter,
		Lessor:          lessor,
		Amount:          amount,
		PaymentDate:     time.Now(),
		NextPaymentDue:  time.Now().AddDate(0, 1, 0), // Assume monthly payments
		EncryptedData:   encryptedData,
		Status:          "pending",
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateRentalTransaction(transaction); err != nil {
		return nil, errors.New("rental transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordRentalTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record rental transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// CreateLeaseTransaction handles lease payment transactions.
func (tm *TransactionManager) CreateLeaseTransaction(leaseAgreementID, lessee, lessor string, amount float64) (*LeaseTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(leaseAgreementID + lessee + lessor + string(amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt lease transaction data")
	}

	transaction := &LeaseTransaction{
		TransactionID:   generateTransactionID(leaseAgreementID),
		LeaseAgreementID: leaseAgreementID,
		Lessee:          lessee,
		Lessor:          lessor,
		Amount:          amount,
		PaymentDate:     time.Now(),
		NextPaymentDue:  time.Now().AddDate(0, 1, 0), // Assume monthly payments
		EncryptedData:   encryptedData,
		Status:          "pending",
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidateLeaseTransaction(transaction); err != nil {
		return nil, errors.New("lease transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordLeaseTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record lease transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// CreatePurchaseTransaction handles purchase transactions for assets.
func (tm *TransactionManager) CreatePurchaseTransaction(assetID, buyer, seller string, amount float64) (*PurchaseTransaction, error) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt transaction data
	encryptedData, _, err := tm.EncryptionService.EncryptData([]byte(assetID + buyer + seller + string(amount)))
	if err != nil {
		return nil, errors.New("failed to encrypt purchase transaction data")
	}

	transaction := &PurchaseTransaction{
		TransactionID: generateTransactionID(assetID),
		AssetID:       assetID,
		Buyer:         buyer,
		Seller:        seller,
		Amount:        amount,
		Timestamp:     time.Now(),
		EncryptedData: encryptedData,
		Status:        "pending",
	}

	// Validate the transaction with Synnergy Consensus
	if err := tm.ConsensusEngine.ValidatePurchaseTransaction(transaction); err != nil {
		return nil, errors.New("purchase transaction validation failed via Synnergy Consensus")
	}

	// Record the transaction in the ledger
	if err := tm.Ledger.RecordPurchaseTransaction(transaction.TransactionID, transaction); err != nil {
		return nil, errors.New("failed to record purchase transaction in the ledger")
	}

	// Mark transaction as complete
	transaction.Status = "completed"
	return transaction, nil
}

// generateTransactionID generates a unique transaction ID based on asset or agreement ID.
func generateTransactionID(referenceID string) string {
	return referenceID + "_" + time.Now().Format("20060102150405")
}
