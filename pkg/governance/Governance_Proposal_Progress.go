package governance


import (
    "fmt"
    "time"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// GovernanceMonitorProposalProgress monitors the progress of a proposal to ensure it aligns with governance standards.
func GovernanceMonitorProposalProgress(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.GovernanceLedger.MonitorProposalProgress(proposalID); err != nil {
        return fmt.Errorf("failed to monitor proposal progress for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Monitoring progress for proposal %s.\n", proposalID)
    return nil
}

// GovernanceValidateProposalProgress validates the ongoing progress of a proposal.
func GovernanceValidateProposalProgress(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.GovernanceLedger.ValidateProposalProgress(proposalID); err != nil {
        return fmt.Errorf("failed to validate proposal progress for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal progress validated for proposal %s.\n", proposalID)
    return nil
}

// GovernanceRevertProposalProgress reverts any changes in proposal progress that are found to be invalid.
func GovernanceRevertProposalProgress(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.GovernanceLedger.RevertProposalProgress(proposalID); err != nil {
        return fmt.Errorf("failed to revert proposal progress for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Proposal progress reverted for proposal %s.\n", proposalID)
    return nil
}

// GovernanceEscrowFunds escrows funds for a proposal as part of governance compliance.
func GovernanceEscrowFunds(proposalID string, amount int, ledgerInstance *ledger.Ledger) error {
    encryptedAmount := encryption.EncryptAmount(amount)
    if err := ledgerInstance.GovernanceLedger.EscrowFunds(proposalID, encryptedAmount); err != nil {
        return fmt.Errorf("failed to escrow funds for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Funds escrowed for proposal %s.\n", proposalID)
    return nil
}

// GovernanceReleaseEscrowedFunds releases escrowed funds for a completed proposal.
func GovernanceReleaseEscrowedFunds(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseEscrowedFunds(proposalID); err != nil {
        return fmt.Errorf("failed to release escrowed funds for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Escrowed funds released for proposal %s.\n", proposalID)
    return nil
}

// GovernanceSetProposalExpiration sets an expiration date for a proposal.
func GovernanceSetProposalExpiration(proposalID string, expirationDate time.Time, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetProposalExpiration(proposalID, expirationDate); err != nil {
        return fmt.Errorf("failed to set expiration date for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Expiration date set for proposal %s.\n", proposalID)
    return nil
}

// GovernanceGetProposalExpiration retrieves the expiration date for a proposal.
func GovernanceGetProposalExpiration(proposalID string, ledgerInstance *ledger.Ledger) (time.Time, error) {
    expirationDate, err := ledgerInstance.GetProposalExpiration(proposalID)
    if err != nil {
        return time.Time{}, fmt.Errorf("failed to get expiration date for proposal %s: %v", proposalID, err)
    }
    return expirationDate, nil
}

// GovernanceExtendProposalExpiration extends the expiration date of a proposal.
func GovernanceExtendProposalExpiration(proposalID string, extension time.Duration, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ExtendProposalExpiration(proposalID, extension); err != nil {
        return fmt.Errorf("failed to extend expiration date for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Expiration date extended for proposal %s.\n", proposalID)
    return nil
}

// GovernanceValidateProposalExpiration checks whether a proposal has expired.
func GovernanceValidateProposalExpiration(proposalID string, ledgerInstance *ledger.Ledger) (bool, error) {
    expired, err := ledgerInstance.ValidateProposalExpiration(proposalID)
    if err != nil {
        return false, fmt.Errorf("failed to validate expiration for proposal %s: %v", proposalID, err)
    }
    return expired, nil
}

// GovernanceFinalizeProposalExpiration finalizes a proposal’s expiration.
func GovernanceFinalizeProposalExpiration(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeProposalExpiration(proposalID); err != nil {
        return fmt.Errorf("failed to finalize expiration for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Expiration finalized for proposal %s.\n", proposalID)
    return nil
}

// GovernanceSetVotingPeriod sets a voting period for a proposal.
func GovernanceSetVotingPeriod(proposalID string, votingPeriod time.Duration, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetVotingPeriod(proposalID, votingPeriod); err != nil {
        return fmt.Errorf("failed to set voting period for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Voting period set for proposal %s.\n", proposalID)
    return nil
}

// GovernanceGetVotingPeriod retrieves the voting period for a proposal.
func GovernanceGetVotingPeriod(proposalID string, ledgerInstance *ledger.Ledger) (time.Duration, error) {
    votingPeriod, err := ledgerInstance.GetVotingPeriod(proposalID)
    if err != nil {
        return 0, fmt.Errorf("failed to get voting period for proposal %s: %v", proposalID, err)
    }
    return votingPeriod, nil
}

// GovernanceExtendVotingPeriod extends the voting period of a proposal.
func GovernanceExtendVotingPeriod(proposalID string, extension time.Duration, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ExtendVotingPeriod(proposalID, extension); err != nil {
        return fmt.Errorf("failed to extend voting period for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Voting period extended for proposal %s.\n", proposalID)
    return nil
}

// GovernanceFinalizeVotingPeriod finalizes the voting period of a proposal.
func GovernanceFinalizeVotingPeriod(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeVotingPeriod(proposalID); err != nil {
        return fmt.Errorf("failed to finalize voting period for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Voting period finalized for proposal %s.\n", proposalID)
    return nil
}

// GovernanceTrackComplianceMetrics tracks compliance metrics for a proposal.
func GovernanceTrackComplianceMetrics(proposalID string, metrics common.ComplianceMetrics, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.TrackComplianceMetrics(proposalID, metrics); err != nil {
        return fmt.Errorf("failed to track compliance metrics for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance metrics tracked for proposal %s.\n", proposalID)
    return nil
}

// GovernanceRecordComplianceMetrics records compliance metrics in the ledger.
func GovernanceRecordComplianceMetrics(proposalID string, metrics common.ComplianceMetrics, ledgerInstance *ledger.Ledger) error {
    encryptedMetrics := encryption.EncryptMetrics(metrics)
    if err := ledgerInstance.RecordComplianceMetrics(proposalID, encryptedMetrics); err != nil {
        return fmt.Errorf("failed to record compliance metrics for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance metrics recorded for proposal %s.\n", proposalID)
    return nil
}

// GovernanceValidateComplianceMetrics validates compliance metrics for a proposal.
func GovernanceValidateComplianceMetrics(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ValidateComplianceMetrics(proposalID); err != nil {
        return fmt.Errorf("failed to validate compliance metrics for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance metrics validated for proposal %s.\n", proposalID)
    return nil
}

// GovernanceFinalizeComplianceMetrics finalizes the compliance metrics for a proposal.
func GovernanceFinalizeComplianceMetrics(proposalID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeComplianceMetrics(proposalID); err != nil {
        return fmt.Errorf("failed to finalize compliance metrics for proposal %s: %v", proposalID, err)
    }
    fmt.Printf("Compliance metrics finalized for proposal %s.\n", proposalID)
    return nil
}

// GovernanceLogGovernanceEvent logs a governance event for audit purposes.
func GovernanceLogGovernanceEvent(event common.GovernanceEvent, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LogGovernanceEvent(event); err != nil {
        return fmt.Errorf("failed to log governance event: %v", err)
    }
    fmt.Println("Governance event logged.")
    return nil
}

// GovernanceMonitorGovernanceEvent monitors ongoing governance events in real-time.
func GovernanceMonitorGovernanceEvent(eventID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.MonitorGovernanceEvent(eventID); err != nil {
        return fmt.Errorf("failed to monitor governance event with ID %s: %v", eventID, err)
    }
    fmt.Printf("Monitoring governance event with ID %s.\n", eventID)
    return nil
}

// GovernanceAuditGovernanceEvent audits past governance events for compliance.
func GovernanceAuditGovernanceEvent(eventID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditGovernanceEvent(eventID); err != nil {
        return fmt.Errorf("failed to audit governance event with ID %s: %v", eventID, err)
    }
    fmt.Printf("Governance event with ID %s audited.\n", eventID)
    return nil
}
