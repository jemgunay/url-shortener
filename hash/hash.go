package hash

import (
	"fmt"
	"time"

	"github.com/speps/go-hashids/v2"
)

// Hasher defines the requirements for a type which can generate hashes.
type Hasher interface {
	Hash(string) (string, error)
}

// Generator generates hashes. Generating hashes is concurrency safe.
type Generator struct {
	// EpochFunc defines how epochs for hash creation are generated.
	EpochFunc func() int64
}

// Ensure Generator satisfies Hasher.
var _ Hasher = Generator{}

// New creates a Generator with the EpochFunc func initialised to return a current nanosecond timestamp.
func New() Generator {
	return Generator{
		EpochFunc: func() int64 {
			// use nano timestamp by default to seed the hash generator with more data
			return time.Now().UnixNano()
		},
	}
}

// Hash generates a unique hash for the given value. The Generator's EpochFunc contributes to the randomness and length
// of the output hashes.
func (g Generator) Hash(val string) (string, error) {
	// define config for hasher
	hashData := hashids.NewData()
	hashData.Salt = val
	hashData.MinLength = 6

	hashID, err := hashids.NewWithData(hashData)
	if err != nil {
		return "", fmt.Errorf("failed to create new hash ID from hash data: %s", err)
	}

	// the length of the provided number is proportional to the length of the output
	outputHash, err := hashID.EncodeInt64([]int64{g.EpochFunc()})
	if err != nil {
		return "", fmt.Errorf("failed to hash timestamp: %s", err)
	}

	return outputHash, nil
}
