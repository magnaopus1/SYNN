package testnet

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)


// NewTestnetSimulation initializes a new testnet simulation.
func NewTestnetSimulation(testnet *common.TestnetNetwork) *common.TestnetSimulation {
    return &common.TestnetSimulation{
        Testnet:         testnet,
        transactionPool: []common.Transaction{},
    }
}

// AddTransactionToPool adds a transaction to the testnet transaction pool.
func (ts *common.TestnetSimulation) AddTransactionToPool(tx common.Transaction) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()

    ts.transactionPool = append(ts.transactionPool, tx)
    fmt.Printf("Transaction added to pool: %+v\n", tx)
}

// SimulateActivity continuously adds and processes transactions, validating sub-blocks and blocks.
func (ts *common.TestnetSimulation) SimulateActivity(interval time.Duration, transactionBatchSize int) {
    for {
        ts.mutex.Lock()
        // Check if there are enough transactions to process
        if len(ts.transactionPool) >= transactionBatchSize {
            // Take a batch of transactions and process them
            transactions := ts.transactionPool[:transactionBatchSize]
            ts.transactionPool = ts.transactionPool[transactionBatchSize:]

            err := ts.Testnet.AddSubBlock(transactions)
            if err != nil {
                fmt.Printf("Failed to add sub-block: %v\n", err)
            }
        } else {
            fmt.Println("Not enough transactions to create a sub-block. Waiting...")
        }
        ts.mutex.Unlock()
        time.Sleep(interval)
    }
}

// GenerateRandomTransactions creates random transactions for simulation purposes.
func (ts *common.TestnetSimulation) GenerateRandomTransactions(numTransactions int) {
    for i := 0; i < numTransactions; i++ {
        tx := common.Transaction{
            From:    fmt.Sprintf("wallet_%d", rand.Intn(1000)),
            To:      fmt.Sprintf("wallet_%d", rand.Intn(1000)),
            Amount:  rand.Float64() * 100,
            Fee:     rand.Float64(),
            Message: fmt.Sprintf("Testnet transaction %d", i),
        }
        ts.AddTransactionToPool(tx)
    }
}

// StartSimulation initializes the process to generate random transactions and process them into sub-blocks and blocks.
func (ts *common.TestnetSimulation) StartSimulation(transactionBatchSize int, interval time.Duration, transactionGenerationRate int) {
    go func() {
        for {
            ts.GenerateRandomTransactions(transactionGenerationRate)
            time.Sleep(interval / 2)
        }
    }()

    go ts.SimulateActivity(interval, transactionBatchSize)
}

// Encryption and Ledger Interaction

// EncryptTransactionPool encrypts the current state of the transaction pool before processing.
func (ts *common.TestnetSimulation) EncryptTransactionPool() ([]string, error) {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()

    var encryptedTransactions []string
    for _, tx := range ts.transactionPool {
        encryptedTx, err := encryption.EncryptData(fmt.Sprintf("%+v", tx), common.EncryptionKey)
        if err != nil {
            return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
        }
        encryptedTransactions = append(encryptedTransactions, encryptedTx)
    }
    return encryptedTransactions, nil
}

// EncryptAndRecordSubBlocks encrypts sub-blocks and stores them in the ledger.
func (ts *common.TestnetSimulation) EncryptAndRecordSubBlocks(subBlocks []common.SubBlock) error {
    for _, subBlock := range subBlocks {
        encryptedSubBlock, err := encryption.EncryptData(fmt.Sprintf("Sub-block %+v", subBlock), common.EncryptionKey)
        if err != nil {
            return fmt.Errorf("failed to encrypt sub-block: %v", err)
        }

        err = ts.Testnet.LedgerInstance.RecordSubBlock(subBlock)
        if err != nil {
            return fmt.Errorf("failed to record sub-block in ledger: %v", err)
        }

        fmt.Printf("Sub-block %d encrypted and recorded in the ledger.\n", subBlock.Index)
    }
    return nil
}

// ValidateSimulation ensures that all blocks in the testnet are valid.
func (ts *common.TestnetSimulation) ValidateSimulation() error {
    ts.mutex.Lock()
    defer ts.mutex.Unlock()

    for _, block := range ts.Testnet.Blocks {
        if !ts.Testnet.ConsensusEngine.ValidateBlock(block) {
            return fmt.Errorf("block %d failed validation", block.Index)
        }
    }

    fmt.Println("All blocks successfully validated in the testnet.")
    return nil
}

