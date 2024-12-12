package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// SynnergyConsensus manages the combination of PoH, PoS, and PoW for sub-block validation
type SynnergyConsensus struct {
	PoH            *PoH           // Proof of History mechanism
	PoS            *PoS           // Proof of Stake mechanism
	PoW            *PoW           // Proof of Work mechanism
	RewardManager  *RewardManager // Reward Manager for PoH, PoS, and PoW rewards
	LedgerInstance *ledger.Ledger // Ledger instance for tracking rewards, blocks, and transactions
	SubBlockCount  int            // Keeps track of sub-blocks processed for the current block
	Validators     []Validator    // List of validators in the system
	Encryption     *Encryption    // Encryption system for proof data
	
	mu             sync.Mutex     // Mutex for concurrency handling

}

// NewSynnergyConsensus initializes the full Synnergy Consensus with PoH, PoS, PoW, and the Reward Manager
func NewSynnergyConsensus(poh *PoH, pos *PoS, pow *PoW, rewardManager *RewardManager, ledgerInstance *ledger.Ledger) *SynnergyConsensus {
	return &SynnergyConsensus{
		PoH:            poh,
		PoS:            pos,
		PoW:            pow,
		RewardManager:  rewardManager,
		LedgerInstance: ledgerInstance,
		SubBlockCount:  0,
	}
}

// ProcessSingleTransaction handles a single transaction with SynnergyConsensus for processing and validation.
func ProcessSingleTransaction(sc *SynnergyConsensus, transactionID string, data string, isCrossChain bool, crossChainData CrossChainTransaction) error {
    if isCrossChain {
        // Process cross-chain transaction
        sc.ProcessTransactions([]Transaction{}, []CrossChainTransaction{crossChainData})
    } else {
        // Process standard transaction
        transaction := Transaction{
            TransactionID: transactionID,
            EncryptedData: data,
            Status:        "Pending",
        }
        sc.ProcessTransactions([]Transaction{transaction}, nil)
    }
    return nil
}


func (sc *SynnergyConsensus) ProcessTransactions() {
    fmt.Println("[Info] Starting transaction processing...")

    var wg sync.WaitGroup

    for {
        // Step 1: Collect transactions from the transaction pool
        transactions := sc.TransactionPool.FetchTransactions()
        crossChainTransactions := sc.CrossChainTransactionPool.FetchTransactions()

        // Step 2: Convert CrossChainTransactions to standard Transactions
        for _, crossChainTx := range crossChainTransactions {
            transactions = append(transactions, ConvertCrossChainToTransaction(crossChainTx))
        }

        if len(transactions) == 0 {
            time.Sleep(100 * time.Millisecond) // No transactions, wait before retrying
            continue
        }

        fmt.Printf("[Info] Collected %d transactions for processing...\n", len(transactions))

        // Step 3: Generate a PoH proof
        pohProof := sc.PoH.GeneratePoHProof(transactions)
        fmt.Printf("[Info] Generated PoH proof with hash: %s\n", pohProof.Hash)

        // Step 4: Create or fetch a new sub-block for the transactions
        subBlock := sc.createOrFetchSubBlock(transactions, pohProof.Hash)
        fmt.Printf("[Info] Created or updated sub-block #%d with %d transactions.\n", subBlock.Index, len(subBlock.Transactions))

        // Step 5: Concurrently validate the sub-block using PoH and PoS
        wg.Add(2)

        go func() {
            defer wg.Done()
            fmt.Println("[Info] Starting PoH validation...")
            if sc.PoH.ValidatePoHProof(pohProof, pohProof.Hash) {
                fmt.Println("[Success] PoH validation completed successfully.")
            } else {
                fmt.Println("[Error] PoH validation failed.")
            }
        }()

        go func() {
            defer wg.Done()
            fmt.Println("[Info] Starting PoS validation...")
            validator, err := sc.PoS.SelectValidator()
            if err != nil {
                fmt.Printf("[Error] Failed to select validator: %v\n", err)
                return
            }

            fmt.Printf("[Info] Validator selected: %s\n", validator.Address)
            sc.RewardManager.DistributePoSRewards(validator)

            if sc.PoS.ValidateSubBlock(subBlock) {
                fmt.Printf("[Success] Sub-block #%d validated by validator %s.\n", subBlock.Index, validator.Address)

                // Add sub-block to ledger
                ledgerSubBlock := ledger.SubBlock{
                    Index:        subBlock.Index,
                    Timestamp:    subBlock.Timestamp,
                    Transactions: convertTransactionsToLedger(subBlock.Transactions),
                    Validator:    validator.Address,
                    PrevHash:     subBlock.PrevHash,
                    Hash:         subBlock.Hash,
                }
                sc.LedgerInstance.BlockchainConsensusCoinLedger.AddSubBlock(ledgerSubBlock)
                fmt.Printf("[Info] Sub-block #%d added to ledger.\n", subBlock.Index)

                // Finalize block if sub-block count threshold is reached
                sc.SubBlockCount++
                if sc.SubBlockCount >= 1000 {
                    fmt.Println("[Info] Finalizing block after 1000 sub-blocks...")
                    sc.FinalizeBlock()
                }
            } else {
                fmt.Printf("[Error] Sub-block #%d validation failed.\n", subBlock.Index)
            }
        }()

        // Step 6: Wait for concurrent validations to complete
        wg.Wait()
        fmt.Println("[Info] Sub-block processing completed. Checking transaction pool for more transactions...")

        // Check transaction pool again after processing
        time.Sleep(100 * time.Millisecond) // Loop every 0.1 seconds
    }
}


func (sc *SynnergyConsensus) FinalizeBlock() {
    log.Printf("[Info] Starting block finalization process...")

    // Step 1: Retrieve all validated sub-blocks from the ledger
    subBlocks := sc.LedgerInstance.BlockchainConsensusCoinLedger.GetFinalizedSubBlocks()
    if len(subBlocks) == 0 {
        log.Println("[Error] No validated sub-blocks available for finalization.")
        return
    }
    log.Printf("[Info] Retrieved %d validated sub-blocks for finalization.\n", len(subBlocks))

    // Step 2: Retrieve the last block from the ledger
    prevBlock, err := sc.LedgerInstance.BlockchainConsensusCoinLedger.GetLastBlock()
    if err != nil {
        log.Printf("[Error] Failed to retrieve the last block: %v\n", err)
        return
    }
    log.Printf("[Info] Retrieved previous block #%d with hash: %s\n", prevBlock.Index, prevBlock.Hash)

    // Step 3: Convert the sub-blocks to blockchain-compatible format
    convertedSubBlocks := convertLedgerToBlockchainSubBlocks(subBlocks)

    // Step 4: Prepare the new block for mining
    blockToMine := &Block{
        Index:     prevBlock.Index + 1,
        Timestamp: time.Now(),
        SubBlocks: convertedSubBlocks,
        PrevHash:  prevBlock.Hash,
    }

    log.Printf("[Info] Prepared new block #%d for mining with %d sub-blocks.\n", blockToMine.Index, len(blockToMine.SubBlocks))

    // Step 5: Perform Proof of Work (PoW) mining
    log.Println("[Info] Initiating Proof of Work (PoW) mining process...")
    if err := sc.PoW.MineBlock(blockToMine); err != nil {
        log.Printf("[Error] Mining process failed for block #%d: %v\n", blockToMine.Index, err)
        return
    }
    log.Printf("[Success] Block #%d mined successfully with hash: %s\n", blockToMine.Index, blockToMine.Hash)

    // Step 6: Convert the mined block to ledger format
    ledgerBlock := convertBlockToLedgerBlock(*blockToMine)

    // Step 7: Add the mined block to the ledger
    log.Println("[Info] Adding mined block to ledger...")
    if err := sc.LedgerInstance.AddBlock(ledgerBlock); err != nil {
        log.Printf("[Error] Failed to add block #%d to ledger: %v\n", blockToMine.Index, err)
        return
    }
    log.Printf("[Success] Block #%d successfully added to ledger.\n", blockToMine.Index)

    // Step 8: Distribute PoW rewards to the miner
    log.Printf("[Info] Distributing Proof of Work rewards to miner: %s\n", sc.PoW.State.MinerAddress)
    if err := sc.RewardManager.DistributePoWRewards(sc.PoW.State.MinerAddress); err != nil {
        log.Printf("[Error] Failed to distribute rewards to miner: %v\n", err)
        return
    }
    log.Printf("[Success] Rewards distributed to miner: %s\n", sc.PoW.State.MinerAddress)

    // Step 9: Reset the sub-block count and clean up
    sc.SubBlockCount = 0
    sc.LedgerInstance.BlockchainConsensusCoinLedger.ClearFinalizedSubBlocks()
    log.Printf("[Info] Sub-block count reset and finalized sub-blocks cleared.\n")

    log.Printf("[Success] Block #%d finalized and recorded successfully.\n", blockToMine.Index)
}




func (sc *SynnergyConsensus) createSubBlock(transactions []Transaction, pohHash string) SubBlock {
    // Step 1: Retrieve current sub-block count
    subBlockIndex := sc.LedgerInstance.BlockchainConsensusCoinLedger.GetSubBlockCount()

    // Step 2: Ensure the PoH hash is valid
    if pohHash == "" {
        log.Printf("[Error] PoH hash is empty while creating sub-block #%d.\n", subBlockIndex)
        return SubBlock{}
    }

    // Step 3: Prepare the sub-block with metadata
    subBlock := SubBlock{
        Index:        subBlockIndex,
        Timestamp:    time.Now(),
        Transactions: transactions,
        Validator:    "", // Will be assigned during PoS validation
        PrevHash:     pohHash,
        Hash:         "", // Final hash computed after PoS validation
    }

    log.Printf("[Info] Created sub-block #%d with %d transactions.\n", subBlock.Index, len(transactions))
    return subBlock
}


// convertBlockToLedgerBlock converts a blockchain.Block to a ledger.Block
func convertBlockToLedgerBlock(block Block) ledger.Block {
    return ledger.Block{
        Index:     block.Index,
        Timestamp: block.Timestamp,
        SubBlocks: convertBlockchainSubBlocksToLedger(block.SubBlocks),
        PrevHash:  block.PrevHash,
        Hash:      block.Hash,
    }
}

// convertBlockchainSubBlocksToLedger converts blockchain.SubBlocks to ledger.SubBlocks
func convertBlockchainSubBlocksToLedger(subBlocks []SubBlock) []ledger.SubBlock {
    converted := make([]ledger.SubBlock, len(subBlocks))
    for i, sb := range subBlocks {
        converted[i] = ledger.SubBlock{
            Index:        sb.Index,
            Timestamp:    sb.Timestamp,
            Transactions: convertBlockchainTransactionsToLedger(sb.Transactions),
            Validator:    sb.Validator,
            PrevHash:     sb.PrevHash,
            Hash:         sb.Hash,
        }
    }
    return converted
}

// convertBlockchainTransactionsToLedger converts blockchain transactions to ledger transactions
func convertBlockchainTransactionsToLedger(transactions []Transaction) []ledger.Transaction {
    ledgerTransactions := make([]ledger.Transaction, len(transactions))
    for i, tx := range transactions {
        ledgerTransactions[i] = ledger.Transaction{
            TransactionID:     tx.TransactionID,
            FromAddress:       tx.FromAddress,
            ToAddress:         tx.ToAddress,
            Amount:            tx.Amount,
            Fee:               tx.Fee,
            TokenStandard:     tx.TokenStandard,
            TokenID:           tx.TokenID,
            Timestamp:         tx.Timestamp,
            SubBlockID:        tx.SubBlockID,
            BlockID:           tx.BlockID,
            ValidatorID:       tx.ValidatorID,
            Signature:         tx.Signature,
            Status:            tx.Status,
            EncryptedData:     tx.EncryptedData,
            DecryptedData:     tx.DecryptedData,
            ExecutionResult:   tx.ExecutionResult,
            FrozenAmount:      tx.FrozenAmount,
            RefundAmount:      tx.RefundAmount,
            ReversalRequested: tx.ReversalRequested,
        }
    }
    return ledgerTransactions
}


func (sc *SynnergyConsensus) ValidateChain() bool {
	log.Println("[Info] Starting blockchain validation...")

	// Step 1: Fetch all blocks from the ledger
	blocks := sc.LedgerInstance.BlockchainConsensusCoinLedger.GetBlocks()
	if len(blocks) == 0 {
		log.Println("[Error] No blocks available for validation.")
		return false
	}

	var wg sync.WaitGroup // WaitGroup to synchronize validation goroutines
	validationResults := make(chan bool) // Channel to collect validation results
	isValid := true                      // Overall blockchain validity flag

	// Helper function to validate a single sub-block
	validateSubBlock := func(ledgerSubBlock ledger.SubBlock, blockIndex int, method string) {
		defer wg.Done() // Decrement WaitGroup counter

		// Convert ledger.SubBlock to blockchain.SubBlock
		subBlock := convertLedgerSubBlockToBlockchainSubBlock(ledgerSubBlock)

		// Validate using PoS or PoH based on the method
		var validationPassed bool
		if method == "PoS" {
			validationPassed = sc.PoS.ValidateSubBlock(subBlock)
		} else if method == "PoH" {
			validationPassed = sc.PoH.ValidatePoHProof(subBlock.PoHProof, subBlock.Validator)
		}

		// Log results and send to channel
		if validationPassed {
			log.Printf("[Success] Sub-block #%d in block #%d validated successfully using %s.\n", subBlock.Index, blockIndex, method)
			validationResults <- true
		} else {
			log.Printf("[Error] Sub-block #%d in block #%d failed validation using %s.\n", subBlock.Index, blockIndex, method)
			validationResults <- false
		}
	}

	// Step 2: Validate sub-blocks concurrently using PoS or PoH
	for _, block := range blocks {
		for _, ledgerSubBlock := range block.SubBlocks {
			wg.Add(1) // Increment WaitGroup counter

			// Determine validation method (PoS for even indexes, PoH for odd indexes)
			if ledgerSubBlock.Index%2 == 0 {
				go validateSubBlock(ledgerSubBlock, block.Index, "PoS")
			} else {
				go validateSubBlock(ledgerSubBlock, block.Index, "PoH")
			}
		}
	}

	// Goroutine to close the results channel after all validations
	go func() {
		wg.Wait()                // Wait for all sub-block validations to complete
		close(validationResults) // Close channel
	}()

	// Collect validation results
	for result := range validationResults {
		if !result {
			isValid = false
		}
	}

	// Step 3: Validate each block using PoW
	for _, block := range blocks {
		// Convert ledger.Block to blockchain.Block
		blockchainBlock := convertLedgerBlockToBlockchainBlock(block)

		// Validate block using Proof of Work
		if !sc.PoW.ValidateBlock(&blockchainBlock) {
			log.Printf("[Error] Block #%d failed Proof of Work validation.\n", block.Index)
			isValid = false
		} else {
			log.Printf("[Success] Block #%d validated successfully using Proof of Work.\n", block.Index)
		}
	}

	// Step 4: Log final blockchain validation result
	if isValid {
		log.Println("[Success] Blockchain validation completed successfully. All blocks and sub-blocks are valid.")
	} else {
		log.Println("[Error] Blockchain validation failed. One or more blocks/sub-blocks are invalid.")
	}

	return isValid
}



// Helper function to convert ledger.SubBlock to blockchain.SubBlock
func convertLedgerSubBlockToBlockchainSubBlock(ledgerSubBlock ledger.SubBlock) SubBlock {
	return SubBlock{
		Index:        ledgerSubBlock.Index,
		Timestamp:    ledgerSubBlock.Timestamp, // No conversion needed, both are time.Time
		Transactions: convertLedgerTransactionsToBlockchain(ledgerSubBlock.Transactions),
		Validator:    ledgerSubBlock.Validator,
		PrevHash:     ledgerSubBlock.PrevHash,
		Hash:         ledgerSubBlock.Hash,
		PoHProof:     convertLedgerPoHProofToBlockchain(ledgerSubBlock.PoHProof), // Convert PoHProof
	}
}

// Helper function to convert ledger.PoHProof to blockchain.PoHProof
func convertLedgerPoHProofToBlockchain(ledgerPoHProof ledger.PoHProof) PoHProof {
	return PoHProof{
		Sequence:  ledgerPoHProof.Sequence,  // Copy the sequence number
		Timestamp: ledgerPoHProof.Timestamp, // Copy the timestamp
		Hash:      ledgerPoHProof.Hash,      // Copy the hash value
	}
}

// Helper function to convert ledger.Block to blockchain.Block
func convertLedgerBlockToBlockchainBlock(ledgerBlock ledger.Block) Block {
	return Block{
		Index:     ledgerBlock.Index,
		Timestamp: ledgerBlock.Timestamp, // Assuming both are of type time.Time
		SubBlocks: convertLedgerSubBlocksToBlockchain(ledgerBlock.SubBlocks),
		PrevHash:  ledgerBlock.PrevHash,
		Hash:      ledgerBlock.Hash,
	}
}

// Helper function to convert ledger.SubBlocks to blockchain.SubBlocks
func convertLedgerSubBlocksToBlockchain(subBlocks []ledger.SubBlock) []SubBlock {
	converted := make([]SubBlock, len(subBlocks))
	for i, sb := range subBlocks {
		converted[i] = convertLedgerSubBlockToBlockchainSubBlock(sb)
	}
	return converted
}

// computeTransactionHash computes the hash of a transaction using its relevant fields.
func computeTransactionHash(tx Transaction) string {
	hashInput := fmt.Sprintf("%s:%s:%s:%.2f:%s",
		tx.TransactionID,
		tx.FromAddress,
		tx.ToAddress,
		tx.Amount,
		tx.Timestamp.String(),
	)

	hash := sha256.New()
	hash.Write([]byte(hashInput))
	return hex.EncodeToString(hash.Sum(nil))
}

// extractTransactionHashes extracts the hashes of a list of transactions.
func extractTransactionHashes(transactions []Transaction) []string {
	var hashes []string
	for _, tx := range transactions {
		hashes = append(hashes, computeTransactionHash(tx)) // Compute hash for each transaction
	}
	return hashes
}

// Converts a slice of ledger.SubBlock to blockchain.SubBlock
func convertLedgerToBlockchainSubBlocks(ledgerSubBlocks []ledger.SubBlock) []SubBlock {
	var blockchainSubBlocks []SubBlock
	for _, ls := range ledgerSubBlocks {
		blockchainSubBlocks = append(blockchainSubBlocks, SubBlock{
			Index:        ls.Index,
			Timestamp:    ls.Timestamp,                                           // No need for conversion, both are time.Time
			Transactions: convertLedgerTransactionsToBlockchain(ls.Transactions), // Convert ledger transactions properly
			Validator:    ls.Validator,
			PrevHash:     ls.PrevHash,
			Hash:         ls.Hash,
		})
	}
	return blockchainSubBlocks
}

// Converts transaction hashes (strings) to blockchain.Transaction
func convertTransactionHashesToTransactions(hashes []string) []Transaction {
	var transactions []Transaction
	for _, hash := range hashes {
		transactions = append(transactions, Transaction{TransactionID: hash}) // Use TransactionID
	}
	return transactions
}

// ValidateTransactionSyntax validates a transaction before it is processed in the blockchain.
func (sc *SynnergyConsensus) ValidateTransactionSyntax(transaction *Transaction) (bool, error) {
	// Basic validation check
	if transaction.TransactionID == "" || transaction.Amount < 0 || transaction.FromAddress == "" || transaction.ToAddress == "" {
		return false, fmt.Errorf("invalid transaction syntax")
	}
	return true, nil
}

func (sc *SynnergyConsensus) ValidateZKProof(proofID string) error {
	log.Printf("[Info] Starting validation for ZK Proof ID: %s", proofID)

	// Step 1: Retrieve the zk-proof from the ledger using the proofID
	zkProof, err := sc.LedgerInstance.CryptographyLedger.GetZkProofByID(proofID)
	if err != nil {
		log.Printf("[Error] ZK Proof not found for proofID: %s. Error: %v", proofID, err)
		return fmt.Errorf("zk-proof not found: %w", err)
	}

	// Step 2: Verify the zk-proof data
	proofDataBytes := []byte(zkProof.ProofData)
	if !sc.verifyZKProofData(proofDataBytes) {
		log.Printf("[Error] ZK Proof validation failed for proofID: %s", proofID)
		return fmt.Errorf("zk-proof validation failed for proofID: %s", proofID)
	}
	log.Printf("[Success] ZK Proof data verified successfully for proofID: %s", proofID)

	// Step 3: Run consensus algorithm to validate the proof with validators
	validators := sc.getValidatorsForConsensus()
	consensusReached := sc.validateZKProofConsensus(zkProof, validators)
	if !consensusReached {
		log.Printf("[Error] Consensus could not be reached for ZK Proof ID: %s", proofID)
		return fmt.Errorf("consensus could not be reached for zk-proof ID: %s", proofID)
	}
	log.Printf("[Success] Consensus reached for ZK Proof ID: %s", proofID)

	// Step 4: Encrypt the proof data
	encryptedProofData, err := sc.Encryption.EncryptData(proofID, proofDataBytes, generateRandomIV())
	if err != nil {
		log.Printf("[Error] Encryption failed for ZK Proof ID: %s. Error: %v", proofID, err)
		return fmt.Errorf("encryption failed for zk-proof ID: %s: %w", proofID, err)
	}
	log.Printf("[Info] Proof data encrypted successfully for proofID: %s", proofID)

	// Step 5: Log the successful zk-proof validation in the ledger
	zkProof.ProofStatus = "validated"
	zkProof.GeneratedAt = time.Now()
	zkProof.ProofData = string(encryptedProofData)

	err = sc.LedgerInstance.CryptographyLedger.RecordZKProofValidation(zkProof.ProofID, "synnergy_validator", true)
	if err != nil {
		log.Printf("[Error] Failed to log ZK Proof validation in ledger for proofID: %s. Error: %v", proofID, err)
		return fmt.Errorf("failed to log zk-proof validation in ledger: %w", err)
	}
	log.Printf("[Success] ZK Proof ID: %s validated and recorded successfully.", proofID)

	return nil
}

func (sc *SynnergyConsensus) validateZKProofConsensus(zkProof ledger.ZKProof, validators []Validator) bool {
	log.Printf("[Info] Initiating consensus validation for ZK Proof ID: %s", zkProof.ProofID)
	consensusReached := true

	for _, validator := range validators {
		if !validator.AgreeOnZkProof(zkProof.ProofID, zkProof.ProofData, 0, 0) {
			log.Printf("[Error] Validator %s did not agree on ZK Proof ID: %s", validator.Address, zkProof.ProofID)
			consensusReached = false
			break
		}
		log.Printf("[Info] Validator %s agreed on ZK Proof ID: %s", validator.Address, zkProof.ProofID)
	}

	return consensusReached
}


func (sc *SynnergyConsensus) verifyZKProofData(proofData []byte) bool {
	log.Printf("[Info] Starting verification for ZK Proof data.")

	// Step 1: Input validation
	if len(proofData) == 0 {
		log.Printf("[Error] ZK Proof data is empty.")
		return false
	}

	// Step 2: Parse and decode proof data
	decodedData, err := decodeZKProofData(proofData)
	if err != nil {
		log.Printf("[Error] Failed to decode ZK Proof data: %v", err)
		return false
	}
	log.Printf("[Info] Decoded ZK Proof data successfully: %x", decodedData)

	// Step 3: Perform cryptographic checks
	isValid := performCryptographicVerification(decodedData)
	if !isValid {
		log.Printf("[Error] Cryptographic verification failed for ZK Proof data.")
		return false
	}
	log.Printf("[Success] Cryptographic verification succeeded for ZK Proof data.")

	// Step 4: Validate proof metadata (e.g., timestamp, format, compliance)
	if !validateZKProofMetadata(decodedData) {
		log.Printf("[Error] Metadata validation failed for ZK Proof data.")
		return false
	}
	log.Printf("[Success] Metadata validation succeeded for ZK Proof data.")

	// Step 5: Log and return success
	log.Printf("[Success] ZK Proof data verification completed successfully.")
	return true
}

func decodeZKProofData(proofData []byte) ([]byte, error) {
	log.Printf("[Info] Decoding ZK Proof data...")
	decodedData, err := base64.StdEncoding.DecodeString(string(proofData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode proof data: %w", err)
	}
	return decodedData, nil
}


func (v *Validator) AgreeOnZkProof(proofID, dataHash string, spaceUsed uint64, timeStored time.Duration) bool {
	log.Printf("[Info] Validator %s is verifying ZK Proof ID: %s", v.Address, proofID)

	// Perform simulated validation; replace with cryptographic verification in production.
	isValid := strings.HasPrefix(dataHash, "valid")
	if isValid {
		log.Printf("[Success] Validator %s agreed on ZK Proof ID: %s", v.Address, proofID)
	} else {
		log.Printf("[Error] Validator %s rejected ZK Proof ID: %s", v.Address, proofID)
	}
	return isValid
}

func generateRandomIV() []byte {
	iv := make([]byte, 16) // Generate a 16-byte initialization vector
	_, err := rand.Read(iv)
	if err != nil {
		log.Fatalf("[Critical Error] Failed to generate random IV: %v", err)
	}
	return iv
}




// ValidateSpaceTimeProof validates a space-time proof with full consensus, security, and logging.
func (sc *SynnergyConsensus) ValidateSpaceTimeProof(proofID, dataHash string, spaceUsed uint64, timeStored time.Duration) error {
	log.Printf("[Info] Starting validation for Space-Time Proof ID: %s", proofID)

	// Step 1: Input validation
	if proofID == "" || dataHash == "" {
		return fmt.Errorf("invalid proof ID or data hash: proofID='%s', dataHash='%s'", proofID, dataHash)
	}
	if spaceUsed <= 0 || timeStored <= 0 {
		return fmt.Errorf("invalid space or time parameters: spaceUsed=%d, timeStored=%s", spaceUsed, timeStored)
	}

	// Step 2: Validate space usage
	const (
		minSpaceUsed = 1024                      // Minimum: 1 KB
		maxSpaceUsed = 1024 * 1024 * 1024 * 100  // Maximum: 100 GB
	)
	if spaceUsed < minSpaceUsed || spaceUsed > maxSpaceUsed {
		return fmt.Errorf("space used (%d bytes) is out of bounds (min: %d bytes, max: %d bytes)", spaceUsed, minSpaceUsed, maxSpaceUsed)
	}
	log.Printf("[Info] Space usage validation passed: %d bytes", spaceUsed)

	// Step 3: Validate time stored
	const (
		minTimeStored = 24 * time.Hour       // Minimum: 1 day
		maxTimeStored = 365 * 24 * time.Hour // Maximum: 1 year
	)
	if timeStored < minTimeStored || timeStored > maxTimeStored {
		return fmt.Errorf("time stored (%s) is out of bounds (min: %s, max: %s)", timeStored, minTimeStored, maxTimeStored)
	}
	log.Printf("[Info] Time stored validation passed: %s", timeStored)

	// Step 4: Verify data hash integrity
	if !sc.verifyDataHash(dataHash) {
		return fmt.Errorf("data hash validation failed: %s", dataHash)
	}
	log.Printf("[Info] Data hash integrity verified: %s", dataHash)

	// Step 5: Retrieve validators and perform consensus
	log.Printf("[Info] Retrieving validators for consensus...")
	validators := sc.getValidatorsForConsensus()
	consensusReached := sc.performConsensus(validators, proofID, dataHash, spaceUsed, timeStored)
	if !consensusReached {
		return fmt.Errorf("consensus could not be reached for proofID: %s", proofID)
	}
	log.Printf("[Info] Consensus reached for proofID: %s", proofID)

	// Step 6: Log the proof validation in the ledger
	log.Printf("[Info] Recording proof validation in the ledger for proofID: %s", proofID)
	err := sc.LedgerInstance.StorageLedger.RecordProofValidation(proofID, "synnergy_validator")
	if err != nil {
		return fmt.Errorf("failed to log proof validation in ledger for proofID %s: %v", proofID, err)
	}

	// Step 7: Encrypt proof data
	log.Printf("[Info] Encrypting proof data for proofID: %s", proofID)
	encryptedProofData, err := sc.Encryption.EncryptData(proofID, []byte(dataHash), []byte("randomIV"))
	if err != nil {
		return fmt.Errorf("failed to encrypt proof data for proofID %s: %v", proofID, err)
	}
	log.Printf("[Success] Proof data encrypted successfully for proofID: %s", proofID)

	// Step 8: Finalize and return success
	log.Printf("[Success] Space-Time Proof validated successfully: proofID=%s, spaceUsed=%d bytes, timeStored=%s", proofID, spaceUsed, timeStored)
	return nil
}


func (sc *SynnergyConsensus) getValidatorsForConsensus() ([]Validator, error) {
	log.Println("[Info] Retrieving validators for consensus...")

	// Step 1: Fetch the active validators from the ledger
	validators, err := sc.LedgerInstance.ConsensusLedger.GetActiveValidators()
	if err != nil {
		log.Printf("[Error] Failed to retrieve validators from the ledger: %v", err)
		return nil, fmt.Errorf("failed to retrieve validators: %w", err)
	}

	// Step 2: Validate the integrity and status of the retrieved validators
	validValidatorList := make([]Validator, 0)
	for _, validator := range validators {
		if validator.ReputationScore >= 95.0 { // Example reputation threshold
			validValidatorList = append(validValidatorList, validator)
		} else {
			log.Printf("[Warning] Validator %s excluded due to low reputation score: %.2f", validator.Address, validator.ReputationScore)
		}
	}

	// Step 3: Ensure a minimum quorum of validators is available
	if len(validValidatorList) < 3 { // Example quorum requirement
		log.Println("[Error] Insufficient number of valid validators for consensus.")
		return nil, fmt.Errorf("insufficient validators for consensus")
	}

	log.Printf("[Info] Successfully retrieved %d validators for consensus.", len(validValidatorList))
	return validValidatorList, nil
}

// AgreeOnSpaceTimeProof simulates a validator agreeing on proof validity based on space, time, and data hash.
func (v *Validator) AgreeOnSpaceTimeProof(proofID, dataHash string, spaceUsed uint64, timeStored time.Duration) bool {
	log.Printf("[Info] Validator %s is validating proofID %s...", v.Address, proofID)

	// Step 1: Validate the integrity of the data hash
	if !isValidDataHash(dataHash) {
		log.Printf("[Error] Validator %s rejected proofID %s due to invalid data hash: %s.", v.Address, proofID, dataHash)
		return false
	}

	// Step 2: Validate space usage within bounds
	const (
		minSpaceUsed = 1024                      // Minimum 1 KB
		maxSpaceUsed = 1024 * 1024 * 1024 * 10   // Maximum 10 GB
	)
	if spaceUsed < minSpaceUsed || spaceUsed > maxSpaceUsed {
		log.Printf("[Error] Validator %s rejected proofID %s due to space used (%d bytes) out of bounds.", v.Address, proofID, spaceUsed)
		return false
	}

	// Step 3: Validate time stored within acceptable limits
	const (
		minTimeStored = 24 * time.Hour       // Minimum 1 day
		maxTimeStored = 365 * 24 * time.Hour // Maximum 1 year
	)
	if timeStored < minTimeStored || timeStored > maxTimeStored {
		log.Printf("[Error] Validator %s rejected proofID %s due to time stored (%s) out of bounds.", v.Address, proofID, timeStored)
		return false
	}

	// Step 4: Perform additional cryptographic validations
	if !v.performCryptographicCheck(proofID, dataHash, spaceUsed, timeStored) {
		log.Printf("[Error] Validator %s rejected proofID %s after cryptographic validation failure.", v.Address, proofID)
		return false
	}

	// Step 5: Log agreement and return success
	log.Printf("[Success] Validator %s agreed on proofID %s with space used: %d bytes and time stored: %s.", v.Address, proofID, spaceUsed, timeStored)
	return true
}


// isValidDataHash simulates a cryptographic validation of the data hash.
func isValidDataHash(dataHash string) bool {
	log.Printf("[Info] Validating data hash: %s", dataHash)
	return strings.HasPrefix(dataHash, "valid") 
}


// performCryptographicCheck simulates a more detailed validation using cryptographic algorithms.
func (v *Validator) performCryptographicCheck(proofID, dataHash string, spaceUsed uint64, timeStored time.Duration) bool {
	log.Printf("[Info] Performing cryptographic validation for proofID %s...", proofID)

	// EReplace with real cryptographic checks (e.g., ZK-SNARKs, Merkle proofs)
	isValid := strings.HasPrefix(dataHash, "valid") && spaceUsed > 0 && timeStored > 0
	if isValid {
		log.Printf("[Success] Cryptographic validation passed for proofID %s.", proofID)
	} else {
		log.Printf("[Error] Cryptographic validation failed for proofID %s.", proofID)
	}
	return isValid
}

// RevalidateSpaceTimeProof checks the space-time proof to ensure consistency over time.
func (sc *SynnergyConsensus) RevalidateSpaceTimeProof(proofID, dataHash string, spaceUsed uint64, timeStored time.Duration) error {
	log.Printf("[Info] Revalidating Space-Time Proof: proofID=%s", proofID)

	// Step 1: Validate input parameters
	if dataHash == "" || spaceUsed <= 0 || timeStored <= 0 {
		return fmt.Errorf("invalid parameters: dataHash='%s', spaceUsed=%d, timeStored=%s", dataHash, spaceUsed, timeStored)
	}

	// Step 2: Validate space used
	const (
		minSpaceUsed = 1024                      // Minimum 1 KB
		maxSpaceUsed = 1024 * 1024 * 1024 * 10   // Maximum 10 GB
	)
	if spaceUsed < minSpaceUsed || spaceUsed > maxSpaceUsed {
		return fmt.Errorf("space used (%d bytes) is out of bounds [%d, %d]", spaceUsed, minSpaceUsed, maxSpaceUsed)
	}

	// Step 3: Validate time stored
	const (
		minTimeStored = 24 * time.Hour       // Minimum 1 day
		maxTimeStored = 365 * 24 * time.Hour // Maximum 1 year
	)
	if timeStored < minTimeStored || timeStored > maxTimeStored {
		return fmt.Errorf("time stored (%s) is out of bounds [%s, %s]", timeStored, minTimeStored, maxTimeStored)
	}

	// Step 4: Verify data hash integrity
	if !strings.HasPrefix(dataHash, "valid") { // Placeholder for actual cryptographic validation
		return fmt.Errorf("data hash validation failed: %s", dataHash)
	}

	// Step 5: Consensus validation
	validators := sc.getValidatorsForConsensus()
	consensusReached := true
	for _, validator := range validators {
		if !validator.AgreeOnSpaceTimeProof(proofID, dataHash, spaceUsed, timeStored) {
			log.Printf("[Error] Validator %s disagreed on proofID %s", validator.Address, proofID)
			consensusReached = false
		}
	}
	if !consensusReached {
		return fmt.Errorf("consensus not reached for proofID: %s", proofID)
	}

	// Step 6: Log revalidation result
	log.Printf("[Success] ProofID %s revalidated with spaceUsed=%d bytes and timeStored=%s", proofID, spaceUsed, timeStored)
	return nil
}


// ValidateSubBlocks validates multiple sub-blocks using both PoS and PoH mechanisms.
func (sc *SynnergyConsensus) ValidateSubBlocks(subBlocks []SubBlock) bool {
	log.Printf("[Info] Validating %d sub-blocks...", len(subBlocks))

	var wg sync.WaitGroup
	subBlockCount := len(subBlocks)
	maxConcurrentValidations := 1000 // Limit concurrent validations
	validationResults := make(chan bool, subBlockCount)

	for i, subBlock := range subBlocks {
		if i >= maxConcurrentValidations {
			break
		}

		wg.Add(1)
		go func(sb SubBlock) {
			defer wg.Done()
			if sc.ShouldUsePoS(sb) {
				log.Printf("[Info] Validating sub-block %d using PoS...", sb.Index)
				if sc.PoS.ValidateSubBlock(sb) {
					log.Printf("[Success] Sub-block %d validated successfully using PoS.", sb.Index)
					validationResults <- true
				} else {
					log.Printf("[Error] Sub-block %d validation failed using PoS.", sb.Index)
					validationResults <- false
				}
			} else {
				log.Printf("[Info] Validating sub-block %d using PoH...", sb.Index)
				pohProof := sc.PoH.GeneratePoHProof()
				if sc.PoH.ValidatePoHProof(pohProof, sb.Validator) {
					log.Printf("[Success] Sub-block %d validated successfully using PoH.", sb.Index)
					validationResults <- true
				} else {
					log.Printf("[Error] Sub-block %d validation failed using PoH.", sb.Index)
					validationResults <- false
				}
			}
		}(subBlock)
	}

	wg.Wait()
	close(validationResults)

	allValid := true
	for result := range validationResults {
		if !result {
			allValid = false
		}
	}

	log.Printf("[Info] Validation completed. All sub-blocks valid: %t", allValid)
	return allValid
}



// shouldUsePoS determines whether to use PoS or PoH for sub-block validation.
func (sc *SynnergyConsensus) ShouldUsePoS(subBlock SubBlock) bool {
    sc.mu.Lock() // Lock to ensure thread-safe access to shared resources
    defer sc.mu.Unlock()

    // Step 1: Check validator availability
    validatorCount := len(sc.Validators)
    if validatorCount == 0 {
        log.Printf("[Warning] No PoS validators available. Defaulting to PoH for sub-block %d.", subBlock.Index)
        return false // Fall back to PoH if no validators are available
    }

    // Step 2: Evaluate PoS availability based on validator count
    const minValidatorsRequired = 5 // Threshold for PoS availability
    posAvailable := validatorCount >= minValidatorsRequired

    // Step 3: Load balancing logic to avoid overloading PoS
    if sc.SubBlockCount > 500 && posAvailable {
        log.Printf("[Info] PoS validator overload detected. Defaulting to PoH for sub-block %d.", subBlock.Index)
        return false // Use PoH if PoS validators are overloaded
    }

    // Step 4: Alternate between PoS and PoH for even and odd sub-block indices
    if subBlock.Index%2 == 0 && posAvailable {
        log.Printf("[Info] Sub-block %d is even-indexed. Using PoS.", subBlock.Index)
        return true // Use PoS for even-indexed sub-blocks
    }

    // Step 5: Default to PoH for all other scenarios
    log.Printf("[Info] Defaulting to PoH for sub-block %d.", subBlock.Index)
    return false
}


// VerifySignature verifies the transaction signature for a given operation under Synnergy Consensus.
func (sc *SynnergyConsensus) VerifySignature(ledgerInstance *ledger.Ledger, operationID, signerID string, signature []byte) (bool, error) {
	// Retrieve the public key associated with the signerID from the ledger instance
	publicKey, err := ledgerInstance.AuthorizationLedger.GetPublicKey(signerID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve public key for signer %s: %v", signerID, err)
	}

	// Create a transaction hash from the operationID to verify the signature
	transactionHash := sha256.Sum256([]byte(operationID))

	// Split the signature into its components, R and S
	r, s, err := parseSignature(signature) // Ensure r, s are *big.Int types
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %v", err)
	}

	// Perform ECDSA signature verification
	verified := ecdsa.Verify(publicKey, transactionHash[:], r, s)
	if !verified {
		return false, errors.New("signature verification failed under consensus")
	}

	return true, nil
}

// parseSignature splits the byte array signature into two big integers, R and S.
func parseSignature(signature []byte) (*big.Int, *big.Int, error) {
	if len(signature) != 64 {
		return nil, nil, errors.New("invalid signature length")
	}

	// Separate the signature into R and S components
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	return r, s, nil
}

// Add this CrossChainTransaction to regular Transaction conversions
func ConvertCrossChainToTransaction(crossChainTx CrossChainTransaction) Transaction {
	return Transaction{
		TransactionID: crossChainTx.TransactionID,
		FromAddress:   crossChainTx.FromAddress,
		ToAddress:     crossChainTx.ToAddress,
		Amount:        crossChainTx.Amount,
		TokenStandard: crossChainTx.TokenSymbol,
		Timestamp:     crossChainTx.Timestamp,
		SubBlockID:    0, // To be assigned later
		BlockID:       0, // To be assigned later
		Status:        crossChainTx.Status,
		EncryptedData: crossChainTx.Data, // Optional metadata
	}
}


func (tp *TransactionPool) FetchTransactions() []Transaction {
    tp.mu.Lock()
    defer tp.mu.Unlock()

    transactions := make([]Transaction, 0, len(tp.Transactions))
    for _, tx := range tp.Transactions {
        transactions = append(transactions, tx)
    }
    tp.Transactions = []Transaction{} // Clear the pool after fetching
    return transactions
}

func (sc *SynnergyConsensus) createOrFetchSubBlock(transactions []Transaction, pohHash string) SubBlock {
    sc.SubBlockLock.Lock()
    defer sc.SubBlockLock.Unlock()

    if len(sc.CurrentSubBlock.Transactions) == 0 {
        // Create a new sub-block
        sc.CurrentSubBlock = SubBlock{
            Index:        sc.SubBlockCount + 1,
            Timestamp:    time.Now(),
            Transactions: transactions,
            PrevHash:     sc.LastSubBlockHash,
            Hash:         pohHash,
        }
    } else {
        // Add transactions to the current sub-block
        sc.CurrentSubBlock.Transactions = append(sc.CurrentSubBlock.Transactions, transactions...)
        sc.CurrentSubBlock.Hash = pohHash // Update hash
    }

    return sc.CurrentSubBlock
}
