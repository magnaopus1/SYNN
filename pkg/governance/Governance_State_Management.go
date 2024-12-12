package governance


import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// GovernanceSetProposalPriority sets the priority for a proposal.
func GovernanceSetProposalPriority(proposalID string, priority int, ledgerInstance *ledger.Ledger) error {
    encryptedPriority := encryption.EncryptInt(priority)
    if err := ledgerInstance.SetProposalPriority(proposalID, encryptedPriority); err != nil {
        return fmt.Errorf("failed to set priority for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Priority set for proposal %s.\n", proposalID)
    return nil
}

// GovernanceGetProposalPriority retrieves the priority for a proposal.
func GovernanceGetProposalPriority(proposalID string, ledgerInstance *ledger.Ledger) (int, error) {
    encryptedPriority, err := ledgerInstance.GetProposalPriority(proposalID)
    if err != nil {
        return 0, fmt.Errorf("failed to get priority for proposal %s: %v", proposalID, err)
    }
    priority := encryption.DecryptInt(encryptedPriority)
    return priority, nil
}

// GovernanceAdjustProposalPriority dynamically adjusts the priority of a proposal.
func GovernanceAdjustProposalPriority(proposalID string, adjustment int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AdjustProposalPriority(proposalID, adjustment); err != nil {
        return fmt.Errorf("failed to adjust priority for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Priority adjusted for proposal %s.\n", proposalID)
    return nil
}

// GovernanceFinalizeProposalPriority finalizes the priority for a proposal.
func GovernanceFinalizeProposalPriority(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeProposalPriority(proposalID); err != nil {
        return fmt.Errorf("failed to finalize priority for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Priority finalized for proposal %s.\n", proposalID)
    return nil
}

// GovernanceLockProposalState locks the state of a proposal.
func GovernanceLockProposalState(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockProposalState(proposalID); err != nil {
        return fmt.Errorf("failed to lock state for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("State locked for proposal %s.\n", proposalID)
    return nil
}

// GovernanceUnlockProposalState unlocks the state of a proposal.
func GovernanceUnlockProposalState(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockProposalState(proposalID); err != nil {
        return fmt.Errorf("failed to unlock state for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("State unlocked for proposal %s.\n", proposalID)
    return nil
}

// GovernanceLogStateChange logs a change in the state of a proposal.
func GovernanceLogStateChange(proposalID string, stateChange common.StateChange, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LogStateChange(proposalID, stateChange); err != nil {
        return fmt.Errorf("failed to log state change for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("State change logged for proposal %s.\n", proposalID)
    return nil
}

// GovernanceMonitorStateChange monitors a state change in a proposal.
func GovernanceMonitorStateChange(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MonitorStateChange(proposalID); err != nil {
        return fmt.Errorf("failed to monitor state change for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Monitoring state change for proposal %s.\n", proposalID)
    return nil
}

// GovernanceValidateStateChange validates the recent state change.
func GovernanceValidateStateChange(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ValidateStateChange(proposalID); err != nil {
        return fmt.Errorf("failed to validate state change for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("State change validated for proposal %s.\n", proposalID)
    return nil
}

// GovernanceRevertStateChange reverts a state change for a proposal if necessary.
func GovernanceRevertStateChange(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertStateChange(proposalID); err != nil {
        return fmt.Errorf("failed to revert state change for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("State change reverted for proposal %s.\n", proposalID)
    return nil
}

// GovernanceRecordProposalAudit records an audit of a proposal.
func GovernanceRecordProposalAudit(proposalID string, audit common.AuditRecord, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RecordProposalAudit(proposalID, audit); err != nil {
        return fmt.Errorf("failed to record audit for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Audit recorded for proposal %s.\n", proposalID)
    return nil
}

// GovernanceFetchProposalAudit fetches an audit record for a proposal.
func GovernanceFetchProposalAudit(proposalID string, ledgerInstance *ledger.Ledger) ([]common.AuditRecord, error) {
    auditRecords, err := ledgerInstance.FetchProposalAudit(proposalID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch audit for proposal %s: %v", proposalID, err)
    }
    return auditRecords, nil
}

// GovernanceStoreProposalAudit stores an audit record for a proposal in a secure way.
func GovernanceStoreProposalAudit(proposalID string, audit common.AuditRecord, ledgerInstance *ledger.Ledger) error {
    encryptedAudit := encryption.EncryptAuditRecord(audit)
    if err := ledgerInstance.StoreProposalAudit(proposalID, encryptedAudit); err != nil {
        return fmt.Errorf("failed to store audit for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Audit stored for proposal %s.\n", proposalID)
    return nil
}

// GovernanceReconcileAuditRecords reconciles the audit records of a proposal.
func GovernanceReconcileAuditRecords(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReconcileAuditRecords(proposalID); err != nil {
        return fmt.Errorf("failed to reconcile audit records for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Audit records reconciled for proposal %s.\n", proposalID)
    return nil
}

// GovernanceFinalizeAuditRecords finalizes the audit records for compliance.
func GovernanceFinalizeAuditRecords(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeAuditRecords(proposalID); err != nil {
        return fmt.Errorf("failed to finalize audit records for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Audit records finalized for proposal %s.\n", proposalID)
    return nil
}

// GovernanceVerifyComplianceStatus verifies the compliance status of a proposal.
func GovernanceVerifyComplianceStatus(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.VerifyComplianceStatus(proposalID); err != nil {
        return fmt.Errorf("failed to verify compliance status for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance status verified for proposal %s.\n", proposalID)
    return nil
}

// GovernanceInitiateComplianceReview initiates a compliance review for a proposal.
func GovernanceInitiateComplianceReview(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitiateComplianceReview(proposalID); err != nil {
        return fmt.Errorf("failed to initiate compliance review for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance review initiated for proposal %s.\n", proposalID)
    return nil
}

// GovernanceConfirmComplianceReview confirms the compliance review results.
func GovernanceConfirmComplianceReview(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmComplianceReview(proposalID); err != nil {
        return fmt.Errorf("failed to confirm compliance review for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance review confirmed for proposal %s.\n", proposalID)
    return nil
}

// GovernanceRevertComplianceReview reverts a compliance review for a proposal.
func GovernanceRevertComplianceReview(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertComplianceReview(proposalID); err != nil {
        return fmt.Errorf("failed to revert compliance review for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance review reverted for proposal %s.\n", proposalID)
    return nil
}

// GovernanceTrackProposalProgress tracks the progress of a proposal through its stages.
func GovernanceTrackProposalProgress(proposalID string, progress common.ProgressRecord, ledgerInstance *ledger.Ledger) error {
    encryptedProgress := encryption.EncryptProgressRecord(progress)
    if err := ledgerInstance.TrackProposalProgress(proposalID, encryptedProgress); err != nil {
        return fmt.Errorf("failed to track progress for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Progress tracked for proposal %s.\n", proposalID)
    return nil
}
