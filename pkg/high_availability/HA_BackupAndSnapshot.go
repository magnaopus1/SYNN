package high_availability

import (
	"fmt"
	"synnergy_network/pkg/ledger"
)


// Utility function to validate non-empty strings
func validateNonEmptyString(value, name string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", name)
	}
	return nil
}

// haDeleteSnapshot deletes a specific snapshot from storage.
func haDeleteSnapshot(snapshotID string, ledgerInstance *ledger.Ledger) error {
	if err := validateNonEmptyString(snapshotID, "Snapshot ID"); err != nil {
		logError("Deleting Snapshot", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.DeleteSnapshot(snapshotID); err != nil {
		logError("Deleting Snapshot", err)
		return fmt.Errorf("failed to delete snapshot %s: %w", snapshotID, err)
	}
	logSuccess("Delete Snapshot", fmt.Sprintf("Snapshot %s deleted.", snapshotID))
	return nil
}

// haListSnapshots lists all available snapshots.
func haListSnapshots(ledgerInstance *ledger.Ledger) ([]ledger.Snapshot, error) {
	snapshots, err := ledgerInstance.HighAvailabilityLedger.ListSnapshots()
	if err != nil {
		logError("Listing Snapshots", err)
		return nil, fmt.Errorf("failed to list snapshots: %w", err)
	}
	logSuccess("List Snapshots", fmt.Sprintf("Retrieved %d snapshots.", len(snapshots)))
	return snapshots, nil
}

// haSetSnapshotFrequency sets the frequency for creating snapshots.
func haSetSnapshotFrequency(frequency int, ledgerInstance *ledger.Ledger) error {
	if err := validatePositiveInt(frequency, "Snapshot Frequency"); err != nil {
		logError("Setting Snapshot Frequency", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetSnapshotFrequency(frequency); err != nil {
		logError("Setting Snapshot Frequency", err)
		return fmt.Errorf("failed to set snapshot frequency: %w", err)
	}
	logSuccess("Set Snapshot Frequency", fmt.Sprintf("Snapshot frequency set to %d.", frequency))
	return nil
}

// haGetSnapshotFrequency retrieves the current snapshot frequency.
func haGetSnapshotFrequency(ledgerInstance *ledger.Ledger) (int, error) {
	frequency, err := ledgerInstance.HighAvailabilityLedger.GetSnapshotFrequency()
	if err != nil {
		logError("Getting Snapshot Frequency", err)
		return 0, fmt.Errorf("failed to get snapshot frequency: %w", err)
	}
	logSuccess("Get Snapshot Frequency", fmt.Sprintf("Retrieved snapshot frequency: %d.", frequency))
	return frequency, nil
}

// haMonitorSnapshotStatus monitors the status of ongoing snapshots.
func haMonitorSnapshotStatus(ledgerInstance *ledger.Ledger) (ledger.SnapshotStatus, error) {
	status, err := ledgerInstance.HighAvailabilityLedger.MonitorSnapshotStatus()
	if err != nil {
		logError("Monitoring Snapshot Status", err)
		return ledger.SnapshotStatus{}, fmt.Errorf("failed to monitor snapshot status: %w", err)
	}
	logSuccess("Monitor Snapshot Status", "Snapshot status retrieved successfully.")
	return status, nil
}

// haEnableDataMirroring enables data mirroring for high availability.
func haEnableDataMirroring(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.EnableDataMirroring(); err != nil {
		logError("Enabling Data Mirroring", err)
		return fmt.Errorf("failed to enable data mirroring: %w", err)
	}
	logSuccess("Enable Data Mirroring", "Data mirroring enabled.")
	return nil
}

// haDisableDataMirroring disables data mirroring.
func haDisableDataMirroring(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.DisableDataMirroring(); err != nil {
		logError("Disabling Data Mirroring", err)
		return fmt.Errorf("failed to disable data mirroring: %w", err)
	}
	logSuccess("Disable Data Mirroring", "Data mirroring disabled.")
	return nil
}

// haSetMirroringFrequency sets the frequency of data mirroring.
func haSetMirroringFrequency(frequency int, ledgerInstance *ledger.Ledger) error {
	if err := validatePositiveInt(frequency, "Mirroring Frequency"); err != nil {
		logError("Setting Mirroring Frequency", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetMirroringFrequency(frequency); err != nil {
		logError("Setting Mirroring Frequency", err)
		return fmt.Errorf("failed to set mirroring frequency: %w", err)
	}
	logSuccess("Set Mirroring Frequency", fmt.Sprintf("Mirroring frequency set to %d.", frequency))
	return nil
}

// haGetMirroringFrequency retrieves the mirroring frequency.
func haGetMirroringFrequency(ledgerInstance *ledger.Ledger) (int, error) {
	frequency, err := ledgerInstance.HighAvailabilityLedger.GetMirroringFrequency()
	if err != nil {
		logError("Getting Mirroring Frequency", err)
		return 0, fmt.Errorf("failed to get mirroring frequency: %w", err)
	}
	logSuccess("Get Mirroring Frequency", fmt.Sprintf("Retrieved mirroring frequency: %d.", frequency))
	return frequency, nil
}

// haMonitorMirroring monitors the data mirroring process.
func haMonitorMirroring(ledgerInstance *ledger.Ledger) (ledger.MirroringStatus, error) {
	status, err := ledgerInstance.HighAvailabilityLedger.MonitorMirroring()
	if err != nil {
		logError("Monitoring Data Mirroring", err)
		return ledger.MirroringStatus{}, fmt.Errorf("failed to monitor data mirroring: %w", err)
	}
	logSuccess("Monitor Data Mirroring", "Data mirroring status retrieved successfully.")
	return status, nil
}

// haEnableQuorum enables quorum for decision-making.
func haEnableQuorum(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.EnableQuorum(); err != nil {
		logError("Enabling Quorum", err)
		return fmt.Errorf("failed to enable quorum: %w", err)
	}
	logSuccess("Enable Quorum", "Quorum enabled.")
	return nil
}

// haDisableQuorum disables quorum.
func haDisableQuorum(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.DisableQuorum(); err != nil {
		logError("Disabling Quorum", err)
		return fmt.Errorf("failed to disable quorum: %w", err)
	}
	logSuccess("Disable Quorum", "Quorum disabled.")
	return nil
}

// haSetQuorumPolicy sets the quorum policy.
func haSetQuorumPolicy(policy string, ledgerInstance *ledger.Ledger) error {
	if err := validateNonEmptyString(policy, "Quorum Policy"); err != nil {
		logError("Setting Quorum Policy", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetQuorumPolicy(policy); err != nil {
		logError("Setting Quorum Policy", err)
		return fmt.Errorf("failed to set quorum policy: %w", err)
	}
	logSuccess("Set Quorum Policy", fmt.Sprintf("Quorum policy set to %s.", policy))
	return nil
}

// haGetQuorumPolicy retrieves the quorum policy.
func haGetQuorumPolicy(ledgerInstance *ledger.Ledger) (string, error) {
	policy, err := ledgerInstance.HighAvailabilityLedger.GetQuorumPolicy()
	if err != nil {
		logError("Getting Quorum Policy", err)
		return "", fmt.Errorf("failed to get quorum policy: %w", err)
	}
	logSuccess("Get Quorum Policy", fmt.Sprintf("Retrieved quorum policy: %s.", policy))
	return policy, nil
}

// haMonitorQuorum monitors quorum activities.
func haMonitorQuorum(ledgerInstance *ledger.Ledger) (ledger.QuorumStatus, error) {
	status, err := ledgerInstance.HighAvailabilityLedger.MonitorQuorum()
	if err != nil {
		logError("Monitoring Quorum", err)
		return ledger.QuorumStatus{}, fmt.Errorf("failed to monitor quorum: %w", err)
	}
	logSuccess("Monitor Quorum", "Quorum status retrieved successfully.")
	return status, nil
}

// haEnableLoadBalancer enables the load balancer for high availability.
func haEnableLoadBalancer(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.EnableLoadBalancer(); err != nil {
		logError("Enabling Load Balancer", err)
		return fmt.Errorf("failed to enable load balancer: %w", err)
	}
	logSuccess("Enable Load Balancer", "Load balancer enabled.")
	return nil
}

// haDisableLoadBalancer disables the load balancer.
func haDisableLoadBalancer(ledgerInstance *ledger.Ledger) error {
	if err := ledgerInstance.HighAvailabilityLedger.DisableLoadBalancer(); err != nil {
		logError("Disabling Load Balancer", err)
		return fmt.Errorf("failed to disable load balancer: %w", err)
	}
	logSuccess("Disable Load Balancer", "Load balancer disabled.")
	return nil
}

// haSetLoadBalancerPolicy sets the load balancing policy.
func haSetLoadBalancerPolicy(policy string, ledgerInstance *ledger.Ledger) error {
	if err := validateNonEmptyString(policy, "Load Balancer Policy"); err != nil {
		logError("Setting Load Balancer Policy", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetLoadBalancerPolicy(policy); err != nil {
		logError("Setting Load Balancer Policy", err)
		return fmt.Errorf("failed to set load balancer policy: %w", err)
	}
	logSuccess("Set Load Balancer Policy", fmt.Sprintf("Load balancer policy set to %s.", policy))
	return nil
}

// haGetLoadBalancerPolicy retrieves the load balancing policy.
func haGetLoadBalancerPolicy(ledgerInstance *ledger.Ledger) (string, error) {
	policy, err := ledgerInstance.HighAvailabilityLedger.GetLoadBalancerPolicy()
	if err != nil {
		logError("Getting Load Balancer Policy", err)
		return "", fmt.Errorf("failed to get load balancer policy: %w", err)
	}
	logSuccess("Get Load Balancer Policy", fmt.Sprintf("Retrieved load balancer policy: %s.", policy))
	return policy, nil
}

// haAddLoadBalancerNode adds a node to the load balancer.
func haAddLoadBalancerNode(nodeID string, ledgerInstance *ledger.Ledger) error {
	if err := validateNonEmptyString(nodeID, "Node ID"); err != nil {
		logError("Adding Load Balancer Node", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.AddLoadBalancerNode(nodeID); err != nil {
		logError("Adding Load Balancer Node", err)
		return fmt.Errorf("failed to add node %s to load balancer: %w", nodeID, err)
	}
	logSuccess("Add Load Balancer Node", fmt.Sprintf("Node %s added to load balancer.", nodeID))
	return nil
}

// haRemoveLoadBalancerNode removes a node from the load balancer.
func haRemoveLoadBalancerNode(nodeID string, ledgerInstance *ledger.Ledger) error {
	if err := validateNonEmptyString(nodeID, "Node ID"); err != nil {
		logError("Removing Load Balancer Node", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.RemoveLoadBalancerNode(nodeID); err != nil {
		logError("Removing Load Balancer Node", err)
		return fmt.Errorf("failed to remove node %s from load balancer: %w", nodeID, err)
	}
	logSuccess("Remove Load Balancer Node", fmt.Sprintf("Node %s removed from load balancer.", nodeID))
	return nil
}

// haListLoadBalancerNodes lists all nodes under the load balancer.
func haListLoadBalancerNodes(ledgerInstance *ledger.Ledger) ([]string, error) {
	nodes, err := ledgerInstance.HighAvailabilityLedger.ListLoadBalancerNodes()
	if err != nil {
		logError("Listing Load Balancer Nodes", err)
		return nil, fmt.Errorf("failed to list load balancer nodes: %w", err)
	}
	logSuccess("List Load Balancer Nodes", fmt.Sprintf("Retrieved %d load balancer nodes.", len(nodes)))
	return nodes, nil
}

// haMonitorLoadBalancer monitors the status of the load balancer.
func haMonitorLoadBalancer(ledgerInstance *ledger.Ledger) (ledger.LoadBalancerStatus, error) {
	status, err := ledgerInstance.HighAvailabilityLedger.MonitorLoadBalancer()
	if err != nil {
		logError("Monitoring Load Balancer", err)
		return ledger.LoadBalancerStatus{}, fmt.Errorf("failed to monitor load balancer: %w", err)
	}
	logSuccess("Monitor Load Balancer", "Load balancer status retrieved successfully.")
	return status, nil
}

// haSetRecoveryTimeout sets the recovery timeout for high-availability events.
func haSetRecoveryTimeout(timeout int, ledgerInstance *ledger.Ledger) error {
	if err := validatePositiveInt(timeout, "Recovery Timeout"); err != nil {
		logError("Setting Recovery Timeout", err)
		return err
	}
	if err := ledgerInstance.HighAvailabilityLedger.SetRecoveryTimeout(timeout); err != nil {
		logError("Setting Recovery Timeout", err)
		return fmt.Errorf("failed to set recovery timeout: %w", err)
	}
	logSuccess("Set Recovery Timeout", fmt.Sprintf("Recovery timeout set to %d seconds.", timeout))
	return nil
}

// haGetRecoveryTimeout retrieves the current recovery timeout.
func haGetRecoveryTimeout(ledgerInstance *ledger.Ledger) (int, error) {
	timeout, err := ledgerInstance.HighAvailabilityLedger.GetRecoveryTimeout()
	if err != nil {
		logError("Getting Recovery Timeout", err)
		return 0, fmt.Errorf("failed to get recovery timeout: %w", err)
	}
	logSuccess("Get Recovery Timeout", fmt.Sprintf("Retrieved recovery timeout: %d seconds.", timeout))
	return timeout, nil
}
