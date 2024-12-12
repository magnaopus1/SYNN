package rollups

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewMultiRollupOracle initializes a new multi-rollup oracle
func NewMultiRollupOracle(oracleID string, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, networkManager *common.NetworkManager) *common.MultiRollupOracle {
	return &common.MultiRollupOracle{
		OracleID:       oracleID,
		Rollups:        make(map[string]*Rollup),
		DataSources:    make(map[string]interface{}),
		Ledger:         ledgerInstance,
		Encryption:     encryptionService,
		NetworkManager: networkManager,
	}
}

// RegisterRollup registers a rollup to be served by the oracle
func (mo *common.MultiRollupOracle) RegisterRollup(rollupID string, rollup *common.Rollup) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	if _, exists := mo.Rollups[rollupID]; exists {
		return errors.New("rollup is already registered with the oracle")
	}

	mo.Rollups[rollupID] = rollup

	// Log the rollup registration in the ledger
	err := mo.Ledger.RecordOracleRollupRegistration(mo.OracleID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup registration: %v", err)
	}

	fmt.Printf("Rollup %s registered with oracle %s\n", rollupID, mo.OracleID)
	return nil
}

// FetchData asynchronously fetches data from external sources for a rollup
func (mo *common.MultiRollupOracle) FetchData(rollupID string, dataSourceID string) (interface{}, error) {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	// Validate rollup exists
	rollup, exists := mo.Rollups[rollupID]
	if !exists {
		return nil, fmt.Errorf("rollup %s not found", rollupID)
	}

	// Fetch data from the external data source (simulated here)
	data, exists := mo.DataSources[dataSourceID]
	if !exists {
		return nil, fmt.Errorf("data source %s not found", dataSourceID)
	}

	// Encrypt data before sending to the rollup
	encryptedData, err := mo.Encryption.EncryptData([]byte(fmt.Sprintf("%v", data)), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Simulate sending data to the rollup (broadcast or direct send)
	err = mo.NetworkManager.BroadcastData(rollup.RollupID, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to send data to rollup %s: %v", rollupID, err)
	}

	// Log data fetch event in the ledger
	err = mo.Ledger.RecordOracleDataFetch(mo.OracleID, rollupID, dataSourceID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log data fetch: %v", err)
	}

	fmt.Printf("Data fetched from source %s for rollup %s\n", dataSourceID, rollupID)
	return data, nil
}

// ValidateOracleData validates the data fetched by the oracle for a specific rollup
func (mo *common.MultiRollupOracle) ValidateOracleData(rollupID string, dataHash string) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	// Validate rollup exists
	rollup, exists := mo.Rollups[rollupID]
	if !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	// Simulate the data validation process (e.g., validate the hash)
	if dataHash == "" {
		return errors.New("invalid data hash")
	}

	// Log the data validation in the ledger
	err := mo.Ledger.RecordOracleDataValidation(mo.OracleID, rollupID, dataHash, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data validation: %v", err)
	}

	fmt.Printf("Data for rollup %s validated by oracle %s\n", rollupID, mo.OracleID)
	return nil
}

// RegisterDataSource registers an external data source for the oracle
func (mo *common.MultiRollupOracle) RegisterDataSource(sourceID string, dataSource interface{}) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	if _, exists := mo.DataSources[sourceID]; exists {
		return errors.New("data source is already registered")
	}

	mo.DataSources[sourceID] = dataSource

	// Log the data source registration in the ledger
	err := mo.Ledger.RecordDataSourceRegistration(mo.OracleID, sourceID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data source registration: %v", err)
	}

	fmt.Printf("Data source %s registered with oracle %s\n", sourceID, mo.OracleID)
	return nil
}

// RemoveRollup unregisters a rollup from the oracle
func (mo *common.MultiRollupOracle) RemoveRollup(rollupID string) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	if _, exists := mo.Rollups[rollupID]; !exists {
		return fmt.Errorf("rollup %s not found", rollupID)
	}

	delete(mo.Rollups, rollupID)

	// Log the rollup removal in the ledger
	err := mo.Ledger.RecordOracleRollupRemoval(mo.OracleID, rollupID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log rollup removal: %v", err)
	}

	fmt.Printf("Rollup %s removed from oracle %s\n", rollupID, mo.OracleID)
	return nil
}

// RemoveDataSource unregisters a data source from the oracle
func (mo *common.MultiRollupOracle) RemoveDataSource(sourceID string) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	if _, exists := mo.DataSources[sourceID]; !exists {
		return fmt.Errorf("data source %s not found", sourceID)
	}

	delete(mo.DataSources, sourceID)

	// Log the data source removal in the ledger
	err := mo.Ledger.RecordDataSourceRemoval(mo.OracleID, sourceID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log data source removal: %v", err)
	}

	fmt.Printf("Data source %s removed from oracle %s\n", sourceID, mo.OracleID)
	return nil
}
