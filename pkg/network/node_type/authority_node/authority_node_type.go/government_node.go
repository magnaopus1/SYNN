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

// GovernmentNodePermissions defines the immutable permissions for a government node.
type GovernmentNodePermissions struct {
	ViewRequestList            bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal    bool
	ReportAuthorityNode        bool
	ViewPrivateTransactions    bool
	DisburseLoans              bool
	FreezeWallets              bool
	ExecuteComplianceContracts bool
	AddCompliance              bool
	RemoveCompliance           bool
	VerifySyn900ID             bool
}

// GovernmentNode represents a government node with specific permissions and functionalities.
type GovernmentNode struct {
	NodeID            string                         // Unique identifier for the node
	KeyManager        *KeyManager                    // Key manager for managing the node's key
	Ledger            *ledger.Ledger                 // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption         // Encryption service for secure communication
	NetworkManager    *network.NetworkManager        // Network manager for handling communication between nodes
	RequestList       map[string]*common.Request     // Request list for the government node
	Permissions       GovernmentNodePermissions      // Fixed permissions for the government node
	Syn900Verifier    *syn900.Verifier               // Syn900 verifier for ID confirmation
	mutex             sync.Mutex                     // Mutex for thread-safe operations
}

// NewGovernmentNode initializes a new government node with the given permissions and dependencies.
func NewGovernmentNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, syn900Verifier *syn900.Verifier) *GovernmentNode {
	// Define the fixed permissions for a government node.
	permissions := GovernmentNodePermissions{
		ViewRequestList:            true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:    true,
		ReportAuthorityNode:        true,
		ViewPrivateTransactions:    true,
		DisburseLoans:              true,
		FreezeWallets:              true,
		ExecuteComplianceContracts: true,
		AddCompliance:              true,
		RemoveCompliance:           true,
		VerifySyn900ID:             true,
	}

	return &GovernmentNode{
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
func (gn *GovernmentNode) VerifyAndUseKey(keyID string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := gn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = gn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Government node %s started successfully with key %s.\n", gn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the government node.
func (gn *GovernmentNode) ViewRequestList() ([]*common.Request, error) {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range gn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the government node to confirm or reject a cancellation request.
func (gn *GovernmentNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := gn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	if confirm {
		err := gn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the government node to confirm or reject a reversal request.
func (gn *GovernmentNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := gn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	if confirm {
		err := gn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the government node to report another authority node.
func (gn *GovernmentNode) ReportAuthorityNode(nodeID string, reason string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	err := gn.Ledger.RecordAuthorityNodeReport(nodeID, reason, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the government node to view private transactions.
func (gn *GovernmentNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := gn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// DisburseLoan allows the government node to disburse secured or unsecured loans from the loan pool.
func (gn *GovernmentNode) DisburseLoan(loanID string, amount float64, loanType string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.DisburseLoans {
		return errors.New("permission denied: cannot disburse loans")
	}

	err := gn.Ledger.RecordLoanDisbursement(loanID, amount, loanType, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %v", err)
	}

	fmt.Printf("Loan %s disbursed successfully by node %s.\n", loanID, gn.NodeID)
	return nil
}

// FreezeWallet allows the government node to freeze a specific wallet.
func (gn *GovernmentNode) FreezeWallet(walletID string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.FreezeWallets {
		return errors.New("permission denied: cannot freeze wallets")
	}

	err := gn.Ledger.RecordWalletFreeze(walletID, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet: %v", err)
	}

	fmt.Printf("Wallet %s frozen successfully by node %s.\n", walletID, gn.NodeID)
	return nil
}

// ExecuteComplianceContract allows the government node to execute a compliance contract.
func (gn *GovernmentNode) ExecuteComplianceContract(contractID string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.ExecuteComplianceContracts {
		return errors.New("permission denied: cannot execute compliance contracts")
	}

	err := gn.Ledger.RecordComplianceContractExecution(contractID, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to execute compliance contract: %v", err)
	}

	fmt.Printf("Compliance contract %s executed by node %s.\n", contractID, gn.NodeID)
	return nil
}

// AddCompliance allows the government node to add new compliance rules.
func (gn *GovernmentNode) AddCompliance(complianceID string, details string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.AddCompliance {
		return errors.New("permission denied: cannot add compliance")
	}

	err := gn.Ledger.RecordNewCompliance(complianceID, details, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to add compliance: %v", err)
	}

	fmt.Printf("New compliance %s added by node %s.\n", complianceID, gn.NodeID)
	return nil
}

// RemoveCompliance allows the government node to remove existing compliance rules.
func (gn *GovernmentNode) RemoveCompliance(complianceID string) error {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.RemoveCompliance {
		return errors.New("permission denied: cannot remove compliance")
	}

	err := gn.Ledger.RecordComplianceRemoval(complianceID, gn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to remove compliance: %v", err)
	}

	fmt.Printf("Compliance %s removed by node %s.\n", complianceID, gn.NodeID)
	return nil
}

// VerifySyn900ID allows the government node to verify and confirm Syn900 ID for identity verification.
func (gn *GovernmentNode) VerifySyn900ID(identityToken string) (bool, error) {
	gn.mutex.Lock()
	defer gn.mutex.Unlock()

	if !gn.Permissions.VerifySyn900ID {
		return false, errors.New("permission denied: cannot verify Syn900 ID")
	}

	isVerified, err := gn.Syn900Verifier.VerifyIdentityToken(identityToken)
	if err != nil {
		return false, fmt.Errorf("failed to verify Syn900 ID: %v", err)
	}

	return isVerified, nil
}
