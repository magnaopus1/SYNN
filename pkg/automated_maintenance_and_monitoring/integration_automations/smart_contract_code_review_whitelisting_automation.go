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
    "synnergy_network_demo/smartcontracts"
)

const (
    CodeReviewCheckInterval    = 8000 * time.Millisecond // Interval for reviewing and whitelisting contracts
    SubBlocksPerBlock           = 1000                   // Number of sub-blocks per block
    MaxContractFailureThreshold = 5                      // Maximum number of contract review failures
)

// SmartContractCodeReviewAutomation automates the process of reviewing and whitelisting smart contracts
type SmartContractCodeReviewAutomation struct {
    consensusSystem      *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance       *ledger.Ledger               // Ledger to store code review actions
    stateMutex           *sync.RWMutex                // Mutex for thread-safe access
    reviewCheckCount     int                          // Counter for contract review check cycles
    failedReviewCount    int                          // Counter for failed smart contract reviews
}

// NewSmartContractCodeReviewAutomation initializes the automation for smart contract code review and whitelisting
func NewSmartContractCodeReviewAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *SmartContractCodeReviewAutomation {
    return &SmartContractCodeReviewAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        reviewCheckCount: 0,
        failedReviewCount: 0,
    }
}

// StartContractCodeReview starts the continuous loop for reviewing and whitelisting smart contracts
func (automation *SmartContractCodeReviewAutomation) StartContractCodeReview() {
    ticker := time.NewTicker(CodeReviewCheckInterval)

    go func() {
        for range ticker.C {
            automation.reviewAndWhitelistContracts()
        }
    }()
}

// reviewAndWhitelistContracts reviews new smart contracts, validates them, and adds them to the whitelist if they pass
func (automation *SmartContractCodeReviewAutomation) reviewAndWhitelistContracts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Step 1: Fetch new smart contracts for review
    contractsToReview, err := automation.consensusSystem.GetSmartContractsForReview()
    if err != nil {
        fmt.Printf("Error fetching contracts for review: %v\n", err)
        return
    }

    // Step 2: Process each contract, review it, and validate it
    for _, contract := range contractsToReview {
        fmt.Printf("Reviewing contract: %s\n", contract.ID)
        
        // Step 3: Encrypt contract data before validation
        encryptedContract, err := automation.encryptContractData(contract)
        if err != nil {
            fmt.Printf("Error encrypting contract: %s - %v\n", contract.ID, err)
            automation.failedReviewCount++
            if automation.failedReviewCount >= MaxContractFailureThreshold {
                automation.triggerReviewMitigation("Contract Review Failure Threshold Exceeded")
            }
            continue
        }

        // Step 4: Validate contract code for security and correctness
        reviewPassed := automation.validateContract(encryptedContract)
        if reviewPassed {
            fmt.Printf("Contract %s passed review. Whitelisting.\n", contract.ID)
            automation.whitelistContract(contract)
        } else {
            fmt.Printf("Contract %s failed review.\n", contract.ID)
            automation.failedReviewCount++
            automation.logReviewResult(contract, "Review Failed")
            if automation.failedReviewCount >= MaxContractFailureThreshold {
                automation.triggerReviewMitigation("Multiple Contract Review Failures")
            }
        }
    }

    automation.reviewCheckCount++
    fmt.Printf("Smart contract code review cycle #%d completed.\n", automation.reviewCheckCount)

    if automation.reviewCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeReviewCycle()
    }
}

// encryptContractData encrypts the smart contract data before performing code review
func (automation *SmartContractCodeReviewAutomation) encryptContractData(contract common.SmartContract) (common.SmartContract, error) {
    fmt.Println("Encrypting smart contract data.")

    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        return contract, fmt.Errorf("failed to encrypt contract data: %v", err)
    }

    contract.EncryptedData = encryptedData
    fmt.Println("Contract data successfully encrypted.")
    return contract, nil
}

// validateContract validates the smart contract for security, logic, and vulnerabilities
func (automation *SmartContractCodeReviewAutomation) validateContract(contract common.SmartContract) bool {
    fmt.Printf("Validating contract %s.\n", contract.ID)

    reviewPassed := automation.consensusSystem.ValidateSmartContract(contract)
    if !reviewPassed {
        fmt.Printf("Validation failed for contract %s.\n", contract.ID)
        return false
    }

    fmt.Printf("Contract %s successfully passed validation.\n", contract.ID)
    return true
}

// whitelistContract adds the smart contract to the system whitelist after passing review
func (automation *SmartContractCodeReviewAutomation) whitelistContract(contract common.SmartContract) {
    success := automation.consensusSystem.WhitelistSmartContract(contract)
    if success {
        fmt.Printf("Contract %s successfully whitelisted.\n", contract.ID)
        automation.logReviewResult(contract, "Whitelisted")
    } else {
        fmt.Printf("Failed to whitelist contract %s.\n", contract.ID)
    }
}

// logReviewResult logs the result of a smart contract review in the ledger
func (automation *SmartContractCodeReviewAutomation) logReviewResult(contract common.SmartContract, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("contract-review-%s", contract.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Review",
        Status:    result,
        Details:   fmt.Sprintf("Review result for contract %s: %s", contract.ID, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with contract review result for contract %s: %s\n", contract.ID, result)
}

// triggerReviewMitigation triggers an action when multiple contract reviews fail
func (automation *SmartContractCodeReviewAutomation) triggerReviewMitigation(reason string) {
    fmt.Printf("Triggering mitigation due to: %s\n", reason)

    success := automation.consensusSystem.TriggerMitigation(reason)
    if success {
        fmt.Printf("Mitigation action successfully triggered: %s\n", reason)
        automation.logMitigationAction(reason, "Mitigation Triggered")
    } else {
        fmt.Printf("Failed to trigger mitigation action: %s\n", reason)
        automation.logMitigationAction(reason, "Mitigation Failed")
    }
}

// logMitigationAction logs any mitigation action taken due to failed reviews
func (automation *SmartContractCodeReviewAutomation) logMitigationAction(reason, status string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("mitigation-action-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Mitigation Action",
        Status:    status,
        Details:   fmt.Sprintf("Mitigation action due to: %s", reason),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with mitigation action: %s - %s\n", reason, status)
}

// finalizeReviewCycle finalizes the smart contract code review cycle and logs the results in the ledger
func (automation *SmartContractCodeReviewAutomation) finalizeReviewCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeContractReviewCycle()
    if success {
        fmt.Println("Smart contract review cycle finalized successfully.")
        automation.logReviewCycleFinalization()
    } else {
        fmt.Println("Error finalizing smart contract review cycle.")
    }
}

// logReviewCycleFinalization logs the finalization of a smart contract review cycle in the ledger
func (automation *SmartContractCodeReviewAutomation) logReviewCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("review-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Smart Contract Review Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with smart contract review cycle finalization.")
}
