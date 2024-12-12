package wallet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewWalletBackupService initializes the WalletBackupService.
func NewWalletBackupService(walletID, walletFilePath string, ledgerInstance *ledger.Ledger) *WalletBackupService {
	return &WalletBackupService{
		walletID:       walletID,
		walletFilePath: walletFilePath,
		ledgerInstance: ledgerInstance,
	}
}

// CreateBackup generates a wallet backup, encrypts it, and saves it to a file.
func (wbs *WalletBackupService) CreateBackup() error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Retrieve balances and keys (if applicable) for the wallet
	balances, err := wbs.fetchBalancesForBackup()
	if err != nil {
		return fmt.Errorf("failed to fetch balances for backup: %v", err)
	}
	keys, err := wbs.fetchKeysForBackup()
	if err != nil {
		return fmt.Errorf("failed to fetch keys for backup: %v", err)
	}

	// Create the backup data structure
	backupData := WalletBackupData{
		WalletID: wbs.walletID,
		Keys:     keys,
		Balances: balances,
	}

	// Marshal the backup data to JSON
	backupJSON, err := json.Marshal(backupData)
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %v", err)
	}

	// Create the encryption instance
	encryptionInstance := &common.Encryption{}

	// Encrypt the backup data (convert JSON to byte slice)
	encryptedBackup, err := encryptionInstance.EncryptData("AES", backupJSON, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt wallet backup: %v", err)
	}

	// Save the encrypted backup to a file
	err = ioutil.WriteFile(wbs.walletFilePath, encryptedBackup, 0644)
	if err != nil {
		return fmt.Errorf("failed to write encrypted backup to file: %v", err)
	}

	fmt.Printf("Wallet backup for walletID %s created successfully.\n", wbs.walletID)
	return nil
}


// RestoreBackup restores the wallet from an encrypted backup file.
func (wbs *WalletBackupService) RestoreBackup() error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	// Read the encrypted backup from the file
	encryptedBackup, err := ioutil.ReadFile(wbs.walletFilePath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %v", err)
	}

	// Create the encryption instance
	encryptionInstance := &common.Encryption{}

	// Decrypt the backup data
	decryptedBackup, err := encryptionInstance.DecryptData(encryptedBackup, common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt wallet backup: %v", err)
	}

	// Unmarshal the decrypted JSON data into WalletBackupData
	var backupData WalletBackupData
	err = json.Unmarshal(decryptedBackup, &backupData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal decrypted backup data: %v", err)
	}

	// Restore keys and balances from the backup data
	err = wbs.restoreBalances(backupData.Balances)
	if err != nil {
		return fmt.Errorf("failed to restore balances from backup: %v", err)
	}
	err = wbs.restoreKeys(backupData.Keys)
	if err != nil {
		return fmt.Errorf("failed to restore keys from backup: %v", err)
	}

	fmt.Printf("Wallet backup for walletID %s restored successfully.\n", wbs.walletID)
	return nil
}


// fetchBalancesForBackup retrieves the current balances for backup.
func (wbs *WalletBackupService) fetchBalancesForBackup() (map[string]float64, error) {
    // Retrieve all balances from the ledger (assuming they are stored as float64 values).
    balances, err := wbs.ledgerInstance.GetAllBalances()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve balances from ledger: %v", err)
    }

    // Create an encryption instance (if needed).
    encryptionInstance := &common.Encryption{}

    // Decrypt all balances for backup (if encryption is applied to balance strings or data).
    decryptedBalances := make(map[string]float64)
    for currency, encryptedBalance := range balances {


        // Convert float64 to a string (optional if you want to encrypt/decrypt as a string)
        balanceStr := fmt.Sprintf("%.2f", encryptedBalance)

        // If you wish to encrypt/decrypt, do so on the string or byte version:
        decryptedBalanceBytes, err := encryptionInstance.DecryptData([]byte(balanceStr), common.EncryptionKey)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt balance for currency %s: %v", currency, err)
        }

        // Convert the decrypted string back to float64
        var balance float64
        _, err = fmt.Sscanf(string(decryptedBalanceBytes), "%f", &balance)
        if err != nil {
            return nil, fmt.Errorf("failed to parse decrypted balance for currency %s: %v", currency, err)
        }

        // Store the decrypted balance in the map
        decryptedBalances[currency] = balance
    }

    return decryptedBalances, nil
}




// fetchKeysForBackup retrieves the private and public keys or mnemonics for the wallet.
func (wbs *WalletBackupService) fetchKeysForBackup() (map[string]string, error) {
	// Retrieve the keys using the ledger instance
	privateKey, publicKey, err := wbs.ledgerInstance.GetWalletKeys(wbs.walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet keys: %v", err)
	}

	// Decrypt the private and public keys
	encryptionInstance := &common.Encryption{}
	decryptedPrivateKey, err := encryptionInstance.DecryptData([]byte(privateKey), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	decryptedPublicKey, err := encryptionInstance.DecryptData([]byte(publicKey), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt public key: %v", err)
	}

	// Store the decrypted keys in a map to return
	keys := map[string]string{
		"privateKey": string(decryptedPrivateKey),
		"publicKey":  string(decryptedPublicKey),
	}

	return keys, nil
}

// restoreBalances restores balances to the ledger after decrypting from backup.
func (wbs *WalletBackupService) restoreBalances(balances map[string]float64) error {
	for currency, balance := range balances {
		// Convert balance to uint64 (assuming the ledger stores balances as uint64)
		balanceUint := uint64(balance)

		// Update the balance in the ledger without encryption, as the method expects uint64
		err := wbs.ledgerInstance.UpdateBalance(wbs.walletID, balanceUint)
		if err != nil {
			return fmt.Errorf("failed to update balance in ledger for %s: %v", currency, err)
		}
	}

	return nil
}



// restoreKeys restores the private and public keys or mnemonics from the backup.
func (wbs *WalletBackupService) restoreKeys(keys map[string]string) error {
	encryptionInstance := &common.Encryption{} // Create encryption instance

	for keyType, key := range keys {
		// Convert key to a byte slice before encryption
		keyBytes := []byte(key)

		// Encrypt the key before storing it back in the ledger
		encryptedKey, err := encryptionInstance.EncryptData("AES", keyBytes, common.EncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt key for %s: %v", keyType, err)
		}

		// Convert encrypted key to a base64 string or handle it as needed
		encryptedKeyStr := base64.StdEncoding.EncodeToString(encryptedKey)

		// Store the encrypted key in the ledger
		// Update: Remove the third argument and only pass walletID and keyType as required by StoreWalletKey
		err = wbs.ledgerInstance.StoreWalletKey(keyType, encryptedKeyStr)
		if err != nil {
			return fmt.Errorf("failed to store key in ledger for %s: %v", keyType, err)
		}
	}

	return nil
}



// DeleteBackup deletes the backup file from the file system.
func (wbs *WalletBackupService) DeleteBackup() error {
	wbs.mutex.Lock()
	defer wbs.mutex.Unlock()

	err := os.Remove(wbs.walletFilePath)
	if err != nil {
		return fmt.Errorf("failed to delete backup file: %v", err)
	}

	fmt.Printf("Wallet backup for walletID %s deleted successfully.\n", wbs.walletID)
	return nil
}
