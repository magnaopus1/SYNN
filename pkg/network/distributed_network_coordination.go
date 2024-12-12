package network

import (
	"encoding/json"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// DistributedNetworkCoordinator manages the coordination of nodes in a decentralized network
type DistributedNetworkCoordinator struct {
    Nodes           map[string]*NodeInfo         // All nodes participating in the network
    ActivePeers     map[string]*PeerConnection   // Active peer connections
    mutex           sync.Mutex                   // For thread-safe access
    ledgerInstance  *ledger.Ledger               // Pointer to the ledger for logging and auditing
    Blockchain      *common.Blockchain           // Pointer to the blockchain (add this field)
	Consensus		*common.SynnergyConsensus
}


// NewDistributedNetworkCoordinator initializes a new network coordinator
func NewDistributedNetworkCoordinator(ledgerInstance *ledger.Ledger) *DistributedNetworkCoordinator {
	return &DistributedNetworkCoordinator{
		Nodes:          make(map[string]*NodeInfo),
		ActivePeers:    make(map[string]*PeerConnection),
		ledgerInstance: ledgerInstance,
	}
}

// RegisterNode registers a new node in the network
func (dnc *DistributedNetworkCoordinator) RegisterNode(nodeID, address string) {
	dnc.mutex.Lock()
	defer dnc.mutex.Unlock()

	if _, exists := dnc.Nodes[nodeID]; exists {
		fmt.Printf("Node %s is already registered.\n", nodeID)
		return
	}

	nodeInfo := &NodeInfo{
		NodeID:  nodeID,
		Address: address,
	}

	dnc.Nodes[nodeID] = nodeInfo
	fmt.Printf("Node %s registered at address %s.\n", nodeID, address)

	// Log the event to the ledger (passing only two arguments as required)
	dnc.ledgerInstance.LogNodeEvent("Node Registered", nodeID)
}


// EstablishConnection connects to a peer node and adds the connection to the active peers
func (dnc *DistributedNetworkCoordinator) EstablishConnection(nodeID string) error {
    dnc.mutex.Lock()
    defer dnc.mutex.Unlock()

    nodeInfo, exists := dnc.Nodes[nodeID]
    if !exists {
        return fmt.Errorf("Node %s not found in the network", nodeID)
    }

    // Create a new connection using the connection pool
    connectionPool := NewConnectionPool(5 * time.Minute) // Example max idle time
    _, err := connectionPool.GetConnection(nodeID, nodeInfo.Address)
    if err != nil {
        return fmt.Errorf("Failed to establish connection with node %s: %v", nodeID, err)
    }

    // Store the ConnectionPool in the PeerConnection
    dnc.ActivePeers[nodeID] = &PeerConnection{
        NodeID:       nodeID,
        Connection:   connectionPool,  // Store the ConnectionPool here
        LastPingTime: time.Now(),
        IsAlive:      true,
    }

    // Log the connection to the ledger (passing only two arguments)
    dnc.ledgerInstance.LogNodeEvent("Connection Established", nodeID)

    fmt.Printf("Connection established with node %s at %s.\n", nodeID, nodeInfo.Address)
    return nil
}



// DisconnectPeer disconnects from a peer node and removes it from the active peers
func (dnc *DistributedNetworkCoordinator) DisconnectPeer(nodeID string) error {
    dnc.mutex.Lock()
    defer dnc.mutex.Unlock()

    peer, exists := dnc.ActivePeers[nodeID]
    if !exists {
        return fmt.Errorf("Peer %s not found in active connections", nodeID)
    }

    // Remove the connection from the pool and mark the peer as disconnected
    peer.Connection.RemoveConnection(nodeID)
    delete(dnc.ActivePeers, nodeID)

    // Log the disconnection to the ledger (passing only two arguments)
    dnc.ledgerInstance.LogNodeEvent("Peer Disconnected", nodeID)

    fmt.Printf("Disconnected from peer node %s.\n", nodeID)
    return nil
}


// BroadcastMessage sends a message to all active peers in the network
func (dnc *DistributedNetworkCoordinator) BroadcastMessage(message string) {
    dnc.mutex.Lock()
    defer dnc.mutex.Unlock()

    for nodeID, peer := range dnc.ActivePeers {
        // Retrieve an active connection from the ConnectionPool
        conn, err := peer.Connection.GetConnection(nodeID, "") // Assuming ConnectionPool provides a method to retrieve a connection
        if err != nil {
            fmt.Printf("Failed to get connection for node %s: %v\n", nodeID, err)
            continue
        }

        // Here we send the message to each peer node via its connection
        fmt.Printf("Broadcasting message to node %s: %s\n", nodeID, message)

        // Send the message to the peer
        _, err = conn.Write([]byte(message))
        if err != nil {
            fmt.Printf("Failed to send message to node %s: %v\n", nodeID, err)
            continue
        }

        fmt.Printf("Message successfully sent to node %s.\n", nodeID)
    }
}


// MonitorPeerHealth periodically pings all active peers to ensure they are alive
func (dnc *DistributedNetworkCoordinator) MonitorPeerHealth() {
	go func() {
		for {
			time.Sleep(30 * time.Second) // Ping interval

			dnc.mutex.Lock()
			for nodeID, peer := range dnc.ActivePeers {
				// Simulate a ping (in a real system, this would be an actual network ping)
				if time.Since(peer.LastPingTime) > 1*time.Minute {
					peer.IsAlive = false
					fmt.Printf("Peer node %s is unresponsive. Last ping was over 1 minute ago.\n", nodeID)

					// Log the event to the ledger (without time argument)
					dnc.ledgerInstance.LogNodeEvent("Peer Unresponsive", nodeID)
				} else {
					peer.IsAlive = true
					peer.LastPingTime = time.Now()
					fmt.Printf("Peer node %s is alive. Ping successful.\n", nodeID)

					// Log the successful ping to the ledger (without time argument)
					dnc.ledgerInstance.LogNodeEvent("Peer Ping Success", nodeID)
				}
			}
			dnc.mutex.Unlock()
		}
	}()
}

// SyncBlockWithPeers sends a block to all peers for validation and propagation
func (dnc *DistributedNetworkCoordinator) SyncBlockWithPeers(block *common.Block) {
    dnc.mutex.Lock()
    defer dnc.mutex.Unlock()

    for nodeID, peer := range dnc.ActivePeers {
        fmt.Printf("Sending block %d to peer node %s for validation.\n", block.Index, nodeID)

        // Serialize the block to JSON before sending it over the network
        blockJSON, err := json.Marshal(block)
        if err != nil {
            fmt.Printf("Failed to serialize block %d: %v\n", block.Index, err)
            continue
        }

        // Retrieve the connection from the connection pool
        conn, err := peer.Connection.GetConnection(nodeID, "")
        if err != nil {
            fmt.Printf("Failed to get connection for peer node %s: %v\n", nodeID, err)
            continue
        }

        // Send the block to the peer
        _, err = conn.Write(blockJSON)
        if err != nil {
            fmt.Printf("Failed to send block %d to peer node %s: %v\n", block.Index, nodeID, err)
            continue
        }

        // Log the sync event to the ledger (only pass two arguments as required)
        dnc.ledgerInstance.LogNodeEvent(fmt.Sprintf("Block %d Sent", block.Index), nodeID)
    }
}


// detectFork checks if the new block causes a fork in the current chain
func (dnc *DistributedNetworkCoordinator) detectFork(newBlock *common.Block) bool {
    // Access the blockchain and get the latest block
    latestBlock, err := dnc.ledgerInstance.GetLatestBlock()
    
    // Check if there was an error in fetching the latest block
    if err != nil {
        fmt.Printf("Error retrieving the latest block: %v\n", err)
        return false // You can handle the error based on your requirements (e.g., return true if an error should indicate a potential fork)
    }

    // Compare the previous block hash of the new block with the current chain's latest block
    if newBlock.PrevHash != latestBlock.Hash {
        fmt.Printf("Fork detected: new block's previous hash %s doesn't match latest block hash %s\n", newBlock.PrevHash, latestBlock.Hash)
        return true
    }
    return false
}



// findLongestChain finds the longest valid chain among peers based on the new block
func (dnc *DistributedNetworkCoordinator) findLongestChain(newBlock *common.Block) *common.Block {
    var longestChain *common.Block

    // Access the latest block to get the chain length from its index
    latestBlock, err := dnc.ledgerInstance.GetLatestBlock()
    if err != nil {
        fmt.Printf("Error retrieving the latest block: %v\n", err)
        return nil
    }
    longestChainLength := latestBlock.Index  // Assuming Index is the block number or height

    // Iterate over active peers and request their chains
    for nodeID, peer := range dnc.ActivePeers {
        fmt.Printf("Requesting chain from peer node %s...\n", nodeID)
        peerChain := dnc.requestChainFromPeer(peer)

        if peerChain != nil && peerChain.Index > longestChainLength {
            fmt.Printf("Peer node %s has a longer chain with block index %d\n", nodeID, peerChain.Index)
            longestChain = peerChain
            longestChainLength = peerChain.Index
        }
    }

    return longestChain
}



// HandleChainReorganization detects and manages chain forks across peers
func (dnc *DistributedNetworkCoordinator) HandleChainReorganization(block *common.Block) {
	fmt.Println("Handling chain reorganization due to detected fork...")

	// Detect fork based on the received block
	if dnc.detectFork(block) {
		// Logic to handle forks and re-sync the chain
		longestChain := dnc.findLongestChain(block)

		if longestChain != nil {
			fmt.Println("Valid chain found. Re-syncing...")
			dnc.syncWithValidChain(longestChain)
		} else {
			fmt.Println("No valid longer chain found. Retaining the current chain.")
		}

		// Log the chain reorganization event to the ledger
		dnc.ledgerInstance.LogEvent("Chain Reorganization Detected")

		// Broadcast the reorganization event to other peers
		dnc.broadcastChainReorganization(block)
	} else {
		fmt.Println("No fork detected. No reorganization required.")
	}
}

// requestChainFromPeer requests the chain from a peer node
func (dnc *DistributedNetworkCoordinator) requestChainFromPeer(peer *PeerConnection) *common.Block {
    // Retrieve the connection from the connection pool
    conn, err := peer.Connection.GetConnection(peer.NodeID, "")
    if err != nil {
        fmt.Printf("Failed to get connection for peer %s: %v\n", peer.NodeID, err)
        return nil
    }

    // Send a request for the peer's chain data
    requestMessage := "REQUEST_CHAIN"
    _, err = conn.Write([]byte(requestMessage))
    if err != nil {
        fmt.Printf("Failed to request chain from peer %s: %v\n", peer.NodeID, err)
        return nil
    }

    // Receive the response (entire chain or the latest block)
    var receivedBlock common.Block
    decoder := json.NewDecoder(conn)
    err = decoder.Decode(&receivedBlock)
    if err != nil {
        fmt.Printf("Failed to decode chain from peer %s: %v\n", peer.NodeID, err)
        return nil
    }

    // Return the received chain data (could be a single block or entire chain)
    return &receivedBlock
}

// validateChain validates the chain according to consensus rules
func (dnc *DistributedNetworkCoordinator) validateChain(chain *common.Block) bool {
    currentBlock := chain
    for currentBlock != nil {
        // Retrieve the previous block using its hash (PrevHash)
        previousBlock := dnc.ledgerInstance.GetBlockByHash(currentBlock.PrevHash)

        if previousBlock == nil {
            fmt.Printf("Chain validation failed: previous block not found.\n")
            return false
        }

        // Convert previousBlock (of type *ledger.Block) to *common.Block if the structures are compatible
        convertedPreviousBlock := &common.Block{
            Hash:         previousBlock.Hash,
            PrevHash:     previousBlock.PrevHash,
            Index:        previousBlock.Index,
            // Add other fields that match between the two structs
        }

        // Check if the previous block exists and if the hash matches
        if currentBlock.PrevHash != convertedPreviousBlock.Hash {
            fmt.Printf("Chain validation failed: Block %d has invalid previous hash\n", currentBlock.Index)
            return false
        }

        // Synnergy Consensus, validate the proof
        if !dnc.Consensus.ValidateChain() {
            fmt.Printf("Chain validation failed: Block %d has invalid Proof of Work\n", currentBlock.Index)
            return false
        }

        // Move to the previous block in the chain
        currentBlock = convertedPreviousBlock
    }

    fmt.Println("Chain validation successful.")
    return true
}




// isLongerChain checks if the peer's chain is longer than the current chain
func (dnc *DistributedNetworkCoordinator) isLongerChain(chain *common.Block) bool {
    // Get the local chain's latest block and handle the error
    localLatestBlock, err := dnc.ledgerInstance.GetLatestBlock()
    if err != nil {
        fmt.Printf("Error retrieving the latest block: %v\n", err)
        return false
    }

    // Compare the index of the peer's block with the local block
    if chain.Index > localLatestBlock.Index {
        fmt.Printf("Peer's chain is longer (peer block index: %d, local block index: %d).\n", chain.Index, localLatestBlock.Index)
        return true
    }

    fmt.Printf("Local chain is longer or equal (peer block index: %d, local block index: %d).\n", chain.Index, localLatestBlock.Index)
    return false
}


// syncWithValidChain replaces the local chain with the valid longer chain
func (dnc *DistributedNetworkCoordinator) syncWithValidChain(newChain *common.Block) {
    fmt.Println("Syncing with the valid longer chain...")

    // Ensure that the new chain is valid
    if dnc.validateChain(newChain) {
        // Convert the validated chain to an array or slice of common blocks
        chainBlocks := dnc.convertToChain(newChain)

        // Access blocks in ledger using a method or a public field
        currentBlocks := dnc.ledgerInstance.GetBlocks()  // Using the GetBlocks method

        // Convert []common.Block to []ledger.Block
        ledgerChainBlocks := convertCommonToLedgerBlocks(chainBlocks)

        // Replace the local chain with the new one manually
        if len(ledgerChainBlocks) > len(currentBlocks) {
            dnc.ledgerInstance.ReplaceChain(ledgerChainBlocks)
            fmt.Println("Local chain successfully synchronized with the longest valid chain.")
        } else {
            fmt.Println("Sync failed: The new chain was not longer.")
        }
    } else {
        fmt.Println("Sync failed: The new chain was invalid.")
    }
}

// convertCommonToLedgerBlocks converts a slice of common.Block to a slice of ledger.Block
func convertCommonToLedgerBlocks(commonBlocks []common.Block) []ledger.Block {
    var ledgerBlocks []ledger.Block
    for _, commonBlock := range commonBlocks {
        // Convert []common.SubBlock to []ledger.SubBlock
        ledgerSubBlocks := convertCommonToLedgerSubBlocks(commonBlock.SubBlocks)

        // Manually map fields from common.Block to ledger.Block
        ledgerBlock := ledger.Block{
            BlockID:     commonBlock.BlockID,
            Index:       commonBlock.Index,
            Timestamp:   commonBlock.Timestamp,
            SubBlocks:   ledgerSubBlocks,  // Use the converted ledger sub-blocks
            PrevHash:    commonBlock.PrevHash,
            Hash:        commonBlock.Hash,
            Nonce:       commonBlock.Nonce,
            Difficulty:  commonBlock.Difficulty,
            MinerReward: commonBlock.MinerReward,
            Validators:  commonBlock.Validators,
            Status:      commonBlock.Status,  // Map the new Status field
        }
        ledgerBlocks = append(ledgerBlocks, ledgerBlock)
    }
    return ledgerBlocks
}

// convertCommonToLedgerSubBlocks converts a slice of common.SubBlock to a slice of ledger.SubBlock
func convertCommonToLedgerSubBlocks(commonSubBlocks []common.SubBlock) []ledger.SubBlock {
    var ledgerSubBlocks []ledger.SubBlock
    for _, commonSubBlock := range commonSubBlocks {
        // Convert []common.Transaction to []ledger.Transaction
        ledgerTransactions := convertCommonToLedgerTransactions(commonSubBlock.Transactions)

        // Convert common.PoHProof to ledger.PoHProof
        ledgerPoHProof := convertCommonToLedgerPoHProof(commonSubBlock.PoHProof)

        // Manually map fields from common.SubBlock to ledger.SubBlock
        ledgerSubBlock := ledger.SubBlock{
            SubBlockID:   commonSubBlock.SubBlockID,
            Index:        commonSubBlock.Index,
            Timestamp:    commonSubBlock.Timestamp,
            Transactions: ledgerTransactions,  // Use the converted ledger transactions
            Validator:    commonSubBlock.Validator,
            PrevHash:     commonSubBlock.PrevHash,
            Hash:         commonSubBlock.Hash,
            PoHProof:     ledgerPoHProof,      // Use the converted ledger PoHProof
            Status:       commonSubBlock.Status,  // Map the new Status field
            Signature:    commonSubBlock.Signature, // Map the new Signature field
        }
        ledgerSubBlocks = append(ledgerSubBlocks, ledgerSubBlock)
    }
    return ledgerSubBlocks
}

// convertCommonToLedgerTransactions converts a slice of common.Transaction to a slice of ledger.Transaction
func convertCommonToLedgerTransactions(commonTransactions []common.Transaction) []ledger.Transaction {
    var ledgerTransactions []ledger.Transaction
    for _, commonTransaction := range commonTransactions {
        // Manually map fields from common.Transaction to ledger.Transaction
        ledgerTransaction := ledger.Transaction{
            TransactionID:    commonTransaction.TransactionID,
            FromAddress:      commonTransaction.FromAddress,
            ToAddress:        commonTransaction.ToAddress,
            Amount:           commonTransaction.Amount,
            Fee:              commonTransaction.Fee,
            TokenStandard:    commonTransaction.TokenStandard,
            TokenID:          commonTransaction.TokenID,
            Timestamp:        commonTransaction.Timestamp,
            SubBlockID:       commonTransaction.SubBlockID,
            BlockID:          commonTransaction.BlockID,
            ValidatorID:      commonTransaction.ValidatorID,
            Signature:        commonTransaction.Signature,
            Status:           commonTransaction.Status,
            EncryptedData:    commonTransaction.EncryptedData,
            DecryptedData:    commonTransaction.DecryptedData,
            ExecutionResult:  commonTransaction.ExecutionResult,
            FrozenAmount:     commonTransaction.FrozenAmount,
            RefundAmount:     commonTransaction.RefundAmount,
            ReversalRequested: commonTransaction.ReversalRequested, // Add this field mapping
        }
        ledgerTransactions = append(ledgerTransactions, ledgerTransaction)
    }
    return ledgerTransactions
}

// convertCommonToLedgerPoHProof converts a common.PoHProof to a ledger.PoHProof
func convertCommonToLedgerPoHProof(commonPoHProof common.PoHProof) ledger.PoHProof {
    // Manually map fields from common.PoHProof to ledger.PoHProof
    return ledger.PoHProof{
        Sequence:  commonPoHProof.Sequence,
        Timestamp: commonPoHProof.Timestamp,
        Hash:      commonPoHProof.Hash,
    }
}







// convertToChain converts a linked list of blocks into a slice of blocks
func (dnc *DistributedNetworkCoordinator) convertToChain(startBlock *common.Block) []common.Block {
    var chain []common.Block
    currentBlock := startBlock

    for currentBlock != nil {
        chain = append(chain, *currentBlock)  // Add the block to the slice

        // Convert *ledger.Block to *common.Block if needed
        previousLedgerBlock := dnc.ledgerInstance.GetBlockByHash(currentBlock.PrevHash)
        if previousLedgerBlock != nil {
            // Assuming manual conversion between ledger.Block and common.Block
            currentBlock = &common.Block{
                Hash:      previousLedgerBlock.Hash,
                PrevHash:  previousLedgerBlock.PrevHash,
                Index:     previousLedgerBlock.Index,
                // Map other fields if necessary
            }
        } else {
            currentBlock = nil
        }
    }

    // Reverse the chain slice so that it starts from the genesis block
    for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
        chain[i], chain[j] = chain[j], chain[i]
    }

    return chain
}



// broadcastChainReorganization broadcasts the reorganization event to all peers
func (dnc *DistributedNetworkCoordinator) broadcastChainReorganization(block *common.Block) {
    // Notify all peers about the reorganization
    for nodeID, peer := range dnc.ActivePeers {
        fmt.Printf("Notifying peer %s of chain reorganization.\n", nodeID)

        // Retrieve the connection from the connection pool
        conn, err := peer.Connection.GetConnection(nodeID, "")
        if err != nil {
            fmt.Printf("Failed to get connection for peer %s: %v\n", nodeID, err)
            continue
        }

        // Prepare the reorganization message
        reorgMessage := fmt.Sprintf("CHAIN_REORG Block %d", block.Index)

        // Send the notification to the peer
        _, err = conn.Write([]byte(reorgMessage))
        if err != nil {
            fmt.Printf("Failed to notify peer %s about reorganization: %v\n", nodeID, err)
        }
    }
}

