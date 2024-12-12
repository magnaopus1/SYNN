package ai_ml_operation

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"synnergy_network/pkg/ledger"
	"time"
)

// ModelAccessRestrict applies restrictions to a model, limiting access based on user roles, usage counts, or other parameters.
func ModelAccessRestrict(modelID string, reason string, ledgerInstance *ledger.Ledger) error {
	// Record restrictions in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordModelRestriction(modelID, reason); err != nil {
		return errors.New("failed to record model access restriction in ledger")
	}

	return nil
}

// ModelPermissionsUpdate allows updating access permissions for a model, modifying who can interact with the model.
func ModelPermissionsUpdate(modelID string, users []string, ledgerInstance *ledger.Ledger) error {
	// Update permissions in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordModelPermissions(modelID, users); err != nil {
		return errors.New("failed to record permissions update in ledger")
	}

	return nil
}

// ModelAccessToken generates and securely stores an access token for a model, allowing controlled temporary access.
func ModelAccessToken(modelID, grantedTo, permissions string, duration int, ledgerInstance *ledger.Ledger) (string, error) {
	tokenID := generateAccessToken(modelID, duration)
	expiry := time.Now().Add(time.Duration(duration) * time.Minute)

	// Store access token in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordAccessToken(tokenID, modelID, grantedTo, permissions, expiry); err != nil {
		return "", errors.New("failed to record access token in ledger")
	}

	return tokenID, nil
}

// ModelAccessList retrieves and decrypts the list of users or entities with access permissions for a model.
func ModelAccessList(modelID string, ledgerInstance *ledger.Ledger) ([]string, error) {
	// Retrieve access list from ledger
	accessList, err := ledgerInstance.AiMLMLedger.GetModelAccessList(modelID)
	if err != nil {
		return nil, errors.New("failed to retrieve model access list from ledger")
	}

	return accessList, nil
}

// ModelAccessRemove revokes access for specific users or entities and updates the ledger.
func ModelAccessRemove(modelID string, userID string, ledgerInstance *ledger.Ledger) error {
	// Retrieve current access list from ledger
	currentAccessList, err := ledgerInstance.AiMLMLedger.GetModelAccessList(modelID)
	if err != nil {
		return errors.New("failed to retrieve current access list")
	}

	// Remove specified user and update access list
	for i, user := range currentAccessList {
		if user == userID {
			currentAccessList = append(currentAccessList[:i], currentAccessList[i+1:]...)
			break
		}
	}

	// Record updated access list in the ledger
	if err := ledgerInstance.AiMLMLedger.UpdateModelAccessList(modelID, currentAccessList); err != nil {
		return errors.New("failed to record access removal in ledger")
	}

	return nil
}

// generateAccessToken creates a time-limited token for model access control
func generateAccessToken(modelID string, duration int) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", modelID, duration)))
	return hex.EncodeToString(hash[:])
}
