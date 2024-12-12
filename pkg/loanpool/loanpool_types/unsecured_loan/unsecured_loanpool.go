package loanpool

import (
	"math/big"

	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/synnergy_consensus"
)

// NewUnsecuredLoanPool initializes the Unsecured Loan Pool with the provided ledger, consensus engine, and encryption service.
func NewUnsecuredLoanPool(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption) *common.UnsecuredLoanPool {
	return &common.UnsecuredLoanPool{
		TotalBalance:      big.NewInt(0),
		LoansDistributed:  big.NewInt(0),
		LoansRepaid:       big.NewInt(0),
		LoansDefaulted:    big.NewInt(0),
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		EncryptionService: encryptionService,
		LoanRecords:       make(map[string]*common.LoanRecord),
	}
}

// ViewFundBalance returns the current balance available in the unsecured loan pool.
func (pool *common.UnsecuredLoanPool) ViewFundBalance() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.TotalBalance)
}

// ViewLoansDistributed returns the total amount of loans that have been distributed from the unsecured loan pool.
func (pool *common.UnsecuredLoanPool) ViewLoansDistributed() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansDistributed)
}

// ViewLoansRepaid returns the total amount of loans that have been repaid to the unsecured loan pool.
func (pool *common.UnsecuredLoanPool) ViewLoansRepaid() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansRepaid)
}

// ViewLoansDefaulted returns the total amount of loans that have defaulted in the unsecured loan pool.
func (pool *common.UnsecuredLoanPool) ViewLoansDefaulted() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansDefaulted)
}
