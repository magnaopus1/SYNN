package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/compliance"
)

// Configuration for KYC/AML enforcement automation
const (
	KYCAMLCheckInterval           = 30 * time.Second // Interval to check KYC/AML compliance
	MaxAllowedSuspiciousActivities = 3               // Maximum allowed suspicious activities before restriction
)

// KYCAMLEnforcementAutomation monitors and enforces KYC and AML compliance requirements
type KYCAMLEnforcementAutomation struct {
	complianceManager     *compliance.ComplianceManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	suspiciousActivityMap map[string]int // Tracks suspicious activity count for each user
	kycVerifiedMap        map[string]bool // Tracks KYC verification status of users
}

// NewKYCAMLEnforcementAutomation initializes the KYC/AML enforcement automation
func NewKYCAMLEnforcementAutomation(complianceManager *compliance.ComplianceManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *KYCAMLEnforcementAutomation {
	return &KYCAMLEnforcementAutomation{
		complianceManager:     complianceManager,
		consensusEngine:       consensusEngine,
		ledgerInstance:        ledgerInstance,
		enforcementMutex:      enforcementMutex,
		suspiciousActivityMap: make(map[string]int),
		kycVerifiedMap:        make(map[string]bool),
	}
}

// StartKYCAMLEnforcement begins continuous monitoring and enforcement of KYC/AML compliance
func (automation *KYCAMLEnforcementAutomation) StartKYCAMLEnforcement() {
	ticker := time.NewTicker(KYCAMLCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkKYCAMLCompliance()
		}
	}()
}

// checkKYCAMLCompliance monitors each user’s KYC/AML status and restricts non-compliant accounts
func (automation *KYCAMLEnforcementAutomation) checkKYCAMLCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyKYCCompliance()
	automation.enforceAMLRestrictions()
}

// verifyKYCCompliance updates the KYC status of users in the network
func (automation *KYCAMLEnforcementAutomation) verifyKYCCompliance() {
	for _, userID := range automation.complianceManager.GetAllUsers() {
		isKYCVerified := automation.complianceManager.IsKYCVerified(userID)
		automation.kycVerifiedMap[userID] = isKYCVerified
		if !isKYCVerified {
			automation.suspiciousActivityMap[userID]++
		} else {
			automation.suspiciousActivityMap[userID] = 0 // Reset if verified
		}
	}
}

// enforceAMLRestrictions restricts users who exceed suspicious activity thresholds
func (automation *KYCAMLEnforcementAutomation) enforceAMLRestrictions() {
	for userID, activityCount := range automation.suspiciousActivityMap {
		if activityCount > MaxAllowedSuspiciousActivities && !automation.kycVerifiedMap[userID] {
			fmt.Printf("KYC/AML enforcement triggered for user %s due to suspicious activities.\n", userID)
			automation.applyRestriction(userID, "Exceeded Suspicious Activity Threshold")
		}
	}
}

// applyRestriction restricts access for users who fail KYC/AML checks
func (automation *KYCAMLEnforcementAutomation) applyRestriction(userID, reason string) {
	err := automation.complianceManager.RestrictAccount(userID)
	if err != nil {
		fmt.Printf("Failed to restrict non-compliant user %s: %v\n", userID, err)
		automation.logComplianceAction(userID, "Restriction Failed", reason)
	} else {
		fmt.Printf("User %s restricted due to %s.\n", userID, reason)
		automation.logComplianceAction(userID, "Restricted", reason)
	}
}

// logComplianceAction securely logs actions related to KYC/AML enforcement
func (automation *KYCAMLEnforcementAutomation) logComplianceAction(userID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, User ID: %s, Reason: %s", action, userID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("kyc-aml-enforcement-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "KYC/AML Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log KYC/AML enforcement action for user %s: %v\n", userID, err)
	} else {
		fmt.Println("KYC/AML enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *KYCAMLEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualKYCAMLVerification allows administrators to manually restrict a non-compliant user
func (automation *KYCAMLEnforcementAutomation) TriggerManualKYCAMLVerification(userID string) {
	fmt.Printf("Manually enforcing KYC/AML compliance for user: %s\n", userID)

	if automation.kycVerifiedMap[userID] {
		fmt.Printf("User %s is already KYC verified.\n", userID)
		automation.logComplianceAction(userID, "Manual Verification Skipped - Already Verified", "Manual Check")
	} else {
		automation.applyRestriction(userID, "Manual Trigger: KYC/AML Non-Compliance Restriction")
	}
}
