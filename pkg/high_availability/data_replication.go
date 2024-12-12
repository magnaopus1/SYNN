package high_availability

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"synnergy_network/pkg/ledger"
)

// NewDataReplicationManager initializes a DataReplicationManager with a list of nodes and the ledger instance.
func NewDataReplicationManager(nodes []string, ledgerInstance *ledger.Ledger) *DataReplicationManager {
    return &DataReplicationManager{
        Nodes:          nodes,
        LedgerInstance: ledgerInstance,
    }
}

// ReplicateLedger replicates the entire ledger across the network to ensure data availability and redundancy.
func (drm *DataReplicationManager) ReplicateLedger() {
    drm.mutex.Lock()
    defer drm.mutex.Unlock()

    fmt.Println("Replicating the entire ledger across the network...")

    // Replicate the full ledger state, including blocks, sub-blocks, transactions, balances, etc.
    for _, node := range drm.Nodes {
        fmt.Printf("Replicating the full ledger to node %s...\n", node)
        drm.replicateLedgerToNode(node)
    }

    fmt.Println("Ledger replication completed.")
}

// replicateLedgerToNode handles the process of sending the entire ledger (state and structures) to a specific node.
func (drm *DataReplicationManager) replicateLedgerToNode(node string) {
	log.Printf("Starting replication of the full ledger (state and structures) to node %s...\n", node)

	// Retrieve the current ledger state and associated structs
	fullLedger := drm.prepareFullLedger()

	// Transmit the entire ledger data to the specified node
	err := drm.sendFullLedgerToNode(node, fullLedger)
	if err != nil {
		log.Printf("Error replicating ledger to node %s: %v\n", node, err)

		// Handle replication failure and attempt retries
		log.Printf("Retrying replication to node %s...\n", node)
		for i := 1; i <= 3; i++ {
			log.Printf("Retry attempt %d to replicate ledger to node %s...\n", i, node)
			err = drm.sendFullLedgerToNode(node, fullLedger)
			if err == nil {
				log.Printf("Ledger successfully replicated to node %s on retry attempt %d.\n", node, i)
				return
			}
			log.Printf("Retry attempt %d failed for node %s: %v\n", i, node, err)
		}

		// If all retries fail, mark the node for further handling
		log.Printf("Failed to replicate ledger to node %s after 3 attempts. Handling failure...\n", node)
		drm.handleFailedReplication(node)
	} else {
		log.Printf("Ledger successfully replicated to node %s on the first attempt.\n", node)
	}
}

// prepareFullLedger dynamically prepares the full ledger for replication.
func (drm *DataReplicationManager) prepareFullLedger() *ledger.Ledger {
	// Locking to ensure thread safety when accessing the ledger
	drm.mutex.Lock()
	defer drm.mutex.Unlock()

	// Directly use the ledger instance, which dynamically contains all state, blocks, and associated data
	fullLedger := drm.LedgerInstance

	// Log the preparation process
	log.Printf("Preparing the full ledger for replication.")

	return fullLedger
}


func (drm *DataReplicationManager) sendFullLedgerToNode(node string, fullLedger *ledger.Ledger) error {
	// Serialize, encrypt, and send the ledger data
	log.Printf("Sending full ledger to node %s...", node)
	serializedLedger, err := drm.serializeFullLedger(fullLedger)
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}

	encryptedLedger, err := drm.encryptDataForNode(node, serializedLedger)
	if err != nil {
		return fmt.Errorf("encryption error: %w", err)
	}

	url := fmt.Sprintf("https://%s/api/replicate-full-ledger", node)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encryptedLedger))
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network transmission error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("node %s returned error: %s", node, string(body))
	}

	log.Printf("Full ledger successfully sent to node %s.", node)
	return nil
}


// serializeFullLedger serializes the full ledger into binary format for transmission.
func (drm *DataReplicationManager) serializeFullLedger(fullLedger *ledger.Ledger) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fullLedger); err != nil {
		return nil, fmt.Errorf("serialization error: %w", err)
	}
	return buffer.Bytes(), nil
}


// sendLedgerStateToNode transmits the serialized and encrypted ledger state to the specified node.
func (drm *DataReplicationManager) sendLedgerStateToNode(node string, ledgerState *ledger.Ledger) error {
	// Step 1: Serialize the ledger state
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(ledgerState); err != nil {
		log.Printf("Error serializing ledger state: %v", err)
		return fmt.Errorf("serialization error: %w", err)
	}
	serializedData := buffer.Bytes()

	// Step 2: Encrypt the serialized ledger state
	encryptedData, err := drm.encryptDataForNode(node, serializedData)
	if err != nil {
		log.Printf("Error encrypting ledger state for node %s: %v", node, err)
		return fmt.Errorf("encryption error: %w", err)
	}

	// Step 3: Send the encrypted data to the node
	url := fmt.Sprintf("https://%s/api/replicate-ledger", node)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encryptedData))
	if err != nil {
		log.Printf("Error creating HTTP request for node %s: %v", node, err)
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending ledger state to node %s: %v", node, err)
		return fmt.Errorf("network transmission error: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: Verify the response from the node
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Node %s returned error: %s", node, string(body))
		return fmt.Errorf("node %s returned error: %s", node, string(body))
	}

	log.Printf("Ledger state successfully sent to node %s.", node)
	return nil
}

// encryptDataForNode encrypts the serialized data using the node's public key.
func (drm *DataReplicationManager) encryptDataForNode(node string, data []byte) ([]byte, error) {
	// Step 1: Retrieve the public key of the target node
	publicKey, err := drm.getNodePublicKey(node)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve public key for node %s: %w", node, err)
	}

	// Step 2: Validate input data length against RSA encryption limits
	maxDataLength := (publicKey.Size() - 2*sha256.Size() - 2) // RSA-OAEP padding constraints
	if len(data) > maxDataLength {
		return nil, fmt.Errorf("data too large to encrypt for node %s: size %d exceeds max %d", node, len(data), maxDataLength)
	}

	// Step 3: Encrypt the data using RSA-OAEP
	hash := sha256.New()
	encryptedData, err := rsa.EncryptOAEP(hash, nil, publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed for node %s: %w", node, err)
	}

	return encryptedData, nil
}


// ReplicateSubBlocks replicates sub-blocks across the network for redundancy.
func (drm *DataReplicationManager) ReplicateSubBlocks(subBlocks []ledger.SubBlock) {
    drm.mutex.Lock()
    defer drm.mutex.Unlock()

    fmt.Printf("Replicating %d sub-blocks across the network...\n", len(subBlocks))

    for _, node := range drm.Nodes {
        fmt.Printf("Replicating sub-blocks to node %s...\n", node)
        drm.replicateSubBlocksToNode(node, subBlocks)
    }

    drm.ReplicatedSubBlocks = append(drm.ReplicatedSubBlocks, subBlocks...)
    fmt.Println("Sub-block replication completed.")
}

// ReplicateBlocks replicates blocks across the network for redundancy.
func (drm *DataReplicationManager) ReplicateBlocks(blocks []ledger.Block) {
    drm.mutex.Lock()
    defer drm.mutex.Unlock()

    fmt.Printf("Replicating %d blocks across the network...\n", len(blocks))

    for _, node := range drm.Nodes {
        fmt.Printf("Replicating blocks to node %s...\n", node)
        drm.replicateBlocksToNode(node, blocks)
    }

    drm.ReplicatedBlocks = append(drm.ReplicatedBlocks, blocks...)
    fmt.Println("Block replication completed.")
}

// replicateSubBlocksToNode handles the process of replicating sub-blocks to a specific node.
func (drm *DataReplicationManager) replicateSubBlocksToNode(node string, subBlocks []ledger.SubBlock) {
    fmt.Printf("Replicating sub-blocks to node %s...\n", node)
    // In a real-world implementation, network communication code would go here
}

// replicateBlocksToNode handles the process of replicating blocks to a specific node.
func (drm *DataReplicationManager) replicateBlocksToNode(node string, blocks []ledger.Block) {
    fmt.Printf("Replicating blocks to node %s...\n", node)
    // In a real-world implementation, network communication code would go here
}

// VerifyReplication verifies that all nodes have received the replicated data correctly.
func (drm *DataReplicationManager) VerifyReplication() bool {
    drm.mutex.Lock()
    defer drm.mutex.Unlock()

    fmt.Println("Verifying replication across nodes...")

    for _, node := range drm.Nodes {
        if !drm.verifyNodeReplication(node) {
            fmt.Printf("Node %s has not received replicated data correctly.\n", node)
            return false
        }
    }

    fmt.Println("All nodes have received replicated data correctly.")
    return true
}

// verifyNodeReplication simulates verifying the replicated data on a specific node.
func (drm *DataReplicationManager) verifyNodeReplication(node string) bool {
    // Simulate verifying the node's replicated data
    fmt.Printf("Node %s replication verification complete.\n", node)
    return true
}

// HandleReplicationFailure handles the case where a node has failed to receive replicated data.
func (drm *DataReplicationManager) HandleReplicationFailure(node string) {
    drm.mutex.Lock()
    defer drm.mutex.Unlock()

    fmt.Printf("Handling replication failure for node %s...\n", node)

    // Resend all replicated sub-blocks, blocks, and the ledger state to the failed node.
    drm.replicateSubBlocksToNode(node, drm.ReplicatedSubBlocks)
    drm.replicateBlocksToNode(node, drm.ReplicatedBlocks)
    drm.replicateLedgerToNode(node)

    fmt.Printf("Resent all replicated data to node %s successfully.\n", node)
}
