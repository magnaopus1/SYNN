package loanpool

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// ApprovalStatus constants representing proposal approval outcomes.
const (
	StatusPending   = "Pending"
	StatusApproved  = "Approved"
	StatusRejected  = "Rejected"
)

// AuthorityNodeType constants representing different types of authority nodes.
const (
	BankAuthority        = "Bank"
	CreditorAuthority    = "Creditor"
	CentralBankAuthority = "CentralBank"
	GovernmentAuthority  = "Government"
)


// NewSecuredLoanApprovalProcess initializes the approval process for secured loan proposals.
func NewSecuredLoanApprovalProcess(ledgerInstance *ledger.Ledger, nodes []*common.AuthorityNode, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.SecuredLoanApprovalProcess {
	return &common.SecuredLoanApprovalProcess{
		Ledger:            ledgerInstance,
		Nodes:             nodes,
		ActiveProposals:   make(map[string]*common.ActiveProposal),
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		RequeueDuration:   48 * time.Hour, // Requeue if not processed within 48 hours
		MaxConfirmations:  5,
		MaxRejections:     5,
	}
}

// StartApprovalProcess starts the approval process for a new secured loan proposal.
func (p *common.SecuredLoanApprovalProcess) StartApprovalProcess(proposal *common.SecuredLoanProposal) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the proposal is already active
	if _, exists := p.ActiveProposals[proposal.ApplicantName]; exists {
		return errors.New("this proposal is already being reviewed")
	}

	// Select initial 5 random authority nodes to review the proposal
	nodes := p.selectRandomNodes(5)

	activeProposal := &common.ActiveProposal{
		ProposalID:       common.GenerateUniqueID(),
		ProposalData:     proposal,
		ConfirmedNodes:   make(map[string]bool),
		RejectedNodes:    make(map[string]bool),
		AssignedNodes:    nodes,
		Status:           StatusPending,
		LastDistribution: time.Now(),
		ProposalDeadline: time.Now().Add(p.RequeueDuration),
	}

	p.ActiveProposals[proposal.ApplicantName] = activeProposal

	// Send the proposal to the initially selected nodes
	p.distributeProposalToNodes(activeProposal)

	return nil
}

// HandleNodeDecision processes a confirmation, rejection, and interest rate from an authority node.
func (p *common.SecuredLoanApprovalProcess) HandleNodeDecision(proposalID, nodeID string, confirmed bool, interestRate float64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	activeProposal, exists := p.getActiveProposalByID(proposalID)
	if !exists {
		return errors.New("proposal not found")
	}

	node := activeProposal.AssignedNodes[nodeID]
	if node == nil {
		return errors.New("node not assigned to this proposal")
	}

	// Ensure the node has reviewed the attached documents
	if !node.HasOpenedDocs {
		return errors.New("node has not opened the attached documents")
	}

	// Process node decision and store interest rate
	node.InterestRate = interestRate
	activeProposal.InterestRates = append(activeProposal.InterestRates, interestRate)
	activeProposal.AverageInterest = p.calculateAverageInterest(activeProposal.InterestRates)

	if confirmed {
		activeProposal.ConfirmedNodes[nodeID] = true
	} else {
		activeProposal.RejectedNodes[nodeID] = true
	}

	// Check if required confirmations or rejections have been met
	if len(activeProposal.ConfirmedNodes) >= p.MaxConfirmations {
		activeProposal.Status = StatusApproved
		p.finalizeProposal(activeProposal)
		return nil
	}

	if len(activeProposal.RejectedNodes) >= p.MaxRejections {
		activeProposal.Status = StatusRejected
		p.finalizeProposal(activeProposal)
		return nil
	}

	// Distribute the proposal to an additional random node after each decision
	additionalNode := p.selectRandomNodeExcluding(activeProposal.AssignedNodes)
	if additionalNode != nil {
		activeProposal.AssignedNodes[additionalNode.NodeID] = additionalNode
		p.sendProposalToNode(additionalNode, activeProposal.ProposalData)
	}

	return nil
}

// RequeueProposals checks for proposals that need to be redistributed after inactivity.
func (p *common.SecuredLoanApprovalProcess) RequeueProposals() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, activeProposal := range p.ActiveProposals {
		if time.Now().After(activeProposal.ProposalDeadline) && activeProposal.Status == StatusPending {
			// Redistribute proposal to new random nodes
			activeProposal.LastDistribution = time.Now()
			activeProposal.ProposalDeadline = time.Now().Add(p.RequeueDuration)
			p.distributeProposalToNodes(activeProposal)
		}
	}
}

// distributeProposalToNodes sends the proposal to the assigned authority nodes.
func (p *common.SecuredLoanApprovalProcess) distributeProposalToNodes(proposal *common.ActiveProposal) {
	for _, node := range proposal.AssignedNodes {
		p.sendProposalToNode(node, proposal.ProposalData)
	}
}

// sendProposalToNode sends an encrypted proposal to a single authority node.
func (p *common.SecuredLoanApprovalProcess) sendProposalToNode(node *common.AuthorityNode, proposal *common.SecuredLoanProposal) error {
	// Encrypt the proposal data
	encryptedProposal, err := p.EncryptionService.EncryptData(fmt.Sprintf("%v", proposal), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt proposal data: %v", err)
	}

	// Send the encrypted proposal over the network
	err = p.NetworkManager.SendProposal(node.NodeID, encryptedProposal)
	if err != nil {
		return fmt.Errorf("failed to send proposal to node %s: %v", node.NodeID, err)
	}

	fmt.Printf("Proposal sent to node %s\n", node.NodeID)
	return nil
}

// finalizeProposal finalizes the proposal and logs the result in the ledger.
func (p *common.SecuredLoanApprovalProcess) finalizeProposal(proposal *common.ActiveProposal) {
	// Log the final result (Approved or Rejected) in the ledger
	err := p.Ledger.RecordProposalDecision(proposal.ProposalID, proposal.Status)
	if err != nil {
		fmt.Printf("Failed to log proposal %s decision in the ledger: %v\n", proposal.ProposalID, err)
		return
	}

	// Remove the proposal from the active queue
	delete(p.ActiveProposals, proposal.ProposalData.ApplicantName)

	fmt.Printf("Proposal %s has been %s with an average interest rate of %.2f.\n", proposal.ProposalID, proposal.Status, proposal.AverageInterest)
}

// calculateAverageInterest calculates the running average of the interest rates from authority nodes.
func (p *common.SecuredLoanApprovalProcess) calculateAverageInterest(rates []float64) float64 {
	var total float64
	for _, rate := range rates {
		total += rate
	}
	return total / float64(len(rates))
}

// selectRandomNodes selects a random set of nodes from the available authority nodes.
func (p *common.SecuredLoanApprovalProcess) selectRandomNodes(count int) map[string]*common.AuthorityNode {
	selectedNodes := make(map[string]*common.AuthorityNode)
	rand.Seed(time.Now().UnixNano())

	for len(selectedNodes) < count {
		node := p.Nodes[rand.Intn(len(p.Nodes))]
		if node.NodeStatus == "Online" {
			selectedNodes[node.NodeID] = node
		}
	}

	return selectedNodes
}

// selectRandomNodeExcluding selects a random authority node, excluding already assigned nodes.
func (p *common.SecuredLoanApprovalProcess) selectRandomNodeExcluding(exclude map[string]*common.AuthorityNode) *common.AuthorityNode {
	rand.Seed(time.Now().UnixNano())

	for {
		node := p.Nodes[rand.Intn(len(p.Nodes))]
		if _, exists := exclude[node.NodeID]; !exists && node.NodeStatus == "Online" {
			return node
		}
	}
}

// getActiveProposalByID retrieves an active proposal by its ID.
func (p *common.SecuredLoanApprovalProcess) getActiveProposalByID(proposalID string) (*common.ActiveProposal, bool) {
	for _, proposal := range p.ActiveProposals {
		if proposal.ProposalID == proposalID {
			return proposal, true
		}
	}
	return nil, false
}
