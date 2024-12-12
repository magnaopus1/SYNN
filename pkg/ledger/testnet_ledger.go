package ledger

import (
	"fmt"
	"time"
)

// RecordContractDeployment logs the deployment of a smart contract in the testnet.
func (l *TestnetLedger) RecordTestnetContractDeployment(contractID, deployer, contractCode string) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the contract deployment record
	deployment := TestnetContractDeployment{
		ContractID:   contractID,
		Deployer:     deployer,
		DeployedAt:   time.Now(),
		ContractCode: contractCode,
	}

	l.ContractDeployments[contractID] = deployment
	return nil
}

// RecordTNContractExecution logs the execution of a smart contract in the testnet.
func (l *TestnetLedger) RecordTNContractExecution(contractID, executor, inputData string) error {
	l.Lock()
	defer l.Unlock()

	// Initialize the map if it's nil
	if l.ContractExecutions == nil {
		l.ContractExecutions = make(map[string][]TestnetContractExecution)
	}

	// Create the execution record
	execution := TestnetContractExecution{
		ContractID: contractID,
		Executor:   executor,
		InputData:  inputData,
		ExecutedAt: time.Now(),
		Status:     "success", // Default status; could be dynamic based on actual execution
	}

	// Append the execution record to the list for the contract
	l.ContractExecutions[contractID] = append(l.ContractExecutions[contractID], execution)

	// Optionally log the execution event
	fmt.Printf("Contract execution recorded: ContractID %s executed by %s\n", contractID, executor)

	return nil
}

// RecordTestnetFaucetClaim logs a faucet claim in the testnet.
func (l *TestnetLedger) RecordTestnetFaucetClaim(claimID, claimer string, amount uint64) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the faucet claim record
	claim := TestnetFaucetClaim{
		ClaimID:   claimID,
		Claimer:   claimer,
		ClaimedAt: time.Now(),
		Amount:    amount,
	}

	// Add the claim to the ledger
	l.TestnetFaucetClaims[claimID] = claim
	return nil
}


// RecordSubBlock logs a sub-block in the testnet.
func (l *TestnetLedger) RecordSubBlock(subBlockID string, subBlockData []Transaction, validator string, prevHash string, pohProof PoHProof) error {
	l.Lock()
	defer l.Unlock()

	// Create and log the sub-block details
	subBlock := SubBlock{
		SubBlockID:   subBlockID,
		Index:        len(l.SubBlocks) + 1,  // Sub-block index is based on length
		Timestamp:    time.Now(),
		Transactions: subBlockData,
		Validator:    validator,
		PrevHash:     prevHash,
		PoHProof:     pohProof,
		Status:       "valid",               // Status is "valid" for logged sub-blocks
	}

	// Append the sub-block to the list of sub-blocks in the testnet
	l.SubBlocks = append(l.SubBlocks, subBlock)
	return nil
}

// RecordBlock logs a finalized block in the testnet.
func (l *TestnetLedger) RecordBlock(blockID string, subBlocks []SubBlock, prevHash string, hash string, nonce int, difficulty int, minerReward float64, validators []string) error {
	l.Lock()
	defer l.Unlock()

	// Create and log the block details based on the sub-blocks
	block := Block{
		BlockID:     blockID,
		Index:       len(l.Blocks) + 1,     // Block index is based on length
		Timestamp:   time.Now(),
		SubBlocks:   subBlocks,             // Sub-blocks associated with the block
		PrevHash:    prevHash,
		Hash:        hash,
		Nonce:       nonce,
		Difficulty:  difficulty,
		MinerReward: minerReward,
		Validators:  validators,
		Status:      "finalized",           // Block is "finalized" once logged
	}

	// Append the block to the list of blocks in the testnet
	l.Blocks = append(l.Blocks, block)
	return nil
}


// RecordTransaction logs a transaction in the testnet.
func (l *TestnetLedger) RecordTransaction(txID string, txData Transaction) error {
	l.Lock()
	defer l.Unlock()

	// Store the transaction record in the transaction cache
	l.TransactionCache[txID] = txData
	return nil
}

// RecordTokenDeployment logs the deployment of a token in the testnet.
func (l *TestnetLedger) RecordTokenDeployment(tokenID, deployer, tokenSymbol string, initialSupply uint64) error {
	l.Lock()
	defer l.Unlock()

	// Create and store the token deployment record
	deployment := TokenDeployment{
		TokenID:      tokenID,
		Deployer:     deployer,
		TokenSymbol:  tokenSymbol,
		InitialSupply: initialSupply,
		DeployedAt:   time.Now(),
	}

	l.TokenDeployments[tokenID] = deployment
	return nil
}
