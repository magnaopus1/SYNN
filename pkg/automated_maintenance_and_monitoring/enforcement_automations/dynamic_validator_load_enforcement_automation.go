package enforcement_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
)

// Configuration for dynamic validator load enforcement automation
const (
	LoadCheckInterval          = 10 * time.Second // Interval to check for validator load balancing
	MaxValidationLoad          = 80               // Maximum validation load percentage per validator
	MaxProcessingTimePerBlock  = 5 * time.Second  // Maximum processing time allowed per block per validator
)

// DynamicValidatorLoadEnforcementAutomation monitors and enforces load balancing among validator nodes
type DynamicValidatorLoadEnforcementAutomation struct {
	networkManager       *network.NetworkManager
	consensusEngine      *consensus.SynnergyConsensus
	ledgerInstance       *ledger.Ledger
	enforcementMutex     *sync.RWMutex
	validatorLoadMap     map[string]int // Tracks load percentage per validator node
	blockProcessingTime  map[string]time.Duration // Tracks processing time per block per validator node
}

// NewDynamicValidatorLoadEnforcementAutomation initializes the dynamic validator load enforcement automation
func NewDynamicValidatorLoadEnforcementAutomation(networkManager *network.NetworkManager, consensusEngine *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, enforcementMutex *sync.RWMutex) *DynamicValidatorLoadEnforcementAutomation {
	return &DynamicValidatorLoadEnforcementAutomation{
		networkManager:       networkManager,
		consensusEngine:      consensusEngine,
		ledgerInstance:       ledgerInstance,
		enforcementMutex:     enforcementMutex,
		validatorLoadMap:     make(map[string]int),
		blockProcessingTime:  make(map[string]time.Duration),
	}
}

// StartLoadEnforcement begins continuous monitoring and enforcement of validator load balancing
func (automation *DynamicValidatorLoadEnforcementAutomation) StartLoadEnforcement() {
	ticker := time.NewTicker(LoadCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkValidatorLoadCompliance()
		}
	}()
}

// checkValidatorLoadCompliance monitors validator node load and processing times to enforce load balancing
func (automation *DynamicValidatorLoadEnforcementAutomation) checkValidatorLoadCompliance() {
	automation.enforcementMutex.Lock()
	defer automation.enforcementMutex.Unlock()

	automation.adjustValidatorLoad()
	automation.checkProcessingTimeCompliance()
}

// adjustValidatorLoad dynamically reallocates tasks to validators based on current load
func (automation *DynamicValidatorLoadEnforcementAutomation) adjustValidatorLoad() {
	for _, validatorID := range automation.networkManager.GetValidators() {
		load := automation.networkManager.GetValidatorLoad(validatorID)

		if load > MaxValidationLoad {
			fmt.Printf("Load balancing required for validator %s.\n", validatorID)
			automation.reallocateTasks(validatorID, "Excessive Load")
		}
	}
}

// checkProcessingTimeCompliance ensures validators are processing blocks within the allowable time limits
func (automation *DynamicValidatorLoadEnforcementAutomation) checkProcessingTimeCompliance() {
	for _, validatorID := range automation.networkManager.GetValidators() {
		processingTime := automation.consensusEngine.GetBlockProcessingTime(validatorID)

		if processingTime > MaxProcessingTimePerBlock {
			fmt.Printf("Processing time violation detected for validator %s.\n", validatorID)
			automation.reallocateTasks(validatorID, "Processing Time Exceeded")
		}
	}
}

// reallocateTasks redistributes validation tasks from overloaded validators to maintain balanced load
func (automation *DynamicValidatorLoadEnforcementAutomation) reallocateTasks(validatorID, reason string) {
	err := automation.networkManager.ReallocateValidatorTasks(validatorID)
	if err != nil {
		fmt.Printf("Failed to reallocate tasks for validator %s: %v\n", validatorID, err)
		automation.logLoadAction(validatorID, "Task Reallocation Failed", reason)
	} else {
		fmt.Printf("Tasks reallocated from validator %s due to %s.\n", validatorID, reason)
		automation.logLoadAction(validatorID, "Tasks Reallocated", reason)
	}
}

// logLoadAction securely logs actions related to validator load balancing
func (automation *DynamicValidatorLoadEnforcementAutomation) logLoadAction(validatorID, action, reason string) {
	entryDetails := fmt.Sprintf("Action: %s, Validator ID: %s, Reason: %s", action, validatorID, reason)
	encryptedDetails := automation.encryptData(entryDetails)

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("validator-load-enforcement-%s-%d", validatorID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Load Enforcement",
		Status:    "Completed",
		Details:   encryptedDetails,
	}

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log load enforcement action for validator %s: %v\n", validatorID, err)
	} else {
		fmt.Println("Load enforcement action successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DynamicValidatorLoadEnforcementAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLoadCheck allows administrators to manually check and reallocate tasks for a specific validator
func (automation *DynamicValidatorLoadEnforcementAutomation) TriggerManualLoadCheck(validatorID string) {
	fmt.Printf("Manually triggering load check for validator: %s\n", validatorID)

	load := automation.networkManager.GetValidatorLoad(validatorID)
	processingTime := automation.consensusEngine.GetBlockProcessingTime(validatorID)

	if load > MaxValidationLoad {
		automation.reallocateTasks(validatorID, "Manual Trigger: Excessive Load")
	} else if processingTime > MaxProcessingTimePerBlock {
		automation.reallocateTasks(validatorID, "Manual Trigger: Processing Time Exceeded")
	} else {
		fmt.Printf("Validator %s is compliant with load balancing policies.\n", validatorID)
		automation.logLoadAction(validatorID, "Manual Compliance Check Passed", "Load Balancing Verified")
	}
}
