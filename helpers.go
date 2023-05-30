package smt

import (
	"math"
	"math/big"
	"strconv"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

// getHashEmptyForDepth calculates the hash value for an empty node at a given depth.
func getHashEmptyForDepth(depth int, zeroLeaf *big.Int) *big.Int {
	h := zeroLeaf
	for i := 0; i < depth; i++ {
		h, _ = poseidon.Hash([]*big.Int{h, h})
	}
	return h
}

// hashChildren computes the hash value of two child nodes.
func hashChildren(left, right *MerkleNode, depth int, zeroLeaf *big.Int) *big.Int {
	leftData := getHashEmptyForDepth(depth-1, zeroLeaf)
	rightData := getHashEmptyForDepth(depth-1, zeroLeaf)

	if left != nil {
		leftData = left.Data
	}

	if right != nil {
		rightData = right.Data
	}

	hash, _ := poseidon.Hash([]*big.Int{leftData, rightData})
	return hash
}

// getPathBit retrieves the bit value of the key at the specified depth.
func getPathBit(key string, depth int) int {
	if len(key) == 0 {
		return 0
	}
	i, _ := strconv.Atoi(key[depth : depth+1])
	return i
}

// getPaddedBinaryString returns a binary string representation of an integer,
// padded with leading zeros to a specified length.
func getPaddedBinaryString(i int, depth int) string {
	binStr := strconv.FormatInt(int64(i), 2)
	for len(binStr) < depth {
		binStr = "0" + binStr
	}
	return binStr
}

// NewDeterministicSparseMerkleTree creates a new deterministic sparse Merkle tree with non-null leaves.
func NewDeterministicSparseMerkleTree(depth int, zeroLeaf *big.Int) *SparseMerkleTree {
	numLeaves := int(math.Pow(2, float64(depth)))
	smt := NewSparseMerkleTree(depth, zeroLeaf)
	for i := 0; i < numLeaves; i++ {
		leaf := big.NewInt(int64(i))
		smt.Insert(i, leaf)
	}

	return smt
}
