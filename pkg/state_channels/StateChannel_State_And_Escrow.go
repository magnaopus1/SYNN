// StateChannel_State_And_Escrow.go

package state_channels

import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// StateChannelFreezeState freezes the current state of a state channel.
func StateChannelFreezeState(channelID string, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.UpdateChannelState(channelID, "frozen")
    if err != nil {
        return fmt.Errorf("failed to freeze state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State frozen for channel %s\n", channelID)
    return nil
}

// StateChannelUnfreezeState unfreezes the state of a state channel.
func StateChannelUnfreezeState(channelID string, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.UpdateChannelState(channelID, "active")
    if err != nil {
        return fmt.Errorf("failed to unfreeze state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State unfrozen for channel %s\n", channelID)
    return nil
}

// StateChannelEscrowState places the state of a channel in escrow.
func StateChannelEscrowState(channelID string, amount int, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptInt(amount)
    err := ledgerInstance.EscrowFunds(channelID, encryptedAmount)
    if err != nil {
        return fmt.Errorf("failed to escrow state for channel %s: %v", channelID, err)
    }
    fmt.Printf("State escrowed for channel %s with amount: %d\n", channelID, amount)
    return nil
}

// StateChannelReleaseStateEscrow releases the escrowed state funds.
func StateChannelReleaseStateEscrow(channelID string, ledgerInstance *ledger.Ledger) error {
    err := ledgerInstance.ReleaseEscrow(channelID)
    if err != nil {
        return fmt.Errorf("failed to release escrow for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrow released for channel %s\n", channelID)
    return nil
}

// StateChannelLogEscrowEvent logs an event related to escrow activities.
func StateChannelLogEscrowEvent(channelID string, eventDetails string, ledgerInstance *ledger.Ledger) error {
    timestamp := time.Now().Format(time.RFC3339)
    if err := ledgerInstance.LogEscrowEvent(channelID, eventDetails, timestamp); err != nil {
        return fmt.Errorf("failed to log escrow event for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrow event logged for channel %s: %s\n", channelID, eventDetails)
    return nil
}

// StateChannelFetchEscrowLog retrieves the escrow log for auditing.
func StateChannelFetchEscrowLog(channelID string, ledgerInstance *ledger.Ledger) ([]string, error) {
    logs, err := ledgerInstance.FetchEscrowLog(channelID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch escrow log for channel %s: %v", channelID, err)
    }
    fmt.Printf("Fetched escrow log for channel %s\n", channelID)
    return logs, nil
}

// StateChannelAuditEscrowLog audits the escrow logs to ensure compliance.
func StateChannelAuditEscrowLog(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditEscrowLog(channelID); err != nil {
        return fmt.Errorf("failed to audit escrow log for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrow log audited for channel %s\n", channelID)
    return nil
}

// StateChannelValidateEscrowLog validates the escrow log entries.
func StateChannelValidateEscrowLog(channelID string, validationCriteria string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ValidateEscrowLog(channelID, validationCriteria); err != nil {
        return fmt.Errorf("failed to validate escrow log for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrow log validated for channel %s\n", channelID)
    return nil
}

// StateChannelRevertEscrowLog reverts the last escrow entry.
func StateChannelRevertEscrowLog(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertEscrowLog(channelID); err != nil {
        return fmt.Errorf("failed to revert escrow log for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrow log reverted for channel %s\n", channelID)
    return nil
}

// StateChannelInitiateDispute initiates a dispute in the state channel.
func StateChannelInitiateDispute(channelID string, disputeReason string, ledgerInstance *ledger.Ledger) error {
    encryptedReason := encryption.EncryptString(disputeReason)
    if err := ledgerInstance.RecordDispute(channelID, encryptedReason); err != nil {
        return fmt.Errorf("failed to initiate dispute for channel %s: %v", channelID, err)
    }
    fmt.Printf("Dispute initiated for channel %s with reason: %s\n", channelID, disputeReason)
    return nil
}

// StateChannelRespondToDispute provides a response to an ongoing dispute.
func StateChannelRespondToDispute(channelID string, response string, ledgerInstance *ledger.Ledger) error {
    encryptedResponse := encryption.EncryptString(response)
    if err := ledgerInstance.RecordDisputeResponse(channelID, encryptedResponse); err != nil {
        return fmt.Errorf("failed to respond to dispute for channel %s: %v", channelID, err)
    }
    fmt.Printf("Responded to dispute for channel %s\n", channelID)
    return nil
}

// StateChannelResolveDispute resolves a dispute in the state channel.
func StateChannelResolveDispute(channelID string, resolutionDetails string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResolveDispute(channelID, resolutionDetails); err != nil {
        return fmt.Errorf("failed to resolve dispute for channel %s: %v", channelID, err)
    }
    fmt.Printf("Dispute resolved for channel %s\n", channelID)
    return nil
}

// StateChannelEscrowDisputeTokens escrows tokens related to a dispute.
func StateChannelEscrowDisputeTokens(channelID string, tokenAmount int, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptInt(tokenAmount)
    if err := ledgerInstance.EscrowTokens(channelID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to escrow tokens for dispute in channel %s: %v", channelID, err)
    }
    fmt.Printf("Tokens escrowed for dispute in channel %s: %d\n", channelID, tokenAmount)
    return nil
}

// StateChannelReleaseDisputeTokens releases escrowed tokens for a resolved dispute.
func StateChannelReleaseDisputeTokens(channelID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseEscrowedTokens(channelID); err != nil {
        return fmt.Errorf("failed to release escrowed tokens for channel %s: %v", channelID, err)
    }
    fmt.Printf("Escrowed tokens released for channel %s\n", channelID)
    return nil
}

// StateChannelFetchDisputeHistory fetches the history of disputes for auditing.
func StateChannelFetchDisputeHistory(channelID string, ledgerInstance *ledger.Ledger) ([]string, error) {
    history, err := ledgerInstance.FetchDisputeHistory(channelID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch dispute history for channel %s: %v", channelID, err)
    }
    fmt.Printf("Fetched dispute history for channel %s\n", channelID)
    return history, nil
}
