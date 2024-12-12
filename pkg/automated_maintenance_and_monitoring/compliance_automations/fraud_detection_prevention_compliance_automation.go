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
    FraudCheckInterval     = 10 * time.Minute  // Interval for checking for fraudulent activity
    FraudDetectionKey      = "fraud_detection_key" // Encryption key for fraud detection operations
    FraudThresholdAmount   = 10000.0           // Example threshold for suspicious transaction amounts
    FraudSuspiciousPattern = 10                // Number of small transactions in a short period to flag as suspicious
)

// FraudDetectionPreventionAutomation automates the detection and prevention of fraudulent activities
type FraudDetectionPreventionAutomation struct {
    ledgerInstance  *ledger.Ledger               // Blockchain ledger for tracking transactions
    consensusEngine *synnergy_consensus.Consensus // Synnergy Consensus engine for validating fraud reports
    stateMutex      *sync.RWMutex                // Mutex for thread-safe ledger access
    apiURL          string                       // API URL for fraud-related endpoints
}

// NewFraudDetectionPreventionAutomation initializes the fraud detection and prevention handler
func NewFraudDetectionPreventionAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Consensus, stateMutex *sync.RWMutex) *FraudDetectionPreventionAutomation {
    return &FraudDetectionPreventionAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartFraudMonitoring initiates continuous monitoring of suspicious activities
func (automation *FraudDetectionPreventionAutomation) StartFraudMonitoring() {
    ticker := time.NewTicker(FraudCheckInterval)
    for range ticker.C {
        fmt.Println("Checking for fraudulent activity...")
        automation.monitorFraudulentTransactions()
    }
}

// monitorFraudulentTransactions retrieves transactions and checks for irregular patterns or possible fraud
func (automation *FraudDetectionPreventionAutomation) monitorFraudulentTransactions() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    transactions := automation.ledgerInstance.GetRecentTransactions() // Get recent transactions from ledger
    for _, tx := range transactions {
        if automation.isFraudulentTransaction(tx) {
            fmt.Printf("Fraudulent transaction detected: %s. Initiating fraud prevention actions.\n", tx.ID)
            automation.blockFraudulentWallet(tx)
            automation.logFraudulentTransaction(tx)
        }
    }
}

// isFraudulentTransaction checks for patterns that indicate possible fraud, such as large amounts or repeated small transactions
func (automation *FraudDetectionPreventionAutomation) isFraudulentTransaction(tx common.Transaction) bool {
    if tx.Amount > FraudThresholdAmount {
        fmt.Printf("Transaction ID %s flagged for exceeding fraud threshold.\n", tx.ID)
        return true
    }

    if automation.checkSuspiciousPattern(tx) {
        fmt.Printf("Transaction ID %s flagged for suspicious small transaction pattern.\n", tx.ID)
        return true
    }

    return false
}

// checkSuspiciousPattern checks for a pattern of small transactions within a short time window
func (automation *FraudDetectionPreventionAutomation) checkSuspiciousPattern(tx common.Transaction) bool {
    smallTxCount := 0
    smallTxThreshold := 500.0   // Example threshold for small transactions
    timeWindow := 1 * time.Hour // Time window to check for suspicious patterns

    history := automation.ledgerInstance.GetTransactionHistory(tx.Sender, timeWindow)
    for _, histTx := range history {
        if histTx.Amount < smallTxThreshold {
            smallTxCount++
        }
    }

    return smallTxCount >= FraudSuspiciousPattern
}

// blockFraudulentWallet blocks a wallet involved in fraudulent activity
func (automation *FraudDetectionPreventionAutomation) blockFraudulentWallet(tx common.Transaction) {
    url := fmt.Sprintf("%s/api/compliance/aml/block_wallet", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"wallet_address": tx.Sender})

    // Encrypt the request before sending it
    encryptedBody, err := encryption.Encrypt(body, []byte(FraudDetectionKey))
    if err != nil {
        fmt.Printf("Error encrypting wallet blocking request for wallet %s: %v\n", tx.Sender, err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error blocking fraudulent wallet %s: %v\n", tx.Sender, err)
        return
    }

    fmt.Printf("Wallet %s blocked successfully due to suspected fraud.\n", tx.Sender)
    automation.updateLedgerForBlockedWallet(tx.Sender)
}

// updateLedgerForBlockedWallet logs the wallet block in the ledger
func (automation *FraudDetectionPreventionAutomation) updateLedgerForBlockedWallet(wallet string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        wallet,
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Blocked",
        Status:    "Fraud Block",
    }

    // Encrypt the ledger entry before adding it to the ledger
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(FraudDetectionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for wallet block: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Validate block with Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for blocked wallet: %s\n", wallet)
}

// logFraudulentTransaction logs the fraudulent transaction in the ledger for audit and investigation purposes
func (automation *FraudDetectionPreventionAutomation) logFraudulentTransaction(tx common.Transaction) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        tx.ID,
        Timestamp: time.Now().Unix(),
        Type:      "Fraudulent Transaction",
        Status:    "Reported",
    }

    // Encrypt the ledger entry before logging it
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(FraudDetectionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for fraudulent transaction: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Validate the entry through Synnergy Consensus
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for fraudulent transaction ID: %s\n", tx.ID)
}

// retrieveReportedFraudulentTransactions retrieves the list of reported fraudulent transactions for further investigation
func (automation *FraudDetectionPreventionAutomation) retrieveReportedFraudulentTransactions() {
    url := fmt.Sprintf("%s/api/compliance/aml/list_reported_transactions", automation.apiURL)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving reported fraudulent transactions: %v\n", err)
        return
    }

    var reportedTransactions []common.FraudulentTransaction
    json.NewDecoder(resp.Body).Decode(&reportedTransactions)
    fmt.Printf("Reported fraudulent transactions: %v\n", reportedTransactions)
}

// unblockWallet unblocks a wallet that has cleared its fraud investigation
func (automation *FraudDetectionPreventionAutomation) unblockWallet(wallet string) {
    url := fmt.Sprintf("%s/api/compliance/aml/unblock_wallet", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"wallet_address": wallet})

    encryptedBody, err := encryption.Encrypt(body, []byte(FraudDetectionKey))
    if err != nil {
        fmt.Printf("Error encrypting wallet unblocking request for wallet %s: %v\n", wallet, err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error unblocking wallet %s: %v\n", wallet, err)
        return
    }

    fmt.Printf("Wallet %s unblocked successfully.\n", wallet)
    automation.updateLedgerForUnblockedWallet(wallet)
}

// updateLedgerForUnblockedWallet logs the wallet unblocking event in the ledger
func (automation *FraudDetectionPreventionAutomation) updateLedgerForUnblockedWallet(wallet string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        wallet,
        Timestamp: time.Now().Unix(),
        Type:      "Wallet Unblocked",
        Status:    "Fraud Cleared",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(FraudDetectionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry for wallet unblock: %v\n", err)
        return
    }

    automation.consensusEngine.ValidateSubBlock(entry) // Synnergy Consensus validation
    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger updated for unblocked wallet: %s\n", wallet)
}
