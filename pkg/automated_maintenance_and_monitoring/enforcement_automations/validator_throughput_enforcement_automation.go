package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/validator"
)

// Configuration for validator throughput enforcement automation
const (
	ThroughputCheckInterval         = 20 * time.Second // Interval to check validator throughput
	MinThroughputRequirement        = 800              // Minimum transactions per interval
	ThroughputViolationThreshold    = 3                // Allowed violations before penalty
	ThroughputRewardMultiplier      = 1.4              // Multiplier for rewards on high throughput
	ThroughputPenaltyMultiplier     = 0.6              // Multiplier for penalties on low throughput
)

// ValidatorThroughputEnforcementAutomation monitors and enforces throughput standards for validators
type ValidatorThroughputEnforcementAutomation struct {
	validatorManager      *validator.ValidatorManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	violationCountMap     map[string]int // Tracks throughput violations for each validator
}

// NewValidatorThroughputEnforcementAutomation initializes the validator throughput enforcement automation
func NewValidatorThroughputEnforcementAutomation(validatorManager *validator.ValidatorManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ValidatorThroughputEnforcementAutomation {
	return &ValidatorThroughputEnforcementAutomation{
		validatorManager:    validatorManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		violationCountMap:   make(map[string]int),
	}
}

// StartValidatorThroughputEnforcement begins continuous monitoring and enforcement of validator throughput
func (automation *ValidatorThroughputEnforcementAutomation) StartValidatorThroughputEnforcement() {
	ticker := time.NewTicker(ThroughputCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateValidatorThroughput()
		}
	}()
}

// evaluateValidatorThroughput checks validator throughput and rewards or penalizes as needed
func (automation *ValidatorThroughputEnforcementAutomation) evaluateValidatorThroughput() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, validatorID := range automation.validatorManager.GetAllValidators() {
		throughput := automation.validatorManager.GetValidatorThroughput(validatorID)

		if throughput >= MinThroughputRequirement {
			automation.rewardValidatorForThroughput(validatorID)
			automation.violationCountMap[validatorID] = 0 // Reset violation count on sufficient throughput
		} else {
			automation.violationCountMap[validatorID]++
			if automation.violationCountMap[validatorID] >= ThroughputViolationThreshold {
				automation.penalizeValidatorForLowThroughput(validatorID)
			}
		}
	}
}

// rewardValidatorForThroughput rewards a validator for maintaining or exceeding the required throughput
func (automation *ValidatorThroughputEnforcementAutomation) rewardValidatorForThroughput(validatorID string) {
	rewardAmount := automation.calculateThroughputReward(validatorID)
	err := automation.ledgerInstance.AddReward(validatorID, rewardAmount)
	if err != nil {
		fmt.Printf("Failed to reward validator %s for throughput: %v\n", validatorID, err)
		automation.logThroughputEnforcementAction(validatorID, "Throughput Reward Failed", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	} else {
		fmt.Printf("Validator %s rewarded with %f for meeting throughput standards.\n", validatorID, rewardAmount)
		automation.logThroughputEnforcementAction(validatorID, "Validator Throughput Rewarded", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	}
}

// penalizeValidatorForLowThroughput penalizes a validator for not meeting the required throughput
func (automation *ValidatorThroughputEnforcementAutomation) penalizeValidatorForLowThroughput(validatorID string) {
	penaltyAmount := automation.calculateThroughputPenalty(validatorID)
	err := automation.ledgerInstance.ApplyPenalty(validatorID, penaltyAmount)
	if err != nil {
		fmt.Printf("Failed to penalize validator %s for low throughput: %v\n", validatorID, err)
		automation.logThroughputEnforcementAction(validatorID, "Throughput Penalty Failed", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	} else {
		fmt.Printf("Validator %s penalized with %f for failing throughput standards.\n", validatorID, penaltyAmount)
		automation.logThroughputEnforcementAction(validatorID, "Validator Throughput Penalized", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	}
}

// calculateThroughputReward calculates the reward amount based on validator throughput
func (automation *ValidatorThroughputEnforcementAutomation) calculateThroughputReward(validatorID string) float64 {
	baseReward := automation.validatorManager.GetBaseReward(validatorID)
	return baseReward * ThroughputRewardMultiplier
}

// calculateThroughputPenalty calculates the penalty amount for validators with low throughput
func (automation *ValidatorThroughputEnforcementAutomation) calculateThroughputPenalty(validatorID string) float64 {
	basePenalty := automation.validatorManager.GetBaseReward(validatorID)
	return basePenalty * ThroughputPenaltyMultiplier
}

// logThroughputEnforcementAction securely logs actions related to throughput enforcement
func (automation *ValidatorThroughputEnforcementAutomation) logThroughputEnforcementAction(validatorID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Validator ID: %s, Details: %s", action, validatorID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-throughput-enforcement-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Throughput Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log throughput enforcement action for validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Throughput enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorThroughputEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualThroughputReward allows administrators to manually reward a validator for high throughput
func (automation *ValidatorThroughputEnforcementAutomation) TriggerManualThroughputReward(validatorID string, amount float64) {
	fmt.Printf("Manually rewarding validator %s with %f for exceptional throughput.\n", validatorID, amount)
	err := automation.ledgerInstance.AddReward(validatorID, amount)
	if err != nil {
		fmt.Printf("Failed to manually reward validator %s: %v\n", validatorID, err)
		automation.logThroughputEnforcementAction(validatorID, "Manual Throughput Reward Failed", fmt.Sprintf("Amount: %f", amount))
	} else {
		automation.logThroughputEnforcementAction(validatorID, "Manual Throughput Rewarded", fmt.Sprintf("Amount: %f", amount))
	}
}
