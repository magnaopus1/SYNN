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
)

// ElectedAuthorityNodePermissions defines what each elected authority node can do.
type ElectedAuthorityNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
}

// ElectedAuthorityNode represents an elected authority node with specific permissions and functionalities.
type ElectedAuthorityNode struct {
	NodeID            string                              // Unique identifier for the node
	KeyManager        *KeyManager                         // Key manager for managing the node's key
	Ledger            *ledger.Ledger                      // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption              // Encryption service for secure communication
	NetworkManager    *network.NetworkManager             // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request          // Request list for the elected authority node
	Permissions       ElectedAuthorityNodePermissions     // Fixed permissions for the elected authority node
	mutex             sync.Mutex                          // Mutex for thread-safe operations
}

// NewElectedAuthorityNode initializes a new elected authority node with the given permissions and dependencies.
func NewElectedAuthorityNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager) *ElectedAuthorityNode {
	// Define the fixed permissions for an elected authority node.
	permissions := ElectedAuthorityNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
	}

	return &ElectedAuthorityNode{
		NodeID:            nodeID,
		KeyManager:        keyManager,
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		RequestList:       make(map[string]*common.Request),
		Permissions:       permissions,
	}
}

// VerifyAndUseKey verifies the node's key before allowing the node to function.
func (ean *ElectedAuthorityNode) VerifyAndUseKey(keyID string) error {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := ean.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = ean.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Elected authority node %s started successfully with key %s.\n", ean.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the elected authority node.
func (ean *ElectedAuthorityNode) ViewRequestList() ([]*common.Request, error) {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !ean.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range ean.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the elected authority node to confirm or reject a cancellation request.
func (ean *ElectedAuthorityNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !ean.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := ean.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := ean.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the elected authority node to confirm or reject a reversal request.
func (ean *ElectedAuthorityNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversal.
	if !ean.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := ean.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := ean.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the elected authority node to report another authority node.
func (ean *ElectedAuthorityNode) ReportAuthorityNode(nodeID string, reason string) error {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !ean.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := ean.Ledger.RecordAuthorityNodeReport(nodeID, reason, ean.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the elected authority node to view private transactions.
func (ean *ElectedAuthorityNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	ean.mutex.Lock()
	defer ean.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !ean.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := ean.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}
