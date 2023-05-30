# Sparse Merkle Tree Implementation in Go

This repository contains an implementation of a Sparse Merkle Tree (SMT) in Go. An SMT is a kind of Merkle Tree, which is a tree of hashes where each non-leaf node is the hash of its children. Sparse Merkle Trees are 'sparse' because they are very large (often having 2^256 nodes) but are almost entirely empty.

This implementation is highly memory efficient because it only stores non-empty leaves, making it ideal for situations where the tree is mostly empty.

## Key Features

- **Sparse Merkle Tree representation**: Contains the root node, depth, and a map of leaves.
- **Merkle path generation and verification**: Allows the generation of a Merkle path for a given leaf key, and can verify a provided Merkle path against an expected root hash.
- **Leaf insertion**: Supports the insertion of leaves at specific indexes.
- **Deterministic Sparse Merkle Tree creation**: Creates deterministic Sparse Merkle Trees with non-null leaves.

## Code Structure

The repository includes the following important components:

- `smt.go`: Contains the main implementation of the Sparse Merkle Tree, including the definition of the tree structure, leaf insertion, and Merkle path generation and verification.
- `helpers.go`: Contains helper functions for the SMT implementation, such as functions for calculating the hash of an empty node, getting a padded binary string of a given integer, and more.

## Installation and Usage

The code is written in Go, so you'll need to have Go installed on your machine to use it. You can download Go from the [official Go website](https://golang.org/).

To use this code in your project, simply import the `smt` package.

```go
import "github.com/pycckuu/smt"
```

You can then create a new Sparse Merkle Tree like so:

```go
tree := smt.NewSparseMerkleTree(depth, zeroLeaf)
```
Where depth is the desired depth of your tree and zeroLeaf is the hash of the zero leaf.

To insert a new leaf into the tree:

```go
tree.Insert(index, value)
```

Where index is the index at which to insert the new leaf and value is the value of the new leaf.

To generate a Merkle path for a given leaf index:

```go
path, err := tree.GenerateMerklePath(index)
```
And to verify a Merkle path against an expected root:

```go

valid := smt.VerifyMerklePath(leafHash, path, expectedRoot)
```

## Contributions
Contributions to the repository are welcome! Please submit a pull request with your changes.

## License
This project is licensed under the MIT License.
