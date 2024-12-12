package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/identity"
)

// Configuration for identity verification enforcement automation
const (
	IdentityVerificationCheckInterval = 20 * time.Second // Interval to check identity verification compliance
	UnverifiedActivityThreshold       = 5                // Number of suspicious activities allowed before action is taken
)

// IdentityVerificationEnforcementAutomation monitors and enforces compliance with identity verification requirements
type IdentityVerificationEnforcementAutomation struct {
	identityManager     *identity.IdentityManager
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	unverifiedActivityMap map[string]int // Tracks activities by unverified identities
	verifiedIdentityMap map[string]bool  // Tracks verified identities
}

// NewIdentityVerificationEnforcementAutomation initializes the identity verification enforcement automation
func NewIdentityVerificationEnforcementAutomation(identityManager *identity.IdentityManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *IdentityVerificationEnforcementAutomation {
	return &IdentityVerificationEnforcementAutomation{
		identityManager:      identityManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		unverifiedActivityMap: make(map[string]int),
		verifiedIdentityMap:  make(map[string]bool),
	}
}

// StartIdentityVerificationEnforcement begins continuous monitoring and enforcement of identity verification
func (automation *IdentityVerificationEnforcementAutomation) StartIdentityVerificationEnforcement() {
	ticker := time.NewTicker(IdentityVerificationCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkIdentityVerificationCompliance()
		}
	}()
}

// checkIdentityVerificationCompliance monitors identities and restricts unverified identities from critical actions
func (automation *IdentityVerificationEnforcementAutomation) checkIdentityVerificationCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyIdentities()
	automation.enforceUnverifiedRestrictions()
}

// verifyIdentities updates the verified status of identities in the network
func (automation *IdentityVerificationEnforcementAutomation) verifyIdentities() {
	for _, userID := range automation.identityManager.GetAllIdentities() {
		isVerified := automation.identityManager.IsIdentityVerified(userID)
		automation.verifiedIdentityMap[userID] = isVerified
		if !isVerified {
			automation.unverifiedActivityMap[userID]++
		} else {
			automation.unverifiedActivityMap[userID] = 0 // Reset if verified
		}
	}
}

// enforceUnverifiedRestrictions restricts unverified identities exceeding activity thresholds
func (automation *IdentityVerificationEnforcementAutomation) enforceUnverifiedRestrictions() {
	for userID, activityCount := range automation.unverifiedActivityMap {
		if activityCount > UnverifiedActivityThreshold && !automation.verifiedIdentityMap[userID] {
			fmt.Printf("Identity verification enforcement triggered for user %s.\n", userID)
			automation.applyRestriction(userID, "Exceeded Unverified Activity Threshold")
		}
	}
}

// applyRestriction restricts access for unverified identities that violate thresholds
func (automation *IdentityVerificationEnforcementAutomation) applyRestriction(userID, reason string) {
	err := automation.identityManager.RestrictIdentity(userID)
	if err != nil {
		fmt.Printf("Failed to restrict unverified identity %s: %v\n", userID, err)
		automation.logIdentityAction(userID, "Restriction Failed", reason)
	} else {
		fmt.Printf("Unverified identity %s restricted due to %s.\n", userID, reason)
		automation.logIdentityAction(userID, "Restricted", reason)
	}
}

// logIdentityAction securely logs actions related to identity verification enforcement
func (automation *IdentityVerificationEnforcementAutomation) logIdentityAction(userID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, User ID: %s, Reason: %s", action, userID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("identity-enforcement-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Identity Verification Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log identity verification enforcement action for user %s: %v\n", userID, err)
	} else {
		fmt.Println("Identity verification enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *IdentityVerificationEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualIdentityVerification allows administrators to manually enforce restrictions on an unverified identity
func (automation *IdentityVerificationEnforcementAutomation) TriggerManualIdentityVerification(userID string) {
	fmt.Printf("Manually enforcing identity verification for user: %s\n", userID)

	if automation.verifiedIdentityMap[userID] {
		fmt.Printf("User %s is already verified.\n", userID)
		automation.logIdentityAction(userID, "Manual Verification Skipped - Already Verified", "Manual Check")
	} else {
		automation.applyRestriction(userID, "Manual Trigger: Unverified Identity Restriction")
	}
}
