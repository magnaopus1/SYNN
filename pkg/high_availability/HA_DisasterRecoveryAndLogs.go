package high_availability

import (
	"fmt"
	"log"
	"synnergy_network/pkg/ledger"
	"time"
)

// Utility function for logging errors.
func logError(action string, err error) {
    log.Printf("[%s] ERROR: %v", time.Now().Format(time.RFC3339), fmt.Errorf("%s: %w", action, err))
}

// Utility function for logging success actions.
func logSuccess(action, message string) {
    log.Printf("[%s] SUCCESS: %s - %s", time.Now().Format(time.RFC3339), action, message)
}


// Centralized validation function for integer inputs
func validatePositiveInt(value int, name string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be a positive integer, got %d", name, value)
	}
	return nil
}

// haInitiateDisasterRecovery initiates the disaster recovery process.
func haInitiateDisasterRecovery(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateDisasterRecovery(); err != nil {
        logError("Initiating Disaster Recovery", err)
        return fmt.Errorf("failed to initiate disaster recovery: %w", err)
    }
    logSuccess("Initiate Disaster Recovery", "Disaster recovery process initiated.")
    return nil
}

// haConfirmDisasterRecovery confirms successful disaster recovery.
func haConfirmDisasterRecovery(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.ConfirmDisasterRecovery(); err != nil {
        logError("Confirming Disaster Recovery", err)
        return fmt.Errorf("failed to confirm disaster recovery: %w", err)
    }
    logSuccess("Confirm Disaster Recovery", "Disaster recovery successfully confirmed.")
    return nil
}

// haCancelDisasterRecovery cancels an ongoing disaster recovery process.
func haCancelDisasterRecovery(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.CancelDisasterRecovery(); err != nil {
        logError("Cancelling Disaster Recovery", err)
        return fmt.Errorf("failed to cancel disaster recovery: %w", err)
    }
    logSuccess("Cancel Disaster Recovery", "Disaster recovery process canceled.")
    return nil
}

// haSetDisasterRecoveryPlan sets the disaster recovery plan.
func haSetDisasterRecoveryPlan(plan string, ledgerInstance *ledger.Ledger) error {
    if plan == "" {
        logError("Setting Disaster Recovery Plan", fmt.Errorf("plan cannot be empty"))
        return fmt.Errorf("disaster recovery plan cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetDisasterRecoveryPlan(plan); err != nil {
        logError("Setting Disaster Recovery Plan", err)
        return fmt.Errorf("failed to set disaster recovery plan: %w", err)
    }
    logSuccess("Set Disaster Recovery Plan", fmt.Sprintf("Disaster recovery plan set to: %s", plan))
    return nil
}

// haGetDisasterRecoveryPlan retrieves the current disaster recovery plan.
func haGetDisasterRecoveryPlan(ledgerInstance *ledger.Ledger) (string, error) {
    plan, err := ledgerInstance.HighAvailabilityLedger.GetDisasterRecoveryPlan()
    if err != nil {
        logError("Getting Disaster Recovery Plan", err)
        return "", fmt.Errorf("failed to get disaster recovery plan: %w", err)
    }
    logSuccess("Get Disaster Recovery Plan", fmt.Sprintf("Retrieved disaster recovery plan: %s", plan))
    return plan, nil
}

// haMonitorDisasterRecovery monitors the disaster recovery status.
func haMonitorDisasterRecovery(ledgerInstance *ledger.Ledger) (string, error) {
    status, err := ledgerInstance.HighAvailabilityLedger.MonitorDisasterRecovery()
    if err != nil {
        logError("Monitoring Disaster Recovery", err)
        return "", fmt.Errorf("failed to monitor disaster recovery: %w", err)
    }
    logSuccess("Monitor Disaster Recovery", fmt.Sprintf("Current disaster recovery status: %s", status))
    return status, nil
}

// haCreateDisasterRecoveryBackup creates a backup for disaster recovery.
func haCreateDisasterRecoveryBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if backupName == "" {
        logError("Creating Disaster Recovery Backup", fmt.Errorf("backup name cannot be empty"))
        return fmt.Errorf("backup name cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.CreateDisasterRecoveryBackup(backupName); err != nil {
        logError("Creating Disaster Recovery Backup", err)
        return fmt.Errorf("failed to create disaster recovery backup: %w", err)
    }
    logSuccess("Create Disaster Recovery Backup", fmt.Sprintf("Disaster recovery backup %s created.", backupName))
    return nil
}

// haRestoreDisasterRecoveryBackup restores data from a disaster recovery backup.
func haRestoreDisasterRecoveryBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if backupName == "" {
        logError("Restoring Disaster Recovery Backup", fmt.Errorf("backup name cannot be empty"))
        return fmt.Errorf("backup name cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.RestoreDisasterRecoveryBackup(backupName); err != nil {
        logError("Restoring Disaster Recovery Backup", err)
        return fmt.Errorf("failed to restore disaster recovery backup: %w", err)
    }
    logSuccess("Restore Disaster Recovery Backup", fmt.Sprintf("Disaster recovery backup %s restored.", backupName))
    return nil
}

// haDeleteDisasterRecoveryBackup deletes a specified disaster recovery backup.
func haDeleteDisasterRecoveryBackup(backupName string, ledgerInstance *ledger.Ledger) error {
    if backupName == "" {
        logError("Deleting Disaster Recovery Backup", fmt.Errorf("backup name cannot be empty"))
        return fmt.Errorf("backup name cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.DeleteDisasterRecoveryBackup(backupName); err != nil {
        logError("Deleting Disaster Recovery Backup", err)
        return fmt.Errorf("failed to delete disaster recovery backup: %w", err)
    }
    logSuccess("Delete Disaster Recovery Backup", fmt.Sprintf("Disaster recovery backup %s deleted.", backupName))
    return nil
}

// haSetDataConsistencyLevel sets the data consistency level for disaster recovery.
func haSetDataConsistencyLevel(level string, ledgerInstance *ledger.Ledger) error {
    if level == "" {
        logError("Setting Data Consistency Level", fmt.Errorf("level cannot be empty"))
        return fmt.Errorf("data consistency level cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetDataConsistencyLevel(level); err != nil {
        logError("Setting Data Consistency Level", err)
        return fmt.Errorf("failed to set data consistency level: %w", err)
    }
    logSuccess("Set Data Consistency Level", fmt.Sprintf("Data consistency level set to: %s", level))
    return nil
}

// haGetDataConsistencyLevel retrieves the data consistency level.
func haGetDataConsistencyLevel(ledgerInstance *ledger.Ledger) (string, error) {
    level, err := ledgerInstance.HighAvailabilityLedger.GetDataConsistencyLevel()
    if err != nil {
        logError("Getting Data Consistency Level", err)
        return "", fmt.Errorf("failed to get data consistency level: %w", err)
    }
    logSuccess("Get Data Consistency Level", fmt.Sprintf("Retrieved data consistency level: %s", level))
    return level, nil
}

// haDisableWriteAheadLog disables write-ahead logging.
func haDisableWriteAheadLog(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.DisableWriteAheadLog(); err != nil {
		logError("Disabling Write-Ahead Log", err)
		return fmt.Errorf("failed to disable write-ahead log: %w", err)
	}
	logSuccess("Disable Write-Ahead Log", "Write-ahead log disabled.")
	return nil
}

// haSetLogRetention sets the log retention period.
func haSetLogRetention(period int, ledgerInstance *ledger.Ledger) error {
	if err := validatePositiveInt(period, "Log retention period"); err != nil {
		logError("Setting Log Retention Period", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetLogRetention(period); err != nil {
		logError("Setting Log Retention Period", err)
		return fmt.Errorf("failed to set log retention period: %w", err)
	}
	logSuccess("Set Log Retention Period", fmt.Sprintf("Log retention period set to %d days.", period))
	return nil
}

// haGetLogRetention retrieves the log retention period.
func haGetLogRetention(ledgerInstance *ledger.Ledger) (int, error) {
	period, err := ledgerInstance.HighAvailabilityLedger.GetLogRetention()
	if err != nil {
		logError("Getting Log Retention Period", err)
		return 0, fmt.Errorf("failed to get log retention period: %w", err)
	}
	logSuccess("Get Log Retention Period", fmt.Sprintf("Retrieved log retention period: %d days.", period))
	return period, nil
}

// haListLogs lists the logs available for review.
func haListLogs(ledgerInstance *ledger.Ledger) ([]string, error) {
	logs, err := ledgerInstance.HighAvailabilityLedger.ListLogs()
	if err != nil {
		logError("Listing Logs", err)
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	logSuccess("List Logs", fmt.Sprintf("Retrieved %d logs for review.", len(logs)))
	return logs, nil
}

// haDeleteLogs deletes specified logs.
func haDeleteLogs(logNames []string, ledgerInstance *ledger.Ledger) error {
	if len(logNames) == 0 {
		logError("Deleting Logs", fmt.Errorf("no logs specified for deletion"))
		return fmt.Errorf("log names list cannot be empty")
	}
	if err := ledgerInstance.HighAvailabilityLedger.DeleteLogs(logNames); err != nil {
		logError("Deleting Logs", err)
		return fmt.Errorf("failed to delete logs: %w", err)
	}
	logSuccess("Delete Logs", fmt.Sprintf("Specified logs deleted: %v.", logNames))
	return nil
}

// haSynchronizeLogs synchronizes logs across nodes for consistency.
func haSynchronizeLogs(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.SynchronizeLogs(); err != nil {
		logError("Synchronizing Logs", err)
		return fmt.Errorf("failed to synchronize logs: %w", err)
	}
	logSuccess("Synchronize Logs", "Logs synchronized across nodes.")
	return nil
}

// haSetHighAvailabilityMode sets the high-availability mode (e.g., active-active).
func haSetHighAvailabilityMode(mode string, ledgerInstance *ledger.Ledger) error {
	if mode == "" {
		logError("Setting High-Availability Mode", fmt.Errorf("mode cannot be empty"))
		return fmt.Errorf("high-availability mode cannot be empty")
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetHighAvailabilityMode(mode); err != nil {
		logError("Setting High-Availability Mode", err)
		return fmt.Errorf("failed to set high-availability mode: %w", err)
	}
	logSuccess("Set High-Availability Mode", fmt.Sprintf("High-availability mode set to %s.", mode))
	return nil
}

// haGetHighAvailabilityMode retrieves the current high-availability mode.
func haGetHighAvailabilityMode(ledgerInstance *ledger.Ledger) (string, error) {
	mode, err := ledgerInstance.HighAvailabilityLedger.GetHighAvailabilityMode()
	if err != nil {
		logError("Getting High-Availability Mode", err)
		return "", fmt.Errorf("failed to get high-availability mode: %w", err)
	}
	logSuccess("Get High-Availability Mode", fmt.Sprintf("Retrieved high-availability mode: %s.", mode))
	return mode, nil
}

// haEnableActiveActive enables the active-active high-availability mode.
func haEnableActiveActive(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.EnableActiveActive(); err != nil {
		logError("Enabling Active-Active Mode", err)
		return fmt.Errorf("failed to enable active-active mode: %w", err)
	}
	logSuccess("Enable Active-Active Mode", "Active-active high-availability mode enabled.")
	return nil
}

// haDisableActiveActive disables the active-active high-availability mode.
func haDisableActiveActive(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.DisableActiveActive(); err != nil {
		logError("Disabling Active-Active Mode", err)
		return fmt.Errorf("failed to disable active-active mode: %w", err)
	}
	logSuccess("Disable Active-Active Mode", "Active-active high-availability mode disabled.")
	return nil
}

// haSetFailoverTimeout sets the failover timeout for high availability.
func haSetFailoverTimeout(timeout int, ledgerInstance *ledger.Ledger) error {
	if err := validatePositiveInt(timeout, "Failover timeout"); err != nil {
		logError("Setting Failover Timeout", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetFailoverTimeout(timeout); err != nil {
		logError("Setting Failover Timeout", err)
		return fmt.Errorf("failed to set failover timeout: %w", err)
	}
	logSuccess("Set Failover Timeout", fmt.Sprintf("Failover timeout set to %d seconds.", timeout))
	return nil
}

// haGetFailoverTimeout retrieves the current failover timeout.
func haGetFailoverTimeout(ledgerInstance *ledger.Ledger) (int, error) {
	timeout, err := ledgerInstance.HighAvailabilityLedger.GetFailoverTimeout()
	if err != nil {
		logError("Getting Failover Timeout", err)
		return 0, fmt.Errorf("failed to get failover timeout: %w", err)
	}
	logSuccess("Get Failover Timeout", fmt.Sprintf("Retrieved failover timeout: %d seconds.", timeout))
	return timeout, nil
}