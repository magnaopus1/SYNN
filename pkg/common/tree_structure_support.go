package common


import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"synnergy_network/pkg/ledger"
	"sync"
)

// MerkleTreeNode represents a node in a Merkle Tree.
type MerkleTreeNode struct {
	Hash     string          // Hash value of the node
	Left     *MerkleTreeNode // Left child node
	Right    *MerkleTreeNode // Right child node
	IsLeaf   bool            // Indicates if it is a leaf node
	Data     []byte          // Original data at the leaf node
}

// MerkleTree represents the structure of a Merkle Tree.
type MerkleTree struct {
	Root  *MerkleTreeNode    // Root node of the Merkle Tree
	mutex sync.Mutex         // Mutex for thread safety
}

// NewMerkleTree constructs a Merkle Tree from a list of data blocks.
func NewMerkleTree(dataBlocks [][]byte) (*MerkleTree, error) {
	if len(dataBlocks) == 0 {
		return nil, errors.New("data blocks cannot be empty")
	}

	var nodes []*MerkleTreeNode
	for _, data := range dataBlocks {
		hash := hashData(data)
		nodes = append(nodes, &MerkleTreeNode{
			Hash:   hash,
			IsLeaf: true,
			Data:   data,
		})
	}

	// Build the Merkle Tree by hashing nodes layer by layer.
	for len(nodes) > 1 {
		var level []*MerkleTreeNode
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				parent := createParentNode(nodes[i], nodes[i+1])
				level = append(level, parent)
			} else {
				level = append(level, nodes[i])
			}
		}
		nodes = level
	}

	return &MerkleTree{Root: nodes[0]}, nil
}

// createParentNode creates a parent node by hashing two child nodes.
func createParentNode(left, right *MerkleTreeNode) *MerkleTreeNode {
	combinedHash := hashData([]byte(left.Hash + right.Hash))
	return &MerkleTreeNode{
		Hash:  combinedHash,
		Left:  left,
		Right: right,
	}
}

// hashData generates a SHA-256 hash of the provided data.
func hashData(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// VerifyData verifies if a given data block is part of the Merkle Tree.
func (mt *MerkleTree) VerifyData(data []byte) bool {
	hash := hashData(data)
	return mt.verifyHash(mt.Root, hash)
}

// verifyHash recursively checks if a given hash exists in the Merkle Tree.
func (mt *MerkleTree) verifyHash(node *MerkleTreeNode, hash string) bool {
	if node == nil {
		return false
	}
	if node.IsLeaf && node.Hash == hash {
		return true
	}
	return mt.verifyHash(node.Left, hash) || mt.verifyHash(node.Right, hash)
}

// Binary Search Tree Implementation

// BinaryTreeNode represents a node in a binary search tree.
type BinaryTreeNode struct {
	Key    int             // Key of the node for comparison
	Value  []byte          // Value stored at the node
	Left   *BinaryTreeNode // Left child
	Right  *BinaryTreeNode // Right child
}

// BinarySearchTree represents the structure of a binary search tree.
type BinarySearchTree struct {
	Root  *BinaryTreeNode
	mutex sync.Mutex
}

// Insert adds a new key-value pair into the binary search tree.
func (bst *BinarySearchTree) Insert(key int, value []byte) {
	bst.mutex.Lock()
	defer bst.mutex.Unlock()

	node := &BinaryTreeNode{Key: key, Value: value}
	if bst.Root == nil {
		bst.Root = node
	} else {
		bst.insertNode(bst.Root, node)
	}
}

// insertNode recursively inserts a node into the correct position in the tree.
func (bst *BinarySearchTree) insertNode(root, node *BinaryTreeNode) {
	if node.Key < root.Key {
		if root.Left == nil {
			root.Left = node
		} else {
			bst.insertNode(root.Left, node)
		}
	} else {
		if root.Right == nil {
			root.Right = node
		} else {
			bst.insertNode(root.Right, node)
		}
	}
}

// Search looks for a node by key in the binary search tree and returns the value if found.
func (bst *BinarySearchTree) Search(key int) ([]byte, error) {
	bst.mutex.Lock()
	defer bst.mutex.Unlock()

	node := bst.searchNode(bst.Root, key)
	if node != nil {
		return node.Value, nil
	}
	return nil, errors.New("key not found in binary search tree")
}

// searchNode recursively searches for a node by key in the binary search tree.
func (bst *BinarySearchTree) searchNode(root *BinaryTreeNode, key int) *BinaryTreeNode {
	if root == nil || root.Key == key {
		return root
	}
	if key < root.Key {
		return bst.searchNode(root.Left, key)
	}
	return bst.searchNode(root.Right, key)
}

// Ledger Integration

// RecordMerkleTreeRoot stores the Merkle tree root in the ledger for verification purposes.
func (mt *MerkleTree) RecordMerkleTreeRoot(ledgerInstance *ledger.Ledger) error {
	return ledgerInstance.BlockchainConsensusCoinLedger.RecordMerkleRoot(mt.Root.Hash)
}

// RecordBinarySearchTreeNode stores a node in the ledger for tracking purposes.
func (bst *BinarySearchTree) RecordBinarySearchTreeNode(ledgerInstance *ledger.Ledger, key int) error {
	value, err := bst.Search(key)
	if err != nil {
		return err
	}
	// Convert the byte slice to a string
	stringValue := string(value)
	return ledgerInstance.RecordBinaryNode(key, stringValue)
}

