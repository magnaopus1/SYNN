package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    P2PEncryptionMonitoringInterval  = 10 * time.Second // Interval for monitoring P2P communication encryption
    MaxEncryptionRetries             = 3                // Maximum retries for encrypting P2P communication
    SubBlocksPerBlock                = 1000             // Number of sub-blocks in a block
    UnauthorizedP2PAlertThreshold    = 5                // Threshold for alerting on unauthorized P2P communication attempts
)

// P2PCommunicationEncryptionProtocol manages the encryption and security of peer-to-peer communications
type P2PCommunicationEncryptionProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging encryption-related events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    encryptionRetryCount     map[string]int               // Counter for retrying encryption on communication
    p2pCommunicationCycleCount int                        // Counter for monitoring cycles
    unauthorizedCommCount    map[string]int               // Tracks unauthorized communication attempts
}

// NewP2PCommunicationEncryptionProtocol initializes the automation for P2P communication encryption
func NewP2PCommunicationEncryptionProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *P2PCommunicationEncryptionProtocol {
    return &P2PCommunicationEncryptionProtocol{
        consensusSystem:           consensusSystem,
        ledgerInstance:            ledgerInstance,
        stateMutex:                stateMutex,
        encryptionRetryCount:      make(map[string]int),
        unauthorizedCommCount:     make(map[string]int),
        p2pCommunicationCycleCount: 0,
    }
}

// StartP2PEncryptionMonitoring starts the continuous loop for monitoring and enforcing P2P communication encryption
func (protocol *P2PCommunicationEncryptionProtocol) StartP2PEncryptionMonitoring() {
    ticker := time.NewTicker(P2PEncryptionMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorP2PCommunicationEncryption()
        }
    }()
}

// monitorP2PCommunicationEncryption monitors the P2P communication and ensures encryption is properly applied
func (protocol *P2PCommunicationEncryptionProtocol) monitorP2PCommunicationEncryption() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch the list of active P2P communications from the consensus system
    p2pCommunications := protocol.consensusSystem.FetchP2PCommunications()

    for _, comm := range p2pCommunications {
        if protocol.isCommunicationEncrypted(comm) {
            fmt.Printf("P2P communication from node %s to node %s is encrypted and secure.\n", comm.SenderID, comm.ReceiverID)
            protocol.logP2PEncryptionEvent(comm, "Encrypted")
        } else {
            fmt.Printf("Unencrypted P2P communication detected between node %s and node %s. Triggering encryption.\n", comm.SenderID, comm.ReceiverID)
            protocol.handleUnencryptedCommunication(comm)
        }
    }

    protocol.p2pCommunicationCycleCount++
    fmt.Printf("P2P communication encryption cycle #%d completed.\n", protocol.p2pCommunicationCycleCount)

    if protocol.p2pCommunicationCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeEncryptionCycle()
    }
}

// isCommunicationEncrypted checks if the P2P communication is encrypted
func (protocol *P2PCommunicationEncryptionProtocol) isCommunicationEncrypted(comm common.P2PCommunication) bool {
    // Example logic to check if the communication is encrypted (this logic can be customized based on encryption keys, protocols)
    return comm.IsEncrypted
}

// handleUnencryptedCommunication handles P2P communications that are unencrypted and attempts to secure them
func (protocol *P2PCommunicationEncryptionProtocol) handleUnencryptedCommunication(comm common.P2PCommunication) {
    protocol.unauthorizedCommCount[comm.SenderID]++

    if protocol.unauthorizedCommCount[comm.SenderID] >= UnauthorizedP2PAlertThreshold {
        fmt.Printf("Multiple unencrypted P2P communications detected for node %s. Taking action.\n", comm.SenderID)
        protocol.blockUnencryptedCommunication(comm)
    } else {
        fmt.Printf("Applying encryption to P2P communication from node %s to node %s.\n", comm.SenderID, comm.ReceiverID)
        protocol.applyEncryptionToCommunication(comm)
    }
}

// applyEncryptionToCommunication applies encryption to an unencrypted P2P communication
func (protocol *P2PCommunicationEncryptionProtocol) applyEncryptionToCommunication(comm common.P2PCommunication) {
    encryptedComm := protocol.encryptCommunicationData(comm)

    // Attempt to secure the P2P communication through the Synnergy Consensus system
    encryptionSuccess := protocol.consensusSystem.ApplyEncryptionToP2PCommunication(encryptedComm)

    if encryptionSuccess {
        fmt.Printf("Encryption successfully applied to P2P communication between node %s and node %s.\n", comm.SenderID, comm.ReceiverID)
        protocol.logP2PEncryptionEvent(comm, "Encrypted")
        protocol.resetEncryptionRetry(comm.SenderID)
    } else {
        fmt.Printf("Error applying encryption to P2P communication between node %s and node %s. Retrying...\n", comm.SenderID, comm.ReceiverID)
        protocol.retryEncryption(comm)
    }
}

// blockUnencryptedCommunication blocks unencrypted P2P communication after repeated unauthorized attempts
func (protocol *P2PCommunicationEncryptionProtocol) blockUnencryptedCommunication(comm common.P2PCommunication) {
    encryptedComm := protocol.encryptCommunicationData(comm)

    // Attempt to block the unencrypted P2P communication through the Synnergy Consensus system
    blockSuccess := protocol.consensusSystem.BlockUnencryptedP2PCommunication(encryptedComm)

    if blockSuccess {
        fmt.Printf("Unencrypted P2P communication blocked between node %s and node %s.\n", comm.SenderID, comm.ReceiverID)
        protocol.logP2PEncryptionEvent(comm, "Blocked")
        protocol.resetEncryptionRetry(comm.SenderID)
    } else {
        fmt.Printf("Error blocking unencrypted P2P communication between node %s and node %s. Retrying...\n", comm.SenderID, comm.ReceiverID)
        protocol.retryEncryption(comm)
    }
}

// retryEncryption retries encryption of P2P communication in case of failure
func (protocol *P2PCommunicationEncryptionProtocol) retryEncryption(comm common.P2PCommunication) {
    protocol.encryptionRetryCount[comm.SenderID]++
    if protocol.encryptionRetryCount[comm.SenderID] < MaxEncryptionRetries {
        protocol.applyEncryptionToCommunication(comm)
    } else {
        fmt.Printf("Max retries reached for encrypting P2P communication between node %s and node %s. Action failed.\n", comm.SenderID, comm.ReceiverID)
        protocol.logEncryptionFailure(comm)
    }
}

// resetEncryptionRetry resets the retry count for encryption actions on a specific P2P communication
func (protocol *P2PCommunicationEncryptionProtocol) resetEncryptionRetry(nodeID string) {
    protocol.encryptionRetryCount[nodeID] = 0
}

// finalizeEncryptionCycle finalizes the encryption monitoring cycle and logs the result in the ledger
func (protocol *P2PCommunicationEncryptionProtocol) finalizeEncryptionCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeEncryptionCycle()
    if success {
        fmt.Println("P2P communication encryption cycle finalized successfully.")
        protocol.logEncryptionCycleFinalization()
    } else {
        fmt.Println("Error finalizing P2P communication encryption cycle.")
    }
}

// logP2PEncryptionEvent logs a P2P communication encryption event into the ledger
func (protocol *P2PCommunicationEncryptionProtocol) logP2PEncryptionEvent(comm common.P2PCommunication, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("p2p-encryption-%s-%s", comm.SenderID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "P2P Communication Encryption Event",
        Status:    eventType,
        Details:   fmt.Sprintf("P2P communication between node %s and node %s was %s.", comm.SenderID, comm.ReceiverID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with P2P encryption event for nodes %s and %s.\n", comm.SenderID, comm.ReceiverID)
}

// logEncryptionFailure logs the failure to encrypt a P2P communication into the ledger
func (protocol *P2PCommunicationEncryptionProtocol) logEncryptionFailure(comm common.P2PCommunication) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("p2p-encryption-failure-%s", comm.SenderID),
        Timestamp: time.Now().Unix(),
        Type:      "P2P Encryption Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to encrypt P2P communication between node %s and node %s after maximum retries.", comm.SenderID, comm.ReceiverID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with P2P encryption failure for nodes %s and %s.\n", comm.SenderID, comm.ReceiverID)
}

// logEncryptionCycleFinalization logs the finalization of a P2P communication encryption cycle into the ledger
func (protocol *P2PCommunicationEncryptionProtocol) logEncryptionCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("encryption-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Encryption Cycle Finalization",
        Status:    "Finalized",
        Details:   "P2P communication encryption cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with encryption cycle finalization.")
}

// encryptCommunicationData encrypts the P2P communication data before applying or blocking communication
func (protocol *P2PCommunicationEncryptionProtocol) encryptCommunicationData(comm common.P2PCommunication) common.P2PCommunication {
    encryptedData, err := encryption.EncryptData(comm.Data)
    if err != nil {
        fmt.Println("Error encrypting P2P communication data:", err)
        return comm
    }

    comm.EncryptedData = encryptedData
    fmt.Println("P2P communication data successfully encrypted for nodes:", comm.SenderID, comm.ReceiverID)
    return comm
}

// triggerEmergencyP2PCommunicationLockdown triggers an emergency lockdown on P2P communication in case of severe security issues
func (protocol *P2PCommunicationEncryptionProtocol) triggerEmergencyP2PCommunicationLockdown(senderID string, receiverID string) {
    fmt.Printf("Emergency communication lockdown triggered for nodes %s and %s.\n", senderID, receiverID)
    comm := protocol.consensusSystem.GetP2PCommunicationByID(senderID, receiverID)
    encryptedData := protocol.encryptCommunicationData(comm)

    success := protocol.consensusSystem.TriggerEmergencyP2PCommunicationLockdown(senderID, receiverID, encryptedData)

    if success {
        protocol.logP2PEncryptionEvent(comm, "Emergency Locked Down")
        fmt.Println("Emergency communication lockdown executed successfully.")
    } else {
        fmt.Println("Emergency communication lockdown failed.")
    }
}
