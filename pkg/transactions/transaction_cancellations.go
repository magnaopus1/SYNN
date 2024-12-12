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

// CancellationRequest represents a request to cancel a transaction.
type CancellationRequest struct {
	ID                string            // Unique request ID
	TransactionID     string            // ID of the transaction to cancel
	UserID            string            // ID of the user requesting cancellation
	Timestamp         time.Time         // Timestamp of the cancellation request
	Status            string            // Current status (e.g., "Pending", "Approved", "Rejected")
	Reason            string            // Reason for the cancellation request
	DocumentEvidence  string            // Any evidence provided for the cancellation
	ContactDetails    string            // Contact details of the user requesting cancellation
	RequiredApprovals int               // Number of required approvals for cancellation
	ApprovalNodes     map[string]bool   // Nodes that have approved
	RejectionNodes    map[string]bool   // Nodes that have rejected
}

// RequestTransactionCancellation allows users to request the cancellation of a transaction.
func (tcm *TransactionCancellationManager) RequestTransactionCancellation(txnID, reason, contactDetails, documentEvidence string) (*CancellationRequest, error) {
    tcm.mutex.Lock()
    defer tcm.mutex.Unlock()

    // Retrieve the transaction from the ledger
    txn, err := tcm.Ledger.GetTransactionByID(txnID)
    if err != nil || txn.ID == "" {
        return nil, errors.New("transaction not found or unauthorized request")
    }

    // Ensure the transaction is within the allowed cancellation time frame (after 24 hours and within the timeout period)
    timeSinceTxn := time.Since(txn.Timestamp)
    if timeSinceTxn < 24*time.Hour || timeSinceTxn > tcm.TimeoutPeriod {
        return nil, errors.New("cancellation request out of allowed time frame")
    }

    // Ensure the user's wallet is SYN900 verified
    verified, err := IsSyn900Verified(txn.From, tcm.Ledger)
    if err != nil || !verified {
        return nil, errors.New("wallet is not SYN900 verified")
    }

    // Freeze the transaction funds
    err = tcm.Ledger.FreezeTransactionAmount(txn.ID, txn.Amount)
    if err != nil {
        return nil, err
    }

    // Create a cancellation request
    cancelRequest := &CancellationRequest{
        ID:                GenerateRequestID(txnID), // Assuming a function GenerateRequestID exists
        TransactionID:     txnID,
        UserID:            txn.From,                 // Use txn.From instead of txn.SenderID
        Timestamp:         time.Now(),
        Status:            "Pending",
        Reason:            reason,
        ContactDetails:    contactDetails,
        DocumentEvidence:  documentEvidence,
        RequiredApprovals: getRequiredApprovals(txn.Amount), // Placeholder function, ensure to define
        ApprovalNodes:     make(map[string]bool),
        RejectionNodes:    make(map[string]bool),
    }

    // Broadcast the cancellation request to authority nodes for verification
    err = tcm.BroadcastTransactionCancellationRequest(cancelRequest)
    if err != nil {
        return nil, err
    }

    // Log and return the cancellation request
    tcm.Logger.Println("Cancellation request broadcasted successfully:", cancelRequest.ID)
    return cancelRequest, nil
}


// getRequiredApprovals determines the number of confirmations or rejections required based on the amount of SYNN in the transaction
func getRequiredCancellationApprovals(amount *big.Int) int {
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

// HandleTransactionCancellationApproval processes approval or rejection of a cancellation request by authority nodes
func (tcm *TransactionCancellationManager) HandleTransactionCancellationApproval(request *CancellationRequest, nodeID string, approval bool) error {
    tcm.mutex.Lock()
    defer tcm.mutex.Unlock()

    if approval {
        request.ApprovalNodes[nodeID] = true
    } else {
        request.RejectionNodes[nodeID] = true
    }

    // Check if the cancellation is approved or rejected
    if len(request.ApprovalNodes) >= request.RequiredApprovals {
        request.Status = "Approved"
        err := tcm.ExecuteTransactionCancellation(request.TransactionID)
        if err != nil {
            return err
        }
    } else if len(request.RejectionNodes) >= request.RequiredApprovals {
        request.Status = "Rejected"
        err := tcm.ReleaseTransactionCancellationFunds(request.TransactionID)
        if err != nil {
            return err
        }
    }

    // Update request status in the ledger
    approved := request.Status == "Approved"
    err := tcm.Ledger.UpdateCancellationRequest(request.TransactionID, approved) // Pass TransactionID and bool
    if err != nil {
        return err
    }

    // Use Println for logging (Logger.Log does not exist in log.Logger)
    tcm.Logger.Println("Cancellation request processed:", request.ID)

    return nil
}


// BroadcastTransactionCancellationRequest sends the cancellation request to authority nodes and handles timeout responses
func (tcm *TransactionCancellationManager) BroadcastTransactionCancellationRequest(request *CancellationRequest) error {
    // Fetch authority nodes from the ledger as []*AuthorityNodeVersion
    authorityNodes, err := tcm.Ledger.GetRandomAuthorityNodes(request.RequiredApprovals)
    if err != nil {
        return fmt.Errorf("failed to fetch authority nodes: %w", err)
    }

    // Ensure there are enough authority nodes
    if len(authorityNodes) < request.RequiredApprovals {
        return fmt.Errorf("not enough authority nodes available to meet required approvals")
    }

    // Send the cancellation request to each selected authority node
    for _, node := range authorityNodes {
        // Convert to common.AuthorityNodeVersion if necessary
        commonNodeVersion := convertToAuthorityNodeVersion(node)

        // Send the encrypted cancellation request to the node with a timeout
        err := tcm.sendCancellationRequestWithTimeout(commonNodeVersion, request)
        if err != nil {
            tcm.Logger.Printf("Failed to send encrypted cancellation request to node %s: %v", node.NodeID, err)
        } else {
            tcm.Logger.Printf("Successfully sent encrypted cancellation request to node %s", node.NodeID)
        }
    }

    return nil
}

// sendCancellationRequestWithTimeout sends a cancellation request to a node and retries if no response is received within the timeout period
func (tcm *TransactionCancellationManager) sendCancellationRequestWithTimeout(node common.AuthorityNodeVersion, request *CancellationRequest) error {
    successChan := make(chan bool)

    go func() {
        err := tcm.SendEncryptedMessage(node, request)
        if err != nil {
            tcm.Logger.Printf("Error sending request to node %s: %v", node.NodeID, err)
            successChan <- false
            return
        }

        // Simulate node response delay and check for timeout
        select {
        case <-time.After(tcm.ResponseTimeout): // No response within timeout period
            successChan <- false
        case success := <-successChan:
            successChan <- success
        }
    }()

    success := <-successChan
    if !success {
        // Re-fetch a new set of random authority nodes
        tcm.Logger.Println("No response from node, resending request to another random node.")
        authorityNodes, err := tcm.Ledger.GetRandomAuthorityNodes(1) // Fetch a single random authority node
        if err != nil || len(authorityNodes) == 0 {
            return fmt.Errorf("failed to fetch a new authority node: %v", err)
        }

        newNode := authorityNodes[0] // Get the new random node

        // Convert to common.AuthorityNodeVersion if necessary
        commonNodeVersion := convertToAuthorityNodeVersion(newNode)

        return tcm.sendCancellationRequestWithTimeout(commonNodeVersion, request)
    }

    return nil
}



// SendEncryptedMessage sends an encrypted cancellation request to a specific authority node
func (tcm *TransactionCancellationManager) SendEncryptedMessage(node common.AuthorityNodeVersion, request *CancellationRequest) error {
    // Serialize the request to JSON
    requestData, err := json.Marshal(request)
    if err != nil {
        return fmt.Errorf("failed to serialize request: %w", err)
    }

    // Define an encryption key (in a real-world case, this key should be securely managed)
    encryptionKey := []byte("your-secure-encryption-key")

    // Encrypt the data (assumes the method expects: a type, a key, and the data to encrypt)
    encryptionType := "AES256" // Example encryption type
    encryptedMessage, err := tcm.Encryption.EncryptData(encryptionType, encryptionKey, requestData)
    if err != nil {
        return fmt.Errorf("encryption failed for node %s: %w", node.NodeID, err)
    }

    // Simulate sending the encrypted message over the network
    fmt.Printf("Sending encrypted cancellation message to node %s: %x\n", node.NodeID, encryptedMessage)

    // Simulate possible failure in the network send
    if node.NodeID == "2" {
        return fmt.Errorf("failed to send message to node %s", node.NodeID) // Simulate a failure for node 2
    }

    return nil
}


// ExecuteTransactionCancellation cancels the transaction by reverting it on the ledger
func (tcm *TransactionCancellationManager) ExecuteTransactionCancellation(transactionID string) error {
    // Retrieve the transaction from the ledger
    txn, err := tcm.Ledger.GetTransactionByID(transactionID)
    if err != nil || txn.From == "" {
        return errors.New("transaction not found")
    }

    // Double validate that the transaction can be cancelled
    if txn.Status != "Completed" {
        return errors.New("only completed transactions can be cancelled")
    }

    // Reverse the transaction: deduct from receiver and add back to sender
    err = tcm.reverseFunds(&txn)
    if err != nil {
        return err
    }

    // Update the transaction state to 'Cancelled' in the ledger
    err = tcm.Ledger.UpdateTransactionStatus(transactionID, "Cancelled") // Assuming UpdateTransactionStatus exists
    if err != nil {
        return err
    }

    // Notify the sender and receiver of the cancellation
    err = tcm.NotifyPartiesOfCancellation(txn.ID, "Transaction has been cancelled")
    if err != nil {
        return err
    }

    // Log the transaction cancellation event
    tcm.Logger.Println("Transaction cancellation executed for:", txn.ID)

    return nil
}

// reverseFunds reverts the transaction by deducting from the receiver and returning funds to the sender.
func (tcm *TransactionCancellationManager) reverseFunds(txn *ledger.TransactionRecord) error {
    // Deduct funds from the receiver
    err := tcm.Ledger.AdjustBalance(txn.To, -txn.Amount)
    if err != nil {
        return fmt.Errorf("failed to deduct funds from receiver: %v", err)
    }

    // Return funds to the sender
    err = tcm.Ledger.AdjustBalance(txn.From, txn.Amount)
    if err != nil {
        return fmt.Errorf("failed to return funds to sender: %v", err)
    }

    return nil
}



// ReleaseTransactionCancellationFunds releases the transaction amount back to the user if cancellation is rejected
func (tcm *TransactionCancellationManager) ReleaseTransactionCancellationFunds(transactionID string) error {
    // Retrieve the transaction from the ledger
    txn, err := tcm.Ledger.GetTransactionByID(transactionID)
    if err != nil {
        return errors.New("transaction not found")
    }

    // Release the frozen transaction amount by passing the transaction ID
    err = tcm.Ledger.ReleaseFrozenFunds(txn.ID)
    if err != nil {
        return fmt.Errorf("failed to release frozen funds: %v", err)
    }

    // Notify the sender and receiver of the rejection
    err = tcm.NotifyPartiesOfCancellation(txn.ID, "Cancellation request rejected. Funds released.")
    if err != nil {
        return fmt.Errorf("failed to notify parties: %v", err)
    }

    // Log the event using log.Println
    tcm.Logger.Println("Funds released for transaction:", txn.ID)
    return nil
}


// NotifyPartiesOfCancellation sends notifications to the sender and receiver about the cancellation status
func (tcm *TransactionCancellationManager) NotifyPartiesOfCancellation(transactionID, message string) error {
    // Retrieve the transaction from the ledger
    txn, err := tcm.Ledger.GetTransactionByID(transactionID)
    if err != nil {
        return errors.New("transaction not found")
    }

    // Notify the sender
    senderNotification := common.Notification{
        Type:      common.EmailNotification,  // Adjust as needed (e.g., webhook)
        Recipient: txn.From,                  // Sender's contact info (e.g., email)
        Message:   message,
    }
    err = common.SendNotification(senderNotification)
    if err != nil {
        return fmt.Errorf("failed to notify sender: %v", err)
    }

    // Notify the receiver
    receiverNotification := common.Notification{
        Type:      common.EmailNotification,  // Adjust as needed (e.g., webhook)
        Recipient: txn.To,                    // Receiver's contact info (e.g., email)
        Message:   message,
    }
    err = common.SendNotification(receiverNotification)
    if err != nil {
        return fmt.Errorf("failed to notify receiver: %v", err)
    }

    // Log the successful notification event
    tcm.Logger.Println("Notification sent for transaction:", txn.ID)
    return nil
}

