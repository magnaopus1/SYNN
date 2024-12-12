package rollups

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/consensus"
)

// NewRollupChallenge initializes a new challenge for rollups
func NewRollupChallenge(challengeID, rollupID, challenger, challengedBlock string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.RollupChallenge {
	return &common.RollupChallenge{
		ChallengeID:     challengeID,
		RollupID:        rollupID,
		Challenger:      challenger,
		ChallengedBlock: challengedBlock,
		Timestamp:       time.Now(),
		IsResolved:      false,
		Ledger:          ledgerInstance,
		Encryption:      encryptionService,
		Consensus:       consensus,
	}
}

// SubmitChallenge allows a challenger to submit a challenge against a rollup's block or state
func (rc *common.RollupChallenge) SubmitChallenge() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.IsResolved {
		return errors.New("challenge already resolved")
	}

	// Encrypt the challenge details before recording
	encryptedChallengedBlock, err := rc.Encryption.EncryptData([]byte(rc.ChallengedBlock), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt challenge data: %v", err)
	}
	rc.ChallengedBlock = string(encryptedChallengedBlock)

	// Log the challenge in the ledger
	err = rc.Ledger.RecordChallengeSubmission(rc.ChallengeID, rc.RollupID, rc.Challenger, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log challenge submission: %v", err)
	}

	fmt.Printf("Challenge %s submitted for rollup %s by challenger %s\n", rc.ChallengeID, rc.RollupID, rc.Challenger)
	return nil
}

// ResolveChallenge resolves the challenge through the consensus mechanism
func (rc *common.RollupChallenge) ResolveChallenge() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.IsResolved {
		return errors.New("challenge already resolved")
	}

	// Use consensus to validate the rollup's state and the challenge
	err := rc.Consensus.ValidateChallenge(rc.ChallengeID, rc.RollupID, rc.ChallengedBlock)
	if err != nil {
		return fmt.Errorf("challenge resolution failed: %v", err)
	}

	rc.IsResolved = true

	// Log the challenge resolution in the ledger
	err = rc.Ledger.RecordChallengeResolution(rc.ChallengeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log challenge resolution: %v", err)
	}

	fmt.Printf("Challenge %s resolved for rollup %s\n", rc.ChallengeID, rc.RollupID)
	return nil
}

// EscalateChallenge escalates the challenge if not resolved through the normal consensus process
func (rc *common.RollupChallenge) EscalateChallenge() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.IsResolved {
		return errors.New("challenge already resolved")
	}

	// Example logic for escalation (this could be to a higher-level protocol or court of arbitration)
	// In practice, this would involve more complex procedures and further validation
	err := rc.Consensus.EscalateChallenge(rc.ChallengeID, rc.RollupID)
	if err != nil {
		return fmt.Errorf("challenge escalation failed: %v", err)
	}

	// Log the escalation event in the ledger
	err = rc.Ledger.RecordChallengeEscalation(rc.ChallengeID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log challenge escalation: %v", err)
	}

	fmt.Printf("Challenge %s escalated for rollup %s\n", rc.ChallengeID, rc.RollupID)
	return nil
}

// RetrieveChallenge retrieves the challenge details by its ID
func (rc *common.RollupChallenge) RetrieveChallenge() (*common.RollupChallenge, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	fmt.Printf("Retrieved challenge %s for rollup %s\n", rc.ChallengeID, rc.RollupID)
	return rc, nil
}
