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

// ExchangeNodePermissions defines what each exchange node can do.
type ExchangeNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	VerifySyn900ID            bool
}

// ExchangeNode represents an exchange node with specific permissions and functionalities.
type ExchangeNode struct {
	NodeID            string                          // Unique identifier for the node
	KeyManager        *KeyManager                     // Key manager for managing the node's key
	Ledger            *ledger.Ledger                  // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption          // Encryption service for secure communication
	NetworkManager    *network.NetworkManager         // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request      // Request list for the exchange node
	Permissions       ExchangeNodePermissions         // Fixed permissions for the exchange node
	Syn900Verifier    *syn900.Verifier                // Syn900 verifier for ID confirmation
	mutex             sync.Mutex                      // Mutex for thread-safe operations
}

// NewExchangeNode initializes a new exchange node with the given permissions and dependencies.
func NewExchangeNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syn900Verifier *syn900.Verifier) *ExchangeNode {
	// Define the fixed permissions for an exchange node.
	permissions := ExchangeNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
		VerifySyn900ID:            true,
	}

	return &ExchangeNode{
		NodeID:            nodeID,
		KeyManager:        keyManager,
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		RequestList:       make(map[string]*common.Request),
		Permissions:       permissions,
		Syn900Verifier:    syn900Verifier,
	}
}

// VerifyAndUseKey verifies the node's key before allowing the node to function.
func (en *ExchangeNode) VerifyAndUseKey(keyID string) error {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := en.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = en.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Exchange node %s started successfully with key %s.\n", en.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the exchange node.
func (en *ExchangeNode) ViewRequestList() ([]*common.Request, error) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !en.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range en.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the exchange node to confirm or reject a cancellation request.
func (en *ExchangeNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !en.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := en.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := en.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the exchange node to confirm or reject a reversal request.
func (en *ExchangeNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversal.
	if !en.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := en.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := en.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the exchange node to report another authority node.
func (en *ExchangeNode) ReportAuthorityNode(nodeID string, reason string) error {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !en.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := en.Ledger.RecordAuthorityNodeReport(nodeID, reason, en.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the exchange node to view private transactions.
func (en *ExchangeNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !en.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := en.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// VerifySyn900ID allows the exchange node to verify and confirm Syn900 ID for identity verification.
func (en *ExchangeNode) VerifySyn900ID(identityToken string) (bool, error) {
	en.mutex.Lock()
	defer en.mutex.Unlock()

	// Check if the node has permission to verify Syn900 ID.
	if !en.Permissions.VerifySyn900ID {
		return false, errors.New("permission denied: cannot verify Syn900 ID")
	}

	// Verify the Syn900 ID.
	isVerified, err := en.Syn900Verifier.VerifyIdentityToken(identityToken)
	if err != nil {
		return false, fmt.Errorf("failed to verify Syn900 ID: %v", err)
	}

	return isVerified, nil
}
