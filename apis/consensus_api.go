package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/ledger"

	"github.com/gorilla/mux"
)

// ConsensusAPI handles all consensus-related API endpoints
type ConsensusAPI struct {
	LedgerInstance *ledger.Ledger
}

// NewConsensusAPI creates a new consensus API instance
func NewConsensusAPI(ledgerInstance *ledger.Ledger) *ConsensusAPI {
	return &ConsensusAPI{
		LedgerInstance: ledgerInstance,
	}
}

// RegisterRoutes registers all consensus API routes
func (api *ConsensusAPI) RegisterRoutes(router *mux.Router) {
	// Difficulty management
	router.HandleFunc("/consensus/difficulty/adjust", api.AdjustDifficulty).Methods("POST")
	router.HandleFunc("/consensus/difficulty/monitor", api.MonitorBlockGeneration).Methods("POST")
	
	// Audit management
	router.HandleFunc("/consensus/audit/enable", api.EnableAudit).Methods("POST")
	router.HandleFunc("/consensus/audit/disable", api.DisableAudit).Methods("POST")
	router.HandleFunc("/consensus/audit/logs", api.GetAuditLogs).Methods("GET")
	
	// Reward distribution
	router.HandleFunc("/consensus/rewards/mode", api.SetRewardMode).Methods("POST")
	router.HandleFunc("/consensus/rewards/mode", api.GetRewardMode).Methods("GET")
	
	// Validator management
	router.HandleFunc("/consensus/validators/participation", api.TrackParticipation).Methods("POST")
	router.HandleFunc("/consensus/validators/selection-mode", api.SetValidatorSelectionMode).Methods("POST")
	router.HandleFunc("/consensus/validators/selection-mode", api.GetValidatorSelectionMode).Methods("GET")
	router.HandleFunc("/consensus/validators/activity", api.ValidateActivity).Methods("POST")
	router.HandleFunc("/consensus/validators/activity/logs", api.GetActivityLogs).Methods("GET")
	
	// PoH management
	router.HandleFunc("/consensus/poh/threshold", api.SetPoHThreshold).Methods("POST")
	router.HandleFunc("/consensus/poh/threshold", api.GetPoHThreshold).Methods("GET")
	
	// Dynamic stake adjustment
	router.HandleFunc("/consensus/stake/dynamic/enable", api.EnableDynamicStake).Methods("POST")
	router.HandleFunc("/consensus/stake/dynamic/disable", api.DisableDynamicStake).Methods("POST")
}

// AdjustDifficulty adjusts consensus difficulty based on time
func (api *ConsensusAPI) AdjustDifficulty(w http.ResponseWriter, r *http.Request) {
	var request struct {
		NewLevel int    `json:"new_level"`
		Reason   string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.ConsensusAdjustDifficultyBasedOnTime(request.NewLevel, request.Reason, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to adjust difficulty: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Difficulty adjusted successfully",
		"new_level": request.NewLevel,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MonitorBlockGeneration logs block generation time
func (api *ConsensusAPI) MonitorBlockGeneration(w http.ResponseWriter, r *http.Request) {
	var request struct {
		BlockID        string `json:"block_id"`
		GenerationTime int64  `json:"generation_time_ms"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	generationTime := time.Duration(request.GenerationTime) * time.Millisecond
	err := consensus.consensusMonitorBlockGenerationTime(request.BlockID, generationTime, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to monitor block generation: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Block generation time logged successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EnableAudit enables consensus audit
func (api *ConsensusAPI) EnableAudit(w http.ResponseWriter, r *http.Request) {
	err := consensus.consensusEnableConsensusAudit(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enable audit: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Consensus audit enabled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DisableAudit disables consensus audit
func (api *ConsensusAPI) DisableAudit(w http.ResponseWriter, r *http.Request) {
	err := consensus.consensusDisableConsensusAudit(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to disable audit: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Consensus audit disabled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLogs retrieves consensus audit logs
func (api *ConsensusAPI) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := consensus.ConsensusFetchConsensusLogs(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"logs": logs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetRewardMode sets the reward distribution mode
func (api *ConsensusAPI) SetRewardMode(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Mode ledger.RewardDistributionMode `json:"mode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.consensusSetRewardDistributionMode(request.Mode, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set reward mode: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Reward distribution mode set successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRewardMode retrieves the current reward distribution mode
func (api *ConsensusAPI) GetRewardMode(w http.ResponseWriter, r *http.Request) {
	mode, err := consensus.consensusGetRewardDistributionMode(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get reward mode: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"mode": mode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TrackParticipation tracks validator participation in consensus
func (api *ConsensusAPI) TrackParticipation(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ValidatorID string `json:"validator_id"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.ConsensusTrackConsensusParticipation(request.ValidatorID, request.Status, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to track participation: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Participation tracked successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetValidatorSelectionMode sets the validator selection mode
func (api *ConsensusAPI) SetValidatorSelectionMode(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Mode ledger.ValidatorSelectionMode `json:"mode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.ConsensusSetValidatorSelectionMode(request.Mode, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set validator selection mode: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Validator selection mode set successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetValidatorSelectionMode retrieves the current validator selection mode
func (api *ConsensusAPI) GetValidatorSelectionMode(w http.ResponseWriter, r *http.Request) {
	mode, err := consensus.ConsensusGetValidatorSelectionMode(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get validator selection mode: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"mode": mode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ValidateActivity validates validator activity
func (api *ConsensusAPI) ValidateActivity(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ValidatorID string `json:"validator_id"`
		Action      string `json:"action"`
		Details     string `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.ConsensusValidateValidatorActivity(request.ValidatorID, request.Action, request.Details, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to validate activity: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Validator activity validated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetActivityLogs retrieves validator activity logs
func (api *ConsensusAPI) GetActivityLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := consensus.ConsensusFetchValidatorActivityLogs(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch activity logs: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"logs": logs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetPoHThreshold sets the PoH participation threshold
func (api *ConsensusAPI) SetPoHThreshold(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Threshold float64 `json:"threshold"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := consensus.ConsensusSetPoHParticipationThreshold(request.Threshold, api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set PoH threshold: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "PoH participation threshold set successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPoHThreshold retrieves the current PoH participation threshold
func (api *ConsensusAPI) GetPoHThreshold(w http.ResponseWriter, r *http.Request) {
	threshold, err := consensus.ConsensusGetPoHParticipationThreshold(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get PoH threshold: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"threshold": threshold,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EnableDynamicStake enables dynamic stake adjustment
func (api *ConsensusAPI) EnableDynamicStake(w http.ResponseWriter, r *http.Request) {
	err := consensus.ConsensusEnableDynamicStakeAdjustment(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enable dynamic stake: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Dynamic stake adjustment enabled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DisableDynamicStake disables dynamic stake adjustment
func (api *ConsensusAPI) DisableDynamicStake(w http.ResponseWriter, r *http.Request) {
	err := consensus.ConsensusDisableDynamicStakeAdjustment(api.LedgerInstance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to disable dynamic stake: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Dynamic stake adjustment disabled successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}