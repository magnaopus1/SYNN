package node_type

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"           // Blockchain ledger-related components
	"synnergy_network_demo/encryption"       // Encryption service for secure data handling
	"synnergy_network_demo/synnergy_consensus" // Synnergy Consensus engine
	"synnergy_network_demo/network"          // Network and communication management
	"synnergy_network_demo/common"           // Common utilities for random generation, keys, etc.
)

// LightningNode represents a node responsible for handling off-chain transactions using payment channels in the Lightning Network.
type LightningNode struct {
	NodeID            string                        // Unique identifier for the node
	PaymentChannels   map[string]*PaymentChannel     // Active payment channels managed by the Lightning Node
	ConsensusEngine   *synnergy_consensus.Engine     // Consensus engine for validating transactions and state changes
	EncryptionService *encryption.Encryption         // Encryption service for securing channel communications and transactions
	NetworkManager    *network.NetworkManager        // Network manager for communicating with other nodes
	LiquidityPool     map[string]float64             // Liquidity for each channel
	mutex             sync.Mutex                     // Mutex for thread-safe operations
	SyncInterval      time.Duration                  // Interval for syncing channels with the blockchain
	FullNodes         []string                       // List of full nodes for on-chain operations
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
}

// PaymentChannel represents an active payment channel managed by the Lightning Node.
type PaymentChannel struct {
	ChannelID      string                      // Unique identifier for the payment channel
	ParticipantA   string                      // First participant in the channel
	ParticipantB   string                      // Second participant in the channel
	BalanceA       float64                     // Balance of Participant A in the channel
	BalanceB       float64                     // Balance of Participant B in the channel
	ChannelState   *ledger.ChannelState        // Current state of the payment channel
	ChannelTimeout time.Time                   // Expiry time for the channel
	IsActive       bool                        // Indicates whether the channel is active
}

// NewLightningNode initializes a new Lightning Node in the network for managing off-chain transactions.
func NewLightningNode(nodeID string, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syncInterval time.Duration, fullNodes []string) *LightningNode {
	return &LightningNode{
		NodeID:            nodeID,
		PaymentChannels:   make(map[string]*PaymentChannel),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SyncInterval:      syncInterval,
		FullNodes:         fullNodes,
		LiquidityPool:     make(map[string]float64),
	}
}

// StartNode begins the operations of the Lightning Node, including channel management and transaction processing.
func (ln *LightningNode) StartNode() error {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	// Start syncing with the blockchain periodically.
	go ln.syncWithFullNodes()

	// Listen for incoming channel requests and transactions.
	go ln.listenForChannelRequests()

	fmt.Printf("Lightning node %s started successfully.\n", ln.NodeID)
	return nil
}

// syncWithFullNodes handles syncing the payment channels with full nodes at regular intervals to settle disputes or close channels.
func (ln *LightningNode) syncWithFullNodes() {
	ticker := time.NewTicker(ln.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		ln.mutex.Lock()
		for _, fullNode := range ln.FullNodes {
			// Sync payment channel states to the blockchain via full nodes.
			ln.syncChannelsWithFullNode(fullNode)
		}
		ln.mutex.Unlock()
	}
}

// syncChannelsWithFullNode synchronizes the active payment channels with a full node.
func (ln *LightningNode) syncChannelsWithFullNode(fullNodeID string) {
	for channelID, channel := range ln.PaymentChannels {
		if channel.IsActive {
			// Encrypt the channel state for secure transmission.
			encryptedState, err := ln.EncryptChannelState(channel.ChannelState)
			if err != nil {
				fmt.Printf("Failed to encrypt channel state for channel %s: %v\n", channelID, err)
				continue
			}

			// Sync the state with the full node.
			err = ln.NetworkManager.SyncChannelWithFullNode(fullNodeID, encryptedState)
			if err != nil {
				fmt.Printf("Failed to sync channel %s with full node %s: %v\n", channelID, fullNodeID, err)
			} else {
				fmt.Printf("Channel %s synced successfully with full node %s.\n", channelID, fullNodeID)
			}
		}
	}
}

// listenForChannelRequests listens for requests to open, close, or update payment channels.
func (ln *LightningNode) listenForChannelRequests() {
	for {
		request, err := ln.NetworkManager.ReceiveChannelRequest()
		if err != nil {
			fmt.Printf("Error receiving channel request: %v\n", err)
			continue
		}

		// Process the channel request.
		err = ln.processChannelRequest(request)
		if err != nil {
			fmt.Printf("Channel request processing failed: %v\n", err)
		}
	}
}

// processChannelRequest handles opening, closing, and updating payment channels.
func (ln *LightningNode) processChannelRequest(request *network.ChannelRequest) error {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	switch request.Action {
	case "open":
		return ln.openPaymentChannel(request.ParticipantA, request.ParticipantB, request.InitialBalanceA, request.InitialBalanceB)
	case "close":
		return ln.closePaymentChannel(request.ChannelID)
	case "update":
		return ln.updatePaymentChannel(request.ChannelID, request.BalanceA, request.BalanceB)
	default:
		return errors.New("invalid channel request action")
	}
}

// openPaymentChannel opens a new payment channel between two participants.
func (ln *LightningNode) openPaymentChannel(participantA, participantB string, balanceA, balanceB float64) error {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	channelID := common.GenerateUniqueID()
	channelState := &ledger.ChannelState{
		ChannelID:    channelID,
		ParticipantA: participantA,
		ParticipantB: participantB,
		BalanceA:     balanceA,
		BalanceB:     balanceB,
	}

	ln.PaymentChannels[channelID] = &PaymentChannel{
		ChannelID:    channelID,
		ParticipantA: participantA,
		ParticipantB: participantB,
		BalanceA:     balanceA,
		BalanceB:     balanceB,
		ChannelState: channelState,
		IsActive:     true,
	}

	// Broadcast channel opening to the network.
	err := ln.NetworkManager.BroadcastChannelOpening(channelID)
	if err != nil {
		return fmt.Errorf("failed to broadcast channel opening: %v", err)
	}

	fmt.Printf("Payment channel %s opened between %s and %s.\n", channelID, participantA, participantB)
	return nil
}

// closePaymentChannel closes an active payment channel and syncs the final state with the blockchain.
func (ln *LightningNode) closePaymentChannel(channelID string) error {
	channel, exists := ln.PaymentChannels[channelID]
	if !exists || !channel.IsActive {
		return errors.New("channel does not exist or is not active")
	}

	// Sync final state with a full node for settlement on the blockchain.
	fullNodeID := ln.selectRandomFullNode()
	err := ln.syncChannelsWithFullNode(fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to sync final state of channel %s with full node %s: %v", channelID, fullNodeID, err)
	}

	channel.IsActive = false
	fmt.Printf("Payment channel %s closed.\n", channelID)
	return nil
}

// updatePaymentChannel updates the balances in an active payment channel.
func (ln *LightningNode) updatePaymentChannel(channelID string, newBalanceA, newBalanceB float64) error {
	channel, exists := ln.PaymentChannels[channelID]
	if !exists || !channel.IsActive {
		return errors.New("channel does not exist or is not active")
	}

	channel.BalanceA = newBalanceA
	channel.BalanceB = newBalanceB
	channel.ChannelState.BalanceA = newBalanceA
	channel.ChannelState.BalanceB = newBalanceB

	// Sync updated state with a full node.
	fullNodeID := ln.selectRandomFullNode()
	err := ln.syncChannelsWithFullNode(fullNodeID)
	if err != nil {
		return fmt.Errorf("failed to sync updated state of channel %s with full node %s: %v", channelID, fullNodeID, err)
	}

	fmt.Printf("Payment channel %s updated with new balances: A=%f, B=%f.\n", channelID, newBalanceA, newBalanceB)
	return nil
}

// EncryptChannelState encrypts the channel state before syncing with full nodes.
func (ln *LightningNode) EncryptChannelState(channelState *ledger.ChannelState) ([]byte, error) {
	encryptedState, err := ln.EncryptionService.EncryptData([]byte(channelState.Serialize()), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt channel state: %v", err)
	}
	return encryptedState, nil
}

// selectRandomFullNode selects a random full node for syncing payment channels.
func (ln *LightningNode) selectRandomFullNode() string {
	ln.mutex.Lock()
	defer ln.mutex.Unlock()

	if len(ln.FullNodes) == 0 {
		fmt.Println("No full nodes available for transaction forwarding.")
		return ""
	}

	randomIndex := common.GenerateRandomInt(len(ln.FullNodes))
	return ln.FullNodes[randomIndex]
}

