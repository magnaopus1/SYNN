package smart_contract

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
)

// NewSmartLegalContractManager initializes the manager for smart legal contracts.
func NewSmartLegalContractManager(ledgerInstance *ledger.Ledger) *common.SmartLegalContractManager {
    return &common.SmartLegalContractManager{
        Contracts:      make(map[string]*common.SmartLegalContract),
        LedgerInstance: ledgerInstance,
    }
}

// DeployLegalContract deploys a new smart legal contract on the blockchain.
func (slcm *common.SmartLegalContractManager) DeployLegalContract(owner, contractTerms string, partiesInvolved []string) (*common.SmartLegalContract, error) {
    slcm.mutex.Lock()
    defer slcm.mutex.Unlock()

    contractID := generateLegalContractID(owner, contractTerms)
    legalContract := &common.SmartLegalContract{
        ID:              contractID,
        ContractTerms:   contractTerms,
        PartiesInvolved: partiesInvolved,
        Signatures:      make(map[string]string),
        State:           make(map[string]interface{}),
        Owner:           owner,
        LegallyBinding:  true,
        LedgerInstance:  slcm.LedgerInstance,
    }

    slcm.Contracts[legalContract.ID] = legalContract
    fmt.Printf("Smart Legal Contract %s deployed by %s.\n", legalContract.ID, owner)

    // Encrypt and store the contract deployment on the ledger
    encryptedContract, err := encryption.EncryptContract(legalContract, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt legal contract: %v", err)
    }

    err = slcm.LedgerInstance.RecordContractDeployment(legalContract.ID, encryptedContract)
    if err != nil {
        return nil, fmt.Errorf("failed to record legal contract deployment in the ledger: %v", err)
    }

    return legalContract, nil
}

// SignLegalContract allows parties to digitally sign the legal contract.
func (slc *common.SmartLegalContract) SignLegalContract(party, signature string) error {
    slc.mutex.Lock()
    defer slc.mutex.Unlock()

    // Check if the party is involved in the contract
    if !slc.isPartyInvolved(party) {
        return fmt.Errorf("party %s is not involved in this contract", party)
    }

    // Check if the party has already signed
    if _, signed := slc.Signatures[party]; signed {
        return fmt.Errorf("party %s has already signed the contract", party)
    }

    // Store the digital signature
    slc.Signatures[party] = signature
    fmt.Printf("Party %s signed the contract %s.\n", party, slc.ID)

    // Record the signature on the ledger
    encryptedSignature, err := encryption.EncryptData(fmt.Sprintf("%s:%s", party, signature), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt signature: %v", err)
    }

    err = slc.LedgerInstance.RecordContractSignature(slc.ID, encryptedSignature)
    if err != nil {
        return fmt.Errorf("failed to record signature in the ledger: %v", err)
    }

    return nil
}

// ExecuteLegalContract executes the legal contract once all parties have signed.
func (slc *common.SmartLegalContract) ExecuteLegalContract() (map[string]interface{}, error) {
    slc.mutex.Lock()
    defer slc.mutex.Unlock()

    // Check if all parties have signed the contract
    if !slc.allPartiesSigned() {
        return nil, fmt.Errorf("not all parties have signed the contract")
    }

    // Simulate contract execution logic
    result := make(map[string]interface{})
    result["status"] = "success"
    result["executed"] = true

    // Record execution
    execution := common.ContractExecution{
        ExecutionID: fmt.Sprintf("%s-exec-%d", slc.ID, time.Now().UnixNano()),
        ContractID:  slc.ID,
        Executor:    slc.Owner,
        Timestamp:   time.Now(),
        Result:      result,
    }
    slc.Executions = append(slc.Executions, execution)

    // Encrypt and store the execution in the ledger
    encryptedExecution, err := encryption.EncryptContractExecution(execution, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract execution: %v", err)
    }

    err = slc.LedgerInstance.RecordContractExecution(execution.ExecutionID, encryptedExecution)
    if err != nil {
        return nil, fmt.Errorf("failed to record contract execution in the ledger: %v", err)
    }

    fmt.Printf("Smart Legal Contract %s executed.\n", slc.ID)
    return result, nil
}

// QueryLegalContract queries the details of the legal contract.
func (slc *common.SmartLegalContract) QueryLegalContract() (map[string]interface{}, error) {
    slc.mutex.Lock()
    defer slc.mutex.Unlock()

    // Encrypt and return the current contract state
    encryptedState, err := encryption.EncryptContractState(slc.State, common.EncryptionKey)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt contract state: %v", err)
    }

    return encryptedState, nil
}

// Helper method to check if a party is involved in the contract.
func (slc *common.SmartLegalContract) isPartyInvolved(party string) bool {
    for _, p := range slc.PartiesInvolved {
        if p == party {
            return true
        }
    }
    return false
}

// Helper method to check if all parties have signed the contract.
func (slc *common.SmartLegalContract) allPartiesSigned() bool {
    return len(slc.Signatures) == len(slc.PartiesInvolved)
}

// generateLegalContractID generates a unique ID for the legal contract based on its owner and terms.
func generateLegalContractID(owner, contractTerms string) string {
    hashInput := fmt.Sprintf("%s%s%d", owner, contractTerms, time.Now().UnixNano())
    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}
