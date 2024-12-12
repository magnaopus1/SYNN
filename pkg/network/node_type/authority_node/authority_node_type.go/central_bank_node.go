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
	"synnergy_network_demo/loanpool"
)

// CentralBankNodePermissions defines what each central bank node is allowed to do.
type CentralBankNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	FreezeWallets             bool
	VerifySyn900Deployment    bool
}

// CentralBankNode represents a central bank node with extended permissions and functionalities specific to central bank operations.
type CentralBankNode struct {
	NodeID            string                             // Unique identifier for the node
	KeyManager        *KeyManager                        // Key manager for verifying and managing central bank node keys
	Ledger            *ledger.Ledger                     // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption             // Encryption service for secure communications
	NetworkManager    *network.NetworkManager            // Network manager for handling communication between nodes
	LoanPool          *loanpool.LoanPool                 // Reference to loan pool for secured/unsecured loan management
	RequestList       map[string]*common.Request         // Individual request list for the central bank node
	Permissions       CentralBankNodePermissions         // Fixed permissions for the central bank node
	WalletFreezeList  map[string]bool                    // Keeps track of frozen wallets
	mutex             sync.Mutex                         // Mutex for thread-safe operations
}

// NewCentralBankNode initializes a new central bank node with the given parameters and permissions.
func NewCentralBankNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, loanPool *loanpool.LoanPool) *CentralBankNode {
	// Define the fixed permissions for a central bank node.
	permissions := CentralBankNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
		DisburseLoans:             true,
		FreezeWallets:             true,
		VerifySyn900Deployment:    true,
	}

	return &CentralBankNode{
		NodeID:            nodeID,
		KeyManager:        keyManager,
		Ledger:            ledger,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		LoanPool:          loanPool,
		RequestList:       make(map[string]*common.Request),
		Permissions:       permissions,
		WalletFreezeList:  make(map[string]bool),
	}
}

// VerifyAndUseKey verifies the central bank node key before allowing the node to function.
func (cbn *CentralBankNode) VerifyAndUseKey(keyID string) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := cbn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = cbn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Central bank node %s started successfully with key %s.\n", cbn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the central bank node.
func (cbn *CentralBankNode) ViewRequestList() ([]*common.Request, error) {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !cbn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range cbn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the central bank node to confirm or reject a cancellation request.
func (cbn *CentralBankNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !cbn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := cbn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := cbn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the central bank node to confirm or reject a reversal request.
func (cbn *CentralBankNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversal.
	if !cbn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := cbn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := cbn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the central bank node to report another authority node.
func (cbn *CentralBankNode) ReportAuthorityNode(nodeID string, reason string) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !cbn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := cbn.Ledger.RecordAuthorityNodeReport(nodeID, reason, cbn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the central bank node to view private transactions.
func (cbn *CentralBankNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !cbn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := cbn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// DisburseLoan allows the central bank node to disburse a loan from the loan pool.
func (cbn *CentralBankNode) DisburseLoan(proposalID string, amount float64) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to disburse loans.
	if !cbn.Permissions.DisburseLoans {
		return errors.New("permission denied: cannot disburse loans")
	}

	// Disburse the loan from the loan pool.
	err := cbn.LoanPool.DisburseLoan(proposalID, amount)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %v", err)
	}

	fmt.Printf("Loan %s disbursed successfully for %.2f.\n", proposalID, amount)
	return nil
}

// FreezeWallet allows the central bank node to freeze a wallet by its address.
func (cbn *CentralBankNode) FreezeWallet(walletAddress string) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to freeze wallets.
	if !cbn.Permissions.FreezeWallets {
		return errors.New("permission denied: cannot freeze wallets")
	}

	// Freeze the wallet.
	cbn.WalletFreezeList[walletAddress] = true

	// Update the ledger with the frozen wallet status.
	err := cbn.Ledger.RecordWalletFreeze(walletAddress)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet in ledger: %v", err)
	}

	fmt.Printf("Wallet %s frozen successfully.\n", walletAddress)
	return nil
}

// VerifySyn900Deployment allows the central bank node to verify a Syn900 deployment.
func (cbn *CentralBankNode) VerifySyn900Deployment(deploymentID string) error {
	cbn.mutex.Lock()
	defer cbn.mutex.Unlock()

	// Check if the node has permission to verify Syn900 deployments.
	if !cbn.Permissions.VerifySyn900Deployment {
		return errors.New("permission denied: cannot verify Syn900 deployment")
	}

	// Use the Syn900 verification system.
	err := syn900.VerifyDeployment(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to verify Syn900 deployment: %v", err)
	}

	fmt.Printf("Syn900 deployment %s verified successfully by central bank node %s.\n", deploymentID, cbn.NodeID)
	return nil
}
