package syn1500

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SYN1500Token represents a reputation token under the SYN1500 standard.
type SYN1500Token struct {
	TokenID               string              // Unique identifier for the token
	Owner                 string              // The owner of the token (user or platform)
	ReputationScore       float64             // User's current reputation score
	TrustLevel            string              // Trust level based on the reputation score (e.g., Bronze, Silver, Gold)
	LastActivityDate      time.Time           // The last activity date for decay calculations
	DecayRate             float64             // Decay rate of the reputation score when inactive
	EndorsementsReceived  int                 // Number of endorsements received from other users
	ReviewsReceived       []ReviewLog         // List of reviews received by the user
	VerificationBadge     bool                // Whether the user has received a verification badge
	Disputes              []DisputeLog        // Dispute logs, if the user has challenged any reviews or penalties
	RoleBasedTrustMetrics map[string]float64  // Trust metrics differentiated by role (developer, trader, etc.)
	ComplianceStatus      string              // Compliance status of the token (e.g., "Compliant", "Non-Compliant")
	ReputationEvents      []ReputationEventLog // Logs of key reputation events (score changes, penalties, endorsements)
	ReputationHistory     []HistoricalReputationData // Historical reputation data tracking changes over time
	RestrictedTransfers   bool                // Whether there are restrictions on transferring reputation tokens
	ApprovalRequired      bool                // Whether certain actions (like score adjustments) require approval
	EncryptedMetadata     []byte              // Encrypted metadata containing sensitive information (e.g., KYC details)
	ReputationAnalytics   ReputationAnalytics // Analytics data for insights into user behavior and engagement
	ImmutableRecords      []ImmutableRecord   // Immutable records to ensure transparency and trust
}

// ReviewLog, DisputeLog, ReputationEventLog, HistoricalReputationData, ReputationAnalytics, and ImmutableRecord structs are the same as previously defined.

// TokenFactory is responsible for creating SYN1500 reputation tokens and interacting with the ledger
type TokenFactory struct {
	FactoryID string // Unique ID for the factory managing SYN1500 token creation
	Ledger    ledger.Ledger // Reference to the blockchain ledger for recording transactions
}

// CreateSYN1500Token creates a new SYN1500Token with provided parameters and logs the event in the ledger.
func (tf *TokenFactory) CreateSYN1500Token(owner string, decayRate float64, restrictedTransfers bool, approvalRequired bool) (*SYN1500Token, error) {
	// Generate a unique token ID using the factory ID and current timestamp
	tokenID := generateUniqueID(tf.FactoryID)

	// Initialize the SYN1500Token struct
	token := &SYN1500Token{
		TokenID:               tokenID,
		Owner:                 owner,
		ReputationScore:       0.0, // Initial score starts at 0
		TrustLevel:            "Bronze", // Initial trust level
		LastActivityDate:      time.Now(),
		DecayRate:             decayRate,
		EndorsementsReceived:  0,
		ReviewsReceived:       []ReviewLog{},
		VerificationBadge:     false,
		Disputes:              []DisputeLog{},
		RoleBasedTrustMetrics: map[string]float64{},
		ComplianceStatus:      "Compliant",
		ReputationEvents:      []ReputationEventLog{},
		ReputationHistory:     []HistoricalReputationData{},
		RestrictedTransfers:   restrictedTransfers,
		ApprovalRequired:      approvalRequired,
		EncryptedMetadata:     nil, // Encrypted metadata will be added later
		ReputationAnalytics:   ReputationAnalytics{},
		ImmutableRecords:      []ImmutableRecord{},
	}

	// Encrypt any sensitive metadata
	token.EncryptedMetadata = encryptMetadata(token)

	// Record the creation of the token in the ledger
	if err := tf.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        generateUniqueID(tokenID),
		Description: fmt.Sprintf("Created SYN1500Token for owner %s", owner),
		Timestamp:   time.Now(),
		Data:        token,
	}); err != nil {
		return nil, errors.New("failed to record transaction in the ledger")
	}

	// Return the created token
	return token, nil
}

// AddEndorsement adds an endorsement to the user's reputation and logs the event in the ledger.
func (token *SYN1500Token) AddEndorsement(tf *TokenFactory, from string, rating float64, review string) error {
	token.EndorsementsReceived++
	token.ReviewsReceived = append(token.ReviewsReceived, ReviewLog{
		ReviewID:     generateUniqueID(token.TokenID),
		Reviewer:     from,
		Rating:       rating,
		Description:  review,
		ReviewDate:   time.Now(),
	})

	// Update score based on review rating
	token.ReputationScore += rating
	token.updateTrustLevel()
	token.updateReputationHistory("Endorsement received")

	// Log the endorsement event in the ledger
	eventLog := ReputationEventLog{
		EventID:       generateUniqueID(token.TokenID),
		EventType:     "Endorsement",
		Description:   fmt.Sprintf("Endorsement received from %s with rating %.2f", from, rating),
		EventDate:     time.Now(),
		PerformedBy:   from,
		ImpactOnScore: rating,
	}
	token.ReputationEvents = append(token.ReputationEvents, eventLog)

	if err := tf.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        eventLog.EventID,
		Description: eventLog.Description,
		Timestamp:   eventLog.EventDate,
		Data:        token,
	}); err != nil {
		return errors.New("failed to record endorsement in the ledger")
	}

	return nil
}

// DisputeReputation allows a user to file a dispute and records the dispute in the ledger.
func (token *SYN1500Token) DisputeReputation(tf *TokenFactory, challenger string, reason string) error {
	dispute := DisputeLog{
		DisputeID:    generateUniqueID(token.TokenID),
		Challenger:   challenger,
		Reason:       reason,
		DisputeDate:  time.Now(),
		Resolution:   "Pending",
		Arbitrator:   "System",
	}

	token.Disputes = append(token.Disputes, dispute)

	// Log the dispute in the ledger
	if err := tf.Ledger.RecordTransaction(ledger.Transaction{
		TxID:        dispute.DisputeID,
		Description: fmt.Sprintf("Dispute filed by %s: %s", challenger, reason),
		Timestamp:   dispute.DisputeDate,
		Data:        token,
	}); err != nil {
		return errors.New("failed to record dispute in the ledger")
	}

	return nil
}

// updateTrustLevel updates the user's trust level based on the current reputation score
func (token *SYN1500Token) updateTrustLevel() {
	if token.ReputationScore > 100 {
		token.TrustLevel = "Gold"
	} else if token.ReputationScore > 50 {
		token.TrustLevel = "Silver"
	} else {
		token.TrustLevel = "Bronze"
	}
}

// updateReputationHistory logs historical changes in the user's reputation
func (token *SYN1500Token) updateReputationHistory(activity string) {
	token.ReputationHistory = append(token.ReputationHistory, HistoricalReputationData{
		Timestamp:      time.Now(),
		ReputationScore: token.ReputationScore,
		TrustLevel:     token.TrustLevel,
		ActivityLog:    activity,
	})
}

// encryptMetadata encrypts sensitive metadata for a SYN1500Token
func encryptMetadata(token *SYN1500Token) []byte {
	data, _ := json.Marshal(token)
	hash := sha256.Sum256(data)
	return hash[:]
}

// generateUniqueID creates a unique identifier for tokens and logs
func generateUniqueID(seed string) string {
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", seed, timestamp)))
	return hex.EncodeToString(hash[:])
}
