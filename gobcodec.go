package metafile

import (
	"bytes"
	"encoding/gob"
)

// GobCodec implements a Codec using gob encoding.
type GobCodec struct {
}

// Ensure that the GobCodec is a valid Codec
var _ Codec = GobCodec{}

// Encode returns a byte slice encoding the value passed as an argument
func (c GobCodec) Encode(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode converts the byte slice to data that it stores in the provided pointer
func (c GobCodec) Decode(data []byte, valuePtr interface{}) error {
	buf := bytes.Buffer{}
	dec := gob.NewDecoder(&buf)

	buf.Write(data)
	return dec.Decode(valuePtr)
}
