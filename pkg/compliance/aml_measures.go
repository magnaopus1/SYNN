package compliance

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "synnergy_network/pkg/common"

    "synnergy_network/pkg/ledger"
)


// NewAMLSystem initializes the AML system with a suspicious activity threshold
func NewAMLSystem(threshold float64, ledgerInstance *ledger.Ledger) *AMLSystem {
    return &AMLSystem{
        SuspiciousActivityThreshold: threshold,
        BlockedWallets:              make(map[string]bool),
        ReportedTransactions:        make(map[string]string),
        LedgerInstance:              ledgerInstance,
    }
}

// MonitorTransaction monitors a transaction to detect suspicious activity
func (aml *AMLSystem) MonitorTransaction(tx common.Transaction) error {
    aml.mutex.Lock()
    defer aml.mutex.Unlock()

    // Check if any of the wallets are blocked
    if aml.BlockedWallets[tx.FromAddress] || aml.BlockedWallets[tx.ToAddress] {
        return fmt.Errorf("transaction involves a blocked wallet: %s or %s", tx.FromAddress, tx.ToAddress)
    }

    // Detect suspicious activity based on the transaction amount
    if tx.Amount > aml.SuspiciousActivityThreshold {
        aml.ReportSuspiciousTransaction(tx)
    }

    return nil
}


// ReportSuspiciousTransaction reports a suspicious transaction
func (aml *AMLSystem) ReportSuspiciousTransaction(tx common.Transaction) {
    txID := generateTransactionID(tx)
    aml.ReportedTransactions[txID] = fmt.Sprintf("Suspicious transaction from %s to %s for amount %.2f", tx.FromAddress, tx.ToAddress, tx.Amount)

    // Encrypt and store the report in the ledger
    encryptionInstance := & common.Encryption{} // Assuming you have an instance of Encryption

    encryptedReport, err := encryptionInstance.EncryptData("AES", []byte(aml.ReportedTransactions[txID]), common.EncryptionKey)
    if err != nil {
        fmt.Printf("Error encrypting suspicious transaction report: %v\n", err)
        return
    }

    // Convert encryptedReport from []byte to string
    encryptedReportStr := string(encryptedReport)

    // RecordTransactionReport expects a string, so we pass the encryptedReportStr
    result, err := aml.LedgerInstance.ComplianceLedger.RecordTransactionReport(txID, encryptedReportStr)
    if err != nil {
        fmt.Printf("Failed to store suspicious transaction report in the ledger: %v\n", err)
    } else {
        fmt.Printf("Reported suspicious transaction: %s, Result: %v\n", txID, result)
    }
}



// BlockWallet blocks a wallet from performing further transactions
func (aml *AMLSystem) BlockWallet(walletAddress string) {
    aml.mutex.Lock()
    defer aml.mutex.Unlock()

    aml.BlockedWallets[walletAddress] = true
    fmt.Printf("Wallet %s has been blocked.\n", walletAddress)
}

// UnblockWallet unblocks a previously blocked wallet
func (aml *AMLSystem) UnblockWallet(walletAddress string) {
    aml.mutex.Lock()
    defer aml.mutex.Unlock()

    if aml.BlockedWallets[walletAddress] {
        delete(aml.BlockedWallets, walletAddress)
        fmt.Printf("Wallet %s has been unblocked.\n", walletAddress)
    } else {
        fmt.Printf("Wallet %s was not blocked.\n", walletAddress)
    }
}

// ListReportedTransactions lists all reported suspicious transactions
func (aml *AMLSystem) ListReportedTransactions() {
    fmt.Println("Reported Suspicious Transactions:")
    for txID, report := range aml.ReportedTransactions {
        fmt.Printf("Transaction ID: %s, Report: %s\n", txID, report)
    }
}

// generateTransactionID generates a unique ID for each transaction
func generateTransactionID(tx common.Transaction) string {
    // Use FromAddress and ToAddress instead of From and To
    hashInput := fmt.Sprintf("%s%s%.2f", tx.FromAddress, tx.ToAddress, tx.Amount)
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

