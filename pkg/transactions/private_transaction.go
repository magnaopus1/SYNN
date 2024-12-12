package transactions

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// PrivateTransaction defines the structure of a private transaction.
type PrivateTransaction struct {
    TransactionID   string            // Unique identifier for the transaction
    Sender          string            // Sender of the transaction
    Receiver        string            // Receiver of the transaction
    Amount          float64           // Transaction amount
    TokenType       string            // Token type (optional, defaults to "SYNN")
    TokenID         string            // Token ID (optional)
    IsPrivate       bool              // Flag indicating if this transaction is private
    EncryptedData   string            // Encrypted transaction details
    AuthorizedNodes map[string]string // Nodes allowed to view private details
    Fee             float64           // Fee for converting to private transaction
}


// PrivateTransactionManager manages the creation and conversion of private transactions.
type PrivateTransactionManager struct {
    mutex             sync.Mutex                    // For thread-safe operations
    Transactions      map[string]*PrivateTransaction // List of all private transactions
    Ledger            *ledger.Ledger                 // Ledger reference for transaction logging
    Consensus         *common.SynnergyConsensus      // Consensus engine for validation
    Encryption        *common.Encryption             // Encryption service for securing transaction details
    AuthorityNodeTypes []string                      // List of authority node types allowed to manage private transactions
    TransactionPool   *TransactionPool               // Pool for holding unconfirmed transactions
}

// NewPrivateTransactionManager initializes a new private transaction manager.
func NewPrivateTransactionManager(ledgerInstance *ledger.Ledger, consensus *common.SynnergyConsensus, encryptionService *common.Encryption, transactionPool *TransactionPool) *PrivateTransactionManager {
    return &PrivateTransactionManager{
        Transactions:      make(map[string]*PrivateTransaction),  // Initialize the transaction map
        Ledger:            ledgerInstance,                        // Ledger instance passed for logging
        Consensus:         consensus,                             // Consensus engine for validation
        Encryption:        encryptionService,                     // Encryption service for securing transaction details
        AuthorityNodeTypes: []string{"CentralBank", "Bank", "Regulator", "Government", "Creditor"}, // Define allowed authority node types
        TransactionPool:   transactionPool,                       // Initialize the transaction pool for unconfirmed transactions
    }
}

// CreatePrivateTransaction creates a new private transaction and adds it to the transaction pool for later inclusion in a sub-block.
func (ptm *PrivateTransactionManager) CreatePrivateTransaction(sender, receiver string, amount *big.Int, tokenType, tokenID string) (*PrivateTransaction, error) {
    ptm.mutex.Lock()
    defer ptm.mutex.Unlock()

    // Set default token type if not provided
    if tokenType == "" {
        tokenType = "SYNN"
    }

    // Generate a unique transaction ID
    transactionID := fmt.Sprintf("tx-%s-to-%s", sender, receiver)

    // Calculate the private transaction fee (1% of the total amount)
    fee := new(big.Int).Div(amount, big.NewInt(100)) // 1% fee

    // Encrypt transaction details
    encryptedData, err := ptm.Encryption.EncryptData("AES", []byte(fmt.Sprintf("%s:%s:%s:%s:%s:%s", sender, receiver, tokenType, amount.String(), tokenID, transactionID)), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt transaction data: %v", err)
    }

    // Convert amount and fee from *big.Int to float64
    amountFloat, _ := new(big.Float).SetInt(amount).Float64()
    feeFloat, _ := new(big.Float).SetInt(fee).Float64()

    // Create the private transaction
    privateTx := &PrivateTransaction{
        TransactionID:   transactionID,
        Sender:          sender,
        Receiver:        receiver,
        Amount:          amountFloat,
        TokenType:       tokenType,
        TokenID:         tokenID,
        IsPrivate:       true,
        EncryptedData:   string(encryptedData),
        AuthorizedNodes: ptm.getAuthorizedNodes(),
        Fee:             feeFloat,
    }

    // Convert to common.Transaction and add to transaction pool
    ptm.TransactionPool.AddTransaction(ptm.ConvertToCommonTransaction(privateTx))

    // Store the transaction in the private transaction map
    ptm.Transactions[transactionID] = privateTx

    // Log transaction creation
    fmt.Printf("Private transaction %s created successfully and added to the transaction pool by %s.\n", transactionID, sender)

    return privateTx, nil
}

// ConvertToStandardTransaction converts a private transaction back to a public/standard transaction.
func (ptm *PrivateTransactionManager) ConvertToStandardTransaction(transactionID string, sender string) error {
    ptm.mutex.Lock()
    defer ptm.mutex.Unlock()

    privateTx, exists := ptm.Transactions[transactionID]
    if !exists {
        return errors.New("transaction not found")
    }

    // Ensure only the sender can convert the transaction back to standard
    if privateTx.Sender != sender {
        return errors.New("only the sender can convert this transaction")
    }

    // Convert privateTx.Amount (float64) to *big.Int
    amountBigInt := new(big.Int)
    amountBigInt.SetString(fmt.Sprintf("%.0f", privateTx.Amount), 10)

    // Calculate the conversion fee (1% of the total amount)
    fee := new(big.Int).Div(amountBigInt, big.NewInt(100)) // 1% fee in *big.Int

    // Convert fee to float64 and add it to the transaction fee
    feeFloat, _ := new(big.Float).SetInt(fee).Float64()
    privateTx.Fee += feeFloat // Add the fee for conversion

    // Update transaction to standard
    privateTx.IsPrivate = false

    // Record the conversion in the ledger (fixing argument count)
    err := ptm.Ledger.RecordTransaction(privateTx.TransactionID, privateTx.Sender, privateTx.Amount)
    if err != nil {
        return fmt.Errorf("failed to record standard transaction in ledger: %v", err)
    }

    fmt.Printf("Private transaction %s converted to standard by %s.\n", transactionID, sender)
    return nil
}



// getAuthorizedNodes retrieves the list of authorized nodes that can view private transaction details.
func (ptm *PrivateTransactionManager) getAuthorizedNodes() map[string]string {
	authorizedNodes := make(map[string]string)

	// In a real system, this would select nodes of certain types (e.g., CentralBank, Regulator, etc.)
	for _, nodeType := range ptm.AuthorityNodeTypes {
		authorizedNodes[fmt.Sprintf("node-%s", nodeType)] = nodeType
	}

	return authorizedNodes
}

// GetTransactionDetails retrieves the details of a private transaction if the requester is authorized.
func (ptm *PrivateTransactionManager) GetTransactionDetails(transactionID, requester string) (*PrivateTransaction, error) {
	ptm.mutex.Lock()
	defer ptm.mutex.Unlock()

	privateTx, exists := ptm.Transactions[transactionID]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	// Check if the requester is the sender, receiver, or an authorized node
	if requester != privateTx.Sender && requester != privateTx.Receiver {
		if _, authorized := privateTx.AuthorizedNodes[requester]; !authorized {
			return nil, errors.New("access denied: requester is not authorized to view this private transaction")
		}
	}

	// Return the transaction details
	return privateTx, nil
}


// ConvertToCommonTransaction converts a PrivateTransaction to a common.Transaction
func (ptm *PrivateTransactionManager) ConvertToCommonTransaction(privateTx *PrivateTransaction) *common.Transaction {
    return &common.Transaction{
        TransactionID: privateTx.TransactionID,
        FromAddress:   privateTx.Sender,
        ToAddress:     privateTx.Receiver,
        Amount:        privateTx.Amount,
        TokenStandard: privateTx.TokenType,
        TokenID:       privateTx.TokenID,
        Status:        "Pending",
        EncryptedData: privateTx.EncryptedData,
        Fee:           privateTx.Fee,
    }
}