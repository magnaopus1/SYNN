package data_management

import (
    "synnergy_network/pkg/common"
    "synnergy_network/pkg/ledger"
)

// MULTILEVEL_BINARY_TREE_SETUP initializes a multilevel binary tree with the given data
func MULTILEVEL_BINARY_TREE_SETUP(levels int, data []string) (*common.BinaryTree, error) {
	if levels <= 0 {
		return nil, errors.New("levels must be positive")
	}
	root := &common.BinaryNode{Data: data[0]}
	tree := &common.BinaryTree{Root: root}
	currentLevel := []*common.BinaryNode{root}
	for l := 1; l < levels; l++ {
		var nextLevel []*common.BinaryNode
		for _, parent := range currentLevel {
			left, right := &common.BinaryNode{}, &common.BinaryNode{}
			parent.Left = left
			parent.Right = right
			nextLevel = append(nextLevel, left, right)
		}
		currentLevel = nextLevel
	}
	return tree, nil
}

// QUERY_MULTILEVEL_BINARY_TREE retrieves a node from a multilevel binary tree based on path
func QUERY_MULTILEVEL_BINARY_TREE(tree *common.BinaryTree, path []bool) (*common.BinaryNode, error) {
	node := tree.Root
	for _, isRight := range path {
		if isRight {
			node = node.Right
		} else {
			node = node.Left
		}
		if node == nil {
			return nil, errors.New("node not found in specified path")
		}
	}
	return node, nil
}

// LINK_MERKLE_NODES_ACROSS_LEVELS links Merkle nodes across tree levels
func LINK_MERKLE_NODES_ACROSS_LEVELS(parent *common.MerkleNode, child *common.MerkleNode) {
	parent.Left = child
	child.Parent = parent
}

// LINK_BINARY_NODES_ACROSS_LEVELS links binary nodes across tree levels
func LINK_BINARY_NODES_ACROSS_LEVELS(parent *common.BinaryNode, child *common.BinaryNode, isRight bool) {
	if isRight {
		parent.Right = child
	} else {
		parent.Left = child
	}
	child.Parent = parent
}

// DEEP_COMPARE_MERKLE_SUBTREES compares two Merkle subtrees for equality
func DEEP_COMPARE_MERKLE_SUBTREES(node1, node2 *common.MerkleNode) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil || node1.Hash != node2.Hash {
		return false
	}
	return DEEP_COMPARE_MERKLE_SUBTREES(node1.Left, node2.Left) && DEEP_COMPARE_MERKLE_SUBTREES(node1.Right, node2.Right)
}

// DEEP_COMPARE_BINARY_SUBTREES compares two binary subtrees for equality
func DEEP_COMPARE_BINARY_SUBTREES(node1, node2 *common.BinaryNode) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil || node1.Data != node2.Data {
		return false
	}
	return DEEP_COMPARE_BINARY_SUBTREES(node1.Left, node2.Left) && DEEP_COMPARE_BINARY_SUBTREES(node1.Right, node2.Right)
}

// MERKLE_HASH_AGGREGATION aggregates hashes of a Merkle subtree starting from a specific node
func MERKLE_HASH_AGGREGATION(node *common.MerkleNode) string {
	if node == nil {
		return ""
	}
	if node.Left == nil && node.Right == nil {
		return node.Hash
	}
	leftHash := MERKLE_HASH_AGGREGATION(node.Left)
	rightHash := MERKLE_HASH_AGGREGATION(node.Right)
	return generateHash(leftHash + rightHash)
}

// NODE_ANCESTOR_LOOKUP finds all ancestors of a node in a binary or Merkle tree
func NODE_ANCESTOR_LOOKUP(node *common.TreeNode) []*common.TreeNode {
	var ancestors []*common.TreeNode
	for node.Parent != nil {
		ancestors = append(ancestors, node.Parent)
		node = node.Parent
	}
	return ancestors
}

// NODE_DESCENDANT_LOOKUP finds all descendants of a node in a binary or Merkle tree
func NODE_DESCENDANT_LOOKUP(node *common.TreeNode) []*common.TreeNode {
	var descendants []*common.TreeNode
	collectDescendants(node, &descendants)
	return descendants
}

// SYNC_MERKLE_TREE_STATE synchronizes the Merkle tree’s state with the ledger
func SYNC_MERKLE_TREE_STATE(tree *common.MerkleTree) error {
	treeHash := MERKLE_HASH_AGGREGATION(tree.Root)
	return common.UpdateLedgerMerkleState(treeHash, time.Now())
}

// SYNC_BINARY_TREE_STATE synchronizes the binary tree’s state with the ledger
func SYNC_BINARY_TREE_STATE(tree *common.BinaryTree) error {
	serializedTree, err := common.Serialize(tree)
	if err != nil {
		return fmt.Errorf("failed to serialize binary tree: %v", err)
	}
	return common.UpdateLedgerBinaryState(serializedTree, time.Now())
}

// ROLLBACK_MERKLE_TREE_CHANGE reverts the Merkle tree to a previous state stored in the ledger
func ROLLBACK_MERKLE_TREE_CHANGE(tree *common.MerkleTree, timestamp time.Time) error {
	savedState, err := common.FetchMerkleStateFromLedger(timestamp)
	if err != nil {
		return fmt.Errorf("failed to fetch Merkle tree state: %v", err)
	}
	tree.Root = savedState.Root
	return nil
}

// ROLLBACK_BINARY_TREE_CHANGE reverts the binary tree to a previous state stored in the ledger
func ROLLBACK_BINARY_TREE_CHANGE(tree *common.BinaryTree, timestamp time.Time) error {
	savedState, err := common.FetchBinaryStateFromLedger(timestamp)
	if err != nil {
		return fmt.Errorf("failed to fetch binary tree state: %v", err)
	}
	tree.Root = savedState.Root
	return nil
}

// LOCK_MERKLE_NODE locks a Merkle node for exclusive access
func LOCK_MERKLE_NODE(node *common.MerkleNode) {
	node.Mutex.Lock()
}

// UNLOCK_MERKLE_NODE unlocks a previously locked Merkle node
func UNLOCK_MERKLE_NODE(node *common.MerkleNode) {
	node.Mutex.Unlock()
}

// LOCK_BINARY_NODE locks a binary tree node for exclusive access
func LOCK_BINARY_NODE(node *common.BinaryNode) {
	node.Mutex.Lock()
}

// UNLOCK_BINARY_NODE unlocks a previously locked binary tree node
func UNLOCK_BINARY_NODE(node *common.BinaryNode) {
	node.Mutex.Unlock()
}

// Helper function: generateHash generates a SHA-256 hash for a given string
func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Helper function: collectDescendants recursively collects all descendants of a node
func collectDescendants(node *common.TreeNode, descendants *[]*common.TreeNode) {
	if node == nil {
		return
	}
	*descendants = append(*descendants, node)
	if binaryNode, ok := node.(*common.BinaryNode); ok {
		collectDescendants(binaryNode.Left, descendants)
		collectDescendants(binaryNode.Right, descendants)
	} else if merkleNode, ok := node.(*common.MerkleNode); ok {
		collectDescendants(merkleNode.Left, descendants)
		collectDescendants(merkleNode.Right, descendants)
	}
}
