package advanced_security

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

var (
	nodeVotingWeights   = make(map[string]int)
	nodeWeightsMutex    sync.Mutex
	rateLimitAPILock    sync.Mutex
	rateLimitingEnabled = false
)

// Structs

type RateLimitConfig struct {
	Limit int `json:"limit"`
}

type DataTransferMetrics struct {
	RateMBps     int       `json:"rate_mbps"`
	PeakRateMBps int       `json:"peak_rate_mbps"`
	Timestamp    time.Time `json:"timestamp"`
}

// SetRateLimitPolicy sets and records the rate limit policy
func SetRateLimitPolicy(policy string) error {
	err := DefineRateLimitPolicy(policy)
	if err != nil {
		return fmt.Errorf("failed to set rate limit policy: %w", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	ledgerInstance := &ledger.Ledger{}

	// Call the method directly, without assigning it to a variable
	ledgerInstance.AdvancedSecurityLedger.RecordRateLimitPolicy(policy, timestamp)

	log.Printf("Rate limit policy set and recorded: %s", policy)
	return nil
}

// DefineRateLimitPolicy defines a rate limit policy for network operations
func DefineRateLimitPolicy(policy string) error {
	var limit int

	switch policy {
	case "low":
		limit = 100
	case "medium":
		limit = 500
	case "high":
		limit = 1000
	default:
		return errors.New("invalid rate limit policy: must be 'low', 'medium', or 'high'")
	}

	if err := setNetworkRateLimit(limit); err != nil {
		return fmt.Errorf("failed to apply rate limit: %w", err)
	}

	log.Printf("Rate limit policy defined: %s with limit %d", policy, limit)
	return nil
}

// SetTransactionThreshold sets and records the transaction threshold
func SetTransactionThreshold(threshold int) error {
	if threshold <= 0 {
		return fmt.Errorf("invalid threshold: must be positive")
	}

	timestamp := time.Now().Format(time.RFC3339)
	ledgerInstance := &ledger.Ledger{}

	// Call the method directly, without assigning it to a variable
	ledgerInstance.AdvancedSecurityLedger.RecordTransactionThreshold(threshold, timestamp)

	log.Printf("Transaction threshold set and recorded: %d", threshold)
	return nil
}

// SetAPIRateLimit sets the rate limit for API requests
func SetAPIRateLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("rate limit must be a positive integer")
	}

	if err := setNetworkRateLimit(limit); err != nil {
		return fmt.Errorf("failed to set API rate limit: %w", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	ledgerInstance := &ledger.Ledger{}

	// Call the method directly without assigning its return value
	ledgerInstance.AdvancedSecurityLedger.RecordAPILimitSet(limit, timestamp)

	log.Printf("API rate limit set and recorded: %d requests/second", limit)
	return nil
}

// EnableRateLimiting enables and records rate limiting on network traffic
func EnableRateLimiting() error {
	rateLimitAPILock.Lock()
	defer rateLimitAPILock.Unlock()

	if rateLimitingEnabled {
		return fmt.Errorf("rate limiting is already enabled")
	}

	if err := ActivateRateLimiting(); err != nil {
		return fmt.Errorf("failed to enable rate limiting: %w", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	ledgerInstance := &ledger.Ledger{}

	// Call the method directly without assigning its return value
	ledgerInstance.AdvancedSecurityLedger.RecordRateLimitingStatus("enabled", timestamp)

	rateLimitingEnabled = true
	log.Println("Rate limiting enabled and recorded successfully.")
	return nil
}

// ActivateRateLimiting enables rate limiting for network traffic
func ActivateRateLimiting() error {
	apiURL := os.Getenv("RATE_LIMIT_API_URL")
	apiKey := os.Getenv("RATE_LIMIT_API_KEY")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("rate-limiting API URL or API key not configured")
	}

	payload := map[string]bool{"rateLimitingEnabled": true}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal activation payload: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/activate", apiURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create rate-limiting request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send activation request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rate-limiting API returned status: %s", resp.Status)
	}

	log.Println("Rate limiting activated successfully.")
	return nil
}

// applyRateLimitToAPISystem applies the rate limit to an external rate-limiting system
func applyRateLimitToAPISystem(limit int) error {
	apiURL := os.Getenv("RATE_LIMIT_API_URL")
	apiKey := os.Getenv("RATE_LIMIT_API_KEY")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("rate-limiting API URL or API key not configured")
	}

	rateLimitConfig := RateLimitConfig{Limit: limit}
	payload, err := json.Marshal(rateLimitConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal rate limit config: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create rate limit request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send rate limit request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rate-limiting API returned status: %s", resp.Status)
	}

	log.Printf("Rate limit successfully applied: %d requests/second", limit)
	return nil
}

// setNetworkRateLimit applies the rate limit to a rate-limiting system via an API call.
func setNetworkRateLimit(limit int) error {
	apiURL := os.Getenv("RATE_LIMIT_API_URL")
	apiKey := os.Getenv("RATE_LIMIT_API_KEY")
	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("rate-limiting API URL or API key not configured")
	}

	// Define the rate limit configuration to send to the API
	rateLimitConfig := RateLimitConfig{Limit: limit}
	payload, err := json.Marshal(rateLimitConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal rate limit configuration: %v", err)
	}

	// Create a new HTTP POST request to apply the rate limit
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/set_limit", apiURL), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request for setting rate limit: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send rate limit request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rate limit API returned non-OK status: %s", resp.Status)
	}

	log.Printf("Rate limit of %d requests per second successfully applied to external system", limit)
	return nil
}

// DisableRateLimiting disables rate limiting on network traffic
func DisableRateLimiting() error {
	err := DeactivateRateLimiting()
	if err != nil {
		return fmt.Errorf("failed to disable rate limiting: %v", err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.AdvancedSecurityLedger.RecordRateLimitingStatus("disabled", timestamp)

	fmt.Println("Rate limiting disabled.")
	return nil
}

// SetTransferRateLimit sets a rate limit on data transfers within the network
func SetTransferRateLimit(limit int) error {
	err := SetDataTransferLimit(limit)
	if err != nil {
		return fmt.Errorf("failed to set transfer rate limit: %v", err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.AdvancedSecurityLedger.RecordTransferRateLimit(limit, timestamp)

	fmt.Printf("Transfer rate limit set to: %d MB/s\n", limit)
	return nil
}

// DeactivateRateLimiting disables rate limiting by updating the rate-limiting system.
func DeactivateRateLimiting() error {
	apiURL := os.Getenv("RATE_LIMIT_API_URL")
	apiKey := os.Getenv("RATE_LIMIT_API_KEY")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("rate-limiting API URL or API key is not configured")
	}

	// Define the request payload to disable rate limiting
	requestBody := map[string]bool{"rateLimitingEnabled": false}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal rate-limiting deactivation payload: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/deactivate", apiURL), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request for deactivating rate limiting: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send rate-limiting deactivation request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rate-limiting deactivation API returned status: %s", resp.Status)
	}

	log.Println("Rate limiting successfully disabled on network traffic.")
	return nil
}

// SetDataTransferLimit applies a data transfer rate limit for network operations
func SetDataTransferLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("data transfer limit must be a positive integer")
	}

	// Assuming interaction with an external data transfer management system
	apiURL := os.Getenv("DATA_TRANSFER_LIMIT_API_URL")
	apiKey := os.Getenv("DATA_TRANSFER_API_KEY")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("data transfer limit API URL or API key is not configured")
	}

	limitConfig := map[string]int{"transferLimit": limit}
	payload, err := json.Marshal(limitConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal data transfer limit configuration: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/set_limit", apiURL), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request for setting data transfer limit: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send data transfer limit request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("data transfer limit API returned non-OK status: %s", resp.Status)
	}

	log.Printf("Data transfer rate limit set to %d MB/s", limit)
	return nil
}

// MonitorDataTransfer monitors data transfer rates for compliance with rate limits.
func MonitorDataTransfer() error {
	// Get data transfer metrics from the common package
	transferData, err := CheckDataTransferRate()
	if err != nil {
		return fmt.Errorf("failed to monitor data transfer: %v", err)
	}

	// Convert common.DataTransferMetrics to ledger.DataTransferMetrics
	ledgerTransferData := ledger.DataTransferMetrics{
		RateMBps:     transferData.RateMBps,
		PeakRateMBps: transferData.PeakRateMBps,
		Timestamp:    transferData.Timestamp,
	}

	// Record data transfer in the ledger with the converted type
	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.DataManagementLedger.RecordDataTransferMonitor(ledgerTransferData, timestamp)

	fmt.Printf("Data transfer monitored: %v MB/s, Peak rate: %v MB/s\n", transferData.RateMBps, transferData.PeakRateMBps)
	return nil
}

// SetConsensusThreshold sets a threshold for consensus decision-making.
func SetConsensusThreshold(threshold int) error {
	err := setConsensusThreshold(threshold)
	if err != nil {
		return fmt.Errorf("failed to set consensus threshold: %v", err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.BlockchainConsensusCoinLedger.RecordConsensusThreshold(threshold, timestamp)

	fmt.Printf("Consensus threshold set to: %d%%\n", threshold)
	return nil
}

// CheckDataTransferRate retrieves current data transfer metrics by interacting with a monitoring system.
func CheckDataTransferRate() (DataTransferMetrics, error) {
	apiURL := os.Getenv("DATA_TRANSFER_MONITORING_API_URL")
	apiKey := os.Getenv("MONITORING_API_KEY")
	if apiURL == "" || apiKey == "" {
		return DataTransferMetrics{}, fmt.Errorf("data transfer monitoring API URL or API key not configured")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/transfer_metrics", apiURL), nil)
	if err != nil {
		return DataTransferMetrics{}, fmt.Errorf("failed to create request for data transfer metrics: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return DataTransferMetrics{}, fmt.Errorf("failed to fetch data transfer metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DataTransferMetrics{}, fmt.Errorf("monitoring API returned status: %s", resp.Status)
	}

	var metrics DataTransferMetrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return DataTransferMetrics{}, fmt.Errorf("failed to decode data transfer metrics: %v", err)
	}

	log.Printf("Current data transfer rate: %d MB/s, Peak rate: %d MB/s", metrics.RateMBps, metrics.PeakRateMBps)
	return metrics, nil
}

// SetConsensusThreshold configures a consensus threshold for decision-making within the network.
func setConsensusThreshold(threshold int) error {
	if threshold < 1 || threshold > 100 {
		return fmt.Errorf("consensus threshold must be between 1 and 100")
	}

	consensusAPI := os.Getenv("CONSENSUS_MANAGEMENT_API_URL")
	apiKey := os.Getenv("CONSENSUS_API_KEY")
	if consensusAPI == "" || apiKey == "" {
		return fmt.Errorf("consensus management API URL or API key not configured")
	}

	// Prepare the request payload
	payload := map[string]int{"consensusThreshold": threshold}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal consensus threshold payload: %v", err)
	}

	// Create HTTP request to set the consensus threshold
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/set_threshold", consensusAPI), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request for setting consensus threshold: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to set consensus threshold: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("consensus management API returned status: %s", resp.Status)
	}

	// Log the setting of the consensus threshold for auditing
	log.Printf("Consensus threshold successfully set to %d%%", threshold)
	return nil
}

// EnableConsensusAnomalyDetection activates anomaly detection specifically for consensus processes.
func EnableConsensusAnomalyDetection() error {
	err := ActivateConsensusAnomalyDetection()
	if err != nil {
		return fmt.Errorf("failed to enable consensus anomaly detection: %v", err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.AdvancedSecurityLedger.RecordConsensusAnomalyDetectionStatus("enabled", timestamp)

	fmt.Println("Consensus anomaly detection enabled.")
	return nil
}

// DisableConsensusAnomalyDetection deactivates anomaly detection for consensus processes.
func DisableConsensusAnomalyDetection() error {
	err := DeactivateConsensusAnomalyDetection()
	if err != nil {
		return fmt.Errorf("failed to disable consensus anomaly detection: %v", err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.AdvancedSecurityLedger.RecordConsensusAnomalyDetectionStatus("disabled", timestamp)

	fmt.Println("Consensus anomaly detection disabled.")
	return nil
}

// ActivateConsensusAnomalyDetection enables anomaly detection for consensus processes by making an API call or setting configuration.
func ActivateConsensusAnomalyDetection() error {
	apiURL := os.Getenv("CONSENSUS_ANOMALY_DETECTION_API_URL")
	apiKey := os.Getenv("CONSENSUS_MONITORING_API_KEY")
	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("consensus anomaly detection API URL or API key not configured")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/enable", apiURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request to enable consensus anomaly detection: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to enable consensus anomaly detection: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("consensus anomaly detection API returned status: %s", resp.Status)
	}

	log.Println("Consensus anomaly detection activated.")
	return nil
}

// DeactivateConsensusAnomalyDetection disables anomaly detection for consensus processes by making an API call or setting configuration.
func DeactivateConsensusAnomalyDetection() error {
	apiURL := os.Getenv("CONSENSUS_ANOMALY_DETECTION_API_URL")
	apiKey := os.Getenv("CONSENSUS_MONITORING_API_KEY")
	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("consensus anomaly detection API URL or API key not configured")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/disable", apiURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request to disable consensus anomaly detection: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to disable consensus anomaly detection: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("consensus anomaly detection API returned status: %s", resp.Status)
	}

	log.Println("Consensus anomaly detection deactivated.")
	return nil
}

// SetNodeVotingWeight sets the voting weight for a specific node.
func SetNodeVotingWeight(nodeID string, weight int) error {
	err := AssignNodeVotingWeight(nodeID, weight)
	if err != nil {
		return fmt.Errorf("failed to set voting weight for node %s: %v", nodeID, err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.RecordNodeVotingWeight(nodeID, weight, timestamp)

	fmt.Printf("Voting weight for node %s set to: %d\n", nodeID, weight)
	return nil
}

// AdjustNodeVotingWeight adjusts the voting weight of a node dynamically.
func AdjustNodeVotingWeight(nodeID string, adjustment int) error {
	newWeight, err := AdjustVotingWeight(nodeID, adjustment)
	if err != nil {
		return fmt.Errorf("failed to adjust voting weight for node %s: %v", nodeID, err)
	}

	ledgerInstance := &ledger.Ledger{}
	timestamp := time.Now().Format(time.RFC3339) // Format timestamp as string
	ledgerInstance.RecordNodeVotingWeight(nodeID, newWeight, timestamp)

	fmt.Printf("Voting weight for node %s adjusted to: %d\n", nodeID, newWeight)
	return nil
}

// AssignNodeVotingWeight sets the initial voting weight for a node.
func AssignNodeVotingWeight(nodeID string, weight int) error {
	if weight < 0 {
		return fmt.Errorf("voting weight must be non-negative")
	}

	nodeWeightsMutex.Lock()
	defer nodeWeightsMutex.Unlock()

	nodeVotingWeights[nodeID] = weight
	log.Printf("Voting weight for node %s assigned to %d", nodeID, weight)
	return nil
}

// AdjustVotingWeight dynamically adjusts the voting weight for a node and returns the new weight.
func AdjustVotingWeight(nodeID string, adjustment int) (int, error) {
	nodeWeightsMutex.Lock()
	defer nodeWeightsMutex.Unlock()

	currentWeight, exists := nodeVotingWeights[nodeID]
	if !exists {
		return 0, fmt.Errorf("node %s does not exist", nodeID)
	}

	newWeight := currentWeight + adjustment
	if newWeight < 0 {
		return 0, fmt.Errorf("resulting voting weight must be non-negative")
	}

	nodeVotingWeights[nodeID] = newWeight
	log.Printf("Voting weight for node %s adjusted to %d", nodeID, newWeight)
	return newWeight, nil
}
