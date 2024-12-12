package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strings"
	"synnergy_network/pkg/ledger"
	"time"
)

// PoS represents the Proof of Stake mechanism, integrated with the Ledger and Reward Manager
type PoS struct {
	State              PoSState                  // Current state of the PoS system
	LedgerInstance     *ledger.Ledger            // Ledger instance for tracking stake and rewards
	RewardManager      *RewardManager            // Reward Manager to handle validator rewards
	MinStake           float64                   // Minimum stake required to validate
	TransactionTypeMap map[string]string         // Maps TransactionType field to token standards
	TransactionFunctionMap   map[string]string         // Maps TransactionFunction field to validation logic keys
	ValidationMap      map[string]func(tx Transaction) bool // Maps keys to actual validation functions
}

// PoSState represents the current state of the PoS system
type PoSState struct {
	Validators   []Validator // List of all validators participating in PoS
	TotalStake   float64     // Total amount of SYNN staked in the network
	LastSelected string      // Last validator selected to validate a sub-block
	Epoch        int         // Epoch number for selecting new validators
	LastUpdated  time.Time   // Timestamp of the last validator selection
}

// NewPoS initializes the PoS system with a list of validators and integrates it with the ledger and reward manager.
func NewPoS(validators []Validator, ledgerInstance *ledger.Ledger, rewardManager *RewardManager, minStake float64) *PoS {
    if len(validators) == 0 {
        log.Fatalf("PoS initialization failed: no validators provided")
    }
    if minStake <= 0 {
        log.Fatalf("PoS initialization failed: minimum stake must be greater than 0")
    }

    totalStake := 0.0
    for _, validator := range validators {
        totalStake += validator.Stake
    }

    log.Printf("[Info] PoS initialized with %d validators and total stake: %.2f", len(validators), totalStake)

    return &PoS{
        State: PoSState{
            Validators:   validators,
            TotalStake:   totalStake,
            LastSelected: "",
            Epoch:        0,
            LastUpdated:  time.Now(),
        },
        LedgerInstance:       ledgerInstance,
        RewardManager:        rewardManager,
        MinStake:             minStake,
        TransactionTypeMap:   NewTransactionTypeMap(),
        TransactionFunctionMap: NewTransactionFunctionMap(),
        ValidationMap:        NewValidationMap(),
    }
}


// SelectValidator selects a validator to validate a sub-block based on their stake.
func (pos *PoS) SelectValidator() (Validator, error) {
    // Validate the PoS state
    if len(pos.State.Validators) == 0 {
        return Validator{}, fmt.Errorf("validator selection failed: no validators available")
    }
    if pos.State.TotalStake <= 0 {
        return Validator{}, fmt.Errorf("validator selection failed: total stake must be greater than 0")
    }

    log.Printf("[Info] Starting validator selection for epoch %d", pos.State.Epoch+1)

    // Perform weighted random selection based on stake
    stakeThreshold := rand.Float64() * pos.State.TotalStake
    cumulativeStake := 0.0
    var selectedValidator Validator

    for _, validator := range pos.State.Validators {
        cumulativeStake += validator.Stake
        if cumulativeStake >= stakeThreshold {
            selectedValidator = validator
            pos.State.LastSelected = validator.Address
            pos.State.LastUpdated = time.Now()
            pos.State.Epoch++

            log.Printf("[Success] Validator %s selected for epoch %d with stake %.2f", 
                validator.Address, pos.State.Epoch, validator.Stake)

            // Record the selection in the ledger
            err := pos.LedgerInstance.BlockchainConsensusCoinLedger.RecordValidatorSelection(validator.Address, pos.State.Epoch)
            if err != nil {
                return Validator{}, fmt.Errorf("failed to record validator selection in ledger: %w", err)
            }

            // Distribute rewards to the selected validator
            err = pos.RewardManager.DistributePoSRewards(validator)
            if err != nil {
                return Validator{}, fmt.Errorf("failed to distribute rewards to validator %s: %w", validator.Address, err)
            }

            return validator, nil
        }
    }

    return Validator{}, fmt.Errorf("validator selection failed due to stake computation error")
}


// AddStake adds more stake to a validator's account and updates the ledger.
func (pos *PoS) AddStake(validatorAddress string, stakeAmount float64) error {
    if validatorAddress == "" {
        return fmt.Errorf("validator address cannot be empty")
    }
    if stakeAmount <= 0 {
        return fmt.Errorf("stake amount must be greater than 0")
    }

    log.Printf("[Info] Adding %.2f SYNN to validator %s's stake", stakeAmount, validatorAddress)

    // Find the validator
    for i, validator := range pos.State.Validators {
        if validator.Address == validatorAddress {
            // Update stake
            pos.State.Validators[i].Stake += stakeAmount
            pos.State.TotalStake += stakeAmount

            log.Printf("[Success] Updated validator %s's stake to %.2f. Total stake: %.2f", 
                validatorAddress, pos.State.Validators[i].Stake, pos.State.TotalStake)

            // Update the ledger
            err := pos.LedgerInstance.BlockchainConsensusCoinLedger.UpdateValidatorStake(validatorAddress, pos.State.Validators[i].Stake)
            if err != nil {
                return fmt.Errorf("failed to update validator stake in ledger: %w", err)
            }

            return nil
        }
    }

    return fmt.Errorf("validator %s not found", validatorAddress)
}


// RemoveStake removes stake from a validator's account and updates the ledger.
func (pos *PoS) RemoveStake(validatorAddress string, stakeAmount float64) error {
    if validatorAddress == "" {
        return fmt.Errorf("validator address cannot be empty")
    }
    if stakeAmount <= 0 {
        return fmt.Errorf("stake amount must be greater than 0")
    }

    log.Printf("[Info] Removing %.2f SYNN from validator %s's stake", stakeAmount, validatorAddress)

    // Find the validator
    for i, validator := range pos.State.Validators {
        if validator.Address == validatorAddress {
            if validator.Stake < stakeAmount {
                return fmt.Errorf("insufficient stake: validator %s only has %.2f SYNN", 
                    validatorAddress, validator.Stake)
            }

            // Update stake
            pos.State.Validators[i].Stake -= stakeAmount
            pos.State.TotalStake -= stakeAmount

            log.Printf("[Success] Reduced validator %s's stake to %.2f. Total stake: %.2f", 
                validatorAddress, pos.State.Validators[i].Stake, pos.State.TotalStake)

            // Update the ledger
            err := pos.LedgerInstance.BlockchainConsensusCoinLedger.UpdateValidatorStake(validatorAddress, pos.State.Validators[i].Stake)
            if err != nil {
                return fmt.Errorf("failed to update validator stake in ledger: %w", err)
            }

            return nil
        }
    }

    return fmt.Errorf("validator %s not found", validatorAddress)
}


// GetValidatorStake retrieves the stake of a specific validator.
func (pos *PoS) GetValidatorStake(validatorAddress string) (float64, error) {
    if validatorAddress == "" {
        return 0, fmt.Errorf("validator address cannot be empty")
    }

    log.Printf("[Info] Retrieving stake for validator: %s", validatorAddress)

    for _, validator := range pos.State.Validators {
        if validator.Address == validatorAddress {
            log.Printf("[Success] Validator %s found with stake: %.2f", validatorAddress, validator.Stake)
            return validator.Stake, nil
        }
    }

    return 0, fmt.Errorf("validator %s not found", validatorAddress)
}


// GetTotalStake returns the total amount of SYNN staked in the network.
func (pos *PoS) GetTotalStake() float64 {
    log.Printf("[Info] Total stake in the network: %.2f", pos.State.TotalStake)
    return pos.State.TotalStake
}


// ValidateSubBlock validates a sub-block using PoS rules.
func (pos *PoS) ValidateSubBlock(subBlock SubBlock) bool {
    if subBlock.Validator == "" {
        log.Printf("[Error] Validator address in sub-block cannot be empty")
        return false
    }

    log.Printf("[Info] Validating sub-block %d by validator %s...", subBlock.Index, subBlock.Validator)

    // Step 1: Select the validator.
    selectedValidator, err := pos.SelectValidator()
    if err != nil {
        log.Printf("[Error] Failed to select validator: %v", err)
        return false
    }
    if selectedValidator.Address != subBlock.Validator {
        log.Printf("[Error] Validator mismatch: Expected %s, Got %s", selectedValidator.Address, subBlock.Validator)
        return false
    }
    log.Printf("[Info] Validator %s selected for sub-block %d.", selectedValidator.Address, subBlock.Index)

    // Step 2: Check if the validator has enough stake.
    validatorStake, err := pos.GetValidatorStake(subBlock.Validator)
    if err != nil {
        log.Printf("[Error] Validator %s not found: %v", subBlock.Validator, err)
        return false
    }
    if validatorStake < pos.MinStake {
        log.Printf("[Error] Validator %s does not have enough stake to validate the sub-block. Required: %.2f, Has: %.2f",
            pos.MinStake, validatorStake, validatorStake)
        return false
    }
    log.Printf("[Info] Validator %s has sufficient stake: %.2f.", subBlock.Validator, validatorStake)

    // Step 3: Verify the validator's signature.
    if !pos.VerifySignature(subBlock.Validator, subBlock.Signature, subBlock.Hash) {
        log.Printf("[Error] Invalid signature for sub-block %d by validator %s.", subBlock.Index, subBlock.Validator)
        return false
    }
    log.Printf("[Info] Signature verified for sub-block %d by validator %s.", subBlock.Index, subBlock.Validator)

    // Step 4: Check the integrity of each transaction in the sub-block.
    for _, tx := range subBlock.Transactions {
        if !pos.ValidateTransaction(tx) {
            log.Printf("[Error] Invalid transaction %s in sub-block %d.", tx.TransactionID, subBlock.Index)
            return false
        }
        log.Printf("[Info] Transaction %s in sub-block %d validated successfully.", tx.TransactionID, subBlock.Index)
    }

    // Step 5: Verify that the sub-block adheres to consensus rules.
    if !pos.CheckConsensusCompliance(subBlock) {
        log.Printf("[Error] Sub-block %d does not comply with PoS consensus rules.", subBlock.Index)
        return false
    }
    log.Printf("[Info] Sub-block %d complies with PoS consensus rules.", subBlock.Index)

    // Step 6: If all checks pass, the sub-block is valid.
    log.Printf("[Success] Sub-block %d by validator %s validated successfully.", subBlock.Index, subBlock.Validator)
    return true
}



// VerifySignature verifies the signature of a validator for a given block hash.
func (pos *PoS) VerifySignature(validatorAddress string, signature string, hash string) bool {
    if validatorAddress == "" {
        log.Printf("[Error] Validator address cannot be empty.")
        return false
    }
    if signature == "" || hash == "" {
        log.Printf("[Error] Signature or hash cannot be empty for validator %s.", validatorAddress)
        return false
    }

    // Retrieve validator details
    validator, err := pos.getValidatorByAddress(validatorAddress)
    if err != nil {
        log.Printf("[Error] Failed to retrieve validator %s: %v", validatorAddress, err)
        return false
    }

    // Ensure the validator's public key is available
    if validator.PublicKey == nil {
        log.Printf("[Error] Validator %s does not have a public key.", validatorAddress)
        return false
    }

    // Split the signature into components (assumes a standardized format).
    signatureParts := strings.Split(signature, ",")
    if len(signatureParts) != 2 {
        log.Printf("[Error] Invalid signature format for validator %s.", validatorAddress)
        return false
    }
    rStr, sStr := signatureParts[0], signatureParts[1]

    // Verify the ECDSA signature
    isValid := VerifyECDSASignature(validator.PublicKey, hash, rStr, sStr)
    if !isValid {
        log.Printf("[Error] Invalid signature for validator %s.", validatorAddress)
        return false
    }

    log.Printf("[Success] Signature verified for validator %s.", validatorAddress)
    return true
}





// getValidatorByAddress retrieves a validator by their address.
func (pos *PoS) getValidatorByAddress(validatorAddress string) (Validator, error) {
    if validatorAddress == "" {
        return Validator{}, fmt.Errorf("validator address cannot be empty")
    }

    for _, validator := range pos.State.Validators {
        if validator.Address == validatorAddress {
            return validator, nil
        }
    }

    return Validator{}, fmt.Errorf("validator %s not found", validatorAddress)
}

// ValidateTransaction checks if a transaction is valid based on PoS rules and the current ledger state.
func (pos *PoS) ValidateTransaction(tx Transaction) bool {
    log.Printf("[Info] Validating transaction %s...", tx.TransactionID)

    // Step 1: Validate inputs
    if tx.TransactionID == "" || tx.FromAddress == "" || tx.Amount <= 0 {
        log.Printf("[Error] Invalid transaction: %v", tx)
        return false
    }

    // Step 2: Identify the TransactionType and TransactionFunction
    transactionType, typeExists := pos.TransactionTypeMap[tx.TransactionType]
    if !typeExists {
        log.Printf("[Error] Unknown TransactionType: %s for transaction %s.", tx.TransactionType, tx.TransactionID)
        return false
    }

    transactionFunction, funcExists := pos.TransactionFunctionMap[tx.TransactionFunction]
    if !funcExists {
        log.Printf("[Error] Unknown TransactionFunction: %s for transaction %s.", tx.TransactionFunction, tx.TransactionID)
        return false
    }

    // Step 3: Derive the validation key
    validationKey := fmt.Sprintf("%s:%s", transactionType, transactionFunction)
    validationFunc, validationExists := pos.ValidationMap[validationKey]
    if !validationExists {
        log.Printf("[Error] No validation logic found for key: %s in transaction %s.", validationKey, tx.TransactionID)
        return false
    }

    // Step 4: Verify the transaction signature
    transactionHash := GenerateHash(fmt.Sprintf("%s:%s:%.2f", tx.TransactionID, tx.FromAddress, tx.Amount))
    if !pos.VerifySignature(tx.FromAddress, tx.Signature, transactionHash) {
        log.Printf("[Error] Invalid signature for transaction %s.", tx.TransactionID)
        return false
    }
    log.Printf("[Info] Signature verified for transaction %s.", tx.TransactionID)

    // Step 5: Ensure the sender has sufficient balance in the ledger
    senderBalance, err := pos.LedgerInstance.GetBalance(tx.FromAddress, tx.TokenStandard)
    if err != nil {
        log.Printf("[Error] Failed to retrieve balance for address %s: %v", tx.FromAddress, err)
        return false
    }
    if senderBalance < tx.Amount {
        log.Printf("[Error] Insufficient balance for transaction %s. Required: %.2f, Available: %.2f", 
            tx.TransactionID, tx.Amount, senderBalance)
        return false
    }
    log.Printf("[Info] Sufficient balance confirmed for transaction %s.", tx.TransactionID)

    // Step 6: Validate transaction fields using the mapped validation function
    if !validationFunc(tx) {
        log.Printf("[Error] Validation logic failed for transaction %s.", tx.TransactionID)
        return false
    }

    log.Printf("[Success] Transaction %s validated successfully.", tx.TransactionID)
    return true
}

// VerifyECDSASignature verifies an ECDSA signature.
func VerifyECDSASignature(publicKey *ecdsa.PublicKey, hash, rStr, sStr string) bool {
    r := new(big.Int)
    s := new(big.Int)
    r.SetString(rStr, 10)
    s.SetString(sStr, 10)

    hashBytes, err := hex.DecodeString(hash)
    if err != nil {
        log.Printf("[Error] Failed to decode hash: %v", err)
        return false
    }

    return ecdsa.Verify(publicKey, hashBytes, r, s)
}



// GenerateHash creates a secure SHA-256 hash for a given string.
func GenerateHash(data string) string {
    if data == "" {
        log.Printf("[Error] Data to hash cannot be empty.")
        return ""
    }

    // Create a new SHA-256 hash object
    hash := sha256.New()

    // Write the data into the hash object
    _, err := hash.Write([]byte(data))
    if err != nil {
        log.Printf("[Error] Failed to write data to hash object: %v", err)
        return ""
    }

    // Compute the hash and return it as a hexadecimal string
    result := fmt.Sprintf("%x", hash.Sum(nil))
    log.Printf("[Info] Hash generated successfully: %s", result)
    return result
}



// CheckConsensusCompliance verifies that the sub-block adheres to PoS consensus rules.
func (pos *PoS) CheckConsensusCompliance(subBlock SubBlock) bool {
    log.Printf("[Info] Checking consensus compliance for sub-block %d...", subBlock.Index)

    // Step 1: Validate sub-block timing
    currentTime := time.Now()
    timeSinceLastBlock := currentTime.Sub(subBlock.Timestamp)
    if timeSinceLastBlock < time.Second*10 { // Example rule: sub-blocks must be at least 10 seconds apart
        log.Printf("[Error] Sub-block %d proposed too soon: %s elapsed since last block.", 
            subBlock.Index, timeSinceLastBlock)
        return false
    }

    // Step 2: Verify the block hash meets the difficulty requirements
    if !pos.CheckBlockHash(subBlock.Hash) {
        log.Printf("[Error] Sub-block %d hash does not meet difficulty requirements.", subBlock.Index)
        return false
    }

    // Step 3: Ensure sub-block links correctly to the chain
    prevBlockHash, err := pos.LedgerInstance.BlockchainConsensusCoinLedger.GetPreviousBlockHash(subBlock.SubBlockID)
    if err != nil {
        log.Printf("[Error] Failed to retrieve previous block hash for sub-block %d: %v", subBlock.Index, err)
        return false
    }
    if subBlock.PrevHash != prevBlockHash {
        log.Printf("[Error] Sub-block %d does not correctly link to previous block. Expected: %s, Got: %s",
            subBlock.Index, prevBlockHash, subBlock.PrevHash)
        return false
    }

    log.Printf("[Success] Sub-block %d complies with PoS consensus rules.", subBlock.Index)
    return true
}


// CheckBlockHash verifies that the sub-block hash meets the network's current difficulty requirements.
func (pos *PoS) CheckBlockHash(hash string) bool {
    if hash == "" {
        log.Printf("[Error] Hash cannot be empty for difficulty check.")
        return false
    }

    // Retrieve current difficulty from the ledger
    currentDifficulty, err := pos.LedgerInstance.BlockchainConsensusCoinLedger.GetCurrentDifficulty()
    if err != nil {
        log.Printf("[Error] Failed to retrieve current difficulty: %v", err)
        return false
    }

    // Generate the required prefix based on the difficulty
    requiredPrefix := generateDifficultyPrefix(currentDifficulty)

    // Check if the hash starts with the required prefix
    if len(hash) < len(requiredPrefix) || hash[:len(requiredPrefix)] != requiredPrefix {
        log.Printf("[Error] Hash does not meet difficulty requirements. Required: %s, Got: %s",
            requiredPrefix, hash[:len(requiredPrefix)])
        return false
    }

    log.Printf("[Success] Hash meets difficulty requirements.")
    return true
}

// generateDifficultyPrefix generates a string of leading zeros based on the current difficulty.
func generateDifficultyPrefix(difficulty int) string {
    if difficulty <= 0 {
        log.Printf("[Warning] Difficulty must be greater than 0. Defaulting to 1.")
        difficulty = 1
    }

    prefix := fmt.Sprintf("%0*s", difficulty, "")
    log.Printf("[Info] Generated difficulty prefix: %s", prefix)
    return prefix
}



