package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"synnergy_network/pkg/ledger"
	"time"
)

// ModelStatusQuery provides the current status of a deployed model by querying its active state in the ledger.
func ModelStatusQuery(modelID string, ledgerInstance *ledger.Ledger) (string, error) {
	// Query model status from ledger
	status, err := ledgerInstance.AiMLMLedger.GetModelStatus(modelID)
	if err != nil {
		return "", errors.New("failed to retrieve model status from ledger")
	}

	// Encrypt the retrieved status
	encryptedStatus, err := encryptData([]byte(status), modelID)
	if err != nil {
		return "", errors.New("failed to encrypt model status data")
	}

	return string(encryptedStatus), nil
}

// ModelDeployCheck verifies the integrity and success of a model deployment, recording the check result in the ledger.
func ModelDeployCheck(modelID string, ledgerInstance *ledger.Ledger) (bool, error) {
	// Retrieve deployment status
	deployStatus, err := ledgerInstance.AiMLMLedger.GetDeploymentStatus(modelID)
	if err != nil {
		return false, errors.New("failed to retrieve deployment status from ledger")
	}

	// Validate deployment status
	isValid := validateDeployment(deployStatus)

	// Create deployment check log
	deployCheckLog := map[string]interface{}{
		"ModelID":      modelID,
		"DeploymentOK": isValid,
		"Timestamp":    time.Now(),
	}

	// Convert log to JSON for encryption
	logData, err := json.Marshal(deployCheckLog)
	if err != nil {
		return false, errors.New("failed to encode deployment check log to JSON")
	}

	// Encrypt the log
	encryptedLog, err := encryptData(logData, modelID)
	if err != nil {
		return false, errors.New("failed to encrypt deployment check log")
	}

	// Record the encrypted log in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDeploymentCheck(modelID, "Checked", string(encryptedLog)); err != nil {
		return false, errors.New("failed to record deployment check in ledger")
	}

	return isValid, nil
}

// ModelTrainStatus provides an update on a model's training status, retrieving details from the ledger.
func ModelTrainStatus(modelID string, ledgerInstance *ledger.Ledger) (string, error) {
	// Retrieve training status from ledger
	trainStatus, err := ledgerInstance.AiMLMLedger.GetTrainingStatus(modelID)
	if err != nil {
		return "", errors.New("failed to retrieve training status from ledger")
	}

	// Encrypt the training status
	encryptedStatus, err := encryptData([]byte(trainStatus), modelID)
	if err != nil {
		return "", errors.New("failed to encrypt training status data")
	}

	return string(encryptedStatus), nil
}

// ModelRunCheck verifies that a model is correctly executing and records the outcome in the ledger for future audits.
func ModelRunCheck(modelID string, ledgerInstance *ledger.Ledger) (bool, error) {
	// Retrieve run status from ledger
	runStatus, err := ledgerInstance.AiMLMLedger.GetRunStatus(modelID)
	if err != nil {
		return false, errors.New("failed to retrieve run status from ledger")
	}

	// Verify run status
	isRunning := verifyRunStatus(runStatus)

	// Create run check log
	runCheckLog := map[string]interface{}{
		"ModelID":     modelID,
		"RunStatusOK": isRunning,
		"Timestamp":   time.Now(),
	}

	// Convert log to JSON for encryption
	logData, err := json.Marshal(runCheckLog)
	if err != nil {
		return false, errors.New("failed to encode run check log to JSON")
	}

	// Encrypt the log
	encryptedLog, err := encryptData(logData, modelID)
	if err != nil {
		return false, errors.New("failed to encrypt run check log")
	}

	// Record the encrypted log in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordRunCheck(modelID, string(encryptedLog)); err != nil {
		return false, errors.New("failed to record run check in ledger")
	}

	return isRunning, nil
}

// validateDeployment checks the validity of a model deployment based on retrieved status details.
func validateDeployment(deployStatus string) bool {
	return deployStatus == "success"
}

// verifyRunStatus confirms if the model run status is active and functioning correctly.
func verifyRunStatus(runStatus string) bool {
	return runStatus == "running"
}
