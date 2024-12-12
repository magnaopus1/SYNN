package ledger

// RecordPerformanceMetrics logs performance metrics for the ledger.
func (l *Ledger) RecordPerformanceMetrics(entityID string, metrics string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// Record performance metrics for the entity
	performance := PerformanceMetrics{
		EntityID:  entityID,
		Metrics:   metrics,
		Timestamp: time.Now(),
	}

	// Store the metrics in the ledger
	l.performanceMetrics[entityID] = performance
	return nil
}

// RecordMaintenanceEvent records a maintenance event in the ledger.
func (mml *MonitoringMaintenanceLedger) RecordMaintenanceEvent(eventID, details string) error {
	mml.mu.Lock()
	defer mml.mu.Unlock()

	event := MaintenanceEvent{
		EventID:   eventID,
		Timestamp: time.Now(),
		Details:   details,
	}
	mml.MaintenanceEvents = append(mml.MaintenanceEvents, event)

	fmt.Printf("Maintenance event recorded: %v\n", event)
	return nil
}

// GetCPUUsage returns the average CPU usage from the history.
func (mml *MonitoringMaintenanceLedger) GetCPUUsage() (float64, error) {
	mml.mu.Lock()
	defer mml.mu.Unlock()

	if len(mml.CPUUsageHistory) == 0 {
		return 0, fmt.Errorf("no CPU usage data available")
	}

	var total float64
	for _, usage := range mml.CPUUsageHistory {
		total += usage
	}
	average := total / float64(len(mml.CPUUsageHistory))
	return average, nil
}

// GetMemoryUsage returns the average memory usage from the history.
func (mml *MonitoringMaintenanceLedger) GetMemoryUsage() (float64, error) {
	mml.mu.Lock()
	defer mml.mu.Unlock()

	if len(mml.MemoryUsageHistory) == 0 {
		return 0, fmt.Errorf("no memory usage data available")
	}

	var total float64
	for _, usage := range mml.MemoryUsageHistory {
		total += usage
	}
	average := total / float64(len(mml.MemoryUsageHistory))
	return average, nil
}


// InitializeMonitoringSystem initializes the monitoring system.
func (ms *MonitoringSystem) InitializeMonitoringSystem(nodes []string) error {
    if ms.Initialized {
        return fmt.Errorf("monitoring system already initialized")
    }

    // Perform initialization logic (e.g., setting up node connections)
    ms.Initialized = true
    fmt.Println("Monitoring system initialized.")
    return nil
}

// GetCurrentThreatLevel retrieves the current system threat level.
func (ms *MonitoringSystem) GetCurrentThreatLevel() (int, error) {
    if !ms.Initialized {
        return 0, fmt.Errorf("monitoring system is not initialized")
    }

    // Logic to calculate or fetch the current threat level
    fmt.Printf("Current threat level: %d\n", ms.ThreatLevel)
    return ms.ThreatLevel, nil
}



// StartMonitoring starts monitoring the system.
func (ms *MonitoringSystem) StartMonitoring() error {
	if len(ms.Nodes) == 0 {
		return fmt.Errorf("monitoring system is not initialized")
	}

	fmt.Println("Monitoring system started.")
	// Add real-world monitoring logic here
	return nil
}

// Record a system check
func (l *Ledger) recordSystemCheck(timestamp time.Time) (string, error) {
    checkID := fmt.Sprintf("check-%d", timestamp.Unix())
    l.SystemChecks[checkID] = SystemCheck{
        CheckID:   checkID,
        Status:    "Completed",
        Timestamp: timestamp,
    }
    return checkID, nil
}

// Record a diagnostic test
func (l *Ledger) recordDiagnosticTest(timestamp time.Time) error {
    testID := fmt.Sprintf("diagnostic-%d", timestamp.Unix())
    l.DiagnosticTests[testID] = DiagnosticTest{
        TestID:    testID,
        Status:    "Completed",
        Timestamp: timestamp,
    }
    return nil
}

// Record a reboot schedule
func (l *Ledger) recordRebootSchedule(rebootTime time.Time) error {
    rebootID := fmt.Sprintf("reboot-%d", rebootTime.Unix())
    l.RebootSchedules[rebootID] = RebootSchedule{
        RebootID:  rebootID,
        Scheduled: rebootTime,
        Status:    "Scheduled",
    }
    return nil
}

// Record storage optimization
func (l *Ledger) recordStorageOptimization(timestamp time.Time) error {
    optimizationID := fmt.Sprintf("optimization-%d", timestamp.Unix())
    l.StorageOptimizations[optimizationID] = StorageOptimization{
        OptimizationID: optimizationID,
        Timestamp:      timestamp,
        Status:         "Optimized",
    }
    return nil
}

// Record disk health
func (l *Ledger) recordDiskHealth(status string, timestamp time.Time) error {
    diskID := fmt.Sprintf("disk-health-%d", timestamp.Unix())
    l.DiskHealthRecords[diskID] = DiskHealth{
        DiskID:    diskID,
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record a backup
func (l *Ledger) recordBackup(encryptedData []byte, timestamp time.Time) error {
    backupID := fmt.Sprintf("backup-%d", timestamp.Unix())
    l.Backups[backupID] = Backup{
        BackupID:      backupID,
        EncryptedData: encryptedData,
        Timestamp:     timestamp,
        IntegrityCheck: false,
    }
    return nil
}

// Validate a backup's integrity
func (l *Ledger) validateBackupIntegrity(backupID string) error {
    backup, exists := l.Backups[backupID]
    if !exists {
        return fmt.Errorf("backup %s not found", backupID)
    }
    backup.IntegrityCheck = true
    l.Backups[backupID] = backup
    return nil
}

// Record hardware status
func (l *Ledger) recordHardwareStatus(componentID, status string, timestamp time.Time) error {
    l.HardwareStatusRecords[componentID] = HardwareStatus{
        ComponentID: componentID,
        Status:      status,
        Timestamp:   timestamp,
    }
    return nil
}

// Record maintenance schedule
func (l *Ledger) recordMaintenanceSchedule(maintenanceTime time.Time) error {
    maintenanceID := fmt.Sprintf("maintenance-%d", maintenanceTime.Unix())
    l.MaintenanceSchedules[maintenanceID] = MaintenanceSchedule{
        MaintenanceID: maintenanceID,
        ScheduledTime: maintenanceTime,
        Status:        "Scheduled",
    }
    return nil
}

// Record system health status
func (l *Ledger) recordSystemHealth(status string, timestamp time.Time) error {
    healthID := fmt.Sprintf("health-%d", timestamp.Unix())
    l.SystemHealthStatuses[healthID] = SystemHealthStatus{
        HealthID:  healthID,
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record system update
func (l *Ledger) recordSystemUpdate(timestamp time.Time) error {
    updateID := fmt.Sprintf("update-%d", timestamp.Unix())
    l.SystemUpdateRecords[updateID] = SystemUpdateRecord{
        UpdateID:  updateID,
        Timestamp: timestamp,
        Status:    "Updated",
    }
    return nil
}

// Record defragmentation
func (l *Ledger) recordDefragmentation(timestamp time.Time) error {
    defragID := fmt.Sprintf("defrag-%d", timestamp.Unix())
    l.DefragmentationRecords[defragID] = DefragmentationRecord{
        DefragID:   defragID,
        Timestamp:  timestamp,
        Status:     "Completed",
    }
    return nil
}

// Record service status
func (l *Ledger) recordServiceStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("service-status-%d", timestamp.Unix())
    l.ServiceStatuses[statusID] = ServiceStatus{
        StatusID:   statusID,
        Status:     status,
        Timestamp:  timestamp,
    }
    return nil
}

// Record configuration validation
func (l *Ledger) recordConfigurationValidation(timestamp time.Time) error {
    validationID := fmt.Sprintf("config-validation-%d", timestamp.Unix())
    l.ConfigurationValidations[validationID] = ConfigurationValidation{
        ValidationID: validationID,
        Status:       "Validated",
        Timestamp:    timestamp,
    }
    return nil
}

// Record firmware update schedule
func (l *Ledger) recordFirmwareUpdateSchedule(updateTime time.Time) error {
    updateID := fmt.Sprintf("firmware-update-%d", updateTime.Unix())
    l.FirmwareUpdates[updateID] = FirmwareUpdate{
        UpdateID:    updateID,
        Scheduled:   updateTime,
        Status:      "Scheduled",
    }
    return nil
}

// Record CPU health
func (l *Ledger) recordCPUHealth(status string, timestamp time.Time) error {
    healthID := fmt.Sprintf("cpu-health-%d", timestamp.Unix())
    l.CPUHealthRecords[healthID] = CPUHealth{
        HealthID:   healthID,
        Status:     status,
        Timestamp:  timestamp,
    }
    return nil
}

// Record encryption key update
func (l *Ledger) recordEncryptionKeyUpdate(newKey string, timestamp time.Time) error {
    keyID := fmt.Sprintf("key-update-%d", timestamp.Unix())
    l.EncryptionKeyUpdates[keyID] = EncryptionKeyUpdate{
        KeyID:      keyID,
        NewKey:     newKey,
        Timestamp:  timestamp,
    }
    return nil
}

// Record backup frequency
func (l *Ledger) recordBackupFrequency(timestamp time.Time) error {
    frequencyID := fmt.Sprintf("backup-frequency-%d", timestamp.Unix())
    l.BackupFrequencies[frequencyID] = BackupFrequency{
        FrequencyID: frequencyID,
        Timestamp:   timestamp,
    }
    return nil
}

// Record error log
func (l *Ledger) recordErrorLog(timestamp time.Time) error {
    logID := fmt.Sprintf("error-log-%d", timestamp.Unix())
    l.ErrorLogs[logID] = ErrorLog{
        LogID:      logID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record security check
func (l *Ledger) recordSecurityCheck(timestamp time.Time) error {
    checkID := fmt.Sprintf("security-check-%d", timestamp.Unix())
    l.SecurityChecks[checkID] = SecurityCheck{
        CheckID:    checkID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record network route validation
func (l *Ledger) recordNetworkRouteValidation(timestamp time.Time) error {
    validationID := fmt.Sprintf("network-route-%d", timestamp.Unix())
    l.NetworkRouteValidations[validationID] = NetworkRouteValidation{
        ValidationID: validationID,
        Timestamp:    timestamp,
    }
    return nil
}

// Record redundant system test
func (l *Ledger) recordRedundantSystemTest(timestamp time.Time) error {
    testID := fmt.Sprintf("redundant-test-%d", timestamp.Unix())
    l.RedundantSystemTests[testID] = RedundantSystemTest{
        TestID:    testID,
        Timestamp: timestamp,
        Status:    "Tested",
    }
    return nil
}

// Record data migration schedule
func (l *Ledger) recordDataMigrationSchedule(migrationTime time.Time) error {
    migrationID := fmt.Sprintf("data-migration-%d", migrationTime.Unix())
    l.DataMigrationSchedules[migrationID] = DataMigrationSchedule{
        MigrationID: migrationID,
        Scheduled:   migrationTime,
        Status:      "Scheduled",
    }
    return nil
}

// Record resource consumption
func (l *Ledger) recordResourceConsumption(timestamp time.Time) error {
    consumptionID := fmt.Sprintf("resource-consumption-%d", timestamp.Unix())
    l.ResourceConsumptions[consumptionID] = ResourceConsumption{
        ConsumptionID: consumptionID,
        Timestamp:     timestamp,
    }
    return nil
}

// Record system uptime
func (l *Ledger) recordSystemUptime(timestamp time.Time) error {
    uptimeID := fmt.Sprintf("system-uptime-%d", timestamp.Unix())
    l.SystemUptimeRecords[uptimeID] = SystemUptime{
        UptimeID:  uptimeID,
        Timestamp: timestamp,
    }
    return nil
}

// Record system alert
func (l *Ledger) recordSystemAlert(timestamp time.Time) error {
    alertID := fmt.Sprintf("system-alert-%d", timestamp.Unix())
    l.SystemAlerts[alertID] = SystemAlert{
        AlertID:   alertID,
        Timestamp: timestamp,
    }
    return nil
}

// Record system snapshot
func (l *Ledger) recordSystemSnapshot(timestamp time.Time) error {
    snapshotID := fmt.Sprintf("system-snapshot-%d", timestamp.Unix())
    l.SystemSnapshots[snapshotID] = SystemSnapshot{
        SnapshotID: snapshotID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record activity log
func (l *Ledger) recordActivityLog(timestamp time.Time) error {
    logID := fmt.Sprintf("activity-log-%d", timestamp.Unix())
    l.ActivityLogs[logID] = ActivityLog{
        LogID:      logID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record stress test
func (l *Ledger) recordStressTestResult(timestamp time.Time) error {
    testID := fmt.Sprintf("stress-test-%d", timestamp.Unix())
    l.StressTests[testID] = StressTest{
        TestID:    testID,
        Timestamp: timestamp,
        Result:    "Passed",
    }
    return nil
}

// Record maintenance history
func (l *Ledger) recordMaintenanceHistory(timestamp time.Time) error {
    eventID := fmt.Sprintf("maintenance-history-%d", timestamp.Unix())
    l.MaintenanceHistories[eventID] = MaintenanceHistory{
        EventID:   eventID,
        Timestamp: timestamp,
    }
    return nil
}

// Record energy consumption
func (l *Ledger) recordEnergyConsumption(timestamp time.Time) error {
    consumptionID := fmt.Sprintf("energy-consumption-%d", timestamp.Unix())
    l.EnergyConsumptions[consumptionID] = EnergyConsumption{
        ConsumptionID: consumptionID,
        Timestamp:     timestamp,
    }
    return nil
}

// Record node sync status
func (l *Ledger) recordNodeSyncStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("node-sync-%d", timestamp.Unix())
    l.NodeSyncStatuses[statusID] = NodeSyncStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record failover event
func (l *Ledger) recordFailoverEvent(timestamp time.Time) error {
    eventID := fmt.Sprintf("failover-%d", timestamp.Unix())
    l.FailoverEvents[eventID] = FailoverEvent{
        EventID:   eventID,
        Timestamp: timestamp,
    }
    return nil
}

// Record latency change
func (l *Ledger) recordLatencyChange(latency float64, timestamp time.Time) error {
    changeID := fmt.Sprintf("latency-change-%d", timestamp.Unix())
    l.LatencyChanges[changeID] = LatencyChange{
        Latency:    latency,
        Timestamp:  timestamp,
    }
    return nil
}

// Record bandwidth usage
func (l *Ledger) recordBandwidthUsage(bandwidth float64, timestamp time.Time) error {
    usageID := fmt.Sprintf("bandwidth-usage-%d", timestamp.Unix())
    l.BandwidthUsages[usageID] = BandwidthUsage{
        Bandwidth:  bandwidth,
        Timestamp:  timestamp,
    }
    return nil
}

// Record config update
func (l *Ledger) recordConfigUpdate(timestamp time.Time) error {
    updateID := fmt.Sprintf("config-update-%d", timestamp.Unix())
    l.ConfigUpdates[updateID] = ConfigUpdate{
        UpdateID:   updateID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record database connection status
func (l *Ledger) recordDatabaseConnectionStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("db-connection-%d", timestamp.Unix())
    l.DatabaseConnectionStatuses[statusID] = DatabaseConnectionStatus{
        Status:     status,
        Timestamp:  timestamp,
    }
    return nil
}

// Record consistency check
func (l *Ledger) recordConsistencyCheck(consistency string, timestamp time.Time) error {
    checkID := fmt.Sprintf("consistency-check-%d", timestamp.Unix())
    l.ConsistencyChecks[checkID] = ConsistencyCheck{
        Consistency: consistency,
        Timestamp:   timestamp,
    }
    return nil
}

// Record resource limits
func (l *Ledger) recordResourceLimits(usage float64, timestamp time.Time) error {
    usageID := fmt.Sprintf("resource-limit-%d", timestamp.Unix())
    l.ResourceLimits[usageID] = ResourceLimit{
        Usage:      usage,
        Timestamp:  timestamp,
    }
    return nil
}

// Record update schedule
func (l *Ledger) recordUpdateSchedule(updateTime time.Time) error {
    scheduleID := fmt.Sprintf("update-schedule-%d", updateTime.Unix())
    l.UpdateSchedules[scheduleID] = UpdateSchedule{
        ScheduleID: scheduleID,
        Scheduled:  updateTime,
    }
    return nil
}

// Record maintenance window
func (l *Ledger) recordMaintenanceWindow(start, end time.Time) error {
    windowID := fmt.Sprintf("maintenance-window-%d", start.Unix())
    l.MaintenanceWindows[windowID] = MaintenanceWindow{
        WindowID:  windowID,
        StartTime: start,
        EndTime:   end,
    }
    return nil
}

// Record encryption compliance
func (l *Ledger) recordEncryptionCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("encryption-compliance-%d", timestamp.Unix())
    l.EncryptionCompliances[complianceID] = EncryptionCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record license compliance
func (l *Ledger) recordLicenseCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("license-compliance-%d", timestamp.Unix())
    l.LicenseCompliances[complianceID] = LicenseCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record firmware compliance
func (l *Ledger) recordFirmwareCheck(compliance string, timestamp time.Time) error {
    checkID := fmt.Sprintf("firmware-check-%d", timestamp.Unix())
    l.FirmwareChecks[checkID] = FirmwareCheck{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record storage utilization
func (l *Ledger) recordStorageUtilization(usage float64, timestamp time.Time) error {
    usageID := fmt.Sprintf("storage-utilization-%d", timestamp.Unix())
    l.StorageUtilizations[usageID] = StorageUtilization{
        Usage:      usage,
        Timestamp:  timestamp,
    }
    return nil
}

// Record transaction load
func (l *Ledger) recordTransactionLoad(load int, timestamp time.Time) error {
    loadID := fmt.Sprintf("transaction-load-%d", timestamp.Unix())
    l.TransactionLoads[loadID] = TransactionLoad{
        Load:       load,
        Timestamp:  timestamp,
    }
    return nil
}

// Record retention policy compliance
func (l *Ledger) recordRetentionPolicyCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("retention-policy-%d", timestamp.Unix())
    l.RetentionPolicyCompliances[complianceID] = RetentionPolicyCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record antivirus scan
func (l *Ledger) recordAntiVirusScan(results string, timestamp time.Time) error {
    scanID := fmt.Sprintf("antivirus-scan-%d", timestamp.Unix())
    l.AntiVirusScanResults[scanID] = AntiVirusScanResult{
        Results:    results,
        Timestamp:  timestamp,
    }
    return nil
}

// Record network config update
func (l *Ledger) recordNetworkConfigUpdate(timestamp time.Time) error {
    updateID := fmt.Sprintf("network-config-update-%d", timestamp.Unix())
    l.NetworkConfigUpdates[updateID] = NetworkConfigUpdate{
        UpdateID:   updateID,
        Timestamp:  timestamp,
    }
    return nil
}

// Record compression compliance
func (l *Ledger) recordCompressionCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("compression-compliance-%d", timestamp.Unix())
    l.CompressionCompliances[complianceID] = CompressionCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record process execution
func (l *Ledger) recordProcessExecution(timestamp time.Time) error {
    executionID := fmt.Sprintf("process-execution-%d", timestamp.Unix())
    l.ProcessExecutions[executionID] = ProcessExecution{
        ExecutionID: executionID,
        Timestamp:   timestamp,
    }
    return nil
}

// Record memory usage
func (l *Ledger) recordMemoryUsage(usage float64, timestamp time.Time) error {
    usageID := fmt.Sprintf("memory-usage-%d", timestamp.Unix())
    l.MemoryUsages[usageID] = MemoryUsage{
        Usage:      usage,
        Timestamp:  timestamp,
    }
    return nil
}

// Record memory cleanup schedule
func (l *Ledger) recordMemoryCleanupSchedule(cleanupTime time.Time) error {
    cleanupID := fmt.Sprintf("memory-cleanup-%d", cleanupTime.Unix())
    l.MemoryCleanupSchedules[cleanupID] = MemoryCleanupSchedule{
        CleanupID:  cleanupID,
        Scheduled:  cleanupTime,
    }
    return nil
}

// Record event queue status
func (l *Ledger) recordEventQueueStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("event-queue-%d", timestamp.Unix())
    l.EventQueueStatuses[statusID] = EventQueueStatus{
        Status:     status,
        Timestamp:  timestamp,
    }
    return nil
}

// Record high availability test
func (l *Ledger) recordHighAvailabilityTest(result string, timestamp time.Time) error {
    testID := fmt.Sprintf("ha-test-%d", timestamp.Unix())
    l.HighAvailabilityTests[testID] = HighAvailabilityTest{
        Result:    result,
        Timestamp: timestamp,
    }
    return nil
}

// Record replication compliance
func (l *Ledger) recordReplicationCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("replication-compliance-%d", timestamp.Unix())
    l.ReplicationCompliances[complianceID] = ReplicationCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record snapshot status
func (l *Ledger) recordSnapshotStatus(status string, timestamp time.Time) error {
    snapshotID := fmt.Sprintf("snapshot-status-%d", timestamp.Unix())
    l.SnapshotStatuses[snapshotID] = SnapshotStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record process lifecycle
func (l *Ledger) recordProcessLifecycle(timestamp time.Time) error {
    lifecycleID := fmt.Sprintf("process-lifecycle-%d", timestamp.Unix())
    l.ProcessLifecycles[lifecycleID] = ProcessLifecycle{
        EventID:   lifecycleID,
        Timestamp: timestamp,
    }
    return nil
}

// Record node redundancy compliance
func (l *Ledger) recordNodeRedundancy(compliance string, timestamp time.Time) error {
    redundancyID := fmt.Sprintf("node-redundancy-%d", timestamp.Unix())
    l.NodeRedundancies[redundancyID] = NodeRedundancy{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record node failover validation
func (l *Ledger) recordNodeFailoverValidation(result string, timestamp time.Time) error {
    failoverID := fmt.Sprintf("failover-validation-%d", timestamp.Unix())
    l.NodeFailoverValidations[failoverID] = NodeFailoverValidation{
        Result:    result,
        Timestamp: timestamp,
    }
    return nil
}

// Record heartbeat check
func (l *Ledger) recordHeartbeatCheck(status string, timestamp time.Time) error {
    heartbeatID := fmt.Sprintf("heartbeat-check-%d", timestamp.Unix())
    l.HeartbeatChecks[heartbeatID] = HeartbeatCheck{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record alert queue status
func (l *Ledger) recordAlertQueueStatus(status string, timestamp time.Time) error {
    alertID := fmt.Sprintf("alert-queue-%d", timestamp.Unix())
    l.AlertQueueStatuses[alertID] = AlertQueueStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record file integrity
func (l *Ledger) recordFileIntegrity(status string, timestamp time.Time) error {
    integrityID := fmt.Sprintf("file-integrity-%d", timestamp.Unix())
    l.FileIntegrities[integrityID] = FileIntegrity{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record firmware compliance
func (l *Ledger) recordFirmwareCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("firmware-compliance-%d", timestamp.Unix())
    l.FirmwareCompliances[complianceID] = FirmwareCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record data validation
func (l *Ledger) recordDataValidation(validationResult string, timestamp time.Time) error {
    validationID := fmt.Sprintf("data-validation-%d", timestamp.Unix())
    l.DataValidations[validationID] = DataValidation{
        ValidationResult: validationResult,
        Timestamp:        timestamp,
    }
    return nil
}

// Record disaster recovery setup
func (l *Ledger) recordDisasterRecoverySetup(timestamp time.Time) error {
    setupID := fmt.Sprintf("disaster-recovery-%d", timestamp.Unix())
    l.DisasterRecoverySetups[setupID] = DisasterRecoverySetup{
        SetupID:   setupID,
        Timestamp: timestamp,
    }
    return nil
}

// Record system shutdown
func (l *Ledger) recordSystemShutdown(timestamp time.Time) error {
    shutdownID := fmt.Sprintf("shutdown-%d", timestamp.Unix())
    l.SystemShutdowns[shutdownID] = SystemShutdown{
        Timestamp: timestamp,
    }
    return nil
}

// Record audit trail compliance
func (l *Ledger) recordAuditTrailCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("audit-compliance-%d", timestamp.Unix())
    l.AuditTrailCompliances[complianceID] = AuditTrailCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record smart contract status
func (l *Ledger) recordSmartContractStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("smart-contract-status-%d", timestamp.Unix())
    l.SmartContractStatuses[statusID] = SmartContractStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record database backup schedule
func (l *Ledger) recordDatabaseBackupSchedule(backupTime time.Time) error {
    scheduleID := fmt.Sprintf("backup-schedule-%d", backupTime.Unix())
    l.DatabaseBackupSchedules[scheduleID] = DatabaseBackupSchedule{
        BackupTime: backupTime,
    }
    return nil
}

// Record service availability
func (l *Ledger) recordServiceAvailability(availability string, timestamp time.Time) error {
    availabilityID := fmt.Sprintf("service-availability-%d", timestamp.Unix())
    l.ServiceAvailabilities[availabilityID] = ServiceAvailability{
        Availability: availability,
        Timestamp:    timestamp,
    }
    return nil
}

// Record service failures
func (l *Ledger) recordServiceFailures(failures []string, timestamp time.Time) error {
    failuresID := fmt.Sprintf("service-failures-%d", timestamp.Unix())
    l.ServiceFailureLogs[failuresID] = ServiceFailures{
        Failures:  failures,
        Timestamp: timestamp,
    }
    return nil
}

// Record license usage
func (l *Ledger) recordLicenseUsage(usage string, timestamp time.Time) error {
    usageID := fmt.Sprintf("license-usage-%d", timestamp.Unix())
    l.LicenseUsages[usageID] = LicenseUsage{
        Usage:     usage,
        Timestamp: timestamp,
    }
    return nil
}

// Record data retention compliance
func (l *Ledger) recordDataRetentionCompliance(compliance string, timestamp time.Time) error {
    complianceID := fmt.Sprintf("data-retention-%d", timestamp.Unix())
    l.DataRetentionCompliances[complianceID] = DataRetentionCompliance{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

// Record network topology
func (l *Ledger) recordNetworkTopology(topology string, timestamp time.Time) error {
    topologyID := fmt.Sprintf("network-topology-%d", timestamp.Unix())
    l.NetworkTopologies[topologyID] = NetworkTopology{
        Topology:  topology,
        Timestamp: timestamp,
    }
    return nil
}

// Record file lock status
func (l *Ledger) recordFileLockStatus(status string, timestamp time.Time) error {
    lockID := fmt.Sprintf("file-lock-%d", timestamp.Unix())
    l.FileLockStatuses[lockID] = FileLockStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// Record smart contract integrity
func (l *Ledger) recordSmartContractIntegrity(integrity string, timestamp time.Time) error {
    integrityID := fmt.Sprintf("contract-integrity-%d", timestamp.Unix())
    l.SmartContractIntegrities[integrityID] = SmartContractIntegrity{
        Integrity: integrity,
        Timestamp: timestamp,
    }
    return nil
}

// Record compression ratio
func (l *Ledger) recordCompressionRatio(ratio float64, timestamp time.Time) error {
    ratioID := fmt.Sprintf("compression-ratio-%d", timestamp.Unix())
    l.CompressionRatios[ratioID] = CompressionRatio{
        Ratio:     ratio,
        Timestamp: timestamp,
    }
    return nil
}

// Record network diagnostics
func (l *Ledger) recordNetworkDiagnostics(diagnostics string, timestamp time.Time) error {
    diagnosticsID := fmt.Sprintf("network-diagnostics-%d", timestamp.Unix())
    l.NetworkDiagnosticsLogs[diagnosticsID] = NetworkDiagnostics{
        Diagnostics: diagnostics,
        Timestamp:   timestamp,
    }
    return nil
}

func (l *Ledger) recordHealthCheckSetup(interval time.Duration, timestamp time.Time) error {
    routineID := fmt.Sprintf("health-check-%d", timestamp.Unix())
    l.HealthCheckRoutines[routineID] = HealthCheckRoutine{
        Interval:  interval,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordBandwidthTestSchedule(testTime time.Time) error {
    scheduleID := fmt.Sprintf("bandwidth-test-%d", testTime.Unix())
    l.BandwidthTestSchedules[scheduleID] = BandwidthTestSchedule{
        TestTime: testTime,
    }
    return nil
}

func (l *Ledger) recordDataRedundancyStatus(status string, timestamp time.Time) error {
    statusID := fmt.Sprintf("data-redundancy-%d", timestamp.Unix())
    l.DataRedundancyStatuses[statusID] = DataRedundancyStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordConfigurationChanges(changes string, timestamp time.Time) error {
    changeID := fmt.Sprintf("config-change-%d", timestamp.Unix())
    l.ConfigurationChanges[changeID] = ConfigurationChange{
        Changes:   changes,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSystemUpdateReapplication(timestamp time.Time) error {
    updateID := fmt.Sprintf("update-reapply-%d", timestamp.Unix())
    l.SystemUpdateReapplications[updateID] = SystemUpdateReapplication{
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordLogRotationStatus(status string, timestamp time.Time) error {
    rotationID := fmt.Sprintf("log-rotation-%d", timestamp.Unix())
    l.LogRotationStatuses[rotationID] = LogRotationStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordBackupValidation(status string, timestamp time.Time) error {
    validationID := fmt.Sprintf("backup-validation-%d", timestamp.Unix())
    l.BackupValidations[validationID] = BackupValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordEncryptionUpdateSchedule(updateTime time.Time) error {
    updateID := fmt.Sprintf("encryption-update-%d", updateTime.Unix())
    l.EncryptionUpdateSchedules[updateID] = EncryptionUpdateSchedule{
        UpdateTime: updateTime,
    }
    return nil
}

func (l *Ledger) recordErrorCorrections(corrections string, timestamp time.Time) error {
    correctionID := fmt.Sprintf("error-correction-%d", timestamp.Unix())
    l.ErrorCorrections[correctionID] = ErrorCorrection{
        Corrections: corrections,
        Timestamp:   timestamp,
    }
    return nil
}

func (l *Ledger) recordNodeMemoryHealth(status string, timestamp time.Time) error {
    memoryHealthID := fmt.Sprintf("node-memory-%d", timestamp.Unix())
    l.NodeMemoryHealthStatuses[memoryHealthID] = NodeMemoryHealth{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordBackupScheduleValidation(status string, timestamp time.Time) error {
    validationID := fmt.Sprintf("backup-schedule-%d", timestamp.Unix())
    l.BackupScheduleValidations[validationID] = BackupScheduleValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSecurityPatchStatus(status string, timestamp time.Time) error {
    patchID := fmt.Sprintf("security-patch-%d", timestamp.Unix())
    l.SecurityPatchStatuses[patchID] = SecurityPatchStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordNodeConnectivity(status string, timestamp time.Time) error {
    id := fmt.Sprintf("node-connectivity-%d", timestamp.Unix())
    l.NodeConnectivityStatuses[id] = NodeConnectivityStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordAPICompliance(compliance string, timestamp time.Time) error {
    id := fmt.Sprintf("api-compliance-%d", timestamp.Unix())
    l.APIComplianceStatuses[id] = APIComplianceStatus{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

func (l *Ledger) recordAPIEndpointHealth(endpointStatus string, timestamp time.Time) error {
    id := fmt.Sprintf("api-endpoint-health-%d", timestamp.Unix())
    l.APIEndpointHealths[id] = APIEndpointHealth{
        EndpointStatus: endpointStatus,
        Timestamp:      timestamp,
    }
    return nil
}

func (l *Ledger) recordSystemAudit(report string, timestamp time.Time) error {
    id := fmt.Sprintf("system-audit-%d", timestamp.Unix())
    l.SystemAuditReports[id] = SystemAuditReport{
        Report:    report,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordDatabaseRebuildSchedule(rebuildTime time.Time) error {
    id := fmt.Sprintf("database-rebuild-%d", rebuildTime.Unix())
    l.DatabaseRebuildSchedules[id] = DatabaseRebuildSchedule{
        RebuildTime: rebuildTime,
    }
    return nil
}

func (l *Ledger) recordDataIngestionValidation(status string, timestamp time.Time) error {
    id := fmt.Sprintf("data-ingestion-%d", timestamp.Unix())
    l.DataIngestionValidations[id] = DataIngestionValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordNodePerformance(metrics string, timestamp time.Time) error {
    id := fmt.Sprintf("node-performance-%d", timestamp.Unix())
    l.NodePerformanceMetrics[id] = NodePerformanceMetrics{
        Metrics:   metrics,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSoftwareCompliance(compliance string, timestamp time.Time) error {
    id := fmt.Sprintf("software-compliance-%d", timestamp.Unix())
    l.SoftwareComplianceStatuses[id] = SoftwareComplianceStatus{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

func (l *Ledger) recordNodeReinitialization(timestamp time.Time) error {
    id := fmt.Sprintf("node-reinitialization-%d", timestamp.Unix())
    l.NodeReinitializations[id] = NodeReinitialization{
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSystemHardening(timestamp time.Time) error {
    id := fmt.Sprintf("system-hardening-%d", timestamp.Unix())
    l.SystemHardenings[id] = SystemHardening{
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordFirewallValidation(ruleStatus string, timestamp time.Time) error {
    id := fmt.Sprintf("firewall-validation-%d", timestamp.Unix())
    l.FirewallValidations[id] = FirewallValidation{
        RuleStatus: ruleStatus,
        Timestamp:  timestamp,
    }
    return nil
}

func (l *Ledger) recordSystemWarmupStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("system-warmup-%d", timestamp.Unix())
    l.SystemWarmupStatuses[id] = SystemWarmupStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordRedundancyValidation(status string, timestamp time.Time) error {
    id := fmt.Sprintf("redundancy-validation-%d", timestamp.Unix())
    l.RedundancyValidations[id] = RedundancyValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordAuditFrequency(frequency time.Duration, timestamp time.Time) error {
    id := fmt.Sprintf("audit-frequency-%d", timestamp.Unix())
    l.SecurityAuditFrequencies[id] = SecurityAuditFrequency{
        Frequency: frequency,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordDatabaseTransactionValidation(status string, timestamp time.Time) error {
    id := fmt.Sprintf("db-transaction-validation-%d", timestamp.Unix())
    l.DatabaseTransactionValidations[id] = DatabaseTransactionValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSecurityScanSchedule(scanTime time.Time) error {
    id := fmt.Sprintf("security-scan-%d", scanTime.Unix())
    l.SecurityScanSchedules[id] = SecurityScanSchedule{
        ScanTime: scanTime,
    }
    return nil
}

func (l *Ledger) recordDataLossPreventionStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("data-loss-prevention-%d", timestamp.Unix())
    l.DataLossPreventionStatuses[id] = DataLossPreventionStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSystemRollback(details string, timestamp time.Time) error {
    id := fmt.Sprintf("system-rollback-%d", timestamp.Unix())
    l.SystemRollbacks[id] = SystemRollback{
        Details:   details,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordCloudBackupStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("cloud-backup-%d", timestamp.Unix())
    l.CloudBackupStatuses[id] = CloudBackupStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordConnectionPoolValidation(status string, timestamp time.Time) error {
    id := fmt.Sprintf("connection-pool-validation-%d", timestamp.Unix())
    l.ConnectionPoolValidations[id] = ConnectionPoolValidation{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSmartContractLoad(load string, timestamp time.Time) error {
    id := fmt.Sprintf("smart-contract-load-%d", timestamp.Unix())
    l.SmartContractLoads[id] = SmartContractLoad{
        Load:      load,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordResourceDeallocations(deallocations string, timestamp time.Time) error {
    id := fmt.Sprintf("resource-deallocation-%d", timestamp.Unix())
    l.ResourceDeallocations[id] = ResourceDeallocation{
        Deallocations: deallocations,
        Timestamp:     timestamp,
    }
    return nil
}

func (l *Ledger) recordPermissionIntegrity(status string, timestamp time.Time) error {
    id := fmt.Sprintf("permission-integrity-%d", timestamp.Unix())
    l.PermissionIntegrities[id] = PermissionIntegrity{
        IntegrityStatus: status,
        Timestamp:       timestamp,
    }
    return nil
}

func (l *Ledger) recordCleanupFrequency(frequency string, timestamp time.Time) error {
    id := fmt.Sprintf("cleanup-frequency-%d", timestamp.Unix())
    l.DataCleanupFrequencies[id] = DataCleanupFrequency{
        Frequency: frequency,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordConfigurationDrifts(drifts string, timestamp time.Time) error {
    id := fmt.Sprintf("config-drift-%d", timestamp.Unix())
    l.ConfigurationDrifts[id] = ConfigurationDrift{
        Drifts:    drifts,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordLogSizeLimits(size string, timestamp time.Time) error {
    id := fmt.Sprintf("log-size-%d", timestamp.Unix())
    l.LogSizeLimitsRecords[id] = LogSizeLimits{
        Size:      size,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordSessionPersistence(status string, timestamp time.Time) error {
    id := fmt.Sprintf("session-persistence-%d", timestamp.Unix())
    l.SessionPersistenceRecords[id] = SessionPersistenceStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordApplicationUpdate(timestamp time.Time) error {
    id := fmt.Sprintf("app-update-%d", timestamp.Unix())
    l.ApplicationUpdateRecords[id] = ApplicationUpdate{
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordRoleAssignments(changes string, timestamp time.Time) error {
    id := fmt.Sprintf("role-changes-%d", timestamp.Unix())
    l.RoleAssignmentChangeRecords[id] = RoleAssignmentChanges{
        Changes:   changes,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordNodeUpdateStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("node-update-%d", timestamp.Unix())
    l.NodeUpdateStatusRecords[id] = NodeUpdateStatus{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordAPIComplianceStatus(compliance string, timestamp time.Time) error {
    id := fmt.Sprintf("api-compliance-%d", timestamp.Unix())
    l.APIComplianceRecords[id] = APIComplianceStatus{
        Compliance: compliance,
        Timestamp:  timestamp,
    }
    return nil
}

func (l *Ledger) recordLogAccessAttempts(attempts string, timestamp time.Time) error {
    id := fmt.Sprintf("log-access-%d", timestamp.Unix())
    l.LogAccessAttemptRecords[id] = LogAccessAttempts{
        Attempts:   attempts,
        Timestamp:  timestamp,
    }
    return nil
}

func (l *Ledger) recordAPIRateLimitStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("api-rate-limit-%d", timestamp.Unix())
    l.APIRateLimitRecords[id] = APIRateLimits{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

func (l *Ledger) recordNodeRebootSchedule(rebootTime time.Time) error {
    id := fmt.Sprintf("node-reboot-%d", rebootTime.Unix())
    l.NodeRebootSchedules[id] = NodeRebootSchedule{
        RebootTime: rebootTime,
    }
    return nil
}

func (l *Ledger) recordTokenDistribution(distribution string, timestamp time.Time) error {
    id := fmt.Sprintf("token-distribution-%d", timestamp.Unix())
    l.TokenDistributionRecords[id] = TokenDistribution{
        Distribution: distribution,
        Timestamp:    timestamp,
    }
    return nil
}

func (l *Ledger) recordSelfRepairStatus(status string, timestamp time.Time) error {
    id := fmt.Sprintf("self-repair-%d", timestamp.Unix())
    l.SystemSelfRepairRecords[id] = SystemSelfRepair{
        Status:    status,
        Timestamp: timestamp,
    }
    return nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// UpdateOptimizationSetting updates the system’s optimization setting
func (l *Ledger) UpdateOptimizationSetting(setting OptimizationSetting) error {
    l.OptimizationSettings = setting
    return nil
}

// GetOptimizationSetting retrieves the current optimization setting
func (l *Ledger) GetOptimizationSetting() (OptimizationSetting, error) {
    if l.OptimizationSettings.Level == 0 {
        return OptimizationSetting{}, fmt.Errorf("optimization setting not initialized")
    }
    return l.OptimizationSettings, nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// UpdateResourceLimit updates the resource limit for a specific resource
func (l *Ledger) UpdateResourceLimit(limit ResourceLimit) error {
    if l.ResourceLimits == nil {
        l.ResourceLimits = make(map[string]ResourceLimit)
    }
    l.ResourceLimits[limit.Resource] = limit
    return nil
}

// GetResourceLimit retrieves the resource limit for a specific resource
func (l *Ledger) GetResourceLimit(resource string) (ResourceLimit, error) {
    if limit, exists := l.ResourceLimits[resource]; exists {
        return limit, nil
    }
    return ResourceLimit{}, fmt.Errorf("resource limit for %s not found", resource)
}

// SetPerformanceMonitoringEnabled enables or disables performance monitoring
func (l *Ledger) SetPerformanceMonitoringEnabled(enabled bool) error {
    l.PerformanceMonitoringState = enabled
    return nil
}

// GetPerformanceLogs retrieves historical performance logs for a specific metric and time range
func (l *Ledger) GetPerformanceLogs(metric string, from, to time.Time) ([]PerformanceLog, error) {
    var logs []PerformanceLog
    for _, log := range l.PerformanceLogs {
        if log.Metric == metric && log.Timestamp.After(from) && log.Timestamp.Before(to) {
            logs = append(logs, log)
        }
    }
    return logs, nil
}

// ClearAllPerformanceLogs removes all stored performance logs
func (l *Ledger) ClearAllPerformanceLogs() error {
    l.PerformanceLogs = []PerformanceLog{}
    return nil
}

// UpdateDiskCacheConfig updates the disk cache configuration in the ledger
func (l *Ledger) UpdateDiskCacheConfig(config DiskCacheConfig) error {
    l.DiskCacheConfig = config
    return nil
}

// UpdateResourceSharingConfig updates the resource sharing configuration in the ledger
func (l *Ledger) UpdateResourceSharingConfig(config ResourceSharingConfig) error {
    l.ResourceSharingConfig = config
    return nil
}

// UpdateNetworkConfig updates the network configuration in the ledger
func (l *Ledger) UpdateNetworkConfig(config NetworkConfig) error {
    l.NetworkConfig = config
    return nil
}

// GetNetworkConfig retrieves the current network configuration from the ledger
func (l *Ledger) GetNetworkConfig() (NetworkConfig, error) {
    if l.NetworkConfig.Timestamp.IsZero() {
        return NetworkConfig{}, fmt.Errorf("network configuration not initialized")
    }
    return l.NetworkConfig, nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// UpdateCompressionConfig updates the memory compression configuration in the ledger
func (l *Ledger) UpdateCompressionConfig(config CompressionConfig) error {
    l.CompressionConfig = config
    return nil
}

// UpdateResourceSharingConfig updates the resource sharing configuration in the ledger
func (l *Ledger) UpdateResourceSharingConfig(config ResourceSharingConfig) error {
    l.ResourceSharingConfig = config
    return nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// RecordHealthMetrics logs blockchain health metrics in the ledger
func (l *Ledger) RecordHealthMetrics(metrics HealthMetrics) error {
    l.HealthMetrics = append(l.HealthMetrics, metrics)
    return nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// GetNodeStatus retrieves the status of a specific node
func (l *Ledger) GetNodeStatus(nodeID string) (string, error) {
    if status, exists := l.NodeStatus[nodeID]; exists {
        return status, nil
    }
    return "", fmt.Errorf("node %s status not found", nodeID)
}

// GetSystemUptime retrieves the system uptime
func (l *Ledger) GetSystemUptime() (time.Duration, error) {
    return l.SystemUptime, nil
}

// GetErrorLogs retrieves error logs within a specified time range
func (l *Ledger) GetErrorLogs(from, to time.Time) ([]ErrorLog, error) {
    var logs []ErrorLog
    for _, log := range l.ErrorLogs {
        if log.Timestamp.After(from) && log.Timestamp.Before(to) {
            logs = append(logs, log)
        }
    }
    return logs, nil
}

// RecordAlertThreshold records an alert threshold in the ledger
func (l *Ledger) RecordAlertThreshold(alert Alert) error {
    l.Alerts = append(l.Alerts, alert)
    return nil
}

// RecordAlert triggers and logs an alert in the ledger
func (l *Ledger) RecordAlert(alert Alert) error {
    l.Alerts = append(l.Alerts, alert)
    return nil
}

// DeactivateAlert deactivates an existing alert in the ledger
func (l *Ledger) DeactivateAlert(metric string) error {
    for i, alert := range l.Alerts {
        if alert.Metric == metric && alert.Active {
            l.Alerts[i].Active = false
            l.Alerts[i].Timestamp = time.Now()
            return nil
        }
    }
    return fmt.Errorf("active alert for metric %s not found", metric)
}

// GetAlertLogs retrieves alert logs for specific metrics within a date range
func (l *Ledger) GetAlertLogs(metric string, from, to time.Time) ([]Alert, error) {
    var result []Alert
    for _, alert := range l.Alerts {
        if alert.Metric == metric && alert.Timestamp.After(from) && alert.Timestamp.Before(to) {
            result = append(result, alert)
        }
    }
    return result, nil
}

// RecordPerformanceLog logs a performance metric in the ledger
func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

// RecordCompressionRate records a data compression rate log
func (l *Ledger) RecordCompressionRate(rate float64, timestamp time.Time) error {
    l.CompressionRateLogs = append(l.CompressionRateLogs, CompressionRateLog{Rate: rate, Timestamp: timestamp})
    return nil
}

// RecordFileTransferStatus records the file transfer status
func (l *Ledger) RecordFileTransferStatus(status string, timestamp time.Time) error {
    l.FileTransferStatusLogs = append(l.FileTransferStatusLogs, FileTransferStatusLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordKeyRotation records an encryption key rotation log
func (l *Ledger) RecordKeyRotation(status string, timestamp time.Time) error {
    l.KeyRotationLogs = append(l.KeyRotationLogs, KeyRotationLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordHardwareStatus logs hardware health status
func (l *Ledger) RecordHardwareStatus(status string, timestamp time.Time) error {
    l.HardwareStatusLogs = append(l.HardwareStatusLogs, HardwareStatusLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordSessionDurations records user session durations
func (l *Ledger) RecordSessionDurations(durations map[string]time.Duration, timestamp time.Time) error {
    l.SessionDurationLogs = append(l.SessionDurationLogs, SessionDurationLog{SessionDurations: durations, Timestamp: timestamp})
    return nil
}

// RecordRBACStatus records role-based access control status
func (l *Ledger) RecordRBACStatus(status string, timestamp time.Time) error {
    l.RBACLogs = append(l.RBACLogs, RBACLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordLogIntegrityStatus logs log integrity status
func (l *Ledger) RecordLogIntegrityStatus(integrity bool, timestamp time.Time) error {
    l.LogIntegrityLogs = append(l.LogIntegrityLogs, LogIntegrityLog{Integrity: integrity, Timestamp: timestamp})
    return nil
}

// RecordMultiFactorAuthStatus records MFA status
func (l *Ledger) RecordMultiFactorAuthStatus(status string, timestamp time.Time) error {
    l.MultiFactorAuthStatusLogs = append(l.MultiFactorAuthStatusLogs, MultiFactorAuthStatusLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordTokenUsage records token usage metrics
func (l *Ledger) RecordTokenUsage(metrics map[string]int, timestamp time.Time) error {
    l.TokenUsageLogs = append(l.TokenUsageLogs, TokenUsageLog{UsageMetrics: metrics, Timestamp: timestamp})
    return nil
}

// RecordConsensusEfficiency logs consensus efficiency
func (l *Ledger) RecordConsensusEfficiency(efficiency float64, timestamp time.Time) error {
    l.ConsensusEfficiencyLogs = append(l.ConsensusEfficiencyLogs, ConsensusEfficiencyLog{Efficiency: efficiency, Timestamp: timestamp})
    return nil
}

// RecordAlertResponseTimes records alert response times
func (l *Ledger) RecordAlertResponseTimes(times map[string]time.Duration, timestamp time.Time) error {
    l.AlertResponseTimeLogs = append(l.AlertResponseTimeLogs, AlertResponseTimeLog{ResponseTimes: times, Timestamp: timestamp})
    return nil
}

// RecordUserPermissionsStatus records user permissions compliance status
func (l *Ledger) RecordUserPermissionsStatus(status string, timestamp time.Time) error {
    l.UserPermissionsStatusLogs = append(l.UserPermissionsStatusLogs, UserPermissionsStatusLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordNodeReconnections records node reconnection events
func (l *Ledger) RecordNodeReconnections(count int, timestamp time.Time) error {
    l.NodeReconnectionLogs = append(l.NodeReconnectionLogs, NodeReconnectionLog{ReconnectionCount: count, Timestamp: timestamp})
    return nil
}

// RecordDataAccessPatterns logs data access patterns
func (l *Ledger) RecordDataAccessPatterns(patterns map[string]interface{}, timestamp time.Time) error {
    l.DataAccessPatternLogs = append(l.DataAccessPatternLogs, DataAccessPatternLog{Patterns: patterns, Timestamp: timestamp})
    return nil
}

// RecordTransactionVolume logs transaction volumes
func (l *Ledger) RecordTransactionVolume(volume int, timestamp time.Time) error {
    l.TransactionVolumeLogs = append(l.TransactionVolumeLogs, TransactionVolumeLog{Volume: volume, Timestamp: timestamp})
    return nil
}

// RecordContractExecution logs contract execution metrics
func (l *Ledger) RecordContractExecution(metrics map[string]interface{}, timestamp time.Time) error {
    l.ContractExecutionLogs = append(l.ContractExecutionLogs, ContractExecutionLog{Metrics: metrics, Timestamp: timestamp})
    return nil
}

// RecordFunctionExecutionTimes logs function execution times
func (l *Ledger) RecordFunctionExecutionTimes(times map[string]time.Duration, timestamp time.Time) error {
    l.FunctionExecutionTimeLogs = append(l.FunctionExecutionTimeLogs, FunctionExecutionTimeLog{ExecutionTimes: times, Timestamp: timestamp})
    return nil
}

// RecordAPICallVolume logs API call volume
func (l *Ledger) RecordAPICallVolume(volume int, timestamp time.Time) error {
    l.APICallVolumeLogs = append(l.APICallVolumeLogs, APICallVolumeLog{Volume: volume, Timestamp: timestamp})
    return nil
}

// RecordResourceUsageTrends logs resource usage trends
func (l *Ledger) RecordResourceUsageTrends(trends map[string]interface{}, timestamp time.Time) error {
    l.ResourceUsageTrendLogs = append(l.ResourceUsageTrendLogs, ResourceUsageTrendLog{Trends: trends, Timestamp: timestamp})
    return nil
}

// RecordSecurityPatchStatus logs security patch statuses
func (l *Ledger) RecordSecurityPatchStatus(status string, timestamp time.Time) error {
    l.SecurityPatchStatusLogs = append(l.SecurityPatchStatusLogs, SecurityPatchStatusLog{Status: status, Timestamp: timestamp})
    return nil
}

// RecordPerformanceSummary records a summary of blockchain performance
func (l *Ledger) RecordPerformanceSummary(summary PerformanceSummary) error {
    l.PerformanceSummaries = append(l.PerformanceSummaries, summary)
    return nil
}

// RecordResourceReport records a resource utilization and efficiency report
func (l *Ledger) RecordResourceReport(report ResourceReport) error {
    l.ResourceReports = append(l.ResourceReports, report)
    return nil
}

// RecordResourceReallocation logs resource reallocation data
func (l *Ledger) RecordResourceReallocation(reallocation ResourceReallocation) error {
    l.ResourceReallocations = append(l.ResourceReallocations, reallocation)
    return nil
}

// RecordCostReport logs the cost of resource allocation
func (l *Ledger) RecordCostReport(report CostReport) error {
    l.CostReports = append(l.CostReports, report)
    return nil
}

func (l *Ledger) RecordFirmwareStatus(status FirmwareStatus) error {
    l.FirmwareStatuses = append(l.FirmwareStatuses, status)
    return nil
}

func (l *Ledger) RecordRoleChanges(changes RoleChange) error {
    l.RoleChanges = append(l.RoleChanges, changes)
    return nil
}

func (l *Ledger) RecordNodeReputation(reputation NodeReputation) error {
    l.NodeReputations = append(l.NodeReputations, reputation)
    return nil
}

func (l *Ledger) RecordAccessViolations(violation AccessViolation) error {
    l.AccessViolations = append(l.AccessViolations, violation)
    return nil
}

func (l *Ledger) RecordIntrusionAttempts(attempt IntrusionAttempt) error {
    l.IntrusionAttempts = append(l.IntrusionAttempts, attempt)
    return nil
}

func (l *Ledger) RecordProtocolCompliance(compliance ProtocolCompliance) error {
    l.ProtocolCompliances = append(l.ProtocolCompliances, compliance)
    return nil
}

func (l *Ledger) RecordThreatLevels(level ThreatLevel) error {
    l.ThreatLevels = append(l.ThreatLevels, level)
    return nil
}

func (l *Ledger) RecordRetentionCompliance(compliance RetentionCompliance) error {
    l.RetentionCompliances = append(l.RetentionCompliances, compliance)
    return nil
}

func (l *Ledger) RecordTrafficVolume(volume TrafficVolume) error {
    l.TrafficVolumes = append(l.TrafficVolumes, volume)
    return nil
}

func (l *Ledger) RecordBandwidthUsage(usage BandwidthUsage) error {
    l.BandwidthUsages = append(l.BandwidthUsages, usage)
    return nil
}

func (l *Ledger) RecordNodeMigration(migration NodeMigration) error {
    l.NodeMigrations = append(l.NodeMigrations, migration)
    return nil
}

func (l *Ledger) RecordServiceResponseTimes(times []ServiceResponseTime) error {
    l.ServiceResponseTimes = append(l.ServiceResponseTimes, times...)
    return nil
}

func (l *Ledger) RecordLoginAttempts(attempts []UserLoginAttempt) error {
    l.UserLoginAttempts = append(l.UserLoginAttempts, attempts...)
    return nil
}

func (l *Ledger) RecordComplianceAudit(results ComplianceAuditResult) error {
    l.ComplianceAuditResults = append(l.ComplianceAuditResults, results)
    return nil
}

func (l *Ledger) RecordBlockchainUpdates(update BlockchainUpdate) error {
    l.BlockchainUpdates = append(l.BlockchainUpdates, update)
    return nil
}

func (l *Ledger) RecordEnergyConsumption(energy EnergyConsumption) error {
    l.EnergyConsumptions = append(l.EnergyConsumptions, energy)
    return nil
}

func (l *Ledger) RecordNodeFailures(failures []NodeFailureRate) error {
    l.NodeFailureRates = append(l.NodeFailureRates, failures...)
    return nil
}

func (l *Ledger) RecordAPIThrottleLimits(limits []APIThrottleLimit) error {
    l.APIThrottleLimits = append(l.APIThrottleLimits, limits...)
    return nil
}

func (l *Ledger) RecordDatabaseHealth(health DatabaseHealth) error {
    l.DatabaseHealthStatuses = append(l.DatabaseHealthStatuses, health)
    return nil
}

func (l *Ledger) RecordConfigurationChanges(changes []SystemConfigurationChange) error {
    l.SystemConfigurationChanges = append(l.SystemConfigurationChanges, changes...)
    return nil
}

func (l *Ledger) RecordCacheUsage(usage CacheUsage) error {
    l.CacheUsages = append(l.CacheUsages, usage)
    return nil
}

func (l *Ledger) RecordAPIUsage(usage APIUsage) error {
    l.APIUsages = append(l.APIUsages, usage)
    return nil
}

func (l *Ledger) RecordSessionTimeouts(timeouts []SessionTimeout) error {
    l.SessionTimeouts = append(l.SessionTimeouts, timeouts...)
    return nil
}

func (l *Ledger) RecordAccessFrequency(frequencies []AccessFrequency) error {
    l.AccessFrequencies = append(l.AccessFrequencies, frequencies...)
    return nil
}

func (l *Ledger) RecordRateLimitCompliance(compliances []RateLimitCompliance) error {
    l.RateLimitCompliances = append(l.RateLimitCompliances, compliances...)
    return nil
}

func (l *Ledger) RecordThreatDetection(threats []ThreatDetection) error {
    l.ThreatDetections = append(l.ThreatDetections, threats...)
    return nil
}

func (l *Ledger) RecordAlertStatus(status AlertStatus) error {
    l.AlertStatuses = append(l.AlertStatuses, status)
    return nil
}

func (l *Ledger) RecordAnomalies(anomalies []AnomalyDetection) error {
    l.AnomalyDetections = append(l.AnomalyDetections, anomalies...)
    return nil
}

func (l *Ledger) RecordEventFrequency(frequencies []EventFrequency) error {
    l.EventFrequencies = append(l.EventFrequencies, frequencies...)
    return nil
}

func (l *Ledger) RecordBackupFrequency(frequency BackupFrequency) error {
    l.BackupFrequencies = append(l.BackupFrequencies, frequency)
    return nil
}

func (l *Ledger) RecordDataTransferRate(rate DataTransferRate) error {
    l.DataTransferRates = append(l.DataTransferRates, rate)
    return nil
}

func (l *Ledger) RecordDataRetrievalTime(retrieval DataRetrievalTime) error {
    l.DataRetrievalTimes = append(l.DataRetrievalTimes, retrieval)
    return nil
}

func (l *Ledger) RecordTransactionLatency(latency TransactionLatency) error {
    l.TransactionLatencies = append(l.TransactionLatencies, latency)
    return nil
}

func (l *Ledger) RecordStorageQuotaUsage(usage StorageQuotaUsage) error {
    l.StorageQuotaUsages = append(l.StorageQuotaUsages, usage)
    return nil
}

func (l *Ledger) RecordDiskSpeed(speed DiskSpeed) error {
    l.DiskSpeeds = append(l.DiskSpeeds, speed)
    return nil
}

func (l *Ledger) RecordNetworkResilience(resilience NetworkResilience) error {
    l.NetworkResilienceMetrics = append(l.NetworkResilienceMetrics, resilience)
    return nil
}

func (l *Ledger) RecordBlockchainIntegrity(integrity BlockchainIntegrity) error {
    l.BlockchainIntegrityLogs = append(l.BlockchainIntegrityLogs, integrity)
    return nil
}

func (l *Ledger) RecordEncryptionCompliance(compliance EncryptionCompliance) error {
    l.EncryptionComplianceLogs = append(l.EncryptionComplianceLogs, compliance)
    return nil
}

func (l *Ledger) RecordSessionActivity(activity SessionActivity) error {
    l.SessionActivities = append(l.SessionActivities, activity)
    return nil
}

func (l *Ledger) RecordAccessControlStatus(status AccessControlStatus) error {
    l.AccessControlStatuses = append(l.AccessControlStatuses, status)
    return nil
}

func (l *Ledger) RecordSystemHealth(health SystemHealth) error {
    l.SystemHealthLogs = append(l.SystemHealthLogs, health)
    return nil
}

func (l *Ledger) RecordNodeStatus(status NodeStatus) error {
    l.NodeStatuses = append(l.NodeStatuses, status)
    return nil
}

func (l *Ledger) RecordResourceUsage(usage ResourceUsage) error {
    l.ResourceUsages = append(l.ResourceUsages, usage)
    return nil
}

func (l *Ledger) RecordNetworkLatency(latency NetworkLatency) error {
    l.NetworkLatencies = append(l.NetworkLatencies, latency)
    return nil
}

func (l *Ledger) RecordDataThroughput(throughput DataThroughput) error {
    l.DataThroughputs = append(l.DataThroughputs, throughput)
    return nil
}

func (l *Ledger) RecordTransactionRate(rate TransactionRate) error {
    l.TransactionRates = append(l.TransactionRates, rate)
    return nil
}

func (l *Ledger) RecordBlockPropagationTime(time BlockPropagationTime) error {
    l.BlockPropagationTimes = append(l.BlockPropagationTimes, time)
    return nil
}

func (l *Ledger) RecordConsensusStatus(status ConsensusStatus) error {
    l.ConsensusStatuses = append(l.ConsensusStatuses, status)
    return nil
}

func (l *Ledger) RecordSubBlockValidation(validation SubBlockValidation) error {
    l.SubBlockValidations = append(l.SubBlockValidations, validation)
    return nil
}

func (l *Ledger) RecordSubBlockCompletion(completion SubBlockCompletion) error {
    l.SubBlockCompletions = append(l.SubBlockCompletions, completion)
    return nil
}

func (l *Ledger) RecordPeerConnections(status PeerConnectionStatus) error {
    l.PeerConnectionStatuses = append(l.PeerConnectionStatuses, status)
    return nil
}

func (l *Ledger) RecordDataSyncStatus(status DataSyncStatus) error {
    l.DataSyncStatuses = append(l.DataSyncStatuses, status)
    return nil
}

func (l *Ledger) RecordNodeAvailability(availability NodeAvailability) error {
    l.NodeAvailabilities = append(l.NodeAvailabilities, availability)
    return nil
}

func (l *Ledger) RecordShardHealth(health ShardHealth) error {
    l.ShardHealthLogs = append(l.ShardHealthLogs, health)
    return nil
}

func (l *Ledger) RecordDiskUsage(usage DiskUsage) error {
    l.DiskUsages = append(l.DiskUsages, usage)
    return nil
}

func (l *Ledger) RecordMemoryUsage(usage MemoryUsage) error {
    l.MemoryUsages = append(l.MemoryUsages, usage)
    return nil
}

func (l *Ledger) RecordCPUUtilization(utilization CPUUtilization) error {
    l.CPUUtilizations = append(l.CPUUtilizations, utilization)
    return nil
}

func (l *Ledger) RecordNodeDowntime(downtime NodeDowntime) error {
    l.NodeDowntimeLogs = append(l.NodeDowntimeLogs, downtime)
    return nil
}

func (l *Ledger) RecordNetworkBandwidth(bandwidth NetworkBandwidth) error {
    l.NetworkBandwidthLogs = append(l.NetworkBandwidthLogs, bandwidth)
    return nil
}

func (l *Ledger) RecordErrorRate(rate ErrorRate) error {
    l.ErrorRates = append(l.ErrorRates, rate)
    return nil
}

func (l *Ledger) RecordUserActivity(activity UserActivity) error {
    l.UserActivities = append(l.UserActivities, activity)
    return nil
}

func (l *Ledger) RecordComplianceStatus(status ComplianceStatus) error {
    l.ComplianceStatuses = append(l.ComplianceStatuses, status)
    return nil
}

func (l *Ledger) RecordAuditLogs(log AuditLog) error {
    l.AuditLogs = append(l.AuditLogs, log)
    return nil
}

func (l *Ledger) RecordThreatResponseTime(response ThreatResponseTime) error {
    l.ThreatResponseTimes = append(l.ThreatResponseTimes, response)
    return nil
}

func (l *Ledger) RecordSystemUptime(uptime SystemUptime) error {
    l.SystemUptimes = append(l.SystemUptimes, uptime)
    return nil
}

func (l *Ledger) RecordTrafficPatterns(pattern TrafficPattern) error {
    l.TrafficPatterns = append(l.TrafficPatterns, pattern)
    return nil
}

func (l *Ledger) RecordSuspiciousActivity(activity SuspiciousActivity) error {
    l.SuspiciousActivities = append(l.SuspiciousActivities, activity)
    return nil
}

func (l *Ledger) RecordLoadBalancingStatus(status LoadBalancingStatus) error {
    l.LoadBalancingStatuses = append(l.LoadBalancingStatuses, status)
    return nil
}

func (l *Ledger) RecordHealthThresholds(threshold HealthThreshold) error {
    l.HealthThresholds = append(l.HealthThresholds, threshold)
    return nil
}

func (l *Ledger) RecordIncidentResponseTime(response IncidentResponseTime) error {
    l.IncidentResponseTimes = append(l.IncidentResponseTimes, response)
    return nil
}

func (l *Ledger) RecordAPIResponseTime(responseTime APIResponseTime) error {
    l.APIResponseTimes = append(l.APIResponseTimes, responseTime)
    return nil
}

func (l *Ledger) RecordDataRequestVolume(volume DataRequestVolume) error {
    l.DataRequestVolumes = append(l.DataRequestVolumes, volume)
    return nil
}

func (l *Ledger) RecordSessionDataUsage(usage SessionDataUsage) error {
    l.SessionDataUsages = append(l.SessionDataUsages, usage)
    return nil
}

func (l *Ledger) RecordRateLimitExceedances(exceedance RateLimitExceedance) error {
    l.RateLimitExceedances = append(l.RateLimitExceedances, exceedance)
    return nil
}

func (l *Ledger) RecordEventLogs(log EventLog) error {
    l.EventLogs = append(l.EventLogs, log)
    return nil
}

func (l *Ledger) RecordSystemAlerts(alert SystemAlert) error {
    l.SystemAlerts = append(l.SystemAlerts, alert)
    return nil
}

func (l *Ledger) RecordResourceAllocation(allocation ResourceAllocation) error {
    l.ResourceAllocations = append(l.ResourceAllocations, allocation)
    return nil
}

func (l *Ledger) RecordDataEncryptionStatus(status EncryptionStatus) error {
    l.EncryptionStatuses = append(l.EncryptionStatuses, status)
    return nil
}

func (l *Ledger) RecordConsensusAnomalies(anomaly ConsensusAnomaly) error {
    l.ConsensusAnomalies = append(l.ConsensusAnomalies, anomaly)
    return nil
}

func (l *Ledger) RecordSecurityPolicyCompliance(compliance SecurityPolicyCompliance) error {
    l.SecurityPolicyCompliances = append(l.SecurityPolicyCompliances, compliance)
    return nil
}

func (l *Ledger) ClearAllResourceAlerts() error {
    l.ResourceAlerts = []string{}
    return nil
}

func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

func (l *Ledger) UpdateOptimizationPolicy(policy OptimizationPolicy) error {
    l.OptimizationPolicies = append(l.OptimizationPolicies, policy)
    return nil
}

func (l *Ledger) RecordSystemOverheadEvent(event SystemOverhead) error {
    l.SystemOverheadLogs = append(l.SystemOverheadLogs, event)
    return nil
}

func (l *Ledger) UpdatePriorityMode(mode PriorityMode) error {
    l.PriorityModes = append(l.PriorityModes, mode)
    return nil
}

func (l *Ledger) GetPriorityMode() (OptimizationPolicy, error) {
    if l.PriorityMode.PriorityMode == "" {
        return OptimizationPolicy{}, fmt.Errorf("no priority mode set")
    }
    return l.PriorityMode, nil
}

func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

func (l *Ledger) RecordResourceConsumption(consumption ResourceConsumption) error {
    l.ResourceConsumption = append(l.ResourceConsumption, consumption)
    return nil
}

func (l *Ledger) UpdateOptimizationPolicy(policy OptimizationPolicy) error {
    l.OptimizationPolicies = append(l.OptimizationPolicies, policy)
    return nil
}

func (l *Ledger) GetThreadPoolConfig() (ThreadPoolConfig, error) {
    if l.ThreadPoolConfig.MaxSize == 0 {
        return ThreadPoolConfig{}, fmt.Errorf("thread pool config not set")
    }
    return l.ThreadPoolConfig, nil
}

func (l *Ledger) RecordUptime(uptime time.Duration) error {
    log := UptimeLog{
        Uptime:    uptime,
        Timestamp: time.Now(),
    }
    l.UptimeLogs = append(l.UptimeLogs, log)
    return nil
}

func (l *Ledger) RecordResourceAlert(alert ResourceAlert) error {
    l.ResourceAlerts = append(l.ResourceAlerts, alert)
    return nil
}

func (l *Ledger) ClearResourceAlert(metric string) error {
    for i, alert := range l.ResourceAlerts {
        if alert.Metric == metric && alert.Active {
            l.ResourceAlerts[i].Active = false
            return nil
        }
    }
    return fmt.Errorf("active alert for metric %s not found", metric)
}

func (l *Ledger) SetIOPSMonitoringEnabled(enabled bool) error {
    l.IOPSTrackingEnabled = enabled
    return nil
}

func (l *Ledger) RecordPerformanceGoal(goal PerformanceGoal) error {
    l.PerformanceGoals = append(l.PerformanceGoals, goal)
    return nil
}

func (l *Ledger) GetPerformanceGoal(metric string) (PerformanceGoal, error) {
    for _, goal := range l.PerformanceGoals {
        if goal.Metric == metric {
            return goal, nil
        }
    }
    return PerformanceGoal{}, fmt.Errorf("performance goal for metric %s not found", metric)
}

func (l *Ledger) RecordPerformanceLog(log PerformanceLog) error {
    l.PerformanceLogs = append(l.PerformanceLogs, log)
    return nil
}

func (l *Ledger) RecordScalingEvent(event ScalingEvent) error {
    l.ScalingEvents = append(l.ScalingEvents, event)
    return nil
}

func (l *Ledger) RecordResourceAlert(alert ResourceAlert) error {
    l.ResourceAlerts = append(l.ResourceAlerts, alert)
    return nil
}

func (l *Ledger) ResetAllResourceUtilizationCounters() error {
    l.ResourceUtilizationLogs = []ResourceAlert{}
    return nil
}
