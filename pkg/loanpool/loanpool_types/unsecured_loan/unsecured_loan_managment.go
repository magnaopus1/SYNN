package loanpool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/synnergy_consensus"
)



// NewUnsecuredLoanManagement initializes a new instance for managing unsecured loans.
func NewUnsecuredLoanManagement(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensusEngine *common.SynnergyConsensus, networkManager *common.NetworkManager, authorityNodes map[string]*common.AuthorityNode) *common.UnsecuredLoanManagement {
	return &common.UnsecuredLoanManagement{
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
		ConsensusEngine:   consensusEngine,
		NetworkManager:    networkManager,
		AuthorityNodes:    authorityNodes,
		LoanBorrowerInfo:  make(map[string]*common.BorrowerDetails),
		TermChangeRequests: make(map[string]*common.BorrowerTermChangeRequest),
	}
}

// UpdateAuthorityWallet allows an authorized authority node to update its wallet address.
func (lm *common.UnsecuredLoanManagement) UpdateAuthorityWallet(nodeID, newWalletAddress string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	node, exists := lm.AuthorityNodes[nodeID]
	if !exists {
		return errors.New("authority node not found")
	}

	// Update the wallet address.
	node.WalletAddress = newWalletAddress

	// Log the wallet update in the ledger.
	err := lm.Ledger.RecordAuthorityWalletUpdate(nodeID, newWalletAddress)
	if err != nil {
		return fmt.Errorf("failed to record authority wallet update: %v", err)
	}

	fmt.Printf("Authority node %s updated wallet address to %s.\n", nodeID, newWalletAddress)
	return nil
}

// UpdateBorrowerDetails allows a borrower to update their personal details.
func (lm *common.UnsecuredLoanManagement) UpdateBorrowerDetails(loanID, name, email, contact, walletAddress string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	borrowerDetails, exists := lm.LoanBorrowerInfo[loanID]
	if !exists {
		return errors.New("borrower details not found for this loan")
	}

	// Update borrower details.
	borrowerDetails.BorrowerName = name
	borrowerDetails.BorrowerEmail = email
	borrowerDetails.BorrowerContact = contact
	borrowerDetails.WalletAddress = walletAddress

	// Log the borrower detail update in the ledger.
	err := lm.Ledger.RecordBorrowerDetailUpdate(loanID, borrowerDetails)
	if err != nil {
		return fmt.Errorf("failed to record borrower detail update: %v", err)
	}

	fmt.Printf("Borrower details for loan %s updated.\n", loanID)
	return nil
}

// RequestTermChange allows a borrower to request a change in loan terms.
func (lm *common.UnsecuredLoanManagement) RequestTermChange(loanID, newTerms string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	// Ensure the loan exists.
	_, exists := lm.LoanBorrowerInfo[loanID]
	if !exists {
		return errors.New("borrower details not found for this loan")
	}

	// Create the term change request.
	termChangeRequest := &common.BorrowerTermChangeRequest{
		LoanID:           loanID,
		RequestedTerms:   newTerms,
		ApprovalStatus:   "Pending",
		ConfirmedNodes:   make(map[string]bool),
		RejectedNodes:    make(map[string]bool),
		AssignedNodes:    lm.selectRandomAuthorityNodes(5),
		LastDistribution: time.Now(),
		RequeueDeadline:  time.Now().Add(24 * time.Hour),
	}

	// Store the term change request.
	lm.TermChangeRequests[loanID] = termChangeRequest

	// Send the term change request to the selected authority nodes for review.
	lm.distributeTermChangeRequestToNodes(termChangeRequest)

	fmt.Printf("Term change request for loan %s submitted for approval.\n", loanID)
	return nil
}

// HandleAuthorityDecision processes a confirmation or rejection from an authority node.
func (lm *common.UnsecuredLoanManagement) HandleAuthorityDecision(loanID, nodeID string, approved bool) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	termChangeRequest, exists := lm.TermChangeRequests[loanID]
	if !exists {
		return errors.New("term change request not found for this loan")
	}

	// Check if node has already responded.
	if termChangeRequest.ConfirmedNodes[nodeID] || termChangeRequest.RejectedNodes[nodeID] {
		return errors.New("this node has already responded to the term change request")
	}

	// Process the decision.
	if approved {
		termChangeRequest.ConfirmedNodes[nodeID] = true
	} else {
		termChangeRequest.RejectedNodes[nodeID] = true
	}

	// Check if required confirmations or rejections have been met (4 out of 5 required).
	if len(termChangeRequest.ConfirmedNodes) >= 4 {
		termChangeRequest.ApprovalStatus = "Accepted"
		lm.finalizeTermChange(loanID, termChangeRequest)
		return nil
	}

	if len(termChangeRequest.RejectedNodes) >= 4 {
		termChangeRequest.ApprovalStatus = "Rejected"
		lm.finalizeTermChange(loanID, termChangeRequest)
		return nil
	}

	// Distribute the term change request to an additional node if needed.
	additionalNode := lm.selectRandomAuthorityNodeExcluding(termChangeRequest.AssignedNodes)
	if additionalNode != nil {
		termChangeRequest.AssignedNodes[additionalNode.NodeID] = additionalNode
		lm.sendTermChangeRequestToNode(additionalNode, termChangeRequest)
	}

	return nil
}

// distributeTermChangeRequestToNodes distributes the term change request to assigned authority nodes.
func (lm *common.UnsecuredLoanManagement) distributeTermChangeRequestToNodes(request *common.BorrowerTermChangeRequest) {
	for _, node := range request.AssignedNodes {
		lm.sendTermChangeRequestToNode(node, request)
	}
}

// sendTermChangeRequestToNode sends a term change request to a single authority node.
func (lm *common.UnsecuredLoanManagement) sendTermChangeRequestToNode(node *common.AuthorityNode, request *common.BorrowerTermChangeRequest) error {
	// Encrypt the request data.
	encryptedData, err := lm.EncryptionService.EncryptData(fmt.Sprintf("%v", request), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt term change request data: %v", err)
	}

	// Send the encrypted request via the network manager.
	err = lm.NetworkManager.SendTermChangeRequest(node.NodeID, encryptedData)
	if err != nil {
		return fmt.Errorf("failed to send term change request to node %s: %v", node.NodeID, err)
	}

	fmt.Printf("Term change request sent to node %s.\n", node.NodeID)
	return nil
}

// finalizeTermChange finalizes the term change request and logs it in the ledger.
func (lm *common.UnsecuredLoanManagement) finalizeTermChange(loanID string, request *common.BorrowerTermChangeRequest) {
	// Log the term change decision (Accepted/Rejected) in the ledger.
	err := lm.Ledger.RecordTermChangeDecision(loanID, request.ApprovalStatus)
	if err != nil {
		fmt.Printf("Failed to record term change decision for loan %s: %v\n", loanID, err)
	}

	// Remove the term change request from active requests.
	delete(lm.TermChangeRequests, loanID)

	fmt.Printf("Term change request for loan %s has been %s.\n", loanID, request.ApprovalStatus)
}

// selectRandomAuthorityNodes selects random authority nodes for reviewing term change requests.
func (lm *common.UnsecuredLoanManagement) selectRandomAuthorityNodes(count int) map[string]*common.AuthorityNode {
	selectedNodes := make(map[string]*common.AuthorityNode)
	for len(selectedNodes) < count {
		node := lm.AuthorityNodes[common.GenerateRandomNodeID()]
		if node.NodeStatus == "Online" {
			selectedNodes[node.NodeID] = node
		}
	}
	return selectedNodes
}

// selectRandomAuthorityNodeExcluding selects a random authority node excluding certain nodes.
func (lm *common.UnsecuredLoanManagement) selectRandomAuthorityNodeExcluding(exclude map[string]*common.AuthorityNode) *common.AuthorityNode {
	for {
		node := lm.AuthorityNodes[common.GenerateRandomNodeID()]
		if _, exists := exclude[node.NodeID]; !exists && node.NodeStatus == "Online" {
			return node
		}
	}
}
