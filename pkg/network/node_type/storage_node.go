package node_type

import (
	"errors"
	"fmt"
	"sync"
	"synnergy_network/pkg/common"
	"time"
)

// StoredFile represents the structure of a file stored on the storage node.
type StoredFile struct {
	FileID          string    // Unique file ID
	FileName        string    // Original file name
	OwnerWallet     string    // Wallet address of the file owner
	EncryptedData   []byte    // Encrypted file data
	StoredAt        time.Time // Timestamp of when the file was stored
	FileSize        int64     // Size of the file in bytes
	IPFSHash        string    // IPFS hash for decentralized storage
	SwarmHash       string    // Swarm hash for decentralized storage
}

// StorageNode represents a storage node responsible for securely storing files and data.
type StorageNode struct {
	NodeID            string                        // Unique identifier for the storage node
	StorageCapacity   int64                         // Maximum storage capacity of the node in bytes
	UsedStorage       int64                         // Amount of used storage in bytes
	StoredFiles       map[string]*StoredFile        // Map of stored files by file ID
	mutex             sync.Mutex                    // Mutex for thread-safe operations
	ConsensusEngine   *synnergy_consensus.Engine    // Consensus engine for validating storage operations
	EncryptionService *encryption.Encryption        // Encryption service for securing file data
	NetworkManager    *network.NetworkManager       // Network manager for communicating with other nodes
	Ledger            *ledger.Ledger                // Reference to the ledger for logging storage transactions
	SNVM              *common.VMInterface   // The Synnergy Network Virtual Machine
	IPFSService       *ipfs_service.IPFSManager     // IPFS service for decentralized file storage
	SwarmService      *swarm_service.SwarmManager   // Swarm service for decentralized file storage
	CacheService      *cache.CacheManager           // Caching service for optimized file access
}

// NewStorageNode initializes a new storage node with IPFS, Swarm, and cache integration.
func NewStorageNode(nodeID string, storageCapacity int64, consensusEngine *synnergy_consensus.Engine, encryptionService *encryption.Encryption, networkManager *network.NetworkManager, ledgerInstance *ledger.Ledger, snvm *synnergy_vm.VirtualMachine, ipfs *ipfs_service.IPFSManager, swarm *swarm_service.SwarmManager, cache *cache.CacheManager) *StorageNode {
	return &StorageNode{
		NodeID:            nodeID,
		StorageCapacity:   storageCapacity,
		UsedStorage:       0,
		StoredFiles:       make(map[string]*StoredFile),
		ConsensusEngine:   consensusEngine,
		EncryptionService: encryptionService,
		NetworkManager:    networkManager,
		Ledger:            ledgerInstance,
		SNVM:              snvm,
		IPFSService:       ipfs,
		SwarmService:      swarm,
		CacheService:      cache,
	}
}

// StoreFile stores a new file on the storage node, distributes it to IPFS and Swarm, and logs the transaction in the ledger.
func (sn *StorageNode) StoreFile(fileName string, fileData []byte, ownerWallet string) (string, error) {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// Calculate file size and check available storage.
	fileSize := int64(len(fileData))
	if sn.UsedStorage+fileSize > sn.StorageCapacity {
		return "", errors.New("insufficient storage capacity")
	}

	// Encrypt the file data.
	encryptedData, err := sn.EncryptionService.EncryptData(fileData, common.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file data: %v", err)
	}

	// Generate a unique file ID.
	fileID := common.GenerateUniqueID()

	// Store the file on IPFS.
	ipfsHash, err := sn.IPFSService.StoreFile(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to store file on IPFS: %v", err)
	}

	// Store the file on Swarm.
	swarmHash, err := sn.SwarmService.StoreFile(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to store file on Swarm: %v", err)
	}

	// Create the stored file entry.
	storedFile := &StoredFile{
		FileID:        fileID,
		FileName:      fileName,
		OwnerWallet:   ownerWallet,
		EncryptedData: encryptedData,
		StoredAt:      time.Now(),
		FileSize:      fileSize,
		IPFSHash:      ipfsHash,
		SwarmHash:     swarmHash,
	}

	// Store the file in the node's map.
	sn.StoredFiles[fileID] = storedFile
	sn.UsedStorage += fileSize

	// Log the storage operation in the ledger.
	err = sn.Ledger.RecordFileStorage(fileID, ownerWallet, fileName, fileSize)
	if err != nil {
		return "", fmt.Errorf("failed to log file storage in ledger: %v", err)
	}

	// Add the file to the cache for quicker retrieval.
	sn.CacheService.AddToCache(fileID, encryptedData)

	fmt.Printf("File %s (ID: %s) stored successfully on node %s with IPFS Hash: %s, Swarm Hash: %s.\n", fileName, fileID, sn.NodeID, ipfsHash, swarmHash)
	return fileID, nil
}

// RetrieveFile retrieves a stored file from cache (if available), IPFS, or Swarm, decrypts it, and returns its contents.
func (sn *StorageNode) RetrieveFile(fileID string, ownerWallet string) ([]byte, error) {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// First, check the cache for the file.
	if cachedData, exists := sn.CacheService.GetFromCache(fileID); exists {
		return cachedData, nil
	}

	// Check if the file exists.
	storedFile, exists := sn.StoredFiles[fileID]
	if !exists {
		return nil, errors.New("file not found")
	}

	// Verify the ownership of the file.
	if storedFile.OwnerWallet != ownerWallet {
		return nil, errors.New("unauthorized access")
	}

	// Attempt to retrieve the file from IPFS.
	fileData, err := sn.IPFSService.RetrieveFile(storedFile.IPFSHash)
	if err != nil {
		// If IPFS retrieval fails, attempt Swarm.
		fileData, err = sn.SwarmService.RetrieveFile(storedFile.SwarmHash)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file from IPFS or Swarm: %v", err)
		}
	}

	// Decrypt the file data before returning it.
	decryptedData, err := sn.EncryptionService.DecryptData(fileData, common.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file data: %v", err)
	}

	// Add the file to the cache for future access.
	sn.CacheService.AddToCache(fileID, decryptedData)

	fmt.Printf("File %s (ID: %s) retrieved successfully from node %s.\n", storedFile.FileName, fileID, sn.NodeID)
	return decryptedData, nil
}

// DeleteFile deletes a stored file from the storage node, IPFS, Swarm, and cache, and logs the transaction in the ledger.
func (sn *StorageNode) DeleteFile(fileID string, ownerWallet string) error {
	sn.mutex.Lock()
	defer sn.mutex.Unlock()

	// Check if the file exists.
	storedFile, exists := sn.StoredFiles[fileID]
	if !exists {
		return errors.New("file not found")
	}

	// Verify the ownership of the file.
	if storedFile.OwnerWallet != ownerWallet {
		return errors.New("unauthorized access")
	}

	// Delete the file from IPFS.
	err := sn.IPFSService.DeleteFile(storedFile.IPFSHash)
	if err != nil {
		return fmt.Errorf("failed to delete file from IPFS: %v", err)
	}

	// Delete the file from Swarm.
	err = sn.SwarmService.DeleteFile(storedFile.SwarmHash)
	if err != nil {
		return fmt.Errorf("failed to delete file from Swarm: %v", err)
	}

	// Remove the file from the node's storage.
	delete(sn.StoredFiles, fileID)
	sn.UsedStorage -= storedFile.FileSize

	// Remove the file from the cache.
	sn.CacheService.RemoveFromCache(fileID)

	// Log the deletion in the ledger.
	err = sn.Ledger.RecordFileDeletion(fileID, ownerWallet)
	if err != nil {
		return fmt.Errorf("failed to log file deletion in ledger: %v", err)
	}

	fmt.Printf("File %s (ID: %s) deleted successfully by node %s.\n", storedFile.FileName, fileID, sn.NodeID)
	return nil
}
