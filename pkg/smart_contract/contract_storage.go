package smart_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// StoredContract represents the data for a stored contract in the blockchain.
type StoredContract struct {
    ContractID         string                 // Unique identifier for the contract
    Code               string                 // Smart contract code (e.g., bytecode or source)
    Owner              string                 // Address of the contract owner
    State              map[string]interface{} // Current state of the contract
    CreationTime       time.Time              // Time when the contract was deployed
    EncryptedCode      []byte                 // Encrypted smart contract code
    EncryptedState     []byte                 // Encrypted state or parameters of the contract
}


// ContractStorageManager manages the storage of contract data on the blockchain.
type ContractStorageManager struct {
	Contracts      map[string]*StoredContract // A map to store contract data by contract ID
	LedgerInstance *ledger.Ledger             // Ledger instance for storing contract data
	mutex          sync.Mutex                 // Mutex for thread-safe operations
}

// NewContractStorageManager initializes the contract storage manager.
func NewContractStorageManager(ledgerInstance *ledger.Ledger) *ContractStorageManager {
    return &ContractStorageManager{
        Contracts:      make(map[string]*StoredContract),
        LedgerInstance: ledgerInstance,
    }
}

// StoreContract securely stores the code and parameters of a smart contract.
func (csm *ContractStorageManager) StoreContract(contractID, owner, code string, parameters map[string]interface{}) error {
    csm.mutex.Lock()
    defer csm.mutex.Unlock()

    // Ensure contract ID is unique
    if _, exists := csm.Contracts[contractID]; exists {
        return fmt.Errorf("contract with ID %s already exists", contractID)
    }

    // Create an instance of the encryption system
    encryptionInstance := &common.Encryption{}

    // Encrypt the contract code and parameters before storage
    encryptedCode, err := encryptionInstance.EncryptData("AES", []byte(code), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt contract code: %v", err)
    }

    encryptedParams, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%v", parameters)), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt contract parameters: %v", err)
    }

    // Create the contract object with the encrypted data
    contract := &StoredContract{
        ContractID:     contractID,
        Owner:          owner,
        Code:           code,
        State:          parameters,
        CreationTime:   time.Now(),
        EncryptedCode:  encryptedCode,
        EncryptedState: encryptedParams,
    }

    // Store the contract in the map and ledger
    csm.Contracts[contractID] = contract
    err = csm.LedgerInstance.RecordContractStorage(contractID, string(encryptedCode), string(encryptedParams))
    if err != nil {
        return fmt.Errorf("failed to record contract in the ledger: %v", err)
    }

    fmt.Printf("Contract %s stored securely by %s.\n", contractID, owner)
    return nil
}


// RetrieveContract retrieves the stored contract data.
func (csm *ContractStorageManager) RetrieveContract(contractID string) (*StoredContract, error) {
    csm.mutex.Lock()
    defer csm.mutex.Unlock()

    contract, exists := csm.Contracts[contractID]
    if !exists {
        return nil, errors.New("contract not found")
    }

    // Step 1: Create an instance of Encryption
    encryptionInstance := &common.Encryption{}

    // Step 2: Decrypt the contract code
    decryptedCodeBytes, err := encryptionInstance.DecryptData(contract.EncryptedCode, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt contract code: %v", err)
    }
    decryptedCode := string(decryptedCodeBytes) // Convert decrypted bytes to string

    // Step 3: Decrypt the contract state/parameters
    decryptedStateBytes, err := encryptionInstance.DecryptData(contract.EncryptedState, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt contract state: %v", err)
    }
    decryptedState := string(decryptedStateBytes) // Convert decrypted bytes to string

    // Step 4: Update the contract with decrypted data
    contract.Code = decryptedCode
    contract.State = parseParameters(decryptedState) // Parse decrypted state

    fmt.Printf("Contract %s retrieved successfully.\n", contractID)
    return contract, nil
}



// ListStoredContracts lists all contracts stored in the blockchain.
func (csm *ContractStorageManager) ListStoredContracts() []*StoredContract {
    csm.mutex.Lock()
    defer csm.mutex.Unlock()

    var contracts []*StoredContract
    for _, contract := range csm.Contracts {
        contracts = append(contracts, contract)
    }

    fmt.Printf("Listing all stored contracts.\n")
    return contracts
}

// DeleteContract securely removes a contract from storage.
func (csm *ContractStorageManager) DeleteContract(contractID string) error {
    csm.mutex.Lock()
    defer csm.mutex.Unlock()

    // Check if the contract exists before deletion
    if _, exists := csm.Contracts[contractID]; !exists {
        return fmt.Errorf("contract with ID %s not found", contractID)
    }

    // Remove the contract from in-memory storage
    delete(csm.Contracts, contractID)

    // Remove the contract from the ledger, passing the contractID and a reason
    err := csm.LedgerInstance.RemoveContractStorage(contractID, "contract deletion")
    if err != nil {
        return fmt.Errorf("failed to remove contract from ledger: %v", err)
    }

    fmt.Printf("Contract %s deleted successfully.\n", contractID)
    return nil
}


// ValidateContractStorage ensures that the stored contract data has not been tampered with.
func (csm *ContractStorageManager) ValidateContractStorage(contractID string) error {
    csm.mutex.Lock()
    defer csm.mutex.Unlock()

    contract, exists := csm.Contracts[contractID]
    if !exists {
        return errors.New("contract not found")
    }

    // Step 1: Create an instance of Encryption
    encryptionInstance := &common.Encryption{}

    // Step 2: Decrypt the contract code
    decryptedCode, err := encryptionInstance.DecryptData([]byte(contract.EncryptedCode), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt contract code: %v", err)
    }

    // Step 3: Retrieve contract storage from the ledger with two arguments (contractID and reason/identifier)
    ledgerRecord, err := csm.LedgerInstance.RetrieveContractStorage(contractID, "validation check")
    if err != nil {
        return fmt.Errorf("failed to retrieve contract storage from ledger: %v", err)
    }

    // Step 4: Calculate the hash of the decrypted contract code
    hash := sha256.New()
    hash.Write(decryptedCode)
    calculatedHash := hex.EncodeToString(hash.Sum(nil))

    // Step 5: Compare the calculated hash with the ledger record
    if ledgerRecord != calculatedHash {
        return errors.New("contract validation failed: data tampered")
    }

    fmt.Printf("Contract %s validated successfully.\n", contractID)
    return nil
}


