package authority_node

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/common"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
	"synnergy_network_demo/loanpool"
	"synnergy_network_demo/syn900"
)

// CreditProviderNodePermissions defines what each credit provider node is allowed to do.
type CreditProviderNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	VerifySyn900Deployment    bool
}

// CreditProviderNode represents a credit provider node with extended permissions and functionalities.
type CreditProviderNode struct {
	NodeID            string                               // Unique identifier for the node
	KeyManager        *KeyManager                          // Key manager for verifying and managing credit provider node keys
	Ledger            *ledger.Ledger                       // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption               // Encryption service for secure communications
	NetworkManager    *network.NetworkManager              // Network manager for handling communication between nodes
	LoanPool          *loanpool.LoanPool                   // Reference to loan pool for secured/unsecured loan management
	RequestList       map[string]*common.Request           // Individual request list for the credit provider node
	Permissions       CreditProviderNodePermissions        // Fixed permissions for the credit provider node
	mutex             sync.Mutex                           // Mutex for thread-safe operations
}

// NewCreditProviderNode initializes a new credit provider node with the given parameters and permissions.
func NewCreditProviderNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, loanPool *loanpool.LoanPool) *CreditProviderNode {
	// Define the fixed permissions for a credit provider node.
	permissions := CreditProviderNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
		DisburseLoans:             true,
		VerifySyn900Deployment:    true,
	}

	return &CreditProviderNode{
		NodeID:            nodeID,
		KeyManager:        keyManager,
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		LoanPool:          loanPool,
		RequestList:       make(map[string]*common.Request),
		Permissions:       permissions,
	}
}

// VerifyAndUseKey verifies the credit provider node key before allowing the node to function.
func (cpn *CreditProviderNode) VerifyAndUseKey(keyID string) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := cpn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = cpn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Credit provider node %s started successfully with key %s.\n", cpn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the credit provider node.
func (cpn *CreditProviderNode) ViewRequestList() ([]*common.Request, error) {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !cpn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range cpn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the credit provider node to confirm or reject a cancellation request.
func (cpn *CreditProviderNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !cpn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := cpn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := cpn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the credit provider node to confirm or reject a reversal request.
func (cpn *CreditProviderNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversal.
	if !cpn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := cpn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := cpn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the credit provider node to report another authority node.
func (cpn *CreditProviderNode) ReportAuthorityNode(nodeID string, reason string) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !cpn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := cpn.Ledger.RecordAuthorityNodeReport(nodeID, reason, cpn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the credit provider node to view private transactions.
func (cpn *CreditProviderNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !cpn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := cpn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// DisburseLoan allows the credit provider node to disburse a loan from the loan pool.
func (cpn *CreditProviderNode) DisburseLoan(proposalID string, amount float64) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to disburse loans.
	if !cpn.Permissions.DisburseLoans {
		return errors.New("permission denied: cannot disburse loans")
	}

	// Disburse the loan from the loan pool.
	err := cpn.LoanPool.DisburseLoan(proposalID, amount)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %v", err)
	}

	fmt.Printf("Loan %s disbursed successfully for %.2f.\n", proposalID, amount)
	return nil
}

// VerifySyn900Deployment allows the credit provider node to verify a Syn900 deployment.
func (cpn *CreditProviderNode) VerifySyn900Deployment(deploymentID string) error {
	cpn.mutex.Lock()
	defer cpn.mutex.Unlock()

	// Check if the node has permission to verify Syn900 deployments.
	if !cpn.Permissions.VerifySyn900Deployment {
		return errors.New("permission denied: cannot verify Syn900 deployment")
	}

	// Use the Syn900 verification system.
	err := syn900.VerifyDeployment(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to verify Syn900 deployment: %v", err)
	}

	fmt.Printf("Syn900 deployment %s verified successfully by credit provider node %s.\n", deploymentID, cpn.NodeID)
	return nil
}
