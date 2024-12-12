package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	SyncProtocolInterval       = 2 * time.Minute  // Interval for checking and syncing nodes
	SyncProtocolLogEncryptionKey = "syncProtocolKey" // Encryption key for sync protocol logs
	MaxSyncBatch               = 50               // Maximum nodes to sync in a single batch
)

// SyncProtocolAutomation handles the continuous syncing of nodes within the network.
type SyncProtocolAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging sync events
	consensusSystem *consensus.SynnergyConsensus // Reference to the Synnergy Consensus system
	syncMutex       *sync.RWMutex                // Mutex for managing concurrent syncs
}

// NewSyncProtocolAutomation creates a new instance of SyncProtocolAutomation.
func NewSyncProtocolAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, syncMutex *sync.RWMutex) *SyncProtocolAutomation {
	return &SyncProtocolAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		syncMutex:       syncMutex,
	}
}

// StartSyncProtocolAutomation begins the automated node syncing process at regular intervals.
func (automation *SyncProtocolAutomation) StartSyncProtocolAutomation() {
	ticker := time.NewTicker(SyncProtocolInterval)
	go func() {
		for range ticker.C {
			automation.syncNodes()
		}
	}()
}

// syncNodes retrieves nodes that require synchronization and syncs them in batches.
func (automation *SyncProtocolAutomation) syncNodes() {
	automation.syncMutex.Lock()
	defer automation.syncMutex.Unlock()

	// Fetch a list of nodes needing sync from the consensus system
	nodesToSync, err := automation.consensusSystem.GetNodesNeedingSync(MaxSyncBatch)
	if err != nil {
		fmt.Printf("Error fetching nodes for sync: %v\n", err)
		return
	}

	// Sync each node and log the results
	for _, node := range nodesToSync {
		err := automation.syncNode(node)
		if err != nil {
			fmt.Printf("Error syncing node %s: %v\n", node.ID, err)
			continue
		}
		automation.logSyncEvent(node)
	}
}

// syncNode performs the actual synchronization of a node with the latest consensus and ledger data.
func (automation *SyncProtocolAutomation) syncNode(node common.Node) error {
	// Fetch the latest sub-blocks and state for synchronization
	subBlocks, err := automation.consensusSystem.GetSubBlocksForSync(node.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch sub-blocks for node %s: %w", node.ID, err)
	}

	// Apply the sub-blocks to the node's local state
	err = automation.consensusSystem.ApplySubBlocks(node.ID, subBlocks)
	if err != nil {
		return fmt.Errorf("failed to apply sub-blocks for node %s: %w", node.ID, err)
	}

	// Finalize the node sync by updating its block state
	err = automation.consensusSystem.UpdateNodeBlockState(node.ID)
	if err != nil {
		return fmt.Errorf("failed to update block state for node %s: %w", node.ID, err)
	}

	return nil
}

// logSyncEvent securely logs the synchronization event in the ledger.
func (automation *SyncProtocolAutomation) logSyncEvent(node common.Node) {
	// Encrypt the sync event details for secure logging
	encryptedSyncDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Node %s synchronized successfully", node.ID)), SyncProtocolLogEncryptionKey)
	if err != nil {
		fmt.Printf("Error encrypting sync log for node %s: %v\n", node.ID, err)
		return
	}

	// Create a ledger entry for the sync event
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("NODE-SYNC-%s-%d", node.ID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Node Sync",
		Status:    "Success",
		Details:   string(encryptedSyncDetails),
	}

	// Log the sync event in the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log sync event for node %s: %v\n", node.ID, err)
	}
}
