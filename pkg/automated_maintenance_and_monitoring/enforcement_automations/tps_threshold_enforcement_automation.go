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

// Configuration for TPS threshold enforcement automation
const (
	TPSCheckInterval       = 10 * time.Second // Interval to check for TPS threshold
	MaxAllowedTPS          = 5000             // Maximum allowed transactions per second
	TPSScalingThreshold    = 0.9              // Scaling trigger when TPS reaches 90% of max
	TPSViolationThreshold  = 3                // Number of times TPS can exceed max before enforcement
)

// TPSThresholdEnforcementAutomation monitors and enforces TPS limits to prevent network overload
type TPSThresholdEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	tpsViolationCount    int // Tracks consecutive TPS violations
}

// NewTPSThresholdEnforcementAutomation initializes the TPS threshold enforcement automation
func NewTPSThresholdEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *TPSThresholdEnforcementAutomation {
	return &TPSThresholdEnforcementAutomation{
		networkManager:    networkManager,
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		enforcementMutex:  enforcementMutex,
		tpsViolationCount: 0,
	}
}

// StartTPSThresholdEnforcement begins continuous monitoring and enforcement of TPS limits
func (automation *TPSThresholdEnforcementAutomation) StartTPSThresholdEnforcement() {
	ticker := time.NewTicker(TPSCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkTPSThreshold()
		}
	}()
}

// checkTPSThreshold monitors the network's TPS and takes action if thresholds are exceeded
func (automation *TPSThresholdEnforcementAutomation) checkTPSThreshold() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	currentTPS := automation.networkManager.GetCurrentTPS()
	fmt.Printf("Current TPS: %d\n", currentTPS)

	if currentTPS > int(TPSScalingThreshold*float64(MaxAllowedTPS)) {
		automation.scaleNetworkResources()
	}

	if currentTPS > MaxAllowedTPS {
		automation.tpsViolationCount++
		fmt.Printf("TPS threshold exceeded! Violation count: %d\n", automation.tpsViolationCount)
	} else {
		automation.tpsViolationCount = 0 // Reset if within limits
	}

	if automation.tpsViolationCount > TPSViolationThreshold {
		automation.restrictNetworkTraffic()
		automation.tpsViolationCount = 0 // Reset after enforcement
	}
}

// scaleNetworkResources scales network resources when TPS reaches a critical threshold
func (automation *TPSThresholdEnforcementAutomation) scaleNetworkResources() {
	err := automation.networkManager.ScaleResources()
	if err != nil {
		fmt.Printf("Failed to scale network resources: %v\n", err)
		automation.logTPSAction("Scaling Failed", fmt.Sprintf("Current TPS: %d", automation.networkManager.GetCurrentTPS()))
	} else {
		fmt.Println("Network resources scaled to manage increased TPS.")
		automation.logTPSAction("Resources Scaled", "Scaling successful due to high TPS")
	}
}

// restrictNetworkTraffic restricts traffic when TPS exceeds allowed limit persistently
func (automation *TPSThresholdEnforcementAutomation) restrictNetworkTraffic() {
	err := automation.networkManager.ApplyTrafficRestrictions()
	if err != nil {
		fmt.Printf("Failed to restrict network traffic: %v\n", err)
		automation.logTPSAction("Traffic Restriction Failed", "Excessive TPS without restriction")
	} else {
		fmt.Println("Network traffic restrictions applied due to excessive TPS.")
		automation.logTPSAction("Traffic Restricted", "TPS exceeded max limit persistently")
	}
}

// logTPSAction securely logs actions related to TPS enforcement
func (automation *TPSThresholdEnforcementAutomation) logTPSAction(action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Details: %s", action, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("tps-threshold-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "TPS Threshold Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log TPS enforcement action: %v\n", err)
	} else {
		fmt.Println("TPS enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *TPSThresholdEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualTPSScaling allows administrators to manually scale network resources
func (automation *TPSThresholdEnforcementAutomation) TriggerManualTPSScaling() {
	fmt.Println("Manually triggering TPS scaling.")
	automation.scaleNetworkResources()
}
