package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/validators"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	ValidatorRotationInterval = 30 * time.Minute // Interval for rotating validators
	MinValidatorUptime        = 0.75             // Minimum uptime threshold for validators to remain active
	RotationBatchSize         = 5                // Number of validators rotated at each interval
)

// ValidatorRotationAutomation manages periodic rotation of validators in the Synnergy Consensus
type ValidatorRotationAutomation struct {
	consensusEngine  *consensus.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	validatorManager *validators.ValidatorManager
	rotationMutex    *sync.RWMutex
}

// NewValidatorRotationAutomation initializes the validator rotation automation
func NewValidatorRotationAutomation(consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, validatorManager *validators.ValidatorManager, rotationMutex *sync.RWMutex) *ValidatorRotationAutomation {
	return &ValidatorRotationAutomation{
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		validatorManager: validatorManager,
		rotationMutex:    rotationMutex,
	}
}

// StartValidatorRotationMonitor starts the validator rotation automation in a continuous loop
func (automation *ValidatorRotationAutomation) StartValidatorRotationMonitor() {
	ticker := time.NewTicker(ValidatorRotationInterval)

	go func() {
		for range ticker.C {
			automation.performValidatorRotation()
		}
	}()
}

// performValidatorRotation triggers the rotation of validators based on uptime and other performance metrics
func (automation *ValidatorRotationAutomation) performValidatorRotation() {
	automation.rotationMutex.Lock()
	defer automation.rotationMutex.Unlock()

	// Fetch current active validators and their uptime
	activeValidators := automation.validatorManager.GetActiveValidators()
	validatorsToRotate := automation.getValidatorsToRotate(activeValidators)

	if len(validatorsToRotate) > 0 {
		automation.rotateValidators(validatorsToRotate)
	} else {
		fmt.Println("No validators require rotation.")
	}
}

// getValidatorsToRotate filters the list of validators to find those with low uptime for rotation
func (automation *ValidatorRotationAutomation) getValidatorsToRotate(validators []validators.Validator) []validators.Validator {
	var validatorsToRotate []validators.Validator

	for _, validator := range validators {
		uptime := automation.validatorManager.GetValidatorUptime(validator.ID)
		if uptime < MinValidatorUptime {
			validatorsToRotate = append(validatorsToRotate, validator)
		}
	}

	// Limit the number of validators rotated to the batch size
	if len(validatorsToRotate) > RotationBatchSize {
		validatorsToRotate = validatorsToRotate[:RotationBatchSize]
	}

	return validatorsToRotate
}

// rotateValidators performs the rotation of selected validators
func (automation *ValidatorRotationAutomation) rotateValidators(validatorsToRotate []validators.Validator) {
	for _, validator := range validatorsToRotate {
		fmt.Printf("Rotating Validator %s due to low uptime.\n", validator.ID)
		automation.validatorManager.RotateValidator(validator.ID)
		automation.logRotationEvent(validator.ID)
	}
}

// logRotationEvent logs the validator rotation event in the ledger for transparency
func (automation *ValidatorRotationAutomation) logRotationEvent(validatorID string) {
	entryDetails := fmt.Sprintf("Validator %s rotated due to low uptime.", validatorID)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := ledger.LedgerEntry{
		ID:        fmt.Sprintf("validator-rotation-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Rotation",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log validator rotation in the ledger for Validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Validator rotation successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ValidatorRotationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualValidatorRotation allows administrators to manually trigger the rotation of specific validators
func (automation *ValidatorRotationAutomation) TriggerManualValidatorRotation(validatorID string) {
	fmt.Printf("Manually rotating Validator %s.\n", validatorID)
	automation.validatorManager.RotateValidator(validatorID)
	automation.logRotationEvent(validatorID)
}
