package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Instruction represents an instruction to be executed by the VM.
type Instruction struct {
	ID            string
	Opcode        string
	Complexity    int
	Payload       interface{}
	RequiresHeavy bool
	Timestamp     time.Time
	Priority      int // Lower value means higher priority
	Metadata      map[string]interface{}
}

type SendInstructionRequest struct {
	Opcode   string
	Payload  []byte
	Metadata map[string]string
}


// VMRouter routes instructions to the appropriate VM (LightVM or HeavyVM) based on instruction complexity, VM load, and node capabilities.
type VMRouter struct {
	LocalNode       *Node             // Now a Node pointer for accessing HasLightVM and HasHeavyVM
	NetworkNodes    []NetworkNode     // List of other nodes in the network
	LightVM         *LightVM          // Local LightVM instance (if any)
	HeavyVM         *HeavyVM          // Local HeavyVM instance (if any)
	LightQueue      chan Instruction  // Queue for LightVM instructions
	HeavyQueue      chan Instruction  // Queue for HeavyVM instructions
	MaxLightLoad    int               // Threshold for light VM load
	MaxHeavyLoad    int               // Threshold for heavy VM load
	currentLightLoad int              // Current number of instructions in LightVM
	currentHeavyLoad int              // Current number of instructions in HeavyVM
	lightVMLoad      int              // Current load on LightVMs
	heavyVMLoad      int              // Current load on HeavyVMs
	mutex            sync.RWMutex     // Use RWMutex to allow read locks for load getters
}

func NewVMRouter(localNode *Node, networkNodes []NetworkNode, lightVM *LightVM, heavyVM *HeavyVM, maxLightLoad, maxHeavyLoad int) *VMRouter {
	router := &VMRouter{
		LocalNode:        localNode,
		NetworkNodes:     networkNodes,
		LightVM:          lightVM,
		HeavyVM:          heavyVM,
		LightQueue:       make(chan Instruction, 1000),
		HeavyQueue:       make(chan Instruction, 1000),
		MaxLightLoad:     maxLightLoad,
		MaxHeavyLoad:     maxHeavyLoad,
		currentLightLoad: 0,
		currentHeavyLoad: 0,
	}

	// Start processing queues only if LocalNode, LightVM, or HeavyVM is available
	if router.LocalNode != nil && (router.LocalNode.HasLightVM || router.LocalNode.HasHeavyVM) {
		router.StartProcessing()
	}

	return router
}



// RouteInstruction routes an instruction to the appropriate VM based on complexity, VM load, and node capabilities.
func (router *VMRouter) RouteInstruction(instr Instruction) error {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	if instr.RequiresHeavy || instr.Complexity > LightVMComplexityThreshold {
		// Instruction requires HeavyVM
		if router.LocalNode.HasHeavyVM && router.currentHeavyLoad < router.MaxHeavyLoad {
			// Route to local HeavyVM
			router.HeavyQueue <- instr
			router.currentHeavyLoad++
			return nil
		} else {
			// Route to a network node running HeavyVM
			err := router.routeToNetworkVM(instr, true)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		// Instruction can be handled by LightVM
		if router.LocalNode.HasLightVM && router.currentLightLoad < router.MaxLightLoad {
			// Route to local LightVM
			router.LightQueue <- instr
			router.currentLightLoad++
			return nil
		} else {
			// Route to a network node running LightVM
			err := router.routeToNetworkVM(instr, false)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func (router *VMRouter) routeToNetworkVM(instr Instruction, requiresHeavy bool) error {
	var availableNodes []*Node

	// Find nodes running the required VM type
	for _, networkNode := range router.NetworkNodes {
		node, ok := networkNode.(*Node) // Cast to *Node
		if !ok {
			continue // Skip if not a *Node type
		}
		if requiresHeavy && node.HasHeavyVM {
			availableNodes = append(availableNodes, node)
		} else if !requiresHeavy && node.HasLightVM {
			availableNodes = append(availableNodes, node)
		}
	}

	if len(availableNodes) == 0 {
		return errors.New("no available nodes running the required VM type")
	}

	// Attempt to send the instruction to each available node
	for _, node := range availableNodes {
		success := sendInstructionToNode(node, instr)
		if success {
			if router.LoggingEnabled() {
				log.Printf("Instruction %s routed to node %s\n", instr.ID, node.Address)
			}
			return nil
		}
	}

	return errors.New("failed to route instruction to any available node")
}


// sendInstructionToNode sends an instruction to the specified node using HTTP with retries and logging.
func sendInstructionToNode(node NetworkNode, instr Instruction) bool {
	const (
		maxRetries     = 3
		retryDelay     = 500 * time.Millisecond
		timeoutSeconds = 3
	)

	// Prepare JSON payload for the instruction
	requestBody, err := json.Marshal(instr)
	if err != nil {
		log.Printf("Failed to marshal instruction: %v", err)
		return false
	}

	// Set up an HTTP client with a timeout
	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}

	// Construct the request URL
	url := fmt.Sprintf("http://%s/execute_instruction", node.GetAddress())

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to create request: %v", attempt, maxRetries, err)
			return false
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the instruction via HTTP POST
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to send instruction to node %s: %v", attempt, maxRetries, node.GetAddress(), err)
			time.Sleep(retryDelay)
			continue
		}

		// Check response status code
		if resp.StatusCode == http.StatusOK {
			log.Printf("Instruction sent successfully to node %s", node.GetAddress())
			return true
		}

		log.Printf("Attempt %d/%d: Node %s returned status code %d; retrying...", attempt, maxRetries, node.GetAddress(), resp.StatusCode)
		resp.Body.Close() // Close response body to avoid resource leaks
		time.Sleep(retryDelay)
	}

	// All attempts failed
	log.Printf("Failed to send instruction to node %s after %d attempts", node.GetAddress(), maxRetries)
	return false
}

func (router *VMRouter) StartProcessing() {
	executionCounter := 0 // Add a counter to track execution loops

	if router.LocalNode != nil && router.LocalNode.HasLightVM && router.LightVM != nil {
		go func() {
			for instr := range router.LightQueue {
				executionCounter++
				log.Printf("Processing LightVM Instruction: %s, Execution Count: %d\n", instr.ID, executionCounter)

				err := router.LightVM.ExecuteInstruction(instr)
				if err != nil {
					log.Printf("Error executing instruction %s in LightVM: %v\n", instr.ID, err)
				}

				router.mutex.Lock()
				router.currentLightLoad--
				router.mutex.Unlock()

				if executionCounter > 100 { // Limit to avoid infinite loop during testing
					log.Println("Breaking LightVM processing loop for testing.")
					break
				}
			}
		}()
	}

	// Similar adjustments for HeavyVM processing loop
	if router.LocalNode != nil && router.LocalNode.HasHeavyVM && router.HeavyVM != nil {
		go func() {
			for instr := range router.HeavyQueue {
				executionCounter++
				log.Printf("Processing HeavyVM Instruction: %s, Execution Count: %d\n", instr.ID, executionCounter)

				err := router.HeavyVM.ExecuteInstruction(instr)
				if err != nil {
					log.Printf("Error executing instruction %s in HeavyVM: %v\n", instr.ID, err)
				}

				router.mutex.Lock()
				router.currentHeavyLoad--
				router.mutex.Unlock()

				if executionCounter > 100 { // Limit for testing
					log.Println("Breaking HeavyVM processing loop for testing.")
					break
				}
			}
		}()
	}
}



// LoggingEnabled returns true if logging is enabled in either VM.
func (router *VMRouter) LoggingEnabled() bool {
	return (router.LightVM != nil && router.LightVM.LoggingEnabled) || (router.HeavyVM != nil && router.HeavyVM.LoggingEnabled)
}

// LightVMComplexityThreshold defines the maximum complexity that the LightVM can handle.
const LightVMComplexityThreshold = 5

// GetLightVMLoad returns the current load on LightVMs.
func (router *VMRouter) GetLightVMLoad() int {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	return router.lightVMLoad
}

// GetHeavyVMLoad returns the current load on HeavyVMs.
func (router *VMRouter) GetHeavyVMLoad() int {
	router.mutex.RLock()
	defer router.mutex.RUnlock()
	return router.heavyVMLoad
}

// IncreaseLightVMLoad increments the LightVM load counter.
func (router *VMRouter) IncreaseLightVMLoad() {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.lightVMLoad++
}

// IncreaseHeavyVMLoad increments the HeavyVM load counter.
func (router *VMRouter) IncreaseHeavyVMLoad() {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.heavyVMLoad++
}

// DecreaseLightVMLoad decrements the LightVM load counter.
func (router *VMRouter) DecreaseLightVMLoad() {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	if router.lightVMLoad > 0 {
		router.lightVMLoad--
	}
}

// DecreaseHeavyVMLoad decrements the HeavyVM load counter.
func (router *VMRouter) DecreaseHeavyVMLoad() {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	if router.heavyVMLoad > 0 {
		router.heavyVMLoad--
	}
}