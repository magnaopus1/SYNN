package authority_node

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/ledger"
	"synnergy_network_demo/network"
	"synnergy_network_demo/common"
)

// AuthorityNodeRewardPool manages and distributes the reward pool for authority nodes.
type AuthorityNodeRewardPool struct {
	TotalPoolAmount    float64                            // Total reward pool amount
	AuthorityNodes     map[string]*common.AuthorityNode   // Active authority nodes
	Ledger             *ledger.Ledger                     // Reference to the blockchain ledger
	NetworkManager     *network.NetworkManager            // Reference to the network manager for discovering active nodes
	mutex              sync.Mutex                         // Mutex for thread-safe operations
	DistributionPeriod time.Duration                      // Distribution interval (default: 24 hours)
}

// NewAuthorityNodeRewardPool initializes a new reward pool for authority nodes.
func NewAuthorityNodeRewardPool(initialAmount float64, ledgerInstance *ledger.Ledger, networkManager *network.NetworkManager, distributionPeriod time.Duration) *AuthorityNodeRewardPool {
	pool := &AuthorityNodeRewardPool{
		TotalPoolAmount:    initialAmount,
		AuthorityNodes:     make(map[string]*common.AuthorityNode),
		Ledger:             ledgerInstance,
		NetworkManager:     networkManager,
		DistributionPeriod: distributionPeriod,
	}

	// Start the automatic distribution process.
	go pool.startRewardDistribution()

	return pool
}

// startRewardDistribution automatically distributes rewards every 24 hours to active online authority nodes.
func (rp *AuthorityNodeRewardPool) startRewardDistribution() {
	ticker := time.NewTicker(rp.DistributionPeriod)
	defer ticker.Stop()

	for range ticker.C {
		rp.distributeRewards()
	}
}

// distributeRewards calculates the reward per node and distributes equally to all online authority nodes.
func (rp *AuthorityNodeRewardPool) distributeRewards() {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()

	// Discover all active online authority nodes.
	activeNodes := rp.NetworkManager.DiscoverActiveAuthorityNodes()

	// Ensure there are active nodes available.
	if len(activeNodes) == 0 {
		fmt.Println("No active authority nodes available for reward distribution.")
		return
	}

	// Calculate the reward to distribute equally among active nodes.
	rewardPerNode := rp.TotalPoolAmount / float64(len(activeNodes))

	// Distribute the rewards.
	for _, node := range activeNodes {
		// Log the reward distribution in the ledger.
		err := rp.Ledger.RecordRewardDistribution(node.NodeID, rewardPerNode)
		if err != nil {
			fmt.Printf("Failed to record reward for node %s: %v\n", node.NodeID, err)
			continue
		}

		// Simulate the reward being added to the authority node's wallet.
		fmt.Printf("Distributed %.2f rewards to node %s.\n", rewardPerNode, node.NodeID)
	}

	// Reset the pool after distribution.
	rp.TotalPoolAmount = 0
}



// ViewRewardPool returns the current total reward pool amount (view-only).
func (rp *AuthorityNodeRewardPool) ViewRewardPool() float64 {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()

	return rp.TotalPoolAmount
}
