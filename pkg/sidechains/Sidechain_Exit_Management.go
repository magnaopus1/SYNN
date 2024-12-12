// Sidechain_Exit_Management.go

package sidechains

import (
    "fmt"
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/encryption"
)

// SidechainProcessExit processes an exit request within the sidechain.
func SidechainProcessExit(exitID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ProcessExit(exitID); err != nil {
        return fmt.Errorf("failed to process exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s processed successfully.\n", exitID)
    return nil
}

// SidechainInitiateExit initiates an exit request within the sidechain.
func SidechainInitiateExit(exitID string, details common.ExitDetails, ledgerInstance *ledger.Ledger) error {
    encryptedDetails := encryption.EncryptExitDetails(details)
    if err := ledgerInstance.InitiateExit(exitID, encryptedDetails); err != nil {
        return fmt.Errorf("failed to initiate exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s initiated.\n", exitID)
    return nil
}

// SidechainConfirmExit confirms an exit request after validation.
func SidechainConfirmExit(exitID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ConfirmExit(exitID); err != nil {
        return fmt.Errorf("failed to confirm exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s confirmed.\n", exitID)
    return nil
}

// SidechainRevertExit reverts an exit process to its prior state.
func SidechainRevertExit(exitID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertExit(exitID); err != nil {
        return fmt.Errorf("failed to revert exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s reverted.\n", exitID)
    return nil
}

// SidechainChallengeExit submits a challenge against an exit.
func SidechainChallengeExit(exitID string, challengeDetails common.ChallengeDetails, ledgerInstance *ledger.Ledger) error {
    encryptedChallenge := encryption.EncryptChallengeDetails(challengeDetails)
    if err := ledgerInstance.ChallengeExit(exitID, encryptedChallenge); err != nil {
        return fmt.Errorf("failed to challenge exit %s: %v", exitID, err)
    }
    fmt.Printf("Exit %s challenged.\n", exitID)
    return nil
}

// SidechainProcessChallenge processes a challenge on an exit request.
func SidechainProcessChallenge(challengeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ProcessChallenge(challengeID); err != nil {
        return fmt.Errorf("failed to process challenge %s: %v", challengeID, err)
    }
    fmt.Printf("Challenge %s processed.\n", challengeID)
    return nil
}

// SidechainFinalizeChallenge finalizes a challenge outcome on an exit request.
func SidechainFinalizeChallenge(challengeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.FinalizeChallenge(challengeID); err != nil {
        return fmt.Errorf("failed to finalize challenge %s: %v", challengeID, err)
    }
    fmt.Printf("Challenge %s finalized.\n", challengeID)
    return nil
}

// SidechainEscrowExitTokens places tokens in escrow for an exit.
func SidechainEscrowExitTokens(tokenID string, amount int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.EscrowTokens(tokenID, amount); err != nil {
        return fmt.Errorf("failed to escrow tokens %s: %v", tokenID, err)
    }
    fmt.Printf("Tokens %s escrowed for exit.\n", tokenID)
    return nil
}

// SidechainReleaseExitTokens releases tokens from escrow upon exit completion.
func SidechainReleaseExitTokens(tokenID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ReleaseEscrowedTokens(tokenID); err != nil {
        return fmt.Errorf("failed to release escrowed tokens %s: %v", tokenID, err)
    }
    fmt.Printf("Escrowed tokens %s released.\n", tokenID)
    return nil
}

// SidechainTrackAssetMovement tracks the movement of assets for exit processes.
func SidechainTrackAssetMovement(assetID string, movementDetails common.MovementDetails, ledgerInstance *ledger.Ledger) error {
    encryptedDetails := encryption.EncryptMovementDetails(movementDetails)
    if err := ledgerInstance.TrackAssetMovement(assetID, encryptedDetails); err != nil {
        return fmt.Errorf("failed to track asset movement %s: %v", assetID, err)
    }
    fmt.Printf("Asset movement %s tracked.\n", assetID)
    return nil
}

// SidechainValidateAssetMovement validates recorded asset movements.
func SidechainValidateAssetMovement(assetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ValidateAssetMovement(assetID); err != nil {
        return fmt.Errorf("failed to validate asset movement %s: %v", assetID, err)
    }
    fmt.Printf("Asset movement %s validated.\n", assetID)
    return nil
}

// SidechainLogAssetMovement logs asset movement within the exit process.
func SidechainLogAssetMovement(assetID string, movementDetails common.MovementDetails, ledgerInstance *ledger.Ledger) error {
    encryptedDetails := encryption.EncryptMovementDetails(movementDetails)
    if err := ledgerInstance.LogAssetMovement(assetID, encryptedDetails); err != nil {
        return fmt.Errorf("failed to log asset movement %s: %v", assetID, err)
    }
    fmt.Printf("Asset movement %s logged.\n", assetID)
    return nil
}

// SidechainAuditAssetMovement audits all recorded movements of assets.
func SidechainAuditAssetMovement(assetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.AuditAssetMovement(assetID); err != nil {
        return fmt.Errorf("failed to audit asset movement %s: %v", assetID, err)
    }
    fmt.Printf("Asset movement %s audited.\n", assetID)
    return nil
}

// SidechainVerifyAssetMovement verifies asset movement records.
func SidechainVerifyAssetMovement(assetID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.VerifyAssetMovement(assetID); err != nil {
        return fmt.Errorf("failed to verify asset movement %s: %v", assetID, err)
    }
    fmt.Printf("Asset movement %s verified.\n", assetID)
    return nil
}

// SidechainLockTransaction locks a transaction to prevent any unauthorized modifications.
func SidechainLockTransaction(transactionID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.LockTransaction(transactionID); err != nil {
        return fmt.Errorf("failed to lock transaction %s: %v", transactionID, err)
    }
    fmt.Printf("Transaction %s locked.\n", transactionID)
    return nil
}

// SidechainUnlockTransaction unlocks a previously locked transaction.
func SidechainUnlockTransaction(transactionID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.UnlockTransaction(transactionID); err != nil {
        return fmt.Errorf("failed to unlock transaction %s: %v", transactionID, err)
    }
    fmt.Printf("Transaction %s unlocked.\n", transactionID)
    return nil
}

// SidechainRevertTransactionLock reverts a transaction lock status.
func SidechainRevertTransactionLock(transactionID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.RevertTransactionLock(transactionID); err != nil {
        return fmt.Errorf("failed to revert transaction lock %s: %v", transactionID, err)
    }
    fmt.Printf("Transaction lock %s reverted.\n", transactionID)
    return nil
}
