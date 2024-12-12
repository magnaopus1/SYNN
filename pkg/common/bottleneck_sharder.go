package common

import (
	"errors"
	"log"
	"sync"
	"synnergy_network/pkg/ledger"
	"time"
)

// BottleneckSharder manages the distribution of instructions to VM instances based on load.
type BottleneckSharder struct {
	LightVMs     []*LightVM            // Pool of light VM instances
	HeavyVMs     []*HeavyVM            // Pool of heavy VM instances
	LightQueue   chan Instruction      // Queue for light instructions
	HeavyQueue   chan Instruction      // Queue for heavy instructions
	LightLoadMap map[*LightVM]int      // Track load per light VM
	HeavyLoadMap map[*HeavyVM]int      // Track load per heavy VM
	mutex        sync.Mutex            // Mutex for thread safety
	maxLightLoad int                   // Max load per LightVM before scaling out
	maxHeavyLoad int                   // Max load per HeavyVM before scaling out
	minLightVMs  int                   // Minimum number of LightVM instances
	minHeavyVMs  int                   // Minimum number of HeavyVM instances
	ledgerInstance       *ledger.Ledger        // Reference to the ledger instance
	synnergyConsensus     *SynnergyConsensus    // Reference to the consensus instance
}

// NewBottleneckSharder initializes a new BottleneckSharder.
func NewBottleneckSharder(
	initialLightVMs int,
	initialHeavyVMs int,
	maxLightLoad int,
	maxHeavyLoad int,
	minLightVMs int,
	minHeavyVMs int,
	ledgerInstance *ledger.Ledger,
	synnergyConsensus *SynnergyConsensus,
) *BottleneckSharder {
	sharder := &BottleneckSharder{
		LightVMs:     make([]*LightVM, 0, initialLightVMs),
		HeavyVMs:     make([]*HeavyVM, 0, initialHeavyVMs),
		LightQueue:   make(chan Instruction, 1000),
		HeavyQueue:   make(chan Instruction, 1000),
		LightLoadMap: make(map[*LightVM]int),
		HeavyLoadMap: make(map[*HeavyVM]int),
		maxLightLoad: maxLightLoad,
		maxHeavyLoad: maxHeavyLoad,
		minLightVMs:  minLightVMs,
		minHeavyVMs:  minHeavyVMs,
		ledgerInstance:       ledgerInstance,         // Initialize ledgerInstance
		synnergyConsensus:    synnergyConsensus,      // Initialize synnergyConsensus
	}

	// Initialize LightVM instances
	for i := 0; i < initialLightVMs; i++ {
		lightVM, _ := NewLightVM(ledgerInstance, synnergyConsensus, true)
		sharder.LightVMs = append(sharder.LightVMs, lightVM)
		sharder.LightLoadMap[lightVM] = 0
		go sharder.startLightVM(lightVM)
	}

	// Initialize HeavyVM instances
	for i := 0; i < initialHeavyVMs; i++ {
		heavyVM, _ := NewHeavyVM(ledgerInstance, synnergyConsensus, true) // Pass correct parameters
		sharder.HeavyVMs = append(sharder.HeavyVMs, heavyVM)
		sharder.HeavyLoadMap[heavyVM] = 0
		go sharder.startHeavyVM(heavyVM)
	}

	// Start load monitoring
	go sharder.monitorLoad(ledgerInstance, synnergyConsensus)

	return sharder
}



// RouteInstruction routes an instruction to the appropriate VM instance based on load.
func (sharder *BottleneckSharder) RouteInstruction(instr Instruction) error {
	if instr.RequiresHeavy || instr.Complexity > LightVMComplexityThreshold {
		// Route to HeavyVM
		sharder.mutex.Lock()
		defer sharder.mutex.Unlock()

		if len(sharder.HeavyVMs) == 0 {
			return errors.New("no available HeavyVM instances")
		}

		// Find the least-loaded HeavyVM
		heavyVM := sharder.getLeastLoadedHeavyVM()
		sharder.HeavyLoadMap[heavyVM]++
		sharder.HeavyQueue <- instr
		return nil
	} else {
		// Route to LightVM
		sharder.mutex.Lock()
		defer sharder.mutex.Unlock()

		if len(sharder.LightVMs) == 0 {
			return errors.New("no available LightVM instances")
		}

		// Find the least-loaded LightVM
		lightVM := sharder.getLeastLoadedLightVM()
		sharder.LightLoadMap[lightVM]++
		sharder.LightQueue <- instr
		return nil
	}
}

// getLeastLoadedLightVM returns the LightVM instance with the least load.
func (sharder *BottleneckSharder) getLeastLoadedLightVM() *LightVM {
	var minLoad int = int(^uint(0) >> 1) // Max int value
	var selectedVM *LightVM
	for vm, load := range sharder.LightLoadMap {
		if load < minLoad {
			minLoad = load
			selectedVM = vm
		}
	}
	return selectedVM
}

// getLeastLoadedHeavyVM returns the HeavyVM instance with the least load.
func (sharder *BottleneckSharder) getLeastLoadedHeavyVM() *HeavyVM {
	var minLoad int = int(^uint(0) >> 1) // Max int value
	var selectedVM *HeavyVM
	for vm, load := range sharder.HeavyLoadMap {
		if load < minLoad {
			minLoad = load
			selectedVM = vm
		}
	}
	return selectedVM
}

// startLightVM starts processing instructions for a LightVM instance.
func (sharder *BottleneckSharder) startLightVM(vm *LightVM) {
	for instr := range sharder.LightQueue {
		err := vm.ExecuteInstruction(instr)
		if err != nil {
			log.Printf("Error executing instruction %s in LightVM: %v\n", instr.ID, err)
		}
		sharder.mutex.Lock()
		sharder.LightLoadMap[vm]--
		sharder.mutex.Unlock()
	}
}

// startHeavyVM starts processing instructions for a HeavyVM instance.
func (sharder *BottleneckSharder) startHeavyVM(vm *HeavyVM) {
	for instr := range sharder.HeavyQueue {
		err := vm.ExecuteInstruction(instr)
		if err != nil {
			log.Printf("Error executing instruction %s in HeavyVM: %v\n", instr.ID, err)
		}
		sharder.mutex.Lock()
		sharder.HeavyLoadMap[vm]--
		sharder.mutex.Unlock()
	}
}

// monitorLoad periodically checks the VM loads and scales out/in as needed.
func (sharder *BottleneckSharder) monitorLoad(ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			sharder.adjustLightVMs(ledgerInstance, consensus) // Adjust LightVMs with arguments
			sharder.adjustHeavyVMs()                           // Adjust HeavyVMs without arguments
		}
	}
}

// adjustLightVMs scales out or in LightVM instances based on load.
func (sharder *BottleneckSharder) adjustLightVMs(ledgerInstance *ledger.Ledger, consensus *SynnergyConsensus) {
	sharder.mutex.Lock()
	defer sharder.mutex.Unlock()

	totalLoad := 0
	for _, load := range sharder.LightLoadMap {
		totalLoad += load
	}

	averageLoad := totalLoad / len(sharder.LightVMs)

	if averageLoad > sharder.maxLightLoad {
		// Scale out
		newVM, _ := NewLightVM(ledgerInstance, consensus, true) // Pass all required arguments
		sharder.LightVMs = append(sharder.LightVMs, newVM)
		sharder.LightLoadMap[newVM] = 0
		go sharder.startLightVM(newVM)
		log.Printf("Scaled out LightVM. Total LightVMs: %d\n", len(sharder.LightVMs))
	} else if averageLoad < sharder.maxLightLoad/2 && len(sharder.LightVMs) > sharder.minLightVMs {
		// Scale in
		vmToRemove := sharder.LightVMs[len(sharder.LightVMs)-1]
		sharder.LightVMs = sharder.LightVMs[:len(sharder.LightVMs)-1]
		delete(sharder.LightLoadMap, vmToRemove)
		// In a real implementation, you'd signal the VM to stop processing
		log.Printf("Scaled in LightVM. Total LightVMs: %d\n", len(sharder.LightVMs))
	}
}

// adjustHeavyVMs scales out or in HeavyVM instances based on load.
func (sharder *BottleneckSharder) adjustHeavyVMs() {
	sharder.mutex.Lock()
	defer sharder.mutex.Unlock()

	totalLoad := 0
	for _, load := range sharder.HeavyLoadMap {
		totalLoad += load
	}

	averageLoad := totalLoad
	if len(sharder.HeavyVMs) > 0 {
		averageLoad = totalLoad / len(sharder.HeavyVMs)
	}

	if averageLoad > sharder.maxHeavyLoad {
		// Scale out
		newVM, err := NewHeavyVM(sharder.ledgerInstance, sharder.synnergyConsensus, true) // Pass correct parameters
		if err != nil {
			log.Printf("Failed to create new HeavyVM: %v", err)
			return
		}
		sharder.HeavyVMs = append(sharder.HeavyVMs, newVM)
		sharder.HeavyLoadMap[newVM] = 0
		go sharder.startHeavyVM(newVM)
		log.Printf("Scaled out HeavyVM. Total HeavyVMs: %d\n", len(sharder.HeavyVMs))
	} else if averageLoad < sharder.maxHeavyLoad/2 && len(sharder.HeavyVMs) > sharder.minHeavyVMs {
		// Scale in
		vmToRemove := sharder.HeavyVMs[len(sharder.HeavyVMs)-1]
		sharder.HeavyVMs = sharder.HeavyVMs[:len(sharder.HeavyVMs)-1]
		delete(sharder.HeavyLoadMap, vmToRemove)
		// In a real implementation, you'd signal the VM to stop processing
		log.Printf("Scaled in HeavyVM. Total HeavyVMs: %d\n", len(sharder.HeavyVMs))
	}
}

