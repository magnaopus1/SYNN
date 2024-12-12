package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/consensus"
)

// NewSidechainManager initializes the sidechain manager
func NewSidechainManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, sidechainNetwork *common.SidechainNetwork, consensus *common.SynnergyConsensus) *common.SidechainManager {
	return &common.SidechainManager{
		Sidechains:      make(map[string]*common.Sidechain),
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		SidechainNetwork: sidechainNetwork,
		Consensus:       consensus,
	}
}

// CreateSidechain creates a new sidechain and records the event in the ledger
func (sm *common.SidechainManager) CreateSidechain(chainID string, parentChainID string) (*common.Sidechain, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.Sidechains[chainID]; exists {
		return nil, errors.New("sidechain already exists")
	}

	// Initialize the new sidechain
	sidechain := NewSidechain(chainID, parentChainID, sm.Ledger, sm.Encryption, sm.Consensus)
	sm.Sidechains[chainID] = sidechain

	// Log the sidechain creation in the ledger
	err := sm.Ledger.RecordSidechainCreation(chainID, parentChainID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log sidechain creation: %v", err)
	}

	fmt.Printf("Sidechain %s created with parent chain %s\n", chainID, parentChainID)
	return sidechain, nil
}

// DeploySidechain deploys the sidechain to the network and records the event in the ledger
func (sm *common.SidechainManager) DeploySidechain(chainID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sidechain, exists := sm.Sidechains[chainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", chainID)
	}

	// Deploy sidechain nodes to the network
	err := sm.SidechainNetwork.DeploySidechain(sidechain)
	if err != nil {
		return fmt.Errorf("failed to deploy sidechain %s: %v", chainID, err)
	}

	// Log the deployment in the ledger
	err = sm.Ledger.RecordSidechainDeployment(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain deployment: %v", err)
	}

	fmt.Printf("Sidechain %s deployed to the network\n", chainID)
	return nil
}

// MonitorSidechain tracks the performance and health of a sidechain, including consensus and encryption
func (sm *common.SidechainManager) MonitorSidechain(chainID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sidechain, exists := sm.Sidechains[chainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", chainID)
	}

	// Perform health checks on the sidechain
	fmt.Printf("Monitoring sidechain %s...\n", chainID)

	// Check the consensus mechanism for this sidechain
	err := sm.Consensus.CheckIntegrity()
	if err != nil {
		return fmt.Errorf("consensus check failed for sidechain %s: %v", chainID, err)
	}

	// Check encryption services
	err = sm.Encryption.CheckEncryptionHealth()
	if err != nil {
		return fmt.Errorf("encryption check failed for sidechain %s: %v", chainID, err)
	}

	// Check node health and performance
	err = sm.SidechainNetwork.CheckSidechainNodes(chainID)
	if err != nil {
		return fmt.Errorf("sidechain node health check failed for sidechain %s: %v", chainID, err)
	}

	// Log the monitoring event
	err = sm.Ledger.RecordSidechainMonitoring(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain monitoring for sidechain %s: %v", chainID, err)
	}

	fmt.Printf("Sidechain %s monitoring completed successfully\n", chainID)
	return nil
}

// UpgradeSidechain handles upgrading a sidechain to a new version, ensuring network compatibility
func (sm *common.SidechainManager) UpgradeSidechain(chainID string, newVersion string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sidechain, exists := sm.Sidechains[chainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", chainID)
	}

	// Perform the upgrade process
	fmt.Printf("Upgrading sidechain %s to version %s...\n", chainID, newVersion)

	// Ensure that nodes in the sidechain are compatible with the new version
	err := sm.SidechainNetwork.UpgradeSidechainNodes(chainID, newVersion)
	if err != nil {
		return fmt.Errorf("failed to upgrade sidechain %s: %v", chainID, err)
	}

	// Log the upgrade in the ledger
	err = sm.Ledger.RecordSidechainUpgrade(chainID, newVersion, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain upgrade: %v", err)
	}

	fmt.Printf("Sidechain %s successfully upgraded to version %s\n", chainID, newVersion)
	return nil
}

// RemoveSidechain decommissions a sidechain and removes it from the network
func (sm *common.SidechainManager) RemoveSidechain(chainID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, exists := sm.Sidechains[chainID]
	if !exists {
		return fmt.Errorf("sidechain %s not found", chainID)
	}

	// Remove the sidechain from the network
	err := sm.SidechainNetwork.RemoveSidechain(chainID)
	if err != nil {
		return fmt.Errorf("failed to remove sidechain %s: %v", chainID, err)
	}

	// Log the removal of the sidechain in the ledger
	err = sm.Ledger.RecordSidechainRemoval(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain removal: %v", err)
	}

	// Delete sidechain from manager's internal list
	delete(sm.Sidechains, chainID)
	fmt.Printf("Sidechain %s removed from the network\n", chainID)
	return nil
}

// RetrieveSidechain retrieves a sidechain by its chain ID
func (sm *common.SidechainManager) RetrieveSidechain(chainID string) (*common.Sidechain, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sidechain, exists := sm.Sidechains[chainID]
	if !exists {
		return nil, fmt.Errorf("sidechain %s not found", chainID)
	}

	fmt.Printf("Retrieved sidechain %s\n", chainID)
	return sidechain, nil
}
