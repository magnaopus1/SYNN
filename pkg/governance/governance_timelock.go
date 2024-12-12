package governance

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewGovernanceTimelock initializes a new Governance Timelock
func NewGovernanceTimelock(ledgerInstance *ledger.Ledger) *GovernanceTimelock {
    return &GovernanceTimelock{
        PendingExecutions: make(map[string]*TimelockExecution),
        LedgerInstance:    ledgerInstance,
    }
}

// AddProposalToTimelock adds a governance proposal to the timelock queue
func (gt *GovernanceTimelock) AddProposalToTimelock(proposalID, creator string, executionDelay time.Duration, encryptedProposal string) error {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    executionTime := time.Now().Add(executionDelay)

    if _, exists := gt.PendingExecutions[proposalID]; exists {
        return fmt.Errorf("proposal %s is already pending execution", proposalID)
    }

    newExecution := &TimelockExecution{
        ProposalID:      proposalID,
        ExecutionTime:   executionTime,
        EncryptedProposal: encryptedProposal,
    }

    gt.PendingExecutions[proposalID] = newExecution

    // Log to the ledger with the creator argument
    err := gt.logTimelockToLedger(proposalID, creator, executionTime)
    if err != nil {
        return fmt.Errorf("failed to log timelock to ledger: %v", err)
    }

    fmt.Printf("Proposal %s added to timelock. Will execute at %s\n", proposalID, executionTime)
    return nil
}


// ExecuteProposal executes a governance proposal after the timelock period has expired
func (gt *GovernanceTimelock) ExecuteProposal(proposalID string, vm *common.VirtualMachine) error {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    execution, exists := gt.PendingExecutions[proposalID]
    if !exists {
        return errors.New("proposal not found in timelock queue")
    }

    if time.Now().Before(execution.ExecutionTime) {
        return errors.New("timelock period has not expired for this proposal")
    }

    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Convert EncryptedProposal (string) to []byte before decryption
    encryptedProposalBytes := []byte(execution.EncryptedProposal)

    // Decrypt the proposal details using the correct DecryptData parameters
    proposalDetails, err := encryptionInstance.DecryptData(encryptedProposalBytes, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt proposal details: %v", err)
    }

    // Execute the proposal using the Virtual Machine
    executionResult, err := vm.ExecuteContract(
        proposalID,                // contractID
        string(proposalDetails),    // contractSource (could be the proposal details, depending on your structure)
        "language",                 // specify the language, e.g., "Solidity", "Rust", etc.
        nil,                        // parameters (if any)
        encryptionInstance,         // Encryption instance
        common.EncryptionKey,       // encryptionKey for secure execution
    )

    if err != nil {
        return fmt.Errorf("virtual machine execution failed for proposal %s: %v", proposalID, err)
    }

    // Log the execution result
    fmt.Printf("Proposal %s executed successfully: %v\n", proposalID, executionResult)

    // Remove from pending executions after successful execution
    delete(gt.PendingExecutions, proposalID)

    // Log the execution to the ledger, including the `creator`
    err = gt.logExecutionToLedger(proposalID, execution.Creator)
    if err != nil {
        return fmt.Errorf("failed to log execution to ledger: %v", err)
    }

    return nil
}



// logTimelockToLedger logs the timelock details to the ledger
func (gt *GovernanceTimelock) logTimelockToLedger(proposalID, creator string, executionTime time.Time) error {
    logEntry := fmt.Sprintf("Proposal %s is timelocked until %s", proposalID, executionTime)

    // Create an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt the log entry using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logEntry), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt timelock log: %v", err)
    }

    // Generate a hash for the proposal
    hash := gt.generateProposalHash(proposalID)

    // Log the encrypted timelock details to the ledger with the creator's address and a creationFee
    err = gt.LedgerInstance.RecordProposal(hash, string(encryptedLog), creator, 0.0) // Added creationFee as 0.0
    if err != nil {
        return fmt.Errorf("failed to log timelock to ledger: %v", err)
    }

    fmt.Printf("Timelock for proposal %s logged to ledger.\n", proposalID)
    return nil
}


// logExecutionToLedger logs the execution details to the ledger
func (gt *GovernanceTimelock) logExecutionToLedger(proposalID, creator string) error {
    logEntry := fmt.Sprintf("Proposal %s executed at %s", proposalID, time.Now())

    // Create an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt the log entry using the encryption instance
    encryptedLog, err := encryptionInstance.EncryptData("AES", []byte(logEntry), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt execution log: %v", err)
    }

    // Generate a hash for the proposal
    hash := gt.generateProposalHash(proposalID)

    // Log the encrypted execution details to the ledger with the creator's address and a creationFee
    err = gt.LedgerInstance.RecordProposal(hash, string(encryptedLog), creator, 0.0) // Added creationFee as 0.0
    if err != nil {
        return fmt.Errorf("failed to log execution to ledger: %v", err)
    }

    fmt.Printf("Execution of proposal %s logged to ledger.\n", proposalID)
    return nil
}




// generateProposalHash generates a hash for a given proposal ID
func (gt *GovernanceTimelock) generateProposalHash(proposalID string) string {
    hashInput := fmt.Sprintf("%s%d", proposalID, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// RemoveExpiredTimelocks removes timelocks that have expired and were not executed
func (gt *GovernanceTimelock) RemoveExpiredTimelocks() {
    gt.mutex.Lock()
    defer gt.mutex.Unlock()

    for id, execution := range gt.PendingExecutions {
        if time.Now().After(execution.ExecutionTime) {
            delete(gt.PendingExecutions, id)
            fmt.Printf("Removed expired timelock for proposal %s.\n", id)
        }
    }
}
