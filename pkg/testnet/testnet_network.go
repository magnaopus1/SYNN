package testnet

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// NewTestnetNetwork initializes a new Testnet network with token support.
func NewTestnetNetwork(validators []common.Validator, ledgerInstance *ledger.Ledger) *common.TestnetNetwork {
    return &common.TestnetNetwork{
        Validators:      validators,
        LedgerInstance:  ledgerInstance,
        ConsensusEngine: consensus.NewSynnergyConsensus(validators, ledgerInstance),
        TokenTestnet:    NewTokenTestnet(ledgerInstance), // Integrating TokenTestnet into the network
        SubBlocks:       []common.SubBlock{},
        Blocks:          []common.Block{},
    }
}

// AddSubBlock adds a new sub-block to the Testnet and validates it using Synnergy consensus.
func (tn *common.TestnetNetwork) AddSubBlock(transactions []common.Transaction) error {
    tn.mutex.Lock()
    defer tn.mutex.Unlock()

    if len(transactions) == 0 {
        return fmt.Errorf("no transactions to process")
    }

    // Generate the previous sub-block hash
    var prevHash string
    if len(tn.SubBlocks) > 0 {
        prevHash = tn.SubBlocks[len(tn.SubBlocks)-1].Hash
    }

    // Create and validate the new sub-block
    subBlock, err := tn.ConsensusEngine.ValidateSubBlock(transactions, prevHash)
    if err != nil {
        return fmt.Errorf("sub-block validation failed: %v", err)
    }

    // Add the new sub-block to the chain
    tn.SubBlocks = append(tn.SubBlocks, subBlock)

    // Check if it's time to create a full block
    if len(tn.SubBlocks) >= MaxSubBlocksPerBlock {
        err := tn.createBlockFromSubBlocks()
        if err != nil {
            return fmt.Errorf("failed to create block: %v", err)
        }
    }

    // Record the sub-block in the ledger
    err = tn.LedgerInstance.RecordSubBlock(subBlock)
    if err != nil {
        return fmt.Errorf("failed to record sub-block in the ledger: %v", err)
    }

    fmt.Printf("Sub-block %d successfully added to the Testnet.\n", subBlock.Index)
    return nil
}

// createBlockFromSubBlocks creates a full block from the accumulated sub-blocks.
func (tn *common.TestnetNetwork) createBlockFromSubBlocks() error {
    // Ensure we have enough sub-blocks to create a block
    if len(tn.SubBlocks) < MaxSubBlocksPerBlock {
        return fmt.Errorf("insufficient sub-blocks to create a full block")
    }

    prevBlockHash := ""
    if len(tn.Blocks) > 0 {
        prevBlockHash = tn.Blocks[len(tn.Blocks)-1].Hash
    }

    // Create the new block using the Synnergy Consensus
    newBlock, err := tn.ConsensusEngine.ValidateBlock(tn.SubBlocks, prevBlockHash)
    if err != nil {
        return fmt.Errorf("block validation failed: %v", err)
    }

    // Add the block to the chain
    tn.Blocks = append(tn.Blocks, newBlock)

    // Clear the sub-blocks for the next block
    tn.SubBlocks = nil

    // Record the block in the ledger
    err = tn.LedgerInstance.RecordBlock(newBlock)
    if err != nil {
        return fmt.Errorf("failed to record block in the ledger: %v", err)
    }

    fmt.Printf("Block %d successfully created and added to the Testnet.\n", newBlock.Index)
    return nil
}

// SimulateNetwork simulates Testnet network activity, adding and validating sub-blocks over time.
func (tn *common.TestnetNetwork) SimulateNetwork(transactionBatches [][]common.Transaction, interval time.Duration) {
    for _, batch := range transactionBatches {
        err := tn.AddSubBlock(batch)
        if err != nil {
            fmt.Printf("Error adding sub-block: %v\n", err)
        }
        time.Sleep(interval)
    }
}

// EncryptSubBlock encrypts a sub-block before storing it in the ledger.
func (tn *common.TestnetNetwork) EncryptSubBlock(subBlock *common.SubBlock) (string, error) {
    encryptedSubBlock, err := encryption.EncryptData(fmt.Sprintf("%+v", subBlock), common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt sub-block: %v", err)
    }
    return encryptedSubBlock, nil
}

// DeployUniversalToken deploys a new token to the Testnet using the TokenTestnet.
func (tn *common.TestnetNetwork) DeployUniversalToken(tokenID string, token common.TokenInterface) error {
    tn.mutex.Lock()
    defer tn.mutex.Unlock()

    return tn.TokenTestnet.DeployUniversalToken(tokenID, token)
}

// TransferTokens transfers tokens between addresses within the Testnet.
func (tn *common.TestnetNetwork) TransferTokens(tokenID, from, to string, amount float64) error {
    tn.mutex.Lock()
    defer tn.mutex.Unlock()

    return tn.TokenTestnet.TransferTokens(tokenID, from, to, amount)
}

// SimulateTokenTransactions simulates token-related transactions over time.
func (tn *common.TestnetNetwork) SimulateTokenTransactions(numTransactions int, interval time.Duration) {
    tn.mutex.Lock()
    defer tn.mutex.Unlock()

    tn.TokenTestnet.SimulateTokenTransactions(numTransactions, interval)
}

// Ledger Integration

// RecordSubBlock logs a sub-block in the ledger, ensuring testnet immutability.
func (l *ledger.Ledger) RecordSubBlock(subBlock common.SubBlock) error {
    encryptedSubBlock, err := encryption.EncryptData(fmt.Sprintf("Sub-block %+v", subBlock), common.EncryptionKey)
    if err != nil {
        return err
    }
    return l.RecordTransaction("SubBlock", "", float64(subBlock.Index), encryptedSubBlock)
}

// RecordBlock logs a full block in the ledger.
func (l *ledger.Ledger) RecordBlock(block common.Block) error {
    encryptedBlock, err := encryption.EncryptData(fmt.Sprintf("Block %+v", block), common.EncryptionKey)
    if err != nil {
        return err
    }
    return l.RecordTransaction("Block", "", float64(block.Index), encryptedBlock)
}
