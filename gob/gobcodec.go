package gob

import (
	"bytes"
	"encoding/gob"
)

// Codec implements a Codec using gob encoding.
type Codec struct {
}

// Ensure that the Codec is a valid Codec
var _ Codec = Codec{}

// Encode returns a byte slice encoding the value passed as an argument
func (c Codec) Encode(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode converts the byte slice to data that it stores in the provided pointer
func (c Codec) Decode(data []byte, valuePtr interface{}) error {
	buf := bytes.Buffer{}
	dec := gob.NewDecoder(&buf)

	buf.Write(data)
	return dec.Decode(valuePtr)
}
