package automations

import (
    "fmt"
    "log"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/security"
)

const (
    ContractCheckInterval      = 3000 * time.Millisecond // Interval for checking contracts for malicious behavior
    SubBlocksPerBlock           = 1000                   // Number of sub-blocks in a block
    MaxAllowedSuspiciousScore   = 50                     // Maximum suspicious score allowed before contract is flagged
    EncryptionErrorLogThreshold = 5                      // Number of encryption errors allowed before stopping execution
)

// MaliciousContractDetectionAutomation automates the detection and prevention of malicious contracts
type MaliciousContractDetectionAutomation struct {
    consensusSystem     *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance      *ledger.Ledger               // Ledger to store detection/prevention logs
    stateMutex          *sync.RWMutex                // Mutex for thread-safe access
    contractCheckCount  int                          // Counter for contract check cycles
    encryptionErrorCount int                         // Tracks the number of encryption errors encountered
}

// NewMaliciousContractDetectionAutomation initializes the automation for malicious contract detection and prevention
func NewMaliciousContractDetectionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *MaliciousContractDetectionAutomation {
    return &MaliciousContractDetectionAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        contractCheckCount: 0,
        encryptionErrorCount: 0,
    }
}

// StartContractDetection starts the continuous loop for scanning and preventing malicious contracts
func (automation *MaliciousContractDetectionAutomation) StartContractDetection() {
    ticker := time.NewTicker(ContractCheckInterval)

    go func() {
        for range ticker.C {
            automation.scanAndPreventMaliciousContracts()
        }
    }()
}

// scanAndPreventMaliciousContracts checks for contracts with suspicious activity and prevents malicious contracts from being executed
func (automation *MaliciousContractDetectionAutomation) scanAndPreventMaliciousContracts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    newContracts := automation.consensusSystem.GetNewContracts() // Fetch newly deployed contracts

    for _, contract := range newContracts {
        fmt.Printf("Checking contract %s for malicious behavior.\n", contract.Address)
        
        suspiciousScore, err := automation.analyzeContractBehavior(contract)
        if err != nil {
            log.Printf("Error analyzing contract behavior for %s: %v", contract.Address, err)
            continue
        }

        if suspiciousScore > MaxAllowedSuspiciousScore {
            fmt.Printf("Contract %s flagged as malicious. Taking preventive action.\n", contract.Address)
            automation.preventMaliciousContractExecution(contract)
        } else {
            fmt.Printf("Contract %s is safe. No malicious activity detected.\n", contract.Address)
            automation.logContractScanResult(contract.Address, "Safe", suspiciousScore)
        }
    }

    automation.contractCheckCount++
    fmt.Printf("Contract check cycle #%d completed.\n", automation.contractCheckCount)

    if automation.contractCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeContractCheckCycle()
    }
}

// analyzeContractBehavior evaluates the behavior of the contract and returns a suspicious score
func (automation *MaliciousContractDetectionAutomation) analyzeContractBehavior(contract common.Contract) (int, error) {
    // Step 1: Encrypt contract data for security
    fmt.Printf("Encrypting contract data for %s.\n", contract.Address)
    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        automation.encryptionErrorCount++
        if automation.encryptionErrorCount >= EncryptionErrorLogThreshold {
            return 0, fmt.Errorf("encryption error threshold exceeded for contract %s", contract.Address)
        }
        return 0, fmt.Errorf("error encrypting contract data for %s: %v", contract.Address, err)
    }
    contract.EncryptedData = encryptedData
    fmt.Printf("Contract data for %s encrypted successfully.\n", contract.Address)

    // Step 2: Analyze the contract for suspicious patterns or behavior using the security package
    return security.EvaluateContractForMaliciousBehavior(contract)
}

// preventMaliciousContractExecution blocks the execution of a flagged malicious contract and logs the result
func (automation *MaliciousContractDetectionAutomation) preventMaliciousContractExecution(contract common.Contract) {
    // Prevent contract execution through the Synnergy Consensus
    success := automation.consensusSystem.BlockContractExecution(contract)
    if success {
        fmt.Printf("Execution of malicious contract %s successfully blocked.\n", contract.Address)
        automation.logContractScanResult(contract.Address, "Blocked", -1)
    } else {
        fmt.Printf("Error blocking execution of contract %s.\n", contract.Address)
    }
}

// logContractScanResult logs the result of the contract scan and action into the ledger
func (automation *MaliciousContractDetectionAutomation) logContractScanResult(contractAddress, result string, score int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    details := fmt.Sprintf("Contract scan result for %s: %s", contractAddress, result)
    if score >= 0 {
        details += fmt.Sprintf(" (Suspicious Score: %d)", score)
    }

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-scan-%s", contractAddress),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Scan",
        Status:    result,
        Details:   details,
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract scan result for %s: %s.\n", contractAddress, result)
}

// finalizeContractCheckCycle finalizes the contract check cycle and logs the result in the ledger
func (automation *MaliciousContractDetectionAutomation) finalizeContractCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeContractCheckCycle()
    if success {
        fmt.Println("Contract check cycle finalized successfully.")
        automation.logContractCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing contract check cycle.")
    }
}

// logContractCheckCycleFinalization logs the finalization of a contract check cycle into the ledger
func (automation *MaliciousContractDetectionAutomation) logContractCheckCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Contract Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with contract check cycle finalization.")
}
