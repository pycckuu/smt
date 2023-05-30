package smt

import (
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/stretchr/testify/assert"
)

var zeroLeaf, _ = poseidon.Hash([]*big.Int{big.NewInt(0)})

func TestNewSparseMerkleTree(t *testing.T) {
	smt := NewSparseMerkleTree(2, zeroLeaf)
	assert.NotNil(t, smt)
	assert.NotNil(t, smt.Root)
	assert.Equal(t, 2, smt.Depth)
	assert.Empty(t, smt.Leaves)

	tests := []struct {
		index        int
		value        *big.Int
		expectedRoot string
	}{
		{
			index:        0,
			value:        big.NewInt(0),
			expectedRoot: "18366138217714291923534712849449091358386817997964088830897385671725623871073",
		},
		{
			index:        1,
			value:        big.NewInt(1),
			expectedRoot: "8029606767784791880250783890079025673413177318829731918696248003381461757603",
		},
		{
			index:        2,
			value:        big.NewInt(2),
			expectedRoot: "16218640559429857690153995608944582618510520484403789436977896806380962629939",
		},
		{
			index:        3,
			value:        big.NewInt(3),
			expectedRoot: "3720616653028013822312861221679392249031832781774563366107458835261883914924",
		},
	}

	initRoot := new(big.Int)
	initRoot.SetString("2186774891605521484511138647132707263205739024356090574223746683689524510919", 10)
	if smt.Root.Data.Cmp(initRoot) != 0 {
		t.Error("Expected root node data to be", initRoot, "got", smt.Root.Data)
	}

	for _, test := range tests {
		smt.Insert(test.index, test.value)
		expectedRoot := new(big.Int)
		expectedRoot.SetString(test.expectedRoot, 10)
		if smt.Root.Data.Cmp(expectedRoot) != 0 {
			t.Error("Expected root node data to be", expectedRoot, "got", smt.Root.Data)
		}
	}
}

func TestInsert(t *testing.T) {
	smt := NewSparseMerkleTree(3, zeroLeaf)

	index := 0
	value := big.NewInt(5)

	smt.Insert(index, value)

	assert.Equal(t, value, smt.Leaves[getPaddedBinaryString(index, smt.Depth)])
}

func TestGetPaddedBinaryString(t *testing.T) {
	assert.Equal(t, "000", getPaddedBinaryString(0, 3))
	assert.Equal(t, "001", getPaddedBinaryString(1, 3))
	assert.Equal(t, "011", getPaddedBinaryString(3, 3))
	assert.Equal(t, "111", getPaddedBinaryString(7, 3))
}

func TestNewDeterministicSparseMerkleTree(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(3, zeroLeaf)
	assert.NotNil(t, smt)
	assert.NotNil(t, smt.Root)
	assert.Equal(t, 3, smt.Depth)
	assert.NotEmpty(t, smt.Leaves)
	assert.Len(t, smt.Leaves, 8)
}

// This test will depend on the poseidon.Hash function behavior.
func TestMerkleNodeHashes(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(3, zeroLeaf)

	// Test the root hash
	expectedRootHash := smt.Root.Data
	actualRootHash := hashChildren(smt.Root.Left, smt.Root.Right, smt.Depth, zeroLeaf)

	assert.Equal(t, expectedRootHash, actualRootHash)
}

func TestGenerateMerklePath(t *testing.T) {
	smt := NewDeterministicSparseMerkleTree(4, zeroLeaf)

	testCases := []struct {
		index       int
		shouldError bool
		description string
	}{
		{0, false, "Should not return an error for index 0"},
		{1, false, "Should not return an error for index 1"},
		{2, false, "Should not return an error for index 2"},
		{3, false, "Should not return an error for index 3"},
		{4, false, "Should not return an error for index 4"},
		{5, false, "Should not return an error for index 5"},
		{20, true, "Should return an error for non-existing index"},
	}

	for _, tc := range testCases {
		_, err := smt.GenerateMerklePath(tc.index)
		if tc.shouldError {
			assert.Error(t, err, tc.description)
		} else {
			assert.NoError(t, err, tc.description)
		}
	}
}

func TestSparseMerkleTree(t *testing.T) {
	depth := 4
	smt := NewDeterministicSparseMerkleTree(depth, zeroLeaf)

	for i := 0; i < (1 << depth); i++ {
		key := getPaddedBinaryString(i, depth)
		value := smt.Leaves[key]
		path, _ := smt.GenerateMerklePath(i)
		valid := VerifyMerklePath(value, path, smt.Root.Data)
		assert.True(t, valid, "The Merkle path should be valid for all leaves")
	}
}
