package loanpool

import (
	"math/big"
	"sync"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/consensus"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
)

// NewSecuredLoanPool initializes the Secured Loan Pool with the provided ledger, consensus engine, and encryption service.
func NewSecuredLoanPool(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *encryption.Encryption) *common.SecuredLoanPool {
	return &common.SecuredLoanPool{
		TotalBalance:     big.NewInt(0),
		LoansDistributed: big.NewInt(0),
		LoansRepaid:      big.NewInt(0),
		LoansDefaulted:   big.NewInt(0),
		Ledger:           ledgerInstance,
		Consensus:        consensusEngine,
		EncryptionService: encryptionService,
		LoanRecords:      make(map[string]*common.LoanRecord),
	}
}

// ViewFundBalance returns the current balance available in the secured loan pool.
func (pool *common.SecuredLoanPool) ViewFundBalance() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.TotalBalance)
}

// ViewLoansDistributed returns the total amount of loans that have been distributed from the loan pool.
func (pool *common.SecuredLoanPool) ViewLoansDistributed() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansDistributed)
}

// ViewLoansRepaid returns the total amount of loans that have been repaid to the loan pool.
func (pool *common.SecuredLoanPool) ViewLoansRepaid() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansRepaid)
}

// ViewLoansDefaulted returns the total amount of loans that have defaulted.
func (pool *common.SecuredLoanPool) ViewLoansDefaulted() *big.Int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return new(big.Int).Set(pool.LoansDefaulted)
}
