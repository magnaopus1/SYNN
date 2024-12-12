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

// ContractInteraction represents an interaction between smart contracts on the blockchain.
type ContractInteraction struct {
	InteractionID            string                 // Unique identifier for the interaction
	InitiatingContractID      string                 // ID of the contract initiating the interaction
	ReceivingContractID       string                 // ID of the contract receiving the interaction
	Caller                    string                 // Address of the caller who initiated the interaction
	FunctionName              string                 // Name of the function called in the contract
	Parameters                map[string]interface{} // Parameters passed to the function
	Timestamp                 time.Time              // Timestamp of when the interaction took place
	ExecutionResult           string                 // Result of the execution (e.g., success, failure)
	InteractionData           string                 // Raw data of the interaction
	EncryptedInteractionData  string                 // Encrypted data of the interaction
}

// ContractInteractionManager manages the interactions between different smart contracts on the blockchain.
type ContractInteractionManager struct {
	Interactions   map[string]*ContractInteraction // A map of contract interactions by interaction ID
	LedgerInstance *ledger.Ledger                  // Ledger instance for storing interaction records
	mutex          sync.Mutex                      // Mutex for thread-safe operations
}


// NewContractInteractionManager initializes the contract interaction manager.
func NewContractInteractionManager(ledgerInstance *ledger.Ledger) *ContractInteractionManager {
    return &ContractInteractionManager{
        Interactions:   make(map[string]*ContractInteraction),
        LedgerInstance: ledgerInstance,
    }
}

// CreateInteraction initializes a new interaction between two smart contracts.
func (cim *ContractInteractionManager) CreateInteraction(initiatingContractID, receivingContractID, interactionData string) (string, error) {
    cim.mutex.Lock()
    defer cim.mutex.Unlock()

    interactionID := cim.generateInteractionID(initiatingContractID, receivingContractID)

    // Ensure interaction ID is unique
    if _, exists := cim.Interactions[interactionID]; exists {
        return "", fmt.Errorf("interaction with ID %s already exists", interactionID)
    }

    interaction := &ContractInteraction{
        InteractionID:        interactionID,
        InitiatingContractID: initiatingContractID,
        ReceivingContractID:  receivingContractID,
        InteractionData:      interactionData,
        Timestamp:            time.Now(),
    }

    // Step 1: Create an instance of Encryption
    encryptionInstance := &common.Encryption{}

    // Step 2: Define the EncryptionKey and encrypt interaction data before storing
    EncryptionKey := []byte("your-32-byte-key-for-aes-encryption")
    encryptedInteractionData, err := encryptionInstance.EncryptData("AES", []byte(interaction.InteractionData), EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt interaction data: %v", err)
    }
    interaction.EncryptedInteractionData = string(encryptedInteractionData)

    // Step 3: Store the interaction in the map
    cim.Interactions[interactionID] = interaction

    // Define placeholders or extract actual values for the additional parameters
    metadata := map[string]interface{}{"description": "Sample metadata for interaction"} // Replace with actual metadata as needed
    initiatorChainID := "chainA" // Placeholder; replace with actual chain ID
    receiverChainID := "chainB" // Placeholder; replace with actual chain ID
    interactionType := "sample_interaction_type" // Define the type of interaction

    // Step 4: Record the interaction in the ledger with all required arguments
    err = cim.LedgerInstance.RecordContractInteraction(
        interactionID,
        interaction.InitiatingContractID,
        interaction.ReceivingContractID,
        interaction.EncryptedInteractionData,
        metadata,
        initiatorChainID,
        receiverChainID,
        interactionType,
    )
    if err != nil {
        return "", fmt.Errorf("failed to record interaction in the ledger: %v", err)
    }

    fmt.Printf("Interaction %s created between contracts %s and %s.\n", interactionID, initiatingContractID, receivingContractID)
    return interactionID, nil
}




// RetrieveInteraction retrieves the interaction data between two smart contracts by interaction ID.
func (cim *ContractInteractionManager) RetrieveInteraction(interactionID string) (*ContractInteraction, error) {
    cim.mutex.Lock()
    defer cim.mutex.Unlock()

    interaction, exists := cim.Interactions[interactionID]
    if !exists {
        return nil, errors.New("interaction not found")
    }

    // Step 1: Create an instance of Encryption
    encryptionInstance := &common.Encryption{}

    // Step 2: Decrypt interaction data
    decryptedData, err := encryptionInstance.DecryptData([]byte(interaction.EncryptedInteractionData), common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt interaction data: %v", err)
    }
    interaction.InteractionData = string(decryptedData)

    fmt.Printf("Retrieved interaction %s between contracts %s and %s.\n", interactionID, interaction.InitiatingContractID, interaction.ReceivingContractID)
    return interaction, nil
}


// generateInteractionID creates a unique ID for each contract interaction.
func (cim *ContractInteractionManager) generateInteractionID(initiatingContractID, receivingContractID string) string {
    hashInput := fmt.Sprintf("%s%s%d", initiatingContractID, receivingContractID, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// ListInteractionsForContract lists all interactions for a specific contract.
func (cim *ContractInteractionManager) ListInteractionsForContract(contractID string) ([]*ContractInteraction, error) {
    cim.mutex.Lock()
    defer cim.mutex.Unlock()

    interactions := []*ContractInteraction{}
    for _, interaction := range cim.Interactions {
        if interaction.InitiatingContractID == contractID || interaction.ReceivingContractID == contractID {
            interactions = append(interactions, interaction)
        }
    }

    if len(interactions) == 0 {
        return nil, errors.New("no interactions found for the given contract ID")
    }

    fmt.Printf("Listing interactions for contract %s.\n", contractID)
    return interactions, nil
}

// ValidateContractInteraction ensures that the interaction is legitimate and not tampered with.
func (cim *ContractInteractionManager) ValidateContractInteraction(interactionID string) error {
    cim.mutex.Lock()
    defer cim.mutex.Unlock()

    // Retrieve the interaction from memory
    interaction, exists := cim.Interactions[interactionID]
    if !exists {
        return errors.New("interaction not found")
    }

    // Step 1: Create an instance of Encryption
    encryptionInstance := &common.Encryption{}

    // Step 2: Decrypt the interaction data
    decryptedData, err := encryptionInstance.DecryptData([]byte(interaction.EncryptedInteractionData), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt interaction data: %v", err)
    }

    // Step 3: Calculate the hash of the decrypted data
    hash := sha256.New()
    hash.Write([]byte(interaction.InitiatingContractID + interaction.ReceivingContractID + string(decryptedData)))
    calculatedHash := hex.EncodeToString(hash.Sum(nil))

    // Step 4: Retrieve the contract interaction record from the ledger
    ledgerRecords, err := cim.LedgerInstance.RetrieveContractInteraction(interactionID)
    if err != nil {
        return fmt.Errorf("failed to retrieve interaction from ledger: %v", err)
    }

    // Step 5: Iterate over the ledger records and find the matching one to validate
    for _, ledgerRecord := range ledgerRecords {
        // Recreate the hash from the ledgerRecord's interaction data
        ledgerRecordHash := sha256.New()
        ledgerRecordHash.Write([]byte(ledgerRecord.InitiatingContractID + ledgerRecord.ReceivingContractID + ledgerRecord.InteractionData))
        ledgerRecordCalculatedHash := hex.EncodeToString(ledgerRecordHash.Sum(nil))

        // Compare the calculated hash from memory to the ledger record hash
        if ledgerRecordCalculatedHash == calculatedHash {
            fmt.Printf("Interaction %s validated successfully.\n", interactionID)
            return nil // Success
        }
    }

    // If no matching ledger record was found, return an error
    return errors.New("interaction validation failed: no matching record found")
}



