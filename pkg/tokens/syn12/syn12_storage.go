package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

)


// SYN12TokenStorage manages storage operations for SYN12 tokens.
type SYN12TokenStorage struct {
	ledgerManager     *ledger.LedgerManager         // Ledger to track issuance and transactions
	encryptionService *encryption.EncryptionService // Encryption for token data
	consensus         *consensus.SynnergyConsensus  // Consensus engine for validation
	tokenStorage      map[string][]byte             // In-memory token storage (as an example)
	mutex             sync.Mutex                    // Mutex for concurrency
}

// NewSYN12TokenStorage initializes the token storage for SYN12 tokens.
func NewSYN12TokenStorage(ledgerManager *ledger.LedgerManager, encryptionService *encryption.EncryptionService, consensus *consensus.SynnergyConsensus) *SYN12TokenStorage {
	return &SYN12TokenStorage{
		ledgerManager:     ledgerManager,
		encryptionService: encryptionService,
		consensus:         consensus,
		tokenStorage:      make(map[string][]byte), // In-memory map for token storage
	}
}

// StoreToken securely stores a SYN12 token with encryption and consensus validation.
func (s *SYN12TokenStorage) StoreToken(tokenID string, tokenData interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate token data structure
	if tokenID == "" {
		return errors.New("token ID cannot be empty")
	}

	// Validate the token via consensus before storage
	if err := s.consensus.ValidateToken(tokenID); err != nil {
		return fmt.Errorf("token validation failed: %v", err)
	}

	// Convert token data to JSON
	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %v", err)
	}

	// Encrypt the token data before storing
	encryptedData, err := s.encryptionService.Encrypt(tokenJSON)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Store encrypted data in memory (can be updated to use persistent storage)
	s.tokenStorage[tokenID] = encryptedData

	// Record token storage event in the ledger
	if err := s.ledgerManager.RecordStorageEvent(tokenID, common.EventTokenStored); err != nil {
		return fmt.Errorf("failed to log storage event in the ledger: %v", err)
	}

	fmt.Printf("Token with ID %s securely stored.\n", tokenID)
	return nil
}

// RetrieveToken retrieves and decrypts the SYN12 token from storage.
func (s *SYN12TokenStorage) RetrieveToken(tokenID string) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if token exists in storage
	encryptedData, exists := s.tokenStorage[tokenID]
	if !exists {
		return nil, errors.New("token not found in storage")
	}

	// Decrypt the token data
	decryptedData, err := s.encryptionService.Decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token data: %v", err)
	}

	// Unmarshal JSON data into an interface{} (for flexibility)
	var tokenData interface{}
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %v", err)
	}

	fmt.Printf("Token with ID %s retrieved.\n", tokenID)
	return tokenData, nil
}

// DeleteToken removes a SYN12 token from storage and records the event in the ledger.
func (s *SYN12TokenStorage) DeleteToken(tokenID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if the token exists before deletion
	if _, exists := s.tokenStorage[tokenID]; !exists {
		return errors.New("token not found in storage")
	}

	// Remove token from in-memory storage
	delete(s.tokenStorage, tokenID)

	// Record deletion event in the ledger
	if err := s.ledgerManager.RecordStorageEvent(tokenID, common.EventTokenDeleted); err != nil {
		return fmt.Errorf("failed to log deletion event in the ledger: %v", err)
	}

	fmt.Printf("Token with ID %s deleted from storage.\n", tokenID)
	return nil
}

// ListStoredTokens lists all the tokens currently stored in the system.
func (s *SYN12TokenStorage) ListStoredTokens() ([]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Extract token IDs from storage map
	var tokenIDs []string
	for tokenID := range s.tokenStorage {
		tokenIDs = append(tokenIDs, tokenID)
	}

	return tokenIDs, nil
}

// ValidateAndStoreSubBlock integrates with the consensus and ledger to store sub-blocks of transactions.
func (s *SYN12TokenStorage) ValidateAndStoreSubBlock(subBlock common.SubBlock) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate sub-block with consensus before storing
	if err := s.consensus.ValidateSubBlock(subBlock); err != nil {
		return fmt.Errorf("sub-block validation failed: %v", err)
	}

	// Encrypt and store the sub-block data
	subBlockData, err := json.Marshal(subBlock)
	if err != nil {
		return fmt.Errorf("failed to marshal sub-block data: %v", err)
	}

	encryptedSubBlockData, err := s.encryptionService.Encrypt(subBlockData)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block data: %v", err)
	}

	// Store encrypted sub-block in ledger
	if err := s.ledgerManager.StoreSubBlock(encryptedSubBlockData); err != nil {
		return fmt.Errorf("failed to store sub-block in ledger: %v", err)
	}

	// Log the sub-block storage event in the ledger
	if err := s.ledgerManager.RecordStorageEvent(subBlock.BlockID, common.EventSubBlockStored); err != nil {
		return fmt.Errorf("failed to log sub-block storage event: %v", err)
	}

	fmt.Printf("Sub-block with ID %s securely stored.\n", subBlock.BlockID)
	return nil
}

// RetrieveSubBlock retrieves and decrypts a sub-block from storage.
func (s *SYN12TokenStorage) RetrieveSubBlock(blockID string) (common.SubBlock, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve encrypted sub-block from ledger
	encryptedSubBlockData, err := s.ledgerManager.RetrieveSubBlock(blockID)
	if err != nil {
		return common.SubBlock{}, fmt.Errorf("failed to retrieve sub-block: %v", err)
	}

	// Decrypt the sub-block data
	subBlockData, err := s.encryptionService.Decrypt(encryptedSubBlockData)
	if err != nil {
		return common.SubBlock{}, fmt.Errorf("failed to decrypt sub-block data: %v", err)
	}

	// Unmarshal JSON into a sub-block structure
	var subBlock common.SubBlock
	if err := json.Unmarshal(subBlockData, &subBlock); err != nil {
		return common.SubBlock{}, fmt.Errorf("failed to unmarshal sub-block data: %v", err)
	}

	fmt.Printf("Sub-block with ID %s retrieved.\n", blockID)
	return subBlock, nil
}
