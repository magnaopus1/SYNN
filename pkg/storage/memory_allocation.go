package storage

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

var memory = make([]byte, 1024) // Example initial size
var memoryStack = make([][]byte, 0)
var memoryLock sync.Mutex

// MemoryAllocate: Allocates a specified number of bytes in memory
func MemoryAllocate(size int) (int, error) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if len(memory)+size > cap(memory) {
        LogMemoryOperation("MemoryAllocate", "Insufficient memory for allocation")
        return 0, errors.New("insufficient memory for allocation")
    }
    oldSize := len(memory)
    memory = append(memory, make([]byte, size)...)
    LogMemoryOperation("MemoryAllocate", fmt.Sprintf("Allocated %d bytes", size))
    return oldSize, nil
}

// MemoryFree: Frees a specified number of bytes from memory, starting from the end
func MemoryFree(size int) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if size > len(memory) {
        LogMemoryOperation("MemoryFree", "Memory free size exceeds allocated memory")
        return errors.New("memory free size exceeds allocated memory")
    }
    memory = memory[:len(memory)-size]
    LogMemoryOperation("MemoryFree", fmt.Sprintf("Freed %d bytes", size))
    return nil
}

// MemoryRead: Reads bytes from memory at a specified position
func MemoryRead(position, size int) ([]byte, error) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if position+size > len(memory) {
        LogMemoryOperation("MemoryRead", "Read exceeds allocated memory")
        return nil, errors.New("read exceeds allocated memory")
    }
    data := memory[position : position+size]
    LogMemoryOperation("MemoryRead", fmt.Sprintf("Read %d bytes from position %d", size, position))
    return data, nil
}

// MemoryWrite: Writes bytes to memory at a specified position
func MemoryWrite(position int, data []byte) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if position+len(data) > len(memory) {
        LogMemoryOperation("MemoryWrite", "Write exceeds allocated memory")
        return errors.New("write exceeds allocated memory")
    }
    copy(memory[position:], data)
    LogMemoryOperation("MemoryWrite", fmt.Sprintf("Wrote %d bytes to position %d", len(data), position))
    return nil
}

// MemoryCopy: Copies bytes within memory from one position to another
func MemoryCopy(from, to, size int) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if from+size > len(memory) || to+size > len(memory) {
        LogMemoryOperation("MemoryCopy", "Copy operation exceeds allocated memory")
        return errors.New("copy operation exceeds allocated memory")
    }
    copy(memory[to:to+size], memory[from:from+size])
    LogMemoryOperation("MemoryCopy", fmt.Sprintf("Copied %d bytes from position %d to %d", size, from, to))
    return nil
}

// MemoryClear: Clears the entire memory by setting all bytes to zero
func MemoryClear() {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    for i := range memory {
        memory[i] = 0
    }
    LogMemoryOperation("MemoryClear", "Memory cleared")
}

// PushToMemoryStack: Pushes data onto the memory stack
func PushToMemoryStack(data []byte) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    memoryStack = append(memoryStack, data)
    LogMemoryOperation("PushToMemoryStack", fmt.Sprintf("Pushed %d bytes to memory stack", len(data)))
}

// PopFromMemoryStack: Pops data from the memory stack
func PopFromMemoryStack() ([]byte, error) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if len(memoryStack) == 0 {
        LogMemoryOperation("PopFromMemoryStack", "Memory stack is empty")
        return nil, errors.New("memory stack is empty")
    }
    data := memoryStack[len(memoryStack)-1]
    memoryStack = memoryStack[:len(memoryStack)-1]
    LogMemoryOperation("PopFromMemoryStack", fmt.Sprintf("Popped %d bytes from memory stack", len(data)))
    return data, nil
}

// MemoryResize: Resizes the memory array to a new size
func MemoryResize(newSize int) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if newSize < len(memory) {
        memory = memory[:newSize]
    } else if newSize > len(memory) {
        memory = append(memory, make([]byte, newSize-len(memory))...)
    }
    LogMemoryOperation("MemoryResize", fmt.Sprintf("Memory resized to %d bytes", newSize))
    return nil
}

// MemorySetValue: Sets a specific value at a specified position in memory
func MemorySetValue(position int, value byte) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if position >= len(memory) {
        LogMemoryOperation("MemorySetValue", "Position exceeds allocated memory")
        return errors.New("position exceeds allocated memory")
    }
    memory[position] = value
    LogMemoryOperation("MemorySetValue", fmt.Sprintf("Set memory at position %d to value %d", position, value))
    return nil
}

// MemoryGetSize: Returns the current size of allocated memory
func MemoryGetSize() int {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    size := len(memory)
    LogMemoryOperation("MemoryGetSize", fmt.Sprintf("Memory size: %d bytes", size))
    return size
}

// MemorySwap: Swaps bytes between two positions in memory
func MemorySwap(pos1, pos2 int) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if pos1 >= len(memory) || pos2 >= len(memory) {
        LogMemoryOperation("MemorySwap", "Swap positions exceed allocated memory")
        return errors.New("swap positions exceed allocated memory")
    }
    memory[pos1], memory[pos2] = memory[pos2], memory[pos1]
    LogMemoryOperation("MemorySwap", fmt.Sprintf("Swapped memory positions %d and %d", pos1, pos2))
    return nil
}

// MemoryCompare: Compares memory content at two positions for a specified length
func MemoryCompare(pos1, pos2, length int) (bool, error) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    if pos1+length > len(memory) || pos2+length > len(memory) {
        LogMemoryOperation("MemoryCompare", "Comparison exceeds allocated memory")
        return false, errors.New("comparison exceeds allocated memory")
    }
    result := string(memory[pos1:pos1+length]) == string(memory[pos2:pos2+length])
    LogMemoryOperation("MemoryCompare", fmt.Sprintf("Memory compared between positions %d and %d for %d bytes", pos1, pos2, length))
    return result, nil
}

// SaveMemorySnapshot: Saves a snapshot of the current memory state
func SaveMemorySnapshot() ([]byte, error) {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    snapshot := make([]byte, len(memory))
    copy(snapshot, memory)
    LogMemoryOperation("SaveMemorySnapshot", "Memory snapshot saved")
    return snapshot, nil
}

// LoadMemorySnapshot: Loads a snapshot into memory, overwriting the current state
func LoadMemorySnapshot(snapshot []byte) error {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    memory = make([]byte, len(snapshot))
    copy(memory, snapshot)
    LogMemoryOperation("LoadMemorySnapshot", "Memory snapshot loaded")
    return nil
}

// QueryMemoryUsage: Returns the current memory usage as a percentage of total capacity
func QueryMemoryUsage() float64 {
    memoryLock.Lock()
    defer memoryLock.Unlock()
    usage := (float64(len(memory)) / float64(cap(memory))) * 100
    LogMemoryOperation("QueryMemoryUsage", fmt.Sprintf("Memory usage: %.2f%%", usage))
    return usage
}

// Helper Functions

// LogMemoryOperation: Logs memory operations with encryption
func LogMemoryOperation(operation string, details string) error {
    encryptedMessage, err := encryption.Encrypt([]byte("Operation: " + operation + " - Details: " + details))
    if err != nil {
        return err
    }
    return common.ledger.LogEvent("MemoryOperation", encryptedMessage)
}
