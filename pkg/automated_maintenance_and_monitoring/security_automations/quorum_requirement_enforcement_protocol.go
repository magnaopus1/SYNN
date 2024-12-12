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
    QuorumMonitoringInterval      = 10 * time.Second // Interval for monitoring quorum requirements
    MaxQuorumEnforcementRetries   = 3                // Maximum retries for enforcing quorum requirements
    SubBlocksPerBlock             = 1000             // Number of sub-blocks in a block
    QuorumThreshold               = 0.67             // Quorum threshold (e.g., 67% of validators)
)

// QuorumRequirementEnforcementProtocol enforces quorum during decision-making processes
type QuorumRequirementEnforcementProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging quorum-related events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    quorumEnforcementRetryCount map[string]int            // Counter for retrying quorum enforcement actions
    quorumMonitoringCycleCount  int                       // Counter for quorum monitoring cycles
    quorumFailureCounter        map[string]int            // Tracks quorum failure attempts
}

// NewQuorumRequirementEnforcementProtocol initializes the protocol for enforcing quorum
func NewQuorumRequirementEnforcementProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *QuorumRequirementEnforcementProtocol {
    return &QuorumRequirementEnforcementProtocol{
        consensusSystem:            consensusSystem,
        ledgerInstance:             ledgerInstance,
        stateMutex:                 stateMutex,
        quorumEnforcementRetryCount: make(map[string]int),
        quorumFailureCounter:        make(map[string]int),
        quorumMonitoringCycleCount:  0,
    }
}

// StartQuorumMonitoring starts the continuous loop for monitoring quorum and enforcing its requirement
func (protocol *QuorumRequirementEnforcementProtocol) StartQuorumMonitoring() {
    ticker := time.NewTicker(QuorumMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorForQuorum()
        }
    }()
}

// monitorForQuorum monitors the network to ensure quorum is met before allowing decisions or consensus operations
func (protocol *QuorumRequirementEnforcementProtocol) monitorForQuorum() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch active validators and their participation from the consensus system
    quorumStatus := protocol.consensusSystem.CheckQuorumStatus()

    if !quorumStatus.IsQuorumMet {
        fmt.Printf("Quorum not met for decision %s. Taking action.\n", quorumStatus.DecisionID)
        protocol.handleQuorumFailure(quorumStatus)
    } else {
        fmt.Printf("Quorum met for decision %s.\n", quorumStatus.DecisionID)
    }

    protocol.quorumMonitoringCycleCount++
    fmt.Printf("Quorum monitoring cycle #%d completed.\n", protocol.quorumMonitoringCycleCount)

    if protocol.quorumMonitoringCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeQuorumMonitoringCycle()
    }
}

// handleQuorumFailure handles quorum failures by enforcing penalties or blocking actions until quorum is met
func (protocol *QuorumRequirementEnforcementProtocol) handleQuorumFailure(quorumStatus common.QuorumStatus) {
    protocol.quorumFailureCounter[quorumStatus.DecisionID]++

    if protocol.quorumFailureCounter[quorumStatus.DecisionID] >= MaxQuorumEnforcementRetries {
        fmt.Printf("Multiple quorum failures detected for decision %s. Blocking decision.\n", quorumStatus.DecisionID)
        protocol.blockDecision(quorumStatus)
    } else {
        fmt.Printf("Issuing warning for quorum failure for decision %s.\n", quorumStatus.DecisionID)
        protocol.warnAboutQuorumFailure(quorumStatus)
    }
}

// warnAboutQuorumFailure issues a warning regarding quorum failure to validators
func (protocol *QuorumRequirementEnforcementProtocol) warnAboutQuorumFailure(quorumStatus common.QuorumStatus) {
    encryptedWarningData := protocol.encryptQuorumData(quorumStatus)

    // Issue a warning through the Synnergy Consensus system
    warningSuccess := protocol.consensusSystem.WarnAboutQuorumFailure(quorumStatus.DecisionID, encryptedWarningData)

    if warningSuccess {
        fmt.Printf("Warning issued for quorum failure in decision %s.\n", quorumStatus.DecisionID)
        protocol.logQuorumEvent(quorumStatus, "Warning Issued")
        protocol.resetQuorumEnforcementRetry(quorumStatus.DecisionID)
    } else {
        fmt.Printf("Error issuing quorum failure warning for decision %s. Retrying...\n", quorumStatus.DecisionID)
        protocol.retryQuorumEnforcementAction(quorumStatus)
    }
}

// blockDecision blocks the decision or consensus operation due to quorum failure
func (protocol *QuorumRequirementEnforcementProtocol) blockDecision(quorumStatus common.QuorumStatus) {
    encryptedBlockData := protocol.encryptQuorumData(quorumStatus)

    // Attempt to block the decision through the Synnergy Consensus system
    blockSuccess := protocol.consensusSystem.BlockDecision(quorumStatus.DecisionID, encryptedBlockData)

    if blockSuccess {
        fmt.Printf("Decision %s blocked due to quorum failure.\n", quorumStatus.DecisionID)
        protocol.logQuorumEvent(quorumStatus, "Decision Blocked")
        protocol.resetQuorumEnforcementRetry(quorumStatus.DecisionID)
    } else {
        fmt.Printf("Error blocking decision %s due to quorum failure. Retrying...\n", quorumStatus.DecisionID)
        protocol.retryQuorumEnforcementAction(quorumStatus)
    }
}

// retryQuorumEnforcementAction retries the quorum enforcement action if it initially fails
func (protocol *QuorumRequirementEnforcementProtocol) retryQuorumEnforcementAction(quorumStatus common.QuorumStatus) {
    protocol.quorumEnforcementRetryCount[quorumStatus.DecisionID]++
    if protocol.quorumEnforcementRetryCount[quorumStatus.DecisionID] < MaxQuorumEnforcementRetries {
        protocol.blockDecision(quorumStatus)
    } else {
        fmt.Printf("Max retries reached for enforcing quorum in decision %s. Action failed.\n", quorumStatus.DecisionID)
        protocol.logQuorumEnforcementFailure(quorumStatus)
    }
}

// resetQuorumEnforcementRetry resets the retry count for quorum enforcement actions on a specific decision
func (protocol *QuorumRequirementEnforcementProtocol) resetQuorumEnforcementRetry(decisionID string) {
    protocol.quorumEnforcementRetryCount[decisionID] = 0
}

// finalizeQuorumMonitoringCycle finalizes the quorum monitoring cycle and logs the result in the ledger
func (protocol *QuorumRequirementEnforcementProtocol) finalizeQuorumMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeQuorumMonitoringCycle()
    if success {
        fmt.Println("Quorum monitoring cycle finalized successfully.")
        protocol.logQuorumCycleFinalization()
    } else {
        fmt.Println("Error finalizing quorum monitoring cycle.")
    }
}

// logQuorumEvent logs a quorum-related event into the ledger
func (protocol *QuorumRequirementEnforcementProtocol) logQuorumEvent(quorumStatus common.QuorumStatus, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("quorum-event-%s-%s", quorumStatus.DecisionID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "Quorum Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Decision %s triggered %s due to quorum failure.", quorumStatus.DecisionID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with quorum event for decision %s.\n", quorumStatus.DecisionID)
}

// logQuorumEnforcementFailure logs the failure to enforce quorum into the ledger
func (protocol *QuorumRequirementEnforcementProtocol) logQuorumEnforcementFailure(quorumStatus common.QuorumStatus) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("quorum-enforcement-failure-%s", quorumStatus.DecisionID),
        Timestamp: time.Now().Unix(),
        Type:      "Quorum Enforcement Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to enforce quorum for decision %s after maximum retries.", quorumStatus.DecisionID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with quorum enforcement failure for decision %s.\n", quorumStatus.DecisionID)
}

// logQuorumCycleFinalization logs the finalization of a quorum monitoring cycle into the ledger
func (protocol *QuorumRequirementEnforcementProtocol) logQuorumCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("quorum-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Quorum Cycle Finalization",
        Status:    "Finalized",
        Details:   "Quorum monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with quorum monitoring cycle finalization.")
}

// encryptQuorumData encrypts the data related to quorum enforcement before taking action or logging events
func (protocol *QuorumRequirementEnforcementProtocol) encryptQuorumData(quorumStatus common.QuorumStatus) common.QuorumStatus {
    encryptedData, err := encryption.EncryptData(quorumStatus.Data)
    if err != nil {
        fmt.Println("Error encrypting quorum data:", err)
        return quorumStatus
    }

    quorumStatus.EncryptedData = encryptedData
    fmt.Println("Quorum data successfully encrypted for decision ID:", quorumStatus.DecisionID)
    return quorumStatus
}

// triggerEmergencyQuorumLockdown triggers an emergency quorum lockdown in case quorum is repeatedly not met
func (protocol *QuorumRequirementEnforcementProtocol) triggerEmergencyQuorumLockdown(decisionID string) {
    fmt.Printf("Emergency quorum lockdown triggered for decision ID: %s.\n", decisionID)
    quorumStatus := protocol.consensusSystem.GetQuorumStatusByID(decisionID)
    encryptedData := protocol.encryptQuorumData(quorumStatus)

    success := protocol.consensusSystem.TriggerEmergencyQuorumLockdown(decisionID, encryptedData)

    if success {
        protocol.logQuorumEvent(quorumStatus, "Emergency Locked Down")
        fmt.Println("Emergency quorum lockdown executed successfully.")
    } else {
        fmt.Println("Emergency quorum lockdown failed.")
    }
}
