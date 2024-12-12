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
	"synnergy_network_demo/syn900"
)

// AuthorityNodePermissions define what each authority node type is allowed to do.
type AuthorityNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	VerifySyn900Deployment    bool
}

// AuthorityNode represents a full node with extended permissions and functionality for authority nodes.
type AuthorityNode struct {
	NodeID            string                         // Unique identifier for the node
	KeyManager        *KeyManager                    // Key manager for verifying and managing authority node keys
	Ledger            *ledger.Ledger                 // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption         // Encryption service for secure communications
	NetworkManager    *network.NetworkManager        // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request     // Individual request list for the authority node
	Permissions       AuthorityNodePermissions       // Fixed permissions for the authority node
	mutex             sync.Mutex                     // Mutex for thread-safe operations
}

// NewAuthorityNode initializes a new authority node with the given parameters and permissions.
func NewAuthorityNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager) *AuthorityNode {
	// Define the fixed permissions for an authority node.
	permissions := AuthorityNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
		VerifySyn900Deployment:    true,
	}

	return &AuthorityNode{
		NodeID:            nodeID,
		KeyManager:        keyManager,
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		RequestList:       make(map[string]*common.Request),
		Permissions:       permissions,
	}
}

// VerifyAndUseKey verifies the authority node key before allowing the node to function.
func (an *AuthorityNode) VerifyAndUseKey(keyID string) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := an.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = an.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Authority node %s started successfully with key %s.\n", an.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the authority node.
func (an *AuthorityNode) ViewRequestList() ([]*common.Request, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !an.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range an.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the authority node to confirm or reject a cancellation request.
func (an *AuthorityNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !an.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := an.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := an.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the authority node to confirm or reject a reversal request.
func (an *AuthorityNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversals.
	if !an.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := an.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := an.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the authority node to report another authority node.
func (an *AuthorityNode) ReportAuthorityNode(nodeID string, reason string) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !an.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := an.Ledger.RecordAuthorityNodeReport(nodeID, reason, an.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the authority node to view private transactions.
func (an *AuthorityNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !an.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := an.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// VerifySyn900Deployment allows the authority node to verify and confirm Syn900 identity deployments.
func (an *AuthorityNode) VerifySyn900Deployment(deploymentID string) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to verify Syn900 deployments.
	if !an.Permissions.VerifySyn900Deployment {
		return errors.New("permission denied: cannot verify Syn900 deployment")
	}

	// Use the Syn900 verification system.
	err := syn900.VerifyDeployment(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to verify Syn900 deployment: %v", err)
	}

	fmt.Printf("Syn900 deployment %s verified successfully by authority node %s.\n", deploymentID, an.NodeID)
	return nil
}
