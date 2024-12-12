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
	IdentityVerificationCheckInterval = 24 * time.Hour // Interval for checking identity verification status
	MaxRetriesForVerification         = 3              // Maximum retries for identity verification
)

// IdentityVerificationRestrictionAutomation monitors and restricts identity verification processes across the network
type IdentityVerificationRestrictionAutomation struct {
	consensusSystem      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	stateMutex           *sync.RWMutex
	userVerificationAttempts map[string]int // Tracks identity verification attempts per user
}

// NewIdentityVerificationRestrictionAutomation initializes and returns an instance of IdentityVerificationRestrictionAutomation
func NewIdentityVerificationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *IdentityVerificationRestrictionAutomation {
	return &IdentityVerificationRestrictionAutomation{
		consensusSystem:         consensusSystem,
		ledgerInstance:          ledgerInstance,
		stateMutex:              stateMutex,
		userVerificationAttempts: make(map[string]int),
	}
}

// StartIdentityVerificationMonitoring starts continuous monitoring of identity verification processes
func (automation *IdentityVerificationRestrictionAutomation) StartIdentityVerificationMonitoring() {
	ticker := time.NewTicker(IdentityVerificationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorIdentityVerification()
		}
	}()
}

// monitorIdentityVerification checks recent identity verification attempts and enforces restrictions
func (automation *IdentityVerificationRestrictionAutomation) monitorIdentityVerification() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent identity verification statuses from Synnergy Consensus
	recentVerifications := automation.consensusSystem.GetRecentIdentityVerifications()

	for _, verification := range recentVerifications {
		// Validate the number of verification attempts and apply restrictions if necessary
		if !automation.validateVerificationAttempts(verification) {
			automation.flagVerificationViolation(verification, "Exceeded maximum number of verification attempts")
		}
	}
}

// validateVerificationAttempts checks if a user has exceeded the maximum allowed verification attempts
func (automation *IdentityVerificationRestrictionAutomation) validateVerificationAttempts(verification common.IdentityVerification) bool {
	currentAttempts := automation.userVerificationAttempts[verification.UserID]
	if currentAttempts+1 > MaxRetriesForVerification {
		return false
	}

	// Update the verification attempt count for the user
	automation.userVerificationAttempts[verification.UserID]++
	return true
}

// flagVerificationViolation flags an identity verification violation and logs it in the ledger
func (automation *IdentityVerificationRestrictionAutomation) flagVerificationViolation(verification common.IdentityVerification, reason string) {
	fmt.Printf("Identity verification violation: User %s, Reason: %s\n", verification.UserID, reason)

	// Log the violation into the ledger
	automation.logVerificationViolation(verification, reason)
}

// logVerificationViolation logs the flagged identity verification violation into the ledger with full details
func (automation *IdentityVerificationRestrictionAutomation) logVerificationViolation(verification common.IdentityVerification, violationReason string) {
	// Encrypt the identity verification data before logging
	encryptedData := automation.encryptVerificationData(verification)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("identity-verification-violation-%s-%d", verification.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Identity Verification Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for identity verification violation. Reason: %s. Encrypted Data: %s", verification.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log identity verification violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Identity verification violation logged for user: %s\n", verification.UserID)
	}
}

// encryptVerificationData encrypts identity verification data before logging for security
func (automation *IdentityVerificationRestrictionAutomation) encryptVerificationData(verification common.IdentityVerification) string {
	data := fmt.Sprintf("User ID: %s, Verification Status: %s, Timestamp: %d", verification.UserID, verification.Status, verification.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting verification data:", err)
		return data
	}
	return string(encryptedData)
}
