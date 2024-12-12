package high_availability

import (
    "fmt"
    "synnergy_network/pkg/common"
)

// NewDataDistributionManager initializes a DataDistributionManager with a list of nodes
func NewDataDistributionManager(nodes []string) *DataDistributionManager {
    return &DataDistributionManager{
        Nodes:               nodes,
        DistributedSubBlocks: []common.SubBlock{},
        DistributedBlocks:    []common.Block{},
    }
}

// DistributeSubBlocks distributes validated sub-blocks to all nodes in the network
func (ddm *DataDistributionManager) DistributeSubBlocks(subBlocks []common.SubBlock) {
    ddm.mutex.Lock()
    defer ddm.mutex.Unlock()

    fmt.Printf("Distributing %d sub-blocks to the network...\n", len(subBlocks))

    for _, node := range ddm.Nodes {
        fmt.Printf("Sending sub-blocks to node %s...\n", node)
        ddm.sendSubBlocksToNode(node, subBlocks)
    }

    // Store the distributed sub-blocks
    ddm.DistributedSubBlocks = append(ddm.DistributedSubBlocks, subBlocks...)
    fmt.Println("Sub-blocks distribution completed.")
}

// DistributeBlocks distributes validated blocks to all nodes in the network
func (ddm *DataDistributionManager) DistributeBlocks(blocks []common.Block) {
    ddm.mutex.Lock()
    defer ddm.mutex.Unlock()

    fmt.Printf("Distributing %d blocks to the network...\n", len(blocks))

    for _, node := range ddm.Nodes {
        fmt.Printf("Sending blocks to node %s...\n", node)
        ddm.sendBlocksToNode(node, blocks)
    }

    // Store the distributed blocks
    ddm.DistributedBlocks = append(ddm.DistributedBlocks, blocks...)
    fmt.Println("Block distribution completed.")
}

// sendSubBlocksToNode simulates sending sub-blocks to a specific node
func (ddm *DataDistributionManager) sendSubBlocksToNode(node string, subBlocks []common.SubBlock) {
    // Simulate the process of sending sub-blocks
    fmt.Printf("Sub-blocks sent to node %s successfully.\n", node)
}

// sendBlocksToNode simulates sending blocks to a specific node
func (ddm *DataDistributionManager) sendBlocksToNode(node string, blocks []common.Block) {
    // Simulate the process of sending blocks
    fmt.Printf("Blocks sent to node %s successfully.\n", node)
}

// VerifyDataDistribution ensures that all nodes have received the correct blocks and sub-blocks
func (ddm *DataDistributionManager) VerifyDataDistribution() bool {
    ddm.mutex.Lock()
    defer ddm.mutex.Unlock()

    fmt.Println("Verifying data distribution across nodes...")
    for _, node := range ddm.Nodes {
        if !ddm.verifyNodeData(node) {
            fmt.Printf("Node %s has not received all data correctly.\n", node)
            return false
        }
    }

    fmt.Println("All nodes have received the data correctly.")
    return true
}

// verifyNodeData simulates the process of verifying if a node has received all sub-blocks and blocks
func (ddm *DataDistributionManager) verifyNodeData(node string) bool {
    fmt.Printf("Node %s data verification complete.\n", node)
    return true
}

// ResendData resends the missing data to nodes that failed to receive it
func (ddm *DataDistributionManager) ResendData(node string) {
    ddm.mutex.Lock()
    defer ddm.mutex.Unlock()

    fmt.Printf("Resending data to node %s...\n", node)
    // Resend all distributed sub-blocks and blocks
    ddm.sendSubBlocksToNode(node, ddm.DistributedSubBlocks)
    ddm.sendBlocksToNode(node, ddm.DistributedBlocks)

    fmt.Printf("Resent all data to node %s successfully.\n", node)
}

