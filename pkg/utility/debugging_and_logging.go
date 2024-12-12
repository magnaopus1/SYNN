package utility

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Loggers and mutex for thread-safe debugging/logging
var (
	debugMode         = false
	opcodeHistory     = []string{}
	opcodeLogging     = false
	stackLogMutex     sync.Mutex
	opcodeLogMutex    sync.Mutex
	conditionalBreaks = make(map[string]bool)
	watchpoints       = make(map[string]interface{})
)

// DebugBreakpoint triggers a breakpoint if in debug mode.
func DebugBreakpoint() {
	if debugMode {
		fmt.Println("Debug Breakpoint Hit!")
	}
}

// TraceOpcodeExecution logs the execution of a specific opcode.
func TraceOpcodeExecution(opcode string) {
	opcodeLogMutex.Lock()
	defer opcodeLogMutex.Unlock()

	if opcodeLogging {
		opcodeHistory = append(opcodeHistory, opcode)
		log.Printf("Tracing Opcode Execution: %s\n", opcode)
	}
}

// LogStackState logs the current state of the stack.
func LogStackState(stack []interface{}) {
	stackLogMutex.Lock()
	defer stackLogMutex.Unlock()

	log.Println("Stack State Dump:")
	for i, v := range stack {
		log.Printf("Stack[%d]: %v\n", i, v)
	}
}

// DumpMemory dumps the entire memory state to a log file.
func DumpMemory(memory map[string]interface{}) {
	file, err := os.Create("memory_dump.log")
	if err != nil {
		log.Printf("Error creating memory dump file: %v\n", err)
		return
	}
	defer file.Close()

	for address, value := range memory {
		file.WriteString(fmt.Sprintf("Memory[%s] = %v\n", address, value))
	}
	log.Println("Memory dumped to memory_dump.log")
}

// SetConditionalBreakpoint sets a conditional breakpoint based on a condition.
func SetConditionalBreakpoint(condition string, active bool) {
	conditionalBreaks[condition] = active
	log.Printf("Conditional Breakpoint Set: %s = %v\n", condition, active)
}

// TrackOpcodeUsage keeps track of opcode usage frequency.
func TrackOpcodeUsage(opcode string, countMap map[string]int) {
	countMap[opcode]++
	log.Printf("Opcode %s used %d times\n", opcode, countMap[opcode])
}

// StepThroughOpcodes allows step-through debugging of opcodes.
func StepThroughOpcodes(opcodes []string) {
	for _, opcode := range opcodes {
		DebugBreakpoint()
		TraceOpcodeExecution(opcode)
	}
}

// ViewCurrentOpcode shows the current opcode without advancing.
func ViewCurrentOpcode(opcodes []string, position int) string {
	if position < 0 || position >= len(opcodes) {
		return "Position out of bounds"
	}
	log.Printf("Current Opcode at position %d: %s\n", position, opcodes[position])
	return opcodes[position]
}

// EnableDebugMode enables the global debug mode.
func EnableDebugMode() {
	debugMode = true
	log.Println("Debug Mode Enabled")
}

// DisableDebugMode disables the global debug mode.
func DisableDebugMode() {
	debugMode = false
	log.Println("Debug Mode Disabled")
}

// ViewOpcodeHistory provides a list of executed opcodes.
func ViewOpcodeHistory() []string {
	log.Println("Viewing Opcode History")
	return opcodeHistory
}

// EnableOpcodeLogging starts logging opcode execution.
func EnableOpcodeLogging() {
	opcodeLogging = true
	log.Println("Opcode Logging Enabled")
}

// DisableOpcodeLogging stops logging opcode execution.
func DisableOpcodeLogging() {
	opcodeLogging = false
	log.Println("Opcode Logging Disabled")
}

// SetWatchpoint sets a watchpoint to monitor specific variable or memory address.
func SetWatchpoint(identifier string, value interface{}) {
	watchpoints[identifier] = value
	log.Printf("Watchpoint Set: %s = %v\n", identifier, value)
}

// RemoveWatchpoint removes a previously set watchpoint.
func RemoveWatchpoint(identifier string) {
	delete(watchpoints, identifier)
	log.Printf("Watchpoint Removed: %s\n", identifier)
}

// Utility function to log with timestamp
func logWithTimestamp(message string) {
	log.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), message)
}
