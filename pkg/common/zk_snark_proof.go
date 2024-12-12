package common

import (
    "errors"
    "fmt"
    "sync"
    "time"
    "synnergy_network/pkg/ledger"
)

// ZkSNARKProof represents a zk-SNARK proof used within the blockchain network
type ZkSNARKProof struct {
    ProofID        string                        // Unique identifier for the zk-SNARK proof
    ProverID       string                        // ID of the prover who generated the proof
    VerifierID     string                        // ID of the verifier who validates the proof
    ProofData      []byte                        // The actual zk-SNARK proof data
    IsValid        bool                          // Whether the proof has been validated
    VerifiedAt     time.Time                     // Timestamp when the proof was validated
    Ledger         *ledger.Ledger                // Reference to the ledger for logging
    Encryption     *Encryption        // Encryption service for securing proof data
    mu             sync.Mutex                    // Mutex for concurrency control
}

// NewZkSNARKProof initializes a new zk-SNARK proof
func NewZkSNARKProof(proofID string, proverID string, verifierID string, ledgerInstance *ledger.Ledger, encryptionService *Encryption) *ZkSNARKProof {
    return &ZkSNARKProof{
        ProofID:    proofID,
        ProverID:   proverID,
        VerifierID: verifierID,
        Ledger:     ledgerInstance,
        Encryption: encryptionService,
    }
}

// GenerateProof generates a zk-SNARK proof based on input data
func (p *ZkSNARKProof) GenerateProof(inputData []byte) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    // Encrypt input data to generate zk-SNARK proof (in reality, a zk-SNARK proving algorithm is used here)
    // Adding a context label (e.g., "zkSNARK") to the EncryptData function
    encryptedData, err := p.Encryption.EncryptData("zkSNARK", inputData, EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt input data for zk-SNARK proof: %v", err)
    }

    // Store the proof data
    p.ProofData = encryptedData

    // Log proof generation in the ledger (convert time.Time to string using RFC3339 format)
    err = p.Ledger.CryptographyLedger.RecordProofGeneration(p.ProofID, p.ProverID, time.Now().Format(time.RFC3339))
    if err != nil {
        return fmt.Errorf("failed to log proof generation: %v", err)
    }

    fmt.Printf("zk-SNARK proof %s generated by prover %s\n", p.ProofID, p.ProverID)
    return nil
}

// VerifyProof verifies the zk-SNARK proof and marks it as valid or invalid
func (p *ZkSNARKProof) VerifyProof(expectedData []byte) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    // Decrypt the proof data to simulate verification (in reality, zk-SNARK verification is done)
    decryptedData, err := p.Encryption.DecryptData(p.ProofData, EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt zk-SNARK proof: %v", err)
    }

    // Check if decrypted data matches expected data (proof validation)
    if string(decryptedData) == string(expectedData) {
        p.IsValid = true
        p.VerifiedAt = time.Now()

        // Log proof validation in the ledger (only two arguments required)
        err = p.Ledger.CryptographyLedger.RecordProofValidation(p.ProofID, p.VerifierID)
        if err != nil {
            return fmt.Errorf("failed to log proof validation: %v", err)
        }

        fmt.Printf("zk-SNARK proof %s successfully verified by verifier %s\n", p.ProofID, p.VerifierID)
    } else {
        return errors.New("zk-SNARK proof verification failed: data mismatch")
    }

    return nil
}


// RetrieveProof retrieves the proof data for inspection or auditing
func (p *ZkSNARKProof) RetrieveProof() ([]byte, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.ProofData == nil {
        return nil, errors.New("no proof data available")
    }

    fmt.Printf("zk-SNARK proof %s retrieved for auditing\n", p.ProofID)
    return p.ProofData, nil
}

// IsProofValid checks if the proof has been validated
func (p *ZkSNARKProof) IsProofValid() bool {
    p.mu.Lock()
    defer p.mu.Unlock()

    return p.IsValid
}