package syn3200

import (
	"time"
	"errors"
	"sync"

)

// BillStorage represents the storage structure for SYN3200 bill tokens.
type BillStorage struct {
	storage          map[string]*Bill
	ledgerService    *ledger.Ledger
	encryptionService *encryption.Encryptor
	consensusService *consensus.SynnergyConsensus
	mutex            sync.Mutex
}

// NewBillStorage creates a new instance of BillStorage.
func NewBillStorage(ledger *ledger.Ledger, encryptor *encryption.Encryptor, consensus *consensus.SynnergyConsensus) *BillStorage {
	return &BillStorage{
		storage:          make(map[string]*Bill),
		ledgerService:    ledger,
		encryptionService: encryptor,
		consensusService:  consensus,
	}
}

// StoreBill stores a new bill in the storage after encrypting and validating it.
func (bs *BillStorage) StoreBill(bill *Bill) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Encrypt the bill before storage.
	encryptedBill, err := bs.encryptionService.EncryptData(bill)
	if err != nil {
		return err
	}

	// Store the encrypted bill in the storage map.
	bs.storage[bill.ID] = encryptedBill.(*Bill)

	// Log the storage event in the ledger.
	if err := bs.ledgerService.LogEvent("BillStored", time.Now(), bill.ID); err != nil {
		return err
	}

	// Validate the bill storage with consensus.
	if err := bs.consensusService.ValidateSubBlock(bill.ID); err != nil {
		return err
	}

	return nil
}

// RetrieveBill retrieves and decrypts a bill from the storage.
func (bs *BillStorage) RetrieveBill(billID string) (*Bill, error) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Retrieve the bill from storage.
	bill, exists := bs.storage[billID]
	if !exists {
		return nil, errors.New("bill not found")
	}

	// Decrypt the bill before returning.
	decryptedBill, err := bs.encryptionService.DecryptData(bill)
	if err != nil {
		return nil, err
	}

	return decryptedBill.(*Bill), nil
}

// UpdateBill updates an existing bill in the storage.
func (bs *BillStorage) UpdateBill(bill *Bill) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Encrypt the updated bill.
	encryptedBill, err := bs.encryptionService.EncryptData(bill)
	if err != nil {
		return err
	}

	// Update the bill in the storage map.
	bs.storage[bill.ID] = encryptedBill.(*Bill)

	// Log the update in the ledger.
	if err := bs.ledgerService.LogEvent("BillUpdated", time.Now(), bill.ID); err != nil {
		return err
	}

	// Validate the bill update with consensus.
	if err := bs.consensusService.ValidateSubBlock(bill.ID); err != nil {
		return err
	}

	return nil
}

// DeleteBill removes a bill from the storage.
func (bs *BillStorage) DeleteBill(billID string) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	// Check if the bill exists in storage.
	_, exists := bs.storage[billID]
	if !exists {
		return errors.New("bill not found")
	}

	// Remove the bill from the storage.
	delete(bs.storage, billID)

	// Log the deletion in the ledger.
	if err := bs.ledgerService.LogEvent("BillDeleted", time.Now(), billID); err != nil {
		return err
	}

	// Validate the deletion with consensus.
	if err := bs.consensusService.ValidateSubBlock(billID); err != nil {
		return err
	}

	return nil
}

// ListAllBills returns a list of all stored bills.
func (bs *BillStorage) ListAllBills() ([]*Bill, error) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	var bills []*Bill
	for _, encryptedBill := range bs.storage {
		// Decrypt each bill before adding to the result.
		decryptedBill, err := bs.encryptionService.DecryptData(encryptedBill)
		if err != nil {
			return nil, err
		}
		bills = append(bills, decryptedBill.(*Bill))
	}

	return bills, nil
}

// EncryptAndStoreBillBatch stores a batch of bills securely.
func (bs *BillStorage) EncryptAndStoreBillBatch(bills []*Bill) error {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	for _, bill := range bills {
		// Encrypt each bill.
		encryptedBill, err := bs.encryptionService.EncryptData(bill)
		if err != nil {
			return err
		}

		// Store the encrypted bill.
		bs.storage[bill.ID] = encryptedBill.(*Bill)

		// Log each storage event in the ledger.
		if err := bs.ledgerService.LogEvent("BillStoredInBatch", time.Now(), bill.ID); err != nil {
			return err
		}

		// Validate the bill storage in the consensus for each bill.
		if err := bs.consensusService.ValidateSubBlock(bill.ID); err != nil {
			return err
		}
	}

	return nil
}
