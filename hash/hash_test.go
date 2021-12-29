package hash

import "testing"

func TestGenerator_Encode(t *testing.T) {
	tests := []struct {
		name      string
		epoch     int64
		rawVal    string
		hashedVal string
	}{
		{
			name:      "success_text",
			epoch:     1000000000000000000,
			rawVal:    "test",
			hashedVal: "1p8Znm3qNwBX",
		},
		{
			name:      "success_empty_hash",
			epoch:     2000000000000000000,
			rawVal:    "test",
			hashedVal: "1eljm8OwNQ0xo",
		},
		{
			name:      "success_url",
			epoch:     3000000000000000000,
			rawVal:    "http://jemgunay.co.uk",
			hashedVal: "5zxqVAkdbvaAp",
		},
		{
			name:      "success_url_long",
			epoch:     4000000000000000000,
			rawVal:    "http://jemgunay.co.uk/this/is/a?test=123456789",
			hashedVal: "GP3Po1BrkVmnq",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := New()
			// override standard epoch generator
			hasher.EpochFunc = func() int64 {
				return tt.epoch
			}

			hashedVal, err := hasher.Hash(tt.rawVal)
			if err != nil {
				t.Fatalf("failed to encode epoch: %s", err)
			}

			if hashedVal != tt.hashedVal {
				t.Fatalf("expected: %s, got: %s", tt.hashedVal, hashedVal)
			}
		})
	}
}
