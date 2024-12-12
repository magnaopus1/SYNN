package ledger

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)

// RecordOracleSubmission records the submission of oracle data into the ledger.
func (l *DataManagementLedger) RecordOracleSubmission(oracleID string, data map[string]interface{}) (string, error) {
	l.Lock()
	defer l.Unlock()

	if l.Oracles == nil {
		l.Oracles = make(map[string]OracleSubmission)
	}

	// Generate a unique ID without passing oracleID
	id := l.generateUniqueID()
	submission := OracleSubmission{
		ID:        id,
		OracleID:  oracleID,
		Data:      data,
		Submitted: time.Now(),
		Verified:  false,
	}

	l.Oracles[id] = submission
	return id, nil
}

// generateUniqueID generates a globally unique identifier using random bytes.
func (l *DataManagementLedger) generateUniqueID() string {
	l.Lock()
	defer l.Unlock()

	// Use random bytes to generate a unique ID
	b := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a counter-based ID in case of an error
		l.UniqueIDCounter++
		return "ID-" + fmt.Sprintf("%d", l.UniqueIDCounter)
	}
	return hex.EncodeToString(b)
}

// RecordOracleVerification verifies the oracle submission and records it in the ledger.
func (l *DataManagementLedger) RecordOracleVerification(submissionID string) error {
	l.Lock()
	defer l.Unlock()

	submission, exists := l.Oracles[submissionID]
	if !exists {
		return errors.New("oracle submission not found")
	}

	submission.Verified = true
	l.Oracles[submissionID] = submission
	return nil
}

// RecordDataTransferMonitor logs data transfer metrics in the ledger with a timestamp.
func (l *DataManagementLedger) RecordDataTransferMonitor(metrics DataTransferMetrics, channelID string) {
    l.Lock()
    defer l.Unlock()

    // Create a new DataTransferRecord
    record := DataTransferRecord{
        RecordID:   l.generateUniqueID(),
        ChannelID:  channelID,
        DataBlocks: nil, // Add DataBlocks if available
        CreatedAt:  time.Now(),
    }

    // Add or update the DataTransferRecord in the map
    l.DataTransferRecords[channelID] = record

    log.Printf("Data transfer metrics recorded for channel %s at %s: %v MB/s, Peak rate: %v MB/s",
        channelID, record.CreatedAt.Format(time.RFC3339), metrics.RateMBps, metrics.PeakRateMBps)
}

// RecordAggregationValidation validates and records an aggregation in the ledger.
func (l *DataManagementLedger) RecordAggregationValidation(aggregatorID string) (string, error) {
	l.Lock()
	defer l.Unlock()

	if l.Aggregations == nil {
		l.Aggregations = make(map[string]AggregationValidation)
	}

	// Generate a unique ID without passing aggregatorID
	id := l.generateUniqueID()
	validation := AggregationValidation{
		ID:         id,
		Aggregator: aggregatorID,
		Validated:  true,
		Timestamp:  time.Now(),
	}

	l.Aggregations[id] = validation
	return id, nil
}



// RecordDataTransmission records the transmission of data between entities.
func (l *DataManagementLedger) RecordDataTransmission(sender, receiver string, data map[string]interface{}) (string, error) {
	l.Lock()
	defer l.Unlock()

	id := l.generateUniqueID() // No argument passed
	transmission := DataTransmission{
		ID:            id,
		Sender:        sender,
		Receiver:      receiver,
		Data:          data,
		TransmittedAt: time.Now(),
	}

	l.DataTransmissions = append(l.DataTransmissions, transmission)
	return id, nil
}
