package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

)


// SYN12SecurityManager manages the security for SYN12 tokens, including encryption, transaction authorization, and validation.
type SYN12SecurityManager struct {
	encryptionService *encryption.EncryptionService // Encryption service for token security
	ledgerManager     *ledger.LedgerManager         // Ledger manager for transaction recording and validation
	consensus         *consensus.SynnergyConsensus  // Consensus engine for validating security events and transactions
	mutex             sync.Mutex                    // Mutex for concurrency control
}

// NewSYN12SecurityManager initializes the security manager with ledger, encryption, and consensus services.
func NewSYN12SecurityManager(encryptionService *encryption.EncryptionService, ledgerManager *ledger.LedgerManager, consensus *consensus.SynnergyConsensus) *SYN12SecurityManager {
	return &SYN12SecurityManager{
		encryptionService: encryptionService,
		ledgerManager:     ledgerManager,
		consensus:         consensus,
	}
}

// EncryptData encrypts data using AES-GCM encryption.
func (sm *SYN12SecurityManager) EncryptData(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(sm.encryptionService.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptData decrypts AES-GCM encrypted data.
func (sm *SYN12SecurityManager) DecryptData(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(sm.encryptionService.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
}

// AuthorizeTokenIssuance checks if a token issuance request is authorized based on the issuer's role and the ledger.
func (sm *SYN12SecurityManager) AuthorizeTokenIssuance(issuerID string, amount uint64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the issuer using consensus
	if err := sm.consensus.ValidateIssuer(issuerID); err != nil {
		return fmt.Errorf("issuer validation failed: %v", err)
	}

	// Authorize token issuance based on ledger rules
	if err := sm.ledgerManager.ValidateIssuance(issuerID, amount); err != nil {
		return fmt.Errorf("token issuance authorization failed: %v", err)
	}

	return nil
}

// AuthorizeTokenBurning ensures only authorized entities (e.g., Central Bank) can burn SYN12 tokens.
func (sm *SYN12SecurityManager) AuthorizeTokenBurning(tokenID string, amount uint64, burnerID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Ensure the burner is authorized to burn tokens
	if !sm.consensus.IsCentralBank(burnerID) {
		return errors.New("only the central bank is authorized to burn tokens")
	}

	// Check ledger for token ownership and balance
	if err := sm.ledgerManager.VerifyOwnership(tokenID, burnerID, amount); err != nil {
		return fmt.Errorf("burner is not authorized: %v", err)
	}

	// Consensus validation for burning tokens
	if err := sm.consensus.ValidateBurning(tokenID, amount, burnerID); err != nil {
		return fmt.Errorf("consensus validation failed for burning: %v", err)
	}

	return nil
}

// AuthorizeRedemption checks if a token redemption request is authorized.
func (sm *SYN12SecurityManager) AuthorizeRedemption(tokenID string, redeemerID string, amount uint64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the redeemer using consensus
	if err := sm.consensus.ValidateRedeemer(tokenID, redeemerID); err != nil {
		return fmt.Errorf("redeemer validation failed: %v", err)
	}

	// Check ledger for token ownership and balance
	if err := sm.ledgerManager.VerifyOwnership(tokenID, redeemerID, amount); err != nil {
		return fmt.Errorf("redeemer does not own enough tokens: %v", err)
	}

	return nil
}

// SecureTransaction ensures that a transaction meets all security requirements before execution.
func (sm *SYN12SecurityManager) SecureTransaction(tokenID, senderID, receiverID string, amount uint64) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate transaction using consensus
	if err := sm.consensus.ValidateTransaction(tokenID, senderID, receiverID, amount); err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Verify ownership and ensure sender has enough tokens
	if err := sm.ledgerManager.VerifyOwnership(tokenID, senderID, amount); err != nil {
		return fmt.Errorf("sender does not have enough tokens: %v", err)
	}

	return nil
}

// LogSecurityEvent logs security-related events in the ledger.
func (sm *SYN12SecurityManager) LogSecurityEvent(eventType, description string) error {
	eventID := common.GenerateUUID()
	eventDetails := fmt.Sprintf("Security Event - Type: %s, Description: %s", eventType, description)

	// Encrypt the event before logging for security
	encryptedEvent, err := sm.EncryptData([]byte(eventDetails))
	if err != nil {
		return fmt.Errorf("failed to encrypt security event: %v", err)
	}

	// Record the encrypted event in the ledger
	if err := sm.ledgerManager.RecordEvent(eventID, eventType, string(encryptedEvent)); err != nil {
		return fmt.Errorf("failed to log security event: %v", err)
	}

	return nil
}
