package interoperability

import (
	"fmt"
	"sync"
	"time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CrossChainTransactionManager manages cross-chain transactions and related operations
type CrossChainTransactionManager struct {
	consensusEngine *common.SynnergyConsensus
	ledgerInstance  *ledger.Ledger
	transactionMutex *sync.RWMutex
}

// NewCrossChainTransactionManager initializes the CrossChainTransactionManager
func NewCrossChainTransactionManager(consensusEngine *common.SynnergyConsensus, ledgerInstance *ledger.Ledger, transactionMutex *sync.RWMutex) *CrossChainTransactionManager {
	return &CrossChainTransactionManager{
		consensusEngine: consensusEngine,
		ledgerInstance:  ledgerInstance,
		transactionMutex: transactionMutex,
	}
}

// crossChainAssetTransfer transfers an asset across chains, logging the details securely in the ledger
func (manager *CrossChainTransactionManager) crossChainAssetTransfer(assetID, sourceChainID, targetChainID string, amount float64) error {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    eventDetails := fmt.Sprintf("Transferring asset %s from %s to %s with amount %.2f", assetID, sourceChainID, targetChainID, amount)
    encryptedDetails := manager.encryptData(eventDetails)

    transfer := CrossChainAssetTransfer{
        TransferID:    fmt.Sprintf("asset-transfer-%s-%d", assetID, time.Now().Unix()),
        AssetID:       assetID,
        SourceChainID: sourceChainID,
        TargetChainID: targetChainID,
        Amount:        amount,
        Status:        "Completed",
        Timestamp:     time.Now(),
    }

    if err := manager.ledgerInstance.logCrossChainAssetTransfer(transfer); err != nil {
        return fmt.Errorf("failed to log cross-chain asset transfer for asset %s: %v", assetID, err)
    }

    return nil
}

// verifyCrossChainTransaction verifies the validity and authenticity of a cross-chain transaction
func (manager *CrossChainTransactionManager) verifyCrossChainTransaction(transactionID string) (bool, error) {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    isValid, err := manager.consensusEngine.ValidateTransaction(transactionID)
    if err != nil || !isValid {
        manager.logTransactionEvent(transactionID, "Verification Failed", "Cross-chain transaction verification failed")
        return false, fmt.Errorf("cross-chain transaction verification failed for transaction %s", transactionID)
    }

    manager.logTransactionEvent(transactionID, "Verified", "Cross-chain transaction verified successfully")
    return true, nil
}

// initiateCrossChainEscrow places assets in escrow for a cross-chain transaction
func (manager *CrossChainTransactionManager) initiateCrossChainEscrow(escrowID, assetID, sourceChainID, targetChainID string, amount float64) error {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    eventDetails := fmt.Sprintf("Initiating escrow %s for asset %s from %s to %s with amount %.2f", escrowID, assetID, sourceChainID, targetChainID, amount)
    encryptedDetails := manager.encryptData(eventDetails)

    escrow := CrossChainEscrow{
        EscrowID:      escrowID,
        AssetID:       assetID,
        SourceChainID: sourceChainID,
        TargetChainID: targetChainID,
        Amount:        amount,
        Status:        "Initiated",
        Timestamp:     time.Now(),
    }

    if err := manager.ledgerInstance.initiateCrossChainEscrow(escrow); err != nil {
        return fmt.Errorf("failed to initiate cross-chain escrow %s: %v", escrowID, err)
    }

    return nil
}

// releaseCrossChainEscrow releases assets from escrow after transaction confirmation
func (manager *CrossChainTransactionManager) releaseCrossChainEscrow(escrowID string) error {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    if err := manager.ledgerInstance.releaseCrossChainEscrow(escrowID); err != nil {
        return fmt.Errorf("failed to release cross-chain escrow %s: %v", escrowID, err)
    }

    eventDetails := fmt.Sprintf("Releasing escrow %s", escrowID)
    manager.logTransactionEvent(escrowID, "Escrow Released", eventDetails)

    return nil
}

// crossChainAssetSwap facilitates a cross-chain asset swap between two parties on different chains
func (manager *CrossChainTransactionManager) crossChainAssetSwap(swapID, assetID1, chainID1, assetID2, chainID2 string, amount1, amount2 float64) error {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    eventDetails := fmt.Sprintf("Asset swap %s: %s on %s for %s on %s, amounts: %.2f and %.2f", swapID, assetID1, chainID1, assetID2, chainID2, amount1, amount2)
    encryptedDetails := manager.encryptData(eventDetails)

    swap := CrossChainAssetSwap{
        SwapID:    swapID,
        AssetID1:  assetID1,
        ChainID1:  chainID1,
        AssetID2:  assetID2,
        ChainID2:  chainID2,
        Amount1:   amount1,
        Amount2:   amount2,
        Status:    "Swapped",
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.logCrossChainAssetSwap(swap); err != nil {
        return fmt.Errorf("failed to log cross-chain asset swap %s: %v", swapID, err)
    }

    return nil
}


// rollbackCrossChainAction rolls back a cross-chain action in case of failure or fraud
func (manager *CrossChainTransactionManager) rollbackCrossChainAction(actionID string) error {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    eventDetails := fmt.Sprintf("Rolling back action %s", actionID)
    encryptedDetails := manager.encryptData(eventDetails)

    rollback := CrossChainActionRollback{
        ActionID:  actionID,
        Details:   encryptedDetails,
        Status:    "Rolled Back",
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.logCrossChainRollback(rollback); err != nil {
        return fmt.Errorf("failed to rollback cross-chain action %s: %v", actionID, err)
    }

    return nil
}

// checkCrossChainBalance checks the cross-chain balance of a specified asset
func (manager *CrossChainTransactionManager) checkCrossChainBalance(assetID, chainID string) (float64, error) {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    balance, err := manager.ledgerInstance.getAssetBalance(assetID, chainID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve balance for asset %s on chain %s: %v", assetID, chainID, err)
    }
    return balance, nil
}

// validateCrossChainContract validates a smart contract involved in a cross-chain transaction
func (manager *CrossChainTransactionManager) validateCrossChainContract(contractID string) (bool, error) {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    contract, err := manager.ledgerInstance.validateCrossChainContract(contractID)
    if err != nil || !contract.IsValid {
        manager.logTransactionEvent(contractID, "Contract Validation Failed", "Cross-chain contract validation failed")
        return false, fmt.Errorf("contract validation failed for contract %s", contractID)
    }

    manager.logTransactionEvent(contractID, "Contract Validated", "Cross-chain contract validated successfully")
    return true, nil
}

// checkInterchainAgreements checks any agreements or arrangements required for a cross-chain transaction
func (manager *CrossChainTransactionManager) checkInterchainAgreements(agreementID string) (bool, error) {
    manager.transactionMutex.Lock()
    defer manager.transactionMutex.Unlock()

    agreement, err := manager.ledgerInstance.getInterchainAgreement(agreementID)
    if err != nil || !agreement.IsValid {
        manager.logTransactionEvent(agreementID, "Agreement Verification Failed", "Interchain agreement verification failed")
        return false, fmt.Errorf("interchain agreement verification failed for agreement %s", agreementID)
    }

    manager.logTransactionEvent(agreementID, "Agreement Verified", "Interchain agreement verified successfully")
    return true, nil
}

// logTransactionEvent logs events related to cross-chain transactions
func (manager *CrossChainTransactionManager) logTransactionEvent(entityID, eventType, details string) {
    eventDetails := fmt.Sprintf("Event: %s, Entity ID: %s, Details: %s", eventType, entityID, details)
    encryptedDetails := manager.encryptData(eventDetails)

    transactionEvent := CrossChainActionRollback{
        ActionID:  entityID,
        Details:   encryptedDetails,
        Status:    "Completed",
        Timestamp: time.Now(),
    }

    if err := manager.ledgerInstance.logCrossChainRollback(transactionEvent); err != nil {
        fmt.Printf("Failed to log transaction event for entity ID %s: %v\n", entityID, err)
    }
}

// encryptData encrypts data for secure logging and storage
func (manager *CrossChainTransactionManager) encryptData(data string) string {
    encryptedData, err := encryption.EncryptData([]byte(data))
    if err != nil {
        fmt.Println("Error encrypting data:", err)
        return data
    }
    return string(encryptedData)
}

// decryptData decrypts stored encrypted data
func (manager *CrossChainTransactionManager) decryptData(encryptedData string) string {
    decryptedData, err := encryption.DecryptData([]byte(encryptedData))
    if err != nil {
        fmt.Println("Error decrypting data:", err)
        return encryptedData
    }
    return string(decryptedData)
}
