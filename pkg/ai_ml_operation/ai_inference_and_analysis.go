package ai_ml_operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"

	"errors"
)

// AIInferenceRun executes an AI model inference with provided data inputs.
func AIInferenceRun(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID string, inputData []byte, transactionID string) (string, error) {
	logAction("AIInferenceRun - Start", fmt.Sprintf("ModelID: %s, NodeID: %s, TransactionID: %s", modelID, nodeID, transactionID))

	// Validate inputs
	if modelID == "" || nodeID == "" || transactionID == "" {
		return "", errors.New("modelID, nodeID, and transactionID cannot be empty")
	}
	if len(inputData) == 0 {
		return "", errors.New("input data cannot be empty")
	}

	// Fetch the AI model
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIInferenceRun - Model Retrieval Failed: %v", err)
		return "", fmt.Errorf("model retrieval failed: %w", err)
	}

	// Execute the inference logic
	inferenceOutput, err := model.RunInference(inputData, nodeID)
	if err != nil {
		log.Printf("AIInferenceRun - Inference Execution Failed: %v", err)
		return "", fmt.Errorf("inference execution failed: %w", err)
	}

	// Encrypt the inference result
	encryptedOutput, err := AIModelEncrypt(inferenceOutput, modelID)
	if err != nil {
		log.Printf("AIInferenceRun - Encryption Failed: %v", err)
		return "", fmt.Errorf("encryption failed for inference output: %w", err)
	}

	// Record the inference result in the ledger
	if err := l.AiMLMLedger.RecordInference(transactionID, modelID, encryptedOutput); err != nil {
		log.Printf("AIInferenceRun - Ledger Recording Failed: %v", err)
		return "", fmt.Errorf("failed to record inference in ledger: %w", err)
	}

	// Process the transaction for validation
	if err := common.ProcessSingleTransaction(sc, transactionID, encryptedOutput); err != nil {
		log.Printf("AIInferenceRun - Transaction Processing Failed: %v", err)
		return "", fmt.Errorf("transaction processing failed: %w", err)
	}

	log.Printf("AIInferenceRun - Success: Encrypted Output Recorded")
	return encryptedOutput, nil
}



// AIAnalysisStart initiates an AI-driven analysis on data and records the action in the ledger.
func AIAnalysisStart(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, transactionID string) error {
	logAction("AIAnalysisStart - Start", fmt.Sprintf("ModelID: %s, NodeID: %s, TransactionID: %s", modelID, nodeID, transactionID))

	// Validate inputs
	if modelID == "" || nodeID == "" || transactionID == "" {
		return errors.New("modelID, nodeID, and transactionID cannot be empty")
	}

	// Retrieve the model for analysis
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIAnalysisStart - Model Retrieval Failed: %v", err)
		return fmt.Errorf("model retrieval failed: %w", err)
	}

	// Start the analysis
	if err := model.StartAnalysis(transactionID); err != nil {
		log.Printf("AIAnalysisStart - Analysis Start Failed: %v", err)
		return fmt.Errorf("analysis start failed: %w", err)
	}

	// Record the analysis start in the ledger
	logID, err := l.AiMLMLedger.RecordAnalysisStart(modelID, nodeID)
	if err != nil {
		log.Printf("AIAnalysisStart - Ledger Recording Failed: %v", err)
		return fmt.Errorf("failed to record analysis start: %w", err)
	}

	// Wrap the transaction for processing
	transaction := common.Transaction{
		TransactionID: transactionID,
		EncryptedData: logID,
		Status:        "Pending",
	}

	// Process the transaction
	if err := sc.ProcessTransactions([]common.Transaction{transaction}); err != nil {
		log.Printf("AIAnalysisStart - Transaction Processing Failed: %v", err)
		return fmt.Errorf("transaction processing failed: %w", err)
	}

	log.Printf("AIAnalysisStart - Success: Analysis Started and Recorded")
	return nil
}



// AIAnalysisStop terminates an ongoing analysis and records the result in the ledger.
func AIAnalysisStop(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, transactionID string) (string, error) {
	logAction("AIAnalysisStop - Start", fmt.Sprintf("ModelID: %s, NodeID: %s, TransactionID: %s", modelID, nodeID, transactionID))

	// Validate inputs
	if modelID == "" || nodeID == "" || transactionID == "" {
		return "", errors.New("modelID, nodeID, and transactionID cannot be empty")
	}

	// Retrieve the model
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIAnalysisStop - Model Retrieval Failed: %v", err)
		return "", fmt.Errorf("model retrieval failed: %w", err)
	}

	// Stop the analysis and get the result
	analysisResult, err := model.StopAnalysis(nodeID)
	if err != nil {
		log.Printf("AIAnalysisStop - Analysis Stop Failed: %v", err)
		return "", fmt.Errorf("analysis stop failed: %w", err)
	}

	// Encrypt the analysis result
	encryptedResult, err := AIModelEncrypt(analysisResult, modelID)
	if err != nil {
		log.Printf("AIAnalysisStop - Encryption Failed: %v", err)
		return "", fmt.Errorf("encryption failed for analysis result: %w", err)
	}

	// Record the analysis stop in the ledger
	if err := l.AiMLMLedger.RecordAnalysisStop(transactionID); err != nil {
		log.Printf("AIAnalysisStop - Ledger Recording Failed: %v", err)
		return "", fmt.Errorf("failed to record analysis stop: %w", err)
	}

	// Process the transaction
	if err := common.ProcessSingleTransaction(sc, transactionID, encryptedResult); err != nil {
		log.Printf("AIAnalysisStop - Transaction Processing Failed: %v", err)
		return "", fmt.Errorf("transaction processing failed: %w", err)
	}

	log.Printf("AIAnalysisStop - Success: Analysis Stopped and Result Recorded")
	return encryptedResult, nil
}



// AIResultDispatch sends the analysis or inference result to specified parties securely.
func AIResultDispatch(l *ledger.Ledger, sc *common.SynnergyConsensus, resultData []byte, recipient string, transactionID string) error {
	logAction("AIResultDispatch - Start", fmt.Sprintf("TransactionID: %s, Recipient: %s", transactionID, recipient))

	// Validate inputs
	if len(resultData) == 0 {
		return errors.New("result data cannot be empty")
	}
	if recipient == "" {
		return errors.New("recipient cannot be empty")
	}
	if transactionID == "" {
		return errors.New("transactionID cannot be empty")
	}

	// Encrypt the result data for dispatch
	encryptedResult, err := AIModelEncrypt(resultData, transactionID)
	if err != nil {
		log.Printf("AIResultDispatch - Encryption Failed: %v", err)
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Dispatch the encrypted result
	if err := sendEncryptedResult(encryptedResult, recipient); err != nil {
		log.Printf("AIResultDispatch - Dispatch Failed: %v", err)
		return fmt.Errorf("result dispatch failed: %w", err)
	}

	// Record the dispatch event in the ledger
	if err := l.AiMLMLedger.RecordResultDispatch(transactionID, recipient, string(encryptedResult)); err != nil {
		log.Printf("AIResultDispatch - Ledger Recording Failed: %v", err)
		return fmt.Errorf("failed to record dispatch in ledger: %w", err)
	}

	// Process the transaction for validation
	if err := common.ProcessSingleTransaction(sc, transactionID, encryptedResult); err != nil {
		log.Printf("AIResultDispatch - Transaction Processing Failed: %v", err)
		return fmt.Errorf("transaction processing failed: %w", err)
	}

	log.Printf("AIResultDispatch - Success: Result Dispatched and Recorded")
	return nil
}


// AIPredictionGenerate generates a predictive result based on model inference.
func AIPredictionGenerate(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID, transactionID string, inputData []byte) (string, error) {
	logAction("AIPredictionGenerate - Start", fmt.Sprintf("ModelID: %s, NodeID: %s, TransactionID: %s", modelID, nodeID, transactionID))

	// Validate inputs
	if modelID == "" || nodeID == "" || transactionID == "" {
		return "", errors.New("modelID, nodeID, and transactionID cannot be empty")
	}
	if len(inputData) == 0 {
		return "", errors.New("input data cannot be empty")
	}

	// Fetch the AI model
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIPredictionGenerate - Model Retrieval Failed: %v", err)
		return "", fmt.Errorf("model retrieval failed: %w", err)
	}

	// Generate the prediction
	prediction, err := model.GeneratePrediction(inputData, nodeID)
	if err != nil {
		log.Printf("AIPredictionGenerate - Prediction Generation Failed: %v", err)
		return "", fmt.Errorf("prediction generation failed: %w", err)
	}

	// Encrypt the prediction result
	encryptedPrediction, err := AIModelEncrypt(prediction, modelID)
	if err != nil {
		log.Printf("AIPredictionGenerate - Encryption Failed: %v", err)
		return "", fmt.Errorf("encryption failed for prediction result: %w", err)
	}

	// Record the prediction in the ledger
	if err := l.AiMLMLedger.RecordPrediction(transactionID, modelID); err != nil {
		log.Printf("AIPredictionGenerate - Ledger Recording Failed: %v", err)
		return "", fmt.Errorf("failed to record prediction in ledger: %w", err)
	}

	// Process the transaction for validation
	if err := common.ProcessSingleTransaction(sc, transactionID, encryptedPrediction); err != nil {
		log.Printf("AIPredictionGenerate - Transaction Processing Failed: %v", err)
		return "", fmt.Errorf("transaction processing failed: %w", err)
	}

	log.Printf("AIPredictionGenerate - Success: Prediction Generated and Recorded")
	return string(encryptedPrediction), nil
}




// AIRecommendationGenerate generates recommendations based on analysis results and model criteria.
func AIRecommendationGenerate(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID string, criteria []byte, transactionID string) (string, error) {
	logAction("AIRecommendationGenerate - Start", fmt.Sprintf("ModelID: %s, TransactionID: %s", modelID, transactionID))

	// Validate inputs
	if modelID == "" || transactionID == "" {
		return "", errors.New("modelID and transactionID cannot be empty")
	}
	if len(criteria) == 0 {
		return "", errors.New("criteria cannot be empty")
	}

	// Fetch the AI model
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIRecommendationGenerate - Model Retrieval Failed: %v", err)
		return "", fmt.Errorf("model retrieval failed: %w", err)
	}

	// Generate the recommendation
	recommendation, err := model.GenerateRecommendation(criteria)
	if err != nil {
		log.Printf("AIRecommendationGenerate - Recommendation Generation Failed: %v", err)
		return "", fmt.Errorf("recommendation generation failed: %w", err)
	}

	// Encrypt the recommendation result
	encryptedRecommendation, err := AIModelEncrypt(recommendation, modelID)
	if err != nil {
		log.Printf("AIRecommendationGenerate - Encryption Failed: %v", err)
		return "", fmt.Errorf("encryption failed for recommendation result: %w", err)
	}

	// Record the recommendation in the ledger
	if err := l.AiMLMLedger.RecordRecommendation(transactionID, string(encryptedRecommendation)); err != nil {
		log.Printf("AIRecommendationGenerate - Ledger Recording Failed: %v", err)
		return "", fmt.Errorf("failed to record recommendation in ledger: %w", err)
	}

	log.Printf("AIRecommendationGenerate - Success: Recommendation Generated and Recorded")
	return string(encryptedRecommendation), nil
}



// AIRecommendationUpdate updates existing recommendations with new data or insights.
func AIRecommendationUpdate(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, nodeID string, newData []byte, transactionID string) ([]byte, error) {
	logAction("AIRecommendationUpdate - Start", fmt.Sprintf("ModelID: %s, NodeID: %s, TransactionID: %s", modelID, nodeID, transactionID))

	// Validate inputs
	if modelID == "" || nodeID == "" || transactionID == "" {
		return nil, errors.New("modelID, nodeID, and transactionID cannot be empty")
	}
	if len(newData) == 0 {
		return nil, errors.New("new data cannot be empty")
	}

	// Fetch the AI model
	model, err := fetchModel(modelID)
	if err != nil {
		log.Printf("AIRecommendationUpdate - Model Retrieval Failed: %v", err)
		return nil, fmt.Errorf("model retrieval failed: %w", err)
	}

	// Update the recommendation
	updatedRecommendation, err := model.UpdateRecommendation(nodeID, newData)
	if err != nil {
		log.Printf("AIRecommendationUpdate - Update Failed: %v", err)
		return nil, fmt.Errorf("recommendation update failed: %w", err)
	}

	// Encrypt the updated recommendation
	encryptedUpdate, err := AIModelEncrypt(updatedRecommendation, modelID)
	if err != nil {
		log.Printf("AIRecommendationUpdate - Encryption Failed: %v", err)
		return nil, fmt.Errorf("encryption failed for recommendation update: %w", err)
	}

	// Record the updated recommendation in the ledger
	if err := l.AiMLMLedger.RecordRecommendationUpdate(transactionID, string(encryptedUpdate)); err != nil {
		log.Printf("AIRecommendationUpdate - Ledger Recording Failed: %v", err)
		return nil, fmt.Errorf("failed to record recommendation update in ledger: %w", err)
	}

	log.Printf("AIRecommendationUpdate - Success: Recommendation Updated and Recorded")
	return []byte(encryptedUpdate), nil
}


// fetchModel retrieves an AI/ML model based on its ModelID by accessing the model's index in the ledger.
func fetchModel(modelID string) (Model, error) {
	logAction("fetchModel - Start", fmt.Sprintf("ModelID: %s", modelID))

	// Retrieve the model index entry from the ledger
	modelIndex, err := GetModelIndex(modelID)
	if err != nil {
		log.Printf("fetchModel - Model Index Retrieval Failed: %v", err)
		return Model{}, fmt.Errorf("model index retrieval failed: %w", err)
	}

	// Retrieve model data from IPFS/Swarm using the link in the index
	modelData, err := fetchModelFromStorage(modelIndex.IPFSLink)
	if err != nil {
		log.Printf("fetchModel - Model Retrieval from Storage Failed: %v", err)
		return Model{}, fmt.Errorf("model retrieval from storage failed: %w", err)
	}

	// Deserialize model data into a Model struct
	model, err := deserializeModelData(modelData)
	if err != nil {
		log.Printf("fetchModel - Model Deserialization Failed: %v", err)
		return Model{}, fmt.Errorf("model deserialization failed: %w", err)
	}

	log.Printf("fetchModel - Success: Model Retrieved and Deserialized for ModelID: %s", modelID)
	return model, nil
}


// GetModelIndex retrieves the model's index record from the ledger.
func GetModelIndex(modelID string) (ledger.ModelIndex, error) {
	logAction("GetModelIndex - Start", fmt.Sprintf("ModelID: %s", modelID))

	modelIndex, err := ledgerInstance.AiMLMLedger.FetchModelIndex(modelID)
	if err != nil {
		log.Printf("GetModelIndex - Retrieval Failed: %v", err)
		return ledger.ModelIndex{}, fmt.Errorf("model index retrieval failed: %w", err)
	}

	log.Printf("GetModelIndex - Success: Model Index Retrieved for ModelID: %s", modelID)
	return modelIndex, nil
}


// fetchModelFromStorage retrieves the model data from IPFS or Swarm using the provided link.
func fetchModelFromStorage(ipfsLink string) ([]byte, error) {
	logAction("fetchModelFromStorage - Start", fmt.Sprintf("IPFSLink: %s", ipfsLink))

	// HTTP client with timeout
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// HTTP GET request
	resp, err := client.Get(ipfsLink)
	if err != nil {
		log.Printf("fetchModelFromStorage - HTTP Request Failed: %v", err)
		return nil, fmt.Errorf("failed to fetch model data from storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("fetchModelFromStorage - Unexpected HTTP Status: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("fetchModelFromStorage - Error Reading Response Body: %v", err)
		return nil, fmt.Errorf("error reading model data: %w", err)
	}

	log.Printf("fetchModelFromStorage - Success: Model Data Retrieved from IPFSLink: %s", ipfsLink)
	return data, nil
}


// deserializeModelData converts raw JSON data into a Model struct.
func deserializeModelData(data []byte) (Model, error) {
	logAction("deserializeModelData - Start", "Deserializing Model Data")

	var model Model
	err := json.Unmarshal(data, &model)
	if err != nil {
		log.Printf("deserializeModelData - Deserialization Failed: %v", err)
		return Model{}, fmt.Errorf("deserialization error: %w", err)
	}

	log.Printf("deserializeModelData - Success: Model Data Deserialized")
	return model, nil
}



// sendEncryptedResult securely dispatches encrypted results to the specified recipient.
func sendEncryptedResult(encryptedData string, recipientURL string) error {
	logAction("sendEncryptedResult - Start", fmt.Sprintf("RecipientURL: %s", recipientURL))

	// Validate inputs
	if encryptedData == "" {
		return errors.New("encrypted data cannot be empty")
	}
	if recipientURL == "" {
		return errors.New("recipient URL cannot be empty")
	}

	// HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Prepare payload and HTTP POST request
	payload := []byte(encryptedData)
	req, err := http.NewRequest("POST", recipientURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("sendEncryptedResult - Request Creation Failed: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set secure headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer some-secure-token") // Replace with a valid token
	req.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano()))

	// Execute HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("sendEncryptedResult - HTTP Request Failed: %v", err)
		return fmt.Errorf("failed to send encrypted data: %w", err)
	}
	defer resp.Body.Close()

	// Validate response
	if resp.StatusCode != http.StatusOK {
		log.Printf("sendEncryptedResult - Unexpected HTTP Status: %d", resp.StatusCode)
		return fmt.Errorf("failed to dispatch encrypted result: received HTTP status %d", resp.StatusCode)
	}

	log.Printf("sendEncryptedResult - Success: Encrypted Result Dispatched to Recipient")
	return nil
}

// logAction logs the entry, parameters, and exit of a specific action with timestamps and context.
func logAction(actionName string, context string) {
	timestamp := time.Now().Format(time.RFC3339) // Standardized timestamp format
	log.Printf("[ACTION] %s | Timestamp: %s | Context: %s", actionName, timestamp, context)
}

