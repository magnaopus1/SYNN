package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/network"
)

const (
    DDoSMonitoringInterval   = 5 * time.Second // Interval for checking network traffic for DDoS attacks
    MaxAllowedRequestRate    = 1000            // Maximum number of allowed requests per second per node
    AttackMitigationDuration = 10 * time.Minute // Duration for DDoS mitigation measures
)

// DDoSProtectionAutomation automates the process of detecting and responding to DDoS attacks
type DDoSProtectionAutomation struct {
    consensusSystem  *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance   *ledger.Ledger               // Ledger for logging DDoS events
    stateMutex       *sync.RWMutex                // Mutex for thread-safe access
    mitigationActive bool                         // Flag to track whether mitigation measures are active
}

// NewDDoSProtectionAutomation initializes the automation for DDoS protection
func NewDDoSProtectionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DDoSProtectionAutomation {
    return &DDoSProtectionAutomation{
        consensusSystem:  consensusSystem,
        ledgerInstance:   ledgerInstance,
        stateMutex:       stateMutex,
        mitigationActive: false,
    }
}

// StartDDoSMonitoring starts the continuous loop for monitoring network traffic for DDoS attacks
func (automation *DDoSProtectionAutomation) StartDDoSMonitoring() {
    ticker := time.NewTicker(DDoSMonitoringInterval)

    go func() {
        for range ticker.C {
            automation.monitorNetworkForDDoS()
        }
    }()
}

// monitorNetworkForDDoS checks network traffic for potential DDoS attacks
func (automation *DDoSProtectionAutomation) monitorNetworkForDDoS() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch network traffic statistics from all nodes
    trafficData := automation.consensusSystem.GetNetworkTrafficData()

    for _, nodeTraffic := range trafficData {
        if nodeTraffic.RequestRate > MaxAllowedRequestRate {
            fmt.Printf("DDoS attack detected on node %s. Request rate: %d requests/sec.\n", nodeTraffic.NodeID, nodeTraffic.RequestRate)
            automation.triggerMitigation(nodeTraffic.NodeID)
            automation.logDDoSEvent(nodeTraffic.NodeID, nodeTraffic.RequestRate)
        }
    }
}

// triggerMitigation activates DDoS mitigation measures
func (automation *DDoSProtectionAutomation) triggerMitigation(nodeID string) {
    if !automation.mitigationActive {
        fmt.Printf("Activating DDoS mitigation for node %s.\n", nodeID)

        // Activate network throttling or other mitigation measures
        automation.consensusSystem.ActivateDDoSMitigation(nodeID, AttackMitigationDuration)

        automation.mitigationActive = true

        // Schedule deactivation of mitigation measures after the set duration
        time.AfterFunc(AttackMitigationDuration, func() {
            automation.deactivateMitigation(nodeID)
        })
    } else {
        fmt.Println("DDoS mitigation is already active.")
    }
}

// deactivateMitigation deactivates the DDoS mitigation measures
func (automation *DDoSProtectionAutomation) deactivateMitigation(nodeID string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    fmt.Printf("Deactivating DDoS mitigation for node %s.\n", nodeID)
    automation.consensusSystem.DeactivateDDoSMitigation(nodeID)

    automation.mitigationActive = false
}

// logDDoSEvent logs the DDoS attack event into the ledger
func (automation *DDoSProtectionAutomation) logDDoSEvent(nodeID string, requestRate int) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ddos-attack-%s-%d", nodeID, time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "DDoS Attack",
        Status:    "Mitigation Activated",
        Details:   fmt.Sprintf("DDoS attack detected on node %s with a request rate of %d requests/sec.", nodeID, requestRate),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with DDoS attack event for node %s.\n", nodeID)
}

// ensureNetworkIntegrity checks the integrity of the network after mitigation
func (automation *DDoSProtectionAutomation) ensureNetworkIntegrity() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    integrityValid := automation.consensusSystem.ValidateNetworkIntegrity()
    if !integrityValid {
        fmt.Println("Network integrity breach detected after DDoS mitigation. Investigating further.")
        automation.triggerMitigationResolution()
    } else {
        fmt.Println("Network integrity is valid after DDoS mitigation.")
    }
}

// triggerMitigationResolution handles the resolution of DDoS attack consequences
func (automation *DDoSProtectionAutomation) triggerMitigationResolution() {
    fmt.Println("Resolving issues caused by DDoS attack...")
    // Add custom logic to restore normal operation, if necessary
}
