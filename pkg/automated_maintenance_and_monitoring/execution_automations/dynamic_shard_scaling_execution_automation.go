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
	ShardCheckInterval         = 10 * time.Minute // Interval for checking if shard scaling is needed
	ShardLoadThreshold         = 0.75             // Threshold for triggering shard scaling
	ShardMaxCapacity           = 10000            // Maximum shard capacity before triggering scaling
	ShardMinCapacity           = 5000             // Minimum shard capacity before triggering downscaling
	ShardScalingLoggingInterval = 15 * time.Minute // Interval for logging shard scaling status to the ledger
)

// DynamicShardScalingAutomation handles the automatic dynamic scaling of shards in the network
type DynamicShardScalingAutomation struct {
	consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for shard management
	ledgerInstance  *ledger.Ledger                        // Ledger to track shard scaling events
	shardMutex      *sync.RWMutex                         // Mutex for thread-safe shard scaling
}

// NewDynamicShardScalingAutomation initializes the dynamic shard scaling automation
func NewDynamicShardScalingAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, shardMutex *sync.RWMutex) *DynamicShardScalingAutomation {
	return &DynamicShardScalingAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		shardMutex:      shardMutex,
	}
}

// StartShardScalingMonitoring begins the automatic shard scaling process
func (automation *DynamicShardScalingAutomation) StartShardScalingMonitoring() {
	ticker := time.NewTicker(ShardCheckInterval)
	go func() {
		for range ticker.C {
			automation.evaluateShardLoad()
		}
	}()

	// Periodically log shard scaling status
	logTicker := time.NewTicker(ShardScalingLoggingInterval)
	go func() {
		for range logTicker.C {
			automation.logShardScalingStatus()
		}
	}()
}

// evaluateShardLoad checks the current shard load and triggers scaling if necessary
func (automation *DynamicShardScalingAutomation) evaluateShardLoad() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	// Fetch the current shard load from the consensus engine
	shardLoads := automation.consensusEngine.GetShardLoads()

	for shardID, shardLoad := range shardLoads {
		fmt.Printf("Shard %s current load: %.2f\n", shardID, shardLoad)

		if shardLoad >= ShardLoadThreshold {
			automation.scaleShard(shardID)
		} else if shardLoad < (ShardLoadThreshold / 2) && shardLoad > ShardMinCapacity {
			automation.downscaleShard(shardID)
		} else {
			fmt.Printf("Shard %s load is within optimal limits.\n", shardID)
		}
	}
}

// scaleShard scales up the specified shard by splitting it
func (automation *DynamicShardScalingAutomation) scaleShard(shardID string) {
	fmt.Printf("Scaling up shard %s...\n", shardID)

	// Trigger shard splitting in the consensus engine
	err := automation.consensusEngine.SplitShard(shardID)
	if err != nil {
		fmt.Printf("Failed to scale shard %s: %v\n", shardID, err)
		return
	}

	// Log the shard scaling event in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-scaling-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Scaling",
		Status:    "Success",
		Details:   fmt.Sprintf("Shard %s was successfully scaled due to high load.", shardID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log shard scaling event for shard %s: %v\n", shardID, err)
	} else {
		fmt.Printf("Shard scaling event for shard %s successfully logged in the ledger.\n", shardID)
	}
}

// downscaleShard scales down the specified shard by merging it with another shard
func (automation *DynamicShardScalingAutomation) downscaleShard(shardID string) {
	fmt.Printf("Downscaling shard %s...\n", shardID)

	// Trigger shard merging in the consensus engine
	err := automation.consensusEngine.MergeShard(shardID)
	if err != nil {
		fmt.Printf("Failed to downscale shard %s: %v\n", shardID, err)
		return
	}

	// Log the shard downscaling event in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("shard-downscaling-%s-%d", shardID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Shard Downscaling",
		Status:    "Success",
		Details:   fmt.Sprintf("Shard %s was successfully downscaled due to low load.", shardID),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err = automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log shard downscaling event for shard %s: %v\n", shardID, err)
	} else {
		fmt.Printf("Shard downscaling event for shard %s successfully logged in the ledger.\n", shardID)
	}
}

// logShardScalingStatus logs the current status of shard loads into the ledger
func (automation *DynamicShardScalingAutomation) logShardScalingStatus() {
	shardLoads := automation.consensusEngine.GetShardLoads()

	for shardID, shardLoad := range shardLoads {
		entry := common.LedgerEntry{
			ID:        fmt.Sprintf("shard-load-status-%s-%d", shardID, time.Now().Unix()),
			Timestamp: time.Now().Unix(),
			Type:      "Shard Load Status",
			Status:    "Success",
			Details:   fmt.Sprintf("Current load for shard %s: %.2f", shardID, shardLoad),
		}

		encryptedDetails := automation.encryptData(entry.Details)
		entry.Details = encryptedDetails

		err := automation.ledgerInstance.AddEntry(entry)
		if err != nil {
			fmt.Printf("Failed to log load status for shard %s: %v\n", shardID, err)
		} else {
			fmt.Printf("Shard load status for shard %s successfully logged in the ledger.\n", shardID)
		}
	}
}

// encryptData encrypts sensitive data before logging it into the ledger
func (automation *DynamicShardScalingAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}

// TriggerManualShardScaling allows administrators to manually trigger shard scaling
func (automation *DynamicShardScalingAutomation) TriggerManualShardScaling(shardID string) {
	fmt.Printf("Manually triggering scaling for shard %s...\n", shardID)
	automation.scaleShard(shardID)
}

// TriggerManualShardDownscaling allows administrators to manually trigger shard downscaling
func (automation *DynamicShardScalingAutomation) TriggerManualShardDownscaling(shardID string) {
	fmt.Printf("Manually triggering downscaling for shard %s...\n", shardID)
	automation.downscaleShard(shardID)
}
