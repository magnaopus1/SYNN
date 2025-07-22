package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/compliance"
	"synnergy_network/pkg/ledger"
)

// ComplianceAPI exposes compliance module features
type ComplianceAPI struct {
	ledger *ledger.Ledger
}

// NewComplianceAPI creates a ComplianceAPI
func NewComplianceAPI(ledgerInstance *ledger.Ledger) *ComplianceAPI {
	return &ComplianceAPI{ledger: ledgerInstance}
}

// RegisterRoutes sets up compliance endpoints
func (api *ComplianceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/compliance/check", api.CheckCompliance).Methods("POST")
	router.HandleFunc("/compliance/kyc/submit", api.SubmitKYC).Methods("POST")
	router.HandleFunc("/compliance/kyc/verify", api.VerifyKYC).Methods("GET")
}

type checkReq struct {
	EntityID string `json:"entity_id"`
}

// CheckCompliance verifies entity compliance
func (api *ComplianceAPI) CheckCompliance(w http.ResponseWriter, r *http.Request) {
	var req checkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ok, err := compliance.CheckCompliance(api.ledger, req.EntityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("compliance check failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"compliant": ok, "entity_id": req.EntityID, "timestamp": time.Now()})
}

type kycSubmitReq struct {
	EntityID string           `json:"entity_id"`
	DataHash string           `json:"data_hash"`
	Status   ledger.KYCStatus `json:"status"`
}

// SubmitKYC stores KYC information
func (api *ComplianceAPI) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	var req kycSubmitReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	kycRecord := ledger.KYCRecord{UserID: req.EntityID, Status: req.Status, DataHash: req.DataHash}
	if err := compliance.SubmitKYC(api.ledger, req.EntityID, kycRecord); err != nil {
		http.Error(w, fmt.Sprintf("submit failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "entity_id": req.EntityID, "timestamp": time.Now()})
}

// VerifyKYC verifies stored KYC information
func (api *ComplianceAPI) VerifyKYC(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entity_id")
	ok, err := compliance.VerifyKYC(api.ledger, entityID)
	if err != nil {
		http.Error(w, fmt.Sprintf("verification failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"verified": ok, "entity_id": entityID, "timestamp": time.Now()})
}
