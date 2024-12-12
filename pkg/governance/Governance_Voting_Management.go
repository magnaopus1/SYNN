package governance

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// GovernanceRevertProposal reverts a specific governance proposal to its previous state.
func GovernanceRevertProposal(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertProposal(proposalID); err != nil {
        return fmt.Errorf("failed to revert proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal %s reverted.\n", proposalID)
    return nil
}

// GovernanceEscrowVoteTokens places tokens in escrow for voting.
func GovernanceEscrowVoteTokens(voteID string, amount int64, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptAmount(amount)
    if err := ledgerInstance.EscrowTokens(voteID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to escrow tokens for vote %s: %v", voteID, err)
    }
    fmt.Printf("Tokens escrowed for vote %s.\n", voteID)
    return nil
}

// GovernanceReleaseVoteTokens releases escrowed tokens after voting.
func GovernanceReleaseVoteTokens(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseEscrowedTokens(voteID); err != nil {
        return fmt.Errorf("failed to release escrowed tokens for vote %s: %v", voteID, err)
    }
    fmt.Printf("Escrowed tokens released for vote %s.\n", voteID)
    return nil
}

// GovernanceFreezeProposal freezes a proposal from further actions.
func GovernanceFreezeProposal(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FreezeProposal(proposalID); err != nil {
        return fmt.Errorf("failed to freeze proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal %s has been frozen.\n", proposalID)
    return nil
}

// GovernanceUnfreezeProposal unfreezes a previously frozen proposal.
func GovernanceUnfreezeProposal(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnfreezeProposal(proposalID); err != nil {
        return fmt.Errorf("failed to unfreeze proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal %s has been unfrozen.\n", proposalID)
    return nil
}

// GovernanceSnapshotProposalState creates a snapshot of the proposal state.
func GovernanceSnapshotProposalState(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SnapshotProposalState(proposalID); err != nil {
        return fmt.Errorf("failed to snapshot proposal state for %s: %v", proposalID, err)
    }
    fmt.Printf("Snapshot of proposal state created for %s.\n", proposalID)
    return nil
}

// GovernanceRestoreProposalState restores a proposal state from a snapshot.
func GovernanceRestoreProposalState(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RestoreProposalState(proposalID); err != nil {
        return fmt.Errorf("failed to restore proposal state for %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal state restored for %s.\n", proposalID)
    return nil
}

// GovernanceLogVotingHistory logs the voting history for future audits.
func GovernanceLogVotingHistory(voteID string, history common.VotingHistory, ledgerInstance *ledger.Ledger) error {
    encryptedHistory := encryption.EncryptVotingHistory(history)
    if err := ledgerInstance.LogVotingHistory(voteID, encryptedHistory); err != nil {
        return fmt.Errorf("failed to log voting history for vote %s: %v", voteID, err)
    }
    fmt.Printf("Voting history logged for vote %s.\n", voteID)
    return nil
}

// GovernanceFetchVotingHistory retrieves the voting history of a vote.
func GovernanceFetchVotingHistory(voteID string, ledgerInstance *ledger.Ledger) (common.VotingHistory, error) {
    history, err := ledgerInstance.FetchVotingHistory(voteID)
    if err != nil {
        return common.VotingHistory{}, fmt.Errorf("failed to fetch voting history for vote %s: %v", voteID, err)
    }
    return history, nil
}

// GovernanceAuditVotingProcess audits the voting process for discrepancies.
func GovernanceAuditVotingProcess(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditVotingProcess(voteID); err != nil {
        return fmt.Errorf("failed to audit voting process for vote %s: %v", voteID, err)
    }
    fmt.Printf("Voting process audited for vote %s.\n", voteID)
    return nil
}

// GovernanceRevertVotingProcess reverts a previously audited voting process.
func GovernanceRevertVotingProcess(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertVotingProcess(voteID); err != nil {
        return fmt.Errorf("failed to revert voting process for vote %s: %v", voteID, err)
    }
    fmt.Printf("Voting process reverted for vote %s.\n", voteID)
    return nil
}

// GovernanceInitiateRecallVote initiates a recall vote process.
func GovernanceInitiateRecallVote(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitiateRecallVote(voteID); err != nil {
        return fmt.Errorf("failed to initiate recall vote for %s: %v", voteID, err)
    }
    fmt.Printf("Recall vote initiated for %s.\n", voteID)
    return nil
}

// GovernanceConfirmRecallVote confirms a recall vote process.
func GovernanceConfirmRecallVote(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmRecallVote(voteID); err != nil {
        return fmt.Errorf("failed to confirm recall vote for %s: %v", voteID, err)
    }
    fmt.Printf("Recall vote confirmed for %s.\n", voteID)
    return nil
}

// GovernanceRevertRecallVote reverts a recall vote.
func GovernanceRevertRecallVote(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertRecallVote(voteID); err != nil {
        return fmt.Errorf("failed to revert recall vote for %s: %v", voteID, err)
    }
    fmt.Printf("Recall vote reverted for %s.\n", voteID)
    return nil
}

// GovernanceFetchDelegateStatus retrieves the status of a delegate.
func GovernanceFetchDelegateStatus(delegateID string, ledgerInstance *ledger.Ledger) (common.DelegateStatus, error) {
    status, err := ledgerInstance.FetchDelegateStatus(delegateID)
    if err != nil {
        return common.DelegateStatus{}, fmt.Errorf("failed to fetch delegate status for %s: %v", delegateID, err)
    }
    return status, nil
}

// GovernanceUpdateDelegateStatus updates the status of a delegate.
func GovernanceUpdateDelegateStatus(delegateID string, status common.DelegateStatus, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UpdateDelegateStatus(delegateID, status); err != nil {
        return fmt.Errorf("failed to update delegate status for %s: %v", delegateID, err)
    }
    fmt.Printf("Delegate status updated for %s.\n", delegateID)
    return nil
}

// GovernanceRemoveDelegate removes a delegate from voting.
func GovernanceRemoveDelegate(delegateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RemoveDelegate(delegateID); err != nil {
        return fmt.Errorf("failed to remove delegate %s: %v", delegateID, err)
    }
    fmt.Printf("Delegate %s removed.\n", delegateID)
    return nil
}

// GovernanceReassignDelegate reassigns a delegate role.
func GovernanceReassignDelegate(delegateID, newDelegateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReassignDelegate(delegateID, newDelegateID); err != nil {
        return fmt.Errorf("failed to reassign delegate %s to %s: %v", delegateID, newDelegateID, err)
    }
    fmt.Printf("Delegate %s reassigned to %s.\n", delegateID, newDelegateID)
    return nil
}

// GovernanceLockVoteTokens locks tokens for a vote.
func GovernanceLockVoteTokens(voteID string, amount int64, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptAmount(amount)
    if err := ledgerInstance.LockTokens(voteID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to lock vote tokens for vote %s: %v", voteID, err)
    }
    fmt.Printf("Vote tokens locked for vote %s.\n", voteID)
    return nil
}

// GovernanceUnlockVoteTokens unlocks tokens after a vote.
func GovernanceUnlockVoteTokens(voteID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockTokens(voteID); err != nil {
        return fmt.Errorf("failed to unlock vote tokens for vote %s: %v", voteID, err)
    }
    fmt.Printf("Vote tokens unlocked for vote %s.\n", voteID)
    return nil
}

// GovernanceTransferVotingPower transfers voting power to another entity.
func GovernanceTransferVotingPower(delegateID, targetID string, power int64, ledgerInstance *ledger.Ledger) error {
    encryptedPower := encryption.EncryptAmount(power)
    if err := ledgerInstance.TransferVotingPower(delegateID, targetID, encryptedPower); err != nil {
        return fmt.Errorf("failed to transfer voting power from %s to %s: %v", delegateID, targetID, err)
    }
    fmt.Printf("Voting power transferred from %s to %s.\n", delegateID, targetID)
    return nil
}

// GovernanceVerifyVotingPower verifies the current voting power of a delegate.
func GovernanceVerifyVotingPower(delegateID string, ledgerInstance *ledger.Ledger) (int64, error) {
    power, err := ledgerInstance.FetchVotingPower(delegateID)
    if err != nil {
        return 0, fmt.Errorf("failed to verify voting power for %s: %v", delegateID, err)
    }
    return power, nil
}
