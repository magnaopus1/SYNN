package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	KYCCheckInterval          = 12 * time.Hour // Interval for checking KYC verification statuses
	MaxKYCVerificationRetries = 3              // Maximum number of retries for KYC identity verification
	KYCExpiryTime             = 365 * 24 * time.Hour // KYC identity verification validity period (1 year)
)

// KYCIdentityVerificationRestrictionAutomation monitors and restricts KYC identity verification processes across the network
type KYCIdentityVerificationRestrictionAutomation struct {
	consensusSystem          *consensus.SynnergyConsensus
	ledgerInstance           *ledger.Ledger
	stateMutex               *sync.RWMutex
	userKYCVerificationAttempts map[string]int // Tracks KYC verification attempts per user
	userKYCExpiryDates       map[string]time.Time // Tracks KYC expiration times per user
}

// NewKYCIdentityVerificationRestrictionAutomation initializes and returns an instance of KYCIdentityVerificationRestrictionAutomation
func NewKYCIdentityVerificationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *KYCIdentityVerificationRestrictionAutomation {
	return &KYCIdentityVerificationRestrictionAutomation{
		consensusSystem:          consensusSystem,
		ledgerInstance:           ledgerInstance,
		stateMutex:               stateMutex,
		userKYCVerificationAttempts: make(map[string]int),
		userKYCExpiryDates:       make(map[string]time.Time),
	}
}

// StartKYCMonitoring starts continuous monitoring of KYC identity verification processes
func (automation *KYCIdentityVerificationRestrictionAutomation) StartKYCMonitoring() {
	ticker := time.NewTicker(KYCCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorKYCVerifications()
		}
	}()
}

// monitorKYCVerifications checks recent KYC identity verification attempts and enforces restrictions
func (automation *KYCIdentityVerificationRestrictionAutomation) monitorKYCVerifications() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent KYC verification statuses from Synnergy Consensus
	recentKYCVerifications := automation.consensusSystem.GetRecentKYCVerifications()

	for _, verification := range recentKYCVerifications {
		// Check if the KYC has expired
		if automation.hasKYCExpired(verification.UserID) {
			automation.flagKYCViolation(verification, "KYC identity verification expired")
		} else if !automation.validateKYCVerificationAttempts(verification) {
			automation.flagKYCViolation(verification, "Exceeded maximum number of KYC verification attempts")
		}
	}
}

// validateKYCVerificationAttempts checks if a user has exceeded the maximum number of allowed verification attempts
func (automation *KYCIdentityVerificationRestrictionAutomation) validateKYCVerificationAttempts(verification common.KYCVerification) bool {
	currentAttempts := automation.userKYCVerificationAttempts[verification.UserID]
	if currentAttempts+1 > MaxKYCVerificationRetries {
		return false
	}

	// Update the verification attempt count for the user
	automation.userKYCVerificationAttempts[verification.UserID]++
	return true
}

// hasKYCExpired checks if a user's KYC verification has expired
func (automation *KYCIdentityVerificationRestrictionAutomation) hasKYCExpired(userID string) bool {
	expiryDate, exists := automation.userKYCExpiryDates[userID]
	if !exists {
		return false
	}
	return time.Now().After(expiryDate)
}

// flagKYCViolation flags a KYC verification violation and logs it in the ledger
func (automation *KYCIdentityVerificationRestrictionAutomation) flagKYCViolation(verification common.KYCVerification, reason string) {
	fmt.Printf("KYC verification violation: User %s, Reason: %s\n", verification.UserID, reason)

	// Log the violation into the ledger
	automation.logKYCViolation(verification, reason)
}

// logKYCViolation logs the flagged KYC verification violation into the ledger with full details
func (automation *KYCIdentityVerificationRestrictionAutomation) logKYCViolation(verification common.KYCVerification, violationReason string) {
	// Encrypt the KYC verification data before logging
	encryptedData := automation.encryptKYCVerificationData(verification)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("kyc-verification-violation-%s-%d", verification.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "KYC Verification Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for KYC verification violation. Reason: %s. Encrypted Data: %s", verification.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log KYC verification violation into ledger: %v\n", err)
	} else {
		fmt.Printf("KYC verification violation logged for user: %s\n", verification.UserID)
	}
}

// encryptKYCVerificationData encrypts KYC verification data before logging for security
func (automation *KYCIdentityVerificationRestrictionAutomation) encryptKYCVerificationData(verification common.KYCVerification) string {
	data := fmt.Sprintf("User ID: %s, Verification Status: %s, Timestamp: %d", verification.UserID, verification.Status, verification.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting KYC verification data:", err)
		return data
	}
	return string(encryptedData)
}
