package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"synnergy_network/pkg/compliance"
	"synnergy_network/pkg/ledger"
)

// ComplianceAPI exposes compliance related operations
type ComplianceAPI struct {
	ledgerInstance *ledger.Ledger
}

// NewComplianceAPI creates a new ComplianceAPI
func NewComplianceAPI(ledgerInstance *ledger.Ledger) *ComplianceAPI {
	return &ComplianceAPI{ledgerInstance: ledgerInstance}
}

// RegisterRoutes registers compliance API routes
func (api *ComplianceAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/compliance/check/{id}", api.CheckCompliance).Methods("GET")
	router.HandleFunc("/compliance/kyc/{id}", api.SubmitKYC).Methods("POST")
	router.HandleFunc("/compliance/record/{id}", api.GetComplianceRecord).Methods("GET")
}

// CheckCompliance verifies if an entity is compliant
func (api *ComplianceAPI) CheckCompliance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ok, err := compliance.CheckCompliance(api.ledgerInstance, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"compliant": ok})
}

// SubmitKYC stores KYC data
func (api *ComplianceAPI) SubmitKYC(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var kyc ledger.KYCRecord
	if err := json.NewDecoder(r.Body).Decode(&kyc); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := compliance.SubmitKYC(api.ledgerInstance, id, kyc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "kyc stored"})
}

// GetComplianceRecord retrieves a compliance record
func (api *ComplianceAPI) GetComplianceRecord(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	rec, err := compliance.RetrieveComplianceRecord(api.ledgerInstance, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(rec)
}
