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

// ReviewLog represents a review that a user has received from another user, contributing to their reputation.
type ReviewLog struct {
	ReviewID     string    // Unique identifier for the review
	Reviewer     string    // ID of the user who provided the review
	Rating       float64   // Rating provided by the reviewer (1-5 stars, for example)
	Description  string    // Description or feedback provided by the reviewer
	ReviewDate   time.Time // Date when the review was made
}

// DisputeLog represents a dispute raised by the user regarding a reputation penalty or review.
type DisputeLog struct {
	DisputeID    string    // Unique identifier for the dispute
	Challenger   string    // ID of the user who raised the dispute
	Reason       string    // Reason for the dispute (e.g., unfair review, fraudulent penalty)
	DisputeDate  time.Time // Date when the dispute was filed
	Resolution   string    // Outcome of the dispute (e.g., "Resolved", "Rejected")
	ResolutionDate time.Time // Date when the dispute was resolved
	Arbitrator   string    // ID of the arbitrator or system handling the dispute
}

// ReputationEventLog captures key events that affect the user's reputation.
type ReputationEventLog struct {
	EventID       string    // Unique identifier for the reputation event
	EventType     string    // Type of event (e.g., "Endorsement", "Penalty", "Score Update")
	Description   string    // Description of the event
	EventDate     time.Time // Date when the event occurred
	PerformedBy   string    // ID of the user or system that triggered the event
	ImpactOnScore float64   // The impact (positive or negative) on the user's reputation score
}

// HistoricalReputationData tracks the evolution of a user's reputation over time.
type HistoricalReputationData struct {
	Timestamp     time.Time // When the reputation was recorded
	ReputationScore float64 // Reputation score at that time
	TrustLevel    string    // Trust level at that time
	ActivityLog   string    // Description of the activity or event that led to the change
}

// ReputationAnalytics offers insights and analytics on the user's reputation performance.
type ReputationAnalytics struct {
	PositiveEndorsements   int     // Total number of positive endorsements received
	NegativeEndorsements   int     // Total number of negative endorsements received
	AverageReviewRating    float64 // Average rating from reviews
	DisputeSuccessRate     float64 // Percentage of disputes that were resolved in favor of the user
	EngagementLevel        string  // User engagement level (e.g., "Low", "Moderate", "High")
	ActivityScore          float64 // Score representing user's activity in the ecosystem
}

// ImmutableRecord stores immutable records for compliance and transparency.
type ImmutableRecord struct {
	RecordID    string    // Unique identifier for the record
	Description string    // Description of the record
	Timestamp   time.Time // When the record was created
}

// TokenFactory is responsible for creating SYN1500 reputation tokens
type TokenFactory struct {
	FactoryID string // Unique ID for the factory managing SYN1500 token creation
}

// CreateSYN1500Token creates a new SYN1500Token with provided parameters
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

	// Return the created token
	return token, nil
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

// AddEndorsement adds an endorsement to the user's reputation.
func (token *SYN1500Token) AddEndorsement(from string, rating float64, review string) {
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
