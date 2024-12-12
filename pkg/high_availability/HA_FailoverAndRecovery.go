package high_availability


import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// haInitiateFailover initiates the failover process.
func haInitiateFailover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateFailover(); err != nil {
        return fmt.Errorf("failed to initiate failover: %v", err)
    }
    fmt.Println("Failover process initiated.")
    return nil
}

// haConfirmFailover confirms the failover process.
func haConfirmFailover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.ConfirmFailover(); err != nil {
        return fmt.Errorf("failed to confirm failover: %v", err)
    }
    fmt.Println("Failover confirmed.")
    return nil
}

// haFetchFailoverStatus retrieves the current status of the failover process.
func haFetchFailoverStatus(ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.HighAvailabilityLedger.FetchFailoverStatus()
    if err != nil {
        return "", fmt.Errorf("failed to fetch failover status: %v", err)
    }
    return status, nil
}

// haSetFailoverThreshold sets the failover threshold.
func haSetFailoverThreshold(threshold int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetFailoverThreshold(threshold); err != nil {
        return fmt.Errorf("failed to set failover threshold: %v", err)
    }
    fmt.Printf("Failover threshold set to %d.\n", threshold)
    return nil
}

// haGetFailoverThreshold retrieves the failover threshold.
func haGetFailoverThreshold(ledgerInstance *ledger.Ledger) (int, error) {
    threshold, err := ledgerInstance.HighAvailabilityLedger.GetFailoverThreshold()
    if err != nil {
        return 0, fmt.Errorf("failed to get failover threshold: %v", err)
    }
    return threshold, nil
}

// haEnableAutoScaling enables auto-scaling for high availability.
func haEnableAutoScaling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableAutoScaling(); err != nil {
        return fmt.Errorf("failed to enable auto-scaling: %v", err)
    }
    fmt.Println("Auto-scaling enabled.")
    return nil
}

// haDisableAutoScaling disables auto-scaling.
func haDisableAutoScaling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableAutoScaling(); err != nil {
        return fmt.Errorf("failed to disable auto-scaling: %v", err)
    }
    fmt.Println("Auto-scaling disabled.")
    return nil
}

// haSetAutoScalingPolicy sets the auto-scaling policy.
func haSetAutoScalingPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetAutoScalingPolicy(policy); err != nil {
        return fmt.Errorf("failed to set auto-scaling policy: %v", err)
    }
    fmt.Println("Auto-scaling policy set.")
    return nil
}

// haGetAutoScalingPolicy retrieves the auto-scaling policy.
func haGetAutoScalingPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetAutoScalingPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get auto-scaling policy: %v", err)
    }
    return policy, nil
}

// haEnableAutoRecovery enables automatic recovery in case of failure.
func haEnableAutoRecovery(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableAutoRecovery(); err != nil {
        return fmt.Errorf("failed to enable auto recovery: %v", err)
    }
    fmt.Println("Auto-recovery enabled.")
    return nil
}

// haDisableAutoRecovery disables automatic recovery.
func haDisableAutoRecovery(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableAutoRecovery(); err != nil {
        return fmt.Errorf("failed to disable auto recovery: %v", err)
    }
    fmt.Println("Auto-recovery disabled.")
    return nil
}

// haSetRecoveryPoint sets a new recovery point.
func haSetRecoveryPoint(pointID, description string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetRecoveryPoint(pointID, description); err != nil {
        return fmt.Errorf("failed to set recovery point: %v", err)
    }
    fmt.Printf("Recovery point %s set.\n", pointID)
    return nil
}

// haRevertToRecoveryPoint reverts to a specified recovery point.
func haRevertToRecoveryPoint(pointID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RevertToRecoveryPoint(pointID); err != nil {
        return fmt.Errorf("failed to revert to recovery point %s: %v", pointID, err)
    }
    fmt.Printf("Reverted to recovery point %s.\n", pointID)
    return nil
}

// haListRecoveryPoints lists all available recovery points.
func haListRecoveryPoints(ledgerInstance *ledger.Ledger) ([]string, error) {
    points, err := ledgerInstance.HighAvailabilityLedger.ListRecoveryPoints()
    if err != nil {
        return nil, fmt.Errorf("failed to list recovery points: %v", err)
    }
    return points, nil
}


// HAInitiateBackup initiates a backup for disaster recovery.
func HAInitiateBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateBackup(backupName); err != nil {
        return fmt.Errorf("failed to initiate backup %s: %v", backupName, err)
    }
    fmt.Printf("Backup %s initiated.\n", backupName)
    return nil
}

// HACompleteBackup marks the backup as complete.
func HACompleteBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.CompleteBackup(backupName); err != nil {
        return fmt.Errorf("failed to complete backup %s: %v", backupName, err)
    }
    fmt.Printf("Backup %s completed.\n", backupName)
    return nil
}

// HARestoreBackup restores from a specified backup.
func HARestoreBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RestoreBackup(backupName); err != nil {
        return fmt.Errorf("failed to restore backup %s: %v", backupName, err)
    }
    fmt.Printf("Backup %s restored.\n", backupName)
    return nil
}

// HAListBackups lists all available backups.
func HAListBackups(ledgerInstance *ledger.Ledger) ([]string, error) {
    backups, err := ledgerInstance.HighAvailabilityLedger.ListBackups()
    if err != nil {
        return nil, fmt.Errorf("failed to list backups: %v", err)
    }
    return backups, nil
}

// HADeleteBackup deletes a specified backup.
func HADeleteBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DeleteBackup(backupName); err != nil {
        return fmt.Errorf("failed to delete backup %s: %v", backupName, err)
    }
    fmt.Printf("Backup %s deleted.\n", backupName)
    return nil
}

// HAMonitorBackupStatus monitors the backup status.
func HAMonitorBackupStatus(ledgerInstance *ledger.Ledger) error {
    status, err := ledgerInstance.HighAvailabilityLedger.MonitorBackupStatus()
    if err != nil {
        return fmt.Errorf("failed to monitor backup status: %v", err)
    }
    fmt.Printf("Backup status: %s.\n", status)
    return nil
}

// HAEnableSnapshot enables the snapshot feature.
func HAEnableSnapshot(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableSnapshot(); err != nil {
        return fmt.Errorf("failed to enable snapshot: %v", err)
    }
    fmt.Println("Snapshot enabled.")
    return nil
}

// HADisableSnapshot disables the snapshot feature.
func HADisableSnapshot(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableSnapshot(); err != nil {
        return fmt.Errorf("failed to disable snapshot: %v", err)
    }
    fmt.Println("Snapshot disabled.")
    return nil
}

// HACreateSnapshot creates a new snapshot of the system state.
func HACreateSnapshot(snapshotName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.CreateSnapshot(snapshotName); err != nil {
        return fmt.Errorf("failed to create snapshot %s: %v", snapshotName, err)
    }
    fmt.Printf("Snapshot %s created.\n", snapshotName)
    return nil
}

// HARestoreSnapshot restores the system state from a snapshot.
func HARestoreSnapshot(snapshotName string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RestoreSnapshot(snapshotName); err != nil {
        return fmt.Errorf("failed to restore snapshot %s: %v", snapshotName, err)
    }
    fmt.Printf("Snapshot %s restored.\n", snapshotName)
    return nil
}
