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

// Configuration for dynamic throttling enforcement automation
const (
	ThrottlingCheckInterval      = 10 * time.Second // Interval to check for throttling adjustments
	MaxTransactionRate           = 10000           // Maximum transactions allowed per minute per node
	MaxNodeLoadPercentage        = 85              // Maximum load percentage before throttling
)

// DynamicThrottlingEnforcementAutomation monitors and enforces throttling based on network activity
type DynamicThrottlingEnforcementAutomation struct {
	networkManager    *network.NetworkManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	nodeTransactionCount map[string]int // Tracks transaction rate per node
	nodeLoadMap       map[string]int    // Tracks load percentage per node
}

// NewDynamicThrottlingEnforcementAutomation initializes the dynamic throttling enforcement automation
func NewDynamicThrottlingEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DynamicThrottlingEnforcementAutomation {
	return &DynamicThrottlingEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		nodeTransactionCount: make(map[string]int),
		nodeLoadMap:          make(map[string]int),
	}
}

// StartThrottlingEnforcement begins continuous monitoring and enforcement of throttling policies
func (automation *DynamicThrottlingEnforcementAutomation) StartThrottlingEnforcement() {
	ticker := time.NewTicker(ThrottlingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkThrottlingCompliance()
		}
	}()
}

// checkThrottlingCompliance monitors node transaction rates and load percentages to adjust throttling dynamically
func (automation *DynamicThrottlingEnforcementAutomation) checkThrottlingCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.adjustTransactionThrottling()
	automation.manageNodeLoad()
}

// adjustTransactionThrottling dynamically throttles nodes based on transaction rate compliance
func (automation *DynamicThrottlingEnforcementAutomation) adjustTransactionThrottling() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		txRate := automation.networkManager.GetNodeTransactionRate(nodeID)

		if txRate > MaxTransactionRate {
			fmt.Printf("Transaction rate violation for node %s.\n", nodeID)
			automation.applyThrottling(nodeID, "Transaction Rate Exceeded")
		}
	}
}

// manageNodeLoad dynamically throttles nodes if their load exceeds maximum allowable percentage
func (automation *DynamicThrottlingEnforcementAutomation) manageNodeLoad() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		nodeLoad := automation.networkManager.GetNodeLoad(nodeID)

		if nodeLoad > MaxNodeLoadPercentage {
			fmt.Printf("Node load percentage exceeded for node %s.\n", nodeID)
			automation.applyThrottling(nodeID, "Node Load Exceeded")
		}
	}
}

// applyThrottling adjusts throttling level for a specific node based on the violation type
func (automation *DynamicThrottlingEnforcementAutomation) applyThrottling(nodeID string, reason string) {
	err := automation.networkManager.ThrottleNode(nodeID)
	if err != nil {
		fmt.Printf("Failed to throttle node %s: %v\n", nodeID, err)
		automation.logThrottlingAction(nodeID, "Throttling Application Failed", reason)
	} else {
		fmt.Printf("Throttling applied to node %s due to %s.\n", nodeID, reason)
		automation.logThrottlingAction(nodeID, "Throttling Applied", reason)
	}
}

// logThrottlingAction securely logs actions related to dynamic throttling enforcement
func (automation *DynamicThrottlingEnforcementAutomation) logThrottlingAction(nodeID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Reason: %s", action, nodeID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("throttling-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Throttling Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log throttling enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Throttling enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DynamicThrottlingEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualThrottlingCheck allows administrators to manually check and apply throttling for a specific node
func (automation *DynamicThrottlingEnforcementAutomation) TriggerManualThrottlingCheck(nodeID string) {
	fmt.Printf("Manually triggering throttling check for node: %s\n", nodeID)

	txRate := automation.networkManager.GetNodeTransactionRate(nodeID)
	nodeLoad := automation.networkManager.GetNodeLoad(nodeID)

	if txRate > MaxTransactionRate {
		automation.applyThrottling(nodeID, "Manual Trigger: Transaction Rate Exceeded")
	} else if nodeLoad > MaxNodeLoadPercentage {
		automation.applyThrottling(nodeID, "Manual Trigger: Node Load Exceeded")
	} else {
		fmt.Printf("Node %s is compliant with throttling policies.\n", nodeID)
		automation.logThrottlingAction(nodeID, "Manual Compliance Check Passed", "Throttling Compliance Verified")
	}
}
