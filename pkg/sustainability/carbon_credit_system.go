package sustainability

import (
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn4300"
)

// NewCarbonCreditSystem initializes the carbon credit system
func NewCarbonCreditSystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *CarbonCreditSystem {
	return &CarbonCreditSystem{
		Tokens:           make(map[string]*syn4300.SYN4300Token),
		Ledger:           ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// IssueCarbonCredits issues new carbon credits to a specified owner
func (ccs *CarbonCreditSystem) IssueCarbonCredits(creditID, tokenID, issuer, owner string, amount float64) (*syn4300.SYN4300Token, error) {
	ccs.mu.Lock()
	defer ccs.mu.Unlock()

	// Encrypt token data
	tokenData := fmt.Sprintf("CreditID: %s, TokenID: %s, Issuer: %s, Owner: %s, Amount: %f", creditID, tokenID, issuer, owner, amount)
	encryptedData, err := ccs.EncryptionService.EncryptData([]byte(tokenData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt carbon credit data: %v", err)
	}

	// Create a new Syn700Token (carbon credit)
	token := &syn4300.SYN4300Token{
		CreditID:     creditID,
		TokenID:      tokenID,
		Issuer:       issuer,
		Owner:        owner,
		IssuedAmount: amount,
		IssuedTime:   time.Now(),
		IsRetired:    false,
	}

	// Add the token to the system
	ccs.Tokens[creditID] = token

	// Log the issuance of carbon credits in the ledger
	err = ccs.Ledger.RecordCarbonCreditIssuance(creditID, tokenID, issuer, owner, amount, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log carbon credit issuance: %v", err)
	}

	fmt.Printf("Carbon credit %s issued by %s to %s for %f tons of CO2\n", creditID, issuer, owner, amount)
	return token, nil
}

// TransferCarbonCredits transfers ownership of a carbon credit to a new owner
func (ccs *CarbonCreditSystem) TransferCarbonCredits(creditID, newOwner string) error {
	ccs.mu.Lock()
	defer ccs.mu.Unlock()

	// Retrieve the carbon credit
	token, exists := ccs.Tokens[creditID]
	if !exists {
		return fmt.Errorf("carbon credit %s not found", creditID)
	}

	// Ensure the token has not been retired
	if token.IsRetired {
		return fmt.Errorf("carbon credit %s has already been retired", creditID)
	}

	// Transfer ownership
	oldOwner := token.Owner
	token.Owner = newOwner

	// Log the transfer in the ledger
	err := ccs.Ledger.RecordCarbonCreditTransfer(creditID, oldOwner, newOwner, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log carbon credit transfer: %v", err)
	}

	fmt.Printf("Carbon credit %s transferred from %s to %s\n", creditID, oldOwner, newOwner)
	return nil
}

// RetireCarbonCredits retires a carbon credit, marking it as used
func (ccs *CarbonCreditSystem) RetireCarbonCredits(creditID string) error {
	ccs.mu.Lock()
	defer ccs.mu.Unlock()

	// Retrieve the carbon credit
	token, exists := ccs.Tokens[creditID]
	if !exists {
		return fmt.Errorf("carbon credit %s not found", creditID)
	}

	// Ensure the token has not already been retired
	if token.IsRetired {
		return fmt.Errorf("carbon credit %s has already been retired", creditID)
	}

	// Mark the token as retired
	token.IsRetired = true
	token.RetiredTime = time.Now()

	// Log the retirement in the ledger
	err := ccs.Ledger.RecordCarbonCreditRetirement(creditID, token.Owner, token.IssuedAmount, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log carbon credit retirement: %v", err)
	}

	fmt.Printf("Carbon credit %s retired by %s for %f tons of CO2\n", creditID, token.Owner, token.IssuedAmount)
	return nil
}

// ViewCarbonCredits allows viewing of the details of a specific carbon credit
func (ccs *CarbonCreditSystem) ViewCarbonCredits(creditID string) (*syn4300.SYN4300Token, error) {
	ccs.mu.Lock()
	defer ccs.mu.Unlock()

	// Retrieve the carbon credit
	token, exists := ccs.Tokens[creditID]
	if !exists {
		return nil, fmt.Errorf("carbon credit %s not found", creditID)
	}

	return token, nil
}

