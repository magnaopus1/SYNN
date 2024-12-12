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
    RegulatoryCheckInterval = 15 * time.Minute // Interval for checking regulatory changes
)

// RegulatoryChangeEnforcementAutomation handles monitoring and enforcing regulatory changes
type RegulatoryChangeEnforcementAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger instance for full compliance
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus for validation
    stateMutex      *sync.RWMutex                // Mutex for thread-safe ledger and state access
    apiURL          string                       // API URL for compliance endpoints
}

// NewRegulatoryChangeEnforcementAutomation initializes the automation for regulatory change enforcement
func NewRegulatoryChangeEnforcementAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *RegulatoryChangeEnforcementAutomation {
    return &RegulatoryChangeEnforcementAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartRegulatoryMonitoring begins continuous monitoring for regulatory changes
func (automation *RegulatoryChangeEnforcementAutomation) StartRegulatoryMonitoring() {
    ticker := time.NewTicker(RegulatoryCheckInterval)
    for range ticker.C {
        fmt.Println("Checking for regulatory updates...")
        automation.monitorRegulatoryChanges()
    }
}

// monitorRegulatoryChanges checks for regulatory updates and enforces them
func (automation *RegulatoryChangeEnforcementAutomation) monitorRegulatoryChanges() {
    url := fmt.Sprintf("%s/api/compliance/check", automation.apiURL)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error checking for regulatory changes: %v\n", err)
        return
    }
    defer resp.Body.Close()

    var complianceChanges []common.ComplianceRule
    err = json.NewDecoder(resp.Body).Decode(&complianceChanges)
    if err != nil {
        fmt.Printf("Error decoding regulatory changes: %v\n", err)
        return
    }

    for _, change := range complianceChanges {
        automation.enforceRegulatoryChange(change)
    }
}

// enforceRegulatoryChange enforces a specific regulatory change across the blockchain
func (automation *RegulatoryChangeEnforcementAutomation) enforceRegulatoryChange(change common.ComplianceRule) {
    fmt.Printf("Enforcing regulatory change: %s\n", change.Description)

    // Retrieve affected contracts and transactions
    affectedContracts := automation.getAffectedContracts(change)
    affectedTransactions := automation.getAffectedTransactions(change)

    // Apply changes to contracts
    for _, contract := range affectedContracts {
        automation.updateContractCompliance(contract, change)
    }

    // Apply changes to transactions
    for _, tx := range affectedTransactions {
        automation.updateTransactionCompliance(tx, change)
    }

    // Log the enforcement action in the ledger
    automation.logRegulatoryChange(change)
}

// getAffectedContracts retrieves all contracts affected by the regulatory change
func (automation *RegulatoryChangeEnforcementAutomation) getAffectedContracts(change common.ComplianceRule) []common.SmartContract {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    contracts := automation.ledgerInstance.GetSmartContracts()
    var affectedContracts []common.SmartContract

    for _, contract := range contracts {
        if automation.isContractAffectedByChange(contract, change) {
            affectedContracts = append(affectedContracts, contract)
        }
    }

    return affectedContracts
}

// isContractAffectedByChange checks if a contract is affected by a regulatory change
func (automation *RegulatoryChangeEnforcementAutomation) isContractAffectedByChange(contract common.SmartContract, change common.ComplianceRule) bool {
    return contract.RulesVersion != change.Version
}

// getAffectedTransactions retrieves all transactions affected by the regulatory change
func (automation *RegulatoryChangeEnforcementAutomation) getAffectedTransactions(change common.ComplianceRule) []common.Transaction {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    transactions := automation.ledgerInstance.GetTransactions()
    var affectedTransactions []common.Transaction

    for _, tx := range transactions {
        if automation.isTransactionAffectedByChange(tx, change) {
            affectedTransactions = append(affectedTransactions, tx)
        }
    }

    return affectedTransactions
}

// isTransactionAffectedByChange checks if a transaction is affected by a regulatory change
func (automation *RegulatoryChangeEnforcementAutomation) isTransactionAffectedByChange(tx common.Transaction, change common.ComplianceRule) bool {
    return tx.RulesVersion != change.Version
}

// updateContractCompliance applies the new regulatory rules to a contract
func (automation *RegulatoryChangeEnforcementAutomation) updateContractCompliance(contract common.SmartContract, change common.ComplianceRule) {
    fmt.Printf("Updating contract ID %s for compliance with new regulations.\n", contract.ID)

    updatedContract := contract
    updatedContract.RulesVersion = change.Version
    body, _ := json.Marshal(updatedContract)

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting contract update: %v\n", err)
        return
    }

    url := fmt.Sprintf("%s/api/compliance/invoke_contract", automation.apiURL)
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error updating contract ID %s: %v\n", contract.ID, err)
        return
    }

    fmt.Printf("Contract ID %s updated successfully for new compliance rules.\n", contract.ID)
}

// updateTransactionCompliance applies the new regulatory rules to a transaction
func (automation *RegulatoryChangeEnforcementAutomation) updateTransactionCompliance(tx common.Transaction, change common.ComplianceRule) {
    fmt.Printf("Updating transaction ID %s for compliance with new regulations.\n", tx.ID)

    updatedTx := tx
    updatedTx.RulesVersion = change.Version
    body, _ := json.Marshal(updatedTx)

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting transaction update: %v\n", err)
        return
    }

    url := fmt.Sprintf("%s/api/compliance/execute", automation.apiURL)
    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error updating transaction ID %s: %v\n", tx.ID, err)
        return
    }

    fmt.Printf("Transaction ID %s updated successfully for new compliance rules.\n", tx.ID)
}

// logRegulatoryChange logs the enforcement of regulatory changes in the ledger
func (automation *RegulatoryChangeEnforcementAutomation) logRegulatoryChange(change common.ComplianceRule) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        change.ID,
        Timestamp: time.Now().Unix(),
        Type:      "RegulatoryChange",
        Status:    "Enforced",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for regulatory change: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Validate entry through Synnergy Consensus before adding
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Regulatory change ID %s enforced and logged in the ledger.\n", change.ID)
}
