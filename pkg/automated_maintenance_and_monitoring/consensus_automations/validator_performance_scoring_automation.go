package consensus_automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    MinStakeRequirement      = 1000             // Minimum stake required for validators
    RewardKey                = "reward_key"     // Encryption key for reward data
    PenaltyKey               = "penalty_key"    // Encryption key for penalty data
    UptimeThreshold          = 0.95             // Uptime threshold for rewards
    PerformanceThreshold     = 0.80             // Performance threshold for rewards
    PenaltyThreshold         = 0.60             // Performance threshold for penalties
)

// ValidatorPerformanceScoringAutomation handles the scoring and ranking of validators in Synnergy Consensus
type ValidatorPerformanceScoringAutomation struct {
    ledgerInstance  *ledger.Ledger                      // Blockchain ledger for logging validator scores
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine to fetch validator data
    stateMutex      *sync.RWMutex                       // Mutex for thread-safe ledger access
}

// NewValidatorPerformanceScoringAutomation initializes the performance scoring automation
func NewValidatorPerformanceScoringAutomation(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *ValidatorPerformanceScoringAutomation {
    return &ValidatorPerformanceScoringAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
    }
}

// StartPerformanceMonitoring begins the automated process of monitoring and scoring validators
func (automation *ValidatorPerformanceScoringAutomation) StartPerformanceMonitoring() {
    ticker := time.NewTicker(PerformanceCheckInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Checking validator performance and updating scores...")
            automation.checkValidatorPerformance()
        }
    }()
}

// checkValidatorPerformance fetches the performance data for each validator and updates their score
func (automation *ValidatorPerformanceScoringAutomation) checkValidatorPerformance() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Get all validators' data from the consensus engine
    validators := automation.consensusEngine.GetValidators()

    for _, validator := range validators {
        // Calculate the performance score based on metrics
        score := automation.calculateValidatorScore(validator)

        // Log and update the score in the ledger
        automation.logValidatorScore(validator, score)

        // Trigger rewards or penalties based on performance score
        if score >= PerformanceThreshold {
            automation.rewardValidator(validator)
        } else if score < PenaltyThreshold {
            automation.penalizeValidator(validator)
        }
    }
}

// calculateValidatorScore computes a score for each validator based on uptime, validations, and stake
func (automation *ValidatorPerformanceScoringAutomation) calculateValidatorScore(validator synnergy_consensus.Validator) float64 {
    // Example formula: (uptime * 0.5) + (successful validations * 0.3) + (stake commitment * 0.2)
    uptimeScore := validator.Uptime * 0.5
    validationScore := float64(validator.SuccessfulValidations) / float64(validator.TotalValidations) * 0.3
    stakeScore := float64(validator.Stake) / float64(MinStakeRequirement) * 0.2

    // Total performance score
    totalScore := uptimeScore + validationScore + stakeScore
    fmt.Printf("Validator %s performance score: %.2f\n", validator.ID, totalScore)
    return totalScore
}

// logValidatorScore logs the validator's performance score in the blockchain ledger
func (automation *ValidatorPerformanceScoringAutomation) logValidatorScore(validator synnergy_consensus.Validator, score float64) {
    scoreLog := common.LedgerEntry{
        ID:        fmt.Sprintf("validator-score-%s-%d", validator.ID, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Performance Score",
        Status:    "Updated",
        Details:   fmt.Sprintf("Validator %s score: %.2f", validator.ID, score),
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(scoreLog, []byte(validator.ID)) // Use the validator's ID as a key for encryption
    if err != nil {
        fmt.Printf("Error encrypting score log for validator %s: %v\n", validator.ID, err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Logged score for validator %s\n", validator.ID)
}

// rewardValidator rewards high-performing validators based on their score
func (automation *ValidatorPerformanceScoringAutomation) rewardValidator(validator synnergy_consensus.Validator) {
    fmt.Printf("Rewarding validator %s for high performance.\n", validator.ID)

    // Calculate and distribute rewards (e.g., tokens, extra stake, etc.)
    rewardAmount := automation.consensusEngine.CalculateReward(validator)

    // Encrypt the reward data before logging it
    rewardLog := common.LedgerEntry{
        ID:        fmt.Sprintf("reward-%s-%d", validator.ID, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Reward",
        Status:    "Rewarded",
        Details:   fmt.Sprintf("Validator %s received reward of %.2f", validator.ID, rewardAmount),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(rewardLog, []byte(RewardKey))
    if err != nil {
        fmt.Printf("Error encrypting reward log for validator %s: %v\n", validator.ID, err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Logged reward for validator %s\n", validator.ID)
}

// penalizeValidator penalizes validators that fall below a certain performance threshold
func (automation *ValidatorPerformanceScoringAutomation) penalizeValidator(validator synnergy_consensus.Validator) {
    fmt.Printf("Penalizing validator %s for poor performance.\n", validator.ID)

    // Penalize the validator (e.g., reduced stake, suspension, etc.)
    penaltyAmount := automation.consensusEngine.CalculatePenalty(validator)

    // Encrypt the penalty data before logging it
    penaltyLog := common.LedgerEntry{
        ID:        fmt.Sprintf("penalty-%s-%d", validator.ID, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Validator Penalty",
        Status:    "Penalized",
        Details:   fmt.Sprintf("Validator %s received penalty of %.2f", validator.ID, penaltyAmount),
    }

    encryptedEntry, err := encryption.EncryptLedgerEntry(penaltyLog, []byte(PenaltyKey))
    if err != nil {
        fmt.Printf("Error encrypting penalty log for validator %s: %v\n", validator.ID, err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Printf("Logged penalty for validator %s\n", validator.ID)
}

// triggerPerformanceCheck allows manual or triggered performance checks for validators
func (automation *ValidatorPerformanceScoringAutomation) TriggerPerformanceCheck() {
    fmt.Println("Manually triggering validator performance check...")
    automation.checkValidatorPerformance()
}

// adjustValidatorScores allows dynamic adjustment of scores based on real-time factors (e.g., load or network stress)
func (automation *ValidatorPerformanceScoringAutomation) AdjustValidatorScores() {
    // Dynamically adjust validator scores based on real-time metrics
    fmt.Println("Adjusting validator scores based on real-time factors...")
    automation.checkValidatorPerformance()
}
