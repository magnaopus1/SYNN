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
)

const (
    AuthorizationCheckInterval = 5000 * time.Millisecond // Interval for checking system addition authorizations
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks per block
)

// SystemAdditionAuthorizationAutomation handles the authorization of new additions to the blockchain system
type SystemAdditionAuthorizationAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store authorization events
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    authorizationCheckCount int                      // Counter for authorization check cycles
}

// NewSystemAdditionAuthorizationAutomation initializes the automation for system addition authorization
func NewSystemAdditionAuthorizationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemAdditionAuthorizationAutomation {
    return &SystemAdditionAuthorizationAutomation{
        consensusSystem:        consensusSystem,
        ledgerInstance:         ledgerInstance,
        stateMutex:             stateMutex,
        authorizationCheckCount: 0,
    }
}

// StartAuthorizationCheck starts the continuous loop for monitoring and enforcing system addition authorization
func (automation *SystemAdditionAuthorizationAutomation) StartAuthorizationCheck() {
    ticker := time.NewTicker(AuthorizationCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndAuthorizeAdditions()
        }
    }()
}

// checkAndAuthorizeAdditions checks for pending system additions and authorizes them based on pre-set criteria
func (automation *SystemAdditionAuthorizationAutomation) checkAndAuthorizeAdditions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch pending additions for authorization
    pendingAdditions, err := automation.consensusSystem.GetPendingSystemAdditions()
    if err != nil {
        fmt.Printf("Error fetching pending system additions: %v\n", err)
        return
    }

    // Process each pending addition and authorize
    for _, addition := range pendingAdditions {
        fmt.Printf("Checking authorization for system addition: %s\n", addition.ID)

        // Encrypt addition data before authorization check
        encryptedAddition, err := automation.encryptAdditionData(addition)
        if err != nil {
            fmt.Printf("Error encrypting data for addition %s: %v\n", addition.ID, err)
            automation.logAuthorizationResult(addition, "Encryption Failed")
            continue
        }

        // Apply authorization checks and authorize if criteria are met
        authorized := automation.applyAuthorizationChecks(encryptedAddition)
        if authorized {
            fmt.Printf("System addition %s authorized successfully.\n", addition.ID)
            automation.logAuthorizationResult(addition, "Authorization Successful")
        } else {
            fmt.Printf("System addition %s failed authorization.\n", addition.ID)
            automation.logAuthorizationResult(addition, "Authorization Failed")
        }
    }

    automation.authorizationCheckCount++
    fmt.Printf("System addition authorization check cycle #%d completed.\n", automation.authorizationCheckCount)

    if automation.authorizationCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeAuthorizationCheckCycle()
    }
}

// encryptAdditionData encrypts the system addition data before performing authorization checks
func (automation *SystemAdditionAuthorizationAutomation) encryptAdditionData(addition common.SystemAddition) (common.SystemAddition, error) {
    fmt.Println("Encrypting system addition data.")

    encryptedData, err := encryption.EncryptData(addition)
    if err != nil {
        return addition, fmt.Errorf("failed to encrypt addition data: %v", err)
    }

    addition.EncryptedData = encryptedData
    fmt.Println("System addition data successfully encrypted.")
    return addition, nil
}

// applyAuthorizationChecks performs checks to determine if the addition meets all authorization criteria
func (automation *SystemAdditionAuthorizationAutomation) applyAuthorizationChecks(addition common.SystemAddition) bool {
    fmt.Printf("Checking if system addition %s meets authorization criteria.\n", addition.ID)

    // Example checks could include version, security compliance, etc.
    authorized := automation.consensusSystem.ValidateSystemAddition(addition)
    if authorized {
        fmt.Printf("System addition %s meets the criteria. Authorizing...\n", addition.ID)
        success := automation.consensusSystem.AuthorizeSystemAddition(addition)
        if success {
            fmt.Printf("System addition %s authorized.\n", addition.ID)
            return true
        } else {
            fmt.Printf("Error authorizing system addition %s.\n", addition.ID)
            return false
        }
    }

    return false
}

// logAuthorizationResult logs the result of a system addition authorization attempt in the ledger
func (automation *SystemAdditionAuthorizationAutomation) logAuthorizationResult(addition common.SystemAddition, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("system-addition-%s", addition.ID),
        Timestamp: time.Now().Unix(),
        Type:      "System Addition Authorization",
        Status:    result,
        Details:   fmt.Sprintf("Authorization result for system addition %s: %s", addition.ID, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with authorization result for system addition %s: %s\n", addition.ID, result)
}

// finalizeAuthorizationCheckCycle finalizes the authorization check cycle and logs the results
func (automation *SystemAdditionAuthorizationAutomation) finalizeAuthorizationCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeAuthorizationCheckCycle()
    if success {
        fmt.Println("System addition authorization check cycle finalized successfully.")
        automation.logAuthorizationCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing system addition authorization check cycle.")
    }
}

// logAuthorizationCheckCycleFinalization logs the finalization of the authorization check cycle in the ledger
func (automation *SystemAdditionAuthorizationAutomation) logAuthorizationCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("authorization-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "System Addition Authorization Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with system addition authorization check cycle finalization.")
}
