package governance

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// GovernanceContract represents the structure of a governance contract
type GovernanceContract struct {
	ContractID     string                  // Unique ID of the governance contract
	Proposals      map[string]common.GovernanceProposal     // Active proposals within the contract
	ExecutionQueue []ExecutionRecord       // Queue of actions to be executed based on the outcome of proposals
	VotingSystem   *DelegatedVoting        // Add the VotingSystem field for managing votes (assuming DelegatedVoting is the voting system)
	mutex          sync.Mutex              // Mutex for thread-safe operations
	LedgerInstance *ledger.Ledger          // Ledger instance for recording governance decisions
}

// NewDelegatedVotingSystem initializes a new delegated voting system
func NewDelegatedVotingSystem(ledgerInstance *ledger.Ledger) *DelegatedVoting {
    return &DelegatedVoting{
        LedgerInstance: ledgerInstance,
        Delegations:    make(map[string]string), // Initialize the delegation map
        Votes:          make(map[string]map[string]bool), // Initialize the vote tracking
    }
}


// NewGovernanceContract initializes a new governance contract
func NewGovernanceContract(contractID string, ledgerInstance *ledger.Ledger) *GovernanceContract {
    return &GovernanceContract{
        ContractID:     contractID,
        Proposals:      make(map[string]common.GovernanceProposal),
        VotingSystem:   NewDelegatedVotingSystem(ledgerInstance),
        ExecutionQueue: []ExecutionRecord{},
        LedgerInstance: ledgerInstance,
    }
}






// CastVote allows a user to cast a vote on a proposal
func (gc *GovernanceContract) CastVote(proposalID, voterID, vote string) error {
    gc.mutex.Lock()
    defer gc.mutex.Unlock()

    proposal, exists := gc.Proposals[proposalID]
    if !exists {
        return errors.New("proposal not found")
    }

    if time.Now().After(proposal.ExpirationTime) {
        return errors.New("voting has closed for this proposal")
    }

    // Create an instance of the Encryption struct
    encryptionInstance := &common.Encryption{}

    // Encrypt the vote using AES
    encryptedVote, err := encryptionInstance.EncryptData("AES", []byte(vote), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt vote: %v", err)
    }

    // Assuming proposal has a field to store encrypted votes by voterID
    if proposal.EncryptedDetails == "" {
        proposal.EncryptedDetails = string(encryptedVote)
    } else {
        proposal.EncryptedDetails += string(encryptedVote) // Append to store multiple votes
    }

    // Update vote counts based on the vote type
    if vote == "yes" {
        proposal.VotesFor++
    } else if vote == "no" {
        proposal.VotesAgainst++
    } else {
        return errors.New("invalid vote option")
    }

    // Update the proposal with the new vote counts
    gc.Proposals[proposalID] = proposal
    fmt.Printf("Voter %s cast a vote on proposal %s\n", voterID, proposalID)

    return nil
}


// TallyVotes tallies the votes for a proposal and adds it to the execution queue if approved
func (gc *GovernanceContract) TallyVotes(proposalID string) error {
    gc.mutex.Lock()
    defer gc.mutex.Unlock()

    // Check if the proposal exists
    proposal, exists := gc.Proposals[proposalID]
    if !exists {
        return errors.New("proposal not found")
    }

    // Ensure that voting is closed
    if time.Now().Before(proposal.ExpirationTime) {
        return errors.New("voting is still open for this proposal")
    }

    // Tally the votes by comparing VotesFor and VotesAgainst
    if proposal.VotesFor > proposal.VotesAgainst {
        // Add to the execution queue if the proposal is approved
        gc.ExecutionQueue = append(gc.ExecutionQueue, ExecutionRecord{
            ProposalID: proposalID,
            Executed:   false,
        })
        fmt.Printf("Proposal %s approved and added to the execution queue\n", proposalID)
    } else {
        fmt.Printf("Proposal %s was not approved by the community\n", proposalID)
    }

    return nil
}


// ExecuteProposal executes an approved proposal by invoking the contract execution in the virtual machine.
func (gc *GovernanceContract) ExecuteProposal(vm *common.VirtualMachine, proposalID string, parameters map[string]interface{}, encryptionInstance *common.Encryption, encryptionKey []byte) error {
    gc.mutex.Lock()
    defer gc.mutex.Unlock()

    // Find the proposal in the execution queue
    for i, record := range gc.ExecutionQueue {
        if record.ProposalID == proposalID && !record.Executed {
            proposal, exists := gc.Proposals[proposalID]
            if !exists {
                return errors.New("proposal not found")
            }

            // Execute the proposal in the virtual machine (assumed to be a contract)
            result, err := vm.ExecuteContract(proposalID, proposal.Description, "solidity", parameters, encryptionInstance, encryptionKey)
            if err != nil {
                return fmt.Errorf("execution failed for proposal %s: %v", proposalID, err)
            }

            // Mark the execution as completed
            gc.ExecutionQueue[i].Executed = true

            // Log the result of the execution (optional)
            fmt.Printf("Proposal %s executed successfully with result: %v\n", proposalID, result)

            return nil
        }
    }

    return errors.New("proposal not found in the execution queue or already executed")
}



// hashProposal generates a hash for a proposal to ensure integrity
func (gc *GovernanceContract) hashProposal(proposalID, title, description string) string {
    hashInput := fmt.Sprintf("%s%s%s", proposalID, title, description)
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
