package common

import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
)

// Distribution percentages
const (
	DevelopmentPoolPercentage       = 0.05
	CharityPoolPercentage           = 0.10
	LoanPoolPercentage              = 0.05
	PassiveIncomePoolPercentage     = 0.05
	ValidatorMinerRewardPoolPercentage = 0.70
	AuthorityNodeHostsRewardPoolPercentage = 0.05
)

// TransactionDistributionManager manages the distribution of transaction fees and rewards.
type TransactionDistributionManager struct {
    ledgerInstance *ledger.Ledger
    mutex          sync.Mutex
}

// NewTransactionDistributionManager initializes a new TransactionDistributionManager.
func NewTransactionDistributionManager(ledgerInstance *ledger.Ledger) *TransactionDistributionManager {
	return &TransactionDistributionManager{
		ledgerInstance: ledgerInstance,
	}
}

// DistributeRewards takes the total transaction fee and distributes it across the specified pools.
func (tdm *TransactionDistributionManager) DistributeRewards(blockID string, totalTransactionFee float64) error {
	tdm.mutex.Lock()
	defer tdm.mutex.Unlock()

	// Calculate distribution amounts
	devPoolAmount := totalTransactionFee * DevelopmentPoolPercentage
	charityPoolAmount := totalTransactionFee * CharityPoolPercentage
	loanPoolAmount := totalTransactionFee * LoanPoolPercentage
	passiveIncomePoolAmount := totalTransactionFee * PassiveIncomePoolPercentage
	validatorMinerRewardAmount := totalTransactionFee * ValidatorMinerRewardPoolPercentage
	nodeHostRewardAmount := totalTransactionFee * AuthorityNodeHostsRewardPoolPercentage

	// Distribute to each pool
	err := tdm.distributeToInternalDevelopmentPool(blockID, devPoolAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to development pool: %v", err)
	}

	err = tdm.distributeToCharityPool(blockID, charityPoolAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to charity pool: %v", err)
	}

	err = tdm.distributeToLoanPool(blockID, loanPoolAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to loan pool: %v", err)
	}

	err = tdm.distributeToPassiveIncomePool(blockID, passiveIncomePoolAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to passive income pool: %v", err)
	}

	err = tdm.distributeToValidatorAndMinerRewardPool(blockID, validatorMinerRewardAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to validator and miner reward pool: %v", err)
	}

	err = tdm.distributeToAuthorityNodeHostsPool(blockID, nodeHostRewardAmount)
	if err != nil {
		return fmt.Errorf("failed to distribute to authority node hosts pool: %v", err)
	}

	return nil
}

// distributeToInternalDevelopmentPool distributes rewards to the internal development pool.
func (tdm *TransactionDistributionManager) distributeToInternalDevelopmentPool(blockID string, amount float64) error {
	return tdm.recordPoolTransaction(blockID, "Internal Development Pool", amount)
}

// distributeToCharityPool distributes rewards to the charity pool.
func (tdm *TransactionDistributionManager) distributeToCharityPool(blockID string, amount float64) error {
	return tdm.recordPoolTransaction(blockID, "Charity Pool", amount)
}

// distributeToLoanPool distributes rewards to the loan pool.
func (tdm *TransactionDistributionManager) distributeToLoanPool(blockID string, amount float64) error {
	return tdm.recordPoolTransaction(blockID, "Loan Pool", amount)
}

// distributeToPassiveIncomePool distributes rewards to the passive income pool, which is distributed to verified wallets every 14 days.
func (tdm *TransactionDistributionManager) distributeToPassiveIncomePool(blockID string, amount float64) error {
    // Ensure SYN900TokenID is defined in the correct scope (use as an argument or calculate it)
    SYN900TokenID := "your-created-token-id" // Replace this with the actual token ID or pass it as an argument

    // Retrieve the list of verified wallets using the correct token ID
    verifiedWallets, err := tdm.ledgerInstance.AccountsWalletLedger.GetVerifiedWallets(SYN900TokenID)
    if err != nil {
        return fmt.Errorf("failed to retrieve verified wallets: %v", err)
    }

    // Calculate the share for each wallet
    share := amount / float64(len(verifiedWallets))

    // Loop through the verified wallets and distribute the share
    for _, walletData := range verifiedWallets {
        // Convert ledger.WalletData to Wallet
        wallet := convertWalletDataToWallet(walletData)
        
        err := tdm.recordWalletTransaction(wallet, "Passive Income Reward", share)
        if err != nil {
            return fmt.Errorf("failed to distribute passive income to wallet %s: %v", wallet.Address, err) // Use the correct field for the wallet ID
        }
    }
    return nil
}

// Convert ledger.WalletData to Wallet (if necessary for other use cases)
func convertWalletDataToWallet(walletData ledger.WalletData) Wallet {
    return Wallet{
        Address: walletData.OwnerAddress, // Use OwnerAddress as Address
        // No conversion for PrivateKey or PublicKey, since WalletData doesn't contain this information
        Ledger: nil, // You may need to handle the Ledger initialization here
    }
}





// distributeToValidatorAndMinerRewardPool distributes rewards to the validator and miner reward pool.
func (tdm *TransactionDistributionManager) distributeToValidatorAndMinerRewardPool(blockID string, amount float64) error {
	return tdm.recordPoolTransaction(blockID, "Validator and Miner Reward Pool", amount)
}

// distributeToAuthorityNodeHostsPool distributes rewards to the authority node hosts pool every 7 days.
func (tdm *TransactionDistributionManager) distributeToAuthorityNodeHostsPool(blockID string, amount float64) error {
    // Retrieve the list of active authority nodes from the ledger
    activeLedgerNodes, err := tdm.ledgerInstance.GetActiveAuthorityNodes()
    if err != nil {
        return fmt.Errorf("failed to retrieve active authority nodes: %v", err)
    }

    // Calculate the share for each active node
    share := amount / float64(len(activeLedgerNodes))

    // Loop through and distribute rewards to each active node
    for _, ledgerNode := range activeLedgerNodes {
        // Convert ledger.Node to Node
        node := convertLedgerNodeToNode(ledgerNode)

        // Record the transaction for each node
        err := tdm.recordNodeTransaction(node, "Authority Node Host Reward", share)
        if err != nil {
            return fmt.Errorf("failed to distribute node reward to node %s: %v", node.Address, err)
        }
    }
    return nil
}

// convertLedgerNodeToNode converts a ledger.Node to a Node type.
func convertLedgerNodeToNode(ledgerNode ledger.Node) Node {
    return Node{
        Address:      ledgerNode.Address,
        Name:         ledgerNode.Name,
        NodeCategory: convertLedgerNodeCategory(ledgerNode.NodeCategory), // Convert NodeCategory
        NodeType:     convertLedgerNodeType(ledgerNode.NodeType),         // Convert NodeType
        NodeKey:      convertLedgerNodeKey(ledgerNode.NodeKey),           // Convert NodeKey
        IsActive:     ledgerNode.IsActive,                                // Assuming IsActive exists in ledger.Node
    }
}


// convertLedgerNodeCategory converts ledger.NodeCategory to NodeCategory.
func convertLedgerNodeCategory(ledgerCategory ledger.NodeCategory) NodeCategory {
    // Assuming both are string-based types, you can directly cast them.
    return NodeCategory(ledgerCategory)  // Cast ledger.NodeCategory to NodeCategory
}

// convertLedgerNodeType converts ledger.NodeType to NodeType.
func convertLedgerNodeType(ledgerType ledger.NodeType) NodeType {
    // Assuming both are string-based types, you can directly cast them.
    return NodeType(ledgerType)  // Cast ledger.NodeType to NodeType
}

// convertLedgerNodeKey converts *ledger.NodeKey to *NodeKey.
func convertLedgerNodeKey(ledgerKey *ledger.NodeKey) *NodeKey {
    // Assuming both are pointers to similar structures, you can directly cast them.
    return (*NodeKey)(ledgerKey)  // Cast ledger.NodeKey to NodeKey
}




// recordPoolTransaction logs the distribution transaction into the ledger.
func (tdm *TransactionDistributionManager) recordPoolTransaction(blockID, poolName string, amount float64) error {
    // Create an encryption instance
    encryptionInstance := &Encryption{}

    // Encrypt the amount using AES algorithm and the predefined encryption key
    encryptedAmount, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%f", amount)), EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt pool transaction: %v", err)
    }

    // Store the encrypted amount in the EncryptedData field
    transaction := ledger.Transaction{
        BlockID:      blockID,
        EncryptedData: string(encryptedAmount), // Store encrypted amount here
        // Other fields that exist in ledger.Transaction should be populated here as needed
    }

    // Record the transaction in the ledger
    return tdm.ledgerInstance.BlockchainConsensusCoinLedger.RecordPoolTransaction(poolName, transaction) // Use poolName as the first argument
}



// recordWalletTransaction records a transaction for a specific wallet in the ledger.
func (tdm *TransactionDistributionManager) recordWalletTransaction(wallet Wallet, description string, amount float64) error {
    // Create an encryption instance
    encryptionInstance := &Encryption{}

    // Encrypt the amount using AES algorithm and the predefined encryption key
    encryptedAmount, err := encryptionInstance.EncryptData("AES", []byte(fmt.Sprintf("%f", amount)), EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt wallet transaction: %v", err)
    }

    // Store the encrypted amount in the EncryptedData field
    transaction := ledger.Transaction{
        FromAddress:  wallet.Address,      // Store wallet address in FromAddress
        EncryptedData: string(encryptedAmount), // Store encrypted amount here
        // Removed the Description field since it's not part of ledger.Transaction
    }

    // Record the transaction in the ledger with the wallet address and the transaction object
    return tdm.ledgerInstance.BlockchainConsensusCoinLedger.RecordWalletTransaction(wallet.Address, transaction)
}





// recordNodeTransaction records a transaction for a specific node host in the ledger.
func (tdm *TransactionDistributionManager) recordNodeTransaction(node Node, description string, amount float64) error {
    // Use the general RecordTransaction method with the node address, description, and amount
    return tdm.ledgerInstance.BlockchainConsensusCoinLedger.RecordTransaction(node.Address, description, amount)
}


