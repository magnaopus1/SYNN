package high_availability

import (
    "fmt"
    "synnergy_network/pkg/common"
    "time"
)


// NewDataCollectionManager initializes a DataCollectionManager with a list of nodes and a data request interval
func NewDataCollectionManager(nodes []string, requestInterval time.Duration) *DataCollectionManager {
    return &DataCollectionManager{
        Nodes:               nodes,
        CollectedTransactions: []common.Transaction{},
        CollectedSubBlocks:    []common.SubBlock{},
        CollectedBlocks:       []common.Block{},
        DataRequestInterval: requestInterval,
    }
}

// StartDataCollection starts collecting data from the network at regular intervals
func (dcm *DataCollectionManager) StartDataCollection() {
    fmt.Printf("Starting data collection every %s...\n", dcm.DataRequestInterval)

    ticker := time.NewTicker(dcm.DataRequestInterval)
    go func() {
        for range ticker.C {
            dcm.mutex.Lock()
            dcm.CollectDataFromNetwork()
            dcm.mutex.Unlock()
        }
    }()
}

// CollectDataFromNetwork simulates collecting transactions, sub-blocks, and blocks from network nodes
func (dcm *DataCollectionManager) CollectDataFromNetwork() {
    fmt.Println("Collecting data from the network...")

    // Simulating data collection from each node
    for _, node := range dcm.Nodes {
        fmt.Printf("Requesting data from node %s...\n", node)

        // Simulate receiving transactions, sub-blocks, and blocks from the node
        newTransactions := dcm.simulateTransactionCollection(node)
        newSubBlocks := dcm.simulateSubBlockCollection(node)
        newBlocks := dcm.simulateBlockCollection(node)

        // Add the new data to the collected data
        dcm.CollectedTransactions = append(dcm.CollectedTransactions, newTransactions...)
        dcm.CollectedSubBlocks = append(dcm.CollectedSubBlocks, newSubBlocks...)
        dcm.CollectedBlocks = append(dcm.CollectedBlocks, newBlocks...)

        fmt.Printf("Collected %d transactions, %d sub-blocks, and %d blocks from node %s.\n", len(newTransactions), len(newSubBlocks), len(newBlocks), node)
    }
}

// simulateTransactionCollection simulates collecting transactions from a node
func (dcm *DataCollectionManager) TransactionCollection(node string) []common.Transaction {
    // Simulate random transactions
    tx1 := common.Transaction{FromAddress: "Alice", ToAddress: "Bob", Amount: 10.5}
    tx2 := common.Transaction{FromAddress: "Charlie", ToAddress: "Dave", Amount: 15.0}
    return []common.Transaction{tx1, tx2}
}

// simulateSubBlockCollection simulates collecting sub-blocks from a node
func (dcm *DataCollectionManager) SubBlockCollection(node string) []common.SubBlock {
    // Simulate random sub-blocks
    subBlock1 := common.SubBlock{Index: 1, Validator: "Validator1", PrevHash: "abc123"}
    subBlock2 := common.SubBlock{Index: 2, Validator: "Validator2", PrevHash: "def456"}
    return []common.SubBlock{subBlock1, subBlock2}
}

// simulateBlockCollection simulates collecting blocks from a node
func (dcm *DataCollectionManager) BlockCollection(node string) []common.Block {
    // Simulate random blocks
    block1 := common.Block{Index: 1, PrevHash: "0000", Hash: "abc123", Nonce: 1234}
    block2 := common.Block{Index: 2, PrevHash: "abc123", Hash: "def456", Nonce: 5678}
    return []common.Block{block1, block2}
}

// DistributeCollectedData distributes the collected data to the other nodes in the network
func (dcm *DataCollectionManager) DistributeCollectedData() {
    dcm.mutex.Lock()
    defer dcm.mutex.Unlock()

    fmt.Println("Distributing collected data to network nodes...")

    // Simulate distributing the collected data to each node
    for _, node := range dcm.Nodes {
        fmt.Printf("Sending collected data to node %s...\n", node)
        dcm.sendDataToNode(node)
    }
}

// sendDataToNode simulates sending data to a specific node
func (dcm *DataCollectionManager) sendDataToNode(node string) {
    // Simulate sending transactions, sub-blocks, and blocks to the node
    fmt.Printf("Sent %d transactions, %d sub-blocks, and %d blocks to node %s.\n", len(dcm.CollectedTransactions), len(dcm.CollectedSubBlocks), len(dcm.CollectedBlocks), node)
}

// ClearCollectedData clears the collected transactions, sub-blocks, and blocks after distribution
func (dcm *DataCollectionManager) ClearCollectedData() {
    dcm.mutex.Lock()
    defer dcm.mutex.Unlock()

    fmt.Println("Clearing collected data...")

    dcm.CollectedTransactions = []common.Transaction{}
    dcm.CollectedSubBlocks = []common.SubBlock{}
    dcm.CollectedBlocks = []common.Block{}

    fmt.Println("Collected data cleared.")
}




// simulateTransactionCollection simulates the collection of transactions from a node
func (dcm *DataCollectionManager) simulateTransactionCollection(node string) []common.Transaction {
    fmt.Printf("Simulating transaction collection from node %s...\n", node)
    return []common.Transaction{} // Return an empty slice for simulation
}

// simulateSubBlockCollection simulates the collection of sub-blocks from a node
func (dcm *DataCollectionManager) simulateSubBlockCollection(node string) []common.SubBlock {
    fmt.Printf("Simulating sub-block collection from node %s...\n", node)
    return []common.SubBlock{} // Return an empty slice for simulation
}

// simulateBlockCollection simulates the collection of blocks from a node
func (dcm *DataCollectionManager) simulateBlockCollection(node string) []common.Block {
    fmt.Printf("Simulating block collection from node %s...\n", node)
    return []common.Block{} // Return an empty slice for simulation
}
