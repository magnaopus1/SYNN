package high_availability


import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// haSetColdStandbyPolicy sets the policy for cold standby resources.
func haSetColdStandbyPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetColdStandbyPolicy(policy); err != nil {
        return fmt.Errorf("failed to set cold standby policy: %v", err)
    }
    fmt.Println("Cold standby policy set.")
    return nil
}

// haEnableFailoverGroups enables failover groups.
func haEnableFailoverGroups(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResourceManagementLedger.EnableResourcePooling(); err != nil {
        return fmt.Errorf("failed to enable failover groups: %v", err)
    }
    fmt.Println("Failover groups enabled.")
    return nil
}

// haDisableFailoverGroups disables failover groups.
func haDisableFailoverGroups(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResourceManagementLedger.DisableResourcePooling(); err != nil {
        return fmt.Errorf("failed to disable failover groups: %v", err)
    }
    fmt.Println("Failover groups disabled.")
    return nil
}

// haSetFailoverGroupPolicy sets the failover group policy.
func haSetFailoverGroupPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetFailoverGroupPolicy(policy); err != nil {
        return fmt.Errorf("failed to set failover group policy: %v", err)
    }
    fmt.Println("Failover group policy set.")
    return nil
}

// haGetFailoverGroupPolicy retrieves the failover group policy.
func haGetFailoverGroupPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetFailoverGroupPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get failover group policy: %v", err)
    }
    return policy, nil
}

// haAddFailoverGroupMember adds a new member to a failover group.
func haAddFailoverGroupMember(member string, groupID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.AddFailoverGroupMember(member, groupID); err != nil {
        return fmt.Errorf("failed to add member to failover group %s: %v", groupID, err)
    }
    fmt.Printf("Member %s added to failover group %s.\n", member, groupID)
    return nil
}

// haRemoveFailoverGroupMember removes a member from a failover group.
func haRemoveFailoverGroupMember(member string, groupID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RemoveFailoverGroupMember(member, groupID); err != nil {
        return fmt.Errorf("failed to remove member from failover group %s: %v", groupID, err)
    }
    fmt.Printf("Member %s removed from failover group %s.\n", member, groupID)
    return nil
}

// haEnableHAProxy enables high availability proxy services.
func haEnableHAProxy(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableHAProxy(); err != nil {
        return fmt.Errorf("failed to enable HA proxy: %v", err)
    }
    fmt.Println("HA proxy enabled.")
    return nil
}

// haDisableHAProxy disables high availability proxy services.
func haDisableHAProxy(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableHAProxy(); err != nil {
        return fmt.Errorf("failed to disable HA proxy: %v", err)
    }
    fmt.Println("HA proxy disabled.")
    return nil
}

// haSetHAProxyPolicy sets the policy for HA proxy.
func haSetHAProxyPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetHAProxyPolicy(policy); err != nil {
        return fmt.Errorf("failed to set HA proxy policy: %v", err)
    }
    fmt.Println("HA proxy policy set.")
    return nil
}

// haGetHAProxyPolicy retrieves the HA proxy policy.
func haGetHAProxyPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetHAProxyPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get HA proxy policy: %v", err)
    }
    return policy, nil
}


// HAEnableResourcePooling enables resource pooling for high availability.
func HAEnableResourcePooling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResourceManagementLedger.EnableResourcePooling(); err != nil {
        return fmt.Errorf("failed to enable resource pooling: %v", err)
    }
    fmt.Println("Resource pooling enabled.")
    return nil
}

// HADisableResourcePooling disables resource pooling.
func HADisableResourcePooling(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResourceManagementLedger.DisableResourcePooling(); err != nil {
        return fmt.Errorf("failed to disable resource pooling: %v", err)
    }
    fmt.Println("Resource pooling disabled.")
    return nil
}

// haSetResourcePoolingPolicy sets the resource pooling policy.
func haSetResourcePoolingPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.ResourceManagementLedger.SetResourcePoolingPolicy(policy); err != nil {
        return fmt.Errorf("failed to set resource pooling policy: %v", err)
    }
    fmt.Println("Resource pooling policy set.")
    return nil
}

// haGetResourcePoolingPolicy retrieves the resource pooling policy.
func haGetResourcePoolingPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.ResourceManagementLedger.GetResourcePoolingPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get resource pooling policy: %v", err)
    }
    return policy, nil
}

// haEnableGeoRedundancy enables geographic redundancy.
func haEnableGeoRedundancy(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableGeoRedundancy(); err != nil {
        return fmt.Errorf("failed to enable geo redundancy: %v", err)
    }
    fmt.Println("Geo redundancy enabled.")
    return nil
}

// haDisableGeoRedundancy disables geographic redundancy.
func haDisableGeoRedundancy(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableGeoRedundancy(); err != nil {
        return fmt.Errorf("failed to disable geo redundancy: %v", err)
    }
    fmt.Println("Geo redundancy disabled.")
    return nil
}

// haSetGeoRedundancyPolicy sets the geographic redundancy policy.
func haSetGeoRedundancyPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetGeoRedundancyPolicy(policy); err != nil {
        return fmt.Errorf("failed to set geo redundancy policy: %v", err)
    }
    fmt.Println("Geo redundancy policy set.")
    return nil
}

// haGetGeoRedundancyPolicy retrieves the geographic redundancy policy.
func haGetGeoRedundancyPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetGeoRedundancyPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get geo redundancy policy: %v", err)
    }
    return policy, nil
}

// haSetDisasterSimulationMode sets the mode for disaster simulation.
func haSetDisasterSimulationMode(mode string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetDisasterSimulationMode(mode); err != nil {
        return fmt.Errorf("failed to set disaster simulation mode: %v", err)
    }
    fmt.Printf("Disaster simulation mode set to %s.\n", mode)
    return nil
}

// haEnableDisasterSimulation enables disaster simulation.
func haEnableDisasterSimulation(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableDisasterSimulation(); err != nil {
        return fmt.Errorf("failed to enable disaster simulation: %v", err)
    }
    fmt.Println("Disaster simulation enabled.")
    return nil
}

// haDisableDisasterSimulation disables disaster simulation.
func haDisableDisasterSimulation(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableDisasterSimulation(); err != nil {
        return fmt.Errorf("failed to disable disaster simulation: %v", err)
    }
    fmt.Println("Disaster simulation disabled.")
    return nil
}

// haInitiateDisasterSimulation initiates a new disaster simulation.
func haInitiateDisasterSimulation(params string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitiateDisasterSimulation(params); err != nil {
        return fmt.Errorf("failed to initiate disaster simulation: %v", err)
    }
    fmt.Println("Disaster simulation initiated.")
    return nil
}

// haConfirmDisasterSimulation confirms a disaster simulation.
func haConfirmDisasterSimulation(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.ConfirmDisasterSimulation(); err != nil {
        return fmt.Errorf("failed to confirm disaster simulation: %v", err)
    }
    fmt.Println("Disaster simulation confirmed.")
    return nil
}

// haSetSimulationParameters sets parameters for disaster simulation.
func haSetSimulationParameters(params string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetSimulationParameters(params); err != nil {
        return fmt.Errorf("failed to set simulation parameters: %v", err)
    }
    fmt.Println("Simulation parameters set.")
    return nil
}

// haGetSimulationParameters retrieves parameters for disaster simulation.
func haGetSimulationParameters(ledgerInstance *ledger.Ledger) (string, error) {
    params, err := ledgerInstance.HighAvailabilityLedger.GetSimulationParameters()
    if err != nil {
        return "", fmt.Errorf("failed to get simulation parameters: %v", err)
    }
    return params, nil
}
