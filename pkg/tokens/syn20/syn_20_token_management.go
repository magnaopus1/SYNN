package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// TRANSFER_TOKEN transfers SYN20 tokens from the caller’s account to a specified recipient.
func (token *Syn20Token) TRANSFER_TOKEN(from, to string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Metadata.Paused {
        return fmt.Errorf("token transfers are currently paused")
    }

    if token.BalanceSheet[from].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", from)
    }

    token.BalanceSheet[from].Sub(token.BalanceSheet[from], amount)
    token.BalanceSheet[to].Add(token.BalanceSheet[to], amount)

    return token.Ledger.RecordTransaction("TokenTransfer", from, to, amount)
}

// TRANSFER_FROM allows a specified spender to transfer tokens on behalf of an owner.
func (token *Syn20Token) TRANSFER_FROM(owner, spender, to string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance := token.Allowances[owner][spender]
    if allowance.Cmp(amount) < 0 {
        return fmt.Errorf("transfer amount exceeds allowance")
    }

    if token.BalanceSheet[owner].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", owner)
    }

    token.Allowances[owner][spender].Sub(allowance, amount)
    token.BalanceSheet[owner].Sub(token.BalanceSheet[owner], amount)
    token.BalanceSheet[to].Add(token.BalanceSheet[to], amount)

    return token.Ledger.RecordTransaction("TransferFrom", owner, to, amount)
}

// BALANCE_OF returns the balance of the specified account.
func (token *Syn20Token) BALANCE_OF(account string) (*big.Int, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    balance, exists := token.BalanceSheet[account]
    if !exists {
        return big.NewInt(0), fmt.Errorf("account %s not found", account)
    }
    return balance, nil
}

// APPROVE grants allowance for a spender to transfer a specified amount of tokens from the owner’s account.
func (token *Syn20Token) APPROVE(owner, spender string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Allowances[owner] == nil {
        token.Allowances[owner] = make(map[string]*big.Int)
    }
    token.Allowances[owner][spender] = amount

    return token.Ledger.RecordLog("AllowanceSet", fmt.Sprintf("Allowance of %s for spender %s set to %d", owner, spender, amount))
}

// ALLOWANCE returns the remaining number of tokens that the spender can transfer on behalf of the owner.
func (token *Syn20Token) ALLOWANCE(owner, spender string) (*big.Int, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    allowance, exists := token.Allowances[owner][spender]
    if !exists {
        return big.NewInt(0), nil
    }
    return allowance, nil
}

// INCREASE_ALLOWANCE increases the allowance for a spender by a specified amount.
func (token *Syn20Token) INCREASE_ALLOWANCE(owner, spender string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Allowances[owner][spender] == nil {
        token.Allowances[owner][spender] = big.NewInt(0)
    }
    token.Allowances[owner][spender].Add(token.Allowances[owner][spender], amount)

    return token.Ledger.RecordLog("AllowanceIncreased", fmt.Sprintf("Allowance of %s for spender %s increased by %d", owner, spender, amount))
}

// DECREASE_ALLOWANCE decreases the allowance for a spender by a specified amount.
func (token *Syn20Token) DECREASE_ALLOWANCE(owner, spender string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Allowances[owner][spender] == nil || token.Allowances[owner][spender].Cmp(amount) < 0 {
        return fmt.Errorf("decrease amount exceeds allowance")
    }
    token.Allowances[owner][spender].Sub(token.Allowances[owner][spender], amount)

    return token.Ledger.RecordLog("AllowanceDecreased", fmt.Sprintf("Allowance of %s for spender %s decreased by %d", owner, spender, amount))
}

// PAUSE halts all token transfers and approvals.
func (token *Syn20Token) PAUSE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.Paused = true
    return token.Ledger.RecordLog("TokenPaused", "Token operations paused")
}

// UNPAUSE resumes all token transfers and approvals.
func (token *Syn20Token) UNPAUSE() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    token.Metadata.Paused = false
    return token.Ledger.RecordLog("TokenUnpaused", "Token operations resumed")
}

// PERMIT enables off-chain approvals by generating signed approvals for spending.
func (token *Syn20Token) PERMIT(owner, spender string, amount *big.Int, signature string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    valid, err := token.Encryption.VerifySignature(owner, signature)
    if err != nil || !valid {
        return fmt.Errorf("invalid signature")
    }

    return token.APPROVE(owner, spender, amount)
}

// SNAPSHOT creates a snapshot of current balances and allowances for audit and reference.
func (token *Syn20Token) SNAPSHOT() error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    snapshotData, err := token.Encryption.Encrypt(fmt.Sprintf("Snapshot at %v: Balances %v, Allowances %v", time.Now(), token.BalanceSheet, token.Allowances))
    if err != nil {
        return fmt.Errorf("snapshot encryption failed: %v", err)
    }

    return token.Ledger.RecordLog("SnapshotCreated", snapshotData)
}

// DELEGATE allows token holders to delegate their voting power to another account.
func (token *Syn20Token) DELEGATE(delegator, delegatee string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.DelegateVote(delegator, delegatee)
}

// TRANSFER_WITH_MEMO transfers tokens from one account to another with an encrypted memo.
func (token *Syn20Token) TRANSFER_WITH_MEMO(from, to string, amount *big.Int, memo string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.Metadata.Paused {
        return fmt.Errorf("token transfers are currently paused")
    }

    if token.BalanceSheet[from].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", from)
    }

    encryptedMemo, err := token.Encryption.Encrypt(memo)
    if err != nil {
        return fmt.Errorf("failed to encrypt memo: %v", err)
    }

    token.BalanceSheet[from].Sub(token.BalanceSheet[from], amount)
    token.BalanceSheet[to].Add(token.BalanceSheet[to], amount)

    return token.Ledger.RecordTransactionWithMemo("TransferWithMemo", from, to, amount, encryptedMemo)
}
