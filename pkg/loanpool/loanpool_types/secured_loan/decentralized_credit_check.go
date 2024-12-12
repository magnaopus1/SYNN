package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/synnergy_consensus"
)

// NewDecentralizedCreditCheck initializes a new decentralized credit check system.
func NewDecentralizedCreditCheck(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.DecentralizedCreditCheck {
	return &common.DecentralizedCreditCheck{
		Ledger:                ledgerInstance,
		Consensus:             consensusEngine,
		WalletSpendingRecords: make(map[string]*common.SpendingRecord),
		CreditScoreDocuments:  make(map[string][]byte),
		EncryptionService:     encryptionService,
	}
}

// TrackWalletSpending checks and records wallet spending and transactions.
func (dcc *common.DecentralizedCreditCheck) TrackWalletSpending(walletAddress string) (*common.SpendingRecord, error) {
	dcc.mutex.Lock()
	defer dcc.mutex.Unlock()

	// Retrieve transactions from the ledger for the given wallet
	transactions, err := dcc.Ledger.GetWalletTransactions(walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions for wallet %s: %v", walletAddress, err)
	}

	// Calculate total spending
	totalSpent := 0.0
	for _, tx := range transactions {
		totalSpent += tx.Amount
	}

	// Create or update the spending record
	spendingRecord, exists := dcc.WalletSpendingRecords[walletAddress]
	if !exists {
		spendingRecord = &common.SpendingRecord{
			WalletAddress: walletAddress,
			TotalSpent:    totalSpent,
			Transactions:  transactions,
			LastUpdated:   time.Now(),
		}
	} else {
		spendingRecord.TotalSpent = totalSpent
		spendingRecord.Transactions = transactions
		spendingRecord.LastUpdated = time.Now()
	}

	dcc.WalletSpendingRecords[walletAddress] = spendingRecord

	// Store the updated record in the ledger
	err = dcc.Ledger.RecordSpending(walletAddress, spendingRecord.TotalSpent, spendingRecord.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to record spending in ledger: %v", err)
	}

	fmt.Printf("Spending record updated for wallet %s. Total spent: %.2f\n", walletAddress, totalSpent)
	return spendingRecord, nil
}

// AttachCreditScore allows a user to attach an encrypted credit score document to their wallet.
func (dcc *common.DecentralizedCreditCheck) AttachCreditScore(walletAddress string, creditScoreDocument []byte, encryptionKey string) error {
	dcc.mutex.Lock()
	defer dcc.mutex.Unlock()

	// Encrypt the credit score document before storing it
	encryptedDocument, err := dcc.EncryptionService.EncryptData(creditScoreDocument, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt credit score document: %v", err)
	}

	dcc.CreditScoreDocuments[walletAddress] = encryptedDocument

	// Store the encrypted document in the ledger
	err = dcc.Ledger.RecordCreditScore(walletAddress, encryptedDocument)
	if err != nil {
		return fmt.Errorf("failed to record encrypted credit score in ledger: %v", err)
	}

	fmt.Printf("Credit score document attached for wallet %s.\n", walletAddress)
	return nil
}

// ViewCreditScore retrieves the attached credit score document for a wallet.
func (dcc *common.DecentralizedCreditCheck) ViewCreditScore(walletAddress string, encryptionKey string) ([]byte, error) {
	dcc.mutex.Lock()
	defer dcc.mutex.Unlock()

	// Retrieve the encrypted document
	encryptedDocument, exists := dcc.CreditScoreDocuments[walletAddress]
	if !exists {
		return nil, errors.New("no credit score document found for this wallet")
	}

	// Decrypt the document before returning it
	decryptedDocument, err := dcc.EncryptionService.DecryptData(encryptedDocument, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credit score document: %v", err)
	}

	return decryptedDocument, nil
}

// ValidateSpendingAndCreditScore uses the Synnergy Consensus mechanism to validate wallet spending and credit score.
func (dcc *common.DecentralizedCreditCheck) ValidateSpendingAndCreditScore(walletAddress string) (bool, error) {
	dcc.mutex.Lock()
	defer dcc.mutex.Unlock()

	// Get the spending record
	spendingRecord, exists := dcc.WalletSpendingRecords[walletAddress]
	if !exists {
		return false, errors.New("spending record not found for this wallet")
	}

	// Get the credit score document (optional)
	_, hasCreditScore := dcc.CreditScoreDocuments[walletAddress]

	// Use Synnergy Consensus to validate both spending and credit score (if available)
	validated, err := dcc.Consensus.ValidateSpending(walletAddress, spendingRecord.TotalSpent)
	if err != nil || !validated {
		return false, fmt.Errorf("failed to validate spending for wallet %s: %v", walletAddress, err)
	}

	// Log the validation result
	if hasCreditScore {
		fmt.Printf("Spending and credit score validated for wallet %s.\n", walletAddress)
	} else {
		fmt.Printf("Spending validated (without credit score) for wallet %s.\n", walletAddress)
	}

	return true, nil
}

// GetSpendingRecord retrieves the spending record for a given wallet.
func (dcc *common.DecentralizedCreditCheck) GetSpendingRecord(walletAddress string) (*common.SpendingRecord, error) {
	dcc.mutex.Lock()
	defer dcc.mutex.Unlock()

	record, exists := dcc.WalletSpendingRecords[walletAddress]
	if !exists {
		return nil, errors.New("spending record not found for this wallet")
	}

	return record, nil
}
