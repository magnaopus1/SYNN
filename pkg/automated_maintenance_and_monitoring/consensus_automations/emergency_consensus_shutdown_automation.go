package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    EmergencyCheckInterval     = 5 * time.Second  // Interval for checking catastrophic failures
    EmergencyShutdownKey       = "emergency_shutdown_key" // Encryption key for emergency shutdown logs
    EmergencyShutdownThreshold = 0.9             // Threshold for triggering emergency shutdown
    ResumeCheckInterval        = 10 * time.Second // Interval for checking if the system can resume
    MaintenanceResumeKey       = "maintenance_resume_key" // Encryption key for maintenance resume logs
    ForcedResumeKey            = "forced_resume_key" // Encryption key for forced resume logs
)

// EmergencyConsensusShutdownAutomation handles catastrophic failure detection, consensus shutdown, and resume operations.
type EmergencyConsensusShutdownAutomation struct {
    ledgerInstance   *ledger.Ledger                    // Blockchain ledger for tracking shutdown events
    consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for monitoring and shutdown
    stateMutex       *sync.RWMutex                     // Mutex for thread-safe ledger access
    isShutDown       bool                              // Tracks whether the system is currently shut down
    forceResume      bool                              // Tracks whether a forced resume is triggered
    lastHealthCheck  time.Time                         // Stores the time of the last health check
}

// NewEmergencyConsensusShutdownAutomation initializes the emergency shutdown and resume automation
func NewEmergencyConsensusShutdownAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *EmergencyConsensusShutdownAutomation {
    return &EmergencyConsensusShutdownAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        isShutDown:      false,
        forceResume:     false,
    }
}

// StartEmergencyMonitoring begins monitoring the network for catastrophic failures and checks if the system can resume
func (automation *EmergencyConsensusShutdownAutomation) StartEmergencyMonitoring() {
    ticker := time.NewTicker(EmergencyCheckInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Monitoring network for catastrophic failures...")
            automation.monitorForCatastrophicFailures()

            if automation.isShutDown {
                fmt.Println("System is shut down. Checking if it can resume...")
                automation.checkForResume()
            }
        }
    }()
}

// monitorForCatastrophicFailures checks for network issues that could trigger an emergency shutdown
func (automation *EmergencyConsensusShutdownAutomation) monitorForCatastrophicFailures() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if automation.isShutDown {
        return
    }

    // Directly access the consensus engine to check for severe issues
    networkHealth := automation.consensusEngine.CheckNetworkHealth()

    if networkHealth < EmergencyShutdownThreshold {
        fmt.Println("Severe network failure detected, initiating emergency shutdown.")
        automation.initiateEmergencyShutdown()
    } else {
        fmt.Println("Network health is stable.")
    }
}

// initiateEmergencyShutdown halts PoH, PoS, PoW operations and freezes the consensus
func (automation *EmergencyConsensusShutdownAutomation) initiateEmergencyShutdown() {
    automation.isShutDown = true
    fmt.Println("Initiating emergency shutdown...")

    // Stop PoH operations
    err := automation.consensusEngine.PoH.HaltPoH()
    if err != nil {
        fmt.Printf("Error halting PoH: %v\n", err)
    } else {
        fmt.Println("PoH operations halted successfully.")
    }

    // Stop PoS operations
    err = automation.consensusEngine.PoS.HaltPoS()
    if err != nil {
        fmt.Printf("Error halting PoS: %v\n", err)
    } else {
        fmt.Println("PoS operations halted successfully.")
    }

    // Stop PoW operations
    err = automation.consensusEngine.PoW.HaltPoW()
    if err != nil {
        fmt.Printf("Error halting PoW: %v\n", err)
    } else {
        fmt.Println("PoW operations halted successfully.")
    }

    // Freeze the ledger to prevent further transactions
    err = automation.ledgerInstance.Freeze()
    if err != nil {
        fmt.Printf("Error freezing ledger: %v\n", err)
    } else {
        fmt.Println("Ledger frozen successfully.")
    }

    // Log the emergency shutdown event in the ledger
    automation.logEmergencyShutdown()

    fmt.Println("Emergency shutdown completed. All consensus operations halted.")
}

// logEmergencyShutdown logs the shutdown event into the ledger securely
func (automation *EmergencyConsensusShutdownAutomation) logEmergencyShutdown() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    shutdownLog := common.LedgerEntry{
        ID:        fmt.Sprintf("emergency-shutdown-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Emergency Shutdown",
        Status:    "Completed",
        Details:   "Consensus was shut down due to a catastrophic failure.",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(shutdownLog, []byte(EmergencyShutdownKey))
    if err != nil {
        fmt.Printf("Error encrypting emergency shutdown log: %v\n", err)
        return
    }

    err = automation.ledgerInstance.AddEntry(encryptedEntry)
    if err != nil {
        fmt.Printf("Error storing the shutdown log: %v\n", err)
    } else {
        fmt.Println("Emergency shutdown log stored in the ledger.")
    }
}

// checkForResume checks whether the system can resume based on network health or triggers a forced resume
func (automation *EmergencyConsensusShutdownAutomation) checkForResume() {
    if automation.forceResume {
        automation.resumeConsensusOperations()
        return
    }

    // Check the system health to see if the issue is resolved
    networkHealth := automation.consensusEngine.CheckNetworkHealth()

    if networkHealth >= EmergencyShutdownThreshold {
        fmt.Println("Network health restored. Resuming consensus operations.")
        automation.resumeConsensusOperations()
    } else {
        fmt.Println("Network health is still below the threshold. Cannot resume yet.")
    }
}

// resumeConsensusOperations resumes PoH, PoS, and PoW operations after a shutdown
func (automation *EmergencyConsensusShutdownAutomation) resumeConsensusOperations() {
    fmt.Println("Resuming consensus operations...")

    // Resume PoH operations
    err := automation.consensusEngine.PoH.ResumePoH()
    if err != nil {
        fmt.Printf("Error resuming PoH: %v\n", err)
    } else {
        fmt.Println("PoH operations resumed successfully.")
    }

    // Resume PoS operations
    err = automation.consensusEngine.PoS.ResumePoS()
    if err != nil {
        fmt.Printf("Error resuming PoS: %v\n", err)
    } else {
        fmt.Println("PoS operations resumed successfully.")
    }

    // Resume PoW operations
    err = automation.consensusEngine.PoW.ResumePoW()
    if err != nil {
        fmt.Printf("Error resuming PoW: %v\n", err)
    } else {
        fmt.Println("PoW operations resumed successfully.")
    }

    // Unfreeze the ledger
    err = automation.ledgerInstance.Unfreeze()
    if err != nil {
        fmt.Printf("Error unfreezing the ledger: %v\n", err)
    } else {
        fmt.Println("Ledger unfrozen successfully.")
    }

    // Log the resume event
    automation.logResume()

    automation.isShutDown = false
    automation.forceResume = false

    fmt.Println("Consensus operations resumed successfully.")
}

// logResume logs the resume event into the ledger securely
func (automation *EmergencyConsensusShutdownAutomation) logResume() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    resumeLog := common.LedgerEntry{
        ID:        fmt.Sprintf("resume-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Consensus Resume",
        Status:    "Completed",
        Details:   "Consensus resumed after maintenance or health check.",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(resumeLog, []byte(MaintenanceResumeKey))
    if err != nil {
        fmt.Printf("Error encrypting resume log: %v\n", err)
        return
    }

    err = automation.ledgerInstance.AddEntry(encryptedEntry)
    if err != nil {
        fmt.Printf("Error storing the resume log: %v\n", err)
    } else {
        fmt.Println("Resume log stored in the ledger.")
    }
}

// ForceResume allows administrators to force the system to resume operations
func (automation *EmergencyConsensusShutdownAutomation) ForceResume() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    fmt.Println("Force resuming consensus operations...")
    automation.forceResume = true

    // Log the forced resume event
    automation.logForcedResume()
}

// logForcedResume logs the forced resume event into the ledger securely
func (automation *EmergencyConsensusShutdownAutomation) logForcedResume() {
    forceResumeLog := common.LedgerEntry{
        ID:        fmt.Sprintf("forced-resume-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Forced Resume",
        Status:    "Completed",
        Details:   "Consensus was resumed by a forced override.",
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(forceResumeLog, []byte(ForcedResumeKey))
    if err != nil {
        fmt.Printf("Error encrypting forced resume log: %v\n", err)
        return
    }

    err = automation.ledgerInstance.AddEntry(encryptedEntry)
    if err != nil {
        fmt.Printf("Error storing the forced resume log: %v\n", err)
    } else {
        fmt.Println("Forced resume log stored in the ledger.")
    }
}
