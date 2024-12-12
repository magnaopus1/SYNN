package automations

import (
	"fmt"
	"net"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	GeographyCheckInterval    = 1 * time.Minute // Interval for checking geographical access restrictions
	RestrictedCountries       = "IL,IR,RU,KP"  // ISO codes for restricted countries: Israel, Iran, Russia, North Korea
	MaxAllowedAccessAttempts  = 3              // Maximum allowed unauthorized access attempts before restriction
)

// RestrictedGeographicalAccessAutomation enforces geographical restrictions on the network
type RestrictedGeographicalAccessAutomation struct {
	consensusSystem       *consensus.SynnergyConsensus
	ledgerInstance        *ledger.Ledger
	stateMutex            *sync.RWMutex
	accessViolationCount  map[string]int // Tracks access violations by IP or country
}

// NewRestrictedGeographicalAccessAutomation initializes and returns an instance of RestrictedGeographicalAccessAutomation
func NewRestrictedGeographicalAccessAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *RestrictedGeographicalAccessAutomation {
	return &RestrictedGeographicalAccessAutomation{
		consensusSystem:      consensusSystem,
		ledgerInstance:       ledgerInstance,
		stateMutex:           stateMutex,
		accessViolationCount: make(map[string]int),
	}
}

// StartGeoAccessMonitoring starts continuous monitoring of geographical access violations
func (automation *RestrictedGeographicalAccessAutomation) StartGeoAccessMonitoring() {
	ticker := time.NewTicker(GeographyCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorGeoAccess()
		}
	}()
}

// monitorGeoAccess checks for geographical access violations and restricts access if necessary
func (automation *RestrictedGeographicalAccessAutomation) monitorGeoAccess() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch geographical access data from Synnergy Consensus
	accessData := automation.consensusSystem.GetGeoAccessData()

	for ip, countryCode := range accessData {
		// Check if the country is restricted
		if automation.isRestrictedCountry(countryCode) {
			automation.flagGeoViolation(ip, countryCode, "Unauthorized geographical access attempt detected")
		}
	}
}

// isRestrictedCountry checks if the given country is restricted from network access
func (automation *RestrictedGeographicalAccessAutomation) isRestrictedCountry(countryCode string) bool {
	// Check if the country code is in the restricted list
	return contains(RestrictedCountries, countryCode)
}

// flagGeoViolation flags an IP or country for geographical access violations and logs it in the ledger
func (automation *RestrictedGeographicalAccessAutomation) flagGeoViolation(ip string, countryCode string, reason string) {
	fmt.Printf("Geographical access violation: IP %s, Country %s, Reason: %s\n", ip, countryCode, reason)

	// Increment the violation count for the IP or country
	automation.accessViolationCount[ip]++

	// Log the violation in the ledger
	automation.logGeoViolation(ip, countryCode, reason)

	// Check if the IP or country has exceeded the allowed number of access violations
	if automation.accessViolationCount[ip] >= MaxAllowedAccessAttempts {
		automation.restrictGeoAccess(ip, countryCode)
	}
}

// logGeoViolation logs the flagged geographical access violation into the ledger with full details
func (automation *RestrictedGeographicalAccessAutomation) logGeoViolation(ip string, countryCode string, violationReason string) {
	// Create a ledger entry for geo access violation
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("geo-access-violation-%s-%d", ip, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Geographical Access Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("IP %s from country %s violated access restrictions. Reason: %s", ip, countryCode, violationReason),
	}

	// Encrypt the log data before adding it to the ledger
	encryptedDetails := automation.encryptGeoData(entry.Details)
	entry.Details = encryptedDetails

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log geographical access violation:", err)
	} else {
		fmt.Println("Geographical access violation logged.")
	}
}

// restrictGeoAccess restricts access for an IP after exceeding allowed violations
func (automation *RestrictedGeographicalAccessAutomation) restrictGeoAccess(ip string, countryCode string) {
	// Add restriction details to the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("geo-access-restriction-%s-%d", ip, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Geographical Access Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("IP %s from country %s has been restricted due to repeated access violations.", ip, countryCode),
	}

	// Encrypt the restriction details before adding it to the ledger
	encryptedDetails := automation.encryptGeoData(entry.Details)
	entry.Details = encryptedDetails

	// Add the restriction entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log geographical access restriction:", err)
	} else {
		fmt.Println("Geographical access restriction applied.")
	}
}

// encryptGeoData encrypts the geographical access data before logging for security
func (automation *RestrictedGeographicalAccessAutomation) encryptGeoData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting geo access data:", err)
		return data
	}
	return string(encryptedData)
}

// contains checks if a string is present in a comma-separated list
func contains(list string, element string) bool {
	return strings.Contains(list, element)
}
