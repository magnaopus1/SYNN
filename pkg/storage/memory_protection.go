package storage

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

var protectedMemoryRegions = make(map[int]int) // Start address mapped to size
var volatileMemoryRegions = make(map[int]int)   // Start address mapped to size
var mappedAddresses = make(map[int]int)         // Start address mapped to size
var memoryProtectionLock sync.Mutex

// LockMemoryRegion: Locks a specified memory region for write protection
func LockMemoryRegion(start, size int) error {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    protectedMemoryRegions[start] = size
    LogMemoryProtection("LockMemoryRegion", fmt.Sprintf("Locked memory region from %d to %d", start, start+size-1))
    return nil
}

// UnlockMemoryRegion: Unlocks a previously locked memory region
func UnlockMemoryRegion(start int) error {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if _, exists := protectedMemoryRegions[start]; !exists {
        LogMemoryProtection("UnlockMemoryRegion", "Memory region not locked")
        return errors.New("memory region not locked")
    }
    delete(protectedMemoryRegions, start)
    LogMemoryProtection("UnlockMemoryRegion", fmt.Sprintf("Unlocked memory region starting at %d", start))
    return nil
}

// MemoryWriteProtected: Verifies if a memory region is write-protected
func MemoryWriteProtected(start, size int) bool {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    for regionStart, regionSize := range protectedMemoryRegions {
        if start >= regionStart && start+size <= regionStart+regionSize {
            LogMemoryProtection("MemoryWriteProtected", "Memory is write-protected")
            return true
        }
    }
    LogMemoryProtection("MemoryWriteProtected", "Memory is not write-protected")
    return false
}

// MemoryWriteVolatile: Sets a memory region as volatile (temporary changes allowed)
func MemoryWriteVolatile(start, size int) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    volatileMemoryRegions[start] = size
    LogMemoryProtection("MemoryWriteVolatile", fmt.Sprintf("Set memory region from %d to %d as volatile", start, start+size-1))
}

// MemoryMapAddress: Maps a memory address to allow access
func MemoryMapAddress(start, size int) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    mappedAddresses[start] = size
    LogMemoryProtection("MemoryMapAddress", fmt.Sprintf("Mapped memory address from %d to %d", start, start+size-1))
}

// UnmapMemoryAddress: Unmaps a previously mapped memory address
func UnmapMemoryAddress(start int) error {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if _, exists := mappedAddresses[start]; !exists {
        LogMemoryProtection("UnmapMemoryAddress", "Memory address not mapped")
        return errors.New("memory address not mapped")
    }
    delete(mappedAddresses, start)
    LogMemoryProtection("UnmapMemoryAddress", fmt.Sprintf("Unmapped memory address starting at %d", start))
    return nil
}

// CheckMemoryAllocation: Confirms whether a memory region is allocated
func CheckMemoryAllocation(start, size int) bool {
    allocated := start+size <= cap(memory)
    LogMemoryProtection("CheckMemoryAllocation", fmt.Sprintf("Memory allocation check: %t", allocated))
    return allocated
}

// GetMemoryRegion: Retrieves a specified memory region's data
func GetMemoryRegion(start, size int) ([]byte, error) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if start+size > len(memory) {
        LogMemoryProtection("GetMemoryRegion", "Requested region exceeds memory bounds")
        return nil, errors.New("requested region exceeds memory bounds")
    }
    data := memory[start : start+size]
    LogMemoryProtection("GetMemoryRegion", fmt.Sprintf("Retrieved memory region from %d to %d", start, start+size-1))
    return data, nil
}

// LogMemoryEvent: Logs a specific memory-related event
func LogMemoryEvent(event, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Event: " + event + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MemoryEvent", encryptedMessage)
}

// MemoryPageAllocate: Allocates a page of memory for isolated operations
func MemoryPageAllocate(size int) ([]byte, error) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if size <= 0 || size > cap(memory) {
        LogMemoryProtection("MemoryPageAllocate", "Invalid page size allocation")
        return nil, errors.New("invalid page size allocation")
    }
    page := make([]byte, size)
    LogMemoryProtection("MemoryPageAllocate", fmt.Sprintf("Allocated memory page of %d bytes", size))
    return page, nil
}

// MemoryPageFree: Frees a previously allocated page of memory
func MemoryPageFree(page []byte) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    pageSize := len(page)
    page = nil
    LogMemoryProtection("MemoryPageFree", fmt.Sprintf("Freed memory page of %d bytes", pageSize))
}

// MemoryRemap: Remaps a memory address to a new region
func MemoryRemap(oldStart, newStart, size int) error {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if _, exists := mappedAddresses[oldStart]; !exists {
        LogMemoryProtection("MemoryRemap", "Old address not mapped")
        return errors.New("old address not mapped")
    }
    delete(mappedAddresses, oldStart)
    mappedAddresses[newStart] = size
    LogMemoryProtection("MemoryRemap", fmt.Sprintf("Remapped address from %d to %d with size %d", oldStart, newStart, size))
    return nil
}

// MemoryCopyToRegion: Copies data to a protected memory region
func MemoryCopyToRegion(start int, data []byte) error {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if _, protected := protectedMemoryRegions[start]; !protected {
        LogMemoryProtection("MemoryCopyToRegion", "Memory region not write-protected")
        return errors.New("memory region not write-protected")
    }
    if start+len(data) > len(memory) {
        LogMemoryProtection("MemoryCopyToRegion", "Copy exceeds memory bounds")
        return errors.New("copy exceeds memory bounds")
    }
    copy(memory[start:], data)
    LogMemoryProtection("MemoryCopyToRegion", fmt.Sprintf("Copied %d bytes to protected region starting at %d", len(data), start))
    return nil
}

// MemoryCopyFromRegion: Copies data from a protected memory region
func MemoryCopyFromRegion(start, size int) ([]byte, error) {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    if _, protected := protectedMemoryRegions[start]; !protected {
        LogMemoryProtection("MemoryCopyFromRegion", "Memory region not write-protected")
        return nil, errors.New("memory region not write-protected")
    }
    if start+size > len(memory) {
        LogMemoryProtection("MemoryCopyFromRegion", "Copy exceeds memory bounds")
        return nil, errors.New("copy exceeds memory bounds")
    }
    data := make([]byte, size)
    copy(data, memory[start:start+size])
    LogMemoryProtection("MemoryCopyFromRegion", fmt.Sprintf("Copied %d bytes from protected region starting at %d", size, start))
    return data, nil
}

// MemoryRegionStatus: Checks the protection status of a memory region
func MemoryRegionStatus(start, size int) string {
    memoryProtectionLock.Lock()
    defer memoryProtectionLock.Unlock()
    var status string
    if _, protected := protectedMemoryRegions[start]; protected {
        status = "Protected"
    } else if _, volatile := volatileMemoryRegions[start]; volatile {
        status = "Volatile"
    } else {
        status = "Unprotected"
    }
    LogMemoryProtection("MemoryRegionStatus", fmt.Sprintf("Region status from %d to %d: %s", start, start+size-1, status))
    return status
}

// Helper Functions

// LogMemoryProtection: Logs memory protection operations with encryption
func LogMemoryProtection(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MemoryProtectionOperation", encryptedMessage)
}
