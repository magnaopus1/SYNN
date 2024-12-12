package network

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewFaultToleranceManager initializes the fault tolerance system
func NewFaultToleranceManager(nodes []string, quorumThreshold int, ledger *ledger.Ledger) *FaultToleranceManager {
    nodeStatus := make(map[string]bool)
    for _, node := range nodes {
        nodeStatus[node] = true // Initialize all nodes as alive
    }

    return &FaultToleranceManager{
        Nodes:           nodes,
        NodeStatus:      nodeStatus,
        NodeState:  make(map[string]string),
        QuorumThreshold: quorumThreshold,
        ledger:          ledger,
    }
}

// NodeFailure detects if a node is down and marks it as such
func (ft *FaultToleranceManager) NodeFailure(node string) {
    ft.mutex.Lock()
    defer ft.mutex.Unlock()

    if !ft.NodeStatus[node] {
        return
    }

    ft.NodeStatus[node] = false
    fmt.Printf("Node %s marked as down.\n", node)
    ft.ledger.LogFaultEvent("NodeDown", node, time.Now())
}

// NodeRecovery marks a node as recovered and back online
func (ft *FaultToleranceManager) NodeRecovery(node string) {
    ft.mutex.Lock()
    defer ft.mutex.Unlock()

    if ft.NodeStatus[node] {
        return
    }

    ft.NodeStatus[node] = true
    fmt.Printf("Node %s recovered and back online.\n", node)
    ft.ledger.LogFaultEvent("NodeRecovered", node, time.Now())
}

// IsQuorumAlive checks if the network has enough nodes for consensus
func (ft *FaultToleranceManager) IsQuorumAlive() bool {
    ft.mutex.Lock()
    defer ft.mutex.Unlock()

    aliveCount := 0
    for _, status := range ft.NodeStatus {
        if status {
            aliveCount++
        }
    }

    if aliveCount >= ft.QuorumThreshold {
        fmt.Println("Quorum is alive, network can proceed.")
        return true
    }

    fmt.Println("Quorum not met, network is in fault state.")
    ft.ledger.LogFaultEvent("QuorumFailure", "", time.Now())
    return false
}

// SyncNode synchronizes the state of the blockchain with a recovered node
func (ft *FaultToleranceManager) SyncNode(node string, latestBlockHash string) error {
    ft.mutex.Lock()
    defer ft.mutex.Unlock()

    // Check if the node is down, if so return an error
    if !ft.NodeStatus[node] {
        return fmt.Errorf("cannot sync node %s, it is down", node)
    }

    // Simulate the syncing process
    fmt.Printf("Starting synchronization process for node %s with latest block hash: %s.\n", node, latestBlockHash)

    // Retrieve the latest block hash from the ledger (simulate this)
    currentBlockHash := ft.ledger.GetLatestBlockHash()

    // Check if the node is already synced
    if currentBlockHash == latestBlockHash {
        fmt.Printf("Node %s is already synchronized with the latest block hash: %s.\n", node, latestBlockHash)
        ft.ledger.LogSyncEvent(fmt.Sprintf("NodeSyncComplete: Node %s is already synced with the latest block hash.", node))
        return nil
    }

    // If not synced, proceed with the sync process
    fmt.Printf("Node %s is out of sync, syncing with block hash: %s (latest: %s).\n", node, latestBlockHash, currentBlockHash)

    // Simulate updating the node to the latest block
    ft.updateNodeState(node, latestBlockHash)

    // Log the synchronization event with node, block hash, and timestamp
    ft.ledger.LogSyncEvent(fmt.Sprintf("NodeSync: Node %s synced with block hash: %s at %s", node, latestBlockHash, time.Now()))

    fmt.Printf("Synchronization of node %s completed successfully.\n", node)
    return nil
}

// updateNodeState updates the latest synced block hash for a given node
func (ft *FaultToleranceManager) updateNodeState(node string, latestBlockHash string) {
    ft.mutex.Lock()
    defer ft.mutex.Unlock()

    // Update the node's state with the latest block hash
    ft.NodeState[node] = latestBlockHash
    fmt.Printf("Node %s successfully updated to latest block hash: %s.\n", node, latestBlockHash)
}

// EncryptData encrypts the data exchanged between nodes during fault recovery
func (ft *FaultToleranceManager) EncryptData(data string, pubKey *common.PublicKey) (string, error) {
    hash := sha256.New()
    hash.Write([]byte(data))
    encryptedData := hex.EncodeToString(hash.Sum(nil))

    fmt.Println("Data encrypted for secure transmission.")
    return encryptedData, nil
}

// DecryptData decrypts the data received during fault recovery
func (ft *FaultToleranceManager) DecryptData(encryptedData string) (string, error) {
    hashBytes, err := hex.DecodeString(encryptedData)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt data: %v", err)
    }

    return string(hashBytes), nil
}

// BroadcastFailureNotification notifies the network of a node failure
func (ft *FaultToleranceManager) BroadcastFailureNotification(node string) {
    fmt.Printf("Broadcasting failure notification for node %s.\n", node)
    ft.ledger.LogFaultEvent("FailureBroadcast", node, time.Now())
}

// BroadcastRecoveryNotification notifies the network of a node recovery
func (ft *FaultToleranceManager) BroadcastRecoveryNotification(node string) {
    fmt.Printf("Broadcasting recovery notification for node %s.\n", node)
    ft.ledger.LogFaultEvent("RecoveryBroadcast", node, time.Now())
}

// MonitorNetwork continuously monitors the network for node failures
func (ft *FaultToleranceManager) MonitorNetwork(interval time.Duration) {
    for {
        time.Sleep(interval)
        ft.mutex.Lock()

        for node, status := range ft.NodeStatus {
            if !status {
                fmt.Printf("Detected failure in node %s.\n", node)
                ft.BroadcastFailureNotification(node)
            }
        }

        ft.mutex.Unlock()
    }
}
