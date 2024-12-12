package syn2500

import (
	"time"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
)

// SYN2500Security handles the encryption, decryption, and security management for SYN2500 tokens.
type SYN2500Security struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewSYN2500Security initializes the security management system with RSA keys.
func NewSYN2500Security() (*SYN2500Security, error) {
	// Generate encryption keys for secure handling of data
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &SYN2500Security{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
}

// EncryptDAOData encrypts sensitive data related to DAO memberships, votes, and proposals.
func (sec *SYN2500Security) EncryptDAOData(data interface{}) ([]byte, error) {
	// Serialize data into JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a new AES encryption key
	aesKey := make([]byte, 32) // AES-256 encryption
	_, err = rand.Read(aesKey)
	if err != nil {
		return nil, err
	}

	// Encrypt the data using AES
	encryptedData, err := aesEncrypt(aesKey, dataBytes)
	if err != nil {
		return nil, err
	}

	// Encrypt the AES key using RSA
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, sec.publicKey, aesKey, nil)
	if err != nil {
		return nil, err
	}

	// Combine the encrypted AES key and the encrypted data
	result := append(encryptedAESKey, encryptedData...)

	return result, nil
}

// DecryptDAOData decrypts previously encrypted DAO data using the private key.
func (sec *SYN2500Security) DecryptDAOData(encryptedData []byte) (interface{}, error) {
	// Extract the encrypted AES key and encrypted data
	aesKeySize := sec.privateKey.PublicKey.Size()
	encryptedAESKey := encryptedData[:aesKeySize]
	encryptedDataPart := encryptedData[aesKeySize:]

	// Decrypt the AES key using RSA
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, sec.privateKey, encryptedAESKey, nil)
	if err != nil {
		return nil, err
	}

	// Decrypt the data using AES
	decryptedData, err := aesDecrypt(aesKey, encryptedDataPart)
	if err != nil {
		return nil, err
	}

	// Deserialize the data from JSON format
	var data interface{}
	err = json.Unmarshal(decryptedData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// VerifyDAOIntegrity verifies the integrity of a given DAO transaction using hashing and consensus.
func (sec *SYN2500Security) VerifyDAOIntegrity(data interface{}, expectedHash string) error {
	// Serialize the data to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Compute the SHA-512 hash of the data
	hash := sha512.Sum512(dataBytes)

	// Convert the computed hash to a hex string
	computedHash := hex.EncodeToString(hash[:])

	// Check if the computed hash matches the expected hash
	if computedHash != expectedHash {
		return errors.New("data integrity verification failed: hash mismatch")
	}

	return nil
}

// ValidateTransaction ensures the transaction is valid and adheres to the security protocol before adding it to the blockchain.
func (sec *SYN2500Security) ValidateTransaction(tx *common.DAOTokenTransaction) error {
	// Hash the transaction data to generate the unique identifier for the transaction
	txHash := sha256.Sum256([]byte(tx.TokenID + tx.NewOwner + tx.Timestamp.String()))
	tx.TransactionHash = hex.EncodeToString(txHash[:])

	// Use Synnergy Consensus to validate the transaction
	err := synconsensus.ValidateSubBlock(tx.TransactionHash)
	if err != nil {
		return errors.New("transaction validation failed through Synnergy Consensus")
	}

	// Store the transaction in the ledger
	err = ledger.StoreTransaction(tx, synconsensus.SubBlockValidation)
	if err != nil {
		return errors.New("transaction storage in ledger failed")
	}

	return nil
}

// Encrypt and Decrypt using AES for securing sensitive data
func aesEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func aesDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GenerateUniqueID creates a unique identifier for secure transactions or data blocks.
func GenerateUniqueID() string {
	timestamp := time.Now().UnixNano()
	hash := sha256.New()
	hash.Write([]byte(string(timestamp)))
	return hex.EncodeToString(hash.Sum(nil))
}
