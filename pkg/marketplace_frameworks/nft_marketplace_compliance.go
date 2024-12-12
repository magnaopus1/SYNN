package marketplace

import (
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

func denyNFTTrade(tradeID string, reason string, ledgerInstance *Ledger) error {
    encryptedReason, err := encryption.EncryptData(reason)
    if err != nil {
        return fmt.Errorf("failed to encrypt denial reason: %v", err)
    }
    denial := NFTTradeDenial{
        TradeID:         tradeID,
        Reason:          reason,
        EncryptedReason: encryptedReason,
        DeniedAt:        time.Now(),
    }
    return ledgerInstance.recordNFTTradeDenial(denial)
}

func trackTradeCompletion(tradeID string, completionStatus bool, ledgerInstance *Ledger) error {
    trade := NFTTrade{
        TradeID:     tradeID,
        Completed:   completionStatus,
        CompletedAt: time.Now(),
    }
    return ledgerInstance.recordNFTTradeCompletion(trade)
}

func enableNFTExchange(ledgerInstance *Ledger) error {
    return ledgerInstance.setNFTExchangeStatus(true)
}

func setNFTExchangeRate(nftID string, rate float64, ledgerInstance *Ledger) error {
    exchangeRate := NFTExchangeRate{
        NFTID: nftID,
        Rate:  rate,
        SetAt: time.Now(),
    }
    return ledgerInstance.recordNFTExchangeRate(exchangeRate)
}

func logNFTExchangeTransaction(transaction ExchangeTransaction, ledgerInstance *Ledger) error {
    transaction.Timestamp = time.Now()
    return ledgerInstance.recordNFTExchangeTransaction(transaction)
}

func generateNFTExchangeReport(from, to time.Time, ledgerInstance *Ledger) (NFTExchangeReport, error) {
    transactions, err := ledgerInstance.getNFTExchangeTransactions(from, to)
    if err != nil {
        return NFTExchangeReport{}, fmt.Errorf("failed to retrieve NFT exchange transactions: %v", err)
    }
    return NFTExchangeReport{
        Transactions: transactions,
        Period:       fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
        GeneratedAt:  time.Now(),
    }, nil
}

func setNFTMintingLimit(nftID string, limit int, ledgerInstance *Ledger) error {
    mintingLimit := NFTMintingLimit{
        NFTID: nftID,
        Limit: limit,
        SetAt: time.Now(),
    }
    return ledgerInstance.recordNFTMintingLimit(mintingLimit)
}

func trackNFTMintingHistory(nftID string, amount int, ledgerInstance *Ledger) error {
    minting := NFTMintingEvent{
        NFTID:    nftID,
        Amount:   amount,
        MintedAt: time.Now(),
    }
    return ledgerInstance.recordNFTMintingEvent(minting)
}

func generateMintingReport(nftID string, from, to time.Time, ledgerInstance *Ledger) (NFTMintingReport, error) {
    events, err := ledgerInstance.getNFTMintingHistory(nftID, from, to)
    if err != nil {
        return NFTMintingReport{}, fmt.Errorf("failed to retrieve minting history: %v", err)
    }
    return NFTMintingReport{
        NFTID:  nftID,
        Events: events,
        Period: fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}


// ApproveMintingAuthorization grants authorization for NFT minting.
func ApproveMintingAuthorization(requestID string) error {
    ledgerInstance := getLedgerInstance()
    return ledgerInstance.UpdateMintingAuthorizationStatus(requestID, "Approved")
}

// DenyMintingAuthorization denies an NFT minting request.
func DenyMintingAuthorization(requestID string, reason string) error {
    ledgerInstance := getLedgerInstance()
    encryptedReason, err := encryption.EncryptData(reason)
    if err != nil {
        return fmt.Errorf("failed to encrypt denial reason: %v", err)
    }
    authorization := MintingAuthorization{
        RequestID:       requestID,
        Status:          "Denied",
        Reason:          reason,
        EncryptedReason: encryptedReason,
        UpdatedAt:       time.Now(),
    }
    return ledgerInstance.RecordMintingAuthorization(authorization)
}

// EnableNFTCustomization enables customization options for an NFT.
func EnableNFTCustomization(nftID string) error {
    ledgerInstance := getLedgerInstance()
    return ledgerInstance.UpdateNFTCustomizationStatus(nftID, true)
}

// SetCustomizationOptions sets customization options for an NFT.
func SetCustomizationOptions(nftID string, options NFTCustomizationOptions) error {
    ledgerInstance := getLedgerInstance()
    options.NFTID = nftID
    options.SetAt = time.Now()
    return ledgerInstance.RecordNFTCustomizationOptions(options)
}

// TrackCustomizationHistory logs each customization event for an NFT.
func TrackCustomizationHistory(nftID string, customization NFTCustomization) error {
    ledgerInstance := getLedgerInstance()
    customization.NFTID = nftID
    customization.Timestamp = time.Now()
    return ledgerInstance.RecordNFTCustomizationHistory(customization)
}

// LogCustomizationEvent logs a specific event related to NFT customization.
func LogCustomizationEvent(nftID string, eventDescription string) error {
    ledgerInstance := getLedgerInstance()
    encryptedDescription, err := encryption.EncryptData(eventDescription)
    if err != nil {
        return fmt.Errorf("failed to encrypt customization event description: %v", err)
    }
    event := CustomizationEvent{
        NFTID:               nftID,
        Description:         eventDescription,
        EncryptedDescription: encryptedDescription,
        LoggedAt:            time.Now(),
    }
    return ledgerInstance.RecordCustomizationEvent(event)
}

// EnableNFTStakeRewards enables staking rewards for NFT holders.
func EnableNFTStakeRewards(nftID string) error {
    ledgerInstance := getLedgerInstance()
    return ledgerInstance.SetNFTStakeRewardsStatus(nftID, true)
}

// CalculateStakeYieldForNFT calculates the staking yield for a specific NFT.
func CalculateStakeYieldForNFT(nftID string, stakedAmount float64) (float64, error) {
    ledgerInstance := getLedgerInstance()
    yieldRate, err := ledgerInstance.GetNFTYieldRate(nftID)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve NFT yield rate: %v", err)
    }
    return stakedAmount * yieldRate, nil
}

// DistributeStakeRewards distributes staking rewards to NFT holders.
func DistributeStakeRewards(nftID string, rewards []StakeReward) error {
    ledgerInstance := getLedgerInstance()
    for _, reward := range rewards {
        reward.NFTID = nftID
        reward.DistributedAt = time.Now()
        if err := ledgerInstance.RecordStakeRewardDistribution(reward); err != nil {
            return fmt.Errorf("failed to distribute reward for %s: %v", reward.HolderID, err)
        }
    }
    return nil
}

// GenerateStakingReport generates a staking report for an NFT within a specific period.
func GenerateStakingReport(nftID string, from, to time.Time) (StakeReport, error) {
    ledgerInstance := getLedgerInstance()
    rewards, err := ledgerInstance.GetStakeRewards(nftID, from, to)
    if err != nil {
        return StakeReport{}, fmt.Errorf("failed to retrieve staking rewards: %v", err)
    }
    return StakeReport{
        NFTID:   nftID,
        Rewards: rewards,
        Period:  fmt.Sprintf("%s - %s", from.Format(time.RFC3339), to.Format(time.RFC3339)),
    }, nil
}
