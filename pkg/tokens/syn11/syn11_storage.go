package syn11

import (
	"errors"
	"sync"
	"time"

  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// Syn11Token defines the central bank digital gilt token (SYN11) structure.
type Syn11Token struct {
	TokenID      string                      // Unique token ID
	Metadata     Syn11Metadata               // Metadata associated with the token
	Issuer       string                      // Issuer of the token (Central Bank or Government Authority)
	Ledger       *ledger.Ledger              // Ledger to track issuance and transfers
	Consensus    *consensus.SynnergyConsensus // Consensus engine for validation
	Compliance   *compliance.KYCAmlService   // KYC/AML Compliance for regulatory requirements
	CentralBank  string                      // Address of the Central Bank (Only this can mint or burn)
	Encrypted    bool                        // Indicates if the token data is encrypted
}

// Syn11Metadata contains metadata about the token.
type Syn11Metadata struct {
	TokenID          string    // Unique identifier for the token
	Name             string    // Full name of the gilt (e.g., "10 Year Central Bank Gilt")
	Symbol           string    // Symbol representing the token (e.g., "CBGILT")
	GiltCode         string    // Internationally recognized code for the gilt
	IssuerName       string    // Name of the issuer
	MaturityDate     time.Time // Maturity date of the gilt
	CouponRate       float64   // Interest rate paid to gilt holders
	CreationDate     time.Time // Timestamp of token creation
	TotalSupply      uint64    // Total supply of tokens
	CirculatingSupply uint64   // Circulating tokens
	LegalCompliance  LegalInfo // Legal and regulatory compliance details
}

// LegalInfo defines legal and compliance information related to the token.
type LegalInfo struct {
	RegulatoryStatus string   // Current regulatory status
	Licenses         []string // Licenses or certifications held
	Compliance       []string // Compliance requirements
}

// Syn11StorageManager handles the storage, retrieval, and updates to Syn11 tokens.
type Syn11StorageManager struct {
	mutex      sync.Mutex
	Ledger     *ledger.Ledger                // Ledger for managing token records
	Consensus  *consensus.SynnergyConsensus  // Consensus engine for validation
	Encryption *encryption.EncryptionService // Encryption service
	Compliance *compliance.KYCAmlService     // Compliance for regulatory and AML/KYC checks
	tokens     map[string]Syn11Token         // In-memory store for Syn11 tokens
}

// NewSyn11StorageManager initializes a new Syn11StorageManager instance.
func NewSyn11StorageManager(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService, complianceService *compliance.KYCAmlService) *Syn11StorageManager {
	return &Syn11StorageManager{
		Ledger:     ledgerInstance,
		Consensus:  consensusEngine,
		Encryption: encryptionService,
		Compliance: complianceService,
		tokens:     make(map[string]Syn11Token),
	}
}

// StoreToken stores a new token in the ledger.
func (ssm *Syn11StorageManager) StoreToken(token Syn11Token) error {
	ssm.mutex.Lock()
	defer ssm.mutex.Unlock()

	// Compliance and Consensus Validation
	if err := ssm.Compliance.VerifyIssuer(token.Issuer); err != nil {
		return err
	}
	if err := ssm.Consensus.ValidateTokenCreation(token.TokenID, token.Issuer); err != nil {
		return err
	}

	// Encrypt sensitive token details
	encryptedName, err := ssm.Encryption.Encrypt([]byte(token.Metadata.Name))
	if err != nil {
		return err
	}
	encryptedSymbol, err := ssm.Encryption.Encrypt([]byte(token.Metadata.Symbol))
	if err != nil {
		return err
	}

	// Update token metadata with encrypted values
	token.Metadata.Name = string(encryptedName)
	token.Metadata.Symbol = string(encryptedSymbol)
	token.Encrypted = true

	// Store the token in the in-memory map and ledger
	ssm.tokens[token.TokenID] = token
	if err := ssm.Ledger.RecordIssuance(token.TokenID, token.Issuer, token.Metadata.TotalSupply); err != nil {
		return err
	}

	return nil
}

// RetrieveToken retrieves a token's details from the ledger.
func (ssm *Syn11StorageManager) RetrieveToken(tokenID string) (*Syn11Token, error) {
	ssm.mutex.Lock()
	defer ssm.mutex.Unlock()

	token, exists := ssm.tokens[tokenID]
	if !exists {
		return nil, errors.New("token not found")
	}

	return &token, nil
}

// BurnToken removes tokens from circulation by updating the ledger and in-memory store.
func (ssm *Syn11StorageManager) BurnToken(tokenID string, amount uint64, burnerID string) error {
	ssm.mutex.Lock()
	defer ssm.mutex.Unlock()

	// Only the central bank can burn tokens
	token, exists := ssm.tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}
	if burnerID != token.CentralBank {
		return errors.New("only the central bank can burn tokens")
	}

	// Consensus Validation for burning
	if err := ssm.Consensus.AuthorizeTokenBurning(tokenID, amount, burnerID); err != nil {
		return err
	}

	// Update the token supply
	token.Metadata.CirculatingSupply -= amount

	// Record the burning in the ledger
	if err := ssm.Ledger.RecordBurning(tokenID, burnerID, amount); err != nil {
		return err
	}

	// Update in-memory token state
	ssm.tokens[tokenID] = token

	return nil
}

// TransferTokens handles token transfers between two accounts.
func (ssm *Syn11StorageManager) TransferTokens(tokenID, from, to string, amount uint64) error {
	ssm.mutex.Lock()
	defer ssm.mutex.Unlock()

	// Validate the existence of the token
	token, exists := ssm.tokens[tokenID]
	if !exists {
		return errors.New("token not found")
	}

	// Consensus Validation for the transfer
	if err := ssm.Consensus.ValidateTransfer(tokenID, from, to, amount); err != nil {
		return err
	}

	// Update the ledger with the transfer details
	if err := ssm.Ledger.TransferTokens(tokenID, from, to, amount); err != nil {
		return err
	}

	return nil
}

// ListTokens returns all tokens stored in the system.
func (ssm *Syn11StorageManager) ListTokens() ([]Syn11Token, error) {
	ssm.mutex.Lock()
	defer ssm.mutex.Unlock()

	var tokenList []Syn11Token
	for _, token := range ssm.tokens {
		tokenList = append(tokenList, token)
	}
	return tokenList, nil
}
