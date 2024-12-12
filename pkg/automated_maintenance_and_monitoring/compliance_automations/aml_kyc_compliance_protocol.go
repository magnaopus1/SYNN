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
)

const (
    AMLCheckInterval  = 10 * time.Minute // Interval for checking AML violations
    KYCCheckInterval  = 30 * time.Minute // Interval for verifying KYC submissions
    EncryptionKey     = "super_secret_key" // Encryption key for sensitive data
)

// AMLKYCComplianceAutomation automates AML and KYC compliance checks
type AMLKYCComplianceAutomation struct {
    ledgerInstance *ledger.Ledger // Blockchain ledger instance
    stateMutex     *sync.RWMutex  // Mutex for thread-safe ledger access
    apiURL         string         // API URL for compliance endpoints
}

// NewAMLKYCComplianceAutomation initializes the AML/KYC automation system
func NewAMLKYCComplianceAutomation(apiURL string, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AMLKYCComplianceAutomation {
    return &AMLKYCComplianceAutomation{
        ledgerInstance: ledgerInstance,
        stateMutex:     stateMutex,
        apiURL:         apiURL,
    }
}

// StartAMLMonitoring starts continuous AML monitoring
func (automation *AMLKYCComplianceAutomation) StartAMLMonitoring() {
    ticker := time.NewTicker(AMLCheckInterval)
    for range ticker.C {
        fmt.Println("Starting AML monitoring...")
        automation.monitorAMLTransactions()
    }
}

// monitorAMLTransactions retrieves sub-blocks and checks each transaction for AML violations
func (automation *AMLKYCComplianceAutomation) monitorAMLTransactions() {
    automation.stateMutex.RLock()
    defer automation.stateMutex.RUnlock()

    subBlocks := automation.ledgerInstance.GetPendingSubBlocks()
    for _, subBlock := range subBlocks {
        for _, tx := range subBlock.Transactions {
            if amlViolation := automation.checkAMLRules(tx); amlViolation {
                fmt.Printf("AML violation detected for transaction ID: %s\n", tx.ID)
                automation.reportAMLViolation(tx)
            }
        }
    }
}

// checkAMLRules determines if a transaction violates AML policies
func (automation *AMLKYCComplianceAutomation) checkAMLRules(tx common.Transaction) bool {
    return automation.isTransactionSuspicious(tx) || automation.isWalletFlagged(tx.Sender) || automation.isHighRiskCountry(tx)
}

// isTransactionSuspicious checks if a transaction is suspicious based on predefined AML thresholds
func (automation *AMLKYCComplianceAutomation) isTransactionSuspicious(tx common.Transaction) bool {
    if tx.Amount > 10000 {
        fmt.Printf("Transaction ID %s flagged: Amount exceeds $10,000.\n", tx.ID)
        return true
    }
    if automation.checkFrequentSmallTransactions(tx) {
        fmt.Printf("Transaction ID %s flagged: Frequent small transactions detected.\n", tx.ID)
        return true
    }
    if automation.isTransactionUnusual(tx) {
        fmt.Printf("Transaction ID %s flagged: Unusual transaction pattern detected.\n", tx.ID)
        return true
    }
    return false
}

// checkFrequentSmallTransactions checks for multiple small transactions within a short time window
func (automation *AMLKYCComplianceAutomation) checkFrequentSmallTransactions(tx common.Transaction) bool {
    smallTxThreshold := 1000.0
    maxTransactionCount := 10
    timeWindow := time.Hour

    history := automation.ledgerInstance.GetTransactionHistory(tx.Sender, timeWindow)
    smallTxCount := 0
    for _, histTx := range history {
        if histTx.Amount < smallTxThreshold {
            smallTxCount++
        }
    }

    return smallTxCount > maxTransactionCount
}

// isTransactionUnusual checks for transactions that deviate from normal behavior
func (automation *AMLKYCComplianceAutomation) isTransactionUnusual(tx common.Transaction) bool {
    history := automation.ledgerInstance.GetTransactionHistory(tx.Sender, time.Hour*24*30) // Last 30 days
    totalAmount := 0.0
    for _, histTx := range history {
        totalAmount += histTx.Amount
    }
    averageTransactionAmount := totalAmount / float64(len(history))
    return tx.Amount > averageTransactionAmount*1.5
}

// isWalletFlagged checks if a wallet has been flagged for suspicious activity
func (automation *AMLKYCComplianceAutomation) isWalletFlagged(walletAddress string) bool {
    flaggedWallets := automation.ledgerInstance.GetFlaggedWallets()
    if _, exists := flaggedWallets[walletAddress]; exists {
        fmt.Printf("Wallet %s is flagged for suspicious activity.\n", walletAddress)
        return true
    }
    return false
}

// isHighRiskCountry checks if a transaction involves a high-risk country
func (automation *AMLKYCComplianceAutomation) isHighRiskCountry(tx common.Transaction) bool {
    highRiskCountries := []string{"North Korea", "Iran", "Syria", "Sudan", "Yemen"}
    return automation.isCountryHighRisk(tx.SenderCountry) || automation.isCountryHighRisk(tx.RecipientCountry)
}

// isCountryHighRisk checks if a country is considered high-risk for AML
func (automation *AMLKYCComplianceAutomation) isCountryHighRisk(country string) bool {
    highRiskCountries := []string{"North Korea", "Iran", "Syria", "Sudan", "Yemen"}
    for _, riskCountry := range highRiskCountries {
        if country == riskCountry {
            return true
        }
    }
    return false
}

// reportAMLViolation reports an AML violation to the API and blocks the wallet if necessary
func (automation *AMLKYCComplianceAutomation) reportAMLViolation(tx common.Transaction) {
    url := fmt.Sprintf("%s/api/compliance/aml/monitor_transaction", automation.apiURL)
    body, _ := json.Marshal(tx)

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting transaction data: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error reporting AML violation: %v\n", err)
        return
    }

    automation.blockWallet(tx.Sender)
}

// blockWallet blocks a wallet due to AML violations
func (automation *AMLKYCComplianceAutomation) blockWallet(walletAddress string) {
    url := fmt.Sprintf("%s/api/compliance/aml/block_wallet", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"wallet_address": walletAddress})

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting wallet data: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error blocking wallet: %v\n", err)
    } else {
        fmt.Printf("Wallet %s blocked due to AML violation.\n", walletAddress)
    }
}

// unblockWallet unblocks a wallet after clearing AML investigations
func (automation *AMLKYCComplianceAutomation) unblockWallet(walletAddress string) {
    url := fmt.Sprintf("%s/api/compliance/aml/unblock_wallet", automation.apiURL)
    body, _ := json.Marshal(map[string]string{"wallet_address": walletAddress})

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting wallet data: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error unblocking wallet: %v\n", err)
    } else {
        fmt.Printf("Wallet %s unblocked.\n", walletAddress)
    }
}

// StartKYCVerification starts continuous KYC verification
func (automation *AMLKYCComplianceAutomation) StartKYCVerification() {
    ticker := time.NewTicker(KYCCheckInterval)
    for range ticker.C {
        fmt.Println("Starting KYC verification...")
        automation.verifyKYCSubmissions()
    }
}

// verifyKYCSubmissions retrieves and processes pending KYC submissions
func (automation *AMLKYCComplianceAutomation) verifyKYCSubmissions() {
    url := fmt.Sprintf("%s/api/compliance/kyc/retrieve", automation.apiURL)
    resp, err := http.Get(url)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error retrieving KYC submissions: %v\n", err)
        return
    }
    defer resp.Body.Close()

    var submissions []common.KYCSubmission
    json.NewDecoder(resp.Body).Decode(&submissions)

    for _, submission := range submissions {
        automation.verifyKYCSubmission(submission)
    }
}

// verifyKYCSubmission processes a KYC submission
func (automation *AMLKYCComplianceAutomation) verifyKYCSubmission(submission common.KYCSubmission) {
    if automation.checkKYCValidity(submission) {
        automation.approveKYC(submission)
    } else {
        automation.rejectKYC(submission)
    }
}

// checkKYCValidity verifies a KYC submission's validity
func (automation *AMLKYCComplianceAutomation) checkKYCValidity(submission common.KYCSubmission) bool {
    return automation.isIdentityDocumentValid(submission.IdentityDocument) &&
        !automation.hasCriminalRecord(submission.IdentityDocument) &&
        automation.verifyAddress(submission.Address) &&
        automation.validateBiometricData(submission.BiometricData)
}

// isIdentityDocumentValid validates an identity document
func (automation *AMLKYCComplianceAutomation) isIdentityDocumentValid(doc common.IdentityDocument) bool {
    if doc.Expiration.Before(time.Now()) {
        fmt.Printf("Identity document expired: %s\n", doc.DocumentNumber)
        return false
    }
    return true
}

// hasCriminalRecord checks for criminal records associated with an identity document
func (automation *AMLKYCComplianceAutomation) hasCriminalRecord(doc common.IdentityDocument) bool {
    criminalRecordService := automation.ledgerInstance.GetCriminalRecordService()
    return criminalRecordService.CheckCriminalRecord(doc.DocumentNumber)
}

// verifyAddress validates an address submission
func (automation *AMLKYCComplianceAutomation) verifyAddress(address common.Address) bool {
    addressVerificationService := automation.ledgerInstance.GetAddressVerificationService()
    return addressVerificationService.VerifyAddress(address.Street, address.City, address.Country, address.PostalCode)
}

// validateBiometricData validates biometric data for a KYC submission
func (automation *AMLKYCComplianceAutomation) validateBiometricData(biometricData common.BiometricData) bool {
    biometricService := automation.ledgerInstance.GetBiometricValidationService()
    return biometricService.ValidateBiometrics(biometricData)
}

// approveKYC approves a KYC submission and adds it to the ledger
func (automation *AMLKYCComplianceAutomation) approveKYC(submission common.KYCSubmission) {
    url := fmt.Sprintf("%s/api/compliance/kyc/verify", automation.apiURL)
    body, _ := json.Marshal(submission)

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting KYC submission: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error approving KYC submission: %v\n", err)
    } else {
        fmt.Printf("KYC submission %s approved.\n", submission.ID)
    }

    automation.addToLedger(submission.ID)
}

// rejectKYC rejects a KYC submission
func (automation *AMLKYCComplianceAutomation) rejectKYC(submission common.KYCSubmission) {
    url := fmt.Sprintf("%s/api/compliance/kyc/reject", automation.apiURL)
    body, _ := json.Marshal(submission)

    encryptedBody, err := encryption.Encrypt(body, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting KYC submission: %v\n", err)
        return
    }

    resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptedBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Printf("Error rejecting KYC submission: %v\n", err)
    } else {
        fmt.Printf("KYC submission %s rejected.\n", submission.ID)
    }
}

// addToLedger adds an entry to the ledger for compliance tracking
func (automation *AMLKYCComplianceAutomation) addToLedger(submissionID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        submissionID,
        Timestamp: time.Now().Unix(),
        Type:      "KYC",
        Status:    "Verified",
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting ledger entry: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Ledger entry added for submission ID: %s\n", submissionID)
}
