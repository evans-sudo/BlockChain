package core

import (
	"blockchainsystem/types"
	"crypto/sha256"
	
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}


type BlockHasher struct {}


func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}