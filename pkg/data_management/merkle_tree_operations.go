package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CREATE_MERKLE_TREE creates a Merkle tree from a list of data hashes
func CREATE_MERKLE_TREE(dataHashes []string) *common.MerkleTree {
	nodes := make([]*common.MerkleNode, len(dataHashes))
	for i, hash := range dataHashes {
		nodes[i] = &common.MerkleNode{Hash: hash}
	}
	tree := &common.MerkleTree{Leaves: nodes}
	tree.Root = buildMerkleTree(nodes)
	return tree
}

// ADD_MERKLE_NODE adds a new node to the Merkle tree and recalculates the root
func ADD_MERKLE_NODE(tree *common.MerkleTree, data string) {
	hash := generateHash(data)
	newNode := &common.MerkleNode{Hash: hash}
	tree.Leaves = append(tree.Leaves, newNode)
	tree.Root = buildMerkleTree(tree.Leaves)
}

// REMOVE_MERKLE_NODE removes a node from the Merkle tree by index and recalculates the root
func REMOVE_MERKLE_NODE(tree *common.MerkleTree, index int) error {
	if index < 0 || index >= len(tree.Leaves) {
		return errors.New("invalid node index")
	}
	tree.Leaves = append(tree.Leaves[:index], tree.Leaves[index+1:]...)
	tree.Root = buildMerkleTree(tree.Leaves)
	return nil
}

// MERKLE_ROOT_CALCULATION calculates and returns the root hash of the Merkle tree
func MERKLE_ROOT_CALCULATION(tree *common.MerkleTree) string {
	tree.Root = buildMerkleTree(tree.Leaves)
	return tree.Root.Hash
}

// VERIFY_MERKLE_PROOF verifies a data entry's inclusion in the Merkle tree using a proof path
func VERIFY_MERKLE_PROOF(rootHash, leafHash string, proofPath []string) bool {
	hash := leafHash
	for _, proof := range proofPath {
		hash = generateHash(hash + proof)
	}
	return hash == rootHash
}

// MERKLE_NODE_LOOKUP finds and returns the index of a node in the Merkle tree with the given hash
func MERKLE_NODE_LOOKUP(tree *common.MerkleTree, hash string) (int, error) {
	for i, node := range tree.Leaves {
		if node.Hash == hash {
			return i, nil
		}
	}
	return -1, errors.New("node not found")
}

// UPDATE_MERKLE_NODE updates a node’s data and recalculates the root
func UPDATE_MERKLE_NODE(tree *common.MerkleTree, index int, newData string) error {
	if index < 0 || index >= len(tree.Leaves) {
		return errors.New("invalid node index")
	}
	tree.Leaves[index].Hash = generateHash(newData)
	tree.Root = buildMerkleTree(tree.Leaves)
	return nil
}

// MERKLE_TREE_TRAVERSE traverses the Merkle tree and returns all node hashes
func MERKLE_TREE_TRAVERSE(tree *common.MerkleTree) []string {
	hashes := []string{}
	collectHashes(tree.Root, &hashes)
	return hashes
}

// MERKLE_PATH_TO_ROOT generates a proof path for a leaf node to the root
func MERKLE_PATH_TO_ROOT(tree *common.MerkleTree, index int) ([]string, error) {
	if index < 0 || index >= len(tree.Leaves) {
		return nil, errors.New("invalid node index")
	}
	return generateProofPath(tree.Leaves, index), nil
}

// COMPARE_MERKLE_TREES compares two Merkle trees and returns true if they have the same structure and hashes
func COMPARE_MERKLE_TREES(tree1, tree2 *common.MerkleTree) bool {
	return compareNodes(tree1.Root, tree2.Root)
}

// EXPORT_MERKLE_TREE serializes the Merkle tree for storage or transmission
func EXPORT_MERKLE_TREE(tree *common.MerkleTree) ([]byte, error) {
	return common.Serialize(tree)
}

// IMPORT_MERKLE_TREE deserializes a Merkle tree from storage or transmission
func IMPORT_MERKLE_TREE(data []byte) (*common.MerkleTree, error) {
	return common.DeserializeMerkleTree(data)
}

// CHECK_MERKLE_INTEGRITY verifies the integrity of a Merkle tree by recalculating hashes and comparing to the root
func CHECK_MERKLE_INTEGRITY(tree *common.MerkleTree) bool {
	calculatedRoot := buildMerkleTree(tree.Leaves)
	return tree.Root.Hash == calculatedRoot.Hash
}

// CALCULATE_SUBTREE_HASH calculates the root hash of a subtree starting from a specified node
func CALCULATE_SUBTREE_HASH(node *common.MerkleNode) string {
	if node == nil {
		return ""
	}
	if node.Left == nil && node.Right == nil {
		return node.Hash
	}
	return generateHash(CALCULATE_SUBTREE_HASH(node.Left) + CALCULATE_SUBTREE_HASH(node.Right))
}

// MERKLE_TREE_BALANCE balances an unbalanced Merkle tree by filling with empty nodes if needed
func MERKLE_TREE_BALANCE(tree *common.MerkleTree) {
	for len(tree.Leaves)%2 != 0 {
		emptyNode := &common.MerkleNode{Hash: ""}
		tree.Leaves = append(tree.Leaves, emptyNode)
	}
	tree.Root = buildMerkleTree(tree.Leaves)
}

// Helper function: buildMerkleTree constructs the tree and returns the root node
func buildMerkleTree(nodes []*common.MerkleNode) *common.MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}
	var newLevel []*common.MerkleNode
	for i := 0; i < len(nodes); i += 2 {
		var left, right *common.MerkleNode
		left = nodes[i]
		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = left
		}
		combinedHash := generateHash(left.Hash + right.Hash)
		newLevel = append(newLevel, &common.MerkleNode{Hash: combinedHash, Left: left, Right: right})
	}
	return buildMerkleTree(newLevel)
}

// Helper function: generateHash generates a SHA-256 hash for a given string
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Helper function: collectHashes collects all hashes in a tree from root to leaves
func collectHashes(node *common.MerkleNode, hashes *[]string) {
	if node == nil {
		return
	}
	*hashes = append(*hashes, node.Hash)
	collectHashes(node.Left, hashes)
	collectHashes(node.Right, hashes)
}

// Helper function: generateProofPath creates a proof path for a node to the root
func generateProofPath(nodes []*common.MerkleNode, index int) []string {
	proof := []string{}
	for len(nodes) > 1 {
		var newLevel []*common.MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			var left, right *common.MerkleNode
			left = nodes[i]
			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				right = left
			}
			if i == index || i+1 == index {
				proof = append(proof, right.Hash)
			}
			combinedHash := generateHash(left.Hash + right.Hash)
			newLevel = append(newLevel, &common.MerkleNode{Hash: combinedHash, Left: left, Right: right})
		}
		nodes = newLevel
		index /= 2
	}
	return proof
}

// Helper function: compareNodes recursively compares two nodes and their children for equality
func compareNodes(node1, node2 *common.MerkleNode) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil || node1.Hash != node2.Hash {
		return false
	}
	return compareNodes(node1.Left, node2.Left) && compareNodes(node1.Right, node2.Right)
}
