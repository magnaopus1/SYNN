package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewOffchainWallet creates a new OffchainWallet with its own set of private and public keys.
func NewOffchainWallet(walletID, privateKey, publicKey string, ledgerInstance *ledger.Ledger) *OffchainWallet {
	return &OffchainWallet{
		WalletID:        walletID,
		PrivateKey:      privateKey,
		PublicKey:       publicKey,
		OffchainBalances: make(map[string]float64),
		ledgerInstance:  ledgerInstance,
	}
}

// CreateOffchainWallet generates a new off-chain wallet with a key pair and initializes it.
func CreateOffchainWallet(walletID string, ledgerInstance *ledger.Ledger) (*OffchainWallet, error) {
    // Generate private/public key pair (for simplicity, using a mock method)
    privateKey, publicKey := generateKeyPair()

    // Initialize the off-chain wallet
    wallet := NewOffchainWallet(walletID, privateKey, publicKey, ledgerInstance)

    // Convert keys to []byte for encryption
    privateKeyBytes := []byte(privateKey) // Assuming privateKey is a string
    publicKeyBytes := []byte(publicKey)   // Assuming publicKey is a string

    // Create the encryption instance inside the function
    encryptionInstance := &common.Encryption{}

    // Encrypt and store the wallet's keys in the ledger for safe recovery
    encryptedPrivateKey, err := encryptionInstance.EncryptData("AES", privateKeyBytes, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt private key: %v", err)
    }

    encryptedPublicKey, err := encryptionInstance.EncryptData("AES", publicKeyBytes, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt public key: %v", err)
    }

    // Convert encrypted keys to hex or base64 encoding to store them as strings
    encryptedPrivateKeyStr := hex.EncodeToString(encryptedPrivateKey)
    encryptedPublicKeyStr := hex.EncodeToString(encryptedPublicKey)

    // Store the encrypted private key and public key using the same StoreWalletKey method
    err = ledgerInstance.StoreWalletKey(walletID, "privateKey:"+encryptedPrivateKeyStr)
    if err != nil {
        return nil, fmt.Errorf("failed to store private key: %v", err)
    }

    err = ledgerInstance.StoreWalletKey(walletID, "publicKey:"+encryptedPublicKeyStr)
    if err != nil {
        return nil, fmt.Errorf("failed to store public key: %v", err)
    }

    fmt.Printf("Offchain wallet for walletID %s created successfully.\n", walletID)
    return wallet, nil
}





// GetBalance retrieves the off-chain balance for a specific currency or token.
func (ow *OffchainWallet) GetBalance(currency string) (float64, error) {
	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	balance, exists := ow.OffchainBalances[currency]
	if !exists {
		return 0, fmt.Errorf("no balance available for currency %s", currency)
	}
	return balance, nil
}

// AddFunds adds funds to the off-chain wallet's balance.
func (ow *OffchainWallet) AddFunds(currency string, amount float64) error {
	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	ow.OffchainBalances[currency] += amount
	fmt.Printf("Added %.2f %s to the offchain wallet %s\n", amount, currency, ow.WalletID)
	return nil
}

// TransferFunds transfers funds between off-chain wallets without touching the ledger.
func (ow *OffchainWallet) TransferFunds(receiver *OffchainWallet, currency string, amount float64) error {
	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	// Check if the sender has enough funds
	if ow.OffchainBalances[currency] < amount {
		return fmt.Errorf("insufficient funds to transfer %.2f %s from wallet %s", amount, currency, ow.WalletID)
	}

	// Deduct from the sender's balance
	ow.OffchainBalances[currency] -= amount

	// Add to the receiver's balance
	receiver.OffchainBalances[currency] += amount

	fmt.Printf("Transferred %.2f %s from wallet %s to wallet %s\n", amount, currency, ow.WalletID, receiver.WalletID)
	return nil
}

// CommitOffchainTransaction records the off-chain transaction to the ledger when necessary.
func (ow *OffchainWallet) CommitOffchainTransaction(currency string, amount float64, recipientWalletID string) error {
    ow.mutex.Lock()
    defer ow.mutex.Unlock()

    // Ensure enough balance is available before committing
    if ow.OffchainBalances[currency] < amount {
        return fmt.Errorf("insufficient funds to commit %.2f %s transaction for wallet %s", amount, currency, ow.WalletID)
    }

    // Deduct funds from off-chain balance
    ow.OffchainBalances[currency] -= amount

    // Create the transaction struct
    transaction := common.Transaction{
        TransactionID: fmt.Sprintf("tx-%d", time.Now().UnixNano()), // Unique Transaction ID
        FromAddress:   ow.WalletID,                                 // Sender's wallet
        ToAddress:     recipientWalletID,                           // Receiver's wallet
        Amount:        amount,                                      // Transaction amount
        Timestamp:     time.Now(),                                  // Current timestamp
        TokenStandard: currency,                                    // Assuming currency as the token standard
        EncryptedData: "",                                          // Encrypted transaction data
    }

    // Convert transaction to a string or byte representation (JSON serialization can be used here)
    transactionData, err := json.Marshal(transaction)
    if err != nil {
        return fmt.Errorf("failed to serialize transaction data: %v", err)
    }

    // Create the encryption instance
    encryptionInstance := &common.Encryption{}
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Define encryption key

    // Encrypt the transaction data
    encryptedTransactionData, err := encryptionInstance.EncryptData("AES", transactionData, encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt off-chain transaction: %v", err)
    }

    // Encode encrypted data to store as a string
    encryptedTransaction := base64.StdEncoding.EncodeToString(encryptedTransactionData)

    // Set encrypted data back to the transaction
    transaction.EncryptedData = encryptedTransaction

    // **Record the transaction to the ledger (only pass the required 3 arguments)**
    err = ow.ledgerInstance.RecordTransaction(transaction.FromAddress, transaction.ToAddress, transaction.Amount)
    if err != nil {
        return fmt.Errorf("failed to record off-chain transaction in the ledger: %v", err)
    }

    fmt.Printf("Committed off-chain transaction for wallet %s to the ledger.\n", ow.WalletID)
    return nil
}



// WithdrawToOnchain transfers funds from the off-chain wallet to an on-chain wallet.
func (ow *OffchainWallet) WithdrawToOnchain(onchainWalletID string, currency string, amount float64) error {
    ow.mutex.Lock()
    defer ow.mutex.Unlock()

    // Check if the off-chain wallet has enough balance
    if ow.OffchainBalances[currency] < amount {
        return fmt.Errorf("insufficient funds to withdraw %.2f %s from wallet %s", amount, currency, ow.WalletID)
    }

    // Deduct the amount from the off-chain wallet
    ow.OffchainBalances[currency] -= amount

    // Create the transaction struct
    transaction := common.Transaction{
        TransactionID:     fmt.Sprintf("tx-%d", time.Now().UnixNano()), // Unique Transaction ID
        FromAddress:       ow.WalletID,                                 // Off-chain wallet as the sender
        ToAddress:         onchainWalletID,                             // On-chain wallet as the receiver
        Amount:            amount,                                      // Transaction amount
        TokenStandard:     currency,                                    // Currency or token standard
        Timestamp:         time.Now(),                                  // Current timestamp
        EncryptedData:     "",                                          // Encrypted transaction data placeholder
    }

    // Convert the transaction to a byte representation (JSON serialization)
    transactionData, err := json.Marshal(transaction)
    if err != nil {
        return fmt.Errorf("failed to serialize transaction data: %v", err)
    }

    // Create the encryption instance
    encryptionInstance := &common.Encryption{}
    encryptionKey := []byte("your-32-byte-key-for-aes-encryption") // Define encryption key

    // Encrypt the transaction data
    encryptedTransactionData, err := encryptionInstance.EncryptData("AES", transactionData, encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt withdrawal transaction: %v", err)
    }

    // Encode the encrypted transaction data to a base64 string for storage
    encryptedTransaction := base64.StdEncoding.EncodeToString(encryptedTransactionData)

    // Set the encrypted data back to the transaction
    transaction.EncryptedData = encryptedTransaction

    // **Record the transaction in the ledger (only pass the required 3 arguments)**
    err = ow.ledgerInstance.RecordTransaction(transaction.FromAddress, transaction.ToAddress, transaction.Amount)
    if err != nil {
        return fmt.Errorf("failed to record withdrawal transaction in the ledger: %v", err)
    }

    fmt.Printf("Withdrew %.2f %s from off-chain wallet %s to on-chain wallet %s.\n", amount, currency, ow.WalletID, onchainWalletID)
    return nil
}



// generateKeyPair is a placeholder for a real-world key generation function (private/public).
func generateKeyPair() (string, string) {
	// This is a simple mock. In real-world applications, proper cryptographic methods must be used.
	privateKey := "mock_private_key"
	publicKey := "mock_public_key"
	return privateKey, publicKey
}
