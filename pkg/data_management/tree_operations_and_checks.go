package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// BINARY_TREE_NODE_SWAP swaps two nodes in a binary tree
func BINARY_TREE_NODE_SWAP(tree *common.BinaryTree, node1, node2 *common.BinaryNode) error {
	if node1 == nil || node2 == nil {
		return errors.New("cannot swap nil nodes")
	}
	node1.Data, node2.Data = node2.Data, node1.Data
	return nil
}

// BALANCED_BINARY_INSERT inserts a new node in a binary tree while keeping it balanced
func BALANCED_BINARY_INSERT(tree *common.BinaryTree, newNode *common.BinaryNode) {
	tree.Root = balancedInsert(tree.Root, newNode)
}

// Helper function: balancedInsert recursively inserts a node in a balanced manner
func balancedInsert(node, newNode *common.BinaryNode) *common.BinaryNode {
	if node == nil {
		return newNode
	}
	if countBinaryNodes(node.Left) <= countBinaryNodes(node.Right) {
		node.Left = balancedInsert(node.Left, newNode)
	} else {
		node.Right = balancedInsert(node.Right, newNode)
	}
	return node
}

// MERKLE_TREE_PROOF_GENERATION generates a proof for a specific leaf node in the Merkle tree
func MERKLE_TREE_PROOF_GENERATION(tree *common.MerkleTree, leafHash string) ([]string, error) {
	path, err := GENERATE_MERKLE_PATH(tree, leafHash)
	if err != nil {
		return nil, fmt.Errorf("proof generation failed: %v", err)
	}
	return path, nil
}

// IMPORT_BINARY_TREE imports a binary tree from a given data source
func IMPORT_BINARY_TREE(data []common.BinaryNodeData) *common.BinaryTree {
	tree := &common.BinaryTree{}
	for _, nodeData := range data {
		newNode := &common.BinaryNode{Data: nodeData.Data}
		BALANCED_BINARY_INSERT(tree, newNode)
	}
	return tree
}

// EXPORT_BINARY_TREE exports a binary tree structure to a list of node data
func EXPORT_BINARY_TREE(tree *common.BinaryTree) []common.BinaryNodeData {
	return exportBinaryNodes(tree.Root)
}

// Helper function: exportBinaryNodes recursively exports binary nodes
func exportBinaryNodes(node *common.BinaryNode) []common.BinaryNodeData {
	if node == nil {
		return []common.BinaryNodeData{}
	}
	left := exportBinaryNodes(node.Left)
	right := exportBinaryNodes(node.Right)
	return append(append(left, common.BinaryNodeData{Data: node.Data}), right...)
}

// INTERSECT_MERKLE_TREES finds the common nodes between two Merkle trees
func INTERSECT_MERKLE_TREES(tree1, tree2 *common.MerkleTree) []string {
	var intersection []string
	findMerkleIntersection(tree1.Root, tree2.Root, &intersection)
	return intersection
}

// Helper function: findMerkleIntersection recursively finds common nodes
func findMerkleIntersection(node1, node2 *common.MerkleNode, intersection *[]string) {
	if node1 == nil || node2 == nil {
		return
	}
	if node1.Hash == node2.Hash {
		*intersection = append(*intersection, node1.Hash)
	}
	findMerkleIntersection(node1.Left, node2.Left, intersection)
	findMerkleIntersection(node1.Right, node2.Right, intersection)
}

// INTERSECT_BINARY_TREES finds common nodes between two binary trees
func INTERSECT_BINARY_TREES(tree1, tree2 *common.BinaryTree) []*common.BinaryNode {
	var intersection []*common.BinaryNode
	findBinaryIntersection(tree1.Root, tree2.Root, &intersection)
	return intersection
}

// Helper function: findBinaryIntersection recursively finds common nodes
func findBinaryIntersection(node1, node2 *common.BinaryNode, intersection *[]*common.BinaryNode) {
	if node1 == nil || node2 == nil {
		return
	}
	if node1.Data == node2.Data {
		*intersection = append(*intersection, node1)
	}
	findBinaryIntersection(node1.Left, node2.Left, intersection)
	findBinaryIntersection(node1.Right, node2.Right, intersection)
}

// MERKLE_HASH_AGGREGATE aggregates hashes in a Merkle tree for a combined root hash
func MERKLE_HASH_AGGREGATE(tree *common.MerkleTree) string {
	return aggregateMerkleHash(tree.Root)
}

// Helper function: aggregateMerkleHash recursively aggregates hashes in a Merkle tree
func aggregateMerkleHash(node *common.MerkleNode) string {
	if node == nil {
		return ""
	}
	if node.Left == nil && node.Right == nil {
		return node.Hash
	}
	leftHash := aggregateMerkleHash(node.Left)
	rightHash := aggregateMerkleHash(node.Right)
	return generateHash(leftHash + rightHash)
}

// BINARY_TREE_SHALLOW_COPY creates a shallow copy of a binary tree
func BINARY_TREE_SHALLOW_COPY(tree *common.BinaryTree) *common.BinaryTree {
	return &common.BinaryTree{Root: tree.Root}
}

// MERKLE_PATH_VERIFICATION verifies a Merkle path from leaf to root
func MERKLE_PATH_VERIFICATION(path []string, rootHash, leafHash string) bool {
	currentHash := leafHash
	for _, hash := range path {
		currentHash = generateHash(currentHash + hash)
	}
	return currentHash == rootHash
}

// COUNT_MERKLE_LEAF_NODES counts the leaf nodes in a Merkle tree
func COUNT_MERKLE_LEAF_NODES(tree *common.MerkleTree) int {
	return countMerkleLeaves(tree.Root)
}

// DELETE_BINARY_TREE deletes all nodes in a binary tree
func DELETE_BINARY_TREE(tree *common.BinaryTree) {
	tree.Root = nil
}

// DELETE_MERKLE_TREE deletes all nodes in a Merkle tree
func DELETE_MERKLE_TREE(tree *common.MerkleTree) {
	tree.Root = nil
}

// MULTILEVEL_MERKLE_TREE_SETUP sets up a multilevel Merkle tree for hierarchical data
func MULTILEVEL_MERKLE_TREE_SETUP(data [][]string) *common.MultiLevelMerkleTree {
	var levels []common.MerkleTree
	for _, levelData := range data {
		leaves := createMerkleLeaves(levelData)
		tree := &common.MerkleTree{Root: buildMerkleTree(leaves)}
		levels = append(levels, *tree)
	}
	return &common.MultiLevelMerkleTree{Levels: levels}
}

// Helper function: createMerkleLeaves creates leaf nodes for Merkle tree setup
func createMerkleLeaves(data []string) []*common.MerkleNode {
	var leaves []*common.MerkleNode
	for _, datum := range data {
		leaves = append(leaves, &common.MerkleNode{Hash: generateHash(datum)})
	}
	return leaves
}

// QUERY_MULTILEVEL_MERKLE_TREE queries data across levels in a multilevel Merkle tree
func QUERY_MULTILEVEL_MERKLE_TREE(tree *common.MultiLevelMerkleTree, dataHash string) (bool, int) {
	for level, merkleTree := range tree.Levels {
		if findMerkleLeafByHash(merkleTree.Root, dataHash) {
			return true, level
		}
	}
	return false, -1
}

// Helper function: findMerkleLeafByHash searches for a hash in a Merkle tree
func findMerkleLeafByHash(node *common.MerkleNode, targetHash string) bool {
	if node == nil {
		return false
	}
	if node.Left == nil && node.Right == nil && node.Hash == targetHash {
		return true
	}
	return findMerkleLeafByHash(node.Left, targetHash) || findMerkleLeafByHash(node.Right, targetHash)
}

// Helper function: generateHash generates a SHA-256 hash for given data
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Helper function for counting binary nodes
func countBinaryNodes(node *common.BinaryNode) int {
	if node == nil {
		return 0
	}
	return 1 + countBinaryNodes(node.Left) + countBinaryNodes(node.Right)
}

// Helper function for counting Merkle leaves
func countMerkleLeaves(node *common.MerkleNode) int {
	if node == nil {
		return 0
	}
	if node.Left == nil && node.Right == nil {
		return 1
	}
	return countMerkleLeaves(node.Left) + countMerkleLeaves(node.Right)
}
