package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
)

const (
	LoadCheckInterval           = 5 * time.Minute // Interval for checking network load conditions
	LoadDistributionThreshold   = 0.75            // Threshold for load balancing, when 75% of resources are used
	MaxNodeCapacity             = 10000           // Example maximum node transaction capacity
	LoadLoggingInterval         = 15 * time.Minute // Interval for logging load status to the ledger
)

// DynamicLoadDistributionAutomation handles dynamic load balancing across nodes in the network
type DynamicLoadDistributionAutomation struct {
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation and integration
	ledgerInstance  *ledger.Ledger                 // Ledger to track load distribution events
	loadMutex       *sync.RWMutex                  // Mutex for thread-safe load management
}

// NewDynamicLoadDistributionAutomation initializes the dynamic load distribution automation
func NewDynamicLoadDistributionAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, loadMutex *sync.RWMutex) *DynamicLoadDistributionAutomation {
	return &DynamicLoadDistributionAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		loadMutex:       loadMutex,
	}
}

// StartLoadMonitoring begins the automatic load monitoring and distribution process
func (automation *DynamicLoadDistributionAutomation) StartLoadMonitoring() {
	ticker := time.NewTicker(LoadCheckInterval)
	go func() {
		for range ticker.C {
			automation.evaluateNetworkLoad()
		}
	}()

	// Periodically log load distribution status
	logTicker := time.NewTicker(LoadLoggingInterval)
	go func() {
		for range logTicker.C {
			automation.logLoadStatus()
		}
	}()
}

// evaluateNetworkLoad checks the current network load and triggers load redistribution if necessary
func (automation *DynamicLoadDistributionAutomation) evaluateNetworkLoad() {
	automation.loadMutex.Lock()
	defer automation.loadMutex.Unlock()

	// Fetch the current load from the consensus engine
	networkLoad := automation.consensusEngine.GetNetworkLoad()

	fmt.Printf("Current network load: %.2f\n", networkLoad)

	// Trigger load distribution if the load exceeds the threshold
	if networkLoad >= LoadDistributionThreshold {
		automation.distributeNetworkLoad()
	} else {
		fmt.Println("Network load is within optimal limits.")
	}
}

// distributeNetworkLoad redistributes the load across available nodes
func (automation *DynamicLoadDistributionAutomation) distributeNetworkLoad() {
	fmt.Println("Distributing network load to balance the system...")

	// Fetch the nodes and their capacities from the consensus engine
	nodes := automation.consensusEngine.GetNodeCapacities()

	// Distribute excess load to less loaded nodes
	for nodeID, nodeCapacity := range nodes {
		if nodeCapacity < MaxNodeCapacity {
			excessLoad := (LoadDistributionThreshold * MaxNodeCapacity) - nodeCapacity
			err := automation.consensusEngine.RedistributeLoad(nodeID, excessLoad)
			if err != nil {
				fmt.Printf("Failed to redistribute load to node %s: %v\n", nodeID, err)
			} else {
				fmt.Printf("Successfully redistributed load to node %s\n", nodeID)
			}
		}
	}

	// Log the load distribution event in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("load-distribution-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Load Distribution",
		Status:    "Success",
		Details:   "Network load redistributed to optimize system performance.",
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log load distribution:", err)
	} else {
		fmt.Println("Load distribution successfully logged in the ledger.")
	}
}

// logLoadStatus logs the current status of network load into the ledger
func (automation *DynamicLoadDistributionAutomation) logLoadStatus() {
	networkLoad := automation.consensusEngine.GetNetworkLoad()

	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("network-load-status-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Network Load Status",
		Status:    "Success",
		Details:   fmt.Sprintf("Current network load: %.2f", networkLoad),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log network load status:", err)
	} else {
		fmt.Println("Network load status successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DynamicLoadDistributionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualLoadDistribution allows administrators to manually trigger load distribution
func (automation *DynamicLoadDistributionAutomation) TriggerManualLoadDistribution() {
	fmt.Println("Manually triggering load distribution...")

	// Trigger load redistribution across nodes
	automation.distributeNetworkLoad()
}
