package cryptography

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// Base64Encode: Encodes data to Base64 format
func Base64Encode(data []byte) (string, error) {
    encoded := base64.StdEncoding.EncodeToString(data)
    LogEncodingOperation("Base64Encode", fmt.Sprintf("Data encoded to Base64, length: %d", len(encoded)))
    return encoded, nil
}

// Base64Decode: Decodes Base64 encoded data
func Base64Decode(encoded string) ([]byte, error) {
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        LogEncodingOperation("Base64Decode", "Base64 decoding failed")
        return nil, errors.New("failed to decode Base64 data")
    }
    LogEncodingOperation("Base64Decode", fmt.Sprintf("Data decoded from Base64, length: %d", len(decoded)))
    return decoded, nil
}

// Base32Encode: Encodes data to Base32 format
func Base32Encode(data []byte) (string, error) {
    encoded := base32.StdEncoding.EncodeToString(data)
    LogEncodingOperation("Base32Encode", fmt.Sprintf("Data encoded to Base32, length: %d", len(encoded)))
    return encoded, nil
}

// Base32Decode: Decodes Base32 encoded data
func Base32Decode(encoded string) ([]byte, error) {
    decoded, err := base32.StdEncoding.DecodeString(encoded)
    if err != nil {
        LogEncodingOperation("Base32Decode", "Base32 decoding failed")
        return nil, errors.New("failed to decode Base32 data")
    }
    LogEncodingOperation("Base32Decode", fmt.Sprintf("Data decoded from Base32, length: %d", len(decoded)))
    return decoded, nil
}

// EncodeToHex: Encodes data to hexadecimal format
func EncodeToHex(data []byte) (string, error) {
    encoded := hex.EncodeToString(data)
    LogEncodingOperation("EncodeToHex", fmt.Sprintf("Data encoded to hexadecimal, length: %d", len(encoded)))
    return encoded, nil
}

// DecodeFromHex: Decodes hexadecimal encoded data
func DecodeFromHex(encoded string) ([]byte, error) {
    decoded, err := hex.DecodeString(encoded)
    if err != nil {
        LogEncodingOperation("DecodeFromHex", "Hexadecimal decoding failed")
        return nil, errors.New("failed to decode hexadecimal data")
    }
    LogEncodingOperation("DecodeFromHex", fmt.Sprintf("Data decoded from hexadecimal, length: %d", len(decoded)))
    return decoded, nil
}

// Helper Functions

// LogEncodingOperation: Logs encoding and decoding operations with encryption
func LogEncodingOperation(operation string, details string) error {
    encryptedMessage, err := common.EncryptData([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return ledger.LogEvent("EncodingOperation", encryptedMessage)
}
