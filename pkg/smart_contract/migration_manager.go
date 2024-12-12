package smart_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewMigrationManager initializes the MigrationManager.
func NewMigrationManager(ledgerInstance *ledger.Ledger) *MigrationManager {
    return &MigrationManager{
        Contracts:      make(map[string]*MigratedContract),
        LedgerInstance: ledgerInstance,
    }
}

// MigrateContract initiates the migration of a smart contract to a new version.
func (mm *MigrationManager) MigrateContract(oldContractID, newCode, owner string, newParameters map[string]interface{}) (*MigratedContract, error) {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()

    // Retrieve the old contract
    oldContract, err := mm.LedgerInstance.RetrieveContract(oldContractID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve old contract: %v", err)
    }

    // Check if the owner matches the old contract's owner
    if oldContract.Deployer != owner {
        return nil, errors.New("unauthorized: only the contract owner can migrate the contract")
    }

    // Generate new contract ID for the migrated version
    newContractID := generateNewContractID(oldContractID, newCode)

    // Create the migrated contract structure
    newContract := &MigratedContract{
        OldContractID: oldContractID,
        NewContractID: newContractID,
        MigrationTime: time.Now(),
        Owner:         owner,
        NewCode:       newCode,
        NewParameters: newParameters,
    }

    // Encrypt the new contract details
    encryptInstance, err := common.NewEncryption(256) // Adjust key size if necessary
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Convert `NewCode` to byte format before encrypting
    encryptedCode, err := encryptInstance.EncryptData(newContract.NewCode, common.EncryptionKey, make([]byte, 12))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt new contract code: %v", err)
    }

    // Convert `NewParameters` to a string before encrypting
    parametersStr := fmt.Sprintf("%v", newContract.NewParameters)
    encryptedParams, err := encryptInstance.EncryptData(parametersStr, common.EncryptionKey, make([]byte, 12))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt new contract parameters: %v", err)
    }

    newContract.EncryptedCode = encryptedCode
    newContract.EncryptedParameters = encryptedParams

    // Store the migrated contract in the manager and ledger
    mm.Contracts[newContractID] = newContract
    err = mm.LedgerInstance.RecordMigration(oldContractID, newContractID, string(encryptedCode))
    if err != nil {
        return nil, fmt.Errorf("failed to record migration in the ledger: %v", err)
    }

    fmt.Printf("Contract %s successfully migrated to %s.\n", oldContractID, newContractID)
    return newContract, nil
}


// RetrieveMigratedContract retrieves the new version of a migrated contract.
func (mm *MigrationManager) RetrieveMigratedContract(newContractID string) (*MigratedContract, error) {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()

    // Check if the migrated contract exists in the manager
    contract, exists := mm.Contracts[newContractID]
    if !exists {
        return nil, errors.New("migrated contract not found")
    }

    // Create a decryption instance
    decryptInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return nil, fmt.Errorf("failed to create decryption instance: %v", err)
    }

    // Decrypt contract code
    decryptedCode, err := decryptInstance.DecryptData([]byte(contract.EncryptedCode), common.EncryptionKey) // Convert string to byte slice
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt contract code: %v", err)
    }

    // Decrypt contract parameters
    decryptedParams, err := decryptInstance.DecryptData([]byte(contract.EncryptedParameters), common.EncryptionKey) // Convert string to byte slice
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt contract parameters: %v", err)
    }

    // Assign decrypted values back to the contract
    contract.NewCode = string(decryptedCode) // Convert byte slice back to string
    contract.NewParameters = parseParameters(string(decryptedParams)) // Use your parsing function to handle params

    fmt.Printf("Migrated contract %s retrieved successfully.\n", newContractID)
    return contract, nil
}



// ListAllMigrations lists all the contracts that have been migrated.
func (mm *MigrationManager) ListAllMigrations() []*MigratedContract {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()

    var migrations []*MigratedContract
    for _, contract := range mm.Contracts {
        migrations = append(migrations, contract)
    }

    fmt.Printf("Listing all migrated contracts.\n")
    return migrations
}

// ValidateMigration ensures that the new migrated contract data matches the original contract data for integrity.
func (mm *MigrationManager) ValidateMigration(newContractID string) error {
    mm.mutex.Lock()
    defer mm.mutex.Unlock()

    // Retrieve the migrated contract
    contract, exists := mm.Contracts[newContractID]
    if !exists {
        return errors.New("migrated contract not found")
    }

    // Decrypt the contract code
    decryptInstance, err := common.NewEncryption(256) // Adjust key size if necessary
    if err != nil {
        return fmt.Errorf("failed to create decryption instance: %v", err)
    }

    decryptedCode, err := decryptInstance.DecryptData([]byte(contract.EncryptedCode), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt contract code: %v", err)
    }

    // Retrieve the migration record from the ledger
    ledgerRecord, err := mm.LedgerInstance.RetrieveMigrationRecord(contract.OldContractID, newContractID)
    if err != nil {
        return fmt.Errorf("failed to retrieve migration record from ledger: %v", err)
    }

    // Calculate the hash of the decrypted contract code
    hash := sha256.New()
    hash.Write(decryptedCode)
    calculatedHash := hex.EncodeToString(hash.Sum(nil))

    // Compare the hashes to ensure integrity
    if ledgerRecord != calculatedHash {
        return errors.New("migration validation failed: data tampered")
    }

    fmt.Printf("Migration for contract %s validated successfully.\n", newContractID)
    return nil
}



// generateNewContractID creates a new contract ID based on the old ID and new contract code.
func generateNewContractID(oldContractID, newCode string) string {
    hashInput := fmt.Sprintf("%s%s%d", oldContractID, newCode, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// parseParameters parses the contract parameters from the encrypted data.
// It expects the encrypted parameters to be in JSON format and returns a map of the parameters.
func parseParameters(params string) map[string]interface{} {
    var parsedParams map[string]interface{}

    // Attempt to unmarshal the JSON data into a map
    err := json.Unmarshal([]byte(params), &parsedParams)
    if err != nil {
        fmt.Printf("Error parsing parameters: %v\n", err)
        return make(map[string]interface{}) // Return empty map in case of an error
    }

    return parsedParams
}
