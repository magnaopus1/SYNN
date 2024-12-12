package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewServer initializes the server with ledger and consensus integration
func NewServer(address, encryptionKey string, ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus) *Server {
	return &Server{
		Address:         address,
		LedgerInstance:  ledgerInstance,
		ConsensusEngine: consensusEngine,
		EncryptionKey:   encryptionKey,
		connections:     make(map[string]net.Conn),
	}
}

// StartServer starts the server and begins accepting secure client connections
func (s *Server) StartServer(certFile, keyFile string) error {
	// Load TLS certificates for secure communication
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificates: %v", err)
	}

	// Listen for incoming connections using TLS
	listener, err := tls.Listen("tcp", s.Address, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		return fmt.Errorf("failed to start TLS listener: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s...\n", s.Address)

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

// handleConnection processes incoming connections securely
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Generate a unique ID for the connection
	clientAddr := conn.RemoteAddr().String()
	s.connectionLock.Lock()
	s.connections[clientAddr] = conn
	s.connectionLock.Unlock()

	fmt.Printf("Accepted connection from %s\n", clientAddr)

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Process encrypted messages from the client
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected.\n", clientAddr)
			} else {
				fmt.Printf("Error reading from client %s: %v\n", clientAddr, err)
			}
			break
		}

		encryptedData := buffer[:n]
		decryptedMessage, err := encryption.DecryptData(encryptedData, []byte(s.EncryptionKey)) // Convert key to []byte
		if err != nil {
			fmt.Printf("Error decrypting message from %s: %v\n", clientAddr, err)
			continue
		}

		fmt.Printf("Received secure message from %s: %s\n", clientAddr, string(decryptedMessage))

		// Handle the decrypted message
		err = s.ProcessMessage(decryptedMessage, conn)
		if err != nil {
			fmt.Printf("Failed to process message: %v\n", err)
		}
	}

	// Remove connection when done
	s.connectionLock.Lock()
	delete(s.connections, clientAddr)
	s.connectionLock.Unlock()
}


// processMessage processes a decrypted message from a client
func (s *Server) ProcessMessage(message []byte, conn net.Conn) error {
	var request NetworkRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		return fmt.Errorf("invalid message format: %v", err)
	}

	switch request.Type {
	case "SUB_BLOCK":
		return s.HandleSubBlockRequest(request.Payload, conn)
	case "BLOCK_FINALIZATION":
		return s.HandleBlockFinalizationRequest(request.Payload, conn)
	default:
		return fmt.Errorf("unknown message type: %s", request.Type)
	}
}


// handleSubBlockRequest processes a request to validate and add a sub-block
func (s *Server) HandleSubBlockRequest(payload json.RawMessage, conn net.Conn) error {
    var subBlock common.SubBlock
    err := json.Unmarshal(payload, &subBlock)
    if err != nil {
        return fmt.Errorf("invalid sub-block format: %v", err)
    }

    // Determine whether to use PoS or PoH for validation
    if s.ConsensusEngine.ShouldUsePoS(subBlock) {
        fmt.Printf("Validating sub-block #%d using PoS.\n", subBlock.Index)

        if !s.ConsensusEngine.PoS.ValidateSubBlock(subBlock) {
            return fmt.Errorf("PoS validation failed for sub-block #%d", subBlock.Index)
        }
    } else {
        fmt.Printf("Validating sub-block #%d using PoH.\n", subBlock.Index)

        pohProof := s.ConsensusEngine.PoH.GeneratePoHProof()
        if !s.ConsensusEngine.PoH.ValidatePoHProof(pohProof, subBlock.Validator) {
            return fmt.Errorf("PoH validation failed for sub-block #%d", subBlock.Index)
        }
    }

    // Convert common.SubBlock to ledger.SubBlock
    ledgerSubBlock := convertSubBlockToLedger(subBlock)

    // Add the sub-block to the ledger
    err = s.LedgerInstance.AddSubBlock(ledgerSubBlock)
    if err != nil {
        return fmt.Errorf("failed to add sub-block to ledger: %v", err)
    }

    fmt.Printf("Sub-block #%d added to the ledger.\n", subBlock.Index)
    return s.sendResponse(conn, "SUB_BLOCK_SUCCESS", "Sub-block added successfully.")
}



// convertSubBlockToLedger converts a common.SubBlock to ledger.SubBlock
func convertSubBlockToLedger(subBlock common.SubBlock) ledger.SubBlock {
	return ledger.SubBlock{
		SubBlockID:   subBlock.SubBlockID,
		Index:        subBlock.Index,
		Timestamp:    subBlock.Timestamp,
		Transactions: convertTransactionsToLedger(subBlock.Transactions), // Ensure transactions are converted
		Validator:    subBlock.Validator,
		PrevHash:     subBlock.PrevHash,
		Hash:         subBlock.Hash,
		PoHProof:     convertPoHProofToLedger(subBlock.PoHProof), // Convert PoHProof if needed
		Status:       subBlock.Status,
		Signature:    subBlock.Signature,
	}
}

// convertTransactionsToLedger converts a slice of common.Transaction to ledger.Transaction
func convertTransactionsToLedger(transactions []common.Transaction) []ledger.Transaction {
	var ledgerTransactions []ledger.Transaction
	for _, tx := range transactions {
		ledgerTx := ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Amount:          tx.Amount,
			Fee:             tx.Fee,
			TokenStandard:   tx.TokenStandard,
			TokenID:         tx.TokenID,
			Timestamp:       tx.Timestamp,
			SubBlockID:      tx.SubBlockID,
			BlockID:         tx.BlockID,
			ValidatorID:     tx.ValidatorID,
			Signature:       tx.Signature,
			Status:          tx.Status,
			EncryptedData:   tx.EncryptedData,
			DecryptedData:   tx.DecryptedData,
			ExecutionResult: tx.ExecutionResult,
			FrozenAmount:    tx.FrozenAmount,
			RefundAmount:    tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		}
		ledgerTransactions = append(ledgerTransactions, ledgerTx)
	}
	return ledgerTransactions
}


// convertPoHProofToLedger converts common.PoHProof to ledger.PoHProof
func convertPoHProofToLedger(pohProof common.PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Sequence:  pohProof.Sequence,
		Timestamp: pohProof.Timestamp,
		Hash:      pohProof.Hash,
	}
}


// handleBlockFinalizationRequest processes a request to finalize a block
func (s *Server) HandleBlockFinalizationRequest(payload json.RawMessage, conn net.Conn) error {
	var block common.Block
	err := json.Unmarshal(payload, &block)
	if err != nil {
		return fmt.Errorf("invalid block format: %v", err)
	}

	// Finalize the block using Proof of Work (PoW)
	s.ConsensusEngine.FinalizeBlock() // No need to assign the result, just call it

	// Convert common.Block to ledger.Block before adding to the ledger
	ledgerBlock := convertBlockToLedgerBlock(block)

	// Add the finalized block to the ledger
	err = s.LedgerInstance.FinalizeBlock(ledgerBlock)
	if err != nil {
		return fmt.Errorf("failed to finalize block in ledger: %v", err)
	}

	fmt.Printf("Block #%d finalized and added to the ledger.\n", block.Index)
	return s.sendResponse(conn, "BLOCK_FINALIZATION_SUCCESS", "Block finalized successfully.")
}



// convertBlockToLedgerBlock converts a common.Block to ledger.Block
func convertBlockToLedgerBlock(commonBlock common.Block) ledger.Block {
	return ledger.Block{
		BlockID:    commonBlock.BlockID,
		Index:      commonBlock.Index,
		Timestamp:  commonBlock.Timestamp,
		SubBlocks:  convertSubBlocks(commonBlock.SubBlocks), // Assuming you need to convert sub-blocks too
		PrevHash:   commonBlock.PrevHash,
		Hash:       commonBlock.Hash,
		Nonce:      commonBlock.Nonce,
		Difficulty: commonBlock.Difficulty,
		MinerReward: commonBlock.MinerReward,
		Validators: commonBlock.Validators,
		Status:     commonBlock.Status,
	}
}


// convertSubBlocks converts a slice of common.SubBlock to ledger.SubBlock
func convertSubBlocks(commonSubBlocks []common.SubBlock) []ledger.SubBlock {
	var ledgerSubBlocks []ledger.SubBlock
	for _, subBlock := range commonSubBlocks {
		ledgerSubBlocks = append(ledgerSubBlocks, ledger.SubBlock{
			SubBlockID:  subBlock.SubBlockID,
			Index:       subBlock.Index,
			Timestamp:   subBlock.Timestamp,
			Transactions: convertTransactions(subBlock.Transactions), // Convert transactions
			Validator:   subBlock.Validator,
			PrevHash:    subBlock.PrevHash,
			Hash:        subBlock.Hash,
			PoHProof:    convertPoHProof(subBlock.PoHProof), // Convert PoHProof
			Status:      subBlock.Status,
			Signature:   subBlock.Signature,
		})
	}
	return ledgerSubBlocks
}


// convertTransactions converts a slice of common.Transaction to a slice of ledger.Transaction
func convertTransactions(commonTxs []common.Transaction) []ledger.Transaction {
	var ledgerTxs []ledger.Transaction
	for _, tx := range commonTxs {
		ledgerTx := ledger.Transaction{
			TransactionID:   tx.TransactionID,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Amount:          tx.Amount,
			Fee:             tx.Fee,
			TokenStandard:   tx.TokenStandard,
			TokenID:         tx.TokenID,
			Timestamp:       tx.Timestamp,
			SubBlockID:      tx.SubBlockID,
			BlockID:         tx.BlockID,
			ValidatorID:     tx.ValidatorID,
			Signature:       tx.Signature,
			Status:          tx.Status,
			EncryptedData:   tx.EncryptedData,
			DecryptedData:   tx.DecryptedData,
			ExecutionResult: tx.ExecutionResult,
			FrozenAmount:    tx.FrozenAmount,
			RefundAmount:    tx.RefundAmount,
			ReversalRequested: tx.ReversalRequested,
		}
		ledgerTxs = append(ledgerTxs, ledgerTx)
	}
	return ledgerTxs
}


// convertPoHProof converts a common.PoHProof to a ledger.PoHProof
func convertPoHProof(commonPoH common.PoHProof) ledger.PoHProof {
	return ledger.PoHProof{
		Sequence:  commonPoH.Sequence,
		Timestamp: commonPoH.Timestamp,
		Hash:      commonPoH.Hash,
	}
}


// sendResponse sends an encrypted response to the client
func (s *Server) sendResponse(conn net.Conn, responseType, message string) error {
	response := NetworkResponse{
		Type:    responseType,
		Message: message,
	}

	// Marshal the response into JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Encrypt the response data using AES
	encryptedResponse, err := encryption.EncryptData("AES", responseData, []byte(s.EncryptionKey)) // Added the algorithm "AES"
	if err != nil {
		return fmt.Errorf("failed to encrypt response: %v", err)
	}

	// Send the encrypted response over the network
	_, err = conn.Write(encryptedResponse)
	if err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}



// StopServer gracefully stops the server by closing all connections
func (s *Server) StopServer() {
	s.connectionLock.Lock()
	defer s.connectionLock.Unlock()

	for clientAddr, conn := range s.connections {
		conn.Close()
		fmt.Printf("Closed connection with client %s\n", clientAddr)
	}

	fmt.Println("Server stopped.")
	os.Exit(0)
}

// createSelfSignedCert generates a self-signed certificate for testing
func CreateSelfSignedCert(certFile, keyFile string) error {
	// Generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// Set up certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Organization"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour), // Valid for 1 day

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Save the certificate to certFile
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	// Save the private key to keyFile
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	return nil
}