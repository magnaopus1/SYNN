package wallet

import (
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

// WalletNotificationService handles sending notifications and alerts to wallet users.
type WalletNotificationService struct {
	ledgerInstance  *ledger.Ledger
	networkManager  *network.NetworkManager // Add the network manager to the struct
	alerts          map[string]string       // Store alerts related to transactions or blocks
	mutex           sync.Mutex
}

// NewWalletNotificationService initializes the notification service for wallets.
func NewWalletNotificationService(ledgerInstance *ledger.Ledger, networkManager *network.NetworkManager) *WalletNotificationService {
	return &WalletNotificationService{
		ledgerInstance: ledgerInstance,
		networkManager: networkManager,  // Initialize the network manager
		alerts:         make(map[string]string),
	}
}



// SendTransactionAlert sends an alert to the wallet owner for a specific transaction.
func (wns *WalletNotificationService) SendTransactionAlert(transactionID, walletAddress string) error {
	// Step 1: Retrieve transaction status from the ledger
	status, err := wns.ledgerInstance.GetTransactionStatus(transactionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction status for alert: %v", err)
	}

	// Step 2: Create the alert message
	alertMessage := fmt.Sprintf("Transaction %s for wallet %s is now %s.", transactionID, walletAddress, status)

	// Step 3: Encrypt the alert message (create an encryption instance in the function)
	encryptionInstance := &common.Encryption{} // Create the encryption instance
	encryptedAlert, err := encryptionInstance.EncryptData("AES", []byte(alertMessage), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction alert: %v", err)
	}

	// Step 4: Send the encrypted alert via the network using the existing `SendEncryptedMessage` function
	err = wns.networkManager.SendEncryptedMessage(walletAddress, string(encryptedAlert))
	if err != nil {
		return fmt.Errorf("failed to send transaction alert: %v", err)
	}

	// Step 5: Log the alert in the wallet alert history
	wns.alerts[transactionID] = string(encryptedAlert)
	return nil
}



// SendBlockConfirmationAlert sends an alert to wallet users when a block has been confirmed.
func (wns *WalletNotificationService) SendBlockConfirmationAlert(blockID string, walletAddresses []string) error {
	// Step 1: Retrieve the block confirmation status from the ledger
	status, err := wns.ledgerInstance.GetBlockStatus(blockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve block confirmation status: %v", err)
	}

	// Step 2: Create the alert message
	alertMessage := fmt.Sprintf("Block %s has been confirmed with status: %s.", blockID, status)

	// Step 3: Encrypt the alert message (using encryption instance)
	encryptionInstance := &common.Encryption{} // Create the encryption instance
	encryptedAlert, err := encryptionInstance.EncryptData("AES", []byte(alertMessage), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block confirmation alert: %v", err)
	}

	// Step 4: Send the encrypted alert to each wallet address using SendEncryptedMessage
	for _, walletAddress := range walletAddresses {
		err := wns.networkManager.SendEncryptedMessage(walletAddress, string(encryptedAlert))
		if err != nil {
			return fmt.Errorf("failed to send block confirmation alert to wallet %s: %v", walletAddress, err)
		}
	}

	// Step 5: Log the alert in the notification service
	wns.alerts[blockID] = string(encryptedAlert)
	return nil
}


// SendBalanceUpdateAlert notifies the user when their wallet balance is updated.
func (wns *WalletNotificationService) SendBalanceUpdateAlert(walletAddress string, newBalance float64) error {
	// Step 1: Create the balance update alert message
	alertMessage := fmt.Sprintf("Your wallet %s balance has been updated to: %.2f.", walletAddress, newBalance)

	// Step 2: Encrypt the alert message (using encryption instance)
	encryptionInstance := &common.Encryption{} // Create the encryption instance
	encryptedAlert, err := encryptionInstance.EncryptData("AES", []byte(alertMessage), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt balance update alert: %v", err)
	}

	// Step 3: Send the encrypted alert to the wallet address using SendEncryptedMessage
	err = wns.networkManager.SendEncryptedMessage(walletAddress, string(encryptedAlert))
	if err != nil {
		return fmt.Errorf("failed to send balance update alert: %v", err)
	}

	// Step 4: Log the alert in the notification service
	alertID := fmt.Sprintf("%s-balance-update-%d", walletAddress, time.Now().UnixNano())
	wns.alerts[alertID] = string(encryptedAlert)

	return nil
}


// GetAlerts retrieves the list of alerts for the wallet.
func (wns *WalletNotificationService) GetAlerts() map[string]string {
	// Return a copy of the alerts to avoid direct modification
	alertCopy := make(map[string]string)
	for key, val := range wns.alerts {
		alertCopy[key] = val
	}
	return alertCopy
}

// SendSubBlockAlert sends a notification when a sub-block is confirmed.
func (wns *WalletNotificationService) SendSubBlockAlert(subBlockID, walletAddress string) error {
	// Step 1: Retrieve the sub-block confirmation status from the ledger
	status, err := wns.ledgerInstance.GetSubBlockStatus(subBlockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve sub-block confirmation status: %v", err)
	}

	// Step 2: Create the sub-block confirmation alert message
	alertMessage := fmt.Sprintf("Sub-block %s is confirmed for wallet %s with status: %s.", subBlockID, walletAddress, status)

	// Step 3: Encrypt the alert message (using encryption instance)
	encryptionInstance := &common.Encryption{} // Create the encryption instance
	encryptedAlert, err := encryptionInstance.EncryptData("AES", []byte(alertMessage), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt sub-block alert: %v", err)
	}

	// Step 4: Send the encrypted alert to the wallet address using SendEncryptedMessage
	err = wns.networkManager.SendEncryptedMessage(walletAddress, string(encryptedAlert))
	if err != nil {
		return fmt.Errorf("failed to send sub-block alert: %v", err)
	}

	// Step 5: Log the alert
	wns.alerts[subBlockID] = string(encryptedAlert)

	return nil
}

