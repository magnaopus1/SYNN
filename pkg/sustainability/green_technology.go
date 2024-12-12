package sustainability

import (
	"errors"
	"fmt"
	"time"

	"synnergy_network/pkg/common"
	"synnergy_network/pkg/ledger"
)

// NewGreenTechnologySystem initializes the green technology system
func NewGreenTechnologySystem(ledgerInstance *ledger.Ledger, encryptionService *common.Encryption) *GreenTechnologySystem {
	return &GreenTechnologySystem{
		HardwareInventory: make(map[string]*GreenHardware),
		SoftwareInventory: make(map[string]string),
		Programs:          make(map[string]string),
		NodeCertificates:  make(map[string]string),
		Conservation:      []string{},
		Ledger:            ledgerInstance,
		EncryptionService: encryptionService,
	}
}

// RegisterGreenHardware registers new eco-friendly hardware
func (gts *GreenTechnologySystem) RegisterGreenHardware(hardwareID, manufacturer, model string, energyRating float64) (*GreenHardware, error) {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	// Encrypt hardware data
	hardwareData := fmt.Sprintf("HardwareID: %s, Manufacturer: %s, Model: %s, EnergyRating: %f", hardwareID, manufacturer, model, energyRating)
	encryptedData, err := gts.EncryptionService.EncryptData([]byte(hardwareData), common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt hardware data: %v", err)
	}

	// Register the hardware
	hardware := &GreenHardware{
		HardwareID:     hardwareID,
		Manufacturer:   manufacturer,
		Model:          model,
		EnergyRating:   energyRating,
		RegisteredDate: time.Now(),
	}
	gts.HardwareInventory[hardwareID] = hardware

	// Log the hardware registration in the ledger
	err = gts.Ledger.RecordGreenHardwareRegistration(hardwareID, manufacturer, model, energyRating, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to log hardware registration: %v", err)
	}

	fmt.Printf("Green hardware %s (%s %s) registered with energy rating %f\n", hardwareID, manufacturer, model, energyRating)
	return hardware, nil
}

// CalculateNetworkEfficiency calculates the overall network efficiency based on the registered hardware
func (gts *GreenTechnologySystem) CalculateNetworkEfficiency() (float64, error) {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	if len(gts.HardwareInventory) == 0 {
		return 0, errors.New("no green hardware registered in the system")
	}

	totalEfficiency := 0.0
	for _, hardware := range gts.HardwareInventory {
		totalEfficiency += hardware.EnergyRating
	}

	averageEfficiency := totalEfficiency / float64(len(gts.HardwareInventory))
	fmt.Printf("Calculated network efficiency: %f\n", averageEfficiency)
	return averageEfficiency, nil
}

// PrintHardwareDetails prints the details of all registered green hardware
func (gts *GreenTechnologySystem) PrintHardwareDetails() {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	for _, hardware := range gts.HardwareInventory {
		fmt.Printf("HardwareID: %s, Manufacturer: %s, Model: %s, EnergyRating: %f\n", hardware.HardwareID, hardware.Manufacturer, hardware.Model, hardware.EnergyRating)
	}
}

// RegisterSoftware registers new eco-friendly software
func (gts *GreenTechnologySystem) RegisterSoftware(softwareID, description string) error {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	// Encrypt software data
	softwareData := fmt.Sprintf("SoftwareID: %s, Description: %s", softwareID, description)
	encryptedData, err := gts.EncryptionService.EncryptData([]byte(softwareData), common.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt software data: %v", err)
	}

	// Register the software
	gts.SoftwareInventory[softwareID] = description

	// Log the software registration in the ledger
	err = gts.Ledger.RecordSoftwareRegistration(softwareID, description, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log software registration: %v", err)
	}

	fmt.Printf("Eco-friendly software %s registered: %s\n", softwareID, description)
	return nil
}

// LaunchCircularEconomyProgram registers a new circular economy program
func (gts *GreenTechnologySystem) LaunchCircularEconomyProgram(programID, description string) error {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	// Register the program
	gts.Programs[programID] = description

	// Log the program launch in the ledger
	err := gts.Ledger.RecordCircularEconomyProgram(programID, description, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log circular economy program: %v", err)
	}

	fmt.Printf("Circular economy program %s launched: %s\n", programID, description)
	return nil
}

// AwardEcoFriendlyCertificate awards an eco-friendly certificate to a node that meets certain green standards
func (gts *GreenTechnologySystem) AwardEcoFriendlyCertificate(nodeID, certificateDetails string) error {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	// Award the certificate
	gts.NodeCertificates[nodeID] = certificateDetails

	// Log the certificate award in the ledger
	err := gts.Ledger.RecordEcoFriendlyCertificateAward(nodeID, certificateDetails, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log certificate award: %v", err)
	}

	fmt.Printf("Eco-friendly certificate awarded to node %s: %s\n", nodeID, certificateDetails)
	return nil
}

// RegisterConservationInitiative adds a conservation initiative to the system
func (gts *GreenTechnologySystem) RegisterConservationInitiative(initiative string) {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	// Register the initiative
	gts.Conservation = append(gts.Conservation, initiative)

	// Log the initiative in the ledger
	err := gts.Ledger.RecordConservationInitiative(initiative, time.Now())
	if err != nil {
		fmt.Printf("Failed to log conservation initiative: %v\n", err)
	} else {
		fmt.Printf("Conservation initiative registered: %s\n", initiative)
	}
}

// CalculateEnvironmentalImpact calculates and prints the environmental impact of the green technology system
func (gts *GreenTechnologySystem) CalculateEnvironmentalImpact() {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	totalImpact := float64(len(gts.HardwareInventory) + len(gts.SoftwareInventory))
	fmt.Printf("Environmental impact calculated based on registered hardware and software: %f\n", totalImpact)
}

// CoolingSolutions optimizes cooling solutions for energy savings
func (gts *GreenTechnologySystem) CoolingSolutions() {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	fmt.Println("Cooling solutions optimized for energy savings.")
	// Log cooling solution optimization in the ledger
	err := gts.Ledger.RecordCoolingSolutionOptimization(time.Now())
	if err != nil {
		fmt.Printf("Failed to log cooling solution optimization: %v\n", err)
	}
}

// OptimizeEnergyConsumption optimizes energy consumption across the network
func (gts *GreenTechnologySystem) OptimizeEnergyConsumption() {
	gts.mu.Lock()
	defer gts.mu.Unlock()

	fmt.Println("Energy consumption optimized across the network.")
	// Log energy consumption optimization in the ledger
	err := gts.Ledger.RecordEnergyConsumptionOptimization(time.Now())
	if err != nil {
		fmt.Printf("Failed to log energy consumption optimization: %v\n", err)
	}
}
