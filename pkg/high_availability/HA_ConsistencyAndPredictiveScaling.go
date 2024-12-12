package high_availability


import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// haSetArchiveRetentionPolicy sets the archive retention policy.
func haSetArchiveRetentionPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetArchiveRetentionPolicy(policy); err != nil {
        return fmt.Errorf("failed to set archive retention policy: %v", err)
    }
    fmt.Printf("Archive retention policy set to %s.\n", policy)
    return nil
}

// haGetArchiveRetentionPolicy retrieves the current archive retention policy.
func haGetArchiveRetentionPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetArchiveRetentionPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get archive retention policy: %v", err)
    }
    return policy, nil
}

// haEnableConsistencyChecks enables periodic consistency checks.
func haEnableConsistencyChecks(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableConsistencyChecks(); err != nil {
        return fmt.Errorf("failed to enable consistency checks: %v", err)
    }
    fmt.Println("Consistency checks enabled.")
    return nil
}

// haDisableConsistencyChecks disables consistency checks.
func haDisableConsistencyChecks(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableConsistencyChecks(); err != nil {
        return fmt.Errorf("failed to disable consistency checks: %v", err)
    }
    fmt.Println("Consistency checks disabled.")
    return nil
}

// haSetConsistencyCheckInterval sets the interval for consistency checks.
func haSetConsistencyCheckInterval(interval int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetConsistencyCheckInterval(interval); err != nil {
        return fmt.Errorf("failed to set consistency check interval: %v", err)
    }
    fmt.Printf("Consistency check interval set to %d seconds.\n", interval)
    return nil
}

// haGetConsistencyCheckInterval retrieves the consistency check interval.
func haGetConsistencyCheckInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.HighAvailabilityLedger.GetConsistencyCheckInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get consistency check interval: %v", err)
    }
    return interval, nil
}

// haInitiateConsistencyCheck initiates a consistency check across the network.
func haInitiateConsistencyCheck(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateConsistencyCheck(); err != nil {
        return fmt.Errorf("failed to initiate consistency check: %v", err)
    }
    fmt.Println("Consistency check initiated.")
    return nil
}

// haListConsistencyCheckResults lists the results of past consistency checks.
func haListConsistencyCheckResults(ledgerInstance *ledger.Ledger) ([]ledger.ConsistencyCheckResult, error) {
    results, err := ledgerInstance.HighAvailabilityLedger.ListConsistencyCheckResults()
    if err != nil {
        return nil, fmt.Errorf("failed to list consistency check results: %v", err)
    }
    return results, nil
}

// haEnablePredictiveScaling enables predictive scaling based on usage patterns.
func haEnablePredictiveScaling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnablePredictiveScaling(); err != nil {
        return fmt.Errorf("failed to enable predictive scaling: %v", err)
    }
    fmt.Println("Predictive scaling enabled.")
    return nil
}

// haDisablePredictiveScaling disables predictive scaling.
func haDisablePredictiveScaling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisablePredictiveScaling(); err != nil {
        return fmt.Errorf("failed to disable predictive scaling: %v", err)
    }
    fmt.Println("Predictive scaling disabled.")
    return nil
}

// haSetPredictiveScalingPolicy sets the policy for predictive scaling.
func haSetPredictiveScalingPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetPredictiveScalingPolicy(policy); err != nil {
        return fmt.Errorf("failed to set predictive scaling policy: %v", err)
    }
    fmt.Printf("Predictive scaling policy set to %s.\n", policy)
    return nil
}

// haGetPredictiveScalingPolicy retrieves the current predictive scaling policy.
func haGetPredictiveScalingPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetPredictiveScalingPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get predictive scaling policy: %v", err)
    }
    return policy, nil
}

// haEnablePredictiveFailover enables predictive failover for high availability.
func haEnablePredictiveFailover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnablePredictiveFailover(); err != nil {
        return fmt.Errorf("failed to enable predictive failover: %v", err)
    }
    fmt.Println("Predictive failover enabled.")
    return nil
}

// haDisablePredictiveFailover disables predictive failover.
func haDisablePredictiveFailover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisablePredictiveFailover(); err != nil {
        return fmt.Errorf("failed to disable predictive failover: %v", err)
    }
    fmt.Println("Predictive failover disabled.")
    return nil
}


// haSetPredictiveFailoverPolicy sets the predictive failover policy.
func haSetPredictiveFailoverPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetPredictiveFailoverPolicy(policy); err != nil {
        return fmt.Errorf("failed to set predictive failover policy: %v", err)
    }
    fmt.Printf("Predictive failover policy set to %s.\n", policy)
    return nil
}

// haGetPredictiveFailoverPolicy retrieves the predictive failover policy.
func haGetPredictiveFailoverPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetPredictiveFailoverPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get predictive failover policy: %v", err)
    }
    return policy, nil
}

// haEnableAdaptiveResourceManagement enables adaptive resource management.
func haEnableAdaptiveResourceManagement(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableAdaptiveResourceManagement(); err != nil {
        return fmt.Errorf("failed to enable adaptive resource management: %v", err)
    }
    fmt.Println("Adaptive resource management enabled.")
    return nil
}

// haDisableAdaptiveResourceManagement disables adaptive resource management.
func haDisableAdaptiveResourceManagement(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableAdaptiveResourceManagement(); err != nil {
        return fmt.Errorf("failed to disable adaptive resource management: %v", err)
    }
    fmt.Println("Adaptive resource management disabled.")
    return nil
}

// haSetAdaptiveResourcePolicy sets the adaptive resource management policy.
func haSetAdaptiveResourcePolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetAdaptiveResourcePolicy(policy); err != nil {
        return fmt.Errorf("failed to set adaptive resource policy: %v", err)
    }
    fmt.Printf("Adaptive resource policy set to %s.\n", policy)
    return nil
}

// haGetAdaptiveResourcePolicy retrieves the adaptive resource policy.
func haGetAdaptiveResourcePolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetAdaptiveResourcePolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get adaptive resource policy: %v", err)
    }
    return policy, nil
}

// haSimulateNodeFailure simulates a node failure for high-availability testing.
func haSimulateNodeFailure(nodeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateNodeFailure(nodeID); err != nil {
        return fmt.Errorf("failed to simulate node failure for %s: %v", nodeID, err)
    }
    fmt.Printf("Node failure simulation completed for node %s.\n", nodeID)
    return nil
}

// haSimulateNetworkPartition simulates a network partition for resilience testing.
func haSimulateNetworkPartition(partitionID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateNetworkPartition(partitionID); err != nil {
        return fmt.Errorf("failed to simulate network partition for %s: %v", partitionID, err)
    }
    fmt.Printf("Network partition simulation completed for partition %s.\n", partitionID)
    return nil
}

// haSimulateDiskFailure simulates a disk failure scenario.
func haSimulateDiskFailure(diskID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateDiskFailure(diskID); err != nil {
        return fmt.Errorf("failed to simulate disk failure for %s: %v", diskID, err)
    }
    fmt.Printf("Disk failure simulation completed for disk %s.\n", diskID)
    return nil
}

// haSimulateMemoryFailure simulates a memory failure scenario.
func haSimulateMemoryFailure(memoryID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateMemoryFailure(memoryID); err != nil {
        return fmt.Errorf("failed to simulate memory failure for %s: %v", memoryID, err)
    }
    fmt.Printf("Memory failure simulation completed for memory %s.\n", memoryID)
    return nil
}

// haSimulateCPUOverload simulates a CPU overload situation.
func haSimulateCPUOverload(cpuID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateCPUOverload(cpuID); err != nil {
        return fmt.Errorf("failed to simulate CPU overload for %s: %v", cpuID, err)
    }
    fmt.Printf("CPU overload simulation completed for CPU %s.\n", cpuID)
    return nil
}

// haSimulateApplicationFailure simulates an application failure for resilience testing.
func haSimulateApplicationFailure(appID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SimulateApplicationFailure(appID); err != nil {
        return fmt.Errorf("failed to simulate application failure for %s: %v", appID, err)
    }
    fmt.Printf("Application failure simulation completed for application %s.\n", appID)
    return nil
}
