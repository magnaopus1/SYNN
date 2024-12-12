package syn11

import (
	
	"fmt"
	"log"
	"sync"
	"time"

  	"synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// SYN11Token represents the core structure for Central Bank Digital Gilt tokens.
type SYN11Token struct {
	TokenID      string                      // Unique token ID
	Metadata     Syn11Metadata               // Metadata associated with the token
	Issuer       string                      // Issuer of the token (Central Bank or Government Authority)
	Ledger       *ledger.Ledger              // Ledger to track issuance and transfers
	Consensus    *consensus.SynnergyConsensus // Consensus engine for validation
	Compliance   *compliance.KYCAmlService   // KYC/AML Compliance for regulatory requirements
	CentralBank  string                      // Address of the Central Bank (Only this can mint or burn)
	Encrypted    bool                        // Indicates if the token data is encrypted
}

// Syn11Metadata defines the metadata for SYN11 digital gilt tokens.
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

// LegalInfo holds the regulatory and legal status of SYN11 gilts.
type LegalInfo struct {
	RegulatoryStatus string   // Current regulatory status
	Licenses         []string // List of licenses or certifications
	Compliance       []string // Compliance requirements
}

// TokenFactory manages the creation, issuance, and burning of SYN11 tokens.
type TokenFactory struct {
	mutex          sync.Mutex
	Ledger         *ledger.Ledger                // Ledger for managing token records
	Consensus      *consensus.SynnergyConsensus  // Consensus engine for validation
	Encryption     *encryption.EncryptionService // Encryption service
	Compliance     *compliance.KYCAmlService     // KYC/AML Compliance service
	CentralBank    string                        // Address of the Central Bank for token minting/burning
}

// NewTokenFactory initializes the TokenFactory.
func NewTokenFactory(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService, complianceService *compliance.KYCAmlService, centralBank string) *TokenFactory {
	return &TokenFactory{
		Ledger:      ledgerInstance,
		Consensus:   consensusEngine,
		Encryption:  encryptionService,
		Compliance:  complianceService,
		CentralBank: centralBank,
	}
}

// IssueToken handles the issuance of new SYN11 tokens.
func (factory *TokenFactory) IssueToken(name, symbol, giltCode, issuerID string, amount uint64, maturityDate time.Time, couponRate float64) (string, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// KYC/AML Compliance Check
	if err := factory.Compliance.VerifyIssuer(issuerID); err != nil {
		return "", fmt.Errorf("KYC/AML verification failed: %w", err)
	}

	// Compliance and Security Checks
	if err := factory.Consensus.ValidateIssuer(issuerID); err != nil {
		return "", fmt.Errorf("issuer validation failed: %w", err)
	}

	if err := factory.Consensus.AuthorizeTokenIssuance(issuerID, amount); err != nil {
		return "", fmt.Errorf("authorization failed: %w", err)
	}

	// Create and Register the Token
	tokenID := fmt.Sprintf("SYN11-%d-%s", time.Now().UnixNano(), giltCode)
	token := SYN11Token{
		TokenID: tokenID,
		Metadata: Syn11Metadata{
			TokenID:          tokenID,
			Name:             name,
			Symbol:           symbol,
			GiltCode:         giltCode,
			IssuerName:       issuerID,
			MaturityDate:     maturityDate,
			CouponRate:       couponRate,
			CreationDate:     time.Now(),
			TotalSupply:      amount,
			CirculatingSupply: amount,
		},
		Issuer: issuerID,
		Ledger: factory.Ledger,
		Consensus: factory.Consensus,
		Compliance: factory.Compliance,
		CentralBank: factory.CentralBank,
	}

	// Store the token in the ledger
	if err := factory.Ledger.RecordIssuance(token.TokenID, issuerID, amount); err != nil {
		return "", fmt.Errorf("ledger recording failed: %w", err)
	}

	log.Printf("Successfully issued token with ID: %s, Name: %s, Symbol: %s", tokenID, name, symbol)
	return tokenID, nil
}

// BurnToken handles the burning of SYN11 tokens, removing them from circulation.
func (factory *TokenFactory) BurnToken(tokenID string, amount uint64, burnerID string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Only the central bank can burn tokens
	if burnerID != factory.CentralBank {
		return fmt.Errorf("only the central bank can burn tokens")
	}

	// Verify Ownership and Compliance
	if err := factory.Consensus.VerifyOwnership(tokenID, burnerID, amount); err != nil {
		return fmt.Errorf("ownership verification failed: %w", err)
	}

	// Security and Authorization
	if err := factory.Consensus.AuthorizeTokenBurning(tokenID, amount, burnerID); err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}

	// Burn the Tokens in the Ledger
	if err := factory.Ledger.BurnTokens(tokenID, amount); err != nil {
		return fmt.Errorf("burning tokens failed: %w", err)
	}

	log.Printf("Successfully burned %d tokens of ID: %s", amount, tokenID)
	return nil
}

// TransferOwnership handles the transfer of token ownership.
func (factory *TokenFactory) TransferOwnership(tokenID, fromID, toID string, amount uint64) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Verify Ownership and Compliance
	if err := factory.Consensus.VerifyOwnership(tokenID, fromID, amount); err != nil {
		return fmt.Errorf("ownership verification failed: %w", err)
	}

	// Execute the transfer in the ledger
	if err := factory.Ledger.TransferTokens(tokenID, fromID, toID, amount); err != nil {
		return fmt.Errorf("token transfer failed: %w", err)
	}

	log.Printf("Successfully transferred %d tokens of ID: %s from %s to %s", amount, tokenID, fromID, toID)
	return nil
}
