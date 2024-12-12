package high_availability


import (
    "fmt"
    "synnergy_network/pkg/ledger"
)

// haInit initializes the high-availability setup.
func haInit(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.InitializeHA(); err != nil {
        return fmt.Errorf("failed to initialize high availability: %v", err)
    }
    fmt.Println("High availability initialized.")
    return nil
}

// haEnableReplication enables replication.
func haEnableReplication(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableReplication(); err != nil {
        return fmt.Errorf("failed to enable replication: %v", err)
    }
    fmt.Println("Replication enabled.")
    return nil
}

// haDisableReplication disables replication.
func haDisableReplication(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableReplication(); err != nil {
        return fmt.Errorf("failed to disable replication: %v", err)
    }
    fmt.Println("Replication disabled.")
    return nil
}

// haFailover performs a failover operation.
func haFailover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.Failover(); err != nil {
        return fmt.Errorf("failover operation failed: %v", err)
    }
    fmt.Println("Failover completed.")
    return nil
}

// haSwitchover performs a switchover operation.
func haSwitchover(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.Switchover(); err != nil {
        return fmt.Errorf("switchover operation failed: %v", err)
    }
    fmt.Println("Switchover completed.")
    return nil
}

// haEnableLoadBalancing enables load balancing.
func haEnableLoadBalancing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableLoadBalancing(); err != nil {
        return fmt.Errorf("failed to enable load balancing: %v", err)
    }
    fmt.Println("Load balancing enabled.")
    return nil
}

// haDisableLoadBalancing disables load balancing.
func haDisableLoadBalancing(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableLoadBalancing(); err != nil {
        return fmt.Errorf("failed to disable load balancing: %v", err)
    }
    fmt.Println("Load balancing disabled.")
    return nil
}

// haSetReplicationFactor sets the replication factor.
func haSetReplicationFactor(factor int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetReplicationFactor(factor); err != nil {
        return fmt.Errorf("failed to set replication factor: %v", err)
    }
    fmt.Println("Replication factor set.")
    return nil
}

// haGetReplicationFactor retrieves the replication factor.
func haGetReplicationFactor(ledgerInstance *ledger.Ledger) (int, error) {
    factor, err := ledgerInstance.HighAvailabilityLedger.GetReplicationFactor()
    if err != nil {
        return 0, fmt.Errorf("failed to get replication factor: %v", err)
    }
    return factor, nil
}

// haEnableCluster enables clustering.
func haEnableCluster(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableCluster(); err != nil {
        return fmt.Errorf("failed to enable clustering: %v", err)
    }
    fmt.Println("Clustering enabled.")
    return nil
}

// haDisableCluster disables clustering.
func haDisableCluster(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableCluster(); err != nil {
        return fmt.Errorf("failed to disable clustering: %v", err)
    }
    fmt.Println("Clustering disabled.")
    return nil
}

// haAddClusterNode adds a node to the cluster.
func haAddClusterNode(nodeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.AddClusterNode(nodeID); err != nil {
        return fmt.Errorf("failed to add node to cluster: %v", err)
    }
    fmt.Printf("Node %s added to cluster.\n", nodeID)
    return nil
}

// haRemoveClusterNode removes a node from the cluster.
func haRemoveClusterNode(nodeID string, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.RemoveClusterNode(nodeID); err != nil {
        return fmt.Errorf("failed to remove node from cluster: %v", err)
    }
    fmt.Printf("Node %s removed from cluster.\n", nodeID)
    return nil
}

// haListClusterNodes lists all nodes in the cluster.
func haListClusterNodes(ledgerInstance *ledger.Ledger) ([]string, error) {
    nodes, err := ledgerInstance.HighAvailabilityLedger.ListClusterNodes()
    if err != nil {
        return nil, fmt.Errorf("failed to list cluster nodes: %v", err)
    }
    return nodes, nil
}

// haSetClusterPolicy sets the policy for clustering.
func haSetClusterPolicy(policy string, ledgerInstance *ledger.Ledger) error {
    if policy == "" {
        return fmt.Errorf("policy cannot be empty")
    }
    clusterPolicy := ledger.ClusterPolicy{
        PolicyName: policy,
        Rules:      "Standard Rules",
    }
    if err := ledgerInstance.HighAvailabilityLedger.SetClusterPolicy(clusterPolicy); err != nil {
        return fmt.Errorf("failed to set cluster policy: %v", err)
    }
    fmt.Println("Cluster policy set.")
    return nil
}

// haGetClusterPolicy retrieves the cluster policy.
func haGetClusterPolicy(ledgerInstance *ledger.Ledger) (string, error) {
    policy, err := ledgerInstance.HighAvailabilityLedger.GetClusterPolicy()
    if err != nil {
        return "", fmt.Errorf("failed to get cluster policy: %v", err)
    }
    return policy.PolicyName, nil
}

// haEnableHeartbeat enables heartbeat monitoring.
func haEnableHeartbeat(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableHeartbeat(); err != nil {
        return fmt.Errorf("failed to enable heartbeat monitoring: %v", err)
    }
    fmt.Println("Heartbeat monitoring enabled.")
    return nil
}

// haDisableHeartbeat disables heartbeat monitoring.
func haDisableHeartbeat(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableHeartbeat(); err != nil {
        return fmt.Errorf("failed to disable heartbeat monitoring: %v", err)
    }
    fmt.Println("Heartbeat monitoring disabled.")
    return nil
}

// haSetHeartbeatInterval sets the interval for heartbeat checks.
func haSetHeartbeatInterval(interval int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetHeartbeatInterval(interval); err != nil {
        return fmt.Errorf("failed to set heartbeat interval: %v", err)
    }
    fmt.Println("Heartbeat interval set.")
    return nil
}

// haGetHeartbeatInterval retrieves the interval for heartbeat checks.
func haGetHeartbeatInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.HighAvailabilityLedger.GetHeartbeatInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get heartbeat interval: %v", err)
    }
    return interval, nil
}

// haMonitorHeartbeat monitors heartbeat signals.
func haMonitorHeartbeat(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.MonitorHeartbeat(); err != nil {
        return fmt.Errorf("failed to monitor heartbeat: %v", err)
    }
    fmt.Println("Heartbeat monitoring in progress.")
    return nil
}

// haEnableHealthCheck enables health checks.
func haEnableHealthCheck(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.EnableHealthCheck(); err != nil {
        return fmt.Errorf("failed to enable health check: %v", err)
    }
    fmt.Println("Health check enabled.")
    return nil
}

// haDisableHealthCheck disables health checks.
func haDisableHealthCheck(ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.DisableHealthCheck(); err != nil {
        return fmt.Errorf("failed to disable health check: %v", err)
    }
    fmt.Println("Health check disabled.")
    return nil
}

// haSetHealthCheckInterval sets the interval for health checks.
func haSetHealthCheckInterval(interval int, ledgerInstance *ledger.Ledger) error {
    if err := ledgerInstance.HighAvailabilityLedger.SetHealthCheckInterval(interval); err != nil {
        return fmt.Errorf("failed to set health check interval: %v", err)
    }
    fmt.Println("Health check interval set.")
    return nil
}

// haGetHealthCheckInterval retrieves the interval for health checks.
func haGetHealthCheckInterval(ledgerInstance *ledger.Ledger) (int, error) {
    interval, err := ledgerInstance.HighAvailabilityLedger.GetHealthCheckInterval()
    if err != nil {
        return 0, fmt.Errorf("failed to get health check interval: %v", err)
    }
    return interval, nil
}
