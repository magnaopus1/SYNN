package syn722

import (
	"errors"
	"sync"
	"time"

)

// SYN722SecurityManager handles encryption, digital signatures, and security operations for SYN722 tokens
type SYN722SecurityManager struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	Ledger           *ledger.Ledger                // Ledger for recording security-related operations
	ConsensusEngine  *consensus.SynnergyConsensus  // Synnergy Consensus for validating security actions
	EncryptionService *encryption.EncryptionService // Encryption service for securing token data
	mutex            sync.Mutex                    // Mutex for secure concurrent access
}

// NewSYN722SecurityManager initializes a new SYN722SecurityManager
func NewSYN722SecurityManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus, encryptionService *encryption.EncryptionService) (*SYN722SecurityManager, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	return &SYN722SecurityManager{
		privateKey:      privateKey,
		publicKey:       &privateKey.PublicKey,
		Ledger:          ledger,
		ConsensusEngine: consensusEngine,
		EncryptionService: encryptionService,
	}, nil
}

// EncryptData encrypts the given data using RSA encryption
func (sm *SYN722SecurityManager) EncryptData(data []byte) (string, string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, sm.publicKey, data, nil)
	if err != nil {
		return "", "", errors.New("failed to encrypt data")
	}

	encryptedString := string(encryptedData)

	// Store the encryption key in the ledger for auditing and recovery purposes
	encryptionKey := sm.storeEncryptionKeyInLedger(encryptedString)

	return encryptedString, encryptionKey, nil
}

// DecryptData decrypts the encrypted data using the stored private key
func (sm *SYN722SecurityManager) DecryptData(encryptedData string, encryptionKey string) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Retrieve the encryption key if necessary (e.g., auditing or after a secure recovery)
	if encryptionKey != "" {
		encryptionKey, err := sm.Ledger.GetEncryptionKey(encryptionKey)
		if err != nil {
			return "", errors.New("failed to retrieve encryption key from ledger")
		}
		encryptedData = encryptionKey
	}

	// Decrypt the data using the private key
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, sm.privateKey, []byte(encryptedData), nil)
	if err != nil {
		return "", errors.New("failed to decrypt data")
	}

	return string(decryptedData), nil
}

// SignData generates a digital signature for the provided data using the private key
func (sm *SYN722SecurityManager) SignData(data []byte) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, sm.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", errors.New("failed to sign data")
	}

	return string(signature), nil
}

// VerifySignature verifies the provided digital signature using the public key
func (sm *SYN722SecurityManager) VerifySignature(data []byte, signature string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	hash := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(sm.publicKey, crypto.SHA256, hash[:], []byte(signature))
	if err != nil {
		return errors.New("signature verification failed")
	}

	return nil
}

// storeEncryptionKeyInLedger securely stores the encryption key in the ledger for future use
func (sm *SYN722SecurityManager) storeEncryptionKeyInLedger(encryptionKey string) string {
	// Store the key in the ledger and return its reference ID for future retrieval
	keyID, err := sm.Ledger.StoreEncryptionKey(encryptionKey)
	if err != nil {
		return ""
	}
	return keyID
}

// ValidateSecurityEvent uses Synnergy Consensus to validate security-related actions (e.g., encryption, decryption)
func (sm *SYN722SecurityManager) ValidateSecurityEvent(eventType, tokenID string, details map[string]interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate the event through Synnergy Consensus
	if err := sm.ConsensusEngine.ValidateSecurityEvent(eventType, tokenID, details); err != nil {
		return errors.New("security event validation failed via Synnergy Consensus")
	}

	return nil
}

// generatePrivateKey generates a new RSA private key for encryption and digital signatures
func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.New("failed to generate private key")
	}
	return privateKey, nil
}

// StoreKeys persists the RSA private and public keys to a secure location for future retrieval
func (sm *SYN722SecurityManager) StoreKeys(privateKeyPath, publicKeyPath string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Store private key
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return errors.New("failed to create private key file")
	}
	defer privateKeyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(sm.privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(privateKeyFile, privateBlock)
	if err != nil {
		return errors.New("failed to encode private key")
	}

	// Store public key
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return errors.New("failed to create public key file")
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(sm.publicKey)
	if err != nil {
		return errors.New("failed to marshal public key")
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	err = pem.Encode(publicKeyFile, publicBlock)
	if err != nil {
		return errors.New("failed to encode public key")
	}

	return nil
}

// LoadKeys loads the RSA private and public keys from the specified files
func (sm *SYN722SecurityManager) LoadKeys(privateKeyPath, publicKeyPath string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Load private key
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return errors.New("failed to read private key file")
	}
	privateBlock, _ := pem.Decode(privateKeyFile)
	if privateBlock == nil {
		return errors.New("failed to decode private key PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return errors.New("failed to parse private key")
	}
	sm.privateKey = privateKey

	// Load public key
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return errors.New("failed to read public key file")
	}
	publicBlock, _ := pem.Decode(publicKeyFile)
	if publicBlock == nil {
		return errors.New("failed to decode public key PEM block")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return errors.New("failed to parse public key")
	}
	sm.publicKey = publicKey.(*rsa.PublicKey)

	return nil
}

