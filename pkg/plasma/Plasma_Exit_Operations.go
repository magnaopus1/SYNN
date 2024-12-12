// Plasma_Exit_Operations.go

package main

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// PlasmaGetProofStatus retrieves the status of an exit proof.
func PlasmaGetProofStatus(proofID string, ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.GetProofStatus(proofID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve proof status: %v", err)
    }
    return status, nil
}

// PlasmaSettleChallenge settles a challenge on the Plasma chain.
func PlasmaSettleChallenge(challengeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SettleChallenge(challengeID); err != nil {
        return fmt.Errorf("failed to settle challenge: %v", err)
    }
    fmt.Printf("Challenge %s settled.\n", challengeID)
    return nil
}

// PlasmaVerifyExit verifies an exit request on the Plasma chain.
func PlasmaVerifyExit(exitID string, ledgerInstance *ledger.Ledger) (bool, error) {
    isValid, err := ledgerInstance.VerifyExit(exitID)
    if err != nil {
        return false, fmt.Errorf("failed to verify exit: %v", err)
    }
    return isValid, nil
}

// PlasmaProcessChallenge processes a challenge for an exit.
func PlasmaProcessChallenge(challengeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ProcessChallenge(challengeID); err != nil {
        return fmt.Errorf("failed to process challenge: %v", err)
    }
    fmt.Printf("Challenge %s processed.\n", challengeID)
    return nil
}

// PlasmaValidateExit validates an exit request.
func PlasmaValidateExit(exitID string, ledgerInstance *ledger.Ledger) (bool, error) {
    isValid, err := ledgerInstance.ValidateExit(exitID)
    if err != nil {
        return false, fmt.Errorf("failed to validate exit: %v", err)
    }
    return isValid, nil
}

// PlasmaRevertExit reverts an exit transaction.
func PlasmaRevertExit(exitID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertExit(exitID); err != nil {
        return fmt.Errorf("failed to revert exit: %v", err)
    }
    fmt.Printf("Exit %s reverted.\n", exitID)
    return nil
}

// PlasmaEscrowDeposit places a deposit into escrow.
func PlasmaEscrowDeposit(depositID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.EscrowDeposit(depositID, amount); err != nil {
        return fmt.Errorf("failed to escrow deposit: %v", err)
    }
    fmt.Printf("Deposit %s escrowed with amount %d.\n", depositID, amount)
    return nil
}

// PlasmaReleaseEscrow releases funds from escrow.
func PlasmaReleaseEscrow(depositID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseEscrow(depositID); err != nil {
        return fmt.Errorf("failed to release escrow: %v", err)
    }
    fmt.Printf("Escrow %s released.\n", depositID)
    return nil
}

// PlasmaInitiateExit initiates an exit for a Plasma transaction.
func PlasmaInitiateExit(accountID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.InitiateExit(accountID); err != nil {
        return fmt.Errorf("failed to initiate exit for account %s: %v", accountID, err)
    }
    fmt.Printf("Exit initiated for account %s.\n", accountID)
    return nil
}

// PlasmaCancelExit cancels a previously requested exit.
func PlasmaCancelExit(exitID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.CancelExit(exitID); err != nil {
        return fmt.Errorf("failed to cancel exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s cancelled.\n", exitID)
    return nil
}

// PlasmaUpdateState updates the state of an account.
func PlasmaUpdateState(accountID string, stateData string, ledgerInstance *ledger.Ledger) error {
    encryptedState := encryption.EncryptData(stateData)
    if err := ledgerInstance.UpdateState(accountID, encryptedState); err != nil {
        return fmt.Errorf("failed to update state for account %s: %v", accountID, err)
    }
    fmt.Printf("State updated for account %s.\n", accountID)
    return nil
}

// PlasmaFreezeAccount freezes an account on the Plasma chain.
func PlasmaFreezeAccount(accountID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FreezeAccount(accountID); err != nil {
        return fmt.Errorf("failed to freeze account %s: %v", accountID, err)
    }
    fmt.Printf("Account %s frozen.\n", accountID)
    return nil
}

// PlasmaUnfreezeAccount unfreezes an account on the Plasma chain.
func PlasmaUnfreezeAccount(accountID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnfreezeAccount(accountID); err != nil {
        return fmt.Errorf("failed to unfreeze account %s: %v", accountID, err)
    }
    fmt.Printf("Account %s unfrozen.\n", accountID)
    return nil
}

// PlasmaSnapshotState captures a snapshot of the Plasma chain state.
func PlasmaSnapshotState(ledgerInstance *ledger.Ledger) (string, error) {
    snapshotID, err := ledgerInstance.SnapshotState()
    if err != nil {
        return "", fmt.Errorf("failed to snapshot state: %v", err)
    }
    fmt.Printf("Snapshot created with ID %s.\n", snapshotID)
    return snapshotID, nil
}

// PlasmaLoadSnapshot loads a snapshot into the Plasma chain state.
func PlasmaLoadSnapshot(snapshotID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LoadSnapshot(snapshotID); err != nil {
        return fmt.Errorf("failed to load snapshot %s: %v", snapshotID, err)
    }
    fmt.Printf("Snapshot %s loaded.\n", snapshotID)
    return nil
}

// PlasmaGetSnapshotData retrieves data from a state snapshot.
func PlasmaGetSnapshotData(snapshotID string, ledgerInstance *ledger.Ledger) (string, error) {
    data, err := ledgerInstance.GetSnapshotData(snapshotID)
    if err != nil {
        return "", fmt.Errorf("failed to get snapshot data for %s: %v", snapshotID, err)
    }
    return data, nil
}

// PlasmaSetRootHash sets the root hash for the Plasma Merkle tree.
func PlasmaSetRootHash(rootHash string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.SetRootHash(rootHash); err != nil {
        return fmt.Errorf("failed to set root hash: %v", err)
    }
    fmt.Printf("Root hash set to %s.\n", rootHash)
    return nil
}

// PlasmaGetRootHash retrieves the current root hash.
func PlasmaGetRootHash(ledgerInstance *ledger.Ledger) (string, error) {
    rootHash, err := ledgerInstance.GetRootHash()
    if err != nil {
        return "", fmt.Errorf("failed to get root hash: %v", err)
    }
    return rootHash, nil
}

// PlasmaMerkleProof generates a Merkle proof for a transaction.
func PlasmaMerkleProof(txID string, ledgerInstance *ledger.Ledger) (string, error) {
    proof, err := ledgerInstance.GenerateMerkleProof(txID)
    if err != nil {
        return "", fmt.Errorf("failed to generate Merkle proof: %v", err)
    }
    return proof, nil
}

// PlasmaMerkleVerify verifies a Merkle proof.
func PlasmaMerkleVerify(txID string, proof string, ledgerInstance *ledger.Ledger) (bool, error) {
    isValid, err := ledgerInstance.VerifyMerkleProof(txID, proof)
    if err != nil {
        return false, fmt.Errorf("failed to verify Merkle proof: %v", err)
    }
    return isValid, nil
}
