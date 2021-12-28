package stub

type Stub struct {
	EncodedVal string
	EncodeErr  error
	DecodedVal string
	DecodeErr  error
}

func (s Stub) Encode(_ string) (string, error) {
	return s.EncodedVal, s.EncodeErr
}

func (s Stub) Decode(_ string) (string, error) {
	return s.DecodedVal, s.DecodeErr
}
