package network

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewRouter initializes a new Router
func NewRouter(ledgerInstance *ledger.Ledger, nodePublicKey string, maxIdleTime time.Duration) *Router {
	return &Router{
		Routes:          make(map[string]string),
		Mutex:           sync.Mutex{},
		LedgerInstance:  ledgerInstance,
		EncryptedRoutes: make(map[string][]byte),
		NodePublicKey:   nodePublicKey,
		Peers:           make(map[string]*Peer),
		ConnectionPool:  NewConnectionPool(maxIdleTime), // Pass maxIdleTime as argument
	}
}


func (r *Router) AddRoute(peerID string, address string, peerPublicKey string) error {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()

    // Create an encryption instance from the common package
    encryption := &common.Encryption{}

    // Encrypt the route using the peer's public key
    encryptedRoute, err := encryption.EncryptData("AES", []byte(address), []byte(peerPublicKey)) // AES encryption example
    if err != nil {
        return fmt.Errorf("failed to encrypt route for peer %s: %v", peerID, err)
    }

    // Store the route and encrypted route in the router's internal maps
    r.Routes[peerID] = address
    r.EncryptedRoutes[peerID] = encryptedRoute

    // Add peer information
    r.Peers[peerID] = &Peer{Address: address, PublicKey: peerPublicKey}

    // Create a RoutingEvent and record the route addition in the ledger
    routingEvent := ledger.RoutingEvent{
        RouteID:    fmt.Sprintf("route-%s", peerID),
        SourceNode: r.NodePublicKey,
        DestNode:   peerID,
        Timestamp:  time.Now(),
    }

    r.LedgerInstance.RecordRoutingEvent(routingEvent)
    fmt.Printf("Added route for peer %s at address %s\n", peerID, address)

    return nil
}



// GetRoute retrieves the decrypted route for a peer
func (r *Router) GetRoute(peerID string) (string, error) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()

    // Check if the route exists for the peer
    encryptedRoute, exists := r.EncryptedRoutes[peerID]
    if !exists {
        return "", errors.New("route not found")
    }

    // Create an encryption instance
    encryption := &common.Encryption{}

    // Decrypt the route using the peer's public key
    decryptedRoute, err := encryption.DecryptData(encryptedRoute, []byte(r.Peers[peerID].PublicKey))
    if err != nil {
        return "", fmt.Errorf("failed to decrypt route for peer %s: %v", peerID, err)
    }

    return string(decryptedRoute), nil
}


// RemoveRoute removes a route from the routing table
func (r *Router) RemoveRoute(peerID string) {
    r.Mutex.Lock()
    defer r.Mutex.Unlock()

    delete(r.Routes, peerID)
    delete(r.EncryptedRoutes, peerID)
    delete(r.Peers, peerID)

    // Create a RoutingEvent for the removal
    routingEvent := ledger.RoutingEvent{
        RouteID:    fmt.Sprintf("route-%s", peerID),
        SourceNode: r.NodePublicKey,
        DestNode:   peerID,
        Timestamp:  time.Now(),
    }

    // Record the route removal in the ledger
    r.LedgerInstance.RecordRoutingEvent(routingEvent)
    fmt.Printf("Removed route for peer %s\n", peerID)
}


// RoutePacket routes a packet to the destination peer by looking up the route and sending the packet
func (r *Router) RoutePacket(packet *Packet) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	route, err := r.GetRoute(packet.DestinationID)
	if err != nil {
		return fmt.Errorf("failed to get route for peer %s: %v", packet.DestinationID, err)
	}

	// Send the packet to the destination via the connection pool
	err = r.ConnectionPool.SendPacket(route, packet)
	if err != nil {
		return fmt.Errorf("failed to route packet to %s: %v", packet.DestinationID, err)
	}

	fmt.Printf("Packet routed to peer %s at address %s\n", packet.DestinationID, route)
	return nil
}

// ValidateRouting ensures that all routes are still valid by checking connectivity with peers
func (r *Router) ValidateRouting() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for peerID, address := range r.Routes {
		// Attempt to connect to each peer to verify the route is still active
		conn, err := net.DialTimeout("tcp", address, time.Second*2)
		if err != nil {
			fmt.Printf("Route validation failed for peer %s: %v\n", peerID, err)
			// Optionally, remove the route if validation fails
			r.RemoveRoute(peerID)
			continue
		}
		conn.Close()
		fmt.Printf("Route to peer %s is valid\n", peerID)
	}
}

// BroadcastRoutes sends the routing table to all connected peers
func (r *Router) BroadcastRoutes() error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	// Serialize and encrypt the routing table
	routingTable, err := r.serializeRoutes()
	if err != nil {
		return fmt.Errorf("failed to serialize routing table: %v", err)
	}

	for peerID, peer := range r.Peers {
		// Create an encryption instance
		encryption := &common.Encryption{}

		// Encrypt the routing table for the peer
		encryptedTable, err := encryption.EncryptData("AES", []byte(routingTable), []byte(peer.PublicKey)) 
		if err != nil {
			fmt.Printf("Failed to encrypt routing table for peer %s: %v\n", peerID, err)
			continue
		}

		// Send the encrypted routing table to the peer in a packet
		packet := &Packet{
			SourceID:      r.NodePublicKey,
			DestinationID: peerID,
			Data:          encryptedTable, // Use Data field to store the encrypted data
		}

		err = r.ConnectionPool.SendPacket(peer.Address, packet)
		if err != nil {
			fmt.Printf("Failed to broadcast routing table to peer %s: %v\n", peerID, err)
		} else {
			fmt.Printf("Broadcasted routing table to peer %s\n", peerID)
		}
	}

	return nil
}


// serializeRoutes serializes the routing table for broadcasting
func (r *Router) serializeRoutes() ([]byte, error) {
	routeData := ""
	for peerID, address := range r.Routes {
		routeData += fmt.Sprintf("%s:%s\n", peerID, address)
	}

	// Hash the routing data for integrity
	hash := sha256.New()
	hash.Write([]byte(routeData))
	routeHash := hex.EncodeToString(hash.Sum(nil))

	return []byte(routeData + "\nHash:" + routeHash), nil
}
