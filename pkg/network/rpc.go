package network

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewRPCServer initializes a new encrypted RPC server
func NewRPCServer(address string, ledgerInstance *ledger.Ledger, router *Router) *RPCServer {
	return &RPCServer{
		Address:        address,
		LedgerInstance: ledgerInstance,
		Router:         router,
	}
}

// Start launches the RPC server with TLS encryption
func (server *RPCServer) Start(certFile, keyFile, caCertFile string) error {
	// Load CA certificate
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return fmt.Errorf("failed to load CA certificate: %v", err)
	}

	// Create CA pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// Create a server using HTTP with TLS
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/synnergy_rpc", server.handleRPCRequest)

	s := &http.Server{
		Addr:      server.Address,
		Handler:   serverMux,
		TLSConfig: tlsConfig,
	}

	fmt.Printf("Starting RPC server at %s...\n", server.Address)
	return s.ListenAndServeTLS(certFile, keyFile)
}

// handleRPCRequest handles incoming RPC requests
func (server *RPCServer) handleRPCRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming request
	var rpcRequest RPCRequest
	err := json.NewDecoder(r.Body).Decode(&rpcRequest)
	if err != nil {
		http.Error(w, "Failed to parse RPC request", http.StatusBadRequest)
		return
	}

	// Create encryption instance
	encryption := &common.Encryption{}

	// Validate request by decrypting the payload
	decryptedData, err := encryption.DecryptData([]byte(rpcRequest.Payload), []byte(rpcRequest.SenderPublicKey))
	if err != nil {
		http.Error(w, "Failed to decrypt RPC request", http.StatusUnauthorized)
		return
	}

	// Execute the RPC call
	var rpcResponse RPCResponse
	switch rpcRequest.Method {
	case "GetBalance":
		rpcResponse, err = server.getBalance(decryptedData)
	case "SendTransaction":
		rpcResponse, err = server.sendTransaction(decryptedData)
	case "GetBlock":
		rpcResponse, err = server.getBlock(decryptedData)
	default:
		http.Error(w, "Invalid RPC method", http.StatusNotFound)
		return
	}

	// Handle potential errors
	if err != nil {
		http.Error(w, fmt.Sprintf("RPC error: %v", err), http.StatusInternalServerError)
		return
	}

	// Encrypt the response before sending it back
	encryptedResponse, err := encryption.EncryptData("AES", []byte(rpcResponse.Data), []byte(rpcRequest.SenderPublicKey))
	if err != nil {
		http.Error(w, "Failed to encrypt RPC response", http.StatusInternalServerError)
		return
	}

	// Send the encrypted response back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RPCResponse{
		Data: string(encryptedResponse),
	})
}

// getBalance retrieves the balance of a given wallet address
func (server *RPCServer) getBalance(data []byte) (RPCResponse, error) {
	// Parse the incoming request
	var request BalanceRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("invalid balance request: %v", err)
	}

	// Retrieve the balance from the ledger, which returns balance and an error
	balance, err := server.LedgerInstance.GetBalance(request.WalletAddress)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to get balance: %v", err)
	}

	// Create the balance response
	response := BalanceResponse{
		WalletAddress: request.WalletAddress,
		Balance:       balance,
	}

	// Serialize the response
	responseData, err := json.Marshal(response)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to marshal balance response: %v", err)
	}

	// Return the response
	return RPCResponse{
		Data: string(responseData),
	}, nil
}

// sendTransaction processes and broadcasts a transaction across the network
func (server *RPCServer) sendTransaction(data []byte) (RPCResponse, error) {
    var tx common.Transaction
    err := json.Unmarshal(data, &tx)
    if err != nil {
        return RPCResponse{}, fmt.Errorf("invalid transaction data: %v", err)
    }

    // Convert common.Transaction to ledger.TransactionRecord
    ledgerTx := ledger.TransactionRecord{
        From:        tx.FromAddress,
        To:          tx.ToAddress,
        Amount:      tx.Amount,
        Fee:         tx.Fee,
        Hash:        tx.TransactionID,  // Use TransactionID as the hash
        Status:      "pending",         // Initial status
        Timestamp:   tx.Timestamp,
    }

    // Verify the transaction using ledger and encryption logic
    if !server.LedgerInstance.VerifyTransaction(ledgerTx) {
        return RPCResponse{}, fmt.Errorf("invalid transaction")
    }

    // Add the transaction to the pending pool and propagate it across the network
    server.LedgerInstance.AddTransaction(tx.FromAddress, tx.ToAddress, tx.Amount)

    server.Router.BroadcastRoutes()

    return RPCResponse{
        Data: "Transaction successfully broadcasted",
    }, nil
}

// getBlock retrieves a specific block from the blockchain ledger by index
func (server *RPCServer) getBlock(data []byte) (RPCResponse, error) {
	var request BlockRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("invalid block request: %v", err)
	}

	// Use GetBlockByIndex to retrieve the block
	block, err := server.LedgerInstance.GetBlockByIndex(request.BlockIndex)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("block not found: %v", err)
	}

	// Marshal the block into a response
	responseData, err := json.Marshal(block)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to marshal block data: %v", err)
	}

	return RPCResponse{
		Data: string(responseData),
	}, nil
}




// NewRPCClient initializes a new RPC client
func NewRPCClient(serverAddress, publicKey, privateKey string) *RPCClient {
	return &RPCClient{
		ServerAddress: serverAddress,
		PublicKey:     publicKey,
		PrivateKey:    privateKey,
	}
}

// Call makes a secure RPC call to the server
func (client *RPCClient) Call(method string, payload interface{}) (RPCResponse, error) {
	// Serialize the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create an encryption instance
	encryption := &common.Encryption{}

	// Convert the client's private key from string to []byte
	privateKeyBytes := []byte(client.PrivateKey)

	// Encrypt the payload using AES (or another encryption method)
	encryptedPayload, err := encryption.EncryptData("AES", payloadBytes, privateKeyBytes)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to encrypt payload: %v", err)
	}

	// Create an RPC request
	rpcRequest := RPCRequest{
		Method:          method,
		Payload:         string(encryptedPayload),
		SenderPublicKey: client.PublicKey,
	}

	// Send the RPC request
	data, err := json.Marshal(rpcRequest)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to marshal RPC request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("https://%s/synnergy_rpc", client.ServerAddress), "application/json", strings.NewReader(string(data)))
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to send RPC request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RPCResponse{}, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Parse the response
	var rpcResponse RPCResponse
	err = json.NewDecoder(resp.Body).Decode(&rpcResponse)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to parse RPC response: %v", err)
	}

	return rpcResponse, nil
}
