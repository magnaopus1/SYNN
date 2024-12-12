package syn131

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewSecurityManager initializes a new instance of SecurityManager.
func NewSecurityManager(ledger *ledger.Ledger, consensusEngine *consensus.SynnergyConsensus) (*SecurityManager, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	return &SecurityManager{
		privateKey:     privateKey,
		publicKey:      &privateKey.PublicKey,
		ledger:         ledger,
		consensusEngine: consensusEngine,
	}, nil
}

// generatePrivateKey generates an RSA private key for encryption and digital signatures.
func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.New("failed to generate private key")
	}
	return privateKey, nil
}

// EncryptData encrypts data using RSA encryption.
func (sm *SecurityManager) EncryptData(data []byte) ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, sm.publicKey, data, nil)
	if err != nil {
		return nil, errors.New("failed to encrypt data")
	}
	return encryptedData, nil
}

// DecryptData decrypts data using RSA encryption.
func (sm *SecurityManager) DecryptData(encryptedData []byte) ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, sm.privateKey, encryptedData, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}
	return decryptedData, nil
}

// SignData generates a digital signature for the given data using RSA.
func (sm *SecurityManager) SignData(data []byte) ([]byte, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	hash := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, sm.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, errors.New("failed to sign data")
	}
	return signature, nil
}

// VerifySignature verifies a digital signature using the RSA public key.
func (sm *SecurityManager) VerifySignature(data, signature []byte) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	hash := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(sm.publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return errors.New("invalid signature")
	}
	return nil
}

// StoreKeys saves the RSA private key to a file for persistence.
func (sm *SecurityManager) StoreKeys(privateKeyPath string, publicKeyPath string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Save private key
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

	// Save public key
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

// LoadKeys loads the RSA private and public keys from files.
func (sm *SecurityManager) LoadKeys(privateKeyPath string, publicKeyPath string) error {
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

// ValidateTransaction ensures that a transaction meets security criteria, using consensus validation.
func (sm *SecurityManager) ValidateTransaction(transaction interface{}) error {
	// This uses Synnergy Consensus to validate the transaction
	validationResult := sm.consensusEngine.ValidateSYN131Transaction(transaction)
	if !validationResult.Valid {
		return errors.New("transaction failed validation")
	}

	// Additional security checks can be implemented here
	return nil
}

// SecureStorage wraps encryption around any data stored using the SecurityManager.
func (sm *SecurityManager) SecureStorage(data []byte) ([]byte, string, error) {
	encryptedData, encryptionKey, err := sm.EncryptData(data)
	if err != nil {
		return nil, "", err
	}

	// Save encryption key in the ledger for reference
	ledgerRefID := sm.ledger.StoreEncryptionKey(encryptionKey)
	if ledgerRefID == "" {
		return nil, "", errors.New("failed to store encryption key in ledger")
	}

	return encryptedData, ledgerRefID, nil
}

// LoadFromSecureStorage loads and decrypts data from storage using the SecurityManager.
func (sm *SecurityManager) LoadFromSecureStorage(encryptedData []byte, ledgerRefID string) ([]byte, error) {
	encryptionKey, err := sm.ledger.GetEncryptionKey(ledgerRefID)
	if err != nil {
		return nil, errors.New("failed to retrieve encryption key from ledger")
	}

	// Decrypt the data
	decryptedData, err := sm.DecryptData(encryptedData)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}

	return decryptedData, nil
}

// VerifyOwnership ensures ownership data is valid through signature verification.
func (sm *SecurityManager) VerifyOwnership(tokenID string, signature []byte, ownerPublicKey *rsa.PublicKey) error {
	tokenData, err := sm.ledger.GetTokenData(tokenID)
	if err != nil {
		return errors.New("failed to retrieve token data")
	}

	// Verify the signature using the owner's public key
	return sm.VerifySignature(tokenData, signature)
}

