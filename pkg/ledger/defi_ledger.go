package ledger

import (
	"errors"
	"fmt"
	"time"
)

// RecordPolicyCreation records the creation of an insurance policy.
func (l *DeFiLedger) RecordPolicyCreation(policyID, holder string, insuredAmount, premium float64, duration time.Duration) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.InsurancePolicies[policyID]; exists {
		return errors.New("policy already exists")
	}

	// Create a new insurance policy
	policy := InsurancePolicy{
		PolicyID:       policyID,
		Holder:         holder,
		InsuredAmount:  insuredAmount,
		Premium:        premium,
		PolicyDuration: duration,
		StartDate:      time.Now(),
		ExpiryDate:     time.Now().Add(duration),
		Status:         "Active",
	}

	l.InsurancePolicies[policyID] = policy
	fmt.Printf("Insurance policy %s created for holder %s\n", policyID, holder)
	return nil
}

// RecordClaimSubmission records a claim submission for an insurance policy.
func (l *DeFiLedger) RecordClaimSubmission(policyID string, claimAmount float64) error {
	l.Lock()
	defer l.Unlock()

	if policy, exists := l.InsurancePolicies[policyID]; exists {
		if policy.Status != "Active" {
			return errors.New("policy is not active")
		}

		// Create a new claim for the policy
		claimID := generateClaimID()
		claim := InsuranceClaim{
			ClaimID:     claimID,
			PolicyID:    policyID,
			ClaimAmount: claimAmount,
			ClaimDate:   time.Now(),
			ClaimStatus: "Pending",
		}

		l.InsuranceClaims[claimID] = claim
		fmt.Printf("Claim submitted for policy %s with claim ID %s\n", policyID, claimID)
		return nil
	}
	return errors.New("policy does not exist")
}

// RecordClaimApproval records the approval of a claim.
func (l *DeFiLedger) RecordClaimApproval(claimID string) error {
	l.Lock()
	defer l.Unlock()

	if claim, exists := l.InsuranceClaims[claimID]; exists {
		claim.ClaimStatus = "Approved"
		l.InsuranceClaims[claimID] = claim
		fmt.Printf("Claim %s approved\n", claimID)
		return nil
	}
	return errors.New("claim does not exist")
}

// RecordClaimRejection records the rejection of a claim.
func (l *DeFiLedger) RecordClaimRejection(claimID string) error {
	l.Lock()
	defer l.Unlock()

	if claim, exists := l.InsuranceClaims[claimID]; exists {
		claim.ClaimStatus = "Rejected"
		l.InsuranceClaims[claimID] = claim
		fmt.Printf("Claim %s rejected\n", claimID)
		return nil
	}
	return errors.New("claim does not exist")
}

// RecordLiquidityPoolCreation records the creation of a liquidity pool.
func (l *DeFiLedger) RecordLiquidityPoolCreation(poolID string, totalLiquidity, rewardRate float64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.LiquidityPools[poolID]; exists {
		return errors.New("liquidity pool already exists")
	}

	// Create a new liquidity pool
	pool := LiquidityPool{
		PoolID:           poolID,
		TotalLiquidity:   totalLiquidity,
		AvailableLiquidity: totalLiquidity,
		RewardRate:       rewardRate,
		CreatedAt:        time.Now(),
		Status:           "Active",
	}

	l.LiquidityPools[poolID] = pool
	fmt.Printf("Liquidity pool %s created with total liquidity %.2f\n", poolID, totalLiquidity)
	return nil
}

// RecordAssetPoolCreation records the creation of an asset pool.
func (l *DeFiLedger) RecordAssetPoolCreation(poolID, assetType string, totalAssets, rewardRate float64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.AssetPools[poolID]; exists {
		return errors.New("asset pool already exists")
	}

	// Create a new asset pool
	pool := AssetPool{
		PoolID:       poolID,
		TotalAssets:  totalAssets,
		AssetType:    assetType,
		RewardRate:   rewardRate,
		CreatedAt:    time.Now(),
		Status:       "Active",
	}

	l.AssetPools[poolID] = pool
	fmt.Printf("Asset pool %s created with total assets %.2f\n", poolID, totalAssets)
	return nil
}

// RecordYieldFarmingInitialization records the initialization of a yield farming pool.
func (l *DeFiLedger) RecordYieldFarmingInitialization(poolID string, tokenPair string, totalLiquidity, rewardRate float64) error {
	l.Lock()
	defer l.Unlock()

	// Ensure YieldFarmingPools exist in the state
	if _, exists := l.YieldFarmingPools[poolID]; exists {
		return errors.New("farming pool already exists")
	}

	// Create a new yield farming pool (FarmingPool)
	farmingPool := FarmingPool{
		PoolID:         poolID,
		TokenPair:      tokenPair,
		TotalLiquidity: totalLiquidity,
		RewardRate:     rewardRate,
		CreatedAt:      time.Now(),
		Status:         "Active",
	}

	// Assign the new farming pool to YieldFarmingPools
	l.YieldFarmingPools[poolID] = farmingPool
	fmt.Printf("Yield farming pool %s created for token pair %s\n", poolID, tokenPair)
	return nil
}




// RecordLoanRequest records a loan request in a lending pool.
func (l *DeFiLedger) RecordLoanRequest(poolID, loanID, lender, borrower string, amount, collateral, interestRate float64, duration time.Duration) error {
	l.Lock()
	defer l.Unlock()

	// Check if the pool exists
	pool, poolExists := l.LendingPools[poolID]
	if !poolExists {
		return fmt.Errorf("lending pool %s does not exist", poolID)
	}

	// Check if the loan already exists
	for _, existingLoan := range pool.ActiveLoans {
		if existingLoan.LoanID == loanID {
			return errors.New("loan already exists")
		}
	}

	// Create a new loan
	loan := Loan{
		LoanID:       loanID,
		Lender:       lender,
		Borrower:     borrower,
		Amount:       amount,
		Collateral:   collateral,
		InterestRate: interestRate,
		Duration:     duration,
		StartDate:    time.Now(),
		ExpiryDate:   time.Now().Add(duration),
		Status:       "Active",
	}

	// Append the new loan to the lending pool
	pool.ActiveLoans = append(pool.ActiveLoans, &loan)
	l.LendingPools[poolID] = pool

	fmt.Printf("Loan %s requested by borrower %s for amount %.2f in pool %s\n", loanID, borrower, amount, poolID)
	return nil
}


// RecordLoanRepayment records the repayment of a loan.
func (l *DeFiLedger) RecordLoanRepayment(poolID, loanID string) error {
	l.Lock()
	defer l.Unlock()

	// Check if the lending pool exists
	pool, poolExists := l.LendingPools[poolID]
	if !poolExists {
		return fmt.Errorf("lending pool %s does not exist", poolID)
	}

	// Find the loan within the lending pool
	for _, loan := range pool.ActiveLoans {
		if loan.LoanID == loanID {
			loan.Status = "Repaid"
			fmt.Printf("Loan %s has been repaid in pool %s\n", loanID, poolID)
			return nil
		}
	}

	return fmt.Errorf("loan %s does not exist in pool %s", loanID, poolID)
}


// RecordSyntheticAssetCreation records the creation of a synthetic asset.
func (l *DeFiLedger) RecordSyntheticAssetCreation(assetID, assetName, underlyingAsset string, collateralRatio, totalSupply float64) error {
	l.Lock()
	defer l.Unlock()

	// Check if the synthetic asset already exists
	if _, exists := l.SyntheticAssets[assetID]; exists {
		return errors.New("synthetic asset already exists")
	}

	// Create a new synthetic asset
	asset := &SyntheticAsset{  // Use pointer to create the asset
		AssetID:         assetID,
		AssetName:       assetName,
		UnderlyingAsset: underlyingAsset,
		CollateralRatio: collateralRatio,
		TotalSupply:     totalSupply,
		CreatedAt:       time.Now(),
		Status:          "Active",
	}

	// Store the pointer to the asset in the map
	l.SyntheticAssets[assetID] = asset
	fmt.Printf("Synthetic asset %s created with total supply %.2f\n", assetName, totalSupply)
	return nil
}


// RecordMintingEvent records a minting event for synthetic assets.
func (l *DeFiLedger) RecordMintingEvent(assetID string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	if asset, exists := l.SyntheticAssets[assetID]; exists {
		asset.TotalSupply += amount
		l.SyntheticAssets[assetID] = asset
		fmt.Printf("Minted %.2f of synthetic asset %s\n", amount, assetID)
		return nil
	}
	return errors.New("synthetic asset does not exist")
}

// RecordBurningEvent records a burning event for synthetic assets.
func (l *DeFiLedger) RecordBurningEvent(assetID string, amount float64) error {
	l.Lock()
	defer l.Unlock()

	if asset, exists := l.SyntheticAssets[assetID]; exists {
		asset.TotalSupply -= amount
		l.SyntheticAssets[assetID] = asset
		fmt.Printf("Burned %.2f of synthetic asset %s\n", amount, assetID)
		return nil
	}
	return errors.New("synthetic asset does not exist")
}

// Helper function to generate a unique claim ID
func generateClaimID() string {
	return fmt.Sprintf("claim_%d", time.Now().UnixNano())
}

// RecordLendingPoolCreation logs the creation of a new lending pool in the ledger
func (l *DeFiLedger) RecordLendingPoolCreation(poolID string, liquidity, interestRate float64) error {
    // Check for errors or constraints, e.g., if the poolID already exists (optional)
    if _, exists := l.LendingPools[poolID]; exists {
        return fmt.Errorf("lending pool %s already exists", poolID)
    }

    // Perform logging or recording actions here
    l.LendingPools[poolID] = &LendingPool{
        PoolID:         poolID,
        TotalLiquidity: liquidity,
        InterestRate:   interestRate,
    }

    fmt.Printf("Lending pool %s with %f liquidity and %f interest rate logged in the ledger\n", poolID, liquidity, interestRate)
    return nil
}


// RecordYieldFarming logs a yield farming event in the ledger
func (l *DeFiLedger) RecordYieldFarming(poolID string, amountStaked, rewardRate float64) error {
    // Ensure the pool exists before recording the event
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool %s not found", poolID)
    }

    // Record the yield farming event (this could be saved to a specific record structure or database)
    pool.TotalStaked += amountStaked
    pool.RewardRate = rewardRate // Update the reward rate if needed

    fmt.Printf("Yield farming recorded: pool %s, amount staked %f, reward rate %f\n", poolID, amountStaked, rewardRate)
    return nil
}


// RecordLiquidityPoolClosure logs the closure of a liquidity pool in the ledger
func (l *DeFiLedger) RecordLiquidityPoolClosure(poolID string) error {
    // Ensure the pool exists before recording closure
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool %s not found", poolID)
    }

    // Mark the pool as closed
    pool.Status = "Closed"

    // Perform any necessary cleanup or removal of the pool
    delete(l.LiquidityPools, poolID)

    fmt.Printf("Liquidity pool %s has been closed and removed from the ledger\n", poolID)
    return nil
}

// RecordFarmingPoolCreation logs the creation of a farming pool in the ledger
func (l *DeFiLedger) RecordFarmingPoolCreation(poolID string, initialLiquidity, rewardRate float64) error {
    // Check if the pool already exists
    if _, exists := l.YieldFarmingPools[poolID]; exists {
        return fmt.Errorf("farming pool %s already exists", poolID)
    }

    // Add the farming pool to the ledger
    l.YieldFarmingPools[poolID] = FarmingPool{
        PoolID:          poolID,
        TotalLiquidity:  initialLiquidity,
        RewardRate:      rewardRate,
        CreatedAt:       time.Now(),
        Status:          "Active",
    }

    // Perform logging
    fmt.Printf("Farming pool %s created with initial liquidity %f and reward rate %f\n", poolID, initialLiquidity, rewardRate)

    return nil
}


// RecordLiquidityStake logs a liquidity stake event in the ledger
func (l *DeFiLedger) RecordLiquidityStake(poolID, userID string, amountStaked float64) error {
    // Ensure the pool exists
    pool, exists := l.YieldFarmingPools[poolID]
    if !exists {
        return fmt.Errorf("farming pool %s not found", poolID)
    }

    // Update the pool's total liquidity staked
    pool.TotalLiquidity += amountStaked
    l.YieldFarmingPools[poolID] = pool // Update the pool in the map

    // Add the staking event to the user's record
    l.StakeRecords[userID] = append(l.StakeRecords[userID], &StakeRecord{
        PoolID:       poolID,
        UserID:       userID,
        AmountStaked: amountStaked,
        Timestamp:    time.Now(),
    })

    // Perform logging
    fmt.Printf("Liquidity of %f staked in pool %s by user %s\n", amountStaked, poolID, userID)

    return nil
}


// RecordRewardClaim logs a reward claim event in the ledger
func (l *DeFiLedger) RecordRewardClaim(userID, farmingID string, rewardAmount float64) error {
    // Ensure the farming record exists
    farmingRecord, exists := l.YieldFarmingRecords[farmingID]
    if (!exists) {
        return fmt.Errorf("farming record %s not found", farmingID)
    }

    // Update the user's reward balance
    farmingRecord.RewardsEarned += rewardAmount

    // Perform logging
    fmt.Printf("User %s claimed %f rewards from farming record %s\n", userID, rewardAmount, farmingID)

    return nil
}



// RecordLiquidityUnstake logs a liquidity unstake event in the ledger
func (l *DeFiLedger) RecordLiquidityUnstake(poolID, userID string, amountUnstaked float64) error {
    // Ensure the pool exists
    pool, exists := l.YieldFarmingPools[poolID]
    if !exists {
        return fmt.Errorf("farming pool %s not found", poolID)
    }

    // Update the pool's total liquidity
    if amountUnstaked > pool.TotalLiquidity {
        return fmt.Errorf("unstake amount exceeds staked liquidity in pool %s", poolID)
    }
    pool.TotalLiquidity -= amountUnstaked
    l.YieldFarmingPools[poolID] = pool // Update the pool in the map

    // Perform logging
    fmt.Printf("Liquidity of %f unstaked from pool %s by user %s\n", amountUnstaked, poolID, userID)

    return nil
}


// RecordFarmingPoolClosure logs the closure of a farming pool in the ledger
func (l *DeFiLedger) RecordFarmingPoolClosure(poolID string) error {
    // Ensure the pool exists
    pool, exists := l.YieldFarmingPools[poolID]
    if !exists {
        return fmt.Errorf("farming pool %s not found", poolID)
    }

    // Update the pool status
    pool.Status = "Closed"
    l.YieldFarmingPools[poolID] = pool // Update the pool in the map

    // Perform logging
    fmt.Printf("Farming pool %s has been closed\n", poolID)

    return nil
}


// UpdatePool modifies the liquidity of a specific pool.
func (l *DeFiLedger) UpdatePool(poolID string, pool LiquidityPool) error {
	l.Lock()
	defer l.Unlock()

	l.LiquidityPools[poolID] = pool
	return nil
}

// GetPool retrieves a liquidity pool by its ID.
func (l *DeFiLedger) GetPool(poolID string) (*LiquidityPool, error) {
	l.Lock()
	defer l.Unlock()

	pool, exists := l.LiquidityPools[poolID]
	if !exists {
		return nil, errors.New("liquidity pool not found")
	}
	return &pool, nil
}

func (l *DeFiLedger) CreateBet(bet Bet) error {
	if _, exists := l.Bets[bet.BetID]; exists {
		return fmt.Errorf("bet ID already exists")
	}
	l.Bets[bet.BetID] = bet
	return nil
}

func (l *DeFiLedger) PlaceBet(betID, user string, amount float64) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	if bet.Status != "Open" {
		return fmt.Errorf("bet is not open for placing bets")
	}
	l.BetParticipants[betID] = append(l.BetParticipants[betID], BetParticipant{BetID: betID, User: user, Amount: amount})
	l.BettingEscrowFundsBalance[betID] += amount
	return nil
}

func (l *DeFiLedger) SetOdds(betID string, odds float64) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	bet.Odds = odds
	l.Bets[betID] = bet
	return nil
}

func (l *DeFiLedger) GetBetStatus(betID string) (string, error) {
	bet, exists := l.Bets[betID]
	if !exists {
		return "", fmt.Errorf("bet not found")
	}
	return bet.Status, nil
}

func (l *DeFiLedger) TrackBet(betID string) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	// Logic to track bet changes (e.g., status updates)
	bet.Status = "Tracking"
	l.Bets[betID] = bet
	return nil
}

func (l *DeFiLedger) DistributeWinnings(betID string) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	if bet.Status != "Closed" {
		return fmt.Errorf("bet must be closed to distribute winnings")
	}

	totalWinnings := l.EscrowFundsBalance[betID]
	participants, exists := l.BetParticipants[betID]
	if !exists || len(participants) == 0 {
		return fmt.Errorf("no participants found for the bet")
	}

	for _, participant := range participants {
		payout := participant.Amount * bet.Odds
		fmt.Printf("Distributing %.2f to %s\n", payout, participant.User)
	}

	delete(l.BettingEscrowFundsBalance, betID)
	return nil
}

func (l *DeFiLedger) EscrowFunds(betID string, amount float64) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	if bet.Status != "Open" {
		return fmt.Errorf("bet is not open for escrowing funds")
	}
	l.BettingEscrowFundsBalance[betID] += amount
	return nil
}

func (l *DeFiLedger) ReleaseEscrowedFunds(betID string) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	if bet.Status != "Closed" {
		return fmt.Errorf("bet must be closed to release escrowed funds")
	}
	delete(l.BettingEscrowFundsBalance, betID)
	return nil
}

func (l *DeFiLedger) AuditBet(betID string) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	// Perform an audit by checking bet integrity and participants
	participants, exists := l.BetParticipants[betID]
	if !exists {
		return fmt.Errorf("no participants found for the bet")
	}
	if len(participants) == 0 {
		return fmt.Errorf("bet has no participants")
	}
	fmt.Printf("Audit for bet ID %s completed successfully.\n", betID)
	return nil
}

func (l *DeFiLedger) MonitorOdds(betID string) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	// Example monitoring logic (could be a callback for real-time updates)
	fmt.Printf("Monitoring odds for bet ID %s...\n", betID)
	return nil
}

func (l *DeFiLedger) GetBetHistory(betID string) ([]BetHistoryRecord, error) {
	history, exists := l.BetHistory[betID]
	if !exists {
		return nil, fmt.Errorf("no history found for bet ID %s", betID)
	}
	return history, nil
}

func (l *DeFiLedger) SetMaxBet(betID string, maxAmount float64) error {
	if maxAmount <= 0 {
		return fmt.Errorf("maximum bet amount must be positive")
	}
	l.MaximumBetLimits[betID] = maxAmount
	return nil
}

func (l *DeFiLedger) GetMaxBet(betID string) (float64, error) {
	maxBet, exists := l.MaximumBetLimits[betID]
	if !exists {
		return 0, fmt.Errorf("maximum bet amount not set for bet ID %s", betID)
	}
	return maxBet, nil
}

func (l *DeFiLedger) SetBetExpiration(betID string, expiration time.Time) error {
	bet, exists := l.Bets[betID]
	if !exists {
		return fmt.Errorf("bet not found")
	}
	bet.Expiration = expiration
	l.Bets[betID] = bet
	return nil
}

func (l *DeFiLedger) GetBetExpiration(betID string) (time.Time, error) {
	bet, exists := l.Bets[betID]
	if !exists {
		return time.Time{}, fmt.Errorf("bet not found")
	}
	return bet.Expiration, nil
}

func (l *DeFiLedger) UpdateConfig(key string, value interface{}) error {
	switch key {
	case "BettingPaused":
		boolValue, ok := value.(bool)
		if !ok {
			return fmt.Errorf("invalid value type for BettingPaused")
		}
		l.Configurations.BettingPaused = boolValue
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}
	return nil
}


func (l *DeFiLedger) CreateCampaign(campaign CrowdfundingCampaign) error {
	if _, exists := l.Campaigns[campaign.CampaignID]; exists {
		return fmt.Errorf("campaign ID already exists")
	}
	l.CrowdfundingCampaigns[campaign.CampaignID] = campaign
	return nil
}

func (l *Ledger) Contribute(campaignID, userID string, amount float64) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	if campaign.Status != "Active" {
		return fmt.Errorf("campaign is not active")
	}
	if campaign.EndTime.Before(time.Now()) {
		return fmt.Errorf("campaign has ended")
	}

	contribution := CrowdfundingContribution{
		CampaignID: campaignID,
		UserID:     userID,
		Amount:     amount,
		Time:       time.Now(),
	}

	l.Contributions[campaignID] = append(l.Contributions[campaignID], contribution)
	campaign.CollectedFunds += amount
	l.Campaigns[campaignID] = campaign
	return nil
}

func (l *DeFiLedger) RefundContributors(campaignID string) error {
	contributions, exists := l.Contributions[campaignID]
	if !exists {
		return fmt.Errorf("no contributions found for campaign ID %s", campaignID)
	}

	for _, contribution := range contributions {
		fmt.Printf("Refunding %.2f to user %s for campaign ID %s.\n", contribution.Amount, contribution.UserID, campaignID)
	}

	delete(l.Contributions, campaignID)
	l.CrowdfundingCampaigns[campaignID] = CrowdfundingCampaign{}
	return nil
}

func (l *DeFiLedger) DistributeFunds(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	if campaign.Status != "Closed" {
		return fmt.Errorf("campaign must be closed to distribute funds")
	}

	fmt.Printf("Distributing %.2f to campaign creator %s.\n", campaign.CollectedFunds, campaign.CreatorID)
	delete(l.Campaigns, campaignID)
	return nil
}

func (l *DeFiLedger) AuditCampaign(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}

	audit := CrowdfundingAuditRecord{
		CampaignID: campaignID,
		Details:    fmt.Sprintf("Audit of campaign %s: %f collected out of %f goal.", campaignID, campaign.CollectedFunds, campaign.GoalAmount),
	}

	l.AuditRecords[campaignID] = append(l.AuditRecords[campaignID], audit)
	return nil
}

func (l *DeFiLedger)TrackCampaignProgress(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}

	progress := campaign.CollectedFunds / campaign.GoalAmount * 100
	fmt.Printf("Campaign %s is %.2f%% funded.\n", campaignID, progress)
	return nil
}

func (l *DeFiLedger) FetchCampaignDetails(campaignID string) (CrowdfundingCampaign, error) {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return CrowdfundingCampaign{}, fmt.Errorf("campaign not found")
	}
	return campaign, nil
}

func (l *DeFiLedger) CloseCampaign(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}

	if campaign.EndTime.Before(time.Now()) || campaign.CollectedFunds >= campaign.GoalAmount {
		campaign.Status = "Closed"
	} else {
		campaign.Status = "Failed"
	}

	l.Campaigns[campaignID] = campaign
	return nil
}

func (l *DeFiLedger) LockFunds(campaignID string, amount float64) error {
	l.EscrowFunds[campaignID] += amount
	return nil
}

func (l *DeFiLedger) UnlockFunds(campaignID string) error {
	delete(l.EscrowFunds, campaignID)
	return nil
}


func (l *DeFiLedger) ExtendCampaignDuration(campaignID string, newEndTime time.Time) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	if newEndTime.Before(campaign.EndTime) {
		return fmt.Errorf("new end time must be later than the current end time")
	}
	campaign.EndTime = newEndTime
	l.Campaigns[campaignID] = campaign
	return nil
}

func (l *DeFiLedger) GetContributionHistory(campaignID string) ([]ContributionRecord, error) {
	history, exists := l.Contributions[campaignID]
	if !exists {
		return nil, fmt.Errorf("no contributions found for campaign ID %s", campaignID)
	}
	return history, nil
}

func (l *DeFiLedger) SetContributionLimits(campaignID string, minLimit, maxLimit float64) error {
	if minLimit < 0 || maxLimit <= 0 || minLimit >= maxLimit {
		return fmt.Errorf("invalid contribution limits")
	}
	l.ContributionLimits[campaignID] = ContributionLimits{Min: minLimit, Max: maxLimit}
	return nil
}

func (l *DeFiLedger) GetContributionLimits(campaignID string) (float64, float64, error) {
	limits, exists := l.ContributionLimits[campaignID]
	if !exists {
		return 0, 0, fmt.Errorf("contribution limits not set for campaign ID %s", campaignID)
	}
	return limits.Min, limits.Max, nil
}

func (l *DeFiLedger) PauseCampaign(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	if campaign.Status != "Active" {
		return fmt.Errorf("only active campaigns can be paused")
	}
	campaign.Status = "Paused"
	l.Campaigns[campaignID] = campaign
	l.PausedCampaigns[campaignID] = true
	return nil
}

func (l *DeFiLedger) ResumeCampaign(campaignID string) error {
	campaign, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	if campaign.Status != "Paused" {
		return fmt.Errorf("campaign is not paused")
	}
	campaign.Status = "Active"
	l.Campaigns[campaignID] = campaign
	delete(l.PausedCampaigns, campaignID)
	return nil
}

func (l *DeFiLedger) MonitorContributionFlow(campaignID string) error {
	_, exists := l.Campaigns[campaignID]
	if !exists {
		return fmt.Errorf("campaign not found")
	}
	// Real-time monitoring logic to be added here.
	fmt.Printf("Monitoring contributions for campaign ID %s...\n", campaignID)
	return nil
}


// CreatePolicy adds a new policy to the ledger
func (l *DeFiLedger) CreatePolicy(policy InsurancePolicy) error {
    if _, exists := l.InsurancePolicies[policy.PolicyID]; exists {
        return fmt.Errorf("policy with ID %s already exists", policy.PolicyID)
    }
    l.InsurancePolicies[policy.PolicyID] = policy
    return nil
}

// ActivatePolicy updates a policy's status to active
func (l *DeFiLedger) ActivatePolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status != "Pending" {
        return fmt.Errorf("policy is not in a pending state")
    }
    policy.Status = "Active"
    policy.StartTime = time.Now()
    policy.EndTime = time.Now().Add(policy.Duration)
    l.InsurancePolicies[policyID] = policy
    return nil
}

// ClaimPolicy processes a claim on an active policy
func (l *DeFiLedger) ClaimPolicy(policyID, claimantID string, claimAmount float64) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists || policy.Status != "Active" {
        return fmt.Errorf("policy with ID %s is not active", policyID)
    }
    if claimAmount > policy.CoverageAmount {
        return fmt.Errorf("claim amount exceeds coverage")
    }
    claim := InsuranceClaim{
        PolicyID:    policyID,
        ClaimantID:  claimantID,
        ClaimAmount: claimAmount,
        Status:      "Pending",
        Timestamp:   time.Now(),
    }
    l.InsuranceClaims[policyID] = claim
    return nil
}

// EscrowFunds places funds into escrow
func (l *DeFiLedger) EscrowFunds(policyID string, amount float64) error {
    if _, exists := l.InsurancePolicies[policyID]; !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    l.EscrowBalances[policyID] += amount
    return nil
}

// ReleaseEscrow releases funds from escrow
func (l *DeFiLedger) ReleaseEscrow(policyID string) error {
    if _, exists := l.EscrowBalances[policyID]; !exists {
        return fmt.Errorf("no escrow found for policy ID %s", policyID)
    }
    delete(l.EscrowBalances, policyID)
    return nil
}

// SetPremium updates the premium for a policy
func (l *DeFiLedger) SetPremium(policyID string, premium float64) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Premium = premium
    l.InsurancePolicies[policyID] = policy
    return nil
}

// CalculatePayout calculates the payout for a claim
func (l *DeFiLedger) CalculatePayout(policyID string, claimAmount float64) (float64, error) {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return 0, fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if claimAmount > policy.CoverageAmount {
        return policy.CoverageAmount, nil
    }
    return claimAmount, nil
}

func (l *DeFiLedger) VerifyClaim(policyID, claimantID string) (bool, error) {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return false, fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status != "Active" || policy.Frozen || policy.Locked {
        return false, fmt.Errorf("policy is not eligible for claims")
    }
    return true, nil
}


func (l *DeFiLedger) DistributePayout(policyID, claimantID string, payoutAmount float64) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status != "Active" || payoutAmount > policy.CoverageAmount {
        return fmt.Errorf("invalid payout conditions")
    }
    // Logic to transfer funds (e.g., via a payment gateway)
    fmt.Printf("Distributed %f to %s for policy %s.\n", payoutAmount, claimantID, policyID)
    return nil
}

func (l *DeFiLedger) TrackPolicyStatus(policyID string) (string, error) {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return "", fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    return policy.Status, nil
}

func (l *DeFiLedger) AuditPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status == "Active" && policy.CoverageAmount > 0 {
        fmt.Printf("Policy %s passed audit.\n", policyID)
    } else {
        return fmt.Errorf("policy %s failed audit", policyID)
    }
    return nil
}

func (l *DeFiLedger) EscrowFunds(policyID string, amount float64) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    l.EscrowBalances[policyID] += amount
    return nil
}

func (l *DeFiLedger) ReleaseEscrow(policyID string) error {
    if _, exists := l.EscrowBalances[policyID]; !exists {
        return fmt.Errorf("no escrow funds found for policy ID %s", policyID)
    }
    delete(l.EscrowBalances, policyID)
    return nil
}

func (l *DeFiLedger) CancelPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status != "Active" {
        return fmt.Errorf("policy must be active to cancel")
    }
    policy.Status = "Canceled"
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) CancelPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    if policy.Status != "Active" {
        return fmt.Errorf("policy must be active to cancel")
    }
    policy.Status = "Canceled"
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) LockPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Locked = true
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) UnlockPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Locked = false
    l.InsurancePolicies[policyID] = policy
    return nil
}


func (l *DeFiLedger) FetchPolicyTerms(policyID string) (string, error) {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return "", fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    return policy.Terms, nil
}

func (l *DeFiLedger) UpdatePolicyTerms(policyID, newTerms string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Terms = newTerms
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) FreezePolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Frozen = true
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) UnfreezePolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.Frozen = false
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) GetClaimHistory(policyID string) ([]ClaimRecord, error) {
    history, exists := l.ClaimHistory[policyID]
    if !exists {
        return nil, fmt.Errorf("no claim history found for policy ID %s", policyID)
    }
    return history, nil
}

func (l *DeFiLedger) SetClaimProcessingFee(policyID string, fee float64) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.ClaimFee = fee
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger) GetClaimProcessingFee(policyID string) (float64, error) {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return 0, fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    return policy.ClaimFee, nil
}

func (l *DeFiLedger) AutoRenewPolicy(policyID string) error {
    policy, exists := l.InsurancePolicies[policyID]
    if !exists {
        return fmt.Errorf("policy with ID %s does not exist", policyID)
    }
    policy.AutoRenew = true
    l.InsurancePolicies[policyID] = policy
    return nil
}

func (l *DeFiLedger)CreateLiquidityPool(pool LiquidityPool) error {
    if _, exists := l.LiquidityPools[pool.PoolID]; exists {
        return fmt.Errorf("liquidity pool with ID %s already exists", pool.PoolID)
    }
    l.LiquidityPools[pool.PoolID] = pool
    return nil
}

func (l *DeFiLedger) DepositToPool(poolID string, amount1, amount2 float64) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    if pool.IsLocked {
        return fmt.Errorf("liquidity pool %s is locked", poolID)
    }
    pool.TotalBalance += amount1 + amount2
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) WithdrawFromPool(poolID string, withdrawalAmount float64) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    if pool.TotalBalance < withdrawalAmount {
        return fmt.Errorf("insufficient balance in pool %s", poolID)
    }
    pool.TotalBalance -= withdrawalAmount
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) SwapTokens(poolID, tokenIn string, amountIn float64) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    if pool.IsLocked {
        return 0, fmt.Errorf("liquidity pool %s is locked", poolID)
    }
    // Example swap logic: 1:1 ratio swap
    amountOut := amountIn * pool.TokenRatio
    pool.FeesAccumulated += amountIn * 0.01 // Example 1% fee
    pool.TotalBalance += amountIn - amountOut
    l.LiquidityPools[poolID] = pool
    return amountOut, nil
}

func (l *DeFiLedger) TrackPoolBalance(poolID string) (float64, float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.Token1, pool.Token2, nil
}

func (l *DeFiLedger) GetTokenPrice(poolID, token string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.TokenRatio, nil
}

func (l *DeFiLedger) AuditPool(poolID string) error {
    if _, exists := l.LiquidityPools[poolID]; !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return nil // Example: pass audit
}

func (l *DeFiLedger) CalculatePoolYield(poolID string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.FeesAccumulated * 0.95, nil // 95% distributed to providers
}

func (l *DeFiLedger) DistributePoolFees(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.FeesAccumulated = 0
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) StakeLP(poolID, userID string, amount float64) error {
    stake := LPStaking{
        PoolID:   poolID,
        UserID:   userID,
        Amount:   amount,
        StakedAt: time.Now(),
    }
    l.LPStakings[poolID] = append(l.LPStakings[poolID], stake)
    return nil
}

func (l *DeFiLedger) UnstakeLP(poolID, userID string, amount float64) error {
    stakings := l.LPStakings[poolID]
    for i, stake := range stakings {
        if stake.UserID == userID && stake.Amount >= amount {
            stake.Amount -= amount
            if stake.Amount == 0 {
                stakings = append(stakings[:i], stakings[i+1:]...)
            }
            l.LPStakings[poolID] = stakings
            return nil
        }
    }
    return fmt.Errorf("no sufficient LP tokens staked by %s in pool %s", userID, poolID)
}

func (l *DeFiLedger) LockPool(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = true
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) UnlockPool(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = false
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) FetchPoolBalance(poolID string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.TotalBalance, nil
}


func (l *DeFiLedger) SetPoolFeeRate(poolID string, feeRate float64) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.FeeRate = feeRate
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger)FetchPoolFeeRate(poolID string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.FeeRate, nil
}

func (l *DeFiLedger) SetPoolTokenRatio(poolID string, ratio float64) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.TokenRatio = ratio
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) FetchPoolTokenRatio(poolID string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.TokenRatio, nil
}

func (l *DeFiLedger) CompoundPoolInterest(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    now := time.Now()
    elapsed := now.Sub(pool.LastCompoundTime).Hours() / 24
    if elapsed < 1 {
        return fmt.Errorf("interest compounding is only allowed after 24 hours")
    }
    compoundInterest := pool.TotalBalance * 0.01 // Example: 1% interest per day
    pool.TotalBalance += compoundInterest
    pool.LastCompoundTime = now
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) SetPoolWithdrawalFee(poolID string, fee float64) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.WithdrawalFee = fee
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) FetchPoolWithdrawalFee(poolID string) (float64, error) {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return 0, fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    return pool.WithdrawalFee, nil
}

func (l *DeFiLedger) PauseSwaps(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.IsSwapsPaused = true
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) ResumeSwaps(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    pool.IsSwapsPaused = false
    l.LiquidityPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) AutoRebalancePool(poolID string) error {
    pool, exists := l.LiquidityPools[poolID]
    if !exists {
        return fmt.Errorf("liquidity pool with ID %s does not exist", poolID)
    }
    if pool.RebalancingActive {
        pool.TokenRatio = 1.0 // Example: Reset to 1:1 ratio for simplicity
        l.LiquidityPools[poolID] = pool
        return nil
    }
    return fmt.Errorf("rebalancing is not active for pool %s", poolID)
}

func (l *DeFiLedger) CreatePredictionEvent(event PredictionEvent) error {
    if _, exists := l.PredictionEvents[event.EventID]; exists {
        return fmt.Errorf("prediction event with ID %s already exists", event.EventID)
    }
    l.PredictionEvents[event.EventID] = event
    return nil
}

func (l *DeFiLedger) PlacePrediction(eventID, userID string, amount float64) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status != "Open" {
        return fmt.Errorf("event %s is not open for predictions", eventID)
    }

    prediction := Prediction{
        EventID: eventID,
        UserID:  userID,
        Amount:  amount,
        Odds:    event.Odds,
        Payout:  amount * event.Odds,
        Status:  "Pending",
    }

    l.Predictions[eventID] = append(l.Predictions[eventID], prediction)
    event.TotalPool += amount
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) SetEventOdds(eventID string, odds float64) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status != "Open" {
        return fmt.Errorf("odds can only be set for open events")
    }
    event.Odds = odds
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) CalculatePayout(eventID string, amount float64) (float64, error) {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return 0, fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    return amount * event.Odds, nil
}

func (l *DeFiLedger)TrackEventOutcome(eventID, outcome string) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status != "Open" {
        return fmt.Errorf("event %s has already been resolved", eventID)
    }
    event.Status = "Resolved"
    event.Outcome = outcome
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) FetchPredictionStatus(eventID string) (string, error) {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return "", fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    return event.Status, nil
}

func (l *DeFiLedger) DistributePredictionRewards(eventID string) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status != "Resolved" {
        return fmt.Errorf("cannot distribute rewards for unresolved event %s", eventID)
    }

    predictions := l.Predictions[eventID]
    for i, prediction := range predictions {
        if prediction.Status == "Pending" && event.Outcome == "Win" {
            predictions[i].Payout = prediction.Amount * event.Odds
            predictions[i].Status = "Paid"
        } else {
            predictions[i].Payout = 0
            predictions[i].Status = "Lost"
        }
    }
    l.Predictions[eventID] = predictions
    return nil
}

func (l *DeFiLedger) AuditPrediction(eventID string) error {
    if _, exists := l.PredictionEvents[eventID]; !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    return nil
}

func (l *DeFiLedger) MonitorEventOutcome(eventID string) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status == "Open" {
        return fmt.Errorf("cannot monitor outcome for open event %s", eventID)
    }
    return nil
}

func (l *DeFiLedger)SetMaxPrediction(eventID string, maxAmount float64) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    event.MaximumAmount = maxAmount
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) FetchMaxPrediction(eventID string) (float64, error) {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return 0, fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    return event.MaximumAmount, nil
}

func (l *DeFiLedger) SetEventExpiration(eventID string, expiration time.Time) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    event.Expiration = expiration
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) FetchEventExpiration(eventID string) (time.Time, error) {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return time.Time{}, fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    return event.Expiration, nil
}

func (l *DeFiLedger) SettleEvent(eventID, outcome string) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if event.Status != "Open" {
        return fmt.Errorf("event %s is not open for settlement", eventID)
    }
    event.Status = "Settled"
    event.Outcome = outcome
    l.PredictionEvents[eventID] = event
    return nil
}

func (l *DeFiLedger) TrackParticipant(eventID, userID string, predictionAmount float64) error {
    event, exists := l.PredictionEvents[eventID]
    if !exists {
        return fmt.Errorf("prediction event with ID %s does not exist", eventID)
    }
    if predictionAmount > event.MaximumAmount {
        return fmt.Errorf("prediction amount exceeds the maximum allowed for event %s", eventID)
    }
    participantPrediction := ParticipantPrediction{
        UserID:  userID,
        Amount:  predictionAmount,
        Odds:    0, // To be updated when event is resolved
        EventID: eventID,
        Status:  "Pending",
    }
    l.ParticipantHistories[userID] = append(l.ParticipantHistories[userID], participantPrediction)
    return nil
}

func (l *DeFiLedger) FetchParticipantHistory(eventID, userID string) ([]ParticipantPrediction, error) {
    history, exists := l.ParticipantHistories[userID]
    if !exists {
        return nil, fmt.Errorf("no prediction history found for user %s", userID)
    }
    eventHistory := []ParticipantPrediction{}
    for _, prediction := range history {
        if prediction.EventID == eventID {
            eventHistory = append(eventHistory, prediction)
        }
    }
    if len(eventHistory) == 0 {
        return nil, fmt.Errorf("no predictions found for user %s in event %s", userID, eventID)
    }
    return eventHistory, nil
}

func (l *DeFiLedger) CreateStakingProgram(program StakingProgram) error {
    if _, exists := l.StakingPrograms[program.ProgramID]; exists {
        return fmt.Errorf("staking program with ID %s already exists", program.ProgramID)
    }
    l.StakingPrograms[program.ProgramID] = program
    return nil
}

func (l *DeFiLedger) StakeTokens(programID, userID string, amount float64) error {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    if program.Status != "Active" {
        return fmt.Errorf("staking program %s is not active", programID)
    }
    if amount < program.MinStake {
        return fmt.Errorf("amount is below the minimum stake requirement")
    }

    participants := l.StakingParticipants[userID]
    for i, participant := range participants {
        if participant.ProgramID == programID {
            participants[i].StakedAmount += amount
            l.StakingParticipants[userID] = participants
            program.TotalStaked += amount
            l.StakingPrograms[programID] = program
            return nil
        }
    }

    participant := StakingParticipant{
        UserID:       userID,
        StakedAmount: amount,
        Rewards:      0,
        ProgramID:    programID,
        Locked:       false,
    }
    l.StakingParticipants[userID] = append(l.StakingParticipants[userID], participant)
    program.TotalStaked += amount
    l.StakingPrograms[programID] = program
    return nil
}

func (l *DeFiLedger) UnstakeTokens(programID, userID string, amount float64) error {
    participants := l.StakingParticipants[userID]
    for i, participant := range participants {
        if participant.ProgramID == programID {
            if participant.Locked {
                return fmt.Errorf("tokens are locked for user %s in program %s", userID, programID)
            }
            if participant.StakedAmount < amount {
                return fmt.Errorf("insufficient staked tokens for user %s in program %s", userID, programID)
            }
            participants[i].StakedAmount -= amount
            l.StakingParticipants[userID] = participants
            program := l.StakingPrograms[programID]
            program.TotalStaked -= amount
            l.StakingPrograms[programID] = program
            return nil
        }
    }
    return fmt.Errorf("no staked tokens found for user %s in program %s", userID, programID)
}

func (l *DeFiLedger) CalculateStakingRewards(programID, userID string) (float64, error) {
    participants := l.StakingParticipants[userID]
    for _, participant := range participants {
        if participant.ProgramID == programID {
            program := l.StakingPrograms[programID]
            rewards := participant.StakedAmount * program.RewardRate
            return rewards, nil
        }
    }
    return 0, fmt.Errorf("no staked tokens found for user %s in program %s", userID, programID)
}

func (l *DeFiLedger) DistributeStakingRewards(programID string) error {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    participants := l.StakingParticipants
    for userID, userParticipants := range participants {
        for i, participant := range userParticipants {
            if participant.ProgramID == programID {
                rewards := participant.StakedAmount * program.RewardRate
                userParticipants[i].Rewards += rewards
            }
        }
    }
    return nil
}

func (l *DeFiLedger) LockTokens(programID, userID string) error {
    participants := l.StakingParticipants[userID]
    for i, participant := range participants {
        if participant.ProgramID == programID {
            participants[i].Locked = true
            l.StakingParticipants[userID] = participants
            return nil
        }
    }
    return fmt.Errorf("no staked tokens found for user %s in program %s", userID, programID)
}

func (l *DeFiLedger) UnlockTokens(programID, userID string) error {
    participants := l.StakingParticipants[userID]
    for i, participant := range participants {
        if participant.ProgramID == programID {
            participants[i].Locked = false
            l.StakingParticipants[userID] = participants
            return nil
        }
    }
    return fmt.Errorf("no staked tokens found for user %s in program %s", userID, programID)
}

func (l *DeFiLedger) AuditStakingProgram(programID string) error {
    if _, exists := l.StakingPrograms[programID]; !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    return nil
}

func (l *Ledger) MonitorStakingProgram(programID string) error {
    if _, exists := l.StakingPrograms[programID]; !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    return nil
}

func (l *DeFiLedger) TakeStakingSnapshot(programID string) error {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    snapshot := StakingSnapshot{
        ProgramID:       programID,
        TotalStaked:     program.TotalStaked,
        ParticipantData: l.StakingParticipants[programID],
        Timestamp:       time.Now(),
    }
    l.StakingSnapshots[programID] = append(l.StakingSnapshots[programID], snapshot)
    return nil
}

func (l *DeFiLedger) FetchStakeAmount(programID, userID string) (float64, error) {
    participants, exists := l.StakingParticipants[programID]
    if !exists {
        return 0, fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    amount, exists := participants[userID]
    if !exists {
        return 0, fmt.Errorf("user %s has no stake in program %s", userID, programID)
    }
    return amount, nil
}

func (l *DeFiLedger) UpdateStakeAmount(programID, userID string, newAmount float64) error {
    participants, exists := l.StakingParticipants[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    if newAmount < 0 {
        return fmt.Errorf("stake amount cannot be negative")
    }
    participants[userID] = newAmount
    l.StakingParticipants[programID] = participants
    return nil
}

func (l *DeFiLedger) FetchRewardHistory(programID, userID string) ([]RewardRecord, error) {
    rewards, exists := l.RewardHistories[userID]
    if !exists {
        return nil, fmt.Errorf("no reward history found for user %s", userID)
    }
    filteredRewards := []RewardRecord{}
    for _, reward := range rewards {
        if reward.ProgramID == programID {
            filteredRewards = append(filteredRewards, reward)
        }
    }
    return filteredRewards, nil
}

func (l *DeFiLedger) DistributeStakingBonuses(programID string) error {
    participants, exists := l.StakingParticipants[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    for userID, amount := range participants {
        bonus := amount * 0.05 // Example bonus: 5% of staked amount
        record := RewardRecord{
            ProgramID:   programID,
            UserID:      userID,
            Reward:      bonus,
            Timestamp:   time.Now(),
            Description: "Staking Bonus",
        }
        l.RewardHistories[userID] = append(l.RewardHistories[userID], record)
    }
    return nil
}

func (l *DeFiLedger) ReclaimStakingRewards(programID, userID string) error {
    rewards, exists := l.RewardHistories[userID]
    if !exists {
        return fmt.Errorf("no rewards found for user %s", userID)
    }
    newRewards := []RewardRecord{}
    for _, reward := range rewards {
        if reward.ProgramID != programID {
            newRewards = append(newRewards, reward)
        }
    }
    l.RewardHistories[userID] = newRewards
    return nil
}

func (l *DeFiLedger) SetMinimumStake(programID string, minAmount float64) error {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    program.MinStake = minAmount
    l.StakingPrograms[programID] = program
    return nil
}

func (l *DeFiLedger) FetchMinimumStake(programID string) (float64, error) {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return 0, fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    return program.MinStake, nil
}

func (l *Ledger) SetRewardDistributionFrequency(programID string, frequency time.Duration) error {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    program.RewardDistributionFreq = frequency
    l.StakingPrograms[programID] = program
    return nil
}

func (l *DeFiLedger) FetchRewardDistributionFrequency(programID string) (time.Duration, error) {
    program, exists := l.StakingPrograms[programID]
    if !exists {
        return 0, fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    return program.RewardDistributionFreq, nil
}

func (l *DeFiLedger) EnableAutoReinvestment(programID, userID string) error {
    participants, exists := l.StakingParticipants[programID]
    if !exists {
        return fmt.Errorf("staking program with ID %s does not exist", programID)
    }
    _, exists = participants[userID]
    if !exists {
        return fmt.Errorf("user %s has no stake in program %s", userID, programID)
    }
    // Placeholder for enabling auto-reinvestment logic
    return nil
}

func (l *DeFiLedger) CreateLoan(loan Loan) error {
    if _, exists := l.Loans[loan.LoanID]; exists {
        return fmt.Errorf("loan with ID %s already exists", loan.LoanID)
    }
    loan.CreatedAt = time.Now()
    loan.DueDate = loan.CreatedAt.Add(loan.Duration)
    loan.RemainingBalance = loan.Principal
    l.Loans[loan.LoanID] = loan
    return nil
}

func (l *DeFiLedger) ApplyInterest(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.Status != "Active" {
        return fmt.Errorf("loan %s is not active", loanID)
    }
    interest := loan.RemainingBalance * (loan.InterestRate / 100)
    loan.RemainingBalance += interest
    l.Loans[loanID] = loan
    return nil
}

func (l *DeFiLedger) RepayLoan(loanID, borrowerID string, repaymentAmount float64) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.BorrowerID != borrowerID {
        return fmt.Errorf("borrower ID does not match for loan %s", loanID)
    }
    if loan.Status != "Active" {
        return fmt.Errorf("loan %s is not active", loanID)
    }
    if repaymentAmount <= 0 {
        return fmt.Errorf("repayment amount must be greater than zero")
    }

    loan.RemainingBalance -= repaymentAmount
    if loan.RemainingBalance <= 0 {
        loan.RemainingBalance = 0
        loan.Status = "Fulfilled"
    }
    l.Loans[loanID] = loan
    l.LoanRepayments[loanID] += repaymentAmount
    return nil
}

func (l *DeFiLedger) LiquidateLoan(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.Status != "Defaulted" {
        return fmt.Errorf("loan %s is not in default status", loanID)
    }
    loan.Status = "Liquidated"
    l.Loans[loanID] = loan
    return nil
}

func (l *DeFiLedger) TrackRepayment(loanID string) (float64, error) {
    repayment, exists := l.LoanRepayments[loanID]
    if !exists {
        return 0, fmt.Errorf("no repayment records found for loan %s", loanID)
    }
    return repayment, nil
}

func (l *DeFiLedger) FetchLoanStatus(loanID string) (string, error) {
    loan, exists := l.Loans[loanID]
    if !exists {
        return "", fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    return loan.Status, nil
}

func (l *DeFiLedger) AdjustInterestRate(loanID string, newRate float64) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    loan.InterestRate = newRate
    l.Loans[loanID] = loan
    return nil
}

func (l *DeFiLedger) EscrowCollateral(loanID, collateral string) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    l.CollateralEscrows[loanID] = collateral
    return nil
}

func (l *DeFiLedger) ReleaseCollateral(loanID string) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if _, exists := l.CollateralEscrows[loanID]; !exists {
        return fmt.Errorf("no collateral escrowed for loan %s", loanID)
    }
    delete(l.CollateralEscrows, loanID)
    return nil
}

func (l *DeFiLedger) VerifyCollateral(loanID, collateral string) (bool, error) {
    storedCollateral, exists := l.CollateralEscrows[loanID]
    if !exists {
        return false, fmt.Errorf("no collateral escrowed for loan %s", loanID)
    }
    return storedCollateral == collateral, nil
}

func (l *DeFiLedger) AuditLoan(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    audit := LoanAuditRecord{
        LoanID:       loanID,
        AuditDetails: "Loan is active and meets all criteria",
        Timestamp:    time.Now(),
    }
    l.LoanAudits[loanID] = append(l.LoanAudits[loanID], audit)
    return nil
}

func (l *DeFiLedger)CalculateLoanHealth(loanID string) (float64, error) {
    loan, exists := l.Loans[loanID]
    if !exists {
        return 0, fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.RemainingBalance == 0 {
        return 100, nil // Loan is fully paid
    }
    health := 100 - ((loan.RemainingBalance / loan.Principal) * 100)
    return health, nil
}

func (l *DeFiLedger) TrackBorrower(borrowerID string) error {
    for _, loan := range l.Loans {
        if loan.BorrowerID == borrowerID {
            return nil
        }
    }
    return fmt.Errorf("no loans found for borrower %s", borrowerID)
}

func (l *DeFiLedger) SetCollateralRequirement(loanID string, requirement float64) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    l.CollateralRequirements[loanID] = requirement
    return nil
}

func (l *DeFiLedger) FetchCollateralRequirement(loanID string) (float64, error) {
    requirement, exists := l.CollateralRequirements[loanID]
    if !exists {
        return 0, fmt.Errorf("no collateral requirement found for loan %s", loanID)
    }
    return requirement, nil
}

func (l *DeFiLedger) UpdateRepaymentSchedule(loanID string, newSchedule []time.Time) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    l.LoanRepaymentSchedules[loanID] = newSchedule
    return nil
}

func (l *DeFiLedger) PauseRepayments(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.Status != "Active" {
        return fmt.Errorf("loan %s is not active", loanID)
    }
    loan.Status = "Paused"
    l.Loans[loanID] = loan
    return nil
}

func (l *DeFiLedger) ResumeRepayments(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    if loan.Status != "Paused" {
        return fmt.Errorf("loan %s is not paused", loanID)
    }
    loan.Status = "Active"
    l.Loans[loanID] = loan
    return nil
}

func (l *DeFiLedger) TrackLatePayments(loanID string) error {
    loan, exists := l.Loans[loanID]
    if !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    // Example logic: Check if payments are overdue
    for _, dueDate := range l.LoanRepaymentSchedules[loanID] {
        if time.Now().After(dueDate) {
            penaltyFee := 50.0 // Example penalty fee
            latePayment := LatePaymentRecord{
                LoanID:     loanID,
                DueDate:    dueDate,
                PaidDate:   time.Now(),
                PenaltyFee: penaltyFee,
            }
            l.LatePayments[loanID] = append(l.LatePayments[loanID], latePayment)
        }
    }
    return nil
}

func (l *DeFiLedger) FetchLatePaymentHistory(loanID string) ([]LatePaymentRecord, error) {
    history, exists := l.LatePayments[loanID]
    if !exists {
        return nil, fmt.Errorf("no late payment history found for loan %s", loanID)
    }
    return history, nil
}

func (l *DeFiLedger) FetchInterestRate(loanID string) (float64, error) {
    loan, exists := l.Loans[loanID]
    if !exists {
        return 0, fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    return loan.InterestRate, nil
}

func (l *DeFiLedger) AdjustInterestRatePeriod(loanID string, period time.Duration) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    l.InterestRatePeriods[loanID] = period
    return nil
}

func (l *DeFiLedger) EnableAutoRepaymentAdjustments(loanID string) error {
    if _, exists := l.Loans[loanID]; !exists {
        return fmt.Errorf("loan with ID %s does not exist", loanID)
    }
    // Logic for auto-repayment adjustments
    return nil
}

func (l *DeFiLedger) MintSyntheticAsset(assetID string, amount float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    asset.TotalSupply += amount
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger)BurnSyntheticAsset(assetID string, amount float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    if asset.TotalSupply < amount {
        return fmt.Errorf("insufficient supply to burn for asset %s", assetID)
    }
    asset.TotalSupply -= amount
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger) UpdateSyntheticAssetPrice(assetID string, newPrice float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    oldPrice := asset.Price
    asset.Price = newPrice
    l.SyntheticAssets[assetID] = asset
    priceChange := SyntheticAssetPriceChange{
        AssetID:    assetID,
        OldPrice:   oldPrice,
        NewPrice:   newPrice,
        ChangeTime: time.Now(),
    }
    l.SyntheticAssetPrices[assetID] = append(l.SyntheticAssetPrices[assetID], priceChange)
    return nil
}

func (l *DeFiLedger) SetAssetCollateral(assetID string, collateralAmount float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    asset.Collateral = collateralAmount
    asset.Collateralized = asset.Collateral >= asset.TotalSupply*asset.Price
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger) VerifyAssetCollateral(assetID string) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    if asset.Collateral < asset.TotalSupply*asset.Price {
        return fmt.Errorf("collateral for asset %s is insufficient", assetID)
    }
    return nil
}

func (l *DeFiLedger) LiquidateSyntheticAsset(assetID string) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    if asset.Collateral >= asset.TotalSupply*asset.Price {
        return fmt.Errorf("asset %s does not meet liquidation criteria", assetID)
    }
    delete(l.SyntheticAssets, assetID)
    return nil
}

func (l *DeFiLedger) DistributeAssetDividends(assetID string, dividendAmount float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    l.AssetDividends[assetID] += dividendAmount
    return nil
}

func (l *DeFiLedger) SetDividendRate(assetID string, rate float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    asset.DividendRate = rate
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger) GetDividendRate(assetID string) (float64, error) {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return 0, fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    return asset.DividendRate, nil
}

func (l *DeFiLedger) TrackMarketCap(assetID string) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    marketCap := asset.TotalSupply * asset.Price
    record := MarketCapRecord{
        AssetID:   assetID,
        MarketCap: marketCap,
        Timestamp: time.Now(),
    }
    l.SyntheticAssetMarketCap[assetID] = append(l.SyntheticAssetMarketCap[assetID], record)
    return nil
}

func (l *DeFiLedger) GetMarketCap(assetID string) (float64, error) {
    records, exists := l.SyntheticAssetMarketCap[assetID]
    if !exists || len(records) == 0 {
        return 0, fmt.Errorf("no market cap records found for asset %s", assetID)
    }
    return records[len(records)-1].MarketCap, nil
}

func (l *DeFiLedger) GetCollateralRatio(assetID string) (float64, error) {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return 0, fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    if asset.TotalSupply == 0 {
        return 0, fmt.Errorf("collateral ratio cannot be calculated for zero supply")
    }
    return asset.Collateral / (asset.TotalSupply * asset.Price), nil
}

func (l *DeFiLedger) SetCollateralRatio(assetID string, newRatio float64) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    asset.CollateralRatio = newRatio
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger)TrackAssetVolatility(assetID string) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    // Example volatility tracking logic based on price changes
    volatility := 0.05 * asset.Price // Placeholder for volatility calculation
    record := VolatilityRecord{
        AssetID:        assetID,
        VolatilityRate: volatility,
        Timestamp:      time.Now(),
    }
    l.SyntheticAssetVolatility[assetID] = append(l.SyntheticAssetVolatility[assetID], record)
    return nil
}

func (l *DeFiLedger) GetAssetVolatility(assetID string) (float64, error) {
    records, exists := l.SyntheticAssetVolatility[assetID]
    if !exists || len(records) == 0 {
        return 0, fmt.Errorf("no volatility records found for asset %s", assetID)
    }
    return records[len(records)-1].VolatilityRate, nil
}

func (l *DeFiLedger) AutoAdjustCollateral(assetID string) error {
    asset, exists := l.SyntheticAssets[assetID]
    if !exists {
        return fmt.Errorf("synthetic asset with ID %s does not exist", assetID)
    }
    // Example adjustment logic: Increase collateral if volatility is high
    volatility, err := l.getAssetVolatility(assetID)
    if err != nil {
        return fmt.Errorf("failed to fetch volatility: %v", err)
    }
    if volatility > 0.1 { // Example threshold
        asset.Collateral *= 1.1 // Increase collateral by 10%
    }
    l.SyntheticAssets[assetID] = asset
    return nil
}

func (l *DeFiLedger) AddLiquidityToPool(poolID string, amount float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.TotalLiquidity += amount
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) RemoveLiquidityFromPool(poolID string, amount float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    if pool.TotalLiquidity < amount {
        return fmt.Errorf("insufficient liquidity to remove from pool %s", poolID)
    }
    pool.TotalLiquidity -= amount
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) StakeTokensInPool(poolID string, userID string, amount float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.StakedTokens[userID] += amount
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) UnstakeTokensFromPool(poolID string, userID string, amount float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    if pool.StakedTokens[userID] < amount {
        return fmt.Errorf("insufficient staked tokens to unstake for user %s in pool %s", userID, poolID)
    }
    pool.StakedTokens[userID] -= amount
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) HarvestYieldFarmRewards(poolID string, userID string) error {
    earnings, exists := l.YieldFarmEarnings[poolID][userID]
    if !exists {
        return fmt.Errorf("no earnings found for user %s in pool %s", userID, poolID)
    }
    if earnings.EarnedRewards == 0 {
        return fmt.Errorf("no rewards to harvest for user %s in pool %s", userID, poolID)
    }
    earnings.EarnedRewards = 0
    earnings.LastHarvest = time.Now()
    l.YieldFarmEarnings[poolID][userID] = earnings
    return nil
}

func (l *DeFiLedger) DistributePoolRewards(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    if pool.RewardBalance == 0 {
        return fmt.Errorf("no rewards available for distribution in pool %s", poolID)
    }
    for userID, staked := range pool.StakedTokens {
        earnings := l.YieldFarmEarnings[poolID][userID]
        earnings.EarnedRewards += staked * (pool.RewardBalance / pool.TotalLiquidity)
        l.YieldFarmEarnings[poolID][userID] = earnings
    }
    pool.RewardBalance = 0
    pool.LastDistributed = time.Now()
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) CalculateAPY(poolID string) (float64, error) {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return 0, fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    apy := (pool.RewardBalance / pool.TotalLiquidity) * 100
    return apy, nil
}

func (l *DeFiLedger) LockPoolFunds(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = true
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger)UnlockPoolFunds(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = false
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) AuditYieldFarmPool(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    // Example security checks for compliance
    if pool.TotalLiquidity <= 0 {
        return fmt.Errorf("pool %s has no liquidity, audit failed", poolID)
    }
    if len(pool.StakedTokens) == 0 {
        return fmt.Errorf("no participants in pool %s, audit failed", poolID)
    }
    return nil
}

func (l *DeFiLedger) LockYieldFarmPool(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = true
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) UnlockYieldFarmPool(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.IsLocked = false
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) TrackPoolPerformance(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    metrics := PoolPerformanceMetrics{
        PoolID:          poolID,
        TotalLiquidity:  pool.TotalLiquidity,
        TotalRewards:    pool.RewardBalance,
        TotalParticipants: len(pool.StakedTokens),
        APY:             (pool.RewardBalance / pool.TotalLiquidity) * 100, // Simplified APY calculation
        LastUpdated:     time.Now(),
    }
    l.YieldFarmPerformance[poolID] = metrics
    return nil
}

func (l *DeFiLedger) GetPoolPerformanceMetrics(poolID string) (PoolPerformanceMetrics, error) {
    metrics, exists := l.YieldFarmPerformance[poolID]
    if !exists {
        return PoolPerformanceMetrics{}, fmt.Errorf("performance metrics not available for pool %s", poolID)
    }
    return metrics, nil
}

func (l *DeFiLedger) IncreasePoolRewards(poolID string, increment float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    pool.RewardBalance += increment
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) DecreasePoolRewards(poolID string, decrement float64) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    if pool.RewardBalance < decrement {
        return fmt.Errorf("insufficient rewards to decrease in pool %s", poolID)
    }
    pool.RewardBalance -= decrement
    l.YieldFarmPools[poolID] = pool
    return nil
}

func (l *DeFiLedger) CompoundPoolRewards(poolID string) error {
    pool, exists := l.YieldFarmPools[poolID]
    if !exists {
        return fmt.Errorf("yield farming pool with ID %s does not exist", poolID)
    }
    if pool.RewardBalance == 0 {
        return fmt.Errorf("no rewards to compound in pool %s", poolID)
    }
    // Example compounding logic: Add a percentage of the reward back to the pool liquidity
    compounded := pool.RewardBalance * 0.1 // 10% of rewards compounded
    pool.TotalLiquidity += compounded
    pool.RewardBalance -= compounded
    l.YieldFarmPools[poolID] = pool
    return nil
}
