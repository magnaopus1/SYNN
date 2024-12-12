package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// CLONE_MERKLE_TREE creates a deep clone of a Merkle tree
func CLONE_MERKLE_TREE(original *common.MerkleTree) *common.MerkleTree {
	clonedTree := &common.MerkleTree{}
	clonedTree.Root = cloneMerkleNode(original.Root)
	return clonedTree
}

// Helper function: cloneMerkleNode recursively clones a Merkle tree node
func cloneMerkleNode(node *common.MerkleNode) *common.MerkleNode {
	if node == nil {
		return nil
	}
	clonedNode := &common.MerkleNode{Hash: node.Hash}
	clonedNode.Left = cloneMerkleNode(node.Left)
	clonedNode.Right = cloneMerkleNode(node.Right)
	return clonedNode
}

// CLONE_BINARY_TREE creates a deep clone of a binary tree
func CLONE_BINARY_TREE(original *common.BinaryTree) *common.BinaryTree {
	clonedTree := &common.BinaryTree{}
	clonedTree.Root = cloneBinaryNode(original.Root)
	return clonedTree
}

// Helper function: cloneBinaryNode recursively clones a binary tree node
func cloneBinaryNode(node *common.BinaryNode) *common.BinaryNode {
	if node == nil {
		return nil
	}
	clonedNode := &common.BinaryNode{Data: node.Data}
	clonedNode.Left = cloneBinaryNode(node.Left)
	clonedNode.Right = cloneBinaryNode(node.Right)
	return clonedNode
}

// MERKLE_TREE_VALIDATION validates the integrity of a Merkle tree by recalculating hashes
func MERKLE_TREE_VALIDATION(tree *common.MerkleTree) bool {
	return validateMerkleNode(tree.Root)
}

// Helper function: validateMerkleNode recursively validates a Merkle node’s hash integrity
func validateMerkleNode(node *common.MerkleNode) bool {
	if node == nil {
		return true
	}
	if node.Left == nil && node.Right == nil {
		return true
	}
	expectedHash := generateHash(node.Left.Hash + node.Right.Hash)
	return node.Hash == expectedHash && validateMerkleNode(node.Left) && validateMerkleNode(node.Right)
}

// BINARY_TREE_VALIDATION checks if a binary tree has a balanced structure and meets constraints
func BINARY_TREE_VALIDATION(tree *common.BinaryTree) bool {
	_, balanced := checkBinaryBalance(tree.Root)
	return balanced
}

// Helper function: checkBinaryBalance checks height balance of a binary tree node
func checkBinaryBalance(node *common.BinaryNode) (int, bool) {
	if node == nil {
		return 0, true
	}
	leftHeight, leftBalanced := checkBinaryBalance(node.Left)
	rightHeight, rightBalanced := checkBinaryBalance(node.Right)
	if !leftBalanced || !rightBalanced || abs(leftHeight-rightHeight) > 1 {
		return 0, false
	}
	return 1 + max(leftHeight, rightHeight), true
}

// BALANCE_MERKLE_TREE ensures a Merkle tree is balanced by adding empty nodes if needed
func BALANCE_MERKLE_TREE(tree *common.MerkleTree) {
	for len(tree.Leaves)%2 != 0 {
		emptyNode := &common.MerkleNode{Hash: ""}
		tree.Leaves = append(tree.Leaves, emptyNode)
	}
	tree.Root = buildMerkleTree(tree.Leaves)
}

// MERKLE_TREE_LEAF_COUNT returns the number of leaves in a Merkle tree
func MERKLE_TREE_LEAF_COUNT(tree *common.MerkleTree) int {
	return countMerkleLeaves(tree.Root)
}

// Helper function: countMerkleLeaves counts the leaves in a Merkle tree recursively
func countMerkleLeaves(node *common.MerkleNode) int {
	if node == nil {
		return 0
	}
	if node.Left == nil && node.Right == nil {
		return 1
	}
	return countMerkleLeaves(node.Left) + countMerkleLeaves(node.Right)
}

// MERKLE_TREE_DEPTH calculates the depth of a Merkle tree
func MERKLE_TREE_DEPTH(tree *common.MerkleTree) int {
	return calculateMerkleDepth(tree.Root)
}

// Helper function: calculateMerkleDepth recursively calculates the depth of a Merkle tree
func calculateMerkleDepth(node *common.MerkleNode) int {
	if node == nil {
		return 0
	}
	return 1 + max(calculateMerkleDepth(node.Left), calculateMerkleDepth(node.Right))
}

// BINARY_TREE_STRUCTURE_ANALYSIS performs structural analysis on a binary tree and returns its height
func BINARY_TREE_STRUCTURE_ANALYSIS(tree *common.BinaryTree) int {
	return binaryTreeHeight(tree.Root)
}

// Helper function: binaryTreeHeight calculates the height of a binary tree node
func binaryTreeHeight(node *common.BinaryNode) int {
	if node == nil {
		return 0
	}
	return 1 + max(binaryTreeHeight(node.Left), binaryTreeHeight(node.Right))
}

// MERKLE_NODE_HASH_UPDATE updates the hash of a specified node and propagates the change to the root
func MERKLE_NODE_HASH_UPDATE(tree *common.MerkleTree, node *common.MerkleNode, newData string) {
	node.Hash = generateHash(newData)
	updateMerkleHashesToRoot(tree.Root, node)
}

// Helper function: updateMerkleHashesToRoot propagates hash updates up to the Merkle tree root
func updateMerkleHashesToRoot(root, target *common.MerkleNode) {
	if root == nil || target == nil {
		return
	}
	if root == target {
		return
	}
	if root.Left != nil {
		updateMerkleHashesToRoot(root.Left, target)
	}
	if root.Right != nil {
		updateMerkleHashesToRoot(root.Right, target)
	}
	root.Hash = generateHash(root.Left.Hash + root.Right.Hash)
}

// BINARY_TREE_NODE_LEVEL returns the level of a specified node in a binary tree
func BINARY_TREE_NODE_LEVEL(tree *common.BinaryTree, target *common.BinaryNode) (int, error) {
	level, found := findBinaryNodeLevel(tree.Root, target, 1)
	if !found {
		return 0, errors.New("node not found in the binary tree")
	}
	return level, nil
}

// Helper function: findBinaryNodeLevel recursively finds the level of a binary tree node
func findBinaryNodeLevel(node, target *common.BinaryNode, level int) (int, bool) {
	if node == nil {
		return 0, false
	}
	if node == target {
		return level, true
	}
	leftLevel, found := findBinaryNodeLevel(node.Left, target, level+1)
	if found {
		return leftLevel, true
	}
	return findBinaryNodeLevel(node.Right, target, level+1)
}

// Helper function: generateHash generates a SHA-256 hash for given data
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Helper functions for numeric operations
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
