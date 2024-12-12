package high_availability

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// NodeData contains information about a node in the network
type NodeData struct {
    NodeID          string      // Unique identifier of the node
    PublicKey       string      // Public key of the node
    Address         string      // Network address of the node
    Metrics         NodeMetrics // Performance metrics of the node
    IsValidator     bool        // Whether the node is a validator in the network
    IsDataReplicated bool       // Whether the node has replicated data
    LastReplicated   time.Time  // Timestamp of the last replication
}

// RedundancyManager handles data replication and ensures redundancy across multiple nodes
type RedundancyManager struct {
    Nodes             map[string]*NodeData  // Node data with redundancy states
    LedgerInstance    *ledger.Ledger        // Ledger for recording redundancy operations
    mutex             sync.Mutex            // Mutex for thread-safe operations
    ReplicationFactor int                   // How many copies of the data should be replicated
}

// NewRedundancyManager initializes the redundancy manager
func NewRedundancyManager(nodes []string, ledger *ledger.Ledger, replicationFactor int) *RedundancyManager {
    nodeData := make(map[string]*NodeData)

    // Initialize each node's data replication status
    for _, node := range nodes {
        nodeData[node] = &NodeData{
            NodeID:           node,
            IsDataReplicated: false,      // Default to not replicated
            LastReplicated:   time.Now(), // Initialize with current time
        }
    }

    return &RedundancyManager{
        Nodes:             nodeData,
        LedgerInstance:    ledger,
        ReplicationFactor: replicationFactor,
    }
}


// StartRedundancyCheck continuously checks if all nodes have the required replication
func (rm *RedundancyManager) StartRedundancyCheck() {
    go func() {
        for {
            rm.mutex.Lock()
            for nodeID, data := range rm.Nodes {
                rm.checkRedundancy(nodeID, data)
            }
            rm.mutex.Unlock()
            time.Sleep(time.Minute * 5)  // Redundancy check every 5 minutes
        }
    }()
}

// checkRedundancy ensures that each node has replicated data across other nodes
func (rm *RedundancyManager) checkRedundancy(nodeID string, data *NodeData) {
    fmt.Printf("Checking redundancy for node %s...\n", nodeID)

    // Simulating the check for replicated data
    if !data.IsDataReplicated || time.Since(data.LastReplicated) > time.Hour {
        fmt.Printf("Data for node %s is not sufficiently replicated. Starting replication.\n", nodeID)
        rm.replicateData(nodeID, data)
    } else {
        fmt.Printf("Node %s is fully redundant and replicated.\n", nodeID)
    }
}

// replicateData replicates the data from one node to other nodes, ensuring redundancy and data integrity.
func (rm *RedundancyManager) replicateData(sourceNodeID string, data *NodeData) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    fmt.Printf("Initiating data replication from node %s...\n", sourceNodeID)
    
    if _, exists := rm.Nodes[sourceNodeID]; !exists {
        return fmt.Errorf("source node %s does not exist in the network", sourceNodeID)
    }
    
    replicationCount := 0
    for _, targetNode := range rm.Nodes {
        // Skip the source node and already replicated nodes
        if targetNode.NodeID == sourceNodeID || targetNode.IsDataReplicated {
            continue
        }

        // Replicate the data to the target node
        err := rm.replicateToNode(targetNode, data)
        if err != nil {
            fmt.Printf("Failed to replicate data to node %s: %v\n", targetNode.NodeID, err)
            continue
        }

        // Mark as replicated
        targetNode.IsDataReplicated = true
        targetNode.LastReplicated = time.Now()
        replicationCount++
        
        // Break if we have reached the replication factor
        if replicationCount >= rm.ReplicationFactor {
            break
        }
    }

    // Ensure that the data was replicated to at least the required number of nodes
    if replicationCount >= rm.ReplicationFactor {
        fmt.Printf("Successfully replicated data from node %s across %d nodes.\n", sourceNodeID, replicationCount)
        
        // Mark the data as replicated on the source node
        data.IsDataReplicated = true
        data.LastReplicated = time.Now()

        // Log the replication event in the ledger (assuming it takes nodeID and a string timestamp)
        replicationTime := time.Now().Format(time.RFC3339)
        rm.LedgerInstance.LogReplication(sourceNodeID, replicationTime)
        
        return nil
    }

    // Log failure to replicate data sufficiently
    fmt.Printf("Insufficient replication for node %s. Only replicated to %d nodes.\n", sourceNodeID, replicationCount)
    return fmt.Errorf("failed to meet replication factor for node %s", sourceNodeID)
}

// replicateToNode handles the actual replication process to a target node
func (rm *RedundancyManager) replicateToNode(targetNode *NodeData, data *NodeData) error {
    // Step 1: Establish a secure connection
    conn, err := rm.establishSecureConnection(targetNode)
    if err != nil {
        return fmt.Errorf("error establishing connection: %v", err)
    }
    defer conn.Close()

    // Step 2: Transfer the data
    err = rm.transferData(conn, data)
    if err != nil {
        return fmt.Errorf("data transfer failed: %v", err)
    }

    // Step 3: Validate the data transfer by comparing hashes
    err = rm.validateDataTransfer(data, targetNode)
    if err != nil {
        return fmt.Errorf("data validation failed: %v", err)
    }

    return nil
}


// establishSecureConnection sets up a TLS-encrypted connection to the target node
func (rm *RedundancyManager) establishSecureConnection(targetNode *NodeData) (net.Conn, error) {
    // Load client certificates
    cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        return nil, fmt.Errorf("failed to load certificates: %v", err)
    }

    // Load CA cert
    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        return nil, fmt.Errorf("failed to load CA cert: %v", err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Configure TLS
    tlsConfig := &tls.Config{
        Certificates:       []tls.Certificate{cert},
        RootCAs:            caCertPool,
        InsecureSkipVerify: false, // Set to false to enforce verification in production
    }

    // Establish a TLS-secured connection
    conn, err := tls.Dial("tcp", targetNode.Address, tlsConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to establish secure connection: %v", err)
    }

    fmt.Printf("Secure connection established with node %s\n", targetNode.NodeID)
    return conn, nil
}

// chunkedTransfer transfers data in chunks over the connection to handle large payloads
func (rm *RedundancyManager) chunkedTransfer(conn net.Conn, data []byte) error {
    const chunkSize = 1024 * 64 // 64KB chunk size
    totalChunks := (len(data) + chunkSize - 1) / chunkSize // Calculate total chunks

    for i := 0; i < totalChunks; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(data) {
            end = len(data)
        }

        chunk := data[start:end]
        _, err := conn.Write(chunk)
        if err != nil {
            return fmt.Errorf("failed to send data chunk: %v", err)
        }
    }

    return nil
}

// transferData securely transfers data to the target node
func (rm *RedundancyManager) transferData(conn net.Conn, data *NodeData) error {
    // Serialize the data to be transferred (could be JSON, Protocol Buffers, etc.)
    serializedData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to serialize data for transfer: %v", err)
    }

    // Transfer the serialized data using chunking
    err = rm.chunkedTransfer(conn, serializedData)
    if err != nil {
        return fmt.Errorf("data transfer failed: %v", err)
    }

    fmt.Printf("Data successfully transferred to node.\n")
    return nil
}

// validateDataTransfer validates data integrity by comparing SHA-256 hashes
func (rm *RedundancyManager) validateDataTransfer(originalData, targetNodeData *NodeData) error {
    // Compute hash of the original data
    originalHash := sha256.New()
    originalBytes, err := json.Marshal(originalData)
    if err != nil {
        return fmt.Errorf("failed to serialize original data for hashing: %v", err)
    }
    originalHash.Write(originalBytes)
    originalHashSum := originalHash.Sum(nil)

    // Compute hash of the target node data
    targetHash := sha256.New()
    targetBytes, err := json.Marshal(targetNodeData)
    if err != nil {
        return fmt.Errorf("failed to serialize target data for hashing: %v", err)
    }
    targetHash.Write(targetBytes)
    targetHashSum := targetHash.Sum(nil)

    // Compare the two hashes
    if !equalHashes(originalHashSum, targetHashSum) {
        return errors.New("data mismatch during validation")
    }

    fmt.Printf("Data successfully validated for node %s.\n", targetNodeData.NodeID)
    return nil
}

// equalHashes checks if two SHA-256 hashes are equal
func equalHashes(a, b []byte) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}




// RecoverNodeData triggers recovery for a node if its data becomes unavailable
func (rm *RedundancyManager) RecoverNodeData(faultyNodeID string) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    fmt.Printf("Recovering data for faulty node %s...\n", faultyNodeID)

    // Simulate recovery by replicating data from other nodes
    recovered := false
    for _, node := range rm.Nodes {
        if node.NodeID != faultyNodeID && node.IsDataReplicated {
            fmt.Printf("Recovering data from node %s to node %s...\n", node.NodeID, faultyNodeID)
            recovered = true
            rm.Nodes[faultyNodeID].IsDataReplicated = true
            rm.Nodes[faultyNodeID].LastReplicated = time.Now()

            // Log the recovery to the ledger, formatting time to a string
            rm.LedgerInstance.LogDataRecovery(faultyNodeID, time.Now().Format(time.RFC3339))
            break
        }
    }

    if recovered {
        fmt.Printf("Data for node %s successfully recovered.\n", faultyNodeID)
    } else {
        fmt.Printf("Failed to recover data for node %s. Insufficient redundancy.\n", faultyNodeID)
    }
}


// IncreaseReplicationFactor dynamically increases the replication factor to enhance fault tolerance
func (rm *RedundancyManager) IncreaseReplicationFactor(newFactor int) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    if newFactor > rm.ReplicationFactor {
        fmt.Printf("Increasing replication factor from %d to %d.\n", rm.ReplicationFactor, newFactor)
        rm.ReplicationFactor = newFactor
    }
}

// GetReplicationStatus provides a detailed replication report for all nodes
func (rm *RedundancyManager) GetReplicationStatus() map[string]*NodeData {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    status := make(map[string]*NodeData)
    for nodeID, data := range rm.Nodes {
        status[nodeID] = data
    }
    return status
}
