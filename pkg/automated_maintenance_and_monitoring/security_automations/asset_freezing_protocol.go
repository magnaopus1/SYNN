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
    FreezeCheckInterval    = 5 * time.Minute // Interval for checking assets to freeze/unfreeze
    SubBlocksPerBlock      = 1000            // Number of sub-blocks in a block
    EmergencyUnfreezeDelay = 10 * time.Second // Delay for emergency unfreezing actions
)

// AssetFreezingUnfreezingProtocolAutomation automates the process of freezing/unfreezing assets based on consensus and conditions
type AssetFreezingUnfreezingProtocolAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging asset freeze/unfreeze actions
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    freezeCycleCount   int                          // Counter for asset freeze cycles
    emergencyTriggered bool                         // Indicator for emergency unfreeze trigger
}

// NewAssetFreezingUnfreezingProtocolAutomation initializes the automation for freezing/unfreezing assets
func NewAssetFreezingUnfreezingProtocolAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *AssetFreezingUnfreezingProtocolAutomation {
    return &AssetFreezingUnfreezingProtocolAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        freezeCycleCount:   0,
        emergencyTriggered: false,
    }
}

// StartAssetFreezeUnfreeze starts the continuous loop for checking and freezing/unfreezing assets
func (automation *AssetFreezingUnfreezingProtocolAutomation) StartAssetFreezeUnfreeze() {
    ticker := time.NewTicker(FreezeCheckInterval)

    go func() {
        for range ticker.C {
            automation.checkAndFreezeUnfreezeAssets()
        }
    }()
}

// checkAndFreezeUnfreezeAssets checks the assets that meet the criteria for freezing or unfreezing and performs the actions
func (automation *AssetFreezingUnfreezingProtocolAutomation) checkAndFreezeUnfreezeAssets() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch the list of assets that need to be frozen/unfrozen based on consensus conditions
    assetList := automation.consensusSystem.GetAssetsForFreezeUnfreeze()

    if len(assetList) > 0 {
        for _, asset := range assetList {
            if asset.NeedsFreezing {
                fmt.Printf("Freezing asset %s.\n", asset.ID)
                automation.freezeAsset(asset)
            } else if asset.NeedsUnfreezing {
                fmt.Printf("Unfreezing asset %s.\n", asset.ID)
                automation.unfreezeAsset(asset)
            }
        }
    } else {
        fmt.Println("No assets need freezing or unfreezing at this time.")
    }

    automation.freezeCycleCount++
    fmt.Printf("Asset freeze/unfreeze cycle #%d executed.\n", automation.freezeCycleCount)

    if automation.freezeCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeFreezeUnfreezeCycle()
    }
}

// freezeAsset encrypts asset details and triggers the freezing process
func (automation *AssetFreezingUnfreezingProtocolAutomation) freezeAsset(asset common.Asset) {
    encryptedAssetData := automation.encryptAssetData(asset)
    freezeSuccess := automation.consensusSystem.FreezeAsset(asset, encryptedAssetData)

    if freezeSuccess {
        fmt.Printf("Asset %s frozen successfully.\n", asset.ID)
        automation.logAssetFreezeEvent(asset)
    } else {
        fmt.Printf("Error freezing asset %s.\n", asset.ID)
    }
}

// unfreezeAsset triggers the unfreezing process for an asset
func (automation *AssetFreezingUnfreezingProtocolAutomation) unfreezeAsset(asset common.Asset) {
    unfreezeSuccess := automation.consensusSystem.UnfreezeAsset(asset)

    if unfreezeSuccess {
        fmt.Printf("Asset %s unfrozen successfully.\n", asset.ID)
        automation.logAssetUnfreezeEvent(asset)
    } else {
        fmt.Printf("Error unfreezing asset %s.\n", asset.ID)
    }
}

// emergencyUnfreeze manually unfreezes an asset in case of an emergency
func (automation *AssetFreezingUnfreezingProtocolAutomation) emergencyUnfreeze(asset common.Asset) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    fmt.Printf("Emergency unfreezing asset %s.\n", asset.ID)
    time.Sleep(EmergencyUnfreezeDelay)

    unfreezeSuccess := automation.consensusSystem.UnfreezeAsset(asset)
    if unfreezeSuccess {
        fmt.Printf("Emergency unfreeze successful for asset %s.\n", asset.ID)
        automation.logAssetUnfreezeEvent(asset)
    } else {
        fmt.Printf("Emergency unfreeze failed for asset %s.\n", asset.ID)
    }
}

// manualIntervention allows administrators to manually freeze or unfreeze an asset
func (automation *AssetFreezingUnfreezingProtocolAutomation) manualIntervention(asset common.Asset, action string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    if action == "freeze" {
        fmt.Printf("Manually freezing asset %s.\n", asset.ID)
        automation.freezeAsset(asset)
    } else if action == "unfreeze" {
        fmt.Printf("Manually unfreezing asset %s.\n", asset.ID)
        automation.unfreezeAsset(asset)
    } else {
        fmt.Println("Invalid action for manual intervention.")
    }
}

// finalizeFreezeUnfreezeCycle finalizes the asset freeze/unfreeze check cycle and logs the result in the ledger
func (automation *AssetFreezingUnfreezingProtocolAutomation) finalizeFreezeUnfreezeCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeFreezeUnfreezeCycle()
    if success {
        fmt.Println("Asset freeze/unfreeze check cycle finalized successfully.")
        automation.logFreezeUnfreezeCycleFinalization()
    } else {
        fmt.Println("Error finalizing asset freeze/unfreeze check cycle.")
    }
}

// logAssetFreezeEvent logs the asset freeze event into the ledger
func (automation *AssetFreezingUnfreezingProtocolAutomation) logAssetFreezeEvent(asset common.Asset) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("asset-freeze-%s", asset.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Asset Freeze",
        Status:    "Completed",
        Details:   fmt.Sprintf("Asset %s frozen successfully.", asset.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with asset freeze event for asset %s.\n", asset.ID)
}

// logAssetUnfreezeEvent logs the asset unfreeze event into the ledger
func (automation *AssetFreezingUnfreezingProtocolAutomation) logAssetUnfreezeEvent(asset common.Asset) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("asset-unfreeze-%s", asset.ID),
        Timestamp: time.Now().Unix(),
        Type:      "Asset Unfreeze",
        Status:    "Completed",
        Details:   fmt.Sprintf("Asset %s unfrozen successfully.", asset.ID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with asset unfreeze event for asset %s.\n", asset.ID)
}

// logFreezeUnfreezeCycleFinalization logs the finalization of an asset freeze/unfreeze cycle into the ledger
func (automation *AssetFreezingUnfreezingProtocolAutomation) logFreezeUnfreezeCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("asset-freeze-unfreeze-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Asset Freeze/Unfreeze Cycle Finalization",
        Status:    "Finalized",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with asset freeze/unfreeze cycle finalization.")
}

// encryptAssetData encrypts the asset data before freezing
func (automation *AssetFreezingUnfreezingProtocolAutomation) encryptAssetData(asset common.Asset) common.Asset {
    encryptedData, err := encryption.EncryptData(asset.Data)
    if err != nil {
        fmt.Println("Error encrypting asset data:", err)
        return asset
    }

    asset.EncryptedData = encryptedData
    fmt.Println("Asset data successfully encrypted.")
    return asset
}
