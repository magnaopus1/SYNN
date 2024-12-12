package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strings"
	"synnergy_network/pkg/ledger"
	"time"
)

// PoW represents the Proof of Work mechanism, integrated with the Ledger
type PoW struct {
	State          PoWState        // The current state of the PoW system
	LedgerInstance *ledger.Ledger  // Ledger instance for recording blocks and rewards
}

// PoWState holds the current state of the Proof of Work system
type PoWState struct {
	Difficulty   int       // Current difficulty for PoW
	BlockReward  float64   // Block reward for the miner
	MinerAddress string    // Address of the miner
	Epoch        int       // Current epoch
	LastHash     string    // Hash of the last mined block
}


// NewPoW initializes the PoW system with the given parameters.
func NewPoW(difficulty int, blockReward float64, minerAddress string, ledgerInstance *ledger.Ledger) *PoW {
    if difficulty <= 0 {
        log.Printf("[Warning] Difficulty must be greater than 0. Defaulting to 1.")
        difficulty = 1
    }
    if blockReward <= 0 {
        log.Printf("[Warning] Block reward must be positive. Defaulting to 10.")
        blockReward = 10.0
    }
    if minerAddress == "" {
        log.Fatalf("[Error] Miner address cannot be empty.")
    }
    if ledgerInstance == nil {
        log.Fatalf("[Error] Ledger instance cannot be nil.")
    }

    log.Printf("[Info] Initializing PoW system with Difficulty=%d, BlockReward=%.2f, MinerAddress=%s.",
        difficulty, blockReward, minerAddress)

    return &PoW{
        State: PoWState{
            Difficulty:   difficulty,
            BlockReward:  blockReward,
            MinerAddress: minerAddress,
            Epoch:        0,
            LastHash:     "genesis_hash",
        },
        LedgerInstance: ledgerInstance,
    }
}



// MineBlock attempts to mine a new block by solving the PoW puzzle.
func (pow *PoW) MineBlock(block *Block) error {
    if block == nil {
        return fmt.Errorf("[Error] Block cannot be nil")
    }

    log.Printf("[Info] Mining block #%d with previous hash: %s.", block.Index, pow.State.LastHash)

    // Step 1: Set block properties
    block.Timestamp = time.Now()
    block.PrevHash = pow.State.LastHash
    block.Difficulty = pow.State.Difficulty

    // Step 2: Validate sub-blocks
    log.Printf("[Info] Validating %d sub-blocks for block #%d...", len(block.SubBlocks), block.Index)
    for _, subBlock := range block.SubBlocks {
        if !pow.ValidateSubBlock(subBlock) {
            return fmt.Errorf("[Error] Validation failed for sub-block #%d in block #%d", subBlock.Index, block.Index)
        }
    }
    log.Printf("[Success] All sub-blocks validated successfully for block #%d.", block.Index)

    // Step 3: Mining process
    start := time.Now()
    miningSuccess := false

    for block.Nonce = 0; block.Nonce < math.MaxInt64; block.Nonce++ {
        block.Hash = pow.calculateBlockHash(block)

        if pow.isValidHash(block.Hash) {
            log.Printf("[Success] Block #%d mined with hash: %s and nonce: %d.", block.Index, block.Hash, block.Nonce)

            // Step 4: Update PoW state
            pow.State.LastHash = block.Hash
            pow.State.Epoch++

            // Step 5: Reward the miner
            block.MinerReward = pow.State.BlockReward
            log.Printf("[Info] Miner %s rewarded %.2f SYNN for block #%d.", pow.State.MinerAddress, block.MinerReward)

            // Step 6: Record block in the ledger
            ledgerBlock := ConvertToLedgerBlock(*block)
            err := pow.LedgerInstance.BlockchainConsensusCoinLedger.RecordMinedBlock(ledgerBlock, pow.State.MinerAddress, pow.State.BlockReward)
            if err != nil {
                return fmt.Errorf("[Error] Failed to record mined block in ledger: %v", err)
            }

            // Step 7: Adjust difficulty dynamically
            pow.AdjustDifficulty(time.Since(start))

            miningSuccess = true
            break
        }
    }

    if !miningSuccess {
        return fmt.Errorf("[Error] Mining failed for block #%d. Maximum nonce reached.", block.Index)
    }

    log.Printf("[Success] Block #%d successfully mined and recorded.", block.Index)
    return nil
}


// ValidateSubBlock validates a single sub-block before mining the block.
func (pow *PoW) ValidateSubBlock(subBlock SubBlock) bool {
    log.Printf("[Info] Validating sub-block #%d...", subBlock.Index)

    // Step 1: Verify hash integrity of the sub-block
    recalculatedHash := pow.calculateSubBlockHash(subBlock)
    if recalculatedHash != subBlock.Hash {
        log.Printf("[Error] Sub-block #%d hash mismatch. Expected: %s, Got: %s", subBlock.Index, recalculatedHash, subBlock.Hash)
        return false
    }
    log.Printf("[Success] Sub-block #%d hash integrity verified.", subBlock.Index)

    // Step 2: Check consensus compliance
    if !pow.CheckSubBlockCompliance(subBlock) {
        log.Printf("[Error] Sub-block #%d failed consensus compliance checks.", subBlock.Index)
        return false
    }
    log.Printf("[Success] Sub-block #%d adheres to consensus rules.", subBlock.Index)

    // Step 3: Validate all transactions within the sub-block
    for _, tx := range subBlock.Transactions {
        if !pow.ValidateTransaction(tx) {
            log.Printf("[Error] Transaction %s in sub-block #%d failed validation.", tx.TransactionID, subBlock.Index)
            return false
        }
    }
    log.Printf("[Success] All transactions in sub-block #%d validated successfully.", subBlock.Index)

    return true
}


// ValidateTransaction checks if a transaction within a PoW block is valid.
func (pow *PoW) ValidateTransaction(tx Transaction) bool {
    log.Printf("[Info] Validating transaction %s for PoW...", tx.TransactionID)

    // Step 1: Validate the transaction type
    transactionType, typeExists := pow.TransactionTypeMap[tx.TransactionType]
    if !typeExists {
        log.Printf("[Error] Unknown TransactionType: %s for transaction %s.", tx.TransactionType, tx.TransactionID)
        return false
    }

    // Step 2: Validate the transaction function
    transactionFunction, funcExists := pow.TransactionFunctionMap[tx.TransactionFunction]
    if !funcExists {
        log.Printf("[Error] Unknown TransactionFunction: %s for transaction %s.", tx.TransactionFunction, tx.TransactionID)
        return false
    }

    // Step 3: Derive the validation key
    validationKey := fmt.Sprintf("%s:%s", transactionType, transactionFunction)
    validationFunc, validationExists := pow.ValidationMap[validationKey]
    if !validationExists {
        log.Printf("[Error] No validation logic found for key: %s in transaction %s.", validationKey, tx.TransactionID)
        return false
    }

    // Step 4: Validate the transaction using the mapped function
    if !validationFunc(tx) {
        log.Printf("[Error] Validation failed for transaction %s using key %s.", tx.TransactionID, validationKey)
        return false
    }

    log.Printf("[Success] Transaction %s validated successfully for PoW.", tx.TransactionID)
    return true
}


// AdjustDifficulty dynamically adjusts the mining difficulty based on block time.
func (pow *PoW) AdjustDifficulty(blockTime time.Duration) {
    const targetTime = 10 * time.Second // Target block time

    if blockTime < targetTime {
        pow.State.Difficulty++
        log.Printf("[Info] Increased difficulty to %d due to faster block time (%.2fs).", pow.State.Difficulty, blockTime.Seconds())
    } else if blockTime > targetTime && pow.State.Difficulty > 1 {
        pow.State.Difficulty--
        log.Printf("[Info] Decreased difficulty to %d due to slower block time (%.2fs).", pow.State.Difficulty, blockTime.Seconds())
    } else {
        log.Printf("[Info] Difficulty remains unchanged at %d.", pow.State.Difficulty)
    }
}



// calculateSubBlockHash computes the hash of a sub-block.
func (pow *PoW) calculateSubBlockHash(subBlock SubBlock) string {
    transactionHashes := extractTransactionHashes(subBlock.Transactions)
    input := fmt.Sprintf("%d:%s:%s:%s", subBlock.Index, subBlock.PrevHash, transactionHashes, subBlock.Timestamp)
    
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}



// checkSubBlockCompliance verifies that the sub-block adheres to the consensus rules.
func (pow *PoW) checkSubBlockCompliance(subBlock SubBlock) bool {
    log.Printf("[Info] Checking consensus compliance for sub-block #%d...", subBlock.Index)

    // Rule 1: Timestamp validation
    maxAllowedTimeDifference := time.Minute
    if time.Since(subBlock.Timestamp) > maxAllowedTimeDifference {
        log.Printf("[Error] Sub-block #%d has a timestamp too old (Timestamp: %s).", subBlock.Index, subBlock.Timestamp)
        return false
    }
    if subBlock.Timestamp.After(time.Now()) {
        log.Printf("[Error] Sub-block #%d has a future timestamp (Timestamp: %s).", subBlock.Index, subBlock.Timestamp)
        return false
    }

    // Rule 2: Previous hash validation
    if subBlock.PrevHash != pow.State.LastHash {
        log.Printf("[Error] Sub-block #%d has an invalid previous hash. Expected: %s, Got: %s.", subBlock.Index, pow.State.LastHash, subBlock.PrevHash)
        return false
    }

    // Rule 3: Transaction count limit
    const maxTransactionsPerSubBlock = 1000
    if len(subBlock.Transactions) > maxTransactionsPerSubBlock {
        log.Printf("[Error] Sub-block #%d exceeds transaction limit. Limit: %d, Got: %d.", subBlock.Index, maxTransactionsPerSubBlock, len(subBlock.Transactions))
        return false
    }

    log.Printf("[Success] Sub-block #%d complies with consensus rules.", subBlock.Index)
    return true
}


// validateTransaction validates a single transaction within a sub-block.
func (pow *PoW) validateTransaction(tx Transaction) bool {
    log.Printf("[Info] Validating transaction %s...", tx.TransactionID)

    // Step 1: Validate Transaction Type
    transactionType, typeExists := pow.TransactionTypeMap[tx.TransactionType]
    if !typeExists {
        log.Printf("[Error] Unknown TransactionType: %s for transaction %s.", tx.TransactionType, tx.TransactionID)
        return false
    }

    // Step 2: Validate Transaction Function
    transactionFunction, funcExists := pow.TransactionFunctionMap[tx.TransactionFunction]
    if !funcExists {
        log.Printf("[Error] Unknown TransactionFunction: %s for transaction %s.", tx.TransactionFunction, tx.TransactionID)
        return false
    }

    // Step 3: Derive the validation key
    validationKey := fmt.Sprintf("%s:%s", transactionType, transactionFunction)
    validationFunc, validationExists := pow.ValidationMap[validationKey]
    if !validationExists {
        log.Printf("[Error] No validation logic found for key: %s in transaction %s.", validationKey, tx.TransactionID)
        return false
    }

    // Step 4: Verify Transaction Signature
    transactionHash := computeTransactionHash(tx)
    if !pow.VerifySignature(tx.FromAddress, tx.Signature, transactionHash) {
        log.Printf("[Error] Transaction %s has an invalid signature.", tx.TransactionID)
        return false
    }

    // Step 5: Validate Sender Balance
    senderBalance, err := pow.LedgerInstance.GetBalance(tx.FromAddress, tx.TokenStandard)
    if err != nil {
        log.Printf("[Error] Failed to retrieve balance for address %s: %v.", tx.FromAddress, err)
        return false
    }
    if senderBalance < tx.Amount {
        log.Printf("[Error] Insufficient balance for transaction %s. Required: %.2f, Available: %.2f.", tx.TransactionID, tx.Amount, senderBalance)
        return false
    }

    // Step 6: Execute Validation Function
    if !validationFunc(tx) {
        log.Printf("[Error] Validation logic failed for transaction %s using key %s.", tx.TransactionID, validationKey)
        return false
    }

    log.Printf("[Success] Transaction %s validated successfully.", tx.TransactionID)
    return true
}


// calculateBlockHash generates a SHA-256 hash for the block.
func (pow *PoW) calculateBlockHash(block *Block) string {
    hashInput := fmt.Sprintf("%d:%s:%s:%d:%s", 
        block.Index, 
        block.Timestamp.Format(time.RFC3339Nano), 
        block.PrevHash, 
        block.Nonce, 
        pow.concatenateSubBlockHashes(block.SubBlocks))

    hash := sha256.New()
    hash.Write([]byte(hashInput))
    return hex.EncodeToString(hash.Sum(nil))
}

// concatenateSubBlockHashes generates a single string from all sub-block hashes.
func (pow *PoW) concatenateSubBlockHashes(subBlocks []SubBlock) string {
    var subBlockHashes []string
    for _, subBlock := range subBlocks {
        subBlockHashes = append(subBlockHashes, subBlock.Hash)
    }
    return strings.Join(subBlockHashes, "")
}




// ValidateBlock validates a mined block by recalculating its hash, checking difficulty, and validating sub-blocks.
func (pow *PoW) ValidateBlock(block *Block) bool {
    log.Printf("[Info] Validating block #%d...", block.Index)

    // Step 1: Recalculate and verify the block hash
    recalculatedHash := pow.calculateBlockHash(block)
    if recalculatedHash != block.Hash {
        log.Printf("[Error] Block #%d hash mismatch. Expected: %s, Got: %s.", block.Index, recalculatedHash, block.Hash)
        return false
    }
    log.Printf("[Success] Block #%d hash verified.", block.Index)

    // Step 2: Verify the block hash meets the difficulty requirement
    if !pow.isValidHash(block.Hash) {
        log.Printf("[Error] Block #%d hash does not meet the difficulty requirement.", block.Index)
        return false
    }
    log.Printf("[Success] Block #%d hash meets difficulty requirements.", block.Index)

    // Step 3: Validate all sub-blocks within the block
    log.Printf("[Info] Validating %d sub-blocks in block #%d...", len(block.SubBlocks), block.Index)
    for _, subBlock := range block.SubBlocks {
        if !pow.validateSubBlock(subBlock) {
            log.Printf("[Error] Sub-block #%d failed validation in block #%d.", subBlock.Index, block.Index)
            return false
        }
    }
    log.Printf("[Success] All sub-blocks in block #%d validated successfully.", block.Index)

    log.Printf("[Success] Block #%d validated successfully.", block.Index)
    return true
}

// validateSubBlock validates a single sub-block before block validation.
func (pow *PoW) validateSubBlock(subBlock SubBlock) bool {
    log.Printf("[Info] Validating sub-block #%d...", subBlock.Index)

    // Step 1: Verify hash integrity
    recalculatedHash := pow.calculateSubBlockHash(subBlock)
    if recalculatedHash != subBlock.Hash {
        log.Printf("[Error] Sub-block #%d hash mismatch. Expected: %s, Got: %s.", subBlock.Index, recalculatedHash, subBlock.Hash)
        return false
    }
    log.Printf("[Success] Sub-block #%d hash integrity verified.", subBlock.Index)

    // Step 2: Verify consensus compliance
    if !pow.checkSubBlockCompliance(subBlock) {
        log.Printf("[Error] Sub-block #%d failed consensus compliance checks.", subBlock.Index)
        return false
    }
    log.Printf("[Success] Sub-block #%d complies with consensus rules.", subBlock.Index)

    // Step 3: Validate transactions within the sub-block
    for _, tx := range subBlock.Transactions {
        if !pow.validateTransaction(tx) {
            log.Printf("[Error] Transaction %s in sub-block #%d failed validation.", tx.TransactionID, subBlock.Index)
            return false
        }
    }
    log.Printf("[Success] All transactions in sub-block #%d validated successfully.", subBlock.Index)

    return true
}



// isValidHash checks if a hash meets the current difficulty requirements.
func (pow *PoW) isValidHash(hash string) bool {
    requiredPrefix := generateDifficultyPrefix(pow.State.Difficulty)
    return len(hash) >= len(requiredPrefix) && hash[:len(requiredPrefix)] == requiredPrefix
}


