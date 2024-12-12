package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/access_control"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	AccessControlCheckInterval = 5 * time.Minute // Interval for checking access control compliance
)

// AccessControlEnforcementAutomation enforces access control policies across the network
type AccessControlEnforcementAutomation struct {
	accessControlManager *access_control.AccessControlManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
}

// NewAccessControlEnforcementAutomation initializes the access control enforcement automation
func NewAccessControlEnforcementAutomation(accessControlManager *access_control.AccessControlManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *AccessControlEnforcementAutomation {
	return &AccessControlEnforcementAutomation{
		accessControlManager: accessControlManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
	}
}

// StartAccessControlEnforcement starts the automation for continuous access control enforcement
func (automation *AccessControlEnforcementAutomation) StartAccessControlEnforcement() {
	ticker := time.NewTicker(AccessControlCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAndEnforceAccessControl()
		}
	}()
}

// checkAndEnforceAccessControl verifies access control policies and enforces them if needed
func (automation *AccessControlEnforcementAutomation) checkAndEnforceAccessControl() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	// Fetch all active users and their access levels
	activeUsers := automation.accessControlManager.GetActiveUsers()

	for _, user := range activeUsers {
		if !automation.accessControlManager.HasValidAccess(user.ID) {
			automation.revokeAccess(user.ID)
		} else {
			fmt.Printf("User %s has valid access.\n", user.ID)
		}
	}
}

// revokeAccess enforces access control by revoking invalid users' access and logs it into the ledger
func (automation *AccessControlEnforcementAutomation) revokeAccess(userID string) {
	fmt.Printf("Revoking access for user %s due to invalid access.\n", userID)
	err := automation.accessControlManager.RevokeAccess(userID)
	if err != nil {
		fmt.Printf("Failed to revoke access for user %s: %v\n", userID, err)
	} else {
		automation.logAccessRevocation(userID)
	}
}

// logAccessRevocation securely logs access revocation events into the ledger
func (automation *AccessControlEnforcementAutomation) logAccessRevocation(userID string) {
	entryDetails := fmt.Sprintf("Access revoked for user %s due to invalid credentials or expired permissions.", userID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("access-revocation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Access Revocation",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log access revocation for user %s in the ledger: %v\n", userID, err)
	} else {
		fmt.Println("Access revocation successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *AccessControlEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualAccessRevocation allows administrators to manually revoke a user's access
func (automation *AccessControlEnforcementAutomation) TriggerManualAccessRevocation(userID string) {
	fmt.Printf("Manually revoking access for user %s.\n", userID)
	err := automation.accessControlManager.RevokeAccess(userID)
	if err != nil {
		fmt.Printf("Failed to manually revoke access for user %s: %v\n", userID, err)
	} else {
		automation.logAccessRevocation(userID)
	}
}
