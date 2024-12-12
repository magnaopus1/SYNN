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
    LoadBalancingInterval    = 10 * time.Second  // Interval for dynamic load balancing checks
    OverloadThreshold        = 0.75              // Load threshold beyond which balancing is triggered
    DifficultyAdjustmentStep = 0.05              // Step size for dynamic difficulty adjustment
)

// DynamicLoadBalancingAutomation automates load balancing across PoH, PoS, and PoW stages
type DynamicLoadBalancingAutomation struct {
    ledgerInstance   *ledger.Ledger                   // Blockchain ledger for tracking consensus actions
    consensusEngine  *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine
    stateMutex       *sync.RWMutex                    // Mutex for thread-safe ledger access
    apiURL           string                           // API URL for consensus operations
}

// NewDynamicLoadBalancingAutomation initializes the automation for dynamic load balancing in Synnergy Consensus
func NewDynamicLoadBalancingAutomation(apiURL string, ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *DynamicLoadBalancingAutomation {
    return &DynamicLoadBalancingAutomation{
        ledgerInstance:  ledgerInstance,
        consensusEngine: consensusEngine,
        stateMutex:      stateMutex,
        apiURL:          apiURL,
    }
}

// StartLoadBalancing initiates the continuous monitoring and balancing of loads across PoH, PoS, and PoW stages
func (automation *DynamicLoadBalancingAutomation) StartLoadBalancing() {
    ticker := time.NewTicker(LoadBalancingInterval)
    go func() {
        for range ticker.C {
            automation.monitorAndBalanceLoad()
        }
    }()
}

// monitorAndBalanceLoad checks system load and dynamically balances tasks across the consensus stages
func (automation *DynamicLoadBalancingAutomation) monitorAndBalanceLoad() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Monitor the system load for each consensus stage
    pohLoad := automation.consensusEngine.GetPoHLoad()
    posLoad := automation.consensusEngine.GetPoSLoad()
    powLoad := automation.consensusEngine.GetPoWLoad()

    fmt.Printf("Monitoring load: PoH Load = %.2f, PoS Load = %.2f, PoW Load = %.2f\n", pohLoad, posLoad, powLoad)

    // Balance the load dynamically based on thresholds
    if pohLoad > OverloadThreshold {
        fmt.Println("PoH overloaded, distributing proof generation tasks.")
        automation.distributePoHTasks()
    }
    if posLoad > OverloadThreshold {
        fmt.Println("PoS overloaded, allocating more stake to validators.")
        automation.addStakeToValidators()
    }
    if powLoad > OverloadThreshold {
        fmt.Println("PoW overloaded, adjusting difficulty.")
        automation.adjustPoWDifficulty(powLoad)
    }

    // Log the load balancing action in the ledger
    automation.logLoadBalancingAction(pohLoad, posLoad, powLoad)
}

// distributePoHTasks distributes PoH proof generation tasks across multiple nodes
func (automation *DynamicLoadBalancingAutomation) distributePoHTasks() {
    url := fmt.Sprintf("%s/api/consensus/poh/generate-multiple", automation.apiURL)
    resp, err := http.Post(url, "application/json", nil)
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error distributing PoH tasks.")
        return
    }
    fmt.Println("PoH proof generation distributed across nodes.")
}

// addStakeToValidators adds more stake to overloaded validators in PoS
func (automation *DynamicLoadBalancingAutomation) addStakeToValidators() {
    url := fmt.Sprintf("%s/api/consensus/pos/add-stake", automation.apiURL)
    payload := map[string]float64{"additional_stake": 100.0} // Example stake adjustment
    reqBody, _ := json.Marshal(payload)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error adding stake to validators.")
        return
    }
    fmt.Println("Stake added to validators to balance PoS load.")
}

// adjustPoWDifficulty adjusts PoW difficulty based on the current load
func (automation *DynamicLoadBalancingAutomation) adjustPoWDifficulty(currentLoad float64) {
    adjustmentFactor := DifficultyAdjustmentStep * (currentLoad - OverloadThreshold)
    url := fmt.Sprintf("%s/api/consensus/pow/adjust-difficulty", automation.apiURL)
    payload := map[string]float64{"adjustment_factor": adjustmentFactor}
    reqBody, _ := json.Marshal(payload)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil || resp.StatusCode != http.StatusOK {
        fmt.Println("Error adjusting PoW difficulty.")
        return
    }
    fmt.Printf("PoW difficulty adjusted by factor: %.2f\n", adjustmentFactor)
}

// logLoadBalancingAction logs the load balancing action and stores it in the ledger
func (automation *DynamicLoadBalancingAutomation) logLoadBalancingAction(pohLoad, posLoad, powLoad float64) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    ledgerEntry := common.LedgerEntry{
        ID:        fmt.Sprintf("load-balancing-log-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Load Balancing",
        Status:    "Executed",
        Details:   fmt.Sprintf("PoH Load: %.2f, PoS Load: %.2f, PoW Load: %.2f", pohLoad, posLoad, powLoad),
    }

    // Encrypt the ledger entry for security purposes
    encryptedEntry, err := encryption.EncryptLedgerEntry(ledgerEntry, []byte(EncryptionKey))
    if err != nil {
        fmt.Printf("Error encrypting load balancing log: %v\n", err)
        return
    }

    automation.ledgerInstance.AddEntry(encryptedEntry)
    fmt.Println("Load balancing action logged in the ledger.")
}

// Additional helper function to ensure system efficiency and load balancing over time
func (automation *DynamicLoadBalancingAutomation) ensureSystemEfficiency() {
    fmt.Println("Ensuring system efficiency and optimal load balancing...")

    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Re-validate chain consistency after load balancing
    err := automation.consensusEngine.ValidateChain()
    if err != nil {
        fmt.Printf("Chain validation failed: %v\n", err)
        automation.monitorAndBalanceLoad() // Trigger load balancing if the system is inefficient
    } else {
        fmt.Println("System load balanced and chain is consistent.")
    }
}
