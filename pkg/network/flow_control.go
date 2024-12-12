package network

import (
	"fmt"
	"time"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"
)

// NewFlowControlManager initializes the FlowControlManager with the provided limits
func NewFlowControlManager(maxPendingTransactions, maxSubBlockSize, maxBlockSize int, ledgerInstance *ledger.Ledger) *FlowControlManager {
	return &FlowControlManager{
		MaxPendingTransactions: maxPendingTransactions,
		MaxSubBlockSize:        maxSubBlockSize,
		MaxBlockSize:           maxBlockSize,
		PendingTransactions:    make([]common.Transaction, 0),
		PendingSubBlocks:       make([]common.SubBlock, 0),
		ledgerInstance:         ledgerInstance,
	}
}

// AddTransaction adds a new transaction to the pending pool, ensuring the limit is respected
func (fcm *FlowControlManager) AddTransaction(tx common.Transaction) error {
    fcm.mutex.Lock()
    defer fcm.mutex.Unlock()

    if len(fcm.PendingTransactions) >= fcm.MaxPendingTransactions {
        return fmt.Errorf("transaction pool full, cannot add transaction")
    }

    fcm.PendingTransactions = append(fcm.PendingTransactions, tx)
    fmt.Printf("Transaction from %s to %s added to pending pool.\n", tx.FromAddress, tx.ToAddress)

    // Log the event to the ledger
    logMessage := fmt.Sprintf("Transaction Added: from %s to %s", tx.FromAddress, tx.ToAddress)
    fcm.ledgerInstance.LogFlowControlEvent(logMessage) // Adjusted to pass a single argument

    return nil
}



// AddSubBlock creates and adds a sub-block from pending transactions
func (fcm *FlowControlManager) AddSubBlock() error {
    fcm.mutex.Lock()
    defer fcm.mutex.Unlock()

    if len(fcm.PendingTransactions) == 0 {
        return fmt.Errorf("no pending transactions to create a sub-block")
    }

    // Limit the number of transactions per sub-block
    var txBatch []common.Transaction
    if len(fcm.PendingTransactions) > fcm.MaxSubBlockSize {
        txBatch = fcm.PendingTransactions[:fcm.MaxSubBlockSize]
        fcm.PendingTransactions = fcm.PendingTransactions[fcm.MaxSubBlockSize:]
    } else {
        txBatch = fcm.PendingTransactions
        fcm.PendingTransactions = nil
    }

    // Create the sub-block
    subBlock := common.SubBlock{
        Index:        len(fcm.PendingSubBlocks),
        Timestamp:    time.Now(),
        Transactions: txBatch,
    }

    // Add sub-block to the pending pool
    fcm.PendingSubBlocks = append(fcm.PendingSubBlocks, subBlock)
    fmt.Printf("Sub-block %d created with %d transactions.\n", subBlock.Index, len(txBatch))

    // Log the creation of the sub-block to the ledger
    logMessage := fmt.Sprintf("SubBlock Created: SubBlock #%d at %s", subBlock.Index, time.Now().Format(time.RFC3339))
    fcm.ledgerInstance.LogFlowControlEvent(logMessage) // Adjusted to pass a single argument

    return nil
}


// CreateBlock consolidates sub-blocks into a full block and resets pending sub-blocks
func (fcm *FlowControlManager) CreateBlock() (*common.Block, error) {
    fcm.mutex.Lock()
    defer fcm.mutex.Unlock()

    if len(fcm.PendingSubBlocks) == 0 {
        return nil, fmt.Errorf("no sub-blocks available to create a block")
    }

    // Limit the number of sub-blocks per block
    var subBlockBatch []common.SubBlock
    if len(fcm.PendingSubBlocks) > fcm.MaxBlockSize {
        subBlockBatch = fcm.PendingSubBlocks[:fcm.MaxBlockSize]
        fcm.PendingSubBlocks = fcm.PendingSubBlocks[fcm.MaxBlockSize:]
    } else {
        subBlockBatch = fcm.PendingSubBlocks
        fcm.PendingSubBlocks = nil
    }

    // Get the list of finalized blocks from the ledger
    finalizedBlocks := fcm.ledgerInstance.GetFinalizedBlocks()

    // Create a new block from the sub-blocks
    block := common.Block{
        Index:     len(finalizedBlocks),
        Timestamp: time.Now(),
        SubBlocks: subBlockBatch,
    }

    // Log the creation of the block to the ledger as a single concatenated message
    logMessage := fmt.Sprintf("Block Created: Block #%d with %d sub-blocks at %s", block.Index, len(subBlockBatch), time.Now().Format(time.RFC3339))
    fcm.ledgerInstance.LogFlowControlEvent(logMessage)

    fmt.Printf("Block %d created with %d sub-blocks.\n", block.Index, len(subBlockBatch))
    return &block, nil
}


// BroadcastTransaction broadcasts a transaction to all connected peers
func (fcm *FlowControlManager) BroadcastTransaction(tx common.Transaction, coordinator *DistributedNetworkCoordinator) error {
    // Broadcast the transaction to all connected nodes
    fmt.Printf("Broadcasting transaction from %s to %s.\n", tx.FromAddress, tx.ToAddress)
    coordinator.BroadcastMessage(fmt.Sprintf("New Transaction: %s -> %s: %.2f", tx.FromAddress, tx.ToAddress, tx.Amount))

    // Log the broadcast to the ledger (single string argument)
    logMessage := fmt.Sprintf("Transaction Broadcast: %s -> %s", tx.FromAddress, tx.ToAddress)
    fcm.ledgerInstance.LogFlowControlEvent(logMessage)

    return nil
}


// BroadcastBlock broadcasts a full block to all connected peers
func (fcm *FlowControlManager) BroadcastBlock(block *common.Block, coordinator *DistributedNetworkCoordinator) error {
    // Broadcast the block to all connected nodes
    fmt.Printf("Broadcasting block %d to peers.\n", block.Index)
    coordinator.BroadcastMessage(fmt.Sprintf("New Block: #%d", block.Index))

    // Log the broadcast to the ledger (single string argument)
    logMessage := fmt.Sprintf("Block Broadcast: Block #%d", block.Index)
    fcm.ledgerInstance.LogFlowControlEvent(logMessage)

    return nil
}


// MonitorNetworkTraffic periodically checks the size of pending pools and creates sub-blocks or blocks as needed
func (fcm *FlowControlManager) MonitorNetworkTraffic(coordinator *DistributedNetworkCoordinator) {
	go func() {
		for {
			time.Sleep(10 * time.Second) // Check every 10 seconds

			fcm.mutex.Lock()
			// Check if there are enough transactions to form a sub-block
			if len(fcm.PendingTransactions) >= fcm.MaxSubBlockSize {
				if err := fcm.AddSubBlock(); err != nil {
					fmt.Println("Error adding sub-block:", err)
				}
			}

			// Check if there are enough sub-blocks to form a full block
			if len(fcm.PendingSubBlocks) >= fcm.MaxBlockSize {
				if block, err := fcm.CreateBlock(); err == nil {
					fcm.BroadcastBlock(block, coordinator)
				} else {
					fmt.Println("Error creating block:", err)
				}
			}
			fcm.mutex.Unlock()
		}
	}()
}

