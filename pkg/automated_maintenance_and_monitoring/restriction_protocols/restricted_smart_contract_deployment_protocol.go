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
	SmartContractDeploymentCheckInterval = 15 * time.Second // Interval for checking restricted smart contract deployments
	MaxDeploymentAttempts                = 3                // Maximum allowed unauthorized deployment attempts
)

// RestrictedSmartContractDeploymentAutomation enforces restrictions on unauthorized smart contract deployments
type RestrictedSmartContractDeploymentAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	deploymentViolation   map[string]int // Tracks unauthorized smart contract deployment attempts by user
}

// NewRestrictedSmartContractDeploymentAutomation initializes an instance of RestrictedSmartContractDeploymentAutomation
func NewRestrictedSmartContractDeploymentAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedSmartContractDeploymentAutomation {
	return &RestrictedSmartContractDeploymentAutomation{
		consensusSystem:     consensusSystem,
		ledgerInstance:      ledgerInstance,
		stateMutex:          stateMutex,
		deploymentViolation: make(map[string]int),
	}
}

// StartDeploymentMonitoring starts continuous monitoring of smart contract deployments
func (automation *RestrictedSmartContractDeploymentAutomation) StartDeploymentMonitoring() {
	ticker := time.NewTicker(SmartContractDeploymentCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorSmartContractDeployments()
		}
	}()
}

// monitorSmartContractDeployments checks for unauthorized smart contract deployment attempts and enforces restrictions
func (automation *RestrictedSmartContractDeploymentAutomation) monitorSmartContractDeployments() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch deployment data from Synnergy Consensus
	deploymentData := automation.consensusSystem.GetSmartContractDeploymentData()

	for userID, deploymentStatus := range deploymentData {
		// Check if the deployment attempt is unauthorized
		if deploymentStatus == "unauthorized" {
			automation.flagDeploymentViolation(userID, "Unauthorized smart contract deployment attempt detected")
		}
	}
}

// flagDeploymentViolation flags an unauthorized smart contract deployment attempt and logs it in the ledger
func (automation *RestrictedSmartContractDeploymentAutomation) flagDeploymentViolation(userID string, reason string) {
	fmt.Printf("Smart contract deployment violation: User ID %s, Reason: %s\n", userID, reason)

	// Increment the violation count for the user
	automation.deploymentViolation[userID]++

	// Log the violation in the ledger
	automation.logDeploymentViolation(userID, reason)

	// Check if the user has exceeded the allowed number of deployment violations
	if automation.deploymentViolation[userID] >= MaxDeploymentAttempts {
		automation.restrictSmartContractDeployment(userID)
	}
}

// logDeploymentViolation logs the flagged deployment violation into the ledger with details
func (automation *RestrictedSmartContractDeploymentAutomation) logDeploymentViolation(userID string, violationReason string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("deployment-violation-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Deployment Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s violated smart contract deployment restrictions. Reason: %s", userID, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptDeploymentData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log smart contract deployment violation:", err)
	} else {
		fmt.Println("Smart contract deployment violation logged.")
	}
}

// restrictSmartContractDeployment restricts smart contract deployment for a user after exceeding violations
func (automation *RestrictedSmartContractDeploymentAutomation) restrictSmartContractDeployment(userID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("deployment-restriction-%s-%d", userID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Smart Contract Deployment Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("User %s has been restricted from deploying smart contracts due to repeated violations.", userID),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptDeploymentData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log smart contract deployment restriction:", err)
	} else {
		fmt.Println("Smart contract deployment restriction applied.")
	}
}

// encryptDeploymentData encrypts the smart contract deployment data before logging for security
func (automation *RestrictedSmartContractDeploymentAutomation) encryptDeploymentData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting deployment data:", err)
		return data
	}
	return string(encryptedData)
}
