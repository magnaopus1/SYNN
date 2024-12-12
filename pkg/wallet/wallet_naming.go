package wallet

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"unicode"
)

// NewWalletNaming initializes a new WalletNaming instance.
func NewWalletNaming(ledgerInstance *ledger.Ledger) *WalletNaming {
	return &WalletNaming{
		ledgerInstance: ledgerInstance,
		walletNames:    make(map[string]string),
	}
}


// AssignNameToWallet assigns a human-readable name to a wallet based on its address.
func (wn *WalletNaming) AssignNameToWallet(walletName, walletAddress string) error {
	// Ensure that the name and address are valid.
	if strings.TrimSpace(walletName) == "" || !wn.validateAddress(walletAddress) {
		return errors.New("invalid wallet name or address")
	}

	// Check if the name is already taken.
	if _, exists := wn.walletNames[walletName]; exists {
		return errors.New("wallet name already exists")
	}

	// Securely encrypt the wallet name before storage.
	encryptionInstance := &common.Encryption{} // Create an encryption instance
	encryptedWalletName, err := encryptionInstance.EncryptData("AES", []byte(walletName), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt wallet name: %v", err)
	}

	// Convert encrypted wallet name (which is a byte slice) to a Base64-encoded string
	encryptedWalletNameStr := base64.StdEncoding.EncodeToString(encryptedWalletName)

	// Add the wallet name to the map and update the ledger.
	wn.walletNames[walletName] = walletAddress
	err = wn.ledgerInstance.RecordWalletNaming(walletAddress, encryptedWalletNameStr) // Pass the Base64 string
	if err != nil {
		return fmt.Errorf("failed to store wallet name in ledger: %v", err)
	}

	return nil
}



// RetrieveWalletAddressByName retrieves the wallet address associated with the provided name.
func (wn *WalletNaming) RetrieveWalletAddressByName(walletName string) (string, error) {
	// Decrypt the name to retrieve the original stored name.
	if address, exists := wn.walletNames[walletName]; exists {
		return address, nil
	}

	return "", errors.New("wallet name not found")
}

// IsAlphaNumeric checks if a string contains only alphanumeric characters.
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// ValidateWalletName ensures the name follows proper format.
func (wn *WalletNaming) ValidateWalletName(name string) bool {
	// Ensure the name is alphanumeric, non-empty, and within reasonable length.
	if len(name) > 0 && len(name) <= 32 && IsAlphaNumeric(name) {
		return true
	}
	return false
}

// ValidateWalletAddress checks the validity of a wallet address.
func (wn *WalletNaming) validateAddress(address string) bool {
	// Ensure the address starts with "0x" and has the correct length (42 characters for Ethereum-like addresses).
	if !(strings.HasPrefix(address, "0x") && len(address) == 42) {
		return false
	}

	// Ensure the characters following "0x" are valid hexadecimal characters.
	addressBody := address[2:] // Skip the "0x" prefix
	for _, char := range addressBody {
		if !strings.ContainsRune("0123456789abcdefABCDEF", char) {
			return false
		}
	}

	return true
}


// RemoveWalletName deletes the wallet name association from the system.
func (wn *WalletNaming) RemoveWalletName(walletName string) error {
	if _, exists := wn.walletNames[walletName]; !exists {
		return errors.New("wallet name does not exist")
	}

	delete(wn.walletNames, walletName)
	return nil
}

// EncryptWalletName encrypts a wallet name for secure storage.
func (wn *WalletNaming) EncryptWalletName(walletName string) (string, error) {
	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Encrypt the wallet name using the encryption instance
	encryptedData, err := encryptionInstance.EncryptData("AES", []byte(walletName), common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt wallet name: %v", err)
	}

	return string(encryptedData), nil
}

// DecryptWalletName decrypts a wallet name for use.
func (wn *WalletNaming) DecryptWalletName(encryptedName string) (string, error) {
	// Create an encryption instance
	encryptionInstance := &common.Encryption{}

	// Decrypt the wallet name using the encryption instance
	decryptedData, err := encryptionInstance.DecryptData([]byte(encryptedName), common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt wallet name: %v", err)
	}

	return string(decryptedData), nil
}


// FetchAllWalletNames returns all wallet names currently stored.
func (wn *WalletNaming) FetchAllWalletNames() ([]string, error) {
	if len(wn.walletNames) == 0 {
		return nil, errors.New("no wallet names available")
	}

	names := make([]string, 0, len(wn.walletNames))
	for name := range wn.walletNames {
		names = append(names, name)
	}

	return names, nil
}

// FetchAllWalletAddresses returns all wallet addresses currently stored.
func (wn *WalletNaming) FetchAllWalletAddresses() ([]string, error) {
	if len(wn.walletNames) == 0 {
		return nil, errors.New("no wallet addresses available")
	}

	addresses := make([]string, 0, len(wn.walletNames))
	for _, address := range wn.walletNames {
		addresses = append(addresses, address)
	}

	return addresses, nil
}
