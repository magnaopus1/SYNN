package plasma

import (
    "errors"
    "fmt"
    "time"

    "synnergy_network/pkg/common"
    "synnergy_network/pkg/encryption"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/cross_chain" // Hypothetical cross-chain package
)


// NewPlasmaCrossChain initializes the cross-chain capabilities in the Plasma childchain.
func NewPlasmaCrossChain(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, crossChainService *cross_chain.Communication) *common.PlasmaCrossChain {
    return &common.PlasmaCrossChain{
        Blocks:         make(map[string]*common.PlasmaBlock),
        SubBlocks:      make(map[string]*common.PlasmaSubBlock),
        Ledger:         ledgerInstance,
        Encryption:     encryptionService,
        CrossChainComm: crossChainService,
    }
}

// CreateSubBlock handles the creation of a new Plasma sub-block for cross-chain transactions.
func (pc *common.PlasmaCrossChain) CreateSubBlock(subBlockID string, txs []*common.PlasmaUTXO, subBlockIndex int) (*common.PlasmaSubBlock, error) {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if _, exists := pc.SubBlocks[subBlockID]; exists {
        return nil, errors.New("sub-block already exists")
    }

    // Create the sub-block with transactions
    subBlock := &common.PlasmaSubBlock{
        SubBlockID:    subBlockID,
        Transactions:  txs,
        SubBlockIndex: subBlockIndex,
        Timestamp:     time.Now(),
        MerkleRoot:    common.GenerateMerkleRoot(txs),
    }

    // Encrypt transactions
    for _, tx := range txs {
        encryptedTx, err := pc.Encryption.EncryptData([]byte(tx.UTXOID), common.EncryptionKey)
        if err != nil {
            return nil, fmt.Errorf("failed to encrypt transaction: %v", err)
        }
        tx.UTXOID = string(encryptedTx)
    }

    // Store the sub-block
    pc.SubBlocks[subBlockID] = subBlock

    // Log the sub-block creation
    err := pc.Ledger.RecordSubBlockCreation(subBlockID, time.Now())
    if err != nil {
        return nil, fmt.Errorf("failed to log sub-block creation: %v", err)
    }

    fmt.Printf("Sub-block %s created with %d transactions\n", subBlockID, len(txs))
    return subBlock, nil
}

// CreateBlock creates a Plasma block that contains 1000 sub-blocks, with cross-chain support.
func (pc *common.PlasmaCrossChain) CreateBlock(blockID, previousBlock string, subBlocks []*common.PlasmaSubBlock) (*common.PlasmaBlock, error) {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if len(subBlocks) != 1000 {
        return nil, errors.New("a Plasma block must contain exactly 1000 sub-blocks")
    }

    // Create the block and assign sub-blocks
    block := &common.PlasmaBlock{
        BlockID:       blockID,
        SubBlocks:     subBlocks,
        PreviousBlock: previousBlock,
        Timestamp:     time.Now(),
        MerkleRoot:    common.GenerateMerkleRoot(subBlocks),
    }

    // Log the block creation in the ledger
    err := pc.Ledger.RecordBlockCreation(blockID, previousBlock, time.Now())
    if err != nil {
        return nil, fmt.Errorf("failed to log block creation: %v", err)
    }

    // Broadcast block information to the cross-chain
    err = pc.CrossChainComm.BroadcastBlock(blockID)
    if err != nil {
        return nil, fmt.Errorf("failed to broadcast block for cross-chain operation: %v", err)
    }

    pc.Blocks[blockID] = block
    fmt.Printf("Block %s created with 1000 sub-blocks\n", blockID)
    return block, nil
}

// CrossChainTransfer initiates a cross-chain transfer of assets or data.
func (pc *common.PlasmaCrossChain) CrossChainTransfer(destinationChain string, tx *common.PlasmaUTXO) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if tx.IsSpent {
        return errors.New("cannot transfer spent UTXO")
    }

    // Encrypt the transaction before transfer
    encryptedTx, err := pc.Encryption.EncryptData([]byte(tx.UTXOID), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt transaction: %v", err)
    }

    // Initiate cross-chain transfer via communication layer
    err = pc.CrossChainComm.InitiateTransfer(destinationChain, encryptedTx)
    if err != nil {
        return fmt.Errorf("failed to initiate cross-chain transfer: %v", err)
    }

    // Log the cross-chain transfer
    err = pc.Ledger.RecordCrossChainTransfer(tx.UTXOID, destinationChain, time.Now())
    if err != nil {
        return fmt.Errorf("failed to log cross-chain transfer: %v", err)
    }

    fmt.Printf("Cross-chain transfer initiated for transaction %s to chain %s\n", tx.UTXOID, destinationChain)
    return nil
}

// ValidateCrossChainTransaction validates the integrity of an incoming cross-chain transaction.
func (pc *common.PlasmaCrossChain) ValidateCrossChainTransaction(txID string, sourceChain string) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    // Fetch cross-chain transaction details
    txDetails, err := pc.CrossChainComm.RetrieveTransaction(txID, sourceChain)
    if err != nil {
        return fmt.Errorf("failed to retrieve cross-chain transaction from chain %s: %v", sourceChain, err)
    }

    // Decrypt transaction data
    decryptedTx, err := pc.Encryption.DecryptData(txDetails.Data, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt cross-chain transaction: %v", err)
    }

    fmt.Printf("Cross-chain transaction %s validated from chain %s\n", txID, sourceChain)
    return nil
}

// FinalizeCrossChainTransaction finalizes a cross-chain transaction and integrates it into Plasma.
func (pc *common.PlasmaCrossChain) FinalizeCrossChainTransaction(txID string, destinationChain string) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    // Search for the transaction in the childchain
    for _, subBlock := range pc.SubBlocks {
        for _, tx := range subBlock.Transactions {
            if tx.UTXOID == txID && !tx.IsSpent {
                tx.IsSpent = true
                tx.SpentTime = time.Now()

                // Log the finalization event
                err := pc.Ledger.RecordCrossChainTransactionFinalization(txID, destinationChain, time.Now())
                if err != nil {
                    return fmt.Errorf("failed to log cross-chain transaction finalization: %v", err)
                }

                fmt.Printf("Cross-chain transaction %s finalized to chain %s\n", txID, destinationChain)
                return nil
            }
        }
    }

    return fmt.Errorf("transaction %s not found or already finalized", txID)
}
