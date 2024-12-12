package governance

import (
    "fmt"
    "time"
    "errors"
    "crypto/sha256"
    "encoding/hex"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewGovernanceExecution initializes the governance execution system
func NewGovernanceExecution(ledgerInstance *ledger.Ledger) *GovernanceExecution {
    return &GovernanceExecution{
        ExecutionQueue:   []ExecutionRecord{},
        LedgerInstance:   ledgerInstance,
    }
}

// ScheduleExecution adds a proposal to the execution queue after it has been approved
func (ge *GovernanceExecution) ScheduleExecution(proposalID string) error {
    ge.mutex.Lock()
    defer ge.mutex.Unlock()

    // Check if the proposal already exists in the queue
    for _, record := range ge.ExecutionQueue {
        if record.ProposalID == proposalID {
            return errors.New("proposal already scheduled for execution")
        }
    }

    // Add the proposal to the execution queue
    ge.ExecutionQueue = append(ge.ExecutionQueue, ExecutionRecord{
        ProposalID: proposalID,
        Executed:   false,
        Timestamp:  time.Now(),
    })

    fmt.Printf("Proposal %s scheduled for execution.\n", proposalID)
    return nil
}

// ExecuteScheduledProposals runs all proposals in the queue that are pending execution using the virtual machine
func (ge *GovernanceExecution) ExecuteScheduledProposals(vm *common.VirtualMachine, contract *GovernanceContract, encryptionInstance *common.Encryption, encryptionKey []byte) error {
    ge.mutex.Lock()
    defer ge.mutex.Unlock()

    for i, record := range ge.ExecutionQueue {
        if !record.Executed {
            // Fetch the proposal from the governance contract
            proposal, exists := contract.Proposals[record.ProposalID]
            if !exists {
                return fmt.Errorf("proposal %s not found", record.ProposalID)
            }

            // Prepare parameters for execution (based on the governance contract logic)
            parameters := map[string]interface{}{
                "ProposalID":    proposal.ProposalID,
                "Title":         proposal.Title,
                "Description":   proposal.Description,
                "VotesFor":      proposal.VotesFor,
                "VotesAgainst":  proposal.VotesAgainst,
                "ExpirationTime": proposal.ExpirationTime,
            }

            // Execute the proposal in the virtual machine
            _, err := vm.ExecuteContract(proposal.ProposalID, proposal.Description, "solidity", parameters, encryptionInstance, encryptionKey)
            if err != nil {
                return fmt.Errorf("failed to execute proposal %s: %v", record.ProposalID, err)
            }

            // Mark the proposal as executed
            ge.ExecutionQueue[i].Executed = true
            ge.ExecutionQueue[i].Timestamp = time.Now()

            // Log the execution in the ledger
            err = ge.logExecutionToLedger(record)
            if err != nil {
                return fmt.Errorf("failed to log execution to the ledger: %v", err)
            }

            fmt.Printf("Proposal %s executed successfully.\n", record.ProposalID)
        }
    }

    return nil
}

// logExecutionToLedger logs the execution of a proposal to the ledger
func (ge *GovernanceExecution) logExecutionToLedger(record ExecutionRecord) error {
    // Create an Encryption instance
    encryptionInstance := &common.Encryption{}

    // Prepare the execution record as a string for encryption
    executionDetails := fmt.Sprintf("%+v", record) 

    // Encrypt the execution record for secure storage (if needed elsewhere)
    _, err := encryptionInstance.EncryptData("AES", []byte(executionDetails), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt execution record: %v", err)
    }

    // Generate a hash for integrity verification
    hash := ge.generateExecutionHash(record)

    // Store the execution record in the ledger (only the hash)
    err = ge.LedgerInstance.RecordExecution(hash)
    if err != nil {
        return fmt.Errorf("failed to record execution in the ledger: %v", err)
    }

    fmt.Printf("Execution of proposal %s logged in the ledger.\n", record.ProposalID)
    return nil
}




// generateExecutionHash generates a hash for the execution record
func (ge *GovernanceExecution) generateExecutionHash(record ExecutionRecord) string {
    hashInput := fmt.Sprintf("%s%d", record.ProposalID, record.Timestamp.UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// CheckExecutionStatus returns the status of a proposal's execution
func (ge *GovernanceExecution) CheckExecutionStatus(proposalID string) (bool, error) {
    ge.mutex.Lock()
    defer ge.mutex.Unlock()

    for _, record := range ge.ExecutionQueue {
        if record.ProposalID == proposalID {
            return record.Executed, nil
        }
    }

    return false, errors.New("proposal not found in execution queue")
}

// RemoveExecutedProposals removes executed proposals from the execution queue
func (ge *GovernanceExecution) RemoveExecutedProposals() {
    ge.mutex.Lock()
    defer ge.mutex.Unlock()

    var updatedQueue []ExecutionRecord
    for _, record := range ge.ExecutionQueue {
        if !record.Executed {
            updatedQueue = append(updatedQueue, record)
        }
    }

    ge.ExecutionQueue = updatedQueue
    fmt.Println("Executed proposals removed from the queue.")
}
