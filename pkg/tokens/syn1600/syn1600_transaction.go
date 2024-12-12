package syn1600

import (
	"errors"
	"time"
)

// TransactionManager handles the creation, validation, and execution of transactions involving SYN1600 tokens.
type TransactionManager struct {
	Ledger ledger.Ledger // The blockchain ledger to store transaction data
}

// CreateSYN1600Transaction creates and stores a new transaction for transferring ownership or distributing royalties.
func (tm *TransactionManager) CreateSYN1600Transaction(sender string, recipient string, tokenID string, amount float64, encryptionKey []byte) (*common.SYN1600Transaction, error) {
	// Step 1: Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return nil, errors.New("failed to retrieve token from ledger: " + err.Error())
	}
	
	// Step 2: Decrypt sensitive token data
	err = tm.decryptTokenData(token.(*common.SYN1600Token), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to decrypt token data: " + err.Error())
	}
	
	// Step 3: Ensure the sender has enough ownership rights for the transaction
	if !tm.validateOwnership(token.(*common.SYN1600Token), sender, amount) {
		return nil, errors.New("insufficient ownership rights for the transaction")
	}

	// Step 4: Create a new transaction struct
	transaction := &common.SYN1600Transaction{
		TransactionID:     generateUniqueID(),
		TokenID:           tokenID,
		Sender:            sender,
		Recipient:         recipient,
		Amount:            amount,
		TransactionType:   "OwnershipTransfer",
		TransactionStatus: "Pending",
		TransactionDate:   time.Now(),
	}

	// Step 5: Validate the transaction using Synnergy Consensus
	err = synnergy.ValidateTransaction(transaction.TransactionID, tokenID)
	if err != nil {
		return nil, errors.New("transaction validation failed: " + err.Error())
	}

	// Step 6: Update the token's ownership records after successful validation
	err = tm.updateOwnership(token.(*common.SYN1600Token), sender, recipient, amount)
	if err != nil {
		return nil, errors.New("failed to update token ownership: " + err.Error())
	}

	// Step 7: Store the transaction in the ledger
	err = tm.Ledger.StoreTransaction(transaction.TransactionID, transaction)
	if err != nil {
		return nil, errors.New("failed to store transaction in ledger: " + err.Error())
	}

	// Step 8: Encrypt token data before storing it back into the ledger
	err = tm.encryptTokenData(token.(*common.SYN1600Token), encryptionKey)
	if err != nil {
		return nil, errors.New("failed to encrypt token data: " + err.Error())
	}

	// Step 9: Update the token in the ledger
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return nil, errors.New("failed to update token in ledger: " + err.Error())
	}

	// Step 10: Log the event of a successful transaction
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "TransactionCreated",
		Description: "Ownership transfer transaction successfully created.",
		EventDate:   time.Now(),
		PerformedBy: "TransactionSystem",
	}
	token.(*common.SYN1600Token).AuditTrail = append(token.(*common.SYN1600Token).AuditTrail, eventLog)

	// Return the transaction object
	return transaction, nil
}

// RetrieveSYN1600Transaction retrieves a specific SYN1600Transaction from the ledger.
func (tm *TransactionManager) RetrieveSYN1600Transaction(transactionID string) (*common.SYN1600Transaction, error) {
	// Step 1: Retrieve the transaction from the ledger
	transaction, err := tm.Ledger.GetTransaction(transactionID)
	if err != nil {
		return nil, errors.New("transaction not found in ledger: " + err.Error())
	}

	// Step 2: Log the event of transaction retrieval
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "TransactionRetrieved",
		Description: "SYN1600 transaction retrieved successfully.",
		EventDate:   time.Now(),
		PerformedBy: "TransactionSystem",
	}

	// Return the transaction
	return transaction.(*common.SYN1600Transaction), nil
}

// DistributeRoyalties creates and executes a transaction for royalty distribution.
func (tm *TransactionManager) DistributeRoyalties(tokenID string, amount float64, encryptionKey []byte) error {
	// Step 1: Retrieve the token from the ledger
	token, err := tm.Ledger.GetToken(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token from ledger: " + err.Error())
	}

	// Step 2: Decrypt token data
	err = tm.decryptTokenData(token.(*common.SYN1600Token), encryptionKey)
	if err != nil {
		return errors.New("failed to decrypt token data: " + err.Error())
	}

	// Step 3: Calculate the distribution for each rights holder based on ownership percentage
	distributionLogs := []common.RoyaltyDistributionLog{}
	for _, ownership := range token.(*common.SYN1600Token).OwnershipRights {
		distributionAmount := (ownership.FractionalShare / 100) * amount
		distributionLog := common.RoyaltyDistributionLog{
			DistributionID:   generateUniqueID(),
			RecipientID:      ownership.OwnerID,
			Amount:           distributionAmount,
			DistributionDate: time.Now(),
			DistributionType: "Automatic",
		}
		distributionLogs = append(distributionLogs, distributionLog)
		token.(*common.SYN1600Token).RevenueDistribution = append(token.(*common.SYN1600Token).RevenueDistribution, distributionLog)
	}

	// Step 4: Store updated token with revenue distribution logs in the ledger
	err = tm.encryptTokenData(token.(*common.SYN1600Token), encryptionKey)
	if err != nil {
		return errors.New("failed to encrypt token data: " + err.Error())
	}
	err = tm.Ledger.UpdateToken(tokenID, token)
	if err != nil {
		return errors.New("failed to update token in ledger: " + err.Error())
	}

	// Step 5: Log the event of royalty distribution
	eventLog := common.EventLog{
		EventID:     generateUniqueID(),
		EventType:   "RoyaltyDistributed",
		Description: "Royalties successfully distributed to rights holders.",
		EventDate:   time.Now(),
		PerformedBy: "TransactionSystem",
	}
	token.(*common.SYN1600Token).AuditTrail = append(token.(*common.SYN1600Token).AuditTrail, eventLog)

	return nil
}

// Helper functions

// validateOwnership checks if the sender has enough ownership rights to perform the transaction.
func (tm *TransactionManager) validateOwnership(token *common.SYN1600Token, owner string, amount float64) bool {
	for _, ownership := range token.OwnershipRights {
		if ownership.OwnerID == owner && ownership.FractionalShare >= amount {
			return true
		}
	}
	return false
}

// updateOwnership updates the ownership of the token after a transaction.
func (tm *TransactionManager) updateOwnership(token *common.SYN1600Token, sender, recipient string, amount float64) error {
	// Reduce the sender's share
	for i, ownership := range token.OwnershipRights {
		if ownership.OwnerID == sender {
			token.OwnershipRights[i].FractionalShare -= amount
		}
	}

	// Increase the recipient's share or add them as a new owner
	found := false
	for i, ownership := range token.OwnershipRights {
		if ownership.OwnerID == recipient {
			token.OwnershipRights[i].FractionalShare += amount
			found = true
			break
		}
	}
	if !found {
		newOwnership := common.OwnershipRecord{
			OwnerID:        recipient,
			FractionalShare: amount,
			PurchaseDate:   time.Now(),
		}
		token.OwnershipRights = append(token.OwnershipRights, newOwnership)
	}
	return nil
}

// encryptTokenData encrypts the sensitive fields of the SYN1600Token before storage.
func (tm *TransactionManager) encryptTokenData(token *common.SYN1600Token, encryptionKey []byte) error {
	for i, ownership := range token.OwnershipRights {
		encryptedOwnerID, err := encryption.Encrypt([]byte(ownership.OwnerID), encryptionKey)
		if err != nil {
			return err
		}
		token.OwnershipRights[i].OwnerID = string(encryptedOwnerID)
	}

	for i, log := range token.RevenueDistribution {
		encryptedRecipientID, err := encryption.Encrypt([]byte(log.RecipientID), encryptionKey)
		if err != nil {
			return err
		}
		token.RevenueDistribution[i].RecipientID = string(encryptedRecipientID)
	}

	token.EncryptedMetadata = encryption.EncryptBytes(token.EncryptedMetadata, encryptionKey)
	return nil
}

// decryptTokenData decrypts the sensitive fields of the SYN1600Token after retrieval.
func (tm *TransactionManager) decryptTokenData(token *common.SYN1600Token, decryptionKey []byte) error {
	for i, ownership := range token.OwnershipRights {
		decryptedOwnerID, err := encryption.Decrypt([]byte(ownership.OwnerID), decryptionKey)
		if err != nil {
			return err
		}
		token.OwnershipRights[i].OwnerID = string(decryptedOwnerID)
	}

	for i, log := range token.RevenueDistribution {
		decryptedRecipientID, err := encryption.Decrypt([]byte(log.RecipientID), decryptionKey)
		if err != nil {
			return err
		}
		token.RevenueDistribution[i].RecipientID = string(decryptedRecipientID)
	}

	token.EncryptedMetadata = encryption.DecryptBytes(token.EncryptedMetadata, decryptionKey)
	return nil
}

// generateUniqueID generates a unique ID for transactions and events.
func generateUniqueID() string {
	return "TXN_" + time.Now().Format("20060102150405")
}
