package monitoring_and_performance

import (
	"fmt"
	"time"
	"synnergy_network/pkg/common"
	"synnergy_network/pkg/encryption"
	"synnergy_network/pkg/ledger"
	"synnergy_network/pkg/network"
)

func MonitorSystemHealth(ledgerInstance *Ledger) error {
    systemHealth := network.GetSystemHealth()
    err := ledgerInstance.RecordSystemHealth(SystemHealth{
        Status:    systemHealth.Status,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("system health monitoring failed: %v", err)
    }
    fmt.Println("System health monitored:", systemHealth)
    return nil
}

func MonitorNodeStatus(ledgerInstance *Ledger) error {
    nodeStatus := network.GetNodeStatus()
    err := ledgerInstance.RecordNodeStatus(NodeStatus{
        NodeID:    nodeStatus.NodeID,
        Status:    nodeStatus.Status,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("node status monitoring failed: %v", err)
    }
    fmt.Println("Node status monitored:", nodeStatus)
    return nil
}

func TrackResourceUsage(ledgerInstance *Ledger) error {
    resourceUsage := network.GetResourceUsage()
    err := ledgerInstance.RecordResourceUsage(ResourceUsage{
        ResourceType: resourceUsage.Type,
        Usage:        resourceUsage.Usage,
        Timestamp:    time.Now(),
    })
    if err != nil {
        return fmt.Errorf("resource usage tracking failed: %v", err)
    }
    fmt.Println("Resource usage tracked:", resourceUsage)
    return nil
}

func MonitorNetworkLatency(ledgerInstance *Ledger) error {
    networkLatency := network.GetNetworkLatency()
    err := ledgerInstance.RecordNetworkLatency(NetworkLatency{
        Latency:   networkLatency.Latency,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("network latency monitoring failed: %v", err)
    }
    fmt.Println("Network latency monitored:", networkLatency)
    return nil
}

func MonitorDataThroughput(ledgerInstance *Ledger) error {
    dataThroughput := network.GetDataThroughput()
    err := ledgerInstance.RecordDataThroughput(DataThroughput{
        Throughput: dataThroughput.Throughput,
        Timestamp:  time.Now(),
    })
    if err != nil {
        return fmt.Errorf("data throughput monitoring failed: %v", err)
    }
    fmt.Println("Data throughput monitored:", dataThroughput)
    return nil
}

func MonitorTransactionRate(ledgerInstance *Ledger) error {
    transactionRate := network.GetTransactionRate()
    err := ledgerInstance.RecordTransactionRate(TransactionRate{
        Rate:      transactionRate.Rate,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("transaction rate monitoring failed: %v", err)
    }
    fmt.Println("Transaction rate monitored:", transactionRate)
    return nil
}

func TrackBlockPropagationTime(ledgerInstance *Ledger) error {
    blockPropagationTime := network.GetBlockPropagationTime()
    err := ledgerInstance.RecordBlockPropagationTime(BlockPropagationTime{
        Time:      blockPropagationTime.Time,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("block propagation time tracking failed: %v", err)
    }
    fmt.Println("Block propagation time tracked:", blockPropagationTime)
    return nil
}

func MonitorConsensusStatus(ledgerInstance *Ledger) error {
    consensusStatus := network.GetConsensusStatus()
    err := ledgerInstance.RecordConsensusStatus(ConsensusStatus{
        Status:    consensusStatus.Status,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("consensus status monitoring failed: %v", err)
    }
    fmt.Println("Consensus status monitored:", consensusStatus)
    return nil
}

func MonitorSubBlockValidation(ledgerInstance *Ledger) error {
    subBlockValidation := network.GetSubBlockValidationStatus()
    err := ledgerInstance.RecordSubBlockValidation(SubBlockValidation{
        ValidationID: subBlockValidation.ID,
        Status:       subBlockValidation.Status,
        Timestamp:    time.Now(),
    })
    if err != nil {
        return fmt.Errorf("sub-block validation monitoring failed: %v", err)
    }
    fmt.Println("Sub-block validation monitored:", subBlockValidation)
    return nil
}

func TrackSubBlockCompletion(ledgerInstance *Ledger) error {
    subBlockCompletion := network.GetSubBlockCompletionTime()
    err := ledgerInstance.RecordSubBlockCompletion(SubBlockCompletion{
        BlockID:    subBlockCompletion.ID,
        Completion: subBlockCompletion.Time,
        Timestamp:  time.Now(),
    })
    if err != nil {
        return fmt.Errorf("sub-block completion tracking failed: %v", err)
    }
    fmt.Println("Sub-block completion tracked:", subBlockCompletion)
    return nil
}

func MonitorPeerConnections(ledgerInstance *Ledger) error {
    peerConnections := network.GetPeerConnectionStatus()
    err := ledgerInstance.RecordPeerConnections(PeerConnectionStatus{
        PeerID:    peerConnections.ID,
        Status:    peerConnections.Status,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("peer connections monitoring failed: %v", err)
    }
    fmt.Println("Peer connections monitored:", peerConnections)
    return nil
}


func MonitorDataSyncStatus(ledgerInstance *Ledger) error {
    dataSyncStatus := network.GetDataSyncStatus()
    err := ledgerInstance.RecordDataSyncStatus(DataSyncStatus{
        Status:    dataSyncStatus.Status,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("data sync status monitoring failed: %v", err)
    }
    fmt.Println("Data sync status monitored:", dataSyncStatus)
    return nil
}

func TrackNodeAvailability(ledgerInstance *Ledger) error {
    nodeAvailability := network.GetNodeAvailability()
    err := ledgerInstance.RecordNodeAvailability(NodeAvailability{
        NodeID:    nodeAvailability.NodeID,
        Available: nodeAvailability.Available,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("node availability tracking failed: %v", err)
    }
    fmt.Println("Node availability tracked:", nodeAvailability)
    return nil
}

func MonitorShardHealth(ledgerInstance *Ledger) error {
    shardHealth := network.GetShardHealthStatus()
    err := ledgerInstance.RecordShardHealth(ShardHealth{
        ShardID:   shardHealth.ShardID,
        Health:    shardHealth.Health,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("shard health monitoring failed: %v", err)
    }
    fmt.Println("Shard health monitored:", shardHealth)
    return nil
}

func TrackDiskUsage(ledgerInstance *Ledger) error {
    diskUsage := network.GetDiskUsage()
    err := ledgerInstance.RecordDiskUsage(DiskUsage{
        TotalSpace:  diskUsage.TotalSpace,
        UsedSpace:   diskUsage.UsedSpace,
        FreeSpace:   diskUsage.FreeSpace,
        Timestamp:   time.Now(),
    })
    if err != nil {
        return fmt.Errorf("disk usage tracking failed: %v", err)
    }
    fmt.Println("Disk usage tracked:", diskUsage)
    return nil
}

func MonitorMemoryUsage(ledgerInstance *Ledger) error {
    memoryUsage := network.GetMemoryUsage()
    err := ledgerInstance.RecordMemoryUsage(MemoryUsage{
        TotalMemory: memoryUsage.TotalMemory,
        UsedMemory:  memoryUsage.UsedMemory,
        FreeMemory:  memoryUsage.FreeMemory,
        Timestamp:   time.Now(),
    })
    if err != nil {
        return fmt.Errorf("memory usage monitoring failed: %v", err)
    }
    fmt.Println("Memory usage monitored:", memoryUsage)
    return nil
}

func TrackCPUUtilization(ledgerInstance *Ledger) error {
    cpuUtilization := network.GetCPUUtilization()
    err := ledgerInstance.RecordCPUUtilization(CPUUtilization{
        Usage:     cpuUtilization.Usage,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("CPU utilization tracking failed: %v", err)
    }
    fmt.Println("CPU utilization tracked:", cpuUtilization)
    return nil
}

func MonitorNodeDowntime(ledgerInstance *Ledger) error {
    nodeDowntime := network.GetNodeDowntime()
    err := ledgerInstance.RecordNodeDowntime(NodeDowntime{
        NodeID:    nodeDowntime.NodeID,
        Downtime:  nodeDowntime.Downtime,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("node downtime monitoring failed: %v", err)
    }
    fmt.Println("Node downtime monitored:", nodeDowntime)
    return nil
}

func MonitorNetworkBandwidth(ledgerInstance *Ledger) error {
    networkBandwidth := network.GetNetworkBandwidth()
    err := ledgerInstance.RecordNetworkBandwidth(NetworkBandwidth{
        Bandwidth: networkBandwidth.Bandwidth,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("network bandwidth monitoring failed: %v", err)
    }
    fmt.Println("Network bandwidth monitored:", networkBandwidth)
    return nil
}

func TrackErrorRate(ledgerInstance *Ledger) error {
    errorRate := network.GetErrorRate()
    err := ledgerInstance.RecordErrorRate(ErrorRate{
        Rate:      errorRate.Rate,
        Timestamp: time.Now(),
    })
    if err != nil {
        return fmt.Errorf("error rate tracking failed: %v", err)
    }
    fmt.Println("Error rate tracked:", errorRate)
    return nil
}
