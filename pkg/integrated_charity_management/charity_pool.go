package integrated_charity_management

import (
    "fmt"
    "time"
)


// NewCharityPool initializes the CharityPool with a starting total balance
func NewCharityPool(initialBalance float64) *CharityPool {
    return &CharityPool{
        externalPool:  0.0,           // Initialize with 0 in the external charity pool
        internalPool:  0.0,           // Initialize with 0 in the internal charity pool
        totalBalance:  initialBalance, // Set the initial total balance
    }
}

// ViewPools provides a view of the current balances in the external and internal charity pools
func (cp *CharityPool) ViewPools() {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    fmt.Printf("External Charity Pool Balance: %.2f\n", cp.externalPool)
    fmt.Printf("Internal Charity Pool Balance: %.2f\n", cp.internalPool)
    fmt.Printf("Total Balance: %.2f\n", cp.totalBalance)
}

// DistributeFunds distributes the total balance 50:50 between the external and internal charity pools every 24 hours
func (cp *CharityPool) DistributeFunds() {
    cp.mutex.Lock()
    defer cp.mutex.Unlock()

    // Calculate 50:50 distribution
    distributionAmount := cp.totalBalance / 2
    cp.externalPool += distributionAmount
    cp.internalPool += distributionAmount

    // After distribution, reset totalBalance to 0 (since it’s fully distributed)
    cp.totalBalance = 0

    fmt.Printf("Distributed %.2f SYNN to both external and internal charity pools.\n", distributionAmount)
    fmt.Printf("External Charity Pool Balance: %.2f\n", cp.externalPool)
    fmt.Printf("Internal Charity Pool Balance: %.2f\n", cp.internalPool)
}

// StartDistributionTimer starts the process of distributing funds every 24 hours
func (cp *CharityPool) StartDistributionTimer() {
    ticker := time.NewTicker(24 * time.Hour)
    go func() {
        for {
            <-ticker.C
            cp.DistributeFunds()
        }
    }()
}

