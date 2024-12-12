package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network"
	"synnergy_network_demo/common"
	"synnergy_network_demo/synnergy_vm"
)

// WatchtowerNode represents a node responsible for overseeing transactions and ensuring compliance with smart contracts, particularly for off-chain transactions (such as in the Lightning Network).
type WatchtowerNode struct {
	NodeID            string                        // Unique identifier for the watchtower node
	MonitoredChannels map[string]*common.Channel    // Monitored off-chain channels (e.g., Lightning Network)
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validation and conflict resolution
	EncryptionService *encryption.Encryption        // Encryption service for securing data transmission
	NetworkManager    *network.NetworkManager       // Manages network communications
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	SNVM              *synnergy_vm.VirtualMachine   // Virtual Machine instance for enforcing smart contract rules
	SyncInterval      time.Duration                 // Interval for syncing monitored data with the network
	Logs              map[string]string             // Logs of transaction activities for audits and compliance
}

// NewWatchtowerNode initializes a new WatchtowerNode in the blockchain network.
func NewWatchtowerNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, snvm *synnergy_vm.VirtualMachine) *WatchtowerNode {
	return &WatchtowerNode{
		NodeID:            nodeID,
		MonitoredChannels: make(map[string]*common.Channel),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SNVM:              snvm,
		SyncInterval:      syncInterval,
		Logs:              make(map[string]string),
	}
}

// StartNode starts the watchtower node's monitoring activities.
func (wn *WatchtowerNode) StartNode() error {
	wn.mutex.Lock()
	defer wn.mutex.Unlock()

	// Start syncing and monitoring off-chain channels and transactions.
	go wn.syncWithNetwork()
	go wn.monitorTransactions()

	fmt.Printf("Watchtower node %s started successfully.\n", wn.NodeID)
	return nil
}

// syncWithNetwork syncs monitored channels and data with the network at regular intervals.
func (wn *WatchtowerNode) syncWithNetwork() {
	ticker := time.NewTicker(wn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		wn.mutex.Lock()
		// Fetch any updates to monitored channels and synchronize data with the network.
		otherWatchtowerNodes := wn.NetworkManager.DiscoverOtherWatchtowerNodes(wn.NodeID)
		for _, node := range otherWatchtowerNodes {
			wn.syncMonitoredChannels(node)
		}
		wn.mutex.Unlock()
	}
}

// syncMonitoredChannels synchronizes the channel states from other Watchtower nodes.
func (wn *WatchtowerNode) syncMonitoredChannels(peerNode string) {
	peerChannels, err := wn.NetworkManager.RequestChannelUpdates(peerNode)
	if err != nil {
		fmt.Printf("Failed to sync channels from node %s: %v\n", peerNode, err)
		return
	}

	// Validate and update the monitored channels with the new state information.
	for _, channel := range peerChannels {
		if wn.ConsensusEngine.ValidateChannelState(channel) {
			wn.MonitoredChannels[channel.ChannelID] = channel
			fmt.Printf("Channel %s synced successfully from node %s.\n", channel.ChannelID, peerNode)
		} else {
			fmt.Printf("Channel %s from node %s failed validation.\n", channel.ChannelID, peerNode)
		}
	}
}

// monitorTransactions monitors all channels for compliance with smart contract terms.
func (wn *WatchtowerNode) monitorTransactions() {
	for {
		// Listen for updates or activity on monitored channels.
		activity, err := wn.NetworkManager.ReceiveChannelActivity()
		if err != nil {
			fmt.Printf("Error receiving channel activity: %v\n", err)
			continue
		}

		// Validate and process the activity.
		err = wn.processChannelActivity(activity)
		if err != nil {
			fmt.Printf("Channel activity processing failed: %v\n", err)
		}
	}
}

// processChannelActivity validates and processes off-chain channel activities to ensure compliance.
func (wn *WatchtowerNode) processChannelActivity(activity *common.ChannelActivity) error {
	wn.mutex.Lock()
	defer wn.mutex.Unlock()

	// Validate the channel activity with the consensus engine.
	if valid, err := wn.ConsensusEngine.ValidateChannelActivity(activity); err != nil || !valid {
		return fmt.Errorf("invalid channel activity: %v", err)
	}

	// Log the valid activity for future audits.
	wn.Logs[activity.ChannelID] = fmt.Sprintf("Valid activity recorded on channel %s at %s.", activity.ChannelID, time.Now().String())
	fmt.Printf("Channel activity on %s processed successfully.\n", activity.ChannelID)
	return nil
}

// handleContractBreach detects breaches of smart contracts and resolves conflicts automatically.
func (wn *WatchtowerNode) handleContractBreach(channel *common.Channel, activity *common.ChannelActivity) error {
	wn.mutex.Lock()
	defer wn.mutex.Unlock()

	// Detect contract breaches.
	if wn.ConsensusEngine.DetectContractBreach(channel, activity) {
		// Resolve conflict and update the channel state.
		err := wn.resolveConflict(channel, activity)
		if err != nil {
			return fmt.Errorf("failed to resolve contract breach: %v", err)
		}
		// Log the breach and resolution.
		wn.Logs[channel.ChannelID] = fmt.Sprintf("Breach detected and resolved for channel %s at %s.", channel.ChannelID, time.Now().String())
	}
	return nil
}

// resolveConflict resolves conflicts when contract terms are violated in a channel.
func (wn *WatchtowerNode) resolveConflict(channel *common.Channel, activity *common.ChannelActivity) error {
	// Use the Synnergy Consensus to resolve the contract breach.
	err := wn.ConsensusEngine.ResolveChannelConflict(channel, activity)
	if err != nil {
		return fmt.Errorf("failed to resolve conflict in channel %s: %v", channel.ChannelID, err)
	}

	fmt.Printf("Conflict resolved for channel %s.\n", channel.ChannelID)
	return nil
}

// EncryptData encrypts sensitive data before transmitting to other nodes.
func (wn *WatchtowerNode) EncryptData(data []byte) ([]byte, error) {
	encryptedData, err := wn.EncryptionService.EncryptData(data, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// DecryptData decrypts incoming encrypted data from other nodes.
func (wn *WatchtowerNode) DecryptData(encryptedData []byte) ([]byte, error) {
	decryptedData, err := wn.EncryptionService.DecryptData(encryptedData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}
