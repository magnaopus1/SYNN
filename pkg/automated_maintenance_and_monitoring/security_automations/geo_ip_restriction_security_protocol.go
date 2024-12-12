package security_automations

import (
    "fmt"
    "net"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
)

const (
    GeoIPCheckInterval        = 30 * time.Second // Interval for checking geo IP restrictions
    MaxGeoViolationThreshold  = 5                // Maximum allowed geo IP violations before action
    GeoBlockEnforcementPeriod = 24 * time.Hour   // Period for blocking restricted geo IP addresses
)

// GeoIPRestrictionAutomation ensures that transactions are restricted based on geographic IP policies
type GeoIPRestrictionAutomation struct {
    consensusSystem   *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance    *ledger.Ledger               // Ledger for logging geo IP restriction events
    stateMutex        *sync.RWMutex                // Mutex for thread-safe access
    restrictedRegions map[string]bool              // Map of restricted countries (based on ISO country codes)
    violationTracker  map[string]int               // Tracks geo IP violations per address
    blockedIPs        map[string]time.Time         // Blocked IPs with timestamps for enforcement period
}

// NewGeoIPRestrictionAutomation initializes the geo IP restriction automation
func NewGeoIPRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *GeoIPRestrictionAutomation {
    return &GeoIPRestrictionAutomation{
        consensusSystem:   consensusSystem,
        ledgerInstance:    ledgerInstance,
        stateMutex:        stateMutex,
        restrictedRegions: make(map[string]bool),
        violationTracker:  make(map[string]int),
        blockedIPs:        make(map[string]time.Time),
    }
}

// StartGeoIPMonitoring starts the continuous loop for geo IP restriction enforcement
func (automation *GeoIPRestrictionAutomation) StartGeoIPMonitoring() {
    ticker := time.NewTicker(GeoIPCheckInterval)

    go func() {
        for range ticker.C {
            automation.monitorGeoIPRestrictions()
        }
    }()
}

// monitorGeoIPRestrictions checks all incoming transactions for geo IP violations
func (automation *GeoIPRestrictionAutomation) monitorGeoIPRestrictions() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    transactionList := automation.consensusSystem.GetPendingTransactions()

    if len(transactionList) > 0 {
        for _, transaction := range transactionList {
            ip := automation.consensusSystem.GetIPForTransaction(transaction)
            if automation.isGeoIPRestricted(ip) {
                fmt.Printf("Geo IP restriction violation for transaction %s from IP %s.\n", transaction.ID, ip)
                automation.handleGeoIPViolation(transaction, ip)
            }
        }
    } else {
        fmt.Println("No pending transactions at this time.")
    }

    automation.enforceGeoIPBlocks()
}

// isGeoIPRestricted checks if an IP address is located in a restricted region
func (automation *GeoIPRestrictionAutomation) isGeoIPRestricted(ip string) bool {
    // Here, a mock IP-to-country resolution is used. In real-world scenarios, integrate with an external geo IP service
    countryCode := automation.resolveCountryCodeFromIP(ip)
    return automation.restrictedRegions[countryCode]
}

// resolveCountryCodeFromIP resolves a country code from an IP address (mock function for demo)
func (automation *GeoIPRestrictionAutomation) resolveCountryCodeFromIP(ip string) string {
    // Example mock data. In practice, use a geo IP API or database to resolve IP to country code.
    mockIPToCountry := map[string]string{
        "192.168.0.1": "US", // United States
        "10.0.0.1":   "CN", // China
        "172.16.0.1": "RU", // Russia
    }

    if countryCode, found := mockIPToCountry[ip]; found {
        return countryCode
    }
    return "UNKNOWN"
}

// handleGeoIPViolation handles actions taken when a geo IP violation occurs
func (automation *GeoIPRestrictionAutomation) handleGeoIPViolation(transaction common.Transaction, ip string) {
    automation.violationTracker[ip]++

    if automation.violationTracker[ip] >= MaxGeoViolationThreshold {
        automation.blockIPAddress(ip)
    }

    automation.logGeoIPViolation(transaction, ip)
}

// blockIPAddress adds an IP address to the blocked list for a defined enforcement period
func (automation *GeoIPRestrictionAutomation) blockIPAddress(ip string) {
    automation.blockedIPs[ip] = time.Now().Add(GeoBlockEnforcementPeriod)
    fmt.Printf("IP address %s blocked for geo IP restriction violations.\n", ip)
    automation.logIPBlock(ip)
}

// enforceGeoIPBlocks checks if blocked IP addresses are still within the enforcement period
func (automation *GeoIPRestrictionAutomation) enforceGeoIPBlocks() {
    now := time.Now()

    for ip, unblockTime := range automation.blockedIPs {
        if now.After(unblockTime) {
            delete(automation.blockedIPs, ip)
            fmt.Printf("IP address %s unblocked after enforcement period.\n", ip)
            automation.logIPUnblock(ip)
        }
    }
}

// logGeoIPViolation logs a geo IP violation into the ledger
func (automation *GeoIPRestrictionAutomation) logGeoIPViolation(transaction common.Transaction, ip string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("geo-ip-violation-%s", transaction.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Geo IP Violation",
        Status:    "Violation",
        Details:   fmt.Sprintf("Transaction %s from IP %s violated geo IP restrictions.", transaction.ID, ip),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with geo IP violation for transaction %s.\n", transaction.ID)
}

// logIPBlock logs an IP block event into the ledger
func (automation *GeoIPRestrictionAutomation) logIPBlock(ip string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ip-block-%s", ip),
        Timestamp: time.Now().Unix(),
        Type:      "IP Block",
        Status:    "Blocked",
        Details:   fmt.Sprintf("IP address %s was blocked due to repeated geo IP violations.", ip),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with IP block event for IP %s.\n", ip)
}

// logIPUnblock logs an IP unblock event into the ledger
func (automation *GeoIPRestrictionAutomation) logIPUnblock(ip string) {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ip-unblock-%s", ip),
        Timestamp: time.Now().Unix(),
        Type:      "IP Unblock",
        Status:    "Unblocked",
        Details:   fmt.Sprintf("IP address %s was unblocked after enforcement period expired.", ip),
    }
    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with IP unblock event for IP %s.\n", ip)
}
