package ledger

import (
	"fmt"
	"time"
)

// RecordOracleRollupRegistration logs the registration of a new Oracle for the rollup.
func (ledger *RollupLedger) RecordOracleRollupRegistration(rollupID, oracleID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		rollup.OracleSources = append(rollup.OracleSources, oracleID)
		ledger.Rollups[rollupID] = rollup // Update the rollup in the ledger
		fmt.Printf("Oracle %s registered for Rollup %s.\n", oracleID, rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordRollupCreation records the creation of a rollup in the ledger.
func (ledger *RollupLedger) RecordRollupCreation(validatorID string) (string, error) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Rollups == nil {
		ledger.Rollups = make(map[string]Rollup)
	}

	id := generateUniqueID() // Assume `generateUniqueID` generates a unique string

	rollup := Rollup{
		RollupID:         id,
		ValidatorAddress: validatorID,
		CreatedAt:        time.Now(),
		Transactions:     []*Transaction{},
		IsFinalized:      false,
	}

	ledger.Rollups[id] = rollup
	fmt.Printf("Rollup %s created by Validator %s.\n", id, validatorID)
	return id, nil
}

// RecordOracleDataFetch logs the fetching of data by an Oracle and updates the ledger.
func (ledger *RollupLedger) RecordOracleDataFetch(oracleID string) {
	ledger.Lock()
	defer ledger.Unlock()

	// Initialize the OracleDataFetchLogs map if nil
	if ledger.OracleDataFetchLogs == nil {
		ledger.OracleDataFetchLogs = make(map[string][]time.Time)
	}

	// Record the fetch event
	ledger.OracleDataFetchLogs[oracleID] = append(ledger.OracleDataFetchLogs[oracleID], time.Now())

	fmt.Printf("Oracle %s fetched new data at %s.\n", oracleID, time.Now().Format(time.RFC3339))
}

// RecordOracleDataValidation logs the validation of data from an Oracle and updates the ledger.
func (ledger *RollupLedger) RecordOracleDataValidation(oracleID, validationResult string) {
	ledger.Lock()
	defer ledger.Unlock()

	// Initialize the OracleDataValidationLogs map if nil
	if ledger.OracleDataValidationLogs == nil {
		ledger.OracleDataValidationLogs = make(map[string][]OracleValidationLog)
	}

	// Create a validation log entry
	validationLog := OracleValidationLog{
		OracleID:         oracleID,
		ValidationResult: validationResult,
		Timestamp:        time.Now(),
	}

	// Record the validation log
	ledger.OracleDataValidationLogs[oracleID] = append(ledger.OracleDataValidationLogs[oracleID], validationLog)

	fmt.Printf("Oracle %s data validated with result: %s at %s.\n", oracleID, validationResult, time.Now().Format(time.RFC3339))
}


// RecordDataSourceRegistration registers a data source for a rollup.
func (ledger *RollupLedger) RecordDataSourceRegistration(rollupID, dataSourceID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		rollup.DataSources = append(rollup.DataSources, dataSourceID)
		ledger.Rollups[rollupID] = rollup // Update the rollup in the ledger
		fmt.Printf("DataSource %s registered for Rollup %s.\n", dataSourceID, rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordOracleRollupRemoval removes an Oracle from a rollup.
func (ledger *RollupLedger) RecordOracleRollupRemoval(rollupID, oracleID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		for i, oracle := range rollup.OracleSources {
			if oracle == oracleID {
				rollup.OracleSources = append(rollup.OracleSources[:i], rollup.OracleSources[i+1:]...)
				ledger.Rollups[rollupID] = rollup // Update the rollup in the ledger
				fmt.Printf("Oracle %s removed from Rollup %s.\n", oracleID, rollupID)
				return
			}
		}
		fmt.Printf("Oracle %s not found in Rollup %s.\n", oracleID, rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordDataSourceRemoval removes a data source from a rollup.
func (ledger *RollupLedger) RecordDataSourceRemoval(rollupID, dataSourceID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		for i, dataSource := range rollup.DataSources {
			if dataSource == dataSourceID {
				rollup.DataSources = append(rollup.DataSources[:i], rollup.DataSources[i+1:]...)
				ledger.Rollups[rollupID] = rollup // Update the rollup in the ledger
				fmt.Printf("DataSource %s removed from Rollup %s.\n", dataSourceID, rollupID)
				return
			}
		}
		fmt.Printf("DataSource %s not found in Rollup %s.\n", dataSourceID, rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordTransactionAddition adds a transaction to the rollup.
func (ledger *RollupLedger) RecordTransactionAddition(transactionRecord *TransactionRecord) {
	ledger.Lock()
	defer ledger.Unlock()

	transaction := &Transaction{
		TransactionID: transactionRecord.TransactionID,
		FromAddress:   transactionRecord.FromAddress,
		ToAddress:     transactionRecord.ToAddress,
		Amount:        transactionRecord.Amount,
		Fee:           transactionRecord.Fee,
		Timestamp:     transactionRecord.Timestamp,
		Status:        transactionRecord.Status,
	}

	ledger.PendingTransactions = append(ledger.PendingTransactions, transaction)

	fmt.Printf("Transaction %s added to the ledger.\n", transactionRecord.TransactionID)
}

// RecordGovernanceProposal records a governance proposal for rollup management.
func (ledger *RollupLedger) RecordGovernanceProposal(proposalID, description string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.GovernanceProposals == nil {
		ledger.GovernanceProposals = make(map[string]GovernanceProposal)
	}

	ledger.GovernanceProposals[proposalID] = GovernanceProposal{
		ProposalID:   proposalID,
		Description:  description,
		CreationTime: time.Now(),
	}

	fmt.Printf("Governance proposal %s added with description: %s.\n", proposalID, description)
}

// RecordRollupLayerFinalization finalizes a layer in the rollup.
func (ledger *RollupLedger) RecordRollupLayerFinalization(rollupID, layerID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		rollup.FinalizedLayers = append(rollup.FinalizedLayers, layerID)
		ledger.Rollups[rollupID] = rollup
		fmt.Printf("Layer %s of Rollup %s finalized.\n", layerID, rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordCrossLayerVerification logs the verification process between rollup layers.
func (ledger *RollupLedger) RecordCrossLayerVerification(layerID, verificationResult string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.LayerVerification == nil {
		ledger.LayerVerification = make(map[string]string)
	}

	ledger.LayerVerification[layerID] = verificationResult

	fmt.Printf("Cross-layer verification of layer %s completed with result: %s\n", layerID, verificationResult)
}

// RecordRollupRebalancing logs the rebalancing of a rollup.
func (ledger *RollupLedger) RecordRollupRebalancing(rollupID, details string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		rollup.RebalanceDetails = details
		ledger.Rollups[rollupID] = rollup
		fmt.Printf("Rollup %s rebalanced: %s\n", rollupID, details)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordRollupRemoval removes a rollup from the ledger.
func (ledger *RollupLedger) RecordRollupRemoval(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if _, exists := ledger.Rollups[rollupID]; exists {
		delete(ledger.Rollups, rollupID)
		fmt.Printf("Rollup %s removed.\n", rollupID)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordStateUpdate records an update in the state of the rollup.
func (ledger *RollupLedger) RecordStateUpdate(rollupID, newState string) {
	ledger.Lock()
	defer ledger.Unlock()

	if rollup, exists := ledger.Rollups[rollupID]; exists {
		rollup.State = newState
		ledger.Rollups[rollupID] = rollup
		fmt.Printf("Rollup %s state updated to: %s\n", rollupID, newState)
	} else {
		fmt.Printf("Rollup %s not found.\n", rollupID)
	}
}

// RecordCrossRollupTransaction logs a cross-rollup transaction.
func (ledger *RollupLedger) RecordCrossRollupTransaction(txID, sourceRollup, targetRollup string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.CrossRollupTransactions == nil {
		ledger.CrossRollupTransactions = make(map[string]CrossRollupTransaction)
	}

	ledger.CrossRollupTransactions[txID] = CrossRollupTransaction{
		TxID:         txID,
		SourceRollup: sourceRollup,
		TargetRollup: targetRollup,
		Timestamp:    time.Now(),
	}

	fmt.Printf("Cross-rollup transaction %s from Rollup %s to Rollup %s recorded.\n", txID, sourceRollup, targetRollup)
}

// RecordLayerFinalization logs the finalization of a layer in a rollup.
func (ledger *RollupLedger) RecordLayerFinalization(layerID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.LayerFinalization == nil {
		ledger.LayerFinalization = make(map[string]time.Time)
	}

	ledger.LayerFinalization[layerID] = time.Now()

	fmt.Printf("Layer %s finalized.\n", layerID)
}

// RecordSharedStateValidation logs the validation of shared state between rollups.
func (ledger *RollupLedger) RecordSharedStateValidation(sharedStateID, validationResult string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.SharedStateValidations == nil {
		ledger.SharedStateValidations = make(map[string]string)
	}

	ledger.SharedStateValidations[sharedStateID] = validationResult

	fmt.Printf("Shared state %s validated with result: %s\n", sharedStateID, validationResult)
}

// RecordLiquidityAddition logs the addition of liquidity to a rollup pool.
func (ledger *RollupLedger) RecordLiquidityAddition(poolID, liquidityAmount string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.LiquidityPools == nil {
		ledger.LiquidityPools = make(map[string]LiquidityPool)
	}

	pool := ledger.LiquidityPools[poolID]
	pool.TotalLiquidity += liquidityAmount
	ledger.LiquidityPools[poolID] = pool

	fmt.Printf("Liquidity of %s added to pool %s.\n", liquidityAmount, poolID)
}

// RecordLiquidityRemoval logs the removal of liquidity from a rollup pool.
func (ledger *RollupLedger) RecordLiquidityRemoval(poolID, liquidityAmount string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.LiquidityPools == nil {
		ledger.LiquidityPools = make(map[string]LiquidityPool)
	}

	pool := ledger.LiquidityPools[poolID]
	pool.TotalLiquidity -= liquidityAmount
	ledger.LiquidityPools[poolID] = pool

	fmt.Printf("Liquidity of %s removed from pool %s.\n", liquidityAmount, poolID)
}


// RecordSpaceTimeProofGeneration logs the generation of space-time proof in the rollup.
func (ledger *RollupLedger) RecordSpaceTimeProofGeneration(rollupID, proofData string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.SpaceTimeProofs == nil {
		ledger.SpaceTimeProofs = make(map[string][]ProofRecord)
	}

	ledger.SpaceTimeProofs[rollupID] = append(ledger.SpaceTimeProofs[rollupID], ProofRecord{
		ProofID:   generateUUID(),
		ProofData: proofData,
		Action:    "generated",
		Timestamp: time.Now(),
	})

	fmt.Printf("Space-time proof generated for Rollup %s: %s\n", rollupID, proofData)
}

// RecordSpaceTimeProofVerification logs the verification of space-time proof in the rollup.
func (ledger *RollupLedger) RecordSpaceTimeProofVerification(rollupID, verificationResult string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.SpaceTimeProofs == nil {
		ledger.SpaceTimeProofs = make(map[string][]ProofRecord)
	}

	ledger.SpaceTimeProofs[rollupID] = append(ledger.SpaceTimeProofs[rollupID], ProofRecord{
		ProofID:   generateUUID(),
		Action:    "verified",
		Details:   verificationResult,
		Timestamp: time.Now(),
	})

	fmt.Printf("Space-time proof for Rollup %s verified with result: %s\n", rollupID, verificationResult)
}

// RecordProofAddition adds proof to a rollup.
func (ledger *RollupLedger) RecordProofAddition(rollupID, proofData string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Proofs == nil {
		ledger.Proofs = make(map[string][]ProofRecord)
	}

	ledger.Proofs[rollupID] = append(ledger.Proofs[rollupID], ProofRecord{
		ProofID:   generateUUID(),
		ProofData: proofData,
		Action:    "added",
		Timestamp: time.Now(),
	})

	fmt.Printf("Proof added to Rollup %s: %s\n", rollupID, proofData)
}

// RecordProofAggregationFinalization logs the finalization of proof aggregation for a rollup.
func (ledger *RollupLedger) RecordProofAggregationFinalization(rollupID, aggregationDetails string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.ProofAggregations == nil {
		ledger.ProofAggregations = make(map[string][]ProofAggregationRecord)
	}

	ledger.ProofAggregations[rollupID] = append(ledger.ProofAggregations[rollupID], ProofAggregationRecord{
		AggregationID:    generateUUID(),
		AggregationDetails: aggregationDetails,
		Action:           "finalized",
		Timestamp:        time.Now(),
	})

	fmt.Printf("Proof aggregation for Rollup %s finalized: %s\n", rollupID, aggregationDetails)
}

// RecordNodeRemoval logs the removal of a node from the rollup.
func (ledger *RollupLedger) RecordNodeRemoval(nodeID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.NodeConnections == nil {
		ledger.NodeConnections = make(map[string][]string)
	}

	for rollupID, nodes := range ledger.NodeConnections {
		for i, id := range nodes {
			if id == nodeID {
				ledger.NodeConnections[rollupID] = append(nodes[:i], nodes[i+1:]...)
				break
			}
		}
	}

	fmt.Printf("Node %s removed from the rollup.\n", nodeID)
}

// RecordRollupSync logs the synchronization of a rollup.
func (ledger *RollupLedger) RecordRollupSync(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.SyncRecords == nil {
		ledger.SyncRecords = make(map[string]time.Time)
	}

	ledger.SyncRecords[rollupID] = time.Now()

	fmt.Printf("Rollup %s synchronized.\n", rollupID)
}

// RecordRollupValidation logs the validation of a rollup.
func (ledger *RollupLedger) RecordRollupValidation(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.ValidationRecords == nil {
		ledger.ValidationRecords = make(map[string][]ValidationRecord)
	}

	ledger.ValidationRecords[rollupID] = append(ledger.ValidationRecords[rollupID], ValidationRecord{
		Action:    "validated",
		Timestamp: time.Now(),
	})

	fmt.Printf("Rollup %s validated.\n", rollupID)
}

// RecordBridgeTransaction logs a bridge transaction between two rollups.
func (ledger *RollupLedger) RecordBridgeTransaction(txID, sourceRollup, targetRollup string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BridgeTransactions == nil {
		ledger.BridgeTransactions = make(map[string]BridgeRecord)
	}

	ledger.BridgeTransactions[txID] = BridgeRecord{
		TransactionID: txID,
		SourceRollup:  sourceRollup,
		TargetRollup:  targetRollup,
		Timestamp:     time.Now(),
	}

	fmt.Printf("Bridge transaction %s from Rollup %s to Rollup %s recorded.\n", txID, sourceRollup, targetRollup)
}

// RecordChallengeSubmission logs a challenge submission on a rollup.
func (ledger *RollupLedger) RecordChallengeSubmission(challengeID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Challenges == nil {
		ledger.Challenges = make(map[string][]ChallengeRecord)
	}

	ledger.Challenges[rollupID] = append(ledger.Challenges[rollupID], ChallengeRecord{
		ChallengeID: challengeID,
		Action:      "submitted",
		Timestamp:   time.Now(),
	})

	fmt.Printf("Challenge %s submitted on Rollup %s.\n", challengeID, rollupID)
}

// RecordChallengeResolution logs the resolution of a challenge on a rollup.
func (ledger *RollupLedger) RecordChallengeResolution(challengeID, rollupID, resolution string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Challenges == nil {
		ledger.Challenges = make(map[string][]ChallengeRecord)
	}

	ledger.Challenges[rollupID] = append(ledger.Challenges[rollupID], ChallengeRecord{
		ChallengeID: challengeID,
		Action:      "resolved",
		Details:     resolution,
		Timestamp:   time.Now(),
	})

	fmt.Printf("Challenge %s on Rollup %s resolved with result: %s.\n", challengeID, rollupID, resolution)
}

// RecordChallengeEscalation logs the escalation of a challenge on a rollup.
func (ledger *RollupLedger) RecordChallengeEscalation(challengeID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Challenges == nil {
		ledger.Challenges = make(map[string][]ChallengeRecord)
	}

	ledger.Challenges[rollupID] = append(ledger.Challenges[rollupID], ChallengeRecord{
		ChallengeID: challengeID,
		Action:      "escalated",
		Timestamp:   time.Now(),
	})

	fmt.Printf("Challenge %s on Rollup %s escalated.\n", challengeID, rollupID)
}


// RecordFeeCalculation logs the calculation of fees for a rollup.
func (ledger *RollupLedger) RecordFeeCalculation(rollupID string, feeAmount float64) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.FeeRecords == nil {
		ledger.FeeRecords = make(map[string][]FeeRecord)
	}

	ledger.FeeRecords[rollupID] = append(ledger.FeeRecords[rollupID], FeeRecord{
		RollupID:   rollupID,
		FeeAmount:  feeAmount,
		Action:     "calculated",
		Timestamp:  time.Now(),
	})

	fmt.Printf("Fee of %.2f calculated for Rollup %s.\n", feeAmount, rollupID)
}

// RecordFeeApplication logs the application of fees on a rollup.
func (ledger *RollupLedger) RecordFeeApplication(rollupID string, feeAmount float64) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.FeeRecords == nil {
		ledger.FeeRecords = make(map[string][]FeeRecord)
	}

	ledger.FeeRecords[rollupID] = append(ledger.FeeRecords[rollupID], FeeRecord{
		RollupID:   rollupID,
		FeeAmount:  feeAmount,
		Action:     "applied",
		Timestamp:  time.Now(),
	})

	fmt.Printf("Fee of %.2f applied to Rollup %s.\n", feeAmount, rollupID)
}

// RecordBaseFeeUpdate logs the update of the base fee for a rollup.
func (ledger *RollupLedger) RecordBaseFeeUpdate(rollupID string, newBaseFee float64) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BaseFees == nil {
		ledger.BaseFees = make(map[string]float64)
	}

	ledger.BaseFees[rollupID] = newBaseFee

	fmt.Printf("Base fee for Rollup %s updated to %.2f.\n", rollupID, newBaseFee)
}

// RecordFeeRefund logs a fee refund on a rollup.
func (ledger *RollupLedger) RecordFeeRefund(rollupID string, refundAmount float64) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.FeeRecords == nil {
		ledger.FeeRecords = make(map[string][]FeeRecord)
	}

	ledger.FeeRecords[rollupID] = append(ledger.FeeRecords[rollupID], FeeRecord{
		RollupID:   rollupID,
		FeeAmount:  refundAmount,
		Action:     "refunded",
		Timestamp:  time.Now(),
	})

	fmt.Printf("Fee of %.2f refunded on Rollup %s.\n", refundAmount, rollupID)
}

// RecordRollupNodeConnection logs the connection of a node to the rollup.
func (ledger *RollupLedger) RecordRollupNodeConnection(nodeID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.NodeConnections == nil {
		ledger.NodeConnections = make(map[string][]string)
	}

	ledger.NodeConnections[rollupID] = append(ledger.NodeConnections[rollupID], nodeID)

	fmt.Printf("Node %s connected to Rollup %s.\n", nodeID, rollupID)
}

// RecordNodeDisconnection logs the disconnection of a node from the rollup.
func (ledger *RollupLedger) RecordNodeDisconnection(nodeID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.NodeConnections == nil || len(ledger.NodeConnections[rollupID]) == 0 {
		fmt.Printf("Node %s is not connected to Rollup %s.\n", nodeID, rollupID)
		return
	}

	// Remove the nodeID from the list of connections
	for i, id := range ledger.NodeConnections[rollupID] {
		if id == nodeID {
			ledger.NodeConnections[rollupID] = append(ledger.NodeConnections[rollupID][:i], ledger.NodeConnections[rollupID][i+1:]...)
			break
		}
	}

	fmt.Printf("Node %s disconnected from Rollup %s.\n", nodeID, rollupID)
}

// RecordBatchCreation logs the creation of a batch on a rollup.
func (ledger *RollupLedger) RecordBatchCreation(batchID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Batches == nil {
		ledger.Batches = make(map[string][]BatchRecord)
	}

	ledger.Batches[rollupID] = append(ledger.Batches[rollupID], BatchRecord{
		BatchID:   batchID,
		Action:    "created",
		Timestamp: time.Now(),
	})

	fmt.Printf("Batch %s created on Rollup %s.\n", batchID, rollupID)
}

// RecordBatchValidation logs the validation of a batch on a rollup.
func (ledger *RollupLedger) RecordBatchValidation(batchID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BatchValidations == nil {
		ledger.BatchValidations = make(map[string][]string)
	}

	ledger.BatchValidations[rollupID] = append(ledger.BatchValidations[rollupID], batchID)

	fmt.Printf("Batch %s validated on Rollup %s.\n", batchID, rollupID)
}

// RecordBatchBroadcast logs the broadcasting of a batch on a rollup.
func (ledger *RollupLedger) RecordBatchBroadcast(batchID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BatchBroadcasts == nil {
		ledger.BatchBroadcasts = make(map[string][]string)
	}

	ledger.BatchBroadcasts[rollupID] = append(ledger.BatchBroadcasts[rollupID], batchID)

	fmt.Printf("Batch %s broadcasted on Rollup %s.\n", batchID, rollupID)
}

// RecordBatchSubmission logs the submission of a batch on a rollup.
func (ledger *RollupLedger) RecordBatchSubmission(batchID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BatchSubmissions == nil {
		ledger.BatchSubmissions = make(map[string][]string)
	}

	ledger.BatchSubmissions[rollupID] = append(ledger.BatchSubmissions[rollupID], batchID)

	fmt.Printf("Batch %s submitted on Rollup %s.\n", batchID, rollupID)
}

// RecordBatchRemoval logs the removal of a batch on a rollup.
func (ledger *RollupLedger) RecordBatchRemoval(batchID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Batches == nil || len(ledger.Batches[rollupID]) == 0 {
		fmt.Printf("Batch %s not found on Rollup %s.\n", batchID, rollupID)
		return
	}

	// Remove the batchID from the list of batches
	for i, batch := range ledger.Batches[rollupID] {
		if batch.BatchID == batchID {
			ledger.Batches[rollupID] = append(ledger.Batches[rollupID][:i], ledger.Batches[rollupID][i+1:]...)
			break
		}
	}

	fmt.Printf("Batch %s removed from Rollup %s.\n", batchID, rollupID)
}

// RecordRollupSubmission logs the submission of a rollup.
func (ledger *RollupLedger) RecordRollupSubmission(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.RollupSubmissions == nil {
		ledger.RollupSubmissions = make(map[string]time.Time)
	}

	ledger.RollupSubmissions[rollupID] = time.Now()

	fmt.Printf("Rollup %s submitted.\n", rollupID)
}

// RecordFraudProofSubmission logs the submission of a fraud proof for a rollup.
func (ledger *RollupLedger) RecordFraudProofSubmission(proofID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.FraudProofs == nil {
		ledger.FraudProofs = make(map[string][]string)
	}

	ledger.FraudProofs[rollupID] = append(ledger.FraudProofs[rollupID], proofID)

	fmt.Printf("Fraud proof %s submitted for Rollup %s.\n", proofID, rollupID)
}



// RecordFraudProofResolution logs the resolution of a fraud proof.
func (ledger *RollupLedger) RecordFraudProofResolution(proofID, rollupID, result string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.FraudProofs == nil {
		ledger.FraudProofs = make(map[string]FraudProofRecord)
	}

	ledger.FraudProofs[proofID] = FraudProofRecord{
		ProofID:  proofID,
		RollupID: rollupID,
		Result:   result,
		Timestamp: time.Now(),
	}

	fmt.Printf("Fraud proof %s for Rollup %s resolved with result: %s.\n", proofID, rollupID, result)
}

// RecordScalingAdjustment logs a scaling adjustment for the rollup.
func (ledger *RollupLedger) RecordScalingAdjustment(rollupID, newScale string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.RollupScalings == nil {
		ledger.RollupScalings = make(map[string]string)
	}

	ledger.RollupScalings[rollupID] = newScale

	fmt.Printf("Scaling adjustment for Rollup %s: new scale %s.\n", rollupID, newScale)
}

// RecordTransactionCancellation logs the cancellation of a transaction in the rollup.
func (ledger *RollupLedger) RecordTransactionCancellation(txID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.CancelledTransactions == nil {
		ledger.CancelledTransactions = make(map[string]string)
	}

	ledger.CancelledTransactions[txID] = rollupID

	fmt.Printf("Transaction %s cancelled in Rollup %s.\n", txID, rollupID)
}

// RecordRollupVerification logs the verification of a rollup.
func (ledger *RollupLedger) RecordRollupVerification(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.VerifiedRollups == nil {
		ledger.VerifiedRollups = make(map[string]time.Time)
	}

	ledger.VerifiedRollups[rollupID] = time.Now()

	fmt.Printf("Rollup %s verified.\n", rollupID)
}

// RecordZKProofGeneration logs the generation of a zk-proof in the rollup.
func (ledger *RollupLedger) RecordZKProofGeneration(proofID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.ZKProofs == nil {
		ledger.ZKProofs = make(map[string]ZKProofRecord)
	}

	ledger.ZKProofs[proofID] = ZKProofRecord{
		ProofID:  proofID,
		RollupID: rollupID,
		GeneratedAt: time.Now(),
	}

	fmt.Printf("zk-Proof %s generated for Rollup %s.\n", proofID, rollupID)
}

// RecordZKProofVerification logs the verification of a zk-proof in the rollup.
func (ledger *RollupLedger) RecordZKProofVerification(proofID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.VerifiedZKProofs == nil {
		ledger.VerifiedZKProofs = make(map[string]time.Time)
	}

	ledger.VerifiedZKProofs[proofID] = time.Now()

	fmt.Printf("zk-Proof %s verified for Rollup %s.\n", proofID, rollupID)
}

// RecordGovernanceUpdate logs a governance update in the rollup.
func (ledger *RollupLedger) RecordGovernanceUpdate(rollupID, update string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.GovernanceUpdates == nil {
		ledger.GovernanceUpdates = make(map[string][]string)
	}

	ledger.GovernanceUpdates[rollupID] = append(ledger.GovernanceUpdates[rollupID], update)

	fmt.Printf("Governance update for Rollup %s: %s.\n", rollupID, update)
}

// RecordGovernanceApplication logs the application of a governance decision in the rollup.
func (ledger *RollupLedger) RecordGovernanceApplication(rollupID, application string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.GovernanceApplications == nil {
		ledger.GovernanceApplications = make(map[string][]string)
	}

	ledger.GovernanceApplications[rollupID] = append(ledger.GovernanceApplications[rollupID], application)

	fmt.Printf("Governance application for Rollup %s: %s.\n", rollupID, application)
}

// RecordGovernanceMonitoring logs the monitoring of a governance decision in the rollup.
func (ledger *RollupLedger) RecordGovernanceMonitoring(rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.GovernanceMonitoring == nil {
		ledger.GovernanceMonitoring = make(map[string]bool)
	}

	ledger.GovernanceMonitoring[rollupID] = true

	fmt.Printf("Governance monitoring initiated for Rollup %s.\n", rollupID)
}

// RecordContractAddition logs the addition of a contract to the rollup.
func (ledger *RollupLedger) RecordContractAddition(contractID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.Contracts == nil {
		ledger.Contracts = make(map[string]string)
	}

	ledger.Contracts[contractID] = rollupID

	fmt.Printf("Contract %s added to Rollup %s.\n", contractID, rollupID)
}

// RecordResultSync logs the synchronization of a rollup result.
func (ledger *RollupLedger) RecordResultSync(resultID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.ResultSyncs == nil {
		ledger.ResultSyncs = make(map[string]string)
	}

	ledger.ResultSyncs[resultID] = rollupID

	fmt.Printf("Result %s synchronized for Rollup %s.\n", resultID, rollupID)
}

// RecordResultVerification logs the verification of a rollup result.
func (ledger *RollupLedger) RecordResultVerification(resultID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.VerifiedResults == nil {
		ledger.VerifiedResults = make(map[string]time.Time)
	}

	ledger.VerifiedResults[resultID] = time.Now()

	fmt.Printf("Result %s verified for Rollup %s.\n", resultID, rollupID)
}

// RecordTransactionPruning logs the pruning of transactions in a rollup.
func (ledger *RollupLedger) RecordTransactionPruning(rollupID, pruningDetails string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.TransactionPruning == nil {
		ledger.TransactionPruning = make(map[string][]string)
	}

	ledger.TransactionPruning[rollupID] = append(ledger.TransactionPruning[rollupID], pruningDetails)

	fmt.Printf("Transactions pruned for Rollup %s: %s.\n", rollupID, pruningDetails)
}

// RecordBridgeFinalization logs the finalization of a bridge transaction for a rollup.
func (ledger *RollupLedger) RecordBridgeFinalization(bridgeID, rollupID string) {
	ledger.Lock()
	defer ledger.Unlock()

	if ledger.BridgeFinalizations == nil {
		ledger.BridgeFinalizations = make(map[string]string)
	}

	ledger.BridgeFinalizations[bridgeID] = rollupID

	fmt.Printf("Bridge transaction %s finalized for Rollup %s.\n", bridgeID, rollupID)
}

