package ledger

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// RecordInference logs an inference request.
func (l *AiMLMLedger) RecordInference(modelID, nodeID, result string) error {
	// Validate inputs
	if modelID == "" || nodeID == "" || result == "" {
		return fmt.Errorf("invalid input: modelID, nodeID, and result must be non-empty")
	}

	l.Lock()
	defer l.Unlock()

	// Generate a unique ID for the inference record
	inferenceID := generateUniqueID()

	// Log the inference
	l.AiMLMLedgerState.Inferences[inferenceID] = InferenceRecord{
		ModelID:    modelID,
		NodeID:     nodeID,
		Timestamp:  time.Now(),
		Result:     result,
		Processed:  true,
	}

	log.Printf("[INFO] Inference logged: ID=%s, ModelID=%s, NodeID=%s, Result=%s", inferenceID, modelID, nodeID, result)
	return nil
}


// RecordAnalysisStart starts an AI analysis.
func (l *AiMLMLedger) RecordAnalysisStart(modelID, nodeID string) (string, error) {
	// Validate inputs
	if modelID == "" || nodeID == "" {
		return "", fmt.Errorf("invalid input: modelID and nodeID must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Generate a unique ID for the analysis
	analysisID := generateUniqueID()

	// Record the analysis start
	l.AiMLMLedgerState.ActiveAnalyses[analysisID] = AnalysisRecord{
		AnalysisID: analysisID,
		ModelID:    modelID,
		NodeID:     nodeID,
		StartTime:  time.Now(),
		Status:     "Active",
	}

	log.Printf("[INFO] Analysis started: ID=%s, ModelID=%s, NodeID=%s", analysisID, modelID, nodeID)
	return analysisID, nil
}


// RecordAnalysisStop stops an ongoing analysis.
func (l *AiMLMLedger) RecordAnalysisStop(analysisID string) error {
	// Validate input
	if analysisID == "" {
		return fmt.Errorf("invalid input: analysisID must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Retrieve the analysis record
	analysis, exists := l.AiMLMLedgerState.ActiveAnalyses[analysisID]
	if !exists {
		return fmt.Errorf("analysis not found for ID: %s", analysisID)
	}

	// Update the analysis status
	analysis.Status = "Completed"
	analysis.StopTime = time.Now()
	l.AiMLMLedgerState.ActiveAnalyses[analysisID] = analysis

	log.Printf("[INFO] Analysis stopped: ID=%s, ModelID=%s, NodeID=%s", analysisID, analysis.ModelID, analysis.NodeID)
	return nil
}


// RecordResultDispatch logs the dispatch of an AI inference result.
func (l *AiMLMLedger) RecordResultDispatch(modelID, nodeID, result string) error {
	// Reuse the RecordInference method to log the result dispatch
	err := l.RecordInference(modelID, nodeID, result)
	if err != nil {
		return fmt.Errorf("failed to dispatch result: %w", err)
	}

	log.Printf("[INFO] Result dispatched: ModelID=%s, NodeID=%s, Result=%s", modelID, nodeID, result)
	return nil
}


// RecordPrediction logs a prediction for a model.
func (l *AiMLMLedger) RecordPrediction(modelID, prediction string) error {
	// Validate input
	if modelID == "" || prediction == "" {
		return fmt.Errorf("invalid input: modelID and prediction must be non-empty")
	}

	// Use RecordInference to log the prediction
	err := l.RecordInference(modelID, "", prediction)
	if err != nil {
		return fmt.Errorf("failed to record prediction: %w", err)
	}

	log.Printf("[INFO] Prediction recorded for ModelID=%s, Prediction=%s", modelID, prediction)
	return nil
}


// RecordRecommendation logs a recommendation by a model.
func (l *AiMLMLedger) RecordRecommendation(modelID, content string) error {
	// Validate input
	if modelID == "" || content == "" {
		return fmt.Errorf("invalid input: modelID and content must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Generate a unique ID for the recommendation
	recommendationID := generateUniqueID()

	// Log the recommendation
	l.AiMLMLedgerState.Recommendations[recommendationID] = Recommendation{
		ModelID:    modelID,
		Timestamp:  time.Now(),
		Content:    content,
		Updated:    false,
	}

	log.Printf("[INFO] Recommendation recorded: ID=%s, ModelID=%s, Content=%s", recommendationID, modelID, content)
	return nil
}


// RecordRecommendationUpdate updates a recommendation.
func (l *AiMLMLedger) RecordRecommendationUpdate(recommendationID, content string) error {
	// Validate input
	if recommendationID == "" || content == "" {
		return fmt.Errorf("invalid input: recommendationID and content must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Retrieve the existing recommendation
	recommendation, exists := l.AiMLMLedgerState.Recommendations[recommendationID]
	if !exists {
		return fmt.Errorf("recommendation not found for ID: %s", recommendationID)
	}

	// Update the recommendation
	recommendation.Content = content
	recommendation.Updated = true
	recommendation.Timestamp = time.Now() // Update the timestamp
	l.AiMLMLedgerState.Recommendations[recommendationID] = recommendation

	log.Printf("[INFO] Recommendation updated: ID=%s, NewContent=%s", recommendationID, content)
	return nil
}


// GetModelIndex retrieves a model's index.
func (l *AiMLMLedger) GetModelIndex(modelID string) (ModelIndex, error) {
	// Validate input
	if modelID == "" {
		return ModelIndex{}, fmt.Errorf("invalid input: modelID must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Retrieve the model index
	modelIndex, exists := l.AiMLMLedgerState.ModelIndex[modelID]
	if !exists {
		return ModelIndex{}, fmt.Errorf("model not found in index for ID: %s", modelID)
	}

	log.Printf("[INFO] Model index retrieved for ModelID=%s", modelID)
	return modelIndex, nil
}


// UpdateModelIndex updates a model's entry in the index.
func (l *AiMLMLedger) UpdateModelIndex(modelIndex ModelIndex) error {
	// Validate input
	if modelIndex.ModelID == "" {
		return fmt.Errorf("invalid input: modelID must be non-empty")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Update the model index
	l.AiMLMLedgerState.ModelIndex[modelIndex.ModelID] = modelIndex

	log.Printf("[INFO] Model index updated: ModelID=%s", modelIndex.ModelID)
	return nil
}



// RecordScaling logs a scaling action for a model in the ledger.
func (l *AiMLMLedger) RecordScaling(transactionID, modelID, direction, encryptedLog string) error {
	// Add a new scaling log entry
	l.ScalingLogs = append(l.ScalingLogs, ScalingLog{
		TransactionID: transactionID,
		ModelID:       modelID,
		Direction:     direction,
		EncryptedLog:  encryptedLog,
		Timestamp:     time.Now(),
	})
	return nil
}



// HasAccess checks if the model grants access based on the provided access key.
func (m *ModelIndex) HasAccess(accessKey string) bool {
	// Check if the accessKey is present in the Permissions map
	_, allowed := m.Permissions[accessKey]
	return allowed
}

// HasAccess is a standalone function that checks if the model grants access based on the provided access key.
func HasAccess(model ModelIndex, accessKey string) bool {
	// Check if the accessKey is present in the Permissions map
	_, allowed := model.Permissions[accessKey]
	return allowed
}

// RecordAiNetworkMaintenance logs maintenance activity.
func (l *AiMLMLedger) RecordAiNetworkMaintenance(action string) error {
	l.TrafficRecords = append(l.TrafficRecords, TrafficRecord{
		ModelID:   "NetworkMaintenance",
		Action:    action,
		Timestamp: time.Now(),
	})
	return nil
}

// GetAllModels retrieves all model indices.
func (l *AiMLMLedger) GetAllModels() []ModelIndex {
	models := make([]ModelIndex, 0, len(l.AiMLMLedgerState.ModelIndex))
	for _, model := range l.AiMLMLedgerState.ModelIndex {
		models = append(models, model)
	}
	return models
}

// RecordProcessing logs data processing for a model.
func (l *AiMLMLedger) RecordProcessing(modelID, nodeID string) error {
	l.AiMLMLedgerState.DataProcessingLogs[generateUniqueID()] = DataProcessingLog{
		ProcessID: generateUniqueID(),
		ModelID:   modelID,
		NodeID:    nodeID,
		StartTime: time.Now(),
		Status:    "Processing",
	}
	return nil
}

// RecordAiService logs an AI service entry.
func (l *AiMLMLedger) RecordAiService(serviceID string, metrics ServiceMetrics) error {
	l.AiMLMLedgerState.Services[serviceID] = AiService{
		ServiceID:   serviceID,
		Status:      "Active",
		Metrics:     metrics,
		LastUpdated: time.Now(),
	}
	return nil
}

// RemoveAiService removes an AI service.
func (l *AiMLMLedger) RemoveAiService(serviceID string) error {
	if _, exists := l.AiMLMLedgerState.Services[serviceID]; !exists {
		return errors.New("service not found")
	}
	delete(l.AiMLMLedgerState.Services, serviceID)
	return nil
}

// GetServiceStatusAndMetrics retrieves service metrics.
func (l *AiMLMLedger) GetServiceStatusAndMetrics(serviceID string) (ServiceMetrics, error) {
	service, exists := l.AiMLMLedgerState.Services[serviceID]
	if !exists {
		return ServiceMetrics{}, errors.New("service not found")
	}
	return service.Metrics, nil
}

// LogServiceMetrics logs metrics for an AI service.
func (l *AiMLMLedger) LogServiceMetrics(serviceID string, metrics ServiceMetrics) error {
	service, exists := l.AiMLMLedgerState.Services[serviceID]
	if !exists {
		return errors.New("service not found")
	}
	service.Metrics = metrics
	service.LastUpdated = time.Now()
	l.AiMLMLedgerState.Services[serviceID] = service
	return nil
}

// GetServiceStatus retrieves an AI service status.
func (l *AiMLMLedger) GetServiceStatus(serviceID string) (string, error) {
	service, exists := l.AiMLMLedgerState.Services[serviceID]
	if !exists {
		return "", errors.New("service not found")
	}
	return service.Status, nil
}

// GetOptimalNodeForModel finds the node with the optimal load for a model.
func (l *AiMLMLedger) GetOptimalNodeForModel(modelID string) (string, error) {
	minLoad := int(^uint(0) >> 1) // Max int value
	var optimalNode string

	for _, model := range l.AiMLMLedgerState.ModelIndex {
		if model.ModelID == modelID && model.Load < minLoad {
			minLoad = model.Load
			optimalNode = model.NodeLocation
		}
	}
	if optimalNode == "" {
		return "", errors.New("no optimal node found")
	}
	return optimalNode, nil
}

// GetModelLoad retrieves the current load for a model.
func (l *AiMLMLedger) GetModelLoad(modelID string) (int, error) {
	model, exists := l.AiMLMLedgerState.ModelIndex[modelID]
	if !exists {
		return 0, errors.New("model not found")
	}
	return model.Load, nil
}

// RedistributeModelTraffic redistributes traffic for load balancing.
func (l *AiMLMLedger) RedistributeModelTraffic(modelID string, targetNode string) error {
	record := TrafficRecord{
		ModelID:   modelID,
		NodeID:    targetNode,
		Action:    "Redistribute",
		Timestamp: time.Now(),
	}
	l.TrafficRecords = append(l.TrafficRecords, record)
	return nil
}

// LogTrafficBalanceOperation logs a traffic balance operation.
func (l *AiMLMLedger) LogTrafficBalanceOperation(modelID string, nodeID string) error {
	l.TrafficRecords = append(l.TrafficRecords, TrafficRecord{
		ModelID:   modelID,
		NodeID:    nodeID,
		Action:    "Balance",
		Timestamp: time.Now(),
	})
	return nil
}

// RecordContainerDeployment records container deployments in the ledger.
func (l *AiMLMLedger) RecordContainerDeployment(containerID, modelID, nodeID string) error {
	l.AiMLMLedgerState.Containers[containerID] = Container{
		ContainerID: containerID,
		ModelID:     modelID,
		NodeID:      nodeID,
		Status:      "Deployed",
	}
	return nil
}

// UpdateContainerStatus updates the status of a container in the ledger.
func (l *AiMLMLedger) UpdateContainerStatus(containerID, status string) error {
	container, exists := l.AiMLMLedgerState.Containers[containerID]
	if !exists {
		return errors.New("container not found")
	}
	container.Status = status
	l.AiMLMLedgerState.Containers[containerID] = container
	return nil
}

// RecordDataProcess logs data processing operations for a model.
func (l *AiMLMLedger) RecordDataProcess(modelID, nodeID string) error {
	logID := generateUniqueID()
	l.AiMLMLedgerState.DataProcessingLogs[logID] = DataProcessingLog{
		ProcessID: logID,
		ModelID:   modelID,
		NodeID:    nodeID,
		StartTime: time.Now(),
		Status:    "Processing",
	}
	return nil
}

// GetModelData retrieves data associated with a model.
func (l *AiMLMLedger) GetModelData(modelID string) ([]DataProcessingLog, error) {
	var logs []DataProcessingLog
	for _, log := range l.AiMLMLedgerState.DataProcessingLogs {
		if log.ModelID == modelID {
			logs = append(logs, log)
		}
	}
	if len(logs) == 0 {
		return nil, errors.New("no data found for model")
	}
	return logs, nil
}


// RecordModelRestriction records restrictions on a model.
func (l *AiMLMLedger) RecordModelRestriction(modelID, reason string) error {
	l.AiMLMLedgerState.ModelRestrictions[modelID] = ModelRestriction{
		ModelID:    modelID,
		Restricted: true,
		Reason:     reason,
		Timestamp:  time.Now(),
	}
	return nil
}

// RecordModelPermissions logs permissions for a model.
func (l *AiMLMLedger) RecordModelPermissions(modelID string, users []string) error {
	l.AiMLMLedgerState.ModelPermissions[modelID] = ModelPermissions{
		ModelID:      modelID,
		AllowedUsers: users,
	}
	return nil
}

// RecordAccessToken logs an access token for a model.
func (l *AiMLMLedger) RecordAccessToken(tokenID, modelID, grantedTo, permissions string, expiry time.Time) error {
	l.AccessTokens[tokenID] = AccessToken{
		TokenID:      tokenID,
		ModelID:      modelID,
		GrantedTo:    grantedTo,
		Permissions:  permissions,
		Expiry:       expiry,
	}
	return nil
}

// GetModelAccessList retrieves the access list for a model.
func (l *AiMLMLedger) GetModelAccessList(modelID string) ([]string, error) {
	accessList, exists := l.AiMLMLedgerState.ModelAccessList[modelID]
	if !exists {
		return nil, errors.New("model access list not found")
	}
	return accessList, nil
}

// UpdateModelAccessList updates the access list for a model.
func (l *AiMLMLedger) UpdateModelAccessList(modelID string, users []string) error {
	l.AiMLMLedgerState.ModelAccessList[modelID] = users
	return nil
}

// RecordModelCheckpoint logs a model checkpoint.
func (l *AiMLMLedger) RecordModelCheckpoint(modelID string, version int, dataHash string) error {
	l.Checkpoints[modelID] = ModelCheckpoint{
		ModelID:    modelID,
		Version:    version,
		CreatedAt:  time.Now(),
		DataHash:   dataHash,
	}
	return nil
}

// GetModelCheckpoint retrieves a model checkpoint.
func (l *AiMLMLedger) GetModelCheckpoint(modelID string) (ModelCheckpoint, error) {
	checkpoint, exists := l.Checkpoints[modelID]
	if !exists {
		return ModelCheckpoint{}, errors.New("checkpoint not found")
	}
	return checkpoint, nil
}

// RecordModelAccessLog logs access to a model.
func (l *AiMLMLedger) RecordModelAccessLog(modelID, userID, action string) error {
	log := AccessLog{
		ModelID:   modelID,
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
	}
	l.AiMLMLedgerState.ModelAccessLogs[modelID] = append(l.AiMLMLedgerState.ModelAccessLogs[modelID], log)
	return nil
}

// RecordModelUsageStatistics updates usage statistics for a model.
func (l *AiMLMLedger) RecordModelUsageStatistics(modelID string, duration time.Duration) error {
	usageStats := l.UsageStatistics[modelID]
	usageStats.UsageCount++
	usageStats.LastUsedAt = time.Now()
	usageStats.UsageDuration += duration
	l.UsageStatistics[modelID] = usageStats
	return nil
}

// RecordModelPerformanceMetrics logs performance metrics.
func (l *AiMLMLedger) RecordModelPerformanceMetrics(modelID string, accuracy, loss float64) error {
	l.PerformanceMetrics[modelID] = PerformanceMetrics{
		ModelID:    modelID,
		Accuracy:   accuracy,
		Loss:       loss,
		LastUpdated: time.Now(),
	}
	return nil
}

// RecordModelHyperparameterTuning logs hyperparameter tuning results.
func (l *AiMLMLedger) RecordModelHyperparameterTuning(modelID string, params map[string]interface{}) error {
	// Log or store tuning parameters, depending on system requirements.
	return nil
}

// RecordModelComplianceCheck logs a compliance check for a model.
func (l *AiMLMLedger) RecordModelComplianceCheck(modelID, status, details string) error {
	l.ComplianceChecks = append(l.ComplianceChecks, ComplianceCheck{
		ModelID:    modelID,
		Status:     status,
		Details:    details,
		Timestamp:  time.Now(),
	})
	return nil
}

// RecordModelSecurityAudit logs a security audit for a model.
func (l *AiMLMLedger) RecordModelSecurityAudit(modelID string, passed bool, findings string) error {
	l.SecurityAudits = append(l.SecurityAudits, SecurityAudit{
		ModelID:   modelID,
		Passed:    passed,
		Timestamp: time.Now(),
		Findings:  findings,
	})
	return nil
}

// RecordResourceAllocation logs a resource allocation.
func (l *AiMLMLedger) RecordResourceAllocation(resourceID, modelID string, amount float64) error {
	l.ResourceAllocations[resourceID] = ResourceAllocation{
		ResourceID:  resourceID,
		ModelID:     modelID,
		Amount:      amount,
		AllocatedAt: time.Now(),
	}
	return nil
}

// RecordStorageAllocation logs a storage allocation.
func (l *AiMLMLedger) RecordStorageAllocation(storageID, modelID string, sizeMB int) error {
	l.StorageAllocations[storageID] = AiModelStorageAllocation{
		StorageID:   storageID,
		ModelID:     modelID,
		SizeMB:      sizeMB,
		AllocatedAt: time.Now(),
	}
	return nil
}

// RecordResourceRelease releases allocated resources.
func (l *AiMLMLedger) RecordResourceRelease(resourceID string) error {
	delete(l.ResourceAllocations, resourceID)
	return nil
}

// RecordCacheData caches data for a model.
func (l *AiMLMLedger) RecordCacheData(modelID, dataID, data string) error {
	l.CacheRecords[dataID] = CacheData{
		ModelID:    modelID,
		DataID:     dataID,
		Data:       data,
		CreatedAt:  time.Now(),
	}
	return nil
}

// RecordCacheClear clears cached data for a model.
func (l *AiMLMLedger) RecordCacheClear(dataID string) error {
	delete(l.CacheRecords, dataID)
	return nil
}

// RecordModelShutdown logs a model shutdown.
func (l *AiMLMLedger) RecordModelShutdown(modelID string) error {
	return l.UpdateModelStatus(modelID, "Shutdown")
}

// RecordModelRestart logs a model restart.
func (l *AiMLMLedger) RecordModelRestart(modelID string) error {
	return l.UpdateModelStatus(modelID, "Restarted")
}

// GetModelStatus retrieves the status of a model.
func (l *AiMLMLedger) GetModelStatus(modelID string) (string, error) {
	model, exists := l.Models[modelID]
	if !exists {
		return "", errors.New("model not found")
	}
	return model.Status, nil
}

// GetDeploymentStatus retrieves the deployment status of a model.
func (l *AiMLMLedger) GetDeploymentStatus(deploymentID string) (string, error) {
	deployment, exists := l.AiMLMLedgerState.DeploymentChecks[deploymentID]
	if !exists {
		return "", errors.New("deployment not found")
	}
	return deployment.Status, nil
}

// RecordDeploymentCheck logs a deployment check.
func (l *AiMLMLedger) RecordDeploymentCheck(deploymentID, status, details string) error {
	l.AiMLMLedgerState.DeploymentChecks[deploymentID] = DeploymentCheck{
		DeploymentID: deploymentID,
		Status:       status,
		Details:      details,
		Timestamp:    time.Now(),
	}
	return nil
}

// GetTrainingStatus retrieves training status for a model.
func (l *AiMLMLedger) GetTrainingStatus(modelID string) (string, error) {
	status, exists := l.TrainingStatus[modelID]
	if !exists {
		return "", errors.New("training status not found")
	}
	return status.Status, nil
}

// GetRunStatus retrieves the running status for a model.
func (l *AiMLMLedger) GetRunStatus(modelID string) (string, error) {
	status, exists := l.RunStatus[modelID]
	if !exists {
		return "", errors.New("run status not found")
	}
	return status.Status, nil
}

// RecordRunCheck records a run check for a model.
func (l *AiMLMLedger) RecordRunCheck(modelID, status string) error {
	l.RunStatus[modelID] = RunStatus{
		ModelID:     modelID,
		Status:      status,
		LastChecked: time.Now(),
	}
	return nil
}

// OnChainTrainModel initiates on-chain training for a model.
func (l *AiMLMLedger) OnChainTrainModel(modelID string) error {
	l.TrainingStatus[modelID] = TrainingStatus{
		ModelID:     modelID,
		Status:      "On-Chain Training Started",
		LastUpdated: time.Now(),
	}
	return nil
}

// RecordModelTraining logs a training event for a model.
func (l *AiMLMLedger) RecordModelTraining(modelID, trainingType string) error {
	l.TrainingStatus[modelID] = TrainingStatus{
		ModelID:     modelID,
		Status:      trainingType + " Training Started",
		LastUpdated: time.Now(),
	}
	return nil
}

// RecordOffChainTraining initiates off-chain training for a model.
func (l *AiMLMLedger) RecordOffChainTraining(modelID string) error {
	l.TrainingStatus[modelID] = TrainingStatus{
		ModelID:     modelID,
		Status:      "Off-Chain Training Started",
		LastUpdated: time.Now(),
	}
	return nil
}

// RecordModelRetrain logs a retraining event.
func (l *AiMLMLedger) RecordModelRetrain(modelID string) error {
	l.TrainingStatus[modelID] = TrainingStatus{
		ModelID:     modelID,
		Status:      "Retraining",
		LastUpdated: time.Now(),
	}
	return nil
}

// RecordEncryption logs encryption activity in the ledger
func (l *AiMLMLedger) RecordEncryption(transactionID, encryptedData string) error {
	l.EncryptionLogs[transactionID] = EncryptionLog{
		TransactionID: transactionID,
		EncryptedData: encryptedData,
		Timestamp:     time.Now(),
	}
	return nil
}

// RecordDecryption logs decryption activity in the ledger
func (l *AiMLMLedger) RecordDecryption(transactionID, decryptedData string) error {
	l.DecryptionLogs[transactionID] = DecryptionLog{
		TransactionID: transactionID,
		DecryptedData: decryptedData,
		Timestamp:     time.Now(),
	}
	return nil
}

// FetchModelIndex retrieves a model's index entry from the ledger based on the model ID.
func (l *AiMLMLedger) FetchModelIndex(modelID string) (ModelIndex, error) {
	model, exists := l.ModelIndexMap[modelID]
	if !exists {
		return ModelIndex{}, errors.New("model index not found")
	}
	return model, nil
}

// AddModelIndex adds a new model index entry to the ledger.
func (l *AiMLMLedger) AddModelIndex(modelIndex ModelIndex) {
	l.ModelIndexMap[modelIndex.ModelID] = modelIndex
}

// AddAccessControl stores access control rules for a model in the ledger.
func (l *AiMLMLedger) AddModelAccessControl(modelID string, accessControlJSON string) error {
	// Ensure the model ID exists in the model index
	if _, exists := l.ModelIndexMap[modelID]; !exists {
		return errors.New("model ID not found in the ledger")
	}

	// Add or update access control rules for the specified model ID
	l.AccessControls[modelID] = accessControlJSON
	return nil
}

// GetServiceBalance retrieves the balance of a specific AI service.
func (l *AiMLMLedger) GetServiceBalance(serviceID string) (float64, error) {
	service, exists := l.State.Services[serviceID]
	if !exists {
		return 0, errors.New("service not found in the ledger")
	}
	return service.Balance, nil
}

// RecordModelDeployment logs a model deployment action in the ledger.
func (l *AiMLMLedger) RecordModelDeployment(transactionID, modelID, description string) error {
	record := ModelActionRecord{
		TransactionID: transactionID,
		ModelID:       modelID,
		Action:        "Deployment",
		Timestamp:     time.Now(),
		Description:   description,
	}
	l.ModelActions = append(l.ModelActions, record)
	return nil
}

// RecordModelUndeployment logs a model undeployment action in the ledger.
func (l *AiMLMLedger) RecordModelUndeployment(transactionID, modelID, description string) error {
	record := ModelActionRecord{
		TransactionID: transactionID,
		ModelID:       modelID,
		Action:        "Undeployment",
		Timestamp:     time.Now(),
		Description:   description,
	}
	l.ModelActions = append(l.ModelActions, record)
	return nil
}

// RecordModelUpdate logs a model update action in the ledger.
func (l *AiMLMLedger) RecordModelUpdate(transactionID, modelID, description string) error {
	record := ModelActionRecord{
		TransactionID: transactionID,
		ModelID:       modelID,
		Action:        "Update",
		Timestamp:     time.Now(),
		Description:   description,
	}
	l.ModelActions = append(l.ModelActions, record)
	return nil
}

