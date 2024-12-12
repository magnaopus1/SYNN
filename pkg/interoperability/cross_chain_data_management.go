package interoperability

import (
	"fmt"
	"time"
	"sync"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CrossChainDataManager handles cross-chain data management functions
type CrossChainDataManager struct {
	consensusEngine  *common.SynnergyConsensus
	ledgerInstance   *ledger.Ledger
	dataMutex        *sync.RWMutex
}

// NewCrossChainDataManager initializes the CrossChainDataManager
func NewCrossChainDataManager(consensusEngine *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, dataMutex *sync.RWMutex) *CrossChainDataManager {
	return &CrossChainDataManager{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		dataMutex:       dataMutex,
	}
}

func (manager *CrossChainDataManager) queryCrossChainData(chainID, query string) (string, error) {
    manager.dataMutex.Lock()
    defer manager.dataMutex.Unlock()

    data, err := manager.ledgerInstance.queryExternalData(chainID, query)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve data from chain %s: %v", chainID, err)
    }

    if err := manager.consensusEngine.validateExternalData(data); err != nil {
        return "", fmt.Errorf("data validation failed for chain %s: %v", chainID, err)
    }

    return manager.encryptData(data), nil
}

func (manager *CrossChainDataManager) validateExternalDataFeed(feedID string) error {
    feedData, err := manager.ledgerInstance.getDataFeed(feedID)
    if err != nil {
        return fmt.Errorf("failed to retrieve data feed %s: %v", feedID, err)
    }

    if err := manager.consensusEngine.validateDataFeed(feedData); err != nil {
        manager.logDataFeedEvent(feedID, "Validation Failed", fmt.Sprintf("Data feed validation failed: %v", err))
        return fmt.Errorf("data feed validation failed: %v", err)
    }

    manager.logDataFeedEvent(feedID, "Validation Success", "Data feed validated successfully")
    return nil
}

func (manager *CrossChainDataManager) reviewExternalDataAccuracy(dataID string) error {
    externalData, err := manager.ledgerInstance.getExternalData(dataID)
    if err != nil {
        return fmt.Errorf("failed to retrieve external data %s: %v", dataID, err)
    }

    if !manager.consensusEngine.verifyDataAccuracy(externalData) {
        return fmt.Errorf("data accuracy check failed for data ID %s", dataID)
    }

    manager.logDataEvent(dataID, "Accuracy Check Passed", "External data is accurate")
    return nil
}

func (manager *CrossChainDataManager) verifyCrossChainSignature(signature, data string) (bool, error) {
    isVerified, err := manager.consensusEngine.verifySignature(signature, data)
    if err != nil || !isVerified {
        manager.logDataEvent(data, "Signature Verification Failed", "Failed to verify signature")
        return false, fmt.Errorf("signature verification failed for data")
    }

    manager.logDataEvent(data, "Signature Verified", "Signature verification successful")
    return true, nil
}

func (manager *CrossChainDataManager) validateCrossChainLicense(licenseID string) error {
    license, err := manager.ledgerInstance.getLicense(licenseID)
    if err != nil {
        return fmt.Errorf("failed to retrieve license %s: %v", licenseID, err)
    }

    if !manager.consensusEngine.validateLicense(license) {
        return fmt.Errorf("license validation failed for license ID %s", licenseID)
    }

    manager.logDataEvent(licenseID, "License Validated", "Cross-chain license is valid")
    return nil
}

func (manager *CrossChainDataManager) revokeCrossChainLicense(licenseID string) error {
    if err := manager.ledgerInstance.RevokeLicense(licenseID); err != nil {
        return fmt.Errorf("failed to revoke license %s: %v", licenseID, err)
    }

    manager.LogDataEvent(licenseID, "License Revoked", "Cross-chain license revoked successfully")
    return nil
}

func (manager *CrossChainDataManager) checkCrossChainStatus(chainID string) (string, error) {
    status, err := manager.ledgerInstance.GetChainStatus(chainID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve status for chain %s: %v", chainID, err)
    }

    manager.LogDataEvent(chainID, "Status Checked", fmt.Sprintf("Chain status: %s", status.Status))
    return manager.encryptData(status.Status), nil
}

func (manager *CrossChainDataManager) retrieveEscrowFunds(transactionID string) error {
    if err := manager.ledgerInstance.ReleaseEscrow(transactionID); err != nil {
        return fmt.Errorf("failed to release escrow funds for transaction %s: %v", transactionID, err)
    }

    manager.LogEscrowEvent(transactionID, "Funds Retrieved", "Escrow funds retrieved successfully")
    return nil
}


// ReturnEscrowFunds returns escrowed funds to the originator in the event of transaction failure
func (manager *CrossChainDataManager) returnEscrowFunds(transactionID string) error {
    if err := manager.ledgerInstance.returnEscrowFunds(transactionID); err != nil {
        return fmt.Errorf("failed to return escrow funds for transaction %s: %v", transactionID, err)
    }

    manager.logEscrowEvent(transactionID, "Funds Returned", "Escrow funds returned to originator")
    return nil
}

// LogDataFeedEvent logs events related to data feeds for traceability
func (manager *CrossChainDataManager) logDataFeedEvent(feedID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Feed ID: %s, Details: %s", eventType, feedID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    event := DataFeedEvent{
        FeedID:    feedID,
        EventType: eventType,
        Details:   encryptedDetails,
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.addDataFeedEvent(event); err != nil {
        fmt.Printf("Failed to log data feed event for feed ID %s: %v\n", feedID, err)
    }
}

// LogDataEvent logs events related to cross-chain data management
func (manager *CrossChainDataManager) logDataEvent(dataID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Data ID: %s, Details: %s", eventType, dataID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    event := DataEvent{
        DataID:    dataID,
        EventType: eventType,
        Details:   encryptedDetails,
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.addDataEvent(event); err != nil {
        fmt.Printf("Failed to log data event for data ID %s: %v\n", dataID, err)
    }
}

// LogEscrowEvent logs events related to escrow management for cross-chain transactions
func (manager *CrossChainDataManager) logEscrowEvent(transactionID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Transaction ID: %s, Details: %s", eventType, transactionID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    event := EscrowEvent{
        TransactionID: transactionID,
        EventType:     eventType,
        Details:       encryptedDetails,
        Timestamp:     time.Now(),
    }

    if err := manager.ledgerInstance.addEscrowEvent(event); err != nil {
        fmt.Printf("Failed to log escrow event for transaction ID %s: %v\n", transactionID, err)
    }
}

// encryptData encrypts data for secure logging and storage
func (manager *CrossChainDataManager) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}

// decryptData decrypts stored encrypted data
func (manager *CrossChainDataManager) decryptData(encryptedData string) string {
    decryptedData, err := encryption.DecryptData([]byte(encryptedData))
    if err != nil {
        fmt.Println("Error decrypting data:", err)
        return encryptedData
    }
    return string(decryptedData)
}
