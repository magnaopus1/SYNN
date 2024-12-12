package execution_automations

import (
	"fmt"
	"sync"
	"time"
	"net"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/encryption"
	"synnergy_network_demo/synnergy_consensus"
	"synnergy_network_demo/network_monitoring"
)

const (
	DDoSCheckInterval         = 15 * time.Second  // Interval for checking DDoS attacks
	RequestThresholdPerSecond = 1000              // Max allowed requests per second before action
	BlacklistDuration         = 30 * time.Minute  // Time an IP is blacklisted after detection
	WhitelistDuration         = 24 * time.Hour    // Duration of whitelist before reevaluation
)

// DDoSMitigationAutomation monitors and mitigates DDoS attacks
type DDoSMitigationAutomation struct {
	consensusEngine    *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for validation
	ledgerInstance     *ledger.Ledger                        // Ledger for logging DDoS mitigation events
	networkMonitor     *network_monitoring.Monitor           // Monitor network activity for unusual patterns
	blacklist          map[string]time.Time                  // Blacklist of malicious IPs and their expiry times
	whitelist          map[string]time.Time                  // Whitelist of trusted IPs
	ddosMutex          *sync.RWMutex                         // Mutex for thread-safe operations
}

// NewDDoSMitigationAutomation initializes the DDoS mitigation automation
func NewDDoSMitigationAutomation(consensusEngine *synnergy_consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, networkMonitor *network_monitoring.Monitor, ddosMutex *sync.RWMutex) *DDoSMitigationAutomation {
	return &DDoSMitigationAutomation{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		networkMonitor:  networkMonitor,
		blacklist:       make(map[string]time.Time),
		whitelist:       make(map[string]time.Time),
		ddosMutex:       ddosMutex,
	}
}

// StartDDoSMitigation begins monitoring the network and mitigating DDoS threats
func (automation *DDoSMitigationAutomation) StartDDoSMitigation() {
	ticker := time.NewTicker(DDoSCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorNetworkForDDoS()
		}
	}()
}

// monitorNetworkForDDoS checks network traffic for signs of DDoS attacks and responds accordingly
func (automation *DDoSMitigationAutomation) monitorNetworkForDDoS() {
	automation.ddosMutex.Lock()
	defer automation.ddosMutex.Unlock()

	// Get the current network traffic statistics
	trafficStats := automation.networkMonitor.GetTrafficStats()

	// Analyze the traffic and detect potential DDoS attacks
	for ip, reqsPerSecond := range trafficStats.RequestsPerSecond {
		if reqsPerSecond > RequestThresholdPerSecond {
			// If IP is not whitelisted and exceeds threshold, blacklist it
			if _, whitelisted := automation.whitelist[ip]; !whitelisted {
				automation.blacklistIP(ip)
			}
		}
	}

	// Clean up expired blacklisted IPs
	automation.cleanupExpiredBlacklistEntries()

	// Log the network traffic analysis
	automation.logDDoSMonitoringEvent(trafficStats)
}

// blacklistIP adds an IP address to the blacklist and blocks it from further requests
func (automation *DDoSMitigationAutomation) blacklistIP(ip string) {
	// Blacklist the IP for the duration
	automation.blacklist[ip] = time.Now().Add(BlacklistDuration)

	// Log the blacklist action in the ledger
	automation.logDDoSAttack(ip)

	// Broadcast the IP blacklist to the consensus network
	automation.consensusEngine.BlockIP(ip)

	fmt.Printf("Blacklisted IP: %s due to excessive requests.\n", ip)
}

// cleanupExpiredBlacklistEntries removes expired entries from the blacklist
func (automation *DDoSMitigationAutomation) cleanupExpiredBlacklistEntries() {
	now := time.Now()

	for ip, expiry := range automation.blacklist {
		if now.After(expiry) {
			delete(automation.blacklist, ip)
			automation.consensusEngine.UnblockIP(ip)
			fmt.Printf("Unblocked IP: %s after blacklist expiry.\n", ip)
		}
	}
}

// logDDoSAttack logs a DDoS attack and mitigation action in the ledger
func (automation *DDoSMitigationAutomation) logDDoSAttack(ip string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ddos-attack-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DDoS Mitigation",
		Status:    "IP Blacklisted",
		Details:   fmt.Sprintf("IP %s was blacklisted due to exceeding request threshold.", ip),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log DDoS mitigation:", err)
	} else {
		fmt.Println("DDoS attack mitigation logged in the ledger.")
	}
}

// logDDoSMonitoringEvent logs the results of DDoS monitoring in the ledger
func (automation *DDoSMitigationAutomation) logDDoSMonitoringEvent(trafficStats network_monitoring.TrafficStats) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("ddos-monitoring-%d", time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DDoS Monitoring",
		Status:    "Success",
		Details:   fmt.Sprintf("Monitored network traffic. Total requests: %d", trafficStats.TotalRequests),
	}

	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log DDoS monitoring:", err)
	} else {
		fmt.Println("DDoS monitoring successfully logged in the ledger.")
	}
}

// encryptData encrypts sensitive data before logging it in the ledger
func (automation *DDoSMitigationAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return data
	}
	return string(encryptedData)
}
