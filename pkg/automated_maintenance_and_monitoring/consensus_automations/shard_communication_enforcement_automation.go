package consensus_automations

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
	ShardCommunicationCheckInterval = 3 * time.Second  // Interval for checking shard communication
	CommunicationFailureThreshold   = 0.8              // Threshold for communication failures triggering enforcement
	ShardCommunicationKey           = "shard_comm_log_key" // Encryption key for shard communication logs
)

// ShardCommunicationEnforcementAutomation ensures shard communication in Synnergy Consensus (PoH, PoS, and PoW)
type ShardCommunicationEnforcementAutomation struct {
	ledgerInstance  *ledger.Ledger                    // Blockchain ledger for tracking shard communication events
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine to monitor and enforce communication
	stateMutex      *sync.RWMutex                     // Mutex for thread-safe ledger access
}

// NewShardCommunicationEnforcementAutomation initializes the shard communication enforcement automation
func NewShardCommunicationEnforcementAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *ShardCommunicationEnforcementAutomation {
	return &ShardCommunicationEnforcementAutomation{
		ledgerInstance:  ledgerInstance,
		consensusEngine: consensusEngine,
		stateMutex:      stateMutex,
	}
}

// StartShardCommunicationMonitoring starts the continuous monitoring of shard communication in Synnergy Consensus
func (automation *ShardCommunicationEnforcementAutomation) StartShardCommunicationMonitoring() {
	ticker := time.NewTicker(ShardCommunicationCheckInterval)

	go func() {
		for range ticker.C {
			fmt.Println("Checking shard communication across PoH, PoS, and PoW...")
			automation.monitorShardCommunication()
		}
	}()
}

// monitorShardCommunication checks the communication between shards in PoH, PoS, and PoW stages
func (automation *ShardCommunicationEnforcementAutomation) monitorShardCommunication() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Check shard communication across PoH
	pohCommHealth := automation.consensusEngine.CheckPoHShardCommunication()

	// Check shard communication across PoS
	posCommHealth := automation.consensusEngine.CheckPoSShardCommunication()

	// Check shard communication across PoW
	powCommHealth := automation.consensusEngine.CheckPoWShardCommunication()

	// If communication failure threshold is exceeded, trigger enforcement
	if pohCommHealth < CommunicationFailureThreshold || posCommHealth < CommunicationFailureThreshold || powCommHealth < CommunicationFailureThreshold {
		fmt.Println("Shard communication failure detected, initiating enforcement...")
		automation.enforceShardCommunication(pohCommHealth, posCommHealth, powCommHealth)
	} else {
		fmt.Println("Shard communication is healthy.")
	}
}

// enforceShardCommunication ensures shards communicate correctly and enforces corrections where necessary
func (automation *ShardCommunicationEnforcementAutomation) enforceShardCommunication(pohCommHealth, posCommHealth, powCommHealth float64) {
	fmt.Printf("Enforcing shard communication across consensus layers. PoH: %.2f, PoS: %.2f, PoW: %.2f\n", pohCommHealth, posCommHealth, powCommHealth)

	// Attempt to correct shard communication issues by invoking consensus engine functions
	automation.consensusEngine.CorrectPoHShardCommunication()
	automation.consensusEngine.CorrectPoSShardCommunication()
	automation.consensusEngine.CorrectPoWShardCommunication()

	// Log the enforcement action in the ledger
	automation.logShardCommunicationEnforcement(pohCommHealth, posCommHealth, powCommHealth)
}

// logShardCommunicationEnforcement logs the shard communication enforcement event in the ledger
func (automation *ShardCommunicationEnforcementAutomation) logShardCommunicationEnforcement(pohCommHealth, posCommHealth, powCommHealth float64) {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	enforcementLog := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-comm-enforcement-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Communication Enforcement",
		Status:    "Enforced",
		Details:   fmt.Sprintf("PoH: %.2f, PoS: %.2f, PoW: %.2f", pohCommHealth, posCommHealth, powCommHealth),
	}

	// Encrypt the ledger entry for security purposes
	encryptedEntry, err := encryption.EncryptLedgerEntry(enforcementLog, []byte(ShardCommunicationKey))
	if err != nil {
		fmt.Printf("Error encrypting shard communication enforcement log: %v\n", err)
		return
	}

	automation.ledgerInstance.AddEntry(encryptedEntry)
	fmt.Println("Shard communication enforcement log stored in the ledger.")
}

// Additional helper function to ensure overall shard communication integrity post-enforcement
func (automation *ShardCommunicationEnforcementAutomation) ensureShardCommunicationIntegrity() {
	fmt.Println("Ensuring shard communication integrity post-enforcement...")

	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Validate overall chain and shard communication consistency
	err := automation.consensusEngine.ValidateShardCommunication()
	if err != nil {
		fmt.Printf("Shard communication validation failed: %v\n", err)
		automation.enforceShardCommunication(0.0, 0.0, 0.0) // Re-trigger enforcement if integrity is not restored
	} else {
		fmt.Println("Shard communication is consistent and fully integrated.")
	}
}
