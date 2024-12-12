package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/hotfix"
)

const (
	HotfixCheckInterval       = 5 * time.Minute // Interval for checking for new hotfixes
	MaxHotfixRetryAttempts    = 3               // Maximum number of attempts to apply a hotfix
	HotfixLedgerEntryTemplate = "hotfix-application-%s" // Template for logging hotfix applications
)

// HotfixExecutionAutomation manages the execution and application of blockchain hotfixes
type HotfixExecutionAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance    *ledger.Ledger                        // Ledger instance for logging hotfix applications
	hotfixManager     *hotfix.Manager                       // Hotfix manager responsible for hotfixes
	hotfixMutex       *sync.RWMutex                         // Mutex for thread-safe hotfix operations
}

// NewHotfixExecutionAutomation initializes a new HotfixExecutionAutomation instance
func NewHotfixExecutionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, hotfixManager *hotfix.Manager, hotfixMutex *sync.RWMutex) *HotfixExecutionAutomation {
	return &HotfixExecutionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		hotfixManager:   hotfixManager,
		hotfixMutex:     hotfixMutex,
	}
}

// StartHotfixMonitoring starts the continuous hotfix monitoring process
func (automation *HotfixExecutionAutomation) StartHotfixMonitoring() {
	ticker := time.NewTicker(HotfixCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkForHotfixes()
		}
	}()
}

// checkForHotfixes checks if there are new hotfixes to apply and processes them
func (automation *HotfixExecutionAutomation) checkForHotfixes() {
	automation.hotfixMutex.Lock()
	defer automation.hotfixMutex.Unlock()

	pendingHotfixes := automation.hotfixManager.GetPendingHotfixes()

	for _, hotfix := range pendingHotfixes {
		automation.applyHotfix(hotfix)
	}
}

// applyHotfix applies a hotfix and logs the process in the ledger
func (automation *HotfixExecutionAutomation) applyHotfix(hotfix *hotfix.Hotfix) {
	fmt.Printf("Applying hotfix %s...\n", hotfix.ID)

	for i := 0; i < MaxHotfixRetryAttempts; i++ {
		err := automation.hotfixManager.ApplyHotfix(hotfix)
		if err != nil {
			fmt.Printf("Failed to apply hotfix %s, attempt %d: %v\n", hotfix.ID, i+1, err)
			continue
		}

		// Log successful hotfix application in the ledger
		automation.logHotfixInLedger(hotfix)
		fmt.Printf("Hotfix %s successfully applied.\n", hotfix.ID)
		return
	}

	fmt.Printf("Failed to apply hotfix %s after %d attempts.\n", hotfix.ID, MaxHotfixRetryAttempts)
}

// logHotfixInLedger securely logs the successful hotfix application into the ledger
func (automation *HotfixExecutionAutomation) logHotfixInLedger(hotfix *hotfix.Hotfix) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf(HotfixLedgerEntryTemplate, hotfix.ID),
		Timestamp: time.Now().Unix(),
		Type:      "Hotfix Application",
		Status:    "Success",
		Details:   fmt.Sprintf("Hotfix %s was applied to the system.", hotfix.ID),
	}

	// Encrypt details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log hotfix %s in the ledger: %v\n", hotfix.ID, err)
	} else {
		fmt.Printf("Hotfix %s successfully logged in the ledger.\n", hotfix.ID)
	}
}

// encryptData encrypts sensitive details before logging in the ledger
func (automation *HotfixExecutionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualHotfixApplication allows administrators to manually apply a hotfix
func (automation *HotfixExecutionAutomation) TriggerManualHotfixApplication(hotfixID string) {
	fmt.Printf("Manually triggering application of hotfix %s...\n", hotfixID)

	hotfix := automation.hotfixManager.GetHotfixByID(hotfixID)
	if hotfix != nil {
		automation.applyHotfix(hotfix)
	} else {
		fmt.Printf("Hotfix %s not found.\n", hotfixID)
	}
}
