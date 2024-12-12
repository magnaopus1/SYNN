package ledger

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"
)

// IncrementRetryCount increments the retry count for a given transaction ID.
func (l *BlockchainConsensusCoinLedger) IncrementRetryCount(txID string) error {
	l.TransactionRetryData[txID]++
	return nil
}

// SetSubblockCapacity sets the maximum capacity for sub-blocks.
func (l *BlockchainConsensusCoinLedger) SetSubblockCapacity(capacity int) error {
	l.SubblockCapacity = capacity
	return nil
}

// GetSubblockCapacity retrieves the current sub-block capacity.
func (l *BlockchainConsensusCoinLedger) GetSubblockCapacity() (int, error) {
	return l.SubblockCapacity, nil
}

// GetSubblockCapacityData retrieves the current capacity usage of all sub-blocks.
func (l *BlockchainConsensusCoinLedger) GetSubblockCapacityData() map[string]int {
	return l.SubblockCapacities
}

// AuditSubblockCapacity audits the capacity usage across sub-blocks.
func (l *BlockchainConsensusCoinLedger) AuditSubblockCapacity() error {
	for id, usage := range l.SubblockCapacities {
		if usage > 100 { // assuming 100 is the maximum allowed capacity per sub-block
			return fmt.Errorf("subblock %s exceeds capacity: %d", id, usage)
		}
	}
	return nil
}

// LogSubblockCapacity logs the capacity for a specific sub-block.
func (l *BlockchainConsensusCoinLedger) LogSubblockCapacity(subblockID string, capacity int) error {
	l.SubblockCapacities[subblockID] = capacity
	l.SubblockCapacityHistory[subblockID] = append(l.SubblockCapacityHistory[subblockID], capacity)
	return nil
}

// GetSubblockCapacityHistory retrieves the historical capacity usage of a specified sub-block.
func (l *BlockchainConsensusCoinLedger) GetSubblockCapacityHistory(subblockID string) ([]int, error) {
	history, exists := l.SubblockCapacityHistory[subblockID]
	if !exists {
		return nil, fmt.Errorf("no capacity history found for subblock %s", subblockID)
	}
	return history, nil
}

// SetSynthronCoinDenomination sets the denomination for Synthron Coin.
func (l *BlockchainConsensusCoinLedger) SetSynthronCoinDenomination(denomination string) error {
	l.SynthronCoinDenomination = denomination
	return nil
}

// GetSynthronCoinDenomination retrieves the current Synthron Coin denomination.
func (l *BlockchainConsensusCoinLedger) GetSynthronCoinDenomination() (string, error) {
	if l.SynthronCoinDenomination == "" {
		return "", fmt.Errorf("Synthron Coin denomination not set")
	}
	return l.SynthronCoinDenomination, nil
}

// AuditCoinDenominations audits the denomination settings.
func (l *BlockchainConsensusCoinLedger) AuditCoinDenominations() error {
	if l.SynthronCoinDenomination == "" {
		return fmt.Errorf("no denomination set for Synthron Coin")
	}
	if len(l.CoinDenominationHistory) == 0 {
		return fmt.Errorf("no history found for Synthron Coin denomination changes")
	}
	return nil
}

// TrackCoinDenominationChange logs denomination changes with a timestamp.
func (l *BlockchainConsensusCoinLedger) TrackCoinDenominationChange(change string) error {
	date := time.Now().Format("2006-01-02")
	l.CoinDenominationHistory[date] = change
	l.SynthronCoinDenomination = change
	return nil
}

// GetCoinDenominationData retrieves historical denomination changes.
func (l *BlockchainConsensusCoinLedger) GetCoinDenominationData() map[string]string {
	return l.CoinDenominationHistory
}

// SetSubblockCacheLimit sets the cache limit for sub-blocks.
func (l *BlockchainConsensusCoinLedger) SetSubblockCacheLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("invalid cache limit; must be positive")
	}
	l.SubblockCacheLimit = limit
	return nil
}

// GetSubblockCacheLimit retrieves the current sub-block cache limit.
func (l *BlockchainConsensusCoinLedger) GetSubblockCacheLimit() (int, error) {
	l.Lock()
	defer l.Unlock()

	if l.SubblockCacheLimit == 0 {
		return 0, fmt.Errorf("subblock cache limit not set")
	}
	return l.SubblockCacheLimit, nil
}

// EnableBlockCompression enables block compression.
func (l *BlockchainConsensusCoinLedger) EnableBlockCompression() error {
	l.Lock()
	defer l.Unlock()
	l.BlockCompressionEnabled = true
	return nil
}

// DisableBlockCompression disables block compression.
func (l *BlockchainConsensusCoinLedger) DisableBlockCompression() error {
	l.Lock()
	defer l.Unlock()
	l.BlockCompressionEnabled = false
	return nil
}

// FetchCompressionStatus retrieves the current block compression status.
func (l *BlockchainConsensusCoinLedger) FetchCompressionStatus() (bool, error) {
	l.Lock()
	defer l.Unlock()
	return l.BlockCompressionEnabled, nil
}

// SetBlockCompressionLevel sets the compression level for blocks.
func (l *BlockchainConsensusCoinLedger) SetBlockCompressionLevel(level int) error {
	l.Lock()
	defer l.Unlock()
	l.BlockCompressionLevel = level
	return nil
}

// GetBlockCompressionLevel retrieves the current block compression level.
func (l *BlockchainConsensusCoinLedger) GetBlockCompressionLevel() (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.BlockCompressionLevel, nil
}

// EnableEncryption enables encryption for blocks.
func (l *BlockchainConsensusCoinLedger) EnableEncryption() error {
	l.Lock()
	defer l.Unlock()
	l.EncryptionEnabled = true
	return nil
}

// DisableEncryption disables encryption for blocks.
func (l *BlockchainConsensusCoinLedger) DisableEncryption() error {
	l.Lock()
	defer l.Unlock()
	l.EncryptionEnabled = false
	return nil
}

// SetEncryptionKey sets the encryption key.
func (l *BlockchainConsensusCoinLedger) SetEncryptionKey(key string) error {
	l.Lock()
	defer l.Unlock()

	if len(key) != 32 {
		return errors.New("encryption key must be 32 bytes long")
	}
	l.EncryptionKey = key
	return nil
}

// GetEncryptionKey retrieves the current encryption key.
func (l *BlockchainConsensusCoinLedger) GetEncryptionKey() (string, error) {
	l.Lock()
	defer l.Unlock()

	if l.EncryptionKey == "" {
		return "", errors.New("no encryption key set")
	}
	return l.EncryptionKey, nil
}

// VerifyEncryptionStatus checks if encryption is enabled.
func (l *BlockchainConsensusCoinLedger) VerifyEncryptionStatus() (bool, error) {
	l.Lock()
	defer l.Unlock()
	return l.EncryptionEnabled, nil
}

// EncryptBlock encrypts the given block data using AES.
func (l *BlockchainConsensusCoinLedger) EncryptBlock(blockData []byte) ([]byte, error) {
	l.Lock()
	defer l.Unlock()

	if len(l.EncryptionKey) != 32 {
		return nil, errors.New("invalid encryption key")
	}
	block, err := aes.NewCipher([]byte(l.EncryptionKey))
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	encryptedData := aesGCM.Seal(nonce, nonce, blockData, nil)
	return encryptedData, nil
}

// DecryptBlock decrypts the given encrypted data using AES.
func (l *BlockchainConsensusCoinLedger) DecryptBlock(encryptedData []byte) ([]byte, error) {
	l.Lock()
	defer l.Unlock()

	if len(l.EncryptionKey) != 32 {
		return nil, errors.New("invalid encryption key")
	}
	block, err := aes.NewCipher([]byte(l.EncryptionKey))
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	decryptedData, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	return decryptedData, err
}

// EnableSubblockCompression enables compression for sub-blocks.
func (l *BlockchainConsensusCoinLedger) EnableSubblockCompression() error {
	l.Lock()
	defer l.Unlock()
	l.SubblockCompressionEnabled = true
	return nil
}

// DisableSubblockCompression disables compression for sub-blocks.
func (l *BlockchainConsensusCoinLedger) DisableSubblockCompression() error {
	l.Lock()
	defer l.Unlock()
	l.SubblockCompressionEnabled = false
	return nil
}

// SetSubblockValidationCriteria sets the validation criteria.
func (l *BlockchainConsensusCoinLedger) SetSubblockValidationCriteria(criteria string) error {
	l.Lock()
	defer l.Unlock()
	l.SubblockValidationCriteria = criteria
	return nil
}

// GetSubblockValidationCriteria retrieves the validation criteria.
func (l *BlockchainConsensusCoinLedger) GetSubblockValidationCriteria() (string, error) {
	l.Lock()
	defer l.Unlock()
	return l.SubblockValidationCriteria, nil
}

// SetValidationInterval sets the interval for validation.
func (l *BlockchainConsensusCoinLedger) SetValidationInterval(interval int) error {
	l.Lock()
	defer l.Unlock()
	l.ValidationInterval = interval
	return nil
}

// GetValidationInterval retrieves the validation interval.
func (l *BlockchainConsensusCoinLedger) GetValidationInterval() (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.ValidationInterval, nil
}

// SetTransactionLimit sets the transaction limit for blocks.
func (l *BlockchainConsensusCoinLedger) SetTransactionLimit(limit int) error {
	l.Lock()
	defer l.Unlock()
	l.BlockTransactionLimit = limit
	return nil
}

// GetTransactionLimit retrieves the transaction limit for blocks.
func (l *BlockchainConsensusCoinLedger) GetTransactionLimit() (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.BlockTransactionLimit, nil
}

// EnableTransactionTracking enables tracking of transactions.
func (l *BlockchainConsensusCoinLedger) EnableTransactionTracking() error {
	l.Lock()
	defer l.Unlock()
	l.TransactionTrackingEnabled = true
	return nil
}

// DisableTransactionTracking disables tracking of transactions.
func (l *BlockchainConsensusCoinLedger) DisableTransactionTracking() error {
	l.Lock()
	defer l.Unlock()
	l.TransactionTrackingEnabled = false
	return nil
}

// GetBlocks returns a slice of all blocks in the ledger.
func (l *BlockchainConsensusCoinLedger) GetBlocks() []Block {
	l.Lock() // Ensure thread-safety with locking
	defer l.Unlock()
	return l.Blocks
}

// GetBlockByIndex retrieves a block from the ledger by its index
func (l *BlockchainConsensusCoinLedger) GetBlockByIndex(index int) (*Block, error) {
	l.Lock() // Ensure thread-safety
	defer l.Unlock()

	// Check if index is within bounds
	if index < 0 || index >= len(l.Blocks) {
		return nil, fmt.Errorf("block with index %d not found", index)
	}

	return &l.Blocks[index], nil
}

// GetBlockCount returns the total number of blocks in the ledger.
func (l *BlockchainConsensusCoinLedger) GetBlockCount() int {
	l.Lock()
	defer l.Unlock()
	return len(l.Blocks)
}

// GetSubBlockByID retrieves a sub-block by its ID from the ledger.
func (l *BlockchainConsensusCoinLedger) GetSubBlockByID(subBlockID string) (SubBlock, error) {
	l.Lock()
	defer l.Unlock()

	for _, subBlock := range l.SubBlocks {
		if subBlock.SubBlockID == subBlockID {
			return subBlock, nil
		}
	}
	return SubBlock{}, fmt.Errorf("sub-block with ID %s not found", subBlockID)
}

// GetLastSubBlock retrieves the most recent sub-block from the ledger.
func (l *BlockchainConsensusCoinLedger) GetLastSubBlock() (*SubBlock, error) {
	l.Lock()
	defer l.Unlock()

	if len(l.SubBlocks) == 0 {
		return nil, fmt.Errorf("no sub-blocks found")
	}
	return &l.SubBlocks[len(l.SubBlocks)-1], nil
}

// LogSubBlock logs a finalized sub-block in the ledger.
func (l *BlockchainConsensusCoinLedger) LogSubBlock(subBlock *SubBlock) error {
	fmt.Printf("Sub-block %d logged at %s with hash %s.\n", subBlock.Index, subBlock.Timestamp, subBlock.Hash)
	return nil
}

// GetSubBlockCount returns the number of sub-blocks in the ledger.
func (l *BlockchainConsensusCoinLedger) GetSubBlockCount() int {
	return len(l.SubBlocks)
}

// AddSubBlock adds a validated sub-block to the ledger.
func (l *BlockchainConsensusCoinLedger) AddSubBlock(subBlock SubBlock) error {
	l.Lock()
	defer l.Unlock()

	// Validate the sub-block
	err := l.validateSubBlock(subBlock)
	if err != nil {
		l.RejectedTransactions = append(l.RejectedTransactions, subBlock.Transactions...)
		return fmt.Errorf("sub-block validation failed: %v", err)
	}

	// Add to the list of sub-blocks
	l.SubBlocks = append(l.SubBlocks, subBlock)

	// If 1000 sub-blocks are reached, create a new block
	if len(l.SubBlocks) == 1000 {
		newBlock, err := l.createBlockFromSubBlocks()
		if err != nil {
			return fmt.Errorf("failed to create block from sub-blocks: %v", err)
		}

		err = l.AddBlock(newBlock)
		if err != nil {
			return fmt.Errorf("failed to add block: %v", err)
		}

		// Reset sub-blocks after a block is created
		l.SubBlocks = []SubBlock{}
	}

	return nil
}

// GetPreviousBlockHash retrieves the hash of the block immediately preceding the provided block.
func (ledger *BlockchainConsensusCoinLedger) GetPreviousBlockHash(currentBlockID string) (string, error) {
	// Retrieve the current block from the ledger.
	currentBlock, err := ledger.GetBlockByID(currentBlockID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve current block: %v", err)
	}

	// Check if this is the first block (genesis block) which has no previous block.
	if currentBlock.PrevHash == "" {
		return "", fmt.Errorf("no previous block found (this might be the genesis block)")
	}

	// Return the hash of the previous block.
	return currentBlock.PrevHash, nil
}

// GetCurrentDifficulty retrieves the current mining difficulty for the blockchain.
func (ledger *BlockchainConsensusCoinLedger) GetCurrentDifficulty() (int, error) {
	// Get the latest block from the ledger.
	latestBlock, err := ledger.GetLatestBlock()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve the latest block: %v", err)
	}

	// Return the difficulty level of the latest block.
	// Assuming the block has a 'Difficulty' field.
	return latestBlock.Difficulty, nil
}

// GetBlockByID retrieves a block by its ID from the ledger.
func (l *BlockchainConsensusCoinLedger) GetBlockByID(blockID string) (Block, error) {
	for _, block := range l.GetBlocks() { // Use the getter method
		if block.BlockID == blockID {
			return block, nil
		}
	}
	return Block{}, fmt.Errorf("block with ID %s not found", blockID)
}

// GetLatestBlock retrieves the most recent block from the ledger.
func (l *BlockchainConsensusCoinLedger) GetLatestBlock() (Block, error) {
	blocks := l.GetBlocks() // Use the getter method
	if len(blocks) == 0 {
		return Block{}, fmt.Errorf("no blocks found in the ledger")
	}

	var latestBlock Block

	for _, block := range blocks {
		if block.Timestamp.After(latestBlock.Timestamp) {
			latestBlock = block
		}
	}

	return latestBlock, nil
}

// AddBlock adds a finalized block to the ledger.
func (l *BlockchainConsensusCoinLedger) AddBlock(block Block) error {
	l.Lock()
	defer l.Unlock()

	// Validate the block before adding it
	expectedIndex := l.BlockchainConsensusCoinState.BlockHeight
	if block.Index != expectedIndex {
		return fmt.Errorf("block validation failed: block index does not match expected value %d", expectedIndex)
	}

	// Block is valid, so add it to the finalized blocks
	l.FinalizedBlocks = append(l.FinalizedBlocks, block)
	l.BlockchainConsensusCoinState.BlockHeight++ // Increment block height
	l.BlockchainConsensusCoinState.LastBlockHash = block.Hash

	// Notify any registered listeners of the new block
	l.notifyBlockListeners(block)
	return nil
}

// notifyBlockListeners iterates over all registered listeners and calls them with the new block.
func (l *BlockchainConsensusCoinLedger) notifyBlockListeners(block Block) {
	for _, listener := range l.BlockListeners {
		listener(block)
	}
}

// RegisterBlockListener registers a listener for new blocks.
func (l *BlockchainConsensusCoinLedger) RegisterBlockListener(listener func(Block)) {
	l.Lock()
	defer l.Unlock()
	l.BlockListeners = append(l.BlockListeners, listener)
}

func (l *BlockchainConsensusCoinLedger) ValidateBlock(block Block) error {
	// Initialize the block index if this is the first block added
	if len(l.FinalizedBlocks) == 0 {
		l.BlockIndex = -1 // Set to -1 so the first block added has index 0
	}

	// Check if the block follows the previous one in terms of indexing
	expectedIndex := l.BlockIndex + 1
	if block.Index != expectedIndex {
		return fmt.Errorf("block index does not match expected value: got %d, expected %d", block.Index, expectedIndex)
	}

	// Check that the previous block hash matches
	if expectedIndex > 0 { // Not the genesis block
		prevBlock := l.FinalizedBlocks[len(l.FinalizedBlocks)-1]
		if block.PrevHash != prevBlock.Hash {
			return errors.New("block previous hash does not match the hash of the last finalized block")
		}
	}

	// Validate transactions in the block
	for _, subBlock := range block.SubBlocks {
		if err := l.validateSubBlock(subBlock); err != nil {
			return fmt.Errorf("block validation failed due to invalid sub-block: %v", err)
		}
	}

	// Verify the block hash
	calculatedHash := l.CalculateBlockHash(block)
	if calculatedHash != block.Hash {
		return errors.New("block hash does not match")
	}

	// If all checks pass, increment block index and return nil
	l.BlockIndex++
	fmt.Printf("Block validated successfully with index %d.\n", block.Index) // Debug statement
	return nil
}

// createBlockFromSubBlocks consolidates 1000 sub-blocks into a finalized block.
func (l *BlockchainConsensusCoinLedger) createBlockFromSubBlocks() (Block, error) {
	// Ensure that there are exactly 1000 sub-blocks
	if len(l.SubBlocks) != 1000 {
		return Block{}, errors.New("sub-block count does not match expected 1000")
	}

	// Create a new block from sub-blocks
	newBlock := Block{
		BlockID:   fmt.Sprintf("block_%d", l.BlockIndex+1),
		Index:     l.BlockIndex + 1,
		SubBlocks: l.SubBlocks,
		Timestamp: time.Now(),
	}

	// Calculate the hash for the new block
	newBlock.Hash = l.CalculateBlockHash(newBlock)

	return newBlock, nil
}

// validateSubBlock securely validates a sub-block.
func (l *BlockchainConsensusCoinLedger) validateSubBlock(subBlock SubBlock) error {
	l.Lock()
	defer l.Unlock()

	// 1. Validate Sub-Block Integrity: Ensure correct structure and fields
	if subBlock.Index <= 0 {
		return fmt.Errorf("invalid sub-block index: %d", subBlock.Index)
	}

	// Validate sub-block timestamp
	if subBlock.Timestamp.IsZero() || subBlock.Timestamp.After(time.Now()) {
		return fmt.Errorf("invalid sub-block timestamp")
	}

	// 2. Check for unique transactions (avoid double-spending)
	for _, tx := range subBlock.Transactions {
		if _, exists := l.TransactionCache[tx.TransactionID]; exists {
			return fmt.Errorf("transaction %s already exists in ledger", tx.TransactionID)
		}

		// Optional: Validate transaction integrity (could include checking signatures, formats, etc.)
		if err := l.ValidateTransaction(tx.TransactionID); err != nil {
			return fmt.Errorf("invalid transaction %s: %v", tx.TransactionID, err)
		}
	}

	// 3. Add transactions to ledger state (assuming they passed validation)
	for _, tx := range subBlock.Transactions {
		l.TransactionCache[tx.TransactionID] = tx

		// Convert Transaction to TransactionRecord and add to TransactionHistory
		transactionRecord := TransactionRecord{
			From:       tx.FromAddress,
			To:         tx.ToAddress,
			Amount:     tx.Amount,
			Fee:        tx.Fee,
			Hash:       tx.TransactionID, // Using TransactionID as a unique identifier (hash)
			BlockIndex: subBlock.Index,   // Using the sub-block index as BlockIndex
		}
		l.BlockchainConsensusCoinState.TransactionHistory = append(l.BlockchainConsensusCoinState.TransactionHistory, transactionRecord)
	}

	// 4. Ensure sub-block links properly to a valid previous sub-block (consistency)
	if len(l.SubBlocks) > 0 {
		lastSubBlock := l.SubBlocks[len(l.SubBlocks)-1]
		if subBlock.PrevHash != lastSubBlock.Hash {
			return fmt.Errorf("sub-block chain is broken. Expected previous hash: %s, got: %s", lastSubBlock.Hash, subBlock.PrevHash)
		}
	}

	return nil
}

// calculateBlockHash generates a SHA-256 hash for the block.
func (l *BlockchainConsensusCoinLedger) CalculateBlockHash(block Block) string {
	hashData := fmt.Sprintf("%d-%v-%s", block.Index, block.SubBlocks, block.Timestamp)
	hash := sha256.Sum256([]byte(hashData))
	return hex.EncodeToString(hash[:])
}

// updateMerkleRoot recalculates the Merkle root for the ledger state.
func (l *BlockchainConsensusCoinLedger) updateMerkleRoot() {
	// Generate a Merkle root from transaction history (for simplicity, concatenating hashes)
	var hashData string
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		hashData += tx.Hash // Use the correct field for transaction ID, which is Hash in TransactionRecord
	}

	// Calculate the Merkle root by hashing the concatenated transaction hashes
	hash := sha256.Sum256([]byte(hashData))
	l.BlockchainConsensusCoinState.MerkleRoot = hex.EncodeToString(hash[:])
}

// GetBlockByHash retrieves a block by its hash by iterating through the slice
func (l *BlockchainConsensusCoinLedger) GetBlockByHash(hash string) *Block {
	for _, block := range l.Blocks {
		if block.Hash == hash {
			return &block // Return a pointer to the found block
		}
	}
	return nil // Return nil if no block is found with the given hash
}

// ReplaceChain replaces the current chain with a new chain of blocks
func (l *BlockchainConsensusCoinLedger) ReplaceChain(newChain []Block) {
	// Only replace the chain if the new chain is longer
	if len(newChain) > len(l.Blocks) {
		l.Blocks = newChain
		fmt.Println("Chain replaced with the longer chain.")
	} else {
		fmt.Println("Chain replacement aborted: the new chain is not longer.")
	}
}

// GetLatestBlockHash fetches the latest block hash from the ledger
func (ledger *BlockchainConsensusCoinLedger) GetLatestBlockHash() string {
	ledger.Lock()
	defer ledger.Unlock()

	if len(ledger.Blocks) == 0 {
		fmt.Println("Ledger is empty, no blocks available.")
		return "" // No blocks available
	}

	// Get the hash of the latest block
	latestBlock := ledger.Blocks[len(ledger.Blocks)-1]
	fmt.Printf("Fetched latest block hash: %s (block timestamp: %s).\n", latestBlock.Hash, latestBlock.Timestamp)

	return latestBlock.Hash
}

// ValidateSubBlock validates the structure and integrity of a sub-block.
func (l *BlockchainConsensusCoinLedger) ValidateSubBlock(subBlockID string) error {
	l.Lock()
	defer l.Unlock()

	// Find the sub-block by ID
	var foundSubBlock *SubBlock
	for i, subBlock := range l.SubBlocks {
		if subBlock.SubBlockID == subBlockID {
			foundSubBlock = &l.SubBlocks[i]
			break
		}
	}

	if foundSubBlock == nil {
		return errors.New("sub-block not found")
	}

	// Perform validation on the sub-block
	foundSubBlock.Status = "validated"
	return nil
}

// RecordConfirmedBlock logs a block confirmation event.
func (l *BlockchainConsensusCoinLedger) RecordConfirmedBlock(blockID string) error {
	l.Lock()
	defer l.Unlock()

	// Find the block by ID
	var foundBlock *Block
	for i, block := range l.FinalizedBlocks {
		if block.BlockID == blockID {
			foundBlock = &l.FinalizedBlocks[i]
			break
		}
	}

	if foundBlock == nil {
		return errors.New("block not found")
	}

	// Update the block status to "confirmed"
	foundBlock.Status = "confirmed"
	return nil
}

// GetBlockStatus returns the status of a block by its ID.
func (l *BlockchainConsensusCoinLedger) GetBlockStatus(blockID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Loop through the finalizedBlocks slice to find the block by blockID
	for _, block := range l.FinalizedBlocks {
		if block.BlockID == blockID {
			return block.Status, nil
		}
	}

	return "", errors.New("block not found")
}

// GetSubBlockStatus returns the status of a sub-block by its ID.
func (l *BlockchainConsensusCoinLedger) GetSubBlockStatus(subBlockID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	// Loop through the subBlocks slice to find the sub-block by subBlockID
	for _, subBlock := range l.SubBlocks {
		if subBlock.SubBlockID == subBlockID {
			return subBlock.Status, nil
		}
	}

	return "", errors.New("sub-block not found")
}

// SetSystemBalance sets the system's global balance (SynthronBalance) for faucet deposits or other special purposes.
func (l *BlockchainConsensusCoinLedger) SetSystemBalance(amount float64) {
	l.Lock()
	defer l.Unlock()

	l.SynthronBalance = amount
	fmt.Printf("Synthron balance updated to: %.2f\n", l.SynthronBalance)
}

// GetSystemBalance retrieves the SynthronBalance for faucet deposits or other special purposes.
func (l *BlockchainConsensusCoinLedger) GetSystemBalance() float64 {
	l.Lock()
	defer l.Unlock()

	return l.SynthronBalance
}

// RecordFaucetClaim allows an account to claim from the SynthronBalance.
func (l *BlockchainConsensusCoinLedger) RecordFaucetClaim(accountsLedger *AccountsWalletLedger, accountID string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the SynthronBalance has enough funds
	if l.SynthronBalance < amount {
		return fmt.Errorf("insufficient funds in the Synthron faucet to claim %.2f", amount)
	}

	// Check if the account exists in AccountsWalletLedger
	accountsLedger.lock.Lock()
	account, exists := accountsLedger.AccountsWalletLedgerState.Accounts[accountID]
	if !exists {
		accountsLedger.lock.Unlock()
		return fmt.Errorf("account %s does not exist", accountID)
	}

	// Update the SynthronBalance and the account's balance
	l.SynthronBalance -= amount
	account.Balance += amount
	accountsLedger.AccountsWalletLedgerState.Accounts[accountID] = account
	accountsLedger.lock.Unlock()

	// Record the faucet claim in the transaction history
	transaction := TransactionRecord{
		Hash:   generateTransactionID(),
		From:   "SynthronBalance",
		To:     accountID,
		Amount: amount,
	}

	l.BlockchainConsensusCoinState.TransactionHistory = append(l.BlockchainConsensusCoinState.TransactionHistory, transaction)

	fmt.Printf("Faucet claim of %.2f by account %s recorded successfully.\n", amount, accountID)
	return nil
}

// RecordTransaction records a general transaction between accounts.
func (l *BlockchainConsensusCoinLedger) RecordTransaction(accountsLedger *AccountsWalletLedger, from string, to string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	// Ensure both accounts exist in AccountsWalletLedger
	accountsLedger.lock.Lock()
	fromAccount, fromExists := accountsLedger.AccountsWalletLedgerState.Accounts[from]
	toAccount, toExists := accountsLedger.AccountsWalletLedgerState.Accounts[to]

	if !fromExists {
		accountsLedger.lock.Unlock()
		return fmt.Errorf("sender account %s does not exist", from)
	}
	if !toExists {
		accountsLedger.lock.Unlock()
		return fmt.Errorf("receiver account %s does not exist", to)
	}

	// Ensure sufficient balance in the sender's account
	if fromAccount.Balance < amount {
		accountsLedger.lock.Unlock()
		return fmt.Errorf("insufficient funds in sender account %s", from)
	}

	// Update the balances
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	accountsLedger.AccountsWalletLedgerState.Accounts[from] = fromAccount
	accountsLedger.AccountsWalletLedgerState.Accounts[to] = toAccount
	accountsLedger.lock.Unlock()

	// Record the transaction in the transaction history
	transaction := TransactionRecord{
		Hash:   generateTransactionID(),
		From:   from,
		To:     to,
		Amount: amount,
	}

	l.BlockchainConsensusCoinState.TransactionHistory = append(l.BlockchainConsensusCoinState.TransactionHistory, transaction)

	fmt.Printf("Transaction of %.2f from %s to %s recorded successfully.\n", amount, from, to)
	return nil
}

// RecordFaucetDeposit adds a faucet deposit to the SynthronBalance.
func (l *BlockchainConsensusCoinLedger) RecordFaucetDeposit(amount float64) error {
	l.Lock()
	defer l.Unlock()

	// Ensure the deposit amount is valid
	if amount <= 0 {
		return errors.New("deposit amount must be greater than zero")
	}

	// Add the amount to the SynthronBalance
	l.SynthronBalance += amount

	// Record the deposit in the transaction history
	transaction := TransactionRecord{
		Hash:   generateTransactionID(), // Ensure Hash is generated
		From:   "faucet_system",
		To:     "SynthronBalance",
		Amount: amount,
	}
	l.BlockchainConsensusCoinState.TransactionHistory = append(l.BlockchainConsensusCoinState.TransactionHistory, transaction)

	fmt.Printf("Faucet deposit of %.2f recorded successfully.\n", amount)
	return nil
}

// RecordConsensusMechanismAddition logs the addition of a new Layer 2 consensus mechanism.
func (l *BlockchainConsensusCoinLedger) RecordConsensusMechanismAddition(mechanismID, mechanismName string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Consensus Mechanism Added: ID: %s, Name: %s", mechanismID, mechanismName)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "ConsensusMechanismAddition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Added",
	})

	fmt.Printf("Layer 2 consensus mechanism %s added successfully.\n", mechanismID)
}

// RecordConsensusTransition logs the transition from one consensus mechanism to another.
func (l *BlockchainConsensusCoinLedger) RecordConsensusTransition(fromMechanism, toMechanism string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Consensus Transition: From %s to %s", fromMechanism, toMechanism)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "ConsensusTransition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Transitioned",
	})

	fmt.Printf("Consensus transitioned from %s to %s.\n", fromMechanism, toMechanism)
}

// RecordStrategyAddition logs the addition of a new strategy to Layer 2 consensus.
func (l *BlockchainConsensusCoinLedger) RecordStrategyAddition(strategyID, strategyName string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Strategy Added: Strategy ID: %s, Name: %s", strategyID, strategyName)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "StrategyAddition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Added",
	})

	fmt.Printf("Layer 2 strategy %s added successfully.\n", strategyID)
}

// RecordStrategyHop logs a transition between two Layer 2 strategies.
func (l *BlockchainConsensusCoinLedger) RecordStrategyHop(fromStrategy, toStrategy string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Strategy Hop: From %s to %s", fromStrategy, toStrategy)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "StrategyHop",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Hopped",
	})

	fmt.Printf("Strategy transitioned from %s to %s.\n", fromStrategy, toStrategy)
}

// RecordConsensusLayerAddition logs the addition of a new consensus layer.
func (l *BlockchainConsensusCoinLedger) RecordConsensusLayerAddition(layerID, layerName string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Consensus Layer Added: Layer ID: %s, Name: %s", layerID, layerName)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "ConsensusLayerAddition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Added",
	})

	fmt.Printf("Layer 2 consensus layer %s added successfully.\n", layerID)
}

// RecordConsensusLayerTransition logs the transition between two consensus layers.
func (l *BlockchainConsensusCoinLedger) RecordConsensusLayerTransition(fromLayer, toLayer string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Layer 2 Consensus Layer Transition: From %s to %s", fromLayer, toLayer)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "ConsensusLayerTransition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Transitioned",
	})

	fmt.Printf("Consensus layer transitioned from %s to %s.\n", fromLayer, toLayer)
}

// RecordCollaborationNodeAddition logs the addition of a new collaboration node in the Layer 2 framework.
func (l *BlockchainConsensusCoinLedger) RecordCollaborationNodeAddition(nodeID, nodeName string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Collaboration Node Added: Node ID: %s, Name: %s", nodeID, nodeName)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "CollaborationNodeAddition",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Added",
	})

	fmt.Printf("Collaboration node %s added successfully.\n", nodeID)
}

// RecordCollaborationTaskAssignment logs the assignment of a task to a collaboration node.
func (l *BlockchainConsensusCoinLedger) RecordCollaborationTaskAssignment(taskID, nodeID string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Collaboration Task Assigned: Task ID: %s, Node ID: %s", taskID, nodeID)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "CollaborationTaskAssignment",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Assigned",
	})

	fmt.Printf("Task %s assigned to node %s successfully.\n", taskID, nodeID)
}

// RecordCollaborationTaskCompletion logs the completion of a task by a collaboration node.
func (l *BlockchainConsensusCoinLedger) RecordCollaborationTaskCompletion(taskID, nodeID string) {
	l.Lock()
	defer l.Unlock()

	details := fmt.Sprintf("Collaboration Task Completed: Task ID: %s, Node ID: %s", taskID, nodeID)

	l.Layer2ConsensusLogs = append(l.Layer2ConsensusLogs, Layer2ConsensusLog{
		EventType: "CollaborationTaskCompletion",
		Timestamp: time.Now(),
		Details:   details,
		Status:    "Completed",
	})

	fmt.Printf("Task %s completed by node %s successfully.\n", taskID, nodeID)
}

// RecordFeeLog logs the transaction fees for auditing purposes.
func (l *BlockchainConsensusCoinLedger) RecordFeeLog(txID string, fee uint64) error {
	l.Lock()
	defer l.Unlock()

	// Fetch transaction, log the fee
	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Convert uint64 fee to float64 to match the transaction structure
	tx.Fee = float64(fee)
	l.TransactionCache[txID] = tx
	return nil
}

// FreezeTransactionAmount locks a certain amount in a transaction temporarily.
func (l *BlockchainConsensusCoinLedger) FreezeTransactionAmount(txID string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Ensure the types of tx.Amount and amount match
	if tx.Amount < amount {
		return errors.New("insufficient transaction amount")
	}

	// Assuming tx should have FrozenAmount field, ensure it's added to the Transaction struct
	tx.FrozenAmount = amount
	l.TransactionCache[txID] = tx
	return nil
}

// GetTransactionHistoryByAddress retrieves all transactions associated with a specific address.
func (l *BlockchainConsensusCoinLedger) GetTransactionHistoryByAddress(address string) ([]*Transaction, error) {
	var transactionHistory []*Transaction

	// Iterate through all transactions to find those associated with the given address
	for _, tx := range l.TransactionCache {
		if tx.FromAddress == address || tx.ToAddress == address {
			// Append the pointer to the transaction
			transactionHistory = append(transactionHistory, &tx)
		}
	}

	if len(transactionHistory) == 0 {
		return nil, fmt.Errorf("no transactions found for address: %s", address)
	}

	return transactionHistory, nil
}

// ValidateTransaction validates the integrity and authenticity of a transaction.
func (l *BlockchainConsensusCoinLedger) ValidateTransaction(txID string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Perform validation using the IsValid method
	if err := tx.IsValid(l); err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Mark the transaction as validated
	tx.Status = "validated"
	l.TransactionCache[txID] = tx
	return nil
}

// IsValid checks the validity of the transaction.
func (tx *Transaction) IsValid(l *BlockchainConsensusCoinLedger, A *AccountsWalletLedger) error {
	// Ensure the amount is positive
	if tx.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	// Check if the sender exists in the ledger and has enough balance
	sender, senderExists := A.AccountsWalletLedgerState.Accounts[tx.FromAddress]
	if !senderExists {
		return fmt.Errorf("sender account not found")
	}

	// Check if the receiver exists in the ledger
	_, receiverExists := A.AccountsWalletLedgerState.Accounts[tx.ToAddress]
	if !receiverExists {
		return fmt.Errorf("receiver account not found")
	}

	// Ensure the sender has sufficient balance to cover the amount and fee
	totalAmount := tx.Amount + tx.Fee
	if sender.Balance < totalAmount {
		return fmt.Errorf("insufficient funds in sender's account")
	}

	// Decode the sender's public key (assuming it's stored as a PEM-encoded string)
	publicKey, err := decodePublicKey(sender.PublicKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %v", err)
	}

	// Create the message to be validated (a simple concatenation of transaction data)
	message := fmt.Sprintf("%s:%s:%f", tx.FromAddress, tx.ToAddress, tx.Amount)

	// Validate the transaction signature using the decoded ECDSA public key
	isValidSignature, err := ValidateSignature(publicKey, tx.Signature, message)
	if err != nil {
		return fmt.Errorf("error validating signature: %v", err)
	}
	if !isValidSignature {
		return fmt.Errorf("invalid transaction signature")
	}

	// Ensure the transaction timestamp is valid (not too old or too far in the future)
	now := time.Now()
	if tx.Timestamp.After(now) || tx.Timestamp.Before(now.Add(-24*time.Hour)) {
		return fmt.Errorf("transaction timestamp is invalid")
	}

	// If the transaction passed all checks, return nil (valid)
	return nil
}

// decodePublicKey decodes a PEM-encoded ECDSA public key string into an *ecdsa.PublicKey.
func decodePublicKey(pemEncodedPublicKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemEncodedPublicKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, fmt.Errorf("not an ECDSA public key")
	}
}

// ValidateSignature verifies that a given signature matches the sender's address and amount.
func ValidateSignature(publicKey *ecdsa.PublicKey, signature string, message string) (bool, error) {
	// 1. Hash the message using SHA-256
	hashedMessage := sha256.Sum256([]byte(message))

	// 2. Decode the signature (assuming it’s provided as a hex string)
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil || len(signatureBytes) != 64 {
		return false, fmt.Errorf("invalid signature format")
	}

	// Split the signature into r and s values (ECDSA signatures consist of two parts: r and s)
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// 3. Verify the signature using the sender's public key, message hash, and the r/s values
	return ecdsa.Verify(publicKey, hashedMessage[:], r, s), nil
}

// LogTransaction logs the transaction ID and signature into the ledger.
func (l *BlockchainConsensusCoinLedger) LogTransaction(txID string, signature string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the transaction exists in the cache
	tx, exists := l.TransactionCache[txID]
	if !exists {
		// Create a new transaction entry if it does not exist
		tx = Transaction{
			TransactionID: txID,
			Signature:     signature,
			// Additional fields could be added here as needed
		}
	} else {
		// Update the signature if the transaction already exists
		tx.Signature = signature
	}

	// Store or update the transaction in the cache
	l.TransactionCache[txID] = tx
	return nil
}

// UpdateCancellationRequest allows updating a transaction cancellation request.
func (l *BlockchainConsensusCoinLedger) UpdateCancellationRequest(txID string, cancel bool) error {
	l.Lock()
	defer l.Unlock()

	request, exists := l.CancellationRequests[txID]
	if !exists {
		return errors.New("cancellation request not found")
	}

	request.Status = "updated"
	l.CancellationRequests[txID] = request
	return nil
}

// GetTransactionByID fetches a transaction by its ID.
func (l *BlockchainConsensusCoinLedger) GetTransactionByID(txID string) (TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	// Fetch the transaction from the cache
	tx, exists := l.TransactionCache[txID]
	if !exists {
		return TransactionRecord{}, errors.New("transaction not found")
	}

	// Convert Transaction to TransactionRecord
	txRecord := TransactionRecord{
		From:       tx.FromAddress,
		To:         tx.ToAddress,
		Amount:     tx.Amount,
		Fee:        tx.Fee,
		Hash:       tx.TransactionID, // Assuming this is the transaction hash
		Status:     tx.Status,
		BlockIndex: 0, // Add appropriate BlockIndex logic here
	}

	// Return the converted TransactionRecord
	return txRecord, nil
}

// ReleaseFrozenFunds releases funds that were frozen in a transaction.
func (l *BlockchainConsensusCoinLedger) ReleaseFrozenFunds(txID string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	tx.FrozenAmount = 0
	l.TransactionCache[txID] = tx
	return nil
}

// RecordPoolTransaction logs a transaction related to a liquidity pool.
func (l *BlockchainConsensusCoinLedger) RecordPoolTransaction(poolID string, tx Transaction) error {
	l.Lock()
	defer l.Unlock()

	// Check if the transaction cache exists, if not initialize
	if l.TransactionCache == nil {
		l.TransactionCache = make(map[string]Transaction)
	}

	// Append the transaction to the pool-specific log
	l.TransactionCache[poolID] = tx

	return nil
}

// RecordWalletTransaction logs a transaction in a user's wallet.
func (l *BlockchainConsensusCoinLedger) RecordWalletTransaction(walletID string, tx Transaction) error {
	l.Lock()
	defer l.Unlock()

	// Check if the transaction cache exists, if not initialize
	if l.TransactionCache == nil {
		l.TransactionCache = make(map[string]Transaction)
	}

	// Append the transaction to the wallet-specific log
	l.TransactionCache[walletID] = tx

	return nil
}

// TransferFunds facilitates the transfer of funds between two accounts.
func (l *AccountsWalletLedger) TransferFunds(fromAccountID, toAccountID string, amount float64) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Fetch the accounts involved
	fromAccount, exists := l.AccountsWalletLedgerState.Accounts[fromAccountID]
	if !exists {
		return errors.New("source account not found")
	}

	toAccount, exists := l.AccountsWalletLedgerState.Accounts[toAccountID]
	if !exists {
		return errors.New("destination account not found")
	}

	// Check if the source account has enough balance
	if fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Perform the transfer by updating balances
	fromAccount.Balance -= amount
	toAccount.Balance += amount

	// Update the accounts in the ledger state
	l.AccountsWalletLedgerState.Accounts[fromAccountID] = fromAccount
	l.AccountsWalletLedgerState.Accounts[toAccountID] = toAccount

	return nil
}

// RecordTransactionFee logs the fee for a specific transaction.
func (l *BlockchainConsensusCoinLedger) RecordTransactionFee(txID string, fee uint64) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Convert fee from uint64 to float64
	tx.Fee = float64(fee)
	l.TransactionCache[txID] = tx
	return nil
}

// RefundTransactionGas refunds gas fees for a specific transaction.
func (l *BlockchainConsensusCoinLedger) RefundTransactionGas(txID string, refundAmount uint64) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// Convert refundAmount from uint64 to float64
	tx.RefundAmount = float64(refundAmount)
	l.TransactionCache[txID] = tx
	return nil
}

// RecordTransactionMetrics logs performance metrics for a transaction.
func (l *BlockchainConsensusCoinLedger) RecordTransactionMetrics(txID string, metrics string) error {
	l.Lock()
	defer l.Unlock()

	// Initialize the transaction metrics map if not already initialized
	if l.TransactionMetrics == nil {
		l.TransactionMetrics = make(map[string]string)
	}

	// Record the metrics for the transaction
	l.TransactionMetrics[txID] = metrics
	return nil
}

// AddTransaction records a new transaction between two accounts and updates their balances.
func (l *BlockchainConsensusCoinLedger) AddTransaction(A *AccountsWalletLedger, from string, to string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if both accounts exist
	fromAccount, fromExists := A.AccountsWalletLedgerState.Accounts[from]
	toAccount, toExists := A.AccountsWalletLedgerState.Accounts[to]

	if !fromExists {
		return fmt.Errorf("sender account %s does not exist", from)
	}
	if !toExists {
		return fmt.Errorf("receiver account %s does not exist", to)
	}

	// Ensure sufficient balance in the sender's account
	if fromAccount.Balance < amount {
		return fmt.Errorf("insufficient funds in sender account %s", from)
	}

	// Update balances
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	A.AccountsWalletLedgerState.Accounts[from] = fromAccount
	A.AccountsWalletLedgerState.Accounts[to] = toAccount

	// Record the transaction
	transaction := TransactionRecord{
		ID:          generateTransactionID(), // Generate a unique transaction ID
		From:        from,
		To:          to,
		Amount:      amount,
		Fee:         0.0,                                                             // Assuming no fee for now
		Status:      "confirmed",                                                     // Assuming confirmed status
		BlockIndex:  l.BlockIndex,                                                    // Assuming BlockIndex is the current ledger block index
		Timestamp:   time.Now(),                                                      // Set the current time
		BlockHeight: l.BlockchainConsensusCoinState.BlockHeight,                      // Set the current block height in the ledger
		ValidatorID: "",                                                              // Validator ID can be added if relevant
		Action:      "Transfer",                                                      // Default action for this transaction
		Details:     fmt.Sprintf("Transfer of %.2f from %s to %s", amount, from, to), // Additional details
	}

	l.BlockchainConsensusCoinState.TransactionHistory = append(l.BlockchainConsensusCoinState.TransactionHistory, transaction)

	fmt.Printf("Transaction from %s to %s for %.2f recorded successfully.\n", from, to, amount)
	return nil
}

// StoreReversalRequest stores a new reversal request in the ledger.
func (l *BlockchainConsensusCoinLedger) StoreReversalRequest(request *ReversalRequest) error {
	l.Lock()
	defer l.Unlock()

	// Check if the reversal request already exists
	if _, exists := l.ReversalRequests[request.TransactionID]; exists {
		return fmt.Errorf("reversal request already exists for transaction ID: %s", request.TransactionID)
	}

	// Store the dereferenced reversal request in the ReversalRequests map
	l.ReversalRequests[request.TransactionID] = *request
	return nil
}

// RemoveTransaction removes a transaction from the ledger.
func (l *BlockchainConsensusCoinLedger) RemoveTransaction(txID string) error {
	l.Lock()
	defer l.Unlock()

	_, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	delete(l.TransactionCache, txID)

	return nil
}

// DistributeFees distributes fees collected across validators, nodes, and system participants.
func (l *BlockchainConsensusCoinLedger) DistributeFees(A *AccountsWalletLedger, txID string, distributionMap map[string]uint64) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	// distributionMap contains participant IDs (validators/nodes) and their respective fee portions
	for participantID, feeAmount := range distributionMap {
		participant, exists := A.AccountsWalletLedgerState.Accounts[participantID]
		if !exists {
			return errors.New("participant account not found")
		}

		// Convert feeAmount from uint64 to float64 before adding it to the balance
		participant.Balance += float64(feeAmount)
		A.AccountsWalletLedgerState.Accounts[participantID] = participant
	}

	// Update transaction status to indicate that fees were distributed
	tx.Status = "fees_distributed"
	l.TransactionCache[txID] = tx

	return nil
}

// RecordTransactionExecution logs the execution of a transaction in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordTransactionExecution(txID string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	tx.Status = "executed"
	l.TransactionCache[txID] = tx
	return nil
}

// GetTransactionsByBlockHeight retrieves transactions from a specific block height.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByBlockHeight(blockHeight int) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		if tx.BlockHeight == blockHeight {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found at the specified block height")
	}
	return transactions, nil
}

// GetTransactionsByValidator retrieves all transactions validated by a specific validator.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByValidator(validatorID string) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		if tx.ValidatorID == validatorID {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found for the specified validator")
	}
	return transactions, nil
}

// GetTransactionsByTimeRange retrieves transactions within a specified time range.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByTimeRange(startTime, endTime time.Time) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		if tx.Timestamp.After(startTime) && tx.Timestamp.Before(endTime) {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found within the specified time range")
	}
	return transactions, nil
}

// GetTransactionsByCondition retrieves transactions based on a custom condition.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByCondition(conditionFunc func(TransactionRecord) bool) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		if conditionFunc(tx) {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found that match the condition")
	}
	return transactions, nil
}

// GetTransactionsByStatus retrieves transactions based on their current status.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByStatus(status string) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		if tx.Status == status {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found with the specified status")
	}
	return transactions, nil
}

// GetTransactionsByFeeRange retrieves transactions based on a range of fees.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByFeeRange(minFee, maxFee uint64) ([]TransactionRecord, error) {
	l.Lock()
	defer l.Unlock()

	var transactions []TransactionRecord
	for _, tx := range l.BlockchainConsensusCoinState.TransactionHistory {
		// Convert uint64 fee values to float64 for comparison
		if float64(tx.Fee) >= float64(minFee) && float64(tx.Fee) <= float64(maxFee) {
			transactions = append(transactions, tx)
		}
	}

	if len(transactions) == 0 {
		return nil, errors.New("no transactions found within the specified fee range")
	}
	return transactions, nil
}

// RecordEscrowTransaction logs a transaction related to an escrow.
func (l *BlockchainConsensusCoinLedger) RecordEscrowTransaction(escrowID string, tx TransactionRecord) error {
	l.Lock()
	defer l.Unlock()

	// Check if the escrow exists in the transaction log
	escrow, exists := l.EscrowTransactions[escrowID]
	if !exists {
		return fmt.Errorf("escrow transaction with ID %s not found", escrowID)
	}

	// Validate that the transaction relates to the escrow
	if escrow.SenderID != tx.From || escrow.ReceiverID != tx.To {
		return fmt.Errorf("transaction parties do not match the escrow parties")
	}

	// Update the escrow transaction status based on transaction details, if applicable
	if tx.Amount >= escrow.Amount && escrow.Status == EscrowStatusPending {
		escrow.Status = EscrowStatusReleased
		escrow.ReleaseTime = time.Now()
	}

	// Update the escrow transaction in the ledger
	l.EscrowTransactions[escrowID] = escrow

	// Since `tx` is of type `TransactionRecord`, you need to store it accordingly
	l.TransactionCache[tx.Hash] = Transaction{ // Convert TransactionRecord to Transaction
		TransactionID: tx.Hash,
		FromAddress:   tx.From,
		ToAddress:     tx.To,
		Amount:        tx.Amount,
		Fee:           tx.Fee,
		Status:        tx.Status,
	}

	return nil
}

// UpdateReversalRequest updates the reversal request for a specific transaction.
func (l *BlockchainConsensusCoinLedger) UpdateReversalRequest(txID string, requestReversal bool) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	tx.ReversalRequested = requestReversal
	l.TransactionCache[txID] = tx

	return nil
}

// GetTransactionsByBlockID retrieves transactions by the block ID.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByBlockID(blockID string) ([]*Transaction, error) {
	l.Lock()
	defer l.Unlock()

	// Iterate over the blocks to find the block by its ID
	for _, block := range l.Blocks {
		if block.BlockID == blockID {
			var transactions []*Transaction

			// Iterate through sub-blocks to gather all transactions
			for _, subBlock := range block.SubBlocks {
				// Convert []Transaction to []*Transaction
				for i := range subBlock.Transactions {
					transactions = append(transactions, &subBlock.Transactions[i])
				}
			}

			return transactions, nil
		}
	}

	return nil, fmt.Errorf("block with ID %s not found", blockID)
}

// FinalizeBlock finalizes a block and adds it to the blockchain.
func (ledger *BlockchainConsensusCoinLedger) FinalizeBlock(block Block) error {
	ledger.Lock()
	defer ledger.Unlock()

	ledger.FinalizedBlocks = append(ledger.FinalizedBlocks, block)
	ledger.BlockchainConsensusCoinState.BlockHeight++
	ledger.BlockchainConsensusCoinState.LastBlockHash = block.Hash
	fmt.Printf("Block #%d finalized with hash %s\n", block.Index, block.Hash)

	// Add the block's transactions to history
	for _, subBlock := range block.SubBlocks {
		for _, tx := range subBlock.Transactions {
			// Convert Transaction to TransactionRecord
			record := TransactionRecord{
				From:        tx.FromAddress,
				To:          tx.ToAddress,
				Amount:      tx.Amount,
				Fee:         tx.Fee,
				Hash:        tx.TransactionID,
				Status:      tx.Status,
				BlockIndex:  block.Index,
				Timestamp:   tx.Timestamp,
				BlockHeight: ledger.BlockchainConsensusCoinState.BlockHeight,
				ValidatorID: tx.ValidatorID, // Assuming ValidatorID is a part of Transaction
			}
			ledger.BlockchainConsensusCoinState.TransactionHistory = append(ledger.BlockchainConsensusCoinState.TransactionHistory, record)
		}
	}

	return nil
}

// RejectBlock rejects a block and adds it to the rejected blocks list.
func (ledger *BlockchainConsensusCoinLedger) RejectBlock(block Block, reason string) error {
	ledger.Lock()
	defer ledger.Unlock()

	ledger.RejectedBlocks = append(ledger.RejectedBlocks, block)
	fmt.Printf("Block #%d rejected: %s\n", block.Index, reason)
	return nil
}

// GetTransactionHistory returns the transaction history for the ledger.
func (ledger *BlockchainConsensusCoinLedger) GetTransactionHistory() ([]TransactionRecord, error) {
	ledger.Lock()
	defer ledger.Unlock()

	return ledger.BlockchainConsensusCoinState.TransactionHistory, nil // Return the history directly
}

// UpdateBalances updates account balances after a transaction is processed.
func (l *BlockchainConsensusCoinLedger) UpdateBalances(A *AccountsWalletLedger, txID string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	fromAccount, exists := A.AccountsWalletLedgerState.Accounts[tx.FromAddress]
	if !exists {
		return errors.New("sender account not found")
	}

	toAccount, exists := A.AccountsWalletLedgerState.Accounts[tx.ToAddress]
	if !exists {
		return errors.New("receiver account not found")
	}

	if fromAccount.Balance < tx.Amount {
		return errors.New("insufficient funds")
	}

	// Transfer the amount
	fromAccount.Balance -= tx.Amount
	toAccount.Balance += tx.Amount

	// Update ledger balances
	A.AccountsWalletLedgerState.Accounts[tx.FromAddress] = fromAccount
	A.AccountsWalletLedgerState.Accounts[tx.ToAddress] = toAccount

	return nil
}

// RecordEscrowLog records an event or action related to an escrow.
func (l *BlockchainConsensusCoinLedger) RecordEscrowLog(escrowID, logDetails string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the escrow log exists, if not initialize
	if l.EscrowLogs == nil {
		l.EscrowLogs = make(map[string][]string)
	}

	// Append the log entry to the escrow-specific log
	l.EscrowLogs[escrowID] = append(l.EscrowLogs[escrowID], logDetails)

	return nil
}

// RecordBlockMetrics records metrics for a full block in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordBlockMetrics(blockID string, metrics string, blockSize int64, transactions int, validatorID string) error {
	// Record the block metrics
	l.BlockMetrics[blockID] = BlockMetric{
		BlockID:      blockID,
		Timestamp:    time.Now(),
		Metrics:      metrics,
		BlockSize:    blockSize,
		Transactions: transactions,
		ValidatorID:  validatorID,
	}

	fmt.Printf("Block metrics recorded for Block %s at %s:\n- Metrics: %s\n- Block Size: %d bytes\n- Transactions: %d\n- Validator: %s\n",
		blockID, time.Now().String(), metrics, blockSize, transactions, validatorID)

	return nil
}

// RecordSubBlockMetrics records metrics for a sub-block in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordSubBlockMetrics(subBlockID string, metrics string, subBlockSize int64, transactions int, parentBlockID string) error {
	// Record the sub-block metrics
	l.SubBlockMetrics[subBlockID] = SubBlockMetric{
		SubBlockID:   subBlockID,
		Timestamp:    time.Now(),
		Metrics:      metrics,
		SubBlockSize: subBlockSize,
		Transactions: transactions,
		ParentBlock:  parentBlockID,
	}

	fmt.Printf("Sub-block metrics recorded for Sub-Block %s at %s:\n- Metrics: %s\n- Sub-Block Size: %d bytes\n- Transactions: %d\n- Parent Block: %s\n",
		subBlockID, time.Now().String(), metrics, subBlockSize, transactions, parentBlockID)

	return nil
}

// GetTransactionsByGasFeeRange retrieves transactions with gas fees within a specified range.
func (l *BlockchainConsensusCoinLedger) GetTransactionsByGasFeeRange(minGasFee, maxGasFee uint64) ([]TransactionRecord, error) {
	var matchingTxs []TransactionRecord

	// Iterate over all transactions in the ledger and check if their gas fee is within the range.
	for _, txRecord := range l.FinalizedTransactions {
		if txRecord.Fee >= float64(minGasFee) && txRecord.Fee <= float64(maxGasFee) {
			matchingTxs = append(matchingTxs, txRecord)
		}
	}

	if len(matchingTxs) == 0 {
		return nil, fmt.Errorf("no transactions found within the gas fee range %d - %d", minGasFee, maxGasFee)
	}

	return matchingTxs, nil
}

// UpdateTransactionStatus updates the status of a transaction in the ledger.
func (l *BlockchainConsensusCoinLedger) UpdateTransactionStatus(transactionID string, newStatus string) error {
	// Retrieve the transaction from the ledger
	txn, exists := l.FinalizedTransactions[transactionID] // Ensure the map field is correctly named
	if !exists {
		return fmt.Errorf("transaction with ID %s not found", transactionID)
	}

	// Update the transaction's status
	txn.Status = newStatus
	l.FinalizedTransactions[transactionID] = txn // Save the updated transaction

	fmt.Printf("Transaction %s status updated to %s\n", transactionID, newStatus)
	return nil
}

// GetTransactionStatus returns the status of a transaction by its ID.
func (l *BlockchainConsensusCoinLedger) GetTransactionStatus(txID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return "", errors.New("transaction not found")
	}

	return tx.Status, nil
}

// RecordTransactionBroadcast logs the broadcast of a transaction.
func (l *BlockchainConsensusCoinLedger) RecordTransactionBroadcast(txID string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	tx.Status = "broadcasted"
	l.TransactionCache[txID] = tx
	return nil
}

// GetReversalRequestByTransactionID retrieves a reversal request by transaction ID.
func (l *BlockchainConsensusCoinLedger) GetReversalRequestByTransactionID(txnID string) (ReversalRequest, error) {
	l.Lock()
	defer l.Unlock()

	// Check if the reversal request exists
	request, exists := l.ReversalRequests[txnID]
	if !exists {
		return ReversalRequest{}, fmt.Errorf("reversal request not found for transaction ID: %s", txnID)
	}

	return request, nil
}

// RecordTransactionSignature records the signature for a transaction.
func (l *BlockchainConsensusCoinLedger) RecordTransactionSignature(txID, signature string) error {
	l.Lock()
	defer l.Unlock()

	tx, exists := l.TransactionCache[txID]
	if !exists {
		return errors.New("transaction not found")
	}

	tx.Signature = signature
	l.TransactionCache[txID] = tx
	return nil
}

// GetTransaction retrieves a transaction from the transaction cache by ID.
func (l *BlockchainConsensusCoinLedger) GetTransaction(txID string) (*Transaction, error) {
	l.Lock()
	defer l.Unlock()

	transaction, exists := l.TransactionCache[txID]
	if !exists {
		return nil, errors.New("transaction not found in cache")
	}
	return &transaction, nil
}

// UpdateTransaction updates a transaction in the transaction cache.
func (l *BlockchainConsensusCoinLedger) UpdateTransaction(txID string, updatedTransaction Transaction) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.TransactionCache[txID]; !exists {
		return errors.New("transaction not found in cache")
	}
	l.TransactionCache[txID] = updatedTransaction
	return nil
}

// SetTransactionRetryLimit sets the retry limit for transactions.
func (l *BlockchainConsensusCoinLedger) SetTransactionRetryLimit(limit int) error {
	l.TransactionRetryLimit = limit
	return nil
}

// GetTransactionRetryLimit retrieves the current transaction retry limit.
func (l *BlockchainConsensusCoinLedger) GetTransactionRetryLimit() (int, error) {
	return l.TransactionRetryLimit, nil
}

// GetTransactionRetryData retrieves retry counts for all transactions.
func (l *BlockchainConsensusCoinLedger) GetTransactionRetryData() map[string]int {
	return l.TransactionRetryData
}

// AuditTransactionRetries audits and logs retry data for all transactions.
func (l *BlockchainConsensusCoinLedger) AuditTransactionRetries() error {
	fmt.Println("Auditing transaction retries:")
	for txID, retryCount := range l.TransactionRetryData {
		fmt.Printf("Transaction ID: %s, Retries: %d\n", txID, retryCount)
	}
	return nil
}

// AddPoHProof adds a Proof of History (PoH) to the consensus state and ledger.
func (l *BlockchainConsensusCoinLedger) AddPoHProof(proof PoHProof) error {
	l.Lock()
	defer l.Unlock()

	// Assuming PoHProofs is part of consensusState or similar
	l.ConsensusState.PoHProofs = append(l.ConsensusState.PoHProofs, proof)
	fmt.Printf("Added PoH Proof: %s\n", proof.Hash)
	return nil
}

// RecordValidatorSelection logs a validator selection in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordValidatorSelection(validatorAddress string, epoch int) error {
	logEntry := fmt.Sprintf("Validator %s selected in epoch %d", validatorAddress, epoch)
	return l.StoreLog(logEntry)
}

// RecordMinedBlock stores details of the mined block in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordMinedBlock(block Block, minerAddress string, reward float64) error {
	logEntry := fmt.Sprintf("Block #%d mined by %s with reward %.2f SYNN", block.Index, minerAddress, reward)
	return l.StoreLog(logEntry)
}

// UpdatePunishment updates the punishment details in the ledger
func (l *BlockchainConsensusCoinLedger) UpdatePunishment(entity string, punishment Punishment) {
	l.Lock()
	defer l.Unlock()

	l.Punishments[entity] = punishment
	fmt.Printf("Punishment updated for %s: %.2f SYNN\n", entity, punishment.Amount)
}

// UpdateValidatorReward updates the reward of a validator.
func (l *BlockchainConsensusCoinLedger) UpdateValidatorReward(validatorID string, reward float64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.ConsensusState.ValidatorRewards[validatorID]; exists {
		l.ConsensusState.ValidatorRewards[validatorID] += reward
	} else {
		l.ConsensusState.ValidatorRewards[validatorID] = reward
	}
	fmt.Printf("Updated reward for Validator ID: %s to %.2f\n", validatorID, l.ConsensusState.ValidatorRewards[validatorID])
	return nil
}

// UpdateMinerReward updates the reward of a miner.
func (l *BlockchainConsensusCoinLedger) UpdateMinerReward(minerID string, reward float64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.ConsensusState.MinerRewards[minerID]; exists {
		l.ConsensusState.MinerRewards[minerID] += reward
	} else {
		l.ConsensusState.MinerRewards[minerID] = reward
	}
	fmt.Printf("Updated reward for Miner ID: %s to %.2f\n", minerID, l.ConsensusState.MinerRewards[minerID])
	return nil
}

// GetFinalizedSubBlocks returns all finalized sub-blocks in the ledger.
func (l *BlockchainConsensusCoinLedger) GetFinalizedSubBlocks() []SubBlock {
	l.Lock()
	defer l.Unlock()

	var allSubBlocks []SubBlock

	// Iterate over finalized blocks and extract the sub-blocks
	for _, block := range l.FinalizedBlocks {
		allSubBlocks = append(allSubBlocks, block.SubBlocks...)
	}

	return allSubBlocks
}

// UpdateValidatorStake updates the stake of a specific validator.
// It increases or decreases the stake based on the amount.
func (l *BlockchainConsensusCoinLedger) UpdateValidatorStake(validatorID string, stakeAmount float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the validator already exists in the ledger
	if _, exists := l.ValidatorStakes[validatorID]; !exists {
		// If the validator doesn't exist, initialize their stake to 0
		l.ValidatorStakes[validatorID] = 0.0
	}

	// Update the validator's stake
	l.ValidatorStakes[validatorID] += stakeAmount

	// Ensure the stake doesn't go below zero
	if l.ValidatorStakes[validatorID] < 0 {
		return fmt.Errorf("validator %s has insufficient stake", validatorID)
	}

	fmt.Printf("Updated stake for validator %s: new stake %.6f\n", validatorID, l.ValidatorStakes[validatorID])
	return nil
}

// StoreEncryptedPunishment stores the encrypted punishment data in the ledger.
func (l *BlockchainConsensusCoinLedger) StoreEncryptedPunishment(entity string, encryptedData string) {
	l.Lock()
	defer l.Unlock()

	l.EncryptedPunishments[entity] = encryptedData
	fmt.Printf("Encrypted punishment data stored for entity %s.\n", entity)
}

// StoreEncryptedReward stores encrypted reward data for a given entity.
func (l *BlockchainConsensusCoinLedger) StoreEncryptedReward(entity string, encryptedReward string) {
	l.Lock()
	defer l.Unlock()

	// Initialize the EncryptedRewards map if it doesn't exist
	if l.EncryptedRewards == nil {
		l.EncryptedRewards = make(map[string]string)
	}

	// Store the encrypted reward for the given entity
	l.EncryptedRewards[entity] = encryptedReward

	fmt.Printf("Encrypted reward for entity %s stored successfully.\n", entity)
}

// UpdateParticipantReward updates the reward amount for a specific participant in the ledger.
func (l *BlockchainConsensusCoinLedger) UpdateParticipantReward(participantID string, rewardAmount float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the participant exists in the map, and update their reward
	if _, exists := l.ParticipantRewards[participantID]; !exists {
		return fmt.Errorf("participant %s does not exist", participantID)
	}

	l.ParticipantRewards[participantID] += rewardAmount
	fmt.Printf("Reward for participant %s updated by %.2f. New total: %.2f\n", participantID, rewardAmount, l.ParticipantRewards[participantID])
	return nil
}

// GetSubBlocks returns the list of sub-blocks.
func (l *BlockchainConsensusCoinLedger) GetSubBlocks() []SubBlock {
	l.Lock()
	defer l.Unlock()
	return l.SubBlocks // Correctly access the unexported subBlocks field
}

func (l *BlockchainConsensusCoinLedger) GetPunitiveMeasureLogs() ([]PunitiveMeasureRecord, error) {
	l.Lock()
	defer l.Unlock()
	return l.PunitiveMeasureLogs, nil
}

func (l *BlockchainConsensusCoinLedger) RevertPunitiveAction(actionID string) error {
	l.Lock()
	defer l.Unlock()
	for i, record := range l.PunitiveMeasureLogs {
		if record.ActionID == actionID {
			l.PunitiveMeasureLogs[i].Status = "reverted"
			return nil
		}
	}
	return fmt.Errorf("action with ID %s not found", actionID)
}

func (l *BlockchainConsensusCoinLedger) SetPunishmentReevaluationInterval(interval time.Duration) error {
	l.Lock()
	defer l.Unlock()
	l.PunishmentReevaluationInterval = interval
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPunishmentReevaluationInterval() (time.Duration, error) {
	l.Lock()
	defer l.Unlock()
	if l.PunishmentReevaluationInterval == 0 {
		return 0, fmt.Errorf("punishment reevaluation interval is not set")
	}
	return l.PunishmentReevaluationInterval, nil
}

func (l *BlockchainConsensusCoinLedger) LogPunishmentAdjustments(actionID, adjustedBy, details string) error {
	l.Lock()
	defer l.Unlock()
	adjustment := PunishmentAdjustmentLog{
		AdjustmentID: fmt.Sprintf("%s-%d", actionID, time.Now().UnixNano()),
		ActionID:     actionID,
		AdjustedBy:   adjustedBy,
		Timestamp:    time.Now(),
		Details:      details,
	}
	l.PunishmentAdjustmentLogs = append(l.PunishmentAdjustmentLogs, adjustment)
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetAdaptiveDifficulty(difficultyLevel int) error {
	l.Lock()
	defer l.Unlock()
	l.AdaptiveDifficulty = difficultyLevel
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetAdaptiveDifficulty() (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.AdaptiveDifficulty, nil
}

func (l *BlockchainConsensusCoinLedger) EnableAdaptiveRewardDistribution() error {
	l.Lock()
	defer l.Unlock()
	l.AdaptiveRewardDistributionEnabled = true
	return nil
}

func (l *BlockchainConsensusCoinLedger) DisableAdaptiveRewardDistribution() error {
	l.Lock()
	defer l.Unlock()
	l.AdaptiveRewardDistributionEnabled = false
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetDifficultyLevel(newLevel int, reason string) error {
	l.Lock()
	defer l.Unlock()
	adjustmentLog := DifficultyAdjustmentLog{
		AdjustmentID:       fmt.Sprintf("adj-%d", time.Now().UnixNano()),
		Timestamp:          time.Now(),
		NewDifficultyLevel: newLevel,
		Reason:             reason,
	}
	l.DifficultyAdjustmentLogs = append(l.DifficultyAdjustmentLogs, adjustmentLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) LogBlockGeneration(blockID string, generationTime time.Duration) error {
	l.Lock()
	defer l.Unlock()
	logEntry := BlockGenerationLog{
		BlockID:        blockID,
		GenerationTime: generationTime,
		Timestamp:      time.Now(),
	}
	l.BlockGenerationLogs = append(l.BlockGenerationLogs, logEntry)
	return nil
}

func (l *BlockchainConsensusCoinLedger) EnableConsensusAudit() error {
	l.Lock()
	defer l.Unlock()
	l.ConsensusMonitoringEnabled = true
	return nil
}

func (l *BlockchainConsensusCoinLedger) DisableConsensusAudit() error {
	l.Lock()
	defer l.Unlock()
	l.ConsensusMonitoringEnabled = false
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetRewardDistributionMode(mode RewardDistributionMode) error {
	l.Lock()
	defer l.Unlock()
	l.RewardDistributionMode = mode
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetRewardDistributionMode() (RewardDistributionMode, error) {
	l.Lock()
	defer l.Unlock()
	return l.RewardDistributionMode, nil
}

func (l *BlockchainConsensusCoinLedger) LogConsensusParticipation(validatorID, status string) error {
	l.Lock()
	defer l.Unlock()
	auditLog := ConsensusAuditLog{
		AuditID:             fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		ValidatorID:         validatorID,
		Timestamp:           time.Now(),
		ParticipationStatus: status,
	}
	l.ConsensusAuditLogs = append(l.ConsensusAuditLogs, auditLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetValidatorSelectionMode(mode ValidatorSelectionMode) error {
	l.Lock()
	defer l.Unlock()
	l.ValidatorSelectionMode = mode
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetValidatorSelectionMode() (ValidatorSelectionMode, error) {
	l.Lock()
	defer l.Unlock()
	return l.ValidatorSelectionMode, nil
}

func (l *BlockchainConsensusCoinLedger) SetPoHParticipationThreshold(threshold float64) error {
	l.Lock()
	defer l.Unlock()
	l.PoHParticipationThreshold = threshold
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPoHParticipationThreshold() (float64, error) {
	l.Lock()
	defer l.Unlock()
	return l.PoHParticipationThreshold, nil
}

func (l *BlockchainConsensusCoinLedger) LogValidatorActivity(validatorID, action, details string) error {
	l.Lock()
	defer l.Unlock()
	activityLog := ValidatorActivityLog{
		ValidatorID: validatorID,
		Action:      action,
		Timestamp:   time.Now(),
		Details:     details,
	}
	l.ValidatorActivityLogs = append(l.ValidatorActivityLogs, activityLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) EnableDynamicStakeAdjustment() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.DynamicStakeAdjustment = true
	return nil
}

func (l *BlockchainConsensusCoinLedger) DisableDynamicStakeAdjustment() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.DynamicStakeAdjustment = false
	return nil
}

// Update RecordStakeChange to accept StakeChangeRecord
func (l *BlockchainConsensusCoinLedger) RecordStakeChange(stakeChangeRecord StakeChangeRecord) error {
	l.StakeChanges = append(l.StakeChanges, stakeChangeRecord)
	return nil
}

func (l *BlockchainConsensusCoinLedger) LogStakeAdjustment(encryptedLog StakeLog) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.StakeLogs = append(l.StakeLogs, encryptedLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetValidatorPenalty(validatorID string, encryptedPenalty []byte) error {
	l.Lock()
	defer l.Unlock()

	penalty := ValidatorPenalty{
		ValidatorID:   validatorID,
		PenaltyAmount: encryptedPenalty,
		Timestamp:     time.Now(),
	}
	l.ValidatorPenalties[validatorID] = penalty
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetValidatorPenalty(validatorID string) ([]byte, error) {
	penaltyRecord, exists := l.ValidatorPenalties[validatorID]
	if !exists {
		return nil, fmt.Errorf("penalty not found for validator %s", validatorID)
	}
	return penaltyRecord.PenaltyAmount, nil
}

func (l *BlockchainConsensusCoinLedger) SetEpochTimeout(timeout time.Duration) error {
	l.Lock()
	defer l.Unlock()
	l.EpochTimeout = timeout
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetEpochTimeout() (time.Duration, error) {
	l.Lock()
	defer l.Unlock()
	return l.EpochTimeout, nil
}

func (l *BlockchainConsensusCoinLedger) RecordEpochTime(epochID string, duration time.Duration) error {
	l.Lock()
	defer l.Unlock()

	epochLog := EpochLog{
		EpochID:   epochID,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	l.EpochLogs = append(l.EpochLogs, epochLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) LogEpochChange(epochLog EpochLog) error {
	l.EpochLogs = append(l.EpochLogs, epochLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetReinforcementPolicy(policy ReinforcementPolicy) error {
	l.Lock()
	defer l.Unlock()
	l.ReinforcementPolicy = policy
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetReinforcementPolicy() (ReinforcementPolicy, error) {
	l.Lock()
	defer l.Unlock()
	return l.ReinforcementPolicy, nil
}

func (l *BlockchainConsensusCoinLedger) evaluateConsensusHealth() error {
	l.Lock()
	defer l.Unlock()

	// Simulate calculating health metrics for consensus
	metrics := []struct {
		name  string
		value float64
	}{
		{"stability", calculateStability()},
		{"performance", calculatePerformance()},
		{"participation", calculateParticipation()},
		{"security", calculateSecurity()},
	}

	// Record each metric in the consensus health logs
	for _, metric := range metrics {
		healthLog := HealthLog{
			HealthID:  fmt.Sprintf("health-%s-%d", metric.name, time.Now().UnixNano()),
			Metric:    metric.name,
			Value:     metric.value,
			Timestamp: time.Now(),
		}
		l.ConsensusHealthLogs = append(l.ConsensusHealthLogs, healthLog)
	}
	return nil
}

// Ledger method to log health metrics
func (l *BlockchainConsensusCoinLedger) LogHealthMetrics(healthLog HealthLog) error {
	l.ConsensusHealthLogs = append(l.ConsensusHealthLogs, healthLog)
	return nil
}

func (l *BlockchainConsensusCoinLedger) EnableValidatorBans() error {
	l.Lock()
	defer l.Unlock()
	l.ValidatorBansEnabled = true
	return nil
}

func (l *BlockchainConsensusCoinLedger) DisableValidatorBans() error {
	l.Lock()
	defer l.Unlock()
	l.ValidatorBansEnabled = false
	return nil
}

func (l *BlockchainConsensusCoinLedger) BanValidator(validatorID string, encryptedReason string) error {
	l.Lock()
	defer l.Unlock()

	if !l.ValidatorBansEnabled {
		return errors.New("validator banning is not enabled")
	}
	l.BannedValidators[validatorID] = ValidatorBanRecord{
		ValidatorID: validatorID,
		Reason:      encryptedReason,
		Timestamp:   time.Now(),
	}
	return nil
}

func (l *BlockchainConsensusCoinLedger) UnbanValidator(validatorID string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.BannedValidators[validatorID]; !exists {
		return fmt.Errorf("validator %s is not currently banned", validatorID)
	}
	delete(l.BannedValidators, validatorID)
	return nil
}

func (l *BlockchainConsensusCoinLedger) FetchBannedValidators() ([]string, error) {
	l.Lock()
	defer l.Unlock()

	var bannedValidators []string
	for validatorID := range l.BannedValidators {
		bannedValidators = append(bannedValidators, validatorID)
	}
	return bannedValidators, nil
}

func (l *BlockchainConsensusCoinLedger) AuditValidatorPunishments() error {
	l.Lock()
	defer l.Unlock()

	if len(l.ValidatorPunishments) == 0 {
		return errors.New("no punishment records available for auditing")
	}

	fmt.Println("Auditing validator punishments:")
	for validatorID, records := range l.ValidatorPunishments {
		fmt.Printf("Validator %s has %d punishment(s):\n", validatorID, len(records))
		for _, record := range records {
			fmt.Printf("- Reason: %s, Level: %d, Time: %s\n", record.Reason, record.PunishmentLevel, record.Timestamp)
		}
	}
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPunishmentHistory(validatorID string) ([]PunishmentRecord, error) {
	l.Lock()
	defer l.Unlock()

	history, exists := l.ValidatorPunishments[validatorID]
	if !exists {
		return nil, fmt.Errorf("no punishment history found for validator %s", validatorID)
	}
	return history, nil
}

func (l *BlockchainConsensusCoinLedger) ResetPunishmentCount(validatorID string) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.ValidatorPunishments[validatorID]; !exists {
		return fmt.Errorf("no punishments found for validator %s", validatorID)
	}

	delete(l.ValidatorPunishments, validatorID)
	return nil
}

func (l *BlockchainConsensusCoinLedger) SetAutoPunishmentRate(rate float64) error {
	l.Lock()
	defer l.Unlock()

	if rate < 0 {
		return errors.New("auto punishment rate must be non-negative")
	}
	l.AutoPunishmentRate = rate
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetAutoPunishmentRate() (float64, error) {
	l.Lock()
	defer l.Unlock()
	return l.AutoPunishmentRate, nil
}

func (l *BlockchainConsensusCoinLedger) RecordValidatorReward(validatorID string, encryptedReward string) error {
	l.Lock()
	defer l.Unlock()

	rewardRecord := RewardRecord{
		ValidatorID:     validatorID,
		EncryptedReward: encryptedReward,
		Timestamp:       time.Now(),
	}

	l.ValidatorRewardRecords[validatorID] = append(l.ValidatorRewardRecords[validatorID], rewardRecord)
	return nil
}

func (l *BlockchainConsensusCoinLedger) AuditRewardDistributions() error {
	l.Lock()
	defer l.Unlock()

	if len(l.ValidatorRewardRecords) == 0 {
		return errors.New("no reward records available for auditing")
	}

	fmt.Println("Auditing reward distributions:")
	for validatorID, rewards := range l.ValidatorRewardRecords {
		fmt.Printf("Validator %s has received %d reward(s):\n", validatorID, len(rewards))
		for _, reward := range rewards {
			fmt.Printf("- Encrypted Reward: %s, Time: %s\n", reward.EncryptedReward, reward.Timestamp)
		}
	}
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetRewardHistory(validatorID string) ([]RewardRecord, error) {
	l.Lock()
	defer l.Unlock()

	history, exists := l.ValidatorRewardHistory[validatorID]
	if !exists {
		return nil, fmt.Errorf("no reward history found for validator %s", validatorID)
	}
	return history, nil
}

func (l *BlockchainConsensusCoinLedger) SetPoHValidationWindow(window time.Duration) error {
	l.Lock()
	defer l.Unlock()

	if window <= 0 {
		return errors.New("PoH validation window must be a positive duration")
	}
	l.PoHValidationWindow = window
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPoHValidationWindow() (time.Duration, error) {
	l.Lock()
	defer l.Unlock()
	return l.PoHValidationWindow, nil
}

func (l *BlockchainConsensusCoinLedger) AuditPoHValidation() error {
	l.Lock()
	defer l.Unlock()

	if len(l.PoHValidationLogs) == 0 {
		return errors.New("no PoH validation logs available for auditing")
	}

	fmt.Println("Auditing PoH validation logs:")
	for _, log := range l.PoHValidationLogs {
		fmt.Printf("Validator: %s, Status: %s, Time: %s\n", log.ValidatorID, log.Status, log.Timestamp)
	}
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPoHValidationLogs() ([]PoHLog, error) {
	l.Lock()
	defer l.Unlock()

	if len(l.PoHValidationLogs) == 0 {
		return nil, errors.New("no PoH validation logs found")
	}
	return l.PoHValidationLogs, nil
}

func (l *BlockchainConsensusCoinLedger) SetPoHFailureThreshold(threshold int) error {
	l.Lock()
	defer l.Unlock()

	if threshold < 0 {
		return errors.New("PoH failure threshold must be non-negative")
	}
	l.PoHFailureThreshold = threshold
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPoHFailureThreshold() (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.PoHFailureThreshold, nil
}

func (l *BlockchainConsensusCoinLedger) SetPoWHalvingInterval(interval time.Duration) error {
	l.Lock()
	defer l.Unlock()

	if interval <= 0 {
		return errors.New("PoW halving interval must be a positive duration")
	}
	l.PoWHalvingInterval = interval
	return nil
}

func (l *BlockchainConsensusCoinLedger) GetPoWHalvingInterval() (time.Duration, error) {
	l.Lock()
	defer l.Unlock()
	return l.PoWHalvingInterval, nil
}

// RecordConsensusThreshold logs the consensus threshold in the ledger with a timestamp.
func (l *BlockchainConsensusCoinLedger) RecordConsensusThreshold(threshold int, timestamp string) {
	l.Lock()
	defer l.Unlock()

	l.ConsensusThreshold = threshold
	l.ConsensusThresholdTimestamp = timestamp
	log.Printf("Consensus threshold set to %d%% at %s", threshold, timestamp)
}

// SyncWithConsensus updates the ledger with the latest state from Synnergy Consensus.
func (l *BlockchainConsensusCoinLedger) SyncWithConsensus() error {
	l.Lock()
	defer l.Unlock()

	// Synchronize with Synnergy Consensus: update block height and last block hash
	if newBlock, err := l.ConsensusState.GetLatestBlock(); err == nil {
		l.BlockchainConsensusCoinState.BlockHeight++
		l.BlockchainConsensusCoinState.LastBlockHash = newBlock.Hash()
		l.FinalizedBlocks = append(l.FinalizedBlocks, newBlock)
	} else {
		return errors.New("failed to sync with consensus")
	}
	return nil
}

// RecordValidatorRegistration registers a new validator and records it in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordValidatorRegistration(validatorID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	if l.Validators == nil {
		l.Validators = make(map[string]Validator)
	}

	id := l.GenerateUniqueID(validatorID)
	validator := Validator{
		ID:           id,
		RegisteredAt: time.Now(),
	}

	l.Validators[id] = validator
	return id, nil
}

// RecordMerkleRoot records the Merkle tree root hash in the ledger.
func (l *BlockchainConsensusCoinLedger) RecordMerkleRoot(rootHash string) error {
	if rootHash == "" {
		return errors.New("empty Merkle root hash")
	}
	l.MerkleRoot = append(l.MerkleRoot, rootHash)
	fmt.Printf("Merkle root recorded: %s\n", rootHash)
	return nil
}

// RecordEscrow logs the creation of an escrow contract
func (l *BlockchainConsensusCoinLedger) RecordEscrow(escrowID, buyer, seller, resourceID string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	escrow := Escrow{
		EscrowID:   escrowID,
		Buyer:      buyer,
		Seller:     seller,
		Amount:     amount,     // Ensure amount is float64
		ResourceID: resourceID, // ResourceID is now required as a parameter
		Status:     "active",   // Initial status is "active"
		Timestamp:  time.Now(), // Timestamp for creation
		IsReleased: false,      // Funds have not been released yet
		IsDisputed: false,      // No disputes yet
	}

	l.EscrowRecords[escrowID] = escrow
	return nil
}

// RecordCacheAction caches a transaction
func (l *BlockchainConsensusCoinLedger) RecordCacheAction(record Transaction) error { // Use Transaction instead of TransactionRecord
	l.Lock()
	defer l.Unlock()

	// Store the transaction in the cache if it doesn't already exist
	if _, exists := l.TransactionCache[record.TransactionID]; !exists {
		l.TransactionCache[record.TransactionID] = record
	}
	return nil
}

// RecordPrefetch retrieves a cached transaction by its ID
func (l *BlockchainConsensusCoinLedger) RecordPrefetch(recordID string) (*Transaction, error) { // Return *Transaction
	l.lock.Lock()
	defer l.lock.Unlock()

	// Retrieve the transaction from the cache if it exists
	if record, exists := l.TransactionCache[recordID]; exists {
		return &record, nil
	}
	return nil, errors.New("transaction not found in cache")
}

// RecordCacheInvalidation invalidates a cached transaction record by ID
func (l *BlockchainConsensusCoinLedger) RecordCacheInvalidation(recordID string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, exists := l.TransactionCache[recordID]; exists {
		delete(l.TransactionCache, recordID) // Remove record from cache
		return nil
	}
	return errors.New("transaction record not found in cache")
}
