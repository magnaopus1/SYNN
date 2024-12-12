package quantum_cryptography

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// ErrorCorrectionEncode: Encodes data with error correction to protect against data corruption
func ErrorCorrectionEncode(data []byte) ([]byte, error) {
    // Apply error correction encoding using a suitable algorithm, e.g., Reed-Solomon or BCH.
    encodedData, err := error_correction.Encode(data)
    if err != nil {
        LogErrorCorrection("ErrorCorrectionEncode", "Encoding failed")
        return nil, errors.New("error correction encoding failed")
    }

    LogErrorCorrection("ErrorCorrectionEncode", fmt.Sprintf("Data encoded with error correction: %d bytes", len(encodedData)))
    return encodedData, nil
}

// ErrorCorrectionDecode: Decodes data with error correction and corrects any detected errors
func ErrorCorrectionDecode(encodedData []byte) ([]byte, error) {
    // Apply error correction decoding, which also detects and corrects errors if possible.
    decodedData, err := error_correction.Decode(encodedData)
    if err != nil {
        LogErrorCorrection("ErrorCorrectionDecode", "Decoding failed")
        return nil, errors.New("error correction decoding failed or data is irrecoverable")
    }

    LogErrorCorrection("ErrorCorrectionDecode", fmt.Sprintf("Data decoded with error correction: %d bytes", len(decodedData)))
    return decodedData, nil
}

// Helper Functions

// LogErrorCorrection: Logs error correction operations with encryption
func LogErrorCorrection(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("ErrorCorrectionOperation", encryptedMessage)
}
