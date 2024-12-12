package wallet

import (
	"crypto/ecdsa"
	"sync"
	"synnergy_network/pkg/ledger"
	"github.com/btcsuite/btcutil/hdkeychain"
)

// MnemonicRecoveryManager manages the recovery process using secondary options.
type MnemonicRecoveryManager struct {
	RecoveryEmail      string
	RecoveryPhoneNumber string
	Syn900Token        string
	IsRecoverySetUp    bool
	ledgerInstance     *ledger.Ledger
	mutex              sync.Mutex
}

// HDWallet manages the Hierarchical Deterministic wallet using BIP32 standard.
type HDWallet struct {
	masterPrivateKey *hdkeychain.ExtendedKey // BIP32 master private key
	ledgerInstance   *ledger.Ledger
	mutex            sync.Mutex
}

// IDTokenWalletRegistrationService manages the registration of the Syn900 token to a wallet.
type IDTokenWalletRegistrationService struct {
	ledgerInstance *ledger.Ledger
	mutex          sync.Mutex
}

// OffchainWallet represents an off-chain wallet that interacts with the ledger but handles transactions off-chain.
type OffchainWallet struct {
	WalletID         string
	PrivateKey       string
	PublicKey        string
	OffchainBalances map[string]float64
	ledgerInstance   *ledger.Ledger
	mutex            sync.Mutex
}

// WalletBackupService manages wallet backups and restoration functionality.
type WalletBackupService struct {
	walletID       string
	walletFilePath string
	ledgerInstance *ledger.Ledger
	mutex          sync.Mutex
}

// WalletBackupData represents the structure of data that gets backed up.
type WalletBackupData struct {
	WalletID   string            `json:"wallet_id"`
	Keys       map[string]string `json:"keys"` // Could be public/private key pair, mnemonic, etc.
	Balances   map[string]float64 `json:"balances"`
}

// WalletBalanceService manages the balances for a given wallet, including fetching and updating balances.
type WalletBalanceService struct {
	walletID       string
	ledgerInstance *ledger.Ledger
	mutex          sync.Mutex
	balances       map[string]float64 // map to store balances of multiple currencies/tokens
}


// WalletDisplayService handles the display and visualization of wallet data.
type WalletDisplayService struct {
	walletID       string
	ledgerInstance *ledger.Ledger
	mutex          sync.Mutex
}



// WalletNaming manages wallet names and provides mapping between human-readable names and wallet addresses.
type WalletNaming struct {
	ledgerInstance *ledger.Ledger      // Ledger to store wallet name mappings and track activities.
	walletNames    map[string]string   // Map storing wallet names to addresses.
	mutex          sync.Mutex
}


// WalletRecovery provides functionality to recover wallets using mnemonic phrases or private keys.
type WalletRecovery struct {
	PrivateKey *ecdsa.PrivateKey // The recovered private key
}

// WalletSigner handles signing and verifying data using the wallet's private key.
type WalletSigner struct {
	PrivateKey *ecdsa.PrivateKey // The private key used for signing
	Ledger     *ledger.Ledger    // Ledger instance for transaction and contract validation
	mutex      sync.Mutex
}



