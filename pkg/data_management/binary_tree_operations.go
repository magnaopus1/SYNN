package data_management

import (
	"errors"
	"fmt"
)

// BinaryTree represents the root node of the binary tree
type BinaryTree struct {
    Root *BinaryTreeNode
}

// TreeNode represents a node in the binary tree
type BinaryTreeNode struct {
    Data   string   
	Value int
    Left   *BinaryTreeNode
    Right  *BinaryTreeNode
    Parent *BinaryTreeNode
}

// CREATE_BINARY_TREE initializes an empty binary tree
func CreateBinaryTree() *BinaryTree {
	return &BinaryTree{}
}

// NewBinaryTree initializes an empty binary tree
func NewBinaryTree() *BinaryTree {
    return &BinaryTree{}
}

// Add inserts a new node into the binary tree
func (bt *BinaryTree) Add(data string) error {
    newNode := &BinaryTreeNode{Data: data}

    if bt.Root == nil {
        bt.Root = newNode
        return nil
    }

    currentNode := bt.Root
    for {
        if data < currentNode.Data {
            // Insert into left subtree
            if currentNode.Left == nil {
                currentNode.Left = newNode
                newNode.Parent = currentNode
                return nil
            }
            currentNode = currentNode.Left
        } else if data > currentNode.Data {
            // Insert into right subtree
            if currentNode.Right == nil {
                currentNode.Right = newNode
                newNode.Parent = currentNode
                return nil
            }
            currentNode = currentNode.Right
        } else {
            return errors.New("duplicate data")
        }
    }
}

// Find searches for a node with the given data
func (bt *BinaryTree) Find(data string) (*BinaryTreeNode, error) {
    currentNode := bt.Root
    for currentNode != nil {
        if data == currentNode.Data {
            return currentNode, nil
        } else if data < currentNode.Data {
            currentNode = currentNode.Left
        } else {
            currentNode = currentNode.Right
        }
    }
    return nil, errors.New("data not found in tree")
}

// Remove deletes a node with the given data from the binary tree
func (bt *BinaryTree) Remove(data string) error {
    nodeToRemove, err := bt.Find(data)
    if err != nil {
        return err
    }

    // Case 1: Node has no children (leaf node)
    if nodeToRemove.Left == nil && nodeToRemove.Right == nil {
        bt.replaceNodeInParent(nodeToRemove, nil)
    } else if nodeToRemove.Left != nil && nodeToRemove.Right != nil {
        // Case 2: Node has two children
        successor := nodeToRemove.Right.findMin()
        nodeToRemove.Data = successor.Data
        bt.replaceNodeInParent(successor, successor.Right)
    } else {
        // Case 3: Node has one child
        if nodeToRemove.Left != nil {
            bt.replaceNodeInParent(nodeToRemove, nodeToRemove.Left)
        } else {
            bt.replaceNodeInParent(nodeToRemove, nodeToRemove.Right)
        }
    }
    return nil
}

// replaceNodeInParent replaces the target node with a new node in the parent
func (bt *BinaryTree) replaceNodeInParent(targetNode, newNode *BinaryTreeNode) {
    if targetNode.Parent == nil {
        bt.Root = newNode
    } else if targetNode == targetNode.Parent.Left {
        targetNode.Parent.Left = newNode
    } else {
        targetNode.Parent.Right = newNode
    }

    if newNode != nil {
        newNode.Parent = targetNode.Parent
    }
}

// findMin finds the node with the minimum value in the subtree
func (node *BinaryTreeNode) findMin() *BinaryTreeNode {
    currentNode := node
    for currentNode.Left != nil {
        currentNode = currentNode.Left
    }
    return currentNode
}

// InOrderTraversal traverses the binary tree in order
func (bt *BinaryTree) InOrderTraversal(node *BinaryTreeNode, visit func(*BinaryTreeNode)) {
    if node != nil {
        bt.InOrderTraversal(node.Left, visit)
        visit(node)
        bt.InOrderTraversal(node.Right, visit)
    }
}

// PreOrderTraversal traverses the binary tree in pre-order
func (bt *BinaryTree) PreOrderTraversal(node *BinaryTreeNode, visit func(*BinaryTreeNode)) {
    if node != nil {
        visit(node)
        bt.PreOrderTraversal(node.Left, visit)
        bt.PreOrderTraversal(node.Right, visit)
    }
}

// PostOrderTraversal traverses the binary tree in post-order
func (bt *BinaryTree) PostOrderTraversal(node *BinaryTreeNode, visit func(*BinaryTreeNode)) {
    if node != nil {
        bt.PostOrderTraversal(node.Left, visit)
        bt.PostOrderTraversal(node.Right, visit)
        visit(node)
    }
}

// DebugPrint displays the binary tree structure for debugging purposes
func (bt *BinaryTree) DebugPrint() {
    if bt.Root == nil {
        fmt.Println("Binary tree is empty.")
        return
    }

    var printNode func(node *BinaryTreeNode, level int)
    printNode = func(node *BinaryTreeNode, level int) {
        if node == nil {
            return
        }
        printNode(node.Right, level+1)
        fmt.Printf("%*s%s\n", level*4, "", node.Data)
        printNode(node.Left, level+1)
    }

    printNode(bt.Root, 0)
}

// InsertBinaryNode inserts a node into the binary tree, keeping it ordered
func InsertBinaryNode(tree *BinaryTree, value int) error {
    newNode := &BinaryTreeNode{Value: value}
    if tree.Root == nil {
        tree.Root = newNode
    } else {
        insertBinaryNodeHelper(tree.Root, newNode)
    }
    return nil
}

func insertBinaryNodeHelper(node, newNode *BinaryTreeNode) {
    if newNode.Value < node.Value {
        if node.Left == nil {
            node.Left = newNode
        } else {
            insertBinaryNodeHelper(node.Left, newNode)
        }
    } else {
        if node.Right == nil {
            node.Right = newNode
        } else {
            insertBinaryNodeHelper(node.Right, newNode)
        }
    }
}


// RemoveBinaryNode removes a node from the binary tree by value
func RemoveBinaryNode(tree *BinaryTree, value int) error {
    if tree.Root == nil {
        return errors.New("tree is empty")
    }
    tree.Root = removeBinaryNodeHelper(tree.Root, value)
    return nil
}

func removeBinaryNodeHelper(node *BinaryTreeNode, value int) *BinaryTreeNode {
    if node == nil {
        return nil
    }
    if value < node.Value {
        node.Left = removeBinaryNodeHelper(node.Left, value)
    } else if value > node.Value {
        node.Right = removeBinaryNodeHelper(node.Right, value)
    } else {
        if node.Left == nil {
            return node.Right
        } else if node.Right == nil {
            return node.Left
        }
        minNode := findMinBinaryNode(node.Right)
        node.Value = minNode.Value
        node.Right = removeBinaryNodeHelper(node.Right, minNode.Value)
    }
    return node
}

// SearchBinaryNode searches for a node with a specific value
func SearchBinaryNode(tree *BinaryTree, value int) (*BinaryTreeNode, error) {
    return searchBinaryNodeHelper(tree.Root, value)
}

func searchBinaryNodeHelper(node *BinaryTreeNode, value int) (*BinaryTreeNode, error) {
    if node == nil {
        return nil, errors.New("node not found")
    }
    if value < node.Value {
        return searchBinaryNodeHelper(node.Left, value)
    } else if value > node.Value {
        return searchBinaryNodeHelper(node.Right, value)
    }
    return node, nil
}


// TraverseBinaryTreeInOrder traverses the tree in-order and returns values
func TraverseBinaryTreeInOrder(tree *BinaryTree) []int {
    values := []int{}
    traverseInOrderHelper(tree.Root, &values)
    return values
}

func traverseInOrderHelper(node *BinaryTreeNode, values *[]int) {
    if node != nil {
        traverseInOrderHelper(node.Left, values)
        *values = append(*values, node.Value)
        traverseInOrderHelper(node.Right, values)
    }
}


// TraverseBinaryTreePreOrder traverses the tree pre-order and returns values
func TraverseBinaryTreePreOrder(tree *BinaryTree) []int {
    values := []int{}
    traversePreOrderHelper(tree.Root, &values)
    return values
}

func traversePreOrderHelper(node *BinaryTreeNode, values *[]int) {
    if node != nil {
        *values = append(*values, node.Value)
        traversePreOrderHelper(node.Left, values)
        traversePreOrderHelper(node.Right, values)
    }
}

// TraverseBinaryTreePostOrder traverses the tree post-order and returns values
func TraverseBinaryTreePostOrder(tree *BinaryTree) []int {
    values := []int{}
    traversePostOrderHelper(tree.Root, &values)
    return values
}

func traversePostOrderHelper(node *BinaryTreeNode, values *[]int) {
    if node != nil {
        traversePostOrderHelper(node.Left, values)
        traversePostOrderHelper(node.Right, values)
        *values = append(*values, node.Value)
    }
}

// FindMinBinaryNode finds the node with the minimum value
func FindMinBinaryNode(tree *BinaryTree) *BinaryTreeNode {
    return findMinBinaryNode(tree.Root)
}

func findMinBinaryNode(node *BinaryTreeNode) *BinaryTreeNode {
    if node.Left == nil {
        return node
    }
    return findMinBinaryNode(node.Left)
}

// FindMaxBinaryNode finds the node with the maximum value
func FindMaxBinaryNode(tree *BinaryTree) *BinaryTreeNode {
    return findMaxBinaryNode(tree.Root)
}

func findMaxBinaryNode(node *BinaryTreeNode) *BinaryTreeNode {
    if node.Right == nil {
        return node
    }
    return findMaxBinaryNode(node.Right)
}


// BalanceBinaryTree balances the binary tree
func BalanceBinaryTree(tree *BinaryTree) {
    nodes := TraverseBinaryTreeInOrder(tree)
    tree.Root = arrayToBalancedTree(nodes, 0, len(nodes)-1)
}

func arrayToBalancedTree(nodes []int, start, end int) *BinaryTreeNode {
    if start > end {
        return nil
    }
    mid := (start + end) / 2
    root := &BinaryTreeNode{Value: nodes[mid]}
    root.Left = arrayToBalancedTree(nodes, start, mid-1)
    root.Right = arrayToBalancedTree(nodes, mid+1, end)
    return root
}

// ArrayToBinaryTree creates a balanced binary tree from a sorted array
func ArrayToBinaryTree(arr []int) *BinaryTree {
    tree := &BinaryTree{}
    tree.Root = arrayToBalancedTree(arr, 0, len(arr)-1)
    return tree
}


// BinaryTreeDepth calculates the depth of the binary tree
func BinaryTreeDepth(tree *BinaryTree) int {
    return calculateDepth(tree.Root)
}

func calculateDepth(node *BinaryTreeNode) int {
    if node == nil {
        return 0
    }
    leftDepth := calculateDepth(node.Left)
    rightDepth := calculateDepth(node.Right)
    if leftDepth > rightDepth {
        return leftDepth + 1
    }
    return rightDepth + 1
}

// BinaryTreeHeight returns the height of the binary tree
func BinaryTreeHeight(tree *BinaryTree) int {
    return BinaryTreeDepth(tree) - 1
}


// VerifyBinaryTree verifies the binary search tree property of the tree
func VerifyBinaryTree(tree *BinaryTree) bool {
    return verifyBinaryTreeHelper(tree.Root, nil, nil)
}

func verifyBinaryTreeHelper(node *BinaryTreeNode, min, max *int) bool {
    if node == nil {
        return true
    }
    if min != nil && node.Value <= *min {
        return false
    }
    if max != nil && node.Value >= *max {
        return false
    }
    return verifyBinaryTreeHelper(node.Left, min, &node.Value) && verifyBinaryTreeHelper(node.Right, &node.Value, max)
}

