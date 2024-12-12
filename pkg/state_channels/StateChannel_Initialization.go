// StateChannel_Initialization.go

package state_channels

import (
    "fmt"
    "errors"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// StateChannelInit initializes a new state channel.
func StateChannelInit(channelID string, config common.ChannelConfig, ledgerInstance *ledger.Ledger) error {
    encryptedConfig := encryption.EncryptChannelConfig(config)
    if err := ledgerInstance.InitializeChannel(channelID, encryptedConfig); err != nil {
        return fmt.Errorf("failed to initialize channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s initialized\n", channelID)
    return nil
}

// StateChannelOpen opens an initialized state channel.
func StateChannelOpen(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.OpenChannel(channelID); err != nil {
        return fmt.Errorf("failed to open channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s opened\n", channelID)
    return nil
}

// StateChannelClose closes an active state channel.
func StateChannelClose(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CloseChannel(channelID); err != nil {
        return fmt.Errorf("failed to close channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s closed\n", channelID)
    return nil
}

// StateChannelDeposit deposits funds into the state channel.
func StateChannelDeposit(channelID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.DepositToChannel(channelID, amount); err != nil {
        return fmt.Errorf("failed to deposit funds to channel %s: %v", channelID, err)
    }
    fmt.Printf("Deposited %d to channel %s\n", amount, channelID)
    return nil
}

// StateChannelWithdraw withdraws funds from the state channel.
func StateChannelWithdraw(channelID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.WithdrawFromChannel(channelID, amount); err != nil {
        return fmt.Errorf("failed to withdraw funds from channel %s: %v", channelID, err)
    }
    fmt.Printf("Withdrew %d from channel %s\n", amount, channelID)
    return nil
}

// StateChannelTransfer transfers funds between participants within the state channel.
func StateChannelTransfer(channelID string, sender string, receiver string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.TransferWithinChannel(channelID, sender, receiver, amount); err != nil {
        return fmt.Errorf("failed to transfer %d from %s to %s in channel %s: %v", amount, sender, receiver, channelID, err)
    }
    fmt.Printf("Transferred %d from %s to %s in channel %s\n", amount, sender, receiver, channelID)
    return nil
}

// StateChannelLock locks a state channel to prevent further transactions.
func StateChannelLock(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockChannel(channelID); err != nil {
        return fmt.Errorf("failed to lock channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s locked\n", channelID)
    return nil
}

// StateChannelUnlock unlocks a previously locked state channel.
func StateChannelUnlock(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockChannel(channelID); err != nil {
        return fmt.Errorf("failed to unlock channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s unlocked\n", channelID)
    return nil
}

// StateChannelCommit commits all changes within the channel, making them final.
func StateChannelCommit(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CommitChannel(channelID); err != nil {
        return fmt.Errorf("failed to commit channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s committed\n", channelID)
    return nil
}

// StateChannelRevert reverts the state channel to a previous committed state.
func StateChannelRevert(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertChannel(channelID); err != nil {
        return fmt.Errorf("failed to revert channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s reverted\n", channelID)
    return nil
}

// StateChannelSnapshot creates a snapshot of the current channel state.
func StateChannelSnapshot(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CreateChannelSnapshot(channelID); err != nil {
        return fmt.Errorf("failed to create snapshot for channel %s: %v", channelID, err)
    }
    fmt.Printf("Snapshot created for channel %s\n", channelID)
    return nil
}

// StateChannelRestore restores a channel to a previous snapshot state.
func StateChannelRestore(channelID string, snapshotID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RestoreChannelSnapshot(channelID, snapshotID); err != nil {
        return fmt.Errorf("failed to restore snapshot for channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s restored to snapshot %s\n", channelID, snapshotID)
    return nil
}

// StateChannelMonitor monitors the channel's activity and health.
func StateChannelMonitor(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MonitorChannelActivity(channelID); err != nil {
        return fmt.Errorf("failed to monitor channel %s: %v", channelID, err)
    }
    fmt.Printf("Monitoring channel %s\n", channelID)
    return nil
}

// StateChannelAudit audits the channel for compliance and performance.
func StateChannelAudit(channelID string, auditParams common.AuditParams, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditChannel(channelID, auditParams); err != nil {
        return fmt.Errorf("failed to audit channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s audited\n", channelID)
    return nil
}

// StateChannelFinalize finalizes the state channel, closing it with a final commitment.
func StateChannelFinalize(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeChannel(channelID); err != nil {
        return fmt.Errorf("failed to finalize channel %s: %v", channelID, err)
    }
    fmt.Printf("Channel %s finalized\n", channelID)
    return nil
}
