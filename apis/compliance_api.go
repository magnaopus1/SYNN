package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/compliance"
	"synnergy_network/pkg/ledger"
)

// ComplianceAPI exposes compliance operations
type ComplianceAPI struct {
	ledger *ledger.Ledger
}

// NewComplianceAPI creates a ComplianceAPI
func NewComplianceAPI(l *ledger.Ledger) *ComplianceAPI {
	return &ComplianceAPI{ledger: l}
}

// RegisterRoutes registers compliance routes
func (api *ComplianceAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/compliance/kyc", api.SubmitKYC).Methods("POST")
	r.HandleFunc("/compliance/kyc/{id}", api.GetKYC).Methods("GET")
	r.HandleFunc("/compliance/check/{id}", api.Check).Methods("GET")
	r.HandleFunc("/compliance/health", api.Health).Methods("GET")
}

type kycRequest struct {
	EntityID string `json:"entity_id"`
	Data     string `json:"data"`
}

// SubmitKYC records a KYC entry in the ledger
func (api *ComplianceAPI) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	var req kycRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	record := ledger.KYCRecord{UserID: req.EntityID}
	if err := compliance.SubmitKYC(api.ledger, req.EntityID, record); err != nil {
		http.Error(w, fmt.Sprintf("failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true})
}

// GetKYC retrieves a KYC record
func (api *ComplianceAPI) GetKYC(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	rec, err := compliance.RetrieveKYCRecord(api.ledger, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("not found: %v", err), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(rec)
}

// Check runs a compliance check
func (api *ComplianceAPI) Check(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ok, err := compliance.CheckCompliance(api.ledger, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("check failed: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"compliant": ok})
}

func (api *ComplianceAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
