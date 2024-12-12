package authority_node

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"synnergy_network_demo/encryption"
	"synnergy_network_demo/network"
	"synnergy_network_demo/common"
)


// NewAuthorityNode initializes a new authority node with a specific type and permissions.
func NewAuthorityNode(nodeID string, nodeType AuthorityNodeType, permissions *common.PermissionSet, encryptionService *encryption.Encryption, networkManager *network.NetworkManager) *AuthorityNode {
	return &AuthorityNode{
		NodeID:            nodeID,
		NodeType:          nodeType,
		Permissions:       permissions,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		RequestList:       []*common.Request{},
	}
}

// ViewRequestsList allows the authority node to view all requests sent to it based on its permissions.
func (an *AuthorityNode) ViewRequestsList() ([]*common.Request, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Check if the node has permission to view requests.
	if !an.Permissions.CanViewRequests {
		return nil, errors.New("authority node does not have permission to view requests")
	}

	fmt.Printf("Node %s of type %s is viewing its request list.\n", an.NodeID, an.NodeType)
	return an.RequestList, nil
}

// AddRequest adds a new request to the authority node's request list.
func (an *AuthorityNode) AddRequest(req *common.Request) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Add the request to the request list.
	an.RequestList = append(an.RequestList, req)
	fmt.Printf("New request %s added to node %s of type %s.\n", req.RequestID, an.NodeID, an.NodeType)
	return nil
}

// RefreshRequestsList updates the request list by fetching new requests from the network.
func (an *AuthorityNode) RefreshRequestsList() error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Simulate fetching new requests for this node type from the network.
	newRequests, err := an.NetworkManager.FetchRequestsForNode(an.NodeID, an.NodeType)
	if err != nil {
		return fmt.Errorf("failed to fetch requests for node %s: %v", an.NodeID, err)
	}

	// Append the new requests to the existing request list.
	an.RequestList = append(an.RequestList, newRequests...)
	fmt.Printf("Request list updated for node %s of type %s. Total requests: %d\n", an.NodeID, an.NodeType, len(an.RequestList))

	return nil
}

// DisplayRequests prints the list of requests to the console (or another display method).
func (an *AuthorityNode) DisplayRequests() {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	fmt.Printf("Displaying requests for node %s of type %s:\n", an.NodeID, an.NodeType)
	for _, req := range an.RequestList {
		fmt.Printf("Request ID: %s, Type: %s, Received At: %s\n", req.RequestID, req.Type, req.ReceivedAt.Format(time.RFC3339))
	}
}
