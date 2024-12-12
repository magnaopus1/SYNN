package layer2_consensus

import (
	"crypto/rand"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// NewCrossConsensusScalingManager initializes the Cross-Consensus Scaling Manager
func NewCrossConsensusScalingManager(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *CrossConsensusScalingManager {
	return &CrossConsensusScalingManager{
		ConsensusMechanisms: make(map[string]*ConsensusMechanism),
		Ledger:              ledgerInstance,
		EncryptionService:   encryptionService,
	}
}

// AddConsensusMechanism adds a new consensus mechanism to the network
func (ccsm *CrossConsensusScalingManager) AddConsensusMechanism(mechanismID, mechanismType string) (*ConsensusMechanism, error) {
	ccsm.mu.Lock()
	defer ccsm.mu.Unlock()

	// Encrypt mechanism data
	mechanismData := fmt.Sprintf("MechanismID: %s, Type: %s", mechanismID, mechanismType)
	iv := make([]byte, 16) // Assuming a 16-byte initialization vector is needed
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	encryptedData, err := ccsm.EncryptionService.EncryptData(mechanismData, common.EncryptionKey, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt consensus mechanism data: %v", err)
	}

	// You should use encryptedData for something here (e.g., storing it), or at least log it for now
	fmt.Printf("Encrypted mechanism data: %x\n", encryptedData)

	// Create the consensus mechanism
	mechanism := &ConsensusMechanism{
		MechanismID:    mechanismID,
		MechanismType:  mechanismType,
		CurrentLoad:    0,
		TransitionCount: 0,
		Active:         false,
		LastTransition: time.Now(),
	}

	// Add the mechanism to the manager
	ccsm.ConsensusMechanisms[mechanismID] = mechanism

	// Log the addition of the consensus mechanism in the ledger
	ccsm.Ledger.BlockchainConsensusCoinLedger.RecordConsensusMechanismAddition(mechanismID, mechanismType)

	fmt.Printf("Consensus mechanism %s of type %s added\n", mechanismID, mechanismType)
	return mechanism, nil
}



// ActivateConsensusMechanism activates a specific consensus mechanism and transitions from the current one
func (ccsm *CrossConsensusScalingManager) ActivateConsensusMechanism(mechanismID string) error {
	ccsm.mu.Lock()
	defer ccsm.mu.Unlock()

	// Retrieve the consensus mechanism
	mechanism, exists := ccsm.ConsensusMechanisms[mechanismID]
	if !exists {
		return fmt.Errorf("consensus mechanism %s not found", mechanismID)
	}

	// Deactivate the current mechanism if one is active
	if ccsm.ActiveMechanism != nil {
		ccsm.ActiveMechanism.Active = false
	}

	// Activate the new mechanism
	mechanism.Active = true
	mechanism.TransitionCount++
	mechanism.LastTransition = time.Now()
	ccsm.ActiveMechanism = mechanism

	// Log the consensus mechanism transition in the ledger with the correct number of arguments
	ccsm.Ledger.BlockchainConsensusCoinLedger.RecordConsensusTransition(mechanismID, mechanism.MechanismType)

	fmt.Printf("Consensus mechanism %s is now active\n", mechanismID)
	return nil
}



// MonitorMechanismLoad monitors the load on a specific consensus mechanism and triggers a transition if necessary
func (ccsm *CrossConsensusScalingManager) MonitorMechanismLoad(mechanismID string, currentLoad float64) error {
	ccsm.mu.Lock()
	defer ccsm.mu.Unlock()

	// Retrieve the consensus mechanism
	mechanism, exists := ccsm.ConsensusMechanisms[mechanismID]
	if !exists {
		return fmt.Errorf("consensus mechanism %s not found", mechanismID)
	}

	// Update the mechanism load
	mechanism.CurrentLoad = currentLoad

	// Check if load exceeds a threshold and trigger a transition if necessary
	if mechanism.CurrentLoad > 0.75 && ccsm.ActiveMechanism.MechanismID != mechanismID {
		fmt.Printf("Triggering transition to %s due to high load\n", mechanismID)
		return ccsm.ActivateConsensusMechanism(mechanismID)
	}

	return nil
}

// GetActiveConsensusMechanism returns the currently active consensus mechanism
func (ccsm *CrossConsensusScalingManager) GetActiveConsensusMechanism() (*ConsensusMechanism, error) {
	ccsm.mu.Lock()
	defer ccsm.mu.Unlock()

	if ccsm.ActiveMechanism == nil {
		return nil, errors.New("no active consensus mechanism")
	}

	return ccsm.ActiveMechanism, nil
}

