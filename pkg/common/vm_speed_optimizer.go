package common

import (
	"sync"
	"time"
	"container/heap"
)



// SpeedOptimizer optimizes the speed of instruction execution.
type SpeedOptimizer struct {
	InstructionCache   map[string]Instruction      // Cache for frequently used instructions
	CacheMutex         sync.RWMutex                // Mutex for cache access
	PriorityQueue      *InstructionPriorityQueue   // Priority queue for instructions
	PriorityQueueMutex sync.Mutex                  // Mutex for priority queue access
	MaxCacheSize       int                         // Maximum size of the instruction cache
	CacheExpiration    time.Duration               // Duration after which cache entries expire
}

// NewSpeedOptimizer initializes a new SpeedOptimizer.
func NewSpeedOptimizer(maxCacheSize int, cacheExpiration time.Duration) *SpeedOptimizer {
	optimizer := &SpeedOptimizer{
		InstructionCache: make(map[string]Instruction),
		PriorityQueue:    NewInstructionPriorityQueue(),
		MaxCacheSize:     maxCacheSize,
		CacheExpiration:  cacheExpiration,
	}
	return optimizer
}

// OptimizeInstruction optimizes the given instruction.
func (optimizer *SpeedOptimizer) OptimizeInstruction(instr Instruction) Instruction {
	// Check cache
	optimizer.CacheMutex.RLock()
	cachedInstr, exists := optimizer.InstructionCache[instr.ID]
	optimizer.CacheMutex.RUnlock()

	if exists && time.Since(cachedInstr.Timestamp) < optimizer.CacheExpiration {
		// Return cached instruction
		return cachedInstr
	}

	// Update priority based on some criteria (e.g., urgency, frequency)
	instr.Priority = optimizer.calculatePriority(instr)

	// Add to cache
	optimizer.addToCache(instr)

	// Add to priority queue
	optimizer.addToPriorityQueue(instr)

	return instr
}

// calculatePriority calculates the priority of an instruction.
func (optimizer *SpeedOptimizer) calculatePriority(instr Instruction) int {
	// Example priority calculation:
	// Higher complexity instructions get higher priority (lower number)
	return 10 - instr.Complexity
}

// addToCache adds an instruction to the cache.
func (optimizer *SpeedOptimizer) addToCache(instr Instruction) {
	optimizer.CacheMutex.Lock()
	defer optimizer.CacheMutex.Unlock()

	if len(optimizer.InstructionCache) >= optimizer.MaxCacheSize {
		// Remove oldest entry
		var oldestKey string
		var oldestTime time.Time = time.Now()
		for key, val := range optimizer.InstructionCache {
			if val.Timestamp.Before(oldestTime) {
				oldestTime = val.Timestamp
				oldestKey = key
			}
		}
		delete(optimizer.InstructionCache, oldestKey)
	}

	instr.Timestamp = time.Now()
	optimizer.InstructionCache[instr.ID] = instr
}

// addToPriorityQueue adds an instruction to the priority queue.
func (optimizer *SpeedOptimizer) addToPriorityQueue(instr Instruction) {
	optimizer.PriorityQueueMutex.Lock()
	defer optimizer.PriorityQueueMutex.Unlock()

	heap.Push(optimizer.PriorityQueue, &instr)
}

// GetNextInstruction retrieves the next instruction from the priority queue.
func (optimizer *SpeedOptimizer) GetNextInstruction() (*Instruction, bool) {
	optimizer.PriorityQueueMutex.Lock()
	defer optimizer.PriorityQueueMutex.Unlock()

	if optimizer.PriorityQueue.Len() == 0 {
		return nil, false
	}

	instr := heap.Pop(optimizer.PriorityQueue).(*Instruction)
	return instr, true
}

// InstructionPriorityQueue implements a priority queue for instructions.
type InstructionPriorityQueue []*Instruction

// NewInstructionPriorityQueue creates a new InstructionPriorityQueue.
func NewInstructionPriorityQueue() *InstructionPriorityQueue {
	pq := &InstructionPriorityQueue{}
	heap.Init(pq)
	return pq
}

// Len implements heap.Interface.Len.
func (pq InstructionPriorityQueue) Len() int {
	return len(pq)
}

// Less implements heap.Interface.Less.
func (pq InstructionPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest priority (highest priority value)
	return pq[i].Priority < pq[j].Priority
}

// Swap implements heap.Interface.Swap.
func (pq InstructionPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push implements heap.Interface.Push.
func (pq *InstructionPriorityQueue) Push(x interface{}) {
	item := x.(*Instruction)
	*pq = append(*pq, item)
}

// Pop implements heap.Interface.Pop.
func (pq *InstructionPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // Avoid memory leak
	*pq = old[0 : n-1]
	return item
}
