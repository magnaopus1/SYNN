package high_availability


import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// haSetReplicaCount sets the number of replicas for high availability.
func haSetReplicaCount(count int, ledgerInstance *ledger.Ledger) error {
    if count <= 0 {
        return fmt.Errorf("replica count must be greater than zero")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetReplicaCount(count); err != nil {
        return fmt.Errorf("failed to set replica count: %v", err)
    }
    fmt.Println("Replica count set.")
    return nil
}

// haGetReplicaCount retrieves the current replica count.
func haGetReplicaCount(ledgerInstance *ledger.Ledger) (int, error) {
    count, err := ledgerInstance.HighAvailabilityLedger.GetReplicaCount()
    if err != nil {
        return 0, fmt.Errorf("failed to get replica count: %v", err)
    }
    return count, nil
}

// haAddReadReplica adds a new read replica.
func haAddReadReplica(replicaID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.AddReadReplica(replicaID); err != nil {
        return fmt.Errorf("failed to add read replica: %v", err)
    }
    fmt.Printf("Read replica %s added.\n", replicaID)
    return nil
}

// haRemoveReadReplica removes a read replica.
func haRemoveReadReplica(replicaID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RemoveReadReplica(replicaID); err != nil {
        return fmt.Errorf("failed to remove read replica: %v", err)
    }
    fmt.Printf("Read replica %s removed.\n", replicaID)
    return nil
}

// haListReadReplicas lists all active read replicas.
func haListReadReplicas(ledgerInstance *ledger.Ledger) ([]string, error) {
    replicas, err := ledgerInstance.HighAvailabilityLedger.ListReadReplicas()
    if err != nil {
        return nil, fmt.Errorf("failed to list read replicas: %v", err)
    }
    return replicas, nil
}

// haEnableDataSynchronization enables data synchronization.
func haEnableDataSynchronization(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableDataSynchronization(); err != nil {
        return fmt.Errorf("failed to enable data synchronization: %v", err)
    }
    fmt.Println("Data synchronization enabled.")
    return nil
}

// haDisableDataSynchronization disables data synchronization.
func haDisableDataSynchronization(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableDataSynchronization(); err != nil {
        return fmt.Errorf("failed to disable data synchronization: %v", err)
    }
    fmt.Println("Data synchronization disabled.")
    return nil
}

// haSetSynchronizationInterval sets the interval for data synchronization.
func haSetSynchronizationInterval(interval int, ledgerInstance *ledger.Ledger) error {
    if interval <= 0 {
        return fmt.Errorf("synchronization interval must be greater than zero")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetSynchronizationInterval(interval); err != nil {
        return fmt.Errorf("failed to set synchronization interval: %v", err)
    }
    fmt.Println("Synchronization interval set.")
    return nil
}

// haGetSynchronizationInterval retrieves the synchronization interval.
func haGetSynchronizationInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.HighAvailabilityLedger.GetSynchronizationInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get synchronization interval: %v", err)
    }
    return interval, nil
}

// haEnableDataCompression enables data compression.
func haEnableDataCompression(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableDataCompression(); err != nil {
        return fmt.Errorf("failed to enable data compression: %v", err)
    }
    fmt.Println("Data compression enabled.")
    return nil
}

// haDisableDataCompression disables data compression.
func haDisableDataCompression(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableDataCompression(); err != nil {
        return fmt.Errorf("failed to disable data compression: %v", err)
    }
    fmt.Println("Data compression disabled.")
    return nil
}

// haSetCompressionLevel sets the compression level.
func haSetCompressionLevel(level int, ledgerInstance *ledger.Ledger) error {
    if level < 1 || level > 9 {
        return fmt.Errorf("compression level must be between 1 and 9")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetCompressionLevel(level); err != nil {
        return fmt.Errorf("failed to set compression level: %v", err)
    }
    fmt.Println("Compression level set.")
    return nil
}

// haGetCompressionLevel retrieves the compression level.
func haGetCompressionLevel(ledgerInstance *ledger.Ledger) (int, error) {
    level, err := ledgerInstance.HighAvailabilityLedger.GetCompressionLevel()
    if err != nil {
        return 0, fmt.Errorf("failed to get compression level: %v", err)
    }
    return level, nil
}

// haSetRedundancyLevel sets the redundancy level.
func haSetRedundancyLevel(level int, ledgerInstance *ledger.Ledger) error {
    if level < 1 {
        return fmt.Errorf("redundancy level must be greater than or equal to 1")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetRedundancyLevel(level); err != nil {
        return fmt.Errorf("failed to set redundancy level: %v", err)
    }
    fmt.Println("Redundancy level set.")
    return nil
}

// haGetRedundancyLevel retrieves the redundancy level.
func haGetRedundancyLevel(ledgerInstance *ledger.Ledger) (int, error) {
    level, err := ledgerInstance.HighAvailabilityLedger.GetRedundancyLevel()
    if err != nil {
        return 0, fmt.Errorf("failed to get redundancy level: %v", err)
    }
    return level, nil
}

// haEnableDataDeduplication enables data deduplication.
func haEnableDataDeduplication(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableDataDeduplication(); err != nil {
        return fmt.Errorf("failed to enable data deduplication: %v", err)
    }
    fmt.Println("Data deduplication enabled.")
    return nil
}

// haDisableDataDeduplication disables data deduplication.
func haDisableDataDeduplication(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableDataDeduplication(); err != nil {
        return fmt.Errorf("failed to disable data deduplication: %v", err)
    }
    fmt.Println("Data deduplication disabled.")
    return nil
}

// haSetDeduplicationPolicy sets the deduplication policy.
func haSetDeduplicationPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetDeduplicationPolicy(policy); err != nil {
        return fmt.Errorf("failed to set deduplication policy: %v", err)
    }
    fmt.Println("Deduplication policy set.")
    return nil
}

// haGetDeduplicationPolicy retrieves the deduplication policy.
func haGetDeduplicationPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetDeduplicationPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get deduplication policy: %v", err)
    }
    return policy, nil
}

// haEnableHotStandby enables hot standby mode.
func haEnableHotStandby(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableHotStandby(); err != nil {
        return fmt.Errorf("failed to enable hot standby mode: %v", err)
    }
    fmt.Println("Hot standby mode enabled.")
    return nil
}

// haDisableHotStandby disables hot standby mode.
func haDisableHotStandby(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableHotStandby(); err != nil {
        return fmt.Errorf("failed to disable hot standby mode: %v", err)
    }
    fmt.Println("Hot standby mode disabled.")
    return nil
}

// haSetHotStandbyPolicy sets the policy for hot standby.
func haSetHotStandbyPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetHotStandbyPolicy(policy); err != nil {
        return fmt.Errorf("failed to set hot standby policy: %v", err)
    }
    fmt.Println("Hot standby policy set.")
    return nil
}

// haGetHotStandbyPolicy retrieves the policy for hot standby.
func haGetHotStandbyPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetHotStandbyPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get hot standby policy: %v", err)
    }
    return policy, nil
}

// haEnableColdStandby enables cold standby mode.
func haEnableColdStandby(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableColdStandby(); err != nil {
        return fmt.Errorf("failed to enable cold standby mode: %v", err)
    }
    fmt.Println("Cold standby mode enabled.")
    return nil
}

// haDisableColdStandby disables cold standby mode.
func haDisableColdStandby(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableColdStandby(); err != nil {
        return fmt.Errorf("failed to disable cold standby mode: %v", err)
    }
    fmt.Println("Cold standby mode disabled.")
    return nil
}
