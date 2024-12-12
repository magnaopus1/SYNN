package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
)

// Configuration for historical record enforcement automation
const (
	HistoricalRecordCheckInterval = 60 * time.Second // Interval to check for historical record compliance
)

// HistoricalRecordEnforcementAutomation ensures that all network transactions, sub-blocks, and blocks are archived for historical integrity
type HistoricalRecordEnforcementAutomation struct {
	consensusEngine     *consensus.SynnergyConsensus
	ledgerInstance      *ledger.Ledger
	enforcementMutex    *sync.RWMutex
	historicalArchive   map[string]bool // Tracks archived records by ID
}

// NewHistoricalRecordEnforcementAutomation initializes the historical record enforcement automation
func NewHistoricalRecordEnforcementAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *HistoricalRecordEnforcementAutomation {
	return &HistoricalRecordEnforcementAutomation{
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		historicalArchive:  make(map[string]bool),
	}
}

// StartHistoricalRecordEnforcement begins continuous monitoring and enforcement of historical record archiving
func (automation *HistoricalRecordEnforcementAutomation) StartHistoricalRecordEnforcement() {
	ticker := time.NewTicker(HistoricalRecordCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkHistoricalRecordCompliance()
		}
	}()
}

// checkHistoricalRecordCompliance verifies that all validated sub-blocks and blocks are archived for historical record
func (automation *HistoricalRecordEnforcementAutomation) checkHistoricalRecordCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.archiveSubBlockRecords()
	automation.archiveBlockRecords()
}

// archiveSubBlockRecords archives records of all validated sub-blocks
func (automation *HistoricalRecordEnforcementAutomation) archiveSubBlockRecords() {
	for _, subBlockID := range automation.consensusEngine.GetValidatedSubBlocks() {
		if !automation.historicalArchive[subBlockID] {
			automation.archiveRecord(subBlockID, "Sub-Block")
		}
	}
}

// archiveBlockRecords archives records of all validated blocks
func (automation *HistoricalRecordEnforcementAutomation) archiveBlockRecords() {
	for _, blockID := range automation.consensusEngine.GetValidatedBlocks() {
		if !automation.historicalArchive[blockID] {
			automation.archiveRecord(blockID, "Block")
		}
	}
}

// archiveRecord archives a validated record (sub-block or block) in the ledger
func (automation *HistoricalRecordEnforcementAutomation) archiveRecord(recordID, recordType string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("historical-record-%s-%d", recordID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Historical Record",
		Status:    "Archived",
		Details:   automation.encryptData(fmt.Sprintf("Archived %s: %s", recordType, recordID)),
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to archive historical record for %s %s: %v\n", recordType, recordID, err)
		automation.logArchiveAction(recordID, recordType, "Archive Failed")
	} else {
		fmt.Printf("%s %s archived successfully.\n", recordType, recordID)
		automation.historicalArchive[recordID] = true
		automation.logArchiveAction(recordID, recordType, "Archived")
	}
}

// logArchiveAction securely logs actions related to historical record archiving
func (automation *HistoricalRecordEnforcementAutomation) logArchiveAction(recordID, recordType, action string) {
	entryDetails := fmt.Sprintf("Action: %s, Record Type: %s, Record ID: %s", action, recordType, recordID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("historical-archive-log-%s-%d", recordID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Archive Log",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log archive action for %s %s: %v\n", recordType, recordID, err)
	} else {
		fmt.Println("Archive action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *HistoricalRecordEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualArchive allows administrators to manually archive a specific sub-block or block
func (automation *HistoricalRecordEnforcementAutomation) TriggerManualArchive(recordID, recordType string) {
	fmt.Printf("Manually archiving %s: %s\n", recordType, recordID)

	if automation.historicalArchive[recordID] {
		fmt.Printf("%s %s is already archived.\n", recordType, recordID)
		automation.logArchiveAction(recordID, recordType, "Manual Archive Attempt Skipped - Already Archived")
	} else {
		automation.archiveRecord(recordID, recordType)
	}
}
