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

// Configuration for bandwidth enforcement automation
const (
	BandwidthCheckInterval     = 10 * time.Second // Interval to check bandwidth usage
	BandwidthLimitThreshold    = 5000             // Maximum allowable bandwidth (in MB) for nodes
	BandwidthAlertThreshold    = 4000             // Alert threshold for high bandwidth usage
	MaxBandwidthLimitViolations = 3               // Max violations before node restriction
)

// BandwidthLimitEnforcementAutomation enforces bandwidth limits across the network
type BandwidthLimitEnforcementAutomation struct {
	networkManager   *network.NetworkManager
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	enforcementMutex *sync.RWMutex
	violationCount   map[string]int // Track bandwidth violations per node
}

// NewBandwidthLimitEnforcementAutomation initializes the bandwidth limit automation
func NewBandwidthLimitEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *BandwidthLimitEnforcementAutomation {
	return &BandwidthLimitEnforcementAutomation{
		networkManager:   networkManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartBandwidthLimitEnforcement begins continuous monitoring and enforcement of bandwidth limits
func (automation *BandwidthLimitEnforcementAutomation) StartBandwidthLimitEnforcement() {
	ticker := time.NewTicker(BandwidthCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkBandwidthUsage()
		}
	}()
}

// checkBandwidthUsage monitors node bandwidth and applies restrictions if limits are exceeded
func (automation *BandwidthLimitEnforcementAutomation) checkBandwidthUsage() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, nodeID := range automation.networkManager.GetActiveNodes() {
		bandwidthUsed := automation.networkManager.GetNodeBandwidthUsage(nodeID)

		if bandwidthUsed >= BandwidthLimitThreshold {
			automation.handleBandwidthViolation(nodeID, bandwidthUsed)
		} else if bandwidthUsed >= BandwidthAlertThreshold {
			fmt.Printf("Warning: Node %s is approaching bandwidth limit with usage of %d MB.\n", nodeID, bandwidthUsed)
			automation.logBandwidthAction(nodeID, "High Bandwidth Usage Warning", bandwidthUsed)
		}
	}
}

// handleBandwidthViolation applies restrictions on nodes exceeding bandwidth limits
func (automation *BandwidthLimitEnforcementAutomation) handleBandwidthViolation(nodeID string, bandwidthUsed int) {
	automation.violationCount[nodeID]++

	if automation.violationCount[nodeID] >= MaxBandwidthLimitViolations {
		err := automation.networkManager.RestrictNodeBandwidth(nodeID)
		if err != nil {
			fmt.Printf("Failed to restrict bandwidth for node %s: %v\n", nodeID, err)
			automation.logBandwidthAction(nodeID, "Failed Bandwidth Restriction", bandwidthUsed)
		} else {
			fmt.Printf("Bandwidth restriction applied to node %s after %d violations.\n", nodeID, automation.violationCount[nodeID])
			automation.logBandwidthAction(nodeID, "Bandwidth Restricted", bandwidthUsed)
			automation.violationCount[nodeID] = 0
		}
	} else {
		fmt.Printf("Node %s has exceeded bandwidth limit with usage of %d MB.\n", nodeID, bandwidthUsed)
		automation.logBandwidthAction(nodeID, "Bandwidth Limit Exceeded", bandwidthUsed)
	}
}

// logBandwidthAction securely logs actions related to bandwidth enforcement
func (automation *BandwidthLimitEnforcementAutomation) logBandwidthAction(nodeID, action string, bandwidthUsed int) {
	entryDetails := fmt.Sprintf("Action: %s, Node: %s, Bandwidth Used: %d MB", action, nodeID, bandwidthUsed)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("bandwidth-enforcement-%s-%d", nodeID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Bandwidth Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log bandwidth enforcement action for node %s: %v\n", nodeID, err)
	} else {
		fmt.Println("Bandwidth enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *BandwidthLimitEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualBandwidthRestriction allows administrators to manually restrict bandwidth for a node
func (automation *BandwidthLimitEnforcementAutomation) TriggerManualBandwidthRestriction(nodeID string) {
	fmt.Printf("Manually triggering bandwidth restriction for node: %s\n", nodeID)

	err := automation.networkManager.RestrictNodeBandwidth(nodeID)
	if err != nil {
		fmt.Printf("Failed to manually restrict bandwidth for node %s: %v\n", nodeID, err)
		automation.logBandwidthAction(nodeID, "Manual Bandwidth Restriction Failed", 0)
	} else {
		fmt.Printf("Manual bandwidth restriction applied to node %s.\n", nodeID)
		automation.logBandwidthAction(nodeID, "Manual Bandwidth Restriction", 0)
	}
}
