package authority_node

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
	"synnergy_network_demo/syn900"
)

// VotingResult represents the result of the voting process.
type VotingResult struct {
	ConfirmedVotes int
	RejectedVotes  int
	Status         string // "Accepted" or "Rejected"
}

// AuthorityNodeAcceptanceManager manages the process for voting on authority node proposals.
type AuthorityNodeAcceptanceManager struct {
	mutex             sync.Mutex
	Ledger            *ledger.Ledger
	EncryptionService *encryption.Encryption
	NetworkManager    *network.NetworkManager
	IdentityVerifier  *syn900.Verifier
	KeyDisburser      *KeyDisburser // Manages key disbursement for accepted proposals
}

// NewAuthorityNodeAcceptanceManager initializes a new manager.
func NewAuthorityNodeAcceptanceManager(ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, verifier *syn900.Verifier, keyDisburser *KeyDisburser) *AuthorityNodeAcceptanceManager {
	return &AuthorityNodeAcceptanceManager{
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		IdentityVerifier:  verifier,
		KeyDisburser:      keyDisburser,
	}
}

// Randomly selects a set number of authority nodes for voting.
func (an *AuthorityNodeAcceptanceManager) selectRandomAuthorityNodes(nodeType string, count int) ([]string, error) {
	nodes := an.NetworkManager.GetAuthorityNodes(nodeType)
	if len(nodes) == 0 {
		return nil, errors.New("no authority nodes available")
	}

	selectedNodes := []string{}
	for len(selectedNodes) < count {
		randomNode := nodes[rand.Intn(len(nodes))]
		if !contains(selectedNodes, randomNode) {
			selectedNodes = append(selectedNodes, randomNode)
		}
	}

	return selectedNodes, nil
}

// ProcessProposal processes the acceptance or rejection of a node proposal.
func (an *AuthorityNodeAcceptanceManager) ProcessProposal(proposal interface{}) (*VotingResult, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	var selectedNodes []string
	var err error
	votingResult := &VotingResult{}

	switch p := proposal.(type) {
	case *AuthorityNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 7)
		votingResult, err = an.conductVoting(selectedNodes, 7, 3)

	case *BankNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 15)
		votingResult, err = an.conductVoting(selectedNodes, 15, 2)

	case *CentralBankNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 25)
		votingResult, err = an.conductVoting(selectedNodes, 25, 2)

	case *CreditProviderNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 10)
		votingResult, err = an.conductVoting(selectedNodes, 10, 2)

	case *ElectedAuthorityNodeProposal:
		// Elected authority node requires both general and authority node votes.
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 10)
		if err != nil {
			return nil, err
		}
		votingResult, err = an.conductVoting(selectedNodes, 200, 4)

	case *ExchangeNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 15)
		votingResult, err = an.conductVoting(selectedNodes, 15, 2)

	case *GovernmentNodeProposal:
		// Government node requires both general and authority node votes.
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 15)
		if err != nil {
			return nil, err
		}
		votingResult, err = an.conductVoting(selectedNodes, 2500, 15)

	case *MilitaryNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 10)
		votingResult, err = an.conductVoting(selectedNodes, 1250, 10)

	case *RegulatorNodeProposal:
		selectedNodes, err = an.selectRandomAuthorityNodes(AuthorityNode, 15)
		votingResult, err = an.conductVoting(selectedNodes, 15, 2)

	default:
		return nil, errors.New("invalid proposal type")
	}

	if err != nil {
		return nil, err
	}

	// Record the result in the ledger.
	if err := an.Ledger.RecordVotingResult(proposal, votingResult); err != nil {
		return nil, fmt.Errorf("failed to record voting result: %v", err)
	}

	// If the proposal is accepted, disburse the key.
	if votingResult.Status == "Accepted" {
		err := an.disburseKey(proposal)
		if err != nil {
			return nil, fmt.Errorf("key disbursement failed: %v", err)
		}
	}

	return votingResult, nil
}

// Conducts the voting process for the selected authority nodes.
func (an *AuthorityNodeAcceptanceManager) conductVoting(selectedNodes []string, requiredVotes int, maxElectedAuthorityVotes int) (*VotingResult, error) {
	confirmedVotes := 0
	rejectedVotes := 0
	electedAuthorityVotes := 0

	for _, node := range selectedNodes {
		vote, nodeType := an.NetworkManager.RequestVote(node)

		if vote == "confirm" {
			if nodeType == ElectedAuthority && electedAuthorityVotes < maxElectedAuthorityVotes {
				confirmedVotes++
				electedAuthorityVotes++
			} else if nodeType != ElectedAuthority {
				confirmedVotes++
			}
		} else if vote == "reject" {
			rejectedVotes++
		}

		// Check if we have enough votes for a decision.
		if confirmedVotes >= requiredVotes {
			return &VotingResult{ConfirmedVotes: confirmedVotes, Status: "Accepted"}, nil
		}
		if rejectedVotes >= requiredVotes {
			return &VotingResult{RejectedVotes: rejectedVotes, Status: "Rejected"}, nil
		}
	}

	return &VotingResult{ConfirmedVotes: confirmedVotes, RejectedVotes: rejectedVotes}, nil
}

// Helper function to check if a node has already been selected.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// disburseKey disburses the key to the applicant upon successful acceptance of their proposal.
func (an *AuthorityNodeAcceptanceManager) disburseKey(proposal interface{}) error {
	var applicantWallet string
	var nodeType string

	switch p := proposal.(type) {
	case *AuthorityNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = AuthorityNode
	case *BankNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = BankNode
	case *CentralBankNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = CentralBankNode
	case *CreditProviderNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = CreditProviderNode
	case *ElectedAuthorityNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = ElectedAuthority
	case *ExchangeNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = ExchangeNode
	case *GovernmentNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = GovernmentNode
	case *MilitaryNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = MilitaryNode
	case *RegulatorNodeProposal:
		applicantWallet = p.ApplicantWallet
		nodeType = RegulatorNode
	default:
		return errors.New("invalid proposal type for key disbursement")
	}

	// Create and disburse the key via the KeyDisburser.
	key, err := an.KeyDisburser.CreateKey(nodeType)
	if err != nil {
		return fmt.Errorf("failed to create key: %v", err)
	}

	// Register the key and mark it as disbursed.
	if err := an.Ledger.RecordKeyDisbursement(applicantWallet, key); err != nil {
		return fmt.Errorf("failed to record key disbursement in ledger: %v", err)
	}

	fmt.Printf("Key disbursed to applicant %s for node type %s.\n", applicantWallet, nodeType)
	return nil
}
