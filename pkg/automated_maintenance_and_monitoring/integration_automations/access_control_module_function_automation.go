package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/access_control"
)

const (
    AccessControlCheckInterval = 1000 * time.Millisecond // Interval for checking access control
)

// AccessControlModuleFunctionAutomation automates the access control enforcement for new modules and functions
type AccessControlModuleFunctionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger to store access control decisions
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
}

// NewAccessControlModuleFunctionAutomation initializes the automation for enforcing access control
func NewAccessControlModuleFunctionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AccessControlModuleFunctionAutomation {
    return &AccessControlModuleFunctionAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
    }
}

// StartAccessControlEnforcement starts the continuous loop for monitoring and enforcing access control
func (automation *AccessControlModuleFunctionAutomation) StartAccessControlEnforcement() {
    ticker := time.NewTicker(AccessControlCheckInterval)

    go func() {
        for range ticker.C {
            automation.enforceAccessControl()
        }
    }()
}

// enforceAccessControl checks access control for new modules and functions and logs the decisions
func (automation *AccessControlModuleFunctionAutomation) enforceAccessControl() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch new system functions and modules to enforce access control on
    newFunctions := access_control.GetNewFunctions()
    newModules := access_control.GetNewModules()

    // Enforce access control for new functions
    for _, function := range newFunctions {
        if automation.isAuthorized(function) {
            fmt.Printf("Function %s is authorized and integrated.\n", function.Name)
            automation.logAccessControlDecision(function.Name, "Authorized")
        } else {
            fmt.Printf("Function %s is unauthorized and denied access.\n", function.Name)
            automation.logAccessControlDecision(function.Name, "Unauthorized")
            continue
        }
        automation.integrateNewFunction(function)
    }

    // Enforce access control for new modules
    for _, module := range newModules {
        if automation.isAuthorized(module) {
            fmt.Printf("Module %s is authorized and integrated.\n", module.Name)
            automation.logAccessControlDecision(module.Name, "Authorized")
        } else {
            fmt.Printf("Module %s is unauthorized and denied access.\n", module.Name)
            automation.logAccessControlDecision(module.Name, "Unauthorized")
            continue
        }
        automation.integrateNewModule(module)
    }
}

// isAuthorized checks if the given function or module is authorized for integration
func (automation *AccessControlModuleFunctionAutomation) isAuthorized(entity access_control.Entity) bool {
    authorized := access_control.CheckAuthorization(entity)
    return authorized
}

// logAccessControlDecision logs the access control decision into the ledger
func (automation *AccessControlModuleFunctionAutomation) logAccessControlDecision(name string, status string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("access-control-%s", name),
        Timestamp: time.Now().Unix(),
        Type:      "Access Control",
        Status:    status,
        Details:   fmt.Sprintf("Access control decision: %s for %s", status, name),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with access control decision for %s.\n", name)
}

// integrateNewFunction integrates the authorized function into the system
func (automation *AccessControlModuleFunctionAutomation) integrateNewFunction(function access_control.Entity) {
    success := automation.consensusSystem.RegisterFunction(function)
    if success {
        fmt.Printf("Function %s successfully integrated into the system.\n", function.Name)
    } else {
        fmt.Printf("Error integrating function %s into the system.\n", function.Name)
    }
}

// integrateNewModule integrates the authorized module into the system
func (automation *AccessControlModuleFunctionAutomation) integrateNewModule(module access_control.Entity) {
    success := automation.consensusSystem.RegisterModule(module)
    if success {
        fmt.Printf("Module %s successfully integrated into the system.\n", module.Name)
    } else {
        fmt.Printf("Error integrating module %s into the system.\n", module.Name)
    }
}
