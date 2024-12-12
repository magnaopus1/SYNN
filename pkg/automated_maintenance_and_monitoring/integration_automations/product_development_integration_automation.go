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
)

const (
    ProductDevelopmentCheckInterval = 5000 * time.Millisecond // Interval for checking new product developments
    SubBlocksPerBlock               = 1000                    // Number of sub-blocks in a block
    EncryptionErrorThreshold        = 5                       // Number of encryption errors allowed before halting the process
)

// ProductDevelopmentIntegrationAutomation automates the detection, validation, and integration of new product developments
type ProductDevelopmentIntegrationAutomation struct {
    consensusSystem        *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance         *ledger.Ledger               // Ledger to store product development logs
    stateMutex             *sync.RWMutex                // Mutex for thread-safe access
    productCheckCount      int                          // Counter for product development check cycles
    encryptionErrorCount   int                          // Tracks the number of encryption errors encountered
}

// NewProductDevelopmentIntegrationAutomation initializes the automation for product development integration
func NewProductDevelopmentIntegrationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ProductDevelopmentIntegrationAutomation {
    return &ProductDevelopmentIntegrationAutomation{
        consensusSystem:      consensusSystem,
        ledgerInstance:       ledgerInstance,
        stateMutex:           stateMutex,
        productCheckCount:    0,
        encryptionErrorCount: 0,
    }
}

// StartProductDevelopmentCheck starts the continuous loop for checking and integrating new product developments
func (automation *ProductDevelopmentIntegrationAutomation) StartProductDevelopmentCheck() {
    ticker := time.NewTicker(ProductDevelopmentCheckInterval)

    go func() {
        for range ticker.C {
            automation.detectAndIntegrateProductDevelopments()
        }
    }()
}

// detectAndIntegrateProductDevelopments checks for newly developed products and integrates them into the system
func (automation *ProductDevelopmentIntegrationAutomation) detectAndIntegrateProductDevelopments() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Step 1: Fetch newly developed products
    newProducts := automation.consensusSystem.GetNewProductDevelopments()

    for _, product := range newProducts {
        fmt.Printf("Detecting new product development: %s\n", product.Name)

        // Step 2: Encrypt product data
        encryptedProduct, err := automation.encryptProductData(product)
        if err != nil {
            log.Printf("Error encrypting product development %s: %v", product.Name, err)
            automation.encryptionErrorCount++
            if automation.encryptionErrorCount >= EncryptionErrorThreshold {
                fmt.Println("Encryption error threshold exceeded. Halting process.")
                return
            }
            continue
        }

        // Step 3: Validate the product development through the Synnergy Consensus
        validationResult, err := automation.validateProductDevelopment(encryptedProduct)
        if err != nil || !validationResult {
            fmt.Printf("Product development %s failed validation: %v\n", product.Name, err)
            automation.logProductDevelopmentResult(product.Name, "Failed Validation")
            continue
        }

        // Step 4: Integrate the validated product into the system
        integrationSuccess := automation.integrateProductDevelopment(encryptedProduct)
        if integrationSuccess {
            fmt.Printf("Product development %s successfully integrated into the system.\n", product.Name)
            automation.logProductDevelopmentResult(product.Name, "Integrated")
        } else {
            fmt.Printf("Failed to integrate product development %s.\n", product.Name)
            automation.logProductDevelopmentResult(product.Name, "Integration Failed")
        }
    }

    automation.productCheckCount++
    fmt.Printf("Product development check cycle #%d completed.\n", automation.productCheckCount)

    if automation.productCheckCount%SubBlocksPerBlock == 0 {
        automation.finalizeProductDevelopmentCheckCycle()
    }
}

// encryptProductData encrypts the product data before validation and integration
func (automation *ProductDevelopmentIntegrationAutomation) encryptProductData(product common.Product) (common.Product, error) {
    fmt.Printf("Encrypting data for product development: %s\n", product.Name)
    
    encryptedData, err := encryption.EncryptData(product)
    if err != nil {
        return product, fmt.Errorf("failed to encrypt product data for %s: %v", product.Name, err)
    }

    product.EncryptedData = encryptedData
    fmt.Printf("Product data for %s encrypted successfully.\n", product.Name)
    return product, nil
}

// validateProductDevelopment validates the product development through the Synnergy Consensus system
func (automation *ProductDevelopmentIntegrationAutomation) validateProductDevelopment(product common.Product) (bool, error) {
    fmt.Printf("Validating product development %s with Synnergy Consensus.\n", product.Name)

    validationResult, err := automation.consensusSystem.ValidateProductDevelopment(product)
    if err != nil {
        return false, fmt.Errorf("error during validation for product development %s: %v", product.Name, err)
    }

    return validationResult, nil
}

// integrateProductDevelopment integrates the validated product into the system
func (automation *ProductDevelopmentIntegrationAutomation) integrateProductDevelopment(product common.Product) bool {
    fmt.Printf("Integrating product development %s into the system.\n", product.Name)

    integrationSuccess := automation.consensusSystem.IntegrateProductDevelopment(product)
    if integrationSuccess {
        return true
    }
    return false
}

// logProductDevelopmentResult logs the result of the product development detection and integration into the ledger
func (automation *ProductDevelopmentIntegrationAutomation) logProductDevelopmentResult(productName, result string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("product-development-%s", productName),
        Timestamp: time.Now().Unix(),
        Type:      "Product Development Integration",
        Status:    result,
        Details:   fmt.Sprintf("Result for product development %s: %s", productName, result),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with result for product development %s: %s.\n", productName, result)
}

// finalizeProductDevelopmentCheckCycle finalizes the product development check cycle and logs the result in the ledger
func (automation *ProductDevelopmentIntegrationAutomation) finalizeProductDevelopmentCheckCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeProductDevelopmentCheckCycle()
    if success {
        fmt.Println("Product development check cycle finalized successfully.")
        automation.logProductDevelopmentCheckCycleFinalization()
    } else {
        fmt.Println("Error finalizing product development check cycle.")
    }
}

// logProductDevelopmentCheckCycleFinalization logs the finalization of a product development check cycle into the ledger
func (automation *ProductDevelopmentIntegrationAutomation) logProductDevelopmentCheckCycleFinalization() {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("product-development-check-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Product Development Check Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with product development check cycle finalization.")
}
