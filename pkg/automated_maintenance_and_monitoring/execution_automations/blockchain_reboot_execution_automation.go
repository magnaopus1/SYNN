package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network"
)

const (
	RebootCheckInterval        = 15 * time.Minute  // Interval for checking if a blockchain reboot is required
	NetworkStabilityThreshold  = 0.7               // Threshold for determining network instability
	RebootTriggerFailureLimit  = 3                 // Number of consecutive failures before triggering reboot
	RebootTimeout              = 10 * time.Minute  // Maximum allowed time for the reboot process
)

// BlockchainRebootAutomation manages the automated reboot of the blockchain system
type BlockchainRebootAutomation struct {
	consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
	ledgerInstance    *ledger.Ledger                        // Ledger for recording reboot actions
	networkManager    *network.Manager                      // Network manager for monitoring and triggering reboots
	stateMutex        *sync.RWMutex                         // Mutex for thread-safe operations
	rebootFailures    int                                   // Count of consecutive reboot failures
}

// NewBlockchainRebootAutomation initializes the blockchain reboot automation
func NewBlockchainRebootAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, networkManager *network.Manager, stateMutex *sync.RWMutex) *BlockchainRebootAutomation {
	return &BlockchainRebootAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		networkManager:   networkManager,
		stateMutex:       stateMutex,
		rebootFailures:   0,
	}
}

// StartRebootMonitor starts the monitoring for potential blockchain reboots based on network stability
func (automation *BlockchainRebootAutomation) StartRebootMonitor() {
	ticker := time.NewTicker(RebootCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkForRebootTrigger()
		}
	}()
}

// checkForRebootTrigger checks network conditions and determines if a blockchain reboot is required
func (automation *BlockchainRebootAutomation) checkForRebootTrigger() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	networkStability := automation.networkManager.CheckNetworkStability()
	fmt.Printf("Network stability: %.2f\n", networkStability)

	if networkStability < NetworkStabilityThreshold {
		fmt.Println("Network instability detected. Triggering blockchain reboot.")
		automation.triggerReboot()
	} else {
		fmt.Println("Network stability is sufficient. No reboot required.")
		automation.rebootFailures = 0 // Reset failure count if network is stable
	}
}

// triggerReboot initiates the blockchain reboot process and logs the result
func (automation *BlockchainRebootAutomation) triggerReboot() {
	err := automation.networkManager.RebootBlockchain(RebootTimeout)
	if err != nil {
		automation.rebootFailures++
		fmt.Printf("Blockchain reboot failed. Failure count: %d\n", automation.rebootFailures)
		automation.logRebootAction("Failed", "Blockchain reboot failed.")
		
		if automation.rebootFailures >= RebootTriggerFailureLimit {
			fmt.Println("Multiple reboot failures detected. Further action may be required.")
			// Optionally trigger additional alerts or escalations
		}
	} else {
		fmt.Println("Blockchain reboot successful.")
		automation.rebootFailures = 0 // Reset failure count on success
		automation.logRebootAction("Success", "Blockchain reboot executed successfully.")
	}
}

// logRebootAction logs the result of the blockchain reboot attempt into the ledger
func (automation *BlockchainRebootAutomation) logRebootAction(status, details string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("blockchain-reboot-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Blockchain Reboot",
		Status:    status,
		Details:   details,
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	if err := automation.ledgerInstance.AddEntry(entry); err != nil {
		fmt.Printf("Error logging reboot action: %v\n", err)
	} else {
		fmt.Println("Reboot action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *BlockchainRebootAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
