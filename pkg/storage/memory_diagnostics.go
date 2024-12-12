package storage

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

var memoryDiagnosticsEnabled bool
var temporaryMemory = make([]byte, 0)
var memoryCache = make(map[string][]byte)
var diagnosticsLock sync.Mutex
var memoryAccessLog = make(map[int]int) // Track access counts by position

// EnableMemoryDiagnostics: Activates memory diagnostics for tracking and logging
func EnableMemoryDiagnostics() {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    memoryDiagnosticsEnabled = true
    LogMemoryDiagnostic("EnableMemoryDiagnostics", "Memory diagnostics enabled")
}

// DisableMemoryDiagnostics: Deactivates memory diagnostics
func DisableMemoryDiagnostics() {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    memoryDiagnosticsEnabled = false
    LogMemoryDiagnostic("DisableMemoryDiagnostics", "Memory diagnostics disabled")
}

// TrackMemoryAccess: Logs each access to memory for diagnostic purposes
func TrackMemoryAccess(position int) {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    if memoryDiagnosticsEnabled {
        memoryAccessLog[position]++
        LogMemoryDiagnostic("TrackMemoryAccess", fmt.Sprintf("Access tracked at position %d", position))
    }
}

// MemoryOptimizeUsage: Optimizes memory usage by clearing unused cache
func MemoryOptimizeUsage() {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    if memoryDiagnosticsEnabled {
        memoryCache = make(map[string][]byte) // Clear unused cache
        LogMemoryDiagnostic("MemoryOptimizeUsage", "Memory cache optimized")
    }
}

// ConfigureMemoryCache: Configures a memory cache with a specified key
func ConfigureMemoryCache(key string, data []byte) {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    if memoryDiagnosticsEnabled {
        memoryCache[key] = data
        LogMemoryDiagnostic("ConfigureMemoryCache", fmt.Sprintf("Cache configured for key %s with %d bytes", key, len(data)))
    }
}

// CheckMemoryOverflow: Checks for memory overflow based on capacity
func CheckMemoryOverflow(limit int) bool {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    overflow := len(memory) > limit
    LogMemoryDiagnostic("CheckMemoryOverflow", fmt.Sprintf("Memory overflow check: %t", overflow))
    return overflow
}

// PrepareMemoryBlock: Prepares a specific block of memory for usage
func PrepareMemoryBlock(size int) ([]byte, error) {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    if size <= 0 || size > cap(memory) {
        LogMemoryDiagnostic("PrepareMemoryBlock", "Invalid memory block size")
        return nil, errors.New("invalid memory block size")
    }
    block := make([]byte, size)
    LogMemoryDiagnostic("PrepareMemoryBlock", fmt.Sprintf("Prepared memory block of %d bytes", size))
    return block, nil
}

// AllocateTemporaryMemory: Allocates a specified amount of temporary memory
func AllocateTemporaryMemory(size int) error {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    if size < 0 {
        LogMemoryDiagnostic("AllocateTemporaryMemory", "Invalid size for temporary memory")
        return errors.New("invalid size for temporary memory")
    }
    temporaryMemory = make([]byte, size)
    LogMemoryDiagnostic("AllocateTemporaryMemory", fmt.Sprintf("Allocated %d bytes of temporary memory", size))
    return nil
}

// FreeTemporaryMemory: Frees the allocated temporary memory
func FreeTemporaryMemory() {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    size := len(temporaryMemory)
    temporaryMemory = make([]byte, 0)
    LogMemoryDiagnostic("FreeTemporaryMemory", fmt.Sprintf("Freed %d bytes of temporary memory", size))
}

// ClearTemporaryMemory: Clears content in the temporary memory without freeing it
func ClearTemporaryMemory() {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    for i := range temporaryMemory {
        temporaryMemory[i] = 0
    }
    LogMemoryDiagnostic("ClearTemporaryMemory", "Cleared temporary memory content")
}

// CheckMemoryLeak: Analyzes memory usage patterns to detect potential leaks
func CheckMemoryLeak() bool {
    diagnosticsLock.Lock()
    defer diagnosticsLock.Unlock()
    leakDetected := len(memoryAccessLog) > cap(memory)/2 // Arbitrary threshold
    LogMemoryDiagnostic("CheckMemoryLeak", fmt.Sprintf("Memory leak detected: %t", leakDetected))
    return leakDetected
}

// Helper Functions

// LogMemoryDiagnostic: Logs memory diagnostic operations with encryption
func LogMemoryDiagnostic(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MemoryDiagnosticOperation", encryptedMessage)
}
