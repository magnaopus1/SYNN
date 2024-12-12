package interoperability

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewOracleService initializes a new OracleService with a ledger instance
func NewOracleService(ledgerInstance *ledger.Ledger) *OracleService {
    return &OracleService{
        DataSources:    make(map[string]OracleDataSource),
        LedgerInstance: ledgerInstance,
    }
}

// AddDataSource adds a new data source for the oracle to pull data from.
func (oracle *OracleService) AddDataSource(sourceName, sourceURL string) error {
    oracle.mutex.Lock()
    defer oracle.mutex.Unlock()

    // Encrypt the source URL.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedURL, err := encryptInstance.EncryptData(sourceURL, common.EncryptionKey, make([]byte, 12)) // Example nonce
    if err != nil {
        return fmt.Errorf("failed to encrypt data source URL: %v", err)
    }

    // Add the encrypted URL to the data source.
    oracle.DataSources[sourceName] = OracleDataSource{
        SourceID:      sourceName,             // Using sourceName as ID here; adjust if needed
        Name:          sourceName,
        URL:           string(encryptedURL),   // Store encrypted URL as a string
        IsActive:      true,                   // Mark as active by default
        LastUpdated:   time.Now(),
        DataFormat:    "JSON",                 // Set format if known; modify as needed
    }

    fmt.Printf("Data source %s added successfully.\n", sourceName)
    return nil
}


// FetchData fetches data from a specific oracle data source.
func (oracle *OracleService) FetchData(sourceName string) (OracleData, error) {
    oracle.mutex.Lock()
    defer oracle.mutex.Unlock()

    // Get the data source by name
    dataSource, exists := oracle.DataSources[sourceName]
    if !exists {
        return OracleData{}, fmt.Errorf("data source %s not found", sourceName)
    }

    // Decrypt the URL stored in dataSource.URL
    decryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return OracleData{}, fmt.Errorf("failed to create decryption instance: %v", err)
    }

    decryptedURLBytes, err := decryptInstance.DecryptData([]byte(dataSource.URL), common.EncryptionKey)
    if err != nil {
        return OracleData{}, fmt.Errorf("failed to decrypt source URL: %v", err)
    }
    sourceURL := string(decryptedURLBytes) // Convert decrypted bytes to string

    // Simulate data fetching (replace this with actual fetch logic)
    fetchedData := fmt.Sprintf("Sample data fetched from %s", sourceURL)

    // Hash the fetched data
    hashedData := oracle.hashData(fetchedData)

    // Populate OracleData with fetched details
    oracleData := OracleData{
        SourceID:  dataSource.SourceID,
        Content:   fetchedData,
        FetchedAt: time.Now(),
        DataFormat: dataSource.DataFormat,
        Status:    "valid",
        Hash:      hashedData,
    }

    // Log fetched data to the ledger
    if err := oracle.logDataToLedger(oracleData); err != nil {
        return OracleData{}, fmt.Errorf("failed to log oracle data to ledger: %v", err)
    }

    fmt.Printf("Data fetched from %s and logged to the ledger.\n", sourceName)
    return oracleData, nil
}


// hashData generates a SHA-256 hash for the oracle data
func (oracle *OracleService) hashData(data string) string {
    hash := sha256.New()
    hash.Write([]byte(data))
    return hex.EncodeToString(hash.Sum(nil))
}

// logDataToLedger logs the fetched oracle data to the ledger.
func (oracle *OracleService) logDataToLedger(data OracleData) error {
    // Prepare the data log string for encryption
    dataLog := fmt.Sprintf("SourceID: %s | Content: %s | FetchedAt: %s | Status: %s",
        data.SourceID, data.Content, data.FetchedAt.String(), data.Status)

    // Create an encryption instance and encrypt the data log.
    encryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    encryptedLog, err := encryptInstance.EncryptData(dataLog, common.EncryptionKey, make([]byte, 12)) // Add a nonce if needed
    if err != nil {
        return fmt.Errorf("failed to encrypt oracle data log: %v", err)
    }

    // Convert encrypted data to base64 if required for storage as a string.
    encryptedLogStr := base64.StdEncoding.EncodeToString(encryptedLog)

    // Record the oracle data in the ledger.
    oracle.LedgerInstance.RecordOracleData(data.SourceID, data.SourceID, "chain_id_example", encryptedLogStr) // Use actual `chainID` as needed

    return nil
}




// RemoveDataSource removes a data source from the oracle
func (oracle *OracleService) RemoveDataSource(sourceName string) error {
    oracle.mutex.Lock()
    defer oracle.mutex.Unlock()

    if _, exists := oracle.DataSources[sourceName]; !exists {
        return fmt.Errorf("data source %s does not exist", sourceName)
    }

    delete(oracle.DataSources, sourceName)
    fmt.Printf("Data source %s removed.\n", sourceName)
    return nil
}
