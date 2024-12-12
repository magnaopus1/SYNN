package loanpool

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/common"
)

// ApprovalStage defines the different stages in the approval process
type ApprovalStage int

const (
	StagePublicVote ApprovalStage = iota // Public vote stage
	StageAuthorityNodes                  // Authority nodes confirmation stage
)


// NewSmallBusinessGrantApprovalProcess initializes a new grant approval process.
func NewSmallBusinessGrantApprovalProcess(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, authorityNodes []common.NodeType) *common.SmallBusinessGrantApprovalProcess {
	return &common.SmallBusinessGrantApprovalProcess{
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		Proposals:         make(map[string]*common.SmallBusinessGrantProposalApproval),
		AuthorityNodes:    authorityNodes,
		PublicVotePeriod:  21 * 24 * time.Hour, // 21 days public vote period
		AuthorityVoteTime: 72 * time.Hour,      // 72 hours for authority vote
	}
}

// StartApprovalProcess initiates the approval process for a small business grant proposal.
func (ap *common.SmallBusinessGrantApprovalProcess) StartApprovalProcess(proposal *common.SmallBusinessGrantProposal) error {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	// Check if the proposal already exists
	if _, exists := ap.Proposals[proposal.BusinessName]; exists {
		return errors.New("proposal is already under approval process")
	}

	// Create the approval object
	proposalApproval := &common.SmallBusinessGrantProposalApproval{
		Proposal:      proposal,
		PublicVotes:   make(map[string]bool),
		Stage:         StagePublicVote,
		VoteStartTime: time.Now(),
		AuthorityVotes: make(map[string]bool),
	}

	// Add the proposal to the process map
	ap.Proposals[proposal.BusinessName] = proposalApproval

	// Log the start of the approval process in the ledger
	err := ap.Ledger.RecordProposalForApproval(proposal)
	if err != nil {
		return fmt.Errorf("failed to record proposal for approval in ledger: %v", err)
	}

	fmt.Printf("Grant proposal for %s has entered Stage 1 (Public Voting).\n", proposal.BusinessName)
	return nil
}

// SubmitPublicVote allows the community to vote on a grant proposal during Stage 1.
func (ap *common.SmallBusinessGrantApprovalProcess) SubmitPublicVote(businessName string, voterAddress string, confirm bool) error {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	// Retrieve the proposal
	proposalApproval, exists := ap.Proposals[businessName]
	if !exists {
		return errors.New("proposal not found")
	}

	// Ensure the proposal is in Stage 1 (Public Vote)
	if proposalApproval.Stage != StagePublicVote {
		return errors.New("proposal is not in the public voting stage")
	}

	// Add the public vote
	proposalApproval.PublicVotes[voterAddress] = confirm

	// Check if public vote period has ended
	if time.Since(proposalApproval.VoteStartTime) > ap.PublicVotePeriod {
		// Tally the votes
		acceptVotes := 0
		rejectVotes := 0
		for _, vote := range proposalApproval.PublicVotes {
			if vote {
				acceptVotes++
			} else {
				rejectVotes++
			}
		}

		// If majority vote is to accept, move to Stage 2 (Authority Node Confirmation)
		if acceptVotes > rejectVotes && len(proposalApproval.PublicVotes) > 0 {
			proposalApproval.Stage = StageAuthorityNodes
			proposalApproval.VoteStartTime = time.Now() // Reset for Stage 2
			fmt.Printf("Grant proposal for %s has moved to Stage 2 (Authority Node Confirmation).\n", businessName)
		} else {
			proposalApproval.Stage = -1 // Rejected by public vote
			fmt.Printf("Grant proposal for %s has been rejected by public vote.\n", businessName)
		}
	}

	return nil
}

// SubmitAuthorityNodeVote allows authority nodes to vote during Stage 2.
func (ap *common.SmallBusinessGrantApprovalProcess) SubmitAuthorityNodeVote(businessName, nodeID string, confirm bool) error {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	// Retrieve the proposal
	proposalApproval, exists := ap.Proposals[businessName]
	if !exists {
		return errors.New("proposal not found")
	}

	// Ensure the proposal is in Stage 2 (Authority Node Confirmation)
	if proposalApproval.Stage != StageAuthorityNodes {
		return errors.New("proposal is not in the authority node voting stage")
	}

	// Check if the node has already voted
	if _, alreadyVoted := proposalApproval.AuthorityVotes[nodeID]; alreadyVoted {
		return errors.New("node has already voted")
	}

	// Add the authority node vote
	proposalApproval.AuthorityVotes[nodeID] = confirm
	if confirm {
		proposalApproval.ConfirmationCount++
	} else {
		proposalApproval.RejectionCount++
	}

	// Check if the required confirmations or rejections have been reached
	if proposalApproval.ConfirmationCount >= 5 {
		proposalApproval.Stage = -1 // Proposal confirmed
		err := ap.Ledger.RecordProposalApproval(proposalApproval.Proposal)
		if err != nil {
			return fmt.Errorf("failed to record approval in ledger: %v", err)
		}
		fmt.Printf("Grant proposal for %s has been confirmed by authority nodes.\n", businessName)
	} else if proposalApproval.RejectionCount >= 5 {
		proposalApproval.Stage = -1 // Proposal rejected
		fmt.Printf("Grant proposal for %s has been rejected by authority nodes.\n", businessName)
	}

	return nil
}

// MonitorAuthorityNodeTimeout automatically reassigns the proposal if an authority node fails to vote in time.
func (ap *common.SmallBusinessGrantApprovalProcess) MonitorAuthorityNodeTimeout() {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	for businessName, proposalApproval := range ap.Proposals {
		if proposalApproval.Stage == StageAuthorityNodes {
			for nodeID := range proposalApproval.AuthorityVotes {
				if time.Since(proposalApproval.VoteStartTime) > 24*time.Hour {
					// Reassign proposal to another authority node
					newNodeID := ap.Consensus.GetRandomAuthorityNode(ap.AuthorityNodes)
					delete(proposalApproval.AuthorityVotes, nodeID) // Remove original node
					ap.SubmitAuthorityNodeVote(businessName, newNodeID, false) // Send to new random node
					fmt.Printf("Reassigned authority node vote for proposal %s.\n", businessName)
				}
			}
		}
	}
}
