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

// Configuration for validator stake enforcement automation
const (
	StakeCheckInterval             = 30 * time.Second // Interval to check validator stakes
	MinimumRequiredStake           = 1000.0           // Minimum required stake for validators
	StakePenaltyMultiplier          = 0.8              // Penalty multiplier for validators below required stake
	StakeRewardMultiplier           = 1.2              // Reward multiplier for validators meeting or exceeding required stake
	StakeViolationThreshold         = 3                // Number of consecutive intervals allowed below minimum stake before penalty
)

// ValidatorStakeEnforcementAutomation monitors and enforces staking requirements for validators
type ValidatorStakeEnforcementAutomation struct {
	validatorManager      *validator.ValidatorManager
	consensusEngine       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	enforcementMutex      *sync.RWMutex
	stakeViolationCount   map[string]int // Tracks consecutive stake violations for each validator
}

// NewValidatorStakeEnforcementAutomation initializes the validator stake enforcement automation
func NewValidatorStakeEnforcementAutomation(validatorManager *validator.ValidatorManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ValidatorStakeEnforcementAutomation {
	return &ValidatorStakeEnforcementAutomation{
		validatorManager:    validatorManager,
		consensusEngine:     consensusEngine,
		ledgerInstance:      ledgerInstance,
		enforcementMutex:    enforcementMutex,
		stakeViolationCount: make(map[string]int),
	}
}

// StartValidatorStakeEnforcement begins continuous monitoring and enforcement of validator staking requirements
func (automation *ValidatorStakeEnforcementAutomation) StartValidatorStakeEnforcement() {
	ticker := time.NewTicker(StakeCheckInterval)

	go func() {
		for range ticker.C {
			automation.enforceStakeRequirements()
		}
	}()
}

// enforceStakeRequirements checks if validators meet staking requirements and rewards or penalizes as needed
func (automation *ValidatorStakeEnforcementAutomation) enforceStakeRequirements() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, validatorID := range automation.validatorManager.GetAllValidators() {
		stakeAmount := automation.validatorManager.GetValidatorStake(validatorID)

		if stakeAmount >= MinimumRequiredStake {
			automation.rewardValidatorForStake(validatorID)
			automation.stakeViolationCount[validatorID] = 0 // Reset violation count for sufficient stake
		} else {
			automation.stakeViolationCount[validatorID]++
			if automation.stakeViolationCount[validatorID] >= StakeViolationThreshold {
				automation.penalizeValidatorForLowStake(validatorID)
			}
		}
	}
}

// rewardValidatorForStake rewards a validator for maintaining the required stake or higher
func (automation *ValidatorStakeEnforcementAutomation) rewardValidatorForStake(validatorID string) {
	rewardAmount := automation.calculateStakeReward(validatorID)
	err := automation.ledgerInstance.AddReward(validatorID, rewardAmount)
	if err != nil {
		fmt.Printf("Failed to reward validator %s for stake: %v\n", validatorID, err)
		automation.logStakeEnforcementAction(validatorID, "Stake Reward Failed", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	} else {
		fmt.Printf("Validator %s rewarded with %f for maintaining sufficient stake.\n", validatorID, rewardAmount)
		automation.logStakeEnforcementAction(validatorID, "Validator Stake Rewarded", fmt.Sprintf("Reward Amount: %f", rewardAmount))
	}
}

// penalizeValidatorForLowStake penalizes a validator for not meeting the minimum required stake
func (automation *ValidatorStakeEnforcementAutomation) penalizeValidatorForLowStake(validatorID string) {
	penaltyAmount := automation.calculateStakePenalty(validatorID)
	err := automation.ledgerInstance.ApplyPenalty(validatorID, penaltyAmount)
	if err != nil {
		fmt.Printf("Failed to penalize validator %s for low stake: %v\n", validatorID, err)
		automation.logStakeEnforcementAction(validatorID, "Stake Penalty Failed", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	} else {
		fmt.Printf("Validator %s penalized with %f for insufficient stake.\n", validatorID, penaltyAmount)
		automation.logStakeEnforcementAction(validatorID, "Validator Stake Penalized", fmt.Sprintf("Penalty Amount: %f", penaltyAmount))
	}
}

// calculateStakeReward calculates the reward amount for maintaining the required stake level
func (automation *ValidatorStakeEnforcementAutomation) calculateStakeReward(validatorID string) float64 {
	baseReward := automation.validatorManager.GetBaseReward(validatorID)
	return baseReward * StakeRewardMultiplier
}

// calculateStakePenalty calculates the penalty amount for validators with low stake
func (automation *ValidatorStakeEnforcementAutomation) calculateStakePenalty(validatorID string) float64 {
	basePenalty := automation.validatorManager.GetBaseReward(validatorID)
	return basePenalty * StakePenaltyMultiplier
}

// logStakeEnforcementAction securely logs actions related to stake enforcement
func (automation *ValidatorStakeEnforcementAutomation) logStakeEnforcementAction(validatorID, action, details string) {
	entryDetails := fmt.Sprintf("Action: %s, Validator ID: %s, Details: %s", action, validatorID, details)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-stake-enforcement-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Stake Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log stake enforcement action for validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Stake enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorStakeEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualStakeReward allows administrators to manually reward a validator for maintaining high stake
func (automation *ValidatorStakeEnforcementAutomation) TriggerManualStakeReward(validatorID string, amount float64) {
	fmt.Printf("Manually rewarding validator %s with %f for high stake.\n", validatorID, amount)
	err := automation.ledgerInstance.AddReward(validatorID, amount)
	if err != nil {
		fmt.Printf("Failed to manually reward validator %s: %v\n", validatorID, err)
		automation.logStakeEnforcementAction(validatorID, "Manual Stake Reward Failed", fmt.Sprintf("Amount: %f", amount))
	} else {
		automation.logStakeEnforcementAction(validatorID, "Manual Stake Rewarded", fmt.Sprintf("Amount: %f", amount))
	}
}
