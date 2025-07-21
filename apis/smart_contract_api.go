package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"synnergy_network/pkg/smart_contract"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/common"

	"github.com/gorilla/mux"
)

// SmartContractAPI handles all smart contract-related API endpoints
type SmartContractAPI struct {
	LedgerInstance        *ledger.Ledger
	ContractManager       *smart_contract.SmartContractManager
	MigrationManager      *smart_contract.MigrationManager
	RicardianManager      *smart_contract.RicardianContractManager
	TemplateMarketplace   *smart_contract.SmartContractTemplateMarketplace
}

// NewSmartContractAPI creates a new smart contract API instance
func NewSmartContractAPI(ledgerInstance *ledger.Ledger) *SmartContractAPI {
	return &SmartContractAPI{
		LedgerInstance:      ledgerInstance,
		ContractManager:     &smart_contract.SmartContractManager{
			Contracts:      make(map[string]*common.SmartContract),
			LedgerInstance: ledgerInstance,
		},
		MigrationManager:    &smart_contract.MigrationManager{
			Contracts:      make(map[string]*smart_contract.MigratedContract),
			LedgerInstance: ledgerInstance,
		},
		RicardianManager:    &smart_contract.RicardianContractManager{
			Contracts:      make(map[string]*smart_contract.RicardianContract),
			LedgerInstance: ledgerInstance,
		},
		TemplateMarketplace: &smart_contract.SmartContractTemplateMarketplace{
			Templates:      make(map[string]*smart_contract.SmartContractTemplate),
			Escrows:        make(map[string]*smart_contract.Escrow),
			EscrowFee:      0.025, // 2.5% marketplace fee
			LedgerInstance: ledgerInstance,
		},
	}
}

// RegisterRoutes registers all smart contract API routes
func (api *SmartContractAPI) RegisterRoutes(router *mux.Router) {
	// Smart contract management
	router.HandleFunc("/smart-contracts/deploy", api.DeployContract).Methods("POST")
	router.HandleFunc("/smart-contracts/execute", api.ExecuteContract).Methods("POST")
	router.HandleFunc("/smart-contracts/{contractID}", api.GetContract).Methods("GET")
	router.HandleFunc("/smart-contracts/{contractID}", api.UpdateContract).Methods("PUT")
	router.HandleFunc("/smart-contracts/{contractID}", api.DeleteContract).Methods("DELETE")
	router.HandleFunc("/smart-contracts", api.ListContracts).Methods("GET")
	
	// Contract execution and state
	router.HandleFunc("/smart-contracts/{contractID}/execute/{function}", api.ExecuteFunction).Methods("POST")
	router.HandleFunc("/smart-contracts/{contractID}/state", api.GetContractState).Methods("GET")
	router.HandleFunc("/smart-contracts/{contractID}/history", api.GetExecutionHistory).Methods("GET")
	
	// Ricardian contracts
	router.HandleFunc("/ricardian-contracts/deploy", api.DeployRicardianContract).Methods("POST")
	router.HandleFunc("/ricardian-contracts/{contractID}", api.GetRicardianContract).Methods("GET")
	router.HandleFunc("/ricardian-contracts/{contractID}/sign", api.SignRicardianContract).Methods("POST")
	router.HandleFunc("/ricardian-contracts/{contractID}/execute", api.ExecuteRicardianContract).Methods("POST")
	router.HandleFunc("/ricardian-contracts", api.ListRicardianContracts).Methods("GET")
	
	// Contract migration
	router.HandleFunc("/smart-contracts/{contractID}/migrate", api.MigrateContract).Methods("POST")
	router.HandleFunc("/smart-contracts/migrations", api.ListMigrations).Methods("GET")
	router.HandleFunc("/smart-contracts/migrations/{migrationID}", api.GetMigration).Methods("GET")
	
	// Template marketplace
	router.HandleFunc("/contract-templates", api.ListTemplates).Methods("GET")
	router.HandleFunc("/contract-templates/{templateID}", api.GetTemplate).Methods("GET")
	router.HandleFunc("/contract-templates", api.CreateTemplate).Methods("POST")
	router.HandleFunc("/contract-templates/{templateID}/purchase", api.PurchaseTemplate).Methods("POST")
	router.HandleFunc("/contract-templates/{templateID}", api.UpdateTemplate).Methods("PUT")
	router.HandleFunc("/contract-templates/{templateID}", api.DeleteTemplate).Methods("DELETE")
	
	// Escrow management
	router.HandleFunc("/escrows", api.ListEscrows).Methods("GET")
	router.HandleFunc("/escrows/{escrowID}", api.GetEscrow).Methods("GET")
	router.HandleFunc("/escrows/{escrowID}/release", api.ReleaseEscrow).Methods("POST")
	router.HandleFunc("/escrows/{escrowID}/dispute", api.DisputeEscrow).Methods("POST")
}

// DeployContract deploys a new smart contract
func (api *SmartContractAPI) DeployContract(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string                 `json:"name"`
		Code        string                 `json:"code"`
		Constructor map[string]interface{} `json:"constructor"`
		Owner       string                 `json:"owner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate contract ID
	contractID := generateContractID(request.Name, request.Owner)

	// Create smart contract
	contract := &common.SmartContract{
		ID:           contractID,
		Name:         request.Name,
		Code:         request.Code,
		Owner:        request.Owner,
		State:        make(map[string]interface{}),
		Executions:   []common.ContractExecution{},
		DeployedAt:   time.Now(),
		IsActive:     true,
	}

	// Deploy to ledger
	if err := api.LedgerInstance.DeploySmartContract(contract); err != nil {
		http.Error(w, fmt.Sprintf("Failed to deploy contract: %v", err), http.StatusInternalServerError)
		return
	}

	// Add to contract manager
	api.ContractManager.Contracts[contractID] = contract

	response := map[string]interface{}{
		"success":     true,
		"message":     "Smart contract deployed successfully",
		"contract_id": contractID,
		"contract":    contract,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExecuteContract executes a smart contract function
func (api *SmartContractAPI) ExecuteContract(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ContractID string                 `json:"contract_id"`
		Function   string                 `json:"function"`
		Parameters map[string]interface{} `json:"parameters"`
		Caller     string                 `json:"caller"`
		GasLimit   uint64                 `json:"gas_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get contract
	contract, exists := api.ContractManager.Contracts[request.ContractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Execute contract
	execution, err := api.executeContractFunction(contract, request.Function, request.Parameters, request.Caller, request.GasLimit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Contract execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Record execution
	contract.Executions = append(contract.Executions, *execution)

	response := map[string]interface{}{
		"success":   true,
		"message":   "Contract executed successfully",
		"execution": execution,
		"result":    execution.Result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetContract retrieves a smart contract by ID
func (api *SmartContractAPI) GetContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"contract": contract,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateContract updates a smart contract
func (api *SmartContractAPI) UpdateContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	var request struct {
		Code     string `json:"code"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Update contract
	contract.Code = request.Code
	contract.IsActive = request.IsActive

	// Update in ledger
	if err := api.LedgerInstance.UpdateSmartContract(contract); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update contract: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"message":  "Contract updated successfully",
		"contract": contract,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteContract deactivates a smart contract
func (api *SmartContractAPI) DeleteContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Deactivate contract
	contract.IsActive = false

	// Update in ledger
	if err := api.LedgerInstance.UpdateSmartContract(contract); err != nil {
		http.Error(w, fmt.Sprintf("Failed to deactivate contract: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Contract deactivated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListContracts returns all smart contracts
func (api *SmartContractAPI) ListContracts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	owner := r.URL.Query().Get("owner")
	activeOnly := r.URL.Query().Get("active") == "true"

	var contracts []*common.SmartContract
	for _, contract := range api.ContractManager.Contracts {
		// Filter by owner if specified
		if owner != "" && contract.Owner != owner {
			continue
		}
		// Filter by active status if specified
		if activeOnly && !contract.IsActive {
			continue
		}
		contracts = append(contracts, contract)
	}

	response := map[string]interface{}{
		"success":   true,
		"contracts": contracts,
		"count":     len(contracts),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExecuteFunction executes a specific function of a contract
func (api *SmartContractAPI) ExecuteFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]
	function := vars["function"]

	var request struct {
		Parameters map[string]interface{} `json:"parameters"`
		Caller     string                 `json:"caller"`
		GasLimit   uint64                 `json:"gas_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	execution, err := api.executeContractFunction(contract, function, request.Parameters, request.Caller, request.GasLimit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Function execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"execution": execution,
		"result":    execution.Result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetContractState returns the current state of a contract
func (api *SmartContractAPI) GetContractState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"contract_id": contractID,
		"state":       contract.State,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetExecutionHistory returns the execution history of a contract
func (api *SmartContractAPI) GetExecutionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	executions := contract.Executions
	if len(executions) > limit {
		executions = executions[len(executions)-limit:]
	}

	response := map[string]interface{}{
		"success":    true,
		"contract_id": contractID,
		"executions": executions,
		"count":      len(executions),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeployRicardianContract deploys a new Ricardian contract
func (api *SmartContractAPI) DeployRicardianContract(w http.ResponseWriter, r *http.Request) {
	var request struct {
		HumanReadable   string   `json:"human_readable"`
		MachineReadable string   `json:"machine_readable"`
		PartiesInvolved []string `json:"parties_involved"`
		Owner           string   `json:"owner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contractID := generateContractID("ricardian", request.Owner)

	contract := &smart_contract.RicardianContract{
		ID:              contractID,
		HumanReadable:   request.HumanReadable,
		MachineReadable: request.MachineReadable,
		PartiesInvolved: request.PartiesInvolved,
		Signatures:      make(map[string]string),
		State:           make(map[string]interface{}),
		Owner:           request.Owner,
		Executions:      []smart_contract.ContractExecution{},
		LedgerInstance:  api.LedgerInstance,
	}

	// Add to manager
	api.RicardianManager.Contracts[contractID] = contract

	response := map[string]interface{}{
		"success":     true,
		"message":     "Ricardian contract deployed successfully",
		"contract_id": contractID,
		"contract":    contract,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRicardianContract retrieves a Ricardian contract
func (api *SmartContractAPI) GetRicardianContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	contract, exists := api.RicardianManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Ricardian contract not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"contract": contract,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SignRicardianContract adds a signature to a Ricardian contract
func (api *SmartContractAPI) SignRicardianContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	var request struct {
		PartyAddress string `json:"party_address"`
		Signature    string `json:"signature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract, exists := api.RicardianManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Ricardian contract not found", http.StatusNotFound)
		return
	}

	// Add signature
	contract.Signatures[request.PartyAddress] = request.Signature

	response := map[string]interface{}{
		"success": true,
		"message": "Contract signed successfully",
		"signatures": len(contract.Signatures),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExecuteRicardianContract executes a Ricardian contract
func (api *SmartContractAPI) ExecuteRicardianContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	var request struct {
		Function   string                 `json:"function"`
		Parameters map[string]interface{} `json:"parameters"`
		Executor   string                 `json:"executor"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract, exists := api.RicardianManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Ricardian contract not found", http.StatusNotFound)
		return
	}

	// Check if all required parties have signed
	allSigned := true
	for _, party := range contract.PartiesInvolved {
		if _, signed := contract.Signatures[party]; !signed {
			allSigned = false
			break
		}
	}

	if !allSigned {
		http.Error(w, "Contract not fully signed by all parties", http.StatusBadRequest)
		return
	}

	// Execute contract
	execution := smart_contract.ContractExecution{
		ExecutionID:   generateExecutionID(),
		ContractID:    contractID,
		FunctionName:  request.Function,
		Parameters:    request.Parameters,
		Result:        map[string]interface{}{"status": "executed"},
		ExecutionTime: time.Now(),
		GasUsed:       1000, // Simplified gas calculation
		Executor:      request.Executor,
	}

	contract.Executions = append(contract.Executions, execution)

	response := map[string]interface{}{
		"success":   true,
		"message":   "Ricardian contract executed successfully",
		"execution": execution,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListRicardianContracts returns all Ricardian contracts
func (api *SmartContractAPI) ListRicardianContracts(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")

	var contracts []*smart_contract.RicardianContract
	for _, contract := range api.RicardianManager.Contracts {
		if owner != "" && contract.Owner != owner {
			continue
		}
		contracts = append(contracts, contract)
	}

	response := map[string]interface{}{
		"success":   true,
		"contracts": contracts,
		"count":     len(contracts),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MigrateContract migrates a contract to a new version
func (api *SmartContractAPI) MigrateContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractID := vars["contractID"]

	var request struct {
		NewCode       string                 `json:"new_code"`
		NewParameters map[string]interface{} `json:"new_parameters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	contract, exists := api.ContractManager.Contracts[contractID]
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Create migration
	migrationID := generateMigrationID()
	migration := &smart_contract.MigratedContract{
		OldContractID:   contractID,
		NewContractID:   generateContractID("migrated", contract.Owner),
		Owner:           contract.Owner,
		NewCode:         request.NewCode,
		NewParameters:   request.NewParameters,
		MigrationTime:   time.Now(),
	}

	// Add to migration manager
	api.MigrationManager.Contracts[migrationID] = migration

	response := map[string]interface{}{
		"success":       true,
		"message":       "Contract migration completed",
		"migration_id":  migrationID,
		"new_contract_id": migration.NewContractID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListMigrations returns all contract migrations
func (api *SmartContractAPI) ListMigrations(w http.ResponseWriter, r *http.Request) {
	var migrations []*smart_contract.MigratedContract
	for _, migration := range api.MigrationManager.Contracts {
		migrations = append(migrations, migration)
	}

	response := map[string]interface{}{
		"success":    true,
		"migrations": migrations,
		"count":      len(migrations),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMigration retrieves a specific migration
func (api *SmartContractAPI) GetMigration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	migrationID := vars["migrationID"]

	migration, exists := api.MigrationManager.Contracts[migrationID]
	if !exists {
		http.Error(w, "Migration not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"migration": migration,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListTemplates returns all contract templates
func (api *SmartContractAPI) ListTemplates(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	
	var templates []*smart_contract.SmartContractTemplate
	for _, template := range api.TemplateMarketplace.Templates {
		if category != "" && template.Name != category { // Simplified category filtering
			continue
		}
		templates = append(templates, template)
	}

	response := map[string]interface{}{
		"success":   true,
		"templates": templates,
		"count":     len(templates),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTemplate retrieves a specific template
func (api *SmartContractAPI) GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["templateID"]

	template, exists := api.TemplateMarketplace.Templates[templateID]
	if !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"template": template,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateTemplate creates a new contract template
func (api *SmartContractAPI) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Creator     string  `json:"creator"`
		Code        string  `json:"code"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	templateID := generateTemplateID()
	template := &smart_contract.SmartContractTemplate{
		ID:          templateID,
		Name:        request.Name,
		Description: request.Description,
		Creator:     request.Creator,
		Code:        request.Code,
		Price:       request.Price,
		Timestamp:   time.Now(),
	}

	api.TemplateMarketplace.Templates[templateID] = template

	response := map[string]interface{}{
		"success":     true,
		"message":     "Template created successfully",
		"template_id": templateID,
		"template":    template,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PurchaseTemplate handles template purchase with escrow
func (api *SmartContractAPI) PurchaseTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["templateID"]

	var request struct {
		Buyer  string  `json:"buyer"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	template, exists := api.TemplateMarketplace.Templates[templateID]
	if !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	// Create escrow
	escrowID := generateEscrowID()
	escrow := &smart_contract.Escrow{
		EscrowID:   escrowID,
		Buyer:      request.Buyer,
		Seller:     template.Creator,
		Amount:     request.Amount,
		ResourceID: templateID,
		Timestamp:  time.Now(),
		IsReleased: false,
		IsDisputed: false,
		Status:     "active",
	}

	api.TemplateMarketplace.Escrows[escrowID] = escrow

	response := map[string]interface{}{
		"success":   true,
		"message":   "Template purchase initiated",
		"escrow_id": escrowID,
		"escrow":    escrow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateTemplate updates a template
func (api *SmartContractAPI) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["templateID"]

	var request struct {
		Description string  `json:"description"`
		Code        string  `json:"code"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	template, exists := api.TemplateMarketplace.Templates[templateID]
	if !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	// Update template
	template.Description = request.Description
	template.Code = request.Code
	template.Price = request.Price

	response := map[string]interface{}{
		"success":  true,
		"message":  "Template updated successfully",
		"template": template,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteTemplate deletes a template
func (api *SmartContractAPI) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateID := vars["templateID"]

	if _, exists := api.TemplateMarketplace.Templates[templateID]; !exists {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	delete(api.TemplateMarketplace.Templates, templateID)

	response := map[string]interface{}{
		"success": true,
		"message": "Template deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListEscrows returns all escrows
func (api *SmartContractAPI) ListEscrows(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var escrows []*smart_contract.Escrow
	for _, escrow := range api.TemplateMarketplace.Escrows {
		if status != "" && escrow.Status != status {
			continue
		}
		escrows = append(escrows, escrow)
	}

	response := map[string]interface{}{
		"success": true,
		"escrows": escrows,
		"count":   len(escrows),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEscrow retrieves a specific escrow
func (api *SmartContractAPI) GetEscrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	escrowID := vars["escrowID"]

	escrow, exists := api.TemplateMarketplace.Escrows[escrowID]
	if !exists {
		http.Error(w, "Escrow not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"escrow":  escrow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ReleaseEscrow releases funds from escrow
func (api *SmartContractAPI) ReleaseEscrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	escrowID := vars["escrowID"]

	escrow, exists := api.TemplateMarketplace.Escrows[escrowID]
	if !exists {
		http.Error(w, "Escrow not found", http.StatusNotFound)
		return
	}

	// Release escrow
	escrow.IsReleased = true
	escrow.Status = "completed"
	escrow.CompletionTime = time.Now()

	response := map[string]interface{}{
		"success": true,
		"message": "Escrow funds released successfully",
		"escrow":  escrow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DisputeEscrow initiates an escrow dispute
func (api *SmartContractAPI) DisputeEscrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	escrowID := vars["escrowID"]

	var request struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	escrow, exists := api.TemplateMarketplace.Escrows[escrowID]
	if !exists {
		http.Error(w, "Escrow not found", http.StatusNotFound)
		return
	}

	// Mark as disputed
	escrow.IsDisputed = true
	escrow.Status = "disputed"

	response := map[string]interface{}{
		"success": true,
		"message": "Escrow dispute initiated",
		"escrow":  escrow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func generateContractID(name, owner string) string {
	return fmt.Sprintf("contract_%s_%s_%d", name, owner[:8], time.Now().UnixNano())
}

func generateExecutionID() string {
	return fmt.Sprintf("exec_%d", time.Now().UnixNano())
}

func generateMigrationID() string {
	return fmt.Sprintf("migration_%d", time.Now().UnixNano())
}

func generateTemplateID() string {
	return fmt.Sprintf("template_%d", time.Now().UnixNano())
}

func generateEscrowID() string {
	return fmt.Sprintf("escrow_%d", time.Now().UnixNano())
}

func (api *SmartContractAPI) executeContractFunction(contract *common.SmartContract, function string, parameters map[string]interface{}, caller string, gasLimit uint64) (*common.ContractExecution, error) {
	// Simplified contract execution
	execution := &common.ContractExecution{
		ExecutionID:   generateExecutionID(),
		ContractID:    contract.ID,
		FunctionName:  function,
		Parameters:    parameters,
		Result:        map[string]interface{}{"status": "success", "output": "function executed"},
		ExecutionTime: time.Now(),
		GasUsed:       1000, // Simplified gas calculation
		Executor:      caller,
	}

	// Update contract state (simplified)
	if contract.State == nil {
		contract.State = make(map[string]interface{})
	}
	contract.State["last_execution"] = execution.ExecutionID
	contract.State["execution_count"] = len(contract.Executions) + 1

	return execution, nil
}