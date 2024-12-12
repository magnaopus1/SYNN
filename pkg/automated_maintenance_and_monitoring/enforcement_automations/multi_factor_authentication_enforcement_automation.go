package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/security"
)

// Configuration for MFA enforcement automation
const (
	MFACheckInterval        = 20 * time.Second // Interval to check for MFA compliance
	MaxFailedAttempts       = 3                // Max failed login attempts before enforcing MFA compliance
)

// MFAEnforcementAutomation monitors and enforces MFA compliance requirements
type MFAEnforcementAutomation struct {
	securityManager       *security.SecurityManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	mfaComplianceMap      map[string]bool  // Tracks MFA compliance status of users
	failedAttemptCount    map[string]int   // Tracks failed login attempts per user
}

// NewMFAEnforcementAutomation initializes the MFA enforcement automation
func NewMFAEnforcementAutomation(securityManager *security.SecurityManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *MFAEnforcementAutomation {
	return &MFAEnforcementAutomation{
		securityManager:       securityManager,
		consensusEngine:       consensusEngine,
		ledgerInstance:        ledgerInstance,
		enforcementMutex:      enforcementMutex,
		mfaComplianceMap:      make(map[string]bool),
		failedAttemptCount:    make(map[string]int),
	}
}

// StartMFAEnforcement begins continuous monitoring and enforcement of MFA compliance
func (automation *MFAEnforcementAutomation) StartMFAEnforcement() {
	ticker := time.NewTicker(MFACheckInterval)

	go func() {
		for range ticker.C {
			automation.checkMFACompliance()
		}
	}()
}

// checkMFACompliance monitors each user's MFA status and restricts access for non-compliant accounts
func (automation *MFAEnforcementAutomation) checkMFACompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyMFAStatus()
	automation.enforceMFARestrictions()
}

// verifyMFAStatus updates the MFA compliance status of users in the network
func (automation *MFAEnforcementAutomation) verifyMFAStatus() {
	for _, userID := range automation.securityManager.GetAllUsers() {
		isMFAEnabled := automation.securityManager.IsMFAEnabled(userID)
		automation.mfaComplianceMap[userID] = isMFAEnabled

		if !isMFAEnabled {
			automation.failedAttemptCount[userID]++
		} else {
			automation.failedAttemptCount[userID] = 0 // Reset if MFA is enabled
		}
	}
}

// enforceMFARestrictions restricts users who exceed failed attempts and do not comply with MFA requirements
func (automation *MFAEnforcementAutomation) enforceMFARestrictions() {
	for userID, attempts := range automation.failedAttemptCount {
		if attempts > MaxFailedAttempts && !automation.mfaComplianceMap[userID] {
			fmt.Printf("MFA enforcement triggered for user %s due to repeated failed attempts.\n", userID)
			automation.applyRestriction(userID, "Exceeded Failed Login Attempts")
		}
	}
}

// applyRestriction restricts access for users who fail to meet MFA compliance standards
func (automation *MFAEnforcementAutomation) applyRestriction(userID, reason string) {
	err := automation.securityManager.RestrictAccess(userID)
	if err != nil {
		fmt.Printf("Failed to restrict non-compliant user %s: %v\n", userID, err)
		automation.logMFAAction(userID, "Restriction Failed", reason)
	} else {
		fmt.Printf("User %s restricted due to %s.\n", userID, reason)
		automation.logMFAAction(userID, "Restricted", reason)
	}
}

// logMFAAction securely logs actions related to MFA enforcement
func (automation *MFAEnforcementAutomation) logMFAAction(userID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, User ID: %s, Reason: %s", action, userID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("mfa-enforcement-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "MFA Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log MFA enforcement action for user %s: %v\n", userID, err)
	} else {
		fmt.Println("MFA enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *MFAEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualMFAEnforcement allows administrators to manually restrict a non-compliant user for failing MFA compliance
func (automation *MFAEnforcementAutomation) TriggerManualMFAEnforcement(userID string) {
	fmt.Printf("Manually enforcing MFA compliance for user: %s\n", userID)

	if automation.mfaComplianceMap[userID] {
		fmt.Printf("User %s is already MFA compliant.\n", userID)
		automation.logMFAAction(userID, "Manual Enforcement Skipped - Already Compliant", "Manual Check")
	} else {
		automation.applyRestriction(userID, "Manual Trigger: MFA Non-Compliance Restriction")
	}
}
