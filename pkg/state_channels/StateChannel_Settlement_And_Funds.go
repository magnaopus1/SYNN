// StateChannel_Settlement_And_Funds.go

package state_channels

import (
    "fmt"
    "errors"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// StateChannelSettle initiates settlement within the state channel.
func StateChannelSettle(channelID string, amount int, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptInt(amount)
    if err := ledgerInstance.SettleChannel(channelID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to settle channel %s: %v", channelID, err)
    }
    fmt.Printf("Settlement initiated for channel %s\n", channelID)
    return nil
}

// StateChannelConfirmSettlement confirms the settlement within the state channel.
func StateChannelConfirmSettlement(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmSettlement(channelID); err != nil {
        return fmt.Errorf("failed to confirm settlement for channel %s: %v", channelID, err)
    }
    fmt.Printf("Settlement confirmed for channel %s\n", channelID)
    return nil
}

// StateChannelRevertSettlement reverts a previous settlement operation.
func StateChannelRevertSettlement(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertSettlement(channelID); err != nil {
        return fmt.Errorf("failed to revert settlement for channel %s: %v", channelID, err)
    }
    fmt.Printf("Settlement reverted for channel %s\n", channelID)
    return nil
}

// StateChannelEscrowFunds places funds in escrow within the state channel.
func StateChannelEscrowFunds(channelID string, amount int, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptInt(amount)
    if err := ledgerInstance.EscrowFunds(channelID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to escrow funds in channel %s: %v", channelID, err)
    }
    fmt.Printf("Funds escrowed in channel %s\n", channelID)
    return nil
}

// StateChannelReleaseFunds releases escrowed funds within the state channel.
func StateChannelReleaseFunds(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseFunds(channelID); err != nil {
        return fmt.Errorf("failed to release escrowed funds in channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrowed funds released in channel %s\n", channelID)
    return nil
}

// StateChannelFreeze freezes the state channel, halting all transactions.
func StateChannelFreeze(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FreezeChannel(channelID); err != nil {
        return fmt.Errorf("failed to freeze channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s frozen\n", channelID)
    return nil
}

// StateChannelUnfreeze unfreezes a previously frozen state channel.
func StateChannelUnfreeze(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnfreezeChannel(channelID); err != nil {
        return fmt.Errorf("failed to unfreeze channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s unfrozen\n", channelID)
    return nil
}

// StateChannelCommitToMainChain commits the state channel’s data to the main chain.
func StateChannelCommitToMainChain(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CommitToMainChain(channelID); err != nil {
        return fmt.Errorf("failed to commit channel %s to main chain: %v", channelID, err)
    }
    fmt.Printf("Channel %s committed to main chain\n", channelID)
    return nil
}

// StateChannelSyncToMainChain syncs the current state channel data to the main chain.
func StateChannelSyncToMainChain(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SyncChannelToMainChain(channelID); err != nil {
        return fmt.Errorf("failed to sync channel %s to main chain: %v", channelID, err)
    }
    fmt.Printf("Channel %s synced to main chain\n", channelID)
    return nil
}

// StateChannelSyncFromMainChain syncs the state channel data from the main chain.
func StateChannelSyncFromMainChain(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SyncChannelFromMainChain(channelID); err != nil {
        return fmt.Errorf("failed to sync channel %s from main chain: %v", channelID, err)
    }
    fmt.Printf("Channel %s synced from main chain\n", channelID)
    return nil
}

// StateChannelInitiateExit initiates an exit from the state channel.
func StateChannelInitiateExit(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitiateExit(channelID); err != nil {
        return fmt.Errorf("failed to initiate exit for channel %s: %v", channelID, err)
    }
    fmt.Printf("Exit initiated for channel %s\n", channelID)
    return nil
}

// StateChannelConfirmExit confirms the exit from the state channel.
func StateChannelConfirmExit(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmExit(channelID); err != nil {
        return fmt.Errorf("failed to confirm exit for channel %s: %v", channelID, err)
    }
    fmt.Printf("Exit confirmed for channel %s\n", channelID)
    return nil
}

// StateChannelRevertExit reverts an exit operation.
func StateChannelRevertExit(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertExit(channelID); err != nil {
        return fmt.Errorf("failed to revert exit for channel %s: %v", channelID, err)
    }
    fmt.Printf("Exit reverted for channel %s\n", channelID)
    return nil
}

// StateChannelUpdateExitStatus updates the status of an exit operation.
func StateChannelUpdateExitStatus(channelID string, status common.ExitStatus, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateExitStatus(channelID, status); err != nil {
        return fmt.Errorf("failed to update exit status for channel %s: %v", channelID, err)
    }
    fmt.Printf("Exit status updated for channel %s\n", channelID)
    return nil
}

// StateChannelMonitorExitStatus monitors the current status of the exit.
func StateChannelMonitorExitStatus(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MonitorExitStatus(channelID); err != nil {
        return fmt.Errorf("failed to monitor exit status for channel %s: %v", channelID, err)
    }
    fmt.Printf("Monitoring exit status for channel %s\n", channelID)
    return nil
}

// StateChannelSetExpiration sets an expiration time for the state channel.
func StateChannelSetExpiration(channelID string, expirationTime time.Time, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetChannelExpiration(channelID, expirationTime); err != nil {
        return fmt.Errorf("failed to set expiration for channel %s: %v", channelID, err)
    }
    fmt.Printf("Expiration set for channel %s\n", channelID)
    return nil
}

// StateChannelExtendExpiration extends the expiration time of the state channel.
func StateChannelExtendExpiration(channelID string, additionalTime time.Duration, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ExtendChannelExpiration(channelID, additionalTime); err != nil {
        return fmt.Errorf("failed to extend expiration for channel %s: %v", channelID, err)
    }
    fmt.Printf("Expiration extended for channel %s\n", channelID)
    return nil
}
