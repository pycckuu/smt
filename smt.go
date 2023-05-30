/*
Package smt implements a Sparse Merkle Tree data structure.

The SparseMerkleTree struct represents a sparse Merkle tree and contains the
root node, depth, and a map of leaves. The MerklePathItem struct represents an
item in the Merkle tree path. MerkleNode represents the individual nodes of the
Merkle Tree.
*/

package smt

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

// SparseMerkleTree represents a sparse Merkle tree.
type SparseMerkleTree struct {
	Root     *MerkleNode         // The root node of the Sparse Merkle Tree.
	Depth    int                 // The depth of the Sparse Merkle Tree.
	Leaves   map[string]*big.Int // The leaves of the Sparse Merkle Tree, where keys are the binary representation of the index.
	ZeroLeaf *big.Int            // Hash of the zero leaf.
}

// MerklePathItem represents an item in the Merkle tree path.
type MerklePathItem struct {
	SiblingHash *big.Int // Hash of the sibling in the Merkle Path.
	IsRight     bool     // Indicates whether this sibling node is a right child.
}

// MerkleNode represents an individual node in the Merkle Tree.
type MerkleNode struct {
	Left  *MerkleNode // Left child of the current node.
	Right *MerkleNode // Right child of the current node.
	Data  *big.Int    // Hash of the current node.
}

// NewSparseMerkleTree creates a new sparse Merkle tree with empty leaves.
func NewSparseMerkleTree(depth int, zeroLeaf *big.Int) *SparseMerkleTree {
	emptyLeaves := make(map[string]*big.Int)
	root := &MerkleNode{Data: getHashEmptyForDepth(depth, zeroLeaf)}
	return &SparseMerkleTree{Root: root, Depth: depth, Leaves: emptyLeaves, ZeroLeaf: zeroLeaf}
}

// Insert inserts a leaf with the given index and value into the tree.
func (smt *SparseMerkleTree) Insert(index int, value *big.Int) {
	key := getPaddedBinaryString(int(index), smt.Depth)
	smt.Leaves[key] = value
	smt.Root = smt.insertIntoNode(smt.Root, key, value, 0, smt.Depth)
}

// insertIntoNode inserts a leaf into the given node at the specified depth.
func (smt *SparseMerkleTree) insertIntoNode(node *MerkleNode, key string, value *big.Int, depth, maxDepth int) *MerkleNode {
	if node == nil {
		node = &MerkleNode{Data: getHashEmptyForDepth(maxDepth-depth, smt.ZeroLeaf)}
	}

	if depth == maxDepth {
		return &MerkleNode{Data: value}
	}

	pathBit := getPathBit(key, depth)
	if pathBit == 0 {
		node.Left = smt.insertIntoNode(node.getLeftChild(depth+1, smt.ZeroLeaf), key, value, depth+1, maxDepth)
	} else {
		node.Right = smt.insertIntoNode(node.getRightChild(depth+1, smt.ZeroLeaf), key, value, depth+1, maxDepth)
	}

	node.Data = hashChildren(node.Left, node.Right, maxDepth-depth, smt.ZeroLeaf)
	return node
}

// getLeftChild returns the left child node of the current node.
func (node *MerkleNode) getLeftChild(depth int, zeroLeaf *big.Int) *MerkleNode {
	if node.Left == nil {
		return &MerkleNode{Data: getHashEmptyForDepth(depth, zeroLeaf), Left: nil, Right: nil}
	}
	return node.Left
}

// getRightChild returns the right child node of the current node.
func (node *MerkleNode) getRightChild(depth int, zeroLeaf *big.Int) *MerkleNode {
	if node.Right == nil {
		return &MerkleNode{Data: getHashEmptyForDepth(depth, zeroLeaf), Left: nil, Right: nil}
	}
	return node.Right
}

// GenerateMerklePath generates a Merkle tree path for the leaf with the given index.
func (smt *SparseMerkleTree) GenerateMerklePath(index int) ([]*MerklePathItem, error) {
	key := getPaddedBinaryString(int(index), smt.Depth)
	if _, exists := smt.Leaves[key]; !exists {
		return nil, fmt.Errorf("no leaf exists at key: %s", key)
	}

	path := make([]*MerklePathItem, smt.Depth)
	current := smt.Root
	for depth := 0; depth < smt.Depth; depth++ {
		pathBit := getPathBit(key, depth)
		if pathBit == 0 {
			path[depth] = &MerklePathItem{
				SiblingHash: current.getRightChild(depth+1, smt.ZeroLeaf).Data,
				IsRight:     true,
			}
			current = current.getLeftChild(depth+1, smt.ZeroLeaf)
		} else {
			path[depth] = &MerklePathItem{
				SiblingHash: current.getLeftChild(depth+1, smt.ZeroLeaf).Data,
				IsRight:     false,
			}
			current = current.getRightChild(depth+1, smt.ZeroLeaf)
		}
	}

	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

// VerifyMerklePath verifies a Merkle tree path against the expected root hash.
func VerifyMerklePath(leafHash *big.Int, path []*MerklePathItem, expectedRoot *big.Int) bool {
	currentHash := leafHash
	for _, item := range path {
		siblingHash := item.SiblingHash

		if item.IsRight {
			currentHash, _ = poseidon.Hash([]*big.Int{currentHash, siblingHash})
		} else {
			currentHash, _ = poseidon.Hash([]*big.Int{siblingHash, currentHash})
		}
	}

	return currentHash.Cmp(expectedRoot) == 0
}
