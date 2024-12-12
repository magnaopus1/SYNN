package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	ConsensusParticipationCheckInterval = 5 * time.Second // Interval for checking consensus participation
	MaxParticipationsPerValidator       = 5000            // Maximum participation per validator in a given period
	ParticipationTimeWindow             = 24 * time.Hour  // Time window to count participations
)

// ConsensusParticipationRestrictionAutomation monitors and restricts participation in the consensus mechanism
type ConsensusParticipationRestrictionAutomation struct {
	consensusSystem            *consensus.SynnergyConsensus
	ledgerInstance             *ledger.Ledger
	stateMutex                 *sync.RWMutex
	validatorParticipationCount map[string]int // Tracks participation count per validator
}

// NewConsensusParticipationRestrictionAutomation initializes and returns an instance of ConsensusParticipationRestrictionAutomation
func NewConsensusParticipationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *ConsensusParticipationRestrictionAutomation {
	return &ConsensusParticipationRestrictionAutomation{
		consensusSystem:            consensusSystem,
		ledgerInstance:             ledgerInstance,
		stateMutex:                 stateMutex,
		validatorParticipationCount: make(map[string]int),
	}
}

// StartConsensusParticipationMonitoring starts continuous monitoring of consensus participation
func (automation *ConsensusParticipationRestrictionAutomation) StartConsensusParticipationMonitoring() {
	ticker := time.NewTicker(ConsensusParticipationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorConsensusParticipation()
		}
	}()
}

// monitorConsensusParticipation checks validator participation in the consensus and enforces participation limits
func (automation *ConsensusParticipationRestrictionAutomation) monitorConsensusParticipation() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent validator participation from Synnergy Consensus
	recentParticipations := automation.consensusSystem.GetRecentValidatorParticipations()

	for _, participation := range recentParticipations {
		// Validate participation limits
		if !automation.validateParticipationLimit(participation) {
			automation.flagParticipationViolation(participation, "Exceeded maximum participation in consensus within the time window")
		}
	}
}

// validateParticipationLimit checks if the validator has exceeded the participation limit within the time window
func (automation *ConsensusParticipationRestrictionAutomation) validateParticipationLimit(participation common.ValidatorParticipation) bool {
	currentCount := automation.validatorParticipationCount[participation.ValidatorID]
	if currentCount+1 > MaxParticipationsPerValidator {
		return false
	}

	// Update the participation count for the validator
	automation.validatorParticipationCount[participation.ValidatorID]++
	return true
}

// flagParticipationViolation flags a validator participation that violates system rules and logs it in the ledger
func (automation *ConsensusParticipationRestrictionAutomation) flagParticipationViolation(participation common.ValidatorParticipation, reason string) {
	fmt.Printf("Validator participation violation: Validator %s, Reason: %s\n", participation.ValidatorID, reason)

	// Log the violation into the ledger
	automation.logParticipationViolation(participation, reason)
}

// logParticipationViolation logs the flagged validator participation violation into the ledger with full details
func (automation *ConsensusParticipationRestrictionAutomation) logParticipationViolation(participation common.ValidatorParticipation, violationReason string) {
	// Encrypt the participation data
	encryptedData := automation.encryptParticipationData(participation)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("participation-violation-%s-%d", participation.ValidatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Consensus Participation Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Validator %s flagged for participation violation. Reason: %s. Encrypted Data: %s", participation.ValidatorID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log participation violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Participation violation logged for validator: %s\n", participation.ValidatorID)
	}
}

// encryptParticipationData encrypts validator participation data before logging for security
func (automation *ConsensusParticipationRestrictionAutomation) encryptParticipationData(participation common.ValidatorParticipation) string {
	data := fmt.Sprintf("Validator ID: %s, Timestamp: %d", participation.ValidatorID, participation.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting participation data:", err)
		return data
	}
	return string(encryptedData)
}
