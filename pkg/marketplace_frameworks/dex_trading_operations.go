package marketplace

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

func initializeDEX(config DEXConfig, ledgerInstance *Ledger) error {
    config.InitializedAt = time.Now()
    return ledgerInstance.recordDEXInitialization(config)
}

func createTradingPair(pair TradingPair, ledgerInstance *Ledger) error {
    pair.CreatedAt = time.Now()
    return ledgerInstance.recordTradingPairCreation(pair)
}

func addLiquidityToPool(providerID, pairID string, amountA, amountB float64, ledgerInstance *Ledger) error {
    liquidity := LiquidityPool{
        ProviderID: providerID,
        PairID:     pairID,
        TokenA:     amountA,
        TokenB:     amountB,
        AddedAt:    time.Now(),
    }
    return ledgerInstance.recordLiquidityAddition(liquidity)
}

func removeLiquidityFromPool(providerID, pairID string, amountA, amountB float64, ledgerInstance *Ledger) error {
    liquidity := LiquidityPool{
        ProviderID: providerID,
        PairID:     pairID,
        TokenA:     -amountA,
        TokenB:     -amountB,
        RemovedAt:  time.Now(),
    }
    return ledgerInstance.recordLiquidityRemoval(liquidity)
}

func executeSwap(pairID string, amountIn float64, tokenIn, tokenOut string, ledgerInstance *Ledger) (float64, error) {
    swapRate, err := ledgerInstance.getSwapRate(pairID, amountIn)
    if err != nil {
        return 0, fmt.Errorf("failed to retrieve swap rate: %v", err)
    }
    amountOut := amountIn * swapRate
    swap := Swap{
        PairID:     pairID,
        AmountIn:   amountIn,
        TokenIn:    tokenIn,
        AmountOut:  amountOut,
        TokenOut:   tokenOut,
        ExecutedAt: time.Now(),
    }
    return amountOut, ledgerInstance.recordSwapExecution(swap)
}

func placeLimitOrder(order Order, ledgerInstance *Ledger) error {
    order.Type = "Limit"
    order.PlacedAt = time.Now()
    return ledgerInstance.recordOrderPlacement(order)
}

func cancelLimitOrder(orderID string, ledgerInstance *Ledger) error {
    return ledgerInstance.updateOrderStatus(orderID, "Cancelled")
}

func getOrderBook(pairID string, ledgerInstance *Ledger) (OrderBook, error) {
    return ledgerInstance.getOrderBook(pairID)
}

func updateTradingFees(pairID string, feeRate float64, ledgerInstance *Ledger) error {
    fee := TradingFee{
        PairID:    pairID,
        FeeRate:   feeRate,
        Timestamp: time.Now(),
    }
    return ledgerInstance.recordTradingFeeUpdate(fee)
}

func retrieveTradeHistory(pairID string, from, to time.Time, ledgerInstance *Ledger) ([]TradeExecution, error) {
    return ledgerInstance.getTradeHistory(pairID, from, to)
}

func calculatePriceImpact(pairID string, amount float64, ledgerInstance *Ledger) (PriceImpact, error) {
    return ledgerInstance.getPriceImpact(pairID, amount)
}


// CheckOrderStatus checks the status of a specific order.
func checkOrderStatus(orderID string, ledgerInstance *Ledger) (string, error) {
    return ledgerInstance.getOrderStatus(orderID)
}

// LogOrderCancellation records the cancellation of a specific order.
func logOrderCancellation(orderID string, ledgerInstance *Ledger) error {
    cancellation := OrderCancellation{
        OrderID:     orderID,
        CancelledAt: time.Now(),
    }
    return ledgerInstance.recordOrderCancellation(cancellation)
}

// EnableTokenPairTrading enables trading for a specific token pair.
func enableTokenPairTrading(pairID string, ledgerInstance *Ledger) error {
    return ledgerInstance.setTradingStatus(pairID, true)
}

// DisableTokenPairTrading disables trading for a specific token pair.
func disableTokenPairTrading(pairID string, ledgerInstance *Ledger) error {
    return ledgerInstance.setTradingStatus(pairID, false)
}

// GetLiquidityPoolInfo retrieves information about a specific liquidity pool.
func getLiquidityPoolInfo(pairID string, ledgerInstance *Ledger) (LiquidityPool, error) {
    return ledgerInstance.getLiquidityPoolInfo(pairID)
}

// SetMinimumTradeAmount sets the minimum trade amount for a specific pair.
func setMinimumTradeAmount(pairID string, amount float64, ledgerInstance *Ledger) error {
    minTrade := MinimumTradeAmount{
        PairID: pairID,
        Amount: amount,
        SetAt:  time.Now(),
    }
    return ledgerInstance.recordMinimumTradeAmount(minTrade)
}

// VerifyLiquidityProvider verifies that a user is a registered liquidity provider.
func verifyLiquidityProvider(userID, pairID string, ledgerInstance *Ledger) (bool, error) {
    return ledgerInstance.getLiquidityProvider(userID, pairID)
}

// DistributePoolRewards distributes rewards to liquidity providers based on their contributions.
func distributePoolRewards(pairID string, reward PoolReward, ledgerInstance *Ledger) error {
    reward.PairID = pairID
    reward.DistributedAt = time.Now()
    return ledgerInstance.recordPoolRewardDistribution(reward)
}
