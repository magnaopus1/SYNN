package loanpool

import (
	"math/big"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewEducationFund initializes the Education Fund with the provided ledger and consensus engine.
func NewEducationFund(ledgerInstance *ledger.Ledger, consensusEngine *common.SynnergyConsensus, encryptionService *common.Encryption) *EducationFund {
	return &EducationFund{
		TotalBalance:      big.NewInt(0),
		GrantsDistributed: big.NewInt(0),
		Ledger:            ledgerInstance,
		Consensus:         consensusEngine,
		EncryptionService: encryptionService,
	}
}

// ViewFundBalance returns the current balance available in the education fund.
func (fund *EducationFund) ViewFundBalance() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.TotalBalance)
}

// ViewGrantsDistributed returns the total amount of grants that have been distributed.
func (fund *EducationFund) ViewGrantsDistributed() *big.Int {
	fund.mutex.Lock()
	defer fund.mutex.Unlock()

	return new(big.Int).Set(fund.GrantsDistributed)
}




