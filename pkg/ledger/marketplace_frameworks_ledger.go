package ledger

import (
	"fmt"
)


// RecordOrderPlacement records a new order placement in the marketplace.
func (ledger *Ledger) RecordOrderPlacement(order Order) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.Orders == nil {
		ledger.State.Marketplace.Orders = make(map[string]*Order)
	}

	ledger.State.Marketplace.Orders[order.OrderID] = &order
	fmt.Printf("Order %s placed by %s.\n", order.OrderID, order.BuyerID)
	return nil
}

// RecordTradeExecution records the execution of a trade between two parties.
func (ledger *Ledger) RecordTradeExecution(trade TradeExecution) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.Trades == nil {
		ledger.State.Marketplace.Trades = make(map[string]*TradeExecution)
	}

	ledger.State.Marketplace.Trades[trade.TradeID] = &trade
	fmt.Printf("Trade %s executed between %s and %s.\n", trade.TradeID, trade.BuyerID, trade.SellerID)
	return nil
}

// RecordOrderCancellation records the cancellation of an order.
func (ledger *Ledger) RecordOrderCancellation(orderID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if _, exists := ledger.State.Marketplace.Orders[orderID]; !exists {
		return fmt.Errorf("order with ID %s not found", orderID)
	}

	delete(ledger.State.Marketplace.Orders, orderID)
	fmt.Printf("Order %s has been canceled.\n", orderID)
	return nil
}

// RecordResourceListing records the listing of a resource for sale or rental.
func (ledger *Ledger) RecordResourceListing(listing ResourceListing) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if ledger.State.Marketplace.ResourceListings == nil {
        ledger.State.Marketplace.ResourceListings = make(map[string]*ResourceListing)
    }

    ledger.State.Marketplace.ResourceListings[listing.ListingID] = &listing
    fmt.Printf("Resource %s listed by %s.\n", listing.ListingID, listing.OwnerID)
    return nil
}


// RecordResourceRental records the rental of a listed resource.
func (ledger *Ledger) RecordResourceRental(rental ResourceRental) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.ResourceRentals == nil {
		ledger.State.Marketplace.ResourceRentals = make(map[string]*ResourceRental)
	}

	ledger.State.Marketplace.ResourceRentals[rental.RentalID] = &rental
	fmt.Printf("Resource %s rented by %s.\n", rental.ListingID, rental.RenterID)
	return nil
}

// RecordResourcePurchase records the purchase of a listed resource.
func (ledger *Ledger) RecordResourcePurchase(purchase ResourcePurchase) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.ResourcePurchases == nil {
		ledger.State.Marketplace.ResourcePurchases = make(map[string]*ResourcePurchase)
	}

	ledger.State.Marketplace.ResourcePurchases[purchase.PurchaseID] = &purchase
	fmt.Printf("Resource %s purchased by %s.\n", purchase.ListingID, purchase.BuyerID)
	return nil
}

// RecordEscrowDispute records a dispute regarding an escrow transaction.
func (ledger *Ledger) RecordEscrowDispute(dispute EscrowDispute) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.EscrowDisputes == nil {
		ledger.State.Marketplace.EscrowDisputes = make(map[string]*EscrowDispute)
	}

	ledger.State.Marketplace.EscrowDisputes[dispute.DisputeID] = &dispute
	fmt.Printf("Escrow dispute %s recorded for transaction %s.\n", dispute.DisputeID, dispute.TransactionID)
	return nil
}

// RecordNewOrder records a new order in the marketplace.
func (ledger *Ledger) RecordNewOrder(order Order) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.Orders == nil {
		ledger.State.Marketplace.Orders = make(map[string]*Order)
	}

	ledger.State.Marketplace.Orders[order.OrderID] = &order
	fmt.Printf("New order %s created.\n", order.OrderID)
	return nil
}

// RecordTradeCompletion records the completion of a trade.
func (ledger *Ledger) RecordTradeCompletion(trade TradeCompletion) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.TradeCompletions == nil {
		ledger.State.Marketplace.TradeCompletions = make(map[string]*TradeCompletion)
	}

	ledger.State.Marketplace.TradeCompletions[trade.TradeID] = &trade
	fmt.Printf("Trade %s completed.\n", trade.TradeID)
	return nil
}

// RecordNewListing records a new item listing in the marketplace.
func (ledger *Ledger) RecordNewListing(listing ItemListing) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.ItemListings == nil {
		ledger.State.Marketplace.ItemListings = make(map[string]*ItemListing)
	}

	ledger.State.Marketplace.ItemListings[listing.ListingID] = &listing
	fmt.Printf("Item %s listed by %s.\n", listing.ListingID, listing.OwnerID)
	return nil
}

// RecordItemPurchase records the purchase of an item.
func (ledger *Ledger) RecordItemPurchase(purchase ItemPurchase) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.ItemPurchases == nil {
		ledger.State.Marketplace.ItemPurchases = make(map[string]*ItemPurchase)
	}

	ledger.State.Marketplace.ItemPurchases[purchase.PurchaseID] = &purchase
	fmt.Printf("Item %s purchased by %s.\n", purchase.ListingID, purchase.BuyerID)
	return nil
}

// RecordItemLease records the lease of an item.
func (ledger *Ledger) RecordItemLease(lease ItemLease) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.Marketplace.ItemLeases == nil {
		ledger.State.Marketplace.ItemLeases = make(map[string]*ItemLease)
	}

	ledger.State.Marketplace.ItemLeases[lease.LeaseID] = &lease
	fmt.Printf("Item %s leased by %s.\n", lease.ListingID, lease.LeaserID)
	return nil
}

// RecordEscrowRelease records the release of an escrow payment.
func (ledger *Ledger) RecordEscrowRelease(escrowID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	fmt.Printf("Escrow payment %s has been released.\n", escrowID)
	return nil
}

// RecordIllegalItemReport records the report of an illegal item listing.
func (ledger *Ledger) RecordIllegalItemReport(report IllegalItemReport) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if ledger.State.Marketplace.IllegalItemReports == nil {
        ledger.State.Marketplace.IllegalItemReports = make(map[string]IllegalItemReport)
    }

    ledger.State.Marketplace.IllegalItemReports[report.ReportID] = report // Store value, not pointer
    fmt.Printf("Illegal item report %s recorded for listing %s.\n", report.ReportID, report.ListingID)
    return nil
}


// RecordItemRemoval records the removal of an item listing.
func (ledger *Ledger) RecordItemRemoval(listingID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if _, exists := ledger.State.Marketplace.ItemListings[listingID]; !exists {
		return fmt.Errorf("listing with ID %s not found", listingID)
	}

	delete(ledger.State.Marketplace.ItemListings, listingID)
	fmt.Printf("Item listing %s has been removed.\n", listingID)
	return nil
}

// RecordNFTListing records a new NFT listing in the marketplace.
func (ledger *Ledger) RecordNFTListing(listing NFTListing) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.NFTMarketplace.Listings == nil {
		ledger.State.NFTMarketplace.Listings = make(map[string]*NFTListing)
	}

	ledger.State.NFTMarketplace.Listings[listing.ListingID] = &listing
	fmt.Printf("NFT %s listed by %s.\n", listing.ListingID, listing.OwnerID)
	return nil
}

// RecordNFTPurchase records the purchase of an NFT.
func (ledger *Ledger) RecordNFTPurchase(purchase NFTPurchase) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.NFTMarketplace.Purchases == nil {
		ledger.State.NFTMarketplace.Purchases = make(map[string]*NFTPurchase)
	}

	ledger.State.NFTMarketplace.Purchases[purchase.PurchaseID] = &purchase
	fmt.Printf("NFT %s purchased by %s.\n", purchase.ListingID, purchase.BuyerID)
	return nil
}

// RecordNFTListingCancellation records the cancellation of an NFT listing.
func (ledger *Ledger) RecordNFTListingCancellation(listingID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if _, exists := ledger.State.NFTMarketplace.Listings[listingID]; !exists {
		return fmt.Errorf("NFT listing with ID %s not found", listingID)
	}

	delete(ledger.State.NFTMarketplace.Listings, listingID)
	fmt.Printf("NFT listing %s has been canceled.\n", listingID)
	return nil
}

// RecordStakingPoolCreation records the creation of a staking pool.
func (ledger *Ledger) RecordStakingPoolCreation(poolID string, pool StakingPool) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	if ledger.State.StakingPools == nil {
		ledger.State.StakingPools = make(map[string]*StakingPool)
	}

	ledger.State.StakingPools[poolID] = &pool
	fmt.Printf("Staking pool %s created.\n", poolID)
	return nil
}


// RecordUnstakeTransaction records an unstaking transaction in a staking pool.
func (ledger *Ledger) RecordStakingPoolUnstakeTransaction(poolID string, unstake UnstakeTransaction) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	pool, exists := ledger.State.StakingPools[poolID]
	if !exists {
		return fmt.Errorf("staking pool with ID %s does not exist", poolID)
	}

	pool.Unstakes = append(pool.Unstakes, unstake)
	ledger.State.StakingPools[poolID] = pool
	fmt.Printf("Unstake transaction %s recorded in pool %s.\n", unstake.TransactionID, poolID)
	return nil
}

// RecordRewardDistribution records the distribution of staking rewards.
func (ledger *Ledger) RecordRewardDistribution(poolID string, reward RewardDistribution) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	pool, exists := ledger.State.StakingPools[poolID]
	if !exists {
		return fmt.Errorf("staking pool with ID %s does not exist", poolID)
	}

	pool.RewardDistributions = append(pool.RewardDistributions, reward)
	ledger.State.StakingPools[poolID] = pool
	fmt.Printf("Reward distribution recorded for pool %s.\n", poolID)
	return nil
}

// RecordStakingPoolClosure records the closure of a staking pool.
func (ledger *Ledger) RecordStakingPoolClosure(poolID string) error {
	ledger.lock.Lock()
	defer ledger.lock.Unlock()

	delete(ledger.State.StakingPools, poolID)
	fmt.Printf("Staking pool %s has been closed.\n", poolID)
	return nil
}

func (l *Ledger) recordMarketplaceInitialization(config AIMarketplaceConfig) error {
    for _, module := range config.InitialModules {
        l.AIModules[module.ID] = module
    }
    return nil
}

func (l *Ledger) recordAIModuleRegistration(module AIModule) error {
    if _, exists := l.AIModules[module.ID]; exists {
        return fmt.Errorf("AI module with ID %s already exists", module.ID)
    }
    l.AIModules[module.ID] = module
    return nil
}

func (l *Ledger) updateAIModule(module AIModule) error {
    if _, exists := l.AIModules[module.ID]; !exists {
        return fmt.Errorf("AI module with ID %s not found", module.ID)
    }
    l.AIModules[module.ID] = module
    return nil
}

func (l *Ledger) getAIModule(moduleID string) (AIModule, error) {
    module, exists := l.AIModules[moduleID]
    if !exists {
        return AIModule{}, fmt.Errorf("AI module with ID %s not found", moduleID)
    }
    return module, nil
}

func (l *Ledger) recordAIModuleRental(rental AIRental) error {
    l.AIModuleRentals[rental.ModuleID] = append(l.AIModuleRentals[rental.ModuleID], rental)
    return nil
}

func (l *Ledger) recordAIModuleUsage(log AILog) error {
    l.AIUsageLogs[log.ModuleID] = append(l.AIUsageLogs[log.ModuleID], log)
    return nil
}

func (l *Ledger) recordAIResourceRequest(request AIResourceRequest) error {
    l.AIResourceRequests[request.ModuleID] = append(l.AIResourceRequests[request.ModuleID], request)
    return nil
}

func (l *Ledger) recordAIModulePermissions(moduleID string, permissions AIPermission) error {
    module, exists := l.AIModules[moduleID]
    if !exists {
        return fmt.Errorf("AI module with ID %s not found", moduleID)
    }
    module.Permissions = permissions
    l.AIModules[moduleID] = module
    return nil
}

func (l *Ledger) recordAITransaction(transaction AITransaction) error {
    l.AITransactions[transaction.ModuleID] = append(l.AITransactions[transaction.ModuleID], transaction)
    return nil
}

func (l *Ledger) getAITransactionHistory(moduleID string) ([]AITransaction, error) {
    transactions, exists := l.AITransactions[moduleID]
    if !exists {
        return nil, fmt.Errorf("no transaction history found for module ID %s", moduleID)
    }
    return transactions, nil
}

func (l *Ledger) recordAIModuleLog(log AILog) error {
    l.AIUsageLogs[log.ModuleID] = append(l.AIUsageLogs[log.ModuleID], log)
    return nil
}

func (l *Ledger) getAIModuleUsage(moduleID string) ([]AILog, error) {
    logs, exists := l.AIUsageLogs[moduleID]
    if !exists {
        return nil, fmt.Errorf("no usage logs found for module ID %s", moduleID)
    }
    return logs, nil
}

func (l *Ledger) recordAIUsageSchedule(schedule AIUsageSchedule) error {
    l.AIUsageSchedules[schedule.ModuleID] = append(l.AIUsageSchedules[schedule.ModuleID], schedule)
    return nil
}

func (l *Ledger) recordAIEvent(event AIEventLog) error {
    l.AIEventLogs[event.ModuleID] = append(l.AIEventLogs[event.ModuleID], event)
    return nil
}

func (l *Ledger) updateNetworkUsageRequest(requestID, status string) error {
    request, exists := l.NetworkUsageRequests[requestID]
    if !exists {
        return fmt.Errorf("network usage request with ID %s not found", requestID)
    }
    request.Status = status
    l.NetworkUsageRequests[requestID] = request
    return nil
}

func (l *Ledger) recordAIResourceAllocation(allocation AIResourceAllocation) error {
    l.AIResourceAllocations[allocation.ModuleID] = allocation
    return nil
}

func (l *Ledger) recordAIModelPerformance(metrics AIMetrics) error {
    l.AIModelMetrics[metrics.ModuleID] = metrics
    return nil
}

func (l *Ledger) getAIModelMetrics(moduleID string) (AIMetrics, error) {
    metrics, exists := l.AIModelMetrics[moduleID]
    if !exists {
        return AIMetrics{}, fmt.Errorf("no metrics found for module ID %s", moduleID)
    }
    return metrics, nil
}

func (l *Ledger) updateAITrainingData(data AITrainingData) error {
    l.AITrainingDataRecords[data.ModuleID] = data
    return nil
}

func (l *Ledger) recordAIModelVersion(version AIModelVersion) error {
    l.AIModelVersionHistory[version.ModuleID] = append(l.AIModelVersionHistory[version.ModuleID], version)
    return nil
}

func (l *Ledger) recordAITask(task AITask) error {
    l.AITasks[task.TaskID] = task
    return nil
}

func (l *Ledger) getAITask(taskID string) (AITask, error) {
    task, exists := l.AITasks[taskID]
    if !exists {
        return AITask{}, fmt.Errorf("task with ID %s not found", taskID)
    }
    return task, nil
}

func (l *Ledger) recordAIReward(reward AIReward) error {
    l.AIRewards[reward.ModuleID] = append(l.AIRewards[reward.ModuleID], reward)
    return nil
}

func (l *Ledger) recordAIPenalty(penalty AIPenalty) error {
    l.AIPenalties[penalty.ModuleID] = append(l.AIPenalties[penalty.ModuleID], penalty)
    return nil
}

func (l *Ledger) recordAIDatasetLink(link AIDatasetLink) error {
    l.AIDatasetLinks[link.ModuleID] = append(l.AIDatasetLinks[link.ModuleID], link)
    return nil
}

func (l *Ledger) removeAIDatasetLink(moduleID, datasetID string) error {
    links, exists := l.AIDatasetLinks[moduleID]
    if !exists {
        return fmt.Errorf("no dataset links found for module ID %s", moduleID)
    }

    updatedLinks := []AIDatasetLink{}
    for _, link := range links {
        if link.DatasetID != datasetID {
            updatedLinks = append(updatedLinks, link)
        }
    }

    l.AIDatasetLinks[moduleID] = updatedLinks
    return nil
}

func (l *Ledger) recordLiquidityFeeAdjustment(fee LiquidityFee) error {
    l.LiquidityFees[fee.PairID] = fee
    return nil
}

func (l *Ledger) recordUserLiquidity(liquidity UserLiquidity) error {
    l.UserLiquidities[liquidity.UserID] = append(l.UserLiquidities[liquidity.UserID], liquidity)
    return nil
}

func (l *Ledger) getTradeHistory(pairID string, from, to time.Time) ([]TradeDetails, error) {
    // Implementation of filtering trades within the specified time range
    return []TradeDetails{}, nil
}

func (l *Ledger) recordSlippageTolerance(settings SlippageSettings) error {
    l.SlippageSettings[settings.PairID] = settings
    return nil
}

func (l *Ledger) recordTradeVolume(volume TradeVolume) error {
    l.TradeVolumes[volume.PairID] = append(l.TradeVolumes[volume.PairID], volume)
    return nil
}

func (l *Ledger) updateDEXTransactionStatus(transactionID, status string) error {
    l.DEXTransactions[transactionID] = status
    return nil
}

func (l *Ledger) getLiquidityLevel(pairID string) (float64, error) {
    // Return liquidity level for the pair
    return 0.0, nil
}

func (l *Ledger) recordTradeExpiry(expiry TradeExpiry) error {
    l.TradeExpiries[expiry.PairID] = expiry
    return nil
}

func (l *Ledger) recordPriceFluctuation(fluctuation PriceFluctuation) error {
    l.PriceFluctuations[fluctuation.PairID] = append(l.PriceFluctuations[fluctuation.PairID], fluctuation)
    return nil
}

func (l *Ledger) recordOrderBookDepth(depth OrderBookDepth) error {
    l.OrderBookDepths[depth.PairID] = depth
    return nil
}

func (l *Ledger) recordFeeStructure(feeStructure FeeStructure) error {
    l.FeeStructures[feeStructure.PairID] = feeStructure
    return nil
}

func (l *Ledger) getFeeHistory(pairID string, from, to time.Time) ([]FeeHistory, error) {
    // Implementation of filtering fee history within the specified time range
    return []FeeHistory{}, nil
}

func (l *Ledger) recordPoolTokenRatio(ratio PoolTokenRatio) error {
    l.PoolTokenRatios[ratio.PairID] = ratio
    return nil
}

func (l *Ledger) getSwapRate(pairID string, amount float64) (float64, error) {
    // Mock logic for swap rate calculation
    if ratio, exists := l.PoolTokenRatios[pairID]; exists {
        return amount * (ratio.TokenARatio / ratio.TokenBRatio), nil
    }
    return 0, fmt.Errorf("pair not found")
}

func (l *Ledger) setCrossPairTrading(pairID string, enabled bool) error {
    l.CrossPairTrading[pairID] = enabled
    return nil
}

func (l *Ledger) recordLiquidityProvision(provision LiquidityProvision) error {
    l.LiquidityProvisions[provision.PairID] = append(l.LiquidityProvisions[provision.PairID], provision)
    return nil
}

func (l *Ledger) getLiquidityYield(pairID string, period time.Duration) (float64, error) {
    // Logic to calculate yield based on past provisions or system settings
    if yield, exists := l.LiquidityYields[pairID]; exists && yield.Period == period {
        return yield.Yield, nil
    }
    return 0, fmt.Errorf("yield data not found for pair: %s", pairID)
}

func (l *Ledger) getLiquidityProvisions(pairID string, from, to time.Time) ([]LiquidityProvision, error) {
    provisions := []LiquidityProvision{}
    for _, provision := range l.LiquidityProvisions[pairID] {
        if provision.Timestamp.After(from) && provision.Timestamp.Before(to) {
            provisions = append(provisions, provision)
        }
    }
    return provisions, nil
}

func (l *Ledger) getLiquidityWithdrawals(pairID string, from, to time.Time) ([]LiquidityWithdrawal, error) {
    withdrawals := []LiquidityWithdrawal{}
    for _, withdrawal := range l.LiquidityWithdrawals[pairID] {
        if withdrawal.Timestamp.After(from) && withdrawal.Timestamp.Before(to) {
            withdrawals = append(withdrawals, withdrawal)
        }
    }
    return withdrawals, nil
}

func (l *Ledger) recordLiquidityWithdrawal(withdrawal LiquidityWithdrawal) error {
    l.LiquidityWithdrawals[withdrawal.PairID] = append(l.LiquidityWithdrawals[withdrawal.PairID], withdrawal)
    return nil
}

func (l *Ledger) updateLiquidityWithdrawalStatus(withdrawalID string, status string) error {
    // Logic to update status for a specific withdrawal request
    return nil
}

func (l *Ledger) recordDEXInitialization(config DEXConfig) error {
    l.DEXConfigurations[config.Name] = config
    return nil
}

func (l *Ledger) recordTradingPairCreation(pair TradingPair) error {
    l.TradingPairs[pair.PairID] = pair
    return nil
}

func (l *Ledger) recordLiquidityAddition(liquidity LiquidityPool) error {
    l.LiquidityPools[liquidity.PairID] = append(l.LiquidityPools[liquidity.PairID], liquidity)
    return nil
}

func (l *Ledger) recordLiquidityRemoval(liquidity LiquidityPool) error {
    l.LiquidityPools[liquidity.PairID] = append(l.LiquidityPools[liquidity.PairID], liquidity)
    return nil
}

func (l *Ledger) getSwapRate(pairID string, amount float64) (float64, error) {
    // Example logic to get a swap rate
    if pool, exists := l.LiquidityPools[pairID]; exists && len(pool) > 0 {
        return amount * 0.98, nil // Example: 2% fee applied
    }
    return 0, fmt.Errorf("pair not found")
}

func (l *Ledger) recordSwapExecution(swap Swap) error {
    l.TradeExecutions[swap.PairID] = append(l.TradeExecutions[swap.PairID], TradeExecution{
        TradeID:    fmt.Sprintf("%s_%d", swap.PairID, time.Now().UnixNano()),
        PairID:     swap.PairID,
        Amount:     swap.AmountIn,
        Price:      swap.AmountOut,
        ExecutedAt: swap.ExecutedAt,
    })
    return nil
}

func (l *Ledger) recordOrderPlacement(order Order) error {
    l.Orders[order.OrderID] = order
    return nil
}

func (l *Ledger) updateOrderStatus(orderID, status string) error {
    if order, exists := l.Orders[orderID]; exists {
        order.Status = status
        l.Orders[orderID] = order
        return nil
    }
    return fmt.Errorf("order not found")
}

func (l *Ledger) getOrderBook(pairID string) (OrderBook, error) {
    orders := []Order{}
    for _, order := range l.Orders {
        if order.PairID == pairID {
            orders = append(orders, order)
        }
    }
    return OrderBook{PairID: pairID, Orders: orders}, nil
}

func (l *Ledger) recordTradingFeeUpdate(fee TradingFee) error {
    l.TradingFees[fee.PairID] = fee
    return nil
}

func (l *Ledger) getTradeHistory(pairID string, from, to time.Time) ([]TradeExecution, error) {
    executions := []TradeExecution{}
    for _, execution := range l.TradeExecutions[pairID] {
        if execution.ExecutedAt.After(from) && execution.ExecutedAt.Before(to) {
            executions = append(executions, execution)
        }
    }
    return executions, nil
}

func (l *Ledger) getPriceImpact(pairID string, amount float64) (PriceImpact, error) {
    // Example logic to calculate price impact
    if _, exists := l.LiquidityPools[pairID]; exists {
        impact := PriceImpact{
            PairID:       pairID,
            TradeAmount:  amount,
            ImpactPercent: amount * 0.01, // Example: 1% per unit amount
        }
        return impact, nil
    }
    return PriceImpact{}, fmt.Errorf("pair not found")
}

func (l *Ledger) getOrderStatus(orderID string) (string, error) {
    status, exists := l.OrderStatuses[orderID]
    if !exists {
        return "", fmt.Errorf("order not found")
    }
    return status, nil
}

func (l *Ledger) recordOrderCancellation(cancellation OrderCancellation) error {
    if _, exists := l.OrderStatuses[cancellation.OrderID]; !exists {
        return fmt.Errorf("order not found")
    }
    l.OrderStatuses[cancellation.OrderID] = "Cancelled"
    return nil
}

func (l *Ledger) setTradingStatus(pairID string, status bool) error {
    if _, exists := l.LiquidityPoolInfo[pairID]; !exists {
        return fmt.Errorf("trading pair not found")
    }
    l.LiquidityPoolInfo[pairID].TradingEnabled = status
    return nil
}

func (l *Ledger) getLiquidityPoolInfo(pairID string) (LiquidityPool, error) {
    pool, exists := l.LiquidityPoolInfo[pairID]
    if !exists {
        return LiquidityPool{}, fmt.Errorf("liquidity pool not found")
    }
    return pool, nil
}

func (l *Ledger) recordMinimumTradeAmount(minTrade MinimumTradeAmount) error {
    l.MinimumTradeAmounts[minTrade.PairID] = minTrade
    return nil
}

func (l *Ledger) getLiquidityProvider(userID, pairID string) (bool, error) {
    providers, exists := l.LiquidityProviders[userID]
    if !exists {
        return false, fmt.Errorf("user not found")
    }
    verified, pairExists := providers[pairID]
    if !pairExists {
        return false, fmt.Errorf("user not a provider for the pair")
    }
    return verified, nil
}

func (l *Ledger) recordPoolRewardDistribution(reward PoolReward) error {
    l.PoolRewardDistributions[reward.PairID] = append(l.PoolRewardDistributions[reward.PairID], reward)
    return nil
}

func (l *Ledger) recordNFTTradeDenial(denial NFTTradeDenial) error {
    l.NFTTradeDenials[denial.TradeID] = denial
    return nil
}

func (l *Ledger) recordNFTTradeCompletion(trade NFTTrade) error {
    l.NFTTrades[trade.TradeID] = trade
    return nil
}

func (l *Ledger) setNFTExchangeStatus(status bool) error {
    l.NFTExchangeEnabled = status
    return nil
}

func (l *Ledger) recordNFTExchangeRate(rate NFTExchangeRate) error {
    l.NFTExchangeRates[rate.NFTID] = rate
    return nil
}

func (l *Ledger) recordNFTExchangeTransaction(transaction ExchangeTransaction) error {
    l.NFTExchangeTransactions[transaction.NFTID] = append(l.NFTExchangeTransactions[transaction.NFTID], transaction)
    return nil
}

func (l *Ledger) getNFTExchangeHistory(nftID string) ([]ExchangeTransaction, error) {
    transactions, exists := l.NFTExchangeTransactions[nftID]
    if !exists {
        return nil, fmt.Errorf("no exchange history for NFT ID %s", nftID)
    }
    return transactions, nil
}

func (l *Ledger) getNFTExchangeTransactions(from, to time.Time) ([]ExchangeTransaction, error) {
    var result []ExchangeTransaction
    for _, transactions := range l.NFTExchangeTransactions {
        for _, transaction := range transactions {
            if transaction.Timestamp.After(from) && transaction.Timestamp.Before(to) {
                result = append(result, transaction)
            }
        }
    }
    return result, nil
}

func (l *Ledger) recordNFTMintingLimit(limit NFTMintingLimit) error {
    l.NFTMintingLimits[limit.NFTID] = limit
    return nil
}

func (l *Ledger) recordNFTMintingEvent(event NFTMintingEvent) error {
    l.NFTMintingEvents[event.NFTID] = append(l.NFTMintingEvents[event.NFTID], event)
    return nil
}

func (l *Ledger) getNFTMintingHistory(nftID string, from, to time.Time) ([]NFTMintingEvent, error) {
    events, exists := l.NFTMintingEvents[nftID]
    if !exists {
        return nil, fmt.Errorf("no minting history for NFT ID %s", nftID)
    }
    var result []NFTMintingEvent
    for _, event := range events {
        if event.MintedAt.After(from) && event.MintedAt.Before(to) {
            result = append(result, event)
        }
    }
    return result, nil
}

// UpdateMintingAuthorizationStatus updates the status of a minting authorization
func (l *Ledger) UpdateMintingAuthorizationStatus(requestID, status string) error {
    if auth, exists := l.MintingAuthorizations[requestID]; exists {
        auth.Status = status
        auth.UpdatedAt = time.Now()
        l.MintingAuthorizations[requestID] = auth
        return nil
    }
    return fmt.Errorf("minting authorization with requestID %s not found", requestID)
}

// RecordMintingAuthorization records a minting authorization
func (l *Ledger) RecordMintingAuthorization(auth MintingAuthorization) error {
    l.MintingAuthorizations[auth.RequestID] = auth
    return nil
}

// UpdateNFTCustomizationStatus updates the customization status for an NFT
func (l *Ledger) UpdateNFTCustomizationStatus(nftID string, enabled bool) error {
    if _, exists := l.NFTCustomizationOptions[nftID]; exists {
        l.NFTStakeRewardsStatus[nftID] = enabled
        return nil
    }
    return fmt.Errorf("NFT with ID %s not found", nftID)
}

// RecordNFTCustomizationOptions records customization options for an NFT
func (l *Ledger) RecordNFTCustomizationOptions(options NFTCustomizationOptions) error {
    l.NFTCustomizationOptions[options.NFTID] = options
    return nil
}

// RecordNFTCustomizationHistory logs a customization in history
func (l *Ledger) RecordNFTCustomizationHistory(customization NFTCustomization) error {
    l.NFTCustomizationHistory[customization.NFTID] = append(l.NFTCustomizationHistory[customization.NFTID], customization)
    return nil
}

// RecordCustomizationEvent records a customization event
func (l *Ledger) RecordCustomizationEvent(event CustomizationEvent) error {
    l.CustomizationEvents[event.NFTID] = append(l.CustomizationEvents[event.NFTID], event)
    return nil
}

// SetNFTStakeRewardsStatus sets staking rewards for an NFT
func (l *Ledger) SetNFTStakeRewardsStatus(nftID string, enabled bool) error {
    l.NFTStakeRewardsStatus[nftID] = enabled
    return nil
}

// GetNFTYieldRate retrieves the staking yield rate for an NFT
func (l *Ledger) GetNFTYieldRate(nftID string) (float64, error) {
    if rate, exists := l.NFTYieldRates[nftID]; exists {
        return rate, nil
    }
    return 0, fmt.Errorf("yield rate not found for NFT %s", nftID)
}

// RecordStakeRewardDistribution records the distribution of stake rewards
func (l *Ledger) RecordStakeRewardDistribution(reward StakeReward) error {
    l.StakeRewards[reward.NFTID] = append(l.StakeRewards[reward.NFTID], reward)
    return nil
}

// GetStakeRewards retrieves stake rewards for an NFT within a time range
func (l *Ledger) GetStakeRewards(nftID string, from, to time.Time) ([]StakeReward, error) {
    var rewards []StakeReward
    for _, reward := range l.StakeRewards[nftID] {
        if reward.DistributedAt.After(from) && reward.DistributedAt.Before(to) {
            rewards = append(rewards, reward)
        }
    }
    return rewards, nil
}

// GetStakeDistribution retrieves stake distribution for an NFT
func (l *Ledger) GetStakeDistribution(nftID string) (StakeDistribution, error) {
    if distribution, exists := l.StakeDistributions[nftID]; exists {
        return distribution, nil
    }
    return StakeDistribution{}, fmt.Errorf("stake distribution for NFT %s not found", nftID)
}

// SetCrossMarketplaceTradeStatus sets the status of cross-marketplace trading
func (l *Ledger) SetCrossMarketplaceTradeStatus(enabled bool) error {
    l.CrossMarketplaceTradeEnabled = enabled
    return nil
}

// RecordCrossMarketplaceRate records exchange rates for cross-marketplace NFT trades
func (l *Ledger) RecordCrossMarketplaceRate(rate CrossMarketplaceRate) error {
    l.CrossMarketplaceRates[rate.NFTID] = rate
    return nil
}

// UpdateCrossMarketplaceListingStatus updates the listing status of an NFT across marketplaces
func (l *Ledger) UpdateCrossMarketplaceListingStatus(nftID, status string) error {
    l.CrossMarketplaceStatuses[nftID] = CrossMarketplaceStatus{
        NFTID:   nftID,
        Listed:  status == "Approved",
        Details: fmt.Sprintf("Status updated to %s", status),
    }
    return nil
}

// RecordCrossMarketplaceTrade records a trade across multiple marketplaces
func (l *Ledger) RecordCrossMarketplaceTrade(trade CrossMarketplaceTrade) error {
    l.CrossMarketplaceTrades[trade.NFTID] = append(l.CrossMarketplaceTrades[trade.NFTID], trade)
    return nil
}

// RecordCrossMarketplaceMetrics records metrics for cross-marketplace trades
func (l *Ledger) RecordCrossMarketplaceMetrics(metrics CrossMarketplaceMetrics) error {
    l.CrossMarketplaceMetrics[metrics.NFTID] = append(l.CrossMarketplaceMetrics[metrics.NFTID], metrics)
    return nil
}

// GetCrossMarketplaceStatus retrieves cross-marketplace listing status of an NFT
func (l *Ledger) GetCrossMarketplaceStatus(nftID string) (CrossMarketplaceStatus, error) {
    if status, exists := l.CrossMarketplaceStatuses[nftID]; exists {
        return status, nil
    }
    return CrossMarketplaceStatus{}, fmt.Errorf("cross-marketplace status for NFT %s not found", nftID)
}

// SetUserRatingSystemStatus enables or disables the user rating system
func (l *Ledger) SetUserRatingSystemStatus(enabled bool) error {
    l.UserRatingSystemEnabled = enabled
    return nil
}

// RecordUserRating records a user rating for an NFT
func (l *Ledger) RecordUserRating(rating UserRating) error {
    l.UserRatings[rating.NFTID] = append(l.UserRatings[rating.NFTID], rating)
    return nil
}

// RecordUserFeedback records user feedback for an NFT transaction
func (l *Ledger) RecordUserFeedback(feedback UserFeedback) error {
    l.UserFeedback[feedback.NFTID] = append(l.UserFeedback[feedback.NFTID], feedback)
    return nil
}

// GetRatingSummary generates a summary of user ratings for an NFT
func (l *Ledger) GetRatingSummary(nftID string) (RatingSummary, error) {
    ratings, exists := l.UserRatings[nftID]
    if !exists || len(ratings) == 0 {
        return RatingSummary{}, fmt.Errorf("no ratings found for NFT %s", nftID)
    }

    var total int
    for _, rating := range ratings {
        total += rating.Rating
    }

    avg := float64(total) / float64(len(ratings))
    return RatingSummary{
        NFTID:         nftID,
        AverageRating: avg,
        TotalRatings:  len(ratings),
    }, nil
}

// RecordUserRatingActivity records a user's activity related to NFT ratings
func (l *Ledger) RecordUserRatingActivity(activity RatingActivity) error {
    l.RatingActivities[activity.UserID] = append(l.RatingActivities[activity.UserID], activity)
    return nil
}

// SetNFTInheritanceStatus enables or disables inheritance for an NFT
func (l *Ledger) SetNFTInheritanceStatus(nftID string, enabled bool) error {
    if l.NFTInheritanceEnabled == nil {
        l.NFTInheritanceEnabled = make(map[string]bool)
    }
    l.NFTInheritanceEnabled[nftID] = enabled
    return nil
}

// RecordNFTInheritanceRights records inheritance rights for an NFT
func (l *Ledger) RecordNFTInheritanceRights(inheritance NFTInheritanceRights) error {
    l.NFTInheritanceRights[inheritance.NFTID] = inheritance
    return nil
}

// GetNFTInheritanceRights retrieves the inheritance rights of an NFT
func (l *Ledger) GetNFTInheritanceRights(nftID string) (NFTInheritanceRights, error) {
    if rights, exists := l.NFTInheritanceRights[nftID]; exists {
        return rights, nil
    }
    return NFTInheritanceRights{}, fmt.Errorf("inheritance rights for NFT %s not found", nftID)
}

// RecordInheritanceActivity logs an inheritance activity
func (l *Ledger) RecordInheritanceActivity(activity InheritanceActivity) error {
    l.InheritanceActivities[activity.NFTID] = append(l.InheritanceActivities[activity.NFTID], activity)
    return nil
}

// GetInheritanceActivities retrieves inheritance activities for an NFT within a time range
func (l *Ledger) GetInheritanceActivities(nftID string, from, to time.Time) ([]InheritanceActivity, error) {
    var activities []InheritanceActivity
    for _, activity := range l.InheritanceActivities[nftID] {
        if activity.Timestamp.After(from) && activity.Timestamp.Before(to) {
            activities = append(activities, activity)
        }
    }
    return activities, nil
}

// SetNFTBundleListingStatus enables or disables NFT bundle listing
func (l *Ledger) SetNFTBundleListingStatus(enabled bool) error {
    l.NFTBundleListingEnabled = enabled
    return nil
}

// RecordNFTBundleCreation records the creation of an NFT bundle
func (l *Ledger) RecordNFTBundleCreation(bundle NFTBundle) error {
    l.NFTBundles[bundle.BundleID] = bundle
    return nil
}

// AddNFTToBundle adds an NFT to an existing bundle
func (l *Ledger) AddNFTToBundle(entry NFTBundleEntry) error {
    bundle, exists := l.NFTBundles[entry.BundleID]
    if !exists {
        return fmt.Errorf("bundle with ID %s not found", entry.BundleID)
    }
    bundle.NFTs = append(bundle.NFTs, entry.NFTID)
    l.NFTBundles[entry.BundleID] = bundle
    return nil
}

// RecordCollectionActivity logs an activity for a specific NFT collection
func (l *Ledger) RecordCollectionActivity(activity CollectionActivity) error {
    l.CollectionActivities[activity.CollectionID] = append(l.CollectionActivities[activity.CollectionID], activity)
    return nil
}

// GetCollectionActivities retrieves activities for a specific NFT collection within a time range
func (l *Ledger) GetCollectionActivities(collectionID string, from, to time.Time) ([]CollectionActivity, error) {
    var activities []CollectionActivity
    for _, activity := range l.CollectionActivities[collectionID] {
        if activity.Timestamp.After(from) && activity.Timestamp.Before(to) {
            activities = append(activities, activity)
        }
    }
    return activities, nil
}

// SetNFTTradingStatus enables or disables NFT trading
func (l *Ledger) SetNFTTradingStatus(enabled bool) error {
    l.NFTTradingEnabled = enabled
    return nil
}

// RecordNFTTradeEvent logs a trade event for an NFT
func (l *Ledger) RecordNFTTradeEvent(trade NFTTradeEvent) error {
    l.NFTTradeEvents[trade.TradeID] = trade
    return nil
}

// UpdateNFTTradeStatus updates the status of an NFT trade
func (l *Ledger) UpdateNFTTradeStatus(tradeID, status string) error {
    if _, exists := l.NFTTradeStatuses[tradeID]; exists {
        l.NFTTradeStatuses[tradeID] = NFTTradeStatus{
            TradeID: tradeID,
            Status:  status,
            UpdatedAt: time.Now(),
        }
        return nil
    }
    return fmt.Errorf("trade with ID %s not found", tradeID)
}

// RecordNFTMarketplaceInitialization records the initialization of an NFT marketplace
func (l *Ledger) RecordNFTMarketplaceInitialization(config NFTMarketplaceConfig) error {
    l.NFTMarketplaceConfigs[config.MarketplaceName] = config
    return nil
}

// RecordNFTMinting records the minting of a new NFT
func (l *Ledger) RecordNFTMinting(nft NFT) error {
    l.NFTs[nft.NFTID] = nft
    return nil
}

// RecordNFTBurn records the burning of an NFT
func (l *Ledger) RecordNFTBurn(burnRecord NFTBurnRecord) error {
    l.BurnedNFTs[burnRecord.NFTID] = burnRecord
    return nil
}

// RecordNFTOwnershipTransfer records the ownership transfer of an NFT
func (l *Ledger) RecordNFTOwnershipTransfer(transfer NFTOwnershipTransfer) error {
    l.NFTOwnershipTransfers[transfer.NFTID] = append(l.NFTOwnershipTransfers[transfer.NFTID], transfer)
    return nil
}

// SetNFTTransferApproval sets the approval status for an NFT transfer
func (l *Ledger) SetNFTTransferApproval(nftID, approvedOwnerID string, approved bool) error {
    // Simulate a transfer approval field
    return nil
}

// RecordNFTListingForSale records an NFT for sale in the marketplace
func (l *Ledger) RecordNFTListingForSale(sale NFTSale) error {
    l.NFTSales[sale.NFTID] = sale
    return nil
}

// RemoveNFTFromSale removes an NFT from the marketplace sale listing
func (l *Ledger) RemoveNFTFromSale(nftID string) error {
    delete(l.NFTSales, nftID)
    return nil
}

// UpdateNFTSalePrice updates the price of an NFT listed for sale
func (l *Ledger) UpdateNFTSalePrice(sale NFTSale) error {
    if existingSale, exists := l.NFTSales[sale.NFTID]; exists {
        existingSale.Price = sale.Price
        existingSale.UpdatedAt = sale.UpdatedAt
        l.NFTSales[sale.NFTID] = existingSale
        return nil
    }
    return fmt.Errorf("NFT with ID %s not listed for sale", sale.NFTID)
}

// RecordNFTBid records a bid for an NFT
func (l *Ledger) RecordNFTBid(bid NFTBid) error {
    l.NFTBids[bid.NFTID] = append(l.NFTBids[bid.NFTID], bid)
    return nil
}

// UpdateNFTBidStatus updates the status of an NFT bid
func (l *Ledger) UpdateNFTBidStatus(bidID, status string) error {
    for _, bids := range l.NFTBids {
        for i, bid := range bids {
            if bid.BidID == bidID {
                bid.Timestamp = time.Now()
                l.NFTBids[bid.NFTID][i] = bid
                return nil
            }
        }
    }
    return fmt.Errorf("bid with ID %s not found", bidID)
}

// RecordNFTAuctionStart records the start of an NFT auction
func (l *Ledger) RecordNFTAuctionStart(auction NFTAuction) error {
    l.NFTAuctions[auction.AuctionID] = auction
    return nil
}

// RecordNFTAuctionEnd records the conclusion of an NFT auction
func (l *Ledger) RecordNFTAuctionEnd(auctionEnd NFTAuctionEnd) error {
    delete(l.NFTAuctions, auctionEnd.AuctionID)
    return nil
}

// GetNFTAuctionStatus retrieves the current status of an NFT auction
func (l *Ledger) GetNFTAuctionStatus(auctionID string) (NFTAuctionStatus, error) {
    if status, exists := l.NFTAuctionStatuses[auctionID]; exists {
        return status, nil
    }
    return NFTAuctionStatus{}, fmt.Errorf("auction with ID %s not found", auctionID)
}

// RecordNFTAuctionEvent logs an event for a specific NFT auction
func (l *Ledger) RecordNFTAuctionEvent(event NFTAuctionEvent) error {
    l.NFTAuctionEvents[event.AuctionID] = append(l.NFTAuctionEvents[event.AuctionID], event)
    return nil
}

// GetNFTAuctionEvents retrieves auction events within a specified time range
func (l *Ledger) GetNFTAuctionEvents(auctionID string, from, to time.Time) ([]NFTAuctionEvent, error) {
    var events []NFTAuctionEvent
    for _, event := range l.NFTAuctionEvents[auctionID] {
        if event.Timestamp.After(from) && event.Timestamp.Before(to) {
            events = append(events, event)
        }
    }
    return events, nil
}

// GetNFTAuctionBids retrieves bids for an auction within a specified time range
func (l *Ledger) GetNFTAuctionBids(auctionID string, from, to time.Time) ([]NFTBid, error) {
    var bids []NFTBid
    for _, bid := range l.NFTBids[auctionID] {
        if bid.Timestamp.After(from) && bid.Timestamp.Before(to) {
            bids = append(bids, bid)
        }
    }
    return bids, nil
}

// GetNFTOwner retrieves the current owner of an NFT
func (l *Ledger) GetNFTOwner(nftID string) (NFTOwnership, error) {
    if owner, exists := l.NFTOwnerships[nftID]; exists {
        return owner, nil
    }
    return NFTOwnership{}, fmt.Errorf("ownership record for NFT %s not found", nftID)
}

// GetNFTMetadata retrieves the metadata of an NFT
func (l *Ledger) GetNFTMetadata(nftID string) (NFTMetadata, error) {
    if metadata, exists := l.NFTMetadata[nftID]; exists {
        return metadata, nil
    }
    return NFTMetadata{}, fmt.Errorf("metadata for NFT %s not found", nftID)
}

// UpdateNFTMetadata updates the metadata of an NFT
func (l *Ledger) UpdateNFTMetadata(metadata NFTMetadata) error {
    l.NFTMetadata[metadata.NFTID] = metadata
    return nil
}

// GetNFTAuthenticity retrieves the authenticity of an NFT
func (l *Ledger) GetNFTAuthenticity(nftID string) (NFTAuthenticity, error) {
    if authenticity, exists := l.NFTAuthenticityRecords[nftID]; exists {
        return authenticity, nil
    }
    return NFTAuthenticity{}, fmt.Errorf("authenticity record for NFT %s not found", nftID)
}

// GetNFTOwnershipHistory retrieves the ownership history of an NFT
func (l *Ledger) GetNFTOwnershipHistory(nftID string) ([]NFTOwnershipHistory, error) {
    return l.NFTOwnershipHistories[nftID], nil
}

// RecordNFTTransferEvent logs a transfer event for an NFT
func (l *Ledger) RecordNFTTransferEvent(event NFTTransferEvent) error {
    l.NFTTransferEvents[event.NFTID] = append(l.NFTTransferEvents[event.NFTID], event)
    return nil
}

// UpdateNFTListingStatus updates the status of an NFT listing
func (l *Ledger) UpdateNFTListingStatus(nftID, status string) error {
    if listing, exists := l.NFTListings[nftID]; exists {
        listing.Status = status
        l.NFTListings[nftID] = listing
        return nil
    }
    return fmt.Errorf("listing for NFT %s not found", nftID)
}

// RecordNFTListing records a new listing for an NFT
func (l *Ledger) RecordNFTListing(listing NFTListing) error {
    l.NFTListings[listing.NFTID] = listing
    return nil
}

// SetNFTStakingStatus enables or disables staking for an NFT
func (l *Ledger) SetNFTStakingStatus(nftID string, enabled bool) error {
    if _, exists := l.NFTStakings[nftID]; !exists && enabled {
        l.NFTStakings[nftID] = NFTStaking{}
    }
    return nil
}

// RecordNFTStake logs staking details for an NFT
func (l *Ledger) RecordNFTStake(stake NFTStaking) error {
    l.NFTStakings[stake.NFTID] = stake
    return nil
}

// RecordNFTUnstake logs unstaking details for an NFT
func (l *Ledger) RecordNFTUnstake(unstake NFTUnstake) error {
    l.NFTUnstakes[unstake.NFTID] = append(l.NFTUnstakes[unstake.NFTID], unstake)
    return nil
}

// RecordNFTStakingEvent logs an event related to NFT staking
func (l *Ledger) RecordNFTStakingEvent(event NFTStakingEvent) error {
    l.NFTStakingEvents[event.NFTID] = append(l.NFTStakingEvents[event.NFTID], event)
    return nil
}

// RecordNFTRoyalty records royalty information for an NFT
func (l *Ledger) RecordNFTRoyalty(royalty NFTRoyalty) error {
    l.NFTRoyalties[royalty.NFTID] = royalty
    return nil
}

// GetNFTRoyalty retrieves royalty information for an NFT
func (l *Ledger) GetNFTRoyalty(nftID string) (NFTRoyalty, error) {
    if royalty, exists := l.NFTRoyalties[nftID]; exists {
        return royalty, nil
    }
    return NFTRoyalty{}, fmt.Errorf("royalty record for NFT %s not found", nftID)
}

// RecordRoyaltyDistribution logs the distribution of royalties for an NFT
func (l *Ledger) RecordRoyaltyDistribution(distribution RoyaltyDistribution) error {
    l.RoyaltyDistributions[distribution.NFTID] = append(l.RoyaltyDistributions[distribution.NFTID], distribution)
    return nil
}

// GetRoyaltyDistributions retrieves royalty distributions for an NFT within a specified time range
func (l *Ledger) GetRoyaltyDistributions(nftID string, from, to time.Time) ([]RoyaltyDistribution, error) {
    var distributions []RoyaltyDistribution
    for _, distribution := range l.RoyaltyDistributions[nftID] {
        if distribution.DistributedAt.After(from) && distribution.DistributedAt.Before(to) {
            distributions = append(distributions, distribution)
        }
    }
    return distributions, nil
}

// GetTotalRoyaltyDistribution retrieves the total royalties distributed for an NFT
func (l *Ledger) GetTotalRoyaltyDistribution(nftID string) (float64, error) {
    if total, exists := l.TotalRoyaltyDistributions[nftID]; exists {
        return total, nil
    }
    return 0, fmt.Errorf("no royalty distributions found for NFT %s", nftID)
}

// SetFractionalOwnershipStatus enables or disables fractional ownership for an NFT
func (l *Ledger) SetFractionalOwnershipStatus(nftID string, enabled bool) error {
    l.FractionalOwnershipEnabled[nftID] = enabled
    return nil
}

// GetFractionalOwnershipDetails retrieves fractional ownership details for an NFT
func (l *Ledger) GetFractionalOwnershipDetails(nftID string) ([]FractionalOwnership, error) {
    return l.FractionalOwnerships[nftID], nil
}

// RecordFractionalOwnershipChange logs changes to fractional ownership of an NFT
func (l *Ledger) RecordFractionalOwnershipChange(change FractionalOwnershipChange) error {
    l.FractionalOwnershipChanges[change.NFTID] = append(l.FractionalOwnershipChanges[change.NFTID], change)
    return nil
}

// SetNFTEscrowStatus enables or disables escrow for an NFT
func (l *Ledger) SetNFTEscrowStatus(nftID string, enabled bool) error {
    l.NFTEscrowEnabled[nftID] = enabled
    return nil
}

// RecordEscrowRelease logs the release of an NFT from escrow
func (l *Ledger) RecordEscrowRelease(release NFTEscrowRelease) error {
    l.NFTEscrows[release.EscrowID] = release
    return nil
}

func (l *Ledger) RecordStoreCreation(storeID, owner, name, category string) error {
    if storeID == "" || owner == "" || name == "" || category == "" {
        return fmt.Errorf("all parameters are required for store creation logging")
    }

    logEntry := fmt.Sprintf("Store Created: ID=%s, Owner=%s, Name=%s, Category=%s", storeID, owner, name, category)
    err := l.Log(logEntry)
    if err != nil {
        return fmt.Errorf("failed to log store creation: %w", err)
    }

    log.Printf("[LEDGER] Store creation logged successfully: %s", logEntry)
    return nil
}
