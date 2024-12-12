package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/contracts"
	"synnergy_network_demo/resources"
)

// Configuration for contract resource limit enforcement automation
const (
	ResourceCheckInterval          = 10 * time.Second // Interval to check resource usage of deployed contracts
	CPUUsageThreshold              = 80              // CPU usage percentage threshold
	MemoryUsageThreshold           = 1024            // Memory usage limit in MB
	StorageUsageThreshold          = 5000            // Storage usage limit in MB
	MaxResourceViolations          = 3               // Max violations before restricting contract
)

// ContractResourceLimitEnforcementAutomation monitors and enforces resource usage limits for deployed contracts
type ContractResourceLimitEnforcementAutomation struct {
	resourceManager   *resources.ResourceManager
	contractManager   *contracts.ContractManager
	consensusEngine   *consensus.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	enforcementMutex  *sync.RWMutex
	violationCount    map[string]int // Tracks resource violations per contract
}

// NewContractResourceLimitEnforcementAutomation initializes the resource limit enforcement automation
func NewContractResourceLimitEnforcementAutomation(resourceManager *resources.ResourceManager, contractManager *contracts.ContractManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *ContractResourceLimitEnforcementAutomation {
	return &ContractResourceLimitEnforcementAutomation{
		resourceManager:  resourceManager,
		contractManager:  contractManager,
		consensusEngine:  consensusEngine,
		ledgerInstance:   ledgerInstance,
		enforcementMutex: enforcementMutex,
		violationCount:   make(map[string]int),
	}
}

// StartResourceLimitEnforcement begins continuous monitoring and enforcement of resource limits
func (automation *ContractResourceLimitEnforcementAutomation) StartResourceLimitEnforcement() {
	ticker := time.NewTicker(ResourceCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkResourceLimits()
		}
	}()
}

// checkResourceLimits monitors each deployed contract's resource usage and enforces limits if thresholds are exceeded
func (automation *ContractResourceLimitEnforcementAutomation) checkResourceLimits() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	for _, contractID := range automation.contractManager.GetDeployedContracts() {
		cpuUsage := automation.resourceManager.GetCPUUsage(contractID)
		memoryUsage := automation.resourceManager.GetMemoryUsage(contractID)
		storageUsage := automation.resourceManager.GetStorageUsage(contractID)

		if cpuUsage > CPUUsageThreshold || memoryUsage > MemoryUsageThreshold || storageUsage > StorageUsageThreshold {
			automation.handleResourceViolation(contractID, cpuUsage, memoryUsage, storageUsage)
		}
	}
}

// handleResourceViolation restricts contracts that exceed resource limits
func (automation *ContractResourceLimitEnforcementAutomation) handleResourceViolation(contractID string, cpuUsage, memoryUsage, storageUsage int) {
	automation.violationCount[contractID]++

	if automation.violationCount[contractID] >= MaxResourceViolations {
		err := automation.contractManager.RestrictContract(contractID)
		if err != nil {
			fmt.Printf("Failed to restrict contract %s for resource violations: %v\n", contractID, err)
			automation.logResourceAction(contractID, "Failed Resource Restriction", cpuUsage, memoryUsage, storageUsage)
		} else {
			fmt.Printf("Contract %s restricted due to repeated resource limit violations.\n", contractID)
			automation.logResourceAction(contractID, "Contract Restricted for Resource Violations", cpuUsage, memoryUsage, storageUsage)
			automation.violationCount[contractID] = 0
		}
	} else {
		fmt.Printf("Resource violation detected for contract %s (CPU: %d%%, Memory: %dMB, Storage: %dMB).\n", contractID, cpuUsage, memoryUsage, storageUsage)
		automation.logResourceAction(contractID, "Resource Violation Detected", cpuUsage, memoryUsage, storageUsage)
	}
}

// logResourceAction securely logs actions related to resource limit enforcement
func (automation *ContractResourceLimitEnforcementAutomation) logResourceAction(contractID, action string, cpuUsage, memoryUsage, storageUsage int) {
	entryDetails := fmt.Sprintf("Action: %s, Contract ID: %s, CPU Usage: %d%%, Memory Usage: %dMB, Storage Usage: %dMB", action, contractID, cpuUsage, memoryUsage, storageUsage)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("resource-enforcement-%s-%d", contractID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Resource Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log resource enforcement action for contract %s: %v\n", contractID, err)
	} else {
		fmt.Println("Resource enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *ContractResourceLimitEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualResourceRestriction allows administrators to manually restrict a contract for excessive resource usage
func (automation *ContractResourceLimitEnforcementAutomation) TriggerManualResourceRestriction(contractID string) {
	fmt.Printf("Manually triggering resource restriction for contract: %s\n", contractID)

	err := automation.contractManager.RestrictContract(contractID)
	if err != nil {
		fmt.Printf("Failed to manually restrict resources for contract %s: %v\n", contractID, err)
		automation.logResourceAction(contractID, "Manual Resource Restriction Failed", 0, 0, 0)
	} else {
		fmt.Printf("Manual resource restriction applied to contract %s.\n", contractID)
		automation.logResourceAction(contractID, "Manual Resource Restriction", 0, 0, 0)
	}
}
