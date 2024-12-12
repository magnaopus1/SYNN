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

const (
	OrderBookManipulationCheckInterval = 10 * time.Second // Interval for checking DEX order book activities
	MaxOrderCancelsPerMinute           = 50               // Maximum number of order cancellations allowed per user per minute
	MinOrderLifetime                   = 1 * time.Minute  // Minimum time an order must stay on the order book before cancellation
)

// DexOrderBookManipulationRestrictionAutomation monitors and restricts potential manipulation of the decentralized exchange order book
type DexOrderBookManipulationRestrictionAutomation struct {
	consensusSystem           *consensus.SynnergyConsensus
	ledgerInstance            *ledger.Ledger
	stateMutex                *sync.RWMutex
	userOrderCancelCount      map[string]int // Tracks the number of order cancellations per user
	userLastOrderTimestamps   map[string]time.Time // Tracks the time when orders were placed by users
}

// NewDexOrderBookManipulationRestrictionAutomation initializes and returns an instance of DexOrderBookManipulationRestrictionAutomation
func NewDexOrderBookManipulationRestrictionAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *DexOrderBookManipulationRestrictionAutomation {
	return &DexOrderBookManipulationRestrictionAutomation{
		consensusSystem:         consensusSystem,
		ledgerInstance:          ledgerInstance,
		stateMutex:              stateMutex,
		userOrderCancelCount:    make(map[string]int),
		userLastOrderTimestamps: make(map[string]time.Time),
	}
}

// StartOrderBookMonitoring starts continuous monitoring of DEX order book activities
func (automation *DexOrderBookManipulationRestrictionAutomation) StartOrderBookMonitoring() {
	ticker := time.NewTicker(OrderBookManipulationCheckInterval)

	go func() {
		for range ticker.C {
			automation.monitorOrderBookActivity()
		}
	}()
}

// monitorOrderBookActivity checks recent order book activities and enforces restrictions on manipulation
func (automation *DexOrderBookManipulationRestrictionAutomation) monitorOrderBookActivity() {
	automation.stateMutex.Lock()
	defer automation.stateMutex.Unlock()

	// Fetch recent order activities from Synnergy Consensus
	recentOrders := automation.consensusSystem.GetRecentOrderBookActivity()

	for _, order := range recentOrders {
		// Validate order cancellation limits and minimum order lifetime
		if !automation.validateOrderCancellationLimit(order) {
			automation.flagOrderBookViolation(order, "Exceeded maximum number of order cancellations per minute")
		} else if !automation.validateOrderLifetime(order) {
			automation.flagOrderBookViolation(order, "Order cancelled before minimum lifetime")
		}
	}
}

// validateOrderCancellationLimit checks if a user has exceeded the order cancellation limit
func (automation *DexOrderBookManipulationRestrictionAutomation) validateOrderCancellationLimit(order common.Order) bool {
	currentCancelCount := automation.userOrderCancelCount[order.UserID]
	if order.Type == "Cancel" && currentCancelCount+1 > MaxOrderCancelsPerMinute {
		return false
	}

	if order.Type == "Cancel" {
		automation.userOrderCancelCount[order.UserID]++
	}
	return true
}

// validateOrderLifetime checks if an order was canceled before the minimum allowed lifetime
func (automation *DexOrderBookManipulationRestrictionAutomation) validateOrderLifetime(order common.Order) bool {
	lastOrderTimestamp := automation.userLastOrderTimestamps[order.UserID]
	if time.Since(lastOrderTimestamp) < MinOrderLifetime {
		return false
	}

	// Update the last order timestamp for the user
	if order.Type == "New" {
		automation.userLastOrderTimestamps[order.UserID] = time.Now()
	}
	return true
}

// flagOrderBookViolation flags an order book activity that violates system rules and logs it in the ledger
func (automation *DexOrderBookManipulationRestrictionAutomation) flagOrderBookViolation(order common.Order, reason string) {
	fmt.Printf("DEX order book violation: User %s, Reason: %s\n", order.UserID, reason)

	// Log the violation into the ledger
	automation.logOrderBookViolation(order, reason)
}

// logOrderBookViolation logs the flagged order book violation into the ledger with full details
func (automation *DexOrderBookManipulationRestrictionAutomation) logOrderBookViolation(order common.Order, violationReason string) {
	// Encrypt the order data before logging
	encryptedData := automation.encryptOrderData(order)

	// Create a ledger entry with the violation details
	entry := common.LedgerEntry{
		ID:        fmt.Sprintf("order-book-violation-%s-%d", order.UserID, time.Now().Unix()),
		Timestamp: time.Now().Unix(),
		Type:      "DEX Order Book Violation",
		Status:    "Flagged",
		Details:   fmt.Sprintf("User %s flagged for order book manipulation. Reason: %s. Encrypted Data: %s", order.UserID, violationReason, encryptedData),
	}

	// Add the entry to the ledger
	err := automation.ledgerInstance.AddEntry(entry)
	if err != nil {
		fmt.Printf("Failed to log order book violation into ledger: %v\n", err)
	} else {
		fmt.Printf("Order book violation logged for user: %s\n", order.UserID)
	}
}

// encryptOrderData encrypts order data before logging for security
func (automation *DexOrderBookManipulationRestrictionAutomation) encryptOrderData(order common.Order) string {
	data := fmt.Sprintf("User ID: %s, Order ID: %s, Type: %s, Timestamp: %d", order.UserID, order.OrderID, order.Type, order.Timestamp)
	encryptedData, err := encryption.EncryptData([]byte(data))
	if err != nil {
		fmt.Println("Error encrypting order data:", err)
		return data
	}
	return string(encryptedData)
}
