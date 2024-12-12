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
    PolicyCheckInterval     = 20 * time.Minute // Interval for checking policy compliance violations
    PolicyEncryptionKey     = "policy_encryption_key" // Encryption key for policy compliance
)

// PolicyComplianceEnforcementAutomation automates policy compliance enforcement
type PolicyComplianceEnforcementAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger instance
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus engine for validating policy restrictions
    stateMutex      *sync.RWMutex                // Mutex for thread-safe ledger access
    apiURL          string                       // API URL for policy compliance endpoints
}

// NewPolicyComplianceEnforcementAutomation initializes the policy compliance automation
func NewPolicyComplianceEnforcementAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *PolicyComplianceEnforcementAutomation {
    return &PolicyComplianceEnforcementAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartPolicyComplianceMonitoring begins continuous monitoring of policy compliance
func (automation *PolicyComplianceEnforcementAutomation) StartPolicyComplianceMonitoring() {
    ticker := time.NewTicker(PolicyCheckInterval)
    for range ticker.C {
        fmt.Println("Starting policy compliance monitoring...")
        automation.monitorPolicyCompliance()
    }
}

// monitorPolicyCompliance checks transactions for potential policy violations
func (automation *PolicyComplianceEnforcementAutomation) monitorPolicyCompliance() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    subBlocks := automation.ledgerInstance.GetPendingSubBlocks()
    for _, subBlock := range subBlocks {
        for _, tx := range subBlock.Transactions {
            if policyViolation := automation.checkPolicyRestrictions(tx); policyViolation {
                fmt.Printf("Policy violation detected for transaction ID: %s\n", tx.ID)
                automation.applyPolicyRestrictions(tx)
            }
        }
    }
}

// checkPolicyRestrictions analyzes a transaction to determine if it violates organizational policies
func (automation *PolicyComplianceEnforcementAutomation) checkPolicyRestrictions(tx common.Transaction) bool {
    if automation.isRestrictedByOperationalPolicy(tx) || 
       automation.isFinancialPolicyViolated(tx) || 
       automation.isRegulatoryPolicyBreached(tx) {
        return true
    }
    return false
}

// isRestrictedByOperationalPolicy checks if a transaction violates operational restrictions
func (automation *PolicyComplianceEnforcementAutomation) isRestrictedByOperationalPolicy(tx common.Transaction) bool {
    restrictedHours := []int{22, 23, 0, 1, 2, 3, 4, 5} // Restrict transactions during specific hours
    currentHour := time.Now().Hour()
    for _, hour := range restrictedHours {
        if currentHour == hour {
            fmt.Printf("Transaction ID %s flagged: Violates operational policy (restricted hours).\n", tx.ID)
            return true
        }
    }

    authorizedRoles := automation.ledgerInstance.GetAuthorizedRoles(tx.Sender)
    if !automation.isRoleAuthorized(authorizedRoles, tx.RequiredRole) {
        fmt.Printf("Transaction ID %s flagged: Sender is not authorized.\n", tx.ID)
        return true
    }

    if automation.isTransactionRegionRestricted(tx.SenderCountry) {
        fmt.Printf("Transaction ID %s flagged: Sender country is restricted.\n", tx.ID)
        return true
    }

    return false
}

// isRoleAuthorized checks if the sender's role is authorized for the transaction
func (automation *PolicyComplianceEnforcementAutomation) isRoleAuthorized(authorizedRoles []string, requiredRole string) bool {
    for _, role := range authorizedRoles {
        if role == requiredRole {
            return true
        }
    }
    return false
}

// isTransactionRegionRestricted checks if the sender's country is restricted from performing transactions
func (automation *PolicyComplianceEnforcementAutomation) isTransactionRegionRestricted(senderCountry string) bool {
    restrictedRegions := []string{"North Korea", "Iran", "Syria", "Sudan"}
    for _, region := range restrictedRegions {
        if senderCountry == region {
            return true
        }
    }
    return false
}

// isFinancialPolicyViolated checks if a transaction violates financial policy restrictions
func (automation *PolicyComplianceEnforcementAutomation) isFinancialPolicyViolated(tx common.Transaction) bool {
    financialThreshold := 50000.0 // Limit for high-value transactions
    if tx.Amount > financialThreshold {
        fmt.Printf("Transaction ID %s flagged: Violates financial policy (amount exceeds limit).\n", tx.ID)
        return true
    }

    if automation.requiresSpecialApproval(tx) {
        fmt.Printf("Transaction ID %s flagged: Requires special approval.\n", tx.ID)
        return true
    }

    if automation.detectSuspiciousActivity(tx.Sender) {
        fmt.Printf("Transaction ID %s flagged: Suspicious financial activity detected.\n", tx.ID)
        return true
    }

    return false
}

// requiresSpecialApproval checks if a transaction requires special approval
func (automation *PolicyComplianceEnforcementAutomation) requiresSpecialApproval(tx common.Transaction) bool {
    highRiskIndustries := []string{"cryptocurrency", "arms trading", "gambling"}
    for _, industry := range highRiskIndustries {
        if tx.Industry == industry && tx.Amount > 25000 {
            return true
        }
    }
    return false
}

// detectSuspiciousActivity checks for suspicious financial activity, such as structuring
func (automation *PolicyComplianceEnforcementAutomation) detectSuspiciousActivity(sender string) bool {
    transactionHistory := automation.ledgerInstance.GetTransactionHistory(sender, time.Hour*24)
    smallTransactionCount := 0
    for _, tx := range transactionHistory {
        if tx.Amount < 1000 {
            smallTransactionCount++
        }
    }

    if smallTransactionCount > 10 {
        return true
    }

    return false
}

// isRegulatoryPolicyBreached checks if a transaction violates regulatory restrictions
func (automation *PolicyComplianceEnforcementAutomation) isRegulatoryPolicyBreached(tx common.Transaction) bool {
    restrictedEntities := automation.ledgerInstance.GetRestrictedEntities()
    if _, exists := restrictedEntities[tx.Sender]; exists {
        fmt.Printf("Transaction ID %s flagged: Sender is a restricted entity.\n", tx.ID)
        return true
    }

    restrictedIndustries := automation.ledgerInstance.GetRestrictedIndustries()
    if _, exists := restrictedIndustries[tx.Industry]; exists {
        fmt.Printf("Transaction ID %s flagged: Involves a restricted industry.\n", tx.ID)
        return true
    }

    if automation.isCrossBorderRestricted(tx.SenderCountry, tx.RecipientCountry) {
        fmt.Printf("Transaction ID %s flagged: Cross-border transaction violates regulatory policy.\n", tx.ID)
        return true
    }

    return false
}

// isCrossBorderRestricted checks if a transaction between two countries is restricted
func (automation *PolicyComplianceEnforcementAutomation) isCrossBorderRestricted(senderCountry, recipientCountry string) bool {
    restrictedCountryPairs := map[string][]string{
        "United States": {"North Korea", "Iran"},
        "European Union": {"Syria", "Sudan"},
    }

    if restrictedRecipients, exists := restrictedCountryPairs[senderCountry]; exists {
        for _, restrictedCountry := range restrictedRecipients {
            if recipientCountry == restrictedCountry {
                return true
            }
        }
    }

    return false
}

// applyPolicyRestrictions applies the necessary restrictions based on policy violations
func (automation *PolicyComplianceEnforcementAutomation) applyPolicyRestrictions(tx common.Transaction) {
    url := fmt.Sprintf("%s/api/compliance/restrictions/apply", automation.apiURL)
    body, _ := json.Marshal(tx)

    encryptedBody, err := encryption.Encrypt(body, []byte(PolicyEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting policy restriction data: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error applying policy restrictions for transaction ID %s: %v\n", tx.ID, err)
        return
    }

    fmt.Printf("Policy restrictions applied for transaction ID %s.\n", tx.ID)
    automation.updateLedgerForPolicyCompliance(tx)
}

// updateLedgerForPolicyCompliance updates the ledger with the applied policy restrictions
func (automation *PolicyComplianceEnforcementAutomation) updateLedgerForPolicyCompliance(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        tx.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Policy Compliance",
        Status:    "Restricted",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(PolicyEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting policy ledger entry: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Synnergy Consensus validation
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for policy compliance on transaction ID: %s\n", tx.ID)
}

// retrievePolicyRestrictions retrieves the current policy restrictions applied to a specific transaction
func (automation *PolicyComplianceEnforcementAutomation) retrievePolicyRestrictions(txID string) {
    url := fmt.Sprintf("%s/api/compliance/restrictions/retrieve", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"transaction_id": txID})

    encryptedBody, err := encryption.Encrypt(body, []byte(PolicyEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting request for retrieving policy restrictions: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving policy restrictions for transaction ID %s: %v\n", txID, err)
        return
    }

    var result common.PolicyRestrictions
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Printf("Policy restrictions for transaction ID %s: %v\n", txID, result)
}

// validatePolicyRestrictions validates if the current transaction restrictions are adhered to
func (automation *PolicyComplianceEnforcementAutomation) validatePolicyRestrictions(txID string) bool {
    url := fmt.Sprintf("%s/api/compliance/restrictions/validate", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"transaction_id": txID})

    encryptedBody, err := encryption.Encrypt(body, []byte(PolicyEncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting request for policy validation: %v\n", err)
        return false
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error validating policy restrictions for transaction ID %s: %v\n", txID, err)
        return false
    }

    var validationResponse common.PolicyValidationResponse
    json.NewDecoder(resp.Body).Decode(&validationResponse)
    fmt.Printf("Policy validation result for transaction ID %s: %v\n", txID, validationResponse)

    return validationResponse.IsValid
}
