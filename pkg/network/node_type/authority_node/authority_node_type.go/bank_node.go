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

// BankNodePermissions define what each bank node is allowed to do.
type BankNodePermissions struct {
	ViewRequestList           bool
	ConfirmOrRejectCancellation bool
	ConfirmOrRejectReversal   bool
	ReportAuthorityNode       bool
	ViewPrivateTransactions   bool
	DisburseLoans             bool
	FreezeWallets             bool
	VerifySyn900Deployment    bool
}

// BankNode represents a bank node with extended permissions and functionality specific to banking operations.
type BankNode struct {
	NodeID            string                          // Unique identifier for the node
	KeyManager        *KeyManager                     // Key manager for verifying and managing bank node keys
	Ledger            *ledger.Ledger                  // Reference to the blockchain ledger
	EncryptionService *encryption.Encryption          // Encryption service for secure communications
	NetworkManager    *network.NetworkManager         // Network manager for handling communication between nodes
	LoanPool          *loanpool.LoanPool              // Reference to loan pool for secured/unsecured loan management
	RequestList       map[string]*common.Request      // Individual request list for the bank node
	Permissions       BankNodePermissions             // Fixed permissions for the bank node
	WalletFreezeList  map[string]bool                 // Keeps track of frozen wallets
	mutex             sync.Mutex                      // Mutex for thread-safe operations
}

// NewBankNode initializes a new bank node with the given parameters and permissions.
func NewBankNode(nodeID string, keyManager *KeyManager, ledger *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, loanPool *loanpool.LoanPool) *BankNode {
	// Define the fixed permissions for a bank node.
	permissions := BankNodePermissions{
		ViewRequestList:           true,
		ConfirmOrRejectCancellation: true,
		ConfirmOrRejectReversal:   true,
		ReportAuthorityNode:       true,
		ViewPrivateTransactions:   true,
		DisburseLoans:             true,
		FreezeWallets:             true,
		VerifySyn900Deployment:    true,
	}

	return &BankNode{
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

// VerifyAndUseKey verifies the bank node key before allowing the node to function.
func (bn *BankNode) VerifyAndUseKey(keyID string) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Verify the key with the KeyManager.
	key, err := bn.KeyManager.ViewAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve key: %v", err)
	}

	// Check if the key is expired.
	if key.IsExpired || time.Now().After(key.ExpirationDate) {
		return errors.New("key has expired")
	}

	// Mark the key as used.
	err = bn.KeyManager.UseAuthorityNodeKey(keyID)
	if err != nil {
		return fmt.Errorf("failed to mark key as used: %v", err)
	}

	fmt.Printf("Bank node %s started successfully with key %s.\n", bn.NodeID, keyID)
	return nil
}

// ViewRequestList retrieves the individual request list for the bank node.
func (bn *BankNode) ViewRequestList() ([]*common.Request, error) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to view the request list.
	if !bn.Permissions.ViewRequestList {
		return nil, errors.New("permission denied: cannot view request list")
	}

	requests := []*common.Request{}
	for _, req := range bn.RequestList {
		requests = append(requests, req)
	}

	return requests, nil
}

// ConfirmCancellationRequest allows the bank node to confirm or reject a cancellation request.
func (bn *BankNode) ConfirmCancellationRequest(requestID string, confirm bool) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to confirm or reject cancellation.
	if !bn.Permissions.ConfirmOrRejectCancellation {
		return errors.New("permission denied: cannot confirm or reject cancellation")
	}

	request, exists := bn.RequestList[requestID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	// Confirm or reject the cancellation request.
	if confirm {
		err := bn.Ledger.RecordCancellation(requestID)
		if err != nil {
			return fmt.Errorf("failed to record cancellation: %v", err)
		}
		fmt.Printf("Cancellation request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Cancellation request %s rejected.\n", requestID)
	}

	return nil
}

// ConfirmReversalRequest allows the bank node to confirm or reject a reversal request.
func (bn *BankNode) ConfirmReversalRequest(requestID string, confirm bool) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to confirm or reject reversal.
	if !bn.Permissions.ConfirmOrRejectReversal {
		return errors.New("permission denied: cannot confirm or reject reversal")
	}

	request, exists := bn.RequestList[requestID]
	if !exists {
		return errors.New("reversal request not found")
	}

	// Confirm or reject the reversal request.
	if confirm {
		err := bn.Ledger.RecordReversal(requestID)
		if err != nil {
			return fmt.Errorf("failed to record reversal: %v", err)
		}
		fmt.Printf("Reversal request %s confirmed.\n", requestID)
	} else {
		fmt.Printf("Reversal request %s rejected.\n", requestID)
	}

	return nil
}

// ReportAuthorityNode allows the bank node to report another authority node.
func (bn *BankNode) ReportAuthorityNode(nodeID string, reason string) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to report other authority nodes.
	if !bn.Permissions.ReportAuthorityNode {
		return errors.New("permission denied: cannot report other authority nodes")
	}

	// Submit the report to the ledger.
	err := bn.Ledger.RecordAuthorityNodeReport(nodeID, reason, bn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to record authority node report: %v", err)
	}

	fmt.Printf("Reported authority node %s for reason: %s\n", nodeID, reason)
	return nil
}

// ViewPrivateTransactions allows the bank node to view private transactions.
func (bn *BankNode) ViewPrivateTransactions(transactionID string) (*common.Transaction, error) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to view private transactions.
	if !bn.Permissions.ViewPrivateTransactions {
		return nil, errors.New("permission denied: cannot view private transactions")
	}

	transaction, err := bn.Ledger.GetPrivateTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private transaction: %v", err)
	}

	return transaction, nil
}

// DisburseLoan allows the bank node to disburse a loan from the loan pool.
func (bn *BankNode) DisburseLoan(proposalID string, amount float64) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to disburse loans.
	if !bn.Permissions.DisburseLoans {
		return errors.New("permission denied: cannot disburse loans")
	}

	// Disburse the loan from the loan pool.
	err := bn.LoanPool.DisburseLoan(proposalID, amount)
	if err != nil {
		return fmt.Errorf("failed to disburse loan: %v", err)
	}

	fmt.Printf("Loan %s disbursed successfully for %.2f.\n", proposalID, amount)
	return nil
}

// FreezeWallet allows the bank node to freeze a wallet by its address.
func (bn *BankNode) FreezeWallet(walletAddress string) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to freeze wallets.
	if !bn.Permissions.FreezeWallets {
		return errors.New("permission denied: cannot freeze wallets")
	}

	// Freeze the wallet.
	bn.WalletFreezeList[walletAddress] = true

	// Update the ledger with the frozen wallet status.
	err := bn.Ledger.RecordWalletFreeze(walletAddress)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet in ledger: %v", err)
	}

	fmt.Printf("Wallet %s has been frozen.\n", walletAddress)
	return nil
}

// VerifySyn900Deployment allows the bank node to verify and confirm Syn900 identity deployments.
func (bn *BankNode) VerifySyn900Deployment(deploymentID string) error {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()

	// Check if the node has permission to verify Syn900 deployments.
	if !bn.Permissions.VerifySyn900Deployment {
		return errors.New("permission denied: cannot verify Syn900 deployment")
	}

	// Use the Syn900 verification system.
	err := syn900.VerifyDeployment(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to verify Syn900 deployment: %v", err)
	}

	fmt.Printf("Syn900 deployment %s verified successfully by bank node %s.\n", deploymentID, bn.NodeID)
	return nil
}
