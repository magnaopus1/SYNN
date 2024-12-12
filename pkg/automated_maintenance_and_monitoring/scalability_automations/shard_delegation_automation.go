package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/common"
)

const (
	DelegationInterval         = 30 * time.Second  // Interval for shard delegation checks
	ShardDelegationLogKey      = "delegationLogKey" // Encryption key for shard delegation logs
	MaxValidatorsPerShard      = 100               // Maximum validators allowed per shard
)

// ShardDelegationAutomation manages the dynamic delegation of validators to shards based on load and requirements.
type ShardDelegationAutomation struct {
	ledgerInstance  *ledger.Ledger               // Reference to the ledger for logging
	consensusSystem *consensus.SynnergyConsensus // Reference to Synnergy Consensus
	shardMutex      *sync.RWMutex                // Mutex for concurrency control
}

// NewShardDelegationAutomation creates a new instance of ShardDelegationAutomation.
func NewShardDelegationAutomation(ledgerInstance *ledger.Ledger, consensusSystem *consensus.SynnergyConsensus, shardMutex *sync.RWMutex) *ShardDelegationAutomation {
	return &ShardDelegationAutomation{
		ledgerInstance:  ledgerInstance,
		consensusSystem: consensusSystem,
		shardMutex:      shardMutex,
	}
}

// StartShardDelegationAutomation starts continuous shard delegation monitoring and balancing.
func (automation *ShardDelegationAutomation) StartShardDelegationAutomation() {
	ticker := time.NewTicker(DelegationInterval)
	go func() {
		for range ticker.C {
			automation.checkAndDelegateValidators()
		}
	}()
}

// checkAndDelegateValidators checks the current state of shard validators and delegates them based on network load and capacity.
func (automation *ShardDelegationAutomation) checkAndDelegateValidators() {
	automation.shardMutex.Lock()
	defer automation.shardMutex.Unlock()

	shardList, err := automation.consensusSystem.GetAllShards()
	if err != nil {
		fmt.Printf("Error fetching shards: %v\n", err)
		return
	}

	for _, shard := range shardList {
		validatorCount, err := automation.consensusSystem.GetValidatorCount(shard.ID)
		if err != nil {
			fmt.Printf("Error fetching validator count for shard %s: %v\n", shard.ID, err)
			continue
		}

		if validatorCount < MaxValidatorsPerShard {
			// Delegate additional validators if the shard is underutilized
			automation.delegateValidatorsToShard(shard.ID, MaxValidatorsPerShard-validatorCount)
		}
	}
}

// delegateValidatorsToShard delegates additional validators to the specified shard to balance the network load.
func (automation *ShardDelegationAutomation) delegateValidatorsToShard(shardID string, requiredValidators int) {
	validators, err := automation.consensusSystem.GetAvailableValidators(requiredValidators)
	if err != nil {
		fmt.Printf("Error fetching available validators: %v\n", err)
		return
	}

	for _, validator := range validators {
		if err := automation.consensusSystem.AssignValidatorToShard(shardID, validator); err != nil {
			fmt.Printf("Error assigning validator %s to shard %s: %v\n", validator, shardID, err)
			continue
		}

		// Log the validator delegation into the ledger
		automation.logValidatorDelegation(shardID, validator)
	}
}

// logValidatorDelegation logs the delegation of a validator to a shard in the ledger.
func (automation *ShardDelegationAutomation) logValidatorDelegation(shardID string, validatorID string) {
	// Encrypt the delegation details
	encryptedDetails, err := encryption.EncryptDataWithKey([]byte(fmt.Sprintf("Validator %s delegated to Shard %s", validatorID, shardID)), ShardDelegationLogKey)
	if err != nil {
		fmt.Printf("Error encrypting delegation details for shard %s and validator %s: %v\n", shardID, validatorID, err)
		return
	}

	// Create a ledger entry for the validator delegation
	ledgerEntry := common.LedgerEntry{
		ID:        fmt.Sprintf("DELEGATION-%s-%s-%d", shardID, validatorID, time.Now().UnixNano()),
		Timestamp: time.Now().Unix(),
		Type:      "Validator Delegation",
		Status:    "Completed",
		Details:   string(encryptedDetails),
	}

	// Add the entry to the ledger
	if err := automation.ledgerInstance.AddEntry(ledgerEntry); err != nil {
		fmt.Printf("Failed to log validator delegation for shard %s: %v\n", shardID, err)
	}
}
