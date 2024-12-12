package syn2800

import (
    "errors"
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "sync"
)

// EventManager manages all events related to the SYN2800 life insurance tokens.
type EventManager struct {
	mutex sync.Mutex
}

// NewEventManager creates a new instance of EventManager.
func NewEventManager() *EventManager {
	return &EventManager{}
}

// TriggerEvent triggers an event related to a life insurance policy token.
func (em *EventManager) TriggerEvent(tokenID string, eventType string, eventData interface{}) error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Retrieve and decrypt the token
	token, err := em.retrieveAndDecryptToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token for event trigger: %v", err)
	}

	// Process the event
	switch eventType {
	case "premiumPayment":
		err = em.handlePremiumPayment(token, eventData)
	case "policyUpdate":
		err = em.handlePolicyUpdate(token, eventData)
	case "claimInitiated":
		err = em.handleClaimInitiated(token, eventData)
	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}

	if err != nil {
		return fmt.Errorf("error handling event: %v", err)
	}

	// Log event in token’s event records
	token.EventLogs = append(token.EventLogs, common.EventLog{
		EventType:   eventType,
		EventData:   eventData,
		Timestamp:   time.Now(),
	})

	// Encrypt and store the updated token back into the ledger
	encryptedTokenData, err := em.encryptTokenData(token)
	if err != nil {
		return fmt.Errorf("failed to encrypt token data after event handling: %v", err)
	}
	if err := ledger.StoreToken(tokenID, encryptedTokenData); err != nil {
		return fmt.Errorf("failed to store updated token in ledger after event handling: %v", err)
	}

	log.Printf("Event %s triggered for token %s", eventType, tokenID)
	return nil
}

// handlePremiumPayment processes premium payments for a SYN2800 life insurance token.
func (em *EventManager) handlePremiumPayment(token *common.SYN2800Token, eventData interface{}) error {
	// Extract relevant data from the eventData
	data, ok := eventData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data for premium payment")
	}

	paymentAmount, ok := data["amount"].(float64)
	if !ok {
		return fmt.Errorf("invalid payment amount in event data")
	}

	// Update premium records and token balance
	token.PremiumPaid += paymentAmount
	token.LastPaymentDate = time.Now()

	// Ensure the policy is still active
	if token.PremiumPaid < token.RequiredPremium {
		return fmt.Errorf("premium payment insufficient to keep policy active")
	}
	token.ActiveStatus = true

	log.Printf("Premium payment of %f processed for token %s", paymentAmount, token.TokenID)
	return nil
}

// handlePolicyUpdate processes updates to a SYN2800 life insurance policy.
func (em *EventManager) handlePolicyUpdate(token *common.SYN2800Token, eventData interface{}) error {
	// Update the token metadata with new policy details from eventData
	updateData, ok := eventData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data for policy update")
	}

	// Example: Updating the coverage amount
	if coverageAmount, ok := updateData["coverageAmount"].(float64); ok {
		token.CoverageAmount = coverageAmount
	}

	// Example: Updating the beneficiary details
	if newBeneficiary, ok := updateData["beneficiary"].(string); ok {
		token.Beneficiary = newBeneficiary
	}

	log.Printf("Policy update processed for token %s", token.TokenID)
	return nil
}

// handleClaimInitiated handles the initiation of a claim against a SYN2800 life insurance policy.
func (em *EventManager) handleClaimInitiated(token *common.SYN2800Token, eventData interface{}) error {
	// Validate the claim initiation event data
	claimData, ok := eventData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data for claim initiation")
	}

	// Ensure the claim is valid based on policy conditions
	claimAmount, ok := claimData["claimAmount"].(float64)
	if !ok {
		return fmt.Errorf("invalid claim amount in event data")
	}

	// Ensure the claim amount does not exceed the coverage amount
	if claimAmount > token.CoverageAmount {
		return fmt.Errorf("claim amount exceeds policy coverage for token %s", token.TokenID)
	}

	// Initiate the claim and update the token’s claim records
	token.Claims = append(token.Claims, common.ClaimRecord{
		ClaimAmount: claimAmount,
		ClaimDate:   time.Now(),
		ClaimStatus: "Pending",
	})

	log.Printf("Claim initiation processed for token %s with claim amount %f", token.TokenID, claimAmount)
	return nil
}

// Encrypt and store token data in the ledger.
func (em *EventManager) encryptTokenData(token *common.SYN2800Token) ([]byte, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	tokenData := serializeTokenData(token)
	return gcm.Seal(nonce, nonce, tokenData, nil), nil
}

// Decrypt the token data from the ledger.
func (em *EventManager) decryptTokenData(encryptedData []byte) (*common.SYN2800Token, error) {
	key := generateEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return deserializeTokenData(decryptedData), nil
}

// Helper to retrieve and decrypt the token from the ledger.
func (em *EventManager) retrieveAndDecryptToken(tokenID string) (*common.SYN2800Token, error) {
	encryptedData, err := ledger.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}
	return em.decryptTokenData(encryptedData)
}

// Helper function to generate an encryption key.
func generateEncryptionKey() []byte {
	return []byte("your-secure-256-bit-key")
}

// Helper function to serialize token data.
func serializeTokenData(token *common.SYN2800Token) []byte {
	data, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("failed to serialize token data: %v", err)
	}
	return data
}

// Helper function to deserialize token data after decryption.
func deserializeTokenData(data []byte) *common.SYN2800Token {
	var token common.SYN2800Token
	if err := json.Unmarshal(data, &token); err != nil {
		log.Fatalf("failed to deserialize token data: %v", err)
	}
	return &token
}
