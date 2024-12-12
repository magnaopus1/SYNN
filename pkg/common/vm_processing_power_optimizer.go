package common

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// ProcessingPowerOptimizer optimizes processing power usage by adjusting worker counts.
type ProcessingPowerOptimizer struct {
	MinWorkers      int             // Minimum number of worker goroutines
	MaxWorkers      int             // Maximum number of worker goroutines
	CurrentWorkers  int             // Current number of worker goroutines
	WorkerChan      chan struct{}   // Channel to signal worker goroutines
	LoadThreshold   float64         // CPU usage threshold to scale up/down
	MonitorInterval time.Duration   // Interval for monitoring system load
	mutex           sync.Mutex      // Mutex for thread safety
	stopChan        chan struct{}   // Channel to signal stopping the optimizer
}

// NewProcessingPowerOptimizer initializes a new ProcessingPowerOptimizer.
func NewProcessingPowerOptimizer(minWorkers, maxWorkers int, loadThreshold float64, monitorInterval time.Duration) *ProcessingPowerOptimizer {
	optimizer := &ProcessingPowerOptimizer{
		MinWorkers:      minWorkers,
		MaxWorkers:      maxWorkers,
		CurrentWorkers:  minWorkers,
		WorkerChan:      make(chan struct{}, maxWorkers),
		LoadThreshold:   loadThreshold,
		MonitorInterval: monitorInterval,
		stopChan:        make(chan struct{}),
	}
	optimizer.startWorkers(minWorkers)
	go optimizer.monitorLoad()
	return optimizer
}

// startWorkers starts the specified number of worker goroutines.
func (optimizer *ProcessingPowerOptimizer) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		optimizer.WorkerChan <- struct{}{}
		go optimizer.worker()
	}
	optimizer.CurrentWorkers += numWorkers
}

// stopWorkers stops the specified number of worker goroutines.
func (optimizer *ProcessingPowerOptimizer) stopWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		optimizer.WorkerChan <- struct{}{}
	}
	optimizer.CurrentWorkers -= numWorkers
}

// worker is a goroutine that processes instructions.
func (optimizer *ProcessingPowerOptimizer) worker() {
	for {
		select {
		case <-optimizer.WorkerChan:
			// Signal to stop the worker
			return
		default:
			// Perform work
			time.Sleep(100 * time.Millisecond) // Simulate work
		}
	}
}

// monitorLoad monitors the system load and adjusts workers accordingly.
func (optimizer *ProcessingPowerOptimizer) monitorLoad() {
	ticker := time.NewTicker(optimizer.MonitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cpuUsage := optimizer.getCPUUsage()
			optimizer.adjustWorkers(cpuUsage)
		case <-optimizer.stopChan:
			return
		}
	}
}

// getCPUUsage retrieves an estimated CPU usage percentage based on active goroutines and available CPUs.
func (optimizer *ProcessingPowerOptimizer) getCPUUsage() float64 {
	// Get the number of CPUs and adjust by max concurrency (GOMAXPROCS)
	numCPUs := runtime.NumCPU()
	concurrency := float64(runtime.GOMAXPROCS(0))

	// Estimate load factor by active goroutines over available CPU concurrency
	loadFactor := float64(runtime.NumGoroutine()) / concurrency

	// Cap the CPU usage estimate to 100%
	return loadFactor * 100 / float64(numCPUs)
}

// adjustWorkers adjusts the number of worker goroutines based on CPU usage.
func (optimizer *ProcessingPowerOptimizer) adjustWorkers(cpuUsage float64) {
	optimizer.mutex.Lock()
	defer optimizer.mutex.Unlock()

	if cpuUsage > optimizer.LoadThreshold && optimizer.CurrentWorkers < optimizer.MaxWorkers {
		// Scale up
		optimizer.startWorkers(1)
		log.Printf("Scaling up workers to %d due to CPU usage %.2f%%\n", optimizer.CurrentWorkers, cpuUsage)
	} else if cpuUsage < optimizer.LoadThreshold/2 && optimizer.CurrentWorkers > optimizer.MinWorkers {
		// Scale down
		optimizer.stopWorkers(1)
		log.Printf("Scaling down workers to %d due to CPU usage %.2f%%\n", optimizer.CurrentWorkers, cpuUsage)
	}
}

// Stop gracefully stops the optimizer and all workers.
func (optimizer *ProcessingPowerOptimizer) Stop() {
	close(optimizer.stopChan)
	optimizer.stopWorkers(optimizer.CurrentWorkers)
}
