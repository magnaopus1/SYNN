package automations

import (
	"fmt"
	"sync"
	"time"
	"synnergy_network_demo/common"
	"synnergy_network_demo/ledger"
	"synnergy_network_demo/consensus"
	"synnergy_network_demo/encryption"
)

// Constants for staking monitoring and restrictions
const (
	StakingCheckInterval         = 5 * time.Second // Frequency of staking checks
	MaxStakingLimit              = 1000000         // Maximum allowable staking amount (example)
	MinStakingLimit              = 100             // Minimum staking amount allowed
	ExcessiveStakingWarning      = "Excessive staking detected for wallet"
	UnauthorizedStakingActivity  = "Unauthorized staking detected for wallet"
)

// StakingLimitRestrictionAutomation handles restrictions on staking activities within the network
type StakingLimitRestrictionAutomation struct {
	consensusSystem    *consensus.SynnergyConsensus
	ledgerInstance     *ledger.Ledger
	stateMutex         *sync.RWMutex
	stakingActivity    map[string]int // Track staking activity for each wallet
	restrictedWallets  map[string]bool // Track wallets with restricted staking
}

// NewStakingLimitRestrictionAutomation initializes the staking restriction automation
func NewStakingLimitRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *StakingLimitRestrictionAutomation {
	return &StakingLimitRestrictionAutomation{
		consensusSystem:   consensusSystem,
		ledgerInstance:    ledgerInstance,
		stateMutex:        stateMutex,
		stakingActivity:   make(map[string]int),
		restrictedWallets: make(map[string]bool),
	}
}

// StartStakingMonitoring begins continuous monitoring of staking activities across wallets
func (automation *StakingLimitRestrictionAutomation) StartStakingMonitoring() {
	ticker := time.NewTicker(StakingCheckInterval)

	go func() {
		for range ticker.C {
			automation.evaluateStakingActivities()
		}
	}()
}

// evaluateStakingActivities checks staking activities for excessive or unauthorized behavior
func (automation *StakingLimitRestrictionAutomation) evaluateStakingActivities() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	stakingData := automation.consensusSystem.GetStakingData()

	for walletID, stakingAmount := range stakingData {
		if stakingAmount > MaxStakingLimit {
			automation.logExcessiveStaking(walletID, stakingAmount)
			automation.restrictWallet(walletID)
		} else if stakingAmount < MinStakingLimit && stakingAmount > 0 {
			automation.logInsufficientStaking(walletID, stakingAmount)
		} else if automation.consensusSystem.IsUnauthorizedStaking(walletID) {
			automation.logUnauthorizedStaking(walletID)
			automation.restrictWallet(walletID)
		}
	}
}

// logExcessiveStaking logs instances of excessive staking in the ledger
func (automation *StakingLimitRestrictionAutomation) logExcessiveStaking(walletID string, stakingAmount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("excessive-staking-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Excessive Staking",
		Status:    "Warning",
		Details:   fmt.Sprintf("Wallet %s staked %d, exceeding the maximum limit.", walletID, stakingAmount),
	}

	// Encrypt the staking details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log excessive staking:", err)
	} else {
		fmt.Println("Excessive staking logged for wallet:", walletID)
	}
}

// logInsufficientStaking logs cases where staking falls below the minimum threshold
func (automation *StakingLimitRestrictionAutomation) logInsufficientStaking(walletID string, stakingAmount int) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("insufficient-staking-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Insufficient Staking",
		Status:    "Warning",
		Details:   fmt.Sprintf("Wallet %s staked only %d, below the minimum required staking.", walletID, stakingAmount),
	}

	// Encrypt the staking details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log insufficient staking:", err)
	} else {
		fmt.Println("Insufficient staking logged for wallet:", walletID)
	}
}

// logUnauthorizedStaking logs unauthorized staking attempts in the ledger
func (automation *StakingLimitRestrictionAutomation) logUnauthorizedStaking(walletID string) {
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("unauthorized-staking-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Unauthorized Staking",
		Status:    "Critical",
		Details:   fmt.Sprintf("Unauthorized staking activity detected for wallet %s.", walletID),
	}

	// Encrypt the unauthorized staking details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log unauthorized staking:", err)
	} else {
		fmt.Println("Unauthorized staking logged for wallet:", walletID)
	}
}

// restrictWallet restricts staking activity for wallets that breach limits
func (automation *StakingLimitRestrictionAutomation) restrictWallet(walletID string) {
	fmt.Printf("Wallet %s has been restricted due to staking limit violations.\n", walletID)

	// Log the restriction in the ledger
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("wallet-restriction-%s-%d", walletID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "Staking Restriction",
		Status:    "Restricted",
		Details:   fmt.Sprintf("Wallet %s restricted from staking due to exceeding staking limits or unauthorized activity.", walletID),
	}

	// Encrypt the restriction details before logging
	encryptedDetails := automation.encryptData(entry.Details)
	entry.Details = encryptedDetails

	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Println("Failed to log wallet restriction:", err)
	} else {
		fmt.Println("Wallet restriction logged for:", walletID)
	}

	// Update the consensus system to restrict staking activity for the wallet
	automation.consensusSystem.RestrictStaking(walletID)
}

// encryptData encrypts sensitive information before storing in the ledger
func (automation *StakingLimitRestrictionAutomation) encryptData(data string) string {
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting staking data:", err)
		return data
	}
	return string(encryptedData)
}

