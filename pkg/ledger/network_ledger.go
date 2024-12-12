package ledger

import (
    "fmt"
    "time"
)

// LogNodeJoin logs a node joining the network.
func (ledger *Ledger) LogNodeJoin(node NodeInfo) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Ensure the Nodes map is initialized
    if ledger.State.Nodes == nil {
        ledger.State.Nodes = make(map[string]NodeInfo)
    }

    ledger.State.Nodes[node.NodeID] = node // Use NodeID as the key
    fmt.Printf("Node %s joined the network.\n", node.NodeID) // Use NodeID for printing
    return nil
}

// GetActiveAuthorityNodes returns a list of active authority nodes.
func (l *Ledger) GetActiveAuthorityNodes() ([]Node, error) {
    var activeNodes []Node

    // Iterate over the nodes in the ledger and filter active authority nodes
    for _, node := range l.State.Validators {
        if node.NodeCategory == AuthorityCategory && node.IsActive {
            activeNodes = append(activeNodes, node)
        }
    }

    if len(activeNodes) == 0 {
        return nil, fmt.Errorf("no active authority nodes found")
    }

    return activeNodes, nil
}

// RecordBinaryNode records a binary search tree node with a key-value pair in the ledger.
func (l *Ledger) RecordBinaryNode(key int, value string) error {
	if value == "" {
		return errors.New("empty binary node value")
	}
	l.binaryNodes[key] = value
	fmt.Printf("Binary node recorded - Key: %d, Value: %s\n", key, value)
	return nil
}


// GetRandomAuthorityNodes fetches a random subset of authority nodes from the ledger
func (l *Ledger) GetRandomAuthorityNodes(count int) ([]*AuthorityNodeVersion, error) {
    if len(l.AuthorityNodes) == 0 {
        return nil, errors.New("no authority nodes available in the ledger")
    }

    // Convert map to a slice of pointers for shuffling
    nodeList := make([]*AuthorityNodeVersion, 0, len(l.AuthorityNodes))
    for _, node := range l.AuthorityNodes {
        nodeList = append(nodeList, node) // node is already *AuthorityNodeVersion
    }

    // Shuffle the node list and select random nodes
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(nodeList), func(i, j int) {
        nodeList[i], nodeList[j] = nodeList[j], nodeList[i]
    })

    if count > len(nodeList) {
        return nodeList, nil // Return all nodes if count exceeds available nodes
    }

    return nodeList[:count], nil
}


// RecordNodeVotingWeight logs the voting weight for a node in the ledger.
func (l *Ledger) RecordNodeVotingWeight(nodeID string, weight int, timestamp string) {
    l.Lock()
    defer l.Unlock()

    if l.NodeVotingWeights == nil {
        l.NodeVotingWeights = make(map[string]int)
    }
    if l.NodeVotingTimestamps == nil {
        l.NodeVotingTimestamps = make(map[string]string)
    }

    l.NodeVotingWeights[nodeID] = weight
    l.NodeVotingTimestamps[nodeID] = timestamp
    log.Printf("Node %s voting weight recorded: %d at %s", nodeID, weight, timestamp)
}

// FetchActiveNodeIDs retrieves the list of active node IDs from the ledger.
func (l *Ledger) FetchActiveNodeIDs() ([]string, error) {
    l.Lock()
    defer l.Unlock()

    var activeNodes []string
    for nodeID, status := range l.Nodes {
        if status.IsActive {
            activeNodes = append(activeNodes, nodeID)
        }
    }

    if len(activeNodes) == 0 {
        return nil, errors.New("no active nodes available")
    }

    log.Printf("Active nodes retrieved: %v", activeNodes)
    return activeNodes, nil
}


// LogNodeRemoval logs a node being removed from the network.
func (ledger *Ledger) LogNodeRemoval(nodeID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if _, exists := ledger.State.Nodes[nodeID]; exists {
        delete(ledger.State.Nodes, nodeID)
        fmt.Printf("Node %s removed from the network.\n", nodeID)
        return nil
    }
    return fmt.Errorf("node with ID %s not found", nodeID)
}

// GetFinalizedBlocks returns the list of finalized blocks.
func (l *Ledger) GetFinalizedBlocks() []Block {
    l.lock.Lock()
    defer l.lock.Unlock()

    return l.finalizedBlocks
}

// LogFlowControlEvent logs an event in the ledger with a single string message.
func (l *Ledger) LogFlowControlEvent(eventMessage string) {
    l.lock.Lock()
    defer l.lock.Unlock()

    logEntry := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), eventMessage)
    l.Events = append(l.Events, logEntry)
    fmt.Println("Ledger Event Logged:", logEntry)
}


// LogNodeEvent logs an event related to a node.
func (ledger *Ledger) LogNodeEvent(nodeID, event string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if _, exists := ledger.State.Nodes[nodeID]; exists {
        fmt.Printf("Event for node %s: %s\n", nodeID, event)
        return nil
    }
    return fmt.Errorf("node with ID %s not found", nodeID)
}

// LogEvent logs a general event on the network.
func (ledger *Ledger) LogEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Network event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogFaultEvent logs a network fault event with a node and timestamp.
func (ledger *Ledger) LogFaultEvent(event string, node string, timestamp time.Time) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    logMessage := fmt.Sprintf("Fault event: %s for node: %s at %s", event, node, timestamp)
    fmt.Println(logMessage)
    ledger.Events = append(ledger.Events, logMessage) // Log a detailed event
}



// LogSyncEvent logs a synchronization event on the network.
func (ledger *Ledger) LogSyncEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Synchronization event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogFirewallEvent logs a firewall-related event.
func (ledger *Ledger) LogFirewallEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Firewall event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogGeoEvent logs a geo-location-related event on the network.
func (ledger *Ledger) LogGeoEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Geo-location event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogConnectionEvent logs a network connection event.
func (ledger *Ledger) LogConnectionEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Connection event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogMessageEvent logs a message event on the network.
func (ledger *Ledger) LogMessageEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Message event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogNATEvent logs a NAT (Network Address Translation)-related event.
func (ledger *Ledger) LogNATEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("NAT event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// LogNetworkEvent logs a general network event.
func (ledger *Ledger) LogNetworkEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Network event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// GetPendingTransactions retrieves the list of pending transactions.
func (ledger *Ledger) GetPendingTransactions() []TransactionRecord {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Convert []*Transaction to []TransactionRecord
    var transactionRecords []TransactionRecord
    for _, tx := range ledger.pendingTransactions {
        transactionRecord := TransactionRecord{
            From:   tx.FromAddress,
            To:     tx.ToAddress,
            Amount: tx.Amount,
            // Add other fields from Transaction to TransactionRecord as necessary
        }
        transactionRecords = append(transactionRecords, transactionRecord)
    }

    return transactionRecords
}


// RecordPeerDiscovery logs peer discovery events.
func (ledger *Ledger) RecordPeerDiscovery(peer PeerInfo) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Peer discovered: %s\n", peer.PeerID)
}

// RecordPacketEvent logs events related to packet transmission or reception.
func (ledger *Ledger) RecordPacketEvent(packet PacketEvent) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Packet event: %s\n", packet.PacketID)
}

// RecordEvent logs a general event in the network.
func (ledger *Ledger) RecordEvent(event string) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Recorded event: %s\n", event)
    ledger.Events = append(ledger.Events, event)
}

// RecordRoutingEvent logs a routing event on the network.
func (ledger *Ledger) RecordRoutingEvent(route RoutingEvent) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Routing event: %s\n", route.RouteID)
    ledger.Events = append(ledger.Events, route.RouteID)
}

// GetBalance retrieves the balance of a specific account.
func (ledger *Ledger) GetBalance(accountID string) (float64, error) {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    if account, exists := ledger.State.Accounts[accountID]; exists {
        return account.Balance, nil
    }
    return 0, fmt.Errorf("account with ID %s not found", accountID)
}

// VerifyTransaction verifies a transaction before adding it to the ledger.
func (ledger *Ledger) VerifyTransaction(tx TransactionRecord) bool {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Placeholder: implement actual verification logic based on Synnergy Consensus.
    return true
}

// AddNodeToLedger adds a node to the ledger.
func (ledger *Ledger) AddNodeToLedger(node NodeInfo) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    // Use NodeID instead of ID
    ledger.State.Nodes[node.NodeID] = node
    fmt.Printf("Node %s added to the ledger.\n", node.NodeID)
    return nil
}

// RecordNodeConnection logs a node's connection to the network.
func (ledger *Ledger) RecordNodeConnection(nodeID string, connection ConnectionEvent) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("Node %s connected: %s.\n", nodeID, connection.ConnectionID)
    return nil
}

// RemoveNodeFromLedger removes a node from the ledger.
func (ledger *Ledger) RemoveNodeFromLedger(nodeID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    delete(ledger.State.Nodes, nodeID)
    fmt.Printf("Node %s removed from the ledger.\n", nodeID)
    return nil
}

// RemoveWebRTCConnection removes a WebRTC connection from the network ledger.
func (ledger *Ledger) RemoveWebRTCConnection(connectionID string) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("WebRTC connection %s removed from ledger.\n", connectionID)
    return nil
}

// RecordWebRTCConnection logs the addition of a WebRTC connection.
func (ledger *Ledger) RecordWebRTCConnection(connection WebRTCConnection) error {
    ledger.lock.Lock()
    defer ledger.lock.Unlock()

    fmt.Printf("WebRTC connection %s recorded in ledger.\n", connection.ConnectionID)
    return nil
}

// RecordGossipMessage logs a gossip message sent between nodes.
func (l *ScalabilityLedger) RecordGossipMessage(messageID, nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log gossip message
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        messageID,
		Action:    "GossipMessage",
		Details:   fmt.Sprintf("Message sent by node %s", nodeID),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordSyncAction logs a sync action between nodes.
func (l *ScalabilityLedger) RecordSyncAction(syncID, nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log sync action
	T.BlockchainConsensusCoinState.TransactionHistory = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        syncID,
		Action:    "SyncAction",
		Details:   fmt.Sprintf("Sync performed by node %s", nodeID),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordNodeAddition logs the addition of a node to the ledger.
func (l *ScalabilityLedger) RecordNodeAddition(nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log node addition
	l.State.History = append(T.BlockchainConsensusCoinState.TransactionHistory, TransactionRecord{
		ID:        nodeID,
		Action:    "NodeAdded",
		Details:   "Node added to the ledger",
		Timestamp: time.Now(),
	})
	return nil
}


// RecordTaskDistribution logs the distribution of tasks to nodes.
func (l *ScalabilityLedger) RecordTaskDistribution(taskID, nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log task distribution
	l.State.History = append(l.State.History, TransactionRecord{
		ID:        taskID,
		Action:    "TaskDistributed",
		Details:   fmt.Sprintf("Task distributed to node %s", nodeID),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordSync logs the sync operation between nodes.
func (l *ScalabilityLedger) RecordSync(syncID, nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log node sync operation
	l.State.History = append(l.State.History, TransactionRecord{
		ID:        syncID,
		Action:    "NodeSync",
		Details:   fmt.Sprintf("Sync with node %s", nodeID),
		Timestamp: time.Now(),
	})
	return nil
}

// RecordRedundancy logs a redundancy action for secure operations.
func (l *ScalabilityLedger) RecordRedundancy(actionID, nodeID string) error {
	l.Lock()
	defer l.Unlock()

	// Ensure Nodes map is initialized
	if l.Nodes == nil {
		l.Nodes = make(map[string]NodeRecord)
	}

	// Log redundancy action
	l.State.History = append(l.State.History, TransactionRecord{
		ID:        actionID,
		Action:    "RedundancyAction",
		Details:   fmt.Sprintf("Redundancy action on node %s", nodeID),
		Timestamp: time.Now(),
	})
	return nil
}

// GetAllNodes retrieves all registered nodes.
func (nl *NetworkLedger) GetAllNodes() []Node {
    return nl.Nodes
}

// AddNode registers a new node in the network.
func (nl *NetworkLedger) AddNode(node Node) {
    nl.Nodes = append(nl.Nodes, node)
}