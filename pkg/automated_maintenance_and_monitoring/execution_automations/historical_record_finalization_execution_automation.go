package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	HistoricalRecordCheckInterval = 10 * time.Minute // Interval for checking and finalizing historical records
	RecordFinalizationThreshold   = 100              // Minimum number of transactions or sub-blocks for finalization
)

// HistoricalRecordFinalizationAutomation manages the finalization of historical records
type HistoricalRecordFinalizationAutomation struct {
	consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for record validation
	ledgerInstance   *ledger.Ledger                        // Ledger to track finalized historical records
	recordMutex      *sync.RWMutex                         // Mutex for thread-safe operations on historical records
}

// NewHistoricalRecordFinalizationAutomation initializes the historical record finalization automation
func NewHistoricalRecordFinalizationAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, recordMutex *sync.RWMutex) *HistoricalRecordFinalizationAutomation {
	return &HistoricalRecordFinalizationAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		recordMutex:     recordMutex,
	}
}

// StartHistoricalRecordFinalizationMonitoring begins the monitoring of historical records to finalize
func (automation *HistoricalRecordFinalizationAutomation) StartHistoricalRecordFinalizationMonitoring() {
	ticker := time.NewTicker(HistoricalRecordCheckInterval)
	go func() {
		for range ticker.C {
			automation.finalizeHistoricalRecords()
		}
	}()
}

// finalizeHistoricalRecords processes historical records and finalizes them based on conditions
func (automation *HistoricalRecordFinalizationAutomation) finalizeHistoricalRecords() {
	automation.recordMutex.Lock()
	defer automation.recordMutex.Unlock()

	historicalRecords := automation.consensusEngine.GetPendingHistoricalRecords()

	for _, record := range historicalRecords {
		if record.TransactionCount() >= RecordFinalizationThreshold {
			automation.finalizeRecord(record.ID)
		}
	}
}

// finalizeRecord finalizes a historical record and logs it in the ledger
func (automation *HistoricalRecordFinalizationAutomation) finalizeRecord(recordID string) {
	fmt.Printf("Finalizing historical record %s...\n", recordID)

	err := automation.consensusEngine.FinalizeRecord(recordID)
	if err != nil {
		fmt.Printf("Failed to finalize historical record %s: %v\n", recordID, err)
		return
	}

	// Log the finalized record in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("historical-record-%s-%d", recordID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Historical Record Finalization",
		Status:    "Success",
		Details:   fmt.Sprintf("Historical record %s successfully finalized.", recordID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log historical record finalization for record %s: %v\n", recordID, err)
	} else {
		fmt.Printf("Historical record %s finalization successfully logged in the ledger.\n", recordID)
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *HistoricalRecordFinalizationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualRecordFinalization allows administrators to manually finalize a historical record
func (automation *HistoricalRecordFinalizationAutomation) TriggerManualRecordFinalization(recordID string) {
	fmt.Printf("Manually triggering finalization for historical record %s...\n", recordID)
	automation.finalizeRecord(recordID)
}

