package common

import (
	"errors"
	"log"
	"sync"
	"time"
)



// InstructionFilter analyzes incoming instructions, determines complexity, and routes them to the appropriate VM queues.
type InstructionFilter struct {
	LightQueue       chan Instruction          // Queue for light VM tasks
	HeavyQueue       chan Instruction          // Queue for heavy VM tasks
	LightThreshold   int                       // Threshold for routing to light VM
	HeavyThreshold   int                       // Threshold for routing to heavy VM
	OpcodeComplexity map[string]int            // Mapping of opcodes to their complexity levels
	OpcodeCategory   map[string]string         // Mapping of opcodes to categories (e.g., "arithmetic", "crypto")
	ComplexityCache  map[string]int            // Cache of instruction complexities
	VMRouter         *VMRouter                 // Reference to VMRouter for load information
	mutex            sync.RWMutex              // Mutex for thread-safe access
	workers          int                       // Number of worker goroutines
	instructionChan  chan Instruction          // Channel for incoming instructions
	stopChan         chan struct{}             // Channel to signal stopping of the filter
	SpeedOptimizer           *SpeedOptimizer
	ProcessingPowerOptimizer *ProcessingPowerOptimizer
}

func NewInstructionFilter(lightThreshold, heavyThreshold, workers int, vmRouter *VMRouter) *InstructionFilter {
	filter := &InstructionFilter{
		LightQueue:       make(chan Instruction, 1000),
		HeavyQueue:       make(chan Instruction, 1000),
		LightThreshold:   lightThreshold,
		HeavyThreshold:   heavyThreshold,
		OpcodeComplexity: make(map[string]int),
		OpcodeCategory:   make(map[string]string),
		ComplexityCache:  make(map[string]int),
		VMRouter:         vmRouter,
		workers:          workers,
		instructionChan:  make(chan Instruction, 1000),
		stopChan:         make(chan struct{}),
	}

	// Initialize opcode complexity and categories
	filter.initializeOpcodeData()

	// Initialize optimizers
	filter.SpeedOptimizer = NewSpeedOptimizer(1000, 10*time.Minute) // Example values
	filter.ProcessingPowerOptimizer = NewProcessingPowerOptimizer(2, 10, 75.0, 5*time.Second)
	filter.ProcessingPowerOptimizer.WorkerChan = make(chan struct{}, workers)

	// Start only the required number of worker goroutines
	for i := 0; i < workers; i++ {
		go filter.worker()
	}

	return filter
}


// initializeOpcodeData initializes the opcode complexity and category mappings.
func (filter *InstructionFilter) initializeOpcodeData() {
	// Example opcode complexity levels (1-10)
	filter.OpcodeComplexity = map[string]int{
		"ADD":      1,
		"SUB":      1,
		"MUL":      2,
		"DIV":      2,
		"SHA256":   8,
		"ECVERIFY": 9,
		"CALL":     5,
		"STORE":    3,
		"LOAD":     3,
		// Add more opcodes as needed
	}

	// Opcode categories
	filter.OpcodeCategory = map[string]string{
		"ADD":      "arithmetic",
		"SUB":      "arithmetic",
		"MUL":      "arithmetic",
		"DIV":      "arithmetic",
		"SHA256":   "cryptography",
		"ECVERIFY": "cryptography",
		"CALL":     "control_flow",
		"STORE":    "storage",
		"LOAD":     "storage",
		// Add more opcodes as needed
	}
}

// SubmitInstruction submits an instruction to the filter for processing.
func (filter *InstructionFilter) SubmitInstruction(instr Instruction) error {
	select {
	case filter.instructionChan <- instr:
		return nil
	default:
		return errors.New("instruction filter is full")
	}
}

// worker is a goroutine that processes instructions from the instructionChan
// and handles signals from the ProcessingPowerOptimizer to adjust processing power dynamically.
func (filter *InstructionFilter) worker() {
	for {
		select {
		case instr, ok := <-filter.instructionChan:
			if !ok {
				log.Println("instructionChan closed, worker stopping")
				return // Stop if instruction channel is closed
			}
			filter.processInstruction(instr)

		case _, ok := <-filter.ProcessingPowerOptimizer.WorkerChan:
			if !ok {
				log.Println("WorkerChan closed by ProcessingPowerOptimizer, worker stopping")
				return // Stop if WorkerChan is closed
			}

		case <-filter.stopChan:
			log.Println("stopChan received, worker stopping")
			return // Stop if stop signal is received
		}
	}
}


// processInstruction analyzes the instruction, optimizes it using SpeedOptimizer, 
// and routes it to the appropriate queue based on complexity and VM load.
func (filter *InstructionFilter) processInstruction(instr Instruction) {
	// Optimize instruction using SpeedOptimizer
	optimizedInstr := filter.SpeedOptimizer.OptimizeInstruction(instr)

	// Analyze optimized instruction complexity
	complexity := filter.getInstructionComplexity(optimizedInstr)

	// Route the optimized instruction based on complexity and VM load
	err := filter.routeInstruction(optimizedInstr, complexity)
	if err != nil {
		log.Printf("Failed to route instruction %s: %v\n", optimizedInstr.ID, err)
	}
}


// getInstructionComplexity returns the complexity of the instruction.
func (filter *InstructionFilter) getInstructionComplexity(instr Instruction) int {
	filter.mutex.RLock()
	if complexity, exists := filter.ComplexityCache[instr.Opcode]; exists {
		filter.mutex.RUnlock()
		return complexity
	}
	filter.mutex.RUnlock()

	// If not in cache, calculate complexity
	complexity := filter.calculateComplexity(instr)

	// Cache the complexity
	filter.mutex.Lock()
	filter.ComplexityCache[instr.Opcode] = complexity
	filter.mutex.Unlock()

	return complexity
}

// calculateComplexity calculates the complexity of an instruction.
func (filter *InstructionFilter) calculateComplexity(instr Instruction) int {
	if instr.RequiresHeavy {
		return filter.HeavyThreshold + 1
	}

	// Use predefined complexity if available
	if complexity, exists := filter.OpcodeComplexity[instr.Opcode]; exists {
		return complexity
	}

	// Default complexity
	return 5 // Medium complexity
}

// routeInstruction routes the instruction to the appropriate VM queue based on complexity and VM load.
func (filter *InstructionFilter) routeInstruction(instr Instruction, complexity int) error {
	// Check if instruction should go to LightVM
	if complexity <= filter.LightThreshold {
		// Check LightVM load
		if filter.VMRouter.GetLightVMLoad() < filter.VMRouter.MaxLightLoad {
			filter.LightQueue <- instr
			filter.VMRouter.IncreaseLightVMLoad() // Track load for LightVM
			return nil
		}
	}

	// Check if instruction should go to HeavyVM
	if complexity > filter.LightThreshold && complexity <= filter.HeavyThreshold {
		// Check HeavyVM load
		if filter.VMRouter.GetHeavyVMLoad() < filter.VMRouter.MaxHeavyLoad {
			filter.HeavyQueue <- instr
			filter.VMRouter.IncreaseHeavyVMLoad() // Track load for HeavyVM
			return nil
		}
	}

	// If both VMs are overloaded, buffer the instruction or return error
	select {
	case filter.HeavyQueue <- instr:
		filter.VMRouter.IncreaseHeavyVMLoad()
		return nil
	default:
		// If HeavyQueue is full, log an error or implement backpressure
		return errors.New("both LightVM and HeavyVM are overloaded")
	}
}




// GetLightQueue returns the channel for light VM tasks.
func (filter *InstructionFilter) GetLightQueue() chan Instruction {
	return filter.LightQueue
}

// GetHeavyQueue returns the channel for heavy VM tasks.
func (filter *InstructionFilter) GetHeavyQueue() chan Instruction {
	return filter.HeavyQueue
}

// Stop stops all workers gracefully by closing relevant channels.
func (filter *InstructionFilter) Stop() {
	log.Println("Stopping all workers in InstructionFilter")
	close(filter.stopChan)  // Signal stop to all workers
	close(filter.instructionChan) // Ensure all instructions stop processing
	filter.ProcessingPowerOptimizer.Stop() // Stops optimizer, which closes WorkerChan
}
