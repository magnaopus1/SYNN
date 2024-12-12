package ai_ml_operation

import (
	"errors"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// ModelScalingUp scales the computational resources of the specified model to accommodate increased demand.
func ModelScalingUp(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID string) error {
	logAction("ModelScalingUp", fmt.Sprintf("Initiating scale-up for model: %s, Transaction: %s", modelID, transactionID))

	// Step 1: Fetch model metadata
	model, err := fetchModel(modelID)
	if err != nil {
		logAction("ModelScalingUpError", fmt.Sprintf("Failed to retrieve model: %s", modelID))
		return fmt.Errorf("failed to retrieve model for scaling up: %w", err)
	}

	// Step 2: Validate scaling possibility
	if !model.CanScaleUp() {
		logAction("ModelScalingUpError", fmt.Sprintf("Model cannot be scaled up further: %s", modelID))
		return errors.New("model cannot be scaled up further")
	}

	// Step 3: Perform scaling operation
	if err := model.ScaleUp(); err != nil {
		logAction("ModelScalingUpError", fmt.Sprintf("Scaling up failed for model: %s", modelID))
		return fmt.Errorf("scaling up operation failed: %w", err)
	}

	// Step 4: Create encrypted scaling log
	encryptedLog, err := createEncryptedLog(modelID, "scale_up")
	if err != nil {
		logAction("ModelScalingUpError", fmt.Sprintf("Failed to encrypt scaling log for model: %s", modelID))
		return fmt.Errorf("failed to encrypt scaling log: %w", err)
	}

	// Step 5: Record scaling operation in the ledger
	if err := l.AiMLMLedger.RecordScaling(transactionID, modelID, "up", encryptedLog); err != nil {
		logAction("ModelScalingUpError", fmt.Sprintf("Failed to record scaling in ledger for model: %s", modelID))
		return fmt.Errorf("failed to record scaling in ledger: %w", err)
	}

	logAction("ModelScalingUpCompleted", fmt.Sprintf("Scale-up completed for model: %s", modelID))
	return nil
}


// ModelScalingDown scales down the computational resources of the specified model when demand decreases.
func ModelScalingDown(l *ledger.Ledger, sc *common.SynnergyConsensus, modelID, transactionID string) error {
	logAction("ModelScalingDown", fmt.Sprintf("Initiating scale-down for model: %s, Transaction: %s", modelID, transactionID))

	// Step 1: Fetch model metadata
	model, err := fetchModel(modelID)
	if err != nil {
		logAction("ModelScalingDownError", fmt.Sprintf("Failed to retrieve model: %s", modelID))
		return fmt.Errorf("failed to retrieve model for scaling down: %w", err)
	}

	// Step 2: Validate scaling possibility
	if !model.CanScaleDown() {
		logAction("ModelScalingDownError", fmt.Sprintf("Model cannot be scaled down further: %s", modelID))
		return errors.New("model cannot be scaled down further")
	}

	// Step 3: Perform scaling operation
	if err := model.ScaleDown(); err != nil {
		logAction("ModelScalingDownError", fmt.Sprintf("Scaling down failed for model: %s", modelID))
		return fmt.Errorf("scaling down operation failed: %w", err)
	}

	// Step 4: Create encrypted scaling log
	encryptedLog, err := createEncryptedLog(modelID, "scale_down")
	if err != nil {
		logAction("ModelScalingDownError", fmt.Sprintf("Failed to encrypt scaling log for model: %s", modelID))
		return fmt.Errorf("failed to encrypt scaling log: %w", err)
	}

	// Step 5: Record scaling operation in the ledger
	if err := l.AiMLMLedger.RecordScaling(transactionID, modelID, "down", encryptedLog); err != nil {
		logAction("ModelScalingDownError", fmt.Sprintf("Failed to record scaling in ledger for model: %s", modelID))
		return fmt.Errorf("failed to record scaling in ledger: %w", err)
	}

	logAction("ModelScalingDownCompleted", fmt.Sprintf("Scale-down completed for model: %s", modelID))
	return nil
}


// CanScaleUp checks if the model can be scaled up.
func (m *Model) CanScaleUp() bool {
	logAction("CanScaleUp", fmt.Sprintf("Checking if model can scale up: %s", m.ModelID))
	return m.CurrentScale < m.MaxScale
}

// CanScaleDown checks if the model can be scaled down.
func (m *Model) CanScaleDown() bool {
	logAction("CanScaleDown", fmt.Sprintf("Checking if model can scale down: %s", m.ModelID))
	return m.CurrentScale > m.MinScale
}

// ScaleUp performs the logic to increase the model’s scale.
func (m *Model) ScaleUp() error {
	if !m.CanScaleUp() {
		return errors.New("cannot scale up beyond maximum scale level")
	}
	logAction("ScaleUp", fmt.Sprintf("Scaling up model: %s", m.ModelID))
	m.CurrentScale++
	return nil
}

// ScaleDown performs the logic to decrease the model’s scale.
func (m *Model) ScaleDown() error {
	if !m.CanScaleDown() {
		return errors.New("cannot scale down below minimum scale level")
	}
	logAction("ScaleDown", fmt.Sprintf("Scaling down model: %s", m.ModelID))
	m.CurrentScale--
	return nil
}