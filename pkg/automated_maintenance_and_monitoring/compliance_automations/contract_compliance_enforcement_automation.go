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
    ContractCheckInterval = 15 * time.Minute  // Interval for checking contract compliance
    EnforcementKey        = "enforce_secret_key" // Encryption key for sensitive compliance data
)

// ContractComplianceAutomation handles the automation of contract compliance enforcement
type ContractComplianceAutomation struct {
    ledgerInstance      *ledger.Ledger               // Ledger instance for contract management
    consensusEngine     *synnergy_consensus.Consensus // Synnergy Consensus Engine
    stateMutex          *sync.RWMutex                // Mutex for thread-safe state and ledger access
    apiURL              string                       // API URL for compliance-related endpoints
}

// NewContractComplianceAutomation initializes the contract compliance automation handler
func NewContractComplianceAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *ContractComplianceAutomation {
    return &ContractComplianceAutomation{
        ledgerInstance: ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:     stateMutex,
        apiURL:         apiURL,
    }
}

// StartContractComplianceMonitoring initiates the continuous monitoring of smart contract compliance
func (automation *ContractComplianceAutomation) StartContractComplianceMonitoring() {
    ticker := time.NewTicker(ContractCheckInterval)
    for range ticker.C {
        fmt.Println("Checking contract compliance...")
        automation.enforceContractCompliance()
    }
}

// enforceContractCompliance enforces compliance rules on smart contracts and ensures adherence to regulations
func (automation *ContractComplianceAutomation) enforceContractCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    contracts := automation.ledgerInstance.GetActiveContracts()
    for _, contract := range contracts {
        if !automation.isContractCompliant(contract) {
            fmt.Printf("Contract ID %s is not compliant. Enforcing compliance...\n", contract.ID)
            automation.enforceCompliance(contract)
        }
    }
}

// isContractCompliant checks if a contract is compliant with legal and policy requirements
func (automation *ContractComplianceAutomation) isContractCompliant(contract common.SmartContract) bool {
    if automation.isContractExpired(contract) || automation.isContractBreachingPolicy(contract) {
        return false
    }
    return true
}

// isContractExpired checks if the contract has passed its expiration date
func (automation *ContractComplianceAutomation) isContractExpired(contract common.SmartContract) bool {
    if contract.Expiration.Before(time.Now()) {
        fmt.Printf("Contract ID %s has expired.\n", contract.ID)
        return true
    }
    return false
}

// isContractBreachingPolicy checks if the contract violates any internal or external policy rules
func (automation *ContractComplianceAutomation) isContractBreachingPolicy(contract common.SmartContract) bool {
    policyViolations := automation.ledgerInstance.CheckPolicyViolations(contract)
    if len(policyViolations) > 0 {
        fmt.Printf("Contract ID %s has policy violations: %v\n", contract.ID, policyViolations)
        return true
    }
    return false
}

// enforceCompliance enforces compliance on a contract by calling the relevant API and blocking non-compliant operations
func (automation *ContractComplianceAutomation) enforceCompliance(contract common.SmartContract) {
    url := fmt.Sprintf("%s/api/compliance/enforce", automation.apiURL)
    body, _ := json.Marshal(contract)

    encryptedBody, err := encryption.Encrypt(body, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting contract data for enforcement: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error enforcing compliance on contract ID %s: %v\n", contract.ID, err)
        return
    }

    fmt.Printf("Compliance enforced for contract ID %s successfully.\n", contract.ID)
    automation.updateLedgerForCompliance(contract)
}

// invokeContract automatically invokes a contract if it meets compliance requirements and certain conditions
func (automation *ContractComplianceAutomation) invokeContract(contract common.SmartContract) {
    url := fmt.Sprintf("%s/api/compliance/invoke_contract", automation.apiURL)
    body, _ := json.Marshal(contract)

    encryptedBody, err := encryption.Encrypt(body, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting contract data for invocation: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error invoking contract ID %s: %v\n", contract.ID, err)
        return
    }

    fmt.Printf("Contract ID %s invoked successfully.\n", contract.ID)
    automation.updateLedgerForInvocation(contract)
}

// retrieveComplianceResults retrieves the compliance result for a specific contract
func (automation *ContractComplianceAutomation) retrieveComplianceResults(contractID string) {
    url := fmt.Sprintf("%s/api/compliance/retrieve", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"contract_id": contractID})

    encryptedBody, err := encryption.Encrypt(body, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting compliance result request: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving compliance results for contract ID %s: %v\n", contractID, err)
        return
    }

    var complianceResult common.ComplianceResult
    json.NewDecoder(resp.Body).Decode(&complianceResult)
    fmt.Printf("Compliance results for contract ID %s: %v\n", contractID, complianceResult)
}

// updateLedgerForCompliance adds an entry to the ledger for tracking enforced compliance actions
func (automation *ContractComplianceAutomation) updateLedgerForCompliance(contract common.SmartContract) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        contract.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Contract Compliance",
        Status:    "Enforced",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for contract compliance: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(contract) // Validate contract enforcement through Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for enforced compliance on contract ID: %s\n", contract.ID)
}

// updateLedgerForInvocation adds an entry to the ledger for tracking contract invocations
func (automation *ContractComplianceAutomation) updateLedgerForInvocation(contract common.SmartContract) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        contract.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Contract Invocation",
        Status:    "Invoked",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EnforcementKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for contract invocation: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(contract) // Validate contract invocation through Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for contract invocation on contract ID: %s\n", contract.ID)
}
