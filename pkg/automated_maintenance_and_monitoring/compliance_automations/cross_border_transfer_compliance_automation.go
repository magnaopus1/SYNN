package compliance_automations

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    CrossBorderCheckInterval = 30 * time.Minute   // Interval for checking cross-border data transfer compliance
)

// CrossBorderTransferComplianceAutomation handles cross-border transfer compliance
type CrossBorderTransferComplianceAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger instance
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus Engine
    stateMutex      *sync.RWMutex                // Mutex for thread-safe state and ledger access
    apiURL          string                       // API URL for compliance checks
}

// NewCrossBorderTransferComplianceAutomation initializes the cross-border transfer compliance handler
func NewCrossBorderTransferComplianceAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *CrossBorderTransferComplianceAutomation {
    return &CrossBorderTransferComplianceAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartCrossBorderComplianceMonitoring initiates continuous monitoring for cross-border data transfers
func (automation *CrossBorderTransferComplianceAutomation) StartCrossBorderComplianceMonitoring() {
    ticker := time.NewTicker(CrossBorderCheckInterval)
    for range ticker.C {
        fmt.Println("Starting cross-border transfer compliance monitoring...")
        automation.monitorCrossBorderTransfers()
    }
}

// monitorCrossBorderTransfers retrieves and validates cross-border data transfer compliance
func (automation *CrossBorderTransferComplianceAutomation) monitorCrossBorderTransfers() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    subBlocks := automation.ledgerInstance.GetValidatedSubBlocks() // Get validated sub-blocks
    for _, subBlock := range subBlocks {
        for _, tx := range subBlock.Transactions {
            if automation.isCrossBorderTransaction(tx) {
                fmt.Printf("Checking cross-border compliance for transaction ID: %s\n", tx.ID)
                automation.validateCrossBorderCompliance(tx)
            }
        }
    }
}

// isCrossBorderTransaction checks if a transaction involves cross-border data transfer
func (automation *CrossBorderTransferComplianceAutomation) isCrossBorderTransaction(tx common.Transaction) bool {
    // Check if sender and recipient are in different countries
    return tx.SenderCountry != tx.RecipientCountry
}

// validateCrossBorderCompliance checks if a cross-border data transfer complies with local and regional laws
func (automation *CrossBorderTransferComplianceAutomation) validateCrossBorderCompliance(tx common.Transaction) {
    url := fmt.Sprintf("%s/api/compliance/check", automation.apiURL)
    body, _ := json.Marshal(map[string]string{
        "transaction_id":   tx.ID,
        "sender_country":   tx.SenderCountry,
        "recipient_country": tx.RecipientCountry,
    })

    encryptedBody, err := encryption.Encrypt(body, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting data for compliance check: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error checking cross-border compliance for transaction ID %s: %v\n", tx.ID, err)
        return
    }

    var complianceResult common.ComplianceResult
    json.NewDecoder(resp.Body).Decode(&complianceResult)

    if !complianceResult.Compliant {
        fmt.Printf("Cross-border data transfer not compliant for transaction ID: %s\n", tx.ID)
        automation.blockNonCompliantTransfer(tx)
    } else {
        fmt.Printf("Transaction ID %s is compliant for cross-border data transfer.\n", tx.ID)
    }

    automation.logComplianceCheck(tx, complianceResult)
}

// blockNonCompliantTransfer blocks a transaction that violates cross-border compliance
func (automation *CrossBorderTransferComplianceAutomation) blockNonCompliantTransfer(tx common.Transaction) {
    fmt.Printf("Blocking non-compliant cross-border transfer for transaction ID: %s\n", tx.ID)
    automation.ledgerInstance.BlockTransaction(tx.ID)
    
    // Log the violation in the blockchain ledger
    automation.logTransferViolation(tx)
}

// logComplianceCheck logs the result of a cross-border compliance check into the ledger
func (automation *CrossBorderTransferComplianceAutomation) logComplianceCheck(tx common.Transaction, result common.ComplianceResult) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        tx.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Cross-border Compliance Check",
        Status:    "Checked",
        Data:      fmt.Sprintf("Compliance check result: %v", result.Compliant),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for cross-border compliance check: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(tx) // Validate through Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with cross-border compliance check for transaction ID: %s\n", tx.ID)
}

// logTransferViolation logs a violation for a non-compliant cross-border transfer into the ledger
func (automation *CrossBorderTransferComplianceAutomation) logTransferViolation(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        tx.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Cross-border Transfer Violation",
        Status:    "Violated",
        Data:      "Non-compliant cross-border data transfer detected",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for cross-border violation: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(tx) // Validate the violation through Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated with cross-border violation for transaction ID: %s\n", tx.ID)
}

// StartComplianceLogRetrieval continuously retrieves cross-border compliance logs from the ledger
func (automation *CrossBorderTransferComplianceAutomation) StartComplianceLogRetrieval() {
    ticker := time.NewTicker(CrossBorderCheckInterval)
    for range ticker.C {
        fmt.Println("Retrieving cross-border compliance logs...")
        automation.retrieveComplianceLogs()
    }
}

// retrieveComplianceLogs retrieves cross-border compliance logs from the ledger for auditing purposes
func (automation *CrossBorderTransferComplianceAutomation) retrieveComplianceLogs() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    complianceLogs := automation.ledgerInstance.GetComplianceLogs("Cross-border Compliance Check")
    for _, log := range complianceLogs {
        fmt.Printf("Compliance log ID %s: %s\n", log.ID, log.Data)
    }
}
