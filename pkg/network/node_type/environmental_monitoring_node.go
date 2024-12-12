package node_type

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"    // Shared components like encryption, consensus, and storage
	"synnergy_network/pkg/ledger"    // Blockchain and ledger-related components
	"synnergy_network/pkg/sensor"    // Sensor and IoT-related components
)

// EnvironmentalMonitoringNode represents a node that integrates real-world environmental data with blockchain operations.
type EnvironmentalMonitoringNode struct {
	NodeID            string                     // Unique identifier for the node
	Blockchain        *ledger.Blockchain         // Local copy of the blockchain ledger
	ConsensusEngine   *common.SynnergyConsensus  // Consensus engine for validating transactions and data integrity
	EncryptionService *common.Encryption         // Encryption service for securing data from sensors and IoT devices
	NetworkManager    *common.NetworkManager     // Manages communication with other nodes and data sources
	SensorManager     *sensor.SensorManager      // Manages integration and data collection from environmental sensors
	mutex             sync.Mutex                 // Mutex for thread-safe operations
	SyncInterval      time.Duration              // Interval for syncing with other nodes
	SNVM              *synnergy_vm.VirtualMachine // Virtual Machine for executing smart contracts

}

// NewEnvironmentalMonitoringNode initializes a new environmental monitoring node in the Synnergy Network.
func NewEnvironmentalMonitoringNode(nodeID string, blockchain *ledger.Blockchain, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption, networkManager *common.NetworkManager, sensorManager *sensor.SensorManager, syncInterval time.Duration) *EnvironmentalMonitoringNode {
	return &EnvironmentalMonitoringNode{
		NodeID:            nodeID,
		Blockchain:        blockchain,
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		SensorManager:     sensorManager,
		SyncInterval:      syncInterval,
	}
}

// StartNode begins the environmental monitoring node's operations, syncing data from sensors, processing environmental data, and triggering blockchain actions.
func (emn *EnvironmentalMonitoringNode) StartNode() error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Begin syncing environmental data and triggering blockchain responses.
	go emn.syncEnvironmentalData()
	go emn.monitorSensorConditions()

	fmt.Printf("Environmental Monitoring Node %s started successfully.\n", emn.NodeID)
	return nil
}

// syncEnvironmentalData handles syncing environmental data from sensors and recording it to the blockchain.
func (emn *EnvironmentalMonitoringNode) syncEnvironmentalData() {
	ticker := time.NewTicker(emn.SyncInterval)
	defer ticker.Stop()

	for range ticker.C {
		emn.mutex.Lock()
		err := emn.recordSensorDataToBlockchain()
		if err != nil {
			fmt.Printf("Error syncing environmental data: %v\n", err)
		}
		emn.mutex.Unlock()
	}
}

// recordSensorDataToBlockchain retrieves data from sensors, encrypts it, and records it immutably on the blockchain.
func (emn *EnvironmentalMonitoringNode) recordSensorDataToBlockchain() error {
	// Retrieve environmental data from connected sensors.
	sensorData, err := emn.SensorManager.CollectData()
	if err != nil {
		return fmt.Errorf("failed to collect sensor data: %v", err)
	}

	// Encrypt the sensor data before recording.
	encryptedData, err := emn.EncryptionService.EncryptData(sensorData, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sensor data: %v", err)
	}

	// Record encrypted sensor data in the blockchain ledger.
	err = emn.Blockchain.RecordData(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to record sensor data to blockchain: %v", err)
	}

	fmt.Printf("Sensor data successfully recorded to the blockchain by node %s.\n", emn.NodeID)
	return nil
}

// monitorSensorConditions listens for environmental conditions that trigger blockchain actions based on pre-set thresholds.
func (emn *EnvironmentalMonitoringNode) monitorSensorConditions() {
	for {
		sensorData, err := emn.SensorManager.CollectData()
		if err != nil {
			fmt.Printf("Error receiving sensor data: %v\n", err)
			continue
		}

		// Validate and process environmental conditions.
		err = emn.processEnvironmentalCondition(sensorData)
		if err != nil {
			fmt.Printf("Error processing environmental data: %v\n", err)
		}
	}
}

// processEnvironmentalCondition evaluates sensor data and triggers smart contracts if certain conditions are met.
func (emn *EnvironmentalMonitoringNode) processEnvironmentalCondition(sensorData []byte) error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Define thresholds or conditions for triggering smart contracts.
	conditionMet, err := emn.SensorManager.CheckCondition(sensorData)
	if err != nil || !conditionMet {
		return fmt.Errorf("condition not met: %v", err)
	}

	// Trigger a smart contract execution based on the environmental data.
	contractID := "EnvironmentalResponseContract" // Example smart contract ID
	err = emn.Blockchain.TriggerSmartContract(contractID, sensorData)
	if err != nil {
		return fmt.Errorf("failed to trigger smart contract: %v", err)
	}

	fmt.Printf("Smart contract %s triggered by environmental data from node %s.\n", contractID, emn.NodeID)
	return nil
}

// Sensor and IoT Management

// addSensor adds a new environmental sensor to the node for data collection and monitoring.
func (emn *EnvironmentalMonitoringNode) addSensor(sensorID string, sensorType string) error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Register the new sensor with the SensorManager.
	err := emn.SensorManager.RegisterSensor(sensorID, sensorType)
	if err != nil {
		return fmt.Errorf("failed to add sensor %s: %v", sensorID, err)
	}

	fmt.Printf("Sensor %s of type %s added successfully to node %s.\n", sensorID, sensorType, emn.NodeID)
	return nil
}

// removeSensor removes an environmental sensor from the node's sensor list.
func (emn *EnvironmentalMonitoringNode) removeSensor(sensorID string) error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Unregister the sensor from the SensorManager.
	err := emn.SensorManager.UnregisterSensor(sensorID)
	if err != nil {
		return fmt.Errorf("failed to remove sensor %s: %v", sensorID, err)
	}

	fmt.Printf("Sensor %s removed successfully from node %s.\n", sensorID, emn.NodeID)
	return nil
}

// Environmental Data Security

// applyDataSecurity ensures encryption protocols are applied to all sensor data before being recorded on the blockchain.
func (emn *EnvironmentalMonitoringNode) applyDataSecurity() error {
	// Ensure encryption protocols are up to date for environmental data.
	err := emn.EncryptionService.ApplySecurity(emn.NodeID)
	if err != nil {
		return fmt.Errorf("failed to apply encryption security for environmental data: %v", err)
	}

	fmt.Printf("Encryption security applied successfully for environmental data on node %s.\n", emn.NodeID)
	return nil
}

// Compliance and Regulatory Monitoring

// checkRegulatoryCompliance ensures the node's environmental data collection and actions comply with legal and environmental regulations.
func (emn *EnvironmentalMonitoringNode) checkRegulatoryCompliance() error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Conduct compliance checks using the SensorManager.
	err := emn.SensorManager.EnsureCompliance(emn.NodeID)
	if err != nil {
		return fmt.Errorf("regulatory compliance check failed: %v", err)
	}

	fmt.Printf("Regulatory compliance check passed for node %s.\n", emn.NodeID)
	return nil
}

// Environmental Data Analytics

// performRealTimeAnalysis conducts real-time analysis on incoming environmental data to detect trends and anomalies.
func (emn *EnvironmentalMonitoringNode) performRealTimeAnalysis() error {
	emn.mutex.Lock()
	defer emn.mutex.Unlock()

	// Analyze real-time sensor data for trends and anomalies.
	analysisResults, err := emn.SensorManager.AnalyzeData()
	if err != nil {
		return fmt.Errorf("real-time analysis failed: %v", err)
	}

	// Log the analysis results.
	for _, result := range analysisResults {
		fmt.Printf("Analysis result: %s - Status: %s\n", result.SensorID, result.Status)
	}

	return nil
}
