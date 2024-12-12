package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn130TransactionManager handles various transactions for SYN130 tokens.
type Syn130TransactionManager struct {
	TransactionLedger *ledger.TransactionLedger     // Ledger to record all transactions
	OwnershipLedger   *ledger.OwnershipLedger       // Ledger to manage ownership details
	Consensus         *consensus.SynnergyConsensus  // Synnergy consensus engine for validation
	EncryptionService *encryption.EncryptionService // Encryption service for secure transactions
	mutex             sync.Mutex                    // Mutex for safe concurrent transactions
}

// NewSyn130TransactionManager initializes a new transaction manager.
func NewSyn130TransactionManager(transactionLedger *ledger.TransactionLedger, ownershipLedger *ledger.OwnershipLedger, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *Syn130TransactionManager {
	return &Syn130TransactionManager{
		TransactionLedger: transactionLedger,
		OwnershipLedger:   ownershipLedger,
		EncryptionService: encryptionService,
		Consensus:         consensusEngine,
	}
}

// LeasePayment processes a lease payment for an asset.
func (tm *Syn130TransactionManager) LeasePayment(leaseID, lessor, lessee string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the lease transaction using Synnergy Consensus
	if err := tm.Consensus.ValidateLeaseTransaction(leaseID, lessee, amount); err != nil {
		return fmt.Errorf("lease transaction validation failed: %v", err)
	}

	// Record the payment transaction in the ledger
	transaction := ledger.TransactionRecord{
		From:    lessee,
		To:      lessor,
		Amount:  amount,
		Date:    time.Now(),
		Details: fmt.Sprintf("Lease payment for LeaseID: %s", leaseID),
	}
	if err := tm.TransactionLedger.RecordTransaction(&transaction); err != nil {
		return fmt.Errorf("failed to record lease payment: %v", err)
	}

	fmt.Printf("Lease payment of %.2f from %s to %s for LeaseID %s recorded successfully\n", amount, lessee, lessor, leaseID)
	return nil
}

// LicensePayment processes a license fee payment for an asset.
func (tm *Syn130TransactionManager) LicensePayment(licenseID, licensor, licensee string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the license transaction using Synnergy Consensus
	if err := tm.Consensus.ValidateLicenseTransaction(licenseID, licensee, amount); err != nil {
		return fmt.Errorf("license transaction validation failed: %v", err)
	}

	// Record the payment transaction in the ledger
	transaction := ledger.TransactionRecord{
		From:    licensee,
		To:      licensor,
		Amount:  amount,
		Date:    time.Now(),
		Details: fmt.Sprintf("License payment for LicenseID: %s", licenseID),
	}
	if err := tm.TransactionLedger.RecordTransaction(&transaction); err != nil {
		return fmt.Errorf("failed to record license payment: %v", err)
	}

	fmt.Printf("License payment of %.2f from %s to %s for LicenseID %s recorded successfully\n", amount, licensee, licensor, licenseID)
	return nil
}

// RentalPayment processes a rental fee payment for an asset.
func (tm *Syn130TransactionManager) RentalPayment(rentalID, lessor, lessee string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the rental transaction using Synnergy Consensus
	if err := tm.Consensus.ValidateRentalTransaction(rentalID, lessee, amount); err != nil {
		return fmt.Errorf("rental transaction validation failed: %v", err)
	}

	// Record the payment transaction in the ledger
	transaction := ledger.TransactionRecord{
		From:    lessee,
		To:      lessor,
		Amount:  amount,
		Date:    time.Now(),
		Details: fmt.Sprintf("Rental payment for RentalID: %s", rentalID),
	}
	if err := tm.TransactionLedger.RecordTransaction(&transaction); err != nil {
		return fmt.Errorf("failed to record rental payment: %v", err)
	}

	fmt.Printf("Rental payment of %.2f from %s to %s for RentalID %s recorded successfully\n", amount, lessee, lessor, rentalID)
	return nil
}

// OwnershipSharding splits ownership of an asset across multiple owners.
func (tm *Syn130TransactionManager) OwnershipSharding(assetID string, newOwners map[string]float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the sharding using Synnergy Consensus
	if err := tm.Consensus.ValidateOwnershipSharding(assetID, newOwners); err != nil {
		return fmt.Errorf("ownership sharding validation failed: %v", err)
	}

	// Update the ownership ledger with new owners and their respective shares
	for owner, share := range newOwners {
		if err := tm.OwnershipLedger.RecordOwnership(assetID, owner, share); err != nil {
			return fmt.Errorf("failed to record ownership sharding: %v", err)
		}
	}

	fmt.Printf("Ownership of asset %s successfully sharded among new owners\n", assetID)
	return nil
}

// FullPurchase transfers full ownership of an asset from the seller to the buyer.
func (tm *Syn130TransactionManager) FullPurchase(assetID, seller, buyer string, amount float64) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Validate the purchase using Synnergy Consensus
	if err := tm.Consensus.ValidateFullPurchaseTransaction(assetID, seller, buyer, amount); err != nil {
		return fmt.Errorf("full purchase transaction validation failed: %v", err)
	}

	// Record the full ownership transfer in the ownership ledger
	if err := tm.OwnershipLedger.TransferFullOwnership(assetID, seller, buyer); err != nil {
		return fmt.Errorf("failed to record full ownership transfer: %v", err)
	}

	// Record the payment transaction in the transaction ledger
	transaction := ledger.TransactionRecord{
		From:    buyer,
		To:      seller,
		Amount:  amount,
		Date:    time.Now(),
		Details: fmt.Sprintf("Full purchase of asset %s from %s to %s", assetID, seller, buyer),
	}
	if err := tm.TransactionLedger.RecordTransaction(&transaction); err != nil {
		return fmt.Errorf("failed to record full purchase payment: %v", err)
	}

	fmt.Printf("Full purchase of asset %s from %s to %s for %.2f recorded successfully\n", assetID, seller, buyer, amount)
	return nil
}

// EncryptTransaction encrypts the transaction details before recording.
func (tm *Syn130TransactionManager) EncryptTransaction(transaction *ledger.TransactionRecord) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Encrypt the transaction details
	encryptedDetails, err := tm.EncryptionService.EncryptData(transaction.Details)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction: %v", err)
	}

	// Update the transaction with encrypted details
	transaction.Details = encryptedDetails
	return nil
}

// DecryptTransaction decrypts the transaction details when required.
func (tm *Syn130TransactionManager) DecryptTransaction(transaction *ledger.TransactionRecord) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	// Decrypt the transaction details
	decryptedDetails, err := tm.EncryptionService.DecryptData(transaction.Details)
	if err != nil {
		return fmt.Errorf("failed to decrypt transaction: %v", err)
	}

	// Update the transaction with decrypted details
	transaction.Details = decryptedDetails
	return nil
}
