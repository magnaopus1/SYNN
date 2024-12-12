// StateChannel_State_Management.go

package state_channels

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// StateChannelCheckExpiration checks if the state channel has expired.
func StateChannelCheckExpiration(channelID string, ledgerInstance *ledger.Ledger) (bool, error) {
    expirationTime, err := ledgerInstance.GetChannelExpiration(channelID)
    if err != nil {
        return false, fmt.Errorf("failed to fetch expiration time for channel %s: %v", channelID, err)
    }
    expired := time.Now().After(expirationTime)
    fmt.Printf("Expiration check for channel %s: expired=%t\n", channelID, expired)
    return expired, nil
}

// StateChannelCloseOnExpiration closes the state channel if it has expired.
func StateChannelCloseOnExpiration(channelID string, ledgerInstance *ledger.Ledger) error {
    expired, err := StateChannelCheckExpiration(channelID, ledgerInstance)
    if err != nil {
        return err
    }
    if expired {
        err := ledgerInstance.UpdateChannelState(channelID, "closed")
        if err != nil {
            return fmt.Errorf("failed to close expired channel %s: %v", channelID, err)
        }
        fmt.Printf("Channel %s closed due to expiration\n", channelID)
    }
    return nil
}

// StateChannelLogStateChange logs state changes within the channel.
func StateChannelLogStateChange(channelID string, stateChangeDetails string, ledgerInstance *ledger.Ledger) error {
    timestamp := time.Now().Format(time.RFC3339)
    if err := ledgerInstance.LogStateChange(channelID, stateChangeDetails, timestamp); err != nil {
        return fmt.Errorf("failed to log state change for channel %s: %v", channelID, err)
    }
    fmt.Printf("State change logged for channel %s: %s\n", channelID, stateChangeDetails)
    return nil
}

// StateChannelFetchStateHistory retrieves the history of state changes.
func StateChannelFetchStateHistory(channelID string, ledgerInstance *ledger.Ledger) ([]string, error) {
    history, err := ledgerInstance.FetchStateHistory(channelID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch state history for channel %s: %v", channelID, err)
    }
    fmt.Printf("Fetched state history for channel %s\n", channelID)
    return history, nil
}

// StateChannelAuditStateHistory audits the state history of a channel.
func StateChannelAuditStateHistory(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditStateHistory(channelID); err != nil {
        return fmt.Errorf("failed to audit state history for channel %s: %v", channelID, err)
    }
    fmt.Printf("State history audited for channel %s\n", channelID)
    return nil
}

// StateChannelRevertStateHistory reverts the last recorded state history entry.
func StateChannelRevertStateHistory(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertStateHistory(channelID); err != nil {
        return fmt.Errorf("failed to revert state history for channel %s: %v", channelID, err)
    }
    fmt.Printf("State history reverted for channel %s\n", channelID)
    return nil
}

// StateChannelVerifyStateHistory verifies the integrity of state history entries.
func StateChannelVerifyStateHistory(channelID string, validationCriteria string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.VerifyStateHistory(channelID, validationCriteria); err != nil {
        return fmt.Errorf("failed to verify state history for channel %s: %v", channelID, err)
    }
    fmt.Printf("State history verified for channel %s\n", channelID)
    return nil
}

// StateChannelFinalizeStateHistory finalizes the state history log.
func StateChannelFinalizeStateHistory(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeStateHistory(channelID); err != nil {
        return fmt.Errorf("failed to finalize state history for channel %s: %v", channelID, err)
    }
    fmt.Printf("State history finalized for channel %s\n", channelID)
    return nil
}

// StateChannelValidateInclusionProof validates inclusion proof for transactions.
func StateChannelValidateInclusionProof(channelID string, proofData string, ledgerInstance *ledger.Ledger) error {
    encryptedProof := encryption.EncryptString(proofData)
    if err := ledgerInstance.ValidateInclusionProof(channelID, encryptedProof); err != nil {
        return fmt.Errorf("failed to validate inclusion proof for channel %s: %v", channelID, err)
    }
    fmt.Printf("Inclusion proof validated for channel %s\n", channelID)
    return nil
}

// StateChannelReconcileState reconciles the state of the channel with the main ledger.
func StateChannelReconcileState(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReconcileChannelState(channelID); err != nil {
        return fmt.Errorf("failed to reconcile state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State reconciled for channel %s\n", channelID)
    return nil
}

// StateChannelTrackInclusionProof tracks an inclusion proof within the channel.
func StateChannelTrackInclusionProof(channelID string, proofData string, ledgerInstance *ledger.Ledger) error {
    encryptedProof := encryption.EncryptString(proofData)
    if err := ledgerInstance.TrackInclusionProof(channelID, encryptedProof); err != nil {
        return fmt.Errorf("failed to track inclusion proof for channel %s: %v", channelID, err)
    }
    fmt.Printf("Inclusion proof tracked for channel %s\n", channelID)
    return nil
}

// StateChannelSubmitInclusionProof submits inclusion proof for the channel's transactions.
func StateChannelSubmitInclusionProof(channelID string, proofData string, ledgerInstance *ledger.Ledger) error {
    encryptedProof := encryption.EncryptString(proofData)
    if err := ledgerInstance.SubmitInclusionProof(channelID, encryptedProof); err != nil {
        return fmt.Errorf("failed to submit inclusion proof for channel %s: %v", channelID, err)
    }
    fmt.Printf("Inclusion proof submitted for channel %s\n", channelID)
    return nil
}

// StateChannelMonitorInclusionProof monitors inclusion proof for activity changes.
func StateChannelMonitorInclusionProof(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MonitorInclusionProof(channelID); err != nil {
        return fmt.Errorf("failed to monitor inclusion proof for channel %s: %v", channelID, err)
    }
    fmt.Printf("Inclusion proof monitored for channel %s\n", channelID)
    return nil
}

// StateChannelSetInclusionProof sets an inclusion proof parameter in the ledger.
func StateChannelSetInclusionProof(channelID string, proofParams string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetInclusionProof(channelID, proofParams); err != nil {
        return fmt.Errorf("failed to set inclusion proof for channel %s: %v", channelID, err)
    }
    fmt.Printf("Inclusion proof set for channel %s\n", channelID)
    return nil
}

// StateChannelLockState locks the current state of the channel.
func StateChannelLockState(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockChannelState(channelID); err != nil {
        return fmt.Errorf("failed to lock state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State locked for channel %s\n", channelID)
    return nil
}

// StateChannelUnlockState unlocks the previously locked state of the channel.
func StateChannelUnlockState(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockChannelState(channelID); err != nil {
        return fmt.Errorf("failed to unlock state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State unlocked for channel %s\n", channelID)
    return nil
}
