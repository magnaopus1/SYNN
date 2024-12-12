package automations

import (
    "fmt"
    "log"
    "sync"
    "time"
    "errors"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
)

const (
    ProductContractCheckInterval = 6000 * time.Millisecond // Interval for checking new product contracts
    SubBlocksPerBlock            = 1000                    // Number of sub-blocks in a block
    EncryptionErrorThreshold     = 5                       // Number of encryption errors allowed before halting the process
)

// NewProductContractDetectionAutomation automates the detection of new product contracts
type NewProductContractDetectionAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store contract integration logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    contractCheckCount     int                          // Counter for contract check cycles
    encryptionErrorCount   int                          // Tracks the number of encryption errors
}

// NewNewProductContractDetectionAutomation initializes the automation for detecting new product contracts
func NewNewProductContractDetectionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *NewProductContractDetectionAutomation {
    return &NewProductContractDetectionAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        contractCheckCount:   0,
        encryptionErrorCount: 0,
    }
}

// StartProductContractCheck starts the continuous loop for checking and integrating new product contracts
func (automation *NewProductContractDetectionAutomation) StartProductContractCheck() {
    ticker := time.NewTicker(ProductContractCheckInterval)

    go func() {
        for range ticker.C {
            automation.detectAndIntegrateProductContracts()
        }
    }()
}

// detectAndIntegrateProductContracts checks for newly deployed product contracts and integrates them into the system
func (automation *NewProductContractDetectionAutomation) detectAndIntegrateProductContracts() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Step 1: Fetch newly deployed product contracts
    newContracts := automation.consensusSystem.GetNewProductContracts()

    for _, contract := range newContracts {
        fmt.Printf("Detecting new product contract: %s\n", contract.Address)

        // Step 2: Encrypt contract data
        encryptedContract, err := automation.encryptContractData(contract)
        if err != nil {
            log.Printf("Error encrypting product contract %s: %v", contract.Address, err)
            automation.encryptionErrorCount++
            if automation.encryptionErrorCount >= EncryptionErrorThreshold {
                fmt.Println("Encryption error threshold exceeded. Halting process.")
                return
            }
            continue
        }

        // Step 3: Validate the contract with the Synnergy Consensus
        validationResult, err := automation.validateProductContract(encryptedContract)
        if err != nil || !validationResult {
            fmt.Printf("Contract %s failed validation: %v\n", contract.Address, err)
            automation.logContractResult(contract.Address, "Failed Validation")
            continue
        }

        // Step 4: Integrate the contract into the system
        integrationSuccess := automation.integrateProductContract(encryptedContract)
        if integrationSuccess {
            fmt.Printf("Product contract %s successfully integrated into the system.\n", contract.Address)
            automation.logContractResult(contract.Address, "Integrated")
        } else {
            fmt.Printf("Failed to integrate product contract %s.\n", contract.Address)
            automation.logContractResult(contract.Address, "Integration Failed")
        }
    }

    automation.contractCheckCount++
    fmt.Printf("Product contract check cycle #%d completed.\n", automation.contractCheckCount)

    if automation.contractCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeContractCheckCycle()
    }
}

// encryptContractData encrypts the product contract data before validation and integration
func (automation *NewProductContractDetectionAutomation) encryptContractData(contract common.Contract) (common.Contract, error) {
    fmt.Printf("Encrypting data for product contract: %s\n", contract.Address)
    
    encryptedData, err := encryption.EncryptData(contract)
    if err != nil {
        return contract, fmt.Errorf("failed to encrypt contract data for %s: %v", contract.Address, err)
    }

    contract.EncryptedData = encryptedData
    fmt.Printf("Contract data for %s encrypted successfully.\n", contract.Address)
    return contract, nil
}

// validateProductContract validates the product contract with the Synnergy Consensus system
func (automation *NewProductContractDetectionAutomation) validateProductContract(contract common.Contract) (bool, error) {
    fmt.Printf("Validating product contract %s with Synnergy Consensus.\n", contract.Address)

    validationResult, err := automation.consensusSystem.ValidateContract(contract)
    if err != nil {
        return false, fmt.Errorf("error during validation for contract %s: %v", contract.Address, err)
    }

    return validationResult, nil
}

// integrateProductContract integrates the validated product contract into the system
func (automation *NewProductContractDetectionAutomation) integrateProductContract(contract common.Contract) bool {
    fmt.Printf("Integrating product contract %s into the system.\n", contract.Address)

    integrationSuccess := automation.consensusSystem.IntegrateProductContract(contract)
    if integrationSuccess {
        return true
    }
    return false
}

// logContractResult logs the result of the contract detection and integration into the ledger
func (automation *NewProductContractDetectionAutomation) logContractResult(contractAddress, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("product-contract-%s", contractAddress),
        Timestamp: time.Now().Unix(),
        Type:      "Product Contract Integration",
        Status:    result,
        Details:   fmt.Sprintf("Result for product contract %s: %s", contractAddress, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with result for product contract %s: %s.\n", contractAddress, result)
}

// finalizeContractCheckCycle finalizes the contract check cycle and logs the result in the ledger
func (automation *NewProductContractDetectionAutomation) finalizeContractCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeProductContractCheckCycle()
    if success {
        fmt.Println("Product contract check cycle finalized successfully.")
        automation.logContractCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing product contract check cycle.")
    }
}

// logContractCheckCycleFinalization logs the finalization of a product contract check cycle into the ledger
func (automation *NewProductContractDetectionAutomation) logContractCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("product-contract-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Product Contract Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with product contract check cycle finalization.")
}
