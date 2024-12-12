package automations

import (
    "fmt"
    "sync"
    "time"
    "errors"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/peer_to_peer"
)

const (
    GossipCheckInterval      = 10 * time.Second  // Interval to check for new gossip messages
    MaxMessageSize           = 1024              // Maximum size for gossip messages (in bytes)
    GossipPropagationTimeout = 30 * time.Second  // Timeout for message propagation
    MaxPeerConnections       = 100               // Maximum number of connected peers
)

// GossipProtocolAutomation handles gossip protocol automation for scalability
type GossipProtocolAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus
    ledgerInstance   *ledger.Ledger
    peerNetwork      *peer_to_peer.PeerNetwork
    stateMutex       *sync.RWMutex
    gossipCycle      int
    peerConnections  int
}

// NewGossipProtocolAutomation initializes the gossip protocol automation
func NewGossipProtocolAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, peerNetwork *peer_to_peer.PeerNetwork, stateMutex *sync.RWMutex) *GossipProtocolAutomation {
    return &GossipProtocolAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        peerNetwork:     peerNetwork,
        stateMutex:      stateMutex,
        gossipCycle:     0,
        peerConnections: 0,
    }
}

// StartGossipMonitoring starts the gossip protocol monitoring in a continuous loop
func (automation *GossipProtocolAutomation) StartGossipMonitoring() {
    ticker := time.NewTicker(GossipCheckInterval)

    go func() {
        for range ticker.C {
            automation.handleGossipMessages()
        }
    }()
}

// handleGossipMessages checks and propagates gossip messages
func (automation *GossipProtocolAutomation) handleGossipMessages() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if automation.peerConnections >= MaxPeerConnections {
        fmt.Println("Maximum peer connections reached. Limiting further connections.")
        return
    }

    newMessages := automation.fetchNewGossipMessages()
    if len(newMessages) == 0 {
        return
    }

    for _, message := range newMessages {
        if len(message.Content) > MaxMessageSize {
            fmt.Printf("Gossip message exceeds size limit (%d bytes). Skipping...\n", len(message.Content))
            continue
        }

        err := automation.propagateMessage(message)
        if err != nil {
            fmt.Printf("Failed to propagate gossip message: %v\n", err)
            automation.logPropagationFailure(message, err)
        } else {
            automation.logPropagationSuccess(message)
        }
    }

    automation.gossipCycle++
    fmt.Printf("Gossip cycle #%d completed.\n", automation.gossipCycle)

    if automation.gossipCycle % 100 == 0 {
        automation.finalizeGossipCycle()
    }
}

// fetchNewGossipMessages fetches any new gossip messages from the network
func (automation *GossipProtocolAutomation) fetchNewGossipMessages() []common.GossipMessage {
    return automation.peerNetwork.GetNewGossipMessages()
}

// propagateMessage propagates the gossip message to peer nodes
func (automation *GossipProtocolAutomation) propagateMessage(message common.GossipMessage) error {
    fmt.Printf("Propagating gossip message: %s\n", message.Content)
    
    encryptedContent, err := encryption.EncryptData([]byte(message.Content))
    if err != nil {
        return fmt.Errorf("failed to encrypt gossip message: %v", err)
    }

    // Send the encrypted message to all connected peers
    propagationSuccess := automation.peerNetwork.SendMessageToPeers(encryptedContent)
    if !propagationSuccess {
        return errors.New("message propagation failed due to network issues")
    }

    return nil
}

// finalizeGossipCycle finalizes the current gossip cycle in the consensus system
func (automation *GossipProtocolAutomation) finalizeGossipCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeGossipCycle()
    if success {
        fmt.Println("Gossip cycle finalized successfully.")
        automation.logCycleFinalization()
    } else {
        fmt.Println("Error finalizing gossip cycle.")
    }
}

// logPropagationSuccess logs the successful propagation of a gossip message in the ledger
func (automation *GossipProtocolAutomation) logPropagationSuccess(message common.GossipMessage) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gossip-propagation-%s", message.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Gossip Propagation",
        Status:    "Success",
        Details:   fmt.Sprintf("Gossip message propagated successfully: %s", message.Content),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with gossip message propagation success for message %s.\n", message.ID)
}

// logPropagationFailure logs the failed propagation of a gossip message in the ledger
func (automation *GossipProtocolAutomation) logPropagationFailure(message common.GossipMessage, err error) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gossip-propagation-failure-%s", message.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Gossip Propagation",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to propagate gossip message %s: %v", message.Content, err),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with gossip message propagation failure for message %s.\n", message.ID)
}

// logCycleFinalization logs the finalization of the gossip cycle in the ledger
func (automation *GossipProtocolAutomation) logCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("gossip-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Gossip Cycle",
        Status:    "Finalized",
        Details:   "Gossip cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with gossip cycle finalization.")
}
