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

// Configuration for dApp security enforcement automation
const (
	SecurityCheckInterval        = 10 * time.Second // Interval to check security compliance of dApps
	MaxSecurityViolations        = 3                // Maximum security violations before restricting dApp
)

// DAppSecurityEnforcementAutomation monitors and enforces security compliance for dApps
type DAppSecurityEnforcementAutomation struct {
	securityManager   *security.SecurityManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	violationCount    map[string]int // Tracks security violations per dApp
}

// NewDAppSecurityEnforcementAutomation initializes the dApp security enforcement automation
func NewDAppSecurityEnforcementAutomation(securityManager *security.SecurityManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DAppSecurityEnforcementAutomation {
	return &DAppSecurityEnforcementAutomation{
		securityManager:  securityManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartSecurityEnforcement begins continuous monitoring and enforcement of dApp security compliance
func (automation *DAppSecurityEnforcementAutomation) StartSecurityEnforcement() {
	ticker := time.NewTicker(SecurityCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkDAppSecurity()
		}
	}()
}

// checkDAppSecurity monitors each dApp's security compliance and applies enforcement if thresholds are exceeded
func (automation *DAppSecurityEnforcementAutomation) checkDAppSecurity() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, dAppID := range automation.securityManager.GetRegisteredDApps() {
		compliant, err := automation.securityManager.ValidateDAppSecurity(dAppID)
		if err != nil {
			fmt.Printf("Error validating security for dApp %s: %v\n", dAppID, err)
			automation.logSecurityAction(dAppID, "Security Validation Error")
			continue
		}

		if !compliant {
			automation.handleSecurityViolation(dAppID)
		}
	}
}

// handleSecurityViolation applies restrictions to dApps with repeated security violations
func (automation *DAppSecurityEnforcementAutomation) handleSecurityViolation(dAppID string) {
	automation.violationCount[dAppID]++

	if automation.violationCount[dAppID] >= MaxSecurityViolations {
		err := automation.securityManager.RestrictDApp(dAppID)
		if err != nil {
			fmt.Printf("Failed to restrict dApp %s due to security violations: %v\n", dAppID, err)
			automation.logSecurityAction(dAppID, "Failed Security Restriction")
		} else {
			fmt.Printf("dApp %s restricted due to repeated security violations.\n", dAppID)
			automation.logSecurityAction(dAppID, "dApp Restricted for Security Violations")
			automation.violationCount[dAppID] = 0
		}
	} else {
		fmt.Printf("Security violation detected for dApp %s.\n", dAppID)
		automation.logSecurityAction(dAppID, "Security Violation Detected")
	}
}

// logSecurityAction securely logs actions related to dApp security enforcement
func (automation *DAppSecurityEnforcementAutomation) logSecurityAction(dAppID, action string) {
	entryDetails := fmt.Sprintf("Action: %s, dApp ID: %s", action, dAppID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("dapp-security-enforcement-%s-%d", dAppID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "dApp Security Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log security enforcement action for dApp %s: %v\n", dAppID, err)
	} else {
		fmt.Println("Security enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DAppSecurityEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualSecurityCheck allows administrators to manually check security compliance for a specific dApp
func (automation *DAppSecurityEnforcementAutomation) TriggerManualSecurityCheck(dAppID string) {
	fmt.Printf("Manually triggering security compliance check for dApp: %s\n", dAppID)

	compliant, err := automation.securityManager.ValidateDAppSecurity(dAppID)
	if err != nil {
		fmt.Printf("Failed to manually check security for dApp %s: %v\n", dAppID, err)
		automation.logSecurityAction(dAppID, "Manual Security Check Failed")
		return
	}

	if !compliant {
		automation.handleSecurityViolation(dAppID)
	} else {
		fmt.Printf("dApp %s is compliant with security standards.\n", dAppID)
		automation.logSecurityAction(dAppID, "Manual Security Check Passed")
	}
}
