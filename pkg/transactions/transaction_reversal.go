package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// ReversalRequest represents a request to reverse a transaction.
type ReversalRequest struct {
	ID                string          // Unique ID for the reversal request
	TransactionID     string          // Unique transaction ID
	WalletID          string          // Wallet ID requesting the reversal
	UserID            string          // User ID requesting the reversal
	Reason            string          // Reason for reversal
	ContactDetails    string          // Contact details of the user
	SubmittedEvidence []byte          // Evidence documents submitted by the user
	Status            string          // Current status of the request (e.g., pending, approved, rejected)
	Timestamp         time.Time       // Time when the request was created
	RequiredApprovals int             // Number of required approvals for reversal
	ApprovalNodes     map[string]bool // List of nodes that approved the reversal
	RejectionNodes    map[string]bool // List of nodes that rejected the reversal
}

// RequestTransactionReversal allows users to request the reversal of a transaction.
func (trm *TransactionReversalManager) RequestTransactionReversal(txnID, reason, contactDetails string, documents []byte) (*ReversalRequest, error) {
	trm.mutex.Lock()
	defer trm.mutex.Unlock()

	// Retrieve the transaction from the ledger
	txn, err := trm.Ledger.GetTransactionByID(txnID)
	if err != nil || txn.ID == "" {
		return nil, errors.New("transaction not found or unauthorized request")
	}

	// Ensure the reversal is within the allowed time frame (after 24 hours and within 28 days)
	timeSinceTxn := time.Since(txn.Timestamp)
	if timeSinceTxn < 24*time.Hour || timeSinceTxn > trm.ReversalTimeLimit {
		return nil, errors.New("reversal request out of allowed time frame")
	}

	// Ensure the user's wallet is SYN900 verified
	verified, err := IsSyn900Verified(txn.From, trm.Ledger) // Corrected call to IsSyn900Verified
	if err != nil || !verified {
		return nil, errors.New("wallet is not SYN900 verified")
	}

	// Freeze the funds in the transaction
	err = trm.Ledger.FreezeTransactionAmount(txn.ID, txn.Amount) // Corrected FreezeTransactionAmount call
	if err != nil {
		return nil, err
	}

	// Create a reversal request
	reversalRequest := &ReversalRequest{
		ID:                GenerateRequestID(txnID), // Assuming a function GenerateRequestID exists
		TransactionID:     txnID,
		UserID:            txn.From, // Use txn.From instead of txn.SenderID
		Timestamp:         time.Now(),
		Status:            "Pending",
		Reason:            reason,
		ContactDetails:    contactDetails,
		SubmittedEvidence: documents,
		RequiredApprovals: getRequiredApprovals(txn.Amount), // Placeholder function, ensure to define
		ApprovalNodes:     make(map[string]bool),
		RejectionNodes:    make(map[string]bool),
	}

	// Broadcast the reversal request to authority nodes for verification
	err = trm.BroadcastReversalRequest(reversalRequest)
	if err != nil {
		return nil, err
	}

	// Log and return the reversal request
	trm.Logger.Println("Reversal request broadcasted successfully:", reversalRequest.ID) // Use Println instead of Log
	return reversalRequest, nil
}

// Helper function to generate a unique request ID (example implementation)
func GenerateRequestID(txnID string) string {
	return fmt.Sprintf("RR-%s-%d", txnID, time.Now().UnixNano()) // Generate unique ID with transaction ID and timestamp
}

// Helper function to determine required approvals based on amount (example implementation)
func getRequiredApprovals(amount float64) int {
	if amount > 10000 {
		return 3
	}
	return 1
}

// getRequiredApprovals calculates the number of confirmations or rejections required based on the SYNN amount
func getRequiredReversalApprovals(amount *big.Int) int {
	switch {
	case amount.Cmp(big.NewInt(25)) <= 0:
		return 3
	case amount.Cmp(big.NewInt(100)) <= 0:
		return 4
	case amount.Cmp(big.NewInt(500)) <= 0:
		return 5
	case amount.Cmp(big.NewInt(5000)) <= 0:
		return 6
	default:
		return 7
	}
}

func (trm *TransactionReversalManager) BroadcastReversalRequest(request *ReversalRequest) error {
	// Fetch authority nodes from the ledger as []*AuthorityNodeVersion
	authorityNodes, err := trm.Ledger.GetRandomAuthorityNodes(request.RequiredApprovals)
	if err != nil {
		return fmt.Errorf("failed to fetch authority nodes: %w", err)
	}

	if len(authorityNodes) < request.RequiredApprovals {
		return fmt.Errorf("not enough authority nodes available to meet required approvals")
	}

	// Send the reversal request to each selected authority node
	for _, node := range authorityNodes {
		// node is *AuthorityNodeVersion

		// Convert to common.AuthorityNodeVersion if necessary
		commonNodeVersion := convertToAuthorityNodeVersion(node)

		// Send the encrypted reversal request to the node
		err := trm.SendEncryptedMessage(commonNodeVersion, request)
		if err != nil {
			trm.Logger.Printf("Failed to send encrypted reversal request to node %s: %v", node.NodeID, err)
		} else {
			trm.Logger.Printf("Successfully sent encrypted reversal request to node %s", node.NodeID)
		}
	}

	return nil
}

// convertToAuthorityNodeVersion converts a *ledger.AuthorityNodeVersion to common.AuthorityNodeVersion
func convertToAuthorityNodeVersion(node *ledger.AuthorityNodeVersion) common.AuthorityNodeVersion {
	return common.AuthorityNodeVersion{
		NodeID:            node.NodeID,                                       // Copy NodeID directly
		SecretKey:         node.SecretKey,                                    // Copy SecretKey directly
		CreatedAt:         node.CreatedAt,                                    // Copy CreatedAt directly
		EncryptedKey:      node.EncryptedKey,                                 // Copy EncryptedKey directly
		AuthorityNodeType: common.AuthorityNodeTypes(node.AuthorityNodeType), // Explicit type conversion
	}
}

// convertToLedgerAuthorityNode converts a common.AuthorityNodeVersion to ledger.AuthorityNode
func convertToLedgerAuthorityNode(node *common.AuthorityNodeVersion) ledger.AuthorityNodeVersion {
	return ledger.AuthorityNodeVersion{
		NodeID:            node.NodeID,                                       // Copy NodeID directly
		SecretKey:         node.SecretKey,                                    // Copy SecretKey directly
		CreatedAt:         node.CreatedAt,                                    // Copy CreatedAt directly
		EncryptedKey:      node.EncryptedKey,                                 // Copy EncryptedKey directly
		AuthorityNodeType: ledger.AuthorityNodeTypes(node.AuthorityNodeType), // Explicit type conversion
	}
}

// SendEncryptedMessage sends an encrypted reversal request to a specific authority node
func (trm *TransactionReversalManager) SendEncryptedMessage(node common.AuthorityNodeVersion, request *ReversalRequest) error {
	// Serialize the request to JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w", err)
	}

	// Define an encryption key (in a real-world case, this key should be securely managed)
	encryptionKey := []byte("your-secure-encryption-key")

	// Encrypt the data (assumes the method expects: a type, a key, and the data to encrypt)
	encryptionType := "AES256" // Example encryption type
	encryptedMessage, err := trm.Encryption.EncryptData(encryptionType, encryptionKey, requestData)
	if err != nil {
		return fmt.Errorf("encryption failed for node %s: %w", node.NodeID, err)
	}

	// Simulate sending the encrypted message over the network
	fmt.Printf("Sending encrypted message to node %s: %x\n", node.NodeID, encryptedMessage)

	// Simulate possible failure in the network send
	if node.NodeID == "2" {
		return fmt.Errorf("failed to send message to node %s", node.NodeID) // Simulate a failure for node 2
	}

	return nil
}

// HandleReversalApproval handles approval or rejection of a reversal request by authority nodes.
func (trm *TransactionReversalManager) HandleReversalApproval(request *ReversalRequest, nodeID string, approval bool) error {
	trm.mutex.Lock()
	defer trm.mutex.Unlock()

	// Update approval or rejection nodes
	if approval {
		request.ApprovalNodes[nodeID] = true
	} else {
		request.RejectionNodes[nodeID] = true
	}

	// Check if the reversal is approved
	if len(request.ApprovalNodes) >= request.RequiredApprovals {
		request.Status = "Approved"
		err := trm.ExecuteTransactionReversal(request.TransactionID)
		if err != nil {
			return err
		}

		// Update the reversal request as approved in the ledger
		err = trm.Ledger.UpdateReversalRequest(request.TransactionID, true) // true means approved
		if err != nil {
			return err
		}

	} else if len(request.RejectionNodes) >= request.RequiredApprovals {
		// Check if the reversal is rejected
		request.Status = "Rejected"
		err := trm.ReleaseTransactionFunds(request.TransactionID)
		if err != nil {
			return err
		}

		// Update the reversal request as rejected in the ledger
		err = trm.Ledger.UpdateReversalRequest(request.TransactionID, false) // false means rejected
		if err != nil {
			return err
		}
	}

	// Log the event
	trm.Logger.Println("Reversal request processed:", request.ID)
	return nil
}

// ExecuteTransactionReversal reverts the transaction in the ledger and returns the funds to the sender.
func (trm *TransactionReversalManager) ExecuteTransactionReversal(transactionID string) error {
	// Retrieve the transaction from the ledger
	txn, err := trm.Ledger.GetTransactionByID(transactionID)
	if err != nil || txn.From == "" {
		return errors.New("transaction not found")
	}

	// Double validate that the transaction can be reversed
	if txn.Status != "Completed" {
		return errors.New("only completed transactions can be reversed")
	}

	// Reverse the transaction: deduct from receiver and add back to sender
	// Pass the pointer to the txn record
	err = trm.reverseFunds(&txn)
	if err != nil {
		return err
	}

	// Update the transaction state to 'Reversed' in the ledger
	err = trm.Ledger.UpdateTransactionStatus(transactionID, "Reversed") // Assuming UpdateTransactionStatus exists
	if err != nil {
		return err
	}

	// Notify the sender and receiver of the reversal
	err = trm.NotifyPartiesOfReversal(txn.ID, "Transaction has been reversed")
	if err != nil {
		return err
	}

	// Log the transaction reversal event
	trm.Logger.Println("Transaction reversal executed for:", txn.ID)

	return nil
}

// reverseFunds reverts the transfer by deducting from the receiver and returning funds to the sender.
func (trm *TransactionReversalManager) reverseFunds(txn *ledger.TransactionRecord) error {
	// Deduct funds from the receiver
	err := trm.Ledger.AdjustBalance(txn.To, -txn.Amount)
	if err != nil {
		return fmt.Errorf("failed to deduct funds from receiver: %v", err)
	}

	// Return funds to the sender
	err = trm.Ledger.AdjustBalance(txn.From, txn.Amount)
	if err != nil {
		return fmt.Errorf("failed to return funds to sender: %v", err)
	}

	return nil
}

// ReleaseTransactionFunds releases the frozen funds if the reversal is rejected.
func (trm *TransactionReversalManager) ReleaseTransactionFunds(transactionID string) error {
	// Retrieve the transaction from the ledger
	txn, err := trm.Ledger.GetTransactionByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Release the frozen transaction amount by passing the transaction ID
	err = trm.Ledger.ReleaseFrozenFunds(txn.ID)
	if err != nil {
		return fmt.Errorf("failed to release frozen funds: %v", err)
	}

	// Notify the sender and receiver of the rejection
	err = trm.NotifyPartiesOfReversal(txn.ID, "Reversal request rejected. Funds released.")
	if err != nil {
		return fmt.Errorf("failed to notify parties: %v", err)
	}

	// Log the event using log.Println
	trm.Logger.Println("Funds released for transaction:", txn.ID)
	return nil
}

// NotifyPartiesOfReversal sends a notification to the sender and receiver of the reversal status
func (trm *TransactionReversalManager) NotifyPartiesOfReversal(transactionID, message string) error {
	// Retrieve the transaction from the ledger
	txn, err := trm.Ledger.GetTransactionByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Notify the sender
	senderNotification := common.Notification{
		Type:      common.EmailNotification, // Adjust as needed (e.g., webhook)
		Recipient: txn.From,                 // Sender's contact info
		Message:   message,
	}
	err = common.SendNotification(senderNotification)
	if err != nil {
		return fmt.Errorf("failed to notify sender: %v", err)
	}

	// Notify the receiver
	receiverNotification := common.Notification{
		Type:      common.EmailNotification, // Adjust as needed
		Recipient: txn.To,                   // Receiver's contact info
		Message:   message,
	}
	err = common.SendNotification(receiverNotification)
	if err != nil {
		return fmt.Errorf("failed to notify receiver: %v", err)
	}

	// Log the notification event
	trm.Logger.Println("Notification sent for transaction:", txn.ID)
	return nil
}

// LogReversalEvent logs events related to the reversal process
func (trm *TransactionReversalManager) LogReversalEvent(transactionID, event string, data interface{}) error {
	trm.mutex.Lock()
	defer trm.mutex.Unlock()

	// Format the event details as a string
	eventDetails := fmt.Sprintf("Event: %s, TransactionID: %s, Data: %+v", event, transactionID, data)

	// Log the event using the standard logger
	trm.Logger.Println("Reversal Event:", eventDetails)

	// Record the formatted event in the ledger for immutability
	err := trm.Ledger.RecordEventInLedger(eventDetails) // Pass the concatenated eventDetails as a single string
	if err != nil {
		return fmt.Errorf("failed to log reversal event in ledger: %v", err)
	}

	return nil
}

// IsSyn900Verified checks if a wallet's SYN900 token is verified and active.
func IsSyn900Verified(walletID string, ledger *ledger.Ledger) (bool, error) {
	// Step 1: Retrieve the SYN900 token associated with the walletID from common package.
	token, err := ledger.GetTokenByWalletID(walletID) // Retrieves from common SYN900Tokens
	if err != nil {
		return false, fmt.Errorf("failed to retrieve SYN900 token for wallet %s: %v", walletID, err)
	}

	// Step 2: Check if the token exists.
	if token == nil {
		return false, fmt.Errorf("no SYN900 token found for wallet %s", walletID)
	}

	// Step 3: Ensure the token status is 'active'.
	if token.Status != "active" {
		return false, fmt.Errorf("SYN900 token for wallet %s is not active. Status: %s", walletID, token.Status)
	}

	// Step 4: Check if the token is expired.
	if time.Now().After(token.Metadata.ExpirationDate) {
		return false, fmt.Errorf("SYN900 token for wallet %s has expired", walletID)
	}

	// Step 5: Perform additional compliance checks on the token metadata (referencing common).
	if !token.Metadata.IsKYCCompliant || !token.Metadata.IsAMLCompliant {
		return false, fmt.Errorf("SYN900 token for wallet %s does not meet KYC/AML compliance", walletID)
	}

	// Step 6: Log the verification for auditing purposes in the ledger.
	err = ledger.RecordEventInLedger(fmt.Sprintf("SYN900 verification successful for wallet %s", walletID))
	if err != nil {
		return false, fmt.Errorf("failed to log verification activity: %v", err)
	}

	// Token is valid and verified.
	return true, nil
}
