package consensus_automations

import (
    "crypto/sha256"
    "fmt"
    "sync"
    "time"
    "synnergy_network_demo/common"
    "synnergy_network_demo/encryption"
    "synnergy_network_demo/ledger"
    "synnergy_network_demo/synnergy_consensus"
)

const (
    HashSwitchCheckInterval    = 10 * time.Minute // Interval for checking whether to switch hash algorithms
    HashSwitchThreshold        = 0.8              // Threshold for triggering the hash switch based on network load
    DefaultHashMethod          = "SHA-256"        // Default mining hash algorithm
    HashSyncBroadcastInterval  = 1 * time.Minute  // Interval for broadcasting hash method synchronization
)

// MiningHashSwitcher manages the dynamic switching between hash encryption methods for mining
type MiningHashSwitcher struct {
    currentHashMethod string                             // Currently used hash method (e.g., SHA-256, Scrypt, Argon2)
    ledgerInstance    *ledger.Ledger                     // Ledger to track mining algorithm changes
    consensusEngine   *synnergy_consensus.SynnergyConsensus // Synnergy Consensus engine for integrating changes
    stateMutex        *sync.RWMutex                      // Mutex for thread-safe operations
    syncChannel       chan string                        // Channel for receiving hash method sync broadcasts
}

// NewMiningHashSwitcher initializes a new MiningHashSwitcher with the default hash method
func NewMiningHashSwitcher(ledgerInstance *ledger.Ledger, consensusEngine *synnergy_consensus.SynnergyConsensus, stateMutex *sync.RWMutex) *MiningHashSwitcher {
    return &MiningHashSwitcher{
        currentHashMethod: DefaultHashMethod,
        ledgerInstance:    ledgerInstance,
        consensusEngine:   consensusEngine,
        stateMutex:        stateMutex,
        syncChannel:       make(chan string, 1), // Channel for receiving global hash sync updates
    }
}

// StartHashSwitchMonitoring begins periodic checks for conditions to switch hash methods
func (switcher *MiningHashSwitcher) StartHashSwitchMonitoring() {
    ticker := time.NewTicker(HashSwitchCheckInterval)

    go func() {
        for range ticker.C {
            fmt.Println("Checking whether to switch the mining hash algorithm...")
            switcher.evaluateHashSwitch()
        }
    }()

    // Start listening for global synchronization changes
    go switcher.listenForHashSyncUpdates()
}

// evaluateHashSwitch checks network conditions and decides if the mining hash method should be switched
func (switcher *MiningHashSwitcher) evaluateHashSwitch() {
    switcher.stateMutex.Lock()
    defer switcher.stateMutex.Unlock()

    // Get the current network load and performance metrics from the consensus engine
    networkLoad := switcher.consensusEngine.GetNetworkLoad()
    blockFinalizationTime := switcher.consensusEngine.GetBlockFinalizationTime()

    fmt.Printf("Network load: %.2f, Block finalization time: %.2f seconds\n", networkLoad, blockFinalizationTime)

    // Check if a switch is necessary based on network conditions
    if networkLoad > HashSwitchThreshold {
        // If network load is too high, switch to a more lightweight hash algorithm
        if switcher.currentHashMethod != "Scrypt" {
            switcher.switchHashMethod("Scrypt")
            switcher.broadcastHashSync("Scrypt")
        }
    } else if blockFinalizationTime > 10.0 {
        // If block finalization is taking too long, switch to Argon2 for higher security
        if switcher.currentHashMethod != "Argon2" {
            switcher.switchHashMethod("Argon2")
            switcher.broadcastHashSync("Argon2")
        }
    } else {
        // Default back to SHA-256 for stable conditions
        if switcher.currentHashMethod != "SHA-256" {
            switcher.switchHashMethod("SHA-256")
            switcher.broadcastHashSync("SHA-256")
        }
    }
}

// switchHashMethod dynamically changes the mining hash method and securely logs the event in the ledger
func (switcher *MiningHashSwitcher) switchHashMethod(newHashMethod string) {
    fmt.Printf("Switching mining hash method from %s to %s.\n", switcher.currentHashMethod, newHashMethod)

    // Apply the new hash method to the consensus engine
    switcher.consensusEngine.SetMiningHashMethod(newHashMethod)

    // Update the current hash method
    switcher.currentHashMethod = newHashMethod

    // Log the change in the ledger securely
    if err := switcher.logHashSwitch(newHashMethod); err != nil {
        fmt.Printf("Error logging hash switch in the ledger: %v\n", err)
    }
}

// logHashSwitch logs the mining hash switch event in the ledger securely
func (switcher *MiningHashSwitcher) logHashSwitch(newHashMethod string) error {
    entry := common.LedgerEntry{
        ID:        fmt.Sprintf("hash-switch-%d", time.Now().Unix()),
        Timestamp: time.Now().Unix(),
        Type:      "Mining Hash Switch",
        Status:    "Completed",
        Details:   fmt.Sprintf("Switched mining hash to %s", newHashMethod),
    }

    // Encrypt the ledger entry for security
    encryptedEntry, err := encryption.EncryptLedgerEntry(entry, []byte(EncryptionKey))
    if err != nil {
        return fmt.Errorf("failed to encrypt mining hash switch log: %v", err)
    }

    if err := switcher.ledgerInstance.AddEntry(encryptedEntry); err != nil {
        return fmt.Errorf("failed to store mining hash switch log in the ledger: %v", err)
    }

    fmt.Println("Mining hash switch logged in the ledger.")
    return nil
}

// ApplyHash calculates the hash based on the current mining hash method
func (switcher *MiningHashSwitcher) ApplyHash(input []byte) []byte {
    switch switcher.currentHashMethod {
    case "SHA-256":
        hash := sha256.Sum256(input)
        return hash[:]
    case "Scrypt":
        hash, err := scrypt.Key(input, []byte("synnergy_salt"), 16384, 8, 1, 32)
        if err != nil {
            fmt.Println("Error calculating Scrypt hash:", err)
        }
        return hash
    case "Argon2":
        hash := argon2.IDKey(input, []byte("synnergy_salt"), 1, 64*1024, 4, 32)
        return hash
    default:
        fmt.Println("Unknown hash method, using SHA-256 by default.")
        hash := sha256.Sum256(input)
        return hash[:]
    }
}

// TriggerManualHashSwitch allows administrators to manually switch the mining hash method
func (switcher *MiningHashSwitcher) TriggerManualHashSwitch(newHashMethod string) {
    fmt.Printf("Manually switching mining hash method to %s.\n", newHashMethod)
    switcher.switchHashMethod(newHashMethod)
    switcher.broadcastHashSync(newHashMethod)
}

// broadcastHashSync sends out a network-wide message to synchronize the hash method change across all nodes
func (switcher *MiningHashSwitcher) broadcastHashSync(newHashMethod string) {
    fmt.Printf("Broadcasting hash method switch to %s globally.\n", newHashMethod)

    // Broadcast the hash method change across the network
    switcher.syncChannel <- newHashMethod

    // Optionally, integrate with network broadcast protocols (e.g., P2P messaging) here to notify other nodes
    // For this demo, we'll use a simple channel communication to simulate a network-wide broadcast.
}

// listenForHashSyncUpdates listens for hash method synchronization updates and applies them
func (switcher *MiningHashSwitcher) listenForHashSyncUpdates() {
    for newHashMethod := range switcher.syncChannel {
        fmt.Printf("Received hash method sync update: switching to %s.\n", newHashMethod)
        switcher.switchHashMethod(newHashMethod)
    }
}
