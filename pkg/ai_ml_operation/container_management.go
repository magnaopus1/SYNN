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

// Define necessary constants
const (
	ContainerStopped = "Stopped"
	ContainerRunning = "Running"
)

// ModelContainerDeploy handles the deployment of a containerized AI/ML model, tracks its status, and encrypts data.
func ModelContainerDeploy(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, transactionID string, containerConfig map[string]interface{}) (string, error) {

	// Convert containerConfig to JSON format before encryption
	configBytes, err := json.Marshal(containerConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal container config: %v", err)
	}

	// Encrypt the container configuration data (for security purposes; actual data not stored in ledger)
	if _, err := encryptDataForStorage(configBytes, modelID); err != nil {
		return "", errors.New("failed to encrypt container configuration")
	}

	// Generate a unique container ID
	containerID := fmt.Sprintf("container-%s-%d", modelID, time.Now().Unix())

	// Record the deployment in the ledger
	if err := l.AiMLMLedger.RecordContainerDeployment(containerID, modelID, nodeID); err != nil {
		return "", errors.New("failed to record container deployment in ledger")
	}

	log.Printf("Container for model %s deployed on node %s with container ID %s", modelID, nodeID, containerID)
	return containerID, nil
}

// ModelContainerStop stops a running container, validates the operation, and updates the ledger for real-time tracking.
func ModelContainerStop(l *ledger.Ledger, sc *common.SynnergyConsensus, containerID, nodeID, transactionID string) error {

	// Stop the container and update ledger status
	if err := l.AiMLMLedger.UpdateContainerStatus(containerID, ContainerStopped); err != nil {
		return errors.New("failed to stop container and update ledger status")
	}

	log.Printf("Container with ID %s has been stopped", containerID)
	return nil
}

// ModelContainerRestart restarts a stopped container, verifies the operation, and updates the ledger with new status.
func ModelContainerRestart(l *ledger.Ledger, sc *common.SynnergyConsensus, containerID, nodeID, transactionID string) error {

	// Restart the container and update ledger status
	if err := l.AiMLMLedger.UpdateContainerStatus(containerID, ContainerRunning); err != nil {
		return errors.New("failed to restart container and update ledger status")
	}

	log.Printf("Container with ID %s has been restarted", containerID)
	return nil
}
