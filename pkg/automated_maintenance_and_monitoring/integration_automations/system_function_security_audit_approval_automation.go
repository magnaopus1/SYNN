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
    SecurityAuditCheckInterval = 3000 * time.Millisecond // Interval for checking system function security audits
    SubBlocksPerBlock          = 1000                    // Number of sub-blocks in a block
)

// SystemFunctionSecurityAuditApprovalAutomation automates the process of auditing and approving system functions based on security standards
type SystemFunctionSecurityAuditApprovalAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store security audit logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    auditCheckCount        int                          // Counter for security audit check cycles
}

// NewSystemFunctionSecurityAuditApprovalAutomation initializes the automation for security audits and approvals
func NewSystemFunctionSecurityAuditApprovalAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SystemFunctionSecurityAuditApprovalAutomation {
    return &SystemFunctionSecurityAuditApprovalAutomation{
        consensusSystem: consensusSystem,
        ledgerInstance:  ledgerInstance,
        stateMutex:      stateMutex,
        auditCheckCount: 0,
    }
}

// StartSecurityAuditMonitoring starts the continuous loop for monitoring and auditing system functions
func (automation *SystemFunctionSecurityAuditApprovalAutomation) StartSecurityAuditMonitoring() {
    ticker := time.NewTicker(SecurityAuditCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorAndAuditFunctions()
        }
    }()
}

// monitorAndAuditFunctions checks for system functions pending security audits and processes them
func (automation *SystemFunctionSecurityAuditApprovalAutomation) monitorAndAuditFunctions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of system functions that require security audits
    functions, err := automation.consensusSystem.GetPendingAuditFunctions()
    if err != nil {
        fmt.Printf("Error fetching functions for security audit: %v\n", err)
        return
    }

    // Process each function for audit and approval
    for _, function := range functions {
        fmt.Printf("Auditing system function: %s\n", function.FunctionID)

        // Encrypt function data for audit
        encryptedFunction, err := automation.encryptFunctionData(function)
        if err != nil {
            fmt.Printf("Error encrypting function data for %s: %v\n", function.FunctionID, err)
            automation.logAuditResult(function, "Encryption Failed")
            continue
        }

        // Perform security audit and approval process
        automation.auditAndApproveFunction(encryptedFunction)
    }

    automation.auditCheckCount++
    fmt.Printf("Security audit check cycle #%d completed.\n", automation.auditCheckCount)

    if automation.auditCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeAuditCycle()
    }
}

// encryptFunctionData encrypts the function data before performing a security audit
func (automation *SystemFunctionSecurityAuditApprovalAutomation) encryptFunctionData(function common.SystemFunction) (common.SystemFunction, error) {
    fmt.Println("Encrypting system function data for security audit.")

    encryptedData, err := encryption.EncryptData(function)
    if err != nil {
        return function, fmt.Errorf("failed to encrypt function data: %v", err)
    }

    function.EncryptedData = encryptedData
    fmt.Println("System function data successfully encrypted for security audit.")
    return function, nil
}

// auditAndApproveFunction performs the security audit and approves or rejects the function based on audit results
func (automation *SystemFunctionSecurityAuditApprovalAutomation) auditAndApproveFunction(function common.SystemFunction) {
    success := automation.consensusSystem.PerformSecurityAudit(function)
    if success {
        fmt.Printf("Function %s passed security audit and is approved.\n", function.FunctionID)
        automation.logAuditResult(function, "Approved")
    } else {
        fmt.Printf("Function %s failed security audit and is rejected.\n", function.FunctionID)
        automation.logAuditResult(function, "Rejected")
    }
}

// logAuditResult logs the audit result for a system function into the ledger for auditability
func (automation *SystemFunctionSecurityAuditApprovalAutomation) logAuditResult(function common.SystemFunction, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-audit-%s", function.FunctionID),
        Timestamp: time.Now().Unix(),
        Type:      "System Function Security Audit",
        Status:    status,
        Details:   fmt.Sprintf("Security audit result for function %s: %s", function.FunctionID, status),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with security audit event for function %s: %s\n", function.FunctionID, status)
}

// finalizeAuditCycle finalizes the security audit check cycle and logs the result in the ledger
func (automation *SystemFunctionSecurityAuditApprovalAutomation) finalizeAuditCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeSecurityAuditCycle()
    if success {
        fmt.Println("Security audit check cycle finalized successfully.")
        automation.logAuditCycleFinalization()
    } else {
        fmt.Println("Error finalizing security audit check cycle.")
    }
}

// logAuditCycleFinalization logs the finalization of the security audit cycle in the ledger
func (automation *SystemFunctionSecurityAuditApprovalAutomation) logAuditCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("security-audit-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Security Audit Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with security audit cycle finalization.")
}
