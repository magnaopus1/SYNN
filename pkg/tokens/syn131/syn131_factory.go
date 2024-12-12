package syn131

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewSyn131Factory initializes a new Syn131Factory instance
func NewSyn131Factory(ledger *ledger.SYN131Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *Syn131Factory {
	return &Syn131Factory{
		Ledger:            ledger,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
	}
}


// CreateToken creates a new SYN131 token for intangible assets
func (sf *Syn131Factory) CreateToken(name, owner string, metadata IntangibleAssetMetadata, terms, contractType string) (*Syn131Token, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	if metadata.AssetID == "" {
		return nil, errors.New("metadata must contain a valid asset ID")
	}

	tokenID := GenerateTokenID(name, owner)
	encryptedTerms, encryptionKey, err := sf.EncryptionService.EncryptData([]byte(terms))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt contract terms: %w", err)
	}

	token := &Syn131Token{
		ID:                             tokenID,
		Name:                           name,
		Owner:                          owner,
		IntangibleAssetID:              metadata.AssetID,
		ContractType:                   contractType,
		Terms:                          terms,
		EncryptedTerms:                 encryptedTerms,
		EncryptionKey:                  encryptionKey,
		Status:                         "active",
		IntangibleAssetMetadata:        metadata,
		IntangibleAssetCategory:        metadata.AssetType,
		IntangibleAssetClassification:  metadata.Description,
		CreatedAt:                      time.Now(),
		UpdatedAt:                      time.Now(),
	}

	if err := sf.ConsensusEngine.ValidateTokenCreation(token); err != nil {
		return nil, fmt.Errorf("token creation validation failed: %w", err)
	}

	if err := sf.Ledger.RecordTokenCreation(tokenID, token); err != nil {
		return nil, fmt.Errorf("failed to record token creation: %w", err)
	}

	transaction := SYN131Transaction{
		TransactionID: GenerateTokenID("txn", name),
		TokenID:       tokenID,
		OperationType: "Create",
		PerformedBy:   owner,
		Timestamp:     time.Now(),
		Status:        "Success",
		Details:       "Token created successfully",
	}
	sf.Ledger.RecordTransaction(transaction)

	return token, nil
}

// UpdateTokenTerms updates the terms of a SYN131 token, including re-encrypting the updated terms
func (sf *Syn131Factory) UpdateTokenTerms(tokenID, newTerms string) error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	token, err := sf.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %w", err)
	}

	encryptedTerms, encryptionKey, err := sf.EncryptionService.EncryptData([]byte(newTerms))
	if err != nil {
		return fmt.Errorf("failed to encrypt new terms: %w", err)
	}

	token.EncryptedTerms = encryptedTerms
	token.EncryptionKey = encryptionKey
	token.UpdatedAt = time.Now()

	if err := sf.ConsensusEngine.ValidateTokenUpdate(token); err != nil {
		return fmt.Errorf("token update validation failed: %w", err)
	}

	if err := sf.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	transaction := SYN131Transaction{
		TransactionID: GenerateTokenID("txn", tokenID),
		TokenID:       tokenID,
		OperationType: "Update",
		PerformedBy:   token.Owner,
		Timestamp:     time.Now(),
		Status:        "Success",
		Details:       "Token terms updated",
	}
	sf.Ledger.RecordTransaction(transaction)

	return nil
}


// TransferTokenOwnership transfers the ownership of a SYN131 token to a new owner
func (sf *Syn131Factory) TransferTokenOwnership(tokenID, newOwner string) error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	token, err := sf.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %w", err)
	}

	token.Owner = newOwner
	token.UpdatedAt = time.Now()

	if err := sf.ConsensusEngine.ValidateTokenTransfer(token); err != nil {
		return fmt.Errorf("token transfer validation failed: %w", err)
	}

	if err := sf.Ledger.UpdateToken(tokenID, token); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	transaction := SYN131Transaction{
		TransactionID: GenerateTokenID("txn", tokenID),
		TokenID:       tokenID,
		OperationType: "Transfer",
		PerformedBy:   newOwner,
		Timestamp:     time.Now(),
		Status:        "Success",
		Details:       "Token ownership transferred",
	}
	sf.Ledger.RecordTransaction(transaction)

	return nil
}


// generateTokenID creates a unique token ID based on name and owner
func GenerateTokenID(name, owner string) string {
	// Simple token ID generation logic; this could be replaced with more complex logic
	return name + "_" + owner + "_" + time.Now().Format("20060102150405")
}
