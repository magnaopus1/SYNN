package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"synnergy_network/pkg/ledger"
	"time"
)

// Wallet represents a user's wallet containing private and public keys.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey // The private key used to sign transactions.
	PublicKey  []byte            // The public key used to receive funds.
	Address    string            // Wallet address derived from the public key.
	Ledger     *ledger.Ledger    // Ledger instance for interaction.
}

// WalletData represents data stored for a user's wallet.
type WalletData struct {
	WalletID       string   // Unique identifier for the wallet
	OwnerAddress   string   // Address of the wallet owner
	Balance        *big.Int // Balance of the wallet in the main currency
	TokenBalances  map[string]*big.Int // Token balances
	TransactionHistory []TransactionRecord // Transaction history associated with the wallet
	IsVerified     bool     // Whether the wallet is verified
	CreatedAt      time.Time // Time the wallet was created
	TokenID   string  // Token ID associated with the wallet (Add this field)
    Verified  bool    // Whether the wallet has been verified (Add this field)

}

// NewWallet creates a new wallet with a unique key pair and address.
func NewWallet(ledgerInstance *ledger.Ledger) (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		Ledger:     ledgerInstance,
	}, nil
}

// generateAddress creates a wallet address by hashing the public key.
func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[:])
}

// EncryptPrivateKey encrypts the wallet's private key for secure storage.
func (w *Wallet) EncryptPrivateKey(passphrase string) ([]byte, error) {
	// Instantiate an Encryption instance
	encryptionInstance := &Encryption{}

	// Convert the private key to bytes
	privateKeyBytes := w.PrivateKey.D.Bytes()

	// Encrypt the private key using the passphrase and the AES algorithm
	return encryptionInstance.EncryptData("AES", privateKeyBytes, []byte(passphrase))
}

// DecryptPrivateKey decrypts the wallet's private key.
func (w *Wallet) DecryptPrivateKey(encryptedData []byte, passphrase string) error {
	// Instantiate an Encryption instance
	encryptionInstance := &Encryption{}

	// Decrypt the encrypted data using the passphrase
	decryptedData, err := encryptionInstance.DecryptData(encryptedData, []byte(passphrase))
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %v", err)
	}

	// Rebuild the private key from the decrypted data
	w.PrivateKey = new(ecdsa.PrivateKey)
	w.PrivateKey.PublicKey.Curve = elliptic.P256()
	w.PrivateKey.D = new(big.Int).SetBytes(decryptedData)
	w.PublicKey = append(w.PrivateKey.PublicKey.X.Bytes(), w.PrivateKey.PublicKey.Y.Bytes()...)

	return nil
}



// CreateTransaction creates a new transaction from this wallet.
func (w *Wallet) CreateTransaction(to string, amount float64, feeManager *NetworkFeeManager, gasUnits int, gasPricePerUnit float64, userTip float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("invalid transaction amount")
	}

	// Ensure wallet has sufficient balance
	balance, err := w.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	if balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Get the next available nonce for this wallet (if needed for later use, handle separately)
	_, err = w.Ledger.AccountsWalletLedger.GetNextNonce(w.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get next nonce: %v", err)
	}

	// Calculate the transaction fee using the NetworkFeeManager
	transactionFee, err := feeManager.ProcessTransactionFee(w.Address, gasUnits, gasPricePerUnit, userTip)
	if err != nil {
		return nil, fmt.Errorf("failed to process transaction fee: %v", err)
	}

	// Create the transaction
	tx := &Transaction{
		FromAddress:  w.Address,           // Sender's address
		ToAddress:    to,                  // Recipient's address
		Amount:       amount,              // Transaction amount
		Fee:          transactionFee.TotalFee, // Total transaction fee
		Timestamp:    time.Now(),          // Current timestamp
		Signature:    "",                  // Signature will be added after signing
		Status:       "pending",           // Initial status
		SubBlockID:   "",                  // Will be populated later
		BlockID:      "",                  // Will be populated later
		ValidatorID:  "",                  // Will be populated later
		TokenStandard: "SYNN",            // Default Token Standard 
		TokenID:      "",                  // Token ID (optional, if applicable)
	}

	// Sign the transaction
	err = w.SignTransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return tx, nil
}




// SignTransaction signs the transaction using the wallet's private key.
func (w *Wallet) SignTransaction(tx *Transaction) error {
	// Convert transaction to a string representation for hashing (this method should be implemented)
	txData := tx.StringForSigning() // Use the helper function to get the string representation
	txHash := sha256.Sum256([]byte(txData)) // Hash the transaction data

	// Sign the hashed transaction using the wallet's private key
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, txHash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Combine r and s into a single byte slice and convert to a hexadecimal string for the signature
	signatureBytes := append(r.Bytes(), s.Bytes()...)
	tx.Signature = hex.EncodeToString(signatureBytes) // Convert the byte slice to a hexadecimal string

	return nil
}

// SendTransaction sends the signed transaction to the ledger for validation.
func (w *Wallet) SendTransaction(tx *Transaction) error {
	// Validate transaction via Synnergy Consensus
	err := w.Ledger.BlockchainConsensusCoinLedger.ValidateTransaction(tx.TransactionID) // TransactionID is the correct field
	if err != nil {
		return fmt.Errorf("failed to validate transaction: %v", err)
	}

	// Log the transaction to the ledger
	err = w.Ledger.BlockchainConsensusCoinLedger.LogTransaction(tx.TransactionID, tx.Signature) // Pass TransactionID and Signature
	if err != nil {
		return fmt.Errorf("failed to log transaction: %v", err)
	}

	fmt.Printf("Transaction %s sent successfully.\n", tx.TransactionID)
	return nil
}


// GetBalance retrieves the wallet's current balance by querying the ledger.
func (w *Wallet) GetBalance() (float64, error) {
	balance, err := w.Ledger.GetBalance(w.Address)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve balance: %v", err)
	}
	return balance, nil
}

// GetTransactionHistory retrieves the wallet's transaction history from the ledger.
func (w *Wallet) GetTransactionHistory() ([]*Transaction, error) {
    records, err := w.Ledger.BlockchainConsensusCoinLedger.GetTransactionHistory()
    if err != nil {
        return nil, fmt.Errorf("failed to get transaction history: %v", err)
    }

    // Convert TransactionRecord to Transaction
    var transactions []*Transaction
    for _, record := range records {
        transaction := &Transaction{
            TransactionID: record.Hash,  // Assuming the hash is used as the Transaction ID
            FromAddress:   record.From,
            ToAddress:     record.To,
            Amount:        record.Amount,
            Fee:           record.Fee,
            Timestamp:     record.Timestamp,
            // Map other fields as necessary
        }
        transactions = append(transactions, transaction)
    }
    
    return transactions, nil
}



