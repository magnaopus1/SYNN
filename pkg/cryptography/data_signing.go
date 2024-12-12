package cryptography

import (
	"errors"
	"sync"
	"synnergy_network/pkg/common"
)

var signingLock sync.Mutex

// ZKPProve: Generates a zero-knowledge proof for a specific statement
func ZKPProve(statement []byte) ([]byte, error) {
    proof, err := zkp.GenerateProof(statement)
    if err != nil {
        LogSigningOperation("ZKPProve", "Proof generation failed")
        return nil, errors.New("failed to generate zero-knowledge proof")
    }
    LogSigningOperation("ZKPProve", "Zero-knowledge proof generated successfully")
    return proof, nil
}

// ZKPVerify: Verifies a zero-knowledge proof for a given statement
func ZKPVerify(statement, proof []byte) (bool, error) {
    valid, err := zkp.VerifyProof(statement, proof)
    if err != nil || !valid {
        LogSigningOperation("ZKPVerify", "Proof verification failed")
        return false, errors.New("zero-knowledge proof verification failed")
    }
    LogSigningOperation("ZKPVerify", "Zero-knowledge proof verified successfully")
    return true, nil
}

// ZKPRangeProof: Generates a range proof to prove that a value is within a certain range
func ZKPRangeProof(value []byte, min, max int) ([]byte, error) {
    proof, err := zkp.GenerateRangeProof(value, min, max)
    if err != nil {
        LogSigningOperation("ZKPRangeProof", "Range proof generation failed")
        return nil, errors.New("failed to generate range proof")
    }
    LogSigningOperation("ZKPRangeProof", "Range proof generated successfully")
    return proof, nil
}

// ZKPZeroKnowledgeSum: Proves that the sum of values equals a target value without revealing the values
func ZKPZeroKnowledgeSum(values [][]byte, targetSum []byte) ([]byte, error) {
    proof, err := zkp.GenerateSumProof(values, targetSum)
    if err != nil {
        LogSigningOperation("ZKPZeroKnowledgeSum", "Sum proof generation failed")
        return nil, errors.New("failed to generate zero-knowledge sum proof")
    }
    LogSigningOperation("ZKPZeroKnowledgeSum", "Zero-knowledge sum proof generated successfully")
    return proof, nil
}

// BLSMultiSign: Creates a BLS multisignature from multiple signatures
func BLSMultiSign(signatures [][]byte) ([]byte, error) {
    multisig, err := blscrypto.Aggregate(signatures)
    if err != nil {
        LogSigningOperation("BLSMultiSign", "Multisignature creation failed")
        return nil, errors.New("failed to create BLS multisignature")
    }
    LogSigningOperation("BLSMultiSign", "BLS multisignature created successfully")
    return multisig, nil
}

// BLSVerifyMultiSign: Verifies a BLS multisignature against multiple public keys and messages
func BLSVerifyMultiSign(publicKeys [][]byte, messages [][]byte, multisig []byte) (bool, error) {
    valid, err := blscrypto.VerifyAggregate(publicKeys, messages, multisig)
    if err != nil || !valid {
        LogSigningOperation("BLSVerifyMultiSign", "Multisignature verification failed")
        return false, errors.New("failed to verify BLS multisignature")
    }
    LogSigningOperation("BLSVerifyMultiSign", "BLS multisignature verified successfully")
    return true, nil
}

// SignatureAggregate: Aggregates multiple signatures into a single aggregate signature
func SignatureAggregate(signatures [][]byte) ([]byte, error) {
    aggregateSig, err := blscrypto.Aggregate(signatures)
    if err != nil {
        LogSigningOperation("SignatureAggregate", "Signature aggregation failed")
        return nil, errors.New("failed to aggregate signatures")
    }
    LogSigningOperation("SignatureAggregate", "Signatures aggregated successfully")
    return aggregateSig, nil
}

// SignatureSplit: Splits an aggregate signature back into individual signatures
func SignatureSplit(aggregateSig []byte) ([][]byte, error) {
    signatures, err := blscrypto.Split(aggregateSig)
    if err != nil {
        LogSigningOperation("SignatureSplit", "Signature splitting failed")
        return nil, errors.New("failed to split aggregate signature")
    }
    LogSigningOperation("SignatureSplit", "Aggregate signature split successfully")
    return signatures, nil
}

// DataSignWithKey: Signs data with a specified private key
func DataSignWithKey(privateKey []byte, data []byte) ([]byte, error) {
    signature, err := blscrypto.Sign(privateKey, data)
    if err != nil {
        LogSigningOperation("DataSignWithKey", "Data signing failed")
        return nil, errors.New("failed to sign data")
    }
    LogSigningOperation("DataSignWithKey", "Data signed successfully")
    return signature, nil
}

// DataVerifyWithKey: Verifies data against a signature and public key
func DataVerifyWithKey(publicKey []byte, data []byte, signature []byte) (bool, error) {
    valid, err := blscrypto.Verify(publicKey, data, signature)
    if err != nil || !valid {
        LogSigningOperation("DataVerifyWithKey", "Data verification failed")
        return false, errors.New("failed to verify data")
    }
    LogSigningOperation("DataVerifyWithKey", "Data verified successfully")
    return true, nil
}

// DigitalSignatureInit: Initializes a digital signature process
func DigitalSignatureInit(data []byte) ([]byte, error) {
    hash := sha256.Sum256(data)
    LogSigningOperation("DigitalSignatureInit", "Digital signature initialization completed")
    return hash[:], nil
}

// DigitalSignatureFinalize: Finalizes a digital signature with a given private key
func DigitalSignatureFinalize(privateKey []byte, hash []byte) ([]byte, error) {
    signature, err := blscrypto.Sign(privateKey, hash)
    if err != nil {
        LogSigningOperation("DigitalSignatureFinalize", "Signature finalization failed")
        return nil, errors.New("failed to finalize digital signature")
    }
    LogSigningOperation("DigitalSignatureFinalize", "Digital signature finalized successfully")
    return signature, nil
}

// Helper Functions

// LogSigningOperation: Logs signing and verification operations with encryption
func LogSigningOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("SigningOperation", encryptedMessage)
}
