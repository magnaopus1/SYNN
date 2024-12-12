package common

import (
	"fmt"
	"math"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// RewardManager handles the distribution of rewards for PoH, PoS, and PoW
type RewardManager struct {
	PoSRewardRate      float64        // Percentage reward for staking in PoS
	PoHRewardRate      float64        // Percentage reward for participating in PoH
	PoWInitialReward   float64        // Initial block reward for PoW
	PoWHalvingInterval int            // Number of blocks before PoW reward halves
	CurrentBlockCount  int            // Current block count to track PoW halving
	RewardPool         float64        // Combined pool for validator and miner rewards
	LedgerInstance     *ledger.Ledger // Instance of the ledger for reward tracking
	PunishmentManager  *PunishmentManager // Reference to the PunishmentManager for enforcing penalties
	mutex              sync.Mutex     // Mutex for thread-safe operations
}

// Punishment represents a punishment with its timestamp
type Punishment struct {
	Amount     float64   // Punishment amount
	Timestamp  time.Time // Time when the punishment was applied
}

// PunishmentManager manages punishments for validators and miners
type PunishmentManager struct {
	PoSPunishmentThreshold  float64        // Threshold for PoS violations (e.g., 24 hours of downtime)
	PoHPunishmentThreshold  float64        // Threshold for PoH inactivity (e.g., 1 missed participation cycle)
	PoWPunishmentThreshold  float64        // Threshold for PoW failure (e.g., 5 failed block attempts)
	PoSPunishmentRate       float64        // Punishment rate for PoS violations (e.g., percentage reduction in stake)
	PoHPunishmentRate       float64        // Punishment rate for PoH violations
	PoWPunishmentRate       float64        // Punishment rate for PoW failures
	PunishmentHistory       map[string][]Punishment // Punishment history for each entity
	PunishmentResetInterval time.Duration  // Time interval for punishment reset (e.g., 90 days)
	LedgerInstance          *ledger.Ledger // Reference to the ledger for recording punishments
	mutex                   sync.Mutex     // Mutex for thread-safe operations
}


// NewRewardManager initializes the RewardManager with given reward rates and ledger instance
func NewRewardManager(ledgerInstance *ledger.Ledger, punishmentManager *PunishmentManager) *RewardManager {
    return &RewardManager{
        PoSRewardRate:      3.0,               // 3% staking reward for PoS
        PoHRewardRate:      1.0,               // 1% participation reward for PoH
        PoWInitialReward:   1024,              // Initial PoW reward per block
        PoWHalvingInterval: 200000,            // Halving reward every 200,000 blocks
        CurrentBlockCount:  0,
        RewardPool:         0.0,               // Initialize with 0 SYNN in the reward pool
        LedgerInstance:     ledgerInstance,
        PunishmentManager:  punishmentManager, // Link to the PunishmentManager
    }
}

// NewPunishmentManager initializes the PunishmentManager with thresholds, punishment rates, and a reset interval of 90 days
func NewPunishmentManager(ledgerInstance *ledger.Ledger) *PunishmentManager {
    return &PunishmentManager{
        PoSPunishmentThreshold:  24.0,                     // 24 hours of downtime triggers PoS punishment
        PoHPunishmentThreshold:  1.0,                      // 1 missed cycle triggers PoH punishment
        PoWPunishmentThreshold:  5.0,                      // 5 failed block attempts trigger PoW punishment
        PoSPunishmentRate:       5.0,                      // 5% reduction in stake for PoS violations
        PoHPunishmentRate:       3.0,                      // 3% reduction for PoH inactivity
        PoWPunishmentRate:       10.0,                     // 10% reduction for PoW failures
        PunishmentHistory:       make(map[string][]Punishment), // Track punishment history
        PunishmentResetInterval: 90 * 24 * time.Hour,       // 90 days punishment reset interval
        LedgerInstance:          ledgerInstance,
    }
}

// EnforcePunishments applies punishments to validators and miners based on violations
func (pm *PunishmentManager) EnforcePunishments(violations map[string]float64, category string) {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    for entity, violationLevel := range violations {
        if violationLevel == 0 {
            continue
        }

        var punishmentAmount float64

        switch category {
        case "PoS":
            if violationLevel >= pm.PoSPunishmentThreshold {
                punishmentAmount = (violationLevel / pm.PoSPunishmentThreshold) * pm.PoSPunishmentRate
                pm.applyPunishment(entity, punishmentAmount, category)
            }
        case "PoH":
            if violationLevel >= pm.PoHPunishmentThreshold {
                punishmentAmount = (violationLevel / pm.PoHPunishmentThreshold) * pm.PoHPunishmentRate
                pm.applyPunishment(entity, punishmentAmount, category)
            }
        case "PoW":
            if violationLevel >= pm.PoWPunishmentThreshold {
                punishmentAmount = (violationLevel / pm.PoWPunishmentThreshold) * pm.PoWPunishmentRate
                pm.applyPunishment(entity, punishmentAmount, category)
            }
        }
    }
}

// applyPunishment updates the ledger with punishment details and records the timestamp
func (pm *PunishmentManager) applyPunishment(entity string, punishmentAmount float64, category string) {
    now := time.Now()

    // Create an instance of the Encryption struct
    encryptionService := &Encryption{}

    // Encrypt the punishment amount
    encryptedPunishment, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", punishmentAmount)), EncryptionKey)
    if err != nil {
        fmt.Printf("Failed to encrypt punishment for %s: %v\n", entity, err)
        return
    }

    // Create a Punishment struct instead of using a string
    punishment := ledger.Punishment{
        Amount:    punishmentAmount,
        Entity:    entity,
        Timestamp: now,
    }

    // Update the ledger with the punishment details
    pm.LedgerInstance.BlockchainConsensusCoinLedger.UpdatePunishment(entity, punishment) // Pass the Punishment struct, not the string
    pm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedPunishment(entity, string(encryptedPunishment))

    // Record the punishment in the history with the current timestamp
    pm.PunishmentHistory[entity] = append(pm.PunishmentHistory[entity], Punishment{Amount: punishmentAmount, Timestamp: now})

    fmt.Printf("%s punished with %.2f SYNN for %s violation.\n", entity, punishmentAmount, category)
}



// ResetPunishments removes punishments older than 90 days
func (pm *PunishmentManager) ResetPunishments() {
    pm.mutex.Lock()
    defer pm.mutex.Unlock()

    now := time.Now()

    for entity, punishments := range pm.PunishmentHistory {
        var updatedPunishments []Punishment

        for _, punishment := range punishments {
            if now.Sub(punishment.Timestamp) < pm.PunishmentResetInterval {
                updatedPunishments = append(updatedPunishments, punishment)
            }
        }

        // Update the entity's punishment history
        pm.PunishmentHistory[entity] = updatedPunishments

        if len(punishments) != len(updatedPunishments) {
            fmt.Printf("Punishments older than 90 days for %s have been reset.\n", entity)
        }
    }
}

// DistributeRewards splits the reward pool between validators and miners based on their contribution
func (rm *RewardManager) DistributeRewards(validators map[string]float64, miners map[string]float64, totalSubBlockContribution float64, totalBlockContribution float64) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    if rm.RewardPool == 0 {
        fmt.Println("No rewards available for distribution.")
        return
    }

    // Split rewards between validators (35%) and miners (65%)
    validatorReward := rm.RewardPool * 0.35
    minerReward := rm.RewardPool * 0.65

    fmt.Printf("Distributing %.2f SYNN to validators and %.2f SYNN to miners.\n", validatorReward, minerReward)

    // Create an encryption instance
    encryptionService := &Encryption{}

    // Distribute validator rewards based on contribution to sub-block validation
    for validator, contribution := range validators {
        if totalSubBlockContribution > 0 {
            rewardShare := (contribution / totalSubBlockContribution) * validatorReward
            encryptedReward, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", rewardShare)), EncryptionKey)
            if err != nil {
                fmt.Printf("Failed to encrypt validator reward: %v\n", err)
                continue
            }
            rm.LedgerInstance.BlockchainConsensusCoinLedger.UpdateValidatorReward(validator, rewardShare)
            rm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedReward(validator, string(encryptedReward))
            fmt.Printf("Validator %s awarded %.2f SYNN for sub-block validation.\n", validator, rewardShare)
        }
    }

    // Distribute miner rewards based on contribution to block mining
    for miner, contribution := range miners {
        if totalBlockContribution > 0 {
            rewardShare := (contribution / totalBlockContribution) * minerReward
            encryptedReward, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", rewardShare)), EncryptionKey)
            if err != nil {
                fmt.Printf("Failed to encrypt miner reward: %v\n", err)
                continue
            }
            rm.LedgerInstance.BlockchainConsensusCoinLedger.UpdateMinerReward(miner, rewardShare)
            rm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedReward(miner, string(encryptedReward))
            fmt.Printf("Miner %s awarded %.2f SYNN for block mining.\n", miner, rewardShare)
        }
    }

    // Reset the reward pool after distribution
    rm.RewardPool = 0.0
}

// DistributePoSRewards distributes staking rewards to PoS validators based on their stake
func (rm *RewardManager) DistributePoSRewards(validator Validator) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    reward := (validator.Stake * rm.PoSRewardRate) / 100
    encryptionService := &Encryption{}

    encryptedReward, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", reward)), EncryptionKey)
    if err != nil {
        fmt.Printf("Failed to encrypt PoS reward: %v\n", err)
        return
    }

    rm.LedgerInstance.BlockchainConsensusCoinLedger.UpdateValidatorStake(validator.Address, validator.Stake+reward)
    rm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedReward(validator.Address, string(encryptedReward))

    fmt.Printf("Distributed %.2f SYNN PoS reward to validator %s.\n", reward, validator.Address)
}

// DistributePoHRewards distributes participation rewards to PoH participants
func (rm *RewardManager) DistributePoHRewards(participant string, participationTime float64) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    reward := (participationTime * rm.PoHRewardRate) / 100
    encryptionService := &Encryption{}

    encryptedReward, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", reward)), EncryptionKey)
    if err != nil {
        fmt.Printf("Failed to encrypt PoH reward: %v\n", err)
        return
    }

    rm.LedgerInstance.BlockchainConsensusCoinLedger.UpdateParticipantReward(participant, reward)
    rm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedReward(participant, string(encryptedReward))

    fmt.Printf("Distributed %.2f SYNN PoH reward to participant %s.\n", reward, participant)
}


// DistributePoWRewards handles block mining rewards based on contribution to sub-blocks
func (rm *RewardManager) DistributePoWRewards(minerAddress string) error {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    currentReward := rm.calculateCurrentPoWReward()

    // Create an instance of the Encryption service
    encryptionService := &Encryption{}

    // Define the encryption key (must be exactly 32 bytes for AES-256)
    encryptionKey := []byte("this-is-a-32-byte-encryption-key!!") // Adjusted to 32 bytes
    if len(encryptionKey) != 32 {
        return fmt.Errorf("encryption key must be 32 bytes long")
    }

    // Encrypt the reward before storing it
    encryptedReward, err := encryptionService.EncryptData("AES", []byte(fmt.Sprintf("%.2f", currentReward)), encryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt PoW reward: %v", err)
    }

    // Update the ledger with the miner's reward
    err = rm.LedgerInstance.BlockchainConsensusCoinLedger.UpdateMinerReward(minerAddress, currentReward)
    if err != nil {
        return fmt.Errorf("failed to update miner reward: %v", err)
    }

    // Store the encrypted reward (no return value expected)
    rm.LedgerInstance.BlockchainConsensusCoinLedger.StoreEncryptedReward(minerAddress, string(encryptedReward))

    fmt.Printf("Distributed %.2f SYNN PoW reward to miner %s.\n", currentReward, minerAddress)

    // Increment block count and check for halving (if applicable)
    rm.CurrentBlockCount++

    return nil
}






// calculateCurrentPoWReward calculates the PoW reward based on the halving interval
func (rm *RewardManager) calculateCurrentPoWReward() float64 {
    halvings := rm.CurrentBlockCount / rm.PoWHalvingInterval // Integer division
    currentReward := rm.PoWInitialReward / math.Pow(2, float64(halvings)) // Use exponentiation instead of bitwise shift
    if currentReward < 1 {
        currentReward = 1 // Ensure reward doesn't go below 1 SYNN
    }
    return currentReward
}





// EnforcePunishments allows the RewardManager to call the PunishmentManager's method when necessary
func (rm *RewardManager) EnforcePunishments(violations map[string]float64, category string) {
    rm.PunishmentManager.EnforcePunishments(violations, category)
}
