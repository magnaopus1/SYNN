package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for transaction speed enforcement automation
const (
	SpeedCheckInterval         = 10 * time.Second // Interval to check transaction speeds
	MaxTransactionTime         = 2 * time.Second  // Maximum allowed time for processing a transaction
	SpeedViolationThreshold    = 3                // Number of violations allowed before enforcement action
	SpeedScalingFactor         = 1.5              // Multiplier for scaling resources in response to slow speeds
)

// TransactionSpeedEnforcementAutomation monitors and enforces transaction processing speed standards
type TransactionSpeedEnforcementAutomation struct {
	networkManager        *network.NetworkManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	speedViolationCount   int // Tracks consecutive speed violations
}

// NewTransactionSpeedEnforcementAutomation initializes the transaction speed enforcement automation
func NewTransactionSpeedEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *TransactionSpeedEnforcementAutomation {
	return &TransactionSpeedEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		speedViolationCount:  0,
	}
}

// StartTransactionSpeedEnforcement begins continuous monitoring and enforcement of transaction speeds
func (automation *TransactionSpeedEnforcementAutomation) StartTransactionSpeedEnforcement() {
	ticker := time.NewTicker(SpeedCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkTransactionSpeed()
		}
	}()
}

// checkTransactionSpeed monitors the network's average transaction speed and takes action if speeds are too slow
func (automation *TransactionSpeedEnforcementAutomation) checkTransactionSpeed() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	avgProcessingTime := automation.networkManager.GetAverageTransactionTime()
	fmt.Printf("Current Average Transaction Processing Time: %s\n", avgProcessingTime)

	if avgProcessingTime > MaxTransactionTime {
		automation.speedViolationCount++
		fmt.Printf("Transaction processing time exceeded limit! Violation count: %d\n", automation.speedViolationCount)
	} else {
		automation.speedViolationCount = 0 // Reset if within limits
	}

	if automation.speedViolationCount > SpeedViolationThreshold {
		automation.scaleNetworkResources()
		automation.speedViolationCount = 0 // Reset after enforcement
	}
}

// scaleNetworkResources scales network resources when transaction speeds fall below acceptable levels
func (automation *TransactionSpeedEnforcementAutomation) scaleNetworkResources() {
	err := automation.networkManager.ScaleResources(SpeedScalingFactor)
	if err != nil {
		fmt.Printf("Failed to scale network resources: %v\n", err)
		automation.logSpeedEnforcementAction("Scaling Failed", fmt.Sprintf("Current Avg Speed: %s", automation.networkManager.GetAverageTransactionTime()))
	} else {
		fmt.Println("Network resources scaled to improve transaction processing speed.")
		automation.logSpeedEnforcementAction("Resources Scaled", "Scaling successful due to slow transaction speeds")
	}
}

// logSpeedEnforcementAction securely logs actions related to transaction speed enforcement
func (automation *TransactionSpeedEnforcementAutomation) logSpeedEnforcementAction(action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Details: %s", action, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("transaction-speed-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Transaction Speed Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log transaction speed enforcement action: %v\n", err)
	} else {
		fmt.Println("Transaction speed enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TransactionSpeedEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualScaling allows administrators to manually scale network resources for transaction speed optimization
func (automation *TransactionSpeedEnforcementAutomation) TriggerManualScaling() {
	fmt.Println("Manually triggering scaling for transaction speed optimization.")
	automation.scaleNetworkResources()
}
