package smart_contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewRicardianContractManager initializes the manager for Ricardian contracts.
func NewRicardianContractManager(ledgerInstance *ledger.Ledger) *RicardianContractManager {
    return &RicardianContractManager{
        Contracts:      make(map[string]*RicardianContract),
        LedgerInstance: ledgerInstance,
    }
}

// DeployRicardianContract deploys a new Ricardian contract to the blockchain.
func (rcm *RicardianContractManager) DeployRicardianContract(owner, humanReadable, machineReadable string, partiesInvolved []string) (*RicardianContract, error) {
    rcm.mutex.Lock()
    defer rcm.mutex.Unlock()

    // Generate a unique contract ID
    contractID := generateRicardianContractID(owner, humanReadable)

    // Create a new Ricardian contract instance
    contract := &RicardianContract{
        ID:              contractID,
        HumanReadable:   humanReadable,
        MachineReadable: machineReadable,
        PartiesInvolved: partiesInvolved,
        Signatures:      make(map[string]string),
        State:           make(map[string]interface{}),
        Owner:           owner,
    }

    // Store the contract in the manager's contracts map
    rcm.Contracts[contract.ID] = contract
    fmt.Printf("Ricardian Contract %s deployed by %s.\n", contract.ID, owner)

    // Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Adjust key size if needed
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Serialize and encrypt the Ricardian contract for storage (if needed)
    serializedContract, err := json.Marshal(contract)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize Ricardian contract: %v", err)
    }

    _, err = encryptionInstance.EncryptData(string(serializedContract), common.EncryptionKey, make([]byte, 12)) // Encrypt without assigning
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt Ricardian contract: %v", err)
    }

    // Record the contract deployment in the ledger
    err = rcm.LedgerInstance.RecordContractDeployment(contract.ID, owner, machineReadable)
    if err != nil {
        return nil, fmt.Errorf("failed to record Ricardian contract deployment in the ledger: %v", err)
    }

    return contract, nil
}


// SignRicardianContract allows parties to digitally sign the Ricardian contract.
func (rc *RicardianContract) SignRicardianContract(party, signature string) error {
    rc.mutex.Lock()
    defer rc.mutex.Unlock()

    // Check if the party is involved in the contract
    if !rc.isPartyInvolved(party) {
        return fmt.Errorf("party %s is not involved in this contract", party)
    }

    // Check if the party has already signed
    if _, signed := rc.Signatures[party]; signed {
        return fmt.Errorf("party %s has already signed the contract", party)
    }

    // Store the digital signature
    rc.Signatures[party] = signature
    fmt.Printf("Party %s signed the Ricardian contract %s.\n", party, rc.ID)

    // Create an encryption instance and encrypt the signature data
    encryptionInstance, err := common.NewEncryption(256) // Adjust key size as needed
    if err != nil {
        return fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Encrypt the signature information
    encryptedSignature, err := encryptionInstance.EncryptData(fmt.Sprintf("%s:%s", party, signature), common.EncryptionKey, make([]byte, 12)) // Add a nonce if needed
    if err != nil {
        return fmt.Errorf("failed to encrypt signature: %v", err)
    }

    // Record the signature in the ledger
    err = rc.LedgerInstance.RecordContractSignature(rc.ID, party, string(encryptedSignature))
    if err != nil {
        return fmt.Errorf("failed to record signature in the ledger: %v", err)
    }

    return nil
}


// ExecuteRicardianContract executes the machine-readable terms of the Ricardian contract.
func (rc *RicardianContract) ExecuteRicardianContract(vm common.VirtualMachine, parameters map[string]interface{}) (map[string]interface{}, error) {
    rc.mutex.Lock()
    defer rc.mutex.Unlock()

    // Ensure all parties have signed the contract
    if !rc.allPartiesSigned() {
        return nil, fmt.Errorf("not all parties have signed the contract")
    }

    // Create an encryption instance
    encryptionInstance, err := common.NewEncryption(256) // Adjust key size if needed
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption instance: %v", err)
    }

    // Serialize contract data for optional encryption during execution
    serializedParams, err := json.Marshal(parameters)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize parameters for encryption: %v", err)
    }
    encryptedParams, err := encryptionInstance.EncryptData(string(serializedParams), common.EncryptionKey, make([]byte, 12))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt execution parameters: %v", err)
    }

    // Execute the contract function using the virtual machine
    result, err := vm.ExecuteContract(rc.MachineReadable, rc.ID, rc.Owner, parameters, encryptionInstance, encryptedParams)
    if err != nil {
        return nil, fmt.Errorf("execution of Ricardian contract failed: %v", err)
    }

    // Record the contract execution details
    execution := ContractExecution{
        ExecutionID:   fmt.Sprintf("%s-exec-%d", rc.ID, time.Now().UnixNano()),
        ContractID:    rc.ID,
        FunctionName:  "ExecuteRicardianContract",
        Parameters:    parameters,
        Result:        result,
        ExecutionTime: time.Now(),
        Executor:      rc.Owner,
        GasUsed:       0.0, // Set gas used based on actual execution if applicable
    }
    rc.Executions = append(rc.Executions, execution)

    // Serialize the execution record for ledger storage
    serializedExecution, err := json.Marshal(execution)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize contract execution: %v", err)
    }

    encryptedExecution, err := encryptionInstance.EncryptData(string(serializedExecution), common.EncryptionKey, make([]byte, 12))
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract execution: %v", err)
    }

    // Prepare execution details as map for ledger record
    executionDetails := map[string]interface{}{
        "ExecutionID": execution.ExecutionID,
        "EncryptedExecution": string(encryptedExecution),
    }

    // Record the encrypted execution in the ledger
    err = rc.LedgerInstance.RecordContractExecution(execution.ExecutionID, executionDetails)
    if err != nil {
        return nil, fmt.Errorf("failed to record contract execution in the ledger: %v", err)
    }

    fmt.Printf("Ricardian Contract %s executed.\n", rc.ID)
    return result, nil
}


// QueryRicardianContract queries the state and details of the Ricardian contract.
func (rc *RicardianContract) QueryRicardianContract() (map[string]interface{}, error) {
    rc.mutex.Lock()
    defer rc.mutex.Unlock()

    // Encrypt the current contract state
    encryptedState, err := common.EncryptContractState(rc.State, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract state: %v", err)
    }

    // Return the encrypted state in a map
    return map[string]interface{}{
        "EncryptedState": encryptedState,
    }, nil
}


// Helper method to check if a party is involved in the contract.
func (rc *RicardianContract) isPartyInvolved(party string) bool {
    for _, p := range rc.PartiesInvolved {
        if p == party {
            return true
        }
    }
    return false
}

// Helper method to check if all parties have signed the contract.
func (rc *RicardianContract) allPartiesSigned() bool {
    return len(rc.Signatures) == len(rc.PartiesInvolved)
}

// generateRicardianContractID generates a unique ID for the Ricardian contract based on the owner and terms.
func generateRicardianContractID(owner, humanReadable string) string {
    hashInput := fmt.Sprintf("%s%s%d", owner, humanReadable, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
