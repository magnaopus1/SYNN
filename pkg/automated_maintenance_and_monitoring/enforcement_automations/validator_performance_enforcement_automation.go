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

// Configuration for validator performance enforcement automation
const (
	PerformanceCheckInterval        = 20 * time.Second // Interval to check validator performance
	MinValidationAccuracy           = 0.98             // Minimum validation accuracy required
	MaxResponseTime                 = 1 * time.Second  // Maximum response time for transaction validation
	PerformanceViolationThreshold   = 2                // Allowed violations before penalty
	PerformanceRewardMultiplier     = 1.5              // Multiplier for rewards on high performance
	PerformancePenaltyMultiplier    = 0.7              // Multiplier for penalties on low performance
)

// ValidatorPerformanceEnforcementAutomation monitors and enforces performance standards for validators
type ValidatorPerformanceEnforcementAutomation struct {
	validatorManager      *validator.ValidatorManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	violationCountMap     map[string]int // Tracks performance violations for each validator
}

// NewValidatorPerformanceEnforcementAutomation initializes the validator performance enforcement automation
func NewValidatorPerformanceEnforcementAutomation(validatorManager *validator.ValidatorManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ValidatorPerformanceEnforcementAutomation {
	return &ValidatorPerformanceEnforcementAutomation{
		validatorManager:   validatorManager,
		consensusEngine:    consensusEngine,
		ledgerInstance:     ledgerInstance,
		enforcementMutex:   enforcementMutex,
		violationCountMap:  make(map[string]int),
	}
}

// StartValidatorPerformanceEnforcement begins continuous monitoring and enforcement of validator performance
func (automation *ValidatorPerformanceEnforcementAutomation) StartValidatorPerformanceEnforcement() {
	ticker := time.NewTicker(PerformanceCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateValidatorPerformance()
		}
	}()
}

// evaluateValidatorPerformance checks validator accuracy and response times, enforcing rewards or penalties
func (automation *ValidatorPerformanceEnforcementAutomation) evaluateValidatorPerformance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, validatorID := range automation.validatorManager.GetAllValidators() {
		accuracy := automation.validatorManager.GetValidationAccuracy(validatorID)
		responseTime := automation.validatorManager.GetResponseTime(validatorID)

		if accuracy >= MinValidationAccuracy && responseTime <= MaxResponseTime {
			automation.rewardValidator(validatorID)
			automation.violationCountMap[validatorID] = 0 // Reset violation count on reward
		} else {
			automation.violationCountMap[validatorID]++
			if automation.violationCountMap[validatorID] >= PerformanceViolationThreshold {
				automation.penalizeValidator(validatorID)
			}
		}
	}
}

// rewardValidator rewards a high-performing validator based on accuracy and response times
func (automation *ValidatorPerformanceEnforcementAutomation) rewardValidator(validatorID string) {
	rewardAmount := automation.calculateReward(validatorID)
	err := automation.ledgerInstance.AddReward(validatorID, rewardAmount)
	if err != nil {
		fmt.Printf("Failed to reward validator %s: %v\n", validatorID, err)
		automation.logPerformanceAction(validatorID, "Reward Failed", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	} else {
		fmt.Printf("Validator %s rewarded with %f for high performance.\n", validatorID, rewardAmount)
		automation.logPerformanceAction(validatorID, "Validator Rewarded", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	}
}

// penalizeValidator penalizes a validator that does not meet performance standards
func (automation *ValidatorPerformanceEnforcementAutomation) penalizeValidator(validatorID string) {
	penaltyAmount := automation.calculatePenalty(validatorID)
	err := automation.ledgerInstance.ApplyPenalty(validatorID, penaltyAmount)
	if err != nil {
		fmt.Printf("Failed to penalize validator %s: %v\n", validatorID, err)
		automation.logPerformanceAction(validatorID, "Penalty Failed", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	} else {
		fmt.Printf("Validator %s penalized with %f for poor performance.\n", validatorID, penaltyAmount)
		automation.logPerformanceAction(validatorID, "Validator Penalized", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	}
}

// calculateReward calculates the reward amount based on high performance metrics
func (automation *ValidatorPerformanceEnforcementAutomation) calculateReward(validatorID string) float64 {
	baseReward := automation.validatorManager.GetBaseReward(validatorID)
	return baseReward * PerformanceRewardMultiplier
}

// calculatePenalty calculates the penalty amount for underperforming validators
func (automation *ValidatorPerformanceEnforcementAutomation) calculatePenalty(validatorID string) float64 {
	basePenalty := automation.validatorManager.GetBaseReward(validatorID)
	return basePenalty * PerformancePenaltyMultiplier
}

// logPerformanceAction securely logs actions related to performance enforcement
func (automation *ValidatorPerformanceEnforcementAutomation) logPerformanceAction(validatorID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Validator ID: %s, Details: %s", action, validatorID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-performance-enforcement-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Performance Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log performance enforcement action for validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Performance enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorPerformanceEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualPerformanceReward allows administrators to manually reward a validator for high performance
func (automation *ValidatorPerformanceEnforcementAutomation) TriggerManualPerformanceReward(validatorID string, amount float64) {
	fmt.Printf("Manually rewarding validator %s with %f for outstanding performance.\n", validatorID, amount)
	err := automation.ledgerInstance.AddReward(validatorID, amount)
	if err != nil {
		fmt.Printf("Failed to manually reward validator %s: %v\n", validatorID, err)
		automation.logPerformanceAction(validatorID, "Manual Reward Failed", fmt.Sprintf("Amount: %f", amount))
	} else {
		automation.logPerformanceAction(validatorID, "Manual Rewarded", fmt.Sprintf("Amount: %f", amount))
	}
}
