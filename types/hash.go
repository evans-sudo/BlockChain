package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}

	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)
	for i :=0; i < 32; i++ {
		b[i] = h[i]
	}

	return b
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

// HashFromBytes converts a 32-byte slice into a Hash
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}
	var h Hash
	copy(h[:], b) // Safely copy the byte slice into the fixed array
	return h
}

// RandomHash generates a random 32-byte hash
func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}

// RandomBytes generates a slice of random bytes of the given length
func RandomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random bytes")
	}
	return b
}
