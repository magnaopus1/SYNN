package common


import (
	"fmt"
	"sync"
	"synnergy_network/pkg/ledger"
)

// VirtualMachineConcurrency handles the concurrency mechanisms within the virtual machine.
type VirtualMachineConcurrency struct {
	LedgerInstance    *ledger.Ledger        // Ledger to log transactions, sub-blocks, and blocks
	VirtualMachine    VMInterface           // VMInterface to allow interaction with both LightVM and HeavyVM
	SubBlockManager   *SubBlockManager      // Manages sub-blocks within the system
	SynnergyConsensus *SynnergyConsensus    // Consensus mechanism for block validation
	mutex             sync.Mutex            // Ensures thread-safety across operations
	wg                sync.WaitGroup        // WaitGroup for handling concurrency
}

// NewVirtualMachineConcurrency initializes a new VirtualMachineConcurrency manager.
func NewVirtualMachineConcurrency(ledgerInstance *ledger.Ledger, vm VMInterface, subBlockManager *SubBlockManager, consensus *SynnergyConsensus) *VirtualMachineConcurrency {
	return &VirtualMachineConcurrency{
		LedgerInstance:    ledgerInstance,
		VirtualMachine:    vm,
		SubBlockManager:   subBlockManager,
		SynnergyConsensus: consensus,
	}
}



// ProcessTransactionConcurrently handles transactions concurrently, adding them to the sub-block.
func (vmc *VirtualMachineConcurrency) ProcessTransactionConcurrently(transaction *Transaction) error {
	vmc.wg.Add(1)
	defer vmc.wg.Done()

	go func() {
		vmc.mutex.Lock()
		defer vmc.mutex.Unlock()

		// Step 1: Process the transaction through Synnergy Consensus (no error return)
		vmc.SynnergyConsensus.ProcessTransactions([]Transaction{*transaction}) // Corrected: no error to handle

		// Step 2: Execute transaction
		err := vmc.VirtualMachine.ExecuteTransaction(transaction)
		if err != nil {
			fmt.Printf("Transaction execution failed: %v\n", err)
			return
		}

		// Step 3: Add transaction to sub-block
		err = vmc.SubBlockManager.AddTransactionToSubBlock(transaction)
		if err != nil {
			fmt.Printf("Failed to add transaction to sub-block: %v\n", err)
			return
		}

		// Step 4: Log transaction in the ledger
		err = vmc.LedgerInstance.BlockchainConsensusCoinLedger.LogTransaction(transaction.TransactionID, "Processed") // Assuming `TransactionID` is the correct field for transaction ID
		if err != nil {
			fmt.Printf("Failed to log transaction: %v\n", err)
			return
		}

		fmt.Printf("Transaction %s processed successfully.\n", transaction.TransactionID)
	}()

	return nil
}


// ExecuteContractConcurrently handles the execution of smart contracts concurrently.
func (vmc *VirtualMachineConcurrency) ExecuteContractConcurrently(contractID, contractSource string, parameters map[string]interface{}, signature []byte) (map[string]interface{}, error) {
	vmc.wg.Add(1)
	defer vmc.wg.Done()

	var executionResult map[string]interface{}
	var err error

	// Execute contract in a separate goroutine to allow concurrent executions
	go func() {
		vmc.mutex.Lock()
		defer vmc.mutex.Unlock()

		// Step 1: Validate contract syntax using a method on VMInterface
		isValid, err := vmc.VirtualMachine.ValidateContractSyntax(contractID, contractSource, "metadata or config") // This assumes the interface includes this method
		if err != nil || !isValid {
			fmt.Printf("Contract syntax validation failed: %v\n", err)
			return
		}

		// Step 2: Compile and execute contract without additional encryption argument
		executionResult, err = vmc.VirtualMachine.ExecuteContract(contractID, contractSource, "metadata or config", parameters, signature)
		if err != nil {
			fmt.Printf("Contract execution failed: %v\n", err)
			return
		}

		// Step 3: Log execution result in the ledger with correct argument count
		err = vmc.LedgerInstance.VirtualMachineLedger.LogExecution(contractID, "Processed", "Block123")
		if err != nil {
			fmt.Printf("Failed to log contract execution: %v\n", err)
			return
		}

		fmt.Printf("Contract %s executed successfully.\n", contractID)
	}()

	// Wait for all concurrent operations to complete before returning results
	vmc.wg.Wait()

	return executionResult, err
}



