package ai_ml_operation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"

	"golang.org/x/exp/rand"
)

// ModelDataPreprocess handles data preprocessing tasks such as normalization and formatting, records actions in the ledger, and encrypts results.
func ModelDataPreprocess(modelID string, rawData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
	preprocessedData := preprocess(rawData)

	// Encrypt the processed data
	encryptedData, err := encryptData(preprocessedData, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt preprocessed data")
	}

	// Log the action in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Preprocess"); err != nil {
		return nil, errors.New("failed to record preprocessing action in ledger")
	}

	return encryptedData, nil
}

// ModelDataPostprocess handles post-processing steps like scaling results and formatting output data, with ledger tracking.
func ModelDataPostprocess(modelID string, processedData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
	postprocessedData := postprocess(processedData)

	// Encrypt the postprocessed data
	encryptedData, err := encryptData(postprocessedData, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt postprocessed data")
	}

	// Record in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Postprocess"); err != nil {
		return nil, errors.New("failed to record post-processing action in ledger")
	}

	return encryptedData, nil
}

// ModelDataCleaning removes noise and invalid data points, ensuring model data quality and integrity.
func ModelDataCleaning(modelID string, rawData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
	cleanedData := cleanData(rawData)

	// Encrypt cleaned data
	encryptedData, err := encryptData(cleanedData, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt cleaned data")
	}

	// Record in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Cleaning"); err != nil {
		return nil, errors.New("failed to record data cleaning in ledger")
	}

	return encryptedData, nil
}

// ModelDataAugment performs data augmentation techniques to enrich the dataset.
func ModelDataAugment(modelID string, baseData []byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
	augmentedData := augmentData(baseData)

	// Encrypt augmented data
	encryptedData, err := encryptData(augmentedData, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt augmented data")
	}

	// Record in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Augment"); err != nil {
		return nil, errors.New("failed to record data augmentation in ledger")
	}

	return encryptedData, nil
}

// ModelDataAccess provides controlled access to encrypted model data for further processing.
func ModelDataAccess(modelID string, ledgerInstance *ledger.Ledger) ([]byte, error) {
	// Access data from ledger and decrypt
	encryptedData, err := ledgerInstance.AiMLMLedger.GetModelData(modelID)
	if err != nil {
		return nil, errors.New("failed to access model data from ledger")
	}

	// Decrypt data
	data, err := decryptData([]byte(encryptedData[0].Status), modelID) // Assuming we're using the Status field for demo purposes
	if err != nil {
		return nil, errors.New("failed to decrypt model data")
	}

	// Ensure data is of type []byte
	byteData, ok := data.([]byte)
	if !ok {
		return nil, errors.New("decrypted data is not of type []byte")
	}

	return byteData, nil
}


// ModelDataLabel allows labeling of data for supervised learning tasks.
func ModelDataLabel(modelID string, labeledData map[string]string, ledgerInstance *ledger.Ledger) ([]byte, error) {
	// Convert labeled data to JSON and encrypt it
	dataJSON, err := json.Marshal(labeledData)
	if err != nil {
		return nil, errors.New("failed to marshal labeled data to JSON")
	}

	encryptedData, err := encryptData(dataJSON, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt labeled data")
	}

	// Record labeling in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Labeling"); err != nil {
		return nil, errors.New("failed to record data labeling in ledger")
	}

	return encryptedData, nil
}

// ModelDataMerge merges multiple datasets and records the action in the ledger.
func ModelDataMerge(modelID string, datasets [][]byte, ledgerInstance *ledger.Ledger) ([]byte, error) {
	mergedData := mergeDatasets(datasets)

	// Encrypt merged data
	encryptedData, err := encryptData(mergedData, modelID)
	if err != nil {
		return nil, errors.New("failed to encrypt merged data")
	}

	// Record merging in the ledger
	if err := ledgerInstance.AiMLMLedger.RecordDataProcess(modelID, "Merge"); err != nil {
		return nil, errors.New("failed to record data merging in ledger")
	}

	return encryptedData, nil
}

// preprocess applies transformations like normalization, scaling, or encoding, preparing data for model training.
func preprocess(data []byte) []byte {
	log.Println("Starting preprocessing...")

	// Decode JSON data for structured processing (assuming JSON input)
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Example normalization: Scale numerical values between 0 and 1
	for key, value := range rawData {
		if num, ok := value.(float64); ok {
			rawData[key] = num / 100.0 // Scaling example
		}
	}

	// Encode back to JSON after preprocessing
	preprocessedData, err := json.Marshal(rawData)
	if err != nil {
		log.Fatalf("Error encoding preprocessed data to JSON: %v", err)
	}

	log.Println("Preprocessing completed.")
	return preprocessedData
}

// postprocess applies transformations like rounding, scaling results, or formatting for final output.
func postprocess(data []byte) []byte {
	log.Println("Starting postprocessing...")

	// Decode JSON data for structured processing (assuming JSON input)
	var processedData map[string]interface{}
	if err := json.Unmarshal(data, &processedData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Example: Round all float values to two decimal points for display
	for key, value := range processedData {
		if num, ok := value.(float64); ok {
			processedData[key] = roundToTwoDecimals(num)
		}
	}

	// Encode back to JSON after postprocessing
	postprocessedData, err := json.Marshal(processedData)
	if err != nil {
		log.Fatalf("Error encoding postprocessed data to JSON: %v", err)
	}

	log.Println("Postprocessing completed.")
	return postprocessedData
}

// cleanData removes missing, invalid, or noisy data points from the dataset.
func cleanData(data []byte) []byte {
	log.Println("Starting data cleaning...")

	// Decode JSON data for structured processing (assuming JSON input)
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Remove entries with missing values or invalid types
	for key, value := range rawData {
		if value == nil {
			delete(rawData, key)
		} else if num, ok := value.(float64); ok && (num < 0 || num > 100) {
			// Example noise removal: remove out-of-range values
			delete(rawData, key)
		}
	}

	// Encode back to JSON after cleaning
	cleanedData, err := json.Marshal(rawData)
	if err != nil {
		log.Fatalf("Error encoding cleaned data to JSON: %v", err)
	}

	log.Println("Data cleaning completed.")
	return cleanedData
}

// augmentData generates additional data samples by applying transformations to existing data.
func augmentData(data []byte) []byte {
	log.Println("Starting data augmentation...")

	// Decode JSON data for structured processing (assuming JSON input)
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Example augmentation: Add small random noise to numerical values
	rand.Seed(uint64(time.Now().UnixNano())) // Cast to uint64 to match the required type
	for key, value := range rawData {
		if num, ok := value.(float64); ok {
			noise := (rand.Float64() - 0.5) * 0.1 // Small random noise
			rawData[key] = num + noise
		}
	}

	// Encode back to JSON after augmentation
	augmentedData, err := json.Marshal(rawData)
	if err != nil {
		log.Fatalf("Error encoding augmented data to JSON: %v", err)
	}

	log.Println("Data augmentation completed.")
	return augmentedData
}




// mergeDatasets combines multiple datasets by merging them into one coherent dataset.
func mergeDatasets(datasets [][]byte) []byte {
	log.Println("Starting dataset merging...")

	var mergedData map[string]interface{}
	mergedData = make(map[string]interface{})

	// Iterate through each dataset and merge it into `mergedData`
	for _, dataset := range datasets {
		var data map[string]interface{}
		if err := json.Unmarshal(dataset, &data); err != nil {
			log.Fatalf("Error decoding dataset JSON: %v", err)
		}

		// Append new keys or overwrite existing keys
		for key, value := range data {
			mergedData[key] = value
		}
	}

	// Encode back to JSON after merging
	mergedDataBytes, err := json.Marshal(mergedData)
	if err != nil {
		log.Fatalf("Error encoding merged data to JSON: %v", err)
	}

	log.Println("Dataset merging completed.")
	return mergedDataBytes
}

// Helper function to round float64 to two decimal places
func roundToTwoDecimals(num float64) float64 {
	return float64(int(num*100+0.5)) / 100
}





// Placeholder unique ID generator
func generateUniqueID() string {
	return fmt.Sprintf("id-%d", time.Now().UnixNano())
}
