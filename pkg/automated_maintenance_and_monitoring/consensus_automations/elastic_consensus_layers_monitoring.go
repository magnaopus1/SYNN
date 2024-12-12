package consensus_automations

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    ElasticityCheckInterval   = 10 * time.Second  // Interval for elasticity monitoring
    ElasticityThreshold       = 0.8               // Threshold for triggering elasticity adjustments
    AdjustmentStep            = 0.05              // Step size for dynamic adjustments
    ElasticityLogKey          = "elasticity_monitoring_log_key" // Encryption key for elasticity logs
)

// ElasticConsensusLayersMonitoringAutomation automates the monitoring of PoH, PoS, and PoW elasticity
type ElasticConsensusLayersMonitoringAutomation struct {
    ledgerInstance  *ledger.Ledger                   // Blockchain ledger for tracking consensus actions
    consensusEngine *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    stateMutex      *sync.RWMutex                    // Mutex for thread-safe ledger access
    apiURL          string                           // API URL for consensus operations
}

// NewElasticConsensusLayersMonitoringAutomation initializes the automation for monitoring consensus elasticity
func NewElasticConsensusLayersMonitoringAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *ElasticConsensusLayersMonitoringAutomation {
    return &ElasticConsensusLayersMonitoringAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartElasticityMonitoring begins monitoring the elasticity of PoH, PoS, and PoW layers
func (automation *ElasticConsensusLayersMonitoringAutomation) StartElasticityMonitoring() {
    ticker := time.NewTicker(ElasticityCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorElasticity()
        }
    }()
}

// monitorElasticity checks the elasticity of PoH, PoS, and PoW layers and adjusts them if needed
func (automation *ElasticConsensusLayersMonitoringAutomation) monitorElasticity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the elasticity metrics for each consensus layer
    pohElasticity := automation.consensusEngine.MeasurePoHElasticity()
    posElasticity := automation.consensusEngine.MeasurePoSElasticity()
    powElasticity := automation.consensusEngine.MeasurePoWElasticity()

    fmt.Printf("Monitoring elasticity: PoH = %.2f, PoS = %.2f, PoW = %.2f\n", pohElasticity, posElasticity, powElasticity)

    // Adjust elasticity if any layer exceeds the threshold
    if pohElasticity > ElasticityThreshold {
        fmt.Println("PoH elasticity exceeded threshold, adjusting PoH layer.")
        automation.adjustPoHElasticity(pohElasticity)
    }
    if posElasticity > ElasticityThreshold {
        fmt.Println("PoS elasticity exceeded threshold, adjusting PoS layer.")
        automation.adjustPoSElasticity(posElasticity)
    }
    if powElasticity > ElasticityThreshold {
        fmt.Println("PoW elasticity exceeded threshold, adjusting PoW layer.")
        automation.adjustPoWElasticity(powElasticity)
    }

    // Log elasticity adjustments in the ledger
    automation.logElasticityMetrics(pohElasticity, posElasticity, powElasticity)
}

// adjustPoHElasticity adjusts the PoH elasticity based on the measured load
func (automation *ElasticConsensusLayersMonitoringAutomation) adjustPoHElasticity(currentElasticity float64) {
    adjustmentFactor := AdjustmentStep * (currentElasticity - ElasticityThreshold)
    url := fmt.Sprintf("%s/api/consensus/poh/generate-multiple", automation.apiURL)

    payload := map[string]float64{"adjustment_factor": adjustmentFactor}
    reqBody, _ := json.Marshal(payload)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error adjusting PoH elasticity.")
        return
    }
    fmt.Printf("PoH elasticity adjusted by factor: %.2f\n", adjustmentFactor)
}

// adjustPoSElasticity adjusts the PoS elasticity by adding/removing validators or adjusting their stakes
func (automation *ElasticConsensusLayersMonitoringAutomation) adjustPoSElasticity(currentElasticity float64) {
    adjustmentFactor := AdjustmentStep * (currentElasticity - ElasticityThreshold)
    url := fmt.Sprintf("%s/api/consensus/pos/add-stake", automation.apiURL)

    payload := map[string]float64{"additional_stake": adjustmentFactor * 100} // Example stake adjustment
    reqBody, _ := json.Marshal(payload)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error adjusting PoS elasticity.")
        return
    }
    fmt.Printf("PoS elasticity adjusted by increasing validator stake by factor: %.2f\n", adjustmentFactor)
}

// adjustPoWElasticity adjusts the PoW elasticity by modifying mining difficulty
func (automation *ElasticConsensusLayersMonitoringAutomation) adjustPoWElasticity(currentElasticity float64) {
    adjustmentFactor := AdjustmentStep * (currentElasticity - ElasticityThreshold)
    url := fmt.Sprintf("%s/api/consensus/pow/adjust-difficulty", automation.apiURL)

    payload := map[string]float64{"adjustment_factor": adjustmentFactor}
    reqBody, _ := json.Marshal(payload)
    
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error adjusting PoW elasticity.")
        return
    }
    fmt.Printf("PoW mining difficulty adjusted by factor: %.2f\n", adjustmentFactor)
}

// logElasticityMetrics logs the elasticity metrics into the ledger
func (automation *ElasticConsensusLayersMonitoringAutomation) logElasticityMetrics(pohElasticity, posElasticity, powElasticity float64) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    ledgerEntry := common.LedgerEntry{
        ID:        fmt.Sprintf("elasticity-monitoring-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Elasticity Monitoring",
        Status:    "Logged",
        Details:   fmt.Sprintf("PoH Elasticity: %.2f, PoS Elasticity: %.2f, PoW Elasticity: %.2f", pohElasticity, posElasticity, powElasticity),
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(ledgerEntry, []byte(ElasticityLogKey))
    if err != nil {
        fmt.Printf("Error encrypting elasticity log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Elasticity monitoring data logged in the ledger.")
}

// Additional helper function to ensure consistent performance of the consensus system
func (automation *ElasticConsensusLayersMonitoringAutomation) ensurePerformanceConsistency() {
    fmt.Println("Ensuring consensus performance consistency...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Validate chain consistency post-elasticity adjustments
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Printf("Chain validation failed: %v\n", err)
        automation.monitorElasticity() // Trigger monitoring again if consistency fails
    } else {
        fmt.Println("Chain and consensus layers are consistent post-elasticity adjustments.")
    }
}
