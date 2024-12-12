package common


import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// VMMemory represents the VM's dedicated memory for contract execution.
type VMMemory struct {
	Memory          []byte                 // Byte slice to hold memory data
	Size            int                    // Size of the memory in bytes
	EncryptionKey   []byte                 // Key for encrypting sensitive memory regions
	EncryptedRanges map[int]int            // Map of address ranges that are encrypted
}

// NewVMMemory initializes VM memory with a specified size and optional encryption key.
func NewVMMemory(size int, encryptionKey []byte) *VMMemory {
	return &VMMemory{
		Memory:          make([]byte, size),
		Size:            size,
		EncryptionKey:   encryptionKey,
		EncryptedRanges: make(map[int]int),
	}
}

// Read reads data from a specified address in memory.
func (vm *VMMemory) Read(address, length int) ([]byte, error) {
	if address+length > vm.Size {
		return nil, fmt.Errorf("memory read out of bounds")
	}
	data := vm.Memory[address : address+length]
	
	// Decrypt if the address range is encrypted
	if _, encrypted := vm.EncryptedRanges[address]; encrypted {
		decryptedData, err := DecryptWithKey(data, vm.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt memory data: %v", err)
		}
		return decryptedData, nil
	}
	return data, nil
}

// Write writes data to a specified address in memory.
func (vm *VMMemory) Write(address int, data []byte) error {
	if address+len(data) > vm.Size {
		return fmt.Errorf("memory write out of bounds")
	}
	copy(vm.Memory[address:], data)
	return nil
}

// WriteEncrypted writes encrypted data to a specified address and marks the range as encrypted.
func (vm *VMMemory) WriteEncrypted(address int, data []byte) error {
	if address+len(data) > vm.Size {
		return fmt.Errorf("memory write out of bounds")
	}

	// Encrypt data before storing in memory
	encryptedData, err := EncryptWithKey(data, vm.EncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt memory data: %v", err)
	}

	copy(vm.Memory[address:], encryptedData)
	vm.EncryptedRanges[address] = address + len(data) // Mark this range as encrypted
	return nil
}

// WriteInt writes an integer value to a specific address in memory.
func (vm *VMMemory) WriteInt(address int, value int64) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(value))
	return vm.Write(address, data)
}

// ReadInt reads an integer value from a specific address in memory.
func (vm *VMMemory) ReadInt(address int) (int64, error) {
	data, err := vm.Read(address, 8)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

// Clear clears the entire memory or specified range if provided.
func (vm *VMMemory) Clear(start, length int) error {
	if start+length > vm.Size {
		return fmt.Errorf("memory clear out of bounds")
	}
	for i := start; i < start+length; i++ {
		vm.Memory[i] = 0
	}
	return nil
}

// MemorySegment defines a segmented region within the memory.
type MemorySegment struct {
	StartAddress int
	Length       int
}

// GetSegment returns a memory segment, ensuring it falls within memory bounds.
func (vm *VMMemory) GetSegment(start, length int) (*MemorySegment, error) {
	if start < 0 || start+length > vm.Size {
		return nil, errors.New("memory segment out of bounds")
	}
	return &MemorySegment{StartAddress: start, Length: length}, nil
}

// Compare compares two sections of memory for equality.
func (vm *VMMemory) Compare(address1, address2, length int) (bool, error) {
	data1, err := vm.Read(address1, length)
	if err != nil {
		return false, err
	}
	data2, err := vm.Read(address2, length)
	if err != nil {
		return false, err
	}
	return bytes.Equal(data1, data2), nil
}

// SecureWipe performs a secure memory wipe for a sensitive data range, ensuring the data is irrecoverable.
func (vm *VMMemory) SecureWipe(start, length int) error {
	if start+length > vm.Size {
		return fmt.Errorf("secure wipe out of bounds")
	}
	for i := start; i < start+length; i++ {
		vm.Memory[i] = 0xFF // Overwrite with 0xFF before zeroing
	}
	return vm.Clear(start, length)
}
