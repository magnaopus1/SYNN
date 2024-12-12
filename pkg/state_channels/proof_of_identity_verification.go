package state_channels

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewIdentityVerificationChannel initializes a new Identity Verification Channel (IVC)
func NewIdentityVerificationChannel(channelID string, participants []string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.IdentityVerificationChannel {
	return &common.IdentityVerificationChannel{
		ChannelID:      channelID,
		Participants:   participants,
		IdentityProofs: make(map[string]*common.IdentityProof),
		IsOpen:         true,
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// OpenChannel opens the identity verification channel
func (ivc *common.IdentityVerificationChannel) OpenChannel() error {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if ivc.IsOpen {
		return errors.New("identity verification channel is already open")
	}

	ivc.IsOpen = true

	// Log the channel opening in the ledger
	err := ivc.Ledger.RecordChannelOpening(ivc.ChannelID, ivc.Participants, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel opening: %v", err)
	}

	fmt.Printf("Identity Verification Channel %s opened with participants: %v\n", ivc.ChannelID, ivc.Participants)
	return nil
}

// CloseChannel closes the identity verification channel after finalizing all proofs
func (ivc *common.IdentityVerificationChannel) CloseChannel() error {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if !ivc.IsOpen {
		return errors.New("identity verification channel is already closed")
	}

	// Ensure all identity proofs are verified before closing
	for _, proof := range ivc.IdentityProofs {
		if !proof.Verified {
			return fmt.Errorf("identity proof for participant %s is not verified", proof.Participant)
		}
	}

	ivc.IsOpen = false

	// Log the channel closure in the ledger
	err := ivc.Ledger.RecordChannelClosure(ivc.ChannelID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log channel closure: %v", err)
	}

	fmt.Printf("Identity Verification Channel %s closed\n", ivc.ChannelID)
	return nil
}

// SubmitIdentityProof allows a participant to submit their identity proof for verification
func (ivc *common.IdentityVerificationChannel) SubmitIdentityProof(participant string, proofDocument string) error {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if !ivc.IsOpen {
		return errors.New("identity verification channel is closed")
	}

	if _, exists := ivc.IdentityProofs[participant]; exists {
		return errors.New("identity proof already submitted for this participant")
	}

	// Hash the identity proof document
	proofHash := sha256.Sum256([]byte(proofDocument))

	// Encrypt the proof hash before storing it
	encryptedProofHash, err := ivc.Encryption.EncryptData(proofHash[:], common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt identity proof: %v", err)
	}

	// Create and store the identity proof
	identityProof := &common.IdentityProof{
		Participant: participant,
		ProofHash:   fmt.Sprintf("%x", encryptedProofHash),
		Timestamp:   time.Now(),
		Verified:    false,
	}

	ivc.IdentityProofs[participant] = identityProof

	// Log the proof submission in the ledger
	err = ivc.Ledger.RecordIdentityProofSubmission(ivc.ChannelID, participant, identityProof.ProofHash, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log identity proof submission: %v", err)
	}

	fmt.Printf("Identity proof submitted by participant %s in channel %s\n", participant, ivc.ChannelID)
	return nil
}

// VerifyIdentityProof verifies a participant's identity proof
func (ivc *common.IdentityVerificationChannel) VerifyIdentityProof(participant string, proofDocument string) error {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if !ivc.IsOpen {
		return errors.New("identity verification channel is closed")
	}

	proof, exists := ivc.IdentityProofs[participant]
	if !exists {
		return fmt.Errorf("no identity proof found for participant %s", participant)
	}

	// Hash the provided proof document to compare with the stored proof
	proofHash := sha256.Sum256([]byte(proofDocument))

	// Verify the proof hash
	if fmt.Sprintf("%x", proofHash) != proof.ProofHash {
		return errors.New("identity proof does not match the stored hash")
	}

	// Mark the proof as verified
	proof.Verified = true

	// Log the proof verification in the ledger
	err := ivc.Ledger.RecordIdentityProofVerification(ivc.ChannelID, participant, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log identity proof verification: %v", err)
	}

	fmt.Printf("Identity proof for participant %s verified in channel %s\n", participant, ivc.ChannelID)
	return nil
}

// RetrieveIdentityProof retrieves a participant's identity proof from the channel
func (ivc *common.IdentityVerificationChannel) RetrieveIdentityProof(participant string) (*common.IdentityProof, error) {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	proof, exists := ivc.IdentityProofs[participant]
	if !exists {
		return nil, fmt.Errorf("no identity proof found for participant %s", participant)
	}

	fmt.Printf("Retrieved identity proof for participant %s in channel %s\n", participant, ivc.ChannelID)
	return proof, nil
}

// RetrieveAllProofs retrieves all identity proofs from the channel
func (ivc *common.IdentityVerificationChannel) RetrieveAllProofs() map[string]*common.IdentityProof {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	return ivc.IdentityProofs
}

// UpdateState securely updates the internal state of the channel
func (ivc *common.IdentityVerificationChannel) UpdateState(key string, value interface{}) error {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if !ivc.IsOpen {
		return errors.New("identity verification channel is closed")
	}

	// Update the channel state
	ivc.State[key] = value

	// Log the state update in the ledger
	err := ivc.Ledger.RecordStateUpdate(ivc.ChannelID, key, value, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log state update: %v", err)
	}

	fmt.Printf("State of channel %s updated: %s = %v\n", ivc.ChannelID, key, value)
	return nil
}

// RetrieveState retrieves the current state of the channel
func (ivc *common.IdentityVerificationChannel) RetrieveState(key string) (interface{}, error) {
	ivc.mu.Lock()
	defer ivc.mu.Unlock()

	if !ivc.IsOpen {
		return nil, errors.New("identity verification channel is closed")
	}

	value, exists := ivc.State[key]
	if !exists {
		return nil, fmt.Errorf("state key %s not found", key)
	}

	fmt.Printf("Retrieved state from channel %s: %s = %v\n", ivc.ChannelID, key, value)
	return value, nil
}
