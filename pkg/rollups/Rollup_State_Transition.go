
package rollups

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// RollupRecordStateChange logs and records a state change in the rollup.
func RollupRecordStateChange(stateID string, changeDetails string, ledgerInstance *ledger.Ledger) error {
    encryptedDetails := encryption.EncryptData(changeDetails)
    if err := ledgerInstance.RecordStateChange(stateID, encryptedDetails); err != nil {
        return fmt.Errorf("failed to record state change for %s: %v", stateID, err)
    }
    fmt.Printf("State change recorded for %s.\n", stateID)
    return nil
}

// RollupProcessInclusionProof processes an inclusion proof in the rollup.
func RollupProcessInclusionProof(proofID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ProcessInclusionProof(proofID); err != nil {
        return fmt.Errorf("failed to process inclusion proof for %s: %v", proofID, err)
    }
    fmt.Printf("Inclusion proof processed for %s.\n", proofID)
    return nil
}

// RollupVerifyInclusion verifies an inclusion proof in the rollup.
func RollupVerifyInclusion(proofID string, ledgerInstance *ledger.Ledger) (bool, error) {
    isValid, err := ledgerInstance.VerifyInclusionProof(proofID)
    if err != nil {
        return false, fmt.Errorf("failed to verify inclusion proof for %s: %v", proofID, err)
    }
    fmt.Printf("Inclusion proof verification result for %s: %v\n", proofID, isValid)
    return isValid, nil
}

// RollupTrackInclusionProof tracks the status of an inclusion proof.
func RollupTrackInclusionProof(proofID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.TrackInclusionProof(proofID); err != nil {
        return fmt.Errorf("failed to track inclusion proof for %s: %v", proofID, err)
    }
    fmt.Printf("Inclusion proof tracked for %s.\n", proofID)
    return nil
}

// RollupFinalizeInclusionProof finalizes an inclusion proof in the rollup.
func RollupFinalizeInclusionProof(proofID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeInclusionProof(proofID); err != nil {
        return fmt.Errorf("failed to finalize inclusion proof for %s: %v", proofID, err)
    }
    fmt.Printf("Inclusion proof finalized for %s.\n", proofID)
    return nil
}

// RollupCommitState commits the current state in the rollup.
func RollupCommitState(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CommitState(stateID); err != nil {
        return fmt.Errorf("failed to commit state for %s: %v", stateID, err)
    }
    fmt.Printf("State committed for %s.\n", stateID)
    return nil
}

// RollupRevertState reverts to the previous state in the rollup.
func RollupRevertState(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertState(stateID); err != nil {
        return fmt.Errorf("failed to revert state for %s: %v", stateID, err)
    }
    fmt.Printf("State reverted for %s.\n", stateID)
    return nil
}

// RollupValidateStateTransition validates a state transition in the rollup.
func RollupValidateStateTransition(transitionData string, ledgerInstance *ledger.Ledger) error {
    encryptedData := encryption.EncryptData(transitionData)
    if err := ledgerInstance.ValidateStateTransition(encryptedData); err != nil {
        return fmt.Errorf("failed to validate state transition: %v", err)
    }
    fmt.Println("State transition validated.")
    return nil
}

// RollupStateLock locks the current rollup state.
func RollupStateLock(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockState(stateID); err != nil {
        return fmt.Errorf("failed to lock state for %s: %v", stateID, err)
    }
    fmt.Printf("State locked for %s.\n", stateID)
    return nil
}

// RollupStateUnlock unlocks the current rollup state.
func RollupStateUnlock(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockState(stateID); err != nil {
        return fmt.Errorf("failed to unlock state for %s: %v", stateID, err)
    }
    fmt.Printf("State unlocked for %s.\n", stateID)
    return nil
}

// RollupQueryState queries the current rollup state.
func RollupQueryState(stateID string, ledgerInstance *ledger.Ledger) (string, error) {
    stateData, err := ledgerInstance.QueryState(stateID)
    if err != nil {
        return "", fmt.Errorf("failed to query state for %s: %v", stateID, err)
    }
    fmt.Printf("State data queried for %s.\n", stateID)
    return stateData, nil
}

// RollupStoreState stores the rollup state.
func RollupStoreState(stateID string, stateData string, ledgerInstance *ledger.Ledger) error {
    encryptedData := encryption.EncryptData(stateData)
    if err := ledgerInstance.StoreState(stateID, encryptedData); err != nil {
        return fmt.Errorf("failed to store state for %s: %v", stateID, err)
    }
    fmt.Printf("State stored for %s.\n", stateID)
    return nil
}

// RollupFetchState retrieves stored state data in the rollup.
func RollupFetchState(stateID string, ledgerInstance *ledger.Ledger) (string, error) {
    stateData, err := ledgerInstance.FetchState(stateID)
    if err != nil {
        return "", fmt.Errorf("failed to fetch state for %s: %v", stateID, err)
    }
    fmt.Printf("State fetched for %s.\n", stateID)
    return stateData, nil
}

// RollupUpdateTransactionRoot updates the transaction root for a batch.
func RollupUpdateTransactionRoot(batchID string, rootHash string, ledgerInstance *ledger.Ledger) error {
    encryptedRootHash := encryption.EncryptData(rootHash)
    if err := ledgerInstance.UpdateTransactionRoot(batchID, encryptedRootHash); err != nil {
        return fmt.Errorf("failed to update transaction root for %s: %v", batchID, err)
    }
    fmt.Printf("Transaction root updated for %s.\n", batchID)
    return nil
}

// RollupVerifyTransactionRoot verifies the transaction root of a batch.
func RollupVerifyTransactionRoot(batchID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.VerifyTransactionRoot(batchID); err != nil {
        return fmt.Errorf("failed to verify transaction root for %s: %v", batchID, err)
    }
    fmt.Printf("Transaction root verified for %s.\n", batchID)
    return nil
}

// RollupSetBatchFinalized sets a batch as finalized in the rollup.
func RollupSetBatchFinalized(batchID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetBatchFinalized(batchID); err != nil {
        return fmt.Errorf("failed to set batch as finalized for %s: %v", batchID, err)
    }
    fmt.Printf("Batch %s set as finalized.\n", batchID)
    return nil
}

// RollupGetBatchFinalized retrieves the finalization status of a batch.
func RollupGetBatchFinalized(batchID string, ledgerInstance *ledger.Ledger) (bool, error) {
    isFinalized, err := ledgerInstance.GetBatchFinalized(batchID)
    if err != nil {
        return false, fmt.Errorf("failed to get finalization status for batch %s: %v", batchID, err)
    }
    fmt.Printf("Finalization status for %s: %v\n", batchID, isFinalized)
    return isFinalized, nil
}

// RollupArchiveBatch archives a completed batch in the rollup.
func RollupArchiveBatch(batchID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ArchiveBatch(batchID); err != nil {
        return fmt.Errorf("failed to archive batch %s: %v", batchID, err)
    }
    fmt.Printf("Batch %s archived.\n", batchID)
    return nil
}

// RollupRestoreBatch restores an archived batch in the rollup.
func RollupRestoreBatch(batchID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RestoreBatch(batchID); err != nil {
        return fmt.Errorf("failed to restore batch %s: %v", batchID, err)
    }
    fmt.Printf("Batch %s restored.\n", batchID)
    return nil
}

// RollupBackupState creates a backup of the current rollup state.
func RollupBackupState(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.BackupState(stateID); err != nil {
        return fmt.Errorf("failed to backup state for %s: %v", stateID, err)
    }
    fmt.Printf("State %s backed up.\n", stateID)
    return nil
}

// RollupRecoverState recovers the rollup state from a backup.
func RollupRecoverState(stateID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RecoverState(stateID); err != nil {
        return fmt.Errorf("failed to recover state for %s: %v", stateID, err)
    }
    fmt.Printf("State %s recovered.\n", stateID)
    return nil
}
