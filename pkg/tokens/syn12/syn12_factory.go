package syn12

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn12Metadata defines the metadata for SYN12 Treasury Bill tokens.
type Syn12Metadata struct {
	TokenID          string    // Unique identifier for the token
	Name             string    // Full name of the T-Bill (e.g., "12 Month Treasury Bill")
	Symbol           string    // Symbol representing the token (e.g., "TBL12")
	TBillCode        string    // Treasury Bill code used internationally
	Issuer           IssuerInfo // Information about the issuer (government, central bank)
	MaturityDate     time.Time // Maturity date of the T-Bill
	DiscountRate     float64   // Discount rate at issuance
	CreationDate     time.Time // Timestamp of token creation
	TotalSupply      uint64    // Total supply of tokens
	CirculatingSupply uint64   // Number of tokens in circulation
	LegalCompliance  LegalInfo // Legal and compliance information
}

// Syn12Token represents a Treasury Bill token with metadata and value.
type Syn12Token struct {
	TokenID      string                      // Unique token ID
	Metadata     Syn12Metadata               // Associated metadata
	Issuer       string                      // Issuer of the token (Central Bank or Government)
	Ledger       *ledger.Ledger              // Ledger to track issuance and transfers
	Consensus    *consensus.SynnergyConsensus // Consensus engine for validation
	Compliance   *compliance.KYCAmlService   // KYC/AML Compliance for regulatory requirements
	CentralBank  string                      // Address of the Central Bank (for minting and burning)
	Encrypted    bool                        // Indicates if the token data is encrypted
}

// TokenFactory manages the issuance, burning, and management of SYN12 tokens.
type TokenFactory struct {
	mutex          sync.Mutex
	Ledger         *ledger.Ledger                // Ledger for tracking tokens
	Consensus      *consensus.SynnergyConsensus  // Consensus for validation
	Encryption     *encryption.EncryptionService // Encryption service
	Compliance     *compliance.KYCAmlService     // Compliance service
	CentralBank    string                        // Central Bank's unique address for minting and burning
}

// NewTokenFactory initializes the TokenFactory for SYN12 tokens.
func NewTokenFactory(ledgerInstance *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService, complianceService *compliance.KYCAmlService, centralBank string) *TokenFactory {
	return &TokenFactory{
		Ledger:      ledgerInstance,
		Consensus:   consensusEngine,
		Encryption:  encryptionService,
		Compliance:  complianceService,
		CentralBank: centralBank,
	}
}

// IssueToken handles the issuance of new SYN12 Treasury Bill tokens.
func (factory *TokenFactory) IssueToken(name, symbol, tbillCode, issuerID string, amount uint64, maturityDate time.Time, discountRate float64) (string, error) {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Compliance and KYC/AML Checks
	if err := factory.Compliance.VerifyIssuer(issuerID); err != nil {
		return "", fmt.Errorf("KYC/AML verification failed for issuer: %v", err)
	}

	// Consensus validation for token issuance
	if err := factory.Consensus.ValidateIssuer(issuerID); err != nil {
		return "", fmt.Errorf("consensus validation failed for issuer: %v", err)
	}

	if err := factory.Consensus.AuthorizeTokenIssuance(issuerID, amount); err != nil {
		return "", fmt.Errorf("token issuance not authorized: %v", err)
	}

	// Create the token and register it in the ledger
	tokenID := fmt.Sprintf("SYN12-%d-%s", time.Now().UnixNano(), tbillCode)
	token := Syn12Token{
		TokenID: tokenID,
		Metadata: Syn12Metadata{
			TokenID:          tokenID,
			Name:             name,
			Symbol:           symbol,
			TBillCode:        tbillCode,
			Issuer:           IssuerInfo{Name: issuerID},
			MaturityDate:     maturityDate,
			DiscountRate:     discountRate,
			CreationDate:     time.Now(),
			TotalSupply:      amount,
			CirculatingSupply: amount,
		},
		Issuer:      issuerID,
		Ledger:      factory.Ledger,
		Consensus:   factory.Consensus,
		Compliance:  factory.Compliance,
		CentralBank: factory.CentralBank,
	}

	// Encrypt token data
	encryptedTokenID, err := factory.Encryption.Encrypt([]byte(token.TokenID))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token ID: %v", err)
	}

	// Record the issuance in the ledger
	if err := factory.Ledger.RecordIssuance(string(encryptedTokenID), issuerID, amount); err != nil {
		return "", fmt.Errorf("failed to record issuance in ledger: %v", err)
	}

	log.Printf("Successfully issued SYN12 token with ID: %s, Name: %s, Symbol: %s", tokenID, name, symbol)
	return tokenID, nil
}

// BurnToken burns SYN12 tokens, reducing the circulating supply.
func (factory *TokenFactory) BurnToken(tokenID string, amount uint64, burnerID string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Only the central bank can burn tokens
	if burnerID != factory.CentralBank {
		return fmt.Errorf("only the central bank can burn tokens")
	}

	// Validate ownership before burning tokens
	if err := factory.Consensus.VerifyOwnership(tokenID, burnerID, amount); err != nil {
		return fmt.Errorf("ownership verification failed: %v", err)
	}

	// Authorize token burning through consensus
	if err := factory.Consensus.AuthorizeTokenBurning(tokenID, amount, burnerID); err != nil {
		return fmt.Errorf("token burning not authorized: %v", err)
	}

	// Burn the tokens and update the ledger
	if err := factory.Ledger.BurnTokens(tokenID, amount); err != nil {
		return fmt.Errorf("failed to burn tokens: %v", err)
	}

	log.Printf("Successfully burned %d tokens of ID: %s", amount, tokenID)
	return nil
}

// TransferOwnership handles token ownership transfer between entities.
func (factory *TokenFactory) TransferOwnership(tokenID, fromID, toID string, amount uint64) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Validate ownership before transferring tokens
	if err := factory.Consensus.VerifyOwnership(tokenID, fromID, amount); err != nil {
		return fmt.Errorf("ownership verification failed: %v", err)
	}

	// Transfer the tokens in the ledger
	if err := factory.Ledger.TransferTokens(tokenID, fromID, toID, amount); err != nil {
		return fmt.Errorf("failed to transfer tokens: %v", err)
	}

	log.Printf("Successfully transferred %d tokens of ID: %s from %s to %s", amount, tokenID, fromID, toID)
	return nil
}

// RedeemToken redeems SYN12 tokens by converting them to fiat upon maturity.
func (factory *TokenFactory) RedeemToken(tokenID, redeemerID string) error {
	factory.mutex.Lock()
	defer factory.mutex.Unlock()

	// Retrieve the token and check maturity
	token, err := factory.Ledger.GetToken(tokenID)
	if err != nil {
		return fmt.Errorf("failed to retrieve token: %v", err)
	}
	if token.Metadata.MaturityDate.After(time.Now()) {
		return fmt.Errorf("token %s is not yet matured for redemption", tokenID)
	}

	// Compliance and KYC/AML Checks
	if err := factory.Compliance.VerifyRedeemer(token.Metadata.Issuer, redeemerID); err != nil {
		return fmt.Errorf("redeemer verification failed: %v", err)
	}

	// Authorize the redemption through consensus
	if err := factory.Consensus.AuthorizeRedemption(tokenID, redeemerID); err != nil {
		return fmt.Errorf("token redemption not authorized: %v", err)
	}

	// Update the ledger to mark the token as redeemed
	if err := factory.Ledger.MarkAsRedeemed(tokenID, redeemerID); err != nil {
		return fmt.Errorf("failed to mark token as redeemed: %v", err)
	}

	log.Printf("Successfully redeemed SYN12 token with ID: %s", tokenID)
	return nil
}
