package syn3900

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"sync"

)

// BenefitManager manages SYN3900 benefit tokens, their allocation, claiming, tracking, and expiration.
type BenefitManager struct {
	ledgerService     *ledger.LedgerService
	encryptionService *encryption.Encryptor
	consensusService  *consensus.SynnergyConsensus
	mutex             sync.Mutex
	benefitRequirements map[string]BenefitRequirement // Map tokenID to BenefitRequirement struct
}

// BenefitRequirement captures the specific requirements needed for a recipient to claim a benefit.
type BenefitRequirement struct {
	IncomeCap       float64  // Maximum earnings allowed to claim the benefit
	MinEarnings     float64  // Minimum earnings required to claim the benefit
	EligibilityDocs []string // Required documents (e.g., proof of income, ID)
	AgeLimit        int      // Minimum age required to claim
}

// NewBenefitManager creates a new BenefitManager instance.
func NewBenefitManager(ledgerService *ledger.LedgerService, encryptionService *encryption.Encryptor, consensusService *consensus.SynnergyConsensus) *BenefitManager {
	return &BenefitManager{
		ledgerService:       ledgerService,
		encryptionService:   encryptionService,
		consensusService:    consensusService,
		benefitRequirements: make(map[string]BenefitRequirement),
	}
}

// SetupBenefitRequirements configures the eligibility requirements for a benefit token.
func (bm *BenefitManager) SetupBenefitRequirements(tokenID string, requirement BenefitRequirement) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Store benefit requirements for the tokenID
	bm.benefitRequirements[tokenID] = requirement

	// Log the setup event in the ledger
	bm.logEvent(tokenID, "BenefitRequirementsSet")
	return nil
}

// AllocateBenefit allocates benefits to a recipient based on predefined eligibility criteria, including earnings check.
func (bm *BenefitManager) AllocateBenefit(tokenID string, recipient string, amount float64, earnings float64, recipientAge int, documents []string) (*BenefitAllocation, error) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the token and eligibility requirements from the ledger
	token, err := bm.RetrieveToken(tokenID)
	if err != nil {
		return nil, err
	}

	// Check if the recipient meets the benefit requirements
	requirements := bm.benefitRequirements[tokenID]
	if !bm.checkEligibility(token, recipient, earnings, recipientAge, documents, requirements) {
		return nil, errors.New("recipient does not meet eligibility criteria")
	}

	// Create benefit allocation record
	allocation := BenefitAllocation{
		AllocationID:    generateUniqueAllocationID(),
		AllocationDate:  time.Now(),
		Recipient:       recipient,
		AmountAllocated: amount,
		UsageBreakdown:  "Initial Allocation",
	}

	// Update token allocation history and balance
	token.AllocationHistory = append(token.AllocationHistory, allocation)
	token.Metadata.Amount -= amount

	// Encrypt and store the updated token
	encryptedToken, err := bm.encryptionService.EncryptData(token)
	if err != nil {
		return nil, err
	}
	if err := bm.ledgerService.StoreData(tokenID, encryptedToken); err != nil {
		return nil, err
	}

	// Validate the allocation using Synnergy Consensus
	if err := bm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return nil, err
	}

	// Log the allocation event in the ledger
	bm.logEvent(tokenID, "BenefitAllocated")

	return &allocation, nil
}

// ClaimBenefit allows a recipient to claim a benefit if the specified conditions, such as income level, are met.
func (bm *BenefitManager) ClaimBenefit(tokenID string, recipient string, claimAmount float64, earnings float64, recipientAge int, documents []string) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := bm.RetrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Check if the recipient meets the eligibility criteria for claiming the benefit
	requirements := bm.benefitRequirements[tokenID]
	if !bm.checkClaimConditions(token, recipient, earnings, recipientAge, documents, requirements) {
		return errors.New("conditions for claiming the benefit have not been met")
	}

	// Create a transaction record for the claim
	transaction := BenefitTransaction{
		TransactionID:   generateUniqueTransactionID(),
		Timestamp:       time.Now(),
		TransactionType: "Claim",
		Amount:          claimAmount,
		Recipient:       recipient,
	}

	// Update the token's transaction history
	token.TransactionHistory = append(token.TransactionHistory, transaction)

	// Encrypt and store the updated token
	encryptedToken, err := bm.encryptionService.EncryptData(token)
	if err != nil {
		return err
	}
	if err := bm.ledgerService.StoreData(tokenID, encryptedToken); err != nil {
		return err
	}

	// Validate the claim using Synnergy Consensus
	if err := bm.consensusService.ValidateSubBlock(tokenID); err != nil {
		return err
	}

	// Log the claim event in the ledger
	bm.logEvent(tokenID, "BenefitClaimed")

	return nil
}

// TrackBenefitUsage provides real-time tracking of the benefit usage and remaining balance.
func (bm *BenefitManager) TrackBenefitUsage(tokenID string) (float64, error) {
	// Retrieve the token from the ledger
	token, err := bm.RetrieveToken(tokenID)
	if err != nil {
		return 0, err
	}

	// Return the remaining balance of the token
	return token.Metadata.Amount, nil
}

// checkEligibility checks if the recipient meets the predefined eligibility criteria for benefit allocation.
func (bm *BenefitManager) checkEligibility(token *Syn3900Token, recipient string, earnings float64, recipientAge int, documents []string, requirements BenefitRequirement) bool {
	// Check earnings criteria
	if earnings > requirements.IncomeCap || earnings < requirements.MinEarnings {
		return false
	}

	// Check age limit
	if recipientAge < requirements.AgeLimit {
		return false
	}

	// Check for required documents
	for _, doc := range requirements.EligibilityDocs {
		if !bm.documentExists(doc, documents) {
			return false
		}
	}

	return true
}

// checkClaimConditions checks if the recipient meets the eligibility criteria for claiming the benefit.
func (bm *BenefitManager) checkClaimConditions(token *Syn3900Token, recipient string, earnings float64, recipientAge int, documents []string, requirements BenefitRequirement) bool {
	return bm.checkEligibility(token, recipient, earnings, recipientAge, documents, requirements)
}

// documentExists checks if the required document exists in the provided list of documents.
func (bm *BenefitManager) documentExists(doc string, providedDocs []string) bool {
	for _, providedDoc := range providedDocs {
		if providedDoc == doc {
			return true
		}
	}
	return false
}

// ExpireBenefit automatically handles the expiration of benefit tokens once their validity period ends.
func (bm *BenefitManager) ExpireBenefit(tokenID string) error {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	// Retrieve the token from the ledger
	token, err := bm.RetrieveToken(tokenID)
	if err != nil {
		return err
	}

	// Check if the benefit has expired
	if token.Metadata.ValidUntil != nil && time.Now().After(*token.Metadata.ValidUntil) {
		// Mark the benefit token as expired
		token.Metadata.Status = "Expired"

		// Encrypt and store the updated token
		encryptedToken, err := bm.encryptionService.EncryptData(token)
		if err != nil {
			return err
		}
		if err := bm.ledgerService.StoreData(tokenID, encryptedToken); err != nil {
			return err
		}

		// Log the expiration event in the ledger
		bm.logEvent(tokenID, "BenefitExpired")

		// Validate the expiration using Synnergy Consensus
		if err := bm.consensusService.ValidateSubBlock(tokenID); err != nil {
			return err
		}
	}

	return nil
}

// RetrieveToken retrieves the SYN3900 token from the ledger and decrypts it.
func (bm *BenefitManager) RetrieveToken(tokenID string) (*Syn3900Token, error) {
	encryptedData, err := bm.ledgerService.RetrieveData(tokenID)
	if err != nil {
		return nil, err
	}

	decryptedToken, err := bm.encryptionService.DecryptData(encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedToken.(*Syn3900Token), nil
}

// generateUniqueAllocationID generates a unique identifier for a benefit allocation.
func generateUniqueAllocationID() string {
	return "allocation-" + time.Now().Format("20060102150405")
}

// generateUniqueTransactionID generates a unique identifier for a benefit transaction.
func generateUniqueTransactionID() string {
	return "transaction-" + time.Now().Format("20060102150405")
}

// logEvent logs an event related to a benefit token in the ledger.
func (bm *BenefitManager) logEvent(tokenID, eventType string) {
	_ = bm.ledgerService.LogEvent(eventType, time.Now(), tokenID)
}
