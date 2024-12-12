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

// Interval for checking sidechain interoperability restrictions
const (
	SidechainInteroperabilityCheckInterval = 10 * time.Second // Adjust this as needed
	SidechainErrorThreshold                = 5                // Maximum allowed sidechain communication errors
)

// SidechainInteroperabilityRestrictionAutomation monitors and restricts sidechain communication if errors or issues are detected
type SidechainInteroperabilityRestrictionAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	sidechainErrorCount   map[string]int // Tracks errors between mainchain and each sidechain
	enabledSidechainLinks map[string]bool // Tracks active sidechain interoperability links
}

// NewSidechainInteroperabilityRestrictionAutomation initializes the automation
func NewSidechainInteroperabilityRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SidechainInteroperabilityRestrictionAutomation {
	return &SidechainInteroperabilityRestrictionAutomation{
		consensusSystem:       consensusSystem,
		ledgerInstance:        ledgerInstance,
		stateMutex:            stateMutex,
		sidechainErrorCount:   make(map[string]int),
		enabledSidechainLinks: make(map[string]bool),
	}
}

// StartMonitoringSidechainInteroperability continuously monitors sidechain interoperability
func (automation *SidechainInteroperabilityRestrictionAutomation) StartMonitoringSidechainInteroperability() {
	ticker := time.NewTicker(SidechainInteroperabilityCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateSidechainConnections()
		}
	}()
}

// evaluateSidechainConnections checks if sidechain communications are operating properly and applies restrictions if necessary
func (automation *SidechainInteroperabilityRestrictionAutomation) evaluateSidechainConnections() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	sidechainLinks := automation.consensusSystem.GetSidechainInteroperabilityStatus()

	for sidechainID, isActive := range sidechainLinks {
		if !isActive {
			automation.sidechainErrorCount[sidechainID]++
			automation.logSidechainError(sidechainID)

			// If sidechain errors exceed the threshold, restrict the link
			if automation.sidechainErrorCount[sidechainID] >= SidechainErrorThreshold {
				automation.restrictSidechainInteroperability(sidechainID)
			}
		} else {
			// Reset error count if sidechain communication is successful
			automation.sidechainErrorCount[sidechainID] = 0
		}
	}
}

// logSidechainError logs sidechain communication errors into the ledger
func (automation *SidechainInteroperabilityRestrictionAutomation) logSidechainError(sidechainID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("sidechain-error-%s-%d", sidechainID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Sidechain Interoperability Error",
		Status:    "Error",
		Details:   fmt.Sprintf("Communication error detected with sidechain %s.", sidechainID),
	}

	// Encrypt the sidechain error details before logging them in the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	// Add the error entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log sidechain error:", err)
	} else {
		fmt.Println("Sidechain communication error logged for:", sidechainID)
	}
}

// restrictSidechainInteroperability restricts further communication with a sidechain if error count exceeds the threshold
func (automation *SidechainInteroperabilityRestrictionAutomation) restrictSidechainInteroperability(sidechainID string) {
	fmt.Printf("Sidechain %s exceeded the error threshold. Communication restricted.\n", sidechainID)

	// Log the restriction event in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("sidechain-restriction-%s-%d", sidechainID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Sidechain Interoperability Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Sidechain %s communication restricted due to exceeding %d errors.", sidechainID, SidechainErrorThreshold),
	}

	// Encrypt restriction details before adding them to the ledger
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log sidechain restriction:", err)
	} else {
		fmt.Println("Sidechain interoperability restricted for:", sidechainID)
	}

	// Inform the consensus system to restrict transactions or communication with the affected sidechain
	automation.consensusSystem.RestrictSidechainInteroperability(sidechainID)
}

// encryptData encrypts sensitive sidechain error/restriction data before storing it in the ledger
func (automation *SidechainInteroperabilityRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting sidechain error details:", err)
		return data
	}
	return string(encryptedData)
}
