package syn130

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Syn130TokenFactory handles the issuance, management, and transfer of Syn130 tokens.
type Syn130TokenFactory struct {
	ledgerManager     *ledger.LedgerManager         // Ledger to track issuance and transfers
	encryptionService *encryption.EncryptionService // Encryption service for secure data
	consensus         *consensus.SynnergyConsensus  // Consensus engine for validation
	mutex             sync.Mutex                    // Mutex for thread-safe operations
}

// Syn130Token represents a token with comprehensive attributes.
type SYN130Token struct {
	ID                    string
	Name                  string
	Owner                 string
	Value                 float64
	Metadata              SYN130Metadata
	SaleHistory           []SaleRecord
	LeaseTerms            []LeaseTerms
	LicenseTerms          []LicenseTerms
	RentalTerms           []RentalTerms
	CoOwnershipAgreements []CoOwnershipAgreement
	AssetType             string
	Classification        string
	CreationDate          time.Time
	LastUpdated           time.Time
	TransactionHistory    []TransactionRecord
	Provenance            []ProvenanceRecord
	IsEncrypted           bool
}

// NewSyn130TokenFactory initializes a new Syn130TokenFactory.
func NewSyn130TokenFactory(ledgerManager *ledger.LedgerManager, encryptionService *encryption.EncryptionService, consensusEngine *consensus.SynnergyConsensus) *Syn130TokenFactory {
	return &Syn130TokenFactory{
		ledgerManager:     ledgerManager,
		encryptionService: encryptionService,
		consensus:         consensusEngine,
	}
}

// IssueToken creates and issues a new Syn130 token.
func (tf *Syn130TokenFactory) IssueToken(name, owner string, value float64, metadata map[string]string) (*Syn130Token, error) {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Validate inputs
	if name == "" || owner == "" || value <= 0 {
		return nil, errors.New("invalid token parameters")
	}

	// Generate Token ID
	tokenID := fmt.Sprintf("SYN130-%d", time.Now().UnixNano())

	// Create token structure
	token := &Syn130Token{
		ID:           tokenID,
		Name:         name,
		Owner:        owner,
		Value:        value,
		Metadata:     metadata,
		CreationDate: time.Now(),
		LastUpdated:  time.Now(),
	}

	// Encrypt the token data
	encryptedToken, err := tf.encryptToken(token)
	if err != nil {
		return nil, fmt.Errorf("token encryption failed: %v", err)
	}
	token.IsEncrypted = true

	// Record the token in the ledger
	if err := tf.ledgerManager.RecordIssuance(token.ID, token.Owner, value); err != nil {
		return nil, fmt.Errorf("ledger recording failed: %v", err)
	}

	// Log the issuance event
	if err := tf.ledgerManager.LogEvent(token.ID, "Token Issued"); err != nil {
		return nil, fmt.Errorf("failed to log issuance event: %v", err)
	}

	fmt.Printf("Successfully issued Syn130 token with ID: %s, Owner: %s\n", tokenID, owner)
	return token, nil
}

// TransferToken transfers ownership of a Syn130 token.
func (tf *Syn130TokenFactory) TransferToken(tokenID, fromOwner, toOwner string, amount float64) error {
	tf.mutex.Lock()
	defer tf.mutex.Unlock()

	// Validate inputs
	if tokenID == "" || fromOwner == "" || toOwner == "" || amount <= 0 {
		return errors.New("invalid transfer parameters")
	}

	// Validate the transaction through the consensus
	if err := tf.consensus.ValidateTransfer(fromOwner, toOwner, amount); err != nil {
		return fmt.Errorf("transfer validation failed: %v", err)
	}

	// Record the transfer in the ledger
	if err := tf.ledgerManager.RecordTransfer(tokenID, fromOwner, toOwner, amount); err != nil {
		return fmt.Errorf("ledger recording failed: %v", err)
	}

	// Log the transfer event
	if err := tf.ledgerManager.LogEvent(tokenID, "Token Transferred"); err != nil {
		return fmt.Errorf("failed to log transfer event: %v", err)
	}

	fmt.Printf("Successfully transferred Syn130 token ID: %s from %s to %s\n", tokenID, fromOwner, toOwner)
	return nil
}

// Encrypt token data for secure storage.
func (tf *Syn130TokenFactory) encryptToken(token *Syn130Token) ([]byte, error) {
	tokenData := fmt.Sprintf("%s:%s:%f", token.ID, token.Owner, token.Value)
	encryptedData, err := tf.encryptionService.Encrypt([]byte(tokenData))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token data: %v", err)
	}
	return encryptedData, nil
}

