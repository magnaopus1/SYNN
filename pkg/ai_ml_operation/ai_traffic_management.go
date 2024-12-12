package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"

	"time"
)

// AITrafficRoute directs traffic to the optimal model node based on workload, model capacity, and system health.
func AITrafficRoute(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID, nodeID string, trafficData map[string]interface{}) (string, error) {

	// Encrypt traffic data after converting it to JSON format
	dataBytes, err := json.Marshal(trafficData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal traffic data for encryption: %w", err)
	}
	encryptedTrafficData, err := encryptDataForStorage(dataBytes, modelID)
	if err != nil {
		return "", errors.New("failed to encrypt traffic data for routing")
	}

	// Fetch optimal node for routing traffic
	nodeID, err = l.AiMLMLedger.GetOptimalNodeForModel(modelID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve optimal node for traffic routing: %w", err)
	}

	// Use encrypted data in logging or further processing
	log.Printf("Traffic for model %s routed to node %s. Encrypted traffic data: %x", modelID, nodeID, encryptedTrafficData)
	return nodeID, nil
}


// ModelTrafficBalance balances traffic across model nodes by redistributing workloads to optimize performance.
func ModelTrafficBalance(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID, nodeID string) error {
	// Define a threshold for load
	const LoadThreshold = 80

	for {
		// Retrieve model's current load
		modelLoad, err := l.AiMLMLedger.GetModelLoad(modelID)
		if err != nil {
			return fmt.Errorf("failed to retrieve model load: %w", err)
		}

		// Redistribute traffic if load exceeds the threshold
		if modelLoad > LoadThreshold {
			targetNode, err := l.AiMLMLedger.GetOptimalNodeForModel(modelID)
			if err != nil {
				return fmt.Errorf("failed to find target node for redistribution: %w", err)
			}
			if err := l.AiMLMLedger.RedistributeModelTraffic(modelID, targetNode); err != nil {
				return fmt.Errorf("failed to redistribute traffic: %w", err)
			}
			log.Printf("Traffic for model %s balanced to node %s", modelID, targetNode)

			// Log traffic balance operation
			if err := l.AiMLMLedger.LogTrafficBalanceOperation(modelID, targetNode); err != nil {
				return fmt.Errorf("failed to log traffic balance operation: %w", err)
			}
		}

		// Re-check model load periodically
		time.Sleep(10 * time.Second)
	}
}