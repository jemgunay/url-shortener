package hash

import (
	"fmt"
	"time"

	"github.com/speps/go-hashids/v2"
)

type Hasher interface {
	Hash(string) (string, error)
}

type Generator struct {
	EpochFunc func() int64
}

func New() Generator {
	return Generator{
		EpochFunc: func() int64 {
			// use nano timestamp by default to seed the hash generator with more data
			return time.Now().UnixNano()
		},
	}
}

func (g Generator) Hash(val string) (string, error) {
	// define config for hasher
	hashData := hashids.NewData()
	hashData.Salt = val
	hashData.MinLength = 6

	hashID, err := hashids.NewWithData(hashData)
	if err != nil {
		return "", fmt.Errorf("failed to create new hash: %s", err)
	}

	outputHash, err := hashID.EncodeInt64([]int64{g.EpochFunc()})
	if err != nil {
		return "", fmt.Errorf("failed to encode timestamp: %s", err)
	}

	return outputHash, nil
}
