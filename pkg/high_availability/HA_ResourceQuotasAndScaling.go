package high_availability

import (
	"encoding/base64"
	"fmt"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

func haSaveSimulationResults(simulationID string, results string, ledgerInstance *ledger.Ledger) error {
    if simulationID == "" || results == "" {
        return fmt.Errorf("simulation ID and results cannot be empty")
    }

    // Encrypt the results
    encryption := &common.Encryption{}
    encryptedResults, err := encryption.EncryptData("AES", []byte(results), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt simulation results: %v", err)
    }

    // Convert the encrypted results to a Base64 string
    encodedResults := base64.StdEncoding.EncodeToString(encryptedResults)

    // Pass the Base64 string to the ledger function
    if err := ledgerInstance.HighAvailabilityLedger.SaveSimulationResults(simulationID, encodedResults); err != nil {
        return fmt.Errorf("failed to save simulation results: %v", err)
    }

    fmt.Println("Simulation results saved.")
    return nil
}


// haDeleteSimulationResults deletes simulation results.
func haDeleteSimulationResults(simulationID string, ledgerInstance *ledger.Ledger) error {
    if simulationID == "" {
        return fmt.Errorf("simulation ID cannot be empty")
    }
    if err := ledgerInstance.HighAvailabilityLedger.DeleteSimulationResults(simulationID); err != nil {
        return fmt.Errorf("failed to delete simulation results: %v", err)
    }
    fmt.Printf("Simulation results for %s deleted.\n", simulationID)
    return nil
}

// haListSimulationResults lists all simulation results.
func haListSimulationResults(ledgerInstance *ledger.Ledger) ([]string, error) {
    results, err := ledgerInstance.HighAvailabilityLedger.ListSimulationResults()
    if err != nil {
        return nil, fmt.Errorf("failed to list simulation results: %v", err)
    }
    return results, nil
}

// haEnableResourceQuotas enables resource quotas.
func haEnableResourceQuotas(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableResourceQuotas(); err != nil {
        return fmt.Errorf("failed to enable resource quotas: %v", err)
    }
    fmt.Println("Resource quotas enabled.")
    return nil
}

// haDisableResourceQuotas disables resource quotas.
func haDisableResourceQuotas(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableResourceQuotas(); err != nil {
        return fmt.Errorf("failed to disable resource quotas: %v", err)
    }
    fmt.Println("Resource quotas disabled.")
    return nil
}

func haSetResourceQuotaLimits(limits string, ledgerInstance *ledger.Ledger) error {
    if limits == "" {
        return fmt.Errorf("resource quota limits cannot be empty")
    }

    // Encrypt the limits using AES encryption
    encryption := &common.Encryption{}
    encryptedLimits, err := encryption.EncryptData("AES", []byte(limits), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt resource quota limits: %v", err)
    }

    // Encode the encrypted data into a Base64 string
    encodedLimits := base64.StdEncoding.EncodeToString(encryptedLimits)

    // Pass the Base64 string to the ledger method
    if err := ledgerInstance.HighAvailabilityLedger.SetResourceQuotaLimits(encodedLimits); err != nil {
        return fmt.Errorf("failed to set resource quota limits: %v", err)
    }

    fmt.Println("Resource quota limits set.")
    return nil
}


// haGetResourceQuotaLimits retrieves resource quota limits.
func haGetResourceQuotaLimits(ledgerInstance *ledger.Ledger) (string, error) {
    limits, err := ledgerInstance.HighAvailabilityLedger.GetResourceQuotaLimits()
    if err != nil {
        return "", fmt.Errorf("failed to get resource quota limits: %v", err)
    }
    return limits, nil
}

// haEnableSelfHealing enables self-healing.
func haEnableSelfHealing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableSelfHealing(); err != nil {
        return fmt.Errorf("failed to enable self-healing: %v", err)
    }
    fmt.Println("Self-healing enabled.")
    return nil
}

// haDisableSelfHealing disables self-healing.
func haDisableSelfHealing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableSelfHealing(); err != nil {
        return fmt.Errorf("failed to disable self-healing: %v", err)
    }
    fmt.Println("Self-healing disabled.")
    return nil
}

func haSetSelfHealingInterval(interval int, ledgerInstance *ledger.Ledger) error {
    if interval <= 0 {
        return fmt.Errorf("interval must be greater than 0")
    }

    // Directly set the interval in the ledger
    if err := ledgerInstance.HighAvailabilityLedger.SetSelfHealingInterval(interval); err != nil {
        return fmt.Errorf("failed to set self-healing interval: %v", err)
    }

    fmt.Println("Self-healing interval set.")
    return nil
}


// haGetSelfHealingInterval retrieves the interval for self-healing actions.
func haGetSelfHealingInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.HighAvailabilityLedger.GetSelfHealingInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get self-healing interval: %v", err)
    }
    return interval, nil
}

// haInitiateSelfHealing starts a self-healing operation.
func haInitiateSelfHealing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateSelfHealing(); err != nil {
        return fmt.Errorf("failed to initiate self-healing: %v", err)
    }
    fmt.Println("Self-healing initiated.")
    return nil
}

// haMonitorSelfHealing monitors ongoing self-healing operations.
func haMonitorSelfHealing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.MonitorSelfHealing(); err != nil {
        return fmt.Errorf("failed to monitor self-healing: %v", err)
    }
    fmt.Println("Monitoring self-healing.")
    return nil
}

// haSetFailbackPriority sets the failback priority for recovery.
func haSetFailbackPriority(priority int, ledgerInstance *ledger.Ledger) error {
    if priority < 0 {
        return fmt.Errorf("priority must be non-negative")
    }

    // Encrypt the priority value as a string
    encryption := &common.Encryption{}
    encryptedPriority, err := encryption.EncryptData("AES", []byte(fmt.Sprintf("%d", priority)), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt failback priority: %v", err)
    }

    // Convert the encrypted data to a Base64-encoded string
    encodedPriority := base64.StdEncoding.EncodeToString(encryptedPriority)

    // Pass the encoded string to the ledger function
    if err := ledgerInstance.HighAvailabilityLedger.SetFailbackPriority(encodedPriority); err != nil {
        return fmt.Errorf("failed to set failback priority: %v", err)
    }

    fmt.Println("Failback priority set.")
    return nil
}

// haGetFailbackPriority retrieves the failback priority for recovery.
func haGetFailbackPriority(ledgerInstance *ledger.Ledger) (int, error) {
    priority, err := ledgerInstance.HighAvailabilityLedger.GetFailbackPriority()
    if err != nil {
        return 0, fmt.Errorf("failed to get failback priority: %v", err)
    }
    return priority, nil
}

// haEnableDataArchiving enables data archiving.
func haEnableDataArchiving(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableDataArchiving(); err != nil {
        return fmt.Errorf("failed to enable data archiving: %v", err)
    }
    fmt.Println("Data archiving enabled.")
    return nil
}

// haDisableDataArchiving disables data archiving.
func haDisableDataArchiving(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableDataArchiving(); err != nil {
        return fmt.Errorf("failed to disable data archiving: %v", err)
    }
    fmt.Println("Data archiving disabled.")
    return nil
}

// haScheduleDataArchiving schedules a data archiving operation.
func haScheduleDataArchiving(schedule string, ledgerInstance *ledger.Ledger) error {
    if schedule == "" {
        return fmt.Errorf("schedule cannot be empty")
    }

    encryption := &common.Encryption{}
    encryptedSchedule, err := encryption.EncryptData("AES", []byte(schedule), common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to encrypt schedule: %v", err)
    }

    // Convert the encrypted data to a Base64-encoded string
    encodedSchedule := base64.StdEncoding.EncodeToString(encryptedSchedule)

    // Pass the encoded string to the ledger function
    if err := ledgerInstance.HighAvailabilityLedger.ScheduleDataArchiving(encodedSchedule); err != nil {
        return fmt.Errorf("failed to schedule data archiving: %v", err)
    }

    fmt.Println("Data archiving scheduled.")
    return nil
}

// haListArchivedData lists all archived data.
func haListArchivedData(ledgerInstance *ledger.Ledger) ([]string, error) {
    archives, err := ledgerInstance.HighAvailabilityLedger.ListArchivedData()
    if err != nil {
        return nil, fmt.Errorf("failed to list archived data: %v", err)
    }
    return archives, nil
}

// haRetrieveArchivedData retrieves specific archived data.
func haRetrieveArchivedData(archiveID string, ledgerInstance *ledger.Ledger) (string, error) {
    data, err := ledgerInstance.HighAvailabilityLedger.RetrieveArchivedData(archiveID)
    if err != nil {
        return "", fmt.Errorf("failed to retrieve archived data: %v", err)
    }
    return data, nil
}

// haDeleteArchivedData deletes specific archived data.
func haDeleteArchivedData(archiveID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DeleteArchivedData(archiveID); err != nil {
        return fmt.Errorf("failed to delete archived data: %v", err)
    }
    fmt.Printf("Archived data for %s deleted.\n", archiveID)
    return nil
}
