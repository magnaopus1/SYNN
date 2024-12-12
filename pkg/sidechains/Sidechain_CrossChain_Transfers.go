// Sidechain_CrossChain_Transfers.go

package sidechains

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// SidechainTrackCrossChainTransfer tracks a cross-chain transfer for a sidechain.
func SidechainTrackCrossChainTransfer(transferID string, details common.TransferDetails, ledgerInstance *ledger.Ledger) error {
    encryptedDetails := encryption.EncryptTransferDetails(details)
    if err := ledgerInstance.RecordCrossChainTransfer(transferID, encryptedDetails); err != nil {
        return fmt.Errorf("failed to track cross-chain transfer %s: %v", transferID, err)
    }
    fmt.Printf("Cross-chain transfer %s tracked successfully.\n", transferID)
    return nil
}

// SidechainMonitorCrossChainTransfer monitors the ongoing cross-chain transfer.
func SidechainMonitorCrossChainTransfer(transferID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.MonitorTransferStatus(transferID)
    if err != nil {
        return "", fmt.Errorf("failed to monitor cross-chain transfer %s: %v", transferID, err)
    }
    fmt.Printf("Cross-chain transfer %s status: %s\n", transferID, status)
    return status, nil
}

// SidechainLogCrossChainTransfer logs details of the cross-chain transfer.
func SidechainLogCrossChainTransfer(transferID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LogTransfer(transferID); err != nil {
        return fmt.Errorf("failed to log cross-chain transfer %s: %v", transferID, err)
    }
    fmt.Printf("Cross-chain transfer %s logged successfully.\n", transferID)
    return nil
}

// SidechainRevertCrossChainTransfer reverts a cross-chain transfer.
func SidechainRevertCrossChainTransfer(transferID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertTransfer(transferID); err != nil {
        return fmt.Errorf("failed to revert cross-chain transfer %s: %v", transferID, err)
    }
    fmt.Printf("Cross-chain transfer %s reverted successfully.\n", transferID)
    return nil
}

// SidechainBridgeStatus checks the current bridge status for cross-chain transfers.
func SidechainBridgeStatus(bridgeID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.FetchBridgeStatus(bridgeID)
    if err != nil {
        return "", fmt.Errorf("failed to fetch bridge status for bridge %s: %v", bridgeID, err)
    }
    return status, nil
}

// SidechainBridgeAudit performs an audit of the bridge used for cross-chain transfers.
func SidechainBridgeAudit(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditBridge(bridgeID); err != nil {
        return fmt.Errorf("bridge audit failed for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge %s audited successfully.\n", bridgeID)
    return nil
}

// SidechainUpdateBridgeParameters updates the parameters of a bridge.
func SidechainUpdateBridgeParameters(bridgeID string, params common.BridgeParameters, ledgerInstance *ledger.Ledger) error {
    encryptedParams := encryption.EncryptBridgeParameters(params)
    if err := ledgerInstance.UpdateBridgeParams(bridgeID, encryptedParams); err != nil {
        return fmt.Errorf("failed to update parameters for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge %s parameters updated successfully.\n", bridgeID)
    return nil
}

// SidechainQueryBridgeStatus queries the current status of a bridge.
func SidechainQueryBridgeStatus(bridgeID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.QueryBridgeStatus(bridgeID)
    if err != nil {
        return "", fmt.Errorf("failed to query bridge status for bridge %s: %v", bridgeID, err)
    }
    return status, nil
}

// SidechainRevertBridgeStatus reverts the last status change on a bridge.
func SidechainRevertBridgeStatus(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertBridgeStatus(bridgeID); err != nil {
        return fmt.Errorf("failed to revert bridge status for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge %s status reverted successfully.\n", bridgeID)
    return nil
}

// SidechainInitializeBridge initializes a new bridge for cross-chain transfers.
func SidechainInitializeBridge(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitializeBridge(bridgeID); err != nil {
        return fmt.Errorf("failed to initialize bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge %s initialized successfully.\n", bridgeID)
    return nil
}

// SidechainFinalizeBridge finalizes the setup of a bridge.
func SidechainFinalizeBridge(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeBridge(bridgeID); err != nil {
        return fmt.Errorf("failed to finalize bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge %s finalized successfully.\n", bridgeID)
    return nil
}

// SidechainSetBridgeFee sets the fee for using the bridge.
func SidechainSetBridgeFee(bridgeID string, fee int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetBridgeFee(bridgeID, fee); err != nil {
        return fmt.Errorf("failed to set bridge fee for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge fee set for bridge %s to %d.\n", bridgeID, fee)
    return nil
}

// SidechainFetchBridgeFee fetches the current bridge fee.
func SidechainFetchBridgeFee(bridgeID string, ledgerInstance *ledger.Ledger) (int, error) {
    fee, err := ledgerInstance.GetBridgeFee(bridgeID)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch bridge fee for bridge %s: %v", bridgeID, err)
    }
    return fee, nil
}

// SidechainAuditBridgeFee audits the bridge fee transactions.
func SidechainAuditBridgeFee(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditBridgeFee(bridgeID); err != nil {
        return fmt.Errorf("bridge fee audit failed for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge fee audited for bridge %s.\n", bridgeID)
    return nil
}

// SidechainConfirmBridgeFee confirms the bridge fee set for transactions.
func SidechainConfirmBridgeFee(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmBridgeFee(bridgeID); err != nil {
        return fmt.Errorf("failed to confirm bridge fee for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge fee confirmed for bridge %s.\n", bridgeID)
    return nil
}

// SidechainRevertBridgeFee reverts to the previous bridge fee.
func SidechainRevertBridgeFee(bridgeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertBridgeFee(bridgeID); err != nil {
        return fmt.Errorf("failed to revert bridge fee for bridge %s: %v", bridgeID, err)
    }
    fmt.Printf("Bridge fee reverted for bridge %s.\n", bridgeID)
    return nil
}
