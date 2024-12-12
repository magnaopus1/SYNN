package cryptography

import (
	"errors"
	"sync"
)

var signatureLock sync.Mutex

// BLSAggregateSign: Aggregates individual BLS signatures into a single aggregate signature
func BLSAggregateSign(privateKeys [][]byte, messages [][]byte) ([]byte, error) {
    signatureLock.Lock()
    defer signatureLock.Unlock()

    if len(privateKeys) == 0 || len(messages) == 0 || len(privateKeys) != len(messages) {
        LogSignatureOperation("BLSAggregateSign", "Mismatch between private keys and messages count")
        return nil, errors.New("invalid input: number of private keys must match number of messages")
    }

    // Create individual signatures for each message
    individualSigs := make([][]byte, len(privateKeys))
    for i := 0; i < len(privateKeys); i++ {
        sig, err := blscrypto.Sign(privateKeys[i], messages[i])
        if err != nil {
            LogSignatureOperation("BLSAggregateSign", fmt.Sprintf("Failed to sign message %d", i))
            return nil, errors.New("failed to sign message")
        }
        individualSigs[i] = sig
    }

    // Aggregate the individual signatures into a single signature
    aggregateSignature, err := blscrypto.Aggregate(individualSigs)
    if err != nil {
        LogSignatureOperation("BLSAggregateSign", "Signature aggregation failed")
        return nil, errors.New("signature aggregation failed")
    }

    LogSignatureOperation("BLSAggregateSign", "BLS signature aggregation completed successfully")
    return aggregateSignature, nil
}

// BLSAggregateVerify: Verifies an aggregated BLS signature against a set of public keys and messages
func BLSAggregateVerify(publicKeys [][]byte, messages [][]byte, aggregateSignature []byte) (bool, error) {
    signatureLock.Lock()
    defer signatureLock.Unlock()

    if len(publicKeys) == 0 || len(messages) == 0 || len(publicKeys) != len(messages) {
        LogSignatureOperation("BLSAggregateVerify", "Mismatch between public keys and messages count")
        return false, errors.New("invalid input: number of public keys must match number of messages")
    }

    // Verify the aggregated signature using all public keys and their respective messages
    valid, err := blscrypto.VerifyAggregate(publicKeys, messages, aggregateSignature)
    if err != nil || !valid {
        LogSignatureOperation("BLSAggregateVerify", "Aggregate signature verification failed")
        return false, errors.New("aggregate signature verification failed")
    }

    LogSignatureOperation("BLSAggregateVerify", "Aggregate signature verified successfully")
    return true, nil
}

