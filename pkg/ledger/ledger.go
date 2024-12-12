package ledger

import (
)

// GetLedgerInstance returns the instance of the current ledger.
func (l *Ledger) GetLedgerInstance() *Ledger {
	l.Lock()
	defer l.Unlock()

	return l
}
