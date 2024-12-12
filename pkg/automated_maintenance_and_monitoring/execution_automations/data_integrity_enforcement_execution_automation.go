package execution_automations

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	DataIntegrityCheckInterval     = 5 * time.Minute // Interval for checking data integrity
	DataVerificationThreshold      = 0.99            // Minimum threshold for data verification accuracy (99%)
	BlockSizeVerificationThreshold = 1000            // Threshold for the number of sub-blocks in a block
)

// DataIntegrityEnforcementAutomation automates data integrity enforcement
type DataIntegrityEnforcementAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance    *ledger.Ledger                        // Ledger for logging data integrity actions
	integrityMutex    *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewDataIntegrityEnforcementAutomation initializes the data integrity enforcement automation
func NewDataIntegrityEnforcementAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, integrityMutex *sync.RWMutex) *DataIntegrityEnforcementAutomation {
	return &DataIntegrityEnforcementAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		integrityMutex:  integrityMutex,
	}
}

// StartDataIntegrityEnforcement starts a continuous loop for data integrity checks and enforcement
func (automation *DataIntegrityEnforcementAutomation) StartDataIntegrityEnforcement() {
	ticker := time.NewTicker(DataIntegrityCheckInterval)

	go func() {
		for range ticker.C {
			automation.verifyDataIntegrityAcrossNetwork()
		}
	}()
}

// verifyDataIntegrityAcrossNetwork checks the integrity of the blockchain data and triggers actions if inconsistencies are found
func (automation *DataIntegrityEnforcementAutomation) verifyDataIntegrityAcrossNetwork() {
	automation.integrityMutex.Lock()
	defer automation.integrityMutex.Unlock()

	// Get the current block and verify its integrity
	currentBlock, err := automation.consensusEngine.GetLatestBlock()
	if err != nil {
		fmt.Println("Error fetching latest block:", err)
		return
	}

	// Perform integrity checks for each sub-block in the block
	for _, subBlock := range currentBlock.SubBlocks {
		if !automation.verifySubBlockIntegrity(subBlock) {
			automation.handleIntegrityViolation(subBlock)
		}
	}

	// Log successful integrity check
	automation.logIntegrityCheckSuccess(currentBlock)
}

// verifySubBlockIntegrity validates a sub-block's data against the consensus rules and ensures data consistency
func (automation *DataIntegrityEnforcementAutomation) verifySubBlockIntegrity(subBlock common.SubBlock) bool {
	// Hash the sub-block data
	hashedData := automation.hashData(subBlock.Data)

	// Verify the hash matches the expected value from the block's header
	return hashedData == subBlock.ExpectedHash
}

// handleIntegrityViolation takes action when a sub-block's data integrity is compromised
func (automation *DataIntegrityEnforcementAutomation) handleIntegrityViolation(subBlock common.SubBlock) {
	// Trigger a re-verification of the entire block through Synnergy Consensus
	err := automation.consensusEngine.RevalidateSubBlock(subBlock)
	if err != nil {
		fmt.Println("Error revalidating sub-block:", err)
		return
	}

	// Log the integrity violation and action taken
	automation.logIntegrityViolation(subBlock)
}

// logIntegrityViolation logs data integrity violations in the ledger
func (automation *DataIntegrityEnforcementAutomation) logIntegrityViolation(subBlock common.SubBlock) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("integrity-violation-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Integrity Violation",
		Status:    "Action Taken",
		Details:   fmt.Sprintf("Sub-block ID %s violated data integrity rules.", subBlock.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log data integrity violation:", err)
	} else {
		fmt.Println("Data integrity violation logged in the ledger.")
	}
}

// logIntegrityCheckSuccess logs the successful verification of the latest block in the ledger
func (automation *DataIntegrityEnforcementAutomation) logIntegrityCheckSuccess(block common.Block) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("integrity-check-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Data Integrity Check",
		Status:    "Success",
		Details:   fmt.Sprintf("Block ID %s passed data integrity verification.", block.ID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log successful integrity check:", err)
	} else {
		fmt.Println("Integrity check successfully logged in the ledger.")
	}
}

// hashData hashes the provided data using SHA-256
func (automation *DataIntegrityEnforcementAutomation) hashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}

// encryptData encrypts sensitive data before logging it in the ledger
func (automation *DataIntegrityEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
