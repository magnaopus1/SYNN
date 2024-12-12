package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"log"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
	"time"
)

// AIServiceRegister registers a new AI service, integrating with the ledger and Synnergy Consensus for validation and recording.
func AIServiceRegister(l *ledger.Ledger, sc *common.SynnergyConsensus, serviceID, serviceName, owner string, serviceData map[string]interface{}) error {

	// Convert serviceData to JSON for encryption
	serviceDataBytes, err := json.Marshal(serviceData)
	if err != nil {
		return errors.New("failed to serialize service data for encryption")
	}

	// Encrypt serialized service data for storage in the ledger
	encryptedData, err := encryptDataForStorage(serviceDataBytes, serviceID)
	if err != nil {
		return errors.New("failed to encrypt service data for registration")
	}

	// Prepare metrics for recording in the ledger
	metrics := ledger.ServiceMetrics{
		ServiceName:   serviceName,
		Owner:         owner,
		EncryptedData: encryptedData,
		CreatedAt:     time.Now(),
		Status:        "Active",
	}

	// Record the new service in the ledger
	if err := l.AiMLMLedger.RecordAiService(serviceID, metrics); err != nil {
		return errors.New("failed to record service registration in the ledger")
	}

	log.Printf("Service %s successfully registered with ID %s", serviceName, serviceID)
	return nil
}


// AIServiceDeregister removes a registered AI service, ensuring ledger update and transaction validation.
func AIServiceDeregister(l *ledger.Ledger, sc *common.SynnergyConsensus, serviceID string) error {

	// Remove the service record from the ledger
	if err := l.AiMLMLedger.RemoveAiService(serviceID); err != nil {
		return errors.New("failed to deregister service in the ledger")
	}

	log.Printf("Service with ID %s successfully deregistered", serviceID)
	return nil
}


// AIServiceBalanceCheck checks the balance for an AI service, returning the service's balance.
func AIServiceBalanceCheck(l *ledger.Ledger, serviceID string) (float64, error) {
	// Retrieve service balance from the ledger
	balance, err := l.AiMLMLedger.GetServiceBalance(serviceID)
	if err != nil {
		return 0, errors.New("failed to retrieve service balance from the ledger")
	}

	log.Printf("Service with ID %s has a balance of %.2f", serviceID, balance)
	return balance, nil
}

// AIServiceMonitor continually checks the operational status and performance of a registered AI service.
func AIServiceMonitor(l *ledger.Ledger, sc *common.SynnergyConsensus, serviceID string) error {

	for {
		// Get the service's performance metrics
		metrics, err := l.AiMLMLedger.GetServiceStatusAndMetrics(serviceID)
		if err != nil {
			return errors.New("failed to retrieve service metrics from the ledger")
		}

		log.Printf("Service ID %s Metrics: %v", serviceID, metrics)

		// Log the current metrics in the ledger for auditing
		if err := l.AiMLMLedger.LogServiceMetrics(serviceID, metrics); err != nil {
			return errors.New("failed to log service metrics in the ledger")
		}

		// Wait for a defined interval before the next status check
		time.Sleep(5 * time.Minute)
	}
}


// AIServiceAvailability checks if an AI service is available for use, based on its status in the ledger.
func AIServiceAvailability(l *ledger.Ledger, serviceID string) (bool, error) {
	// Retrieve the current service status from the ledger
	status, err := l.AiMLMLedger.GetServiceStatus(serviceID)
	if err != nil {
		return false, errors.New("failed to retrieve service status from the ledger")
	}

	isAvailable := status == "active"
	log.Printf("Service ID %s availability: %t", serviceID, isAvailable)
	return isAvailable, nil
}
