package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/tokens/syn2600"
	
	"github.com/gorilla/mux"
)

// SYN2600API provides API endpoints for SYN2600 Legal Document Token operations
type SYN2600API struct {
	ledgerInstance *ledger.Ledger
	consensus      *common.SynnergyConsensus
	mutex          *common.SynnergyMutex
}

// NewSYN2600API creates a new SYN2600API instance
func NewSYN2600API(ledgerInstance *ledger.Ledger, consensus *common.SynnergyConsensus, mutex *common.SynnergyMutex) *SYN2600API {
	return &SYN2600API{
		ledgerInstance: ledgerInstance,
		consensus:      consensus,
		mutex:          mutex,
	}
}

// RegisterRoutes registers all SYN2600 token endpoints
func (api *SYN2600API) RegisterRoutes(router *mux.Router) {
	// Token Factory & Creation (15 endpoints)
	router.HandleFunc("/syn2600/factory/create", api.CreateToken).Methods("POST")
	router.HandleFunc("/syn2600/factory/batch", api.CreateBatchTokens).Methods("POST")
	router.HandleFunc("/syn2600/factory/template", api.CreateFromTemplate).Methods("POST")
	router.HandleFunc("/syn2600/factory/validate", api.ValidateTokenCreation).Methods("POST")
	router.HandleFunc("/syn2600/factory/estimate", api.EstimateCreationCost).Methods("POST")
	router.HandleFunc("/syn2600/factory/config", api.GetFactoryConfig).Methods("GET")
	router.HandleFunc("/syn2600/factory/config", api.UpdateFactoryConfig).Methods("PUT")
	router.HandleFunc("/syn2600/factory/status", api.GetFactoryStatus).Methods("GET")
	router.HandleFunc("/syn2600/factory/templates", api.ListTemplates).Methods("GET")
	router.HandleFunc("/syn2600/factory/templates/{templateId}", api.GetTemplate).Methods("GET")
	router.HandleFunc("/syn2600/factory/templates", api.CreateTemplate).Methods("POST")
	router.HandleFunc("/syn2600/factory/templates/{templateId}", api.UpdateTemplate).Methods("PUT")
	router.HandleFunc("/syn2600/factory/templates/{templateId}", api.DeleteTemplate).Methods("DELETE")
	router.HandleFunc("/syn2600/factory/metrics", api.GetFactoryMetrics).Methods("GET")
	router.HandleFunc("/syn2600/factory/history", api.GetCreationHistory).Methods("GET")

	// Token Management (20 endpoints)
	router.HandleFunc("/syn2600/management/tokens/{tokenId}", api.GetToken).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens", api.ListTokens).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}", api.UpdateToken).Methods("PUT")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/activate", api.ActivateToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/deactivate", api.DeactivateToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/suspend", api.SuspendToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/resume", api.ResumeToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/freeze", api.FreezeToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/unfreeze", api.UnfreezeToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/metadata", api.UpdateTokenMetadata).Methods("PUT")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/metadata", api.GetTokenMetadata).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/owner", api.TransferOwnership).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/history", api.GetTokenHistory).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/status", api.GetTokenStatus).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens/{tokenId}/validate", api.ValidateToken).Methods("POST")
	router.HandleFunc("/syn2600/management/tokens/search", api.SearchTokens).Methods("GET")
	router.HandleFunc("/syn2600/management/tokens/filter", api.FilterTokens).Methods("POST")
	router.HandleFunc("/syn2600/management/batch/update", api.BatchUpdateTokens).Methods("PUT")
	router.HandleFunc("/syn2600/management/batch/action", api.BatchTokenAction).Methods("POST")
	router.HandleFunc("/syn2600/management/stats", api.GetManagementStats).Methods("GET")

	// Token Storage & Retrieval (10 endpoints)
	router.HandleFunc("/syn2600/storage/store", api.StoreToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/retrieve/{tokenId}", api.RetrieveToken).Methods("GET")
	router.HandleFunc("/syn2600/storage/exists/{tokenId}", api.CheckTokenExists).Methods("GET")
	router.HandleFunc("/syn2600/storage/backup/{tokenId}", api.BackupToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/restore/{tokenId}", api.RestoreToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/export/{tokenId}", api.ExportToken).Methods("GET")
	router.HandleFunc("/syn2600/storage/import", api.ImportToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/archive/{tokenId}", api.ArchiveToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/unarchive/{tokenId}", api.UnarchiveToken).Methods("POST")
	router.HandleFunc("/syn2600/storage/stats", api.GetStorageStats).Methods("GET")

	// Security & Access Control (15 endpoints)
	router.HandleFunc("/syn2600/security/encrypt/{tokenId}", api.EncryptToken).Methods("POST")
	router.HandleFunc("/syn2600/security/decrypt/{tokenId}", api.DecryptToken).Methods("POST")
	router.HandleFunc("/syn2600/security/verify/{tokenId}", api.VerifyTokenSecurity).Methods("POST")
	router.HandleFunc("/syn2600/security/permissions/{tokenId}", api.SetPermissions).Methods("POST")
	router.HandleFunc("/syn2600/security/permissions/{tokenId}", api.GetPermissions).Methods("GET")
	router.HandleFunc("/syn2600/security/access/{tokenId}/grant", api.GrantAccess).Methods("POST")
	router.HandleFunc("/syn2600/security/access/{tokenId}/revoke", api.RevokeAccess).Methods("POST")
	router.HandleFunc("/syn2600/security/access/{tokenId}/check", api.CheckAccess).Methods("GET")
	router.HandleFunc("/syn2600/security/keys/{tokenId}/rotate", api.RotateKeys).Methods("POST")
	router.HandleFunc("/syn2600/security/keys/{tokenId}", api.GetKeyInfo).Methods("GET")
	router.HandleFunc("/syn2600/security/audit/{tokenId}", api.SecurityAudit).Methods("POST")
	router.HandleFunc("/syn2600/security/violations/{tokenId}", api.GetSecurityViolations).Methods("GET")
	router.HandleFunc("/syn2600/security/policies", api.GetSecurityPolicies).Methods("GET")
	router.HandleFunc("/syn2600/security/policies", api.UpdateSecurityPolicies).Methods("PUT")
	router.HandleFunc("/syn2600/security/status", api.GetSecurityStatus).Methods("GET")

	// Transaction Operations (15 endpoints)
	router.HandleFunc("/syn2600/transactions/transfer", api.TransferToken).Methods("POST")
	router.HandleFunc("/syn2600/transactions/{txId}", api.GetTransaction).Methods("GET")
	router.HandleFunc("/syn2600/transactions", api.ListTransactions).Methods("GET")
	router.HandleFunc("/syn2600/transactions/{txId}/status", api.GetTransactionStatus).Methods("GET")
	router.HandleFunc("/syn2600/transactions/{txId}/confirm", api.ConfirmTransaction).Methods("POST")
	router.HandleFunc("/syn2600/transactions/{txId}/cancel", api.CancelTransaction).Methods("POST")
	router.HandleFunc("/syn2600/transactions/validate", api.ValidateTransaction).Methods("POST")
	router.HandleFunc("/syn2600/transactions/estimate", api.EstimateTransactionFee).Methods("POST")
	router.HandleFunc("/syn2600/transactions/batch", api.BatchTransactions).Methods("POST")
	router.HandleFunc("/syn2600/transactions/pending", api.GetPendingTransactions).Methods("GET")
	router.HandleFunc("/syn2600/transactions/history/{tokenId}", api.GetTokenTransactionHistory).Methods("GET")
	router.HandleFunc("/syn2600/transactions/analytics", api.GetTransactionAnalytics).Methods("GET")
	router.HandleFunc("/syn2600/transactions/{txId}/receipt", api.GetTransactionReceipt).Methods("GET")
	router.HandleFunc("/syn2600/transactions/search", api.SearchTransactions).Methods("GET")
	router.HandleFunc("/syn2600/transactions/stats", api.GetTransactionStats).Methods("GET")

	// Events & Notifications (10 endpoints)
	router.HandleFunc("/syn2600/events/{tokenId}", api.GetTokenEvents).Methods("GET")
	router.HandleFunc("/syn2600/events/subscribe", api.SubscribeToEvents).Methods("POST")
	router.HandleFunc("/syn2600/events/unsubscribe", api.UnsubscribeFromEvents).Methods("POST")
	router.HandleFunc("/syn2600/events/emit", api.EmitCustomEvent).Methods("POST")
	router.HandleFunc("/syn2600/events/history", api.GetEventHistory).Methods("GET")
	router.HandleFunc("/syn2600/events/types", api.GetEventTypes).Methods("GET")
	router.HandleFunc("/syn2600/events/filters", api.SetEventFilters).Methods("POST")
	router.HandleFunc("/syn2600/events/filters", api.GetEventFilters).Methods("GET")
	router.HandleFunc("/syn2600/events/notifications", api.GetNotifications).Methods("GET")
	router.HandleFunc("/syn2600/events/stats", api.GetEventStats).Methods("GET")

	// Compliance & Validation (15 endpoints)
	router.HandleFunc("/syn2600/compliance/check/{tokenId}", api.CheckCompliance).Methods("POST")
	router.HandleFunc("/syn2600/compliance/validate/{tokenId}", api.ValidateCompliance).Methods("POST")
	router.HandleFunc("/syn2600/compliance/report/{tokenId}", api.GetComplianceReport).Methods("GET")
	router.HandleFunc("/syn2600/compliance/rules", api.GetComplianceRules).Methods("GET")
	router.HandleFunc("/syn2600/compliance/rules", api.UpdateComplianceRules).Methods("PUT")
	router.HandleFunc("/syn2600/compliance/violations/{tokenId}", api.GetComplianceViolations).Methods("GET")
	router.HandleFunc("/syn2600/compliance/audit/{tokenId}", api.ComplianceAudit).Methods("POST")
	router.HandleFunc("/syn2600/compliance/certify/{tokenId}", api.CertifyCompliance).Methods("POST")
	router.HandleFunc("/syn2600/compliance/remediation/{tokenId}", api.InitiateRemediation).Methods("POST")
	router.HandleFunc("/syn2600/compliance/status/{tokenId}", api.GetComplianceStatus).Methods("GET")
	router.HandleFunc("/syn2600/compliance/history/{tokenId}", api.GetComplianceHistory).Methods("GET")
	router.HandleFunc("/syn2600/compliance/frameworks", api.GetSupportedFrameworks).Methods("GET")
	router.HandleFunc("/syn2600/compliance/export/{tokenId}", api.ExportComplianceData).Methods("GET")
	router.HandleFunc("/syn2600/compliance/alerts", api.GetComplianceAlerts).Methods("GET")
	router.HandleFunc("/syn2600/compliance/stats", api.GetComplianceStats).Methods("GET")
}

// Token Factory & Creation Endpoints

func (api *SYN2600API) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AssetDetails string    `json:"asset_details"`
		Owner        string    `json:"owner"`
		Shares       float64   `json:"shares"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call real module function
	tokenID, err := syn2600.CreateNewInvestorToken(req.AssetDetails, req.Owner, req.Shares, req.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"tokenId":  tokenID,
		"message":  "SYN2600 token created successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) CreateBatchTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Tokens []struct {
			AssetDetails string    `json:"asset_details"`
			Owner        string    `json:"owner"`
			Shares       float64   `json:"shares"`
			ExpiryDate   time.Time `json:"expiry_date"`
		} `json:"tokens"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdTokens := make([]string, 0)
	for _, token := range req.Tokens {
		tokenID, err := syn2600.CreateNewInvestorToken(token.AssetDetails, token.Owner, token.Shares, token.ExpiryDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create token %d: %v", len(createdTokens)+1, err), http.StatusInternalServerError)
			return
		}
		createdTokens = append(createdTokens, tokenID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_ids": createdTokens,
		"message": "SYN2600 tokens created successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TemplateID string `json:"template_id" validate:"required"`
		Owner      string `json:"owner" validate:"required"`
		ExpiryDate time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenID, err := syn2600.CreateNewInvestorTokenFromTemplate(req.TemplateID, req.Owner, req.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create token from template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenID,
		"message": "SYN2600 token created from template successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ValidateTokenCreation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := syn2600.ValidateTokenCreation(req.AssetDetails, req.Owner, req.Shares, req.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate token creation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": valid,
		"message": "Token creation validation result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) EstimateCreationCost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cost, err := syn2600.EstimateTokenCreationCost(req.AssetDetails, req.Owner, req.Shares, req.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to estimate token creation cost: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"cost": cost,
		"message": "Token creation cost estimated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetFactoryConfig(w http.ResponseWriter, r *http.Request) {
	config, err := syn2600.GetFactoryConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get factory config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (api *SYN2600API) UpdateFactoryConfig(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateFactoryConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update factory config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Factory config updated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetFactoryStatus(w http.ResponseWriter, r *http.Request) {
	status, err := syn2600.GetFactoryStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get factory status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api *SYN2600API) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := syn2600.ListTemplates()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list templates: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func (api *SYN2600API) GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	template, err := syn2600.GetTemplate(templateId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get template %s: %v", templateId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

func (api *SYN2600API) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var template syn2600.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.CreateTemplate(template)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create template: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"template_id": template.ID,
		"message": "Template created successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	var template syn2600.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateTemplate(templateId, template)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update template %s: %v", templateId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"template_id": templateId,
		"message": "Template updated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	err := syn2600.DeleteTemplate(templateId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete template %s: %v", templateId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"template_id": templateId,
		"message": "Template deleted successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetFactoryMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := syn2600.GetFactoryMetrics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get factory metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (api *SYN2600API) GetCreationHistory(w http.ResponseWriter, r *http.Request) {
	history, err := syn2600.GetCreationHistory()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get creation history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// Token Management Endpoints

func (api *SYN2600API) GetToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	token, err := syn2600.GetToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (api *SYN2600API) ListTokens(w http.ResponseWriter, r *http.Request) {
	tokens, err := syn2600.ListTokens()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list tokens: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (api *SYN2600API) UpdateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var token syn2600.Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateToken(tokenId, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token updated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ActivateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.ActivateToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to activate token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token activated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) DeactivateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.DeactivateToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to deactivate token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token deactivated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) SuspendToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.SuspendToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to suspend token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token suspended successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ResumeToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.ResumeToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to resume token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token resumed successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) FreezeToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.FreezeToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to freeze token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token frozen successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) UnfreezeToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.UnfreezeToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unfreeze token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token un-frozen successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) UpdateTokenMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var metadata map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&metadata); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateTokenMetadata(tokenId, metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update token %s metadata: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token metadata updated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetTokenMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	metadata, err := syn2600.GetTokenMetadata(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get token %s metadata: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func (api *SYN2600API) TransferOwnership(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var newOwner string
	if err := json.NewDecoder(r.Body).Decode(&newOwner); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.TransferOwnership(tokenId, newOwner)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transfer ownership of token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token ownership transferred successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetTokenHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	history, err := syn2600.GetTokenHistory(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get history for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (api *SYN2600API) GetTokenStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	status, err := syn2600.GetTokenStatus(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get status for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api *SYN2600API) ValidateToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var validationReq struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&validationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := syn2600.ValidateToken(tokenId, validationReq.AssetDetails, validationReq.Owner, validationReq.Shares, validationReq.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": valid,
		"token_id": tokenId,
		"message": "Token validation result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) SearchTokens(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	tokens, err := syn2600.SearchTokens(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search tokens: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (api *SYN2600API) FilterTokens(w http.ResponseWriter, r *http.Request) {
	var filter syn2600.TokenFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := syn2600.FilterTokens(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to filter tokens: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (api *SYN2600API) BatchUpdateTokens(w http.ResponseWriter, r *http.Request) {
	var batchReq struct {
		Tokens []struct {
			ID           string                 `json:"id"`
			Updates      map[string]interface{} `json:"updates"`
			NewOwner     string                 `json:"new_owner"`
			NewExpiry    time.Time              `json:"new_expiry"`
		} `json:"tokens"`
	}

	if err := json.NewDecoder(r.Body).Decode(&batchReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := make([]map[string]interface{}, 0)
	for _, item := range batchReq.Tokens {
		result, err := syn2600.BatchUpdateToken(item.ID, item.Updates, item.NewOwner, item.NewExpiry)
		if err != nil {
			results = append(results, map[string]interface{}{"id": item.ID, "success": false, "error": err.Error()})
		} else {
			results = append(results, map[string]interface{}{"id": item.ID, "success": true, "result": result})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (api *SYN2600API) BatchTokenAction(w http.ResponseWriter, r *http.Request) {
	var batchReq struct {
		Actions []struct {
			TokenID string `json:"token_id" validate:"required"`
			Action  string `json:"action" validate:"required"` // e.g., "activate", "deactivate", "suspend", "resume", "freeze", "unfreeze"
		} `json:"actions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&batchReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results := make([]map[string]interface{}, 0)
	for _, action := range batchReq.Actions {
		var err error
		switch action.Action {
		case "activate":
			err = syn2600.ActivateToken(action.TokenID)
		case "deactivate":
			err = syn2600.DeactivateToken(action.TokenID)
		case "suspend":
			err = syn2600.SuspendToken(action.TokenID)
		case "resume":
			err = syn2600.ResumeToken(action.TokenID)
		case "freeze":
			err = syn2600.FreezeToken(action.TokenID)
		case "unfreeze":
			err = syn2600.UnfreezeToken(action.TokenID)
		default:
			results = append(results, map[string]interface{}{"token_id": action.TokenID, "success": false, "error": "Unknown action"})
			continue
		}

		if err != nil {
			results = append(results, map[string]interface{}{"token_id": action.TokenID, "success": false, "error": err.Error()})
		} else {
			results = append(results, map[string]interface{}{"token_id": action.TokenID, "success": true, "message": fmt.Sprintf("Token %s action successful", action.Action)})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (api *SYN2600API) GetManagementStats(w http.ResponseWriter, r *http.Request) {
	stats, err := syn2600.GetManagementStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get management stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Token Storage & Retrieval Endpoints

func (api *SYN2600API) StoreToken(w http.ResponseWriter, r *http.Request) {
	var token syn2600.Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenId, err := syn2600.StoreToken(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to store token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token stored successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) RetrieveToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	token, err := syn2600.RetrieveToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (api *SYN2600API) CheckTokenExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	exists, err := syn2600.CheckTokenExists(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check token %s existence: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"exists": exists,
		"message": "Token existence check result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) BackupToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.BackupToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to backup token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token backed up successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) RestoreToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.RestoreToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to restore token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token restored successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ExportToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	exportData, err := syn2600.ExportToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exportData)
}

func (api *SYN2600API) ImportToken(w http.ResponseWriter, r *http.Request) {
	var importData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&importData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenId, err := syn2600.ImportToken(importData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to import token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token imported successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ArchiveToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.ArchiveToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to archive token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token archived successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) UnarchiveToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.UnarchiveToken(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unarchive token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Token un-archived successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetStorageStats(w http.ResponseWriter, r *http.Request) {
	stats, err := syn2600.GetStorageStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get storage stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Security & Access Control Endpoints

func (api *SYN2600API) EncryptToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var encryptionReq struct {
		Plaintext string `json:"plaintext" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&encryptionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ciphertext, err := syn2600.EncryptToken(tokenId, encryptionReq.Plaintext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encrypt token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"ciphertext": ciphertext,
		"message": "Token encrypted successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) DecryptToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var decryptionReq struct {
		Ciphertext string `json:"ciphertext" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&decryptionReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	plaintext, err := syn2600.DecryptToken(tokenId, decryptionReq.Ciphertext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decrypt token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"plaintext": plaintext,
		"message": "Token decrypted successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) VerifyTokenSecurity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var verificationReq struct {
		Plaintext string `json:"plaintext" validate:"required"`
		Ciphertext string `json:"ciphertext" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&verificationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	verified, err := syn2600.VerifyTokenSecurity(tokenId, verificationReq.Plaintext, verificationReq.Ciphertext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify token %s security: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": verified,
		"token_id": tokenId,
		"message": "Token security verification result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) SetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var permissions map[string][]string
	if err := json.NewDecoder(r.Body).Decode(&permissions); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.SetPermissions(tokenId, permissions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set permissions for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Permissions set successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	permissions, err := syn2600.GetPermissions(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get permissions for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (api *SYN2600API) GrantAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var accessReq struct {
		Address string `json:"address" validate:"required"`
		Actions []string `json:"actions" validate:"required,min=1"`
	}

	if err := json.NewDecoder(r.Body).Decode(&accessReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.GrantAccess(tokenId, accessReq.Address, accessReq.Actions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to grant access for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Access granted successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var revokeReq struct {
		Address string `json:"address" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&revokeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.RevokeAccess(tokenId, revokeReq.Address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to revoke access for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Access revoked successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) CheckAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var checkReq struct {
		Address string `json:"address" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hasAccess, err := syn2600.CheckAccess(tokenId, checkReq.Address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check access for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"address": checkReq.Address,
		"has_access": hasAccess,
		"message": "Access check result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) RotateKeys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	err := syn2600.RotateKeys(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to rotate keys for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Keys rotated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetKeyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	keyInfo, err := syn2600.GetKeyInfo(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get key info for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyInfo)
}

func (api *SYN2600API) SecurityAudit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var auditReq struct {
		Action string `json:"action" validate:"required"`
		Details string `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&auditReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.SecurityAudit(tokenId, auditReq.Action, auditReq.Details)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to audit security for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Security audit initiated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetSecurityViolations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	violations, err := syn2600.GetSecurityViolations(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get security violations for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violations)
}

func (api *SYN2600API) GetSecurityPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := syn2600.GetSecurityPolicies()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get security policies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

func (api *SYN2600API) UpdateSecurityPolicies(w http.ResponseWriter, r *http.Request) {
	var policies map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&policies); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateSecurityPolicies(policies)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update security policies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Security policies updated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetSecurityStatus(w http.ResponseWriter, r *http.Request) {
	status, err := syn2600.GetSecurityStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get security status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Transaction Operations Endpoints

func (api *SYN2600API) TransferToken(w http.ResponseWriter, r *http.Request) {
	var transferReq struct {
		FromTokenID string `json:"from_token_id" validate:"required"`
		ToTokenID   string `json:"to_token_id" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,min=0.000000000000000001"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txId, err := syn2600.TransferToken(transferReq.FromTokenID, transferReq.ToTokenID, transferReq.Amount, transferReq.Metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to transfer token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tx_id": txId,
		"message": "Token transferred successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	tx, err := syn2600.GetTransaction(txId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction %s: %v", txId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (api *SYN2600API) ListTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	var filter syn2600.TransactionFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txs, err := syn2600.ListTransactions(filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func (api *SYN2600API) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	status, err := syn2600.GetTransactionStatus(txId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get status for transaction %s: %v", txId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api *SYN2600API) ConfirmTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	err := syn2600.ConfirmTransaction(txId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to confirm transaction %s: %v", txId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tx_id": txId,
		"message": "Transaction confirmed successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	err := syn2600.CancelTransaction(txId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel transaction %s: %v", txId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tx_id": txId,
		"message": "Transaction cancelled successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ValidateTransaction(w http.ResponseWriter, r *http.Request) {
	var validateReq struct {
		FromTokenID string `json:"from_token_id" validate:"required"`
		ToTokenID   string `json:"to_token_id" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,min=0.000000000000000001"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&validateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := syn2600.ValidateTransaction(validateReq.FromTokenID, validateReq.ToTokenID, validateReq.Amount, validateReq.Metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": valid,
		"message": "Transaction validation result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) EstimateTransactionFee(w http.ResponseWriter, r *http.Request) {
	var estimateReq struct {
		FromTokenID string `json:"from_token_id" validate:"required"`
		ToTokenID   string `json:"to_token_id" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,min=0.000000000000000001"`
		Metadata    map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&estimateReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fee, err := syn2600.EstimateTransactionFee(estimateReq.FromTokenID, estimateReq.ToTokenID, estimateReq.Amount, estimateReq.Metadata)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to estimate transaction fee: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"fee": fee,
		"message": "Transaction fee estimated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) BatchTransactions(w http.ResponseWriter, r *http.Request) {
	var batchReq struct {
		Transactions []struct {
			FromTokenID string `json:"from_token_id" validate:"required"`
			ToTokenID   string `json:"to_token_id" validate:"required"`
			Amount      float64 `json:"amount" validate:"required,min=0.000000000000000001"`
			Metadata    map[string]interface{} `json:"metadata"`
		} `json:"transactions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&batchReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txIds, err := syn2600.BatchTransactions(batchReq.Transactions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to batch transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tx_ids": txIds,
		"message": "Transactions batched successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	pending, err := syn2600.GetPendingTransactions()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get pending transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pending)
}

func (api *SYN2600API) GetTokenTransactionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	history, err := syn2600.GetTokenTransactionHistory(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction history for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (api *SYN2600API) GetTransactionAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := syn2600.GetTransactionAnalytics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction analytics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (api *SYN2600API) GetTransactionReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	receipt, err := syn2600.GetTransactionReceipt(txId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get receipt for transaction %s: %v", txId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func (api *SYN2600API) SearchTransactions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	
	var filter syn2600.TransactionFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	txs, err := syn2600.SearchTransactions(query, filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to search transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func (api *SYN2600API) GetTransactionStats(w http.ResponseWriter, r *http.Request) {
	stats, err := syn2600.GetTransactionStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Events & Notifications Endpoints

func (api *SYN2600API) GetTokenEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	events, err := syn2600.GetTokenEvents(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get events for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (api *SYN2600API) SubscribeToEvents(w http.ResponseWriter, r *http.Request) {
	var subscribeReq struct {
		EventTypes []string `json:"event_types" validate:"required,min=1"`
		CallbackURL string `json:"callback_url" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&subscribeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.SubscribeToEvents(subscribeReq.EventTypes, subscribeReq.CallbackURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to subscribe to events: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Subscribed to events successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) UnsubscribeFromEvents(w http.ResponseWriter, r *http.Request) {
	var unsubscribeReq struct {
		EventTypes []string `json:"event_types" validate:"required,min=1"`
	}

	if err := json.NewDecoder(r.Body).Decode(&unsubscribeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UnsubscribeFromEvents(unsubscribeReq.EventTypes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unsubscribe from events: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Unsubscribed from events successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) EmitCustomEvent(w http.ResponseWriter, r *http.Request) {
	var eventReq struct {
		TokenID string `json:"token_id" validate:"required"`
		EventType string `json:"event_type" validate:"required"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&eventReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.EmitCustomEvent(eventReq.TokenID, eventReq.EventType, eventReq.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to emit custom event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": eventReq.TokenID,
		"event_type": eventReq.EventType,
		"message": "Custom event emitted successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetEventHistory(w http.ResponseWriter, r *http.Request) {
	history, err := syn2600.GetEventHistory()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get event history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (api *SYN2600API) GetEventTypes(w http.ResponseWriter, r *http.Request) {
	types, err := syn2600.GetEventTypes()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get event types: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

func (api *SYN2600API) SetEventFilters(w http.ResponseWriter, r *http.Request) {
	var filters syn2600.EventFilter
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.SetEventFilters(filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set event filters: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Event filters set successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetEventFilters(w http.ResponseWriter, r *http.Request) {
	filters, err := syn2600.GetEventFilters()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get event filters: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filters)
}

func (api *SYN2600API) GetNotifications(w http.ResponseWriter, r *http.Request) {
	notifications, err := syn2600.GetNotifications()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get notifications: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (api *SYN2600API) GetEventStats(w http.ResponseWriter, r *http.Request) {
	stats, err := syn2600.GetEventStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get event stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Compliance & Validation Endpoints

func (api *SYN2600API) CheckCompliance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var complianceReq struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&complianceReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	compliant, err := syn2600.CheckCompliance(tokenId, complianceReq.AssetDetails, complianceReq.Owner, complianceReq.Shares, complianceReq.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check compliance for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": compliant,
		"token_id": tokenId,
		"message": "Compliance check result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) ValidateCompliance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var validationReq struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&validationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := syn2600.ValidateCompliance(tokenId, validationReq.AssetDetails, validationReq.Owner, validationReq.Shares, validationReq.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate compliance for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": valid,
		"token_id": tokenId,
		"message": "Compliance validation result",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetComplianceReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	report, err := syn2600.GetComplianceReport(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance report for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (api *SYN2600API) GetComplianceRules(w http.ResponseWriter, r *http.Request) {
	rules, err := syn2600.GetComplianceRules()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance rules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (api *SYN2600API) UpdateComplianceRules(w http.ResponseWriter, r *http.Request) {
	var rules map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rules); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.UpdateComplianceRules(rules)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update compliance rules: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Compliance rules updated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetComplianceViolations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	violations, err := syn2600.GetComplianceViolations(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance violations for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(violations)
}

func (api *SYN2600API) ComplianceAudit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var auditReq struct {
		Action string `json:"action" validate:"required"`
		Details string `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&auditReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.ComplianceAudit(tokenId, auditReq.Action, auditReq.Details)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to audit compliance for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Compliance audit initiated",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) CertifyCompliance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var certificationReq struct {
		AssetDetails string `json:"asset_details" validate:"required"`
		Owner        string `json:"owner" validate:"required"`
		Shares       float64 `json:"shares" validate:"required,min=0.000000000000000001"`
		ExpiryDate   time.Time `json:"expiry_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&certificationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.CertifyCompliance(tokenId, certificationReq.AssetDetails, certificationReq.Owner, certificationReq.Shares, certificationReq.ExpiryDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to certify compliance for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Compliance certified successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) InitiateRemediation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	var remediationReq struct {
		Issue string `json:"issue" validate:"required"`
		Details string `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&remediationReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := syn2600.InitiateRemediation(tokenId, remediationReq.Issue, remediationReq.Details)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initiate remediation for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token_id": tokenId,
		"message": "Remediation initiated successfully",
		"timestamp": time.Now(),
	})
}

func (api *SYN2600API) GetComplianceStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	status, err := syn2600.GetComplianceStatus(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance status for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api *SYN2600API) GetComplianceHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	history, err := syn2600.GetComplianceHistory(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance history for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (api *SYN2600API) GetSupportedFrameworks(w http.ResponseWriter, r *http.Request) {
	frameworks, err := syn2600.GetSupportedFrameworks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get supported frameworks: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(frameworks)
}

func (api *SYN2600API) ExportComplianceData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenId := vars["tokenId"]

	exportData, err := syn2600.ExportComplianceData(tokenId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export compliance data for token %s: %v", tokenId, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exportData)
}

func (api *SYN2600API) GetComplianceAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := syn2600.GetComplianceAlerts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance alerts: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func (api *SYN2600API) GetComplianceStats(w http.ResponseWriter, r *http.Request) {
	stats, err := syn2600.GetComplianceStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get compliance stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}