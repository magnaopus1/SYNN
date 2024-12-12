package sidechains

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
	"synnergy_network/pkg/consensus"
)

// NewSidechainInteroperability initializes a new interoperability manager
func NewSidechainInteroperability(mainChain *network.MainChain, ledgerInstance *ledger.Ledger, encryptionService *encryption.Encryption, consensus *common.SynnergyConsensus) *common.SidechainInteroperability {
	return &common.SidechainInteroperability{
		MainChain:         mainChain,
		SidechainNetworks: make(map[string]*common.SidechainNetwork),
		Ledger:            ledgerInstance,
		Encryption:        encryptionService,
		Consensus:         consensus,
	}
}

// RegisterSidechain registers a new sidechain for interoperability
func (si *common.SidechainInteroperability) RegisterSidechain(chainID string, sidechainNetwork *common.SidechainNetwork) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	if _, exists := si.SidechainNetworks[chainID]; exists {
		return errors.New("sidechain already registered for interoperability")
	}

	si.SidechainNetworks[chainID] = sidechainNetwork

	// Log the sidechain registration in the ledger
	err := si.Ledger.RecordSidechainRegistration(chainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log sidechain registration: %v", err)
	}

	fmt.Printf("Sidechain %s registered for interoperability\n", chainID)
	return nil
}

// TransferAsset transfers assets between sidechains or between sidechain and mainchain
func (si *common.SidechainInteroperability) TransferAsset(sourceChainID, destinationChainID, txID string, amount float64) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	sourceNetwork, exists := si.SidechainNetworks[sourceChainID]
	if !exists {
		return fmt.Errorf("source sidechain %s not found", sourceChainID)
	}

	// Retrieve the source transaction
	tx, err := sourceNetwork.RetrieveTransaction(txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction from source chain: %v", err)
	}

	// Encrypt the transaction data before transferring
	encryptedTxID, err := si.Encryption.EncryptData([]byte(txID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt transaction data: %v", err)
	}

	if destinationChainID == "mainchain" {
		// Transfer to the mainchain
		err := si.MainChain.TransferData(sourceNetwork.Nodes[tx.SourceNodeID], si.MainChain.Nodes["mainnode"], encryptedTxID)
		if err != nil {
			return fmt.Errorf("failed to transfer data to mainchain: %v", err)
		}

		// Log the transfer in the ledger
		err = si.Ledger.RecordAssetTransfer(sourceChainID, "mainchain", txID, amount, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log asset transfer: %v", err)
		}
		fmt.Printf("Asset transfer of %f from sidechain %s to mainchain completed\n", amount, sourceChainID)

	} else {
		// Transfer between sidechains
		destNetwork, exists := si.SidechainNetworks[destinationChainID]
		if !exists {
			return fmt.Errorf("destination sidechain %s not found", destinationChainID)
		}

		err := si.Consensus.ValidateTransaction(txID, amount)
		if err != nil {
			return fmt.Errorf("transaction validation failed: %v", err)
		}

		// Perform the transfer
		err = sourceNetwork.NetworkManager.TransferData(sourceNetwork.Nodes[tx.SourceNodeID], destNetwork.Nodes[tx.DestinationNodeID], encryptedTxID)
		if err != nil {
			return fmt.Errorf("failed to transfer data between sidechains: %v", err)
		}

		// Log the transfer in the ledger
		err = si.Ledger.RecordAssetTransfer(sourceChainID, destinationChainID, txID, amount, time.Now())
		if err != nil {
			return fmt.Errorf("failed to log asset transfer: %v", err)
		}

		fmt.Printf("Asset transfer of %f from sidechain %s to sidechain %s completed\n", amount, sourceChainID, destinationChainID)
	}

	return nil
}

// ValidateCrossChainTransaction validates a transaction across chains (sidechains or mainchain)
func (si *common.SidechainInteroperability) ValidateCrossChainTransaction(chainID, txID string) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	network, exists := si.SidechainNetworks[chainID]
	if !exists {
		return fmt.Errorf("sidechain network %s not found", chainID)
	}

	// Validate the transaction using Synnergy Consensus
	err := si.Consensus.ValidateTransaction(txID, 0) // Assuming amount is handled elsewhere
	if err != nil {
		return fmt.Errorf("transaction validation failed: %v", err)
	}

	// Log the validation in the ledger
	err = si.Ledger.RecordTransactionValidation(txID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log cross-chain transaction validation: %v", err)
	}

	fmt.Printf("Cross-chain transaction %s validated on sidechain %s\n", txID, chainID)
	return nil
}

// SyncBlockAcrossChains syncs a block from one chain to another
func (si *common.SidechainInteroperability) SyncBlockAcrossChains(sourceChainID, destinationChainID, blockID string) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	sourceNetwork, exists := si.SidechainNetworks[sourceChainID]
	if !exists {
		return fmt.Errorf("source sidechain %s not found", sourceChainID)
	}

	// Retrieve the block from the source sidechain
	block, err := sourceNetwork.RetrieveBlock(blockID)
	if err != nil {
		return fmt.Errorf("failed to retrieve block %s from source sidechain: %v", blockID, err)
	}

	// Encrypt the block before syncing
	encryptedBlockID, err := si.Encryption.EncryptData([]byte(blockID), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt block: %v", err)
	}

	if destinationChainID == "mainchain" {
		// Sync to the mainchain
		err = si.MainChain.SyncBlock(encryptedBlockID, blockID)
		if err != nil {
			return fmt.Errorf("failed to sync block %s to mainchain: %v", blockID, err)
		}
	} else {
		// Sync to another sidechain
		destNetwork, exists := si.SidechainNetworks[destinationChainID]
		if !exists {
			return fmt.Errorf("destination sidechain %s not found", destinationChainID)
		}

		err = sourceNetwork.NetworkManager.TransferData(sourceNetwork.Nodes["sourceNode"], destNetwork.Nodes["destinationNode"], encryptedBlockID)
		if err != nil {
			return fmt.Errorf("failed to sync block %s to destination sidechain: %v", blockID, err)
		}
	}

	// Log the sync operation in the ledger
	err = si.Ledger.RecordBlockSync(blockID, sourceChainID, destinationChainID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log block sync: %v", err)
	}

	fmt.Printf("Block %s synced from sidechain %s to %s\n", blockID, sourceChainID, destinationChainID)
	return nil
}
