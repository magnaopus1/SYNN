package syn20

import (
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// RECLAIM_TOKENS allows the token issuer to reclaim tokens from inactive accounts.
func (token *Syn20Token) RECLAIM_TOKENS(account string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.BalanceSheet[account].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", account)
    }

    token.BalanceSheet[account].Sub(token.BalanceSheet[account], amount)
    token.BalanceSheet[token.Metadata.TokenOwner].Add(token.BalanceSheet[token.Metadata.TokenOwner], amount)

    return token.Ledger.RecordTransaction("TokensReclaimed", account, token.Metadata.TokenOwner, amount)
}

// BULK_APPROVE grants approval for multiple spenders to transfer tokens on behalf of an owner.
func (token *Syn20Token) BULK_APPROVE(owner string, approvals map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for spender, amount := range approvals {
        if token.Allowances[owner] == nil {
            token.Allowances[owner] = make(map[string]*big.Int)
        }
        token.Allowances[owner][spender] = amount
    }
    return token.Ledger.RecordLog("BulkApproval", fmt.Sprintf("Bulk approval set by %s", owner))
}

// BULK_TRANSFER transfers tokens from the owner's account to multiple recipients.
func (token *Syn20Token) BULK_TRANSFER(from string, transfers map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for to, amount := range transfers {
        if token.BalanceSheet[from].Cmp(amount) < 0 {
            return fmt.Errorf("insufficient balance in account %s for transfer to %s", from, to)
        }

        token.BalanceSheet[from].Sub(token.BalanceSheet[from], amount)
        token.BalanceSheet[to].Add(token.BalanceSheet[to], amount)
    }
    return token.Ledger.RecordLog("BulkTransfer", fmt.Sprintf("Bulk transfer executed by %s", from))
}

// TAXATION deducts a specified tax amount from each transaction and transfers it to the tax account.
func (token *Syn20Token) TAXATION(from, to string, amount *big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    taxRate := token.Metadata.TaxRate
    taxAmount := new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(int64(taxRate*100))), big.NewInt(10000))
    netAmount := new(big.Int).Sub(amount, taxAmount)

    token.BalanceSheet[from].Sub(token.BalanceSheet[from], amount)
    token.BalanceSheet[to].Add(token.BalanceSheet[to], netAmount)
    token.BalanceSheet[token.Metadata.TaxAccount].Add(token.BalanceSheet[token.Metadata.TaxAccount], taxAmount)

    return token.Ledger.RecordTransaction("TaxApplied", from, to, netAmount)
}

// STAKING locks tokens in the staking contract for rewards over time.
func (token *Syn20Token) STAKING(account string, amount *big.Int, duration time.Duration) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.BalanceSheet[account].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", account)
    }

    token.BalanceSheet[account].Sub(token.BalanceSheet[account], amount)
    err := token.Ledger.AddStake(account, amount, duration)
    if err != nil {
        return fmt.Errorf("failed to stake tokens: %v", err)
    }

    return token.Ledger.RecordLog("TokensStaked", fmt.Sprintf("Tokens staked by %s for %v", account, duration))
}

// VESTING locks tokens in a vesting contract to be released over a defined schedule.
func (token *Syn20Token) VESTING(account string, amount *big.Int, vestingSchedule map[time.Time]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    if token.BalanceSheet[account].Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance in account %s", account)
    }

    token.BalanceSheet[account].Sub(token.BalanceSheet[account], amount)
    return token.Ledger.AddVesting(account, amount, vestingSchedule)
}

// REWARD_DISTRIBUTION distributes rewards to multiple accounts based on their balances or stakes.
func (token *Syn20Token) REWARD_DISTRIBUTION(rewards map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for account, reward := range rewards {
        token.BalanceSheet[account].Add(token.BalanceSheet[account], reward)
    }
    return token.Ledger.RecordLog("RewardDistributed", "Rewards distributed to eligible accounts")
}

// CHECK_FREEZE_STATUS checks if a specified account is frozen.
func (token *Syn20Token) CHECK_FREEZE_STATUS(account string) (bool, error) {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    return token.Ledger.IsAccountFrozen(account)
}

// FREEZE_ACCOUNT freezes a specified account, preventing all transactions.
func (token *Syn20Token) FREEZE_ACCOUNT(account string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.SetAccountStatus(account, true)
    if err != nil {
        return fmt.Errorf("failed to freeze account %s: %v", account, err)
    }

    return token.Ledger.RecordLog("AccountFrozen", fmt.Sprintf("Account %s frozen", account))
}

// THAW_ACCOUNT unfreezes a specified account, allowing transactions.
func (token *Syn20Token) THAW_ACCOUNT(account string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.SetAccountStatus(account, false)
    if err != nil {
        return fmt.Errorf("failed to thaw account %s: %v", account, err)
    }

    return token.Ledger.RecordLog("AccountThawed", fmt.Sprintf("Account %s thawed", account))
}

// BATCH_TRANSFER performs a batch transfer of tokens to multiple accounts.
func (token *Syn20Token) BATCH_TRANSFER(from string, transfers map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for to, amount := range transfers {
        if token.BalanceSheet[from].Cmp(amount) < 0 {
            return fmt.Errorf("insufficient balance in account %s for transfer to %s", from, to)
        }
        token.BalanceSheet[from].Sub(token.BalanceSheet[from], amount)
        token.BalanceSheet[to].Add(token.BalanceSheet[to], amount)
    }
    return token.Ledger.RecordLog("BatchTransfer", fmt.Sprintf("Batch transfer from %s", from))
}

// BATCH_APPROVE grants multiple approvals for various spenders on behalf of an owner.
func (token *Syn20Token) BATCH_APPROVE(owner string, approvals map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for spender, amount := range approvals {
        if token.Allowances[owner] == nil {
            token.Allowances[owner] = make(map[string]*big.Int)
        }
        token.Allowances[owner][spender] = amount
    }
    return token.Ledger.RecordLog("BatchApproval", fmt.Sprintf("Batch approval set by %s", owner))
}

// BATCH_MINT mints tokens to multiple accounts in a single operation.
func (token *Syn20Token) BATCH_MINT(mints map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for account, amount := range mints {
        token.BalanceSheet[account].Add(token.BalanceSheet[account], amount)
        token.Metadata.TotalSupply.Add(token.Metadata.TotalSupply, amount)
    }
    return token.Ledger.RecordLog("BatchMint", "Tokens minted to multiple accounts")
}

// BATCH_BURN burns tokens from multiple accounts in a single operation.
func (token *Syn20Token) BATCH_BURN(burns map[string]*big.Int) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    for account, amount := range burns {
        if token.BalanceSheet[account].Cmp(amount) < 0 {
            return fmt.Errorf("insufficient balance in account %s", account)
        }
        token.BalanceSheet[account].Sub(token.BalanceSheet[account], amount)
        token.Metadata.TotalSupply.Sub(token.Metadata.TotalSupply, amount)
    }
    return token.Ledger.RecordLog("BatchBurn", "Tokens burned from multiple accounts")
}

// UPGRADE upgrades the token to a new version, handling migration of balances and allowances.
func (token *Syn20Token) UPGRADE(newVersion string) error {
    token.mutex.Lock()
    defer token.mutex.Unlock()

    err := token.Ledger.MigrateToken(token.TokenName, newVersion)
    if err != nil {
        return fmt.Errorf("failed to upgrade token to version %s: %v", newVersion, err)
    }
    return token.Ledger.RecordLog("TokenUpgraded", fmt.Sprintf("Token upgraded to version %s", newVersion))
}
