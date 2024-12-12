package syn3100

import (
	"errors"
	"time"
	"sync"

)

// SYN3100Token represents an employment contract token.
type SYN3100Token struct {
	TokenID         string          `json:"token_id"`
	ContractMetadata ContractMetadata `json:"contract_metadata"`
	IssueDate       time.Time       `json:"issue_date"`
	mutex           sync.Mutex
	ledgerService   *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
}

// NewSYN3100Token creates a new instance of an employment contract token.
func NewSYN3100Token(tokenID string, metadata ContractMetadata, ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *SYN3100Token {
	return &SYN3100Token{
		TokenID:         tokenID,
		ContractMetadata: metadata,
		IssueDate:       time.Now(),
		ledgerService:   ledger,
		encryptionService: encryptor,
		consensusService: consensus,
	}
}

// ContractMetadata represents the metadata of an employment contract.
type ContractMetadata struct {
	ContractID    string    `json:"contract_id"`
	EmployeeID    string    `json:"employee_id"`
	EmployerID    string    `json:"employer_id"`
	Position      string    `json:"position"`
	Salary        float64   `json:"salary"`
	ContractType  string    `json:"contract_type"`
	StartDate     time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Benefits      []string  `json:"benefits"`
	ContractTerms string    `json:"contract_terms"`
	ActiveStatus  bool      `json:"active_status"`
	CompanyDetails CompanyDetails `json:"company_details"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CompanyDetails represents details about the company offering the contract.
type CompanyDetails struct {
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
	Location    string `json:"location"`
	Industry    string `json:"industry"`
}

// OwnershipVerification represents the ownership verification logic for employment contracts.
type OwnershipVerification struct {
	ContractID        string    `json:"contract_id"`
	EmployeeID        string    `json:"employee_id"`
	Verified          bool      `json:"verified"`
	VerifiedAt        time.Time `json:"verified_at"`
	VerificationToken string    `json:"verification_token"`
}

// TokenFactory manages the creation and operations of SYN3100 tokens.
type TokenFactory struct {
	ledgerService    *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex            sync.Mutex
	tokens           map[string]*SYN3100Token
}

// NewTokenFactory creates a new TokenFactory instance.
func NewTokenFactory(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *TokenFactory {
	return &TokenFactory{
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
		tokens:           make(map[string]*SYN3100Token),
	}
}

// CreateContractToken creates a new SYN3100 employment contract token.
func (tf *TokenFactory) CreateContractToken(contractID, employeeID, employerID, position string, salary float64, contractType string, startDate time.Time, endDate *time.Time, benefits []string, contractTerms string, companyDetails CompanyDetails) (*SYN3100Token, error) {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Generate token ID (or use the contract ID as the token ID).
	tokenID := generateTokenID(contractID)

	// Create contract metadata.
	metadata := ContractMetadata{
		ContractID:   contractID,
		EmployeeID:   employeeID,
		EmployerID:   employerID,
		Position:     position,
		Salary:       salary,
		ContractType: contractType,
		StartDate:    startDate,
		EndDate:      endDate,
		Benefits:     benefits,
		ContractTerms: contractTerms,
		ActiveStatus:  true,
		CompanyDetails: companyDetails,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create the SYN3100 token.
	token := NewSYN3100Token(tokenID, metadata, tf.ledgerService, tf.encryptionService, tf.consensusService)

	// Encrypt the token before storing.
	encryptedToken, err := tf.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}

	// Log the creation in the ledger.
	if err := tf.ledgerService.LogEvent("ContractTokenCreated", time.Now(), tokenID); err != nil {
		return nil, err
	}

	// Store the token in the factory.
	tf.tokens[tokenID] = encryptedToken.(*SYN3100Token)

	// Validate the token with Synnergy Consensus.
	if err := tf.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	return token, nil
}

// VerifyOwnership verifies the ownership of a contract token.
func (tf *TokenFactory) VerifyOwnership(contractID, employeeID string) (*OwnershipVerification, error) {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Retrieve the token by contract ID.
	token, err := tf.retrieveToken(contractID)
	if err != nil {
		return nil, err
	}

	// Verify the ownership.
	verified := token.ContractMetadata.EmployeeID == employeeID

	// Generate verification token.
	verificationToken := generateVerificationToken(contractID)

	// Create ownership verification.
	ownershipVerification := &OwnershipVerification{
		ContractID:        contractID,
		EmployeeID:        employeeID,
		Verified:          verified,
		VerifiedAt:        time.Now(),
		VerificationToken: verificationToken,
	}

	// Log the verification in the ledger.
	if err := tf.ledgerService.LogEvent("OwnershipVerified", time.Now(), contractID); err != nil {
		return nil, err
	}

	// Validate the verification with Synnergy Consensus.
	if err := tf.consensusService.ValidateSubBlock(contractID); err != nil {
		return nil, err
	}

	return ownershipVerification, nil
}

// retrieveToken is a helper function to retrieve a token by its ID.
func (tf *TokenFactory) retrieveToken(contractID string) (*SYN3100Token, error) {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	token, exists := tf.tokens[contractID]
	if !exists {
		return nil, errors.New("contract token not found")
	}

	// Decrypt the token before returning.
	decryptedToken, err := tf.encryptionService.DecryptData(token)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*SYN3100Token), nil
}

// Helper function to generate a unique token ID based on the contract ID.
func generateTokenID(contractID string) string {
	return "TOKEN_" + contractID
}

// Helper function to generate a unique ownership verification token.
func generateVerificationToken(contractID string) string {
	return "VER_" + contractID + "_" + time.Now().Format("20060102150405")
}
