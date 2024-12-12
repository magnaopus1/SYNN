package storage

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewSwarmManager initializes a new SwarmManager
func NewSwarmManager(ledgerInstance *ledger.Ledger) *SwarmManager {
    return &SwarmManager{
        LedgerInstance: ledgerInstance,
        SwarmNodes:     make(map[string]SwarmNode),
        StorageMap:     make(map[string]StorageEntry),
    }
}

// AddSwarmNode adds a new node to the swarm network
func (sm *SwarmManager) AddSwarmNode(nodeID, url string, capacity int64) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if _, exists := sm.SwarmNodes[nodeID]; exists {
        return fmt.Errorf("Swarm node %s already exists", nodeID)
    }

    newNode := SwarmNode{
        NodeID:   nodeID,
        URL:      url,
        Capacity: capacity,
        UsedSpace: 0,
        Status:   "active",
    }

    sm.SwarmNodes[nodeID] = newNode

    // Log the node addition to the ledger
    return sm.logNodeToLedger(newNode, "node_added")
}

// RemoveSwarmNode removes a node from the swarm
func (sm *SwarmManager) RemoveSwarmNode(nodeID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    node, exists := sm.SwarmNodes[nodeID]
    if !exists {
        return fmt.Errorf("Swarm node %s not found", nodeID)
    }

    delete(sm.SwarmNodes, nodeID)

    // Log the node removal to the ledger
    return sm.logNodeToLedger(node, "node_removed")
}

// StoreInSwarm stores encrypted data in a swarm node
func (sm *SwarmManager) StoreInSwarm(ownerID string, data []byte, expirationDuration time.Duration) (string, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Select a swarm node with enough capacity
    swarmNode, err := sm.selectNodeForStorage(len(data))
    if err != nil {
        return "", err
    }

    // Generate unique storage ID
    storageID := common.GenerateStorageID(ownerID, time.Now().UnixNano())

    // Encrypt the data
    encryptedData, err := encryption.EncryptData(string(data), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt storage data: %v", err)
    }

    // Create the storage entry
    newEntry := common.StorageEntry{
        StorageID:    storageID,
        OwnerID:      ownerID,
        Data:         []byte(encryptedData),
        CreatedAt:    time.Now(),
        ExpiresAt:    time.Now().Add(expirationDuration),
        Location:     swarmNode.NodeID,
    }

    // Store in the selected swarm node
    sm.StorageMap[storageID] = newEntry
    swarmNode.UsedSpace += int64(len(data))

    // Log storage to ledger
    err = sm.logStorageToLedger(newEntry)
    if err != nil {
        return "", fmt.Errorf("failed to log storage to ledger: %v", err)
    }

    fmt.Printf("Stored %d bytes in swarm node %s for storage ID %s.\n", len(data), swarmNode.NodeID, storageID)
    return storageID, nil
}

// RetrieveFromSwarm retrieves encrypted data from the swarm and decrypts it
func (sm *SwarmManager) RetrieveFromSwarm(storageID string) ([]byte, error) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    entry, exists := sm.StorageMap[storageID]
    if !exists {
        return nil, fmt.Errorf("storage entry %s not found", storageID)
    }

    if time.Now().After(entry.ExpiresAt) {
        return nil, fmt.Errorf("storage entry %s has expired", storageID)
    }

    // Decrypt the data
    decryptedData, err := encryption.DecryptData(string(entry.Data), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt storage data: %v", err)
    }

    fmt.Printf("Retrieved storage entry %s from swarm node %s.\n", storageID, entry.Location)
    return []byte(decryptedData), nil
}

// selectNodeForStorage selects a swarm node based on available capacity
func (sm *SwarmManager) selectNodeForStorage(dataSize int) (*SwarmNode, error) {
    for _, node := range sm.SwarmNodes {
        if node.Capacity-node.UsedSpace > int64(dataSize) {
            return &node, nil
        }
    }
    return nil, fmt.Errorf("no swarm nodes available with sufficient capacity")
}

// logStorageToLedger logs the creation of a storage entry to the ledger
func (sm *SwarmManager) logStorageToLedger(entry StorageEntry) error {
    storageRecord := fmt.Sprintf("StorageID: %s, OwnerID: %s, Location: %s, CreatedAt: %s, ExpiresAt: %s",
        entry.StorageID, entry.OwnerID, entry.Location, entry.CreatedAt.String(), entry.ExpiresAt.String())

    encryptedRecord, err := common.EncryptData(storageRecord, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt storage log: %v", err)
    }

    return sm.LedgerInstance.RecordTransaction(entry.StorageID, "storage_creation", entry.OwnerID, encryptedRecord)
}

// logNodeToLedger logs the addition or removal of a swarm node to the ledger
func (sm *SwarmManager) logNodeToLedger(node SwarmNode, action string) error {
    nodeRecord := fmt.Sprintf("NodeID: %s, URL: %s, Capacity: %d, Status: %s", 
        node.NodeID, node.URL, node.Capacity, node.Status)

    encryptedRecord, err := common.EncryptData(nodeRecord, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt node log: %v", err)
    }

    return sm.LedgerInstance.RecordTransaction(node.NodeID, action, "", encryptedRecord)
}


SWARM_UPLOAD: Uploads a file to Swarm and returns a Swarm hash as the file identifier.
SWARM_DOWNLOAD: Retrieves a file from Swarm using its hash.
SWARM_PIN: Pins a file in Swarm to ensure it remains accessible.
SWARM_UNPIN: Unpins a file in Swarm, making it available for deletion.
SWARM_METADATA: Fetches metadata for a file stored on Swarm.
SWARM_LIST_FILES: Lists files stored on the node within Swarm.
SWARM_ENCRYPT: Encrypts data before uploading to Swarm to secure private files.
SWARM_DECRYPT: Decrypts encrypted data retrieved from Swarm.
SWARM_SET_RETRIEVAL: Configures file retrieval settings, like speed, cost, or priority.