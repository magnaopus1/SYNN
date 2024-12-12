package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// MERGE_MERKLE_TREES merges two Merkle trees by creating a new root node
func MERGE_MERKLE_TREES(tree1, tree2 *common.MerkleTree) *common.MerkleTree {
	newRoot := &common.MerkleNode{
		Left:  tree1.Root,
		Right: tree2.Root,
		Hash:  generateHash(tree1.Root.Hash + tree2.Root.Hash),
	}
	return &common.MerkleTree{Root: newRoot}
}

// MERGE_BINARY_TREES merges two binary trees by creating a new root node
func MERGE_BINARY_TREES(tree1, tree2 *common.BinaryTree) *common.BinaryTree {
	newRoot := &common.BinaryNode{
		Left: tree1.Root,
		Right: tree2.Root,
		Data: "MergedRoot",
	}
	return &common.BinaryTree{Root: newRoot}
}

// GENERATE_MERKLE_PATH generates a proof path from a leaf node to the Merkle tree root
func GENERATE_MERKLE_PATH(tree *common.MerkleTree, targetHash string) ([]string, error) {
	path := []string{}
	if !findMerklePath(tree.Root, targetHash, &path) {
		return nil, errors.New("target hash not found in Merkle tree")
	}
	return path, nil
}

// Helper function: findMerklePath recursively finds the path to a target hash
func findMerklePath(node *common.MerkleNode, targetHash string, path *[]string) bool {
	if node == nil {
		return false
	}
	if node.Hash == targetHash {
		return true
	}
	if findMerklePath(node.Left, targetHash, path) || findMerklePath(node.Right, targetHash, path) {
		*path = append(*path, node.Hash)
		return true
	}
	return false
}

// CHECK_TREE_SYMMETRY checks if a binary tree is symmetric
func CHECK_TREE_SYMMETRY(tree *common.BinaryTree) bool {
	return isSymmetric(tree.Root.Left, tree.Root.Right)
}

// Helper function: isSymmetric recursively checks symmetry of two binary nodes
func isSymmetric(left, right *common.BinaryNode) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil || left.Data != right.Data {
		return false
	}
	return isSymmetric(left.Left, right.Right) && isSymmetric(left.Right, right.Left)
}

// CHECK_SUBTREE_PRESENCE checks if a binary tree contains a subtree with the same structure and data
func CHECK_SUBTREE_PRESENCE(mainTree, subtree *common.BinaryTree) bool {
	return isSubtree(mainTree.Root, subtree.Root)
}

// Helper function: isSubtree recursively checks if a node contains a subtree
func isSubtree(main, sub *common.BinaryNode) bool {
	if main == nil {
		return false
	}
	if isIdentical(main, sub) {
		return true
	}
	return isSubtree(main.Left, sub) || isSubtree(main.Right, sub)
}

// Helper function: isIdentical checks if two binary nodes are identical
func isIdentical(node1, node2 *common.BinaryNode) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil || node1.Data != node2.Data {
		return false
	}
	return isIdentical(node1.Left, node2.Left) && isIdentical(node1.Right, node2.Right)
}

// CONVERT_BINARY_TO_MERKLE converts a binary tree to a Merkle tree by hashing node data
func CONVERT_BINARY_TO_MERKLE(binaryTree *common.BinaryTree) *common.MerkleTree {
	return &common.MerkleTree{Root: convertBinaryToMerkleNode(binaryTree.Root)}
}

// Helper function: convertBinaryToMerkleNode recursively converts binary nodes to Merkle nodes
func convertBinaryToMerkleNode(node *common.BinaryNode) *common.MerkleNode {
	if node == nil {
		return nil
	}
	return &common.MerkleNode{
		Hash:  generateHash(node.Data),
		Left:  convertBinaryToMerkleNode(node.Left),
		Right: convertBinaryToMerkleNode(node.Right),
	}
}

// CONVERT_MERKLE_TO_BINARY converts a Merkle tree to a binary tree by storing hash values as data
func CONVERT_MERKLE_TO_BINARY(merkleTree *common.MerkleTree) *common.BinaryTree {
	return &common.BinaryTree{Root: convertMerkleToBinaryNode(merkleTree.Root)}
}

// Helper function: convertMerkleToBinaryNode recursively converts Merkle nodes to binary nodes
func convertMerkleToBinaryNode(node *common.MerkleNode) *common.BinaryNode {
	if node == nil {
		return nil
	}
	return &common.BinaryNode{
		Data:  node.Hash,
		Left:  convertMerkleToBinaryNode(node.Left),
		Right: convertMerkleToBinaryNode(node.Right),
	}
}

// COUNT_BINARY_NODES counts the total number of nodes in a binary tree
func COUNT_BINARY_NODES(tree *common.BinaryTree) int {
	return countBinaryNodes(tree.Root)
}

// Helper function: countBinaryNodes recursively counts nodes in a binary tree
func countBinaryNodes(node *common.BinaryNode) int {
	if node == nil {
		return 0
	}
	return 1 + countBinaryNodes(node.Left) + countBinaryNodes(node.Right)
}

// COUNT_MERKLE_NODES counts the total number of nodes in a Merkle tree
func COUNT_MERKLE_NODES(tree *common.MerkleTree) int {
	return countMerkleNodes(tree.Root)
}

// Helper function: countMerkleNodes recursively counts nodes in a Merkle tree
func countMerkleNodes(node *common.MerkleNode) int {
	if node == nil {
		return 0
	}
	return 1 + countMerkleNodes(node.Left) + countMerkleNodes(node.Right)
}

// FIND_BINARY_NODE_DEPTH finds the depth of a specific binary node
func FIND_BINARY_NODE_DEPTH(tree *common.BinaryTree, target *common.BinaryNode) (int, error) {
	depth, found := binaryNodeDepth(tree.Root, target, 0)
	if !found {
		return -1, errors.New("node not found in binary tree")
	}
	return depth, nil
}

// Helper function: binaryNodeDepth recursively finds the depth of a binary node
func binaryNodeDepth(node, target *common.BinaryNode, depth int) (int, bool) {
	if node == nil {
		return 0, false
	}
	if node == target {
		return depth, true
	}
	if leftDepth, found := binaryNodeDepth(node.Left, target, depth+1); found {
		return leftDepth, true
	}
	return binaryNodeDepth(node.Right, target, depth+1)
}

// UPDATE_BINARY_NODE_VALUE updates the data of a specified binary node
func UPDATE_BINARY_NODE_VALUE(node *common.BinaryNode, newValue string) {
	node.Data = newValue
}

// PRUNE_BINARY_TREE removes all leaf nodes from a binary tree
func PRUNE_BINARY_TREE(tree *common.BinaryTree) {
	tree.Root = pruneBinaryLeaves(tree.Root)
}

// Helper function: pruneBinaryLeaves removes leaf nodes recursively
func pruneBinaryLeaves(node *common.BinaryNode) *common.BinaryNode {
	if node == nil || (node.Left == nil && node.Right == nil) {
		return nil
	}
	node.Left = pruneBinaryLeaves(node.Left)
	node.Right = pruneBinaryLeaves(node.Right)
	return node
}

// BINARY_TREE_LEAF_NODES retrieves all leaf nodes in a binary tree
func BINARY_TREE_LEAF_NODES(tree *common.BinaryTree) []*common.BinaryNode {
	var leaves []*common.BinaryNode
	collectBinaryLeaves(tree.Root, &leaves)
	return leaves
}

// Helper function: collectBinaryLeaves recursively collects leaf nodes
func collectBinaryLeaves(node *common.BinaryNode, leaves *[]*common.BinaryNode) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		*leaves = append(*leaves, node)
	}
	collectBinaryLeaves(node.Left, leaves)
	collectBinaryLeaves(node.Right, leaves)
}

// FIND_MERKLE_LEAF_NODE finds a specific leaf node in a Merkle tree by hash value
func FIND_MERKLE_LEAF_NODE(tree *common.MerkleTree, targetHash string) (*common.MerkleNode, error) {
	return findMerkleLeaf(tree.Root, targetHash)
}

// Helper function: findMerkleLeaf recursively searches for a leaf node with the given hash
func findMerkleLeaf(node *common.MerkleNode, targetHash string) (*common.MerkleNode, error) {
	if node == nil {
		return nil, errors.New("node not found")
	}
	if node.Left == nil && node.Right == nil && node.Hash == targetHash {
		return node, nil
	}
	if leftResult, err := findMerkleLeaf(node.Left, targetHash); err == nil {
		return leftResult, nil
	}
	return findMerkleLeaf(node.Right, targetHash)
}

// MERKLE_TREE_HASH calculates the root hash of a Merkle tree
func MERKLE_TREE_HASH(tree *common.MerkleTree) string {
	return tree.Root.Hash
}

// Helper function: generateHash generates a SHA-256 hash for given data
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
