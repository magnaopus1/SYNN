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

// Configuration for load balancing enforcement automation
const (
	LoadBalancingCheckInterval     = 15 * time.Second // Interval to check node load and performance
	MaxTransactionLoadPerNode      = 10000            // Maximum allowed transactions per node before redistributing
	HighLoadThreshold              = 8000             // Threshold for high load, triggering early redistribution
)

// LoadBalancingEnforcementAutomation monitors and enforces load balancing across nodes
type LoadBalancingEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	nodeLoadMap          map[string]int // Tracks transaction load for each node
	highLoadNodeCount    map[string]int // Tracks instances of high load for each node
}

// NewLoadBalancingEnforcementAutomation initializes the load balancing enforcement automation
func NewLoadBalancingEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *LoadBalancingEnforcementAutomation {
	return &LoadBalancingEnforcementAutomation{
		networkManager:      networkManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		nodeLoadMap:         make(map[string]int),
		highLoadNodeCount:   make(map[string]int),
	}
}

// StartLoadBalancingEnforcement begins continuous monitoring and enforcement of load balancing across nodes
func (automation *LoadBalancingEnforcementAutomation) StartLoadBalancingEnforcement() {
	ticker := time.NewTicker(LoadBalancingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkLoadBalancingCompliance()
		}
	}()
}

// checkLoadBalancingCompliance monitors each node's load and redistributes transactions if necessary
func (automation *LoadBalancingEnforcementAutomation) checkLoadBalancingCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.assessNodeLoad()
	automation.redistributeLoad()
}

// assessNodeLoad updates the load status of each node in the network
func (automation *LoadBalancingEnforcementAutomation) assessNodeLoad() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		load := automation.networkManager.GetNodeTransactionLoad(nodeID)
		automation.nodeLoadMap[nodeID] = load

		if load > HighLoadThreshold {
			automation.highLoadNodeCount[nodeID]++
		} else {
			automation.highLoadNodeCount[nodeID] = 0 // Reset if load is below threshold
		}
	}
}

// redistributeLoad triggers load redistribution for nodes that exceed load thresholds
func (automation *LoadBalancingEnforcementAutomation) redistributeLoad() {
	for nodeID, load := range automation.nodeLoadMap {
		if load > MaxTransactionLoadPerNode || automation.highLoadNodeCount[nodeID] > 2 {
			fmt.Printf("Load balancing enforcement triggered for node %s due to high transaction load.\n", nodeID)
			automation.applyLoadRedistribution(nodeID, load)
		}
	}
}

// applyLoadRedistribution redistributes transactions from an overloaded node to ensure network stability
func (automation *LoadBalancingEnforcementAutomation) applyLoadRedistribution(nodeID string, load int) {
	err := automation.networkManager.RedistributeLoad(nodeID)
	if err != nil {
		fmt.Printf("Failed to redistribute load for node %s: %v\n", nodeID, err)
		automation.logLoadBalancingAction(nodeID, "Redistribution Failed", fmt.Sprintf("Load: %d", load))
	} else {
		fmt.Printf("Load successfully redistributed from node %s.\n", nodeID)
		automation.logLoadBalancingAction(nodeID, "Redistributed", fmt.Sprintf("Load: %d", load))
	}
}

// logLoadBalancingAction securely logs actions related to load balancing enforcement
func (automation *LoadBalancingEnforcementAutomation) logLoadBalancingAction(nodeID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Node ID: %s, Details: %s", action, nodeID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("load-balancing-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Load Balancing Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log load balancing enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Load balancing enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *LoadBalancingEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLoadRedistribution allows administrators to manually enforce load redistribution for a specific node
func (automation *LoadBalancingEnforcementAutomation) TriggerManualLoadRedistribution(nodeID string) {
	fmt.Printf("Manually triggering load redistribution for node: %s\n", nodeID)

	load := automation.networkManager.GetNodeTransactionLoad(nodeID)
	if load > HighLoadThreshold {
		automation.applyLoadRedistribution(nodeID, load)
	} else {
		fmt.Printf("Node %s is within acceptable load limits.\n", nodeID)
		automation.logLoadBalancingAction(nodeID, "Manual Redistribution Skipped", "Load Within Limits")
	}
}
