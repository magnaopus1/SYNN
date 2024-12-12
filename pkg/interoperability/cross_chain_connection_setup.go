package interoperability

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)


// NewCrossChainSetup initializes a new CrossChainSetup instance
func NewCrossChainSetup(ledgerInstance *ledger.Ledger) *CrossChainSetup {
    return &CrossChainSetup{
        Connections:    make(map[string]string),
        LedgerInstance: ledgerInstance,
    }
}

// SetupConnection allows the setup of a connection between this blockchain and another by specifying the blockchain's URL.
func (ccs *CrossChainSetup) SetupConnection(chainName, connectionURL string) error {
    ccs.mutex.Lock()
    defer ccs.mutex.Unlock()

    // Create an encryption instance and encrypt the connection URL.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedURL, err := encryptInstance.EncryptData(connectionURL, common.EncryptionKey, nil)
    if err != nil {
        return fmt.Errorf("failed to encrypt connection URL: %v", err)
    }

    // Record the encrypted connection URL in the map.
    ccs.Connections[chainName] = string(encryptedURL)

    // Log the connection to the ledger.
    err = ccs.logConnectionToLedger(chainName, string(encryptedURL))
    if err != nil {
        return fmt.Errorf("failed to log connection to ledger: %v", err)
    }

    fmt.Printf("Connection to blockchain %s set up with URL %s.\n", chainName, connectionURL)
    return nil
}
// GetConnection retrieves the connection URL for a specific blockchain.
func (ccs *CrossChainSetup) GetConnection(chainName string) (string, error) {
    ccs.mutex.Lock()
    defer ccs.mutex.Unlock()

    encryptedURL, exists := ccs.Connections[chainName]
    if !exists {
        return "", fmt.Errorf("connection to blockchain %s not found", chainName)
    }

    // Create an encryption instance to decrypt the URL.
    decryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return "", fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Decrypt the URL for use, passing only the required arguments.
    decryptedURL, err := decryptInstance.DecryptData([]byte(encryptedURL), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt connection URL: %v", err)
    }

    return string(decryptedURL), nil
}

// RemoveConnection removes a connection to a blockchain
func (ccs *CrossChainSetup) RemoveConnection(chainName string) error {
    ccs.mutex.Lock()
    defer ccs.mutex.Unlock()

    if _, exists := ccs.Connections[chainName]; !exists {
        return fmt.Errorf("connection to blockchain %s does not exist", chainName)
    }

    delete(ccs.Connections, chainName)
    fmt.Printf("Connection to blockchain %s removed.\n", chainName)
    return nil
}

// logConnectionToLedger logs the connection details to the ledger for immutability.
func (ccs *CrossChainSetup) logConnectionToLedger(chainName, encryptedURL string) error {
    connectionLog := fmt.Sprintf("Blockchain: %s | EncryptedURL: %s", chainName, encryptedURL)

    // Create an encryption instance to encrypt the connection log.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the log with the required arguments (message, encryption key, and additional nonce or IV).
    nonce := make([]byte, 12) // Example nonce (adjust as necessary)
    encryptedLog, err := encryptInstance.EncryptData(connectionLog, common.EncryptionKey, nonce)
    if err != nil {
        return fmt.Errorf("failed to encrypt connection log: %v", err)
    }

    // Define protocolID and setupDetails, and use the expected arguments for RecordCrossChainSetup.
    protocolID := "protocol_id_example" // Replace with actual protocol ID
    setupDetails := string(encryptedLog) // Convert encryptedLog to string if required as setup details

    // Call RecordCrossChainSetup and handle any errors without returning its result, as it has no return type.
    ccs.LedgerInstance.RecordCrossChainSetup(protocolID, chainName, "TargetChain", setupDetails)

    return nil
}

// generateConnectionHash generates a hash for the connection based on the blockchain name and timestamp
func (ccs *CrossChainSetup) generateConnectionHash(chainName string) string {
    hashInput := fmt.Sprintf("%s%d", chainName, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
