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

// Configuration for validator incentive enforcement automation
const (
	IncentiveCheckInterval        = 15 * time.Second // Interval to check validator incentives and penalties
	RewardMultiplier              = 1.2              // Multiplier for regular rewards
	PenaltyMultiplier             = 0.5              // Multiplier for penalties on inactive validators
	MinimumActivityRequirement    = 100              // Minimum activity level required for rewards
	MaximumInactivityTolerance    = 3                // Number of intervals allowed for inactivity before penalty
)

// ValidatorIncentiveEnforcementAutomation monitors and enforces incentive policies for validators
type ValidatorIncentiveEnforcementAutomation struct {
	validatorManager      *validator.ValidatorManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	inactiveValidatorMap  map[string]int // Tracks consecutive inactivity count for each validator
}

// NewValidatorIncentiveEnforcementAutomation initializes the validator incentive enforcement automation
func NewValidatorIncentiveEnforcementAutomation(validatorManager *validator.ValidatorManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ValidatorIncentiveEnforcementAutomation {
	return &ValidatorIncentiveEnforcementAutomation{
		validatorManager:     validatorManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		inactiveValidatorMap: make(map[string]int),
	}
}

// StartValidatorIncentiveEnforcement begins continuous monitoring and enforcement of validator incentives
func (automation *ValidatorIncentiveEnforcementAutomation) StartValidatorIncentiveEnforcement() {
	ticker := time.NewTicker(IncentiveCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateValidatorIncentives()
		}
	}()
}

// evaluateValidatorIncentives assesses validator activity and enforces rewards or penalties accordingly
func (automation *ValidatorIncentiveEnforcementAutomation) evaluateValidatorIncentives() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, validatorID := range automation.validatorManager.GetAllValidators() {
		activityLevel := automation.validatorManager.GetValidatorActivity(validatorID)

		if activityLevel >= MinimumActivityRequirement {
			automation.rewardValidator(validatorID)
			automation.inactiveValidatorMap[validatorID] = 0 // Reset inactivity count on reward
		} else {
			automation.inactiveValidatorMap[validatorID]++
			if automation.inactiveValidatorMap[validatorID] >= MaximumInactivityTolerance {
				automation.penalizeValidator(validatorID)
			}
		}
	}
}

// rewardValidator rewards a validator based on activity level and predefined incentives
func (automation *ValidatorIncentiveEnforcementAutomation) rewardValidator(validatorID string) {
	rewardAmount := automation.calculateReward(validatorID)
	err := automation.ledgerInstance.AddReward(validatorID, rewardAmount)
	if err != nil {
		fmt.Printf("Failed to reward validator %s: %v\n", validatorID, err)
		automation.logIncentiveAction(validatorID, "Reward Failed", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	} else {
		fmt.Printf("Validator %s rewarded with %f.\n", validatorID, rewardAmount)
		automation.logIncentiveAction(validatorID, "Validator Rewarded", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	}
}

// penalizeValidator penalizes an inactive validator by reducing their rewards or imposing other penalties
func (automation *ValidatorIncentiveEnforcementAutomation) penalizeValidator(validatorID string) {
	penaltyAmount := automation.calculatePenalty(validatorID)
	err := automation.ledgerInstance.ApplyPenalty(validatorID, penaltyAmount)
	if err != nil {
		fmt.Printf("Failed to penalize validator %s: %v\n", validatorID, err)
		automation.logIncentiveAction(validatorID, "Penalty Failed", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	} else {
		fmt.Printf("Validator %s penalized with %f.\n", validatorID, penaltyAmount)
		automation.logIncentiveAction(validatorID, "Validator Penalized", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	}
}

// calculateReward calculates the reward amount based on validator activity and incentives
func (automation *ValidatorIncentiveEnforcementAutomation) calculateReward(validatorID string) float64 {
	baseReward := automation.validatorManager.GetBaseReward(validatorID)
	return baseReward * RewardMultiplier
}

// calculatePenalty calculates the penalty amount for an inactive validator
func (automation *ValidatorIncentiveEnforcementAutomation) calculatePenalty(validatorID string) float64 {
	basePenalty := automation.validatorManager.GetBaseReward(validatorID)
	return basePenalty * PenaltyMultiplier
}

// logIncentiveAction securely logs actions related to validator incentives and penalties
func (automation *ValidatorIncentiveEnforcementAutomation) logIncentiveAction(validatorID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Validator ID: %s, Details: %s", action, validatorID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-incentive-enforcement-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Incentive Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log incentive enforcement action for validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Incentive enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorIncentiveEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualReward allows administrators to manually reward a validator
func (automation *ValidatorIncentiveEnforcementAutomation) TriggerManualReward(validatorID string, amount float64) {
	fmt.Printf("Manually rewarding validator %s with %f.\n", validatorID, amount)
	err := automation.ledgerInstance.AddReward(validatorID, amount)
	if err != nil {
		fmt.Printf("Failed to manually reward validator %s: %v\n", validatorID, err)
		automation.logIncentiveAction(validatorID, "Manual Reward Failed", fmt.Sprintf("Amount: %f", amount))
	} else {
		automation.logIncentiveAction(validatorID, "Manual Rewarded", fmt.Sprintf("Amount: %f", amount))
	}
}
