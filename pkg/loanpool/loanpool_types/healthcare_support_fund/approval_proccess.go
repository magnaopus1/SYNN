package loanpool

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// ApprovalStatus constants representing proposal approval outcomes.
const (
	StatusPending   = "Pending"
	StatusApproved  = "Approved"
	StatusRejected  = "Rejected"
)


// NewHealthcareSupportFundApprovalProcess initializes the approval process for healthcare support fund proposals.
func NewHealthcareSupportFundApprovalProcess(ledgerInstance *ledger.Ledger, nodes []*common.AuthorityNodeTypes, encryptionService *common.Encryption) *HealthcareSupportFundApprovalProcess {
	return &HealthcareSupportFundApprovalProcess{
		Ledger:            ledgerInstance,
		Nodes:             nodes,
		ActiveProposals:   make(map[string]*HealthcareSupportFundActiveProposal),
		EncryptionService: encryptionService,
		RequeueDuration:   1 * 24 * time.Hour, // Requeue if not processed within 7 days
		MaxConfirmations:  5,
		MaxRejections:     5,
	}
}

// StartApprovalProcess starts the approval process for a new healthcare support fund proposal.
func (p *HealthcareSupportFundApprovalProcess) StartApprovalProcess(proposal *HealthcareSupportFundProposal) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the proposal is already active
	if _, exists := p.ActiveProposals[proposal.ApplicantName]; exists {
		return errors.New("this proposal is already being reviewed")
	}

	// Select initial 5 random authority nodes to review the proposal
	nodes := p.selectRandomNodes(5)

	activeProposal := &HealthcareSupportFundActiveProposal{
		ProposalID:       GenerateUniqueID(),
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

// HandleNodeDecision processes a confirmation or rejection from an authority node.
func (p *HealthcareSupportFundApprovalProcess) HandleNodeDecision(proposalID, nodeID string, confirmed bool) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	activeProposal, exists := p.getActiveProposalByID(proposalID)
	if !exists {
		return errors.New("proposal not found")
	}

	// Check if node has already responded
	if activeProposal.ConfirmedNodes[nodeID] || activeProposal.RejectedNodes[nodeID] {
		return errors.New("this node has already responded to the proposal")
	}

	// Process node decision
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
func (p *HealthcareSupportFundApprovalProcess) RequeueProposals() {
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
func (p *HealthcareSupportFundApprovalProcess) distributeProposalToNodes(proposal *HealthcareSupportFundActiveProposal) {
	for _, node := range proposal.AssignedNodes {
		p.sendProposalToNode(node, proposal.ProposalData)
	}
}

// sendProposalToNode sends an encrypted proposal to a single authority node.
func (p *HealthcareSupportFundApprovalProcess) sendProposalToNode(node *common.AuthorityNodeTypes, proposal *HealthcareSupportFundProposal) error {
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
func (p *HealthcareSupportFundApprovalProcess) finalizeProposal(proposal *HealthcareSupportFundActiveProposal) {
	// Log the final result (Approved or Rejected) in the ledger
	err := p.Ledger.RecordProposalDecision(proposal.ProposalID, proposal.Status)
	if err != nil {
		fmt.Printf("Failed to log proposal %s decision in the ledger: %v\n", proposal.ProposalID, err)
		return
	}

	// Remove the proposal from the active queue
	delete(p.ActiveProposals, proposal.ProposalData.ApplicantName)

	fmt.Printf("Proposal %s has been %s.\n", proposal.ProposalID, proposal.Status)
}

// selectRandomNodes selects a random set of nodes from the available authority nodes.
func (p *HealthcareSupportFundApprovalProcess) selectRandomNodes(count int) map[string]*common.AuthorityNodeTypes {
	selectedNodes := make(map[string]*common.AuthorityNodeTypes)
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
func (p *HealthcareSupportFundApprovalProcess) selectRandomNodeExcluding(exclude map[string]*common.AuthorityNodeTypes) *common.AuthorityNodeTypes {
	rand.Seed(time.Now().UnixNano())

	for {
		node := p.Nodes[rand.Intn(len(p.Nodes))]
		if _, exists := exclude[node.NodeID]; !exists && node.NodeStatus == "Online" {
			return node
		}
	}
}

// getActiveProposalByID retrieves an active proposal by its ID.
func (p *HealthcareSupportFundApprovalProcess) getActiveProposalByID(proposalID string) (*HealthcareSupportFundActiveProposal, bool) {
	for _, proposal := range p.ActiveProposals {
		if proposal.ProposalID == proposalID {
			return proposal, true
		}
	}
	return nil, false
}
