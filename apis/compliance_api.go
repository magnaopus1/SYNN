package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/compliance"
	"synnergy_network/pkg/ledger"
)

// ComplianceAPI provides compliance operations endpoints
type ComplianceAPI struct {
	LedgerInstance *ledger.Ledger
}

// NewComplianceAPI creates a new ComplianceAPI instance
func NewComplianceAPI(ledgerInstance *ledger.Ledger) *ComplianceAPI {
	return &ComplianceAPI{LedgerInstance: ledgerInstance}
}

// RegisterRoutes registers compliance routes
func (api *ComplianceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/compliance/kyc", api.SubmitKYC).Methods("POST")
	router.HandleFunc("/compliance/kyc/{id}", api.VerifyKYC).Methods("GET")
	router.HandleFunc("/compliance/record/{id}", api.GetComplianceRecord).Methods("GET")
	router.HandleFunc("/compliance/restrict/{id}", api.ApplyRestriction).Methods("POST")
}

// SubmitKYC handles KYC submission
func (api *ComplianceAPI) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EntityID string           `json:"entity_id"`
		Record   ledger.KYCRecord `json:"record"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := compliance.SubmitKYC(api.LedgerInstance, req.EntityID, req.Record); err != nil {
		http.Error(w, fmt.Sprintf("submit failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// VerifyKYC verifies a KYC record
func (api *ComplianceAPI) VerifyKYC(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ok, err := compliance.VerifyKYC(api.LedgerInstance, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("verify failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"verified": ok})
}

// GetComplianceRecord retrieves compliance record
func (api *ComplianceAPI) GetComplianceRecord(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	rec, err := compliance.RetrieveComplianceRecord(api.LedgerInstance, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("retrieve failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rec)
}

// ApplyRestriction applies compliance restriction
func (api *ComplianceAPI) ApplyRestriction(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := compliance.ApplyRestrictions(api.LedgerInstance, id, req.Reason); err != nil {
		http.Error(w, fmt.Sprintf("restriction failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
