package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

const (
	AssetListingCheckInterval = 1 * time.Second  // Interval for checking asset listings
	IllegalAssetFlagLimit     = 5                // Number of reports before asset flagged
	AssetOwnershipValidation  = true             // Ensure asset ownership before listing
)

// AssetListingRestrictionAutomation manages asset listings and enforces restrictions across the network
type AssetListingRestrictionAutomation struct {
	consensusSystem        *consensus.SynnergyConsensus
	ledgerInstance         *ledger.Ledger
	stateMutex             *sync.RWMutex
	flaggedAssetReports    map[string]int  // Tracks the number of flags per asset
	illegalAssetCategories []string        // Categories of assets restricted from listing
}

// NewAssetListingRestrictionAutomation initializes and returns an instance of AssetListingRestrictionAutomation
func NewAssetListingRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex, illegalCategories []string) *AssetListingRestrictionAutomation {
	return &AssetListingRestrictionAutomation{
		consensusSystem:        consensusSystem,
		ledgerInstance:         ledgerInstance,
		stateMutex:             stateMutex,
		flaggedAssetReports:    make(map[string]int),
		illegalAssetCategories: illegalCategories,
	}
}

// StartAssetListingMonitoring starts real-time monitoring of asset listings for compliance and restrictions
func (automation *AssetListingRestrictionAutomation) StartAssetListingMonitoring() {
	ticker := time.NewTicker(AssetListingCheckInterval)

	go func() {
		for range ticker.C {
			automation.checkAssetListings()
		}
	}()
}

// checkAssetListings continuously checks new asset listings and enforces restrictions based on the network’s protocols
func (automation *AssetListingRestrictionAutomation) checkAssetListings() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Retrieve the recent asset listings from the consensus
	recentListings := automation.consensusSystem.GetRecentAssetListings()

	for _, asset := range recentListings {
		// Validate asset ownership and check for restricted categories
		if !automation.validateAssetOwnership(asset) {
			automation.flagIllegalAsset(asset, "Ownership validation failed")
			continue
		}

		if automation.isRestrictedCategory(asset.Category) {
			automation.flagIllegalAsset(asset, "Restricted category")
		}
	}
}

// validateAssetOwnership checks if the user owns the asset being listed
func (automation *AssetListingRestrictionAutomation) validateAssetOwnership(asset common.Asset) bool {
	if AssetOwnershipValidation {
		// Validate the ownership of the asset through Synnergy Consensus
		return automation.consensusSystem.ValidateAssetOwnership(asset.Owner, asset.ID)
	}
	return true
}

// isRestrictedCategory checks if the asset belongs to a restricted category
func (automation *AssetListingRestrictionAutomation) isRestrictedCategory(category string) bool {
	for _, restrictedCategory := range automation.illegalAssetCategories {
		if category == restrictedCategory {
			return true
		}
	}
	return false
}

// flagIllegalAsset flags an asset listing that violates listing protocols and logs it in the ledger
func (automation *AssetListingRestrictionAutomation) flagIllegalAsset(asset common.Asset, reason string) {
	fmt.Printf("Illegal asset detected: %s, Reason: %s\n", asset.ID, reason)

	// Track the number of reports for this asset
	automation.flaggedAssetReports[asset.ID]++

	if automation.flaggedAssetReports[asset.ID] >= IllegalAssetFlagLimit {
		fmt.Printf("Asset %s flagged for repeated violations.\n", asset.ID)
		automation.logViolationToLedger(asset, reason)
	}
}

// logViolationToLedger logs the flagged asset violation in the ledger for auditing and transparency
func (automation *AssetListingRestrictionAutomation) logViolationToLedger(asset common.Asset, violationReason string) {
	// Encrypt asset details before logging
	encryptedData := automation.encryptAssetData(asset)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("asset-violation-%s-%d", asset.ID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Asset Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("Asset (%s) flagged for violation. Reason: %s. Encrypted Data: %s", asset.ID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log asset violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Asset violation logged for asset: %s\n", asset.ID)
	}
}

// encryptAssetData encrypts asset data before logging for security and compliance
func (automation *AssetListingRestrictionAutomation) encryptAssetData(asset common.Asset) string {
	data := fmt.Sprintf("ID: %s, Owner: %s, Category: %s, Price: %.2f", asset.ID, asset.Owner, asset.Category, asset.Price)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting asset data:", err)
		return data
	}
	return string(encryptedData)
}
