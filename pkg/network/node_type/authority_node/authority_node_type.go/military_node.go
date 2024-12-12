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

// MilitaryNodePermissions defines the immutable permissions for a military node.
type MilitaryNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	VerifySyn900ID             bool
}

// MilitaryNode represents a military node with specific permissions and functionalities.
type MilitaryNode struct {
	NodeID            string                        // Unique identifier for the node
	KeyManager        *KeyManager                   // Key manager for managing the node's key
	Ledger            *ledger.Ledger                // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption        // Encryption service for secure communication
	NetworkManager    *network.NetworkManager       // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request    // Request list for the military node
	Permissions       MilitaryNodePermissions       // Fixed permissions for the military node
	Syn900Verifier    *syn900.Verifier              // Syn900 verifier for ID confirmation
	mutex             sync.Mutex                    // Mutex for thread-safe operations
}

// NewMilitaryNode initializes a new military node with the given permissions and dependencies.
func NewMilitaryNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syn900Verifier *syn900.Verifier) *MilitaryNode {
	// Define the fixed permissions for a military node.
	permissions := MilitaryNodePermissions{
		ViewRequestList:            true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:    true,
		ReportAuthorityNode:        true,
		ViewPrivateTransactions:    true,
		VerifySyn900ID:             true,
	}

	return &MilitaryNode{
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
func (mn *MilitaryNode) VerifyAndUseKey(keyID string) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := mn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = mn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Military node %s started successfully with key %s.\n", mn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the military node.
func (mn *MilitaryNode) ViewRequestList() ([]*common.Request, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range mn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the military node to confirm or reject a cancellation request.
func (mn *MilitaryNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := mn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	if confirm {
		err := mn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the military node to confirm or reject a reversal request.
func (mn *MilitaryNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := mn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	if confirm {
		err := mn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the military node to report another authority node.
func (mn *MilitaryNode) ReportAuthorityNode(nodeID string, reason string) error {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	err := mn.Ledger.RecordAuthorityNodeReport(nodeID, reason, mn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the military node to view private transactions.
func (mn *MilitaryNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := mn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// VerifySyn900ID allows the military node to verify and confirm Syn900 ID for identity verification.
func (mn *MilitaryNode) VerifySyn900ID(identityToken string) (bool, error) {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()

	if !mn.Permissions.VerifySyn900ID {
		return false, errors.New("permission denied: cannot verify Syn900 ID")
	}

	isVerified, err := mn.Syn900Verifier.VerifyIdentityToken(identityToken)
	if err != nil {
		return false, fmt.Errorf("failed to verify Syn900 ID: %v", err)
	}

	return isVerified, nil
}
