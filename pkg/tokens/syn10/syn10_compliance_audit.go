package syn10

import (
	"fmt"
	"time"
    "synnergy_network/pkg/common"

)

func (token *SYN10Token) checkComplianceAudit(transactionID string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Step 2: Check compliance standards
    auditResult, err := token.CheckSYN10Compliance(transactionID)
    if err != nil || !auditResult {
        // Log failure
        reason := fmt.Sprintf("compliance check failed for transaction %s: %v", transactionID, err)
        _ = token.Ledger.LogAuditFailure(transactionID, reason)
        return false, fmt.Errorf(reason)
    }

    // Step 3: Record audit log
    auditLog := fmt.Sprintf("Audit passed for transaction %s at %v", transactionID, time.Now())
    err = token.Ledger.RecordComplianceAuditLog(transactionID, auditLog)
    if err != nil {
        return false, fmt.Errorf("failed to record compliance audit log: %v", err)
    }

    return true, nil
}


func (token *SYN10Token) ReviewComplianceAudit(transactionID string) (string, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Step 1: Fetch transaction history
    history, err := token.Ledger.GetSYN10TransactionHistory(transactionID) // Fixed call
    if err != nil {
        return "", fmt.Errorf("failed to retrieve transaction history: %v", err)
    }

    // Step 2: Verify compliance standards in history
    complianceVerified, err := token.VerifySYN10History(history) // Corrected method call
    if err != nil {
        return "", fmt.Errorf("compliance verification failed: %v", err)
    }

    if !complianceVerified {
        return "", fmt.Errorf("compliance review failed for transaction %s", transactionID)
    }

    return fmt.Sprintf("Compliance review passed for transaction %s: %v", transactionID, history), nil
}


func (token *SYN10Token) CheckSYN10Compliance(transactionID string, encryptedData string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Step 1: Convert SYN10 transaction to standard transaction
    standardTransaction, err := token.ConvertToStandardTransaction(transactionID, encryptedData)
    if err != nil {
        reason := fmt.Sprintf("failed to convert SYN10 transaction %s: %v", transactionID, err)
        _ = token.Ledger.LogAuditFailure(transactionID, reason)
        return false, fmt.Errorf(reason)
    }

    // Step 2: Validate transaction using consensus
    err = token.Consensus.ProcessSingleTransaction(
        token.Consensus,
        standardTransaction.TransactionID,
        standardTransaction.EncryptedData,
        false,
        common.CrossChainTransaction{},
    )
    if err != nil {
        reason := fmt.Sprintf("transaction %s failed consensus validation: %v", transactionID, err)
        _ = token.Ledger.LogAuditFailure(transactionID, reason)
        return false, fmt.Errorf(reason)
    }

    // Step 3: Check compliance standards
    auditResult, err := token.Compliance.CheckCompliance(transactionID)
    if err != nil || !auditResult {
        reason := fmt.Sprintf("compliance check failed for transaction %s: %v", transactionID, err)
        _ = token.Ledger.LogAuditFailure(transactionID, reason)
        return false, fmt.Errorf(reason)
    }

    // Step 4: Log successful compliance check
    logEntry := fmt.Sprintf("Compliance check passed for transaction %s at %v", transactionID, time.Now())
    err = token.Ledger.RecordComplianceAuditLog(transactionID, logEntry)
    if err != nil {
        return false, fmt.Errorf("failed to log compliance audit: %v", err)
    }

    return true, nil
}


func (token *SYN10Token) ConvertToStandardTransaction(transactionID string, encryptedData string) (*common.Transaction, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Check if the transaction exists in the ledger
    transactionHistory, exists := token.Ledger.TransactionHistory[transactionID]
    if !exists {
        return nil, fmt.Errorf("transaction ID %s not found in SYN10 ledger", transactionID)
    }

    // Create a new standard transaction from the SYN10 transaction
    standardTransaction := &common.Transaction{
        TransactionID: transactionID,
        EncryptedData: encryptedData,
        Status:        "Pending",
        Metadata: map[string]interface{}{
            "source": "SYN10",
            "history": transactionHistory,
        },
    }

    return standardTransaction, nil
}

func (token *SYN10Token) ValidateSYN10Transaction(sc *common.SynnergyConsensus, transactionID string, encryptedData string, isCrossChain bool, crossChainData common.CrossChainTransaction) error {
    // Step 1: Convert SYN10 transaction to standard transaction
    standardTransaction, err := token.ConvertToStandardTransaction(transactionID, encryptedData)
    if err != nil {
        return fmt.Errorf("failed to convert SYN10 transaction: %v", err)
    }

    // Step 2: Process the transaction using SynnergyConsensus
    if isCrossChain {
        return sc.ProcessSingleTransaction(transactionID, encryptedData, isCrossChain, crossChainData)
    }

    return sc.ProcessSingleTransaction(standardTransaction.TransactionID, standardTransaction.EncryptedData, false, common.CrossChainTransaction{})
}

func (token *SYN10Token) ValidateSYN10Transactions(sc *common.SynnergyConsensus, transactions []string, encryptedData []string) error {
    if len(transactions) != len(encryptedData) {
        return fmt.Errorf("mismatch between transaction IDs and encrypted data length")
    }

    var standardTransactions []common.Transaction
    for i, transactionID := range transactions {
        // Convert each SYN10 transaction to a standard transaction
        standardTransaction, err := token.ConvertToStandardTransaction(transactionID, encryptedData[i])
        if err != nil {
            return fmt.Errorf("failed to convert SYN10 transaction: %v", err)
        }

        standardTransactions = append(standardTransactions, *standardTransaction)
    }

    // Process the batch of transactions using SynnergyConsensus
    sc.ProcessTransactions(standardTransactions, nil)
    return nil
}



func (token *SYN10Token) LogComplianceActivity(transactionID string, activity string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Step 1: Generate log entry
    logEntry := fmt.Sprintf("Transaction %s: %s at %v", transactionID, activity, time.Now())
    encryptedLog, err := common.EncryptData(logEntry)
    if err != nil {
        return fmt.Errorf("failed to encrypt compliance activity log: %v", err)
    }

    // Step 2: Record encrypted log in ledger
    err = token.Ledger.RecordComplianceAuditLog(transactionID, encryptedLog)
    if err != nil {
        return fmt.Errorf("failed to record compliance activity log: %v", err)
    }

    return nil
}


func (token *SYN10Token) VerifySYN10History(history []string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    // Example: Iterate over history and ensure compliance rules are met
    for _, entry := range history {
        if !token.Compliance.VerifyHistoryEntry(entry) {
            return false, fmt.Errorf("history entry failed compliance: %s", entry)
        }
    }

    return true, nil
}





