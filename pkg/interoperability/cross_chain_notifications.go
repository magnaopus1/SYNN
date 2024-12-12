package interoperability

import (
	"fmt"
	"time"
	"sync"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CrossChainNotificationManager manages cross-chain notifications and asset tracking
type CrossChainNotificationManager struct {
	consensusEngine   *common.SynnergyConsensus
	ledgerInstance    *ledger.Ledger
	notificationMutex *sync.RWMutex
}

// NewCrossChainNotificationManager initializes the CrossChainNotificationManager
func NewCrossChainNotificationManager(consensusEngine *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, notificationMutex *sync.RWMutex) *CrossChainNotificationManager {
	return &CrossChainNotificationManager{
		consensusEngine:   consensusEngine,
		ledgerInstance:    ledgerInstance,
		notificationMutex: notificationMutex,
	}
}

// notifyAssetDeparture logs the departure of an asset from the Synnergy Network to an external chain
func (manager *CrossChainNotificationManager) notifyAssetDeparture(assetID, destinationChainID string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    eventDetails := fmt.Sprintf("Asset %s departing to chain %s", assetID, destinationChainID)
    encryptedDetails := manager.encryptData(eventDetails)

    assetLog := CrossChainAssetLog{
        AssetID:         assetID,
        ChainID:         destinationChainID,
        TransactionType: "Departure",
        Details:         encryptedDetails,
        Timestamp:       time.Now(),
        Status:          "Completed",
    }

    if err := manager.ledgerInstance.logCrossChainAsset(assetLog); err != nil {
        return fmt.Errorf("failed to log asset departure for asset %s: %v", assetID, err)
    }

    return nil
}

// confirmAssetArrival confirms the arrival of an asset on the Synnergy Network from an external chain
func (manager *CrossChainNotificationManager) confirmAssetArrival(assetID, originChainID string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    eventDetails := fmt.Sprintf("Asset %s arrived from chain %s", assetID, originChainID)
    encryptedDetails := manager.encryptData(eventDetails)

    assetLog := CrossChainAssetLog{
        AssetID:         assetID,
        ChainID:         originChainID,
        TransactionType: "Arrival",
        Details:         encryptedDetails,
        Timestamp:       time.Now(),
        Status:          "Completed",
    }

    if err := manager.ledgerInstance.logCrossChainAsset(assetLog); err != nil {
        return fmt.Errorf("failed to confirm asset arrival for asset %s: %v", assetID, err)
    }

    return nil
}

// freezeCrossChainAsset freezes an asset involved in a cross-chain transaction
func (manager *CrossChainNotificationManager) freezeCrossChainAsset(assetID string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    if err := manager.ledgerInstance.freezeAsset(assetID); err != nil {
        return fmt.Errorf("failed to freeze asset %s: %v", assetID, err)
    }

    eventDetails := fmt.Sprintf("Asset %s frozen for cross-chain transaction", assetID)
    encryptedDetails := manager.encryptData(eventDetails)

    assetLog := CrossChainAssetLog{
        AssetID:         assetID,
        TransactionType: "Freeze",
        Details:         encryptedDetails,
        Timestamp:       time.Now(),
        Status:          "Frozen",
    }

    if err := manager.ledgerInstance.logCrossChainAsset(assetLog); err != nil {
        return fmt.Errorf("failed to log freeze event for asset %s: %v", assetID, err)
    }

    return nil
}

// unfreezeCrossChainAsset unfreezes an asset involved in a cross-chain transaction
func (manager *CrossChainNotificationManager) unfreezeCrossChainAsset(assetID string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    if err := manager.ledgerInstance.unfreezeAsset(assetID); err != nil {
        return fmt.Errorf("failed to unfreeze asset %s: %v", assetID, err)
    }

    eventDetails := fmt.Sprintf("Asset %s unfrozen after cross-chain transaction", assetID)
    encryptedDetails := manager.encryptData(eventDetails)

    assetLog := CrossChainAssetLog{
        AssetID:         assetID,
        TransactionType: "Unfreeze",
        Details:         encryptedDetails,
        Timestamp:       time.Now(),
        Status:          "Unfrozen",
    }

    if err := manager.ledgerInstance.logCrossChainAsset(assetLog); err != nil {
        return fmt.Errorf("failed to log unfreeze event for asset %s: %v", assetID, err)
    }

    return nil
}

// trackAssetHistory logs historical data related to cross-chain asset transactions
func (manager *CrossChainNotificationManager) trackAssetHistory(assetID, transactionDetails string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    encryptedDetails := manager.encryptData(transactionDetails)

    assetHistory := AssetHistory{
        AssetID:           assetID,
        TransactionID:     fmt.Sprintf("txn-%s-%d", assetID, time.Now().Unix()),
        TransactionDetails: encryptedDetails,
        Timestamp:         time.Now(),
    }

    if err := manager.ledgerInstance.trackAssetHistory(assetHistory); err != nil {
        return fmt.Errorf("failed to track asset history for asset %s: %v", assetID, err)
    }

    return nil
}


// verifyAssetHistory verifies the historical records of an asset to ensure consistency and integrity
func (manager *CrossChainNotificationManager) verifyAssetHistory(assetID string) (bool, error) {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    historyRecords, err := manager.ledgerInstance.getAssetHistory(assetID)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve asset history for asset %s: %v", assetID, err)
    }

    isVerified := manager.consensusEngine.VerifyHistoryIntegrity(historyRecords)
    if !isVerified {
        manager.logNotificationEvent(assetID, "History Verification Failed", "Asset history verification failed")
        return false, fmt.Errorf("asset history verification failed for asset %s", assetID)
    }

    manager.logNotificationEvent(assetID, "History Verified", "Asset history verified successfully")
    return true, nil
}

// notifyCrossChainEvent notifies the network of significant cross-chain events related to assets
func (manager *CrossChainNotificationManager) notifyCrossChainEvent(eventType, assetID, eventDetails string) error {
    manager.notificationMutex.Lock()
    defer manager.notificationMutex.Unlock()

    encryptedDetails := manager.encryptData(eventDetails)

    crossChainEvent := CrossChainEvent{
        EventID:    fmt.Sprintf("event-%s-%d", assetID, time.Now().Unix()),
        AssetID:    assetID,
        EventType:  eventType,
        Details:    encryptedDetails,
        Timestamp:  time.Now(),
    }

    if err := manager.ledgerInstance.addCrossChainEvent(crossChainEvent); err != nil {
        return fmt.Errorf("failed to notify cross-chain event for asset %s: %v", assetID, err)
    }

    return nil
}

// logNotificationEvent logs events related to cross-chain asset notifications
func (manager *CrossChainNotificationManager) logNotificationEvent(assetID, eventType, details string) {
    encryptedDetails := manager.encryptData(details)

    notificationEvent := CrossChainEvent{
        EventID:    fmt.Sprintf("notification-%s-%d", assetID, time.Now().Unix()),
        AssetID:    assetID,
        EventType:  eventType,
        Details:    encryptedDetails,
        Timestamp:  time.Now(),
    }

    if err := manager.ledgerInstance.logNotificationEvent(assetID, notificationEvent); err != nil {
        fmt.Printf("Failed to log notification event for asset ID %s: %v\n", assetID, err)
    }
}

// encryptData encrypts data for secure logging and storage
func (manager *CrossChainNotificationManager) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}

// decryptData decrypts stored encrypted data
func (manager *CrossChainNotificationManager) decryptData(encryptedData string) string {
    decryptedData, err := encryption.DecryptData([]byte(encryptedData))
    if err != nil {
        fmt.Println("Error decrypting data:", err)
        return encryptedData
    }
    return string(decryptedData)
}
