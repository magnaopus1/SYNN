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

// Configuration for decentralization threshold enforcement
const (
	DecentralizationCheckInterval = 20 * time.Second // Interval to check decentralization metrics
	MinValidatorNodePercentage    = 60               // Minimum percentage of nodes required to be validators
	MaxCentralizedNodeControl     = 10               // Maximum allowable nodes controlled by a single entity
)

// DecentralizationThresholdEnforcementAutomation monitors and enforces decentralization thresholds across the network
type DecentralizationThresholdEnforcementAutomation struct {
	networkManager    *network.NetworkManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	entityControlMap  map[string]int // Tracks the number of nodes controlled by each entity
	violationCount    map[string]int // Tracks decentralization violations per entity
}

// NewDecentralizationThresholdEnforcementAutomation initializes the decentralization threshold enforcement automation
func NewDecentralizationThresholdEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DecentralizationThresholdEnforcementAutomation {
	return &DecentralizationThresholdEnforcementAutomation{
		networkManager:   networkManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		entityControlMap: make(map[string]int),
		violationCount:   make(map[string]int),
	}
}

// StartDecentralizationEnforcement begins continuous monitoring and enforcement of decentralization compliance
func (automation *DecentralizationThresholdEnforcementAutomation) StartDecentralizationEnforcement() {
	ticker := time.NewTicker(DecentralizationCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkDecentralizationThresholds()
		}
	}()
}

// checkDecentralizationThresholds monitors the network’s node distribution and validator participation to enforce decentralization standards
func (automation *DecentralizationThresholdEnforcementAutomation) checkDecentralizationThresholds() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyValidatorNodePercentage()
	automation.checkEntityControl()
}

// verifyValidatorNodePercentage checks if the percentage of validator nodes meets the minimum requirement
func (automation *DecentralizationThresholdEnforcementAutomation) verifyValidatorNodePercentage() {
	totalNodes := automation.networkManager.GetTotalNodes()
	validatorNodes := automation.networkManager.GetValidatorNodesCount()

	validatorPercentage := (validatorNodes * 100) / totalNodes
	if validatorPercentage < MinValidatorNodePercentage {
		fmt.Printf("Decentralization violation: Validator nodes are below the required %d%% threshold.\n", MinValidatorNodePercentage)
		automation.logDecentralizationAction("Validator Node Percentage Violation", validatorPercentage)
	}
}

// checkEntityControl ensures no single entity controls an excessive portion of the network
func (automation *DecentralizationThresholdEnforcementAutomation) checkEntityControl() {
	for _, nodeID := range automation.networkManager.GetAllNodes() {
		entityID := automation.networkManager.GetNodeEntityID(nodeID)
		automation.entityControlMap[entityID]++

		if automation.entityControlMap[entityID] > MaxCentralizedNodeControl {
			fmt.Printf("Decentralization violation: Entity %s controls more than %d nodes.\n", entityID, MaxCentralizedNodeControl)
			automation.handleCentralizationViolation(entityID)
		}
	}
}

// handleCentralizationViolation restricts additional nodes for entities controlling an excessive number
func (automation *DecentralizationThresholdEnforcementAutomation) handleCentralizationViolation(entityID string) {
	automation.violationCount[entityID]++

	if automation.violationCount[entityID] > MaxCentralizedNodeControl {
		err := automation.networkManager.RestrictAdditionalNodes(entityID)
		if err != nil {
			fmt.Printf("Failed to restrict nodes for entity %s due to centralization violation: %v\n", entityID, err)
			automation.logDecentralizationAction("Failed Node Restriction for Centralization", entityID)
		} else {
			fmt.Printf("Restricted additional nodes for entity %s due to repeated centralization violations.\n", entityID)
			automation.logDecentralizationAction("Node Restriction for Centralization Violation", entityID)
			automation.violationCount[entityID] = 0
		}
	}
}

// logDecentralizationAction securely logs actions related to decentralization enforcement
func (automation *DecentralizationThresholdEnforcementAutomation) logDecentralizationAction(action string, detail interface{}) {
	entryDetails := fmt.Sprintf("Action: %s, Details: %v", action, detail)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("decentralization-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Decentralization Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log decentralization enforcement action: %v\n", err)
	} else {
		fmt.Println("Decentralization enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DecentralizationThresholdEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualDecentralizationCheck allows administrators to manually check decentralization compliance for specific entities
func (automation *DecentralizationThresholdEnforcementAutomation) TriggerManualDecentralizationCheck(entityID string) {
	fmt.Printf("Manually triggering decentralization compliance check for entity: %s\n", entityID)

	controlCount := automation.entityControlMap[entityID]
	if controlCount > MaxCentralizedNodeControl {
		automation.handleCentralizationViolation(entityID)
	} else {
		fmt.Printf("Entity %s is compliant with decentralization standards.\n", entityID)
		automation.logDecentralizationAction("Manual Compliance Check Passed", entityID)
	}
}
