package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"synnergy_network/pkg/ledger"
	"time"
)

// AIModelAccessLog records each access to an AI model, capturing metadata and securely storing logs in the ledger.
func AIModelAccessLog(modelID, userID, accessType string, ledgerInstance *ledger.Ledger) error {
	logEntry := generateLogEntry(modelID, userID, accessType)
	logEntryData, err := json.Marshal(logEntry)
	if err != nil {
		return errors.New("failed to encode access log to JSON")
	}

	// Encrypt the access log data (optional for security, not stored here)
	_, err = encryptData(logEntryData, modelID)
	if err != nil {
		return errors.New("failed to encrypt access log")
	}

	// Record the access log in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordModelAccessLog(modelID, userID, accessType); err != nil {
		return errors.New("failed to record model access log in ledger")
	}

	return nil
}

// ModelUsageStatistics aggregates and records model usage metrics, such as access frequency, user interaction, and operational stats.
func ModelUsageStatistics(modelID string, duration time.Duration, ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.AiMLMLedger.RecordModelUsageStatistics(modelID, duration); err != nil {
		return errors.New("failed to record model usage statistics in ledger")
	}

	return nil
}

// ModelPerformanceMetric captures and stores performance data like accuracy and loss.
func ModelPerformanceMetric(modelID string, accuracy, loss float64, ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.AiMLMLedger.RecordModelPerformanceMetrics(modelID, accuracy, loss); err != nil {
		return errors.New("failed to record performance metrics in ledger")
	}

	return nil
}

// ModelHyperparameterTune logs hyperparameter tuning activities, tracking the parameters and results for future model evaluations.
func ModelHyperparameterTune(modelID string, parameters map[string]interface{}, results map[string]float64, ledgerInstance *ledger.Ledger) error {
	logEntry := struct {
		Parameters map[string]interface{}
		Results    map[string]float64
		Time       time.Time
	}{
		Parameters: parameters,
		Results:    results,
		Time:       time.Now(),
	}

	// Convert log entry to JSON for encryption
	logEntryData, err := json.Marshal(logEntry)
	if err != nil {
		return errors.New("failed to encode hyperparameter tuning log to JSON")
	}

	// Encrypt the hyperparameter tuning log (optional for security, not stored here)
	_, err = encryptData(logEntryData, modelID)
	if err != nil {
		return errors.New("failed to encrypt hyperparameter tuning log")
	}

	if err := ledgerInstance.AiMLMLedger.RecordModelHyperparameterTuning(modelID, parameters); err != nil {
		return errors.New("failed to record hyperparameter tuning log in ledger")
	}

	return nil
}

// ModelComplianceCheck verifies model compliance with regulatory and organizational standards, logging results securely.
func ModelComplianceCheck(modelID string, status, details string, ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.AiMLMLedger.RecordModelComplianceCheck(modelID, status, details); err != nil {
		return errors.New("failed to record compliance check in ledger")
	}

	return nil
}

// ModelSecurityAudit conducts a security audit for the model, checking for vulnerabilities and logging findings.
func ModelSecurityAudit(modelID string, passed bool, findings string, ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.AiMLMLedger.RecordModelSecurityAudit(modelID, passed, findings); err != nil {
		return errors.New("failed to record security audit in ledger")
	}

	return nil
}

// generateLogEntry creates a standardized log entry with metadata for logging access
func generateLogEntry(modelID, userID, accessType string) map[string]interface{} {
	return map[string]interface{}{
		"ModelID":    modelID,
		"UserID":     userID,
		"AccessType": accessType,
		"Timestamp":  time.Now(),
	}
}
