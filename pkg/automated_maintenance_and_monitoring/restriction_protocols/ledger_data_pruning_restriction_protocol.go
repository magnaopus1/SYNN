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
	DataPruningCheckInterval = 6 * time.Hour // Interval for checking ledger data pruning needs
	MaxAllowedLedgerSize     = 10 * 1024 * 1024 * 1024 // Maximum allowed size of the ledger before pruning (10 GB)
	PruningBatchSize         = 5000           // Number of ledger entries pruned in a single batch
)

// LedgerDataPruningRestrictionAutomation manages and restricts the pruning of old ledger data based on size and conditions
type LedgerDataPruningRestrictionAutomation struct {
	consensusSystem   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	stateMutex        *sync.RWMutex
}

// NewLedgerDataPruningRestrictionAutomation initializes and returns an instance of LedgerDataPruningRestrictionAutomation
func NewLedgerDataPruningRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *LedgerDataPruningRestrictionAutomation {
	return &LedgerDataPruningRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
	}
}

// StartPruningMonitoring starts continuous monitoring for data pruning conditions
func (automation *LedgerDataPruningRestrictionAutomation) StartPruningMonitoring() {
	ticker := time.NewTicker(DataPruningCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorLedgerDataPruning()
		}
	}()
}

// monitorLedgerDataPruning checks if the ledger has exceeded the maximum allowed size and enforces pruning
func (automation *LedgerDataPruningRestrictionAutomation) monitorLedgerDataPruning() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch the current size of the ledger
	currentLedgerSize := automation.ledgerInstance.GetLedgerSize()

	if currentLedgerSize > MaxAllowedLedgerSize {
		// Trigger pruning process if size exceeds the limit
		automation.pruneLedgerData()
	} else {
		fmt.Printf("Ledger size within limit: %.2f GB\n", float64(currentLedgerSize)/(1024*1024*1024))
	}
}

// pruneLedgerData performs the pruning of old ledger entries
func (automation *LedgerDataPruningRestrictionAutomation) pruneLedgerData() {
	// Fetch the old entries to be pruned based on the pruning batch size
	oldEntries := automation.ledgerInstance.GetOldestEntries(PruningBatchSize)

	// Prune the selected entries
	err := automation.ledgerInstance.PruneEntries(oldEntries)
	if err != nil {
		fmt.Println("Error during ledger data pruning:", err)
		automation.logPruningFailure(err)
		return
	}

	// Log the successful pruning operation
	automation.logPruningSuccess(len(oldEntries))
}

// logPruningSuccess logs a successful ledger data pruning operation
func (automation *LedgerDataPruningRestrictionAutomation) logPruningSuccess(entriesPruned int) {
	// Create a ledger entry for pruning success
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ledger-pruning-success-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Ledger Pruning",
		Status:    "Completed",
		Details:   fmt.Sprintf("Successfully pruned %d ledger entries.", entriesPruned),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptPruningData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log ledger pruning success:", err)
	} else {
		fmt.Println("Ledger pruning success logged.")
	}
}

// logPruningFailure logs a failure during the pruning process
func (automation *LedgerDataPruningRestrictionAutomation) logPruningFailure(err error) {
	// Create a ledger entry for pruning failure
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ledger-pruning-failure-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Ledger Pruning Failure",
		Status:    "Failed",
		Details:   fmt.Sprintf("Ledger pruning failed: %v", err),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptPruningData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	errLog := automation.ledgerInstance.AddEntry(entry)
	if errLog != nil {
		fmt.Println("Failed to log ledger pruning failure:", errLog)
	} else {
		fmt.Println("Ledger pruning failure logged.")
	}
}

// encryptPruningData encrypts the pruning data before logging for security
func (automation *LedgerDataPruningRestrictionAutomation) encryptPruningData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting pruning data:", err)
		return data
	}
	return string(encryptedData)
}
