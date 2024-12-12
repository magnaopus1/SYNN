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
    PoWResistanceMonitoringInterval  = 10 * time.Second // Interval for monitoring PoW resistance issues
    MaxResistanceRetries             = 3                // Maximum retries for enforcing PoW resistance actions
    SubBlocksPerBlock                = 1000             // Number of sub-blocks in a block
    HashRateThreshold                = 70.0             // Threshold for total hash rate concentration in percentage
)

// ProofOfWorkResistanceProtocol manages PoW resistance, focusing on preventing centralization of hash power and 51% attacks
type ProofOfWorkResistanceProtocol struct {
    consensusSystem          *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance           *ledger.Ledger               // Ledger for logging PoW resistance-related events
    stateMutex               *sync.RWMutex                // Mutex for thread-safe access
    resistanceRetryCount     map[string]int               // Counter for retrying PoW resistance actions
    powMonitoringCycleCount  int                          // Counter for monitoring cycles
    highHashRateEntityCount  map[string]int               // Tracks entities with unusually high hash rates
}

// NewProofOfWorkResistanceProtocol initializes the PoW resistance protocol
func NewProofOfWorkResistanceProtocol(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ProofOfWorkResistanceProtocol {
    return &ProofOfWorkResistanceProtocol{
        consensusSystem:         consensusSystem,
        ledgerInstance:          ledgerInstance,
        stateMutex:              stateMutex,
        resistanceRetryCount:    make(map[string]int),
        highHashRateEntityCount: make(map[string]int),
        powMonitoringCycleCount: 0,
    }
}

// StartPoWResistanceMonitoring starts the continuous loop for monitoring and enforcing PoW resistance
func (protocol *ProofOfWorkResistanceProtocol) StartPoWResistanceMonitoring() {
    ticker := time.NewTicker(PoWResistanceMonitoringInterval)

    go func() {
        for range ticker.C {
            protocol.monitorForResistanceIssues()
        }
    }()
}

// monitorForResistanceIssues monitors the PoW mining system for hash rate concentration, centralization, or other vulnerabilities
func (protocol *ProofOfWorkResistanceProtocol) monitorForResistanceIssues() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    // Fetch hash rate data and other mining activities from the consensus system
    miningActivities := protocol.consensusSystem.FetchMiningActivities()

    for _, activity := range miningActivities {
        if protocol.isHashRateConcentrationTooHigh(activity) {
            fmt.Printf("High hash rate concentration detected for entity %s. Taking action.\n", activity.EntityID)
            protocol.handleHighHashRateConcentration(activity)
        } else {
            fmt.Printf("Hash rate concentration within safe limits for entity %s.\n", activity.EntityID)
        }
    }

    protocol.powMonitoringCycleCount++
    fmt.Printf("PoW resistance monitoring cycle #%d completed.\n", protocol.powMonitoringCycleCount)

    if protocol.powMonitoringCycleCount%SubBlocksPerBlock == 0 {
        protocol.finalizeResistanceMonitoringCycle()
    }
}

// isHashRateConcentrationTooHigh checks if an entity's hash rate concentration exceeds the safe threshold
func (protocol *ProofOfWorkResistanceProtocol) isHashRateConcentrationTooHigh(activity common.MiningActivity) bool {
    // Check if the hash rate of the entity exceeds the defined threshold (70% for example)
    return activity.HashRatePercentage >= HashRateThreshold
}

// handleHighHashRateConcentration handles high hash rate concentration by either limiting hash power or applying penalties
func (protocol *ProofOfWorkResistanceProtocol) handleHighHashRateConcentration(activity common.MiningActivity) {
    protocol.highHashRateEntityCount[activity.EntityID]++

    if protocol.highHashRateEntityCount[activity.EntityID] >= MaxResistanceRetries {
        fmt.Printf("Multiple high hash rate incidents detected for entity %s. Applying hash power limit.\n", activity.EntityID)
        protocol.limitHashPower(activity)
    } else {
        fmt.Printf("Issuing warning to entity %s for high hash rate concentration.\n", activity.EntityID)
        protocol.warnEntityAboutHashRate(activity)
    }
}

// warnEntityAboutHashRate issues a warning to an entity with high hash rate concentration
func (protocol *ProofOfWorkResistanceProtocol) warnEntityAboutHashRate(activity common.MiningActivity) {
    encryptedWarningData := protocol.encryptResistanceData(activity)

    // Issue a warning through the Synnergy Consensus system
    warningSuccess := protocol.consensusSystem.WarnEntityAboutHashRate(activity.EntityID, encryptedWarningData)

    if warningSuccess {
        fmt.Printf("Warning issued to entity %s for high hash rate concentration.\n", activity.EntityID)
        protocol.logResistanceEvent(activity, "Warning Issued")
        protocol.resetResistanceRetry(activity.EntityID)
    } else {
        fmt.Printf("Error issuing warning to entity %s. Retrying...\n", activity.EntityID)
        protocol.retryResistanceAction(activity)
    }
}

// limitHashPower limits the hash rate of an entity to prevent over-concentration of mining power
func (protocol *ProofOfWorkResistanceProtocol) limitHashPower(activity common.MiningActivity) {
    encryptedLimitData := protocol.encryptResistanceData(activity)

    // Attempt to limit the entity's hash power through the Synnergy Consensus system
    limitSuccess := protocol.consensusSystem.LimitHashPower(activity.EntityID, encryptedLimitData)

    if limitSuccess {
        fmt.Printf("Hash power limit applied to entity %s.\n", activity.EntityID)
        protocol.logResistanceEvent(activity, "Hash Power Limited")
        protocol.resetResistanceRetry(activity.EntityID)
    } else {
        fmt.Printf("Error limiting hash power for entity %s. Retrying...\n", activity.EntityID)
        protocol.retryResistanceAction(activity)
    }
}

// retryResistanceAction retries the PoW resistance action if it initially fails
func (protocol *ProofOfWorkResistanceProtocol) retryResistanceAction(activity common.MiningActivity) {
    protocol.resistanceRetryCount[activity.EntityID]++
    if protocol.resistanceRetryCount[activity.EntityID] < MaxResistanceRetries {
        protocol.limitHashPower(activity)
    } else {
        fmt.Printf("Max retries reached for limiting hash power of entity %s. Action failed.\n", activity.EntityID)
        protocol.logResistanceFailure(activity)
    }
}

// resetResistanceRetry resets the retry count for resistance actions on a specific entity
func (protocol *ProofOfWorkResistanceProtocol) resetResistanceRetry(entityID string) {
    protocol.resistanceRetryCount[entityID] = 0
}

// finalizeResistanceMonitoringCycle finalizes the PoW resistance monitoring cycle and logs the result in the ledger
func (protocol *ProofOfWorkResistanceProtocol) finalizeResistanceMonitoringCycle() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    success := protocol.consensusSystem.FinalizeResistanceMonitoringCycle()
    if success {
        fmt.Println("PoW resistance monitoring cycle finalized successfully.")
        protocol.logResistanceCycleFinalization()
    } else {
        fmt.Println("Error finalizing PoW resistance monitoring cycle.")
    }
}

// logResistanceEvent logs a PoW resistance-related event into the ledger
func (protocol *ProofOfWorkResistanceProtocol) logResistanceEvent(activity common.MiningActivity, eventType string) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("pow-resistance-event-%s-%s", activity.EntityID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "PoW Resistance Event",
        Status:    eventType,
        Details:   fmt.Sprintf("Entity %s triggered %s due to high hash rate concentration.", activity.EntityID, eventType),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with PoW resistance event for entity %s.\n", activity.EntityID)
}

// logResistanceFailure logs the failure to enforce PoW resistance into the ledger
func (protocol *ProofOfWorkResistanceProtocol) logResistanceFailure(activity common.MiningActivity) {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("pow-resistance-failure-%s", activity.EntityID),
        Timestamp: time.Now().Unix(),
        Type:      "PoW Resistance Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("Failed to limit hash power of entity %s after maximum retries.", activity.EntityID),
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with PoW resistance failure for entity %s.\n", activity.EntityID)
}

// logResistanceCycleFinalization logs the finalization of a PoW resistance monitoring cycle into the ledger
func (protocol *ProofOfWorkResistanceProtocol) logResistanceCycleFinalization() {
    protocol.stateMutex.Lock()
    defer protocol.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("pow-resistance-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Resistance Cycle Finalization",
        Status:    "Finalized",
        Details:   "PoW resistance monitoring cycle finalized successfully.",
    }

    protocol.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with PoW resistance monitoring cycle finalization.")
}

// encryptResistanceData encrypts the data related to PoW resistance issues before applying actions or logging events
func (protocol *ProofOfWorkResistanceProtocol) encryptResistanceData(activity common.MiningActivity) common.MiningActivity {
    encryptedData, err := encryption.EncryptData(activity.ActivityData)
    if err != nil {
        fmt.Println("Error encrypting resistance data:", err)
        return activity
    }

    activity.EncryptedData = encryptedData
    fmt.Println("Resistance data successfully encrypted for entity ID:", activity.EntityID)
    return activity
}

// triggerEmergencyResistanceLockdown triggers an emergency PoW resistance lockdown in case of critical mining centralization
func (protocol *ProofOfWorkResistanceProtocol) triggerEmergencyResistanceLockdown(entityID string) {
    fmt.Printf("Emergency PoW resistance lockdown triggered for entity ID: %s.\n", entityID)
    activity := protocol.consensusSystem.GetMiningActivityByID(entityID)
    encryptedData := protocol.encryptResistanceData(activity)

    success := protocol.consensusSystem.TriggerEmergencyResistanceLockdown(entityID, encryptedData)

    if success {
        protocol.logResistanceEvent(activity, "Emergency Locked Down")
        fmt.Println("Emergency PoW resistance lockdown executed successfully.")
    } else {
        fmt.Println("Emergency PoW resistance lockdown failed.")
    }
}
