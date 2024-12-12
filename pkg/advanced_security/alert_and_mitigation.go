package advanced_security

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

const (
    maxRetries         = 3
    retryInterval      = 2 * time.Second
    nodeEndpointFormat = "https://%s/api/v1/alerts" //I need to adjust
)

// FirmwareMetadata holds expected firmware integrity data
type FirmwareMetadata struct {
    ExpectedHash string
    Version      string
}
// FirmwareCheckManager manages the state of firmware integrity checks
type FirmwareCheckManager struct {
    active bool
    mutex  sync.Mutex
}

// NewFirmwareCheckManager initializes a new firmware check manager
func NewFirmwareCheckManager() *FirmwareCheckManager {
    return &FirmwareCheckManager{active: true}
}

// APIUsageStats stores detailed usage data for an API
type APIUsageStats struct {
    TotalRequests     int
    LastAccessed      time.Time
    AvgResponseTime   float64
    PeakUsageDetected bool
    ResponseTimes     []float64 // Stores individual response times for better avg calculation
}

// APIUsageManager manages API usage statistics and monitors for potential abuse
type APIUsageManager struct {
    usageData map[string]*APIUsageStats
    mutex     sync.Mutex
    alertThreshold int // Threshold for peak usage detection
}

// NewAPIUsageManager initializes a new API usage manager
func NewAPIUsageManager() *APIUsageManager {
    return &APIUsageManager{
        usageData: make(map[string]*APIUsageStats),
    }
}

// SendSecurityAlert sends a security alert when the threshold is exceeded
func SendSecurityAlert(message string) error {
    encryption := common.Encryption{}
    key := []byte("your_16_or_24_or_32_byte_key")

    encryptedMessage, err := encryption.EncryptData("AES", []byte(message), key)
    if err != nil {
        return fmt.Errorf("failed to encrypt security alert: %w", err)
    }

    encodedMessage := base64.StdEncoding.EncodeToString(encryptedMessage)
    err = BroadcastAlert(encodedMessage)
    if err != nil {
        return fmt.Errorf("failed to send security alert: %w", err)
    }

    ledgerInstance := &ledger.Ledger{}
    ledgerInstance.AdvancedSecurityLedger.RecordAlertSent(time.Now().Format("20060102150405"), message)

    log.Println("Security alert sent.")
    return nil
}


// sendToNode securely sends the alert message to a node using HTTPS with TLS
func sendToNode(node, message string) error {
    endpoint := fmt.Sprintf("https://%s/api/alert", node)
    client := &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
                InsecureSkipVerify: false,
            },
        },
    }

    for i := 0; i < 3; i++ {
        req, err := http.NewRequest("POST", endpoint, strings.NewReader(message))
        if err != nil {
            log.Printf("Error creating request for node %s: %v", node, err)
            return err
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer your_auth_token")

        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error sending alert to node %s (attempt %d/3): %v", node, i+1, err)
            time.Sleep(2 * time.Second)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            log.Printf("Alert successfully sent to node %s", node)
            return nil
        } else {
            body, _ := ioutil.ReadAll(resp.Body)
            log.Printf("Failed to send alert to node %s, status: %d, response: %s", node, resp.StatusCode, string(body))
            if resp.StatusCode >= 400 && resp.StatusCode < 500 {
                return fmt.Errorf("client error on node %s: %s", node, string(body))
            }
            time.Sleep(2 * time.Second)
        }
    }

    return fmt.Errorf("failed to send alert after multiple attempts")
}

// SetThreatMitigationPlan sets the plan for mitigating specific security threats
func SetThreatMitigationPlan(planID, description string) error {
    err := DefineMitigationPlan(planID, description)
    if err != nil {
        return fmt.Errorf("failed to define threat mitigation plan: %w", err)
    }

    ledgerInstance := &ledger.Ledger{}
    ledgerInstance.AdvancedSecurityLedger.RecordMitigationPlanSet(planID, description)

    log.Println("Threat mitigation plan successfully set and recorded.")
    return nil
}

// DefineMitigationPlan validates and defines a mitigation plan with necessary details
func DefineMitigationPlan(planID, description string) error {
    ledgerInstance := &ledger.Ledger{}

    if ledgerInstance.AdvancedSecurityLedger.MitigationPlanExists(planID) {
        return fmt.Errorf("mitigation plan with ID %s already exists", planID)
    }

    ledgerInstance.AdvancedSecurityLedger.RecordMitigationPlanSet(planID, description)
    log.Printf("Mitigation plan defined with ID: %s, Description: %s", planID, description)
    return nil
}

// ActivateThreatMitigation activates the mitigation plan in response to a detected threat
func ActivateThreatMitigation(planID string) error {
    err := ActivateMitigation(planID)
    if err != nil {
        return fmt.Errorf("failed to activate threat mitigation: %w", err)
    }

    ledgerInstance := &ledger.Ledger{}
    ledgerInstance.AdvancedSecurityLedger.RecordMitigationActivated(planID)

    log.Println("Threat mitigation activated.")
    return nil
}

// ActivateMitigation validates and activates a specific mitigation plan
func ActivateMitigation(planID string) error {
	ledgerInstance := &ledger.Ledger{}

	// Ensure the plan exists
	if !ledgerInstance.AdvancedSecurityLedger.MitigationPlanExists(planID) {
		return fmt.Errorf("mitigation plan with ID %s does not exist", planID)
	}

	// Check if the plan is already active
	active, err := ledgerInstance.AdvancedSecurityLedger.IsMitigationPlanActive(planID)
	if err != nil {
		return fmt.Errorf("failed to check if mitigation plan %s is active: %v", planID, err)
	}
	if active {
		return fmt.Errorf("mitigation plan %s is already active", planID)
	}
	ledgerInstance.AdvancedSecurityLedger.ActivateMitigationPlan(planID)

	log.Printf("Mitigation plan %s successfully activated.", planID)
	return nil
}


// DeactivateThreatMitigation deactivates the mitigation plan after the threat is resolved.
func DeactivateThreatMitigation(planID string) error {
	ledgerInstance := &ledger.Ledger{}

	// Ensure the plan exists
	if !ledgerInstance.AdvancedSecurityLedger.MitigationPlanExists(planID) {
		return fmt.Errorf("mitigation plan with ID %s does not exist", planID)
	}

	// Check if the plan is active
	active, err := ledgerInstance.AdvancedSecurityLedger.IsMitigationPlanActive(planID)
	if err != nil {
		return fmt.Errorf("failed to check if mitigation plan %s is active: %v", planID, err)
	}
	if !active {
		return fmt.Errorf("mitigation plan %s is not active", planID)
	}

	// Deactivate the plan (directly call the method without handling a return value)
	ledgerInstance.AdvancedSecurityLedger.DeactivateMitigationPlan(planID)

	log.Printf("Mitigation plan %s successfully deactivated.", planID)
	return nil
}


// VerifyMitigationEffectiveness verifies the effectiveness of the applied mitigation plan
func VerifyMitigationEffectiveness(planID string) error {
	ledgerInstance := &ledger.Ledger{}

	// Ensure the plan exists
	if !ledgerInstance.AdvancedSecurityLedger.MitigationPlanExists(planID) {
		return fmt.Errorf("mitigation plan with ID %s does not exist", planID)
	}

	// Fetch effectiveness metrics
	metrics, err := ledgerInstance.AdvancedSecurityLedger.GetMitigationMetrics(planID)
	if err != nil {
		return fmt.Errorf("failed to retrieve metrics for plan %s: %v", planID, err)
	}

	// Analyze effectiveness
	effectiveness := analyzeEffectiveness(metrics)

	// Record effectiveness in the ledger (direct call without error handling)
	ledgerInstance.AdvancedSecurityLedger.RecordMitigationEffectiveness(planID, effectiveness)

	log.Printf("Effectiveness for plan %s recorded as: %s", planID, effectiveness)
	return nil
}


// analyzeEffectiveness processes metrics to determine the effectiveness level
func analyzeEffectiveness(metrics ledger.MitigationMetrics) string {
    switch {
    case metrics.IncidentReductionRate > 75 && metrics.PerformanceImprovementScore > 80:
        return "High"
    case metrics.IncidentReductionRate > 50 && metrics.PerformanceImprovementScore > 60:
        return "Moderate"
    default:
        return "Low"
    }
}

// SetEventAlertPolicy defines the policy for handling alert events within the network
func SetEventAlertPolicy(policy string) error {
    ledgerInstance := &ledger.Ledger{}

    // Validate the policy
    if policy == "" {
        return fmt.Errorf("policy cannot be empty")
    }

    // Update policy in the ledger
    err := ledgerInstance.AdvancedSecurityLedger.UpdateCurrentAlertPolicy(policy)
    if err != nil {
        return fmt.Errorf("failed to set event alert policy: %v", err)
    }

    // Record the policy change
    err = ledgerInstance.AdvancedSecurityLedger.RecordAlertPolicySet(policy, time.Now())
    if err != nil {
        return fmt.Errorf("failed to record alert policy in ledger: %v", err)
    }

    log.Printf("Event alert policy set to: %s", policy)
    return nil
}

// EnableEventMonitoring enables monitoring of events that could trigger security alerts
func EnableEventMonitoring() error {
    ledgerInstance := &ledger.Ledger{}

    // Enable monitoring
    err := ledgerInstance.AdvancedSecurityLedger.EnableEventMonitoring()
    if err != nil {
        return fmt.Errorf("failed to enable event monitoring: %v", err)
    }

    // Record monitoring status
    timestamp := time.Now().Format(time.RFC3339)
    ledgerInstance.AdvancedSecurityLedger.RecordEventMonitoringStatus("enabled", timestamp)

    log.Println("Event monitoring enabled.")
    return nil
}


// EnableMonitoring activates the monitoring system to watch for security events
func EnableMonitoring() error {
    ledgerInstance := &ledger.Ledger{}

    // Check if monitoring is already enabled
    status, err := ledgerInstance.AdvancedSecurityLedger.GetEventMonitoringStatus()
    if err != nil {
        return fmt.Errorf("failed to retrieve monitoring status: %w", err)
    }

    if status == "enabled" {
        return fmt.Errorf("event monitoring is already enabled")
    }

    // Enable the monitoring system
    ledgerInstance.AdvancedSecurityLedger.SetEventMonitoringStatus(true)

    log.Println("Event monitoring system has been activated.")
    return nil
}



// LogThreatEvent logs any security threat event detected by the network
func LogThreatEvent(event string) error {
    ledgerInstance := &ledger.Ledger{}

    // Record the threat event in the ledger
    ledgerInstance.AdvancedSecurityLedger.RecordThreatEvent(event, time.Now())

    log.Printf("Threat event logged: %s\n", event)
    return nil
}


// EnableFirmwareCheck enables firmware integrity checks for added security
func EnableFirmwareCheck() error {
    // Activate the firmware check
    err := ActivateFirmwareCheck()
    if err != nil {
        return fmt.Errorf("failed to enable firmware check: %w", err)
    }

    // Initialize ledger instance
    ledgerInstance := &ledger.Ledger{}
    timestamp := time.Now().Format(time.RFC3339)

    // Record the firmware check status in the ledger
    ledgerInstance.AdvancedSecurityLedger.RecordFirmwareCheckStatus("enabled", timestamp)

    log.Println("Firmware check enabled.")
    return nil
}

// LoadFirmwareMetadata loads expected firmware data for integrity checks
func LoadFirmwareMetadata() (*FirmwareMetadata, error) {
    // Load metadata from a secure source
    return &FirmwareMetadata{
        ExpectedHash: "your_expected_sha256_hash", // Replace with real hash
        Version:      "1.0.0",                    // Replace with real version
    }, nil
}

// GetFirmwareFilePath returns the file path of the firmware binary
func GetFirmwareFilePath() string {
    // Return the actual firmware file path
    return "/path/to/firmware.bin" // Replace with actual path
}

// ActivateFirmwareCheck enables firmware integrity checks for added security
func ActivateFirmwareCheck() error {
    metadata, err := LoadFirmwareMetadata()
    if err != nil {
        return fmt.Errorf("failed to load firmware metadata: %w", err)
    }

    firmwarePath := GetFirmwareFilePath()
    firmwareData, err := ioutil.ReadFile(firmwarePath)
    if err != nil {
        return fmt.Errorf("failed to read firmware file: %w", err)
    }

    firmwareHash := sha256.Sum256(firmwareData)
    firmwareHashString := hex.EncodeToString(firmwareHash[:])

    if firmwareHashString != metadata.ExpectedHash {
        return fmt.Errorf("firmware integrity check failed: hash mismatch")
    }

    if metadata.Version != "1.0.0" { // Replace with real logic
        return fmt.Errorf("incompatible firmware version: expected %s, got %s", metadata.Version, "1.0.0")
    }

    log.Printf("Firmware integrity verified successfully. Version: %s, Hash: %s\n", metadata.Version, firmwareHashString)
    return nil
}

// DisableFirmwareCheck disables firmware integrity checks
func DisableFirmwareCheck() error {
    firmwareManager := &FirmwareCheckManager{active: true}

    // Deactivate the firmware check
    err := firmwareManager.DeactivateFirmwareCheck()
    if err != nil {
        return fmt.Errorf("failed to disable firmware check: %w", err)
    }

    // Initialize the ledger instance
    ledgerInstance := &ledger.Ledger{}
    timestamp := time.Now().Format(time.RFC3339)

    // Record the firmware check status in the ledger
    ledgerInstance.AdvancedSecurityLedger.RecordFirmwareCheckStatus("disabled", timestamp)

    log.Println("Firmware check disabled.")
    return nil
}


// DeactivateFirmwareCheck disables firmware integrity checks
func (f *FirmwareCheckManager) DeactivateFirmwareCheck() error {
    f.mutex.Lock()
    defer f.mutex.Unlock()

    if !f.active {
        return errors.New("firmware check is already inactive")
    }

    f.active = false
    log.Printf("Firmware check deactivated at %s", time.Now().Format(time.RFC3339))
    return nil
}
// IsFirmwareCheckActive checks if firmware integrity checks are active
func (f *FirmwareCheckManager) IsFirmwareCheckActive() bool {
    f.mutex.Lock()
    defer f.mutex.Unlock()

    // Log and return the status of the firmware check
    log.Printf("Firmware integrity check active status: %t", f.active)
    return f.active
}

// MonitorAPIUsage tracks the usage patterns of APIs to detect potential abuse.
func MonitorAPIUsage(api string, responseTime float64) error {
    // Validate input
    if api == "" || responseTime <= 0 {
        return fmt.Errorf("invalid API name or response time provided")
    }

    // Initialize APIUsageManager
    usageManager := NewAPIUsageManager()

    // Track the API usage and handle potential errors
    usageStats, err := usageManager.TrackAPIUsage(api, responseTime)
    if err != nil {
        return fmt.Errorf("failed to monitor API usage: %w", err)
    }

    // Record API usage in the ledger
    ledgerInstance := &ledger.Ledger{}

    // Record only API name and total requests in the ledger
    ledgerInstance.AdvancedSecurityLedger.RecordAPIUsage(api, usageStats.TotalRequests)

    log.Printf("API usage for %s monitored successfully. Total Requests: %d", api, usageStats.TotalRequests)
    return nil
}


// TrackAPIUsage monitors usage patterns for a specific API and returns usage statistics
func (a *APIUsageManager) TrackAPIUsage(api string, responseTime float64) (APIUsageStats, error) {
    a.mutex.Lock()
    defer a.mutex.Unlock()

    // Initialize API stats if they don't exist
    if _, exists := a.usageData[api]; !exists {
        a.usageData[api] = &APIUsageStats{
            TotalRequests:     0,
            LastAccessed:      time.Now(),
            AvgResponseTime:   0,
            PeakUsageDetected: false,
            ResponseTimes:     []float64{},
        }
    }

    // Update the API usage stats
    usage := a.usageData[api]
    usage.TotalRequests++
    usage.LastAccessed = time.Now()
    usage.ResponseTimes = append(usage.ResponseTimes, responseTime)
    usage.AvgResponseTime = calculateAvgResponseTime(usage.ResponseTimes)

    // Detect abnormal usage patterns based on thresholds
    if usage.TotalRequests > a.alertThreshold || responseTime > (usage.AvgResponseTime * 2) {
        usage.PeakUsageDetected = true
        log.Printf("Potential abuse detected on API: %s. Total requests: %d, Avg Response Time: %.2f ms, Last Response Time: %.2f ms",
            api, usage.TotalRequests, usage.AvgResponseTime, responseTime)
    } else {
        usage.PeakUsageDetected = false
    }

    log.Printf("API usage stats updated for %s: Total Requests: %d, Avg Response Time: %.2f ms, Last Accessed: %s",
        api, usage.TotalRequests, usage.AvgResponseTime, usage.LastAccessed.Format(time.RFC3339))

    return *usage, nil
}

// calculateAvgResponseTime calculates the average response time based on recorded response times
func calculateAvgResponseTime(responseTimes []float64) float64 {
    if len(responseTimes) == 0 {
        return 0
    }
    total := 0.0
    for _, time := range responseTimes {
        total += time
    }
    return total / float64(len(responseTimes))
}
