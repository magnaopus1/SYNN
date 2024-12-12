package automations

import (
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/consensus"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/storage/ipfs"
)

const (
    IPFSIntegrationInterval = 1 * time.Hour  // Interval for monitoring IPFS integration
    MaxIPFSRetries          = 3              // Maximum retries for failed IPFS operations
    SubBlocksPerBlock       = 1000           // Number of sub-blocks in a block
)

// IPFSIntegrationAutomation manages the integration of IPFS for decentralized storage
type IPFSIntegrationAutomation struct {
    consensusSystem    *consensus.SynnergyConsensus // Reference to SynnergyConsensus struct
    ledgerInstance     *ledger.Ledger               // Ledger for logging IPFS integration events
    stateMutex         *sync.RWMutex                // Mutex for thread-safe access
    ipfsRetryCount     map[string]int               // Counter for retrying failed IPFS operations
    ipfsCycleCount     int                          // Counter for IPFS monitoring cycles
}

// NewIPFSIntegrationAutomation initializes the automation for IPFS integration
func NewIPFSIntegrationAutomation(consensusSystem *consensus.SynnergyConsensus, ledgerInstance *ledger.Ledger, stateMutex *sync.RWMutex) *IPFSIntegrationAutomation {
    return &IPFSIntegrationAutomation{
        consensusSystem:    consensusSystem,
        ledgerInstance:     ledgerInstance,
        stateMutex:         stateMutex,
        ipfsRetryCount:     make(map[string]int),
        ipfsCycleCount:     0,
    }
}

// StartIPFSIntegrationMonitoring starts the continuous loop for monitoring IPFS integration
func (automation *IPFSIntegrationAutomation) StartIPFSIntegrationMonitoring() {
    ticker := time.NewTicker(IPFSIntegrationInterval)

    go func() {
        for range ticker.C {
            automation.monitorIPFSIntegration()
        }
    }()
}

// monitorIPFSIntegration checks files for decentralized storage via IPFS
func (automation *IPFSIntegrationAutomation) monitorIPFSIntegration() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    // Fetch IPFS integration reports
    fileReports := automation.consensusSystem.FetchIPFSIntegrationReports()

    for _, report := range fileReports {
        if automation.isIPFSIntegrationRequired(report) {
            fmt.Printf("IPFS integration required for file %s. Initiating IPFS upload.\n", report.FileID)
            automation.applyIPFSIntegration(report)
        } else {
            fmt.Printf("No IPFS integration required for file %s.\n", report.FileID)
        }
    }

    automation.ipfsCycleCount++
    fmt.Printf("IPFS integration cycle #%d completed.\n", automation.ipfsCycleCount)

    if automation.ipfsCycleCount%SubBlocksPerBlock == 0 {
        automation.finalizeIPFSCycle()
    }
}

// isIPFSIntegrationRequired checks if a file should be uploaded to IPFS based on criteria
func (automation *IPFSIntegrationAutomation) isIPFSIntegrationRequired(report common.IPFSIntegrationReport) bool {
    // Example criteria: Files that are flagged for decentralization or large in size
    if report.ShouldDecentralize {
        fmt.Printf("File %s flagged for IPFS upload.\n", report.FileID)
        return true
    }
    return false
}

// applyIPFSIntegration uploads the file to IPFS and handles decentralized storage
func (automation *IPFSIntegrationAutomation) applyIPFSIntegration(report common.IPFSIntegrationReport) {
    encryptedFileData := automation.encryptIPFSData(report)

    // Upload file to IPFS (using the IPFS package for handling integration)
    cid, err := ipfs.UploadToIPFS(encryptedFileData.FileData)
    if err != nil {
        fmt.Printf("Error uploading file %s to IPFS: %v. Retrying...\n", report.FileID, err)
        automation.retryIPFSIntegration(report)
        return
    }

    fmt.Printf("File %s successfully uploaded to IPFS. CID: %s\n", report.FileID, cid)
    automation.logIPFSEvent(report, "Uploaded to IPFS", cid)
    automation.resetIPFSRetry(report.FileID)
}

// retryIPFSIntegration retries the IPFS upload process in case of failure
func (automation *IPFSIntegrationAutomation) retryIPFSIntegration(report common.IPFSIntegrationReport) {
    automation.ipfsRetryCount[report.FileID]++
    if automation.ipfsRetryCount[report.FileID] < MaxIPFSRetries {
        automation.applyIPFSIntegration(report)
    } else {
        fmt.Printf("Max retries reached for IPFS integration of file %s. Upload failed.\n", report.FileID)
        automation.logIPFSFailure(report)
    }
}

// resetIPFSRetry resets the retry count for IPFS integration actions
func (automation *IPFSIntegrationAutomation) resetIPFSRetry(fileID string) {
    automation.ipfsRetryCount[fileID] = 0
}

// finalizeIPFSCycle finalizes the IPFS integration cycle and logs the result in the ledger
func (automation *IPFSIntegrationAutomation) finalizeIPFSCycle() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    success := automation.consensusSystem.FinalizeIPFSCycle()
    if success {
        fmt.Println("IPFS integration cycle finalized successfully.")
        automation.logIPFSCycleFinalization()
    } else {
        fmt.Println("Error finalizing IPFS integration cycle.")
    }
}

// logIPFSEvent logs an IPFS upload event into the ledger
func (automation *IPFSIntegrationAutomation) logIPFSEvent(report common.IPFSIntegrationReport, eventType, cid string) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ipfs-event-%s-%s", report.FileID, eventType),
        Timestamp: time.Now().Unix(),
        Type:      "IPFS Integration Event",
        Status:    eventType,
        Details:   fmt.Sprintf("File %s uploaded to IPFS. CID: %s", report.FileID, cid),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with IPFS event for file %s.\n", report.FileID)
}

// logIPFSFailure logs the failure of an IPFS upload attempt into the ledger
func (automation *IPFSIntegrationAutomation) logIPFSFailure(report common.IPFSIntegrationReport) {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ipfs-failure-%s", report.FileID),
        Timestamp: time.Now().Unix(),
        Type:      "IPFS Integration Failure",
        Status:    "Failed",
        Details:   fmt.Sprintf("IPFS upload failed for file %s after maximum retries.", report.FileID),
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Printf("Ledger updated with IPFS failure for file %s.\n", report.FileID)
}

// logIPFSCycleFinalization logs the finalization of an IPFS integration cycle into the ledger
func (automation *IPFSIntegrationAutomation) logIPFSCycleFinalization() {
    automation.stateMutex.Lock()
    defer automation.stateMutex.Unlock()

    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("ipfs-cycle-finalization-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "IPFS Integration Cycle Finalization",
        Status:    "Finalized",
        Details:   "IPFS integration cycle finalized successfully.",
    }

    automation.ledgerInstance.AddEntry(entry)
    fmt.Println("Ledger updated with IPFS integration cycle finalization.")
}

// encryptIPFSData encrypts IPFS-related data before upload or logging
func (automation *IPFSIntegrationAutomation) encryptIPFSData(report common.IPFSIntegrationReport) common.IPFSIntegrationReport {
    encryptedData, err := encryption.EncryptData(report.FileData)
    if err != nil {
        fmt.Println("Error encrypting IPFS file data:", err)
        return report
    }

    report.EncryptedData = encryptedData
    fmt.Println("IPFS file data successfully encrypted for file:", report.FileID)
    return report
}
