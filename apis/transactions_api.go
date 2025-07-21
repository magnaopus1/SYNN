package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"synnergy_network/pkg/transactions"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"

	"github.com/gorilla/mux"
)

// TransactionsAPI handles all transaction-related API endpoints
type TransactionsAPI struct {
	TransactionPool *transactions.TransactionPool
	LedgerInstance  *ledger.Ledger
}

// NewTransactionsAPI creates a new transactions API instance
func NewTransactionsAPI(txPool *transactions.TransactionPool, ledgerInstance *ledger.Ledger) *TransactionsAPI {
	return &TransactionsAPI{
		TransactionPool: txPool,
		LedgerInstance:  ledgerInstance,
	}
}

// RegisterRoutes registers all transaction API routes
func (api *TransactionsAPI) RegisterRoutes(router *mux.Router) {
	// Transaction pool management
	router.HandleFunc("/transactions/pool/add", api.AddTransaction).Methods("POST")
	router.HandleFunc("/transactions/pool/remove", api.RemoveTransaction).Methods("DELETE")
	router.HandleFunc("/transactions/pool/get", api.GetTransaction).Methods("GET")
	router.HandleFunc("/transactions/pool/list", api.ListTransactions).Methods("GET")
	router.HandleFunc("/transactions/pool/size", api.GetPoolSize).Methods("GET")
	router.HandleFunc("/transactions/pool/clear", api.ClearPool).Methods("POST")
	
	// Sub-block management
	router.HandleFunc("/transactions/subblocks/create", api.CreateSubBlock).Methods("POST")
	router.HandleFunc("/transactions/subblocks/add-to-ledger", api.AddSubBlockToLedger).Methods("POST")
	router.HandleFunc("/transactions/subblocks/list", api.ListPendingSubBlocks).Methods("GET")
	
	// Transaction status and history
	router.HandleFunc("/transactions/status/{txID}", api.GetTransactionStatus).Methods("GET")
	router.HandleFunc("/transactions/history/{address}", api.GetTransactionHistory).Methods("GET")
}

// AddTransaction adds a new transaction to the pool
func (api *TransactionsAPI) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction common.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.TransactionPool.AddTransaction(&transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Transaction added to pool successfully",
		"transaction_id": transaction.TransactionID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RemoveTransaction removes a transaction from the pool
func (api *TransactionsAPI) RemoveTransaction(w http.ResponseWriter, r *http.Request) {
	txID := r.URL.Query().Get("txID")
	if txID == "" {
		http.Error(w, "txID parameter is required", http.StatusBadRequest)
		return
	}

	api.TransactionPool.RemoveTransaction(txID)

	response := map[string]interface{}{
		"success": true,
		"message": "Transaction removed from pool successfully",
		"transaction_id": txID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTransaction retrieves a transaction from the pool by ID
func (api *TransactionsAPI) GetTransaction(w http.ResponseWriter, r *http.Request) {
	txID := r.URL.Query().Get("txID")
	if txID == "" {
		http.Error(w, "txID parameter is required", http.StatusBadRequest)
		return
	}

	transaction, err := api.TransactionPool.GetTransaction(txID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction: %v", err), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"transaction": transaction,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListTransactions returns all transactions in the pool
func (api *TransactionsAPI) ListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := api.TransactionPool.ListTransactions()

	response := map[string]interface{}{
		"success": true,
		"transactions": transactions,
		"count": len(transactions),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPoolSize returns the current size of the transaction pool
func (api *TransactionsAPI) GetPoolSize(w http.ResponseWriter, r *http.Request) {
	size := api.TransactionPool.PoolSize()

	response := map[string]interface{}{
		"success": true,
		"pool_size": size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ClearPool clears all transactions from the pool
func (api *TransactionsAPI) ClearPool(w http.ResponseWriter, r *http.Request) {
	api.TransactionPool.ClearPool()

	response := map[string]interface{}{
		"success": true,
		"message": "Transaction pool cleared successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSubBlock creates a new sub-block from pending transactions
func (api *TransactionsAPI) CreateSubBlock(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SubBlockID string `json:"sub_block_id"`
		TxCount    int    `json:"tx_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	subBlock, err := api.TransactionPool.CreateSubBlock(request.SubBlockID, request.TxCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create sub-block: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Sub-block created successfully",
		"sub_block": subBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddSubBlockToLedger adds a sub-block to the ledger
func (api *TransactionsAPI) AddSubBlockToLedger(w http.ResponseWriter, r *http.Request) {
	var subBlock common.SubBlock

	if err := json.NewDecoder(r.Body).Decode(&subBlock); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.TransactionPool.AddSubBlockToLedger(&subBlock)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add sub-block to ledger: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Sub-block added to ledger successfully",
		"sub_block_id": subBlock.SubBlockID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListPendingSubBlocks returns all pending sub-blocks
func (api *TransactionsAPI) ListPendingSubBlocks(w http.ResponseWriter, r *http.Request) {
	subBlocks := api.TransactionPool.ListPendingSubBlocks()

	response := map[string]interface{}{
		"success": true,
		"sub_blocks": subBlocks,
		"count": len(subBlocks),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTransactionStatus returns the status of a specific transaction
func (api *TransactionsAPI) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txID := vars["txID"]

	status, err := api.LedgerInstance.GetTransactionStatus(txID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction status: %v", err), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"transaction_id": txID,
		"status": status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTransactionHistory returns the transaction history for an address
func (api *TransactionsAPI) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	// Parse optional query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	history, err := api.LedgerInstance.GetTransactionHistory(address, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction history: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"address": address,
		"history": history,
		"count": len(history),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}