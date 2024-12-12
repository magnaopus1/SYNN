package common


import (
	"sync"
	"time"
	"synnergy_network/pkg/ledger"
)




// SnapshotManager handles the creation and restoration of blockchain snapshots.
type SnapshotManager struct {
	LedgerInstance   *ledger.Ledger          // Ledger instance to record all snapshots
	SnapshotStorage  map[string][]byte       // Stores encrypted snapshots
	CurrentState     *BlockchainState        // Current blockchain state
	mutex            sync.Mutex              // Ensures thread-safe snapshot creation and restoration
}

// BlockchainState represents the state of the blockchain at a given point in time.
type BlockchainState struct {
	BlockHeight int64        // The block height at the time of the snapshot
	Hash        string       // The hash of the blockchain state
	Timestamp   time.Time    // The timestamp of the blockchain snapshot
}






