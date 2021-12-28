package stub

import "github.com/jemgunay/url-shortener/hash"

// Stub satisfies Hasher and is used to stub out the hash Generator.
type Stub struct {
	Val string
	Err error
}

// Ensure Stub satisfies Hasher.
var _ hash.Hasher = Stub{}

// Hash returns the Stub's Val and Err fields.
func (s Stub) Hash(_ string) (string, error) {
	return s.Val, s.Err
}
