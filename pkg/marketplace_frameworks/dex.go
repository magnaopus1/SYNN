package marketplace

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)


// NewDEXManager initializes a new decentralized exchange (DEX) manager
func NewDEXManager(ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption) *common.DEXManager {
	return &common.DEXManager{
		Orders:            make(map[string]*common.Order),
		CompletedOrders:   make(map[string]*common.Order),
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// PlaceOrder places a new order for trading
func (dex *common.DEXManager) PlaceOrder(orderID, trader, assetIn, assetOut string, amountIn, amountOut float64, orderType string) (*common.Order, error) {
	dex.mu.Lock()
	defer dex.mu.Unlock()

	// Encrypt order data
	orderData := fmt.Sprintf("OrderID: %s, Trader: %s, AssetIn: %s, AssetOut: %s, AmountIn: %f, AmountOut: %f, OrderType: %s", orderID, trader, assetIn, assetOut, amountIn, amountOut, orderType)
	encryptedData, err := dex.EncryptionService.EncryptData([]byte(orderData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt order data: %v", err)
	}

	// Create the new order
	order := &common.Order{
		OrderID:   orderID,
		Trader:    trader,
		AssetIn:   assetIn,
		AmountIn:  amountIn,
		AssetOut:  assetOut,
		AmountOut: amountOut,
		OrderType: orderType,
		OrderTime: time.Now(),
		IsFilled:  false,
	}

	// Add the order to the DEX
	dex.Orders[orderID] = order

	// Log the order in the ledger
	err = dex.Ledger.RecordNewOrder(orderID, trader, assetIn, assetOut, amountIn, amountOut, orderType, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log order in ledger: %v", err)
	}

	fmt.Printf("Order %s placed by trader %s to %s %f of %s for %f of %s\n", orderID, trader, orderType, amountIn, assetIn, amountOut, assetOut)
	return order, nil
}

// MatchOrders attempts to match a buy order with a sell order and completes the trade
func (dex *common.DEXManager) MatchOrders(buyOrderID, sellOrderID string) (string, error) {
	dex.mu.Lock()
	defer dex.mu.Unlock()

	// Retrieve the buy and sell orders
	buyOrder, buyExists := dex.Orders[buyOrderID]
	sellOrder, sellExists := dex.Orders[sellOrderID]
	if !buyExists || !sellExists {
		return "", errors.New("both buy and sell orders must exist")
	}

	if buyOrder.IsFilled || sellOrder.IsFilled {
		return "", errors.New("one of the orders is already filled")
	}

	// Ensure the asset pairs and amounts match
	if buyOrder.AssetOut != sellOrder.AssetIn || buyOrder.AmountOut != sellOrder.AmountIn {
		return "", errors.New("order assets or amounts do not match")
	}

	// Create a transaction ID for the completed trade
	transactionID := generateUniqueID()

	// Mark both orders as filled and associate the transaction ID
	buyOrder.IsFilled = true
	buyOrder.TransactionID = transactionID
	sellOrder.IsFilled = true
	sellOrder.TransactionID = transactionID

	// Move orders to the completed orders pool
	dex.CompletedOrders[buyOrderID] = buyOrder
	dex.CompletedOrders[sellOrderID] = sellOrder

	// Remove them from active orders
	delete(dex.Orders, buyOrderID)
	delete(dex.Orders, sellOrderID)

	// Log the completed trade in the ledger
	err := dex.Ledger.RecordTradeCompletion(transactionID, buyOrderID, sellOrderID, buyOrder.Trader, sellOrder.Trader, buyOrder.AssetIn, buyOrder.AmountIn, sellOrder.AssetIn, sellOrder.AmountIn, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to log trade completion: %v", err)
	}

	fmt.Printf("Trade completed between buy order %s and sell order %s. Transaction ID: %s\n", buyOrderID, sellOrderID, transactionID)
	return transactionID, nil
}

// CancelOrder allows a user to cancel an order if it hasn't been filled yet
func (dex *common.DEXManager) CancelOrder(orderID string) error {
	dex.mu.Lock()
	defer dex.mu.Unlock()

	// Retrieve the order
	order, exists := dex.Orders[orderID]
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}

	if order.IsFilled {
		return fmt.Errorf("order %s is already filled and cannot be canceled", orderID)
	}

	// Remove the order from the active orders
	delete(dex.Orders, orderID)

	// Log the order cancellation in the ledger
	err := dex.Ledger.RecordOrderCancellation(orderID, order.Trader, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log order cancellation: %v", err)
	}

	fmt.Printf("Order %s canceled by trader %s\n", orderID, order.Trader)
	return nil
}

// GetOrderDetails retrieves the details of a specific order
func (dex *common.DEXManager) GetOrderDetails(orderID string) (*common.Order, error) {
	dex.mu.Lock()
	defer dex.mu.Unlock()

	order, exists := dex.Orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order %s not found", orderID)
	}

	return order, nil
}

// GetCompletedOrderDetails retrieves the details of a completed order
func (dex *common.DEXManager) GetCompletedOrderDetails(orderID string) (*common.Order, error) {
	dex.mu.Lock()
	defer dex.mu.Unlock()

	order, exists := dex.CompletedOrders[orderID]
	if !exists {
		return nil, fmt.Errorf("completed order %s not found", orderID)
	}

	return order, nil
}

