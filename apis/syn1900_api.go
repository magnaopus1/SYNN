package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/tokens/syn1900"
	"synnergy_network/pkg/common"
)

// SYN1900API handles all SYN1900 Loyalty Program Token related API endpoints
type SYN1900API struct {
	factory     *syn1900.TokenFactory
	storage     *syn1900.TokenStorageService
	security    *syn1900.SecurityService
	management  *syn1900.ManagementService
	transaction *syn1900.TokenTransactionService
	events      *syn1900.EventManager
	compliance  *syn1900.ComplianceManager
}

func NewSYN1900API() *SYN1900API {
	factory := &syn1900.TokenFactory{}
	storage := &syn1900.TokenStorageService{}
	security := &syn1900.SecurityService{}
	management := &syn1900.ManagementService{}
	transaction := &syn1900.TokenTransactionService{}
	events := &syn1900.EventManager{}
	compliance := &syn1900.ComplianceManager{}

	return &SYN1900API{
		factory:     factory,
		storage:     storage,
		security:    security,
		management:  management,
		transaction: transaction,
		events:      events,
		compliance:  compliance,
	}
}

func (api *SYN1900API) RegisterRoutes(router *mux.Router) {
	// Core token management endpoints
	router.HandleFunc("/syn1900/tokens", api.CreateToken).Methods("POST")
	router.HandleFunc("/syn1900/tokens/{tokenID}", api.GetToken).Methods("GET")
	router.HandleFunc("/syn1900/tokens", api.ListTokens).Methods("GET")
	router.HandleFunc("/syn1900/tokens/{tokenID}/revoke", api.RevokeToken).Methods("POST")
	router.HandleFunc("/syn1900/tokens/{tokenID}/transfer", api.TransferToken).Methods("POST")
	router.HandleFunc("/syn1900/tokens/{tokenID}/update", api.UpdateToken).Methods("PUT")
	router.HandleFunc("/syn1900/tokens/{tokenID}/delete", api.DeleteToken).Methods("DELETE")

	// Management endpoints
	router.HandleFunc("/syn1900/management/{tokenID}/credit/transfer", api.TransferCredit).Methods("POST")
	router.HandleFunc("/syn1900/management/{tokenID}/credit/revoke", api.RevokeCredit).Methods("POST")
	router.HandleFunc("/syn1900/management/{tokenID}/credit/validate", api.ValidateCredit).Methods("GET")
	router.HandleFunc("/syn1900/management/{tokenID}/metadata", api.UpdateMetadata).Methods("PUT")
	router.HandleFunc("/syn1900/management/{tokenID}/renew", api.RenewCredit).Methods("POST")
	router.HandleFunc("/syn1900/management/{tokenID}/summary", api.GenerateSummary).Methods("GET")

	// Security endpoints
	router.HandleFunc("/syn1900/security/{tokenID}/verify-signature", api.VerifyDigitalSignature).Methods("POST")
	router.HandleFunc("/syn1900/security/{tokenID}/encrypt", api.EncryptTokenMetadata).Methods("POST")
	router.HandleFunc("/syn1900/security/{tokenID}/decrypt", api.DecryptTokenMetadata).Methods("POST")
	router.HandleFunc("/syn1900/security/{tokenID}/hash", api.HashToken).Methods("GET")
	router.HandleFunc("/syn1900/security/{tokenID}/validate-integrity", api.ValidateTokenIntegrity).Methods("POST")
	router.HandleFunc("/syn1900/security/{tokenID}/revoke", api.RevokeTokenSecurity).Methods("POST")

	// Transaction endpoints
	router.HandleFunc("/syn1900/transactions", api.RecordTransaction).Methods("POST")
	router.HandleFunc("/syn1900/transactions/{transactionID}", api.GetTransaction).Methods("GET")
	router.HandleFunc("/syn1900/transactions/{transactionID}/update", api.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/syn1900/transactions", api.ListTransactions).Methods("GET")
	router.HandleFunc("/syn1900/transactions/{transactionID}/validate", api.ValidateTransaction).Methods("GET")
	router.HandleFunc("/syn1900/transactions/{transactionID}/revoke", api.RevokeTransaction).Methods("POST")

	// Event management endpoints
	router.HandleFunc("/syn1900/events/{tokenID}/log", api.LogEvent).Methods("POST")
	router.HandleFunc("/syn1900/events/{tokenID}", api.FetchTokenEvents).Methods("GET")
	router.HandleFunc("/syn1900/events/{tokenID}/report", api.GenerateEventReport).Methods("GET")
	router.HandleFunc("/syn1900/events/{tokenID}/completion", api.AddEventForCompletion).Methods("POST")
	router.HandleFunc("/syn1900/events/{tokenID}/transfer", api.AddEventForTransfer).Methods("POST")
	router.HandleFunc("/syn1900/events/{tokenID}/revocation", api.AddEventForRevocation).Methods("POST")

	// Compliance endpoints
	router.HandleFunc("/syn1900/compliance/{tokenID}/verify", api.VerifyTokenCompliance).Methods("GET")
	router.HandleFunc("/syn1900/compliance/{tokenID}/revoke", api.RevokeNonCompliantToken).Methods("POST")
	router.HandleFunc("/syn1900/compliance/audit", api.AuditTokenCompliance).Methods("GET")
	router.HandleFunc("/syn1900/compliance/report", api.GenerateComplianceReport).Methods("GET")

	// Utility endpoints
	router.HandleFunc("/syn1900/health", api.HealthCheck).Methods("GET")
	router.HandleFunc("/syn1900/metrics", api.GetMetrics).Methods("GET")
}

// Token Management Handlers

type CreateTokenRequest struct {
	IssuerID        string     `json:"issuer_id"`
	RecipientID     string     `json:"recipient_id"`
	CourseID        string     `json:"course_id"`
	CourseName      string     `json:"course_name"`
	CreditValue     float64    `json:"credit_value"`
	Metadata        string     `json:"metadata"`
	Signature       []byte     `json:"signature"`
	ExpirationDate  *time.Time `json:"expiration_date,omitempty"`
}

func (api *SYN1900API) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	token, err := api.factory.CreateSYN1900Token(
		req.IssuerID,
		req.RecipientID,
		req.CourseID,
		req.CourseName,
		req.CreditValue,
		req.Metadata,
		req.Signature,
		req.ExpirationDate,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create token: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the token
	err = api.storage.SaveToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to store token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  token,
	})
}

func (api *SYN1900API) GetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	decryptionKey := r.URL.Query().Get("decryption_key")
	if decryptionKey == "" {
		http.Error(w, "Decryption key required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	token, err := api.storage.GetToken(tokenID, decryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve token: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  token,
	})
}

func (api *SYN1900API) ListTokens(w http.ResponseWriter, r *http.Request) {
	// Call the real module function
	tokens, err := api.storage.ListAllTokens()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list tokens: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"tokens": tokens,
	})
}

type RevokeTokenRequest struct {
	Reason string `json:"reason"`
}

func (api *SYN1900API) RevokeToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req RevokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.factory.RevokeSYN1900Token(tokenID, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token revoked successfully",
	})
}

type TransferTokenRequest struct {
	NewRecipientID string `json:"new_recipient_id"`
	TransferType   string `json:"transfer_type"`
}

func (api *SYN1900API) TransferToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req TransferTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.factory.TransferSYN1900Token(tokenID, req.NewRecipientID, req.TransferType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transfer token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token transferred successfully",
	})
}

func (api *SYN1900API) UpdateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// First get the token
	decryptionKey := r.URL.Query().Get("decryption_key")
	if decryptionKey == "" {
		http.Error(w, "Decryption key required", http.StatusBadRequest)
		return
	}

	token, err := api.storage.GetToken(tokenID, decryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve token: %v", err), http.StatusNotFound)
		return
	}

	// Call the real module function
	err = api.storage.UpdateToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token updated successfully",
	})
}

func (api *SYN1900API) DeleteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	err := api.storage.DeleteToken(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token deleted successfully",
	})
}

// Management Handlers

type TransferCreditRequest struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
}

func (api *SYN1900API) TransferCredit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req TransferCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.management.TransferCredit(tokenID, req.FromID, req.ToID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transfer credit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Credit transferred successfully",
	})
}

type RevokeCreditRequest struct {
	RevocationReason string `json:"revocation_reason"`
}

func (api *SYN1900API) RevokeCredit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req RevokeCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.management.RevokeCredit(tokenID, req.RevocationReason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke credit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Credit revoked successfully",
	})
}

func (api *SYN1900API) ValidateCredit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	isValid, err := api.management.ValidateCredit(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate credit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"is_valid": isValid,
	})
}

type UpdateMetadataRequest struct {
	NewMetadata string `json:"new_metadata"`
}

func (api *SYN1900API) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req UpdateMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.management.UpdateMetadata(tokenID, req.NewMetadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Metadata updated successfully",
	})
}

type RenewCreditRequest struct {
	NewExpirationDate time.Time `json:"new_expiration_date"`
}

func (api *SYN1900API) RenewCredit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req RenewCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.management.RenewCredit(tokenID, req.NewExpirationDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to renew credit: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Credit renewed successfully",
	})
}

func (api *SYN1900API) GenerateSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	summary, err := api.management.GenerateSummary(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate summary: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"summary": summary,
	})
}

// Security Handlers

func (api *SYN1900API) VerifyDigitalSignature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	isValid, err := api.security.VerifyDigitalSignature(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify signature: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "success",
		"signature_valid":    isValid,
	})
}

type EncryptMetadataRequest struct {
	EncryptionKey string `json:"encryption_key"`
}

func (api *SYN1900API) EncryptTokenMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req EncryptMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.security.EncryptTokenMetadata(tokenID, req.EncryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encrypt metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Metadata encrypted successfully",
	})
}

type DecryptMetadataRequest struct {
	DecryptionKey string `json:"decryption_key"`
}

func (api *SYN1900API) DecryptTokenMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req DecryptMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	metadata, err := api.security.DecryptTokenMetadata(tokenID, req.DecryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decrypt metadata: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"metadata": metadata,
	})
}

func (api *SYN1900API) HashToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	hash, err := api.security.HashToken(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to hash token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"hash":   hash,
	})
}

type ValidateIntegrityRequest struct {
	StoredHash string `json:"stored_hash"`
}

func (api *SYN1900API) ValidateTokenIntegrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req ValidateIntegrityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	isValid, err := api.security.ValidateTokenIntegrity(tokenID, req.StoredHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate integrity: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "success",
		"integrity_valid": isValid,
	})
}

type RevokeTokenSecurityRequest struct {
	Reason string `json:"reason"`
}

func (api *SYN1900API) RevokeTokenSecurity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req RevokeTokenSecurityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.security.RevokeToken(tokenID, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke token security: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token security revoked successfully",
	})
}

// Transaction Handlers

func (api *SYN1900API) RecordTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction common.SYN1900Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	encryptionKey := r.URL.Query().Get("encryption_key")
	if encryptionKey == "" {
		http.Error(w, "Encryption key required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.transaction.RecordTransaction(transaction, encryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to record transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transaction recorded successfully",
	})
}

func (api *SYN1900API) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transactionID"]

	decryptionKey := r.URL.Query().Get("decryption_key")
	if decryptionKey == "" {
		http.Error(w, "Decryption key required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	transaction, err := api.transaction.GetTransaction(transactionID, decryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"transaction": transaction,
	})
}

func (api *SYN1900API) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction common.SYN1900Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	encryptionKey := r.URL.Query().Get("encryption_key")
	if encryptionKey == "" {
		http.Error(w, "Encryption key required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.transaction.UpdateTransaction(transaction, encryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transaction updated successfully",
	})
}

func (api *SYN1900API) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// Call the real module function
	transactions, err := api.transaction.ListAllTransactions()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "success",
		"transactions": transactions,
	})
}

func (api *SYN1900API) ValidateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transactionID"]

	// First get the transaction
	decryptionKey := r.URL.Query().Get("decryption_key")
	if decryptionKey == "" {
		http.Error(w, "Decryption key required", http.StatusBadRequest)
		return
	}

	transaction, err := api.transaction.GetTransaction(transactionID, decryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction: %v", err), http.StatusNotFound)
		return
	}

	// Call the real module function
	err = api.transaction.ValidateTransaction(transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction validation failed: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transaction is valid",
	})
}

type RevokeTransactionRequest struct {
	Reason string `json:"reason"`
}

func (api *SYN1900API) RevokeTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transactionID"]

	var req RevokeTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.transaction.RevokeTransaction(transactionID, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transaction revoked successfully",
	})
}

// Event Handlers

type LogEventRequest struct {
	EventType   string `json:"event_type"`
	Description string `json:"description"`
}

func (api *SYN1900API) LogEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req LogEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.events.LogEvent(tokenID, req.EventType, req.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to log event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Event logged successfully",
	})
}

func (api *SYN1900API) FetchTokenEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	events, err := api.events.FetchTokenEvents(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch events: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"events": events,
	})
}

func (api *SYN1900API) GenerateEventReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	report, err := api.events.GenerateEventReport(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate event report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"report": report,
	})
}

type AddEventForCompletionRequest struct {
	RecipientID string `json:"recipient_id"`
}

func (api *SYN1900API) AddEventForCompletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req AddEventForCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.events.AddEventForCompletion(tokenID, req.RecipientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add completion event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Completion event added successfully",
	})
}

type AddEventForTransferRequest struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
}

func (api *SYN1900API) AddEventForTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req AddEventForTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.events.AddEventForTransfer(tokenID, req.FromID, req.ToID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add transfer event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Transfer event added successfully",
	})
}

type AddEventForRevocationRequest struct {
	Reason string `json:"reason"`
}

func (api *SYN1900API) AddEventForRevocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req AddEventForRevocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.events.AddEventForRevocation(tokenID, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add revocation event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Revocation event added successfully",
	})
}

// Compliance Handlers

func (api *SYN1900API) VerifyTokenCompliance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	// Call the real module function
	err := api.compliance.VerifyTokenCompliance(tokenID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Compliance verification failed: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token is compliant",
	})
}

type RevokeNonCompliantRequest struct {
	Reason string `json:"reason"`
}

func (api *SYN1900API) RevokeNonCompliantToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenID := vars["tokenID"]

	var req RevokeNonCompliantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the real module function
	err := api.compliance.RevokeNonCompliantToken(tokenID, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke non-compliant token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Non-compliant token revoked successfully",
	})
}

func (api *SYN1900API) AuditTokenCompliance(w http.ResponseWriter, r *http.Request) {
	issuerID := r.URL.Query().Get("issuer_id")
	recipientID := r.URL.Query().Get("recipient_id")

	if issuerID == "" || recipientID == "" {
		http.Error(w, "Issuer ID and Recipient ID are required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	auditLogs, err := api.compliance.AuditTokenCompliance(issuerID, recipientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to audit compliance: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"audit_logs": auditLogs,
	})
}

func (api *SYN1900API) GenerateComplianceReport(w http.ResponseWriter, r *http.Request) {
	issuerID := r.URL.Query().Get("issuer_id")
	recipientID := r.URL.Query().Get("recipient_id")

	if issuerID == "" || recipientID == "" {
		http.Error(w, "Issuer ID and Recipient ID are required", http.StatusBadRequest)
		return
	}

	// Call the real module function
	report, err := api.compliance.GenerateComplianceReport(issuerID, recipientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate compliance report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"report": report,
	})
}

// Utility Handlers

func (api *SYN1900API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "syn1900-api",
		"timestamp": time.Now(),
	})
}

func (api *SYN1900API) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"metrics": "Metrics collection would be implemented here",
	})
}