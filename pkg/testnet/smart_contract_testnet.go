package testnet

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewSmartContractTestnet initializes the smart contract testnet with the ledger and testnet network.
func NewSmartContractTestnet(testnet *common.TestnetNetwork) *common.SmartContractTestnet {
    return &common.SmartContractTestnet{
        Testnet:         testnet,
        ContractManager: smart_contract.NewSmartContractManager(testnet.LedgerInstance),
        transactionPool: []common.Transaction{},
    }
}

// DeploySmartContract deploys a new smart contract to the testnet and records it on the ledger.
func (sct *common.SmartContractTestnet) DeploySmartContract(owner, code string, parameters map[string]interface{}) (*common.SmartContract, error) {
    sct.mutex.Lock()
    defer sct.mutex.Unlock()

    contract, err := sct.ContractManager.DeployContract(owner, code, parameters)
    if err != nil {
        return nil, fmt.Errorf("failed to deploy smart contract: %v", err)
    }

    fmt.Printf("Smart Contract %s deployed successfully by %s.\n", contract.ID, owner)
    return contract, nil
}

// ExecuteSmartContract simulates the execution of a smart contract in the testnet environment.
func (sct *common.SmartContractTestnet) ExecuteSmartContract(contractID, sender string, parameters map[string]interface{}) (map[string]interface{}, error) {
    sct.mutex.Lock()
    defer sct.mutex.Unlock()

    contract, exists := sct.ContractManager.Contracts[contractID]
    if !exists {
        return nil, fmt.Errorf("contract with ID %s not found", contractID)
    }

    result, err := contract.ExecuteContract(sender, parameters)
    if err != nil {
        return nil, fmt.Errorf("execution failed for contract %s: %v", contractID, err)
    }

    fmt.Printf("Smart Contract %s executed by %s with result: %+v\n", contractID, sender, result)
    return result, nil
}

// SimulateSmartContractActivity continuously deploys and executes random smart contracts for testing purposes.
func (sct *common.SmartContractTestnet) SimulateSmartContractActivity(interval time.Duration, batchSize int) {
    go func() {
        for {
            // Deploy and execute smart contracts randomly
            sct.mutex.Lock()
            owner := fmt.Sprintf("wallet_%d", rand.Intn(1000))
            code := "contract_code_placeholder"
            parameters := map[string]interface{}{"param": rand.Intn(100)}

            contract, err := sct.DeploySmartContract(owner, code, parameters)
            if err != nil {
                fmt.Printf("Failed to deploy contract: %v\n", err)
                sct.mutex.Unlock()
                time.Sleep(interval)
                continue
            }

            sct.mutex.Unlock()
            time.Sleep(interval)

            // Execute the contract
            sct.ExecuteSmartContract(contract.ID, owner, parameters)
            time.Sleep(interval)
        }
    }()
}

// EncryptAndRecordContract encrypts the contract and records it on the ledger.
func (sct *common.SmartContractTestnet) EncryptAndRecordContract(contract *common.SmartContract) error {
    encryptedContract, err := encryption.EncryptContract(contract, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt contract: %v", err)
    }

    err = sct.Testnet.LedgerInstance.ledger.RecordContractDeployment(contract.ID, encryptedContract)
    if err != nil {
        return fmt.Errorf("failed to record contract deployment in the ledger: %v", err)
    }

    fmt.Printf("Smart Contract %s encrypted and recorded in the ledger.\n", contract.ID)
    return nil
}

// ValidateSmartContracts validates the integrity of all deployed smart contracts.
func (sct *common.SmartContractTestnet) ValidateSmartContracts() error {
    sct.mutex.Lock()
    defer sct.mutex.Unlock()

    for _, contract := range sct.ContractManager.Contracts {
        if !sct.validateContract(contract) {
            return fmt.Errorf("smart contract %s failed validation", contract.ID)
        }
    }

    fmt.Println("All smart contracts validated successfully.")
    return nil
}

// validateContract validates the integrity of an individual smart contract.
func (sct *common.SmartContractTestnet) validateContract(contract *common.SmartContract) bool {
    hash := sha256.New()
    hashInput := fmt.Sprintf("%s%s%d", contract.ID, contract.Owner, time.Now().UnixNano())
    hash.Write([]byte(hashInput))

    contractHash := hex.EncodeToString(hash.Sum(nil))
    return contractHash == contract.ID
}

// RecordContractExecution stores the execution details of a contract to the ledger.
func (sct *common.SmartContractTestnet) RecordContractExecution(execution common.ContractExecution) error {
    encryptedExecution, err := encryption.EncryptContractExecution(execution, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt contract execution: %v", err)
    }

    err = sct.Testnet.LedgerInstance.RecordContractExecution(execution.ExecutionID, encryptedExecution)
    if err != nil {
        return fmt.Errorf("failed to record contract execution in the ledger: %v", err)
    }

    fmt.Printf("Contract execution %s recorded in the ledger.\n", execution.ExecutionID)
    return nil
}

// SimulateSmartContractTransactions generates random contract-related transactions and processes them.
func (sct *common.SmartContractTestnet) SimulateSmartContractTransactions(numTransactions int, interval time.Duration) {
    go func() {
        for {
            for i := 0; i < numTransactions; i++ {
                sct.mutex.Lock()
                tx := common.Transaction{
                    From:   fmt.Sprintf("wallet_%d", rand.Intn(1000)),
                    To:     fmt.Sprintf("contract_%d", rand.Intn(100)),
                    Amount: rand.Float64() * 100,
                    Fee:    rand.Float64(),
                    Message: fmt.Sprintf("Contract transaction %d", i),
                }
                sct.transactionPool = append(sct.transactionPool, tx)
                sct.mutex.Unlock()

                fmt.Printf("Transaction added to pool: %+v\n", tx)
            }
            time.Sleep(interval)
        }
    }()
}
