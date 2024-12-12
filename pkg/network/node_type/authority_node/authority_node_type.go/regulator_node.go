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

// RegulatorNodePermissions defines the immutable permissions for a regulator node.
type RegulatorNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	ViewAndDisburseLoans       bool
	FreezeWallets              bool
	ExecuteComplianceContracts bool
	AddCompliance              bool
	RemoveCompliance           bool
	VerifySyn900ID             bool
}

// RegulatorNode represents a regulator node with specific permissions and functionalities.
type RegulatorNode struct {
	NodeID            string                        // Unique identifier for the node
	KeyManager        *KeyManager                   // Key manager for managing the node's key
	Ledger            *ledger.Ledger                // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption        // Encryption service for secure communication
	NetworkManager    *network.NetworkManager       // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request    // Request list for the regulator node
	Permissions       RegulatorNodePermissions      // Fixed permissions for the regulator node
	Syn900Verifier    *syn900.Verifier              // Syn900 verifier for ID confirmation
	mutex             sync.Mutex                    // Mutex for thread-safe operations
}

// NewRegulatorNode initializes a new regulator node with the given permissions and dependencies.
func NewRegulatorNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syn900Verifier *syn900.Verifier) *RegulatorNode {
	// Define the fixed permissions for a regulator node.
	permissions := RegulatorNodePermissions{
		ViewRequestList:            true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:    true,
		ReportAuthorityNode:        true,
		ViewPrivateTransactions:    true,
		ViewAndDisburseLoans:       true,
		FreezeWallets:              true,
		ExecuteComplianceContracts: true,
		AddCompliance:              true,
		RemoveCompliance:           true,
		VerifySyn900ID:             true,
	}

	return &RegulatorNode{
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
func (rn *RegulatorNode) VerifyAndUseKey(keyID string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := rn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = rn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Regulator node %s started successfully with key %s.\n", rn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the regulator node.
func (rn *RegulatorNode) ViewRequestList() ([]*common.Request, error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range rn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the regulator node to confirm or reject a cancellation request.
func (rn *RegulatorNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := rn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	if confirm {
		err := rn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the regulator node to confirm or reject a reversal request.
func (rn *RegulatorNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := rn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	if confirm {
		err := rn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the regulator node to report another authority node.
func (rn *RegulatorNode) ReportAuthorityNode(nodeID string, reason string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	err := rn.Ledger.RecordAuthorityNodeReport(nodeID, reason, rn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the regulator node to view private transactions.
func (rn *RegulatorNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := rn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// ViewAndDisburseLoans allows the regulator node to view and disburse secured and unsecured loans from the loan pool.
func (rn *RegulatorNode) ViewAndDisburseLoans(loanID string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ViewAndDisburseLoans {
		return errors.New("permission denied: cannot view or disburse loans")
	}

	err := rn.Ledger.DisburseLoan(loanID)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %v", err)
	}

	fmt.Printf("Loan %s disbursed successfully.\n", loanID)
	return nil
}

// FreezeWallet allows the regulator node to freeze a wallet in the network.
func (rn *RegulatorNode) FreezeWallet(walletID string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.FreezeWallets {
		return errors.New("permission denied: cannot freeze wallets")
	}

	err := rn.Ledger.FreezeWallet(walletID)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet: %v", err)
	}

	fmt.Printf("Wallet %s has been frozen.\n", walletID)
	return nil
}

// ExecuteComplianceContract allows the regulator node to execute compliance contracts on the network.
func (rn *RegulatorNode) ExecuteComplianceContract(contractID string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.ExecuteComplianceContracts {
		return errors.New("permission denied: cannot execute compliance contracts")
	}

	err := rn.Ledger.ExecuteComplianceContract(contractID)
	if err != nil {
		return fmt.Errorf("failed to execute compliance contract: %v", err)
	}

	fmt.Printf("Compliance contract %s executed successfully.\n", contractID)
	return nil
}

// AddCompliance allows the regulator node to add new compliance regulations to the network.
func (rn *RegulatorNode) AddCompliance(complianceID, details string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.AddCompliance {
		return errors.New("permission denied: cannot add new compliance")
	}

	err := rn.Ledger.AddCompliance(complianceID, details)
	if err != nil {
		return fmt.Errorf("failed to add compliance: %v", err)
	}

	fmt.Printf("Compliance %s added successfully.\n", complianceID)
	return nil
}

// RemoveCompliance allows the regulator node to remove compliance regulations from the network.
func (rn *RegulatorNode) RemoveCompliance(complianceID string) error {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.RemoveCompliance {
		return errors.New("permission denied: cannot remove compliance")
	}

	err := rn.Ledger.RemoveCompliance(complianceID)
	if err != nil {
		return fmt.Errorf("failed to remove compliance: %v", err)
	}

	fmt.Printf("Compliance %s removed successfully.\n", complianceID)
	return nil
}

// VerifySyn900ID allows the regulator node to verify and confirm Syn900 ID for identity verification.
func (rn *RegulatorNode) VerifySyn900ID(identityToken string) (bool, error) {
	rn.mutex.Lock()
	defer rn.mutex.Unlock()

	if !rn.Permissions.VerifySyn900ID {
		return false, errors.New("permission denied: cannot verify Syn900 ID")
	}

	isVerified, err := rn.Syn900Verifier.VerifyIdentityToken(identityToken)
	if err != nil {
		return false, fmt.Errorf("failed to verify Syn900 ID: %v", err)
	}

	return isVerified, nil
}
