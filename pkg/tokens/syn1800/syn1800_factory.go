package syn1800

import (
	"time"
	"fmt"
)

// SYN1800Token represents a carbon footprint token under the SYN1800 standard.
type SYN1800Token struct {
	TokenID               string            // Unique identifier for the token
	Owner                 string            // The owner of the carbon token (individual or organization)
	CarbonAmount          float64           // Amount of carbon (negative for emissions, positive for offsets)
	IssueDate             time.Time         // Date the token was issued
	Description           string            // Description of the carbon activity (e.g., emission or offset source)
	Source                string            // Source of emission or offset verification (e.g., "Reforestation", "Energy Efficiency")
	VerificationStatus    string            // Status of the verification (e.g., "Verified", "Pending", "Rejected")
	OffsetProjects        []OffsetProject   // Projects related to carbon offsetting (if applicable)
	RewardRecords         []RewardLog       // Reward records for incentivized activities
	CarbonFootprintLogs   []EmissionOffsetLog // Logs of all emission and offset activities
	ImmutableRecords      []ImmutableRecord // Immutable records for compliance and traceability
	NetBalance            float64           // Net carbon balance (sum of all emissions and offsets)
	SourceVerificationLog []VerificationLog // Log of the source verification for emission or offset
	RestrictedTransfers   bool              // Indicates if transfers of this token are restricted
	ApprovalRequired      bool              // Indicates if certain actions require third-party approval (e.g., high-value offsets)
	EncryptedMetadata     []byte            // Encrypted metadata related to sensitive information (e.g., contracts, offset details)
}

// OffsetProject represents a carbon offset project related to the token.
type OffsetProject struct {
	ProjectID     string    // Unique identifier for the project
	ProjectName   string    // Name of the offset project (e.g., "Reforestation in Amazon")
	OffsetAmount  float64   // Amount of carbon offset by this project
	Verification  string    // Verification status of the project (e.g., "Verified by Third-Party")
	StartDate     time.Time // Project start date
	EndDate       time.Time // Project end date
	Description   string    // Detailed description of the project
	OffsetType    string    // Type of offset (e.g., "Renewable Energy", "Carbon Capture")
}

// RewardLog represents rewards issued for carbon reduction activities.
type RewardLog struct {
	RewardID      string    // Unique identifier for the reward
	RecipientID   string    // ID of the recipient who earned the reward
	RewardAmount  float64   // Amount of the reward (in tokens or other units)
	Activity      string    // The activity for which the reward was issued (e.g., "Tree Planting", "Energy Efficiency")
	RewardDate    time.Time // Date the reward was issued
	Verification  string    // Verification status of the rewarded activity
}

// EmissionOffsetLog represents a record of emission or offset activity.
type EmissionOffsetLog struct {
	LogID         string    // Unique identifier for the log entry
	ActivityType  string    // Type of activity ("Emission" or "Offset")
	Amount        float64   // Amount of carbon emitted or offset
	ActivityDate  time.Time // Date the activity occurred
	Description   string    // Description of the activity
	VerifiedBy    string    // Name of the entity or individual who verified the activity
}

// VerificationLog captures source verification details for emissions or offsets.
type VerificationLog struct {
	VerificationID string    // Unique identifier for the verification log
	Source         string    // Source of verification (e.g., "Third-party verifier")
	VerificationDate time.Time // Date of the verification
	Description    string    // Description of the verification process or entity
	VerifiedAmount float64   // Amount verified as emission or offset
	Status         string    // Verification status (e.g., "Verified", "Pending")
}

// ImmutableRecord stores immutable records for compliance and transparency.
type ImmutableRecord struct {
	RecordID    string    // Unique identifier for the immutable record
	Description string    // Description of the event or record
	Timestamp   time.Time // Time the record was created
}

// GetNetBalance calculates the net balance of carbon (positive for net offsets, negative for net emissions).
func (token *SYN1800Token) GetNetBalance() float64 {
	var netBalance float64
	for _, log := range token.CarbonFootprintLogs {
		netBalance += log.Amount
	}
	return netBalance
}

// AddEmissionLog adds a new emission log to the token's CarbonFootprintLogs.
func (token *SYN1800Token) AddEmissionLog(amount float64, description string, verifiedBy string) {
	newLog := EmissionOffsetLog{
		LogID:        generateUniqueID(),
		ActivityType: "Emission",
		Amount:       amount,
		ActivityDate: time.Now(),
		Description:  description,
		VerifiedBy:   verifiedBy,
	}
	token.CarbonFootprintLogs = append(token.CarbonFootprintLogs, newLog)
	token.NetBalance -= amount
}

// AddOffsetLog adds a new offset log to the token's CarbonFootprintLogs.
func (token *SYN1800Token) AddOffsetLog(amount float64, description string, verifiedBy string) {
	newLog := EmissionOffsetLog{
		LogID:        generateUniqueID(),
		ActivityType: "Offset",
		Amount:       amount,
		ActivityDate: time.Now(),
		Description:  description,
		VerifiedBy:   verifiedBy,
	}
	token.CarbonFootprintLogs = append(token.CarbonFootprintLogs, newLog)
	token.NetBalance += amount
}

// AddRewardRecord adds a new reward log entry for incentivized carbon reduction activities.
func (token *SYN1800Token) AddRewardRecord(recipientID string, amount float64, activity string, verifiedBy string) {
	newReward := RewardLog{
		RewardID:     generateUniqueID(),
		RecipientID:  recipientID,
		RewardAmount: amount,
		Activity:     activity,
		RewardDate:   time.Now(),
		Verification: verifiedBy,
	}
	token.RewardRecords = append(token.RewardRecords, newReward)
}

// TokenFactory handles the creation and management of SYN1800 tokens (Carbon Footprint Tokens)
type TokenFactory struct {
	ledger *ledger.Ledger
}

// NewTokenFactory initializes a new TokenFactory instance.
func NewTokenFactory(ledger *ledger.Ledger) *TokenFactory {
	return &TokenFactory{ledger: ledger}
}

// CreateSYN1800Token creates a new SYN1800Token and integrates it into the ledger.
func (factory *TokenFactory) CreateSYN1800Token(owner string, carbonAmount float64, description string, source string, verificationStatus string) (*common.SYN1800Token, error) {
	// Create a new token with initial parameters.
	token := &SYN1800Token{
		TokenID:            generateUniqueID(),
		Owner:              owner,
		CarbonAmount:       carbonAmount,
		IssueDate:          time.Now(),
		Description:        description,
		Source:             source,
		VerificationStatus: verificationStatus,
		NetBalance:         carbonAmount, // Initial balance matches the initial carbon amount
		ImmutableRecords:   []common.ImmutableRecord{},
		CarbonFootprintLogs: []common.EmissionOffsetLog{},
	}

	// Add an initial immutable record for the token creation
	token.ImmutableRecords = append(token.ImmutableRecords, common.ImmutableRecord{
		RecordID:    generateUniqueID(),
		Description: "Token Creation",
		Timestamp:   time.Now(),
	})

	// Integrate with the ledger
	err := factory.ledger.AddTokenToLedger(token)
	if err != nil {
		return nil, fmt.Errorf("failed to add token to ledger: %v", err)
	}

	// Encrypt sensitive metadata before storing it
	encryptedMetadata, err := encryptMetadata(token)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt metadata: %v", err)
	}
	token.EncryptedMetadata = encryptedMetadata

	return token, nil
}

// AddEmissionLog adds an emission log to the token and updates the ledger.
func (factory *TokenFactory) AddEmissionLog(tokenID string, amount float64, description string, verifiedBy string) error {
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	syn1800Token.AddEmissionLog(amount, description, verifiedBy)

	// Update the ledger with the new log and balance
	err = factory.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger: %v", err)
	}

	return nil
}

// AddOffsetLog adds an offset log to the token and updates the ledger.
func (factory *TokenFactory) AddOffsetLog(tokenID string, amount float64, description string, verifiedBy string) error {
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	syn1800Token.AddOffsetLog(amount, description, verifiedBy)

	// Update the ledger with the new log and balance
	err = factory.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger: %v", err)
	}

	return nil
}

// AddRewardLog adds a reward log for incentivized carbon reduction activities and updates the ledger.
func (factory *TokenFactory) AddRewardLog(tokenID string, recipientID string, amount float64, activity string, verifiedBy string) error {
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return fmt.Errorf("token not found: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return fmt.Errorf("invalid token type")
	}

	syn1800Token.AddRewardRecord(recipientID, amount, activity, verifiedBy)

	// Update the ledger with the new reward log
	err = factory.ledger.UpdateTokenInLedger(syn1800Token)
	if err != nil {
		return fmt.Errorf("failed to update ledger: %v", err)
	}

	return nil
}

// GetNetBalance fetches the net carbon balance for a specific token from the ledger.
func (factory *TokenFactory) GetNetBalance(tokenID string) (float64, error) {
	token, err := factory.ledger.GetTokenByID(tokenID)
	if err != nil {
		return 0, fmt.Errorf("token not found: %v", err)
	}

	// Cast to SYN1800Token type
	syn1800Token, ok := token.(*common.SYN1800Token)
	if !ok {
		return 0, fmt.Errorf("invalid token type")
	}

	return syn1800Token.GetNetBalance(), nil
}

// encryptMetadata handles the encryption of sensitive token metadata.
func encryptMetadata(token *common.SYN1800Token) ([]byte, error) {
	// Placeholder encryption logic. Replace with your real encryption implementation.
	return crypto.Encrypt([]byte(fmt.Sprintf("%v", token)), "encryption-key")
}

// generateUniqueID generates a unique ID for each entity (token, log, etc.).
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
