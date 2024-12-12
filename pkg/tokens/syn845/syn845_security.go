package syn845

import (
	"errors"
	"sync"
	"time"

)

// SecureStorage defines the structure for managing secure storage of sensitive data
type SecureStorage struct {
	Key               []byte
	Nonce             []byte
	Storage           map[string]string
	Ledger            *ledger.Ledger               // Ledger for recording storage actions
	ConsensusEngine   *consensus.SynnergyConsensus // Consensus engine for validating secure storage actions
}

// NewSecureStorage initializes a new SecureStorage instance with ledger and consensus integration
func NewSecureStorage(password string, salt []byte, ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) (*SecureStorage, error) {
	key, err := generateKey(password, salt)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return &SecureStorage{
		Key:             key,
		Nonce:           nonce,
		Storage:         make(map[string]string),
		Ledger:          ledger,
		ConsensusEngine: consensusEngine,
	}, nil
}

// generateKey generates a key from the password and salt using scrypt
func generateKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
}

// Encrypt encrypts data using AES-GCM
func (ss *SecureStorage) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(ss.Key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, ss.Nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-GCM
func (ss *SecureStorage) Decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(ss.Key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, ss.Nonce, data, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Store securely stores the data, validates via consensus, and records the action in the ledger
func (ss *SecureStorage) Store(key, value string) error {
	// Encrypt the value before storing
	encryptedValue, err := ss.Encrypt(value)
	if err != nil {
		return err
	}

	// Validate the storage action via Synnergy Consensus
	if err := ss.ConsensusEngine.ValidateStorageAction(key, value); err != nil {
		return errors.New("storage action validation failed via Synnergy Consensus")
	}

	// Store the encrypted value
	ss.Storage[key] = encryptedValue

	// Record the storage action in the ledger
	if err := ss.Ledger.RecordStorageAction(key, encryptedValue); err != nil {
		return errors.New("failed to record storage action in the ledger")
	}

	return nil
}

// Retrieve securely retrieves the data, decrypts it, and validates via consensus
func (ss *SecureStorage) Retrieve(key string) (string, error) {
	// Retrieve the encrypted value from storage
	encryptedValue, exists := ss.Storage[key]
	if !exists {
		return "", errors.New("key does not exist")
	}

	// Decrypt the value
	decryptedValue, err := ss.Decrypt(encryptedValue)
	if err != nil {
		return "", err
	}

	// Validate the retrieval action via Synnergy Consensus
	if err := ss.ConsensusEngine.ValidateRetrieveAction(key, decryptedValue); err != nil {
		return "", errors.New("retrieve action validation failed via Synnergy Consensus")
	}

	return decryptedValue, nil
}

// Delete securely deletes the data, validates via consensus, and records the action in the ledger
func (ss *SecureStorage) Delete(key string) error {
	// Validate the deletion action via Synnergy Consensus
	if err := ss.ConsensusEngine.ValidateDeleteAction(key); err != nil {
		return errors.New("delete action validation failed via Synnergy Consensus")
	}

	// Delete the entry from storage
	delete(ss.Storage, key)

	// Record the deletion action in the ledger
	if err := ss.Ledger.RecordDeleteAction(key); err != nil {
		return errors.New("failed to record deletion action in the ledger")
	}

	return nil
}

// HashPassword securely hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}
