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

// Configuration for gas usage enforcement automation
const (
	GasUsageCheckInterval      = 10 * time.Second // Interval to check gas usage compliance
	MaxAllowableGasPerTx       = 500000           // Maximum gas allowed per transaction
	MaxTotalGasPerBlock        = 500000000        // Maximum total gas allowed per block
	MaxAverageGasPerSubBlock   = 5000000          // Maximum average gas per sub-block
)

// GasUsageEnforcementAutomation monitors and enforces gas usage limits across the network
type GasUsageEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	transactionGasMap    map[string]int // Tracks gas usage per transaction
	blockGasUsageMap     map[int]int    // Tracks gas usage per block
	subBlockGasUsageMap  map[int]int    // Tracks gas usage per sub-block
}

// NewGasUsageEnforcementAutomation initializes the gas usage enforcement automation
func NewGasUsageEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *GasUsageEnforcementAutomation {
	return &GasUsageEnforcementAutomation{
		networkManager:      networkManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		transactionGasMap:   make(map[string]int),
		blockGasUsageMap:    make(map[int]int),
		subBlockGasUsageMap: make(map[int]int),
	}
}

// StartGasUsageEnforcement begins continuous monitoring and enforcement of gas usage compliance
func (automation *GasUsageEnforcementAutomation) StartGasUsageEnforcement() {
	ticker := time.NewTicker(GasUsageCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkGasUsageCompliance()
		}
	}()
}

// checkGasUsageCompliance monitors transactions, sub-blocks, and blocks to enforce gas usage limits
func (automation *GasUsageEnforcementAutomation) checkGasUsageCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.verifyTransactionGasUsage()
	automation.verifySubBlockGasUsage()
	automation.verifyBlockGasUsage()
}

// verifyTransactionGasUsage ensures that individual transactions do not exceed the gas limit
func (automation *GasUsageEnforcementAutomation) verifyTransactionGasUsage() {
	for txID, gasUsage := range automation.transactionGasMap {
		if gasUsage > MaxAllowableGasPerTx {
			fmt.Printf("Gas usage violation for transaction %s.\n", txID)
			automation.applyGasRestriction(txID, "Transaction Gas Limit Exceeded")
		}
	}
}

// verifySubBlockGasUsage ensures that sub-blocks do not exceed the average gas usage limit
func (automation *GasUsageEnforcementAutomation) verifySubBlockGasUsage() {
	for subBlockID, gasUsage := range automation.subBlockGasUsageMap {
		if gasUsage > MaxAverageGasPerSubBlock {
			fmt.Printf("Gas usage violation for sub-block %d.\n", subBlockID)
			automation.applySubBlockRestriction(subBlockID, "Sub-Block Gas Limit Exceeded")
		}
	}
}

// verifyBlockGasUsage ensures that entire blocks do not exceed the total gas usage limit
func (automation *GasUsageEnforcementAutomation) verifyBlockGasUsage() {
	for blockID, gasUsage := range automation.blockGasUsageMap {
		if gasUsage > MaxTotalGasPerBlock {
			fmt.Printf("Gas usage violation for block %d.\n", blockID)
			automation.applyBlockRestriction(blockID, "Block Gas Limit Exceeded")
		}
	}
}

// applyGasRestriction restricts transactions that exceed gas limits
func (automation *GasUsageEnforcementAutomation) applyGasRestriction(txID, reason string) {
	err := automation.networkManager.RestrictTransaction(txID)
	if err != nil {
		fmt.Printf("Failed to restrict transaction %s due to gas limit violation: %v\n", txID, err)
		automation.logGasAction(txID, "Transaction Restriction Failed", reason)
	} else {
		fmt.Printf("Transaction %s restricted due to %s.\n", txID, reason)
		automation.logGasAction(txID, "Transaction Restricted", reason)
	}
}

// applySubBlockRestriction restricts sub-blocks that exceed average gas limits
func (automation *GasUsageEnforcementAutomation) applySubBlockRestriction(subBlockID int, reason string) {
	err := automation.consensusEngine.RestrictSubBlock(subBlockID)
	if err != nil {
		fmt.Printf("Failed to restrict sub-block %d due to gas limit violation: %v\n", subBlockID, err)
		automation.logGasAction(fmt.Sprintf("SubBlock-%d", subBlockID), "Sub-Block Restriction Failed", reason)
	} else {
		fmt.Printf("Sub-block %d restricted due to %s.\n", subBlockID, reason)
		automation.logGasAction(fmt.Sprintf("SubBlock-%d", subBlockID), "Sub-Block Restricted", reason)
	}
}

// applyBlockRestriction restricts blocks that exceed total gas limits
func (automation *GasUsageEnforcementAutomation) applyBlockRestriction(blockID int, reason string) {
	err := automation.consensusEngine.RestrictBlock(blockID)
	if err != nil {
		fmt.Printf("Failed to restrict block %d due to gas limit violation: %v\n", blockID, err)
		automation.logGasAction(fmt.Sprintf("Block-%d", blockID), "Block Restriction Failed", reason)
	} else {
		fmt.Printf("Block %d restricted due to %s.\n", blockID, reason)
		automation.logGasAction(fmt.Sprintf("Block-%d", blockID), "Block Restricted", reason)
	}
}

// logGasAction securely logs actions related to gas usage enforcement
func (automation *GasUsageEnforcementAutomation) logGasAction(entityID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Entity ID: %s, Reason: %s", action, entityID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("gas-usage-enforcement-%s-%d", entityID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Gas Usage Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log gas usage enforcement action for entity %s: %v\n", entityID, err)
	} else {
		fmt.Println("Gas usage enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *GasUsageEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualGasCheck allows administrators to manually check gas usage compliance for a specific entity
func (automation *GasUsageEnforcementAutomation) TriggerManualGasCheck(entityID string, entityType string) {
	fmt.Printf("Manually triggering gas usage compliance check for entity: %s\n", entityID)

	switch entityType {
	case "transaction":
		gasUsage := automation.transactionGasMap[entityID]
		if gasUsage > MaxAllowableGasPerTx {
			automation.applyGasRestriction(entityID, "Manual Trigger: Transaction Gas Limit Exceeded")
		}
	case "sub-block":
		subBlockID := automation.parseID(entityID)
		gasUsage := automation.subBlockGasUsageMap[subBlockID]
		if gasUsage > MaxAverageGasPerSubBlock {
			automation.applySubBlockRestriction(subBlockID, "Manual Trigger: Sub-Block Gas Limit Exceeded")
		}
	case "block":
		blockID := automation.parseID(entityID)
		gasUsage := automation.blockGasUsageMap[blockID]
		if gasUsage > MaxTotalGasPerBlock {
			automation.applyBlockRestriction(blockID, "Manual Trigger: Block Gas Limit Exceeded")
		}
	default:
		fmt.Printf("Unknown entity type for gas usage compliance check: %s\n", entityType)
		automation.logGasAction(entityID, "Manual Compliance Check Failed", "Invalid Entity Type")
	}
}

// parseID is a helper function to parse sub-block or block ID from string
func (automation *GasUsageEnforcementAutomation) parseID(entityID string) int {
	// Implement parsing logic if needed; return integer representation of ID
	return 0 // Placeholder - Implement ID parsing logic if necessary
}
