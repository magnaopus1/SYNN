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

// Configuration for emergency failover enforcement automation
const (
	FailoverCheckInterval        = 5 * time.Second  // Interval to check for failover conditions
	MaxNodeFailurePercentage     = 30               // Maximum allowable percentage of failed nodes before failover
	MaxNetworkLatencyThreshold   = 200              // Maximum allowable network latency (ms)
)

// EmergencyFailoverEnforcementAutomation monitors and enforces emergency failover protocols for the network
type EmergencyFailoverEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	nodeFailureMap       map[string]bool // Tracks failed nodes
	networkLatency       int              // Tracks current network latency
}

// NewEmergencyFailoverEnforcementAutomation initializes the emergency failover enforcement automation
func NewEmergencyFailoverEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *EmergencyFailoverEnforcementAutomation {
	return &EmergencyFailoverEnforcementAutomation{
		networkManager:      networkManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		nodeFailureMap:      make(map[string]bool),
		networkLatency:      0,
	}
}

// StartFailoverEnforcement begins continuous monitoring and enforcement of failover protocols
func (automation *EmergencyFailoverEnforcementAutomation) StartFailoverEnforcement() {
	ticker := time.NewTicker(FailoverCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkFailoverConditions()
		}
	}()
}

// checkFailoverConditions monitors network health, latency, and node failures to trigger failover if needed
func (automation *EmergencyFailoverEnforcementAutomation) checkFailoverConditions() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.checkNodeFailures()
	automation.checkNetworkLatency()
}

// checkNodeFailures verifies if the percentage of failed nodes exceeds the maximum allowable threshold
func (automation *EmergencyFailoverEnforcementAutomation) checkNodeFailures() {
	totalNodes := automation.networkManager.GetTotalNodes()
	failedNodes := automation.getFailedNodeCount()

	failurePercentage := (failedNodes * 100) / totalNodes
	if failurePercentage > MaxNodeFailurePercentage {
		fmt.Printf("Emergency failover triggered due to %d%% node failures.\n", failurePercentage)
		automation.triggerFailover("Excessive Node Failures")
	}
}

// checkNetworkLatency monitors network latency and triggers failover if latency exceeds threshold
func (automation *EmergencyFailoverEnforcementAutomation) checkNetworkLatency() {
	latency := automation.networkManager.GetNetworkLatency()

	if latency > MaxNetworkLatencyThreshold {
		fmt.Printf("Emergency failover triggered due to high network latency: %d ms.\n", latency)
		automation.triggerFailover("High Network Latency")
	}
}

// getFailedNodeCount retrieves the count of nodes that have failed based on the node failure map
func (automation *EmergencyFailoverEnforcementAutomation) getFailedNodeCount() int {
	failedCount := 0
	for _, status := range automation.nodeFailureMap {
		if status {
			failedCount++
		}
	}
	return failedCount
}

// triggerFailover activates failover procedures to maintain network stability
func (automation *EmergencyFailoverEnforcementAutomation) triggerFailover(reason string) {
	err := automation.networkManager.ActivateFailover()
	if err != nil {
		fmt.Printf("Failed to activate failover due to %s: %v\n", reason, err)
		automation.logFailoverAction("Failover Activation Failed", reason)
	} else {
		fmt.Printf("Failover activated due to %s.\n", reason)
		automation.logFailoverAction("Failover Activated", reason)
	}
}

// logFailoverAction securely logs actions related to emergency failover enforcement
func (automation *EmergencyFailoverEnforcementAutomation) logFailoverAction(action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Reason: %s", action, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("failover-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Failover Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log failover enforcement action: %v\n", err)
	} else {
		fmt.Println("Failover enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *EmergencyFailoverEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualFailoverCheck allows administrators to manually check and trigger failover for network stability
func (automation *EmergencyFailoverEnforcementAutomation) TriggerManualFailoverCheck() {
	fmt.Println("Manually triggering emergency failover check.")

	if automation.getFailedNodeCount()*100/automation.networkManager.GetTotalNodes() > MaxNodeFailurePercentage {
		automation.triggerFailover("Manual Trigger: Excessive Node Failures")
	} else if automation.networkManager.GetNetworkLatency() > MaxNetworkLatencyThreshold {
		automation.triggerFailover("Manual Trigger: High Network Latency")
	} else {
		fmt.Println("Network is stable; no failover required.")
		automation.logFailoverAction("Manual Check Completed", "No Failover Triggered")
	}
}
