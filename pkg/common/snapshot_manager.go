package common


import (
	"encoding/json"
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewSnapshotManager initializes a new SnapshotManager.
func NewSnapshotManager(ledgerInstance *ledger.Ledger) *SnapshotManager {
    return &SnapshotManager{
        LedgerInstance:  ledgerInstance,
        SnapshotStorage: make(map[string][]byte),
        CurrentState:    &BlockchainState{},
    }
}

// TakeSnapshot creates a new snapshot of the blockchain's current state and encrypts it.
func (sm *SnapshotManager) TakeSnapshot(snapshotID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Step 1: Serialize the current blockchain state
    stateData, err := sm.serializeState(sm.CurrentState)
    if err != nil {
        return fmt.Errorf("failed to serialize blockchain state: %v", err)
    }

    // Step 2: Encrypt the snapshot using an instance of Encryption
    encryptionInstance := &Encryption{}
    encryptedSnapshot, err := encryptionInstance.EncryptData("AES", stateData, EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt snapshot: %v", err)
    }

    // Step 3: Store the encrypted snapshot
    sm.SnapshotStorage[snapshotID] = encryptedSnapshot

    // Step 4: Log the snapshot creation into the ledger
    logEntry := fmt.Sprintf("Snapshot %s created at %s", snapshotID, time.Now().String())
    err = sm.LedgerInstance.VirtualMachineLedger.LogEntry(snapshotID, logEntry) // Now passing two arguments: snapshotID and logEntry
    if err != nil {
        return fmt.Errorf("failed to log snapshot in the ledger: %v", err)
    }

    fmt.Printf("Snapshot %s successfully created and stored.\n", snapshotID)
    return nil
}



// RestoreSnapshot restores the blockchain state from a given snapshot.
func (sm *SnapshotManager) RestoreSnapshot(snapshotID string) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    // Step 1: Retrieve the encrypted snapshot
    encryptedSnapshot, exists := sm.SnapshotStorage[snapshotID]
    if !exists {
        return fmt.Errorf("snapshot %s does not exist", snapshotID)
    }

    // Step 2: Decrypt the snapshot using an instance of Encryption
    encryptionInstance := &Encryption{}
    decryptedSnapshot, err := encryptionInstance.DecryptData(encryptedSnapshot, EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt snapshot: %v", err)
    }

    // Step 3: Deserialize the state data
    restoredState, err := sm.deserializeState(decryptedSnapshot)
    if err != nil {
        return fmt.Errorf("failed to deserialize snapshot: %v", err)
    }

    // Step 4: Update the current state with the restored state
    sm.CurrentState = restoredState

    // Step 5: Log the restoration into the ledger
    logEntry := fmt.Sprintf("Snapshot %s restored at %s", snapshotID, time.Now().String())
    err = sm.LedgerInstance.VirtualMachineLedger.LogEntry(snapshotID, logEntry) // Now passing two arguments: snapshotID and logEntry
    if err != nil {
        return fmt.Errorf("failed to log snapshot restoration in the ledger: %v", err)
    }

    fmt.Printf("Snapshot %s successfully restored.\n", snapshotID)
    return nil
}




// serializeState serializes the blockchain state into a byte array for storage.
func (sm *SnapshotManager) serializeState(state *BlockchainState) ([]byte, error) {
    // Use JSON serialization
    serializedData, err := json.Marshal(state)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize blockchain state: %v", err)
    }
    return serializedData, nil
}

// deserializeState deserializes a byte array back into a BlockchainState struct.
func (sm *SnapshotManager) deserializeState(data []byte) (*BlockchainState, error) {
    var restoredState BlockchainState
    // Use JSON deserialization
    err := json.Unmarshal(data, &restoredState)
    if err != nil {
        return nil, fmt.Errorf("failed to deserialize blockchain state: %v", err)
    }
    return &restoredState, nil
}

// PruneSnapshots removes older snapshots based on a retention policy (e.g., max number of snapshots or age).
func (sm *SnapshotManager) PruneSnapshots(maxSnapshots int) error {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if len(sm.SnapshotStorage) <= maxSnapshots {
        return nil // No pruning needed
    }

    // Prune the oldest snapshots until the number of snapshots is within the limit
    for len(sm.SnapshotStorage) > maxSnapshots {
        for snapshotID := range sm.SnapshotStorage {
            delete(sm.SnapshotStorage, snapshotID)
            break // Remove one snapshot at a time
        }
    }

    // Log the pruning operation in the ledger
    logEntry := fmt.Sprintf("Snapshot pruning performed at %s, remaining snapshots: %d", time.Now().String(), len(sm.SnapshotStorage))
    err := sm.LedgerInstance.VirtualMachineLedger.LogEntry("SnapshotPruning", logEntry) // Passing two arguments: identifier and log message
    if err != nil {
        return fmt.Errorf("failed to log snapshot pruning: %v", err)
    }

    fmt.Printf("Snapshot pruning completed. Retained %d snapshots.\n", maxSnapshots)
    return nil
}
