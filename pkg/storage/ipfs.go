package storage

import (
    "fmt"
    "os/exec"
    "synnergy_network/pkg/ledger"
    "synnergy_network/pkg/common"
)

// NewIPFSManager initializes a new IPFSManager
func NewIPFSManager(ledgerInstance *ledger.Ledger) *IPFSManager {
    return &IPFSManager{
        LedgerInstance: ledgerInstance,
    }
}

// AddFileToIPFS adds an encrypted file to IPFS and returns the CID (Content Identifier)
func (ipfs *IPFSManager) AddFileToIPFS(filePath string) (string, error) {
    ipfs.mutex.Lock()
    defer ipfs.mutex.Unlock()

    // Encrypt the file before uploading to IPFS
    encryptedFilePath, err := encryption.EncryptFile(filePath, common.EncryptionKey)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt file: %v", err)
    }

    // Run the IPFS add command
    cmd := exec.Command("ipfs", "add", encryptedFilePath)
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to add file to IPFS: %v", err)
    }

    cid := parseCIDFromOutput(string(output))
    if cid == "" {
        return "", fmt.Errorf("failed to parse CID from IPFS output")
    }

    // Log the operation in the ledger
    err = ipfs.LedgerInstance.RecordFileOperation(cid, "ipfs_add", filePath, "")
    if err != nil {
        return "", fmt.Errorf("failed to log IPFS add operation to ledger: %v", err)
    }

    fmt.Printf("File %s added to IPFS with CID: %s\n", filePath, cid)
    return cid, nil
}

// RetrieveFileFromIPFS retrieves a file from IPFS by its CID and decrypts it
func (ipfs *IPFSManager) RetrieveFileFromIPFS(cid string, outputPath string) error {
    ipfs.mutex.Lock()
    defer ipfs.mutex.Unlock()

    // Run the IPFS cat command to retrieve the file
    cmd := exec.Command("ipfs", "cat", cid)
    output, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("failed to retrieve file from IPFS: %v", err)
    }

    // Save the encrypted file locally
    err = encryption.SaveEncryptedFile(outputPath, output)
    if err != nil {
        return fmt.Errorf("failed to save encrypted file: %v", err)
    }

    // Decrypt the file
    err = encryption.DecryptFile(outputPath, common.EncryptionKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt file: %v", err)
    }

    // Log the operation in the ledger
    err = ipfs.LedgerInstance.RecordFileOperation(cid, "ipfs_retrieve", outputPath, "")
    if err != nil {
        return fmt.Errorf("failed to log IPFS retrieve operation to ledger: %v", err)
    }

    fmt.Printf("File with CID %s retrieved and decrypted successfully.\n", cid)
    return nil
}

// DeleteFileFromIPFS simulates a deletion by marking the operation in the ledger (IPFS is immutable)
func (ipfs *IPFSManager) DeleteFileFromIPFS(cid string) error {
    ipfs.mutex.Lock()
    defer ipfs.mutex.Unlock()

    // Note: IPFS is immutable, so files cannot be truly deleted. We mark the deletion in the ledger.
    err := ipfs.LedgerInstance.RecordFileOperation(cid, "ipfs_delete", "", "")
    if err != nil {
        return fmt.Errorf("failed to log IPFS delete operation to ledger: %v", err)
    }

    fmt.Printf("File with CID %s marked as deleted in the ledger.\n", cid)
    return nil
}

// parseCIDFromOutput parses the CID from the IPFS add command output
func parseCIDFromOutput(output string) string {
    // Example output: added QmT5NvUtoM5nXKS4T3LNRbqGosjxghAwTx7tNCHeLA2eGQ myfile.txt
    var cid string
    _, err := fmt.Sscanf(output, "added %s", &cid)
    if err != nil {
        return ""
    }
    return cid
}


IPFS_ADD_FILE: Adds a file to IPFS and returns a unique content identifier (CID).
IPFS_GET_FILE: Retrieves a file from IPFS using its CID.
IPFS_PIN_FILE: Pins a file on IPFS to ensure persistence on a specific node.
IPFS_UNPIN_FILE: Unpins a file from IPFS, making it available for garbage collection.
IPFS_HASH: Generates a hash for a file or data before storing it on IPFS.
IPFS_METADATA: Retrieves metadata for a file stored in IPFS.
IPFS_REMOVE_FILE: Removes a file from IPFS (only for the node that added it).
IPFS_LIST_FILES: Lists files stored on the node, filtered by certain criteria (e.g., pinned, unpinned).
IPFS_GET_NODE_INFO: Retrieves information about the current IPFS node (e.g., node ID, storage space).