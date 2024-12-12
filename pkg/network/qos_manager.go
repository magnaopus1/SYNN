package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"sort"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewQoSManager initializes a new QoS Manager with default parameters
func NewQoSManager(bandwidthLimit int, nodePublicKey string, ledgerInstance *ledger.Ledger) *QoSManager {
	return &QoSManager{
		BandwidthLimit: bandwidthLimit,
		PriorityQueues: make(map[int][]*QoSPacket),
		LedgerInstance: ledgerInstance,
		NodePublicKey:  nodePublicKey,
		ConnectedPeers: make(map[string]*Peer),
	}
}

// AddPeer adds a new peer to the list of connected peers
func (qos *QoSManager) AddPeer(address, publicKey string) {
	qos.Mutex.Lock()
	defer qos.Mutex.Unlock()

	if _, exists := qos.ConnectedPeers[address]; !exists {
		qos.ConnectedPeers[address] = &Peer{Address: address, PublicKey: publicKey}
		fmt.Printf("Added peer %s to the network.\n", address)

		// Create a PeerInfo object and pass it to RecordPeerDiscovery
		peerInfo := ledger.PeerInfo{
			Address:   address,
			PublicKey: publicKey,
		}
		qos.LedgerInstance.RecordPeerDiscovery(peerInfo)
	}
}


// SendPacket sends an encrypted packet to a peer, respecting QoS settings
func (qos *QoSManager) SendPacket(payload []byte, priority int, destination string) error {
	qos.Mutex.Lock()
	defer qos.Mutex.Unlock()

	// Ensure the packet has valid priority and destination
	if priority < 0 || priority > 10 {
		return fmt.Errorf("invalid priority level")
	}
	peer, exists := qos.ConnectedPeers[destination]
	if !exists {
		return fmt.Errorf("peer not found: %s", destination)
	}

	// Create an instance of Encryption
	encryption := &common.Encryption{}

	// Encrypt the payload using the peer's public key (or AES key)
	encryptedPayload, err := encryption.EncryptData("AES", payload, []byte(peer.PublicKey))  // Assuming the PublicKey is a valid AES key
	if err != nil {
		return fmt.Errorf("failed to encrypt payload: %v", err)
	}

	// Create a new QoS packet and add it to the queue
	packet := &QoSPacket{
		Payload:     encryptedPayload,
		Priority:    priority,
		Timestamp:   time.Now(),
		Destination: destination,
	}
	qos.PriorityQueues[priority] = append(qos.PriorityQueues[priority], packet)
	fmt.Printf("Added packet to priority %d queue for %s.\n", priority, destination)

	packetEvent := ledger.PacketEvent{
		PacketID:    "unique-packet-id", 
		Timestamp:   packet.Timestamp,
		Data:        encryptedPayload,
		Destination: destination,
		Size:        len(payload),
		Priority:    priority,
	}
	qos.LedgerInstance.RecordPacketEvent(packetEvent)

	return nil
}




// ProcessPackets processes the queued packets according to priority and bandwidth limits
func (qos *QoSManager) ProcessPackets() {
	qos.Mutex.Lock()
	defer qos.Mutex.Unlock()

	// Sort packets by priority (higher priority first)
	var allPackets []*QoSPacket
	for _, queue := range qos.PriorityQueues {
		allPackets = append(allPackets, queue...)
	}
	sort.Slice(allPackets, func(i, j int) bool {
		return allPackets[i].Priority < allPackets[j].Priority
	})

	totalBytesSent := 0

	// Process packets until bandwidth limit is reached
	for _, packet := range allPackets {
		if totalBytesSent >= qos.BandwidthLimit*1024 { // Bandwidth limit in KBps
			break
		}

		// Send the packet to its destination
		conn, err := net.Dial("tcp", packet.Destination)
		if err != nil {
			fmt.Printf("Failed to connect to peer %s: %v\n", packet.Destination, err)
			continue
		}

		_, err = conn.Write(packet.Payload)
		if err != nil {
			fmt.Printf("Failed to send packet to peer %s: %v\n", packet.Destination, err)
			conn.Close()
			continue
		}
		conn.Close()

		totalBytesSent += len(packet.Payload)
		fmt.Printf("Sent packet to %s (Priority: %d, Bytes: %d)\n", packet.Destination, packet.Priority, len(packet.Payload))

		// Remove the packet from the queue
		qos.removePacketFromQueue(packet)
	}
}

// removePacketFromQueue removes a processed packet from the priority queue
func (qos *QoSManager) removePacketFromQueue(packet *QoSPacket) {
	for i, p := range qos.PriorityQueues[packet.Priority] {
		if p == packet {
			qos.PriorityQueues[packet.Priority] = append(qos.PriorityQueues[packet.Priority][:i], qos.PriorityQueues[packet.Priority][i+1:]...)
			break
		}
	}
}

// MonitorQoS monitors the network for QoS violations and takes appropriate actions
func (qos *QoSManager) MonitorQoS() {
	qos.Mutex.Lock()
	defer qos.Mutex.Unlock()

	// Monitor each peer's bandwidth usage and take actions if they exceed the bandwidth limit
	for _, peer := range qos.ConnectedPeers {
		bandwidthUsed := qos.calculatePeerBandwidthUsage(peer.Address)
		if bandwidthUsed > qos.BandwidthLimit {
			fmt.Printf("Warning: Peer %s exceeded bandwidth limit (%d KBps)\n", peer.Address, qos.BandwidthLimit)
			// Optionally, you could reduce the priority of this peer's packets or take other action
		}
	}
}

// calculatePeerBandwidthUsage calculates the bandwidth usage for a given peer
func (qos *QoSManager) calculatePeerBandwidthUsage(peerAddress string) int {
	totalBytes := 0
	for _, queue := range qos.PriorityQueues {
		for _, packet := range queue {
			if packet.Destination == peerAddress {
				totalBytes += len(packet.Payload)
			}
		}
	}
	return totalBytes / 1024 // Return in KB
}

// RecordQoSViolation records any QoS violations in the ledger
func (qos *QoSManager) RecordQoSViolation(peerAddress string, violationType string) {
	// Create a violation record string
	violationRecord := fmt.Sprintf("QoS violation: %s by peer %s at %s", violationType, peerAddress, time.Now().Format(time.RFC3339))

	// Record the violation event in the ledger (pass only the violationRecord string)
	qos.LedgerInstance.RecordEvent(violationRecord)

	// Optionally, print the violation record
	fmt.Println(violationRecord)
}


// HashQoSPacket creates a hash of the QoS packet for integrity verification
func (qos *QoSManager) HashQoSPacket(packet *QoSPacket) string {
	data := fmt.Sprintf("%s%d%s", packet.Destination, packet.Priority, packet.Timestamp.String())
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
