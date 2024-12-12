package syn2900

import (

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"time"
	"sync"
)

// SYN2900Token represents an insurance policy token with its associated metadata, coverage details, and claims history.
type SYN2900Token struct {
	TokenID         string            `json:"token_id"`          // Unique ID for the token
	PolicyID        string            `json:"policy_id"`         // Associated policy ID
	Issuer          string            `json:"issuer"`            // The insurance company or entity issuing the policy
	Owner           string            `json:"owner"`             // Current owner or policyholder
	Coverages       []CoverageDetail  `json:"coverages"`         // Details of what is covered under this policy
	Premium         float64           `json:"premium"`           // The premium cost of the policy
	StartDate       time.Time         `json:"start_date"`        // Start date of the policy
	EndDate         time.Time         `json:"end_date"`          // End date of the policy
	ActiveStatus    bool              `json:"active_status"`     // Whether the policy is currently active or expired
	TransactionLogs []TransactionLog  `json:"transaction_logs"`  // History of transactions related to the token
	Claims          []ClaimRecord     `json:"claims"`            // Claims filed under this policy
}

// CoverageDetail provides details about the different risks or areas covered under the insurance policy.
type CoverageDetail struct {
	CoverageType  string  `json:"coverage_type"`   // Type of coverage (e.g., Property, Health, Cybersecurity)
	CoverageLimit float64 `json:"coverage_limit"`  // Maximum amount that can be claimed under this coverage type
	Deductible    float64 `json:"deductible"`      // Amount that must be paid out-of-pocket before the insurance kicks in
}

// TransactionLog keeps a record of all transactions and updates related to the SYN2900 token.
type TransactionLog struct {
	TransactionID string    `json:"transaction_id"`
	Timestamp     time.Time `json:"timestamp"`
	Details       string    `json:"details"`
}

// ClaimRecord represents a claim filed against the insurance policy.
type ClaimRecord struct {
	ClaimID      string    `json:"claim_id"`      // Unique claim identifier
	ClaimDate    time.Time `json:"claim_date"`    // Date the claim was filed
	ClaimAmount  float64   `json:"claim_amount"`  // Amount claimed
	ClaimStatus  string    `json:"claim_status"`  // Status of the claim (e.g., Pending, Approved, Rejected)
	PayoutAmount float64   `json:"payout_amount"` // Amount paid out on the claim
}

// SYN2900Claim represents a claim interaction for the SYN2900Token standard.
type SYN2900Claim struct {
	ClaimID     string    `json:"claim_id"`
	TokenID     string    `json:"token_id"`
	ClaimAmount float64   `json:"claim_amount"`
	ClaimDate   time.Time `json:"claim_date"`
	Status      string    `json:"status"`
}

// SYN2900PolicyUpdate captures any modifications made to an active insurance policy.
type SYN2900PolicyUpdate struct {
	TokenID   string    `json:"token_id"`
	UpdateID  string    `json:"update_id"`
	UpdatedBy string    `json:"updated_by"`
	UpdateAt  time.Time `json:"updated_at"`
	Changes   []string  `json:"changes"`
}

// InsuranceFactory manages the creation and issuance of SYN2900 insurance tokens.
type InsuranceFactory struct {
	tokenStore map[string]common.SYN2900Token
	mu         sync.Mutex
}

// NewInsuranceFactory creates a new instance of InsuranceFactory.
func NewInsuranceFactory() *InsuranceFactory {
	return &InsuranceFactory{
		tokenStore: make(map[string]common.SYN2900Token),
	}
}

// CreateInsuranceToken generates a new SYN2900 insurance token.
func (f *InsuranceFactory) CreateInsuranceToken(
	policyID, issuer, owner string,
	coverages []common.CoverageDetail,
	premium float64,
	startDate, endDate time.Time,
) (string, error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	tokenID := generateTokenID()
	token := common.SYN2900Token{
		TokenID:         tokenID,
		PolicyID:        policyID,
		Issuer:          issuer,
		Owner:           owner,
		Coverages:       coverages,
		Premium:         premium,
		StartDate:       startDate,
		EndDate:         endDate,
		ActiveStatus:    true,
		TransactionLogs: []common.TransactionLog{},
		Claims:          []common.ClaimRecord{},
	}

	// Encrypt token data
	encryptedToken, err := encryptTokenData(token)
	if err != nil {
		return "", err
	}

	// Save token to the ledger
	err = ledger.StoreToken(encryptedToken)
	if err != nil {
		return "", err
	}

	// Store in factory's memory
	f.tokenStore[tokenID] = token

	return tokenID, nil
}

// GetInsuranceToken retrieves an insurance token by its token ID.
func (f *InsuranceFactory) GetInsuranceToken(tokenID string) (common.SYN2900Token, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	token, exists := f.tokenStore[tokenID]
	if !exists {
		return common.SYN2900Token{}, errors.New("token not found")
	}

	return token, nil
}

// Encrypt and decrypt functions using AES encryption
var encryptionKey = []byte("your-encryption-key-32-bytes-long") // Replace with secure key

// encryptTokenData encrypts the token data for secure storage.
func encryptTokenData(token common.SYN2900Token) (string, error) {
	plaintext, err := common.SerializeStruct(token)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptTokenData decrypts the token data for use.
func decryptTokenData(encrypted string) (common.SYN2900Token, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return common.SYN2900Token{}, errors.New("invalid ciphertext")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	var token common.SYN2900Token
	err = common.DeserializeStruct(plaintext, &token)
	if err != nil {
		return common.SYN2900Token{}, err
	}

	return token, nil
}

// Generate a unique token ID for the SYN2900Token.
func generateTokenID() string {
	return "SYN2900-" + common.GenerateUniqueID()
}

// LogTransaction logs a transaction in the insurance token's transaction history.
func (f *InsuranceFactory) LogTransaction(tokenID, details string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	token, exists := f.tokenStore[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	log := common.TransactionLog{
		TransactionID: common.GenerateUniqueID(),
		Timestamp:     time.Now(),
		Details:       details,
	}

	token.TransactionLogs = append(token.TransactionLogs, log)
	f.tokenStore[tokenID] = token

	// Save the updated token to the ledger
	encryptedToken, err := encryptTokenData(token)
	if err != nil {
		return err
	}

	err = ledger.UpdateToken(tokenID, encryptedToken)
	if err != nil {
		return err
	}

	return nil
}
